package storemerchandising

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

func (s *Service) SetStocktakeReason(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, id, lineID, reason string, note ...string) (*StocktakeView, error) {
	reason = strings.TrimSpace(reason)
	if !validStocktakeReason(reason) || len(note) > 1 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	comment := ""
	if len(note) == 1 {
		comment = strings.TrimSpace(note[0])
	}
	if utf8.RuneCountInString(comment) > 256 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		store, err := storeScope(tx, principal, storeID, "franchiseStocktake:post", authorization.AccessUpdate)
		if err != nil {
			return err
		}
		return s.setStocktakeReasonTx(tx, principal, store, id, lineID, reason, comment)
	})
	if err != nil {
		return nil, err
	}
	return readStocktake(s.db.WithContext(ctx), storeID, id)
}

func validStocktakeReason(reason string) bool {
	switch reason {
	case "COUNT_DIFFERENCE", "COUNT_OMISSION", "RECORD_ERROR", "PACKAGE_MISMATCH", "OTHER":
		return true
	}
	return false
}

func (s *Service) setStocktakeReasonTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, store *gen.Store, id, lineID, reason, note string) error {
	if _, err := stocktakeInStatus(tx, store.ID, id, gen.StocktakeStatusReview); err != nil {
		return err
	}
	line, err := stocktakeLine(tx, id, lineID)
	if err != nil {
		return err
	}
	if line.CountedQuantity == nil || *line.CountedQuantity == line.SnapshotQuantity {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if sameStocktakeReason(line, reason, note) {
		return nil
	}
	var reasonNote *string
	if note != "" {
		reasonNote = &note
	}
	if err := tx.Model(line).Updates(map[string]any{"reason_code": reason, "reason_note": reasonNote}).Error; err != nil {
		return err
	}
	return s.stocktakeAudit(tx, principal, store, id, "franchiseStocktake:post")
}

func sameStocktakeReason(line *gen.StoreStocktakeLine, reason, note string) bool {
	if line.ReasonCode == nil || *line.ReasonCode != reason {
		return false
	}
	if line.ReasonNote == nil {
		return note == ""
	}
	return *line.ReasonNote == note
}

func (s *Service) PostStocktake(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, id string) (*StocktakeView, error) {
	var outcome error
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		store, err := storeScope(tx, principal, storeID, "franchiseStocktake:post", authorization.AccessUpdate)
		if err != nil {
			return err
		}
		outcome, err = s.postStocktakeTx(tx, principal, store, id)
		return err
	})
	if err != nil {
		return nil, err
	}
	if outcome != nil {
		return nil, outcome
	}
	return readStocktake(s.db.WithContext(ctx), storeID, id)
}

func (s *Service) postStocktakeTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, store *gen.Store, id string) (error, error) {
	sheet, err := lockedStocktake(tx, store.ID, id)
	if err != nil {
		return nil, err
	}
	if sheet.Status == gen.StocktakeStatusPosted {
		return nil, nil
	}
	if sheet.Status != gen.StocktakeStatusReview {
		return nil, auth.NewError(auth.CodeConflict)
	}
	lines, err := stocktakeLines(tx, id)
	if err != nil {
		return nil, err
	}
	outcome, err := validateStocktakePost(tx, store.ID, lines)
	if err != nil {
		return nil, err
	}
	if outcome != nil {
		return outcome, tx.Model(sheet).Update("status", gen.StocktakeStatusCounting).Error
	}
	if err := writeStocktakeDifferences(tx, store.ID, lines); err != nil {
		return nil, err
	}
	if err := tx.Model(sheet).Updates(map[string]any{"status": gen.StocktakeStatusPosted,
		"posted_at": time.Now(), "posted_by_id": principal.AccountID}).Error; err != nil {
		return nil, err
	}
	return nil, s.stocktakeAudit(tx, principal, store, id, "franchiseStocktake:post")
}

