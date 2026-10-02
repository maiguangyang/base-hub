package storemerchandising

import (
	"context"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type StockUnpack struct {
	StoreID, BatchID, SourcePackageID, RequestKey string
	Quantity                                      int64
}

func (s *Service) UnpackStock(ctx context.Context, principal *auth.WorkspacePrincipal, input StockUnpack) (*gen.StoreStockMovement, error) {
	if !validStockUnpack(input) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var movement gen.StoreStockMovement
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.unpackStockTx(tx, principal, input, &movement)
	})
	return &movement, err
}

func validStockUnpack(input StockUnpack) bool {
	return input.StoreID != "" && input.BatchID != "" && input.SourcePackageID != "" && input.Quantity > 0 && validRequestKey(input.RequestKey)
}

type unpackResolved struct {
	batch          *gen.StoreInventoryBatch
	source, target *gen.ProductPackage
	targetQuantity int64
}

func (s *Service) unpackStockTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, input StockUnpack, movement *gen.StoreStockMovement) error {
	store, err := storeScope(tx, principal, input.StoreID, "franchiseStock:manage", authorization.AccessUpdate)
	if err != nil {
		return err
	}
	found, err := replayMovement(tx, store.ID, gen.StockMovementKindUnpack, input.RequestKey)
	if err != nil {
		return err
	}
	if found != nil {
		*movement = *found
		return validateUnpackReplay(movement, input)
	}
	resolved, err := resolveUnpack(tx, store.ID, input)
	if err != nil {
		return err
	}
	if err := lockUnpackBalances(tx, resolved.batch.ID, resolved.source.ID, resolved.target.ID); err != nil {
		return err
	}
	if err := subtractBalance(tx, resolved.batch.ID, resolved.source.ID, input.Quantity); err != nil {
		return err
	}
	if err := addBalance(tx, resolved.batch.ID, resolved.target.ID, resolved.targetQuantity); err != nil {
		return err
	}
	*movement = unpackMovement(store.ID, input, resolved)
	if err := tx.Create(movement).Error; err != nil {
		return err
	}
	return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: "franchiseStock:manage",
		ResourceType: "storeStockMovement", ResourceID: movement.ID, ResultCode: "SUCCESS"})
}

func lockUnpackBalances(tx *gorm.DB, batchID, sourcePackageID, targetPackageID string) error {
	packageIDs := []string{sourcePackageID, targetPackageID}
	sort.Strings(packageIDs)
	for _, packageID := range packageIDs {
		if _, err := getBalance(tx, batchID, packageID); err != nil {
			return err
		}
	}
	return nil
}

func resolveUnpack(tx *gorm.DB, storeID string, input StockUnpack) (*unpackResolved, error) {
	batch, listing, err := stockBatchAndListing(tx, storeID, input.BatchID)
	if err != nil {
		return nil, err
	}
	source, err := stockPackageForListing(tx, input.SourcePackageID, listing.SkuID)
	if err != nil {
		return nil, err
	}
	if source.ContainsPackageID == nil || source.ContainsQuantity == nil {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	factor := *source.ContainsQuantity
	if factor <= 1 || input.Quantity > math.MaxInt32/factor {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var target gen.ProductPackage
	if err := tx.Where("id = ? AND sku_id = ? AND package_set_version = ?", *source.ContainsPackageID, listing.SkuID, source.PackageSetVersion).First(&target).Error; err != nil {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	return &unpackResolved{batch: batch, source: source, target: &target, targetQuantity: input.Quantity * factor}, nil
}

func unpackMovement(storeID string, input StockUnpack, resolved *unpackResolved) gen.StoreStockMovement {
	return gen.StoreStockMovement{ID: uuid.Must(uuid.NewV4()).String(), StoreID: storeID, BatchID: resolved.batch.ID,
		Kind: gen.StockMovementKindUnpack, RequestKey: input.RequestKey, SourceQuantity: &input.Quantity,
		TargetQuantity: &resolved.targetQuantity, FactorSnapshot: resolved.source.ContainsQuantity,
		PackageSetVersion: resolved.source.PackageSetVersion, SourcePackageID: &resolved.source.ID,
		TargetPackageID: &resolved.target.ID, OccurredAt: time.Now()}
}

func validRequestKey(value string) bool {
	return value != "" && value == strings.TrimSpace(value) && len(value) <= 128
}

func replayMovement(tx *gorm.DB, storeID string, kind gen.StockMovementKind, requestKey string) (*gen.StoreStockMovement, error) {
	var found gen.StoreStockMovement
	err := tx.Where("store_id = ? AND kind = ? AND request_key = ?", storeID, kind, requestKey).First(&found).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &found, nil
}

func validateUnpackReplay(found *gen.StoreStockMovement, input StockUnpack) error {
	if found.BatchID != input.BatchID || found.SourcePackageID == nil || *found.SourcePackageID != input.SourcePackageID || found.SourceQuantity == nil || *found.SourceQuantity != input.Quantity {
		return auth.NewError(auth.CodeConflict)
	}
	return nil
}

func validateReceiptReplay(tx *gorm.DB, found *gen.StoreStockMovement, input StockReceipt) error {
	if found.TargetPackageID == nil || *found.TargetPackageID != input.PackageID || found.TargetQuantity == nil || *found.TargetQuantity != input.Quantity {
		return auth.NewError(auth.CodeConflict)
	}
	var batch gen.StoreInventoryBatch
	if err := tx.Where("id = ? AND listing_id = ?", found.BatchID, input.ListingID).First(&batch).Error; err != nil {
		return auth.NewError(auth.CodeConflict)
	}
	if input.BatchNumber != "" && batch.BatchNumber != input.BatchNumber {
		return auth.NewError(auth.CodeConflict)
	}
	if !sameReceiptBatchTerms(&batch, input) {
		return auth.NewError(auth.CodeConflict)
	}
	return nil
}

func sameReceiptBatchTerms(batch *gen.StoreInventoryBatch, input StockReceipt) bool {
	return sameTime(batch.ProducedAt, input.ProducedAt) && sameTime(batch.ExpiresAt, input.ExpiresAt) &&
		sameOptionalString(batch.SourceReference, input.SourceReference)
}

func sameOptionalString(a, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}
