package storemerchandising

import (
	"context"
	"testing"
)

func TestOfferAndPriceHistoryFiltersApplyBeforePagination(t *testing.T) {
	service, _, principal := fixture(t)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	offer, err := service.SetPrice(ctx, principal, listing.ID, "piece", 1000, "INITIAL")
	merchandisingNoError(t, err)
	_, err = service.SetPrice(ctx, principal, listing.ID, "piece", 1200, "PROMOTION")
	merchandisingNoError(t, err)
	packageID, enabled := "piece", true
	offers, err := service.OffersFiltered(ctx, principal, "store", listing.ID, OfferFilter{PackageID: &packageID, Enabled: &enabled}, 1, 1)
	merchandisingNoError(t, err)
	if offers.Total != 1 || len(offers.Data) != 1 || offers.Data[0].ID != offer.ID {
		t.Fatalf("filtered offers = %+v", offers)
	}
	q := "PROMOTION"
	history, err := service.PriceHistoryFiltered(ctx, principal, "store", offer.ID, PriceHistoryFilter{Q: &q}, 1, 1)
	merchandisingNoError(t, err)
	if history.Total != 1 || len(history.Data) != 1 || history.Data[0].ReasonCode != q {
		t.Fatalf("filtered price history = %+v", history)
	}
}
