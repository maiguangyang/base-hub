package storemerchandising

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/dbup"
	"base-engine/src/services/audit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func mysqlStockFixture(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal) {
	t.Helper()
	rootDSN := os.Getenv("PRODUCT_TEST_MYSQL_DSN")
	if rootDSN == "" {
		t.Skip("PRODUCT_TEST_MYSQL_DSN is required for the disposable MySQL test")
	}
	if !strings.HasSuffix(rootDSN, "/") {
		t.Fatal("PRODUCT_TEST_MYSQL_DSN must end in / with no database selected")
	}
	admin, err := gorm.Open(mysql.Open(rootDSN+"mysql?parseTime=true"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	database := fmt.Sprintf("product_stock_test_%d", time.Now().UnixNano())
	if err := admin.Exec("CREATE DATABASE `" + database + "`").Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := admin.Exec("DROP DATABASE `" + database + "`").Error; err != nil {
			t.Errorf("drop database: %v", err)
		}
	})
	db, err := gorm.Open(mysql.Open(rootDSN+database+"?parseTime=true"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range []any{&gen.Organization{}, &gen.Store{}, &gen.ProductBrand{}, &gen.SpecificationDefinition{}, &gen.SpecificationValue{}, &gen.ProductSpecificationChoice{}, &gen.ProductSkuSpecificationValue{}, &gen.ProductCategory{}, &gen.Product{}, &gen.ProductSku{},
		&gen.ProductPackage{}, &gen.StoreListing{}, &gen.StorePackageOffer{}, &gen.StorePriceRevision{}, &gen.StoreInventoryBatch{}, &gen.StoreStockBalance{},
		&gen.StoreStockMovement{}, &gen.StoreStocktake{}, &gen.StoreStocktakeLine{}, &gen.StorePromotion{}, &gen.AuditLog{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
	}
	if err := dbup.EnsureProductIndexes(db); err != nil {
		t.Fatal(err)
	}
	seedMySQLStock(t, db)
	orgID := "org"
	principal := &auth.WorkspacePrincipal{AccountID: "operator", SessionID: "session", WorkspaceType: auth.WorkspaceTypeFranchise,
		OrganizationID: &orgID, Permissions: map[string]struct{}{"franchiseProduct:manage": {}, "franchiseStock:manage": {}, "franchiseStock:read": {}}, StoreIDs: map[string]struct{}{"store": {}}}
	return NewService(db, audit.NewService()), db, principal
}

func seedMySQLStock(t *testing.T, db *gorm.DB) {
	t.Helper()
	child := "piece"
	factor := int64(12)
	for _, item := range []any{
		&gen.Organization{ID: "hq", Code: "HQ", Name: "HQ", Type: gen.OrganizationTypeHeadquarters, Status: gen.OrganizationStatusActive},
		&gen.Organization{ID: "org", Code: "ORG", Name: "ORG", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive},
		&gen.Store{ID: "store", Code: "S", Name: "Store", OrganizationID: "org", Lifecycle: gen.StoreLifecycleActive},
		&gen.ProductCategory{ID: "category", OrganizationID: "hq", Name: "Food", Enabled: true},
		&gen.Product{ID: "product", Name: "Food", OrganizationID: "hq", CategoryID: "category", Enabled: true},
		&gen.ProductSku{ID: "sku", Name: "Large", ProductID: "product", Enabled: true},
		&gen.ProductPackage{ID: child, Name: "Piece", SkuID: "sku", PackageSetVersion: 1, Enabled: true},
		&gen.ProductPackage{ID: "bag", Name: "Bag", SkuID: "sku", PackageSetVersion: 1, ContainsPackageID: &child, ContainsQuantity: &factor, Enabled: true},
	} {
		if err := db.Create(item).Error; err != nil {
			t.Fatal(err)
		}
	}
}

func TestStockMySQLConcurrentUnpack(t *testing.T) {
	service, db, principal := mysqlStockFixture(t)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := service.ReceiveStock(ctx, principal, StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "bag", Quantity: 1, BatchNumber: "B1", RequestKey: "receive"})
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for _, key := range []string{"unpack-one", "unpack-two"} {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			_, err := service.UnpackStock(ctx, principal, StockUnpack{StoreID: "store", BatchID: receipt.BatchID, SourcePackageID: "bag", Quantity: 1, RequestKey: key})
			results <- err
		}(key)
	}
	wg.Wait()
	close(results)
	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("successful unpacks = %d, want 1", successes)
	}
	assertMySQLUnpackBalances(t, db, receipt.BatchID)
}

