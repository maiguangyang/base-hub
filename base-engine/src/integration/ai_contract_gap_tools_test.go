package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"base-engine/gen"
	"base-engine/src/services/ai"
	aitools "base-engine/src/services/ai/tools"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool/toolconfirmation"
)

type protectedGapToolContext struct{ agent.Context }

func (protectedGapToolContext) ToolConfirmation() *toolconfirmation.ToolConfirmation { return nil }

type protectedGapReadCase struct {
	tool, session, root, want, forbidden string
	args                                 map[string]any
}

func TestContractGapReadToolsUseProtectedHandler(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedContractGapTargets(fixture)
	for _, tc := range protectedGapReadCases() {
		t.Run(tc.tool, func(t *testing.T) { testProtectedGapRead(t, fixture, tc) })
	}
}

func testProtectedGapRead(t *testing.T, fixture *securityFixture, tc protectedGapReadCase) {
	result, err := runProtectedGapTool(t, fixture, tc.session, tc.tool, tc.args)
	if err != nil {
		t.Fatal(err)
	}
	value, ok := result[tc.root]
	if !ok || value == nil {
		t.Fatalf("missing protected result %s: %#v", tc.root, result)
	}
	encoded, _ := json.Marshal(value)
	if !bytes.Contains(encoded, []byte(tc.want)) || tc.forbidden != "" && bytes.Contains(encoded, []byte(tc.forbidden)) {
		t.Fatalf("unexpected scoped result for %s: %s", tc.tool, encoded)
	}
	shape := gapExpectedShape(tc.tool, tc.root)
	if shape == nil || !sameGapShape(value, shape) {
		t.Fatalf("unexpected output fields for %s: %s", tc.tool, encoded)
	}
	assertGapReadInvalidArguments(t, fixture, tc)
}

func assertGapReadInvalidArguments(t *testing.T, fixture *securityFixture, tc protectedGapReadCase) {
	t.Helper()
	invalid := map[string]any{"id": []string{"wrong-type"}}
	if _, paginated := tc.args["page"]; paginated {
		invalid = map[string]any{"page": "wrong-type", "pageSize": 10}
	}
	if output, err := runProtectedGapTool(t, fixture, tc.session, tc.tool, invalid); err == nil {
		t.Fatalf("invalid tool arguments were accepted: %#v", output)
	}
}

func TestContractGapReadToolsRejectForeignTargets(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedContractGapTargets(fixture)
	for _, tc := range []protectedGapReadCase{
		{tool: "HqMembershipInvitation", session: "session-hq", args: map[string]any{"id": "invitation-a"}},
		{tool: "FranchiseAccount", session: "session-a-1", args: map[string]any{"id": "account-foreign"}},
		{tool: "FranchiseAuditLog", session: "session-a-1", args: map[string]any{"id": "audit-b"}},
		{tool: "FranchiseMembershipInvitation", session: "session-a-1", args: map[string]any{"id": "invitation-b"}},
		{tool: "FranchiseOperatorMembership", session: "session-a-1", args: map[string]any{"id": "membership-owner-b"}},
		{tool: "FranchiseOperatorRole", session: "session-a-1", args: map[string]any{"id": "role-owner-b"}},
		{tool: "FranchisePermission", session: "session-a-1", args: map[string]any{"id": "permission-hqMembership-read"}},
		{tool: "FranchiseStore", session: "session-a-1", args: map[string]any{"id": "store-b"}},
	} {
		t.Run(tc.tool, func(t *testing.T) {
			if output, err := runProtectedGapTool(t, fixture, tc.session, tc.tool, tc.args); err == nil {
				t.Fatalf("foreign target was returned: %#v", output)
			}
		})
	}
}

func seedContractGapTargets(fixture *securityFixture) {
	fixture.createAll(
		[]gen.MembershipInvitation{{ID: "invitation-hq", MembershipID: "membership-hq", InvitedByAccountID: "account-shared", ExpiresAt: time.Now().Add(time.Hour)}},
		[]gen.Account{{ID: "account-foreign", Phone: "13800000003", DisplayName: "Foreign", Status: gen.AccountStatusActive, CredentialVersion: 1}},
		[]gen.OperatorMembership{{ID: "membership-foreign", AccountID: "account-foreign", OrganizationID: "org-b", Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores}},
	)
}

