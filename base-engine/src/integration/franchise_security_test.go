/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package integration_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
)

// TestFranchiseSecurityMatrix 通过真实 gqlgen HTTP 边界验证租户、总部和强退矩阵。
func TestFranchiseSecurityMatrix(t *testing.T) {
	t.Run("tenant data isolation", func(t *testing.T) { testTenantIsolation(t, newSecurityFixture(t)) })
	t.Run("membership role filter", func(t *testing.T) { testMembershipRoleFilter(t, newSecurityFixture(t)) })
	t.Run("membership store filter", func(t *testing.T) { testMembershipStoreFilter(t, newSecurityFixture(t)) })
	t.Run("generated validation error code", func(t *testing.T) { testGeneratedValidationErrorCode(t, newSecurityFixture(t)) })
	t.Run("selected store scope", func(t *testing.T) { testSelectedStoreScope(t, newSecurityFixture(t)) })
	t.Run("headquarters governance", func(t *testing.T) { testHeadquartersGovernance(t, newSecurityFixture(t)) })
	t.Run("role and owner invariants", func(t *testing.T) { testRoleAndOwnerInvariants(t, newSecurityFixture(t)) })
	t.Run("suspension events and reconnect", func(t *testing.T) { testSuspensionAndReconnect(t, newSecurityFixture(t)) })
	t.Run("handler coverage probe", func(t *testing.T) { assertHandlerCoverage(t) })
}

func testMembershipRoleFilter(t *testing.T, fixture *securityFixture) {
	response := fixture.execute("session-a-1", `query {
		operatorMemberships(filter: {roles: {id: "role-custom-a"}}) { total data { id } }
	}`)
	assertNoErrors(t, response)
	assertResultTotal(t, response, "operatorMemberships", 1)
	if !strings.Contains(response.Body, "membership-staff-a") || strings.Contains(response.Body, "membership-owner-a") {
		t.Fatalf("unexpected tenant role filter result: %s", response.Body)
	}
}

func testMembershipStoreFilter(t *testing.T, fixture *securityFixture) {
	fixture.createAll(
		[]gen.Store{{ID: "store-a-second", Code: "A2", Name: "A Second", Lifecycle: gen.StoreLifecycleActive, OrganizationID: "org-a"}},
		[]membershipStore{{"membership-staff-a", "store-a-second"}, {"membership-owner-a", "store-a"}},
	)
	response := fixture.execute("session-a-1", `query {
		operatorMemberships(filter: {stores: {id: "store-a-second"}}) { total data { id } }
	}`)
	assertNoErrors(t, response)
	assertResultTotal(t, response, "operatorMemberships", 1)
	if !strings.Contains(response.Body, "membership-staff-a") || strings.Contains(response.Body, "membership-owner-a") {
		t.Fatalf("unexpected tenant store filter result: %s", response.Body)
	}
}

func testTenantIsolation(t *testing.T, fixture *securityFixture) {
	query := `query {
		accounts { total data { id membershipsIds sessionsIds reviewedStoresIds sentMembershipInvitationsIds auditLogsIds sessions { id organizationId } reviewedStores { id } sentMembershipInvitations { id } auditLogs { id organizationId } } }
		organizations { total data { id stores { id } roles { id } memberships { id } sessionsIds auditLogsIds sessions { id } auditLogs { id } } }
		stores { total data { id auditLogsIds organization { id } members { id organization { id } } reviewedByAccount { id } auditLogs { id } } }
		operatorMemberships { total data { id invitationsIds account { id } organization { id } roles { id organization { id } } stores { id organization { id } } invitations { id membership { id } invitedByAccount { id } } } }
		operatorRoles { total data { id organization { id } members { id } membersIds permissions { id scope } permissionsIds } }
		permissions { total data { id rolesIds } }
		sessions { total data { id auditLogsIds account { id } organization { id } auditLogs { id } } }
		auditLogs { total data { id organizationId actorAccount { id } session { id } organization { id } store { id } } }
	}`
	response := fixture.execute("session-a-1", query)
	assertNoErrors(t, response)
	for _, forbidden := range []string{"org-b", "store-b", "membership-owner-b", "role-owner-b", "session-b", "invitation-b", "audit-b"} {
		if strings.Contains(response.Body, forbidden) {
			t.Fatalf("tenant B identifier leaked: %s body=%s", forbidden, response.Body)
		}
	}
	assertResultTotal(t, response, "accounts", 2)
	assertResultTotal(t, response, "organizations", 1)
	assertResultTotal(t, response, "stores", 1)
	assertResultTotal(t, response, "operatorMemberships", 2)
	assertResultTotal(t, response, "operatorRoles", 2)
	assertResultTotal(t, response, "sessions", 1)
	assertResultTotal(t, response, "auditLogs", 1)
	viewer := fixture.execute("session-a-1", `query { viewer { account { membershipsIds } } }`)
	assertNoErrors(t, viewer)
	if strings.Contains(viewer.Body, "membership-owner-b") || strings.Contains(viewer.Body, "membership-hq") {
		t.Fatalf("viewer relationship IDs escaped tenant scope: %s", viewer.Body)
	}

	assertCode(t, fixture.execute("session-a-1", `query { store(id: "store-b") { id } }`), auth.CodePermissionDenied)
	mutations := []string{
		`mutation { createStore(input: {code: "X", name: "X", lifecycle: DRAFT, organizationId: "org-b"}) { id } }`,
		`mutation { updateStore(id: "store-b", input: {name: "Changed"}) { id } }`,
		`mutation { deleteStores(id: ["store-b"]) }`,
		`mutation { recoveryStores(id: ["store-b"]) }`,
	}
	for _, mutation := range mutations {
		assertCode(t, fixture.execute("session-a-1", mutation), auth.CodePermissionDenied)
	}
	filtered := fixture.execute("session-a-1", `query { stores(filter: {organizationId: "org-b"}) { total data { id } } }`)
	assertNoErrors(t, filtered)
	assertResultTotal(t, filtered, "stores", 0)
}

