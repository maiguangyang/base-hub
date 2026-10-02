package productcatalog

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestProductVariantsUseSelectedValuesAndKeepIdentity(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	plain := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{Name: "Plain", CategoryID: category.ID}))
	var plainSkus []gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", plain.ID).Find(&plainSkus).Error)
	if len(plainSkus) != 1 || plainSkus[0].ID == "" {
		t.Fatalf("default SKU = %+v", plainSkus)
	}
	color := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "Color"))
	size := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "Size"))
	red := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, color.ID, "Red"))
	blue := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, color.ID, "Blue"))
	small := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, size.ID, "Small"))
	large := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, size.ID, "Large"))
	input := ProductInput{Name: "Shirt", CategoryID: category.ID, Selections: []SpecificationSelection{
		{SpecificationID: color.ID, ValueIDs: []string{red.ID, blue.ID}},
		{SpecificationID: size.ID, ValueIDs: []string{small.ID, large.ID}},
	}}
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, input))
	if product.ID == "" {
		t.Fatal("missing product ID")
	}
	var first []gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).Order("id").Find(&first).Error)
	if len(first) != 4 {
		t.Fatalf("SKU count = %d", len(first))
	}
	_, err := service.UpdateProduct(ctx, principal, product.ID, input)
	catalogNoError(t, err)
	var second []gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).Order("id").Find(&second).Error)
	if len(second) != 4 {
		t.Fatalf("SKU count after save = %d", len(second))
	}
	assertSameSkuIdentity(t, first, second)
}

func TestGeneratedSkuNameContainsOnlySpecificationValues(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "方便食品"}))
	flavor := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "口味"))
	origin := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "产地"))
	mild := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, flavor.ID, "微辣"))
	korea := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, origin.ID, "韩国"))
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{
		Name: "泡面", CategoryID: category.ID, Selections: []SpecificationSelection{
			{SpecificationID: flavor.ID, ValueIDs: []string{mild.ID}},
			{SpecificationID: origin.ID, ValueIDs: []string{korea.ID}},
		},
	}))
	var sku gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).First(&sku).Error)
	if sku.Name != "微辣 / 韩国" {
		t.Fatalf("generated SKU name = %q, want specification values only", sku.Name)
	}
}

func TestProductSaveOverwritesLegacySkuNameFromSpecificationValues(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "方便食品"}))
	flavor := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "口味"))
	mild := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, flavor.ID, "微辣"))
	input := ProductInput{Name: "泡面", CategoryID: category.ID,
		Selections: []SpecificationSelection{{SpecificationID: flavor.ID, ValueIDs: []string{mild.ID}}}}
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, input))
	var sku gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).First(&sku).Error)
	catalogNoError(t, db.Model(&sku).Update("name", "泡面 / 微辣").Error)
	_, err := service.UpdateProduct(ctx, principal, product.ID, input)
	catalogNoError(t, err)
	catalogNoError(t, db.First(&sku, "id = ?", sku.ID).Error)
	if sku.Name != "微辣" {
		t.Fatalf("SKU name after product save = %q", sku.Name)
	}
}

func TestGeneratedSkuWithoutSpecificationsUsesDefaultName(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "方便食品"}))
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{Name: "泡面", CategoryID: category.ID}))
	var sku gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).First(&sku).Error)
	if sku.Name != "默认规格" {
		t.Fatalf("default SKU name = %q", sku.Name)
	}
}

func assertSameSkuIdentity(t *testing.T, first, second []gen.ProductSku) {
	t.Helper()
	for i := range first {
		if first[i].ID != second[i].ID {
			t.Fatalf("SKU identity changed: %+v -> %+v", first[i], second[i])
		}
	}
}

func TestNineVariantCombinationsRetireAndRestoreWithoutChangingIDs(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	color := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "Color"))
	size := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "Size"))
	input := ProductInput{Name: "Shirt", CategoryID: category.ID, Selections: []SpecificationSelection{{SpecificationID: color.ID}, {SpecificationID: size.ID}}}
	input.Selections[0].ValueIDs = createVariantValues(t, service, ctx, principal, color.ID, "Red", "Blue", "Green")
	input.Selections[1].ValueIDs = createVariantValues(t, service, ctx, principal, size.ID, "Small", "Medium", "Large")
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, input))
	var skus []gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).Find(&skus).Error)
	if len(skus) != 9 {
		t.Fatalf("want nine combinations, got %d", len(skus))
	}
	originalIDs := map[string]string{}
	for _, sku := range skus {
		originalIDs[sku.ID] = sku.ID
	}
	trimmed := input
	trimmed.Selections = []SpecificationSelection{
		{SpecificationID: color.ID, ValueIDs: input.Selections[0].ValueIDs[:2]},
		{SpecificationID: size.ID, ValueIDs: input.Selections[1].ValueIDs},
	}
	_, err := service.UpdateProduct(ctx, principal, product.ID, trimmed)
	catalogNoError(t, err)
	var disabled int64
	catalogNoError(t, db.Model(&gen.ProductSku{}).Where("product_id = ? AND enabled = ? AND selection_retired = ?", product.ID, false, true).Count(&disabled).Error)
	if disabled != 3 {
		t.Fatalf("retired SKUs = %d", disabled)
	}
	assertRetiredSkuCannotBeEnabledDirectly(t, service, db, ctx, principal, product.ID)
	_, err = service.UpdateProduct(ctx, principal, product.ID, input)
	catalogNoError(t, err)
	var restored []gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).Find(&restored).Error)
	if len(restored) != 9 {
		t.Fatalf("restored SKU count = %d", len(restored))
	}
	for _, sku := range restored {
		assertRestoredSku(t, sku, originalIDs)
	}
}

