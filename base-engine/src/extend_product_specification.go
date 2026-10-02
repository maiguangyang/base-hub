package src

import (
	"context"

	"base-engine/gen"
	"base-engine/src/services/productcatalog"
)

func specificationView(item *gen.SpecificationDefinition) *gen.HqSpecificationView {
	return &gen.HqSpecificationView{ID: item.ID, Name: item.Name, Enabled: item.Enabled}
}

func specificationValueView(item *gen.SpecificationValue) *gen.HqSpecificationValueView {
	return &gen.HqSpecificationValueView{ID: item.ID, SpecificationID: item.SpecificationID, Name: item.Name, Enabled: item.Enabled}
}

func (r *QueryResolver) HqSpecifications(ctx context.Context, q *string, enabled *bool, page, perPage int) (*gen.HqSpecificationPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.ProductCatalog.Specifications(ctx, principal, productcatalog.CatalogFilter{Q: q, Enabled: enabled}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.HqSpecificationPage{Data: make([]*gen.HqSpecificationView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		result.Data = append(result.Data, specificationView(item))
	}
	return result, nil
}

func (r *QueryResolver) HqSpecificationValues(ctx context.Context, specificationID *string, q *string, enabled *bool, page, perPage int) (*gen.HqSpecificationValuePage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.ProductCatalog.SpecificationValues(ctx, principal, specificationID, productcatalog.CatalogFilter{Q: q, Enabled: enabled}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.HqSpecificationValuePage{Data: make([]*gen.HqSpecificationValueView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		result.Data = append(result.Data, specificationValueView(item))
	}
	return result, nil
}

func (r *MutationResolver) HqCreateSpecification(ctx context.Context, name string) (*gen.HqSpecificationView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.CreateSpecification(ctx, principal, name)
	if err != nil {
		return nil, err
	}
	return specificationView(item), nil
}

func (r *MutationResolver) HqUpdateSpecification(ctx context.Context, id, name string) (*gen.HqSpecificationView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.UpdateSpecification(ctx, principal, id, name)
	if err != nil {
		return nil, err
	}
	return specificationView(item), nil
}

func (r *MutationResolver) HqSetSpecificationEnabled(ctx context.Context, id string, enabled bool) (*gen.HqSpecificationView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.SetSpecificationEnabled(ctx, principal, id, enabled)
	if err != nil {
		return nil, err
	}
	return specificationView(item), nil
}

func (r *MutationResolver) HqDeleteSpecification(ctx context.Context, id string) (bool, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return false, err
	}
	if err := r.Services.ProductCatalog.DeleteSpecification(ctx, principal, id); err != nil {
		return false, err
	}
	return true, nil
}

func (r *MutationResolver) HqCreateSpecificationValue(ctx context.Context, specificationID, name string) (*gen.HqSpecificationValueView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.CreateSpecificationValue(ctx, principal, specificationID, name)
	if err != nil {
		return nil, err
	}
	return specificationValueView(item), nil
}

func (r *MutationResolver) HqUpdateSpecificationValue(ctx context.Context, id, name string) (*gen.HqSpecificationValueView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.UpdateSpecificationValue(ctx, principal, id, name)
	if err != nil {
		return nil, err
	}
	return specificationValueView(item), nil
}

func (r *MutationResolver) HqSetSpecificationValueEnabled(ctx context.Context, id string, enabled bool) (*gen.HqSpecificationValueView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.ProductCatalog.SetSpecificationValueEnabled(ctx, principal, id, enabled)
	if err != nil {
		return nil, err
	}
	return specificationValueView(item), nil
}

func (r *MutationResolver) HqReorderSpecificationValues(ctx context.Context, specificationID string, orderedIds []string) (bool, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return false, err
	}
	if err := r.Services.ProductCatalog.ReorderSpecificationValues(ctx, principal, specificationID, orderedIds); err != nil {
		return false, err
	}
	return true, nil
}

func (r *MutationResolver) HqDeleteSpecificationValue(ctx context.Context, id string) (bool, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return false, err
	}
	if err := r.Services.ProductCatalog.DeleteSpecificationValue(ctx, principal, id); err != nil {
		return false, err
	}
	return true, nil
}
