package src

import (
	"context"

	"base-engine/gen"
)

func (r *QueryResolver) FranchiseMerchandisingStores(ctx context.Context, page, perPage int) (*gen.FranchiseMerchandisingStorePage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.Merchandising.ScopedStores(ctx, principal, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchiseMerchandisingStorePage{Data: make([]*gen.FranchiseMerchandisingStoreView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		result.Data = append(result.Data, &gen.FranchiseMerchandisingStoreView{ID: item.ID, Name: item.Name})
	}
	return result, nil
}
