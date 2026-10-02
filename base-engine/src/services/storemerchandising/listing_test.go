package storemerchandising

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func fixture(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range []any{&gen.Organization{}, &gen.Store{}, &gen.ProductCategory{}, &gen.Product{}, &gen.ProductSku{}, &gen.ProductPackage{}, &gen.StoreListing{}, &gen.StorePackageOffer{}, &gen.StorePriceRevision{}, &gen.StoreInventoryBatch{}, &gen.StoreStockBalance{}, &gen.StoreStockMovement{}, &gen.StoreStocktake{}, &gen.StoreStocktakeLine{}, &gen.StorePromotion{}, &gen.StorePromotionTarget{}, &gen.CustomerBenefitPolicy{}, &gen.AuditLog{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	hq := gen.Organization{ID: "hq", Code: "HQ", Name: "HQ", Type: gen.OrganizationTypeHeadquarters, Status: gen.OrganizationStatusActive}
	franchise := gen.Organization{ID: "org", Code: "ORG", Name: "ORG", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}
	for _, org := range []*gen.Organization{&hq, &franchise} {
		if err := db.Create(org).Error; err != nil {
			t.Fatal(err)
		}
	}
	store := gen.Store{ID: "store", Code: "S", Name: "Store", OrganizationID: "org", Lifecycle: gen.StoreLifecycleActive}
	if err := db.Create(&store).Error; err != nil {
		t.Fatal(err)
	}
	category := gen.ProductCategory{ID: "category", OrganizationID: "hq", Name: "Food", Enabled: true}
	product := gen.Product{ID: "product", Name: "Food", OrganizationID: "hq", CategoryID: "category", Enabled: true}
	sku := gen.ProductSku{ID: "sku", Name: "Large", ProductID: "product", Enabled: true}
	pack := gen.ProductPackage{ID: "piece", Name: "Piece", SkuID: "sku", PackageSetVersion: 1, Enabled: true}
	for _, item := range []any{&category, &product, &sku, &pack} {
		if err := db.Create(item).Error; err != nil {
			t.Fatal(err)
		}
	}
	orgID := "org"
	principal := &auth.WorkspacePrincipal{AccountID: "operator", SessionID: "session", WorkspaceType: auth.WorkspaceTypeFranchise, OrganizationID: &orgID,
		Permissions: map[string]struct{}{"franchiseProduct:read": {}, "franchiseProduct:manage": {}, "franchiseStock:manage": {}, "franchiseStock:read": {}, "franchisePromotion:manage": {}, "franchisePromotion:read": {}}, StoreIDs: map[string]struct{}{"store": {}}}
	return NewService(db, audit.NewService()), db, principal
}

func merchandisingNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListingSelectionRespectsStoreScope(t *testing.T) {
	service, _, principal := fixture(t)
	ctx := context.Background()
	listing, err := service.SetListing(ctx, principal, "store", "sku", true)
	if err != nil {
		t.Fatal(err)
	}
	if !listing.Enabled || listing.StoreID != "store" || listing.SkuID != "sku" {
		t.Fatalf("listing = %+v", listing)
	}
	replayed, err := service.SetListing(ctx, principal, "store", "sku", true)
	if err != nil || replayed.ID != listing.ID {
		t.Fatalf("replay = %+v, %v", replayed, err)
	}
	principal.StoreIDs = map[string]struct{}{}
	if _, err := service.SetListing(ctx, principal, "store", "sku", false); auth.ErrorCode(err) != auth.CodeStoreScopeDenied {
		t.Fatalf("store scope: %v", err)
	}
}

func TestCatalogListingsReturnsOnlyAuthorizedStoreSelections(t *testing.T) {
	service, db, principal := fixture(t)
	ctx := context.Background()
	selected, err := service.SetListing(ctx, principal, "store", "sku", true)
	if err != nil {
		t.Fatal(err)
	}
	otherStore := gen.Store{ID: "other-store", Code: "OTHER", Name: "Other", OrganizationID: "org", Lifecycle: gen.StoreLifecycleActive}
	if err := db.Create(&otherStore).Error; err != nil {
		t.Fatal(err)
	}
	otherListing := gen.StoreListing{ID: "other-listing", StoreID: otherStore.ID, SkuID: "sku", Enabled: false}
	if err := db.Create(&otherListing).Error; err != nil {
		t.Fatal(err)
	}
	items, err := service.CatalogListings(ctx, principal, "store", []string{"sku"})
	if err != nil || items["sku"].ID != selected.ID {
		t.Fatalf("store selection: %+v, %v", items, err)
	}
	if _, err := service.CatalogListings(ctx, principal, otherStore.ID, []string{"sku"}); auth.ErrorCode(err) != auth.CodeStoreScopeDenied {
		t.Fatalf("other store scope: %v", err)
	}
}

func TestMerchandisingStoreDirectorySupportsFeatureOnlyRolesAndStoreScope(t *testing.T) {
	service, db, principal := fixture(t)
	other := gen.Store{ID: "other-store", Code: "OTHER", Name: "Other", OrganizationID: "org", Lifecycle: gen.StoreLifecycleActive}
	if err := db.Create(&other).Error; err != nil {
		t.Fatal(err)
	}
	principal.Permissions = map[string]struct{}{"franchiseStock:read": {}}
	page, err := service.ScopedStores(context.Background(), principal, 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Data) != 1 || page.Data[0].ID != "store" {
		t.Fatalf("scoped stores = %+v", page)
	}
	if _, err := service.Catalog(context.Background(), principal, "store", CatalogFilter{}, 1, 20); err != nil {
		t.Fatalf("stock role catalog: %v", err)
	}
	if _, err := service.Offers(context.Background(), principal, "store", "", 1, 20); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("stock role offers: %v", err)
	}
	principal.Permissions = map[string]struct{}{}
	if _, err := service.ScopedStores(context.Background(), principal, 1, 20); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("permissionless stores: %v", err)
	}
}

