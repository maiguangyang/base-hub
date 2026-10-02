package storemerchandising

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
)

type PromotionView struct {
	Promotion *gen.StorePromotion
	Targets   []*gen.StorePromotionTarget
	Latest    bool
}

func (s *Service) Promotions(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, page, perPage int) (*Page[PromotionView], error) {
	return s.PromotionsFiltered(ctx, principal, storeID, PromotionFilter{}, page, perPage)
}

type PromotionFilter struct {
	Kind    *gen.StorePromotionKind
	Enabled *bool
}

func (s *Service) PromotionsFiltered(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, filter PromotionFilter, page, perPage int) (*Page[PromotionView], error) {
	if _, err := storeScope(s.db.WithContext(ctx), principal, storeID, "franchisePromotion:read", authorization.AccessRead); err != nil {
		return nil, err
	}
	offset, limit, err := bounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&gen.StorePromotion{}).Where("store_id = ?", storeID)
	if filter.Kind != nil {
		query = query.Where("kind = ?", *filter.Kind)
	}
	if filter.Enabled != nil {
		query = query.Where("enabled = ?", *filter.Enabled)
	}
	result := &Page[PromotionView]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	var promotions []*gen.StorePromotion
	if err := query.Order("starts_at DESC, id").Offset(offset).Limit(limit).Find(&promotions).Error; err != nil {
		return nil, err
	}
	for _, promotion := range promotions {
		var targets []*gen.StorePromotionTarget
		if err := s.db.WithContext(ctx).Where("promotion_id = ?", promotion.ID).Find(&targets).Error; err != nil {
			return nil, err
		}
		latest, err := s.promotionIsLatest(ctx, promotion)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, PromotionView{Promotion: promotion, Targets: targets, Latest: latest})
	}
	return result, nil
}

func (s *Service) Promotion(ctx context.Context, principal *auth.WorkspacePrincipal, id string) (*PromotionView, error) {
	var promotion gen.StorePromotion
	if err := s.db.WithContext(ctx).First(&promotion, "id = ?", id).Error; err != nil {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	if _, err := storeScope(s.db.WithContext(ctx), principal, promotion.StoreID, "franchisePromotion:read", authorization.AccessRead); err != nil {
		return nil, err
	}
	var targets []*gen.StorePromotionTarget
	if err := s.db.WithContext(ctx).Where("promotion_id = ?", id).Find(&targets).Error; err != nil {
		return nil, err
	}
	latest, err := s.promotionIsLatest(ctx, &promotion)
	if err != nil {
		return nil, err
	}
	return &PromotionView{Promotion: &promotion, Targets: targets, Latest: latest}, nil
}

func (s *Service) promotionIsLatest(ctx context.Context, promotion *gen.StorePromotion) (bool, error) {
	var latest int64
	err := s.db.WithContext(ctx).Model(&gen.StorePromotion{}).Where("store_id = ? AND rule_key = ?", promotion.StoreID, promotion.RuleKey).Select("MAX(version)").Scan(&latest).Error
	return promotion.Version == latest, err
}
