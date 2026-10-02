package src

import (
	"context"
	"math"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/customer"
)

func pointEntryGraphQL(entry *gen.CustomerPointEntry) (*gen.HqCustomerPointEntryView, error) {
	if entry.Delta < -math.MaxInt32 || entry.Delta > math.MaxInt32 {
		return nil, auth.NewError(auth.CodeInternalError)
	}
	return &gen.HqCustomerPointEntryView{
		ID: entry.ID, MemberID: entry.MemberID, Delta: int(entry.Delta),
		Source: entry.Source, OperationKind: entry.OperationKind, ReasonCode: entry.ReasonCode,
		Note: entry.Note, ReversesID: entry.ReversesID,
		CreatedAt: time.UnixMilli(entry.CreatedAt),
	}, nil
}

func (r *QueryResolver) HqCustomerPointEntries(ctx context.Context, memberID string, page, perPage int) (*gen.HqCustomerPointPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	found, err := r.Services.Customers.ListPointEntries(ctx, principal, memberID, int64(page), int64(perPage))
	if err != nil {
		return nil, err
	}
	if found.Total > math.MaxInt32 {
		return nil, auth.NewError(auth.CodeInternalError)
	}
	result := &gen.HqCustomerPointPage{
		Data:  make([]*gen.HqCustomerPointEntryView, 0, len(found.Data)),
		Total: int(found.Total), CurrentPage: page, PerPage: perPage,
	}
	for _, item := range found.Data {
		view, err := pointEntryGraphQL(&item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func pointMutationGraphQL(result *customer.PointMutationResult) (*gen.HqCustomerPointMutationView, error) {
	view, err := pointEntryGraphQL(result.CustomerPointEntry)
	if err != nil {
		return nil, err
	}
	if result.CurrentBalance < 0 || result.CurrentBalance > math.MaxInt32 {
		return nil, auth.NewError(auth.CodeInternalError)
	}
	return &gen.HqCustomerPointMutationView{Entry: view, CurrentBalance: int(result.CurrentBalance)}, nil
}

func (r *MutationResolver) HqGrantCustomerPoints(ctx context.Context, memberID string, points int, reasonCode gen.CustomerPointReasonCode, note, requestKey string) (*gen.HqCustomerPointMutationView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	entry, err := r.Services.Customers.GrantPoints(ctx, principal, memberID, int64(points), reasonCode, note, requestKey)
	if err != nil {
		return nil, err
	}
	return pointMutationGraphQL(entry)
}

func (r *MutationResolver) HqReverseCustomerPoints(ctx context.Context, entryID string, reasonCode gen.CustomerPointReasonCode, note, requestKey string) (*gen.HqCustomerPointMutationView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	entry, err := r.Services.Customers.ReversePoints(ctx, principal, entryID, reasonCode, note, requestKey)
	if err != nil {
		return nil, err
	}
	return pointMutationGraphQL(entry)
}

func (r *MutationResolver) HqCorrectCustomerPoints(ctx context.Context, memberID, evidenceReference, requestKey string) (*gen.HqCustomerPointMutationView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	entry, err := r.Services.Customers.CorrectPointMismatch(ctx, principal, memberID, evidenceReference, requestKey)
	if err != nil {
		return nil, err
	}
	return pointMutationGraphQL(entry)
}
