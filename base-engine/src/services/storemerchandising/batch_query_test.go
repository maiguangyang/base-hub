package storemerchandising

import (
	"context"
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
	"strings"
	"testing"
	"time"
)

func TestStockListFiltersApplyBeforePagination(t *testing.T) {
	service, db, principal := fixture(t)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	first, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "EARLY", Quantity: 2, RequestKey: "filter-early"})
	merchandisingNoError(t, err)
	second, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "LATER", Quantity: 2, RequestKey: "filter-later"})
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Model(&gen.StoreInventoryBatch{}).Where("id = ?", first.BatchID).Update("expires_at", time.Now().Add(-time.Hour)).Error)
	q, sellable := "LATER", true
	page, err := service.BatchesFiltered(ctx, principal, "store", listing.ID, BatchFilter{Q: &q, Sellable: &sellable}, 1, 1)
	merchandisingNoError(t, err)
	assertBatchFilterResult(t, page, second.BatchID)
	sellable = false
	page, err = service.BatchesFiltered(ctx, principal, "store", listing.ID, BatchFilter{Sellable: &sellable}, 1, 1)
	merchandisingNoError(t, err)
	assertBatchFilterResult(t, page, first.BatchID)
	_, err = service.AdjustStock(ctx, principal, StockAdjustment{StoreID: "store", BatchID: second.BatchID,
		PackageID: "piece", Delta: 1, ReasonCode: "COUNT", RequestKey: "filter-adjust"})
	merchandisingNoError(t, err)
	kind := gen.StockMovementKindCountAdjustment
	movements, err := service.MovementsFiltered(ctx, principal, "store", second.BatchID, MovementFilter{Kind: &kind}, 1, 1)
	merchandisingNoError(t, err)
	if movements.Total != 1 || len(movements.Data) != 1 || movements.Data[0].Kind != kind {
		t.Fatalf("filtered movements = %+v", movements)
	}
}

func assertBatchFilterResult(t *testing.T, page *Page[BatchView], wantID string) {
	t.Helper()
	if page.Total != 1 || len(page.Data) != 1 || page.Data[0].Batch.ID != wantID {
		t.Fatalf("filtered batches = %+v, want %s", page, wantID)
	}
}

func TestBatchesFilteredListsStoreWideBatchesAndMatchesNames(t *testing.T) {
	service, db, principal := fixture(t)
	ctx := context.Background()
	first, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Create(&gen.Product{ID: "product-two", Name: "Noodles", OrganizationID: "hq", CategoryID: "category", Enabled: true}).Error)
	merchandisingNoError(t, db.Create(&gen.ProductSku{ID: "sku-two", Name: "Mild", ProductID: "product-two", Enabled: true}).Error)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "piece-two", Name: "Bag", SkuID: "sku-two", PackageSetVersion: 1, Enabled: true}).Error)
	second, err := service.SetListing(ctx, principal, "store", "sku-two", true)
	merchandisingNoError(t, err)
	one, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: first.ID, PackageID: "piece", BatchNumber: "FOOD-1", Quantity: 1, RequestKey: "batch-one"})
	merchandisingNoError(t, err)
	two, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: second.ID, PackageID: "piece-two", BatchNumber: "NOODLE-1", Quantity: 2, RequestKey: "batch-two"})
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Model(&gen.StoreListing{}).Where("id = ?", second.ID).Update("enabled", false).Error)
	merchandisingNoError(t, db.Model(&gen.ProductSku{}).Where("id = ?", "sku-two").Update("enabled", false).Error)

	page, err := service.BatchesFiltered(ctx, principal, "store", "", BatchFilter{}, 1, 20)
	merchandisingNoError(t, err)
	if page.Total != 2 || len(page.Data) != 2 {
		t.Fatalf("store batches = %+v", page)
	}
	if page.Data[0].Batch.ID != one.BatchID && page.Data[1].Batch.ID != one.BatchID {
		t.Fatalf("first listing batch missing: %+v", page.Data)
	}
	if page.Data[0].Batch.ID != two.BatchID && page.Data[1].Batch.ID != two.BatchID {
		t.Fatalf("disabled historical batch missing: %+v", page.Data)
	}
	q := "Noodles"
	filtered, err := service.BatchesFiltered(ctx, principal, "store", "", BatchFilter{Q: &q}, 1, 20)
	merchandisingNoError(t, err)
	if filtered.Total != 1 || filtered.Data[0].Batch.ID != two.BatchID {
		t.Fatalf("name filter = %+v", filtered)
	}
}

func TestBatchesFilteredLoadsBalancesOncePerPage(t *testing.T) {
	service, db, principal := fixture(t)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	for index, number := range []string{"B-1", "B-2", "B-3"} {
		_, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "piece", BatchNumber: number, Quantity: int64(index + 1), RequestKey: "balance-" + number})
		merchandisingNoError(t, err)
	}
	balanceQueries := 0
	name := "test:count-stock-balance-queries"
	merchandisingNoError(t, db.Callback().Query().After("gorm:query").Register(name, func(tx *gorm.DB) {
		if strings.Contains(tx.Statement.SQL.String(), "store_stock_balances") {
			balanceQueries++
		}
	}))
	t.Cleanup(func() { _ = db.Callback().Query().Remove(name) })

	page, err := service.BatchesFiltered(ctx, principal, "store", "", BatchFilter{}, 1, 20)
	merchandisingNoError(t, err)
	if len(page.Data) != 3 || balanceQueries != 1 {
		t.Fatalf("rows=%d balance queries=%d", len(page.Data), balanceQueries)
	}
}

func TestBatchesFilteredRejectsUnassignedStore(t *testing.T) {
	service, _, principal := fixture(t)
	principal.StoreIDs = map[string]struct{}{}
	_, err := service.BatchesFiltered(context.Background(), principal, "store", "", BatchFilter{}, 1, 20)
	if auth.ErrorCode(err) != auth.CodeStoreScopeDenied {
		t.Fatalf("scope error: %v", err)
	}
}
