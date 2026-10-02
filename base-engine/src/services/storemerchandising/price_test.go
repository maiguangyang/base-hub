package storemerchandising

import (
	"context"
	"fmt"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestPriceHistoryPaginatesAndKeepsStoreScope(t *testing.T) {
	service, _, principal := fixture(t)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	if err != nil {
		t.Fatal(err)
	}
	var offer *gen.StorePackageOffer
	for index := 0; index < 25; index++ {
		offer, err = service.SetPrice(ctx, principal, listing.ID, "piece", int64(1000+index), fmt.Sprintf("CHANGE_%d", index))
		if err != nil {
			t.Fatal(err)
		}
	}
	page, err := service.PriceHistory(ctx, principal, "store", offer.ID, 2, 20)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 25 || len(page.Data) != 5 {
		t.Fatalf("history page = %+v", page)
	}
	if _, err := service.PriceHistory(ctx, principal, "foreign-store", offer.ID, 1, 20); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("foreign store history: %v", err)
	}
}

func TestPriceRevisionAndPackageOwnership(t *testing.T) {
	service, db, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	if err != nil {
		t.Fatal(err)
	}
	offer, err := service.SetPrice(context.Background(), principal, listing.ID, "piece", 1200, "INITIAL")
	if err != nil {
		t.Fatal(err)
	}
	if offer.PriceFen != 1200 {
		t.Fatalf("price = %d", offer.PriceFen)
	}
	updated, err := service.SetPrice(context.Background(), principal, listing.ID, "piece", 900, "PROMO")
	if err != nil || updated.ID != offer.ID || updated.PriceFen != 900 {
		t.Fatalf("updated = %+v, %v", updated, err)
	}
	assertPriceRevisions(t, db, offer.ID)
	assertInvalidStorePrices(t, service, principal, listing.ID)
}

func TestSetPriceAllowsEmptyReason(t *testing.T) {
	service, db, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	if err != nil {
		t.Fatal(err)
	}
	offer, err := service.SetPrice(context.Background(), principal, listing.ID, "piece", 1200, "   ")
	if err != nil {
		t.Fatal(err)
	}
	var revision gen.StorePriceRevision
	if err := db.Where("offer_id = ?", offer.ID).First(&revision).Error; err != nil {
		t.Fatal(err)
	}
	if revision.ReasonCode != "" {
		t.Fatalf("reason code = %q", revision.ReasonCode)
	}
}

func assertPriceRevisions(t *testing.T, db *gorm.DB, offerID string) {
	t.Helper()
	var revisions []gen.StorePriceRevision
	if err := db.Where("offer_id = ?", offerID).Order("effective_at").Find(&revisions).Error; err != nil {
		t.Fatal(err)
	}
	if len(revisions) != 2 || revisions[1].PreviousPriceFen == nil || *revisions[1].PreviousPriceFen != 1200 {
		t.Fatalf("revisions = %+v", revisions)
	}
}

func assertInvalidStorePrices(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, listingID string) {
	t.Helper()
	if _, err := service.SetPrice(context.Background(), principal, listingID, "foreign-package", 100, "BAD"); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("foreign package: %v", err)
	}
	if _, err := service.SetPrice(context.Background(), principal, listingID, "piece", -1, "BAD"); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("negative price: %v", err)
	}
}
