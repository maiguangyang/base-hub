package customer

import (
	"context"
	"strings"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type StoreCouponMemberView struct{ ID, MemberNumber, PhoneMasked string }

func (s *Service) LookupStoreCouponMember(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, identifier string) (*StoreCouponMemberView, error) {
	store, err := couponStoreScope(s.db.WithContext(ctx), principal, storeID, "franchiseCoupon:read", authorization.AccessRead)
	if err != nil {
		return nil, err
	}
	identifier = strings.TrimSpace(identifier)
	query := s.db.WithContext(ctx).Model(&gen.CustomerMember{})
	if _, err := uuid.FromString(identifier); err == nil {
		query = query.Where("id = ?", identifier)
	} else {
		phone, err := NormalizeCNPhone(identifier)
		if err != nil {
			return nil, auth.NewError(auth.CodeValidationFailed)
		}
		query = query.Where("phone = ?", phone)
	}
	if !s.lookupLimiter.Allow(principal.AccountID+":"+store.ID, s.now()) {
		return nil, auth.NewError(auth.CodeRateLimited)
	}
	hqID, err := couponHeadquartersID(s.db.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	var member gen.CustomerMember
	err = query.Where("organization_id = ? AND status = ?", hqID, gen.CustomerMemberStatusActive).First(&member).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	view := &StoreCouponMemberView{ID: member.ID, MemberNumber: member.MemberNumber}
	if member.Phone != nil {
		view.PhoneMasked = MaskCNPhone(*member.Phone)
	}
	return view, nil
}

func (s *Service) ListStoreCouponTemplates(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, page, perPage int64) (CouponTemplatePage, error) {
	return s.ListStoreCouponTemplatesFiltered(ctx, principal, storeID, CouponTemplateFilter{}, page, perPage)
}

func (s *Service) ListStoreCouponTemplatesFiltered(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, filter CouponTemplateFilter, page, perPage int64) (CouponTemplatePage, error) {
	if _, err := couponStoreScope(s.db.WithContext(ctx), principal, storeID, "franchiseCoupon:read", authorization.AccessRead); err != nil {
		return CouponTemplatePage{}, err
	}
	if !validCouponPage(page, perPage) {
		return CouponTemplatePage{}, auth.NewError(auth.CodeValidationFailed)
	}
	query := s.db.WithContext(ctx).Model(&gen.CustomerCouponTemplate{}).Where("applicable_store_id = ? AND issuer_scope = ?", storeID, gen.CouponIssuerScopeStore)
	query = applyCouponTemplateFilter(query, filter)
	var result CouponTemplatePage
	if err := query.Count(&result.Total).Error; err != nil {
		return result, err
	}
	err := query.Order("created_at DESC, id DESC").Offset(int((page - 1) * perPage)).Limit(int(perPage)).Find(&result.Data).Error
	return result, err
}

func (s *Service) ListStoreCouponGrants(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, templateID string, page, perPage int64) (CouponGrantPage, error) {
	return s.ListStoreCouponGrantsFiltered(ctx, principal, storeID, templateID, CouponGrantFilter{}, page, perPage)
}

func (s *Service) ListStoreCouponGrantsFiltered(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, templateID string, filter CouponGrantFilter, page, perPage int64) (CouponGrantPage, error) {
	if _, err := couponStoreScope(s.db.WithContext(ctx), principal, storeID, "franchiseCoupon:read", authorization.AccessRead); err != nil {
		return CouponGrantPage{}, err
	}
	if !validCouponPage(page, perPage) {
		return CouponGrantPage{}, auth.NewError(auth.CodeValidationFailed)
	}
	var template gen.CustomerCouponTemplate
	if err := s.db.WithContext(ctx).Where("id = ? AND applicable_store_id = ? AND issuer_scope = ?", templateID, storeID, gen.CouponIssuerScopeStore).First(&template).Error; err != nil {
		return CouponGrantPage{}, auth.NewError(auth.CodePermissionDenied)
	}
	query := s.db.WithContext(ctx).Model(&gen.CustomerCouponGrant{}).Where("template_id = ?", templateID)
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	var result CouponGrantPage
	if err := query.Count(&result.Total).Error; err != nil {
		return result, err
	}
	err := query.Order("issued_at DESC, id DESC").Offset(int((page - 1) * perPage)).Limit(int(perPage)).Find(&result.Data).Error
	return result, err
}
