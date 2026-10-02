package customer

import (
	"context"
	"testing"
	"time"

	"base-engine/gen"
	"gorm.io/gorm"
)

const couponDistributionTestNow int64 = 1_800_000_000_000

func TestCouponDistributionSchedulesEnabledHQTemplate(t *testing.T) {
	service, db, principal, _ := couponFixture(t)
	clearCouponDistributionJobs(t, db)
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	input := validCouponInput()
	input.EffectiveAt = nil
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	assertCouponDistributionJob(t, db, "TEMPLATE_FANOUT:"+template.ID, gen.CustomerCouponDistributionJobStatusPending, couponDistributionTestNow)
	if _, err := service.CreateCouponTemplate(context.Background(), principal, input); err != nil {
		t.Fatal(err)
	}
	assertCouponDistributionJobCount(t, db, "TEMPLATE_FANOUT:"+template.ID, 1)
}

func TestCouponDistributionSkipsDisabledTemplateUntilEnabled(t *testing.T) {
	service, db, principal, _ := couponFixture(t)
	clearCouponDistributionJobs(t, db)
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	input := validCouponInput()
	input.Enabled = false
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	assertCouponDistributionJobCount(t, db, "TEMPLATE_FANOUT:"+template.ID, 0)
	if _, err := service.SetCouponTemplateEnabled(context.Background(), principal, template.ID, true); err != nil {
		t.Fatal(err)
	}
	assertCouponDistributionJob(t, db, "TEMPLATE_FANOUT:"+template.ID, gen.CustomerCouponDistributionJobStatusPending, couponDistributionTestNow)
}

func TestCouponDistributionSchedulesEnabledStoreTemplate(t *testing.T) {
	service, db, _, principal, _ := storeCouponFixture(t)
	clearCouponDistributionJobs(t, db)
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	input := validCouponInput()
	template, err := service.CreateStoreCouponTemplate(context.Background(), principal, "store", input)
	if err != nil {
		t.Fatal(err)
	}
	assertCouponDistributionJob(t, db, "TEMPLATE_FANOUT:"+template.ID, gen.CustomerCouponDistributionJobStatusPending, couponDistributionTestNow)
}

func TestCouponDistributionSchedulesMemberCreationAndUnsuspension(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	member, err := service.CreateMember(context.Background(), principal, customerCreateInput("distribution-member"))
	if err != nil {
		t.Fatal(err)
	}
	key := "MEMBER_CATCHUP:" + member.ID
	assertCouponDistributionJob(t, db, key, gen.CustomerCouponDistributionJobStatusPending, couponDistributionTestNow)
	if err := db.Model(&gen.CustomerCouponDistributionJob{}).Where("request_key = ?", key).Updates(map[string]any{
		"status": gen.CustomerCouponDistributionJobStatusCompleted, "cursor_created_at": couponDistributionTestNow, "cursor_key": "done", "attempts": 3,
	}).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetMemberStatus(context.Background(), principal, member.ID, gen.CustomerMemberStatusSuspended); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetMemberStatus(context.Background(), principal, member.ID, gen.CustomerMemberStatusActive); err != nil {
		t.Fatal(err)
	}
	var job gen.CustomerCouponDistributionJob
	if err := db.Where("request_key = ?", key).First(&job).Error; err != nil {
		t.Fatal(err)
	}
	if job.Status != gen.CustomerCouponDistributionJobStatusPending || job.CursorCreatedAt != nil || job.CursorKey != nil || job.Attempts != 0 {
		t.Fatalf("reactivated member job was not reset: %+v", job)
	}
}

func TestCouponDistributionDoesNotResetActiveMemberReplay(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	member, err := service.CreateMember(context.Background(), principal, customerCreateInput("active-replay"))
	if err != nil {
		t.Fatal(err)
	}
	key := "MEMBER_CATCHUP:" + member.ID
	if err := db.Model(&gen.CustomerCouponDistributionJob{}).Where("request_key = ?", key).Update("status", gen.CustomerCouponDistributionJobStatusCompleted).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetMemberStatus(context.Background(), principal, member.ID, gen.CustomerMemberStatusActive); err != nil {
		t.Fatal(err)
	}
	assertCouponDistributionJob(t, db, key, gen.CustomerCouponDistributionJobStatusCompleted, couponDistributionTestNow)
}

