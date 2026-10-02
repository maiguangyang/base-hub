package customer

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestCouponTemplateTimeAndValidityBoundaries(t *testing.T) {
	base := validCouponInput()
	for _, test := range []struct {
		name              string
		effectiveAt, days int64
		valid             bool
	}{
		{"12-digit timestamp", 999_999_999_999, 1, false},
		{"minimum timestamp and one day", minCouponUnixMillis, 1, true},
		{"maximum timestamp and 365 days", maxCouponUnixMillis, 365, true},
		{"zero days", minCouponUnixMillis, 0, false},
		{"366 days", minCouponUnixMillis, 366, false},
		{"14-digit timestamp", 10_000_000_000_000, 1, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			input := base
			input.EffectiveAt, input.DaysAfterActivation = &test.effectiveAt, test.days
			if got := validCouponTemplate(input); got != test.valid {
				t.Fatalf("validCouponTemplate() = %t, want %t", got, test.valid)
			}
		})
	}
}

func TestCustomerCouponTemplateBlankEffectiveIsImmediateAndIdempotent(t *testing.T) {
	service, _, principal, _ := couponFixture(t)
	nowMillis := int64(1_800_000_000_123)
	service.now = func() time.Time { return time.UnixMilli(nowMillis) }
	input := validCouponInput()
	input.EffectiveAt, input.TotalIssueLimit = nil, 0
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	if template.CreatedAt != nowMillis || template.EffectiveAt != nowMillis || template.TotalIssueLimit != 0 {
		t.Fatalf("normalized template = %+v", template)
	}
	service.now = func() time.Time { return time.UnixMilli(nowMillis + 10_000) }
	replay, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil || replay.ID != template.ID {
		t.Fatalf("blank effective replay = %+v, %v", replay, err)
	}
}

func TestCustomerCouponEffectiveTimeGate(t *testing.T) {
	service, db, principal, member := couponFixture(t)
	effectiveAt := int64(1_800_000_000_000)
	distributionEndsAt := effectiveAt + 1_000
	input := validCouponInput()
	input.EffectiveAt, input.DistributionEndsAt, input.PerMemberLimit = &effectiveAt, &distributionEndsAt, 2
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return time.UnixMilli(effectiveAt - 1) }
	assertCouponGrantConflict(t, service, principal, template.ID, member.ID, "before")
	assertCouponCounts(t, db, template.ID, 0, 0)
	service.now = func() time.Time { return time.UnixMilli(effectiveAt) }
	grant, err := service.GrantCoupon(context.Background(), principal, template.ID, member.ID, "exact")
	if err != nil || grant.IssuedAt != effectiveAt {
		t.Fatalf("grant at effective time = %+v, %v", grant, err)
	}
	service.now = func() time.Time { return time.UnixMilli(distributionEndsAt) }
	assertCouponGrantConflict(t, service, principal, template.ID, member.ID, "at-cutoff")
	replay, err := service.GrantCoupon(context.Background(), principal, template.ID, member.ID, "exact")
	if err != nil || replay.ID != grant.ID {
		t.Fatalf("expired template replay = %+v, %v", replay, err)
	}
	assertCouponCounts(t, db, template.ID, 1, 1)
}

func assertCouponGrantConflict(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, templateID, memberID, key string) {
	t.Helper()
	if _, err := service.GrantCoupon(context.Background(), principal, templateID, memberID, key); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("grant %s = %v", key, err)
	}
}

func assertCouponCounts(t *testing.T, db *gorm.DB, templateID string, grants, issued int64) {
	t.Helper()
	var grantCount int64
	if err := db.Model(&gen.CustomerCouponGrant{}).Where("template_id = ?", templateID).Count(&grantCount).Error; err != nil {
		t.Fatal(err)
	}
	var template gen.CustomerCouponTemplate
	if err := db.First(&template, "id = ?", templateID).Error; err != nil {
		t.Fatal(err)
	}
	if grantCount != grants || template.IssuedCount != issued {
		t.Fatalf("coupon counts = grants:%d issued:%d", grantCount, template.IssuedCount)
	}
}

func TestCustomerCouponUnlimitedTotalStillHonorsPerMemberLimit(t *testing.T) {
	service, _, principal, member := couponFixture(t)
	input := validCouponInput()
	input.TotalIssueLimit, input.PerMemberLimit = 0, 2
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"unlimited-1", "unlimited-2"} {
		if _, err := service.GrantCoupon(context.Background(), principal, template.ID, member.ID, key); err != nil {
			t.Fatalf("grant %s = %v", key, err)
		}
	}
	assertCouponGrantConflict(t, service, principal, template.ID, member.ID, "unlimited-3")
}

