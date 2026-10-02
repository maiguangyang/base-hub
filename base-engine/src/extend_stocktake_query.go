package src

import (
	"context"
	"time"

	"base-engine/gen"
	"base-engine/src/services/storemerchandising"
)

func (r *QueryResolver) FranchiseStocktake(ctx context.Context, storeID, id string) (*gen.FranchiseStocktakeView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.Stocktake(ctx, principal, storeID, id)
	if err != nil {
		return nil, err
	}
	return stocktakeView(item)
}

func (r *QueryResolver) FranchiseStocktakes(ctx context.Context, storeID string, status *gen.StocktakeStatus, listingID, batchID *string, hasDifference *bool, from, to *time.Time, page, perPage int, includeHistory *bool) (*gen.FranchiseStocktakePage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	filter := storemerchandising.StocktakeFilter{}
	if status != nil {
		filter.Status = *status
	}
	if listingID != nil {
		filter.ListingID = *listingID
	}
	if batchID != nil {
		filter.BatchID = *batchID
	}
	filter.HasDifference, filter.From, filter.To, filter.IncludeHistory = hasDifference, from, to, includeHistory
	items, err := r.Services.Merchandising.Stocktakes(ctx, principal, storeID, filter, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchiseStocktakePage{Data: make([]*gen.FranchiseStocktakeView, 0, len(items.Data)),
		Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		view, err := stocktakeView(&item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (r *QueryResolver) FranchiseStocktakeBatchChoices(ctx context.Context, storeID, listingID string, page, perPage int) (*gen.FranchiseStocktakeBatchChoicePage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.Merchandising.StocktakeBatchChoices(ctx, principal, storeID, listingID, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchiseStocktakeBatchChoicePage{Data: make([]*gen.FranchiseStocktakeBatchChoice, 0, len(items.Data)),
		Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		result.Data = append(result.Data, &gen.FranchiseStocktakeBatchChoice{ID: item.ID, ListingID: item.ListingID,
			BatchNumber: item.BatchNumber, ExpiresAt: item.ExpiresAt})
	}
	return result, nil
}
