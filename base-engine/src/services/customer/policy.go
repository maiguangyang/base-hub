package customer

import (
	"context"
	"math"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PolicyInput struct {
	DiscountEnabled, PurchaseEarnEnabled, RedemptionEnabled         bool
	PromotionWithHqCoupon, PromotionWithStoreCoupon                 bool
	MemberPriceWithPromotion, PointsWithPromotion, PointsWithCoupon bool
	DiscountBasisPoints, EarnAmountFen, EarnPoints                  int64
	RedeemPoints, RedeemAmountFen                                   int64
	MaxRedemptionBasisPoints, MaxRedemptionPoints                   int64
	ManualGrantMaxSingle, ManualGrantMaxDaily                       int64
}

type PolicyView struct {
	PolicyInput
	Version int64
}

func defaultPolicy() PolicyView {
	return PolicyView{PolicyInput: PolicyInput{RedeemPoints: 100, RedeemAmountFen: 100}}
}

func validPolicy(input PolicyInput) bool {
	return policyNumbersInRange(input) && policyRulesValid(input)
}

func policyNumbersInRange(input PolicyInput) bool {
	values := []int64{
		input.DiscountBasisPoints, input.EarnAmountFen, input.EarnPoints,
		input.RedeemPoints, input.RedeemAmountFen, input.MaxRedemptionBasisPoints,
		input.MaxRedemptionPoints, input.ManualGrantMaxSingle, input.ManualGrantMaxDaily,
	}
	for _, value := range values {
		if value < 0 || value > math.MaxInt32 {
			return false
		}
	}
	return true
}

func policyRulesValid(input PolicyInput) bool {
	if !manualGrantRuleValid(input) {
		return false
	}
	if input.DiscountEnabled && (input.DiscountBasisPoints < 1 || input.DiscountBasisPoints > 10000) {
		return false
	}
	if input.PurchaseEarnEnabled && (input.EarnAmountFen == 0 || input.EarnPoints == 0) {
		return false
	}
	return redemptionRuleValid(input)
}

func manualGrantRuleValid(input PolicyInput) bool {
	if input.ManualGrantMaxSingle == 0 && input.ManualGrantMaxDaily == 0 {
		return true
	}
	return input.ManualGrantMaxSingle > 0 && input.ManualGrantMaxDaily >= input.ManualGrantMaxSingle
}

func redemptionRuleValid(input PolicyInput) bool {
	if !input.RedemptionEnabled {
		return true
	}
	return input.RedeemPoints > 0 && input.RedeemAmountFen > 0 &&
		input.MaxRedemptionBasisPoints >= 1 && input.MaxRedemptionBasisPoints <= 10000 &&
		input.MaxRedemptionPoints > 0
}

func policyView(record *gen.CustomerBenefitPolicy) PolicyView {
	return PolicyView{Version: record.Version, PolicyInput: PolicyInput{
		DiscountEnabled: record.DiscountEnabled, DiscountBasisPoints: record.DiscountBasisPoints,
		PurchaseEarnEnabled: record.PurchaseEarnEnabled, EarnAmountFen: record.EarnAmountFen,
		EarnPoints: record.EarnPoints, RedemptionEnabled: record.RedemptionEnabled,
		RedeemPoints: record.RedeemPoints, RedeemAmountFen: record.RedeemAmountFen,
		MaxRedemptionBasisPoints: record.MaxRedemptionBasisPoints,
		MaxRedemptionPoints:      record.MaxRedemptionPoints,
		ManualGrantMaxSingle:     record.ManualGrantMaxSingle, ManualGrantMaxDaily: record.ManualGrantMaxDaily,
		PromotionWithHqCoupon: record.PromotionWithHqCoupon, PromotionWithStoreCoupon: record.PromotionWithStoreCoupon,
		MemberPriceWithPromotion: record.MemberPriceWithPromotion, PointsWithPromotion: record.PointsWithPromotion,
		PointsWithCoupon: record.PointsWithCoupon,
	}}
}

func (s *Service) GetPolicy(ctx context.Context, principal *auth.WorkspacePrincipal) (PolicyView, error) {
	hqID, err := headquartersID(principal, "hqCustomerPolicy:read", authorization.AccessRead)
	if err != nil {
		return PolicyView{}, err
	}
	var record gen.CustomerBenefitPolicy
	err = s.db.WithContext(ctx).Where("organization_id = ?", hqID).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return defaultPolicy(), nil
	}
	if err != nil {
		return PolicyView{}, err
	}
	return policyView(&record), nil
}

func (s *Service) SavePolicy(ctx context.Context, principal *auth.WorkspacePrincipal, expectedVersion int64, input PolicyInput) (PolicyView, error) {
	hqID, err := headquartersID(principal, "hqCustomerPolicy:manage", authorization.AccessUpdate)
	if err != nil {
		return PolicyView{}, err
	}
	if expectedVersion < 0 || expectedVersion > math.MaxInt32-1 || !validPolicy(input) {
		return PolicyView{}, auth.NewError(auth.CodeValidationFailed)
	}
	var result PolicyView
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		record, err := savePolicyRecord(tx, hqID, expectedVersion, input)
		if err != nil {
			return err
		}
		result = policyView(&record)
		return s.auditPolicy(tx, principal, &record, expectedVersion)
	})
	return result, err
}

