package storemerchandising

import (
	"context"
	"testing"
	"time"

	"base-engine/gen"
)

func TestBatchesFilteredKeepsHistoricalMetadataAndBalances(t *testing.T) {
	service, db, principal := fixture(t)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "piece", BatchNumber: "B1", Quantity: 7, RequestKey: "history"})
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Model(&gen.ProductPackage{}).Where("id = ?", "piece").Update("enabled", false).Error)
	page, err := service.BatchesFiltered(ctx, principal, "store", "", BatchFilter{}, 1, 20)
	merchandisingNoError(t, err)
	view := page.Data[0]
	assertHistoricalBatchMetadata(t, view)
	merchandisingNoError(t, db.Delete(&gen.ProductSku{}, "id = ?", "sku").Error)
	page, err = service.BatchesFiltered(ctx, principal, "store", "", BatchFilter{}, 1, 20)
	merchandisingNoError(t, err)
	if len(page.Data) != 1 || page.Data[0].SkuID != "sku" || page.Data[0].SkuName != nil || page.Data[0].ProductName != nil {
		t.Fatalf("missing names must remain nullable: %+v", page)
	}
}

func assertHistoricalBatchMetadata(t *testing.T, view BatchView) {
	t.Helper()
	if view.SkuID != "sku" || view.ProductName == nil || *view.ProductName != "Food" || view.SkuName == nil || *view.SkuName != "Large" {
		t.Fatalf("batch metadata = %+v", view)
	}
	if len(view.Packages) != 1 || view.Packages[0].Enabled || len(view.Balances) != 1 || view.Balances[0].Quantity != 7 {
		t.Fatalf("historical package/balance = %+v", view)
	}
}

func TestBatchesFilteredStablePageOrderAndForeignListing(t *testing.T) {
	service, db, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	expiry := time.Now().Add(time.Hour)
	merchandisingNoError(t, db.Create([]gen.StoreInventoryBatch{{ID: "batch-b", ListingID: listing.ID, BatchNumber: "B2", ExpiresAt: &expiry}, {ID: "batch-a", ListingID: listing.ID, BatchNumber: "B1", ExpiresAt: &expiry}}).Error)
	for index, want := range []string{"batch-a", "batch-b"} {
		page, err := service.BatchesFiltered(context.Background(), principal, "store", "", BatchFilter{}, index+1, 1)
		merchandisingNoError(t, err)
		if page.Total != 2 || len(page.Data) != 1 || page.Data[0].Batch.ID != want {
			t.Fatalf("stable page %d = %+v", index+1, page)
		}
	}
	page, err := service.BatchesFiltered(context.Background(), principal, "store", "foreign-listing", BatchFilter{}, 1, 20)
	merchandisingNoError(t, err)
	if page.Total != 0 || len(page.Data) != 0 {
		t.Fatalf("foreign listing produced rows: %+v", page)
	}
}
