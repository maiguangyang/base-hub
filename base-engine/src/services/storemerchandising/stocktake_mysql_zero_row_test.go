package storemerchandising

import (
	"context"
	"sync"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestStocktakeMySQLConcurrentZeroRowInsertion(t *testing.T) {
	service, db, principal := mysqlStockFixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	batchID, sheetID := preparedMySQLZeroRowSheet(t, service, db, principal)
	postErr, adjustErr := raceMySQLZeroRowInsertion(service, principal, batchID, sheetID)
	if adjustErr != nil {
		t.Fatalf("concurrent zero-row adjustment failed: %v", adjustErr)
	}
	if postErr != nil && auth.ErrorCode(postErr) != auth.CodeConflict {
		t.Fatalf("concurrent zero-row posting failed unexpectedly: %v", postErr)
	}
	assertMySQLZeroRowOutcome(t, db, batchID, sheetID, postErr)
}

func preparedMySQLZeroRowSheet(t *testing.T, service *Service, db *gorm.DB, principal *auth.WorkspacePrincipal) (string, string) {
	t.Helper()
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	receipt, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "zero-row-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store", BatchIDs: []string{receipt.BatchID}, RequestKey: "zero-row-sheet"})
	merchandisingNoError(t, err)
	for _, line := range sheet.Lines {
		count := int64(2)
		if line.PackageID == "bag" {
			count = 0
		}
		_, err = service.RecordStocktakeLine(ctx, principal, "store", sheet.ID, line.ID, count)
		merchandisingNoError(t, err)
	}
	_, err = service.SubmitStocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	var absentBalance int64
	merchandisingNoError(t, db.Model(&gen.StoreStockBalance{}).Where("batch_id = ? AND package_id = ?", receipt.BatchID, "bag").Count(&absentBalance).Error)
	if absentBalance != 0 {
		t.Fatalf("bag balance existed before concurrent insertion: %d", absentBalance)
	}
	for _, line := range sheet.Lines {
		if line.PackageID == "piece" {
			_, err = service.SetStocktakeReason(ctx, principal, "store", sheet.ID, line.ID, "COUNT_DIFFERENCE")
			merchandisingNoError(t, err)
		}
	}
	return receipt.BatchID, sheet.ID
}

func raceMySQLZeroRowInsertion(service *Service, principal *auth.WorkspacePrincipal, batchID, sheetID string) (error, error) {
	ctx := context.Background()
	start := make(chan struct{})
	var wg sync.WaitGroup
	var postErr, adjustErr error
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		_, postErr = service.PostStocktake(ctx, principal, "store", sheetID)
	}()
	go func() {
		defer wg.Done()
		<-start
		_, adjustErr = service.AdjustStock(ctx, principal, StockAdjustment{StoreID: "store", BatchID: batchID,
			PackageID: "bag", Delta: 1, ReasonCode: "CORRECTION", RequestKey: "zero-row-insert"})
	}()
	close(start)
	wg.Wait()
	return postErr, adjustErr
}

func assertMySQLZeroRowOutcome(t *testing.T, db *gorm.DB, batchID, sheetID string, postErr error) {
	t.Helper()
	merchandisingNoError(t, assertBatchBalanced(db, batchID))
	var movements int64
	merchandisingNoError(t, db.Model(&gen.StoreStockMovement{}).Where("stocktake_line_id IS NOT NULL").Count(&movements).Error)
	if postErr == nil && movements != 1 || postErr != nil && movements != 0 {
		t.Fatalf("post=%v, stocktake movements=%d", postErr, movements)
	}
	var updated gen.StoreStocktake
	merchandisingNoError(t, db.First(&updated, "id = ?", sheetID).Error)
	if postErr == nil {
		if updated.Status != gen.StocktakeStatusPosted {
			t.Fatalf("successful post status = %s", updated.Status)
		}
		return
	}
	assertMySQLZeroRowRecount(t, db, sheetID, updated)
}

func assertMySQLZeroRowRecount(t *testing.T, db *gorm.DB, sheetID string, updated gen.StoreStocktake) {
	t.Helper()
	if updated.Status != gen.StocktakeStatusCounting {
		t.Fatalf("conflicted post status = %s", updated.Status)
	}
	var bagLine gen.StoreStocktakeLine
	merchandisingNoError(t, db.First(&bagLine, "stocktake_id = ? AND package_id = ?", sheetID, "bag").Error)
	if !bagLine.NeedsRecount || bagLine.CountedQuantity != nil || bagLine.SnapshotQuantity != 1 {
		t.Fatalf("concurrent insert did not require recount: %+v", bagLine)
	}
}
