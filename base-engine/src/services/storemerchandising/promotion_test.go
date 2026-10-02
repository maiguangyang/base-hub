package storemerchandising

import (
	"context"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
)

func TestPromotionShapeAndHeadquartersStackingCeiling(t *testing.T) {
	service, db, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	if err != nil {
		t.Fatal(err)
	}
	offer, err := service.SetPrice(context.Background(), principal, listing.ID, "piece", 1200, "INITIAL")
	if err != nil {
		t.Fatal(err)
	}
	price := int64(1000)
	input := PromotionInput{StoreID: "store", TimeZone: "Asia/Shanghai", Kind: gen.StorePromotionKindItemPrice, StartsAt: time.Now().Add(time.Hour), EndsAt: time.Now().Add(2 * time.Hour), FixedPriceFen: &price,
		Targets: []PromotionTargetInput{{OfferID: offer.ID, RequiredQuantity: 1}}, StackWithHqCoupon: true}
	if _, err := service.SavePromotion(context.Background(), principal, input); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("HQ stacking denied: %v", err)
	}
	policy := gen.CustomerBenefitPolicy{ID: "policy", OrganizationID: "hq", Version: 1, PromotionWithHqCoupon: true}
	if err := db.Create(&policy).Error; err != nil {
		t.Fatal(err)
	}
	created, err := service.SavePromotion(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	if created.Version != 1 || created.Kind != gen.StorePromotionKindItemPrice {
		t.Fatalf("promotion = %+v", created)
	}
	input.ThresholdFen = &price
	if _, err := service.SavePromotion(context.Background(), principal, input); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("irrelevant field: %v", err)
	}
}

func TestPromotionDraftKeepsOldVersionUntilActivation(t *testing.T) {
	service, db, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	offer, err := service.SetPrice(context.Background(), principal, listing.ID, "piece", 1200, "INITIAL")
	merchandisingNoError(t, err)
	price := int64(1000)
	input := PromotionInput{StoreID: "store", TimeZone: "Asia/Shanghai", Kind: gen.StorePromotionKindItemPrice,
		StartsAt: time.Now().Add(time.Hour), EndsAt: time.Now().Add(2 * time.Hour), FixedPriceFen: &price,
		Targets: []PromotionTargetInput{{OfferID: offer.ID, RequiredQuantity: 1}}}
	first, err := service.SavePromotion(context.Background(), principal, input)
	merchandisingNoError(t, err)
	_, err = service.SetPromotionEnabled(context.Background(), principal, first.ID, true)
	merchandisingNoError(t, err)
	input.PromotionID = &first.ID
	second, err := service.SavePromotion(context.Background(), principal, input)
	merchandisingNoError(t, err)
	var before gen.StorePromotion
	if err := db.First(&before, "id = ?", first.ID).Error; err != nil {
		t.Fatal(err)
	}
	if !before.Enabled || second.Enabled {
		t.Fatalf("draft disabled live version: first=%v second=%v", before.Enabled, second.Enabled)
	}
	_, err = service.SetPromotionEnabled(context.Background(), principal, second.ID, true)
	merchandisingNoError(t, err)
	if err := db.First(&before, "id = ?", first.ID).Error; err != nil {
		t.Fatal(err)
	}
	if before.Enabled {
		t.Fatal("old version remained enabled after activation")
	}
}

func TestPromotionActivationChecksCurrentHeadquartersPolicy(t *testing.T) {
	service, db, principal := fixture(t)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	offer, err := service.SetPrice(context.Background(), principal, listing.ID, "piece", 1200, "INITIAL")
	merchandisingNoError(t, err)
	policy := gen.CustomerBenefitPolicy{ID: "policy", OrganizationID: "hq", Version: 1, PromotionWithHqCoupon: true}
	merchandisingNoError(t, db.Create(&policy).Error)
	price := int64(1000)
	input := PromotionInput{StoreID: "store", TimeZone: "Asia/Shanghai", Kind: gen.StorePromotionKindItemPrice,
		StartsAt: time.Now().Add(time.Hour), EndsAt: time.Now().Add(2 * time.Hour), FixedPriceFen: &price,
		Targets: []PromotionTargetInput{{OfferID: offer.ID, RequiredQuantity: 1}}, StackWithHqCoupon: true}
	draft, err := service.SavePromotion(context.Background(), principal, input)
	merchandisingNoError(t, err)
	merchandisingNoError(t, db.Model(&policy).Update("promotion_with_hq_coupon", false).Error)
	if _, err := service.SetPromotionEnabled(context.Background(), principal, draft.ID, true); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("old draft exceeded current policy: %v", err)
	}
}