func assertMySQLUnpackBalances(t *testing.T, db *gorm.DB, batchID string) {
	t.Helper()
	var balances []gen.StoreStockBalance
	if err := db.Where("batch_id = ?", batchID).Find(&balances).Error; err != nil {
		t.Fatal(err)
	}
	quantities := map[string]int64{}
	for _, balance := range balances {
		quantities[balance.PackageID] = balance.Quantity
	}
	if quantities["bag"] != 0 || quantities["piece"] != 12 {
		t.Fatalf("balances = %+v", quantities)
	}
}

func TestStockBarcodeMySQLIndex(t *testing.T) {
	service, db, principal := mysqlStockFixture(t)
	seedBarcodeIndexSample(t, db)
	var plans []struct {
		Key          *string
		PossibleKeys *string `gorm:"column:possible_keys"`
		Type         string
	}
	merchandisingNoError(t, db.Raw("EXPLAIN FORMAT=TRADITIONAL SELECT * FROM product_packages WHERE active_barcode = ? LIMIT 2", "0012345678905").Scan(&plans).Error)
	if len(plans) != 1 || plans[0].Key == nil || *plans[0].Key != "uidx_product_package_active_barcode" {
		t.Fatalf("actual EXPLAIN: %+v", plans)
	}
	t.Logf("actual index=%s access=%s", *plans[0].Key, plans[0].Type)
	listing, err := service.SetListing(context.Background(), principal, "store", "sku", true)
	merchandisingNoError(t, err)
	result, err := service.StockPackageByBarcode(context.Background(), principal, "store", "0012345678905")
	merchandisingNoError(t, err)
	if result.Target == nil || result.Target.ConversionChain[2].BaseQuantity != 40 {
		t.Fatalf("conversion: %+v", result)
	}
	input := StockReceipt{StoreID: "store", ListingID: listing.ID, PackageID: "box", Barcode: barcodeString("0012345678905"), Quantity: 2, RequestKey: "two-boxes"}
	movement, err := service.ReceiveStock(context.Background(), principal, input)
	merchandisingNoError(t, err)
	var balances []gen.StoreStockBalance
	merchandisingNoError(t, db.Where("batch_id = ?", movement.BatchID).Find(&balances).Error)
	if len(balances) != 1 || balances[0].PackageID != "box" || balances[0].Quantity != 2 {
		t.Fatalf("package quantities: %+v", balances)
	}
}
func seedBarcodeIndexSample(t *testing.T, db *gorm.DB) {
	t.Helper()
	merchandisingNoError(t, db.Model(&gen.ProductSku{}).Where("id = ?", "sku").Update("published_package_set_version", 1).Error)
	merchandisingNoError(t, db.Model(&gen.ProductPackage{}).Where("id = ?", "bag").Update("contains_quantity", 5).Error)
	child, factor := "bag", int64(8)
	merchandisingNoError(t, db.Create(&gen.ProductPackage{ID: "box", Name: "Box", SkuID: "sku", PackageSetVersion: 1, Enabled: true, ContainsPackageID: &child, ContainsQuantity: &factor, Barcode: barcodeString("0012345678905")}).Error)
	child, factor = "piece", 2
	samples := make([]gen.ProductPackage, 1000)
	for i := range samples {
		samples[i] = gen.ProductPackage{ID: fmt.Sprintf("sample-%d", i), SkuID: "sku", Name: "Sample", PackageSetVersion: 1, Enabled: true, ContainsPackageID: &child, ContainsQuantity: &factor, Barcode: barcodeString(fmt.Sprintf("sample-code-%d", i))}
	}
	merchandisingNoError(t, db.CreateInBatches(samples, 100).Error)
	merchandisingNoError(t, db.Exec("ANALYZE TABLE product_packages").Error)
}
