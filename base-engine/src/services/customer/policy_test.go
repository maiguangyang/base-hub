package customer

import (
	"context"
	"math"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func validCustomerPolicy() PolicyInput {
	return PolicyInput{
		DiscountEnabled: true, DiscountBasisPoints: 9500,
		PurchaseEarnEnabled: true, EarnAmountFen: 100, EarnPoints: 1,
		RedemptionEnabled: true, RedeemPoints: 100, RedeemAmountFen: 100,
		MaxRedemptionBasisPoints: 5000, MaxRedemptionPoints: 5000,
		ManualGrantMaxSingle: 1000, ManualGrantMaxDaily: 5000,
	}
}

func TestCustomerPolicyValidation(t *testing.T) {
	service, _, principal := newCustomerServiceFixture(t)
	ctx := context.Background()
	cases := []struct {
		name string
		edit func(*PolicyInput)
	}{
		{"discount range", func(p *PolicyInput) { p.DiscountBasisPoints = 10001 }},
		{"earn denominator", func(p *PolicyInput) { p.EarnAmountFen = 0 }},
		{"redeem numerator", func(p *PolicyInput) { p.RedeemPoints = 0 }},
		{"redemption cap", func(p *PolicyInput) { p.MaxRedemptionBasisPoints = 10001 }},
		{"single cap", func(p *PolicyInput) { p.ManualGrantMaxSingle = 0 }},
		{"daily cap", func(p *PolicyInput) { p.ManualGrantMaxDaily = 1 }},
		{"graphql int", func(p *PolicyInput) { p.EarnPoints = math.MaxInt32 + 1 }},
		{"int64 overflow", func(p *PolicyInput) { p.EarnAmountFen = math.MaxInt64 }},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			input := validCustomerPolicy()
			testCase.edit(&input)
			if _, err := service.SavePolicy(ctx, principal, 0, input); auth.ErrorCode(err) != auth.CodeValidationFailed {
				t.Fatalf("SavePolicy error = %v", err)
			}
		})
	}
}

func TestCustomerPolicyAllowsZeroManualGrantCapsToDisableGrants(t *testing.T) {
	service, _, principal := newCustomerServiceFixture(t)
	input := validCustomerPolicy()
	input.ManualGrantMaxSingle = 0
	input.ManualGrantMaxDaily = 0

	saved, err := service.SavePolicy(context.Background(), principal, 0, input)
	if err != nil {
		t.Fatal(err)
	}
	if saved.ManualGrantMaxSingle != 0 || saved.ManualGrantMaxDaily != 0 {
		t.Fatalf("manual grant caps = %d/%d", saved.ManualGrantMaxSingle, saved.ManualGrantMaxDaily)
	}
}

func TestCustomerPolicyVersionAndScope(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	ctx := context.Background()
	assertCustomerPolicyDraft(t, service, principal)
	first, err := service.SavePolicy(ctx, principal, 0, validCustomerPolicy())
	if err != nil {
		t.Fatal(err)
	}
	if first.Version != 1 {
		t.Fatalf("first save = %+v, %v", first, err)
	}
	if _, err := service.SavePolicy(ctx, principal, 0, validCustomerPolicy()); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("stale save = %v", err)
	}
	input := validCustomerPolicy()
	input.DiscountBasisPoints = 9000
	second, err := service.SavePolicy(ctx, principal, 1, input)
	if err != nil {
		t.Fatal(err)
	}
	if second.Version != 2 || second.DiscountBasisPoints != 9000 {
		t.Fatalf("second save = %+v, %v", second, err)
	}
	var count int64
	if err := db.Table("customer_benefit_policies").Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("policy rows = %d, %v", count, err)
	}
}

func assertCustomerPolicyDraft(t *testing.T, service *Service, principal *auth.WorkspacePrincipal) {
	t.Helper()
	draft, err := service.GetPolicy(context.Background(), principal)
	if err != nil {
		t.Fatal(err)
	}
	if draft.Version != 0 || draft.RedemptionEnabled || draft.RedeemPoints != 100 || draft.RedeemAmountFen != 100 {
		t.Fatalf("draft = %+v", draft)
	}
}

func TestCustomerPolicyScope(t *testing.T) {
	service, _, principal := newCustomerServiceFixture(t)
	principal.WorkspaceType = auth.WorkspaceTypeFranchise
	if _, err := service.GetPolicy(context.Background(), principal); auth.ErrorCode(err) != auth.CodeWorkspaceForbidden {
		t.Fatalf("franchise read = %v", err)
	}
}

func TestPolicyStackingDefaultsAndVersionedSave(t *testing.T) {
	service, _, principal := newCustomerServiceFixture(t)
	draft, err := service.GetPolicy(context.Background(), principal)
	if err != nil {
		t.Fatal(err)
	}
	if draft.PromotionWithHqCoupon || draft.PromotionWithStoreCoupon || draft.PointsWithPromotion {
		t.Fatalf("unsafe default = %+v", draft)
	}
	input := validCustomerPolicy()
	input.PromotionWithHqCoupon = true
	input.PointsWithCoupon = true
	saved, err := service.SavePolicy(context.Background(), principal, 0, input)
	if err != nil {
		t.Fatal(err)
	}
	if !saved.PromotionWithHqCoupon || !saved.PointsWithCoupon || saved.PromotionWithStoreCoupon {
		t.Fatalf("saved = %+v", saved)
	}
}

func TestPolicyTighteningDisablesIncompatiblePromotions(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	if err := db.AutoMigrate(&gen.StorePromotion{}); err != nil {
		t.Fatal(err)
	}
	input := validCustomerPolicy()
	input.PromotionWithHqCoupon = true
	if _, err := service.SavePolicy(context.Background(), principal, 0, input); err != nil {
		t.Fatal(err)
	}
	promotion := gen.StorePromotion{ID: "promotion-1", StoreID: "store", RuleKey: "rule-1", Version: 1,
		Kind: gen.StorePromotionKindItemPrice, TimeZone: "Asia/Shanghai", Enabled: true, StackWithHqCoupon: true}
	if err := db.Create(&promotion).Error; err != nil {
		t.Fatal(err)
	}
	input.PromotionWithHqCoupon = false
	if _, err := service.SavePolicy(context.Background(), principal, 1, input); err != nil {
		t.Fatal(err)
	}
	if err := db.First(&promotion, "id = ?", promotion.ID).Error; err != nil {
		t.Fatal(err)
	}
	if promotion.Enabled {
		t.Fatal("promotion kept a stacking option above the HQ ceiling")
	}
}
