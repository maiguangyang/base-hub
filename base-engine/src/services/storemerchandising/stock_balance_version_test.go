package storemerchandising

import (
	"context"
	"base-engine/gen"
	"gorm.io/gorm"
	"testing"
)

func TestStockBalanceVersionAdvancesAcrossMutations(t *testing.T) {
	service, db, principal := fixture(t)
	ctx := context.Background()
	child, factor := "piece", int64(12)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "bag", Name: "Bag", SkuID: "sku",
		PackageSetVersion: 1, ContainsPackageID: &child, ContainsQuantity: &factor, Enabled: true}).Error)
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	receipt, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "bag", BatchNumber: "VERSION", Quantity: 2, RequestKey: "version-receive"})
	merchandisingNoError(t, err)
	version := stockVersion(t, db, receipt.BatchID, "bag")
	if version != 1 {
		t.Fatalf("receipt version = %d, want 1", version)
	}
	_, err = service.AdjustStock(ctx, principal, StockAdjustment{StoreID: "store", BatchID: receipt.BatchID,
		PackageID: "bag", Delta: 1, ReasonCode: "COUNT", RequestKey: "version-adjust"})
	merchandisingNoError(t, err)
	if got := stockVersion(t, db, receipt.BatchID, "bag"); got != version+1 {
		t.Fatalf("adjust version = %d, want %d", got, version+1)
	}
	_, err = service.UnpackStock(ctx, principal, StockUnpack{StoreID: "store", BatchID: receipt.BatchID,
		SourcePackageID: "bag", Quantity: 1, RequestKey: "version-unpack"})
	merchandisingNoError(t, err)
	if got := stockVersion(t, db, receipt.BatchID, "bag"); got != version+2 {
		t.Fatalf("unpack source version = %d, want %d", got, version+2)
	}
	if got := stockVersion(t, db, receipt.BatchID, "piece"); got != 1 {
		t.Fatalf("unpack target version = %d, want 1", got)
	}
}

func stockVersion(t *testing.T, db *gorm.DB, batchID, packageID string) int64 {
	t.Helper()
	var balance gen.StoreStockBalance
	merchandisingNoError(t, db.Where("batch_id = ? AND package_id = ?", batchID, packageID).First(&balance).Error)
	return balance.Version
}
