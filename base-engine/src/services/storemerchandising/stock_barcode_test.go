package storemerchandising

import (
	"context"
	"errors"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/dbup"
	"gorm.io/gorm"
)

func barcodeStockFixture(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal) {
	t.Helper()
	service, db, principal := fixture(t)
	for _, item := range []any{&gen.ProductBrand{}, &gen.SpecificationDefinition{}, &gen.SpecificationValue{}, &gen.ProductSpecificationChoice{}, &gen.ProductSkuSpecificationValue{}} {
		merchandisingNoError(t, db.AutoMigrate(item))
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	merchandisingNoError(t, dbup.EnsureProductIndexes(db))
	merchandisingNoError(t, db.Model(&gen.ProductSku{}).Where("id = ?", "sku").Update("published_package_set_version", 1).Error)
	merchandisingNoError(t, db.Model(&gen.ProductPackage{}).Where("id = ?", "piece").Update("barcode", "00123").Error)
	return service, db, principal
}
func TestStockBarcodeReadyForReadOnlyOperator(t *testing.T) {
	service, db, principal := barcodeStockFixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	child, factor := "piece", int64(5)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "bag", SkuID: "sku", Name: "Bag", ContainsPackageID: &child, ContainsQuantity: &factor, Barcode: barcodeString("bag-code"), PackageSetVersion: 1, Enabled: true}).Error)
	child, factor = "bag", 8
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "box", SkuID: "sku", Name: "Box", ContainsPackageID: &child, ContainsQuantity: &factor, Barcode: barcodeString("0012345678905"), PackageSetVersion: 1, Enabled: true}).Error)
	principal.Permissions = map[string]struct{}{"franchiseStock:read": {}}
	result, err := service.StockPackageByBarcode(context.Background(), principal, "store", " 0012345678905 ")
	merchandisingNoError(t, err)
	if result.Status != gen.FranchiseStockPackageRecognitionStatusReady || result.Target == nil {
		t.Fatalf("recognition: %+v", result)
	}
	target := result.Target
	if target.ListingID != listing.ID || target.PackageID != "box" || len(target.ConversionChain) != 3 || target.ConversionChain[2].BaseQuantity != 40 {
		t.Fatalf("target: %+v", target)
	}
	if _, err := service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "box", Quantity: 2, RequestKey: "read-only"}); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("read role received: %v", err)
	}
}
func TestStockBarcodeStatusesAndScope(t *testing.T) {
	service, db, principal := barcodeStockFixture(t)
	ctx := context.Background()
	assertStockBarcodeStatus(t, service, principal, "missing", gen.FranchiseStockPackageRecognitionStatusNotFound)
	assertStockBarcodeStatus(t, service, principal, "00123", gen.FranchiseStockPackageRecognitionStatusNotSelected)
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Model(listing).Update("enabled", false).Error)
	assertStockBarcodeStatus(t, service, principal, "00123", gen.FranchiseStockPackageRecognitionStatusUnavailable)
	merchandisingNoError(t, db.Model(listing).Update("enabled", true).Error)
	merchandisingNoError(t, db.Model(&gen.ProductSku{}).Where("id = ?", "sku").Update("published_package_set_version", 2).Error)
	assertStockBarcodeStatus(t, service, principal, "00123", gen.FranchiseStockPackageRecognitionStatusUnavailable)
	merchandisingNoError(t, db.Model(&gen.ProductPackage{}).Where("id = ?", "piece").Update("enabled", false).Error)
	assertStockBarcodeStatus(t, service, principal, "00123", gen.FranchiseStockPackageRecognitionStatusNotFound)
	principal.StoreIDs = map[string]struct{}{}
	if result, err := service.StockPackageByBarcode(ctx, principal, "store", "00123"); auth.ErrorCode(err) != auth.CodeStoreScopeDenied || result != nil {
		t.Fatalf("scope leaked: %+v %v", result, err)
	}
}
func assertStockBarcodeStatus(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, code string, want gen.FranchiseStockPackageRecognitionStatus) {
	t.Helper()
	result, err := service.StockPackageByBarcode(context.Background(), principal, "store", code)
	merchandisingNoError(t, err)
	if result.Status != want || result.Target != nil {
		t.Fatalf("recognition: %+v want %s", result, want)
	}
}
func barcodeString(value string) *string { return &value }
func TestReceiptBarcodeFirstWriteAndReplay(t *testing.T) {
	service, db, principal := barcodeStockFixture(t)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	input := StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "piece", Barcode: barcodeString("stale"), Quantity: 2, RequestKey: "barcode-receipt"}
	_, err = service.ReceiveStock(ctx, principal, input)
	if auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("stale accepted: %v", err)
	}
	for _, item := range []any{&gen.StoreInventoryBatch{}, &gen.StoreStockBalance{}, &gen.StoreStockMovement{}} {
		var count int64
		merchandisingNoError(t, db.Model(item).Count(&count).Error)
		if count != 0 {
			t.Fatalf("first-write rejection changed inventory: %d", count)
		}
	}
	input.Barcode = barcodeString("00123")
	first, err := service.ReceiveStock(ctx, principal, input)
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Model(&gen.ProductPackage{}).Where("id = ?", "piece").Updates(map[string]any{"barcode": "replacement", "enabled": false}).Error)
	replay, err := service.ReceiveStock(ctx, principal, input)
	merchandisingNoError(t, err)
	if first.ID != replay.ID {
		t.Fatal("replay created another movement")
	}
	var balance gen.StoreStockBalance
	merchandisingNoError(t, db.Where("batch_id = ?", first.BatchID).First(&balance).Error)
	if balance.Quantity != 2 {
		t.Fatalf("replay doubled stock: %d", balance.Quantity)
	}
	input.Quantity = 3
	if _, err := service.ReceiveStock(ctx, principal, input); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("changed replay accepted: %v", err)
	}
}

