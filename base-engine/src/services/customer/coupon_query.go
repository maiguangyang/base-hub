package customer

import (
	"context"
	"strings"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type CouponTemplatePage struct {
	Data  []gen.CustomerCouponTemplate
	Total int64
}

type CouponGrantPage struct {
	Data  []gen.CustomerCouponGrant
	Total int64
}

type CouponTemplateFilter struct {
	Q       *string
	Enabled *bool
}

type CouponGrantFilter struct {
	TemplateID *string
	Status     *gen.CustomerCouponGrantStatus
}

func applyCouponTemplateFilter(query *gorm.DB, filter CouponTemplateFilter) *gorm.DB {
	if filter.Q != nil && strings.TrimSpace(*filter.Q) != "" {
		q := "%" + strings.TrimSpace(*filter.Q) + "%"
		query = query.Where("(title LIKE ? OR code LIKE ?)", q, q)
	}
	if filter.Enabled != nil {
		query = query.Where("enabled = ?", *filter.Enabled)
	}
	return query
}

func validCouponPage(page, perPage int64) bool {
	return page >= 1 && page <= 1_000_000 && perPage >= 1 && perPage <= 50
}

func (s *Service) ListCouponTemplates(ctx context.Context, principal *auth.WorkspacePrincipal, page, perPage int64) (CouponTemplatePage, error) {
	return s.ListCouponTemplatesFiltered(ctx, principal, CouponTemplateFilter{}, page, perPage)
}

func (s *Service) ListCouponTemplatesFiltered(ctx context.Context, principal *auth.WorkspacePrincipal, filter CouponTemplateFilter, page, perPage int64) (CouponTemplatePage, error) {
	hqID, err := headquartersID(principal, "hqCustomerCoupon:read", authorization.AccessRead)
	if err != nil {
		return CouponTemplatePage{}, err
	}
	if !validCouponPage(page, perPage) {
		return CouponTemplatePage{}, auth.NewError(auth.CodeValidationFailed)
	}
	query := s.db.WithContext(ctx).Model(&gen.CustomerCouponTemplate{}).Where("organization_id = ? AND issuer_scope = ?", hqID, gen.CouponIssuerScopeHeadquarters)
	query = applyCouponTemplateFilter(query, filter)
	var result CouponTemplatePage
	if err := query.Count(&result.Total).Error; err != nil {
		return CouponTemplatePage{}, err
	}
	err = query.Order("created_at DESC, id DESC").
		Offset(int((page - 1) * perPage)).Limit(int(perPage)).Find(&result.Data).Error
	return result, err
}

func (s *Service) ListCouponGrants(ctx context.Context, principal *auth.WorkspacePrincipal, memberID string, page, perPage int64) (CouponGrantPage, error) {
	return s.ListCouponGrantsFiltered(ctx, principal, memberID, CouponGrantFilter{}, page, perPage)
}

func (s *Service) ListCouponGrantsFiltered(ctx context.Context, principal *auth.WorkspacePrincipal, memberID string, filter CouponGrantFilter, page, perPage int64) (CouponGrantPage, error) {
	hqID, err := headquartersID(principal, "hqCustomerCoupon:read", authorization.AccessRead)
	if err != nil {
		return CouponGrantPage{}, err
	}
	if !validCouponPage(page, perPage) {
		return CouponGrantPage{}, auth.NewError(auth.CodeValidationFailed)
	}
	var member gen.CustomerMember
	err = s.db.WithContext(ctx).Where("id = ? AND organization_id = ?", memberID, hqID).First(&member).Error
	if err == gorm.ErrRecordNotFound {
		return CouponGrantPage{}, auth.NewError(auth.CodePermissionDenied)
	}
	if err != nil {
		return CouponGrantPage{}, err
	}
	query := s.db.WithContext(ctx).Model(&gen.CustomerCouponGrant{}).
		Joins("JOIN customer_coupon_templates ON customer_coupon_templates.id = customer_coupon_grants.template_id").
		Where("customer_coupon_grants.member_id = ? AND customer_coupon_templates.organization_id = ? AND customer_coupon_templates.issuer_scope = ?", memberID, hqID, gen.CouponIssuerScopeHeadquarters)
	if filter.TemplateID != nil {
		query = query.Where("customer_coupon_grants.template_id = ?", *filter.TemplateID)
	}
	if filter.Status != nil {
		query = query.Where("customer_coupon_grants.status = ?", *filter.Status)
	}
	var result CouponGrantPage
	if err := query.Count(&result.Total).Error; err != nil {
		return CouponGrantPage{}, err
	}
	err = query.Order("customer_coupon_grants.issued_at DESC, customer_coupon_grants.id DESC").
		Offset(int((page - 1) * perPage)).Limit(int(perPage)).Find(&result.Data).Error
	return result, err
}
