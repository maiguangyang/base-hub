package dbup

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"base-engine/gen"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestProductIndexesMySQL(t *testing.T) {
	db := productMySQLTestDB(t)
	if err := EnsureProductIndexes(db); err != nil {
		t.Fatal(err)
	}
	if err := EnsureProductIndexes(db); err != nil {
		t.Fatalf("second migration: %v", err)
	}
	assertProductMySQLConstraints(t, db)
}

func productMySQLTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	rootDSN := os.Getenv("PRODUCT_TEST_MYSQL_DSN")
	if rootDSN == "" {
		t.Skip("PRODUCT_TEST_MYSQL_DSN is required for the disposable MySQL test")
	}
	if !strings.HasSuffix(rootDSN, "/") {
		t.Fatal("PRODUCT_TEST_MYSQL_DSN must end in / with no database selected")
	}
	admin, err := gorm.Open(mysql.Open(rootDSN+"mysql?parseTime=true"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	database := fmt.Sprintf("product_test_%d", time.Now().UnixNano())
	if err := admin.Exec("CREATE DATABASE `" + database + "`").Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := admin.Exec("DROP DATABASE `" + database + "`").Error; err != nil {
			t.Errorf("drop database: %v", err)
		}
	})
	db, err := gorm.Open(mysql.Open(rootDSN+database+"?parseTime=true"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range []any{&gen.ProductBrand{}, &gen.SpecificationDefinition{}, &gen.SpecificationValue{}, &gen.ProductSpecificationChoice{}, &gen.ProductSkuSpecificationValue{}, &gen.ProductCategory{}, &gen.Product{}, &gen.ProductSku{}, &gen.ProductPackage{},
		&gen.StoreListing{}, &gen.StorePackageOffer{}, &gen.StorePriceRevision{}, &gen.StoreInventoryBatch{}, &gen.StoreStockBalance{},
		&gen.StoreStockMovement{}, &gen.StoreStocktake{}, &gen.StoreStocktakeLine{}, &gen.StorePromotion{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
	}
	return db
}

func assertProductMySQLConstraints(t *testing.T, db *gorm.DB) {
	t.Helper()
	base := gen.ProductPackage{ID: "base", Name: "Unit", SkuID: "sku", Barcode: productPointer("CODE"), PackageSetVersion: 1, Enabled: true}
	if err := db.Create(&base).Error; err != nil {
		t.Fatal(err)
	}
	secondBase := gen.ProductPackage{ID: "second", Name: "Other", SkuID: "sku", PackageSetVersion: 1, Enabled: true}
	if err := db.Create(&secondBase).Error; err == nil {
		t.Fatal("duplicate base package accepted")
	}
	barcode := gen.ProductPackage{ID: "barcode", Name: "Box", SkuID: "sku", Barcode: productPointer("CODE"), ContainsPackageID: &base.ID, ContainsQuantity: pointerInt64(2), PackageSetVersion: 1, Enabled: true}
	if err := db.Create(&barcode).Error; err == nil {
		t.Fatal("duplicate active barcode accepted")
	}
	if err := db.Model(&base).Update("enabled", false).Error; err != nil {
		t.Fatal(err)
	}
	barcode.ID = "reused"
	if err := db.Create(&barcode).Error; err != nil {
		t.Fatalf("retired barcode reuse: %v", err)
	}
}

func productPointer(value string) *string { return &value }
func pointerInt64(value int64) *int64     { return &value }
