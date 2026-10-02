package storemerchandising

import (
	"context"
	"testing"

	"base-engine/gen"
)

func TestStockBatchQueryMySQL(t *testing.T) {
	service, db, principal := mysqlStockFixture(t)
	ctx := context.Background()
	first, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Create(&gen.Product{ID: "batch-product-two", Name: "Noodles", OrganizationID: "hq", CategoryID: "category", Enabled: true}).Error)
	merchandisingNoError(t, db.Create(&gen.ProductSku{ID: "batch-sku-two", Name: "Mild", ProductID: "batch-product-two", Enabled: true}).Error)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "batch-piece-two", Name: "Bag", SkuID: "batch-sku-two", PackageSetVersion: 1, Enabled: true}).Error)
	second, err := service.SetListing(ctx, principal, "store", "batch-sku-two", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: first.ID, PackageID: "piece", BatchNumber: "FOOD-1", Quantity: 1, RequestKey: "batch-one"})
	merchandisingNoError(t, err)
	two, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: second.ID, PackageID: "batch-piece-two", BatchNumber: "NOODLE-1", Quantity: 2, RequestKey: "batch-two"})
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Model(&gen.StoreListing{}).Where("id = ?", second.ID).Update("enabled", false).Error)
	merchandisingNoError(t, db.Model(&gen.ProductSku{}).Where("id = ?", "batch-sku-two").Update("enabled", false).Error)
	page, err := service.BatchesFiltered(ctx, principal, "store", "", BatchFilter{}, 1, 20)
	merchandisingNoError(t, err)
	if page.Total != 2 || len(page.Data) != 2 {
		t.Fatalf("store-wide batches: %+v", page)
	}
	q := "Noodles"
	filtered, err := service.BatchesFiltered(ctx, principal, "store", "", BatchFilter{Q: &q}, 1, 20)
	merchandisingNoError(t, err)
	assertBatchFilterResult(t, filtered, two.BatchID)
	var plan []map[string]any
	term := "%Food%"
	err = db.Raw(`EXPLAIN FORMAT=TRADITIONAL SELECT store_inventory_batches.*
 FROM store_inventory_batches
 JOIN store_listings ON store_listings.id = store_inventory_batches.listing_id
 LEFT JOIN product_skus ON product_skus.id = store_listings.sku_id
 LEFT JOIN products ON products.id = product_skus.product_id
 WHERE store_listings.store_id = ? AND (store_inventory_batches.batch_number LIKE ? OR product_skus.name LIKE ? OR products.name LIKE ?)
 ORDER BY store_inventory_batches.expires_at, store_inventory_batches.id LIMIT 20`, "store", term, term, term).Scan(&plan).Error
	merchandisingNoError(t, err)
	if len(plan) == 0 {
		t.Fatal("empty MySQL execution plan")
	}
	t.Logf("batch query plan: %+v", plan)
}
