package storemerchandising

import (
	"context"
	"fmt"
	"testing"

	"base-engine/gen"
)

func TestStockBatchNumberMySQLSkipsConsecutiveEquivalentNumbersAcrossProducts(t *testing.T) {
	service, db, principal := mysqlStockFixture(t)
	merchandisingNoError(t, db.Exec("ALTER TABLE store_inventory_batches MODIFY COLUMN batch_number VARCHAR(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL").Error)
	first, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	prefix := batchNumberPrefix(t)
	input := StockReceipt{StoreID: "store", ListingID: first.ID, PackageID: "piece", Quantity: 1}
	for index, digits := range []string{"０００１", "0０0２"} {
		input.BatchNumber, input.RequestKey = prefix+digits, fmt.Sprintf("manual-%d", index)
		receiveNumber(t, service, db, principal, input)
	}
	merchandisingNoError(t, db.Create(&gen.ProductSku{ID: "second-sku", Name: "Second", ProductID: "product", Enabled: true}).Error)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "second-pack", Name: "Piece", SkuID: "second-sku", PackageSetVersion: 1, Enabled: true}).Error)
	second, err := service.SetListing(context.Background(), principal, "store", "second-sku", true)
	merchandisingNoError(t, err)
	input.ListingID, input.PackageID, input.BatchNumber = second.ID, "second-pack", ""
	for _, sequence := range []int{3, 4} {
		input.RequestKey = fmt.Sprintf("automatic-%d", sequence)
		batch := receiveNumber(t, service, db, principal, input)
		if want := fmt.Sprintf("%s%04d", prefix, sequence); batch.BatchNumber != want {
			t.Fatalf("number %s want %s", batch.BatchNumber, want)
		}
	}
}
