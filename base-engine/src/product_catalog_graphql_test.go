package src_test

import (
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestGraphQLProductCatalogAndGeneratedIsolation(t *testing.T) {
	db := openResolverTestDB(t)
	for _, model := range []any{&gen.ProductCategory{}, &gen.Product{}, &gen.ProductSku{}, &gen.ProductPackage{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	cfg := resolverTestConfig()
	seedResolverPermission(t, db, "permission-product-read", "hqProductCatalog:read", gen.PermissionScopeSystem)
	seedResolverPermission(t, db, "permission-product-manage", "hqProductCatalog:manage", gen.PermissionScopeSystem)
	hq := seedResolverWorkspace(t, db, cfg, "hq", gen.OrganizationTypeHeadquarters, gen.RoleKindHqSuperAdmin)
	franchise := seedResolverWorkspace(t, db, cfg, "franchise", gen.OrganizationTypeFranchise, gen.RoleKindFranchiseOwner)
	handler := governanceHandler(db, cfg)
	created := executeResolverGraphQL(t, handler, hq, `mutation { hqCreateProductCategory(input:{name:"Food"}) { id name enabled } }`)
	if strings.Contains(created.Body.String(), `"errors"`) || !strings.Contains(created.Body.String(), `"name":"Food"`) {
		t.Fatalf("create: %s", created.Body.String())
	}
	denied := executeResolverGraphQL(t, handler, franchise, `mutation { hqCreateProductCategory(input:{name:"Foreign"}) { id } }`)
	assertGraphQLErrorCode(t, denied, auth.CodePermissionDenied)
	generated := executeResolverGraphQL(t, handler, hq, `query { productCategories { total } }`)
	assertGraphQLErrorCode(t, generated, auth.CodePermissionDenied)
}

func TestGraphQLHeadquartersReordersSpecificationValuesWithoutOpeningGeneratedWrites(t *testing.T) {
	db := openResolverTestDB(t)
	for _, model := range []any{&gen.SpecificationDefinition{}, &gen.SpecificationValue{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			if err := db.Exec("DROP INDEX IF EXISTS `" + index + "`").Error; err != nil {
				t.Fatal(err)
			}
		}
	}
	cfg := resolverTestConfig()
	seedResolverPermission(t, db, "permission-spec-read", "hqProductCatalog:read", gen.PermissionScopeSystem)
	seedResolverPermission(t, db, "permission-spec-manage", "hqProductCatalog:manage", gen.PermissionScopeSystem)
	hq := seedResolverWorkspace(t, db, cfg, "hq", gen.OrganizationTypeHeadquarters, gen.RoleKindHqSuperAdmin)
	franchise := seedResolverWorkspace(t, db, cfg, "franchise", gen.OrganizationTypeFranchise, gen.RoleKindFranchiseOwner)
	for _, item := range []any{
		&gen.SpecificationDefinition{ID: "taste", OrganizationID: "hq-organization", Name: "口味", Enabled: true},
		&gen.SpecificationValue{ID: "mild", SpecificationID: "taste", Name: "微辣", Enabled: true},
		&gen.SpecificationValue{ID: "hot", SpecificationID: "taste", Name: "特辣", Enabled: true},
	} {
		if err := db.Create(item).Error; err != nil {
			t.Fatal(err)
		}
	}
	handler := governanceHandler(db, cfg)
	mutate := `mutation { hqReorderSpecificationValues(specificationId:"taste",orderedIds:["hot","mild"]) }`
	denied := executeResolverGraphQL(t, handler, franchise, mutate)
	assertGraphQLErrorCode(t, denied, auth.CodePermissionDenied)
	updated := executeResolverGraphQL(t, handler, hq, mutate)
	assertProductCatalogResponse(t, updated.Body.String(), `"hqReorderSpecificationValues":true`)
	listed := executeResolverGraphQL(t, handler, hq, `query { hqSpecificationValues(specificationId:"taste",page:1,perPage:20){data{id}} }`)
	assertProductCatalogResponse(t, listed.Body.String(), `"data":[{"id":"hot"},{"id":"mild"}]`)
	generated := executeResolverGraphQL(t, handler, hq, `mutation { updateSpecificationValue(id:"hot",input:{weight:9}){id weight} }`)
	assertGraphQLErrorCode(t, generated, auth.CodePermissionDenied)
}

func TestGraphQLFranchiseCatalogSelectionIsStoreScoped(t *testing.T) {
	db := openResolverTestDB(t)
	for _, model := range []any{&gen.ProductCategory{}, &gen.Product{}, &gen.ProductSku{}, &gen.ProductPackage{}, &gen.StoreListing{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	cfg := resolverTestConfig()
	seedResolverPermission(t, db, "permission-franchise-product-read", "franchiseProduct:read", gen.PermissionScopeTenant)
	seedResolverPermission(t, db, "permission-franchise-product-manage", "franchiseProduct:manage", gen.PermissionScopeTenant)
	seedResolverWorkspace(t, db, cfg, "hq", gen.OrganizationTypeHeadquarters, gen.RoleKindHqSuperAdmin)
	a := seedResolverWorkspace(t, db, cfg, "franchise-a", gen.OrganizationTypeFranchise, gen.RoleKindFranchiseOwner)
	b := seedResolverWorkspace(t, db, cfg, "franchise-b", gen.OrganizationTypeFranchise, gen.RoleKindFranchiseOwner)
	for _, roleID := range []string{"franchise-a-role", "franchise-b-role"} {
		for _, permissionID := range []string{"permission-franchise-product-read", "permission-franchise-product-manage"} {
			if err := db.Create(&resolverPermissionRole{PermissionID: permissionID, OperatorRoleID: roleID}).Error; err != nil {
				t.Fatal(err)
			}
		}
	}
	seedGraphQLCatalogStores(t, db)
	handler := governanceHandler(db, cfg)
	selected := executeResolverGraphQL(t, handler, a, `mutation { franchiseSetStoreListing(storeId:"store-a", skuId:"sku", enabled:true) { id enabled } }`)
	assertProductCatalogResponse(t, selected.Body.String(), `"enabled":true`)
	owned := executeResolverGraphQL(t, handler, a, `query { franchiseCatalog(storeId:"store-a",page:1,perPage:20){data{id productName categoryId listingId listingEnabled packages{id name}}} }`)
	assertProductCatalogResponse(t, owned.Body.String(), `"listingEnabled":true`)
	assertProductCatalogResponse(t, owned.Body.String(), `"productName":"Food"`)
	assertProductCatalogResponse(t, owned.Body.String(), `"categoryId":"category"`)
	foreign := executeResolverGraphQL(t, handler, b, `query { franchiseCatalog(storeId:"store-a",page:1,perPage:20){data{id listingId}} }`)
	assertGraphQLErrorCode(t, foreign, auth.CodePermissionDenied)
	unselected := executeResolverGraphQL(t, handler, b, `query { franchiseCatalog(storeId:"store-b",page:1,perPage:20){data{id listingId listingEnabled}} }`)
	assertProductCatalogResponse(t, unselected.Body.String(), `"listingId":null`)
}

func assertProductCatalogResponse(t *testing.T, body, expected string) {
	t.Helper()
	if strings.Contains(body, `"errors"`) || !strings.Contains(body, expected) {
		t.Fatalf("catalog response: %s, want %s", body, expected)
	}
}

func seedGraphQLCatalogStores(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, item := range []any{
		&gen.Store{ID: "store-a", Code: "SA", Name: "Store A", OrganizationID: "franchise-a-organization", Lifecycle: gen.StoreLifecycleActive},
		&gen.Store{ID: "store-b", Code: "SB", Name: "Store B", OrganizationID: "franchise-b-organization", Lifecycle: gen.StoreLifecycleActive},
		&gen.ProductCategory{ID: "category", Name: "Food", OrganizationID: "hq-organization", Enabled: true},
		&gen.Product{ID: "product", Name: "Food", OrganizationID: "hq-organization", CategoryID: "category", Enabled: true},
		&gen.ProductSku{ID: "sku", Name: "Spicy", ProductID: "product", Enabled: true},
		&gen.ProductPackage{ID: "piece", Name: "Piece", SkuID: "sku", PackageSetVersion: 1, Enabled: true},
	} {
		if err := db.Create(item).Error; err != nil {
			t.Fatal(err)
		}
	}
}
