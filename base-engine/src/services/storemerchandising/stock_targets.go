package storemerchandising

import (
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func stockBatchAndListing(tx *gorm.DB, storeID, batchID string) (*gen.StoreInventoryBatch, *gen.StoreListing, error) {
	var batch gen.StoreInventoryBatch
	if err := tx.Joins("JOIN store_listings ON store_listings.id = store_inventory_batches.listing_id").
		Where("store_inventory_batches.id = ? AND store_listings.store_id = ?", batchID, storeID).First(&batch).Error; err != nil {
		return nil, nil, auth.NewError(auth.CodePermissionDenied)
	}
	if err := assertBatchBalanced(tx, batch.ID); err != nil {
		return nil, nil, err
	}
	var listing gen.StoreListing
	if err := tx.First(&listing, "id = ?", batch.ListingID).Error; err != nil {
		return nil, nil, err
	}
	return &batch, &listing, nil
}

func stockPackageForListing(tx *gorm.DB, packageID, skuID string) (*gen.ProductPackage, error) {
	var pack gen.ProductPackage
	if err := tx.Where("id = ? AND sku_id = ?", packageID, skuID).First(&pack).Error; err != nil {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	return &pack, nil
}