func testGeneratedValidationErrorCode(t *testing.T, fixture *securityFixture) {
	response := fixture.execute("session-a-1", `mutation {
		createStore(input: {code: "A1", name: "Duplicate", lifecycle: DRAFT, organizationId: "org-a"}) { id }
	}`)
	assertCode(t, response, auth.CodeValidationFailed)
}

func testHeadquartersGovernance(t *testing.T, fixture *securityFixture) {
	created := fixture.execute("session-hq", `mutation { createStore(input: {code: "HQ2", name: "Second", lifecycle: DRAFT, organizationId: "org-hq"}) { id lifecycle organization { id } } }`)
	assertNoErrors(t, created)
	if !strings.Contains(created.Body, `"lifecycle":"ACTIVE"`) || !strings.Contains(created.Body, `"id":"org-hq"`) {
		t.Fatalf("direct store was not forced ACTIVE in HQ: %s", created.Body)
	}
	var payload struct {
		CreateStore struct {
			ID string `json:"id"`
		} `json:"createStore"`
	}
	decodeData(t, created, &payload)
	assertAuditCount(t, fixture, "hqStore:create", 1)
	assertNoErrors(t, fixture.execute("session-hq", `mutation { updateStore(id: "`+payload.CreateStore.ID+`", input: {name: "Renamed"}) { id } }`))
	assertAuditCount(t, fixture, "hqStore:update", 1)
	assertNoErrors(t, fixture.execute("session-hq", `mutation { deleteStores(id: ["`+payload.CreateStore.ID+`"]) }`))
	assertAuditCount(t, fixture, "hqStore:delete", 1)
	var deletedStore gen.Store
	if err := fixture.db.Unscoped().First(&deletedStore, "id = ?", payload.CreateStore.ID).Error; err != nil || deletedStore.IsDelete == nil || *deletedStore.IsDelete != 2 {
		t.Fatalf("store was not soft deleted: %#v err=%v", deletedStore, err)
	}
	assertCode(t, fixture.execute("session-hq", `mutation { recoveryStores(id: ["`+payload.CreateStore.ID+`"]) }`), auth.CodePermissionDenied)
	for _, storeID := range []string{"store-hq", "store-a", "store-b"} {
		assertNoErrors(t, fixture.execute("session-hq", `query { store(id: "`+storeID+`") { id organization { id } } }`))
	}
	remainingStores := fixture.execute("session-hq", `query { stores { total data { id } } }`)
	assertNoErrors(t, remainingStores)
	if strings.Contains(remainingStores.Body, payload.CreateStore.ID) {
		t.Fatalf("deleted store remained visible: %s", remainingStores.Body)
	}
	response := fixture.execute("session-hq", `query { organizations { total data { id stores { id } } } stores { total data { id organization { id } } } auditLogs { total data { id organizationId } } }`)
	assertNoErrors(t, response)
	assertResultTotal(t, response, "organizations", 3)
	assertResultTotal(t, response, "stores", 3)
	assertResultTotal(t, response, "auditLogs", 5)
	assertCode(t, fixture.execute("session-hq", `mutation { updateStore(id: "store-a", input: {name: "Denied"}) { id } }`), auth.CodePermissionDenied)
	assertCode(t, fixture.execute("session-hq", `mutation { updateOperatorMembership(id: "membership-owner-a", input: {weight: 2}) { id } }`), auth.CodePermissionDenied)
	assertCode(t, fixture.execute("session-hq", `mutation { updateOperatorRole(id: "role-custom-a", input: {weight: 2}) { id } }`), auth.CodePermissionDenied)
	assertCode(t, fixture.execute("session-hq", `mutation { createOperatorRole(input: {name: "Bypass", kind: CUSTOM, organizationId: "org-a"}) { id } }`), auth.CodePermissionDenied)
	assertCode(t, fixture.execute("session-a-1", `mutation { reviewStore(input: {storeId: "store-a", approved: true}) { id } }`), auth.CodeWorkspaceForbidden)
}

