package storemerchandising

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type StockAdjustment struct {
	StoreID, BatchID, PackageID, ReasonCode, RequestKey string
	Delta                                               int64
	Loss                                                bool
}

func (s *Service) AdjustStock(ctx context.Context, principal *auth.WorkspacePrincipal, input StockAdjustment) (*gen.StoreStockMovement, error) {
	input.ReasonCode = strings.TrimSpace(input.ReasonCode)
	if !validStockAdjustment(input) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	kind := gen.StockMovementKindCountAdjustment
	if input.Loss {
		kind = gen.StockMovementKindLoss
	}
	var movement gen.StoreStockMovement
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.adjustStockTx(tx, principal, input, kind, &movement)
	})
	return &movement, err
}

func validStockAdjustment(input StockAdjustment) bool {
	if !validAdjustmentTarget(input) {
		return false
	}
	if input.Delta == 0 || input.Delta < -math.MaxInt32 || input.Delta > math.MaxInt32 {
		return false
	}
	if len(input.ReasonCode) == 0 || len(input.ReasonCode) > 64 {
		return false
	}
	return validRequestKey(input.RequestKey) && (!input.Loss || input.Delta < 0)
}

func validAdjustmentTarget(input StockAdjustment) bool {
	return input.StoreID != "" && input.BatchID != "" && input.PackageID != ""
}

func (s *Service) adjustStockTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, input StockAdjustment, kind gen.StockMovementKind, movement *gen.StoreStockMovement) error {
	store, err := storeScope(tx, principal, input.StoreID, "franchiseStock:manage", authorization.AccessUpdate)
	if err != nil {
		return err
	}
	found, err := replayMovement(tx, store.ID, kind, input.RequestKey)
	if err != nil {
		return err
	}
	if found != nil {
		*movement = *found
		return validateAdjustReplay(movement, input)
	}
	batch, listing, err := stockBatchAndListing(tx, store.ID, input.BatchID)
	if err != nil {
		return err
	}
	pack, err := stockPackageForListing(tx, input.PackageID, listing.SkuID)
	if err != nil {
		return err
	}
	if err := applyAdjustment(tx, batch.ID, pack.ID, input.Delta); err != nil {
		return err
	}
	*movement = adjustmentMovement(store.ID, batch.ID, pack, kind, input)
	if err := tx.Create(movement).Error; err != nil {
		return err
	}
	return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: "franchiseStock:manage",
		ResourceType: "storeStockMovement", ResourceID: movement.ID, ResultCode: "SUCCESS", Metadata: audit.Metadata{ReasonCode: input.ReasonCode}})
}

func applyAdjustment(tx *gorm.DB, batchID, packageID string, delta int64) error {
	if delta > 0 {
		return addBalance(tx, batchID, packageID, delta)
	}
	return subtractBalance(tx, batchID, packageID, -delta)
}

func adjustmentMovement(storeID, batchID string, pack *gen.ProductPackage, kind gen.StockMovementKind, input StockAdjustment) gen.StoreStockMovement {
	result := gen.StoreStockMovement{ID: uuid.Must(uuid.NewV4()).String(), StoreID: storeID, BatchID: batchID,
		Kind: kind, RequestKey: input.RequestKey, PackageSetVersion: pack.PackageSetVersion, ReasonCode: &input.ReasonCode, OccurredAt: time.Now()}
	if input.Delta > 0 {
		result.TargetPackageID = &pack.ID
		result.TargetQuantity = &input.Delta
	} else {
		quantity := -input.Delta
		result.SourcePackageID = &pack.ID
		result.SourceQuantity = &quantity
	}
	return result
}

func validateAdjustReplay(found *gen.StoreStockMovement, input StockAdjustment) error {
	if found.BatchID != input.BatchID || found.ReasonCode == nil || *found.ReasonCode != input.ReasonCode {
		return auth.NewError(auth.CodeConflict)
	}
	if !sameAdjustmentQuantity(found, input) {
		return auth.NewError(auth.CodeConflict)
	}
	return nil
}

func sameAdjustmentQuantity(found *gen.StoreStockMovement, input StockAdjustment) bool {
	if input.Delta > 0 {
		return found.TargetPackageID != nil && *found.TargetPackageID == input.PackageID &&
			found.TargetQuantity != nil && *found.TargetQuantity == input.Delta
	}
	return found.SourcePackageID != nil && *found.SourcePackageID == input.PackageID &&
		found.SourceQuantity != nil && *found.SourceQuantity == -input.Delta
}
