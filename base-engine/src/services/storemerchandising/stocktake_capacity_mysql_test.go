package storemerchandising

import (
	"context"
	"fmt"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestStocktakeMySQLTwoHundredLinePosting(t *testing.T) {
	service, db, principal := mysqlStockFixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	for index := range 100 {
		_, err = service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
			PackageID: "piece", BatchNumber: fmt.Sprintf("B-%03d", index), Quantity: 1,
			RequestKey: fmt.Sprintf("capacity-receive-%03d", index)})
		merchandisingNoError(t, err)
	}
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listing.ID}, RequestKey: "capacity-sheet"})
	merchandisingNoError(t, err)
	if len(sheet.Lines) != 200 {
		t.Fatalf("lines = %d, want 200", len(sheet.Lines))
	}
	child, factor := "piece", int64(2)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "capacity-extra-pack", SkuID: "sku", Name: "Extra",
		ContainsPackageID: &child, ContainsQuantity: &factor, PackageSetVersion: 1, Enabled: true}).Error)
	_, err = service.AddStocktakeLine(ctx, principal, "store", sheet.ID, sheet.Lines[0].BatchID, "capacity-extra-pack")
	if auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("line added past 200-line limit: %v", err)
	}
	merchandisingNoError(t, db.Model(&gen.StoreStocktakeLine{}).Where("stocktake_id = ?", sheet.ID).
		Update("counted_quantity", 0).Error)
	_, err = service.SubmitStocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Model(&gen.StoreStocktakeLine{}).Where("stocktake_id = ? AND package_id = ?", sheet.ID, "piece").
		Update("reason_code", "COUNT_DIFFERENCE").Error)
	started := time.Now()
	_, err = service.PostStocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	t.Logf("200-line posting elapsed: %s", time.Since(started))
	var movements int64
	merchandisingNoError(t, db.Model(&gen.StoreStockMovement{}).Where("stocktake_line_id IS NOT NULL").Count(&movements).Error)
	if movements != 100 {
		t.Fatalf("count adjustments = %d, want 100", movements)
	}
	for _, line := range sheet.Lines {
		merchandisingNoError(t, assertBatchBalanced(db, line.BatchID))
	}
	assertStocktakeOverLimit(t, service, db, principal, listing.ID)
}

func assertStocktakeOverLimit(t *testing.T, service *Service, db *gorm.DB, principal *auth.WorkspacePrincipal, listingID string) {
	t.Helper()
	ctx := context.Background()
	_, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listingID,
		PackageID: "piece", BatchNumber: "B-100", Quantity: 1, RequestKey: "capacity-receive-100"})
	merchandisingNoError(t, err)
	_, err = service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listingID}, RequestKey: "capacity-too-large"})
	if auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("oversized stocktake accepted: %v", err)
	}
	var oversized int64
	merchandisingNoError(t, db.Model(&gen.StoreStocktake{}).Where("request_key = ?", "capacity-too-large").Count(&oversized).Error)
	if oversized != 0 {
		t.Fatal("oversized stocktake committed a partial sheet")
	}
}
