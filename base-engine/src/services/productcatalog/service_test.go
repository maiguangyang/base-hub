package productcatalog

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func catalogFixture(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range []any{
		&gen.Organization{}, &gen.ProductCategory{}, &gen.ProductBrand{}, &gen.Product{}, &gen.SpecificationDefinition{}, &gen.SpecificationValue{}, &gen.ProductSpecificationChoice{}, &gen.ProductSkuSpecificationValue{},
		&gen.ProductPackageTemplate{}, &gen.ProductSku{}, &gen.ProductPackage{}, &gen.StoreListing{}, &gen.StorePackageOffer{}, &gen.StorePromotion{},
		&gen.StorePromotionTarget{}, &gen.StoreInventoryBatch{}, &gen.StoreStockMovement{}, &gen.AuditLog{},
	} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			if err := db.Exec("DROP INDEX IF EXISTS `" + index + "`").Error; err != nil {
				t.Fatal(err)
			}
		}
	}
	hq := gen.Organization{ID: "hq", Code: "HQ", Name: "HQ", Type: gen.OrganizationTypeHeadquarters, Status: gen.OrganizationStatusActive}
	if err := db.Create(&hq).Error; err != nil {
		t.Fatal(err)
	}
	hqID := hq.ID
	principal := &auth.WorkspacePrincipal{
		AccountID: "admin", SessionID: "session", WorkspaceType: auth.WorkspaceTypeHeadquarters,
		OrganizationID: &hqID, Permissions: map[string]struct{}{
			"hqProductCatalog:read": {}, "hqProductCatalog:manage": {},
		},
	}
	return NewService(db, audit.NewService()), db, principal
}

func catalogNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCategoryCreateBelongsToHeadquarters(t *testing.T) {
	service, _, principal := catalogFixture(t)
	item, err := service.CreateCategory(context.Background(), principal, CategoryInput{Name: "Food"})
	if err != nil {
		t.Fatal(err)
	}
	if item.OrganizationID != "hq" || item.Name != "Food" || item.ID == "" {
		t.Fatalf("category = %+v", item)
	}
	principal.WorkspaceType = auth.WorkspaceTypeFranchise
	if _, err := service.CreateCategory(context.Background(), principal, CategoryInput{Name: "Other"}); auth.ErrorCode(err) != auth.CodeWorkspaceForbidden {
		t.Fatalf("franchise create code = %v", err)
	}
}

func TestCatalogFiltersApplyBeforePagination(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	if _, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}); err != nil {
		t.Fatal(err)
	}
	archived, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Cold Food"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetCategoryEnabled(ctx, principal, archived.ID, false); err != nil {
		t.Fatal(err)
	}
	term, enabled := "Food", false
	page, err := service.Categories(ctx, principal, CatalogFilter{Q: &term, Enabled: &enabled}, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Data) != 1 || page.Data[0].ID != archived.ID {
		t.Fatalf("filtered page = %+v", page)
	}
}

func TestCategoryParentAndNameConstraints(t *testing.T) {
	service, db, principal := catalogFixture(t)
	parent, err := service.CreateCategory(context.Background(), principal, CategoryInput{Name: "Food"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.CreateCategory(context.Background(), principal, CategoryInput{Name: " Food "}); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("duplicate name: %v", err)
	}
	child, err := service.CreateCategory(context.Background(), principal, CategoryInput{Name: "Snack", ParentID: &parent.ID})
	if err != nil {
		t.Fatal(err)
	}
	if child.ParentID == nil || *child.ParentID != parent.ID {
		t.Fatalf("parent = %v", child.ParentID)
	}
	other := gen.ProductCategory{ID: "other-category", OrganizationID: "other-hq", Name: "External", Enabled: true}
	if err := db.Create(&other).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.CreateCategory(context.Background(), principal, CategoryInput{Name: "Wrong", ParentID: &other.ID}); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("foreign parent: %v", err)
	}
	if err := service.MoveCategory(context.Background(), principal, parent.ID, &child.ID); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("cycle: %v", err)
	}
	if _, err := service.UpdateCategory(context.Background(), principal, child.ID, "Food", nil); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("duplicate category rename: %v", err)
	}
}

