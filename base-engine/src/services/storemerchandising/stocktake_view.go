package storemerchandising

import (
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func stocktakeViews(tx *gorm.DB, sheets []gen.StoreStocktake, includeHistory bool) ([]StocktakeView, error) {
	views := make([]StocktakeView, 0, len(sheets))
	if len(sheets) == 0 {
		return views, nil
	}
	ids := make([]string, 0, len(sheets))
	for _, sheet := range sheets {
		ids = append(ids, sheet.ID)
	}
	var lines []gen.StoreStocktakeLine
	if err := tx.Where("stocktake_id IN ?", ids).Order("stocktake_id, batch_id, package_id").Find(&lines).Error; err != nil {
		return nil, err
	}
	history, batches, packages, err := stocktakeViewRelations(tx, lines, includeHistory)
	if err != nil {
		return nil, err
	}
	bySheet := make(map[string][]gen.StoreStocktakeLine, len(sheets))
	for _, line := range lines {
		bySheet[line.StocktakeID] = append(bySheet[line.StocktakeID], line)
	}
	for _, sheet := range sheets {
		view, err := stocktakeViewFromRows(&sheet, bySheet[sheet.ID], history, batches, packages)
		if err != nil {
			return nil, err
		}
		views = append(views, *view)
	}
	if err := populateStocktakeAddLineChoices(tx, views); err != nil {
		return nil, err
	}
	return views, nil
}

func stocktakeViewRelations(tx *gorm.DB, lines []gen.StoreStocktakeLine, includeHistory bool) (
	map[string][]StocktakeCountEvent, map[string]gen.StoreInventoryBatch, map[string]gen.ProductPackage, error) {
	history := map[string][]StocktakeCountEvent{}
	var err error
	if includeHistory {
		history, err = stocktakeCountHistory(tx, lines)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	if len(lines) == 0 {
		return history, map[string]gen.StoreInventoryBatch{}, map[string]gen.ProductPackage{}, nil
	}
	batches, packages, err := stocktakeLineReferences(tx, lines)
	return history, batches, packages, err
}

func stocktakeViewFromRows(sheet *gen.StoreStocktake, lines []gen.StoreStocktakeLine, history map[string][]StocktakeCountEvent,
	batches map[string]gen.StoreInventoryBatch, packages map[string]gen.ProductPackage) (*StocktakeView, error) {
	view := &StocktakeView{ID: sheet.ID, StoreID: sheet.StoreID, Status: sheet.Status,
		RequestKey: sheet.RequestKey, InitiatedByAccountID: sheet.InitiatedByAccountID,
		PostedByID: sheet.PostedByID, StartedAt: sheet.StartedAt, ReviewedAt: sheet.ReviewedAt,
		PostedAt: sheet.PostedAt, CanceledAt: sheet.CanceledAt, Lines: make([]StocktakeLineView, 0, len(lines))}
	bookVisible := stocktakeBookVisible(sheet)
	for _, line := range lines {
		batch, batchExists := batches[line.BatchID]
		_, packageExists := packages[line.PackageID]
		if !batchExists || !packageExists {
			return nil, auth.NewError(auth.CodeConflict)
		}
		item := StocktakeLineView{ID: line.ID, BatchID: line.BatchID, PackageID: line.PackageID,
			ListingID: batch.ListingID, BatchNumber: line.BatchNumberSnapshot, ExpiresAt: line.ExpiresAtSnapshot,
			PackageName: line.PackageNameSnapshot, PackageEnabled: line.PackageEnabledSnapshot,
			PackageSetVersion: line.PackageSetVersion, CountedQuantity: line.CountedQuantity,
			CountHistory: history[line.ID], NeedsRecount: line.NeedsRecount}
		if bookVisible {
			quantity := line.SnapshotQuantity
			item.SnapshotQuantity = &quantity
			item.ReasonCode, item.ReasonNote = line.ReasonCode, line.ReasonNote
		}
		view.Lines = append(view.Lines, item)
	}
	return view, nil
}

func stocktakeBookVisible(sheet *gen.StoreStocktake) bool {
	if sheet.Status == gen.StocktakeStatusCounting {
		return false
	}
	return sheet.Status != gen.StocktakeStatusCanceled || sheet.ReviewedAt != nil
}

func stocktakeLineReferences(tx *gorm.DB, lines []gen.StoreStocktakeLine) (map[string]gen.StoreInventoryBatch, map[string]gen.ProductPackage, error) {
	batchIDs, packageIDs := make([]string, 0, len(lines)), make([]string, 0, len(lines))
	for _, line := range lines {
		batchIDs = append(batchIDs, line.BatchID)
		packageIDs = append(packageIDs, line.PackageID)
	}
	var batches []gen.StoreInventoryBatch
	if err := tx.Where("id IN ?", batchIDs).Find(&batches).Error; err != nil {
		return nil, nil, err
	}
	var packages []gen.ProductPackage
	if err := tx.Where("id IN ?", packageIDs).Find(&packages).Error; err != nil {
		return nil, nil, err
	}
	batchByID := make(map[string]gen.StoreInventoryBatch, len(batches))
	for _, batch := range batches {
		batchByID[batch.ID] = batch
	}
	packageByID := make(map[string]gen.ProductPackage, len(packages))
	for _, pack := range packages {
		packageByID[pack.ID] = pack
	}
	return batchByID, packageByID, nil
}
