package storemerchandising

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestStocktakePostReportsLedgerMismatchSeparatelyFromRecount(t *testing.T) {
	service, db, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "mismatch-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(context.Background(), principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listing.ID}, RequestKey: "mismatch-sheet"})
	merchandisingNoError(t, err)
	line := sheet.Lines[0]
	_, err = service.RecordStocktakeLine(context.Background(), principal, "store", sheet.ID, line.ID, 3)
	merchandisingNoError(t, err)
	_, err = service.SubmitStocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Model(&gen.StoreStockBalance{}).Where("batch_id = ? AND package_id = ?", line.BatchID, line.PackageID).
		Update("quantity", 4).Error)
	_, err = service.PostStocktake(context.Background(), principal, "store", sheet.ID)
	if code := auth.ErrorCode(err); code != auth.CodeStocktakeReconciliationRequired {
		t.Fatalf("ledger mismatch code = %s, want reconciliation required", code)
	}
}