func TestCatalogCategoriesSupportsPromotionReadRole(t *testing.T) {
	service, _, principal := fixture(t)
	principal.Permissions = map[string]struct{}{"franchisePromotion:read": {}}
	page, err := service.CatalogCategories(context.Background(), principal, "store", 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Data) != 1 || page.Data[0].ID != "category" {
		t.Fatalf("promotion categories = %+v", page)
	}
}

func TestMerchandisingStoreDirectoryExcludesStoresNotReadyForTrading(t *testing.T) {
	service, db, principal := fixture(t)
	principal.AllStores = true
	for _, store := range []gen.Store{
		{ID: "draft-store", Code: "DRAFT", Name: "Draft", OrganizationID: "org", Lifecycle: gen.StoreLifecycleDraft},
		{ID: "rejected-store", Code: "REJECTED", Name: "Rejected", OrganizationID: "org", Lifecycle: gen.StoreLifecycleRejected},
	} {
		if err := db.Create(&store).Error; err != nil {
			t.Fatal(err)
		}
	}
	page, err := service.ScopedStores(context.Background(), principal, 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Data) != 1 || page.Data[0].ID != "store" {
		t.Fatalf("trading store directory = %+v", page)
	}
}

func TestSelectedSkuRemainsVisibleForHistoricalStockAfterHeadquartersDisablesIt(t *testing.T) {
	service, db, principal := fixture(t)
	ctx := context.Background()
	_, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	if err := db.Model(&gen.ProductSku{}).Where("id = ?", "sku").Update("enabled", false).Error; err != nil {
		t.Fatal(err)
	}
	page, err := service.Catalog(ctx, principal, "store", CatalogFilter{}, 1, 20)
	merchandisingNoError(t, err)
	if page.Total != 1 || page.Data[0].ID != "sku" {
		t.Fatalf("historical catalog = %+v", page)
	}
	packages, err := service.CatalogPackages(ctx, principal, "store", "sku")
	merchandisingNoError(t, err)
	if len(packages) != 1 || packages[0].ID != "piece" {
		t.Fatalf("historical packages = %+v", packages)
	}
	published, err := service.CatalogPublished(ctx, "sku")
	if err != nil || published {
		t.Fatalf("published = %v, %v", published, err)
	}
}

