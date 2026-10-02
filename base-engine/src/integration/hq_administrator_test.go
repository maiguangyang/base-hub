/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-22
 */
package integration_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
)

func TestHQMembershipQueryIsRestrictedToHeadquarters(t *testing.T) {
	fixture := newSecurityFixture(t)
	response := fixture.execute("session-hq", `query {
		operatorMemberships { total data { id organizationId } }
	}`)
	assertNoErrors(t, response)
	assertResultTotal(t, response, "operatorMemberships", 1)
	if response.Body != "" && containsAny(response.Body, "membership-owner-a", "membership-staff-a", "membership-owner-b") {
		t.Fatalf("HQ membership query leaked franchise memberships: %s", response.Body)
	}
}

func TestHQMembershipQuerySearchesAccountIdentity(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedHQCustomAdministrator(fixture, "search-target", []string{"hqMembership:read"})
	if err := fixture.db.Model(&gen.Account{}).Where("id = ?", "account-hq-search-target").Updates(map[string]interface{}{
		"display_name": "Search Needle",
		"phone":        "13988887777",
	}).Error; err != nil {
		t.Fatal(err)
	}
	fixture.createAll(
		[]gen.Account{{ID: "account-search-foreign", Phone: "13988886666", DisplayName: "Search Needle", Status: gen.AccountStatusActive, CredentialVersion: 1}},
		[]gen.OperatorMembership{{
			ID: "membership-search-foreign", AccountID: "account-search-foreign", OrganizationID: "org-a",
			Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores,
		}},
	)

	byName := fixture.execute("session-hq", `query {
		operatorMemberships(q: "Search Needle") { total data { id account { displayName phone } } }
	}`)
	assertNoErrors(t, byName)
	assertResultTotal(t, byName, "operatorMemberships", 1)
	if !strings.Contains(byName.Body, "membership-hq-search-target") || strings.Contains(byName.Body, "membership-search-foreign") {
		t.Fatalf("unexpected HQ membership name search result: %s", byName.Body)
	}

	byPhone := fixture.execute("session-hq", `query {
		operatorMemberships(q: "13988887777") { total data { id account { displayName phone } } }
	}`)
	assertNoErrors(t, byPhone)
	assertResultTotal(t, byPhone, "operatorMemberships", 1)
	if !strings.Contains(byPhone.Body, "membership-hq-search-target") || strings.Contains(byPhone.Body, "membership-search-foreign") {
		t.Fatalf("unexpected HQ membership phone search result: %s", byPhone.Body)
	}
}

func TestHQMembershipQueryFiltersByRole(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedHQCustomAdministrator(fixture, "filter-other", []string{"hqMembership:read"})
	fixture.createAll([]membershipRole{{"membership-hq", "role-hq"}})

	response := fixture.execute("session-hq", `query {
		operatorMemberships(filter: {roles: {id: "role-hq"}}) { total data { id } }
	}`)
	assertNoErrors(t, response)
	assertResultTotal(t, response, "operatorMemberships", 1)
	var payload struct {
		OperatorMemberships struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		} `json:"operatorMemberships"`
	}
	decodeData(t, response, &payload)
	if len(payload.OperatorMemberships.Data) != 1 {
		t.Fatalf("HQ role filter returned duplicate memberships: %s", response.Body)
	}
	if !strings.Contains(response.Body, "membership-hq") || strings.Contains(response.Body, "membership-hq-filter-other") {
		t.Fatalf("unexpected HQ role filter result: %s", response.Body)
	}
}

func TestHQMembershipQueryCombinesSearchAndRoleFilter(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedHQCustomAdministrator(fixture, "filter-other", []string{"hqMembership:read"})

	response := fixture.execute("session-hq", `query {
		operatorMemberships(q: "filter-other", filter: {roles: {id: "role-hq"}}) { total data { id } }
	}`)
	assertNoErrors(t, response)
	assertResultTotal(t, response, "operatorMemberships", 0)
}

func TestHQRoleReaderCanLoadSystemPermissionCatalog(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedHQCustomAdministrator(fixture, "catalog-reader", []string{"hqRole:read"})
	response := fixture.execute("session-hq-catalog-reader", `query {
		permissions(current_page: 1, per_page: 200, filter: {scope: SYSTEM}) { total data { id scope } }
	}`)
	assertNoErrors(t, response)
	if response.Body != "" && strings.Contains(response.Body, `"scope":"TENANT"`) {
		t.Fatalf("permission catalog leaked tenant permissions: %s", response.Body)
	}
}

func TestHQAdministratorCannotManageSelfOrSuper(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedHQCustomAdministrator(fixture, "limited", []string{"hqMembership:read", "hqMembership:update"})
	self := fixture.execute("session-hq-limited", `mutation {
		updateOperatorMembership(id: "membership-hq-limited", input: {rolesIds: ["role-hq-limited"]}) { id }
	}`)
	assertCode(t, self, auth.Code("SELF_MEMBERSHIP_CHANGE_DENIED"))

	superTarget := fixture.execute("session-hq-limited", `mutation {
		updateOperatorMembership(id: "membership-hq", input: {rolesIds: ["role-hq-limited"]}) { id }
	}`)
	assertCode(t, superTarget, auth.CodePermissionDelegationDenied)

	suspendSuper := fixture.execute("session-hq-limited", `mutation {
		changeMembershipStatus(input: {membershipId: "membership-hq", status: SUSPENDED}) { id }
	}`)
	assertCode(t, suspendSuper, auth.CodePermissionDelegationDenied)
}

