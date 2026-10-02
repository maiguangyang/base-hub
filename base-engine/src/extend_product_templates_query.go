package src

import (
	"context"

	"base-engine/gen"
	"base-engine/src/services/productcatalog"
)

func packageTemplateView(item *gen.ProductPackageTemplate) (*gen.HqProductPackageTemplateView, error) {
	quantity, err := catalogOptionalInt(item.ContainsQuantity)
	if err != nil {
		return nil, err
	}
	return &gen.HqProductPackageTemplateView{ID: item.ID, Name: item.Name,
		ContainsPackageID: item.ContainsPackageID, ContainsQuantity: quantity, Enabled: item.Enabled}, nil
}

func (r *QueryResolver) HqProductPackageTemplates(ctx context.Context, q *string, enabled *bool, containsPackageID *string, page, perPage int) (*gen.HqProductPackageTemplatePage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.ProductCatalog.PackageTemplates(ctx, principal, productcatalog.CatalogFilter{Q: q, Enabled: enabled, ContainsPackageID: containsPackageID}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.HqProductPackageTemplatePage{Data: make([]*gen.HqProductPackageTemplateView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		view, err := packageTemplateView(item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}