func TestFranchiseCatalogFiltersBeforePagination(t *testing.T) {
	service, db, principal := fixture(t)
	ctx := context.Background()
	_, err := service.SetListing(ctx, principal, "store", "sku", true)
	merchandisingNoError(t, err)
	if err := db.Create(&gen.ProductSku{ID: "other-sku", ProductID: "product", Name: "Small", Enabled: true}).Error; err != nil {
		t.Fatal(err)
	}
	term, categoryID, selected := "Large", "category", true
	page, err := service.Catalog(ctx, principal, "store", CatalogFilter{Q: &term, CategoryID: &categoryID, ListingEnabled: &selected}, 1, 1)
	merchandisingNoError(t, err)
	if page.Total != 1 || len(page.Data) != 1 || page.Data[0].ID != "sku" {
		t.Fatalf("filtered catalog = %+v", page)
	}
	selected = false
	page, err = service.Catalog(ctx, principal, "store", CatalogFilter{CategoryID: &categoryID, ListingEnabled: &selected}, 1, 1)
	merchandisingNoError(t, err)
	if page.Total != 1 || page.Data[0].ID != "other-sku" {
		t.Fatalf("unselected catalog = %+v", page)
	}
	categoryID = "foreign-category"
	page, err = service.Catalog(ctx, principal, "store", CatalogFilter{CategoryID: &categoryID}, 1, 1)
	merchandisingNoError(t, err)
	if page.Total != 0 {
		t.Fatalf("foreign category catalog = %+v", page)
	}
}

func TestFranchiseCatalogOrderUsesListingStateThenLatestUpdate(t *testing.T) {
	service, db, principal := fixture(t)
	ctx := context.Background()
	for _, sku := range []gen.ProductSku{
		{ID: "enabled-recent", ProductID: "product", Name: "Zulu Enabled", Enabled: true},
		{ID: "disabled-recent", ProductID: "product", Name: "Alpha Disabled", Enabled: true},
		{ID: "unselected-recent", ProductID: "product", Name: "Beta Unselected", Enabled: true},
	} {
		merchandisingNoError(t, db.Create(&sku).Error)
	}
	for _, selection := range []struct {
		skuID   string
		enabled bool
	}{
		{skuID: "sku", enabled: true},
		{skuID: "enabled-recent", enabled: true},
		{skuID: "disabled-recent", enabled: false},
	} {
		_, err := service.SetListing(ctx, principal, "store", selection.skuID, selection.enabled)
		merchandisingNoError(t, err)
	}
	for skuID, updatedAt := range map[string]int64{"sku": 900, "enabled-recent": 50, "disabled-recent": 1000, "unselected-recent": 500} {
		merchandisingNoError(t, db.Model(&gen.ProductSku{}).Where("id = ?", skuID).UpdateColumn("updated_at", updatedAt).Error)
	}
	for skuID, updatedAt := range map[string]int64{"sku": 100, "enabled-recent": 300, "disabled-recent": 400} {
		merchandisingNoError(t, db.Model(&gen.StoreListing{}).Where("store_id = ? AND sku_id = ?", "store", skuID).UpdateColumn("updated_at", updatedAt).Error)
	}
	page, err := service.Catalog(ctx, principal, "store", CatalogFilter{}, 1, 20)
	merchandisingNoError(t, err)
	ids := make([]string, 0, len(page.Data))
	for _, sku := range page.Data {
		ids = append(ids, sku.ID)
	}
	want := []string{"enabled-recent", "sku", "unselected-recent", "disabled-recent"}
	if len(ids) != len(want) {
		t.Fatalf("catalog order = %v, want %v", ids, want)
	}
	for index := range want {
		if ids[index] != want[index] {
			t.Fatalf("catalog order = %v, want %v", ids, want)
		}
	}
}
