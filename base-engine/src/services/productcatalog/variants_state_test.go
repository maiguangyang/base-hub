package productcatalog

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestDisabledSelectedValueCanStayButCannotCreateNewCombination(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	flavor := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "Flavor"))
	hot := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, flavor.ID, "Hot"))
	input := ProductInput{Name: "Noodles", CategoryID: category.ID, Selections: []SpecificationSelection{{SpecificationID: flavor.ID, ValueIDs: []string{hot.ID}}}}
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, input))
	_, err := service.SetSpecificationValueEnabled(ctx, principal, hot.ID, false)
	catalogNoError(t, err)
	_, err = service.UpdateProduct(ctx, principal, product.ID, input)
	catalogNoError(t, err)
	_, err = service.SetSpecificationEnabled(ctx, principal, flavor.ID, false)
	catalogNoError(t, err)
	_, err = service.UpdateProduct(ctx, principal, product.ID, input)
	catalogNoError(t, err)
	size := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "Size"))
	small := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, size.ID, "Small"))
	input.Selections = append(input.Selections, SpecificationSelection{SpecificationID: size.ID, ValueIDs: []string{small.ID}})
	if _, err := service.UpdateProduct(ctx, principal, product.ID, input); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("new combination from disabled value: %v", err)
	}
	var count int64
	catalogNoError(t, db.Model(&gen.ProductSku{}).Where("product_id = ?", product.ID).Count(&count).Error)
	if count != 1 {
		t.Fatalf("SKU count after rejected edit = %d", count)
	}
}

func TestProductEditPreservesExistingDisabledSkuWithoutOverride(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	input := ProductInput{Name: "Noodles", CategoryID: category.ID}
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, input))
	var sku gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).First(&sku).Error)
	_, err := service.SetSkuEnabled(ctx, principal, sku.ID, false)
	catalogNoError(t, err)
	input.Name = "Rice noodles"
	_, err = service.UpdateProduct(ctx, principal, product.ID, input)
	catalogNoError(t, err)
	catalogNoError(t, db.First(&sku, "id = ?", sku.ID).Error)
	if sku.Enabled || sku.Name != "默认规格" {
		t.Fatalf("SKU after product edit = %+v", sku)
	}
}

func TestRemovedManuallyDisabledSkuCannotBeEnabledDirectly(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	input := ProductInput{Name: "Noodles", CategoryID: category.ID}
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, input))
	var old gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).First(&old).Error)
	_, err := service.SetSkuEnabled(ctx, principal, old.ID, false)
	catalogNoError(t, err)
	spec := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "Flavor"))
	value := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, spec.ID, "Hot"))
	input.Selections = []SpecificationSelection{{SpecificationID: spec.ID, ValueIDs: []string{value.ID}}}
	_, err = service.UpdateProduct(ctx, principal, product.ID, input)
	catalogNoError(t, err)
	if _, err := service.SetSkuEnabled(ctx, principal, old.ID, true); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("direct activation of removed SKU = %v", err)
	}
}

func TestManualDisableOfRetiredSkuReturnsCurrentState(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	input := ProductInput{Name: "Noodles", CategoryID: category.ID}
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, input))
	var old gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).First(&old).Error)
	spec := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "Flavor"))
	value := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, spec.ID, "Hot"))
	input.Selections = []SpecificationSelection{{SpecificationID: spec.ID, ValueIDs: []string{value.ID}}}
	_, err := service.UpdateProduct(ctx, principal, product.ID, input)
	catalogNoError(t, err)
	result := mustCatalog[*gen.ProductSku](t)(service.SetSkuEnabled(ctx, principal, old.ID, false))
	if result.SelectionRetired {
		t.Fatalf("returned SKU retains stale retired state: %+v", result)
	}
}
