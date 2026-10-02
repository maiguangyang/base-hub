package customer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestCouponDistributionFanoutUsesDeterministicFiniteCap(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	principal.Permissions["hqCustomerCoupon:manage"] = struct{}{}
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	seedCouponDistributionMembers(t, db, 3)
	input := validCouponInput()
	input.TotalIssueLimit = 2
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.ProcessNextCouponDistributionBatch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !result.Completed || result.Issued != 2 || result.Skipped != 1 {
		t.Fatalf("finite-cap result = %+v", result)
	}
	assertGrantCountForTemplate(t, db, template.ID, 2)
	for _, memberID := range []string{"distribution-member-000", "distribution-member-001"} {
		assertAutomaticCouponGrant(t, db, template.ID, memberID)
	}
	var lastCount int64
	if err := db.Model(&gen.CustomerCouponGrant{}).Where("template_id = ? AND member_id = ?", template.ID, "distribution-member-002").Count(&lastCount).Error; err != nil || lastCount != 0 {
		t.Fatalf("member after cap grant count = %d, %v", lastCount, err)
	}
}

func TestCouponDistributionOverlappingCatchupPreservesFiniteCapMemberOrder(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	principal.Permissions["hqCustomerCoupon:manage"] = struct{}{}
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	seedCouponDistributionMembers(t, db, 2)
	lateInput := customerCreateInput("distribution-late-member")
	lateInput.Phone = "13800000003"
	lateMember, err := service.CreateMember(context.Background(), principal, lateInput)
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow + 1) }
	input := validCouponInput()
	input.TotalIssueLimit = 1
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	first := mustProcessCouponDistributionBatch(t, service)
	second := mustProcessCouponDistributionBatch(t, service)
	assertCouponDistributionBatch(t, first, gen.CustomerCouponDistributionJobKindMemberCatchup, 0, 1)
	assertCouponDistributionBatch(t, second, gen.CustomerCouponDistributionJobKindTemplateFanout, 1, 2)
	assertGrantCountForTemplateMember(t, db, template.ID, "distribution-member-000", 1)
	assertGrantCountForTemplateMember(t, db, template.ID, lateMember.ID, 0)
}

