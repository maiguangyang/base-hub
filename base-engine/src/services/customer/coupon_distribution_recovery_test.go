package customer

import (
	"context"
	"errors"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestCouponDistributionFailureReturnsJobToPendingBackoff(t *testing.T) {
	service, db, principal := couponDistributionServiceWithMembers(t, 1)
	template, err := service.CreateCouponTemplate(context.Background(), principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Delete(&gen.CustomerCouponTemplate{}, "id = ?", template.ID).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.ProcessNextCouponDistributionBatch(context.Background()); err == nil {
		t.Fatal("expected missing-template processing failure")
	}
	job := loadCouponDistributionJob(t, db, templateCouponDistributionKey(template.ID))
	assertFailedCouponDistributionJob(t, job)
	result, err := service.ProcessNextCouponDistributionBatch(context.Background())
	if err != nil || result.Claimed {
		t.Fatalf("backed-off job was claimable: result=%+v err=%v", result, err)
	}
}

func TestCouponDistributionExpiredLeaseIsRecovered(t *testing.T) {
	service, db, principal := couponDistributionServiceWithMembers(t, 1)
	template, err := service.CreateCouponTemplate(context.Background(), principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	key := templateCouponDistributionKey(template.ID)
	if err := db.Model(&gen.CustomerCouponDistributionJob{}).Where("request_key = ?", key).Updates(map[string]any{
		"status": gen.CustomerCouponDistributionJobStatusRunning, "lease_token": "expired", "lease_expires_at": couponDistributionTestNow - 1,
	}).Error; err != nil {
		t.Fatal(err)
	}
	result, err := service.ProcessNextCouponDistributionBatch(context.Background())
	if err != nil || !result.Completed || result.Issued != 1 {
		t.Fatalf("expired lease recovery result=%+v err=%v", result, err)
	}
	assertCouponDistributionJob(t, db, key, gen.CustomerCouponDistributionJobStatusCompleted, couponDistributionTestNow)
}

func TestCouponDistributionLostLeaseRollsBackItemWithoutRetry(t *testing.T) {
	service, db, principal := couponDistributionServiceWithMembers(t, 1)
	template, err := service.CreateCouponTemplate(context.Background(), principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	job, candidate := claimCouponJobWithSingleCandidate(t, service, couponDistributionTestNow)
	if err := db.Model(&gen.CustomerCouponDistributionJob{}).Where("id = ?", job.ID).Update("lease_token", "replacement").Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.processCouponDistributionCandidate(context.Background(), job, candidate, couponDistributionTestNow); !errors.Is(err, errCouponDistributionLeaseRelinquished) {
		t.Fatalf("expected relinquished lease, got %v", err)
	}
	assertGrantCountForTemplate(t, db, template.ID, 0)
	stored := loadCouponDistributionJob(t, db, job.RequestKey)
	assertLostCouponDistributionLeaseState(t, stored)
}

func TestCouponDistributionCursorCommitRenewsLease(t *testing.T) {
	service, db, principal := couponDistributionServiceWithMembers(t, 1)
	current := couponDistributionTestNow
	service.now = func() time.Time { return time.UnixMilli(current) }
	_, err := service.CreateCouponTemplate(context.Background(), principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	job, candidate := claimCouponJobWithSingleCandidate(t, service, current)
	current += 5_000
	if _, err := service.processCouponDistributionCandidate(context.Background(), job, candidate, couponDistributionTestNow); err != nil {
		t.Fatal(err)
	}
	stored := loadCouponDistributionJob(t, db, job.RequestKey)
	want := current + couponDistributionLeaseDuration.Milliseconds()
	assertRenewedCouponDistributionLease(t, stored, candidate.ID, want)
}

func TestCouponDistributionClaimUpdatesOnlySelectedExpiredJob(t *testing.T) {
	service, db, principal := couponDistributionServiceWithMembers(t, 1)
	first := validCouponInput()
	first.RequestKey, first.Code = "claim-first", "CLAIM-FIRST"
	firstTemplate, err := service.CreateCouponTemplate(context.Background(), principal, first)
	if err != nil {
		t.Fatal(err)
	}
	second := validCouponInput()
	second.RequestKey, second.Code = "claim-second", "CLAIM-SECOND"
	secondTemplate, err := service.CreateCouponTemplate(context.Background(), principal, second)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&gen.CustomerCouponDistributionJob{}).Where("request_key IN ?", []string{
		templateCouponDistributionKey(firstTemplate.ID), templateCouponDistributionKey(secondTemplate.ID),
	}).Updates(map[string]any{"status": gen.CustomerCouponDistributionJobStatusRunning,
		"lease_token": "expired", "lease_expires_at": couponDistributionTestNow - 1}).Error; err != nil {
		t.Fatal(err)
	}
	claimed, err := service.claimNextCouponDistributionJob(context.Background(), couponDistributionTestNow)
	if err != nil || claimed == nil {
		t.Fatalf("claim=%+v err=%v", claimed, err)
	}
	var stillExpired int64
	if err := db.Model(&gen.CustomerCouponDistributionJob{}).Where("lease_token = ?", "expired").Count(&stillExpired).Error; err != nil || stillExpired != 1 {
		t.Fatalf("untargeted expired jobs=%d err=%v", stillExpired, err)
	}
}

func couponDistributionServiceWithMembers(t *testing.T, count int) (*Service, *gorm.DB, *auth.WorkspacePrincipal) {
	t.Helper()
	service, db, principal := newCustomerServiceFixture(t)
	principal.Permissions["hqCustomerCoupon:manage"] = struct{}{}
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	seedCouponDistributionMembers(t, db, count)
	return service, db, principal
}

func loadCouponDistributionJob(t *testing.T, db *gorm.DB, key string) gen.CustomerCouponDistributionJob {
	t.Helper()
	var job gen.CustomerCouponDistributionJob
	if err := db.Where("request_key = ?", key).First(&job).Error; err != nil {
		t.Fatal(err)
	}
	return job
}

func claimCouponJobWithSingleCandidate(t *testing.T, service *Service, nowMillis int64) (*gen.CustomerCouponDistributionJob, couponDistributionCandidate) {
	t.Helper()
	job, err := service.claimNextCouponDistributionJob(context.Background(), nowMillis)
	if err != nil || job == nil {
		t.Fatalf("claim job=%+v err=%v", job, err)
	}
	candidates, _, err := service.couponDistributionCandidates(context.Background(), job, nowMillis)
	if err != nil || len(candidates) != 1 {
		t.Fatalf("candidates=%+v err=%v", candidates, err)
	}
	return job, candidates[0]
}

func assertFailedCouponDistributionJob(t *testing.T, job gen.CustomerCouponDistributionJob) {
	t.Helper()
	validState := job.Status == gen.CustomerCouponDistributionJobStatusPending && job.Attempts == 1 &&
		job.AvailableAt == couponDistributionTestNow+1_000 && job.LeaseToken == nil && job.LeaseExpiresAt == nil
	validCode := job.LastErrorCode != nil && *job.LastErrorCode == "COUPON_DISTRIBUTION_FAILED"
	if !validState || !validCode {
		t.Fatalf("failed job state = %+v", job)
	}
}

func assertLostCouponDistributionLeaseState(t *testing.T, job gen.CustomerCouponDistributionJob) {
	t.Helper()
	if job.Attempts != 0 || job.LeaseToken == nil || *job.LeaseToken != "replacement" {
		t.Fatalf("lost lease mutated job = %+v", job)
	}
}

func assertRenewedCouponDistributionLease(t *testing.T, job gen.CustomerCouponDistributionJob, cursorKey string, want int64) {
	t.Helper()
	validLease := job.LeaseExpiresAt != nil && *job.LeaseExpiresAt == want
	validCursor := job.CursorKey != nil && *job.CursorKey == cursorKey
	if !validLease || !validCursor {
		t.Fatalf("renewed job = %+v, want lease %d", job, want)
	}
}
