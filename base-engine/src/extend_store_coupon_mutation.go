package src

import (
	"context"

	"base-engine/gen"
	"base-engine/src/services/customer"
)

func (r *MutationResolver) FranchiseCreateCouponTemplate(ctx context.Context, input gen.FranchiseCreateCouponTemplateInput) (*gen.FranchiseCouponTemplateView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	value := customer.FixedAmountTemplateInput{Code: input.Code, Title: input.Title, RequestKey: input.RequestKey,
		AmountFen: int64(input.AmountFen), MinSpendFen: int64(input.MinSpendFen), EffectiveAt: couponInputMillis(input.EffectiveAt),
		DistributionEndsAt: couponInputMillis(input.DistributionEndsAt), DaysAfterActivation: int64(input.DaysAfterActivation),
		PerMemberLimit: int64(input.PerMemberLimit), TotalIssueLimit: couponInputLimit(input.TotalIssueLimit), Enabled: input.Enabled}
	item, err := r.Services.Customers.CreateStoreCouponTemplate(ctx, principal, input.StoreID, value)
	if err != nil {
		return nil, err
	}
	return storeCouponTemplateView(item)
}
func (r *MutationResolver) FranchiseSetCouponTemplateEnabled(ctx context.Context, storeID, templateID string, enabled bool) (*gen.FranchiseCouponTemplateView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Customers.SetStoreCouponTemplateEnabled(ctx, principal, storeID, templateID, enabled)
	if err != nil {
		return nil, err
	}
	return storeCouponTemplateView(item)
}
func (r *MutationResolver) FranchiseGrantCoupon(ctx context.Context, storeID, templateID, memberID, requestKey string) (*gen.FranchiseCouponGrantView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Customers.GrantStoreCoupon(ctx, principal, storeID, templateID, memberID, requestKey)
	if err != nil {
		return nil, err
	}
	return storeCouponGrantView(item)
}
func (r *MutationResolver) FranchiseRevokeCoupon(ctx context.Context, storeID, grantID, reasonCode string) (*gen.FranchiseCouponGrantView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Customers.RevokeStoreCoupon(ctx, principal, storeID, grantID, reasonCode)
	if err != nil {
		return nil, err
	}
	return storeCouponGrantView(item)
}
