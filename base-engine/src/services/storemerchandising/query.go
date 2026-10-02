package storemerchandising

import (
	"context"
	"strings"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CatalogFilter struct {
	Q, CategoryID             *string
	ListingEnabled, Published *bool
}

type CatalogProductDetail struct {
	ProductName string
	CategoryID  string
}

type Page[T any] struct {
	Data          []T
	Total         int64
	Page, PerPage int
}

func bounds(page, perPage int) (int, int, error) {
	if page < 1 || perPage < 1 || perPage > 100 {
		return 0, 0, auth.NewError(auth.CodeValidationFailed)
	}
	return (page - 1) * perPage, perPage, nil
}

func (s *Service) Catalog(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, filter CatalogFilter, page, perPage int) (*Page[*gen.ProductSku], error) {
	if _, err := storeScope(s.db.WithContext(ctx), principal, storeID, scopedReadAction(principal, "franchiseProduct:read", "franchiseStock:read", "franchiseStocktake:read", "franchiseStocktake:record", "franchisePromotion:read"), authorization.AccessRead); err != nil {
		return nil, err
	}
	offset, limit, err := bounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&gen.ProductSku{}).Joins("JOIN products ON products.id = product_skus.product_id").Joins("JOIN organizations ON organizations.id = products.organization_id").
		Where("organizations.type = ? AND ((product_skus.enabled = ? AND products.enabled = ? AND organizations.status = ?) OR EXISTS (SELECT 1 FROM store_listings WHERE store_listings.store_id = ? AND store_listings.sku_id = product_skus.id))",
			gen.OrganizationTypeHeadquarters, true, true, gen.OrganizationStatusActive, storeID)
	query = applyStoreCatalogFilter(query, storeID, filter)
	result := &Page[*gen.ProductSku]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	listingEnabled := "SELECT enabled FROM store_listings WHERE store_id = ? AND sku_id = product_skus.id ORDER BY updated_at DESC, id DESC LIMIT 1"
	listingUpdatedAt := "SELECT updated_at FROM store_listings WHERE store_id = ? AND sku_id = product_skus.id ORDER BY updated_at DESC, id DESC LIMIT 1"
	order := clause.Expr{SQL: "COALESCE((" + listingEnabled + "), 0) DESC, COALESCE((" + listingUpdatedAt + "), product_skus.updated_at, product_skus.created_at) DESC, product_skus.id", Vars: []any{storeID, storeID}, WithoutParentheses: true}
	if err := query.Order(clause.OrderBy{Expression: order}).Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) CatalogProductDetails(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, productIDs []string) (map[string]CatalogProductDetail, error) {
	if _, err := storeScope(s.db.WithContext(ctx), principal, storeID, scopedReadAction(principal, "franchiseProduct:read", "franchiseStock:read", "franchiseStocktake:read", "franchiseStocktake:record", "franchisePromotion:read"), authorization.AccessRead); err != nil {
		return nil, err
	}
	result := make(map[string]CatalogProductDetail, len(productIDs))
	if len(productIDs) == 0 {
		return result, nil
	}
	var products []*gen.Product
	if err := s.db.WithContext(ctx).Model(&gen.Product{}).
		Joins("JOIN organizations ON organizations.id = products.organization_id").
		Where("organizations.type = ? AND products.id IN ?", gen.OrganizationTypeHeadquarters, productIDs).
		Find(&products).Error; err != nil {
		return nil, err
	}
	for _, product := range products {
		result[product.ID] = CatalogProductDetail{ProductName: product.Name, CategoryID: product.CategoryID}
	}
	return result, nil
}