func assertRetiredSkuCannotBeEnabledDirectly(t *testing.T, service *Service, db *gorm.DB, ctx context.Context, principal *auth.WorkspacePrincipal, productID string) {
	t.Helper()
	var retired gen.ProductSku
	catalogNoError(t, db.Where("product_id = ? AND selection_retired = ?", productID, true).First(&retired).Error)
	if _, err := service.SetSkuEnabled(ctx, principal, retired.ID, true); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("direct activation of retired SKU = %v", err)
	}
}

func createVariantValues(t *testing.T, service *Service, ctx context.Context, principal *auth.WorkspacePrincipal, specificationID string, names ...string) []string {
	t.Helper()
	ids := make([]string, 0, len(names))
	for _, name := range names {
		value := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, specificationID, name))
		ids = append(ids, value.ID)
	}
	return ids
}

func assertRestoredSku(t *testing.T, sku gen.ProductSku, originalIDs map[string]string) {
	t.Helper()
	if !sku.Enabled || originalIDs[sku.ID] != sku.ID {
		t.Fatalf("identity changed: %+v", sku)
	}
}

func TestVariantValidationIsAtomicAndPackageTemplatesCopyPerSku(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	spec := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, "Flavor"))
	a := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, spec.ID, "Mild"))
	b := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, spec.ID, "Hot"))
	box := createPackageChain(t, service, principal)
	input := ProductInput{Name: "Noodles", CategoryID: category.ID, Selections: []SpecificationSelection{{SpecificationID: spec.ID, ValueIDs: []string{a.ID, b.ID}}},
		DefaultPackageTemplateID: &box.ID, SkuOverrides: []SkuOverride{{ValueIDs: []string{b.ID}, Enabled: true, DisableDefaultPackage: true}}}
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, input))
	var skus []gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).Find(&skus).Error)
	if len(skus) != 2 {
		t.Fatalf("SKU count = %d", len(skus))
	}
	var packages int64
	catalogNoError(t, db.Model(&gen.ProductPackage{}).Where("sku_id IN ?", []string{skus[0].ID, skus[1].ID}).Count(&packages).Error)
	if packages != 3 {
		t.Fatalf("copied packages = %d, want one three-unit chain", packages)
	}
	input.Selections[0].ValueIDs = []string{a.ID, a.ID}
	if _, err := service.UpdateProduct(ctx, principal, product.ID, input); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("duplicate selection error = %v", err)
	}
	var persisted int64
	catalogNoError(t, db.Model(&gen.ProductSpecificationChoice{}).Where("product_id = ?", product.ID).Count(&persisted).Error)
	if persisted != 2 {
		t.Fatalf("selection changed after rejected update: %d", persisted)
	}
	input.Selections[0].ValueIDs = []string{a.ID, b.ID}
	_, err := service.SetSpecificationValueEnabled(ctx, principal, b.ID, false)
	catalogNoError(t, err)
	_, err = service.UpdateProduct(ctx, principal, product.ID, input)
	catalogNoError(t, err)
}

func TestDisabledNewVariantKeepsCopiedPackageTemplate(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	box := createPackageChain(t, service, principal)
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{
		Name: "Noodles", CategoryID: category.ID, DefaultPackageTemplateID: &box.ID,
		SkuOverrides: []SkuOverride{{Enabled: false}},
	}))
	var sku gen.ProductSku
	catalogNoError(t, db.Where("product_id = ?", product.ID).First(&sku).Error)
	if sku.Enabled || sku.PublishedPackageSetVersion != 1 {
		t.Fatalf("disabled SKU lost its package version: %+v", sku)
	}
	packages := packagesForSku(t, db, sku.ID)
	if len(packages) != 3 {
		t.Fatalf("copied package count = %d", len(packages))
	}
}

func TestVariantLimitRejectsAllWrites(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	input := ProductInput{Name: "Too Many", CategoryID: category.ID}
	for i := 0; i < 7; i++ {
		spec := mustCatalog[*gen.SpecificationDefinition](t)(service.CreateSpecification(ctx, principal, string(rune('A'+i))))
		selection := SpecificationSelection{SpecificationID: spec.ID}
		for _, name := range []string{"One", "Two"} {
			value := mustCatalog[*gen.SpecificationValue](t)(service.CreateSpecificationValue(ctx, principal, spec.ID, name))
			selection.ValueIDs = append(selection.ValueIDs, value.ID)
		}
		input.Selections = append(input.Selections, selection)
	}
	if _, err := service.CreateProduct(ctx, principal, input); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("limit error = %v", err)
	}
	var count int64
	catalogNoError(t, db.Model(&gen.Product{}).Where("name = ?", input.Name).Count(&count).Error)
	if count != 0 {
		t.Fatalf("rejected product was written: %d", count)
	}
}
