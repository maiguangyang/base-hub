package src

import (
	"context"

	"base-engine/gen"
	"base-engine/src/services/productcatalog"
)

func (r *QueryResolver) HqProductCategories(ctx context.Context, q *string, enabled *bool, parentID *string, page, perPage int) (*gen.HqProductCategoryPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.ProductCatalog.Categories(ctx, principal, productcatalog.CatalogFilter{Q: q, Enabled: enabled, ParentID: parentID}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.HqProductCategoryPage{Data: make([]*gen.HqProductCategoryView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		view, err := categoryView(item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (r *QueryResolver) HqProducts(ctx context.Context, categoryID, brandID, q *string, enabled *bool, page, perPage int) (*gen.HqProductPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.ProductCatalog.Products(ctx, principal, categoryID, productcatalog.CatalogFilter{Q: q, Enabled: enabled, BrandID: brandID}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.HqProductPage{Data: make([]*gen.HqProductView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		result.Data = append(result.Data, productView(item))
	}
	return result, nil
}

func (r *QueryResolver) HqProductSkus(ctx context.Context, productID string, q *string, enabled *bool, page, perPage int) (*gen.HqProductSkuPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.ProductCatalog.Skus(ctx, principal, productID, productcatalog.CatalogFilter{Q: q, Enabled: enabled}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.HqProductSkuPage{Data: make([]*gen.HqProductSkuView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		view, err := skuView(item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (r *QueryResolver) HqProductPackages(ctx context.Context, skuID string, q *string, enabled *bool, page, perPage int) (*gen.HqProductPackagePage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.ProductCatalog.Packages(ctx, principal, skuID, productcatalog.CatalogFilter{Q: q, Enabled: enabled}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.HqProductPackagePage{Data: make([]*gen.HqProductPackageView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		view, err := packageView(item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}
