package storemerchandising

import (
	"context"
	"base-engine/auth"
	"base-engine/gen"
	"testing"
)

func TestStocktakePermissionsAllowSameAccountToPost(t *testing.T) {
	service, _, principal := fixture(t)
	principal.Permissions["franchiseStocktake:read"] = struct{}{}
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(context.Background(), principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "permission-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(context.Background(), principal, CreateStocktakeInput{StoreID: "store", ListingIDs: []string{listing.ID}, RequestKey: "permission-sheet"})
	merchandisingNoError(t, err)
	principal.Permissions = map[string]struct{}{"franchiseStocktake:read": {}}
	blind, err := service.Stocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	if blind.Lines[0].SnapshotQuantity != nil {
		t.Fatalf("read role saw book count: %+v", blind.Lines[0])
	}
	if _, err := service.RecordStocktakeLine(context.Background(), principal, "store", sheet.ID, sheet.Lines[0].ID, 3); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("read-only record: %v", err)
	}
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	_, err = service.RecordStocktakeLine(context.Background(), principal, "store", sheet.ID, sheet.Lines[0].ID, 3)
	merchandisingNoError(t, err)
	_, err = service.SubmitStocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	if _, err := service.PostStocktake(context.Background(), principal, "store", sheet.ID); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("record-only post: %v", err)
	}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	posted, err := service.PostStocktake(context.Background(), principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	if posted.Status != gen.StocktakeStatusPosted || posted.PostedByID == nil || *posted.PostedByID != principal.AccountID {
		t.Fatalf("same-account posting: %+v", posted)
	}
	principal.StoreIDs = map[string]struct{}{}
	if _, err := service.Stocktake(context.Background(), principal, "store", sheet.ID); auth.ErrorCode(err) != auth.CodeStoreScopeDenied {
		t.Fatalf("cross-store read: %v", err)
	}
}
