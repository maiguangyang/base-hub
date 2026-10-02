package src

import (
	"context"
	"math"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/customer"
)

func policyGraphQL(view customer.PolicyView) (*gen.HqCustomerBenefitPolicyView, error) {
	values := []int64{
		view.Version, view.DiscountBasisPoints, view.EarnAmountFen, view.EarnPoints,
		view.RedeemPoints, view.RedeemAmountFen, view.MaxRedemptionBasisPoints,
		view.MaxRedemptionPoints, view.ManualGrantMaxSingle, view.ManualGrantMaxDaily,
	}
	for _, value := range values {
		if value < 0 || value > math.MaxInt32 {
			return nil, auth.NewError(auth.CodeInternalError)
		}
	}
	return &gen.HqCustomerBenefitPolicyView{
		Version: int(view.Version), DiscountEnabled: view.DiscountEnabled,
		DiscountBasisPoints: int(view.DiscountBasisPoints),
		PurchaseEarnEnabled: view.PurchaseEarnEnabled,
		EarnAmountFen:       int(view.EarnAmountFen), EarnPoints: int(view.EarnPoints),
		RedemptionEnabled: view.RedemptionEnabled, RedeemPoints: int(view.RedeemPoints),
		RedeemAmountFen:          int(view.RedeemAmountFen),
		MaxRedemptionBasisPoints: int(view.MaxRedemptionBasisPoints),
		MaxRedemptionPoints:      int(view.MaxRedemptionPoints),
		ManualGrantMaxSingle:     int(view.ManualGrantMaxSingle),
		ManualGrantMaxDaily:      int(view.ManualGrantMaxDaily),
		PromotionWithHqCoupon:    view.PromotionWithHqCoupon,
		PromotionWithStoreCoupon: view.PromotionWithStoreCoupon,
		MemberPriceWithPromotion: view.MemberPriceWithPromotion,
		PointsWithPromotion:      view.PointsWithPromotion,
		PointsWithCoupon:         view.PointsWithCoupon,
	}, nil
}

func (r *QueryResolver) HqCustomerBenefitPolicy(ctx context.Context) (*gen.HqCustomerBenefitPolicyView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	view, err := r.Services.Customers.GetPolicy(ctx, principal)
	if err != nil {
		return nil, err
	}
	return policyGraphQL(view)
}

func (r *MutationResolver) HqSaveCustomerBenefitPolicy(ctx context.Context, expectedVersion int, input gen.HqCustomerBenefitPolicyInput) (*gen.HqCustomerBenefitPolicyView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	view, err := r.Services.Customers.SavePolicyWithOptionalStacking(ctx, principal, int64(expectedVersion), customer.PolicyInput{
		DiscountEnabled:     input.DiscountEnabled,
		DiscountBasisPoints: int64(input.DiscountBasisPoints),
		PurchaseEarnEnabled: input.PurchaseEarnEnabled,
		EarnAmountFen:       int64(input.EarnAmountFen), EarnPoints: int64(input.EarnPoints),
		RedemptionEnabled: input.RedemptionEnabled, RedeemPoints: int64(input.RedeemPoints),
		RedeemAmountFen:          int64(input.RedeemAmountFen),
		MaxRedemptionBasisPoints: int64(input.MaxRedemptionBasisPoints),
		MaxRedemptionPoints:      int64(input.MaxRedemptionPoints),
		ManualGrantMaxSingle:     int64(input.ManualGrantMaxSingle),
		ManualGrantMaxDaily:      int64(input.ManualGrantMaxDaily),
	}, customer.PolicyStackingOverrides{
		PromotionWithHqCoupon:    input.PromotionWithHqCoupon,
		PromotionWithStoreCoupon: input.PromotionWithStoreCoupon,
		MemberPriceWithPromotion: input.MemberPriceWithPromotion,
		PointsWithPromotion:      input.PointsWithPromotion,
		PointsWithCoupon:         input.PointsWithCoupon,
	})
	if err != nil {
		return nil, err
	}
	return policyGraphQL(view)
}
