package storemerchandising

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestStockBatchNumberMySQLConcurrentReceipts(t *testing.T) {
	for _, replay := range []bool{false, true} {
		t.Run(fmt.Sprintf("replay-%v", replay), func(t *testing.T) {
			service, db, principal := mysqlStockFixture(t)
			first, err := service.SetListing(context.Background(), principal, "store", "sku", true)
			merchandisingNoError(t, err)
			merchandisingNoError(t, db.Create(&gen.ProductSku{ID: "second-sku", Name: "Second", ProductID: "product", Enabled: true}).Error)
			merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "second-pack", Name: "Piece", SkuID: "second-sku", PackageSetVersion: 1, Enabled: true}).Error)
			second, err := service.SetListing(context.Background(), principal, "store", "second-sku", true)
			merchandisingNoError(t, err)
			inputs := make([]StockReceipt, 12)
			for index := range inputs {
				inputs[index] = StockReceipt{StoreID: "store", ListingID: first.ID, PackageID: "piece", Quantity: 1, RequestKey: fmt.Sprintf("receipt-%d", index)}
				if !replay && index%2 == 1 {
					inputs[index].ListingID, inputs[index].PackageID = second.ID, "second-pack"
				}
				if replay {
					inputs[index].RequestKey = "same-receipt"
				}
			}
			runNumberedConcurrent(t, service, principal, inputs)
			want := 12
			if replay {
				want = 1
			}
			assertNumberedBatchCount(t, db, want, batchNumberPrefix(t))
		})
	}
}

func runNumberedConcurrent(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, inputs []StockReceipt) {
	t.Helper()
	start, results := make(chan struct{}), make(chan error, len(inputs))
	var workers sync.WaitGroup
	for _, input := range inputs {
		workers.Add(1)
		go func(input StockReceipt) {
			defer workers.Done()
			<-start
			_, err := service.ReceiveStock(context.Background(), principal, input)
			results <- err
		}(input)
	}
	close(start)
	workers.Wait()
	close(results)
	for err := range results {
		merchandisingNoError(t, err)
	}
}

func assertNumberedBatchCount(t *testing.T, db *gorm.DB, want int, prefix string) {
	t.Helper()
	var batches []gen.StoreInventoryBatch
	merchandisingNoError(t, db.Order("batch_number").Find(&batches).Error)
	if len(batches) != want {
		t.Fatalf("batch count: %d want %d", len(batches), want)
	}
	for index, batch := range batches {
		if batch.BatchNumber != fmt.Sprintf("%s%04d", prefix, index+1) {
			t.Fatalf("sequence: %s", batch.BatchNumber)
		}
	}
	var movements int64
	merchandisingNoError(t, db.Model(&gen.StoreStockMovement{}).Count(&movements).Error)
	if movements != int64(want) {
		t.Fatalf("movement count: %d want %d", movements, want)
	}
}

func TestStockBatchNumberMySQLReservesCollationEquivalentManualNumber(t *testing.T) {
	for _, test := range []struct {
		initial      string
		digits       string
		changedTerms bool
	}{{"b", "0001", false}, {"b", "0001", true}, {"Ḃ", "0001", false}, {"Ḃ", "0001", true},
		{"B", "０００１", false}, {"B", "０００１", true}, {"B", "0０0１", false}, {"B", "0０0１", true}} {
		t.Run(fmt.Sprintf("prefix-%s/digits-%s/changed-terms-%v", test.initial, test.digits, test.changedTerms), func(t *testing.T) {
			service, db, principal := mysqlStockFixture(t)
			merchandisingNoError(t, db.Exec("ALTER TABLE store_inventory_batches MODIFY COLUMN batch_number VARCHAR(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL").Error)
			listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
			merchandisingNoError(t, err)
			prefix := batchNumberPrefix(t)
			input := StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "piece", Quantity: 1,
				BatchNumber: test.initial + prefix[1:] + test.digits, RequestKey: "manual"}
			manual := receiveNumber(t, service, db, principal, input)
			if manual.BatchNumber != input.BatchNumber {
				t.Fatalf("manual number changed: %s", manual.BatchNumber)
			}
			input.BatchNumber, input.RequestKey = "", "automatic"
			if test.changedTerms {
				source := "new-delivery"
				input.SourceReference = &source
			}
			automatic := receiveNumber(t, service, db, principal, input)
			if automatic.ID == manual.ID || automatic.BatchNumber != prefix+"0002" {
				t.Fatalf("automatic receipt reused manual batch: %+v", automatic)
			}
			if replay := receiveNumber(t, service, db, principal, input); replay.ID != automatic.ID {
				t.Fatalf("replay created batch: %s", replay.ID)
			}
			var balances []gen.StoreStockBalance
			merchandisingNoError(t, db.Find(&balances).Error)
			if len(balances) != 2 {
				t.Fatalf("expected separate batch balances, got %d", len(balances))
			}
			for _, balance := range balances {
				if balance.Quantity != 1 {
					t.Fatalf("batch stock incorrectly combined: %+v", balance)
				}
			}
		})
	}
}