func validateStocktakePost(tx *gorm.DB, storeID string, lines []gen.StoreStocktakeLine) (error, error) {
	if err := requireCountedStocktakeLines(lines); err != nil {
		return nil, err
	}
	if err := requireStocktakeReasons(lines); err != nil {
		return nil, err
	}
	if err := assertStocktakeTargets(tx, storeID, lines); err != nil {
		return nil, stocktakeIntegrityError(err)
	}
	if err := assertStocktakeLedgers(tx, lines); err != nil {
		return nil, stocktakeIntegrityError(err)
	}
	missing, err := stocktakeHasUncountedBalance(tx, lines)
	if err != nil {
		return nil, err
	}
	if missing {
		return auth.NewError(auth.CodeConflict), nil
	}
	conflict, err := refreshStocktakeLines(tx, lines)
	if err != nil {
		return nil, err
	}
	if conflict {
		return auth.NewError(auth.CodeConflict), nil
	}
	return nil, nil
}

func requireStocktakeReasons(lines []gen.StoreStocktakeLine) error {
	for _, line := range lines {
		if *line.CountedQuantity != line.SnapshotQuantity && (line.ReasonCode == nil || *line.ReasonCode == "") {
			return auth.NewError(auth.CodeValidationFailed)
		}
	}
	return nil
}

func stocktakeIntegrityError(err error) error {
	if auth.ErrorCode(err) == auth.CodeConflict {
		return auth.NewError(auth.CodeStocktakeReconciliationRequired)
	}
	return err
}

func assertStocktakeTargets(tx *gorm.DB, storeID string, lines []gen.StoreStocktakeLine) error {
	batches, packages, err := stocktakeLineReferences(tx, lines)
	if err != nil {
		return err
	}
	listingIDs := make([]string, 0, len(batches))
	for _, batch := range batches {
		listingIDs = append(listingIDs, batch.ListingID)
	}
	var listings []gen.StoreListing
	if err := tx.Where("id IN ? AND store_id = ?", listingIDs, storeID).Find(&listings).Error; err != nil {
		return err
	}
	byID := make(map[string]gen.StoreListing, len(listings))
	for _, listing := range listings {
		byID[listing.ID] = listing
	}
	for _, line := range lines {
		if !validStocktakeLineTarget(line, batches, packages, byID) {
			return auth.NewError(auth.CodeConflict)
		}
	}
	return nil
}

func validStocktakeLineTarget(line gen.StoreStocktakeLine, batches map[string]gen.StoreInventoryBatch,
	packages map[string]gen.ProductPackage, listings map[string]gen.StoreListing) bool {
	batch, batchExists := batches[line.BatchID]
	pack, packageExists := packages[line.PackageID]
	listing, listingExists := listings[batch.ListingID]
	return batchExists && packageExists && listingExists && pack.SkuID == listing.SkuID &&
		pack.PackageSetVersion == line.PackageSetVersion
}

func assertStocktakeLedgers(tx *gorm.DB, lines []gen.StoreStocktakeLine) error {
	checked := map[string]bool{}
	for _, line := range lines {
		if checked[line.BatchID] {
			continue
		}
		if err := assertBatchBalanced(tx, line.BatchID); err != nil {
			return err
		}
		checked[line.BatchID] = true
	}
	return nil
}

func writeStocktakeDifferences(tx *gorm.DB, storeID string, lines []gen.StoreStocktakeLine) error {
	for _, line := range lines {
		delta := *line.CountedQuantity - line.SnapshotQuantity
		if delta == 0 {
			continue
		}
		if err := applyAdjustment(tx, line.BatchID, line.PackageID, delta); err != nil {
			return err
		}
		movement := gen.StoreStockMovement{ID: uuid.Must(uuid.NewV4()).String(), Kind: gen.StockMovementKindCountAdjustment,
			StoreID: storeID, BatchID: line.BatchID, StocktakeLineID: &line.ID,
			PackageSetVersion: line.PackageSetVersion, RequestKey: "stocktake:" + line.StocktakeID + ":" + line.ID,
			ReasonCode: line.ReasonCode, OccurredAt: time.Now()}
		if delta > 0 {
			movement.TargetPackageID, movement.TargetQuantity = &line.PackageID, &delta
		} else {
			quantity := -delta
			movement.SourcePackageID, movement.SourceQuantity = &line.PackageID, &quantity
		}
		if err := tx.Create(&movement).Error; err != nil {
			return err
		}
	}
	return nil
}
