package customer

import (
	"strings"

	"base-engine/gen"
	"gorm.io/gorm"
)

func disablePromotionsAboveNewPolicy(tx *gorm.DB, old *gen.CustomerBenefitPolicy, next PolicyInput) error {
	conditions := make([]string, 0, 3)
	if old.PromotionWithHqCoupon && !next.PromotionWithHqCoupon {
		conditions = append(conditions, "stack_with_hq_coupon = 1")
	}
	if old.PromotionWithStoreCoupon && !next.PromotionWithStoreCoupon {
		conditions = append(conditions, "stack_with_store_coupon = 1")
	}
	if old.MemberPriceWithPromotion && !next.MemberPriceWithPromotion {
		conditions = append(conditions, "stack_with_member_price = 1")
	}
	if len(conditions) == 0 {
		return nil
	}
	return tx.Model(&gen.StorePromotion{}).Where("enabled = ?", true).
		Where("("+strings.Join(conditions, " OR ")+")").Update("enabled", false).Error
}
