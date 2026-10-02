package storemerchandising

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type CreateStocktakeInput struct {
	StoreID, RequestKey  string
	ListingIDs, BatchIDs []string
}

const maxStocktakeLines = 200

func (s *Service) CreateStocktake(ctx context.Context, principal *auth.WorkspacePrincipal, input CreateStocktakeInput) (*StocktakeView, error) {
	if input.StoreID == "" || !validRequestKey(input.RequestKey) || len(input.ListingIDs)+len(input.BatchIDs) == 0 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	digest, err := stocktakeScopeDigest(input)
	if err != nil {
		return nil, err
	}
	var id string
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		id, err = s.createStocktakeTx(tx, principal, input, digest)
		return err
	})
	if err != nil {
		if duplicateStocktakeRequest(err) {
			id, lookupErr := existingStocktakeID(s.db.WithContext(ctx), input.StoreID, input.RequestKey, digest)
			if lookupErr != nil {
				return nil, lookupErr
			}
			if id != "" {
				return readStocktake(s.db.WithContext(ctx), input.StoreID, id)
			}
		}
		return nil, err
	}
	return readStocktake(s.db.WithContext(ctx), input.StoreID, id)
}

func duplicateStocktakeRequest(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 &&
		strings.Contains(mysqlErr.Message, "uidx_store_stocktake_request")
}

func (s *Service) createStocktakeTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, input CreateStocktakeInput, digest string) (string, error) {
	store, err := storeScope(tx, principal, input.StoreID, "franchiseStocktake:record", authorization.AccessCreate)
	if err != nil {
		return "", err
	}
	id, err := existingStocktakeID(tx, store.ID, input.RequestKey, digest)
	if err != nil || id != "" {
		return id, err
	}
	batches, err := stocktakeBatches(tx, store.ID, input)
	if err != nil {
		return "", err
	}
	sheet := gen.StoreStocktake{ID: uuid.Must(uuid.NewV4()).String(), StoreID: store.ID,
		InitiatedByAccountID: principal.AccountID, Status: gen.StocktakeStatusCounting,
		RequestKey: input.RequestKey, ScopeDigest: digest, StartedAt: time.Now()}
	if err := tx.Create(&sheet).Error; err != nil {
		return "", err
	}
	if err := createStocktakeLines(tx, &sheet, batches); err != nil {
		return "", err
	}
	return sheet.ID, s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: "franchiseStocktake:record",
		ResourceType: "storeStocktake", ResourceID: sheet.ID, ResultCode: "SUCCESS"})
}

func existingStocktakeID(tx *gorm.DB, storeID, requestKey, digest string) (string, error) {
	var old gen.StoreStocktake
	err := tx.Where("store_id = ? AND request_key = ?", storeID, requestKey).First(&old).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if old.ScopeDigest != digest {
		return "", auth.NewError(auth.CodeConflict)
	}
	return old.ID, nil
}

func stocktakeScopeDigest(input CreateStocktakeInput) (string, error) {
	listings, batches := append([]string(nil), input.ListingIDs...), append([]string(nil), input.BatchIDs...)
	sort.Strings(listings)
	sort.Strings(batches)
	for _, ids := range [][]string{listings, batches} {
		for i, id := range ids {
			if id == "" || i > 0 && ids[i-1] == id {
				return "", auth.NewError(auth.CodeValidationFailed)
			}
		}
	}
	payload, err := json.Marshal([][]string{listings, batches})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func stocktakeBatches(tx *gorm.DB, storeID string, input CreateStocktakeInput) ([]gen.StoreInventoryBatch, error) {
	if err := validateStocktakeScope(tx, storeID, input); err != nil {
		return nil, err
	}
	var batches []gen.StoreInventoryBatch
	query := tx.Table("store_inventory_batches AS b").Joins("JOIN store_listings AS l ON l.id = b.listing_id").
		Where("l.store_id = ?", storeID)
	if len(input.ListingIDs) > 0 && len(input.BatchIDs) > 0 {
		query = query.Where("(b.listing_id IN ? OR b.id IN ?)", input.ListingIDs, input.BatchIDs)
	} else if len(input.ListingIDs) > 0 {
		query = query.Where("b.listing_id IN ?", input.ListingIDs)
	} else {
		query = query.Where("b.id IN ?", input.BatchIDs)
	}
	if err := query.Select("b.*").Find(&batches).Error; err != nil {
		return nil, err
	}
	if len(batches) == 0 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	return batches, nil
}

func validateStocktakeScope(tx *gorm.DB, storeID string, input CreateStocktakeInput) error {
	if len(input.ListingIDs) > 0 {
		var count int64
		if err := tx.Model(&gen.StoreListing{}).Where("store_id = ? AND id IN ?", storeID, input.ListingIDs).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(input.ListingIDs)) {
			return auth.NewError(auth.CodeValidationFailed)
		}
	}
	if len(input.BatchIDs) > 0 {
		var count int64
		if err := tx.Table("store_inventory_batches AS b").Joins("JOIN store_listings AS l ON l.id = b.listing_id").
			Where("l.store_id = ? AND b.id IN ?", storeID, input.BatchIDs).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(input.BatchIDs)) {
			return auth.NewError(auth.CodeValidationFailed)
		}
	}
	return nil
}

func createStocktakeLines(tx *gorm.DB, sheet *gen.StoreStocktake, batches []gen.StoreInventoryBatch) error {
	lineCount := 0
	for _, batch := range batches {
		var listing gen.StoreListing
		if err := tx.First(&listing, "id = ?", batch.ListingID).Error; err != nil {
			return err
		}
		var balances []gen.StoreStockBalance
		if err := tx.Where("batch_id = ?", batch.ID).Find(&balances).Error; err != nil {
			return err
		}
		byPackage, ids := map[string]gen.StoreStockBalance{}, make([]string, 0, len(balances))
		for _, balance := range balances {
			byPackage[balance.PackageID] = balance
			ids = append(ids, balance.PackageID)
		}
		var packages []gen.ProductPackage
		if err := tx.Where("sku_id = ? AND (enabled = ? OR id IN ?)", listing.SkuID, true, ids).Find(&packages).Error; err != nil {
			return err
		}
		for _, pack := range packages {
			lineCount++
			if lineCount > maxStocktakeLines {
				return auth.NewError(auth.CodeValidationFailed)
			}
			balance := byPackage[pack.ID]
			line := gen.StoreStocktakeLine{ID: uuid.Must(uuid.NewV4()).String(), StocktakeID: sheet.ID,
				BatchID: batch.ID, PackageID: pack.ID, PackageSetVersion: pack.PackageSetVersion,
				BatchNumberSnapshot: batch.BatchNumber, ExpiresAtSnapshot: batch.ExpiresAt,
				PackageNameSnapshot: pack.Name, PackageEnabledSnapshot: pack.Enabled,
				SnapshotQuantity: balance.Quantity, SnapshotVersion: balance.Version}
			if err := tx.Create(&line).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
