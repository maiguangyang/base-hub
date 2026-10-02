package storemerchandising

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func batchNumberPrefix(t *testing.T) string {
	t.Helper()
	zone, err := time.LoadLocation("Asia/Shanghai")
	merchandisingNoError(t, err)
	return "B" + time.Now().In(zone).Format("20060102") + "-"
}

func receiveNumber(t *testing.T, service *Service, db *gorm.DB, principal *auth.WorkspacePrincipal, input StockReceipt) gen.StoreInventoryBatch {
	t.Helper()
	movement, err := service.ReceiveStock(context.Background(), principal, input)
	merchandisingNoError(t, err)
	var batch gen.StoreInventoryBatch
	merchandisingNoError(t, db.First(&batch, "id = ?", movement.BatchID).Error)
	return batch
}

func numberedReceipt(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal, StockReceipt) {
	t.Helper()
	service, db, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	return service, db, principal, StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "piece", Quantity: 1, RequestKey: "receipt-first"}
}

func TestStockBatchNumberDailySequenceAndUUIDIdentity(t *testing.T) {
	service, db, principal, input := numberedReceipt(t)
	prefix := batchNumberPrefix(t)
	first := receiveNumber(t, service, db, principal, input)
	input.RequestKey = "receipt-second"
	second := receiveNumber(t, service, db, principal, input)
	if first.BatchNumber != prefix+"0001" || second.BatchNumber != prefix+"0002" {
		t.Fatalf("numbers: %s, %s", first.BatchNumber, second.BatchNumber)
	}
	if _, err := uuid.FromString(first.ID); err != nil {
		t.Fatalf("internal ID no longer UUID: %s", first.ID)
	}
}

func TestStockBatchNumberAcrossProductsAndStores(t *testing.T) {
	service, db, principal, input := numberedReceipt(t)
	prefix := batchNumberPrefix(t)
	receiveNumber(t, service, db, principal, input)
	merchandisingNoError(t, db.Create(&gen.ProductSku{ID: "sku-two", Name: "Other", ProductID: "product", Enabled: true}).Error)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "pack-two", Name: "Piece", SkuID: "sku-two", Enabled: true, PackageSetVersion: 1}).Error)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku-two", true)
	merchandisingNoError(t, err)
	input.ListingID, input.PackageID, input.RequestKey = listing.ID, "pack-two", "other-product"
	if got := receiveNumber(t, service, db, principal, input).BatchNumber; got != prefix+"0002" {
		t.Fatal(got)
	}
	merchandisingNoError(t, db.Create(&gen.Store{ID: "store-two", Code: "S2", Name: "Second", OrganizationID: "org", Lifecycle: gen.StoreLifecycleActive}).Error)
	principal.StoreIDs["store-two"] = struct{}{}
	listing, err = service.SetListing(context.Background(), principal, "store-two", "sku-two", true)
	merchandisingNoError(t, err)
	input.StoreID, input.ListingID, input.RequestKey = "store-two", listing.ID, "other-store"
	if got := receiveNumber(t, service, db, principal, input).BatchNumber; got != prefix+"0001" {
		t.Fatal(got)
	}
}

func TestStockBatchNumberPreservesManualAndSkipsReservedSequence(t *testing.T) {
	service, db, principal, input := numberedReceipt(t)
	prefix := batchNumberPrefix(t)
	for index, number := range []string{"supplier-2026", prefix + "0003", prefix + "9999", "B20000101-9999"} {
		input.BatchNumber, input.RequestKey = number, fmt.Sprintf("manual-%d", index)
		if got := receiveNumber(t, service, db, principal, input).BatchNumber; got != number {
			t.Fatal(got)
		}
	}
	input.BatchNumber, input.RequestKey = "", "automatic-after-manual"
	if got := receiveNumber(t, service, db, principal, input).BatchNumber; got != prefix+"10000" {
		t.Fatal(got)
	}
}

func TestStockBatchNumberReplayAndRejectedReceiptDoNotConsumeSequence(t *testing.T) {
	service, db, principal, input := numberedReceipt(t)
	prefix := batchNumberPrefix(t)
	first := receiveNumber(t, service, db, principal, input)
	if replay := receiveNumber(t, service, db, principal, input); replay.ID != first.ID {
		t.Fatalf("replay created batch: %s", replay.ID)
	}
	input.PackageID, input.RequestKey = "invalid-package", "invalid"
	if _, err := service.ReceiveStock(context.Background(), principal, input); err == nil {
		t.Fatal("invalid package accepted")
	}
	input.PackageID, input.RequestKey = "piece", "new-receipt"
	if got := receiveNumber(t, service, db, principal, input).BatchNumber; got != prefix+"0002" {
		t.Fatal(got)
	}
}

func TestStockBatchNumberResetsAtShanghaiMidnight(t *testing.T) {
	_, db, _, input := numberedReceipt(t)
	merchandisingNoError(t, db.Create(&gen.StoreInventoryBatch{ID: "previous-day", ListingID: input.ListingID, BatchNumber: "B20260929-9999"}).Error)
	for _, test := range []struct {
		hour, minute int
		want         string
	}{{15, 59, "B20260929-10000"}, {16, 0, "B20260930-0001"}} {
		now := time.Date(2026, 9, 29, test.hour, test.minute, 0, 0, time.UTC)
		number, err := nextStockBatchNumber(db, "store", now)
		merchandisingNoError(t, err)
		if number != test.want {
			t.Fatalf("%v: %s want %s", now, number, test.want)
		}
	}
}

func TestStockBatchNumberTransactionRollbackDoesNotConsumeSequence(t *testing.T) {
	service, db, principal, input := numberedReceipt(t)
	callback := "fail_numbered_receipt_audit"
	merchandisingNoError(t, db.Callback().Create().Before("gorm:create").Register(callback, func(tx *gorm.DB) {
		if tx.Statement.Schema != nil && tx.Statement.Schema.Table == "audit_logs" {
			tx.AddError(errors.New("audit failure"))
		}
	}))
	_, err := service.ReceiveStock(context.Background(), principal, input)
	merchandisingNoError(t, db.Callback().Create().Remove(callback))
	if err == nil {
		t.Fatal("audit failure committed receipt")
	}
	var batches, movements int64
	merchandisingNoError(t, db.Model(&gen.StoreInventoryBatch{}).Count(&batches).Error)
	merchandisingNoError(t, db.Model(&gen.StoreStockMovement{}).Count(&movements).Error)
	if batches != 0 || movements != 0 {
		t.Fatalf("rollback left %d batches and %d movements", batches, movements)
	}
	if got := receiveNumber(t, service, db, principal, input).BatchNumber; got != batchNumberPrefix(t)+"0001" {
		t.Fatal(got)
	}
}
