/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package integration_test

import (
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

// TestStorePermissionContract 验证总部有限读取、跨店读取与门店软删除契约。
func TestStorePermissionContract(t *testing.T) {
	fixture := newSecurityFixture(t)
	allStores := fixture.execute("session-hq", `query { stores { total data { id organization { id } } } }`)
	assertNoErrors(t, allStores)
	assertResultTotal(t, allStores, "stores", 3)
	revokeHQCrossStoreRead(t, fixture)

	directOnly := fixture.execute("session-hq", `query { stores { total data { id } } }`)
	assertNoErrors(t, directOnly)
	assertResultTotal(t, directOnly, "stores", 1)
	if response := fixture.execute("session-hq", `query { organization(id: "org-hq") { stores { id } } }`); len(response.Errors) != 0 {
		t.Fatalf("HQ relationship read failed: %s", response.Body)
	}
	assertCode(t, fixture.execute("session-hq", `query { store(id: "store-a") { id } }`), auth.CodePermissionDenied)

	assertNoErrors(t, fixture.execute("session-hq", `mutation { deleteStores(id: ["store-hq"]) }`))
	assertCode(t, fixture.execute("session-hq", `mutation { deleteStores(id: ["store-a"]) }`), auth.CodePermissionDenied)
	assertAuditCount(t, fixture, "hqStore:delete", 1)

	testFranchiseStoreDeletion(t, fixture)
}

func revokeHQCrossStoreRead(t *testing.T, fixture *securityFixture) {
	t.Helper()
	result := fixture.db.Model(&gen.Permission{}).Where("action = ?", "store:read_all").Update("scope", gen.PermissionScopeTenant)
	if result.Error != nil || result.RowsAffected != 1 {
		t.Fatalf("remove cross-store permission rows=%d err=%v", result.RowsAffected, result.Error)
	}
}

func testFranchiseStoreDeletion(t *testing.T, fixture *securityFixture) {
	t.Helper()
	stores := []gen.Store{
		{ID: "store-a-draft-delete", Code: "AD", Name: "Draft", OrganizationID: "org-a", Lifecycle: gen.StoreLifecycleDraft},
		{ID: "store-a-rejected-delete", Code: "AR", Name: "Rejected", OrganizationID: "org-a", Lifecycle: gen.StoreLifecycleRejected},
		{ID: "store-a-pending-delete", Code: "AP", Name: "Pending", OrganizationID: "org-a", Lifecycle: gen.StoreLifecyclePendingApproval},
	}
	fixture.createAll(stores)
	assertNoErrors(t, fixture.execute("session-a-1", `mutation { deleteStores(id: ["store-a-draft-delete", "store-a-rejected-delete"]) }`))
	assertCode(t, fixture.execute("session-a-1", `mutation { deleteStores(id: ["store-a"]) }`), auth.CodeConflict)
	assertCode(t, fixture.execute("session-a-1", `mutation { deleteStores(id: ["store-a-pending-delete"]) }`), auth.CodeConflict)
	assertAuditCount(t, fixture, "store:delete", 2)
}
