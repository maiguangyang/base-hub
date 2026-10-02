package storemerchandising

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestStocktakeListFiltersAndHidesCountingBook(t *testing.T) {
	service, _, principal := fixture(t)
	principal.Permissions["franchiseStocktake:read"] = struct{}{}
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "list-seed"})
	merchandisingNoError(t, err)
	created, err := service.CreateStocktake(context.Background(), principal, CreateStocktakeInput{StoreID: "store", ListingIDs: []string{listing.ID}, RequestKey: "list-sheet"})
	merchandisingNoError(t, err)
	assertStocktakeCreatedAtFilter(t, service, principal, created)
	page, err := service.Stocktakes(context.Background(), principal, "store", StocktakeFilter{Status: gen.StocktakeStatusCounting,
		ListingID: listing.ID}, 1, 20)
	merchandisingNoError(t, err)
	if page.Total != 1 || len(page.Data) != 1 || page.Data[0].ID != created.ID || page.Data[0].Lines[0].SnapshotQuantity != nil {
		t.Fatalf("blind list = %+v", page)
	}
	hasDifference := true
	page, err = service.Stocktakes(context.Background(), principal, "store", StocktakeFilter{HasDifference: &hasDifference}, 1, 20)
	merchandisingNoError(t, err)
	if page.Total != 0 {
		t.Fatalf("difference filter exposed counting sheet: %+v", page)
	}
	_, err = service.RecordStocktakeLine(context.Background(), principal, "store", created.ID, created.Lines[0].ID, 2)
	merchandisingNoError(t, err)
	_, err = service.SubmitStocktake(context.Background(), principal, "store", created.ID)
	merchandisingNoError(t, err)
	page, err = service.Stocktakes(context.Background(), principal, "store", StocktakeFilter{HasDifference: &hasDifference}, 1, 20)
	merchandisingNoError(t, err)
	if page.Total != 1 || page.Data[0].ID != created.ID {
		t.Fatalf("difference filter = %+v", page)
	}
	hasDifference = false
	page, err = service.Stocktakes(context.Background(), principal, "store", StocktakeFilter{HasDifference: &hasDifference}, 1, 20)
	merchandisingNoError(t, err)
	if page.Total != 0 {
		t.Fatalf("no-difference filter = %+v", page)
	}
	principal.Permissions = map[string]struct{}{"franchiseStocktake:record": {}}
	if _, err := service.Stocktakes(context.Background(), principal, "store", StocktakeFilter{}, 1, 20); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("list without read permission: %v", err)
	}
}

func TestStocktakeListIncludesCountingAddLineChoices(t *testing.T) {
	service, db, principal := fixture(t)
	principal.Permissions["franchiseStocktake:read"] = struct{}{}
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	receipt, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "list-choice-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store",
		BatchIDs: []string{receipt.BatchID}, RequestKey: "list-choice-sheet"})
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "list-choice-extra", Name: "Extra", SkuID: "sku",
		PackageSetVersion: 1, Enabled: true}).Error)
	page, err := service.Stocktakes(ctx, principal, "store", StocktakeFilter{}, 1, 20)
	merchandisingNoError(t, err)
	detail, err := service.Stocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	if len(page.Data) != 1 || len(detail.AddLineChoices) != 1 || len(page.Data[0].AddLineChoices) != 1 ||
		page.Data[0].AddLineChoices[0] != detail.AddLineChoices[0] {
		t.Fatalf("list choices = %+v, detail choices = %+v", page.Data, detail.AddLineChoices)
	}
}

func TestStocktakeListLoadsMultipleSheetsInBoundedQueries(t *testing.T) {
	service, db, principal := fixture(t)
	principal.Permissions["franchiseStocktake:read"] = struct{}{}
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "bounded-list-seed"})
	merchandisingNoError(t, err)
	countedSheetID := ""
	for _, key := range []string{"bounded-list-a", "bounded-list-b", "bounded-list-c"} {
		created, createErr := service.CreateStocktake(ctx, principal, CreateStocktakeInput{
			StoreID: "store", ListingIDs: []string{listing.ID}, RequestKey: key,
		})
		merchandisingNoError(t, createErr)
		if key == "bounded-list-a" {
			countedSheetID = created.ID
			_, err = service.RecordStocktakeLine(ctx, principal, "store", created.ID, created.Lines[0].ID, 2)
			merchandisingNoError(t, err)
		}
	}
	queryCount := 0
	const callback = "test:stocktake-list-query-count"
	merchandisingNoError(t, db.Callback().Query().Before("gorm:query").Register(callback, func(*gorm.DB) { queryCount++ }))
	t.Cleanup(func() { _ = db.Callback().Query().Remove(callback) })
	page, err := service.Stocktakes(ctx, principal, "store", StocktakeFilter{}, 1, 20)
	merchandisingNoError(t, err)
	if page.Total != 3 || len(page.Data) != 3 {
		t.Fatalf("list = %+v", page)
	}
	assertBoundedStocktakeHistory(t, page.Data, countedSheetID)
	if queryCount > 8 {
		t.Fatalf("list used %d queries for three sheets, want at most 8", queryCount)
	}
	withoutHistory := false
	queryCount = 0
	page, err = service.Stocktakes(ctx, principal, "store", StocktakeFilter{IncludeHistory: &withoutHistory}, 1, 20)
	merchandisingNoError(t, err)
	if queryCount > 7 {
		t.Fatalf("list without history used %d queries, want at most 7", queryCount)
	}
	for _, item := range page.Data {
		if len(item.Lines) != 1 || len(item.Lines[0].CountHistory) != 0 {
			t.Fatalf("list without history retained audit logs: %+v", item)
		}
	}
}

func assertBoundedStocktakeHistory(t *testing.T, items []StocktakeView, countedSheetID string) {
	t.Helper()
	for _, item := range items {
		if len(item.Lines) != 1 || item.Lines[0].SnapshotQuantity != nil {
			t.Fatalf("counting list leaked or lost lines: %+v", item)
		}
		if item.ID == countedSheetID {
			if item.Lines[0].CountedQuantity == nil || *item.Lines[0].CountedQuantity != 2 || len(item.Lines[0].CountHistory) != 1 {
				t.Fatalf("counted line or history lost: %+v", item)
			}
		} else if item.Lines[0].CountedQuantity != nil || len(item.Lines[0].CountHistory) != 0 {
			t.Fatalf("another sheet inherited count history: %+v", item)
		}
	}
}
