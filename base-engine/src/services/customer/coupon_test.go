package customer

import (
	"context"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func couponFixture(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal, *gen.CustomerMember) {
	t.Helper()
	service, db, principal := newCustomerServiceFixture(t)
	for _, action := range []string{"hqCustomerCoupon:read", "hqCustomerCoupon:manage", "hqCustomerCoupon:grant", "hqCustomerCoupon:revoke"} {
		principal.Permissions[action] = struct{}{}
	}
	member, err := service.CreateMember(context.Background(), principal, customerCreateInput("member-one"))
	if err != nil {
		t.Fatal(err)
	}
	return service, db, principal, member
}

func validCouponInput() FixedAmountTemplateInput {
	effectiveAt := int64(1_700_000_000_000)
	return FixedAmountTemplateInput{
		Code: "WELCOME", Title: "Welcome", AmountFen: 500,
		RequestKey:  "template-one",
		MinSpendFen: 1000, DaysAfterActivation: 30, EffectiveAt: &effectiveAt,
		PerMemberLimit: 1, TotalIssueLimit: 2, Enabled: true,
	}
}

func TestCustomerCouponTemplate(t *testing.T) {
	service, _, principal, _ := couponFixture(t)
	ctx := context.Background()
	for _, mutate := range []func(*FixedAmountTemplateInput){
		func(v *FixedAmountTemplateInput) { v.AmountFen = 0 },
		func(v *FixedAmountTemplateInput) { v.MinSpendFen = -1 },
		func(v *FixedAmountTemplateInput) { v.DaysAfterActivation = 0 },
		func(v *FixedAmountTemplateInput) { v.PerMemberLimit = 0 },
		func(v *FixedAmountTemplateInput) { v.TotalIssueLimit = -1 },
		func(v *FixedAmountTemplateInput) { v.AmountFen = 1 << 31 },
	} {
		input := validCouponInput()
		mutate(&input)
		if _, err := service.CreateCouponTemplate(ctx, principal, input); auth.ErrorCode(err) != auth.CodeValidationFailed {
			t.Fatalf("invalid coupon template: %v", err)
		}
	}
	first, err := service.CreateCouponTemplate(ctx, principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	assertCreatedTime(t, first.CreatedAt)
	replay, err := service.CreateCouponTemplate(ctx, principal, validCouponInput())
	if err != nil || replay.ID != first.ID {
		t.Fatalf("template replay = %+v, %v", replay, err)
	}
	changedTime := validCouponInput()
	changedEffectiveAt := *changedTime.EffectiveAt + 1
	changedTime.EffectiveAt = &changedEffectiveAt
	if _, err := service.CreateCouponTemplate(ctx, principal, changedTime); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("changed effective time replay = %v", err)
	}
	duplicate := validCouponInput()
	duplicate.RequestKey = "different-request"
	if _, err := service.CreateCouponTemplate(ctx, principal, duplicate); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("duplicate code = %v", err)
	}
	if _, err := service.SetCouponTemplateEnabled(ctx, principal, first.ID, false); err != nil {
		t.Fatal(err)
	}
}