func TestStockBarcodeRejectsBrokenChainsAndAmbiguousCode(t *testing.T) {
	service, db, principal := barcodeStockFixture(t)
	_, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	child, factor := "missing", int64(5)
	pack := gen.ProductPackage{ID: "broken", SkuID: "sku", Name: "Box", PackageSetVersion: 1, Enabled: true, ContainsPackageID: &child, ContainsQuantity: &factor, Barcode: barcodeString("broken-code")}
	merchandisingNoError(t, db.Create(&pack).Error)
	assertStockBarcodeStatus(t, service, principal, "broken-code", gen.FranchiseStockPackageRecognitionStatusUnavailable)
	merchandisingNoError(t, db.Model(&pack).Update("contains_package_id", pack.ID).Error)
	assertStockBarcodeStatus(t, service, principal, "broken-code", gen.FranchiseStockPackageRecognitionStatusUnavailable)
	merchandisingNoError(t, db.Migrator().DropIndex(&gen.ProductPackage{}, "uidx_product_package_active_barcode"))
	merchandisingNoError(t, db.Model(&pack).Update("barcode", "00123").Error)
	assertStockBarcodeStatus(t, service, principal, "00123", gen.FranchiseStockPackageRecognitionStatusConflict)
}
func TestStockBarcodeRejectsDisabledProductAndInvalidInput(t *testing.T) {
	service, db, principal := barcodeStockFixture(t)
	merchandisingNoError(t, db.Model(&gen.Product{}).Where("id = ?", "product").Update("enabled", false).Error)
	assertStockBarcodeStatus(t, service, principal, "00123", gen.FranchiseStockPackageRecognitionStatusUnavailable)
	for _, code := range []string{"", "   ", string(make([]byte, 65))} {
		if result, err := service.StockPackageByBarcode(context.Background(), principal, "store", code); auth.ErrorCode(err) != auth.CodeValidationFailed || result != nil {
			t.Fatalf("invalid recognition: %+v %v", result, err)
		}
	}
	principal.Permissions = map[string]struct{}{}
	if result, err := service.StockPackageByBarcode(context.Background(), principal, "store", "00123"); auth.ErrorCode(err) != auth.CodePermissionDenied || result != nil {
		t.Fatalf("permissionless recognition: %+v %v", result, err)
	}
}

func TestStockBarcodePreservesDatabaseErrors(t *testing.T) {
	service, db, principal := barcodeStockFixture(t)
	failure := errors.New("database unavailable")
	merchandisingNoError(t, db.Callback().Query().Before("gorm:query").Register("recognition_database_error", func(tx *gorm.DB) {
		if tx.Statement.Table == "product_skus" {
			tx.AddError(failure)
		}
	}))
	defer db.Callback().Query().Remove("recognition_database_error")
	result, err := service.StockPackageByBarcode(context.Background(), principal, "store", "00123")
	if !errors.Is(err, failure) || result != nil {
		t.Fatalf("database error became recognition: %v %v", result, err)
	}
}
