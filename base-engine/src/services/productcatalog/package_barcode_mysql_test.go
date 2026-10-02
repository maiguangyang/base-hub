package productcatalog

import (
	"context"
	"sync"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/dbup"
	"gorm.io/gorm"
)

func barcodeMySQLFixture(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal) {
	t.Helper()
	service, db, principal := categoryMySQLFixture(t)
	for _, model := range []any{&gen.ProductBrand{}, &gen.Product{}, &gen.ProductSku{}, &gen.ProductPackage{},
		&gen.SpecificationDefinition{}, &gen.SpecificationValue{}, &gen.ProductSpecificationChoice{}, &gen.ProductSkuSpecificationValue{},
		&gen.StoreListing{}, &gen.StorePackageOffer{}, &gen.StorePriceRevision{}, &gen.StorePromotion{}, &gen.StorePromotionTarget{},
		&gen.StoreInventoryBatch{}, &gen.StoreStockBalance{}, &gen.StoreStockMovement{}, &gen.StoreStocktake{}, &gen.StoreStocktakeLine{}} {
		catalogNoError(t, db.AutoMigrate(model))
	}
	catalogNoError(t, dbup.EnsureProductIndexes(db))
	return service, db, principal
}

func TestPackageBarcodeMySQLConcurrentPublish(t *testing.T) {
	service, db, principal := barcodeMySQLFixture(t)
	ctx := context.Background()
	left, right := seedBarcodePublishCompetition(t, service, db, principal)
	start := make(chan struct{})
	results := make(chan struct {
		id  string
		err error
	}, 2)
	var wg sync.WaitGroup
	for _, id := range []string{left.ID, right.ID} {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			<-start
			_, err := service.PublishPackageSet(ctx, principal, id, 2)
			results <- struct {
				id  string
				err error
			}{id, err}
		}(id)
	}
	close(start)
	wg.Wait()
	close(results)
	successes := 0
	for result := range results {
		if result.err == nil {
			successes++
			continue
		}
		if auth.ErrorCode(result.err) != auth.CodeConflict {
			t.Fatalf("publish error: %v", result.err)
		}
		assertBarcodePublicationRollback(t, db, result.id)
	}
	if successes != 1 {
		t.Fatalf("successful publishes: %d", successes)
	}
	var count int64
	catalogNoError(t, db.Model(&gen.ProductPackage{}).Where("active_barcode = ?", "shared-successor").Count(&count).Error)
	if count != 1 {
		t.Fatalf("active duplicates: %d", count)
	}
}

func assertBarcodePublicationRollback(t *testing.T, db *gorm.DB, skuID string) {
	t.Helper()
	var sku gen.ProductSku
	catalogNoError(t, db.First(&sku, "id = ?", skuID).Error)
	if sku.PublishedPackageSetVersion != 1 {
		t.Fatal("failed publish changed version")
	}
	var original, draft gen.ProductPackage
	catalogNoError(t, db.Where("sku_id = ? AND package_set_version = ?", skuID, 1).First(&original).Error)
	catalogNoError(t, db.Where("sku_id = ? AND package_set_version = ?", skuID, 2).First(&draft).Error)
	if !original.Enabled || draft.Enabled {
		t.Fatal("failed publish changed package state")
	}
	var offer gen.StorePackageOffer
	var promotion gen.StorePromotion
	catalogNoError(t, db.First(&offer, "package_id = ?", original.ID).Error)
	catalogNoError(t, db.First(&promotion, "id = ?", original.ID).Error)
	if !offer.Enabled || offer.PriceFen != 100 || !promotion.Enabled {
		t.Fatal("failed publish changed offers or promotions")
	}
}

func TestPackageBarcodeMySQLUniqueErrorMapping(t *testing.T) {
	service, db, principal := barcodeMySQLFixture(t)
	sku, unit := barcodePackage(t, service, principal, "00123")
	factor := int64(5)
	err := db.Select("*").Create(&gen.ProductPackage{ID: "duplicate", SkuID: sku.ID, Name: "box", Barcode: unit.Barcode,
		ContainsPackageID: &unit.ID, ContainsQuantity: &factor, PackageSetVersion: 1, Enabled: true}).Error
	if err == nil || auth.ErrorCode(packageWriteError(err)) != auth.CodeConflict {
		t.Fatalf("unique error mapping: %v", err)
	}
}

func seedBarcodePublishCompetition(t *testing.T, service *Service, db *gorm.DB, principal *auth.WorkspacePrincipal) (*gen.ProductSku, *gen.ProductSku) {
	t.Helper()
	ctx := context.Background()
	left, leftPack := barcodePackage(t, service, principal, "left-original")
	var originalProduct gen.Product
	catalogNoError(t, db.First(&originalProduct, "id = ?", left.ProductID).Error)
	otherProduct, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Other", CategoryID: originalProduct.CategoryID})
	catalogNoError(t, err)
	right, err := service.CreateSku(ctx, principal, SkuInput{ProductID: otherProduct.ID, Name: "Other"})
	catalogNoError(t, err)
	rightPack, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: right.ID, Name: "unit", Barcode: "right-original"})
	catalogNoError(t, err)
	for _, sku := range []*gen.ProductSku{left, right} {
		_, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "unit", Barcode: "shared-successor", PackageSetVersion: 2})
		catalogNoError(t, err)
	}
	for _, pack := range []*gen.ProductPackage{leftPack, rightPack} {
		catalogNoError(t, db.Create(&gen.StorePackageOffer{ID: pack.ID, ListingID: pack.SkuID, PackageID: pack.ID, PriceFen: 100, Enabled: true}).Error)
		catalogNoError(t, db.Create(&gen.StorePromotion{ID: pack.ID, StoreID: "store", RuleKey: pack.ID, Version: 1, Enabled: true}).Error)
		catalogNoError(t, db.Create(&gen.StorePromotionTarget{ID: pack.ID, PromotionID: pack.ID, OfferID: pack.ID, RequiredQuantity: 1}).Error)
	}
	return left, right
}

func TestPackageBarcodeMySQLSameSkuDraftSerialization(t *testing.T) {
	service, _, principal := barcodeMySQLFixture(t)
	sku, _ := barcodePackage(t, service, principal, "original")
	ctx := context.Background()
	base, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "unit", PackageSetVersion: 2})
	catalogNoError(t, err)
	box, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "box", PackageSetVersion: 2, ContainsPackageID: &base.ID, ContainsQuantity: 5})
	catalogNoError(t, err)
	start := make(chan struct{})
	results := make(chan error, 2)
	for _, id := range []string{base.ID, box.ID} {
		go func(id string) {
			<-start
			_, err := service.SetPackageBarcode(ctx, principal, id, "same-draft-code")
			results <- err
		}(id)
	}
	close(start)
	successes := 0
	for range 2 {
		err := <-results
		if err == nil {
			successes++
		} else if auth.ErrorCode(err) != auth.CodeConflict {
			t.Fatalf("draft error: %v", err)
		}
	}
	if successes != 1 {
		t.Fatalf("successful draft updates: %d", successes)
	}
}
