package storemerchandising

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestStocktakeReplayRejectsChangedScope(t *testing.T) {
	service, _, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "replay-seed"})
	merchandisingNoError(t, err)
	input := CreateStocktakeInput{StoreID: "store", ListingIDs: []string{listing.ID}, RequestKey: "replay-sheet"}
	first, err := service.CreateStocktake(context.Background(), principal, input)
	merchandisingNoError(t, err)
	second, err := service.CreateStocktake(context.Background(), principal, input)
	merchandisingNoError(t, err)
	if first.ID != second.ID {
		t.Fatalf("replay created two sheets: %s, %s", first.ID, second.ID)
	}
	input.ListingIDs = nil
	input.BatchIDs = []string{"missing-batch"}
	if _, err := service.CreateStocktake(context.Background(), principal, input); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("changed replay scope: %v", err)
	}
}

func TestStocktakeRejectsMixedUnauthorizedScope(t *testing.T) {
	service, _, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "mixed-seed"})
	merchandisingNoError(t, err)
	_, err = service.CreateStocktake(context.Background(), principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listing.ID, "other-store-listing"}, RequestKey: "mixed-sheet"})
	if auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("mixed scope: %v", err)
	}
}

func TestStocktakeCreateStartsBlindCount(t *testing.T) {
	service, _, principal := fixture(t)
	principal.Permissions["franchiseStocktake:read"] = struct{}{}
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "stocktake-seed"})
	merchandisingNoError(t, err)
	item, err := service.CreateStocktake(context.Background(), principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listing.ID}, RequestKey: "count-1"})
	merchandisingNoError(t, err)
	if item.Status != gen.StocktakeStatusCounting || len(item.Lines) != 1 {
		t.Fatalf("unexpected sheet: %+v", item)
	}
	if item.Lines[0].CountedQuantity != nil || item.Lines[0].SnapshotQuantity != nil {
		t.Fatalf("counting response leaked book quantity: %+v", item.Lines[0])
	}
}

func TestStocktakeSubmitDistinguishesUncountedFromZero(t *testing.T) {
	service, _, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "zero-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(context.Background(), principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listing.ID}, RequestKey: "zero-count"})
	merchandisingNoError(t, err)
	if _, err := service.SubmitStocktake(context.Background(), principal, "store", sheet.ID); err == nil {
		t.Fatal("uncounted line submitted")
	}
	_, err = service.RecordStocktakeLine(context.Background(), principal, "store", sheet.ID, sheet.Lines[0].ID, 0)
	merchandisingNoError(t, err)
	review, err := service.SubmitStocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	if review.Status != gen.StocktakeStatusReview || review.Lines[0].CountedQuantity == nil ||
		*review.Lines[0].CountedQuantity != 0 || review.Lines[0].SnapshotQuantity == nil || *review.Lines[0].SnapshotQuantity != 3 {
		t.Fatalf("review state: %+v", review)
	}
}

func TestStocktakeReturnAndCancelTransitions(t *testing.T) {
	service, _, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "return-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(context.Background(), principal, CreateStocktakeInput{StoreID: "store", ListingIDs: []string{listing.ID}, RequestKey: "return-sheet"})
	merchandisingNoError(t, err)
	line := sheet.Lines[0]
	_, err = service.RecordStocktakeLine(context.Background(), principal, "store", sheet.ID, line.ID, 2)
	merchandisingNoError(t, err)
	_, err = service.SubmitStocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	returned, err := service.ReturnStocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	if returned.Status != gen.StocktakeStatusCounting || returned.Lines[0].SnapshotQuantity != nil {
		t.Fatalf("returned = %+v", returned)
	}
	canceled, err := service.CancelStocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	if canceled.Status != gen.StocktakeStatusCanceled || canceled.Lines[0].SnapshotQuantity != nil {
		t.Fatalf("canceled = %+v", canceled)
	}
	if _, err = service.RecordStocktakeLine(context.Background(), principal, "store", sheet.ID, line.ID, 1); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("canceled sheet editable: %v", err)
	}
}

