package src

import (
	"context"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/storemerchandising"
)

func (r *QueryResolver) FranchiseCatalog(ctx context.Context, storeID string, q, categoryID *string, listingEnabled, published *bool, page, perPage int) (*gen.FranchiseCatalogPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.Merchandising.Catalog(ctx, principal, storeID,
		storemerchandising.CatalogFilter{Q: q, CategoryID: categoryID, ListingEnabled: listingEnabled, Published: published}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchiseCatalogPage{Data: make([]*gen.FranchiseCatalogSkuView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	skuIDs := make([]string, 0, len(items.Data))
	for _, item := range items.Data {
		skuIDs = append(skuIDs, item.ID)
	}
	productIDs := make([]string, 0, len(items.Data))
	for _, item := range items.Data {
		productIDs = append(productIDs, item.ProductID)
	}
	productDetails, err := r.Services.Merchandising.CatalogProductDetails(ctx, principal, storeID, productIDs)
	if err != nil {
		return nil, err
	}
	listings, err := r.Services.Merchandising.CatalogListings(ctx, principal, storeID, skuIDs)
	if err != nil {
		return nil, err
	}
	for _, item := range items.Data {
		view, err := r.catalogSkuWithDetails(ctx, principal, storeID, item, productDetails[item.ProductID], listings[item.ID])
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (r *QueryResolver) catalogSkuWithDetails(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, item *gen.ProductSku, product storemerchandising.CatalogProductDetail, listing *gen.StoreListing) (*gen.FranchiseCatalogSkuView, error) {
	view := catalogSkuView(item, product)
	published, err := r.Services.Merchandising.CatalogPublished(ctx, item.ID)
	if err != nil {
		return nil, err
	}
	view.Published = published
	if listing != nil {
		view.ListingID, view.ListingEnabled = &listing.ID, &listing.Enabled
	}
	packages, err := r.Services.Merchandising.CatalogPackages(ctx, principal, storeID, item.ID)
	if err != nil {
		return nil, err
	}
	view.Packages, err = catalogPackageViews(packages)
	return view, err
}

func (r *QueryResolver) FranchiseCatalogCategories(ctx context.Context, storeID string, page, perPage int) (*gen.FranchiseCatalogCategoryPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.Merchandising.CatalogCategories(ctx, principal, storeID, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchiseCatalogCategoryPage{Data: make([]*gen.FranchiseCatalogCategoryView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		result.Data = append(result.Data, &gen.FranchiseCatalogCategoryView{ID: item.ID, Name: item.Name, ParentID: item.ParentID})
	}
	return result, nil
}

func (r *QueryResolver) FranchiseStoreListings(ctx context.Context, storeID string, page, perPage int) (*gen.FranchiseListingPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.Merchandising.Listings(ctx, principal, storeID, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchiseListingPage{Data: make([]*gen.FranchiseListingView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		result.Data = append(result.Data, listingView(item))
	}
	return result, nil
}

func (r *QueryResolver) FranchiseStoreOffers(ctx context.Context, storeID string, listingID, packageID *string, enabled *bool, page, perPage int) (*gen.FranchiseOfferPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	selectedListingID := ""
	if listingID != nil {
		selectedListingID = *listingID
	}
	items, err := r.Services.Merchandising.OffersFiltered(ctx, principal, storeID, selectedListingID,
		storemerchandising.OfferFilter{PackageID: packageID, Enabled: enabled}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchiseOfferPage{Data: make([]*gen.FranchiseOfferView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		view, err := r.franchiseOfferWithRevisions(ctx, principal, storeID, item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (r *QueryResolver) FranchiseStockBatches(ctx context.Context, storeID string, listingID *string, q *string, sellable *bool, page, perPage int) (*gen.FranchiseStockBatchPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	selectedListingID := ""
	if listingID != nil {
		selectedListingID = *listingID
	}
	items, err := r.Services.Merchandising.BatchesFiltered(ctx, principal, storeID, selectedListingID, storemerchandising.BatchFilter{Q: q, Sellable: sellable}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchiseStockBatchPage{Data: make([]*gen.FranchiseStockBatchView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		view, err := batchView(item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (r *QueryResolver) FranchiseOfferPriceRevisions(ctx context.Context, storeID, offerID string, q *string, from, to *time.Time, page, perPage int) (*gen.FranchisePriceRevisionPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.Merchandising.PriceHistoryFiltered(ctx, principal, storeID, offerID,
		storemerchandising.PriceHistoryFilter{Q: q, From: from, To: to}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchisePriceRevisionPage{Data: make([]*gen.FranchisePriceRevisionView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		oldPrice, err := catalogOptionalInt(item.PreviousPriceFen)
		if err != nil {
			return nil, err
		}
		newPrice, err := catalogInt(item.PriceFen)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, &gen.FranchisePriceRevisionView{ID: item.ID, OldPriceFen: oldPrice, NewPriceFen: newPrice, ReasonCode: item.ReasonCode, EffectiveAt: item.EffectiveAt})
	}
	return result, nil
}

func (r *QueryResolver) FranchiseStockMovements(ctx context.Context, storeID, batchID string, kind *gen.StockMovementKind, from, to *time.Time, page, perPage int) (*gen.FranchiseStockMovementPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.Merchandising.MovementsFiltered(ctx, principal, storeID, batchID, storemerchandising.MovementFilter{Kind: kind, From: from, To: to}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchiseStockMovementPage{Data: make([]*gen.FranchiseStockMovementView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		view, err := movementView(item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (r *QueryResolver) FranchiseStockPackageByBarcode(ctx context.Context, storeID, barcode string) (*gen.FranchiseStockPackageRecognition, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	return r.Services.Merchandising.StockPackageByBarcode(ctx, principal, storeID, barcode)
}
