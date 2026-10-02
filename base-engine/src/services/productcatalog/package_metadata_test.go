package productcatalog

import (
	"context"
	"testing"
)

func TestPackageMetadataEditKeepsBarcodeAndConversionIdentity(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
	catalogNoError(t, err)
	product, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID})
	catalogNoError(t, err)
	sku, err := service.CreateSku(ctx, principal, SkuInput{Name: "Spicy", ProductID: product.ID})
	catalogNoError(t, err)
	unit, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "unit", Barcode: "unit-001"})
	catalogNoError(t, err)
	price := int64(1900)
	updated, err := service.UpdatePackageMetadata(ctx, principal, unit.ID, "single unit", &price)
	catalogNoError(t, err)
	if updated.ID != unit.ID || updated.Name != "single unit" || updated.Barcode == nil || *updated.Barcode != "unit-001" || updated.PackageSetVersion != unit.PackageSetVersion || updated.SuggestedPriceFen == nil || *updated.SuggestedPriceFen != price {
		t.Fatalf("package metadata edit changed identity: %+v", updated)
	}
	_, err = service.UpdatePackageMetadata(ctx, principal, unit.ID, "invalid", nil)
	catalogNoError(t, err)
}
