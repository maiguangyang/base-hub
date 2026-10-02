package src

import (
	"context"
	"math"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/customer"
)

func couponTemplateGraphQL(template *gen.CustomerCouponTemplate) (*gen.HqCustomerCouponTemplateView, error) {
	values := []int64{
		template.AmountFen, template.MinSpendFen, template.DaysAfterActivation,
		template.PerMemberLimit, template.TotalIssueLimit, template.IssuedCount,
	}
	for _, value := range values {
		if value < 0 || value > math.MaxInt32 {
			return nil, auth.NewError(auth.CodeInternalError)
		}
	}
	return &gen.HqCustomerCouponTemplateView{
		ID: template.ID, Code: template.Code, Title: template.Title,
		AmountFen: int(template.AmountFen), MinSpendFen: int(template.MinSpendFen),
		DaysAfterActivation: int(template.DaysAfterActivation),
		EffectiveAt:         time.UnixMilli(template.EffectiveAt),
		DistributionEndsAt:  couponTimePointer(template.DistributionEndsAt),
		PerMemberLimit:      int(template.PerMemberLimit), TotalIssueLimit: int(template.TotalIssueLimit),
		IssuedCount: int(template.IssuedCount), Enabled: template.Enabled,
		CreatedAt: time.UnixMilli(template.CreatedAt),
	}, nil
}

func couponGrantGraphQL(grant *gen.CustomerCouponGrant) (*gen.HqCustomerCouponGrantView, error) {
	if grant.AmountFen < 0 || grant.AmountFen > math.MaxInt32 ||
		grant.MinSpendFen < 0 || grant.MinSpendFen > math.MaxInt32 ||
		grant.DaysAfterActivation < 0 || grant.DaysAfterActivation > math.MaxInt32 {
		return nil, auth.NewError(auth.CodeInternalError)
	}
	return &gen.HqCustomerCouponGrantView{
		ID: grant.ID, TemplateID: grant.TemplateID, MemberID: grant.MemberID,
		Status: grant.Status, AmountFen: int(grant.AmountFen), MinSpendFen: int(grant.MinSpendFen),
		DaysAfterActivation: int(grant.DaysAfterActivation), IssuedAt: time.UnixMilli(grant.IssuedAt),
		ActivatedAt: couponTimePointer(grant.ActivatedAt), ExpiresAt: couponTimePointer(grant.ExpiresAt),
		RevokedAt: couponTimePointer(grant.RevokedAt),
	}, nil
}

func couponTimePointer(value *int64) *time.Time {
	if value == nil {
		return nil
	}
	converted := time.UnixMilli(*value)
	return &converted
}

func couponInputMillis(value *time.Time) *int64 {
	if value == nil {
		return nil
	}
	converted := value.UnixMilli()
	return &converted
}

func couponInputLimit(value *int) int64 {
	if value == nil {
		return 0
	}
	return int64(*value)
}

func (r *QueryResolver) HqCustomerCouponTemplates(ctx context.Context, q *string, enabled *bool, page, perPage int) (*gen.HqCustomerCouponTemplatePage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	found, err := r.Services.Customers.ListCouponTemplatesFiltered(ctx, principal, customer.CouponTemplateFilter{Q: q, Enabled: enabled}, int64(page), int64(perPage))
	if err != nil {
		return nil, err
	}
	if found.Total > math.MaxInt32 {
		return nil, auth.NewError(auth.CodeInternalError)
	}
	result := &gen.HqCustomerCouponTemplatePage{
		Data:  make([]*gen.HqCustomerCouponTemplateView, 0, len(found.Data)),
		Total: int(found.Total), CurrentPage: page, PerPage: perPage,
	}
	for _, item := range found.Data {
		view, err := couponTemplateGraphQL(&item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (r *QueryResolver) HqCustomerCouponGrants(ctx context.Context, memberID string, templateID *string, status *gen.CustomerCouponGrantStatus, page, perPage int) (*gen.HqCustomerCouponGrantPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	found, err := r.Services.Customers.ListCouponGrantsFiltered(ctx, principal, memberID, customer.CouponGrantFilter{TemplateID: templateID, Status: status}, int64(page), int64(perPage))
	if err != nil {
		return nil, err
	}
	if found.Total > math.MaxInt32 {
		return nil, auth.NewError(auth.CodeInternalError)
	}
	result := &gen.HqCustomerCouponGrantPage{
		Data:  make([]*gen.HqCustomerCouponGrantView, 0, len(found.Data)),
		Total: int(found.Total), CurrentPage: page, PerPage: perPage,
	}
	for _, item := range found.Data {
		view, err := couponGrantGraphQL(&item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (r *MutationResolver) HqCreateCustomerCouponTemplate(ctx context.Context, input gen.HqCreateCustomerCouponTemplateInput) (*gen.HqCustomerCouponTemplateView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	template, err := r.Services.Customers.CreateCouponTemplate(ctx, principal, customer.FixedAmountTemplateInput{
		Code: input.Code, Title: input.Title, RequestKey: input.RequestKey, AmountFen: int64(input.AmountFen),
		MinSpendFen: int64(input.MinSpendFen), EffectiveAt: couponInputMillis(input.EffectiveAt),
		DistributionEndsAt: couponInputMillis(input.DistributionEndsAt), DaysAfterActivation: int64(input.DaysAfterActivation),
		PerMemberLimit: int64(input.PerMemberLimit), TotalIssueLimit: couponInputLimit(input.TotalIssueLimit),
		Enabled: input.Enabled,
	})
	if err != nil {
		return nil, err
	}
	return couponTemplateGraphQL(template)
}

func (r *MutationResolver) HqSetCustomerCouponTemplateEnabled(ctx context.Context, templateID string, enabled bool) (*gen.HqCustomerCouponTemplateView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	template, err := r.Services.Customers.SetCouponTemplateEnabled(ctx, principal, templateID, enabled)
	if err != nil {
		return nil, err
	}
	return couponTemplateGraphQL(template)
}

func (r *MutationResolver) HqGrantCustomerCoupon(ctx context.Context, templateID, memberID, requestKey string) (*gen.HqCustomerCouponGrantView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	grant, err := r.Services.Customers.GrantCoupon(ctx, principal, templateID, memberID, requestKey)
	if err != nil {
		return nil, err
	}
	return couponGrantGraphQL(grant)
}

func (r *MutationResolver) HqRevokeCustomerCoupon(ctx context.Context, grantID, reasonCode string) (*gen.HqCustomerCouponGrantView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	grant, err := r.Services.Customers.RevokeCoupon(ctx, principal, grantID, reasonCode)
	if err != nil {
		return nil, err
	}
	return couponGrantGraphQL(grant)
}