func TestProductMetadataCanBeMaintainedWithoutChangingIdentity(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
	if err != nil {
		t.Fatal(err)
	}
	product, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID})
	if err != nil {
		t.Fatal(err)
	}
	description := "New description"
	updated, err := service.UpdateProduct(ctx, principal, product.ID, ProductInput{Name: "Noodles", CategoryID: category.ID, Description: &description})
	if err != nil || updated.ID != product.ID || updated.Description == nil || *updated.Description != description {
		t.Fatalf("product update: %+v, %v", updated, err)
	}
}

func TestSkuMetadataCanBeMaintainedWithoutChangingIdentity(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
	if err != nil {
		t.Fatal(err)
	}
	product, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID})
	if err != nil {
		t.Fatal(err)
	}
	sku, err := service.CreateSku(ctx, principal, SkuInput{Name: "Spicy", ProductID: product.ID})
	if err != nil {
		t.Fatal(err)
	}
	allergens := "Milk"
	updatedSku, err := service.UpdateSku(ctx, principal, sku.ID, SkuInput{Name: "Spicy", ProductID: product.ID, Allergens: &allergens})
	if err != nil || updatedSku.ID != sku.ID || updatedSku.Allergens == nil || *updatedSku.Allergens != allergens {
		t.Fatalf("sku update: %+v, %v", updatedSku, err)
	}
	if _, err := service.UpdateSku(ctx, principal, sku.ID, SkuInput{Name: "Moved", ProductID: "other-product"}); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("sku product reassignment: %v", err)
	}
}

func TestProductSkuAndPackageHierarchy(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
	if err != nil {
		t.Fatal(err)
	}
	product, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID})
	if err != nil {
		t.Fatal(err)
	}
	sku, err := service.CreateSku(ctx, principal, SkuInput{Name: "Spicy", ProductID: product.ID})
	if err != nil {
		t.Fatal(err)
	}
	unit, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "piece", Barcode: "001"})
	if err != nil {
		t.Fatal(err)
	}
	bag, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "bag", ContainsPackageID: &unit.ID, ContainsQuantity: 12})
	if err != nil {
		t.Fatal(err)
	}
	assertUnsetPackageBarcode(t, bag)
	box, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "box", ContainsPackageID: &bag.ID, ContainsQuantity: 10})
	if err != nil {
		t.Fatal(err)
	}
	if box.ContainsPackageID == nil || *box.ContainsPackageID != bag.ID {
		t.Fatalf("box = %+v", box)
	}
	assertInvalidPackageDefinitions(t, service, ctx, principal, sku.ID, unit.ID)
}

func assertUnsetPackageBarcode(t *testing.T, pack *gen.ProductPackage) {
	t.Helper()
	if pack.Barcode != nil {
		t.Fatalf("unexpected automatic barcode: %+v", pack.Barcode)
	}
}

func assertInvalidPackageDefinitions(t *testing.T, service *Service, ctx context.Context, principal *auth.WorkspacePrincipal, skuID, unitID string) {
	t.Helper()
	if _, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: skuID, Name: "bad", ContainsPackageID: &unitID, ContainsQuantity: 0}); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("zero factor: %v", err)
	}
	if _, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: skuID, Name: "duplicate", Barcode: "001"}); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("barcode: %v", err)
	}
}

func TestPackageSetHasOneBaseAndVersionedContainment(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
	if err != nil {
		t.Fatal(err)
	}
	product, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Tea", CategoryID: category.ID})
	if err != nil {
		t.Fatal(err)
	}
	sku, err := service.CreateSku(ctx, principal, SkuInput{Name: "Green", ProductID: product.ID})
	if err != nil {
		t.Fatal(err)
	}
	base, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "unit", PackageSetVersion: 1})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "second unit", PackageSetVersion: 1}); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("second base: %v", err)
	}
	if _, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "new version", PackageSetVersion: 2}); err != nil {
		t.Fatal(err)
	}
	if _, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "cross version", PackageSetVersion: 2, ContainsPackageID: &base.ID, ContainsQuantity: 6}); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("cross version: %v", err)
	}
}