func mustProcessCouponDistributionBatch(t *testing.T, service *Service) CouponDistributionBatchResult {
	t.Helper()
	result, err := service.ProcessNextCouponDistributionBatch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func assertCouponDistributionBatch(t *testing.T, result CouponDistributionBatchResult, kind gen.CustomerCouponDistributionJobKind, issued, skipped int) {
	t.Helper()
	if result.Kind != string(kind) || result.Issued != issued || result.Skipped != skipped {
		t.Fatalf("distribution batch = %+v; want kind=%s issued=%d skipped=%d", result, kind, issued, skipped)
	}
}

func assertGrantCountForTemplateMember(t *testing.T, db *gorm.DB, templateID, memberID string, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(&gen.CustomerCouponGrant{}).Where("template_id = ? AND member_id = ?", templateID, memberID).Count(&count).Error; err != nil || count != want {
		t.Fatalf("template/member grant count = %d, %v; want %d", count, err, want)
	}
}

func TestCouponDistributionUnlimitedFanoutContinuesAcrossBatches(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	principal.Permissions["hqCustomerCoupon:manage"] = struct{}{}
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	seedCouponDistributionMembers(t, db, 101)
	input := validCouponInput()
	input.TotalIssueLimit = 0
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	first, err := service.ProcessNextCouponDistributionBatch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if first.Completed || first.Issued != 100 {
		t.Fatalf("first unlimited batch = %+v", first)
	}
	second, err := service.ProcessNextCouponDistributionBatch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !second.Completed || second.Issued != 1 {
		t.Fatalf("second unlimited batch = %+v", second)
	}
	assertGrantCountForTemplate(t, db, template.ID, 101)
}

func TestCouponDistributionCompletesWhenFiniteCapExhaustsFullBatch(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	principal.Permissions["hqCustomerCoupon:manage"] = struct{}{}
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	seedCouponDistributionMembers(t, db, 100)
	input := validCouponInput()
	input.TotalIssueLimit = 1
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.ProcessNextCouponDistributionBatch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !result.Completed || result.Issued != 1 || result.Skipped != 99 {
		t.Fatalf("full finite-cap batch = %+v", result)
	}
	assertCouponDistributionJob(t, db, templateCouponDistributionKey(template.ID), gen.CustomerCouponDistributionJobStatusCompleted, couponDistributionTestNow)
}

func TestCouponDistributionStopsAtCutoffDuringBatch(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	principal.Permissions["hqCustomerCoupon:manage"] = struct{}{}
	seedCouponDistributionMembers(t, db, 2)
	current := couponDistributionTestNow
	calls := 0
	service.now = func() time.Time {
		calls++
		if calls >= 4 {
			current = couponDistributionTestNow + 1
		}
		return time.UnixMilli(current)
	}
	cutoff := couponDistributionTestNow + 1
	input := validCouponInput()
	input.DistributionEndsAt = &cutoff
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	calls = 0
	result, err := service.ProcessNextCouponDistributionBatch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !result.Completed || result.Issued != 1 || result.Skipped != 1 {
		t.Fatalf("cutoff batch result = %+v", result)
	}
	assertGrantCountForTemplate(t, db, template.ID, 1)
}

func TestCouponDistributionCatchupExcludesIneligibleAndExistingTemplates(t *testing.T) {
	service, db, hq, storePrincipal, existingMember := storeCouponFixture(t)
	clearCouponDistributionJobs(t, db)
	service.now = func() time.Time { return time.UnixMilli(couponDistributionTestNow) }
	hqTemplate, err := service.CreateCouponTemplate(context.Background(), hq, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	storeInput := validCouponInput()
	storeInput.Code, storeInput.RequestKey = "STORE", "store"
	storeTemplate, err := service.CreateStoreCouponTemplate(context.Background(), storePrincipal, "store", storeInput)
	if err != nil {
		t.Fatal(err)
	}
	createIneligibleCouponTemplates(t, service, hq)
	clearCouponDistributionJobs(t, db)
	if _, err := service.GrantCoupon(context.Background(), hq, hqTemplate.ID, existingMember.ID, "existing-manual"); err != nil {
		t.Fatal(err)
	}
	if err := scheduleMemberCouponCatchup(db, existingMember, couponDistributionTestNow); err != nil {
		t.Fatal(err)
	}
	result, err := service.ProcessNextCouponDistributionBatch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !result.Completed || result.Issued != 1 || result.Skipped != 1 {
		t.Fatalf("eligibility catchup result = %+v", result)
	}
	assertGrantCountForMember(t, db, existingMember.ID, 2)
	assertAutomaticCouponGrant(t, db, storeTemplate.ID, existingMember.ID)
}

func createIneligibleCouponTemplates(t *testing.T, service *Service, principal *auth.WorkspacePrincipal) {
	t.Helper()
	future := int64(1_900_000_000_000)
	endedEffective, endedAt := int64(1_700_000_000_000), int64(1_750_000_000_000)
	inputs := []FixedAmountTemplateInput{validCouponInput(), validCouponInput(), validCouponInput()}
	inputs[0].Code, inputs[0].RequestKey, inputs[0].EffectiveAt = "FUTURE", "future", &future
	inputs[1].Code, inputs[1].RequestKey, inputs[1].EffectiveAt, inputs[1].DistributionEndsAt = "ENDED", "ended", &endedEffective, &endedAt
	inputs[2].Code, inputs[2].RequestKey, inputs[2].Enabled = "DISABLED", "disabled", false
	for _, input := range inputs {
		if _, err := service.CreateCouponTemplate(context.Background(), principal, input); err != nil {
			t.Fatal(err)
		}
	}
}

func seedCouponDistributionMembers(t *testing.T, db *gorm.DB, count int) {
	t.Helper()
	members := make([]gen.CustomerMember, 0, count)
	for index := 0; index < count; index++ {
		id := fmt.Sprintf("distribution-member-%03d", index)
		members = append(members, gen.CustomerMember{
			ID: id, MemberNumber: fmt.Sprintf("DM%03d", index), RequestKey: id,
			Status: gen.CustomerMemberStatusActive, OrganizationID: "hq",
			NoticeVersion: "v1", ProcessingBasisCode: "service", EvidenceReference: "fixture",
			CreatedAt: 1_700_000_000_000 + int64(index),
		})
	}
	if err := db.CreateInBatches(members, 100).Error; err != nil {
		t.Fatal(err)
	}
}

func assertGrantCountForTemplate(t *testing.T, db *gorm.DB, templateID string, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(&gen.CustomerCouponGrant{}).Where("template_id = ?", templateID).Count(&count).Error; err != nil || count != want {
		t.Fatalf("template %s grant count = %d, %v; want %d", templateID, count, err, want)
	}
}

func assertGrantCountForMember(t *testing.T, db *gorm.DB, memberID string, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(&gen.CustomerCouponGrant{}).Where("member_id = ?", memberID).Count(&count).Error; err != nil || count != want {
		t.Fatalf("member %s grant count = %d, %v; want %d", memberID, count, err, want)
	}
}