func TestHQRoleCannotEscalateOrDeleteInUse(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedHQCustomAdministrator(fixture, "role-manager", []string{"hqRole:read", "hqRole:update", "hqRole:delete"})
	fixture.createAll(
		[]gen.OperatorRole{{ID: "role-hq-stronger", Name: "Stronger", Kind: gen.RoleKindCustom, OrganizationID: "org-hq"}},
		[]permissionRole{{"permission-account-update", "role-hq-stronger"}},
	)
	renameStronger := fixture.execute("session-hq-role-manager", `mutation {
		updateOperatorRole(id: "role-hq-stronger", input: {name: "Renamed"}) { id }
	}`)
	assertCode(t, renameStronger, auth.CodePermissionDelegationDenied)

	escalate := fixture.execute("session-hq-role-manager", `mutation {
		updateOperatorRole(id: "role-hq-role-manager", input: {permissionsIds: ["permission-account-update"]}) { id }
	}`)
	assertCode(t, escalate, auth.CodePermissionDelegationDenied)

	deleteInUse := fixture.execute("session-hq-role-manager", `mutation {
		deleteOperatorRoles(id: ["role-hq-role-manager"])
	}`)
	assertCode(t, deleteInUse, auth.Code("ROLE_IN_USE"))
}

func TestHQRecoveryMutationsAreClosed(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedHQCustomAdministrator(fixture, "recovery-manager", []string{"hqMembership:update", "hqRole:update"})
	deleted := int64(2)
	fixture.createAll(
		[]gen.Account{{ID: "account-hq-deleted-super", Phone: "13900000009", DisplayName: "Deleted Super", Status: gen.AccountStatusActive, CredentialVersion: 1}},
		[]gen.OperatorMembership{{
			ID: "membership-hq-deleted-super", AccountID: "account-hq-deleted-super", OrganizationID: "org-hq",
			Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores, IsDelete: &deleted,
		}},
		[]membershipRole{{"membership-hq-deleted-super", "role-hq"}},
		[]gen.OperatorRole{{
			ID: "role-hq-deleted-strong", Name: "Deleted Strong", Kind: gen.RoleKindCustom,
			OrganizationID: "org-hq", IsDelete: &deleted,
		}},
		[]permissionRole{{"permission-account-update", "role-hq-deleted-strong"}},
	)

	recoverMembership := fixture.execute("session-hq-recovery-manager", `mutation {
		recoveryOperatorMemberships(id: ["membership-hq-deleted-super"])
	}`)
	assertCode(t, recoverMembership, auth.CodePermissionDenied)

	recoverRole := fixture.execute("session-hq-recovery-manager", `mutation {
		recoveryOperatorRoles(id: ["role-hq-deleted-strong"])
	}`)
	assertCode(t, recoverRole, auth.CodePermissionDenied)

	var membership gen.OperatorMembership
	if err := fixture.db.Unscoped().Where("id = ?", "membership-hq-deleted-super").First(&membership).Error; err != nil {
		t.Fatal(err)
	}
	if membership.IsDelete == nil || *membership.IsDelete != deleted {
		t.Fatalf("deleted membership was recovered: %#v", membership.IsDelete)
	}
	var role gen.OperatorRole
	if err := fixture.db.Unscoped().Where("id = ?", "role-hq-deleted-strong").First(&role).Error; err != nil {
		t.Fatal(err)
	}
	if role.IsDelete == nil || *role.IsDelete != deleted {
		t.Fatalf("deleted role was recovered: %#v", role.IsDelete)
	}
}

func seedHQCustomAdministrator(fixture *securityFixture, suffix string, actions []string) {
	fixture.t.Helper()
	accountID := "account-hq-" + suffix
	membershipID := "membership-hq-" + suffix
	roleID := "role-hq-" + suffix
	fixture.createAll(
		[]gen.Account{{ID: accountID, Phone: "139" + suffix, DisplayName: suffix, Status: gen.AccountStatusActive, CredentialVersion: 1}},
		[]gen.OperatorMembership{{ID: membershipID, AccountID: accountID, OrganizationID: "org-hq", Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores}},
		[]gen.OperatorRole{{ID: roleID, Name: suffix, Kind: gen.RoleKindCustom, OrganizationID: "org-hq"}},
		[]membershipRole{{membershipID, roleID}},
	)
	for _, action := range actions {
		fixture.createAll([]permissionRole{{"permission-" + strings.ReplaceAll(action, ":", "-"), roleID}})
	}
	seedHQSession(fixture, suffix, accountID)
}

func seedHQSession(fixture *securityFixture, suffix, accountID string) {
	fixture.t.Helper()
	sessionID := "session-hq-" + suffix
	organizationID := "org-hq"
	fixture.createAll([]gen.Session{{
		ID: sessionID, AccountID: accountID, OrganizationID: &organizationID,
		WorkspaceType: gen.WorkspaceTypeHeadquarters, CredentialVersion: 1,
		ExpiresAt: time.Now().Add(time.Hour), LastSeenAt: time.Now(),
	}})
	claims := auth.SessionClaims{
		SessionID: sessionID, AccountID: accountID, OrganizationID: &organizationID,
		WorkspaceType: auth.WorkspaceTypeHeadquarters, CredentialVersion: 1,
	}
	token, err := auth.SignSessionClaims(fixture.cfg, claims)
	if err != nil {
		fixture.t.Fatal(err)
	}
	fixture.cookies[sessionID] = &http.Cookie{Name: auth.SessionCookieName, Value: token}
}

func containsAny(value string, candidates ...string) bool {
	for _, candidate := range candidates {
		if strings.Contains(value, candidate) {
			return true
		}
	}
	return false
}
