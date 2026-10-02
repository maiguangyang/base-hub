package storemerchandising

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
)

func (s *Service) CatalogCategories(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, page, perPage int) (*Page[*gen.ProductCategory], error) {
	if _, err := storeScope(s.db.WithContext(ctx), principal, storeID, scopedReadAction(principal, "franchiseProduct:read", "franchisePromotion:read"), authorization.AccessRead); err != nil {
		return nil, err
	}
	offset, limit, err := bounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&gen.ProductCategory{}).Joins("JOIN organizations ON organizations.id = product_categories.organization_id").
		Where("organizations.type = ? AND organizations.status = ?", gen.OrganizationTypeHeadquarters, gen.OrganizationStatusActive)
	result := &Page[*gen.ProductCategory]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Order("product_categories.sort_order, product_categories.name, product_categories.id").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}
