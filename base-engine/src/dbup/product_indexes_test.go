package dbup

import (
	"testing"
	"time"

	"base-engine/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestStocktakeIndexUniqueness(t *testing.T) {
	db := productIndexTestDB(t)
	first := gen.StoreStocktake{ID: "count-one", StoreID: "store", InitiatedByAccountID: "actor",
		Status: gen.StocktakeStatusCounting, RequestKey: "same-request", StartedAt: time.Now()}
	if err := db.Create(&first).Error; err != nil {
		t.Fatal(err)
	}
	second := first
	second.ID = "count-two"
	if err := db.Create(&second).Error; err == nil {
		t.Fatal("duplicate store request key accepted")
	}
	line := gen.StoreStocktakeLine{ID: "line-one", StocktakeID: first.ID, BatchID: "batch", PackageID: "piece", PackageSetVersion: 1}
	if err := db.Create(&line).Error; err != nil {
		t.Fatal(err)
	}
	line.ID = "line-two"
	if err := db.Create(&line).Error; err == nil {
		t.Fatal("duplicate batch/package line accepted")
	}
}

func TestProductIndexesPreventDuplicateBusinessKeys(t *testing.T) {
	db := productIndexTestDB(t)
	if err := db.Create(&gen.ProductBrand{ID: "brand-one", OrganizationID: "hq", Name: "House", Enabled: true}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.ProductBrand{ID: "brand-two", OrganizationID: "hq", Name: "House", Enabled: true}).Error; err == nil {
		t.Fatal("duplicate headquarters brand accepted")
	}
	if err := db.Create(&gen.StoreListing{ID: "one", StoreID: "store", SkuID: "sku"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.StoreListing{ID: "two", StoreID: "store", SkuID: "sku"}).Error; err == nil {
		t.Fatal("duplicate store listing accepted")
	}
	barcode := "123"
	if err := db.Create(&gen.ProductPackage{ID: "p1", SkuID: "sku", Name: "piece", Barcode: &barcode, PackageSetVersion: 1, Enabled: true}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.ProductPackage{ID: "p2", SkuID: "sku2", Name: "piece", Barcode: &barcode, PackageSetVersion: 1, Enabled: true}).Error; err == nil {
		t.Fatal("duplicate barcode accepted")
	}
	otherBarcode := "456"
	if err := db.Create(&gen.ProductPackage{ID: "p3", SkuID: "sku", Name: "second base", Barcode: &otherBarcode, PackageSetVersion: 1, Enabled: true}).Error; err == nil {
		t.Fatal("duplicate base package accepted")
	}
	if err := db.Model(&gen.ProductPackage{}).Where("id = ?", "p1").Update("enabled", false).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.ProductPackage{ID: "p4", SkuID: "sku2", Name: "new base", Barcode: &barcode, PackageSetVersion: 1, Enabled: true}).Error; err != nil {
		t.Fatalf("retired barcode could not be reused: %v", err)
	}
}

func TestSpecificationIndexesPreventDuplicateNames(t *testing.T) {
	db := productIndexTestDB(t)
	if err := db.Create(&gen.SpecificationDefinition{ID: "spec-one", OrganizationID: "hq", Name: "Color", Enabled: true}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.SpecificationDefinition{ID: "spec-two", OrganizationID: "hq", Name: "Color", Enabled: true}).Error; err == nil {
		t.Fatal("duplicate specification name accepted")
	}
	if err := db.Create(&gen.SpecificationValue{ID: "red-one", SpecificationID: "spec-one", Name: "Red", Enabled: true}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.SpecificationValue{ID: "red-two", SpecificationID: "spec-one", Name: "Red", Enabled: true}).Error; err == nil {
		t.Fatal("duplicate specification value accepted")
	}
}

func productIndexTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:product-indexes?mode=memory&cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, IgnoreRelationshipsWhenMigrating: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range []any{&gen.ProductBrand{}, &gen.SpecificationDefinition{}, &gen.SpecificationValue{}, &gen.ProductSpecificationChoice{}, &gen.ProductSkuSpecificationValue{}, &gen.ProductCategory{}, &gen.ProductPackage{}, &gen.StoreListing{}, &gen.StorePackageOffer{}, &gen.StoreInventoryBatch{}, &gen.StoreStockBalance{}, &gen.StoreStockMovement{}, &gen.StoreStocktake{}, &gen.StoreStocktakeLine{}, &gen.StorePromotion{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	if err := EnsureProductIndexes(db); err != nil {
		t.Fatal(err)
	}
	if err := EnsureProductIndexes(db); err != nil {
		t.Fatal(err)
	}
	return db
}
