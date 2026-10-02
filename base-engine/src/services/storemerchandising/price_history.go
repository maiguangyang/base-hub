package storemerchandising

import (
	"context"
	"strings"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
)

func (s *Service) PriceHistory(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, offerID string, page, perPage int) (*Page[*gen.StorePriceRevision], error) {
	return s.PriceHistoryFiltered(ctx, principal, storeID, offerID, PriceHistoryFilter{}, page, perPage)
}

type PriceHistoryFilter struct {
	Q    *string
	From *time.Time
	To   *time.Time
}

func (s *Service) PriceHistoryFiltered(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, offerID string, filter PriceHistoryFilter, page, perPage int) (*Page[*gen.StorePriceRevision], error) {
	if _, err := storeScope(s.db.WithContext(ctx), principal, storeID, "franchiseProduct:read", authorization.AccessRead); err != nil {
		return nil, err
	}
	offset, limit, err := bounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&gen.StorePriceRevision{}).
		Joins("JOIN store_package_offers ON store_package_offers.id = store_price_revisions.offer_id").
		Joins("JOIN store_listings ON store_listings.id = store_package_offers.listing_id").
		Where("store_listings.store_id = ? AND store_package_offers.id = ?", storeID, offerID)
	if filter.Q != nil && strings.TrimSpace(*filter.Q) != "" {
		query = query.Where("store_price_revisions.reason_code LIKE ?", "%"+strings.TrimSpace(*filter.Q)+"%")
	}
	if filter.From != nil {
		query = query.Where("store_price_revisions.effective_at >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("store_price_revisions.effective_at <= ?", *filter.To)
	}
	result := &Page[*gen.StorePriceRevision]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Order("store_price_revisions.effective_at DESC, store_price_revisions.id DESC").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}
