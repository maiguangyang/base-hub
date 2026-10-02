package src

import (
	"context"

	"base-engine/gen"
	"base-engine/src/services/storemerchandising"
)

func (r *MutationResolver) FranchiseCreateStocktake(ctx context.Context, input gen.FranchiseCreateStocktakeInput) (*gen.FranchiseStocktakeView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.CreateStocktake(ctx, principal, storemerchandising.CreateStocktakeInput{
		StoreID: input.StoreID, ListingIDs: input.ListingIds, BatchIDs: input.BatchIds, RequestKey: input.RequestKey})
	if err != nil {
		return nil, err
	}
	return stocktakeView(item)
}

func (r *MutationResolver) FranchiseAddStocktakeLine(ctx context.Context, storeID, id, batchID, packageID string) (*gen.FranchiseStocktakeView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.AddStocktakeLine(ctx, principal, storeID, id, batchID, packageID)
	if err != nil {
		return nil, err
	}
	return stocktakeView(item)
}

func (r *MutationResolver) FranchiseRecordStocktakeLine(ctx context.Context, storeID, id, lineID string, quantity int) (*gen.FranchiseStocktakeView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.RecordStocktakeLine(ctx, principal, storeID, id, lineID, int64(quantity))
	if err != nil {
		return nil, err
	}
	return stocktakeView(item)
}

func (r *MutationResolver) FranchiseSubmitStocktake(ctx context.Context, storeID, id string) (*gen.FranchiseStocktakeView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.SubmitStocktake(ctx, principal, storeID, id)
	if err != nil {
		return nil, err
	}
	return stocktakeView(item)
}

func (r *MutationResolver) FranchiseSetStocktakeReason(ctx context.Context, storeID, id, lineID, reasonCode string, note *string) (*gen.FranchiseStocktakeView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	comment := ""
	if note != nil {
		comment = *note
	}
	item, err := r.Services.Merchandising.SetStocktakeReason(ctx, principal, storeID, id, lineID, reasonCode, comment)
	if err != nil {
		return nil, err
	}
	return stocktakeView(item)
}

func (r *MutationResolver) FranchiseReturnStocktake(ctx context.Context, storeID, id string) (*gen.FranchiseStocktakeView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.ReturnStocktake(ctx, principal, storeID, id)
	if err != nil {
		return nil, err
	}
	return stocktakeView(item)
}

func (r *MutationResolver) FranchiseCancelStocktake(ctx context.Context, storeID, id string) (*gen.FranchiseStocktakeView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.CancelStocktake(ctx, principal, storeID, id)
	if err != nil {
		return nil, err
	}
	return stocktakeView(item)
}

func (r *MutationResolver) FranchisePostStocktake(ctx context.Context, storeID, id string) (*gen.FranchiseStocktakeView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.PostStocktake(ctx, principal, storeID, id)
	if err != nil {
		return nil, err
	}
	return stocktakeView(item)
}
