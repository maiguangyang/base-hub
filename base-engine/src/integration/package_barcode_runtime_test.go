package integration_test

import (
	"context"
	"iter"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

type barcodeWriteModel struct {
	tool, operation string
	readTool        string
	readArgs        map[string]any
	args            map[string]any
	targets         []string
	calls           int
}

func (*barcodeWriteModel) Name() string { return "package-barcode-approval" }
func (script *barcodeWriteModel) GenerateContent(_ context.Context, _ *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		script.calls++
		yield(&model.LLMResponse{Content: &genai.Content{Role: genai.RoleModel, Parts: []*genai.Part{script.part()}}, TurnComplete: true, UsageMetadata: &genai.GenerateContentResponseUsageMetadata{PromptTokenCount: 20, CandidatesTokenCount: 3, TotalTokenCount: 23}}, nil)
	}
}
func (script *barcodeWriteModel) part() *genai.Part {
	switch script.calls {
	case 1:
		if script.readTool == "" {
			script.readTool = "HqProductPackages"
			script.readArgs = map[string]any{"skuId": "barcode-sku", "page": 1, "perPage": 20}
		}
		return gapFunctionCall("select-read", "select_tools", map[string]any{"names": []string{script.readTool}})
	case 2:
		return gapFunctionCall("read", script.readTool, script.readArgs)
	case 3:
		return gapFunctionCall("plan", "propose_plan", map[string]any{"steps": []any{map[string]any{
			"operationId": script.operation, "toolId": script.tool, "arguments": script.args, "targetIds": script.targets, "maxCalls": 1, "sequence": 1}}})
	case 4:
		return genai.NewPartFromText("Preview ready")
	case 5:
		return gapFunctionCall("select-write", "select_tools", map[string]any{"names": []string{script.tool}})
	case 6:
		return gapFunctionCall("write", script.tool, script.args)
	default:
		return genai.NewPartFromText("Finished")
	}
}
func TestPackageBarcodeProtectedApproval(t *testing.T) {
	cases := []struct {
		name, tool, operation string
		args                  map[string]any
		targets               []string
		success               bool
	}{
		{"setter", "HqSetProductPackageBarcode", "graphql.mutation.hqSetProductPackageBarcode", map[string]any{"id": "barcode-package", "barcode": "0012345678905"}, []string{"barcode-package"}, true},
		{"create-conflict", "HqCreateProductPackage", "graphql.mutation.hqCreateProductPackage", map[string]any{"input": map[string]any{"skuId": "barcode-sku", "name": "Box", "containsPackageId": "barcode-draft", "containsQuantity": 5, "barcode": "occupied", "packageSetVersion": 2}}, []string{}, false},
		{"publish-conflict", "HqPublishProductPackageSet", "graphql.mutation.hqPublishProductPackageSet", map[string]any{"skuId": "barcode-sku", "version": 2}, []string{"barcode-sku"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := barcodeRuntimeFixture(t)
			if _, err := runProtectedGapTool(t, fixture, "session-hq", "HqProductPackages", map[string]any{"skuId": "barcode-sku", "page": 1, "perPage": 20}); err != nil {
				t.Fatal(err)
			}
			script := &barcodeWriteModel{tool: tc.tool, operation: tc.operation, args: tc.args, targets: tc.targets}
			router := gapDeleteServiceRouter(t, fixture, script)
			preview := callAIIntegrationWithSession(t, router, fixture, "session-hq", "/api/ai/preview", `{"prompt":"Update the selected package barcode"}`)
			assertRuntimeBarcode(t, fixture, "original")
			assertAuditCount(t, fixture, "hqProductCatalog:manage", 0)
			token := previewTokenFromSSE(t, preview.Body.String())
			run := callAIIntegrationWithSession(t, router, fixture, "session-hq", "/api/ai/run", `{"previewToken":"`+token+`"}`)
			if script.calls < 6 {
				t.Fatalf("write not attempted: %s", run.Body.String())
			}
			if tc.success {
				if !strings.Contains(run.Body.String(), `"status":"SUCCESS"`) {
					t.Fatalf("approved write failed: %s", run.Body.String())
				}
				assertRuntimeBarcode(t, fixture, "0012345678905")
				assertAuditCount(t, fixture, "hqProductCatalog:manage", 1)
			} else {
				if !strings.Contains(run.Body.String(), `"status":"FAILED"`) {
					t.Fatalf("missing conflict: %s", run.Body.String())
				}
				assertRuntimeBarcode(t, fixture, "original")
				assertAuditCount(t, fixture, "hqProductCatalog:manage", 0)
				assertBarcodeConflictResponse(t, fixture, tc.tool)
			}
		})
	}
}
func assertRuntimeBarcode(t *testing.T, fixture *securityFixture, want string) {
	t.Helper()
	var pack gen.ProductPackage
	if err := fixture.db.First(&pack, "id = ?", "barcode-package").Error; err != nil {
		t.Fatal(err)
	}
	if pack.Barcode == nil || *pack.Barcode != want || !pack.Enabled {
		t.Fatalf("package changed: %+v", pack)
	}
	var sku gen.ProductSku
	if err := fixture.db.First(&sku, "id = ?", "barcode-sku").Error; err != nil {
		t.Fatal(err)
	}
	if sku.PublishedPackageSetVersion != 1 {
		t.Fatalf("published version changed: %+v", sku)
	}
}
func barcodeRuntimeFixture(t *testing.T) *securityFixture {
	t.Helper()
	fixture := newSecurityFixture(t)
	for _, item := range []any{&gen.ProductCategory{}, &gen.Product{}, &gen.ProductSku{}, &gen.ProductPackage{}, &gen.StoreListing{}, &gen.StorePackageOffer{}, &gen.StorePromotion{}, &gen.StorePromotionTarget{}} {
		if err := fixture.db.AutoMigrate(item); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			fixture.db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	fixture.createAll([]gen.Permission{
		{ID: "barcode-read", Name: "read", Action: "hqProductCatalog:read", Module: "test", Scope: gen.PermissionScopeSystem},
		{ID: "barcode-manage", Name: "manage", Action: "hqProductCatalog:manage", Module: "test", Scope: gen.PermissionScopeSystem}},
		[]permissionRole{{"barcode-read", "role-hq"}, {"barcode-manage", "role-hq"}},
		[]gen.ProductCategory{{ID: "barcode-category", OrganizationID: "org-hq", Name: "Food", Enabled: true}},
		[]gen.Product{{ID: "barcode-product", OrganizationID: "org-hq", CategoryID: "barcode-category", Name: "Noodles", Enabled: true}},
		[]gen.ProductSku{{ID: "barcode-sku", ProductID: "barcode-product", Name: "Spicy", PublishedPackageSetVersion: 1, Enabled: true}, {ID: "barcode-other-sku", ProductID: "barcode-product", Name: "Other", PublishedPackageSetVersion: 1, Enabled: true}},
		[]gen.ProductPackage{
			{ID: "barcode-package", SkuID: "barcode-sku", Name: "Bag", Barcode: stringPointer("original"), PackageSetVersion: 1, Enabled: true},
			{ID: "barcode-other", SkuID: "barcode-other-sku", Name: "Bag", Barcode: stringPointer("occupied"), PackageSetVersion: 1, Enabled: true},
			{ID: "barcode-draft", SkuID: "barcode-sku", Name: "Bag", Barcode: stringPointer("occupied"), PackageSetVersion: 2, Enabled: false}})
	// GORM's enabled default must not activate the deliberately conflicting draft.
	if err := fixture.db.Model(&gen.ProductPackage{}).Where("id = ?", "barcode-draft").Update("enabled", false).Error; err != nil {
		t.Fatal(err)
	}
	return fixture
}

func assertBarcodeConflictResponse(t *testing.T, fixture *securityFixture, tool string) {
	document := `mutation {hqPublishProductPackageSet(skuId:"barcode-sku",version:2){id}}`
	if tool == "HqCreateProductPackage" {
		document = `mutation {hqCreateProductPackage(input:{skuId:"barcode-sku",name:"Box",containsPackageId:"barcode-draft",containsQuantity:5,barcode:"occupied",packageSetVersion:2}){id}}`
	}
	assertCode(t, fixture.execute("session-hq", document), auth.CodeConflict)
}
