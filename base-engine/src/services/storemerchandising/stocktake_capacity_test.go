package storemerchandising

import (
	"context"
	"fmt"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestStocktakeAtCapacityCanCountNewBalancePackage(t *testing.T) {
	service, db, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	ctx := context.Background()
	sheet, firstBatchID := seedCapacityStocktake(t, service, db, principal)
	var err error
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "capacity-unused", SkuID: "sku", Name: "Unused",
		PackageSetVersion: 1, Enabled: true}).Error)
	if _, err = service.AddStocktakeLine(ctx, principal, "store", sheet.ID, firstBatchID, "capacity-unused"); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("unbacked line exceeded creation limit: %v", err)
	}
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "capacity-historical", SkuID: "sku", Name: "Historical",
		PackageSetVersion: 1, Enabled: false}).Error)
	_, err = service.AdjustStock(ctx, principal, StockAdjustment{StoreID: "store", BatchID: firstBatchID,
		PackageID: "capacity-historical", Delta: 2, ReasonCode: "CORRECTION", RequestKey: "capacity-new-balance-adjust"})
	merchandisingNoError(t, err)
	if _, err = service.SubmitStocktake(ctx, principal, "store", sheet.ID); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("new balance without a count line: %v", err)
	}
	updated, err := service.AddStocktakeLine(ctx, principal, "store", sheet.ID, firstBatchID, "capacity-historical")
	merchandisingNoError(t, err)
	if len(updated.Lines) != 201 {
		t.Fatalf("lines after balance-backed add = %d, want 201", len(updated.Lines))
	}
	var newLineID string
	for _, line := range updated.Lines {
		if line.PackageID == "capacity-historical" {
			newLineID = line.ID
		}
	}
	if newLineID == "" {
		t.Fatal("new balance package line missing")
	}
	_, err = service.RecordStocktakeLine(ctx, principal, "store", sheet.ID, newLineID, 1)
	merchandisingNoError(t, err)
	_, err = service.SubmitStocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	_, err = service.SetStocktakeReason(ctx, principal, "store", sheet.ID, newLineID, "COUNT_DIFFERENCE")
	merchandisingNoError(t, err)
	posted, err := service.PostStocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	if posted.Status != gen.StocktakeStatusPosted {
		t.Fatalf("status = %s, want POSTED", posted.Status)
	}
	assertCapacityBalance(t, db, firstBatchID)
}

func seedCapacityStocktake(t *testing.T, service *Service, db *gorm.DB, principal *auth.WorkspacePrincipal) (*StocktakeView, string) {
	t.Helper()
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	child, factor := "piece", int64(2)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "capacity-case", SkuID: "sku", Name: "Case",
		ContainsPackageID: &child, ContainsQuantity: &factor, PackageSetVersion: 1, Enabled: true}).Error)
	firstBatchID := ""
	for index := range 100 {
		receipt, receiveErr := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
			PackageID: "piece", BatchNumber: fmt.Sprintf("B-%03d", index), Quantity: 1,
			RequestKey: fmt.Sprintf("capacity-new-balance-receive-%03d", index)})
		merchandisingNoError(t, receiveErr)
		if index == 0 {
			firstBatchID = receipt.BatchID
		}
	}
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listing.ID}, RequestKey: "capacity-new-balance-sheet"})
	merchandisingNoError(t, err)
	if len(sheet.Lines) != 200 {
		t.Fatalf("initial lines = %d, want 200", len(sheet.Lines))
	}
	merchandisingNoError(t, db.Model(&gen.StoreStocktakeLine{}).Where("stocktake_id = ?", sheet.ID).
		Update("counted_quantity", gorm.Expr("snapshot_quantity")).Error)
	return sheet, firstBatchID
}

func assertCapacityBalance(t *testing.T, db *gorm.DB, batchID string) {
	t.Helper()
	var balance gen.StoreStockBalance
	merchandisingNoError(t, db.Where("batch_id = ? AND package_id = ?", batchID, "capacity-historical").First(&balance).Error)
	if balance.Quantity != 1 {
		t.Fatalf("new package balance = %d, want 1", balance.Quantity)
	}
}
