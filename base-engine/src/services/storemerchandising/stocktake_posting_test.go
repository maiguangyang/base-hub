package storemerchandising

import (
	"context"
	"base-engine/gen"
	"testing"
)

func TestStocktakePostWritesDifferenceOnce(t *testing.T) {
	service, db, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "post-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(context.Background(), principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listing.ID}, RequestKey: "post-count"})
	merchandisingNoError(t, err)
	line := sheet.Lines[0]
	_, err = service.RecordStocktakeLine(context.Background(), principal, "store", sheet.ID, line.ID, 3)
	merchandisingNoError(t, err)
	_, err = service.RecordStocktakeLine(context.Background(), principal, "store", sheet.ID, line.ID, 2)
	merchandisingNoError(t, err)
	review, err := service.SubmitStocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	assertStocktakeCountHistory(t, review.Lines[0].CountHistory, principal.AccountID)
	_, err = service.SetStocktakeReason(context.Background(), principal, "store", sheet.ID, line.ID, "COUNT_OMISSION", "货架后方漏盘")
	merchandisingNoError(t, err)
	posted, err := service.PostStocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	if posted.Status != gen.StocktakeStatusPosted {
		t.Fatalf("status = %s", posted.Status)
	}
	assertStocktakeNote(t, posted.Lines[0].ReasonNote)
	_, err = service.PostStocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	var movements []gen.StoreStockMovement
	merchandisingNoError(t, db.Where("stocktake_line_id = ?", line.ID).Find(&movements).Error)
	if len(movements) != 1 || movements[0].SourceQuantity == nil || *movements[0].SourceQuantity != 1 {
		t.Fatalf("movements = %+v", movements)
	}
	var balance gen.StoreStockBalance
	merchandisingNoError(t, db.Where("batch_id = ? AND package_id = ?", line.BatchID, line.PackageID).First(&balance).Error)
	if balance.Quantity != 2 {
		t.Fatalf("balance = %d, want 2", balance.Quantity)
	}
}

func assertStocktakeCountHistory(t *testing.T, history []StocktakeCountEvent, actor string) {
	t.Helper()
	if len(history) != 2 {
		t.Fatalf("count history = %+v", history)
	}
	if history[0].CountedQuantity != 3 || history[1].CountedQuantity != 2 || history[1].ActorAccountID != actor {
		t.Fatalf("count history = %+v", history)
	}
}

func assertStocktakeNote(t *testing.T, note *string) {
	t.Helper()
	if note == nil || *note != "货架后方漏盘" {
		t.Fatalf("reason note = %+v", note)
	}
}

func TestStocktakePostKeepsDistinctMovementsForTwoPackages(t *testing.T) {
	service, db, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	child, factor := "piece", int64(12)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "bag", Name: "Bag", SkuID: "sku",
		PackageSetVersion: 1, ContainsPackageID: &child, ContainsQuantity: &factor, Enabled: true}).Error)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "two-pack-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(context.Background(), principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listing.ID}, RequestKey: "two-pack-count"})
	merchandisingNoError(t, err)
	if len(sheet.Lines) != 2 {
		t.Fatalf("lines = %d", len(sheet.Lines))
	}
	for _, line := range sheet.Lines {
		count := int64(2)
		if line.PackageID == "bag" {
			count = 1
		}
		_, err = service.RecordStocktakeLine(context.Background(), principal, "store", sheet.ID, line.ID, count)
		merchandisingNoError(t, err)
	}
	_, err = service.SubmitStocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	for _, line := range sheet.Lines {
		_, err = service.SetStocktakeReason(context.Background(), principal, "store", sheet.ID, line.ID, "COUNT_DIFFERENCE")
		merchandisingNoError(t, err)
	}
	_, err = service.PostStocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	var movements []gen.StoreStockMovement
	merchandisingNoError(t, db.Where("request_key LIKE ?", "stocktake:%").Find(&movements).Error)
	if len(movements) != 2 || movements[0].ID == movements[1].ID || movements[0].ID == "" {
		t.Fatalf("distinct stocktake movements missing: %+v", movements)
	}
}
