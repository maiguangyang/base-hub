package src

import (
	"context"

	"base-engine/gen"
	"base-engine/src/services/productcatalog"
)

func brandView(item *gen.ProductBrand) *gen.HqProductBrandView {
	return &gen.HqProductBrandView{ID: item.ID, Name: item.Name, Enabled: item.Enabled}
}

func (r *QueryResolver) HqProductBrands(ctx context.Context, q *string, enabled *bool, page, perPage int) (*gen.HqProductBrandPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil { return nil, err }
	items, err := r.Services.ProductCatalog.Brands(ctx, principal, productcatalog.CatalogFilter{Q: q, Enabled: enabled}, page, perPage)
	if err != nil { return nil, err }
	total, err := catalogInt(items.Total)
	if err != nil { return nil, err }
	result := &gen.HqProductBrandPage{Data: make([]*gen.HqProductBrandView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data { result.Data = append(result.Data, brandView(item)) }
	return result, nil
}

func (r *MutationResolver) HqCreateProductBrand(ctx context.Context, name string) (*gen.HqProductBrandView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil { return nil, err }
	item, err := r.Services.ProductCatalog.CreateBrand(ctx, principal, name)
	if err != nil { return nil, err }
	return brandView(item), nil
}

func (r *MutationResolver) HqUpdateProductBrand(ctx context.Context, id, name string) (*gen.HqProductBrandView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil { return nil, err }
	item, err := r.Services.ProductCatalog.UpdateBrand(ctx, principal, id, name)
	if err != nil { return nil, err }
	return brandView(item), nil
}

func (r *MutationResolver) HqSetProductBrandEnabled(ctx context.Context, id string, enabled bool) (*gen.HqProductBrandView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil { return nil, err }
	item, err := r.Services.ProductCatalog.SetBrandEnabled(ctx, principal, id, enabled)
	if err != nil { return nil, err }
	return brandView(item), nil
}

func (r *MutationResolver) HqDeleteProductBrand(ctx context.Context, id string) (bool, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil { return false, err }
	if err := r.Services.ProductCatalog.DeleteBrand(ctx, principal, id); err != nil { return false, err }
	return true, nil
}

func (r *MutationResolver) HqDeleteProduct(ctx context.Context, id string) (bool, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil { return false, err }
	if err := r.Services.ProductCatalog.DeleteProduct(ctx, principal, id); err != nil { return false, err }
	return true, nil
}
