package storemerchandising

import (
	"context"
	"sync"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestStocktakeMySQLABAPostConflict(t *testing.T) {
	service, db, principal := mysqlStockFixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	receipt, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "aba-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store", BatchIDs: []string{receipt.BatchID}, RequestKey: "aba-sheet"})
	merchandisingNoError(t, err)
	for _, line := range sheet.Lines {
		count := int64(0)
		if line.PackageID == "piece" {
			count = 2
		}
		_, err = service.RecordStocktakeLine(ctx, principal, "store", sheet.ID, line.ID, count)
		merchandisingNoError(t, err)
	}
	_, err = service.SubmitStocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	for _, line := range sheet.Lines {
		if line.PackageID != "piece" {
			continue
		}
		_, err = service.SetStocktakeReason(ctx, principal, "store", sheet.ID, line.ID, "COUNT_DIFFERENCE")
		merchandisingNoError(t, err)
	}
	stocktakeABAAdjustments(t, service, principal, receipt.BatchID)
	if _, err := service.PostStocktake(ctx, principal, "store", sheet.ID); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("ABA change posted: %v", err)
	}
	var movements int64
	merchandisingNoError(t, db.Model(&gen.StoreStockMovement{}).Where("stocktake_line_id IS NOT NULL").Count(&movements).Error)
	if movements != 0 {
		t.Fatalf("partial posting wrote %d movements", movements)
	}
	var updated gen.StoreStocktake
	merchandisingNoError(t, db.First(&updated, "id = ?", sheet.ID).Error)
	if updated.Status != gen.StocktakeStatusCounting {
		t.Fatalf("status = %s", updated.Status)
	}
}

func TestStocktakeMySQLConcurrentPostAndAdjustment(t *testing.T) {
	service, db, principal := mysqlStockFixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	ctx := context.Background()
	batchID, sheetID := preparedMySQLReviewSheet(t, service, principal)
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
		_, adjustErr = service.AdjustStock(ctx, principal, StockAdjustment{StoreID: "store",
			BatchID: batchID, PackageID: "piece", Delta: 1, ReasonCode: "CORRECTION", RequestKey: "race-adjust"})
	}()
	close(start)
	wg.Wait()
	if adjustErr != nil {
		t.Fatalf("normal stock adjustment failed: %v", adjustErr)
	}
	if postErr != nil && auth.ErrorCode(postErr) != auth.CodeConflict {
		t.Fatalf("post failed unexpectedly: %v", postErr)
	}
	merchandisingNoError(t, assertBatchBalanced(db, batchID))
	var movements int64
	merchandisingNoError(t, db.Model(&gen.StoreStockMovement{}).Where("stocktake_line_id IS NOT NULL").Count(&movements).Error)
	if postErr == nil && movements != 1 || postErr != nil && movements != 0 {
		t.Fatalf("post=%v, movements=%d", postErr, movements)
	}
}

func preparedMySQLReviewSheet(t *testing.T, service *Service, principal *auth.WorkspacePrincipal) (string, string) {
	t.Helper()
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	receipt, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "race-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store", BatchIDs: []string{receipt.BatchID}, RequestKey: "race-sheet"})
	merchandisingNoError(t, err)
	for _, line := range sheet.Lines {
		count := int64(0)
		if line.PackageID == "piece" {
			count = 2
		}
		_, err = service.RecordStocktakeLine(ctx, principal, "store", sheet.ID, line.ID, count)
		merchandisingNoError(t, err)
	}
	_, err = service.SubmitStocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	for _, line := range sheet.Lines {
		if line.PackageID != "piece" {
			continue
		}
		_, err = service.SetStocktakeReason(ctx, principal, "store", sheet.ID, line.ID, "COUNT_DIFFERENCE")
		merchandisingNoError(t, err)
	}
	return receipt.BatchID, sheet.ID
}

func stocktakeABAAdjustments(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, batchID string) {
	t.Helper()
	for _, delta := range []int64{1, -1} {
		key := "aba-plus"
		if delta < 0 {
			key = "aba-minus"
		}
		_, err := service.AdjustStock(context.Background(), principal, StockAdjustment{StoreID: "store", BatchID: batchID,
			PackageID: "piece", Delta: delta, ReasonCode: "CORRECTION", RequestKey: key})
		merchandisingNoError(t, err)
	}
}

func TestStocktakeMySQLConcurrentCreateReplaysOneSheet(t *testing.T) {
	service, db, principal := mysqlStockFixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "create-race-seed"})
	merchandisingNoError(t, err)
	input := CreateStocktakeInput{StoreID: "store", ListingIDs: []string{listing.ID}, RequestKey: "create-race-sheet"}
	start := make(chan struct{})
	results := make(chan *StocktakeView, 2)
	errors := make(chan error, 2)
	for range 2 {
		go func() {
			<-start
			view, err := service.CreateStocktake(ctx, principal, input)
			results <- view
			errors <- err
		}()
	}
	close(start)
	first, second := <-results, <-results
	firstErr, secondErr := <-errors, <-errors
	if firstErr != nil || secondErr != nil {
		t.Fatalf("concurrent replay errors: %v, %v", firstErr, secondErr)
	}
	if first.ID != second.ID {
		t.Fatalf("duplicate sheet IDs: %s, %s", first.ID, second.ID)
	}
	var count int64
	merchandisingNoError(t, db.Model(&gen.StoreStocktake{}).Count(&count).Error)
	if count != 1 {
		t.Fatalf("sheet count = %d", count)
	}
}
