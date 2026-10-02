package src

import (
	"context"
	"math"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/customer"
)

func customerMemberGraphQL(view customer.MemberView) (*gen.HqCustomerMemberView, error) {
	if view.PointsBalance < 0 || view.PointsBalance > math.MaxInt32 ||
		view.PendingCouponCount > math.MaxInt32 {
		return nil, auth.NewError(auth.CodeInternalError)
	}
	return &gen.HqCustomerMemberView{
		ID: view.ID, MemberNumber: view.MemberNumber, PhoneMasked: view.PhoneMasked,
		Status: view.Status, PointsBalance: int(view.PointsBalance),
		PointsFrozen: view.PointsFrozen, PendingCouponCount: int(view.PendingCouponCount),
		CreatedAt: view.CreatedAt,
	}, nil
}

func createdCustomerGraphQL(member *gen.CustomerMember) (*gen.HqCustomerMemberView, error) {
	phoneMasked := ""
	if member.Phone != nil {
		phoneMasked = customer.MaskCNPhone(*member.Phone)
	}
	return customerMemberGraphQL(customer.MemberView{
		ID: member.ID, MemberNumber: member.MemberNumber, PhoneMasked: phoneMasked,
		Status: member.Status, PointsBalance: member.PointsBalance,
		PointsFrozen: member.PointsFrozen, PendingCouponCount: 0,
		CreatedAt: time.UnixMilli(member.CreatedAt),
	})
}

func (r *QueryResolver) HqCustomerMembers(ctx context.Context, phone *string, status *gen.CustomerMemberStatus, page, perPage int) (*gen.HqCustomerMemberPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	search := customer.MemberSearch{Status: status, Page: int64(page), PerPage: int64(perPage)}
	if phone != nil {
		search.Phone = *phone
	}
	found, err := r.Services.Customers.SearchMembers(ctx, principal, search)
	if err != nil {
		return nil, err
	}
	if found.Total > math.MaxInt32 {
		return nil, auth.NewError(auth.CodeInternalError)
	}
	result := &gen.HqCustomerMemberPage{
		Data:  make([]*gen.HqCustomerMemberView, 0, len(found.Data)),
		Total: int(found.Total), CurrentPage: page, PerPage: perPage,
	}
	for _, item := range found.Data {
		view, err := customerMemberGraphQL(item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (r *QueryResolver) HqCustomerMember(ctx context.Context, id string) (*gen.HqCustomerMemberView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	view, err := r.Services.Customers.GetMember(ctx, principal, id)
	if err != nil {
		return nil, err
	}
	return customerMemberGraphQL(*view)
}

func (r *QueryResolver) HqCustomerSensitivePhone(ctx context.Context, id string) (string, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return "", err
	}
	return r.Services.Customers.SensitivePhone(ctx, principal, id)
}

func (r *MutationResolver) HqCreateCustomerMember(ctx context.Context, input gen.HqCreateCustomerMemberInput) (*gen.HqCustomerMemberView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	member, err := r.Services.Customers.CreateMember(ctx, principal, customer.CreateMemberInput{
		Phone: input.Phone, RequestKey: input.RequestKey,
	})
	if err != nil {
		return nil, err
	}
	return createdCustomerGraphQL(member)
}

func (r *MutationResolver) HqSetCustomerMemberStatus(ctx context.Context, id string, status gen.CustomerMemberStatus) (*gen.HqCustomerMemberView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	member, err := r.Services.Customers.SetMemberStatus(ctx, principal, id, status)
	if err != nil {
		return nil, err
	}
	return customerMemberGraphQL(*member)
}

func (r *MutationResolver) HqRequestCustomerCancellation(ctx context.Context, id, identityEvidence, basisCode string) (*gen.HqCustomerMemberView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	member, err := r.Services.Customers.RequestCancellation(ctx, principal, id, identityEvidence, basisCode)
	if err != nil {
		return nil, err
	}
	return customerMemberGraphQL(*member)
}

func (r *MutationResolver) HqCompleteCustomerCancellation(ctx context.Context, id, dispositionReference string) (*gen.HqCustomerMemberView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	member, err := r.Services.Customers.CompleteCancellation(ctx, principal, id, dispositionReference)
	if err != nil {
		return nil, err
	}
	return createdCustomerGraphQL(member)
}
