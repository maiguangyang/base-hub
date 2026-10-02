package productcatalog

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestBrandAndProductUseHeadquartersRelationship(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	brand := mustCatalog[*gen.ProductBrand](t)(service.CreateBrand(ctx, principal, " House Brand "))
	if brand.Name != "House Brand" || brand.OrganizationID != "hq" {
		t.Fatalf("brand = %+v", brand)
	}
	if _, err := service.CreateBrand(ctx, principal, "House Brand"); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("duplicate brand error = %v", err)
	}
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID, BrandID: &brand.ID}))
	if product.BrandID == nil || *product.BrandID != brand.ID {
		t.Fatalf("product = %+v", product)
	}
	if err := service.DeleteBrand(ctx, principal, brand.ID); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("referenced brand delete = %v", err)
	}
	mustCatalogErrorFree(t, service.DeleteProduct(ctx, principal, product.ID))
	mustCatalogErrorFree(t, service.DeleteBrand(ctx, principal, brand.ID))
}

func TestProductRejectsForeignAndDisabledBrands(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
	if err != nil {
		t.Fatal(err)
	}
	brand, err := service.CreateBrand(ctx, principal, "House Brand")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetBrandEnabled(ctx, principal, brand.ID, false); err != nil {
		t.Fatal(err)
	}
	if _, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID, BrandID: &brand.ID}); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("disabled brand error = %v", err)
	}
	foreign := "outside"
	if _, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID, BrandID: &foreign}); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("foreign brand error = %v", err)
	}
}

func TestProductCanKeepDisabledBrandWhenEditing(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
	if err != nil {
		t.Fatal(err)
	}
	brand, err := service.CreateBrand(ctx, principal, "House Brand")
	if err != nil {
		t.Fatal(err)
	}
	product, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID, BrandID: &brand.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetBrandEnabled(ctx, principal, brand.ID, false); err != nil {
		t.Fatal(err)
	}
	updated, err := service.UpdateProduct(ctx, principal, product.ID, ProductInput{Name: "Rice", CategoryID: category.ID, BrandID: &brand.ID})
	if err != nil || updated.Brand == nil || updated.Brand.Name != brand.Name {
		t.Fatalf("updated = %+v, %v", updated, err)
	}
}
