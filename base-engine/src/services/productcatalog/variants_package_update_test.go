package productcatalog

import (
	"context"
	"testing"

	"base-engine/gen"
)

func TestExistingSkuSwitchesPackageTemplateWithHistoryPreserved(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	oldTemplate := createPackageChain(t, service, principal)
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{
		Name: "Noodles", CategoryID: category.ID, DefaultPackageTemplateID: &oldTemplate.ID,
	}))
	var sku gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).First(&sku).Error)
	oldPackages := packagesForSku(t, db, sku.ID)
	newTemplate := mustCatalog[*gen.ProductPackageTemplate](t)(service.CreatePackageTemplate(ctx, principal, PackageTemplateInput{Name: "bottle"}))
	_, err := service.UpdateProduct(ctx, principal, product.ID, ProductInput{
		Name: "Noodles", CategoryID: category.ID, DefaultPackageTemplateID: &oldTemplate.ID,
		SkuOverrides: []SkuOverride{{Enabled: true, PackageTemplateID: &newTemplate.ID}},
	})
	catalogNoError(t, err)
	catalogNoError(t, db.First(&sku, "id = ?", sku.ID).Error)
	if sku.PublishedPackageSetVersion != 2 {
		t.Fatalf("published version = %d", sku.PublishedPackageSetVersion)
	}
	packages := packagesForSku(t, db, sku.ID)
	if len(packages) != len(oldPackages)+1 {
		t.Fatalf("package history length = %d", len(packages))
	}
	for _, pack := range packages {
		if pack.Enabled != (pack.PackageSetVersion == 2) {
			t.Fatalf("mixed active package versions: %+v", packages)
		}
	}
}

func TestExistingSkuCanStopUsingPackages(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	template := createPackageChain(t, service, principal)
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{
		Name: "Noodles", CategoryID: category.ID, DefaultPackageTemplateID: &template.ID,
	}))
	var sku gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).First(&sku).Error)
	packages := packagesForSku(t, db, sku.ID)
	seedRetirementOffer(t, db, sku.ID, packages[len(packages)-1].ID)
	_, err := service.UpdateProduct(ctx, principal, product.ID, ProductInput{
		Name: "Noodles", CategoryID: category.ID, DefaultPackageTemplateID: &template.ID,
		SkuOverrides: []SkuOverride{{Enabled: true, DisableDefaultPackage: true}},
	})
	catalogNoError(t, err)
	for _, pack := range packagesForSku(t, db, sku.ID) {
		if pack.Enabled {
			t.Fatalf("package stayed active: %+v", pack)
		}
	}
	var offer gen.StorePackageOffer
	catalogNoError(t, db.First(&offer, "id = ?", "offer").Error)
	if offer.Enabled {
		t.Fatal("retired package offer stayed active")
	}
}

func TestRepeatedExistingSkuTemplateSelectionDoesNotRepublish(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID}))
	template := createPackageChain(t, service, principal)
	input := ProductInput{Name: "Noodles", CategoryID: category.ID, SkuOverrides: []SkuOverride{{Enabled: true, PackageTemplateID: &template.ID}}}
	_, err := service.UpdateProduct(ctx, principal, product.ID, input)
	catalogNoError(t, err)
	_, err = service.UpdateProduct(ctx, principal, product.ID, input)
	catalogNoError(t, err)
	var sku gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).First(&sku).Error)
	if sku.PublishedPackageSetVersion != 1 || len(packagesForSku(t, db, sku.ID)) != 3 {
		t.Fatalf("same template was republished: version=%d", sku.PublishedPackageSetVersion)
	}
}

func TestDisabledExistingSkuCanSwitchTemplateWithoutBeingEnabled(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID}))
	var sku gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).First(&sku).Error)
	_, err := service.SetSkuEnabled(ctx, principal, sku.ID, false)
	catalogNoError(t, err)
	template := createPackageChain(t, service, principal)
	_, err = service.UpdateProduct(ctx, principal, product.ID, ProductInput{
		Name: "Noodles", CategoryID: category.ID,
		SkuOverrides: []SkuOverride{{Enabled: false, PackageTemplateID: &template.ID}},
	})
	catalogNoError(t, err)
	catalogNoError(t, db.First(&sku, "id = ?", sku.ID).Error)
	if sku.Enabled || sku.PublishedPackageSetVersion != 1 || len(packagesForSku(t, db, sku.ID)) != 3 {
		t.Fatalf("disabled SKU package change failed: %+v", sku)
	}
}
