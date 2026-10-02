package src

import (
	"base-engine/gen"
	"base-engine/src/services/storemerchandising"
)

func catalogSkuView(item *gen.ProductSku, product storemerchandising.CatalogProductDetail) *gen.FranchiseCatalogSkuView {
	return &gen.FranchiseCatalogSkuView{ID: item.ID, ProductID: item.ProductID, ProductName: product.ProductName, CategoryID: product.CategoryID, Name: item.Name}
}
func catalogPackageViews(items []*gen.ProductPackage) ([]*gen.FranchiseCatalogPackageView, error) {
	result := make([]*gen.FranchiseCatalogPackageView, 0, len(items))
	for _, item := range items {
		quantity, err := catalogOptionalInt(item.ContainsQuantity)
		if err != nil {
			return nil, err
		}
		price, err := catalogOptionalInt(item.SuggestedPriceFen)
		if err != nil {
			return nil, err
		}
		version, err := catalogInt(item.PackageSetVersion)
		if err != nil {
			return nil, err
		}
		result = append(result, &gen.FranchiseCatalogPackageView{ID: item.ID, Name: item.Name, Enabled: item.Enabled,
			Barcode: item.Barcode, ContainsPackageID: item.ContainsPackageID,
			ContainsQuantity: quantity, SuggestedPriceFen: price, PackageSetVersion: version})
	}
	return result, nil
}
func listingView(item *gen.StoreListing) *gen.FranchiseListingView {
	return &gen.FranchiseListingView{ID: item.ID, StoreID: item.StoreID, SkuID: item.SkuID, Enabled: item.Enabled, SelectedAt: item.SelectedAt}
}
func offerView(item *gen.StorePackageOffer) (*gen.FranchiseOfferView, error) {
	price, err := catalogInt(item.PriceFen)
	if err != nil {
		return nil, err
	}
	return &gen.FranchiseOfferView{ID: item.ID, ListingID: item.ListingID, PackageID: item.PackageID, PriceFen: price, Enabled: item.Enabled}, nil
}
func batchView(item storemerchandising.BatchView) (*gen.FranchiseStockBatchView, error) {
	packages, err := catalogPackageViews(item.Packages)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchiseStockBatchView{ID: item.Batch.ID, ListingID: item.Batch.ListingID, SkuID: item.SkuID, ProductName: item.ProductName, SkuName: item.SkuName, Packages: packages, BatchNumber: item.Batch.BatchNumber,
		ProducedAt: item.Batch.ProducedAt, ExpiresAt: item.Batch.ExpiresAt, Sellable: item.Sellable,
		Balances: make([]*gen.FranchiseStockBalanceView, 0, len(item.Balances))}
	for _, balance := range item.Balances {
		quantity, err := catalogInt(balance.Quantity)
		if err != nil {
			return nil, err
		}
		result.Balances = append(result.Balances, &gen.FranchiseStockBalanceView{PackageID: balance.PackageID, Quantity: quantity})
	}
	return result, nil
}
func movementView(item *gen.StoreStockMovement) (*gen.FranchiseStockMovementView, error) {
	source, err := catalogOptionalInt(item.SourceQuantity)
	if err != nil {
		return nil, err
	}
	target, err := catalogOptionalInt(item.TargetQuantity)
	if err != nil {
		return nil, err
	}
	factor, err := catalogOptionalInt(item.FactorSnapshot)
	if err != nil {
		return nil, err
	}
	version, err := catalogInt(item.PackageSetVersion)
	if err != nil {
		return nil, err
	}
	return &gen.FranchiseStockMovementView{ID: item.ID, StoreID: item.StoreID, BatchID: item.BatchID, Kind: item.Kind,
		RequestKey: item.RequestKey, SourcePackageID: item.SourcePackageID, TargetPackageID: item.TargetPackageID,
		SourceQuantity: source, TargetQuantity: target, FactorSnapshot: factor, PackageSetVersion: version,
		ReasonCode: item.ReasonCode, OccurredAt: item.OccurredAt}, nil
}
