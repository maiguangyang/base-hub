package src

import (
	"context"

	"base-engine/gen"
	"base-engine/src/services/storemerchandising"
)

func (r *MutationResolver) FranchiseSetStoreListing(ctx context.Context, storeID, skuID string, enabled bool) (*gen.FranchiseListingView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.SetListing(ctx, principal, storeID, skuID, enabled)
	if err != nil {
		return nil, err
	}
	return listingView(item), nil
}
func (r *MutationResolver) FranchiseSetStorePrice(ctx context.Context, listingID, packageID string, priceFen int, reasonCode string) (*gen.FranchiseOfferView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.SetPrice(ctx, principal, listingID, packageID, int64(priceFen), reasonCode)
	if err != nil {
		return nil, err
	}
	return offerView(item)
}
func (r *MutationResolver) FranchiseReceiveStock(ctx context.Context, input gen.FranchiseReceiveStockInput) (*gen.FranchiseStockMovementView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	value := storemerchandising.StockReceipt{Barcode: input.Barcode, StoreID: input.StoreID, ListingID: input.ListingID, PackageID: input.PackageID,
		ProducedAt: input.ProducedAt, ExpiresAt: input.ExpiresAt, SourceReference: input.SourceReference,
		Quantity: int64(input.Quantity), RequestKey: input.RequestKey}
	if input.BatchNumber != nil {
		value.BatchNumber = *input.BatchNumber
	}
	item, err := r.Services.Merchandising.ReceiveStock(ctx, principal, value)
	if err != nil {
		return nil, err
	}
	return movementView(item)
}
func (r *MutationResolver) FranchiseUnpackStock(ctx context.Context, input gen.FranchiseUnpackStockInput) (*gen.FranchiseStockMovementView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.UnpackStock(ctx, principal, storemerchandising.StockUnpack{StoreID: input.StoreID,
		BatchID: input.BatchID, SourcePackageID: input.SourcePackageID, Quantity: int64(input.Quantity), RequestKey: input.RequestKey})
	if err != nil {
		return nil, err
	}
	return movementView(item)
}
func (r *MutationResolver) FranchiseAdjustStock(ctx context.Context, input gen.FranchiseAdjustStockInput) (*gen.FranchiseStockMovementView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.AdjustStock(ctx, principal, storemerchandising.StockAdjustment{StoreID: input.StoreID,
		BatchID: input.BatchID, PackageID: input.PackageID, Delta: int64(input.Delta), Loss: input.Loss,
		ReasonCode: input.ReasonCode, RequestKey: input.RequestKey})
	if err != nil {
		return nil, err
	}
	return movementView(item)
}
