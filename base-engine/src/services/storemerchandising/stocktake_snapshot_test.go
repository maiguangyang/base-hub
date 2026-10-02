package storemerchandising

import (
	"context"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
)

func TestStocktakeKeepsBatchAndPackageLabelsFromCreation(t *testing.T) {
	service, db, principal := fixture(t)
	principal.Permissions["franchiseStocktake:read"] = struct{}{}
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	expiry := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	_, err = service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", ExpiresAt: &expiry, Quantity: 3, RequestKey: "snapshot-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listing.ID}, RequestKey: "snapshot-sheet"})
	merchandisingNoError(t, err)
	original := sheet.Lines[0]
	merchandisingNoError(t, db.Model(&gen.StoreInventoryBatch{}).Where("id = ?", original.BatchID).
		Updates(map[string]any{"batch_number": "B2", "expires_at": expiry.Add(24 * time.Hour)}).Error)
	merchandisingNoError(t, db.Model(&gen.ProductPackage{}).Where("id = ?", original.PackageID).
		Updates(map[string]any{"name": "Renamed", "enabled": false}).Error)
	view, err := service.Stocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	current := view.Lines[0]
	if current.BatchNumber != original.BatchNumber || current.PackageName != original.PackageName ||
		current.PackageEnabled != original.PackageEnabled || current.ExpiresAt == nil || original.ExpiresAt == nil ||
		current.ExpiresAt.Unix() != original.ExpiresAt.Unix() {
		t.Fatalf("historical labels changed: before=%+v after=%+v", original, current)
	}
}

func TestStocktakePostRejectsChangedPackageOwnership(t *testing.T) {
	service, db, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "ownership-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listing.ID}, RequestKey: "ownership-sheet"})
	merchandisingNoError(t, err)
	line := sheet.Lines[0]
	_, err = service.RecordStocktakeLine(ctx, principal, "store", sheet.ID, line.ID, 2)
	merchandisingNoError(t, err)
	_, err = service.SubmitStocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	_, err = service.SetStocktakeReason(ctx, principal, "store", sheet.ID, line.ID, "COUNT_DIFFERENCE")
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Create(&gen.ProductSku{ID: "other-sku", ProductID: "product", Name: "Other", Enabled: true}).Error)
	merchandisingNoError(t, db.Model(&gen.ProductPackage{}).Where("id = ?", line.PackageID).Update("sku_id", "other-sku").Error)
	if _, err = service.PostStocktake(ctx, principal, "store", sheet.ID); auth.ErrorCode(err) != auth.CodeStocktakeReconciliationRequired {
		t.Fatalf("changed package ownership posted: %v", err)
	}
	var movements int64
	merchandisingNoError(t, db.Model(&gen.StoreStockMovement{}).Where("stocktake_line_id IS NOT NULL").Count(&movements).Error)
	if movements != 0 {
		t.Fatalf("rejected post wrote %d movements", movements)
	}
}

func TestStocktakePostConflictKeepsRecountBlind(t *testing.T) {
	service, _, principal := fixture(t)
	principal.Permissions["franchiseStocktake:read"] = struct{}{}
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	receipt, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "blind-conflict-seed"})
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B2", Quantity: 3, RequestKey: "blind-conflict-seed-two"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listing.ID}, RequestKey: "blind-conflict-sheet"})
	merchandisingNoError(t, err)
	for _, line := range sheet.Lines {
		_, err = service.RecordStocktakeLine(ctx, principal, "store", sheet.ID, line.ID, 2)
		merchandisingNoError(t, err)
	}
	_, err = service.SubmitStocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	for _, line := range sheet.Lines {
		_, err = service.SetStocktakeReason(ctx, principal, "store", sheet.ID, line.ID, "COUNT_DIFFERENCE", "prior difference")
		merchandisingNoError(t, err)
	}
	_, err = service.AdjustStock(ctx, principal, StockAdjustment{StoreID: "store", BatchID: receipt.BatchID,
		PackageID: "piece", Delta: 1, ReasonCode: "CORRECTION", RequestKey: "blind-conflict-adjust"})
	merchandisingNoError(t, err)
	if _, err = service.PostStocktake(ctx, principal, "store", sheet.ID); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	view, err := service.Stocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	if view.Status != gen.StocktakeStatusCounting {
		t.Fatalf("status = %s", view.Status)
	}
	for _, line := range view.Lines {
		if line.SnapshotQuantity != nil || line.ReasonCode != nil || line.ReasonNote != nil {
			t.Fatalf("recount leaked prior book or difference: %+v", line)
		}
	}
}