func savePolicyRecord(tx *gorm.DB, hqID string, expectedVersion int64, input PolicyInput) (gen.CustomerBenefitPolicy, error) {
	var record gen.CustomerBenefitPolicy
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("organization_id = ?", hqID).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return createPolicyRecord(tx, hqID, expectedVersion, input)
	}
	if err != nil {
		return record, err
	}
	if record.Version != expectedVersion {
		return record, auth.NewError(auth.CodeConflict)
	}
	previous := record
	updated := tx.Model(&record).Where("version = ?", expectedVersion).Updates(policyUpdates(input, expectedVersion+1))
	if updated.Error != nil {
		return record, updated.Error
	}
	if updated.RowsAffected != 1 {
		return record, auth.NewError(auth.CodeConflict)
	}
	if err := disablePromotionsAboveNewPolicy(tx, &previous, input); err != nil {
		return record, err
	}
	return policyRecord(record.ID, hqID, expectedVersion+1, input), nil
}

func createPolicyRecord(tx *gorm.DB, hqID string, expectedVersion int64, input PolicyInput) (gen.CustomerBenefitPolicy, error) {
	if expectedVersion != 0 {
		return gen.CustomerBenefitPolicy{}, auth.NewError(auth.CodeConflict)
	}
	record := policyRecord(uuid.Must(uuid.NewV4()).String(), hqID, 1, input)
	record.CreatedAt = time.Now().UnixMilli()
	insert := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&record)
	if insert.Error != nil {
		return record, insert.Error
	}
	if insert.RowsAffected != 1 {
		return record, auth.NewError(auth.CodeConflict)
	}
	return record, nil
}

func (s *Service) auditPolicy(tx *gorm.DB, principal *auth.WorkspacePrincipal, record *gen.CustomerBenefitPolicy, before int64) error {
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: principal.OrganizationID, Action: "hqCustomerPolicy:manage",
		ResourceType: "customerBenefitPolicy", ResourceID: record.ID, ResultCode: "SUCCESS",
		Metadata: audit.MetadataForPrincipal(principal, audit.Metadata{
			BeforeVersion: uint64(before), AfterVersion: uint64(record.Version),
		}),
	})
}

func policyRecord(id, hqID string, version int64, input PolicyInput) gen.CustomerBenefitPolicy {
	return gen.CustomerBenefitPolicy{
		ID: id, OrganizationID: hqID, Version: version,
		DiscountEnabled: input.DiscountEnabled, DiscountBasisPoints: input.DiscountBasisPoints,
		PurchaseEarnEnabled: input.PurchaseEarnEnabled, EarnAmountFen: input.EarnAmountFen,
		EarnPoints: input.EarnPoints, RedemptionEnabled: input.RedemptionEnabled,
		RedeemPoints: input.RedeemPoints, RedeemAmountFen: input.RedeemAmountFen,
		MaxRedemptionBasisPoints: input.MaxRedemptionBasisPoints,
		MaxRedemptionPoints:      input.MaxRedemptionPoints,
		ManualGrantMaxSingle:     input.ManualGrantMaxSingle, ManualGrantMaxDaily: input.ManualGrantMaxDaily,
		PromotionWithHqCoupon: input.PromotionWithHqCoupon, PromotionWithStoreCoupon: input.PromotionWithStoreCoupon,
		MemberPriceWithPromotion: input.MemberPriceWithPromotion, PointsWithPromotion: input.PointsWithPromotion,
		PointsWithCoupon: input.PointsWithCoupon,
	}
}

func policyUpdates(input PolicyInput, version int64) map[string]any {
	return map[string]any{
		"version": version, "discount_enabled": input.DiscountEnabled,
		"discount_basis_points": input.DiscountBasisPoints,
		"purchase_earn_enabled": input.PurchaseEarnEnabled,
		"earn_amount_fen":       input.EarnAmountFen, "earn_points": input.EarnPoints,
		"redemption_enabled": input.RedemptionEnabled, "redeem_points": input.RedeemPoints,
		"redeem_amount_fen":           input.RedeemAmountFen,
		"max_redemption_basis_points": input.MaxRedemptionBasisPoints,
		"max_redemption_points":       input.MaxRedemptionPoints,
		"manual_grant_max_single":     input.ManualGrantMaxSingle,
		"manual_grant_max_daily":      input.ManualGrantMaxDaily,
		"promotion_with_hq_coupon":    input.PromotionWithHqCoupon,
		"promotion_with_store_coupon": input.PromotionWithStoreCoupon,
		"member_price_with_promotion": input.MemberPriceWithPromotion,
		"points_with_promotion":       input.PointsWithPromotion,
		"points_with_coupon":          input.PointsWithCoupon,
	}
}
