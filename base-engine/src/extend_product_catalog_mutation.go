package src

import (
	"context"

	"base-engine/gen"
	"base-engine/src/services/productcatalog"
)

func (r *MutationResolver) HqSetProductPackageBarcode(ctx context.Context, id, barcode string) (*gen.HqProductPackageView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.SetPackageBarcode(ctx, principal, id, barcode)
	if err != nil {
		return nil, err
	}
	return packageView(item)
}

func (r *MutationResolver) HqCreateProductCategory(ctx context.Context, input gen.HqCreateProductCategoryInput) (*gen.HqProductCategoryView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.CreateCategory(ctx, principal, productcatalog.CategoryInput{Name: input.Name, ParentID: input.ParentID})
	if err != nil {
		return nil, err
	}
	return categoryView(item)
}
func (r *MutationResolver) HqUpdateProductCategory(ctx context.Context, id, name string, parentID *string) (*gen.HqProductCategoryView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.UpdateCategory(ctx, principal, id, name, parentID)
	if err != nil {
		return nil, err
	}
	return categoryView(item)
}
func (r *MutationResolver) HqSetProductCategoryEnabled(ctx context.Context, id string, enabled bool) (*gen.HqProductCategoryView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.SetCategoryEnabled(ctx, principal, id, enabled)
	if err != nil {
		return nil, err
	}
	return categoryView(item)
}
func (r *MutationResolver) HqCreateProduct(ctx context.Context, input gen.HqCreateProductInput) (*gen.HqProductView, error) {
	if err := validateProductInputPointers(input); err != nil {
		return nil, err
	}
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.CreateProduct(ctx, principal, productInput(input))
	if err != nil {
		return nil, err
	}
	return productView(item), nil
}
func (r *MutationResolver) HqUpdateProduct(ctx context.Context, id string, input gen.HqCreateProductInput) (*gen.HqProductView, error) {
	if err := validateProductInputPointers(input); err != nil {
		return nil, err
	}
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.UpdateProduct(ctx, principal, id, productInput(input))
	if err != nil {
		return nil, err
	}
	return productView(item), nil
}
func (r *MutationResolver) HqSetProductEnabled(ctx context.Context, id string, enabled bool) (*gen.HqProductView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.SetProductEnabled(ctx, principal, id, enabled)
	if err != nil {
		return nil, err
	}
	return productView(item), nil
}
func (r *MutationResolver) HqSetProductMainImage(ctx context.Context, productID, attachmentID string) (*gen.HqProductView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.SetProductMainImage(ctx, principal, productID, attachmentID)
	if err != nil {
		return nil, err
	}
	return productView(item), nil
}
func (r *MutationResolver) HqRemoveProductMainImage(ctx context.Context, productID string) (*gen.HqProductView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.RemoveProductMainImage(ctx, principal, productID)
	if err != nil {
		return nil, err
	}
	return productView(item), nil
}
func (r *MutationResolver) HqSetProductSkuEnabled(ctx context.Context, id string, enabled bool) (*gen.HqProductSkuView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.SetSkuEnabled(ctx, principal, id, enabled)
	if err != nil {
		return nil, err
	}
	return skuView(item)
}
func (r *MutationResolver) HqCreateProductPackage(ctx context.Context, input gen.HqCreateProductPackageInput) (*gen.HqProductPackageView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	value := productcatalog.PackageInput{SkuID: input.SkuID, Name: input.Name, ContainsPackageID: input.ContainsPackageID}
	if input.Barcode != nil {
		value.Barcode = *input.Barcode
	}
	if input.ContainsQuantity != nil {
		value.ContainsQuantity = int64(*input.ContainsQuantity)
	}
	if input.SuggestedPriceFen != nil {
		price := int64(*input.SuggestedPriceFen)
		value.SuggestedPriceFen = &price
	}
	if input.PackageSetVersion != nil {
		value.PackageSetVersion = int64(*input.PackageSetVersion)
	}
	item, err := r.Services.ProductCatalog.CreatePackage(ctx, principal, value)
	if err != nil {
		return nil, err
	}
	return packageView(item)
}
func (r *MutationResolver) HqRetireProductPackage(ctx context.Context, id string) (*gen.HqProductPackageView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.RetirePackage(ctx, principal, id)
	if err != nil {
		return nil, err
	}
	return packageView(item)
}

func (r *MutationResolver) HqPublishProductPackageSet(ctx context.Context, skuID string, version int) (*gen.HqProductPackageView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.PublishPackageSet(ctx, principal, skuID, int64(version))
	if err != nil {
		return nil, err
	}
	return packageView(item)
}

func (r *MutationResolver) HqUpdateProductPackageMetadata(ctx context.Context, id, name string, suggestedPriceFen *int) (*gen.HqProductPackageView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	var price *int64
	if suggestedPriceFen != nil {
		value := int64(*suggestedPriceFen)
		price = &value
	}
	item, err := r.Services.ProductCatalog.UpdatePackageMetadata(ctx, principal, id, name, price)
	if err != nil {
		return nil, err
	}
	return packageView(item)
}
