package customer

import (
	"context"
	"math"
	"strings"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FixedAmountTemplateInput struct {
	Code, Title, RequestKey             string
	AmountFen, MinSpendFen              int64
	EffectiveAt, DistributionEndsAt     *int64
	DaysAfterActivation, PerMemberLimit int64
	TotalIssueLimit                     int64
	Enabled                             bool
}

const (
	minCouponUnixMillis int64 = 1_000_000_000_000
	maxCouponUnixMillis int64 = 9_999_999_999_999
)

func validCouponUnixMillis(value int64) bool {
	return value >= minCouponUnixMillis && value <= maxCouponUnixMillis
}

func validCouponTemplate(input FixedAmountTemplateInput) bool {
	return validTemplateText(input.Code, 32) && validTemplateText(input.Title, 128) &&
		validTemplateText(input.RequestKey, 128) && validCouponTemplateNumbers(input) &&
		validCouponTemplateTimes(input)
}

func validCouponTemplateNumbers(input FixedAmountTemplateInput) bool {
	return input.AmountFen > 0 && input.AmountFen <= math.MaxInt32 &&
		input.MinSpendFen >= 0 && input.MinSpendFen <= math.MaxInt32 &&
		input.DaysAfterActivation >= 1 && input.DaysAfterActivation <= 365 &&
		input.PerMemberLimit > 0 && input.PerMemberLimit <= math.MaxInt32 &&
		input.TotalIssueLimit >= 0 && input.TotalIssueLimit <= math.MaxInt32
}

func validCouponTemplateTimes(input FixedAmountTemplateInput) bool {
	if input.EffectiveAt != nil && !validCouponUnixMillis(*input.EffectiveAt) {
		return false
	}
	if input.DistributionEndsAt != nil && !validCouponUnixMillis(*input.DistributionEndsAt) {
		return false
	}
	return input.EffectiveAt == nil || input.DistributionEndsAt == nil || *input.DistributionEndsAt > *input.EffectiveAt
}

func validTemplateText(value string, maximum int) bool {
	return strings.TrimSpace(value) != "" && len(value) <= maximum
}

func (s *Service) CreateCouponTemplate(ctx context.Context, principal *auth.WorkspacePrincipal, input FixedAmountTemplateInput) (*gen.CustomerCouponTemplate, error) {
	hqID, err := headquartersID(principal, "hqCustomerCoupon:manage", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	if !validCouponTemplate(input) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	nowMillis := s.now().UnixMilli()
	effectiveAt, err := normalizedCouponTemplateTimes(input, nowMillis)
	if err != nil {
		return nil, err
	}
	template := &gen.CustomerCouponTemplate{
		ID: uuid.Must(uuid.NewV4()).String(), OrganizationID: hqID,
		IssuerScope: gen.CouponIssuerScopeHeadquarters,
		CreatedAt:   nowMillis, EffectiveAt: effectiveAt, DistributionEndsAt: input.DistributionEndsAt,
		Code: input.Code, Title: input.Title, RequestKey: input.RequestKey, AmountFen: input.AmountFen,
		MinSpendFen: input.MinSpendFen, DaysAfterActivation: input.DaysAfterActivation,
		PerMemberLimit: input.PerMemberLimit, TotalIssueLimit: input.TotalIssueLimit,
		Enabled: input.Enabled,
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var createErr error
		template, createErr = s.createCouponTemplateInTransaction(tx, principal, template, input)
		return createErr
	})
	if auth.ErrorCode(err) == auth.CodeConflict {
		replay, findErr := findCouponTemplateByRequest(s.db.WithContext(ctx), hqID, input.RequestKey)
		if findErr != nil {
			return nil, findErr
		}
		if replay != nil && sameTemplateIntent(replay, input) {
			return replay, nil
		}
	}
	return template, err
}

func normalizedCouponTemplateTimes(input FixedAmountTemplateInput, nowMillis int64) (int64, error) {
	if !validCouponUnixMillis(nowMillis) {
		return 0, auth.NewError(auth.CodeValidationFailed)
	}
	effectiveAt := nowMillis
	if input.EffectiveAt != nil {
		effectiveAt = *input.EffectiveAt
	}
	if input.DistributionEndsAt != nil && *input.DistributionEndsAt <= effectiveAt {
		return 0, auth.NewError(auth.CodeValidationFailed)
	}
	return effectiveAt, nil
}

func (s *Service) createCouponTemplateInTransaction(tx *gorm.DB, principal *auth.WorkspacePrincipal, template *gen.CustomerCouponTemplate, input FixedAmountTemplateInput) (*gen.CustomerCouponTemplate, error) {
	replay, err := findCouponTemplateByRequest(tx, template.OrganizationID, input.RequestKey)
	if err != nil {
		return nil, err
	}
	if replay != nil {
		if !sameTemplateIntent(replay, input) {
			return nil, auth.NewError(auth.CodeConflict)
		}
		return replay, nil
	}
	insert := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(template)
	if insert.Error != nil {
		return nil, insert.Error
	}
	if insert.RowsAffected != 1 {
		return nil, auth.NewError(auth.CodeConflict)
	}
	if err := scheduleTemplateCouponDistribution(tx, template, template.CreatedAt); err != nil {
		return nil, err
	}
	return template, s.auditCouponTemplate(tx, principal, template.ID, "create")
}

func findCouponTemplateByRequest(tx *gorm.DB, hqID, key string) (*gen.CustomerCouponTemplate, error) {
	var template gen.CustomerCouponTemplate
	err := tx.Where("organization_id = ? AND issuer_scope = ? AND request_key = ?", hqID, gen.CouponIssuerScopeHeadquarters, key).First(&template).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func sameTemplateIntent(template *gen.CustomerCouponTemplate, input FixedAmountTemplateInput) bool {
	return couponEffectiveIntentMatches(template, input.EffectiveAt) &&
		sameOptionalMillis(template.DistributionEndsAt, input.DistributionEndsAt) &&
		template.Code == input.Code && template.Title == input.Title &&
		template.AmountFen == input.AmountFen && template.MinSpendFen == input.MinSpendFen &&
		template.DaysAfterActivation == input.DaysAfterActivation &&
		template.PerMemberLimit == input.PerMemberLimit && template.TotalIssueLimit == input.TotalIssueLimit &&
		template.Enabled == input.Enabled
}

func couponEffectiveIntentMatches(template *gen.CustomerCouponTemplate, effectiveAt *int64) bool {
	if effectiveAt == nil {
		return template.EffectiveAt == template.CreatedAt
	}
	return template.EffectiveAt == *effectiveAt
}

func sameOptionalMillis(left, right *int64) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func couponTemplateDistributable(template *gen.CustomerCouponTemplate, nowMillis int64) bool {
	return template.Enabled && template.EffectiveAt <= nowMillis &&
		(template.DistributionEndsAt == nil || nowMillis < *template.DistributionEndsAt)
}

func autoCouponRequestKey(templateID, memberID string) string {
	return "AUTO:" + templateID + ":" + memberID
}

func (s *Service) SetCouponTemplateEnabled(ctx context.Context, principal *auth.WorkspacePrincipal, templateID string, enabled bool) (*gen.CustomerCouponTemplate, error) {
	hqID, err := headquartersID(principal, "hqCustomerCoupon:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	var template gen.CustomerCouponTemplate
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND organization_id = ? AND issuer_scope = ?", templateID, hqID, gen.CouponIssuerScopeHeadquarters).First(&template).Error; err != nil {
			return err
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
		return s.auditCouponTemplate(tx, principal, template.ID, "toggle")
	})
	return &template, err
}

func (s *Service) auditCouponTemplate(tx *gorm.DB, principal *auth.WorkspacePrincipal, templateID, reasonCode string) error {
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: principal.OrganizationID, Action: "hqCustomerCoupon:manage",
		ResourceType: "customerCouponTemplate", ResourceID: templateID, ResultCode: "SUCCESS",
		Metadata: audit.MetadataForPrincipal(principal, audit.Metadata{ReasonCode: reasonCode}),
	})
}