func TestCustomerCouponAutomaticGrantSkipsExistingAndUsesSystemAudit(t *testing.T) {
	service, db, principal, member := couponFixture(t)
	input := validCouponInput()
	input.Code, input.RequestKey, input.PerMemberLimit = "AUTO-HQ", "auto-hq-template", 2
	template, err := service.CreateCouponTemplate(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	manual, err := service.GrantCoupon(context.Background(), principal, template.ID, member.ID, "manual-before-auto")
	if err != nil {
		t.Fatal(err)
	}
	skipped, created, err := automaticGrant(t, service, db, template, member.ID)
	if err != nil || created || skipped == nil || skipped.ID != manual.ID {
		t.Fatalf("automatic skip = grant:%+v created:%t err:%v", skipped, created, err)
	}
	otherInput := customerCreateInput("automatic-member")
	otherInput.Phone = "13800000002"
	other, err := service.CreateMember(context.Background(), principal, otherInput)
	if err != nil {
		t.Fatal(err)
	}
	automatic, created, err := automaticGrant(t, service, db, template, other.ID)
	if err != nil || !created {
		t.Fatalf("automatic grant = %+v created:%t err:%v", automatic, created, err)
	}
	assertAutomaticCouponAudit(t, db, automatic.ID, "hqCustomerCoupon:grant", template.OrganizationID, nil)
}

func TestCustomerCouponAutomaticDuplicateRequiresExactReplay(t *testing.T) {
	service, db, principal, member := couponFixture(t)
	template, err := service.CreateCouponTemplate(context.Background(), principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	duplicateErr := errors.New("duplicate key")
	if _, _, err := service.resolveAutomaticCouponGrantDuplicate(context.Background(), duplicateErr, template.ID, member.ID); !errors.Is(err, duplicateErr) {
		t.Fatalf("missing exact replay error = %v", err)
	}
	key := autoCouponRequestKey(template.ID, member.ID)
	grant := gen.CustomerCouponGrant{ID: "auto-duplicate", TemplateID: template.ID, MemberID: member.ID,
		RequestKey: key, Status: gen.CustomerCouponGrantStatusPendingActivation, AmountFen: template.AmountFen,
		MinSpendFen: template.MinSpendFen, DaysAfterActivation: template.DaysAfterActivation,
		IssuedAt: time.Now().UnixMilli(), CreatedAt: time.Now().UnixMilli()}
	if err := db.Create(&grant).Error; err != nil {
		t.Fatal(err)
	}
	replay, created, err := service.resolveAutomaticCouponGrantDuplicate(context.Background(), duplicateErr, template.ID, member.ID)
	if err != nil || created || replay == nil || replay.ID != grant.ID {
		t.Fatalf("exact duplicate replay = %+v created:%t err:%v", replay, created, err)
	}
}

func automaticGrant(t *testing.T, service *Service, db *gorm.DB, template *gen.CustomerCouponTemplate, memberID string) (*gen.CustomerCouponGrant, bool, error) {
	t.Helper()
	var grant *gen.CustomerCouponGrant
	var created bool
	err := db.Transaction(func(tx *gorm.DB) error {
		var err error
		grant, created, err = service.createCouponGrantInTransaction(tx, couponGrantInput{
			TemplateID: template.ID, MemberID: memberID, RequestKey: autoCouponRequestKey(template.ID, memberID),
			ExpectedOrganizationID: template.OrganizationID, Automatic: true, NowMillis: time.Now().UnixMilli(),
		})
		return err
	})
	return grant, created, err
}

func assertAutomaticCouponAudit(t *testing.T, db *gorm.DB, grantID, action, organizationID string, storeID *string) {
	t.Helper()
	var log gen.AuditLog
	if err := db.Where("resource_type = ? AND resource_id = ?", "customerCouponGrant", grantID).First(&log).Error; err != nil {
		t.Fatal(err)
	}
	if log.Action != action || log.ActorAccountID != nil || log.SessionID != nil {
		t.Fatalf("automatic audit actor/action = %+v", log)
	}
	if log.OrganizationID == nil || *log.OrganizationID != organizationID || !sameOptionalString(log.StoreID, storeID) {
		t.Fatalf("automatic audit scope = %+v", log)
	}
	if log.MetadataJSON == nil || !strings.Contains(*log.MetadataJSON, `"source":"AUTO_COUPON_DISTRIBUTION"`) {
		t.Fatalf("automatic audit metadata = %+v", log)
	}
}

func sameOptionalString(left, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}
