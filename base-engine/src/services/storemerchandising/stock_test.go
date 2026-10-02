package storemerchandising

import (
	"context"
	"math"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestUnpackRejectsGraphQLQuantityOverflowWithoutChangingStock(t *testing.T) {
	service, db, principal := fixture(t)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	contained, factor := "piece", int64(1_000_000_000)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "large-box", Name: "Large box", SkuID: "sku",
		PackageSetVersion: 1, ContainsPackageID: &contained, ContainsQuantity: &factor, Enabled: true}).Error)
	receipt, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "large-box", Quantity: 3, RequestKey: "large-receipt"})
	merchandisingNoError(t, err)
	if _, err := service.UnpackStock(ctx, principal, StockUnpack{StoreID: "store", BatchID: receipt.BatchID,
		SourcePackageID: "large-box", Quantity: 3, RequestKey: "large-unpack"}); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("expected pre-commit overflow rejection, got %v", err)
	}
	var box gen.StoreStockBalance
	merchandisingNoError(t, db.Where("batch_id = ? AND package_id = ?", receipt.BatchID, "large-box").First(&box).Error)
	if box.Quantity != 3 {
		t.Fatalf("source balance changed: %d", box.Quantity)
	}
	var movements int64
	merchandisingNoError(t, db.Model(&gen.StoreStockMovement{}).Where("request_key = ?", "large-unpack").Count(&movements).Error)
	if movements != 0 {
		t.Fatalf("overflow movement committed: %d", movements)
	}
}

func TestReceiveRejectsBalanceBeyondGraphQLInt(t *testing.T) {
	service, db, principal := fixture(t)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	first, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B-MAX", Quantity: math.MaxInt32, RequestKey: "first-max-receipt"})
	merchandisingNoError(t, err)
	if _, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B-MAX", Quantity: 1, RequestKey: "second-max-receipt"}); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("expected pre-commit balance limit, got %v", err)
	}
	var balance gen.StoreStockBalance
	merchandisingNoError(t, db.Where("batch_id = ? AND package_id = ?", first.BatchID, "piece").First(&balance).Error)
	if balance.Quantity != math.MaxInt32 {
		t.Fatalf("balance changed: %d", balance.Quantity)
	}
}

func TestDraftPackageCannotBeReceivedOrPriced(t *testing.T) {
	service, db, principal := fixture(t)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	draft := gen.ProductPackage{ID: "draft", Name: "Draft", SkuID: "sku", PackageSetVersion: 2, Enabled: false}
	merchandisingNoError(t, db.Select("*").Create(&draft).Error)
	if _, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: draft.ID, Quantity: 1, RequestKey: "draft-receipt"}); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("draft package receipt: %v", err)
	}
	if _, err := service.SetPrice(ctx, principal, listing.ID, draft.ID, 100, "INITIAL"); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("draft package offer: %v", err)
	}
}

func TestStockReceiptRejectsConflictingProductionDateForBatchNumber(t *testing.T) {
	service, _, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	if err != nil {
		t.Fatal(err)
	}
	first := time.Now().Add(-48 * time.Hour).UTC().Truncate(time.Second)
	second := first.Add(24 * time.Hour)
	receipt := StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "piece", BatchNumber: "B1", Quantity: 1, ProducedAt: &first, RequestKey: "receive-first"}
	if _, err := service.ReceiveStock(context.Background(), principal, receipt); err != nil {
		t.Fatal(err)
	}
	receipt.ProducedAt, receipt.RequestKey = &second, "receive-second"
	if _, err := service.ReceiveStock(context.Background(), principal, receipt); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("conflicting production date: %v", err)
	}
}

func TestStockReceiptRejectsConflictingSourceForBatchNumber(t *testing.T) {
	service, db, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	first, second := "delivery-1", "delivery-2"
	receipt := StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "piece", BatchNumber: "B1",
		Quantity: 1, SourceReference: &first, RequestKey: "first-delivery"}
	_, err = service.ReceiveStock(context.Background(), principal, receipt)
	merchandisingNoError(t, err)
	receipt.SourceReference, receipt.RequestKey = &second, "second-delivery"
	if _, err := service.ReceiveStock(context.Background(), principal, receipt); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("conflicting batch source: %v", err)
	}
	var count int64
	merchandisingNoError(t, db.Model(&gen.StoreStockMovement{}).Count(&count).Error)
	if count != 1 {
		t.Fatalf("unexpected receipt movements: %d", count)
	}
}

func TestStockReceiptReplayRejectsChangedBatchTerms(t *testing.T) {
	service, _, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	produced := time.Now().Add(-48 * time.Hour).UTC().Truncate(time.Second)
	expires := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second)
	source := "delivery-1"
	receipt := StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "piece", BatchNumber: "B1",
		Quantity: 1, ProducedAt: &produced, ExpiresAt: &expires, SourceReference: &source, RequestKey: "same-key"}
	_, err = service.ReceiveStock(context.Background(), principal, receipt)
	merchandisingNoError(t, err)
	changed := receipt
	otherProduced := produced.Add(time.Hour)
	changed.ProducedAt = &otherProduced
	if _, err := service.ReceiveStock(context.Background(), principal, changed); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("changed production date replay: %v", err)
	}
	changed = receipt
	otherExpiry := expires.Add(time.Hour)
	changed.ExpiresAt = &otherExpiry
	if _, err := service.ReceiveStock(context.Background(), principal, changed); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("changed expiry replay: %v", err)
	}
	changed = receipt
	otherSource := "delivery-2"
	changed.SourceReference = &otherSource
	if _, err := service.ReceiveStock(context.Background(), principal, changed); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("changed source replay: %v", err)
	}
}

