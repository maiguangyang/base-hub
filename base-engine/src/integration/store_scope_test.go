/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package integration_test

import (
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func testSelectedStoreScope(t *testing.T, fixture *securityFixture) {
	hidden := gen.Store{ID: "store-a-hidden", Code: "A2", Name: "Hidden", Lifecycle: gen.StoreLifecycleActive, OrganizationID: "org-a"}
	candidate := gen.Account{ID: "account-candidate", Phone: "13800000003", DisplayName: "Candidate", Status: gen.AccountStatusActive, CredentialVersion: 1}
	fixture.createAll([]gen.Store{hidden}, []gen.Account{candidate}, []membershipStore{
		{"membership-staff-a", "store-a"}, {"membership-owner-a", "store-a-hidden"},
	})
	if err := fixture.db.Model(&gen.OperatorMembership{}).
		Where("id IN ?", []string{"membership-staff-a", "membership-owner-a"}).
		Update("store_access_mode", gen.StoreAccessModeSelectedStores).Error; err != nil {
		t.Fatal(err)
	}

	assertSelectedStoreRelations(t, fixture, hidden.ID)
	assertSelectedStoreWrites(t, fixture)
}

func assertSelectedStoreRelations(t *testing.T, fixture *securityFixture, hiddenID string) {
	response := fixture.execute("session-a-2", `query {
		organizations { data { stores { id } storesIds } }
		stores { data { id } }
		operatorMembership(id: "membership-owner-a") { stores { id } storesIds }
	}`)
	assertNoErrors(t, response)
	if strings.Contains(response.Body, hiddenID) {
		t.Fatalf("selected-store relation leaked hidden store: %s", response.Body)
	}
}

func assertSelectedStoreWrites(t *testing.T, fixture *securityFixture) {
	update := fixture.execute("session-a-2", `mutation {
		updateOperatorMembership(id: "membership-owner-a", input: {
			storeAccessMode: SELECTED_STORES, storesIds: ["store-a-hidden"]
		}) { id }
	}`)
	assertCode(t, update, auth.CodeStoreScopeDenied)
	allStores := fixture.execute("session-a-2", `mutation {
		updateOperatorMembership(id: "membership-owner-a", input: {storeAccessMode: ALL_STORES}) { id }
	}`)
	assertCode(t, allStores, auth.CodeStoreScopeDenied)

	invite := fixture.execute("session-a-2", `mutation {
		inviteOperator(input: {
			phone: "13800000003", displayName: "Candidate", roleIds: ["role-custom-a"],
			storeAccessMode: SELECTED_STORES, storeIds: ["store-a-hidden"]
		}) { membership { id } }
	}`)
	assertCode(t, invite, auth.CodeStoreScopeDenied)
}