func applyStoreCatalogFilter(query *gorm.DB, storeID string, filter CatalogFilter) *gorm.DB {
	if filter.Q != nil && strings.TrimSpace(*filter.Q) != "" {
		term := "%" + strings.TrimSpace(*filter.Q) + "%"
		query = query.Where("(product_skus.name LIKE ? OR products.name LIKE ? OR EXISTS (SELECT 1 FROM product_packages WHERE sku_id = product_skus.id AND enabled = true AND package_set_version = product_skus.published_package_set_version AND barcode = ?))", term, term, strings.TrimSpace(*filter.Q))
	}
	if filter.CategoryID != nil && *filter.CategoryID != "" {
		query = query.Where("products.category_id = ?", *filter.CategoryID)
	}
	if filter.ListingEnabled != nil {
		condition := "EXISTS (SELECT 1 FROM store_listings WHERE store_id = ? AND sku_id = product_skus.id AND enabled = ?)"
		if !*filter.ListingEnabled {
			condition = "NOT " + condition
		}
		query = query.Where(condition, storeID, true)
	}
	if filter.Published != nil {
		condition := "(product_skus.enabled = ? AND products.enabled = ? AND organizations.status = ?)"
		if !*filter.Published {
			condition = "NOT " + condition
		}
		query = query.Where(condition, true, true, gen.OrganizationStatusActive)
	}
	return query
}

func (s *Service) Listings(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, page, perPage int) (*Page[*gen.StoreListing], error) {
	if _, err := storeScope(s.db.WithContext(ctx), principal, storeID, scopedReadAction(principal, "franchiseProduct:read", "franchisePromotion:read"), authorization.AccessRead); err != nil {
		return nil, err
	}
	offset, limit, err := bounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&gen.StoreListing{}).Where("store_id = ?", storeID)
	result := &Page[*gen.StoreListing]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Order("selected_at DESC, id").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) Offers(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, listingID string, page, perPage int) (*Page[*gen.StorePackageOffer], error) {
	return s.OffersFiltered(ctx, principal, storeID, listingID, OfferFilter{}, page, perPage)
}

type OfferFilter struct {
	PackageID *string
	Enabled   *bool
}

func (s *Service) OffersFiltered(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, listingID string, filter OfferFilter, page, perPage int) (*Page[*gen.StorePackageOffer], error) {
	if _, err := storeScope(s.db.WithContext(ctx), principal, storeID, scopedReadAction(principal, "franchiseProduct:read", "franchisePromotion:read"), authorization.AccessRead); err != nil {
		return nil, err
	}
	offset, limit, err := bounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&gen.StorePackageOffer{}).Joins("JOIN store_listings ON store_listings.id = store_package_offers.listing_id").Where("store_listings.store_id = ? AND store_package_offers.listing_id = ?", storeID, listingID)
	if listingID == "" {
		query = s.db.WithContext(ctx).Model(&gen.StorePackageOffer{}).Joins("JOIN store_listings ON store_listings.id = store_package_offers.listing_id").Where("store_listings.store_id = ?", storeID)
	}
	if filter.PackageID != nil {
		query = query.Where("store_package_offers.package_id = ?", *filter.PackageID)
	}
	if filter.Enabled != nil {
		query = query.Where("store_package_offers.enabled = ?", *filter.Enabled)
	}
	result := &Page[*gen.StorePackageOffer]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Order("store_package_offers.id").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) Movements(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, batchID string, page, perPage int) (*Page[*gen.StoreStockMovement], error) {
	return s.MovementsFiltered(ctx, principal, storeID, batchID, MovementFilter{}, page, perPage)
}

type MovementFilter struct {
	Kind *gen.StockMovementKind
	From *time.Time
	To   *time.Time
}

func (s *Service) MovementsFiltered(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, batchID string, filter MovementFilter, page, perPage int) (*Page[*gen.StoreStockMovement], error) {
	if _, err := storeScope(s.db.WithContext(ctx), principal, storeID, "franchiseStock:read", authorization.AccessRead); err != nil {
		return nil, err
	}
	offset, limit, err := bounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&gen.StoreStockMovement{}).Where("store_id = ? AND batch_id = ?", storeID, batchID)
	if filter.Kind != nil {
		query = query.Where("kind = ?", *filter.Kind)
	}
	if filter.From != nil {
		query = query.Where("occurred_at >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("occurred_at <= ?", *filter.To)
	}
	result := &Page[*gen.StoreStockMovement]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Order("occurred_at DESC, id").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}
