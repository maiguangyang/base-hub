package storemerchandising

import (
	"context"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type StocktakeLineView struct {
	ID, BatchID, PackageID              string
	ListingID, BatchNumber, PackageName string
	ExpiresAt                           *time.Time
	PackageEnabled                      bool
	PackageSetVersion                   int64
	CountedQuantity                     *int64
	SnapshotQuantity                    *int64
	ReasonCode                          *string
	ReasonNote                          *string
	CountHistory                        []StocktakeCountEvent
	NeedsRecount                        bool
}

type StocktakeCountEvent struct {
	ActorAccountID  string
	CountedQuantity int64
	CountedAt       time.Time
}

type StocktakeView struct {
	ID, StoreID                      string
	RequestKey, InitiatedByAccountID string
	PostedByID                       *string
	StartedAt                        time.Time
	ReviewedAt, PostedAt, CanceledAt *time.Time
	Status                           gen.StocktakeStatus
	Lines                            []StocktakeLineView
	AddLineChoices                   []StocktakeAddLineChoice
}

type StocktakeAddLineChoice struct {
	BatchID, PackageID, PackageName string
	PackageEnabled                  bool
}

type StocktakeFilter struct {
	Status             gen.StocktakeStatus
	ListingID, BatchID string
	HasDifference      *bool
	From, To           *time.Time
	IncludeHistory     *bool
}

func (s *Service) StocktakeBatchChoices(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, listingID string, page, perPage int) (*Page[gen.StoreInventoryBatch], error) {
	tx := s.db.WithContext(ctx)
	action := scopedReadAction(principal, "franchiseStocktake:read", "franchiseStocktake:record")
	if _, err := storeScope(tx, principal, storeID, action, authorization.AccessRead); err != nil {
		return nil, err
	}
	offset, limit, err := bounds(page, perPage)
	if err != nil {
		return nil, err
	}
	var listing gen.StoreListing
	if err := tx.Where("id = ? AND store_id = ?", listingID, storeID).First(&listing).Error; err != nil {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	query := tx.Model(&gen.StoreInventoryBatch{}).Where("listing_id = ?", listing.ID)
	result := &Page[gen.StoreInventoryBatch]{Page: page, PerPage: perPage, Data: []gen.StoreInventoryBatch{}}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Order("expires_at, id").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) Stocktakes(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, filter StocktakeFilter, page, perPage int) (*Page[StocktakeView], error) {
	tx := s.db.WithContext(ctx)
	if _, err := storeScope(tx, principal, storeID, "franchiseStocktake:read", authorization.AccessRead); err != nil {
		return nil, err
	}
	offset, limit, err := bounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := filteredStocktakeQuery(tx, storeID, filter)
	result := &Page[StocktakeView]{Page: page, PerPage: perPage, Data: []StocktakeView{}}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	var sheets []gen.StoreStocktake
	if err := query.Order("started_at DESC, id").Offset(offset).Limit(limit).Find(&sheets).Error; err != nil {
		return nil, err
	}
	includeHistory := filter.IncludeHistory == nil || *filter.IncludeHistory
	result.Data, err = stocktakeViews(tx, sheets, includeHistory)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func filteredStocktakeQuery(tx *gorm.DB, storeID string, filter StocktakeFilter) *gorm.DB {
	query := tx.Model(&gen.StoreStocktake{}).Where("store_id = ?", storeID)
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	query = dateFilteredStocktakeQuery(query, filter)
	if filter.ListingID != "" {
		query = query.Where(`EXISTS (SELECT 1 FROM store_stocktake_lines AS line
			JOIN store_inventory_batches AS batch ON batch.id = line.batch_id
			WHERE line.stocktake_id = store_stocktakes.id AND batch.listing_id = ?)`, filter.ListingID)
	}
	if filter.BatchID != "" {
		query = query.Where(`EXISTS (SELECT 1 FROM store_stocktake_lines AS line
			WHERE line.stocktake_id = store_stocktakes.id AND line.batch_id = ?)`, filter.BatchID)
	}
	if filter.HasDifference != nil {
		query = query.Where("status IN ?", []gen.StocktakeStatus{gen.StocktakeStatusReview, gen.StocktakeStatusPosted})
		expression := `EXISTS (SELECT 1 FROM store_stocktake_lines AS line
			WHERE line.stocktake_id = store_stocktakes.id AND line.counted_quantity IS NOT NULL
			AND line.counted_quantity <> line.snapshot_quantity)`
		if *filter.HasDifference {
			query = query.Where(expression)
		} else {
			query = query.Where("NOT " + expression)
		}
	}
	return query
}

func (s *Service) Stocktake(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, id string) (*StocktakeView, error) {
	tx := s.db.WithContext(ctx)
	if _, err := storeScope(tx, principal, storeID, "franchiseStocktake:read", authorization.AccessRead); err != nil {
		return nil, err
	}
	return readStocktake(tx, storeID, id)
}

func readStocktake(tx *gorm.DB, storeID, id string) (*StocktakeView, error) {
	var sheet gen.StoreStocktake
	if err := tx.Where("id = ? AND store_id = ?", id, storeID).First(&sheet).Error; err != nil {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	views, err := stocktakeViews(tx, []gen.StoreStocktake{sheet}, true)
	if err != nil {
		return nil, err
	}
	return &views[0], nil
}
