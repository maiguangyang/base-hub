/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package store

import (
	"context"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestStoreLifecycle 验证创建、提交、批准和拒绝的唯一合法状态流。
func TestStoreLifecycle(t *testing.T) {
	service, franchise, hq := newStoreFixture(t)
	first, err := service.Create(context.Background(), franchise, CreateInput{Code: "S1", Name: "First"})
	if err != nil || first.Lifecycle != gen.StoreLifecycleDraft {
		t.Fatalf("create result = %#v, %v", first, err)
	}
	if _, err := service.Submit(context.Background(), franchise, first.ID); err != nil {
		t.Fatal(err)
	}
	approved, err := service.Review(context.Background(), hq, ReviewInput{StoreID: first.ID, Approved: true})
	if err != nil || approved.Lifecycle != gen.StoreLifecycleActive {
		t.Fatalf("approve result = %#v, %v", approved, err)
	}
	if _, err := service.Review(context.Background(), hq, ReviewInput{StoreID: first.ID, Approved: true}); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("duplicate review error = %v", err)
	}
	if _, err := service.Submit(context.Background(), franchise, first.ID); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("invalid resubmit error = %v", err)
	}
	second, _ := service.Create(context.Background(), franchise, CreateInput{Code: "S2", Name: "Second"})
	service.Submit(context.Background(), franchise, second.ID)
	hq.Permissions = map[string]struct{}{"store:reject": {}}
	rejected, err := service.Review(context.Background(), hq, ReviewInput{StoreID: second.ID, Approved: false, RejectionReason: "MISSING_LICENSE"})
	if err != nil || rejected.Lifecycle != gen.StoreLifecycleRejected {
		t.Fatalf("reject result = %#v, %v", rejected, err)
	}
}

// TestHQDirectStoreIsRestrictedToHQOrganization 验证总部直营门店强制 ACTIVE 且不能指定加盟组织。
func TestHQDirectStoreIsRestrictedToHQOrganization(t *testing.T) {
	service, _, hq := newStoreFixture(t)
	hq.Permissions = map[string]struct{}{"hqStore:create": {}}
	created, err := service.Create(context.Background(), hq, CreateInput{Code: "HQ1", Name: "Direct"})
	if err != nil || created.OrganizationID != *hq.OrganizationID || created.Lifecycle != gen.StoreLifecycleActive {
		t.Fatalf("HQ create result = %#v, %v", created, err)
	}
	franchiseID := "franchise-org"
	if _, err := service.Create(context.Background(), hq, CreateInput{Code: "BAD", Name: "Bad", OrganizationID: &franchiseID}); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("cross organization create error = %v", err)
	}
}

// TestStoreReviewValidation 验证退回原因拒绝空白和超过实体上限的值。
func TestStoreReviewValidation(t *testing.T) {
	service, _, hq := newStoreFixture(t)
	hq.Permissions = map[string]struct{}{"store:reject": {}}
	for _, reason := range []string{" ", strings.Repeat("R", 513), strings.Repeat("字", 513)} {
		if _, err := service.Review(context.Background(), hq, ReviewInput{StoreID: "store", RejectionReason: reason}); auth.ErrorCode(err) != auth.CodeValidationFailed {
			t.Fatalf("reason length %d error = %v", len(reason), err)
		}
	}
	validUnicodeReason := strings.Repeat("字", 512)
	if _, err := service.Review(context.Background(), hq, ReviewInput{StoreID: "store", RejectionReason: validUnicodeReason}); auth.ErrorCode(err) == auth.CodeValidationFailed {
		t.Fatalf("valid 512-rune unicode reason should not fail validation: %v", err)
	}
}

func newStoreFixture(t *testing.T) (*Service, *auth.WorkspacePrincipal, *auth.WorkspacePrincipal) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range []any{&gen.Organization{}, &gen.Store{}, &gen.AuditLog{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	franchiseOrg := gen.Organization{ID: "franchise-org", Code: "F", Name: "F", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}
	hqOrg := gen.Organization{ID: "hq-org", Code: "HQ", Name: "HQ", Type: gen.OrganizationTypeHeadquarters, Status: gen.OrganizationStatusActive}
	db.Create(&franchiseOrg)
	db.Create(&hqOrg)
	franchise := &auth.WorkspacePrincipal{AccountID: "owner", SessionID: "fs", WorkspaceType: auth.WorkspaceTypeFranchise, OrganizationID: &franchiseOrg.ID, Permissions: map[string]struct{}{"store:create": {}, "store:submit": {}}, AllStores: true, StoreIDs: map[string]struct{}{}}
	hq := &auth.WorkspacePrincipal{AccountID: "admin", SessionID: "hs", WorkspaceType: auth.WorkspaceTypeHeadquarters, OrganizationID: &hqOrg.ID, Permissions: map[string]struct{}{"store:approve": {}}}
	return NewService(db, audit.NewService()), franchise, hq
}
