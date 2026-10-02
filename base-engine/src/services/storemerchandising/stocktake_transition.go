package storemerchandising

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

func (s *Service) ReturnStocktake(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, id string) (*StocktakeView, error) {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		store, err := storeScope(tx, principal, storeID, "franchiseStocktake:post", authorization.AccessUpdate)
		if err != nil {
			return err
		}
		sheet, err := lockedStocktake(tx, store.ID, id)
		if err != nil {
			return err
		}
		if sheet.Status != gen.StocktakeStatusReview {
			return auth.NewError(auth.CodeConflict)
		}
		if err := tx.Model(&gen.StoreStocktakeLine{}).Where("stocktake_id = ?", id).
			Updates(map[string]any{"reason_code": nil, "reason_note": nil}).Error; err != nil {
			return err
		}
		if err := tx.Model(sheet).Updates(map[string]any{"status": gen.StocktakeStatusCounting, "reviewed_at": nil}).Error; err != nil {
			return err
		}
		return s.stocktakeAudit(tx, principal, store, id, "franchiseStocktake:post")
	})
	if err != nil {
		return nil, err
	}
	return readStocktake(s.db.WithContext(ctx), storeID, id)
}

func (s *Service) CancelStocktake(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, id string) (*StocktakeView, error) {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		store, err := storeScope(tx, principal, storeID, "franchiseStocktake:post", authorization.AccessUpdate)
		if err != nil {
			return err
		}
		sheet, err := lockedStocktake(tx, store.ID, id)
		if err != nil {
			return err
		}
		if sheet.Status == gen.StocktakeStatusCanceled {
			return nil
		}
		if sheet.Status == gen.StocktakeStatusPosted {
			return auth.NewError(auth.CodeConflict)
		}
		if err := tx.Model(sheet).Updates(map[string]any{"status": gen.StocktakeStatusCanceled, "canceled_at": time.Now()}).Error; err != nil {
			return err
		}
		return s.stocktakeAudit(tx, principal, store, id, "franchiseStocktake:post")
	})
	if err != nil {
		return nil, err
	}
	return readStocktake(s.db.WithContext(ctx), storeID, id)
}

func (s *Service) AddStocktakeLine(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, id, batchID, packageID string) (*StocktakeView, error) {
	if batchID == "" || packageID == "" {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.addStocktakeLineTx(tx, principal, storeID, id, batchID, packageID)
	})
	if err != nil {
		return nil, err
	}
	return readStocktake(s.db.WithContext(ctx), storeID, id)
}

func (s *Service) addStocktakeLineTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, storeID, id, batchID, packageID string) error {
	store, err := storeScope(tx, principal, storeID, "franchiseStocktake:record", authorization.AccessUpdate)
	if err != nil {
		return err
	}
	if _, err := stocktakeInStatus(tx, store.ID, id, gen.StocktakeStatusCounting); err != nil {
		return err
	}
	created, err := createAdditionalStocktakeLine(tx, store.ID, id, batchID, packageID)
	if err != nil || !created {
		return err
	}
	return s.stocktakeAudit(tx, principal, store, id, "franchiseStocktake:record")
}

func createAdditionalStocktakeLine(tx *gorm.DB, storeID, id, batchID, packageID string) (bool, error) {
	var existing gen.StoreStocktakeLine
	err := tx.Where("stocktake_id = ? AND batch_id = ? AND package_id = ?", id, batchID, packageID).First(&existing).Error
	if err == nil {
		return false, nil
	}
	if err != gorm.ErrRecordNotFound {
		return false, err
	}
	var lineCount int64
	if err := tx.Model(&gen.StoreStocktakeLine{}).Where("stocktake_id = ?", id).Count(&lineCount).Error; err != nil {
		return false, err
	}
	line, hasBalance, err := newAdditionalStocktakeLine(tx, storeID, id, batchID, packageID)
	if err != nil {
		return false, err
	}
	if lineCount >= maxStocktakeLines && !hasBalance {
		return false, auth.NewError(auth.CodeValidationFailed)
	}
	return true, tx.Create(line).Error
}

func newAdditionalStocktakeLine(tx *gorm.DB, storeID, id, batchID, packageID string) (*gen.StoreStocktakeLine, bool, error) {
	var batchLine gen.StoreStocktakeLine
	if err := tx.Where("stocktake_id = ? AND batch_id = ?", id, batchID).First(&batchLine).Error; err != nil {
		return nil, false, auth.NewError(auth.CodeValidationFailed)
	}
	batch, listing, err := stockBatchAndListing(tx, storeID, batchID)
	if err != nil {
		return nil, false, err
	}
	pack, err := stockPackageForListing(tx, packageID, listing.SkuID)
	if err != nil {
		return nil, false, err
	}
	balance, err := getBalance(tx, batchID, packageID)
	if err != nil {
		return nil, false, err
	}
	if !pack.Enabled && balance == nil {
		return nil, false, auth.NewError(auth.CodeValidationFailed)
	}
	line := &gen.StoreStocktakeLine{ID: uuid.Must(uuid.NewV4()).String(), StocktakeID: id,
		BatchID: batchID, PackageID: packageID, PackageSetVersion: pack.PackageSetVersion,
		BatchNumberSnapshot: batch.BatchNumber, ExpiresAtSnapshot: batch.ExpiresAt,
		PackageNameSnapshot: pack.Name, PackageEnabledSnapshot: pack.Enabled}
	if balance != nil {
		line.SnapshotQuantity, line.SnapshotVersion = balance.Quantity, balance.Version
	}
	return line, balance != nil, nil
}
