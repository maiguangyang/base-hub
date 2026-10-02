package storemerchandising

import (
	"context"
	"sync"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestStocktakeMySQLConcurrentPostAndReverseOrderUnpack(t *testing.T) {
	service, db, principal := mysqlStockFixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	batchID, sheetID := preparedMySQLReverseOrderSheet(t, service, db, principal)
	ctx := context.Background()
	start := make(chan struct{})
	var wg sync.WaitGroup
	var postErr, unpackErr error
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		_, postErr = service.PostStocktake(ctx, principal, "store", sheetID)
	}()
	go func() {
		defer wg.Done()
		<-start
		_, unpackErr = service.UnpackStock(ctx, principal, StockUnpack{StoreID: "store", BatchID: batchID,
			SourcePackageID: "z-box", Quantity: 1, RequestKey: "race-unpack"})
	}()
	close(start)
	wg.Wait()
	if unpackErr != nil {
		t.Fatalf("unpack failed: %v", unpackErr)
	}
	if postErr != nil && auth.ErrorCode(postErr) != auth.CodeConflict {
		t.Fatalf("post failed unexpectedly: %v", postErr)
	}
	merchandisingNoError(t, assertBatchBalanced(db, batchID))
}

func preparedMySQLReverseOrderSheet(t *testing.T, service *Service, db *gorm.DB, principal *auth.WorkspacePrincipal) (string, string) {
	t.Helper()
	child := "piece"
	factor := int64(12)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "z-box", Name: "Box", SkuID: "sku", PackageSetVersion: 1,
		ContainsPackageID: &child, ContainsQuantity: &factor, Enabled: true}).Error)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	receipt, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "z-box", BatchNumber: "B1", Quantity: 2, RequestKey: "race-box"})
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "race-piece"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store", BatchIDs: []string{receipt.BatchID}, RequestKey: "race-sheet"})
	merchandisingNoError(t, err)
	for _, line := range sheet.Lines {
		count := int64(0)
		switch line.PackageID {
		case "piece":
			count = 3
		case "z-box":
			count = 1
		}
		_, err = service.RecordStocktakeLine(ctx, principal, "store", sheet.ID, line.ID, count)
		merchandisingNoError(t, err)
	}
	_, err = service.SubmitStocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	for _, line := range sheet.Lines {
		if line.PackageID == "z-box" {
			_, err = service.SetStocktakeReason(ctx, principal, "store", sheet.ID, line.ID, "COUNT_DIFFERENCE")
			merchandisingNoError(t, err)
		}
	}
	return receipt.BatchID, sheet.ID
}