func runProtectedGapTool(t *testing.T, fixture *securityFixture, sessionID, toolID string, args map[string]any) (map[string]any, error) {
	t.Helper()
	runtime := ai.FixedToolRuntime{
		Call: func(_ agent.Context, spec ai.ToolSpec, variables json.RawMessage) (ai.FixedResponse, error) {
			return ai.ProtectedCall(t.Context(), fixture.handler, ai.FixedRequest{Method: http.MethodPost, Path: "/graphql", Document: spec.Document,
				Variables: variables, Cookie: fixture.cookies[sessionID].String(), Origin: "https://admin.example.com"})
		},
		DeliverSecret: func(agent.Context, ai.SecretPayload) error { t.Fatal("unexpected secret delivery"); return nil },
	}
	items, err := aitools.Build(runtime)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.Name() == toolID {
			return item.(interface {
				Run(agent.Context, any) (map[string]any, error)
			}).Run(protectedGapToolContext{}, args)
		}
	}
	t.Fatalf("missing tool %s", toolID)
	return nil, nil
}

func protectedGapReadCases() []protectedGapReadCase {
	return []protectedGapReadCase{
		{"HqAccount", "session-hq", "account", "account-shared", "", map[string]any{"id": "account-shared"}},
		{"FranchiseAccount", "session-a-1", "account", "account-shared", "", map[string]any{"id": "account-shared"}},
		{"HqAccounts", "session-hq", "accounts", "account-shared", "", map[string]any{"page": 1, "pageSize": 10}},
		{"FranchiseAccounts", "session-a-1", "accounts", "account-staff", "account-foreign", map[string]any{"page": 1, "pageSize": 10}},
		{"HqAuditLog", "session-hq", "auditLog", "audit-a", "", map[string]any{"id": "audit-a"}},
		{"FranchiseAuditLog", "session-a-1", "auditLog", "audit-a", "", map[string]any{"id": "audit-a"}},
		{"HqMembershipInvitation", "session-hq", "membershipInvitation", "invitation-hq", "", map[string]any{"id": "invitation-hq"}},
		{"FranchiseMembershipInvitation", "session-a-1", "membershipInvitation", "invitation-a", "", map[string]any{"id": "invitation-a"}},
		{"HqMembershipInvitations", "session-hq", "membershipInvitations", "invitation-hq", "invitation-a", map[string]any{"page": 1, "pageSize": 10}},
		{"FranchiseMembershipInvitations", "session-a-1", "membershipInvitations", "invitation-a", "invitation-b", map[string]any{"page": 1, "pageSize": 10}},
		{"HqOperatorMembership", "session-hq", "operatorMembership", "membership-hq", "", map[string]any{"id": "membership-hq"}},
		{"FranchiseOperatorMembership", "session-a-1", "operatorMembership", "membership-staff-a", "", map[string]any{"id": "membership-staff-a"}},
		{"HqOperatorRole", "session-hq", "operatorRole", "role-hq", "", map[string]any{"id": "role-hq"}},
		{"FranchiseOperatorRole", "session-a-1", "operatorRole", "role-owner-a", "", map[string]any{"id": "role-owner-a"}},
		{"HqPermission", "session-hq", "permission", "permission-hqMembership-read", "", map[string]any{"id": "permission-hqMembership-read"}},
		{"FranchisePermission", "session-a-1", "permission", "permission-operatorRole-read", "", map[string]any{"id": "permission-operatorRole-read"}},
		{"HqPermissionByGrant", "session-hq", "permission", "permission-hqMembership-read", "", map[string]any{"id": "permission-hqMembership-read"}},
		{"HqStore", "session-hq", "store", "store-hq", "", map[string]any{"id": "store-hq"}},
		{"HqDirectStore", "session-hq", "store", "store-hq", "", map[string]any{"id": "store-hq"}},
		{"FranchiseStore", "session-a-1", "store", "store-a", "", map[string]any{"id": "store-a"}},
	}
}
