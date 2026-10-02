package productcatalog

import (
	"context"
	"base-engine/gen"
	"base-engine/src/dbup"
	"testing"
)

func TestCatalogCodeRemovalMySQL(t *testing.T) {
	service, db, principal := barcodeMySQLFixture(t)
	ctx := context.Background()
	category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
	catalogNoError(t, err)
	product, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Kept", CategoryID: category.ID})
	catalogNoError(t, err)
	for _, table := range []string{"products", "product_skus"} {
		catalogNoError(t, db.Exec("ALTER TABLE "+table+" ADD COLUMN code VARCHAR(13) NOT NULL DEFAULT 'old-code'").Error)
		catalogNoError(t, db.Exec("CREATE UNIQUE INDEX idx_"+table+"_code ON "+table+" (code)").Error)
	}
	catalogNoError(t, dbup.EnsureProductIndexes(db))
	catalogNoError(t, dbup.EnsureProductIndexes(db))
	var kept gen.Product
	catalogNoError(t, db.First(&kept, "id = ?", product.ID).Error)
	if kept.Name != "Kept" || db.Migrator().HasColumn(&gen.Product{}, "code") {
		t.Fatal("catalog migration changed identity or retained code")
	}
	_, err = service.CreateProduct(ctx, principal, ProductInput{Name: "New", CategoryID: category.ID})
	catalogNoError(t, err)
}
