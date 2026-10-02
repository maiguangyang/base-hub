/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package integration_test

import (
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
)

func testRoleAndOwnerInvariants(t *testing.T, fixture *securityFixture) {
	assertDirectMembershipWritesDenied(t, fixture)
	assertMembershipGovernance(t, fixture)
	assertLastOwnerGuards(t, fixture)
}

func assertDirectMembershipWritesDenied(t *testing.T, fixture *securityFixture) {
	directMembership := fixture.execute("session-a-1", `mutation {
		createOperatorMembership(input: {
			status: INVITED, storeAccessMode: ALL_STORES,
			accountId: "account-staff", organizationId: "org-a"
		}) { id }
	}`)
	assertCode(t, directMembership, auth.CodePermissionDenied)
	assertCode(t, fixture.execute("session-a-1", `mutation {
		updateOperatorMembership(id: "membership-staff-a", input: {invitationsIds: ["invitation-b"]}) { id }
	}`), auth.CodePermissionDenied)
	fixture.createAll([]gen.Store{{ID: "store-a-draft", Code: "A-DRAFT", Name: "Draft", Lifecycle: gen.StoreLifecycleDraft, OrganizationID: "org-a"}})
	assertCode(t, fixture.execute("session-a-1", `mutation {
		updateStore(id: "store-a-draft", input: {auditLogsIds: ["audit-b"]}) { id }
	}`), auth.CodePermissionDenied)
	assertCode(t, fixture.execute("session-a-1", `mutation {
		createStore(input: {code: "A-TAMPER", name: "Tamper", lifecycle: DRAFT, organizationId: "org-a", reviewedByAccount: {id: "account-shared"}}) { id }
	}`), auth.CodePermissionDenied)
	assertCode(t, fixture.execute("session-a-1", `mutation { updateStore(id: "store-a", input: {name: "Bypass"}) { id } }`), auth.CodeConflict)
	directInvitation := fixture.execute("session-a-1", `mutation {
		createMembershipInvitation(input: {
			membershipId: "membership-staff-a", invitedByAccountId: "account-shared",
			expiresAt: "2026-09-24T00:00:00Z"
		}) { id }
	}`)
	assertCode(t, directInvitation, auth.CodePermissionDenied)
}

func assertMembershipGovernance(t *testing.T, fixture *securityFixture) {
	fixture.createAll([]gen.MembershipInvitation{{ID: "invitation-staff-a", MembershipID: "membership-staff-a", InvitedByAccountID: "account-shared", ExpiresAt: time.Now().Add(time.Hour)}})
	assertNoErrors(t, fixture.execute("session-a-1", `mutation { deleteMembershipInvitations(id: ["invitation-staff-a"]) }`))
	assertAuditCount(t, fixture, "membershipInvitation:delete", 1)
	bypass := fixture.execute("session-a-1", `mutation { updateOperatorMembership(id: "membership-owner-a", input: {isDelete: 2}) { id isDelete } }`)
	assertNoErrors(t, bypass)
	if !strings.Contains(bypass.Body, `"isDelete":1`) {
		t.Fatalf("update permission changed soft-delete control: %s", bypass.Body)
	}
	assertCode(t, fixture.execute("session-a-1", `mutation {
		updateOperatorMembership(id: "membership-owner-a", input: {storesIds: ["store-a"]}) { id }
	}`), auth.CodeValidationFailed)
	assertCode(t, fixture.execute("session-a-1", `mutation {
		updateOperatorMembership(id: "membership-staff-a", input: {storeAccessMode: SELECTED_STORES, storesIds: []}) { id }
	}`), auth.CodeValidationFailed)
	staffAccess := fixture.execute("session-a-1", `mutation {
		updateOperatorMembership(id: "membership-staff-a", input: {
			storeAccessMode: SELECTED_STORES, storesIds: ["store-a"], rolesIds: ["role-custom-a"]
		}) { id storeAccessMode storesIds rolesIds }
	}`)
	assertNoErrors(t, staffAccess)
	allStores := fixture.execute("session-a-1", `mutation {
		updateOperatorMembership(id: "membership-staff-a", input: {storeAccessMode: ALL_STORES}) { id storeAccessMode storesIds }
	}`)
	assertNoErrors(t, allStores)
	if !strings.Contains(allStores.Body, `"storeAccessMode":"ALL_STORES","storesIds":[]`) {
		t.Fatalf("ALL_STORES retained selected store associations: %s", allStores.Body)
	}
	roleUpdate := fixture.execute("session-a-1", `mutation { updateOperatorRole(id: "role-custom-a", input: {permissionsIds: ["permission-store-read"]}) { id } }`)
	assertNoErrors(t, roleUpdate)
	assertAuditCount(t, fixture, "operatorRole:update", 1)
	assign := fixture.execute("session-a-1", `mutation { updateOperatorRole(id: "role-custom-a", input: {permissionsIds: ["permission-organization-suspend"]}) { id } }`)
	assertCode(t, assign, auth.CodePermissionDenied)
}

func assertLastOwnerGuards(t *testing.T, fixture *securityFixture) {
	removeRole := fixture.execute("session-a-1", `mutation { updateOperatorMembership(id: "membership-owner-a", input: {rolesIds: ["role-custom-a"]}) { id } }`)
	assertCode(t, removeRole, auth.CodeLastOwnerRequired)
	lastOwner := fixture.execute("session-a-1", `mutation { changeMembershipStatus(input: {membershipId: "membership-owner-a", status: SUSPENDED}) { id } }`)
	assertCode(t, lastOwner, auth.CodeLastOwnerRequired)
	deleteOwner := fixture.execute("session-a-1", `mutation { deleteOperatorMemberships(id: ["membership-owner-a"]) }`)
	assertCode(t, deleteOwner, auth.CodeLastOwnerRequired)
	assertNoErrors(t, fixture.execute("session-a-1", `mutation { deleteOperatorMemberships(id: ["membership-staff-a"]) }`))
	updateDeleted := fixture.execute("session-a-1", `mutation { updateOperatorMembership(id: "membership-staff-a", input: {weight: 2}) { id } }`)
	assertCode(t, updateDeleted, auth.CodePermissionDenied)
	recoverActive := fixture.execute("session-a-1", `mutation { recoveryOperatorMemberships(id: ["membership-owner-a"]) }`)
	assertCode(t, recoverActive, auth.CodePermissionDenied)
	assertDeletedMembershipConcealed(t, fixture)
	assertCode(t, fixture.execute("session-a-2", `query { viewer { account { id } } }`), auth.CodeMembershipInactive)
	assertAuditCount(t, fixture, "operatorMembership:delete", 1)
}

func assertDeletedMembershipConcealed(t *testing.T, fixture *securityFixture) {
	t.Helper()
	relations := fixture.execute("session-a-1", `query {
		operatorRole(id: "role-custom-a") { members { id } }
		store(id: "store-a") { members { id } }
	}`)
	assertNoErrors(t, relations)
	if strings.Contains(relations.Body, "membership-staff-a") {
		t.Fatalf("soft-deleted membership leaked through relation: %s", relations.Body)
	}
}
