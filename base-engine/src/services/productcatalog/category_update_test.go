package productcatalog

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestProductCanKeepDisabledCategoryWhileEditing(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID}))
	_, err := service.SetCategoryEnabled(ctx, principal, category.ID, false)
	catalogNoError(t, err)
	description := "new description"
	updated := mustCatalog[*gen.Product](t)(service.UpdateProduct(ctx, principal, product.ID, ProductInput{
		Name: "Noodles", CategoryID: category.ID, Description: &description,
	}))
	if updated.Description == nil || *updated.Description != description {
		t.Fatalf("product was not updated: %+v", updated)
	}
	if _, err := service.CreateProduct(ctx, principal, ProductInput{Name: "New", CategoryID: category.ID}); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("new product accepted disabled category: %v", err)
	}
}
