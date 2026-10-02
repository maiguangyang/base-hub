/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-22
 */
package integration_test

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestHQAdministratorAndRoleWritesUseHqAuditActions(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedHQCustomAdministrator(fixture, "audit-target", []string{"hqMembership:read"})
	fixture.createAll([]gen.Account{{ID: "account-hq-invite", Phone: "13900000008", DisplayName: "Invite", Status: gen.AccountStatusActive, CredentialVersion: 1}})

	invite := fixture.execute("session-hq", `mutation {
		inviteOperator(input: {phone: "13900000008", displayName: "Ignored", roleIds: ["role-hq-audit-target"], storeAccessMode: ALL_STORES, storeIds: []}) {
			membership { id }
		}
	}`)
	assertNoErrors(t, invite)
	assertAuditCount(t, fixture, "hqAdministrator:invite", 1)

	updateMember := fixture.execute("session-hq", `mutation {
		updateOperatorMembership(id: "membership-hq-audit-target", input: {rolesIds: ["role-hq-audit-target"]}) { id }
	}`)
	assertNoErrors(t, updateMember)
	assertAuditCount(t, fixture, "hqAdministrator:assign_roles", 1)
	assertHQAuditMetadata(t, fixture, "hqAdministrator:assign_roles", "membership-hq", "roleIds", "role-hq-audit-target")

	updateRole := fixture.execute("session-hq", `mutation {
		updateOperatorRole(id: "role-hq-audit-target", input: {permissionsIds: ["permission-hqMembership-read"]}) { id }
	}`)
	assertNoErrors(t, updateRole)
	assertAuditCount(t, fixture, "hqRole:update", 1)
	assertHQAuditMetadata(t, fixture, "hqRole:update", "membership-hq", "permissionIds", "permission-hqMembership-read")
}

func TestHQAdministratorDenialWritesStructuredSecurityLog(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedHQCustomAdministrator(fixture, "denied-logger", []string{"hqMembership:update"})
	var output bytes.Buffer
	previousOutput := log.Writer()
	previousFlags := log.Flags()
	log.SetOutput(&output)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(previousOutput)
		log.SetFlags(previousFlags)
	})

	response := fixture.execute("session-hq-denied-logger", `mutation {
		updateOperatorMembership(id: "membership-hq-denied-logger", input: {rolesIds: ["role-hq-denied-logger"]}) { id }
	}`)
	assertCode(t, response, auth.Code("SELF_MEMBERSHIP_CHANGE_DENIED"))
	logged := output.String()
	for _, fragment := range []string{
		"security_event=authorization_denied",
		"actor_account_id=account-hq-denied-logger",
		"actor_membership_id=membership-hq-denied-logger",
		"action=hqAdministrator:assign_roles",
		"resource_id=membership-hq-denied-logger",
		"result_code=SELF_MEMBERSHIP_CHANGE_DENIED",
	} {
		if !strings.Contains(logged, fragment) {
			t.Fatalf("security log missing %q: %s", fragment, logged)
		}
	}
}

func assertHQAuditMetadata(t *testing.T, fixture *securityFixture, action, actorMembershipID, listKey, listValue string) {
	t.Helper()
	var record gen.AuditLog
	if err := fixture.db.Where("action = ?", action).Order("created_at DESC").First(&record).Error; err != nil {
		t.Fatal(err)
	}
	if record.MetadataJSON == nil {
		t.Fatalf("audit %s metadata missing", action)
	}
	var metadata map[string]any
	if err := json.Unmarshal([]byte(*record.MetadataJSON), &metadata); err != nil {
		t.Fatal(err)
	}
	if metadata["actorMembershipId"] != actorMembershipID || metadata["requestId"] == "" {
		t.Fatalf("audit %s principal metadata = %#v", action, metadata)
	}
	values, ok := metadata[listKey].([]any)
	if !ok || len(values) != 1 || values[0] != listValue {
		t.Fatalf("audit %s %s = %#v", action, listKey, metadata[listKey])
	}
}
