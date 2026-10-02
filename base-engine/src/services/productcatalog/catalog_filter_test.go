package productcatalog

import (
	"context"
	"testing"
)

func TestCatalogRelationshipFiltersApplyBeforePagination(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	parent, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Drinks"})
	catalogNoError(t, err)
	child, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Coffee", ParentID: &parent.ID})
	catalogNoError(t, err)
	brand, err := service.CreateBrand(ctx, principal, "House")
	catalogNoError(t, err)
	_, err = service.CreateProduct(ctx, principal, ProductInput{Name: "No Brand", CategoryID: child.ID})
	catalogNoError(t, err)
	wanted, err := service.CreateProduct(ctx, principal, ProductInput{Name: "House Blend", CategoryID: child.ID, BrandID: &brand.ID})
	catalogNoError(t, err)
	categories, err := service.Categories(ctx, principal, CatalogFilter{ParentID: &parent.ID}, 1, 1)
	catalogNoError(t, err)
	if categories.Total != 1 || len(categories.Data) != 1 || categories.Data[0].ID != child.ID {
		t.Fatalf("filtered categories = %+v", categories)
	}
	products, err := service.Products(ctx, principal, &child.ID, CatalogFilter{BrandID: &brand.ID}, 1, 1)
	catalogNoError(t, err)
	if products.Total != 1 || len(products.Data) != 1 || products.Data[0].ID != wanted.ID {
		t.Fatalf("filtered products = %+v", products)
	}
}