func testSuspensionAndReconnect(t *testing.T, fixture *securityFixture) {
	ctxA1, cancelA1, channelA1 := subscribe(t, fixture, "session-a-1")
	defer cancelA1()
	ctxA2, cancelA2, channelA2 := subscribe(t, fixture, "session-a-2")
	defer cancelA2()
	ctxB, cancelB, channelB := subscribe(t, fixture, "session-b")
	defer cancelB()
	_ = ctxA1
	_ = ctxA2
	_ = ctxB

	response := fixture.execute("session-hq", `mutation { suspendOrganization(input: {organizationId: "org-a", reasonCode: "CONTRACT_ENDED"}) { id status } }`)
	assertNoErrors(t, response)
	assertSessionEvent(t, channelA1, "session-a-1")
	assertSessionEvent(t, channelA2, "session-a-2")
	assertNoSessionEvent(t, channelB)

	for _, sessionID := range []string{"session-a-1", "session-a-2"} {
		var session gen.Session
		if err := fixture.db.First(&session, "id = ?", sessionID).Error; err != nil || session.RevokedAt == nil {
			t.Fatalf("session %s was not revoked: %#v err=%v", sessionID, session, err)
		}
	}
	for _, sessionID := range []string{"session-b", "session-hq"} {
		var session gen.Session
		if err := fixture.db.First(&session, "id = ?", sessionID).Error; err != nil || session.RevokedAt != nil {
			t.Fatalf("unrelated session %s changed: %#v err=%v", sessionID, session, err)
		}
		assertNoErrors(t, fixture.execute(sessionID, `query { viewer { account { id } } }`))
	}
	reconnect := fixture.execute("session-a-1", `query { viewer { account { id } } }`)
	assertCode(t, reconnect, auth.CodeSessionRevoked)
}

func subscribe(t *testing.T, fixture *securityFixture, sessionID string) (context.Context, context.CancelFunc, <-chan *gen.SessionEvent) {
	t.Helper()
	ctx, cancel := context.WithCancel(auth.WithPrincipal(context.Background(), fixture.principals[sessionID]))
	channel, err := fixture.resolver.Subscription().SessionEvents(ctx)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	return ctx, cancel, channel
}

func assertSessionEvent(t *testing.T, channel <-chan *gen.SessionEvent, sessionID string) {
	t.Helper()
	select {
	case event := <-channel:
		if event == nil || event.SessionID != sessionID || event.Code != gen.SessionEventCodeOrganizationSuspended {
			t.Fatalf("unexpected event: %#v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("session event timeout")
	}
}

func assertNoSessionEvent(t *testing.T, channel <-chan *gen.SessionEvent) {
	t.Helper()
	select {
	case event := <-channel:
		t.Fatalf("unrelated event: %#v", event)
	case <-time.After(30 * time.Millisecond):
	}
}

func assertNoErrors(t *testing.T, response graphQLResponse) {
	t.Helper()
	if len(response.Errors) > 0 || response.Data == nil {
		t.Fatalf("GraphQL failed: %s", response.Body)
	}
}

func assertResultTotal(t *testing.T, response graphQLResponse, field string, expected int) {
	t.Helper()
	var result struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(response.Data[field], &result); err != nil || result.Total != expected {
		t.Fatalf("%s total=%d want=%d err=%v body=%s", field, result.Total, expected, err, response.Body)
	}
}

func decodeData(t *testing.T, response graphQLResponse, target any) {
	t.Helper()
	encoded, err := json.Marshal(response.Data)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(encoded, target); err != nil {
		t.Fatal(err)
	}
}

func assertAuditCount(t *testing.T, fixture *securityFixture, action string, expected int64) {
	t.Helper()
	var count int64
	if err := fixture.db.Model(&gen.AuditLog{}).Where("action = ?", action).Count(&count).Error; err != nil || count != expected {
		t.Fatalf("audit %s count=%d want=%d err=%v", action, count, expected, err)
	}
}
