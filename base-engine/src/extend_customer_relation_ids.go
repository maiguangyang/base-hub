package src

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
)

// Generated relationship ID projections bypass the generated data handlers.
// Keep them closed until the customer read surface has scoped resolvers.
func denyCustomerRelationIDs(ctx context.Context) ([]string, error) {
	if _, err := requireResolverPrincipal(ctx); err != nil {
		return nil, err
	}
	return nil, auth.NewError(auth.CodePermissionDenied)
}

func (r *OrganizationResolver) CustomerMembersIds(ctx context.Context, _ *gen.Organization) ([]string, error) {
	return denyCustomerRelationIDs(ctx)
}

func (r *OrganizationResolver) CustomerBenefitPoliciesIds(ctx context.Context, _ *gen.Organization) ([]string, error) {
	return denyCustomerRelationIDs(ctx)
}

func (r *OrganizationResolver) CustomerDailyPointGrantBudgetsIds(ctx context.Context, _ *gen.Organization) ([]string, error) {
	return denyCustomerRelationIDs(ctx)
}

func (r *OrganizationResolver) CustomerPointEntriesIds(ctx context.Context, _ *gen.Organization) ([]string, error) {
	return denyCustomerRelationIDs(ctx)
}

func (r *OrganizationResolver) CustomerCouponTemplatesIds(ctx context.Context, _ *gen.Organization) ([]string, error) {
	return denyCustomerRelationIDs(ctx)
}

func (r *CustomerMemberResolver) PointEntriesIds(ctx context.Context, _ *gen.CustomerMember) ([]string, error) {
	return denyCustomerRelationIDs(ctx)
}

func (r *CustomerMemberResolver) CouponGrantsIds(ctx context.Context, _ *gen.CustomerMember) ([]string, error) {
	return denyCustomerRelationIDs(ctx)
}

func (r *CustomerCouponTemplateResolver) GrantsIds(ctx context.Context, _ *gen.CustomerCouponTemplate) ([]string, error) {
	return denyCustomerRelationIDs(ctx)
}
