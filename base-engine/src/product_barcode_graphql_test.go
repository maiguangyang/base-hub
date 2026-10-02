package src_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/ai"
)

func TestProductBarcodeGraphQLAndProtectedParity(t *testing.T) {
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
	seedResolverPermission(t, db, "barcode-manage", "hqProductCatalog:manage", gen.PermissionScopeSystem)
	hq := seedResolverWorkspace(t, db, cfg, "hq", gen.OrganizationTypeHeadquarters, gen.RoleKindHqSuperAdmin)
	franchise := seedResolverWorkspace(t, db, cfg, "franchise", gen.OrganizationTypeFranchise, gen.RoleKindFranchiseOwner)
	code := "original"
	for _, item := range []any{
		&gen.ProductCategory{ID: "barcode-category", OrganizationID: "hq-organization", Name: "Food", Enabled: true},
		&gen.Product{ID: "barcode-product", OrganizationID: "hq-organization", CategoryID: "barcode-category", Name: "Noodles", Enabled: true},
		&gen.ProductSku{ID: "barcode-sku", ProductID: "barcode-product", Name: "Spicy", PublishedPackageSetVersion: 1, Enabled: true},
		&gen.ProductPackage{ID: "barcode-package", SkuID: "barcode-sku", Name: "Bag", Barcode: &code, PackageSetVersion: 1, Enabled: true},
	} {
		if err := db.Select("*").Create(item).Error; err != nil {
			t.Fatal(err)
		}
	}
	handler := governanceHandler(db, cfg)
	document := `mutation { hqSetProductPackageBarcode(id:"barcode-package",barcode:"0012345678905") {id barcode} }`
	denied := executeResolverGraphQL(t, handler, franchise, document)
	assertGraphQLErrorCode(t, denied, auth.CodePermissionDenied)
	updated := executeResolverGraphQL(t, handler, hq, document)
	assertProductCatalogResponse(t, updated.Body.String(), `"barcode":"0012345678905"`)
	request := ai.FixedRequest{Method: http.MethodPost, Path: "/graphql", Document: document,
		Variables: json.RawMessage(`{}`), Cookie: hq.String(), Origin: "https://admin.example.com"}
	response, err := ai.ProtectedCall(context.Background(), handler, request)
	if err != nil {
		t.Fatal(err)
	}
	assertProductCatalogResponse(t, string(response.Body), `"barcode":"0012345678905"`)
	assertGeneratedAndClearedBarcode(t, handler, hq, franchise, request)
	bad := executeResolverGraphQL(t, handler, hq, `mutation { hqSetProductPackageBarcode(id:"barcode-package",barcode:"`+strings.Repeat("a", 65)+`") {id} }`)
	assertGraphQLErrorCode(t, bad, auth.CodeValidationFailed)
}

func assertGeneratedAndClearedBarcode(t *testing.T, handler http.Handler, hq, franchise *http.Cookie, request ai.FixedRequest) {
	t.Helper()
	cleared := executeResolverGraphQL(t, handler, hq, `mutation { hqSetProductPackageBarcode(id:"barcode-package",barcode:"") {id barcode} }`)
	assertProductCatalogResponse(t, cleared.Body.String(), `"barcode":null`)
	generate := `query { hqGenerateProductPackageBarcode }`
	generated := executeResolverGraphQL(t, handler, hq, generate)
	assertProductCatalogResponse(t, generated.Body.String(), `"hqGenerateProductPackageBarcode":"KH`)
	assertGraphQLErrorCode(t, executeResolverGraphQL(t, handler, franchise, generate), auth.CodePermissionDenied)
	request.Document = generate
	response, err := ai.ProtectedCall(context.Background(), handler, request)
	if err != nil {
		t.Fatal(err)
	}
	assertProductCatalogResponse(t, string(response.Body), `"hqGenerateProductPackageBarcode":"KH`)
}
