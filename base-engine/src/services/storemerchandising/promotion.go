package storemerchandising

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) SavePromotion(ctx context.Context, principal *auth.WorkspacePrincipal, input PromotionInput) (*gen.StorePromotion, error) {
	if err := validatePromotion(input); err != nil {
		return nil, err
	}
	var promotion gen.StorePromotion
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		store, err := storeScope(tx, principal, input.StoreID, "franchisePromotion:manage", authorization.AccessCreate)
		if err != nil {
			return err
		}
		if err := validatePromotionPolicy(tx, input); err != nil {
			return err
		}
		if err := validatePromotionTargets(tx, store.ID, input.Targets); err != nil {
			return err
		}
		ruleKey, version, err := nextPromotionVersion(tx, store.ID, input.PromotionID)
		if err != nil {
			return err
		}
		promotion = promotionRecord(input, ruleKey, version)
		if err := tx.Create(&promotion).Error; err != nil {
			return err
		}
		for _, target := range input.Targets {
			item := gen.StorePromotionTarget{ID: uuid.Must(uuid.NewV4()).String(), PromotionID: promotion.ID, OfferID: target.OfferID, RequiredQuantity: target.RequiredQuantity}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: "franchisePromotion:manage",
			ResourceType: "storePromotion", ResourceID: promotion.ID, ResultCode: "SUCCESS"})
	})
	return &promotion, err
}

func validatePromotionTargets(tx *gorm.DB, storeID string, targets []PromotionTargetInput) error {
	for _, target := range targets {
		var offer gen.StorePackageOffer
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Joins("JOIN store_listings ON store_listings.id = store_package_offers.listing_id").
			Joins("JOIN product_packages ON product_packages.id = store_package_offers.package_id").
			Where("store_package_offers.id = ? AND store_listings.store_id = ? AND store_listings.enabled = ? AND store_package_offers.enabled = ? AND product_packages.enabled = ?", target.OfferID, storeID, true, true, true).First(&offer).Error
		if err != nil {
			return auth.NewError(auth.CodeValidationFailed)
		}
	}
	return nil
}

func nextPromotionVersion(tx *gorm.DB, storeID string, previousID *string) (string, int64, error) {
	if previousID == nil {
		return uuid.Must(uuid.NewV4()).String(), 1, nil
	}
	var previous gen.StorePromotion
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND store_id = ?", *previousID, storeID).First(&previous).Error; err != nil {
		return "", 0, auth.NewError(auth.CodePermissionDenied)
	}
	var latest int64
	if err := tx.Model(&gen.StorePromotion{}).Where("store_id = ? AND rule_key = ?", storeID, previous.RuleKey).Select("MAX(version)").Scan(&latest).Error; err != nil {
		return "", 0, err
	}
	if previous.Version != latest {
		return "", 0, auth.NewError(auth.CodeConflict)
	}
	return previous.RuleKey, latest + 1, nil
}

func promotionRecord(input PromotionInput, ruleKey string, version int64) gen.StorePromotion {
	return gen.StorePromotion{ID: uuid.Must(uuid.NewV4()).String(), StoreID: input.StoreID, RuleKey: ruleKey, Version: version,
		TimeZone: input.TimeZone, Kind: input.Kind, Enabled: false, StartsAt: input.StartsAt, EndsAt: input.EndsAt,
		ThresholdFen: input.ThresholdFen, ThresholdQuantity: input.ThresholdQuantity, DiscountFen: input.DiscountFen,
		DiscountBasisPoints: input.DiscountBasisPoints, FixedPriceFen: input.FixedPriceFen,
		StackWithHqCoupon: input.StackWithHqCoupon, StackWithStoreCoupon: input.StackWithStoreCoupon,
		StackWithMemberPrice: input.StackWithMemberPrice}
}

func (s *Service) SetPromotionEnabled(ctx context.Context, principal *auth.WorkspacePrincipal, id string, enabled bool) (*gen.StorePromotion, error) {
	var item gen.StorePromotion
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, "id = ?", id).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		store, err := storeScope(tx, principal, item.StoreID, "franchisePromotion:manage", authorization.AccessUpdate)
		if err != nil {
			return err
		}
		if enabled {
			if err := validatePromotionActivationPolicy(tx, &item); err != nil {
				return err
			}
			if err := validateExistingPromotionTargets(tx, item.StoreID, item.ID); err != nil {
				return err
			}
			var latest int64
			if err := tx.Model(&gen.StorePromotion{}).Where("store_id = ? AND rule_key = ?", item.StoreID, item.RuleKey).Select("MAX(version)").Scan(&latest).Error; err != nil {
				return err
			}
			if item.Version != latest {
				return auth.NewError(auth.CodeConflict)
			}
			if err := tx.Model(&gen.StorePromotion{}).Where("store_id = ? AND rule_key = ? AND id <> ?", item.StoreID, item.RuleKey, item.ID).Update("enabled", false).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&item).Update("enabled", enabled).Error; err != nil {
			return err
		}
		item.Enabled = enabled
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: "franchisePromotion:manage",
			ResourceType: "storePromotion", ResourceID: item.ID, ResultCode: "SUCCESS"})
	})
	return &item, err
}

func validatePromotionActivationPolicy(tx *gorm.DB, item *gen.StorePromotion) error {
	if !item.EndsAt.After(time.Now()) {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return validatePromotionPolicy(tx, PromotionInput{StackWithHqCoupon: item.StackWithHqCoupon,
		StackWithStoreCoupon: item.StackWithStoreCoupon, StackWithMemberPrice: item.StackWithMemberPrice})
}

func validateExistingPromotionTargets(tx *gorm.DB, storeID, promotionID string) error {
	var targets []gen.StorePromotionTarget
	if err := tx.Where("promotion_id = ?", promotionID).Find(&targets).Error; err != nil {
		return err
	}
	if len(targets) == 0 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	items := make([]PromotionTargetInput, 0, len(targets))
	for _, target := range targets {
		items = append(items, PromotionTargetInput{OfferID: target.OfferID, RequiredQuantity: target.RequiredQuantity})
	}
	return validatePromotionTargets(tx, storeID, items)
}