func TestStockReceiptUnpackConservesUnitsAndReplays(t *testing.T) {
	service, db, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	if err != nil {
		t.Fatal(err)
	}
	contained := "piece"
	factor := int64(12)
	bag := gen.ProductPackage{ID: "bag", Name: "Bag", SkuID: "sku", PackageSetVersion: 1, ContainsPackageID: &contained, ContainsQuantity: &factor, Enabled: true}
	if err := db.Create(&bag).Error; err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	receipt := StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "bag", BatchNumber: "B1", Quantity: 2, RequestKey: "receipt-1"}
	movement, err := service.ReceiveStock(ctx, principal, receipt)
	if err != nil {
		t.Fatal(err)
	}
	if movement.Kind != gen.StockMovementKindReceive {
		t.Fatalf("movement = %+v", movement)
	}
	replayed, err := service.ReceiveStock(ctx, principal, receipt)
	if err != nil || replayed.ID != movement.ID {
		t.Fatalf("replay = %+v, %v", replayed, err)
	}
	unpacked, err := service.UnpackStock(ctx, principal, StockUnpack{StoreID: "store", BatchID: movement.BatchID, SourcePackageID: "bag", Quantity: 1, RequestKey: "unpack-1"})
	if err != nil {
		t.Fatal(err)
	}
	if unpacked.TargetQuantity == nil || *unpacked.TargetQuantity != 12 {
		t.Fatalf("unpack = %+v", unpacked)
	}
	assertStockBalances(t, db, movement.BatchID)
	assertUnpackRejectsOverdrawAndForeignStore(t, service, principal, movement.BatchID)
}

func assertStockBalances(t *testing.T, db *gorm.DB, batchID string) {
	t.Helper()
	var balances []gen.StoreStockBalance
	if err := db.Where("batch_id = ?", batchID).Find(&balances).Error; err != nil {
		t.Fatal(err)
	}
	quantities := map[string]int64{}
	for _, balance := range balances {
		quantities[balance.PackageID] = balance.Quantity
	}
	if quantities["bag"] != 1 || quantities["piece"] != 12 {
		t.Fatalf("balances = %+v", quantities)
	}
}

func assertUnpackRejectsOverdrawAndForeignStore(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, batchID string) {
	t.Helper()
	if _, err := service.UnpackStock(context.Background(), principal, StockUnpack{StoreID: "store", BatchID: batchID, SourcePackageID: "bag", Quantity: 2, RequestKey: "unpack-2"}); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("overdraw: %v", err)
	}
	principal.StoreIDs = map[string]struct{}{}
	if _, err := service.UnpackStock(context.Background(), principal, StockUnpack{StoreID: "store", BatchID: batchID, SourcePackageID: "bag", Quantity: 1, RequestKey: "unpack-3"}); auth.ErrorCode(err) != auth.CodeStoreScopeDenied {
		t.Fatalf("scope: %v", err)
	}
}

func TestStockAdjustmentRequiresReasonAndBalance(t *testing.T) {
	service, _, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	if err != nil {
		t.Fatal(err)
	}
	initial, err := service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "piece", Quantity: 5, RequestKey: "receive", BatchNumber: "B"})
	if err != nil {
		t.Fatal(err)
	}
	input := StockAdjustment{StoreID: "store", BatchID: initial.BatchID, PackageID: "piece", Delta: -2, ReasonCode: "DAMAGED", RequestKey: "loss"}
	adjusted, err := service.AdjustStock(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	if adjusted.SourceQuantity == nil || *adjusted.SourceQuantity != 2 {
		t.Fatalf("adjusted = %+v", adjusted)
	}
	if _, err := service.AdjustStock(context.Background(), principal, input); err != nil {
		t.Fatalf("replay: %v", err)
	}
	input.RequestKey = "overdraw"
	input.Delta = -10
	if _, err := service.AdjustStock(context.Background(), principal, input); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("overdraw: %v", err)
	}
	input.RequestKey = "missing-reason"
	input.Delta = 1
	input.ReasonCode = ""
	if _, err := service.AdjustStock(context.Background(), principal, input); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("reason: %v", err)
	}
}

func TestStockWriteBlocksTamperedBalance(t *testing.T) {
	service, db, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	if err != nil {
		t.Fatal(err)
	}
	initial, err := service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "piece", Quantity: 5, RequestKey: "receive", BatchNumber: "B"})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&gen.StoreStockBalance{}).Where("batch_id = ? AND package_id = ?", initial.BatchID, "piece").Update("quantity", 99).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.AdjustStock(context.Background(), principal, StockAdjustment{StoreID: "store", BatchID: initial.BatchID, PackageID: "piece", Delta: -1, ReasonCode: "COUNT", RequestKey: "adjust"}); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("tampered balance: %v", err)
	}
}
