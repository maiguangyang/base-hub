package productcatalog

import (
	"context"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestRetiredPackageDisablesOffersAndPromotionsButKeepsHistory(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
	catalogNoError(t, err)
	product, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID})
	catalogNoError(t, err)
	sku, err := service.CreateSku(ctx, principal, SkuInput{Name: "Spicy", ProductID: product.ID})
	catalogNoError(t, err)
	base, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "Piece"})
	catalogNoError(t, err)
	pack, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "Box", ContainsPackageID: &base.ID, ContainsQuantity: 12})
	catalogNoError(t, err)
	if _, err := service.RetirePackage(ctx, principal, base.ID); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("retire active child: %v", err)
	}
	seedRetirementOffer(t, db, sku.ID, pack.ID)
	_, err = service.RetirePackage(ctx, principal, pack.ID)
	catalogNoError(t, err)
	var offer gen.StorePackageOffer
	if err := db.First(&offer, "id = ?", "offer").Error; err != nil {
		t.Fatal(err)
	}
	if offer.Enabled {
		t.Fatal("retired package offer remains enabled")
	}
	var promotion gen.StorePromotion
	if err := db.First(&promotion, "id = ?", "promotion").Error; err != nil {
		t.Fatal(err)
	}
	if promotion.Enabled {
		t.Fatal("promotion for retired package remains enabled")
	}
	_, err = service.RetirePackage(ctx, principal, base.ID)
	catalogNoError(t, err)
}

func TestPackageSetSwitchStagesThenPublishesAtomically(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
	catalogNoError(t, err)
	product, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID})
	catalogNoError(t, err)
	sku, err := service.CreateSku(ctx, principal, SkuInput{Name: "Spicy", ProductID: product.ID})
	catalogNoError(t, err)
	old, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "Old piece", Barcode: "SHARED-CODE"})
	catalogNoError(t, err)
	seedRetirementOffer(t, db, sku.ID, old.ID)
	base, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "New piece", Barcode: "SHARED-CODE", PackageSetVersion: 2})
	catalogNoError(t, err)
	box, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "New box", PackageSetVersion: 2,
		ContainsPackageID: &base.ID, ContainsQuantity: 12})
	catalogNoError(t, err)
	if base.Enabled || box.Enabled {
		t.Fatalf("draft package set became active before publish: %+v %+v", base, box)
	}
	if _, err := service.PublishPackageSet(ctx, principal, sku.ID, 2); err != nil {
		t.Fatal(err)
	}
	var packages []gen.ProductPackage
	catalogNoError(t, db.Where("sku_id = ?", sku.ID).Find(&packages).Error)
	for _, pack := range packages {
		if pack.Enabled != (pack.PackageSetVersion == 2) {
			t.Fatalf("mixed package versions: %+v", packages)
		}
	}
	var offer gen.StorePackageOffer
	catalogNoError(t, db.First(&offer, "id = ?", "offer").Error)
	var promotion gen.StorePromotion
	catalogNoError(t, db.First(&promotion, "id = ?", "promotion").Error)
	if offer.Enabled || promotion.Enabled {
		t.Fatalf("old offers still active: offer=%v promotion=%v", offer.Enabled, promotion.Enabled)
	}
	var published gen.ProductSku
	catalogNoError(t, db.First(&published, "id = ?", sku.ID).Error)
	if published.PublishedPackageSetVersion != 2 {
		t.Fatalf("published version = %d", published.PublishedPackageSetVersion)
	}
	if _, err := service.PublishPackageSet(ctx, principal, sku.ID, 1); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("old version republished: %v", err)
	}
}

func TestDraftPackageCanBeEditedAndRebuiltBeforePublication(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
	catalogNoError(t, err)
	product, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID})
	catalogNoError(t, err)
	sku, err := service.CreateSku(ctx, principal, SkuInput{Name: "Spicy", ProductID: product.ID})
	catalogNoError(t, err)
	_, err = service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "Old piece"})
	catalogNoError(t, err)
	base, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "Draft piece", PackageSetVersion: 2})
	catalogNoError(t, err)
	box, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "Draft box", PackageSetVersion: 2,
		ContainsPackageID: &base.ID, ContainsQuantity: 6})
	catalogNoError(t, err)
	updated, err := service.UpdatePackageMetadata(ctx, principal, base.ID, "Correct piece", nil)
	catalogNoError(t, err)
	if updated.Name != "Correct piece" {
		t.Fatalf("draft metadata not updated: %+v", updated)
	}
	if _, err := service.RetirePackage(ctx, principal, base.ID); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("deleted referenced draft base: %v", err)
	}
	_, err = service.RetirePackage(ctx, principal, box.ID)
	catalogNoError(t, err)
	_, err = service.RetirePackage(ctx, principal, base.ID)
	catalogNoError(t, err)
	var remaining int64
	catalogNoError(t, db.Model(&gen.ProductPackage{}).Where("package_set_version = ? AND sku_id = ?", 2, sku.ID).Count(&remaining).Error)
	if remaining != 0 {
		t.Fatalf("draft packages remain: %d", remaining)
	}
	_, err = service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "Replacement piece", PackageSetVersion: 2})
	catalogNoError(t, err)
}

func seedRetirementOffer(t *testing.T, db *gorm.DB, skuID, packageID string) {
	t.Helper()
	for _, item := range []any{
		&gen.StoreListing{ID: "listing", StoreID: "store", SkuID: skuID, Enabled: true},
		&gen.StorePackageOffer{ID: "offer", ListingID: "listing", PackageID: packageID, PriceFen: 100, Enabled: true},
		&gen.StorePromotion{ID: "promotion", StoreID: "store", RuleKey: "rule", Version: 1, Kind: gen.StorePromotionKindItemPrice, TimeZone: "Asia/Seoul", Enabled: true, StartsAt: time.Now(), EndsAt: time.Now().Add(time.Hour)},
		&gen.StorePromotionTarget{ID: "target", PromotionID: "promotion", OfferID: "offer", RequiredQuantity: 1},
	} {
		if err := db.Create(item).Error; err != nil {
			t.Fatal(err)
		}
	}
}

func TestDisablingSkuOrProductStopsOffersAndPromotions(t *testing.T) {
	for _, disableProduct := range []bool{false, true} {
		t.Run(map[bool]string{false: "sku", true: "product"}[disableProduct], func(t *testing.T) {
			service, db, principal := catalogFixture(t)
			ctx := context.Background()
			category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
			catalogNoError(t, err)
			product, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID})
			catalogNoError(t, err)
			sku, err := service.CreateSku(ctx, principal, SkuInput{Name: "Spicy", ProductID: product.ID})
			catalogNoError(t, err)
			pack, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "Piece"})
			catalogNoError(t, err)
			seedRetirementOffer(t, db, sku.ID, pack.ID)
			if disableProduct {
				_, err = service.SetProductEnabled(ctx, principal, product.ID, false)
			} else {
				_, err = service.SetSkuEnabled(ctx, principal, sku.ID, false)
			}
			catalogNoError(t, err)
			var offer gen.StorePackageOffer
			if err := db.First(&offer, "id = ?", "offer").Error; err != nil {
				t.Fatal(err)
			}
			var promotion gen.StorePromotion
			if err := db.First(&promotion, "id = ?", "promotion").Error; err != nil {
				t.Fatal(err)
			}
			if offer.Enabled || promotion.Enabled {
				t.Fatalf("disabled catalog remains sellable: offer=%v promotion=%v", offer.Enabled, promotion.Enabled)
			}
		})
	}
}
