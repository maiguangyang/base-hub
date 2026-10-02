package authorization

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateStoreWithGraphQLBusinessStatus(t *testing.T) {
	ctx, resolver, organizationID := storeBusinessStatusContext(t)
	status := gen.StoreBusinessStatusOpen
	// gqlgen unmarshals the optional enum into a pointer before calling the handler.
	created, err := createStore(ctx, resolver, map[string]interface{}{
		"code": "STR001", "name": "新门店", "organizationId": organizationID, "businessStatus": &status,
		"businessHours": "09:00 - 22:00", "supportDineIn": true, "storeArea": nil, "tableCount": nil,
	})
	if err != nil {
		t.Fatalf("createStore failed: %v", err)
	}
	assertCreatedStoreBusinessStatus(t, created, status)
	var count int64
	if err := gen.GetTransaction(ctx).Model(&gen.AuditLog{}).Count(&count).Error; err != nil || count != 1 {
		t.Fatalf("audit count=%d error=%v", count, err)
	}
	closed := gen.StoreBusinessStatusClosed
	updated, err := updateStore(ctx, resolver, created.ID, map[string]interface{}{"businessStatus": &closed})
	if err != nil || updated.BusinessStatus == nil || *updated.BusinessStatus != closed {
		t.Fatalf("updateStore failed: store=%#v error=%v", updated, err)
	}
	var stored gen.Store
	if err := gen.GetTransaction(ctx).First(&stored, "id = ?", created.ID).Error; err != nil || stored.BusinessStatus == nil || *stored.BusinessStatus != closed {
		t.Fatalf("business status was not persisted: store=%#v error=%v", stored, err)
	}
}

func storeBusinessStatusContext(t *testing.T) (context.Context, *gen.GeneratedResolver, string) {
	t.Helper()
	t.Setenv("TABLE_NAME_PREFIX", "")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range []any{&gen.Organization{}, &gen.Store{}, &gen.AuditLog{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			if err := db.Exec("DROP INDEX IF EXISTS `" + index + "`").Error; err != nil {
				t.Fatal(err)
			}
		}
	}
	org := gen.Organization{ID: "org", Code: "F001", Name: "Franchise", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}
	if err := db.Create(&org).Error; err != nil {
		t.Fatal(err)
	}
	principal := &auth.WorkspacePrincipal{AccountID: "owner", SessionID: "session", WorkspaceType: auth.WorkspaceTypeFranchise,
		OrganizationID: &org.ID, AllStores: true, Permissions: map[string]struct{}{"store:create": {}, "store:update": {}}}
	resolver := &gen.GeneratedResolver{DB: gen.NewDB(db)}
	ctx := gen.EnrichContextWithMutations(auth.WithPrincipal(context.Background(), principal), resolver)
	t.Cleanup(func() { gen.RollbackMutationContext(ctx, resolver) })
	return ctx, resolver, org.ID
}

func assertCreatedStoreBusinessStatus(t *testing.T, created *gen.Store, status gen.StoreBusinessStatus) {
	t.Helper()
	if created.Lifecycle != gen.StoreLifecycleDraft || created.BusinessStatus == nil || *created.BusinessStatus != status {
		t.Fatalf("unexpected created store: %#v", created)
	}
}
