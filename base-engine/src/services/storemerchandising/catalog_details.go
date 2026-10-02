package storemerchandising

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
)

// CatalogPackages includes retired definitions so historical stock remains identifiable.
func (s *Service) CatalogPackages(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, skuID string) ([]*gen.ProductPackage, error) {
	if _, err := storeScope(s.db.WithContext(ctx), principal, storeID, scopedReadAction(principal, "franchiseProduct:read", "franchiseStock:read", "franchiseStocktake:read", "franchiseStocktake:record", "franchisePromotion:read"), authorization.AccessRead); err != nil {
		return nil, err
	}
	if err := publishedSku(s.db.WithContext(ctx), skuID); err != nil {
		var count int64
		if queryErr := s.db.WithContext(ctx).Model(&gen.StoreListing{}).Where("store_id = ? AND sku_id = ?", storeID, skuID).Count(&count).Error; queryErr != nil {
			return nil, queryErr
		}
		if count == 0 {
			return nil, err
		}
	}
	var packages []*gen.ProductPackage
	err := s.db.WithContext(ctx).Where("sku_id = ?", skuID).Order("package_set_version DESC, name, id").Find(&packages).Error
	return packages, err
}

func (s *Service) CatalogPublished(ctx context.Context, skuID string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&gen.ProductSku{}).Joins("JOIN products ON products.id = product_skus.product_id").
		Joins("JOIN organizations ON organizations.id = products.organization_id").
		Where("product_skus.id = ? AND product_skus.enabled = ? AND products.enabled = ? AND organizations.type = ? AND organizations.status = ?",
			skuID, true, true, gen.OrganizationTypeHeadquarters, gen.OrganizationStatusActive).Count(&count).Error
	return count > 0, err
}

// CatalogListings returns the authorized store's selections for one catalog page.
func (s *Service) CatalogListings(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, skuIDs []string) (map[string]*gen.StoreListing, error) {
	if _, err := storeScope(s.db.WithContext(ctx), principal, storeID, scopedReadAction(principal, "franchiseProduct:read", "franchiseStock:read", "franchiseStocktake:read", "franchiseStocktake:record", "franchisePromotion:read"), authorization.AccessRead); err != nil {
		return nil, err
	}
	result := make(map[string]*gen.StoreListing, len(skuIDs))
	if len(skuIDs) == 0 {
		return result, nil
	}
	var listings []*gen.StoreListing
	if err := s.db.WithContext(ctx).Where("store_id = ? AND sku_id IN ?", storeID, skuIDs).Find(&listings).Error; err != nil {
		return nil, err
	}
	for _, listing := range listings {
		result[listing.SkuID] = listing
	}
	return result, nil
}

// OfferPriceRevisions keeps historical prices behind the same store and listing scope as offers.
func (s *Service) OfferPriceRevisions(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, listingID, offerID string) ([]*gen.StorePriceRevision, error) {
	if _, err := storeScope(s.db.WithContext(ctx), principal, storeID, "franchiseProduct:read", authorization.AccessRead); err != nil {
		return nil, err
	}
	var revisions []*gen.StorePriceRevision
	err := s.db.WithContext(ctx).Model(&gen.StorePriceRevision{}).
		Joins("JOIN store_package_offers ON store_package_offers.id = store_price_revisions.offer_id").
		Joins("JOIN store_listings ON store_listings.id = store_package_offers.listing_id").
		Where("store_listings.store_id = ? AND store_listings.id = ? AND store_package_offers.id = ?", storeID, listingID, offerID).
		Order("store_price_revisions.effective_at DESC, store_price_revisions.id DESC").Limit(20).Find(&revisions).Error
	return revisions, err
}