func TestCouponDistributionScheduleRollsBackWithTemplateAudit(t *testing.T) {
	service, db, principal, _ := couponFixture(t)
	clearCouponDistributionJobs(t, db)
	if err := db.Exec(`CREATE TRIGGER fail_coupon_template_audit BEFORE INSERT ON audit_logs WHEN NEW.resource_type = 'customerCouponTemplate' BEGIN SELECT RAISE(ABORT, 'forced audit failure'); END`).Error; err != nil {
		t.Fatal(err)
	}
	input := validCouponInput()
	input.RequestKey = "rollback-template"
	if _, err := service.CreateCouponTemplate(context.Background(), principal, input); err == nil {
		t.Fatal("expected forced audit failure")
	}
	var templates, jobs int64
	if err := db.Model(&gen.CustomerCouponTemplate{}).Where("request_key = ?", input.RequestKey).Count(&templates).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&gen.CustomerCouponDistributionJob{}).Count(&jobs).Error; err != nil {
		t.Fatal(err)
	}
	if templates != 0 || jobs != 0 {
		t.Fatalf("rollback left template=%d jobs=%d", templates, jobs)
	}
}

func TestCouponDistributionFanoutIssuesToExistingMember(t *testing.T) {
	service, db, principal, member := couponFixture(t)
	clearCouponDistributionJobs(t, db)
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	template, err := service.CreateCouponTemplate(context.Background(), principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.ProcessNextCouponDistributionBatch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !result.Claimed || !result.Completed || result.Issued != 1 || result.Skipped != 0 {
		t.Fatalf("fanout result = %+v", result)
	}
	assertAutomaticCouponGrant(t, db, template.ID, member.ID)
}

func TestCouponDistributionMemberCatchupIncludesHQAndStoreTemplates(t *testing.T) {
	service, db, hq, storePrincipal, _ := storeCouponFixture(t)
	clearCouponDistributionJobs(t, db)
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	hqInput := validCouponInput()
	hqTemplate, err := service.CreateCouponTemplate(context.Background(), hq, hqInput)
	if err != nil {
		t.Fatal(err)
	}
	storeInput := validCouponInput()
	storeInput.Code, storeInput.RequestKey = "STORE-WELCOME", "store-welcome"
	storeTemplate, err := service.CreateStoreCouponTemplate(context.Background(), storePrincipal, "store", storeInput)
	if err != nil {
		t.Fatal(err)
	}
	clearCouponDistributionJobs(t, db)
	memberInput := customerCreateInput("later-member")
	memberInput.Phone = "13800000002"
	member, err := service.CreateMember(context.Background(), hq, memberInput)
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.ProcessNextCouponDistributionBatch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !result.Claimed || !result.Completed || result.Issued != 2 {
		t.Fatalf("catchup result = %+v", result)
	}
	assertAutomaticCouponGrant(t, db, hqTemplate.ID, member.ID)
	assertAutomaticCouponGrant(t, db, storeTemplate.ID, member.ID)
}

func clearCouponDistributionJobs(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&gen.CustomerCouponDistributionJob{}).Error; err != nil {
		t.Fatal(err)
	}
}

func assertCouponDistributionJob(t *testing.T, db *gorm.DB, key string, status gen.CustomerCouponDistributionJobStatus, availableAt int64) {
	t.Helper()
	var job gen.CustomerCouponDistributionJob
	if err := db.Where("request_key = ?", key).First(&job).Error; err != nil {
		t.Fatal(err)
	}
	if job.Status != status || job.AvailableAt != availableAt {
		t.Fatalf("job %s = %+v", key, job)
	}
}

func assertCouponDistributionJobCount(t *testing.T, db *gorm.DB, key string, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(&gen.CustomerCouponDistributionJob{}).Where("request_key = ?", key).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != want {
		t.Fatalf("job %s count = %d, want %d", key, count, want)
	}
}

func assertAutomaticCouponGrant(t *testing.T, db *gorm.DB, templateID, memberID string) {
	t.Helper()
	var grant gen.CustomerCouponGrant
	if err := db.Where("template_id = ? AND member_id = ?", templateID, memberID).First(&grant).Error; err != nil {
		t.Fatal(err)
	}
	if grant.RequestKey != autoCouponRequestKey(templateID, memberID) || grant.IssuedAt != couponDistributionTestNow {
		t.Fatalf("automatic grant = %+v", grant)
	}
}
