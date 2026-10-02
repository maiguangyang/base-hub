package storemerchandising

import (
	"context"
	"testing"

	"base-engine/gen"
	"gorm.io/gorm"
)

func TestUnpackLocksLowerPackageBeforeHigherSource(t *testing.T) {
	service, db, principal := fixture(t)
	child := "piece"
	factor := int64(12)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "z-box", Name: "Box", SkuID: "sku", PackageSetVersion: 1,
		ContainsPackageID: &child, ContainsQuantity: &factor, Enabled: true}).Error)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	receipt, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "z-box", Quantity: 1, BatchNumber: "B1", RequestKey: "receive"})
	merchandisingNoError(t, err)

	var lockOrder []string
	merchandisingNoError(t, db.Callback().Query().After("gorm:query").Register("test:unpack-balance-order", func(tx *gorm.DB) {
		if tx.Statement.Table != "store_stock_balances" {
			return
		}
		for _, variable := range tx.Statement.Vars {
			if packageID, ok := variable.(string); ok && (packageID == "piece" || packageID == "z-box") {
				lockOrder = append(lockOrder, packageID)
				return
			}
		}
	}))
	_, err = service.UnpackStock(ctx, principal, StockUnpack{StoreID: "store", BatchID: receipt.BatchID,
		SourcePackageID: "z-box", Quantity: 1, RequestKey: "unpack"})
	merchandisingNoError(t, err)
	if len(lockOrder) < 2 || lockOrder[0] != "piece" || lockOrder[1] != "z-box" {
		t.Fatalf("balance lock order = %v, want piece before z-box", lockOrder)
	}
}
