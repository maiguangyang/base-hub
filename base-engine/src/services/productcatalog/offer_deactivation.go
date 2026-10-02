package productcatalog

import (
	"base-engine/gen"
	"gorm.io/gorm"
)

func disableOffersForSkus(tx *gorm.DB, skuIDs *gorm.DB) error {
	listingIDs := tx.Model(&gen.StoreListing{}).Select("id").Where("sku_id IN (?)", skuIDs)
	offerIDs := tx.Model(&gen.StorePackageOffer{}).Select("id").Where("listing_id IN (?)", listingIDs)
	promotionIDs := tx.Model(&gen.StorePromotionTarget{}).Select("promotion_id").Where("offer_id IN (?)", offerIDs)
	if err := tx.Model(&gen.StorePromotion{}).Where("id IN (?)", promotionIDs).Update("enabled", false).Error; err != nil {
		return err
	}
	return tx.Model(&gen.StorePackageOffer{}).Where("listing_id IN (?)", listingIDs).Update("enabled", false).Error
}
