package customer

import (
	"context"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) CreateStoreCouponTemplate(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, input FixedAmountTemplateInput) (*gen.CustomerCouponTemplate, error) {
	if !validCouponTemplate(input) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	nowMillis := s.now().UnixMilli()
	effectiveAt, err := normalizedCouponTemplateTimes(input, nowMillis)
	if err != nil {
		return nil, err
	}
	template, hqID, err := s.createStoreCouponTemplateTransaction(ctx, principal, storeID, input, nowMillis, effectiveAt)
	if auth.ErrorCode(err) != auth.CodeConflict {
		return template, err
	}
	replay, findErr := findStoreCouponTemplateByRequest(s.db.WithContext(ctx), hqID, storeID, input.RequestKey)
	if findErr != nil {
		return nil, findErr
	}
	if replay != nil && sameTemplateIntent(replay, input) {
		return replay, nil
	}
	return template, err
}

func (s *Service) createStoreCouponTemplateTransaction(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, input FixedAmountTemplateInput, nowMillis, effectiveAt int64) (*gen.CustomerCouponTemplate, string, error) {
	var template gen.CustomerCouponTemplate
	var hqID string
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		store, err := couponStoreScope(tx, principal, storeID, "franchiseCoupon:manage", authorization.AccessCreate)
		if err != nil {
			return err
		}
		hqID, err = couponHeadquartersID(tx)
		if err != nil {
			return err
		}
		existing, err := findStoreCouponTemplateByRequest(tx, hqID, storeID, input.RequestKey)
		if err != nil {
			return err
		}
		if existing != nil {
			if !sameTemplateIntent(existing, input) {
				return auth.NewError(auth.CodeConflict)
			}
			template = *existing
			return nil
		}
		template = gen.CustomerCouponTemplate{ID: uuid.Must(uuid.NewV4()).String(), OrganizationID: hqID, ApplicableStoreID: &store.ID,
			IssuerScope: gen.CouponIssuerScopeStore, CreatedAt: nowMillis, EffectiveAt: effectiveAt, DistributionEndsAt: input.DistributionEndsAt, Code: input.Code, Title: input.Title,
			RequestKey: input.RequestKey, AmountFen: input.AmountFen, MinSpendFen: input.MinSpendFen,
			DaysAfterActivation: input.DaysAfterActivation, PerMemberLimit: input.PerMemberLimit,
			TotalIssueLimit: input.TotalIssueLimit, Enabled: input.Enabled}
		insert := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&template)
		if insert.Error != nil {
			return insert.Error
		}
		if insert.RowsAffected != 1 {
			return auth.NewError(auth.CodeConflict)
		}
		return s.finishStoreCouponTemplateCreate(tx, principal, store, &template)
	})
	return &template, hqID, err
}

func findStoreCouponTemplateByRequest(tx *gorm.DB, hqID, storeID, key string) (*gen.CustomerCouponTemplate, error) {
	var template gen.CustomerCouponTemplate
	err := tx.Where("organization_id = ? AND applicable_store_id = ? AND issuer_scope = ? AND request_key = ?",
		hqID, storeID, gen.CouponIssuerScopeStore, key).First(&template).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (s *Service) SetStoreCouponTemplateEnabled(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, templateID string, enabled bool) (*gen.CustomerCouponTemplate, error) {
	var template gen.CustomerCouponTemplate
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		store, err := couponStoreScope(tx, principal, storeID, "franchiseCoupon:manage", authorization.AccessUpdate)
		if err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND applicable_store_id = ? AND issuer_scope = ?", templateID, store.ID, gen.CouponIssuerScopeStore).First(&template).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if template.Enabled == enabled {
			return nil
		}
		if err := tx.Model(&template).Update("enabled", enabled).Error; err != nil {
			return err
		}
		template.Enabled = enabled
		if enabled {
			if err := scheduleTemplateCouponDistribution(tx, &template, s.now().UnixMilli()); err != nil {
				return err
			}
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: "franchiseCoupon:manage",
			ResourceType: "customerCouponTemplate", ResourceID: template.ID, ResultCode: "SUCCESS"})
	})
	return &template, err
}