func TestCustomerCouponRevocationUsesUnixMillis(t *testing.T) {
	service, db, principal, member := couponFixture(t)
	ctx := context.Background()
	template, err := service.CreateCouponTemplate(ctx, principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	grant, err := service.GrantCoupon(ctx, principal, template.ID, member.ID, "grant-1")
	if err != nil {
		t.Fatal(err)
	}
	revokedAt := int64(1_800_000_000_321)
	service.now = func() time.Time { return time.UnixMilli(revokedAt) }
	revoked, err := service.RevokeCoupon(ctx, principal, grant.ID, "OPERATOR_REVOCATION")
	if err != nil {
		t.Fatal(err)
	}
	var persisted gen.CustomerCouponGrant
	if err := db.First(&persisted, "id = ?", grant.ID).Error; err != nil {
		t.Fatal(err)
	}
	if revoked.RevokedAt == nil || *revoked.RevokedAt != revokedAt || persisted.RevokedAt == nil || *persisted.RevokedAt != revokedAt {
		t.Fatalf("revoked_at = returned:%v persisted:%v, want %d", revoked.RevokedAt, persisted.RevokedAt, revokedAt)
	}
}

func TestCustomerCouponGrantLimits(t *testing.T) {
	service, db, principal, member := couponFixture(t)
	ctx := context.Background()
	template, err := service.CreateCouponTemplate(ctx, principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	first, err := service.GrantCoupon(ctx, principal, template.ID, member.ID, "grant-1")
	if err != nil {
		t.Fatal(err)
	}
	replay, err := service.GrantCoupon(ctx, principal, template.ID, member.ID, "grant-1")
	if err != nil || replay.ID != first.ID {
		t.Fatalf("grant replay = %+v, %v", replay, err)
	}
	if _, err := service.GrantCoupon(ctx, principal, template.ID, member.ID, "grant-2"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("per member cap = %v", err)
	}
	if _, err := service.RevokeCoupon(ctx, principal, first.ID, "OPERATOR_REVOCATION"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.GrantCoupon(ctx, principal, template.ID, member.ID, "grant-3"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("revocation restored cap = %v", err)
	}
	var persisted gen.CustomerCouponTemplate
	if err := db.First(&persisted, "id = ?", template.ID).Error; err != nil || persisted.IssuedCount != 1 {
		t.Fatalf("issued count = %d, %v", persisted.IssuedCount, err)
	}
}

func TestCustomerCouponGrantStillAllowedWhenOnlyPointsAreFrozen(t *testing.T) {
	service, db, principal, member := couponFixture(t)
	template, err := service.CreateCouponTemplate(context.Background(), principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Model(member).Update("points_frozen", true).Error; err != nil {
		t.Fatal(err)
	}
	grant, err := service.GrantCoupon(context.Background(), principal, template.ID, member.ID, "frozen-points-coupon")
	if err != nil || grant == nil || grant.Status != gen.CustomerCouponGrantStatusPendingActivation {
		t.Fatalf("coupon issuance with frozen points = %+v, %v", grant, err)
	}
}

func TestCustomerCouponGrantReplayCannotCrossHeadquarters(t *testing.T) {
	service, _, principal, member := couponFixture(t)
	ctx := context.Background()
	template, err := service.CreateCouponTemplate(ctx, principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.GrantCoupon(ctx, principal, template.ID, member.ID, "grant-1"); err != nil {
		t.Fatal(err)
	}
	foreign := *principal
	foreignOrganizationID := "another-hq"
	foreign.OrganizationID = &foreignOrganizationID
	if _, err := service.GrantCoupon(ctx, &foreign, template.ID, member.ID, "grant-1"); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("cross-headquarters grant replay = %v", err)
	}
}

func TestCustomerCouponRevocationAuditUsesStableAction(t *testing.T) {
	service, db, principal, member := couponFixture(t)
	ctx := context.Background()
	template, err := service.CreateCouponTemplate(ctx, principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	grant, err := service.GrantCoupon(ctx, principal, template.ID, member.ID, "grant-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.RevokeCoupon(ctx, principal, grant.ID, "OPERATOR_REVOCATION"); err != nil {
		t.Fatal(err)
	}
	var logs []gen.AuditLog
	if err := db.Where("resource_type = ? AND resource_id = ? AND action <> ?", "customerCouponGrant", grant.ID, "hqCustomerCoupon:grant").Find(&logs).Error; err != nil {
		t.Fatal(err)
	}
	if len(logs) != 1 || logs[0].Action != "hqCustomerCoupon:revoke" ||
		logs[0].MetadataJSON == nil || !strings.Contains(*logs[0].MetadataJSON, `"reasonCode":"OPERATOR_REVOCATION"`) {
		t.Fatalf("revocation audit = %+v", logs)
	}
}

func TestCustomerCouponPendingActivation(t *testing.T) {
	service, _, principal, member := couponFixture(t)
	template, err := service.CreateCouponTemplate(context.Background(), principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	grant, err := service.GrantCoupon(context.Background(), principal, template.ID, member.ID, "grant-1")
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return time.Now().AddDate(1, 0, 0) }
	var persisted gen.CustomerCouponGrant
	if err := service.db.First(&persisted, "id = ?", grant.ID).Error; err != nil {
		t.Fatal(err)
	}
	if persisted.Status != gen.CustomerCouponGrantStatusPendingActivation || persisted.ActivatedAt != nil || persisted.ExpiresAt != nil {
		t.Fatalf("premature activation: %+v", grant)
	}
}

func TestCustomerCouponTotalCapAndStatus(t *testing.T) {
	service, db, principal, member := couponFixture(t)
	ctx := context.Background()
	input := validCouponInput()
	input.TotalIssueLimit = 1
	template, err := service.CreateCouponTemplate(ctx, principal, input)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetCouponTemplateEnabled(ctx, principal, template.ID, false); err != nil {
		t.Fatal(err)
	}
	if _, err := service.GrantCoupon(ctx, principal, template.ID, member.ID, "disabled"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("disabled grant = %v", err)
	}
	if _, err := service.SetCouponTemplateEnabled(ctx, principal, template.ID, true); err != nil {
		t.Fatal(err)
	}
	if _, err := service.GrantCoupon(ctx, principal, template.ID, member.ID, "first"); err != nil {
		t.Fatal(err)
	}
	otherInput := customerCreateInput("member-two")
	otherInput.Phone = "13900000002"
	other, err := service.CreateMember(ctx, principal, otherInput)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.GrantCoupon(ctx, principal, template.ID, other.ID, "second"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("total cap = %v", err)
	}
	if err := db.Model(member).Update("status", gen.CustomerMemberStatusCancelPending).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.GrantCoupon(ctx, principal, template.ID, member.ID, "pending"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("pending member grant = %v", err)
	}
}
