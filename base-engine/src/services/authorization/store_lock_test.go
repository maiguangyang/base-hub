/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestLoadAuthorizedStoreLocksWrites 验证生命周期检查与更新共享同一行锁。
func TestLoadAuthorizedStoreLocksWrites(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(&gen.Store{}); err != nil {
		t.Fatal(err)
	}
	store := gen.Store{ID: "store-1", Code: "S1", Name: "Store", Lifecycle: gen.StoreLifecycleDraft, OrganizationID: "org-1"}
	if err := database.Create(&store).Error; err != nil {
		t.Fatal(err)
	}
	locked := false
	if err := database.Callback().Query().Before("gorm:query").Register("test:store_lock", func(tx *gorm.DB) {
		_, locked = tx.Statement.Clauses["FOR"]
	}); err != nil {
		t.Fatal(err)
	}
	organizationID := "org-1"
	principal := &auth.WorkspacePrincipal{
		AccountID: "account-1", WorkspaceType: auth.WorkspaceTypeFranchise,
		OrganizationID: &organizationID, AllStores: true,
		Permissions: map[string]struct{}{"store:update": {}},
	}
	ctx := auth.WithPrincipal(context.Background(), principal)
	resolver := &gen.GeneratedResolver{DB: gen.NewDB(database)}
	if _, _, err := loadAuthorizedStore(ctx, resolver, store.ID, "update", AccessUpdate); err != nil {
		t.Fatal(err)
	}
	if !locked {
		t.Fatal("store write authorization did not request a row lock")
	}
}
