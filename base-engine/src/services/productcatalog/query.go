package productcatalog

import (
	"context"
	"strings"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type CatalogFilter struct {
	Q                 *string
	Enabled           *bool
	ParentID          *string
	BrandID           *string
	ContainsPackageID *string
}

func applyCatalogFilter(query *gorm.DB, table string, filter CatalogFilter) *gorm.DB {
	if filter.Q != nil && strings.TrimSpace(*filter.Q) != "" {
		query = query.Where(table+".name LIKE ?", "%"+strings.TrimSpace(*filter.Q)+"%")
	}
	if filter.Enabled != nil {
		query = query.Where(table+".enabled = ?", *filter.Enabled)
	}
	return query
}

type CatalogPage[T any] struct {
	Data          []T
	Total         int64
	Page, PerPage int
}

func pageBounds(page, perPage int) (int, int, error) {
	if page < 1 || perPage < 1 || perPage > 100 {
		return 0, 0, auth.NewError(auth.CodeValidationFailed)
	}
	return (page - 1) * perPage, perPage, nil
}

func (s *Service) Categories(ctx context.Context, principal *auth.WorkspacePrincipal, filter CatalogFilter, page, perPage int) (*CatalogPage[*gen.ProductCategory], error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:read", authorization.AccessRead)
	if err != nil {
		return nil, err
	}
	offset, limit, err := pageBounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := applyCatalogFilter(s.db.WithContext(ctx).Model(&gen.ProductCategory{}).Where("organization_id = ?", hqID), "product_categories", filter)
	if filter.ParentID != nil {
		query = query.Where("product_categories.parent_id = ?", *filter.ParentID)
	}
	result := &CatalogPage[*gen.ProductCategory]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Order("sort_order, name, id").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) Products(ctx context.Context, principal *auth.WorkspacePrincipal, categoryID *string, filter CatalogFilter, page, perPage int) (*CatalogPage[*gen.Product], error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:read", authorization.AccessRead)
	if err != nil {
		return nil, err
	}
	offset, limit, err := pageBounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&gen.Product{}).Where("organization_id = ?", hqID)
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if filter.BrandID != nil {
		query = query.Where("brand_id = ?", *filter.BrandID)
	}
	if filter.Q != nil && strings.TrimSpace(*filter.Q) != "" {
		term := strings.TrimSpace(*filter.Q)
		query = query.Where("(products.name LIKE ? OR EXISTS (SELECT 1 FROM product_skus JOIN product_packages ON product_packages.sku_id = product_skus.id WHERE product_skus.product_id = products.id AND product_packages.enabled = true AND product_packages.package_set_version = product_skus.published_package_set_version AND product_packages.barcode = ?))", "%"+term+"%", term)
	}
	filter.Q = nil
	query = applyCatalogFilter(query, "products", filter)
	result := &CatalogPage[*gen.Product]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Preload("Brand").Preload("SpecificationChoices").Order("name, id").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) Skus(ctx context.Context, principal *auth.WorkspacePrincipal, productID string, filter CatalogFilter, page, perPage int) (*CatalogPage[*gen.ProductSku], error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:read", authorization.AccessRead)
	if err != nil {
		return nil, err
	}
	offset, limit, err := pageBounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&gen.ProductSku{}).Joins("JOIN products ON products.id = product_skus.product_id").Where("product_skus.product_id = ? AND products.organization_id = ?", productID, hqID)
	query = applyCatalogFilter(query, "product_skus", filter)
	result := &CatalogPage[*gen.ProductSku]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Preload("SpecificationValues").Order("product_skus.name, product_skus.id").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) Packages(ctx context.Context, principal *auth.WorkspacePrincipal, skuID string, filter CatalogFilter, page, perPage int) (*CatalogPage[*gen.ProductPackage], error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:read", authorization.AccessRead)
	if err != nil {
		return nil, err
	}
	offset, limit, err := pageBounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&gen.ProductPackage{}).Joins("JOIN product_skus ON product_skus.id = product_packages.sku_id").Joins("JOIN products ON products.id = product_skus.product_id").Where("product_packages.sku_id = ? AND products.organization_id = ?", skuID, hqID)
	query = applyCatalogFilter(query, "product_packages", filter)
	result := &CatalogPage[*gen.ProductPackage]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Order("product_packages.package_set_version, product_packages.name, product_packages.id").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}