func TestStocktakeAddLineRequiresSameBatchAndPermission(t *testing.T) {
	service, db, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	movement, err := service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "add-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(context.Background(), principal, CreateStocktakeInput{StoreID: "store", BatchIDs: []string{movement.BatchID}, RequestKey: "add-sheet"})
	merchandisingNoError(t, err)
	pack := gen.ProductPackage{ID: "later-pack", Name: "Later", SkuID: "sku", PackageSetVersion: 1, Enabled: true}
	merchandisingNoError(t, db.Create(&pack).Error)
	updated, err := service.AddStocktakeLine(context.Background(), principal, "store", sheet.ID, movement.BatchID, pack.ID)
	merchandisingNoError(t, err)
	if len(updated.Lines) != 2 {
		t.Fatalf("lines = %+v", updated.Lines)
	}
	principal.Permissions = map[string]struct{}{"franchiseStocktake:read": {}}
	if _, err := service.AddStocktakeLine(context.Background(), principal, "store", sheet.ID, movement.BatchID, "unknown"); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("read-only add: %v", err)
	}
}

func TestStocktakeSubmitRejectsNewPackageBalanceOutsideSheet(t *testing.T) {
	service, db, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	receipt, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "new-package-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store",
		BatchIDs: []string{receipt.BatchID}, RequestKey: "new-package-sheet"})
	merchandisingNoError(t, err)
	_, err = service.RecordStocktakeLine(ctx, principal, "store", sheet.ID, sheet.Lines[0].ID, 3)
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "historical", Name: "Historical", SkuID: "sku",
		PackageSetVersion: 1, Enabled: false}).Error)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "unused-historical", Name: "Unused", SkuID: "sku",
		PackageSetVersion: 1, Enabled: false}).Error)
	_, err = service.AdjustStock(ctx, principal, StockAdjustment{StoreID: "store", BatchID: receipt.BatchID,
		PackageID: "historical", Delta: 2, ReasonCode: "CORRECTION", RequestKey: "new-package-adjust"})
	merchandisingNoError(t, err)
	if _, err = service.SubmitStocktake(ctx, principal, "store", sheet.ID); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("new package balance passed submit: %v", err)
	}
	var saved gen.StoreStocktake
	merchandisingNoError(t, db.First(&saved, "id = ?", sheet.ID).Error)
	if saved.Status != gen.StocktakeStatusCounting {
		t.Fatalf("new package balance advanced sheet: %s", saved.Status)
	}
	principal.Permissions["franchiseStocktake:read"] = struct{}{}
	current, err := service.Stocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	if len(current.AddLineChoices) != 1 || current.AddLineChoices[0].PackageID != "historical" ||
		current.AddLineChoices[0].PackageEnabled {
		t.Fatalf("unusable add-line choices: %+v", current.AddLineChoices)
	}
	if current.Lines[0].SnapshotQuantity != nil {
		t.Fatalf("candidate query leaked book quantity: %+v", current.Lines[0])
	}
}

func TestStocktakePostReturnsToCountingForNewPackageBalance(t *testing.T) {
	service, db, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	receipt, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "new-post-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store",
		BatchIDs: []string{receipt.BatchID}, RequestKey: "new-post-sheet"})
	merchandisingNoError(t, err)
	_, err = service.RecordStocktakeLine(ctx, principal, "store", sheet.ID, sheet.Lines[0].ID, 3)
	merchandisingNoError(t, err)
	_, err = service.SubmitStocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "post-historical", Name: "Historical", SkuID: "sku",
		PackageSetVersion: 1, Enabled: false}).Error)
	_, err = service.AdjustStock(ctx, principal, StockAdjustment{StoreID: "store", BatchID: receipt.BatchID,
		PackageID: "post-historical", Delta: 2, ReasonCode: "CORRECTION", RequestKey: "new-post-adjust"})
	merchandisingNoError(t, err)
	if _, err = service.PostStocktake(ctx, principal, "store", sheet.ID); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("new package balance passed post: %v", err)
	}
	var saved gen.StoreStocktake
	merchandisingNoError(t, db.First(&saved, "id = ?", sheet.ID).Error)
	if saved.Status != gen.StocktakeStatusCounting {
		t.Fatalf("new package balance did not return to counting: %s", saved.Status)
	}
	var movements []gen.StoreStockMovement
	merchandisingNoError(t, db.Where("stocktake_line_id IS NOT NULL").Find(&movements).Error)
	if len(movements) != 0 {
		t.Fatalf("partial stocktake posting: %+v", movements)
	}
}
