package storemerchandising

import (
	"math"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PromotionTargetInput struct {
	OfferID          string
	RequiredQuantity int64
}
type PromotionInput struct {
	StoreID, TimeZone                                                                string
	PromotionID                                                                      *string
	Kind                                                                             gen.StorePromotionKind
	StartsAt, EndsAt                                                                 time.Time
	ThresholdFen, ThresholdQuantity, DiscountFen, DiscountBasisPoints, FixedPriceFen *int64
	StackWithHqCoupon, StackWithStoreCoupon, StackWithMemberPrice                    bool
	Targets                                                                          []PromotionTargetInput
}

func validatePromotion(input PromotionInput) error {
	if !validPromotionTiming(input) || !validPromotionTargets(input.Targets) || !validPromotionNumbers(input) || !validPromotionShape(input) {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}

func validPromotionTiming(input PromotionInput) bool {
	if input.StoreID == "" || input.TimeZone == "" || len(input.TimeZone) > 64 {
		return false
	}
	if !input.EndsAt.After(input.StartsAt) || !input.EndsAt.After(time.Now()) {
		return false
	}
	_, err := time.LoadLocation(input.TimeZone)
	return err == nil
}

func validPromotionTargets(targets []PromotionTargetInput) bool {
	if len(targets) == 0 || len(targets) > 100 {
		return false
	}
	seen := make(map[string]struct{}, len(targets))
	for _, target := range targets {
		if target.OfferID == "" || target.RequiredQuantity <= 0 || target.RequiredQuantity > math.MaxInt32 {
			return false
		}
		if _, exists := seen[target.OfferID]; exists {
			return false
		}
		seen[target.OfferID] = struct{}{}
	}
	return true
}

func validPromotionNumbers(input PromotionInput) bool {
	for _, value := range []*int64{input.ThresholdFen, input.ThresholdQuantity, input.DiscountFen, input.DiscountBasisPoints, input.FixedPriceFen} {
		if value != nil && (*value < 0 || *value > math.MaxInt32) {
			return false
		}
	}
	return input.DiscountBasisPoints == nil || (*input.DiscountBasisPoints > 0 && *input.DiscountBasisPoints < 10000)
}

func validPromotionShape(input PromotionInput) bool {
	single := len(input.Targets) == 1 && input.Targets[0].RequiredQuantity == 1
	switch input.Kind {
	case gen.StorePromotionKindItemPrice:
		return validItemPrice(input, single)
	case gen.StorePromotionKindAmountThreshold:
		return validAmountThreshold(input)
	case gen.StorePromotionKindQuantityThreshold:
		return validQuantityThreshold(input)
	case gen.StorePromotionKindBundlePrice:
		return validBundlePrice(input)
	case gen.StorePromotionKindMemberPrice:
		return validMemberPrice(input, single)
	default:
		return false
	}
}

func validItemPrice(input PromotionInput, single bool) bool {
	return single && exactlyOne(input.FixedPriceFen, input.DiscountBasisPoints) &&
		input.ThresholdFen == nil && input.ThresholdQuantity == nil && input.DiscountFen == nil
}

func validAmountThreshold(input PromotionInput) bool {
	return positive(input.ThresholdFen) && positive(input.DiscountFen) &&
		*input.DiscountFen <= *input.ThresholdFen && input.ThresholdQuantity == nil &&
		input.DiscountBasisPoints == nil && input.FixedPriceFen == nil
}

func validQuantityThreshold(input PromotionInput) bool {
	return positive(input.ThresholdQuantity) && positive(input.DiscountFen) &&
		input.ThresholdFen == nil && input.DiscountBasisPoints == nil && input.FixedPriceFen == nil
}

func validBundlePrice(input PromotionInput) bool {
	return len(input.Targets) >= 2 && positive(input.FixedPriceFen) && input.ThresholdFen == nil &&
		input.ThresholdQuantity == nil && input.DiscountFen == nil && input.DiscountBasisPoints == nil
}

func validMemberPrice(input PromotionInput, single bool) bool {
	return single && positive(input.FixedPriceFen) && input.ThresholdFen == nil &&
		input.ThresholdQuantity == nil && input.DiscountFen == nil && input.DiscountBasisPoints == nil
}

func positive(value *int64) bool { return value != nil && *value > 0 }
func exactlyOne(a, b *int64) bool {
	return (a == nil) != (b == nil) && ((a == nil || *a > 0) && (b == nil || *b > 0))
}

func validatePromotionPolicy(tx *gorm.DB, input PromotionInput) error {
	var policy gen.CustomerBenefitPolicy
	err := tx.Clauses(clause.Locking{Strength: "SHARE"}).Joins("JOIN organizations ON organizations.id = customer_benefit_policies.organization_id").Where("organizations.type = ?", gen.OrganizationTypeHeadquarters).First(&policy).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if input.StackWithHqCoupon && !policy.PromotionWithHqCoupon {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if input.StackWithStoreCoupon && !policy.PromotionWithStoreCoupon {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if input.StackWithMemberPrice && !policy.MemberPriceWithPromotion {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}
