package customer

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type PolicyStackingOverrides struct {
	PromotionWithHqCoupon, PromotionWithStoreCoupon                 *bool
	MemberPriceWithPromotion, PointsWithPromotion, PointsWithCoupon *bool
}

// SavePolicyWithOptionalStacking preserves existing stacking rules for older clients.
func (s *Service) SavePolicyWithOptionalStacking(ctx context.Context, principal *auth.WorkspacePrincipal, expectedVersion int64, input PolicyInput, options PolicyStackingOverrides) (PolicyView, error) {
	hqID, err := headquartersID(principal, "hqCustomerPolicy:manage", authorization.AccessUpdate)
	if err != nil {
		return PolicyView{}, err
	}
	var record gen.CustomerBenefitPolicy
	err = s.db.WithContext(ctx).Where("organization_id = ?", hqID).First(&record).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return PolicyView{}, err
	}
	if err == nil {
		input.PromotionWithHqCoupon = record.PromotionWithHqCoupon
		input.PromotionWithStoreCoupon = record.PromotionWithStoreCoupon
		input.MemberPriceWithPromotion = record.MemberPriceWithPromotion
		input.PointsWithPromotion = record.PointsWithPromotion
		input.PointsWithCoupon = record.PointsWithCoupon
	}
	if options.PromotionWithHqCoupon != nil {
		input.PromotionWithHqCoupon = *options.PromotionWithHqCoupon
	}
	if options.PromotionWithStoreCoupon != nil {
		input.PromotionWithStoreCoupon = *options.PromotionWithStoreCoupon
	}
	if options.MemberPriceWithPromotion != nil {
		input.MemberPriceWithPromotion = *options.MemberPriceWithPromotion
	}
	if options.PointsWithPromotion != nil {
		input.PointsWithPromotion = *options.PointsWithPromotion
	}
	if options.PointsWithCoupon != nil {
		input.PointsWithCoupon = *options.PointsWithCoupon
	}
	return s.SavePolicy(ctx, principal, expectedVersion, input)
}
