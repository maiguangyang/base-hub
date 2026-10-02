package src

import (
	"context"

	"base-engine/gen"
	"base-engine/src/services/productcatalog"
)

func packageTemplateInput(input gen.HqProductPackageTemplateInput) productcatalog.PackageTemplateInput {
	value := productcatalog.PackageTemplateInput{Name: input.Name, ContainsPackageID: input.ContainsPackageID}
	if input.ContainsQuantity != nil {
		value.ContainsQuantity = int64(*input.ContainsQuantity)
	}
	return value
}

func (r *MutationResolver) HqCreateProductPackageTemplate(ctx context.Context, input gen.HqProductPackageTemplateInput) (*gen.HqProductPackageTemplateView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.CreatePackageTemplate(ctx, principal, packageTemplateInput(input))
	if err != nil {
		return nil, err
	}
	return packageTemplateView(item)
}

func (r *MutationResolver) HqUpdateProductPackageTemplate(ctx context.Context, id string, input gen.HqProductPackageTemplateInput) (*gen.HqProductPackageTemplateView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.UpdatePackageTemplate(ctx, principal, id, packageTemplateInput(input))
	if err != nil {
		return nil, err
	}
	return packageTemplateView(item)
}

func (r *MutationResolver) HqSetProductPackageTemplateEnabled(ctx context.Context, id string, enabled bool) (*gen.HqProductPackageTemplateView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.SetPackageTemplateEnabled(ctx, principal, id, enabled)
	if err != nil {
		return nil, err
	}
	return packageTemplateView(item)
}

func (r *MutationResolver) HqDeleteProductPackageTemplate(ctx context.Context, id string) (bool, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return false, err
	}
	if err := r.Services.ProductCatalog.DeletePackageTemplate(ctx, principal, id); err != nil {
		return false, err
	}
	return true, nil
}
