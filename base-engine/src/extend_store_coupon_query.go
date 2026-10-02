package src

import (
	"context"

	"base-engine/gen"
	"base-engine/src/services/customer"
)

func (r *QueryResolver) FranchiseCouponMember(ctx context.Context, storeID, identifier string) (*gen.FranchiseCouponMemberView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	view, err := r.Services.Customers.LookupStoreCouponMember(ctx, principal, storeID, identifier)
	if err != nil || view == nil {
		return nil, err
	}
	return &gen.FranchiseCouponMemberView{ID: view.ID, MemberNumber: view.MemberNumber, PhoneMasked: view.PhoneMasked}, nil
}

func (r *QueryResolver) FranchiseCouponTemplates(ctx context.Context, storeID string, q *string, enabled *bool, page, perPage int) (*gen.FranchiseCouponTemplatePage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.Customers.ListStoreCouponTemplatesFiltered(ctx, principal, storeID, customer.CouponTemplateFilter{Q: q, Enabled: enabled}, int64(page), int64(perPage))
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchiseCouponTemplatePage{Data: make([]*gen.FranchiseCouponTemplateView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		view, err := storeCouponTemplateView(&item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (r *QueryResolver) FranchiseCouponGrants(ctx context.Context, storeID, templateID string, status *gen.CustomerCouponGrantStatus, page, perPage int) (*gen.FranchiseCouponGrantPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.Customers.ListStoreCouponGrantsFiltered(ctx, principal, storeID, templateID, customer.CouponGrantFilter{Status: status}, int64(page), int64(perPage))
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchiseCouponGrantPage{Data: make([]*gen.FranchiseCouponGrantView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		view, err := storeCouponGrantView(&item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}
