package storemerchandising

import (
	"context"
	"strings"
	"testing"

	"base-engine/auth"
)

func TestStocktakeReasonAcceptsUnicodeCharacterLimit(t *testing.T) {
	service, _, principal := fixture(t)
	principal.Permissions["franchiseStocktake:record"] = struct{}{}
	principal.Permissions["franchiseStocktake:post"] = struct{}{}
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	_, err = service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID,
		PackageID: "piece", BatchNumber: "B1", Quantity: 3, RequestKey: "reason-seed"})
	merchandisingNoError(t, err)
	sheet, err := service.CreateStocktake(ctx, principal, CreateStocktakeInput{StoreID: "store",
		ListingIDs: []string{listing.ID}, RequestKey: "reason-sheet"})
	merchandisingNoError(t, err)
	line := sheet.Lines[0]
	_, err = service.RecordStocktakeLine(ctx, principal, "store", sheet.ID, line.ID, 2)
	merchandisingNoError(t, err)
	_, err = service.SubmitStocktake(ctx, principal, "store", sheet.ID)
	merchandisingNoError(t, err)
	note := strings.Repeat("复", 200)
	view, err := service.SetStocktakeReason(ctx, principal, "store", sheet.ID, line.ID, "COUNT_OMISSION", note)
	merchandisingNoError(t, err)
	if view.Lines[0].ReasonNote == nil || *view.Lines[0].ReasonNote != note {
		t.Fatalf("stored Unicode note = %+v", view.Lines[0].ReasonNote)
	}
	_, err = service.SetStocktakeReason(ctx, principal, "store", sheet.ID, line.ID, "COUNT_OMISSION", strings.Repeat("复", 257))
	if auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("overlong note accepted: %v", err)
	}
}
