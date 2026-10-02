package integration_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/dbup"
)

func stockBarcodeRuntimeFixture(t *testing.T) *securityFixture {
	t.Helper()
	fixture := barcodeRuntimeFixture(t)
	for _, item := range []any{&gen.ProductBrand{}, &gen.SpecificationDefinition{}, &gen.SpecificationValue{}, &gen.ProductSpecificationChoice{}, &gen.ProductSkuSpecificationValue{}, &gen.StorePriceRevision{}, &gen.StoreInventoryBatch{}, &gen.StoreStockBalance{}, &gen.StoreStockMovement{}, &gen.StoreStocktake{}, &gen.StoreStocktakeLine{}} {
		if err := fixture.db.AutoMigrate(item); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			fixture.db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	if err := dbup.EnsureProductIndexes(fixture.db); err != nil {
		t.Fatal(err)
	}
	fixture.createAll([]gen.Permission{
		{ID: "stock-barcode-read", Name: "read", Action: "franchiseStock:read", Module: "test", Scope: gen.PermissionScopeTenant},
		{ID: "stock-barcode-manage", Name: "manage", Action: "franchiseStock:manage", Module: "test", Scope: gen.PermissionScopeTenant}},
		[]permissionRole{{"stock-barcode-read", "role-owner-a"}, {"stock-barcode-read", "role-custom-a"}, {"stock-barcode-manage", "role-owner-a"}},
		[]gen.StoreListing{{ID: "barcode-listing", SkuID: "barcode-sku", StoreID: "store-a", Enabled: true}})
	return fixture
}
func TestStockBarcodeProtectedReadScopeAndPermissions(t *testing.T) {
	fixture := stockBarcodeRuntimeFixture(t)
	result, err := runProtectedGapTool(t, fixture, "session-a-2", "FranchiseStockPackageByBarcode", map[string]any{"storeId": "store-a", "barcode": "original"})
	if err != nil {
		t.Fatal(err)
	}
	root, ok := result["franchiseStockPackageByBarcode"].(map[string]any)
	if !ok || root["status"] != "READY" {
		t.Fatalf("read-only recognition: %+v", result)
	}
	if _, err := runProtectedGapTool(t, fixture, "session-a-2", "FranchiseStockPackageByBarcode", map[string]any{"storeId": "store-b", "barcode": "original"}); err == nil {
		t.Fatal("foreign store leaked")
	}
	assertCode(t, fixture.execute("session-a-2", `mutation {franchiseReceiveStock(input:{storeId:"store-a",listingId:"barcode-listing",packageId:"barcode-package",barcode:"original",quantity:2,requestKey:"read-only"}){id}}`), auth.CodePermissionDenied)
	script := stockReceiptModel("original")
	router := gapDeleteServiceRouter(t, fixture, script)
	preview := callAIIntegrationWithSession(t, router, fixture, "session-a-2", "/api/ai/preview", `{"prompt":"Receive scanned packages"}`)
	if strings.Contains(preview.Body.String(), "previewToken") {
		t.Fatalf("read-only operator approved receipt: %s", preview.Body.String())
	}
	assertAuditCount(t, fixture, "franchiseStock:manage", 0)
}
func stockReceiptModel(barcode string) *barcodeWriteModel {
	return &barcodeWriteModel{tool: "FranchiseReceiveStock", operation: "graphql.mutation.franchiseReceiveStock",
		readTool: "FranchiseStockPackageByBarcode", readArgs: map[string]any{"storeId": "store-a", "barcode": "original"}, targets: []string{},
		args: map[string]any{"input": map[string]any{"storeId": "store-a", "listingId": "barcode-listing", "packageId": "barcode-package", "barcode": barcode, "quantity": 2}}}
}
func TestStockBarcodeApprovedReceiptRechecksBarcode(t *testing.T) {
	for _, code := range []string{"original", "stale"} {
		t.Run(code, func(t *testing.T) {
			fixture := stockBarcodeRuntimeFixture(t)
			script := stockReceiptModel(code)
			router := gapDeleteServiceRouter(t, fixture, script)
			preview := callAIIntegrationWithSession(t, router, fixture, "session-a-1", "/api/ai/preview", `{"prompt":"Receive two scanned packages"}`)
			token := previewTokenFromSSE(t, preview.Body.String())
			assertAuditCount(t, fixture, "franchiseStock:manage", 0)
			run := callAIIntegrationWithSession(t, router, fixture, "session-a-1", "/api/ai/run", `{"previewToken":"`+token+`"}`)
			if script.calls < 6 {
				t.Fatalf("receipt never attempted: %s", run.Body.String())
			}
			wantStatus, wantCount := "FAILED", int64(0)
			if code == "original" {
				wantStatus, wantCount = "SUCCESS", 1
			}
			if !strings.Contains(run.Body.String(), `"status":"`+wantStatus+`"`) {
				t.Fatalf("receipt outcome: %s", run.Body.String())
			}
			assertAuditCount(t, fixture, "franchiseStock:manage", wantCount)
			var count int64
			if err := fixture.db.Model(&gen.StoreStockMovement{}).Count(&count).Error; err != nil {
				t.Fatal(err)
			}
			if count != wantCount {
				t.Fatalf("receipt movements %d want %d", count, wantCount)
			}
			if code == "original" {
				assertBarcodeReceiptReplay(t, fixture)
			}
		})
	}
}
func assertBarcodeReceiptReplay(t *testing.T, fixture *securityFixture) {
	t.Helper()
	if err := fixture.db.Model(&gen.ProductPackage{}).Where("id = ?", "barcode-package").Updates(map[string]any{"barcode": "replacement", "enabled": false}).Error; err != nil {
		t.Fatal(err)
	}
	var movement gen.StoreStockMovement
	if err := fixture.db.First(&movement).Error; err != nil {
		t.Fatal(err)
	}
	var batch gen.StoreInventoryBatch
	if err := fixture.db.First(&batch, "id = ?", movement.BatchID).Error; err != nil {
		t.Fatal(err)
	}
	zone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}
	wantNumber := "B" + time.Now().In(zone).Format("20060102") + "-0001"
	if batch.BatchNumber != wantNumber {
		t.Fatalf("approved AI receipt number %s want %s", batch.BatchNumber, wantNumber)
	}
	result := fixture.execute("session-a-1", fmt.Sprintf(`mutation {franchiseReceiveStock(input:{storeId:"store-a",listingId:"barcode-listing",packageId:"barcode-package",barcode:"original",quantity:2,requestKey:"%s"}){id}}`, movement.RequestKey))
	if len(result.Errors) > 0 {
		t.Fatalf("replay failed: %s", result.Body)
	}
	var count int64
	if err := fixture.db.Model(&gen.StoreStockMovement{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("replay duplicated: %d", count)
	}
}
func stockBatchProtectedFixture(t *testing.T) *securityFixture {
	t.Helper()
	fixture := stockBarcodeRuntimeFixture(t)
	fixture.createAll(
		[]gen.StoreListing{
			{ID: "batch-other-listing", StoreID: "store-a", SkuID: "barcode-other-sku", Enabled: true},
			{ID: "batch-foreign-listing", StoreID: "store-b", SkuID: "barcode-sku", Enabled: true}},
		[]gen.StoreInventoryBatch{
			{ID: "batch-a-1", ListingID: "barcode-listing", BatchNumber: "B-1"},
			{ID: "batch-a-2", ListingID: "batch-other-listing", BatchNumber: "B-2"},
			{ID: "batch-b", ListingID: "batch-foreign-listing", BatchNumber: "FOREIGN"}},
		[]gen.StoreStockBalance{{ID: "batch-a-balance", BatchID: "batch-a-1", PackageID: "barcode-package", Quantity: 2, Version: 1}},
	)
	return fixture
}

func TestStockBatchesProtectedToolStoreWideResult(t *testing.T) {
	fixture := stockBatchProtectedFixture(t)
	args := map[string]any{"storeId": "store-a", "page": 1, "perPage": 20}
	result, err := runProtectedGapTool(t, fixture, "session-a-2", "FranchiseStockBatches", args)
	if err != nil {
		t.Fatal(err)
	}
	direct := fixture.execute("session-a-2", `query {franchiseStockBatches(storeId:"store-a",page:1,perPage:20){data{id listingId skuId productName skuName batchNumber producedAt expiresAt sellable packages{id name enabled containsPackageId containsQuantity packageSetVersion} balances{packageId quantity}} total currentPage perPage}}`)
	if len(direct.Errors) != 0 {
		t.Fatalf("direct batch query failed: %s", direct.Body)
	}
	var expected any
	if err := json.Unmarshal(direct.Data["franchiseStockBatches"], &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result["franchiseStockBatches"], expected) {
		t.Fatalf("tool projection differs: %#v versus %#v", result, expected)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"batch-a-1", "batch-a-2", "barcode-package", `"productName":"Noodles"`, `"skuName":"Spicy"`, `"skuName":"Other"`, `"quantity":2`} {
		if !strings.Contains(string(encoded), id) {
			t.Fatalf("missing %s in %s", id, encoded)
		}
	}
	if strings.Contains(string(encoded), "batch-b") {
		t.Fatalf("foreign batch leaked: %s", encoded)
	}
	assertBatchReadDidNotWrite(t, fixture)
}

func assertBatchReadDidNotWrite(t *testing.T, fixture *securityFixture) {
	t.Helper()
	var movements int64
	err := fixture.db.Model(&gen.StoreStockMovement{}).Count(&movements).Error
	if err != nil || movements != 0 {
		t.Fatalf("read wrote stock movements: count=%d error=%v", movements, err)
	}
}

func TestStockBatchesProtectedToolRejectsScopePermissionAndInvalidInput(t *testing.T) {
	fixture := stockBatchProtectedFixture(t)
	for _, args := range []map[string]any{
		{"storeId": "store-b", "page": 1, "perPage": 20},
		{"storeId": "store-a", "page": "wrong-type", "perPage": 20},
		{"storeId": "store-a", "page": 1, "perPage": 20, "unexpected": true},
	} {
		if _, err := runProtectedGapTool(t, fixture, "session-a-2", "FranchiseStockBatches", args); err == nil {
			t.Fatalf("invalid or foreign input accepted: %#v", args)
		}
	}
	if err := fixture.db.Where("permission_id = ? AND operator_role_id = ?", "stock-barcode-read", "role-custom-a").Delete(&permissionRole{}).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := runProtectedGapTool(t, fixture, "session-a-2", "FranchiseStockBatches", map[string]any{"storeId": "store-a", "page": 1, "perPage": 20}); err == nil {
		t.Fatal("read without stock permission succeeded")
	}
	assertCode(t, fixture.execute("session-a-2", `query {franchiseStockBatches(storeId:"store-a",page:1,perPage:20){total}}`), auth.CodePermissionDenied)
}
