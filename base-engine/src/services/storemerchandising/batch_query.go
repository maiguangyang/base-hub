package storemerchandising

import (
	"context"
	"strings"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type BatchView struct {
	Batch                *gen.StoreInventoryBatch
	SkuID                string
	ProductName, SkuName *string
	Packages             []*gen.ProductPackage
	Balances             []*gen.StoreStockBalance
	Sellable             bool
}

type BatchFilter struct {
	Q        *string
	Sellable *bool
}

func (s *Service) Batches(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, listingID string, page, perPage int) (*Page[BatchView], error) {
	return s.BatchesFiltered(ctx, principal, storeID, listingID, BatchFilter{}, page, perPage)
}

func (s *Service) BatchesFiltered(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, listingID string, filter BatchFilter, page, perPage int) (*Page[BatchView], error) {
	db := s.db.WithContext(ctx)
	if _, err := storeScope(db, principal, storeID, "franchiseStock:read", authorization.AccessRead); err != nil {
		return nil, err
	}
	offset, limit, err := bounds(page, perPage)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	query := applyBatchFilter(batchBaseQuery(db, storeID, listingID, filter), filter, now)
	result := &Page[BatchView]{Data: []BatchView{}, Page: page, PerPage: perPage}
	if err := query.Session(&gorm.Session{}).Select("*").Count(&result.Total).Error; err != nil {
		return nil, err
	}
	var batches []*gen.StoreInventoryBatch
	if err := query.Order("store_inventory_batches.expires_at, store_inventory_batches.id").Offset(offset).Limit(limit).Find(&batches).Error; err != nil {
		return nil, err
	}
	if len(batches) == 0 {
		return result, nil
	}
	relations, err := loadBatchRelations(db, batches)
	if err != nil {
		return nil, err
	}
	result.Data, err = assembleBatchViews(batches, relations, now)
	return result, err
}

func batchBaseQuery(db *gorm.DB, storeID, listingID string, filter BatchFilter) *gorm.DB {
	query := db.Model(&gen.StoreInventoryBatch{}).Select("store_inventory_batches.*").
		Joins("JOIN store_listings ON store_listings.id = store_inventory_batches.listing_id").Where("store_listings.store_id = ?", storeID)
	if listingID != "" {
		query = query.Where("store_inventory_batches.listing_id = ?", listingID)
	}
	if filter.Q != nil && strings.TrimSpace(*filter.Q) != "" {
		query = query.Joins("LEFT JOIN product_skus ON product_skus.id = store_listings.sku_id").Joins("LEFT JOIN products ON products.id = product_skus.product_id")
	}
	return query
}

func applyBatchFilter(query *gorm.DB, filter BatchFilter, now time.Time) *gorm.DB {
	if filter.Q != nil && strings.TrimSpace(*filter.Q) != "" {
		term := "%" + strings.TrimSpace(*filter.Q) + "%"
		query = query.Where("(store_inventory_batches.batch_number LIKE ? OR product_skus.name LIKE ? OR products.name LIKE ?)", term, term, term)
	}
	if filter.Sellable == nil {
		return query
	}
	if *filter.Sellable {
		return query.Where("(store_inventory_batches.expires_at IS NULL OR store_inventory_batches.expires_at > ?)", now)
	}
	return query.Where("store_inventory_batches.expires_at <= ?", now)
}

type batchRelations struct {
	listings map[string]*gen.StoreListing
	skus     map[string]*gen.ProductSku
	products map[string]*gen.Product
	packages map[string][]*gen.ProductPackage
	balances map[string][]*gen.StoreStockBalance
}

func loadBatchRelations(db *gorm.DB, batches []*gen.StoreInventoryBatch) (*batchRelations, error) {
	listingSet := map[string]bool{}
	batchIDs, listingIDs := []string{}, []string{}
	for _, batch := range batches {
		batchIDs = append(batchIDs, batch.ID)
		if !listingSet[batch.ListingID] {
			listingIDs = append(listingIDs, batch.ListingID)
			listingSet[batch.ListingID] = true
		}
	}
	var listings []*gen.StoreListing
	if err := db.Where("id IN ?", listingIDs).Find(&listings).Error; err != nil {
		return nil, err
	}
	relations := &batchRelations{listings: map[string]*gen.StoreListing{}, skus: map[string]*gen.ProductSku{}, products: map[string]*gen.Product{}, packages: map[string][]*gen.ProductPackage{}, balances: map[string][]*gen.StoreStockBalance{}}
	skuSet := map[string]bool{}
	skuIDs := []string{}
	for _, listing := range listings {
		relations.listings[listing.ID] = listing
		if !skuSet[listing.SkuID] {
			skuIDs = append(skuIDs, listing.SkuID)
			skuSet[listing.SkuID] = true
		}
	}
	if err := loadBatchProducts(db, relations, skuIDs); err != nil {
		return nil, err
	}
	var packages []*gen.ProductPackage
	if err := db.Where("sku_id IN ?", skuIDs).Order("package_set_version DESC, name, id").Find(&packages).Error; err != nil {
		return nil, err
	}
	for _, pack := range packages {
		relations.packages[pack.SkuID] = append(relations.packages[pack.SkuID], pack)
	}
	if err := loadBatchBalances(db, relations, batchIDs); err != nil {
		return nil, err
	}
	return relations, nil
}

func loadBatchBalances(db *gorm.DB, relations *batchRelations, batchIDs []string) error {
	var balances []*gen.StoreStockBalance
	if err := db.Where("batch_id IN ?", batchIDs).Order("batch_id, package_id").Find(&balances).Error; err != nil {
		return err
	}
	for _, balance := range balances {
		relations.balances[balance.BatchID] = append(relations.balances[balance.BatchID], balance)
	}
	return nil
}

func loadBatchProducts(db *gorm.DB, relations *batchRelations, skuIDs []string) error {
	var skus []*gen.ProductSku
	if err := db.Where("id IN ?", skuIDs).Find(&skus).Error; err != nil {
		return err
	}
	productIDs := []string{}
	seen := map[string]bool{}
	for _, sku := range skus {
		relations.skus[sku.ID] = sku
		if !seen[sku.ProductID] {
			productIDs = append(productIDs, sku.ProductID)
			seen[sku.ProductID] = true
		}
	}
	if len(productIDs) == 0 {
		return nil
	}
	var products []*gen.Product
	if err := db.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		return err
	}
	for _, product := range products {
		relations.products[product.ID] = product
	}
	return nil
}

func assembleBatchViews(batches []*gen.StoreInventoryBatch, relations *batchRelations, now time.Time) ([]BatchView, error) {
	result := make([]BatchView, 0, len(batches))
	for _, batch := range batches {
		listing := relations.listings[batch.ListingID]
		if listing == nil {
			return nil, auth.NewError(auth.CodeValidationFailed)
		}
		view := BatchView{Batch: batch, SkuID: listing.SkuID, Packages: relations.packages[listing.SkuID], Balances: relations.balances[batch.ID], Sellable: batch.ExpiresAt == nil || batch.ExpiresAt.After(now)}
		if sku := relations.skus[listing.SkuID]; sku != nil {
			name := sku.Name
			view.SkuName = &name
			if product := relations.products[sku.ProductID]; product != nil {
				name := product.Name
				view.ProductName = &name
			}
		}
		result = append(result, view)
	}
	return result, nil
}
