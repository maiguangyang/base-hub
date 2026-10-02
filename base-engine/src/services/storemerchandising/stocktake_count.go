package storemerchandising

import (
	"context"
	"math"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func lockedStocktake(tx *gorm.DB, storeID, id string) (*gen.StoreStocktake, error) {
	var sheet gen.StoreStocktake
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND store_id = ?", id, storeID).First(&sheet).Error; err != nil {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	return &sheet, nil
}

func (s *Service) RecordStocktakeLine(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, id, lineID string, quantity int64) (*StocktakeView, error) {
	if !validStocktakeQuantity(quantity) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		store, err := storeScope(tx, principal, storeID, "franchiseStocktake:record", authorization.AccessUpdate)
		if err != nil {
			return err
		}
		if _, err := stocktakeInStatus(tx, storeID, id, gen.StocktakeStatusCounting); err != nil {
			return err
		}
		line, err := stocktakeLine(tx, id, lineID)
		if err != nil {
			return err
		}
		if line.CountedQuantity != nil && *line.CountedQuantity == quantity && !line.NeedsRecount {
			return nil
		}
		countedAt := time.Now()
		if err := tx.Model(line).Updates(map[string]any{"counted_quantity": quantity,
			"counted_at": countedAt, "needs_recount": false}).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: "franchiseStocktake:record",
			ResourceType: "storeStocktakeLine", ResourceID: line.ID, ResultCode: "SUCCESS",
			Metadata: audit.Metadata{CountedQuantity: &quantity, CountedAt: countedAt.Format(time.RFC3339Nano)}})
	})
	if err != nil {
		return nil, err
	}
	return readStocktake(s.db.WithContext(ctx), storeID, id)
}

func validStocktakeQuantity(quantity int64) bool { return quantity >= 0 && quantity <= math.MaxInt32 }

func stocktakeInStatus(tx *gorm.DB, storeID, id string, status gen.StocktakeStatus) (*gen.StoreStocktake, error) {
	sheet, err := lockedStocktake(tx, storeID, id)
	if err != nil {
		return nil, err
	}
	if sheet.Status != status {
		return nil, auth.NewError(auth.CodeConflict)
	}
	return sheet, nil
}

func stocktakeLine(tx *gorm.DB, id, lineID string) (*gen.StoreStocktakeLine, error) {
	var line gen.StoreStocktakeLine
	if err := tx.Where("id = ? AND stocktake_id = ?", lineID, id).First(&line).Error; err != nil {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	return &line, nil
}

func (s *Service) SubmitStocktake(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, id string) (*StocktakeView, error) {
	var outcome error
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		store, err := storeScope(tx, principal, storeID, "franchiseStocktake:record", authorization.AccessUpdate)
		if err != nil {
			return err
		}
		outcome, err = s.submitStocktakeTx(tx, principal, store, id)
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

func (s *Service) submitStocktakeTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, store *gen.Store, id string) (error, error) {
	sheet, err := stocktakeInStatus(tx, store.ID, id, gen.StocktakeStatusCounting)
	if err != nil {
		return nil, err
	}
	lines, err := stocktakeLines(tx, id)
	if err != nil {
		return nil, err
	}
	if err := requireCountedStocktakeLines(lines); err != nil {
		return nil, err
	}
	missing, err := stocktakeHasUncountedBalance(tx, lines)
	if err != nil {
		return nil, err
	}
	if missing {
		return auth.NewError(auth.CodeConflict), nil
	}
	changed, err := refreshStocktakeLines(tx, lines)
	if err != nil {
		return nil, err
	}
	if changed {
		return auth.NewError(auth.CodeConflict), nil
	}
	if err := tx.Model(sheet).Updates(map[string]any{"status": gen.StocktakeStatusReview, "reviewed_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	return nil, s.stocktakeAudit(tx, principal, store, id, "franchiseStocktake:record")
}

func stocktakeLines(tx *gorm.DB, id string) ([]gen.StoreStocktakeLine, error) {
	var lines []gen.StoreStocktakeLine
	if err := tx.Where("stocktake_id = ?", id).Order("batch_id, package_id").Find(&lines).Error; err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	return lines, nil
}

func requireCountedStocktakeLines(lines []gen.StoreStocktakeLine) error {
	for _, line := range lines {
		if line.CountedQuantity == nil {
			return auth.NewError(auth.CodeValidationFailed)
		}
	}
	return nil
}

// Lock the selected batches' balance ranges before advancing the sheet. A new
// package balance must become a count line even when the package is disabled.
func stocktakeHasUncountedBalance(tx *gorm.DB, lines []gen.StoreStocktakeLine) (bool, error) {
	lineTargets := make(map[string]map[string]bool)
	batchIDs := make([]string, 0)
	for _, line := range lines {
		if lineTargets[line.BatchID] == nil {
			lineTargets[line.BatchID] = make(map[string]bool)
			batchIDs = append(batchIDs, line.BatchID)
		}
		lineTargets[line.BatchID][line.PackageID] = true
	}
	var balances []gen.StoreStockBalance
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("batch_id IN ?", batchIDs).Order("batch_id, package_id").Find(&balances).Error; err != nil {
		return false, err
	}
	for _, balance := range balances {
		if !lineTargets[balance.BatchID][balance.PackageID] {
			return true, nil
		}
	}
	return false, nil
}

func refreshStocktakeLines(tx *gorm.DB, lines []gen.StoreStocktakeLine) (bool, error) {
	changed := false
	for _, line := range lines {
		stale, err := refreshStocktakeLine(tx, &line)
		if err != nil {
			return false, err
		}
		changed = changed || stale
	}
	return changed, nil
}

func refreshStocktakeLine(tx *gorm.DB, line *gen.StoreStocktakeLine) (bool, error) {
	balance, err := getBalance(tx, line.BatchID, line.PackageID)
	if err != nil {
		return false, err
	}
	var quantity, version int64
	if balance != nil {
		quantity, version = balance.Quantity, balance.Version
	}
	if version == line.SnapshotVersion {
		return false, nil
	}
	return true, tx.Model(line).Updates(map[string]any{"snapshot_quantity": quantity,
		"snapshot_version": version, "counted_quantity": nil, "counted_at": nil,
		"reason_code": nil, "reason_note": nil, "needs_recount": true}).Error
}

func (s *Service) stocktakeAudit(tx *gorm.DB, principal *auth.WorkspacePrincipal, store *gen.Store, id, action string) error {
	return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: action,
		ResourceType: "storeStocktake", ResourceID: id, ResultCode: "SUCCESS"})
}
