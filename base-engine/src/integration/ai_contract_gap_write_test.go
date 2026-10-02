package integration_test

import (
	"encoding/json"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/ai"
	aitools "base-engine/src/services/ai/tools"
)

type protectedGapWriteCase struct {
	tool, session, ownID, foreignID, action string
	permissions                             []string
}

func contractGapWriteCases() []protectedGapWriteCase {
	return []protectedGapWriteCase{
		{"HqDeleteMembershipInvitations", "session-hq", "invitation-hq", "invitation-a", "hqMembership:delete", []string{"hqMembership:read", "hqMembership:delete"}},
		{"FranchiseDeleteMembershipInvitations", "session-a-1", "invitation-a", "invitation-b", "membershipInvitation:delete", []string{"operatorMembership:read", "membershipInvitation:delete"}},
	}
}

func TestContractGapDeleteToolsRequireBoundApprovalAndProtectedScope(t *testing.T) {
	for _, tc := range contractGapWriteCases() {
		t.Run(tc.tool, func(t *testing.T) { testContractGapDeleteTool(t, tc) })
	}
}

func testContractGapDeleteTool(t *testing.T, tc protectedGapWriteCase) {
	fixture := newSecurityFixture(t)
	seedContractGapTargets(fixture)
	approval, principal := contractGapApproval(t, fixture, tc)
	invalid := contractGapDeleteDraft(tc.tool, tc.ownID, tc.foreignID)
	if _, _, err := approval.ValidateDraft(invalid, principal, nil); err == nil {
		t.Fatal("approval accepted target IDs different from mutation arguments")
	}
	invalid = contractGapDeleteDraft(tc.tool, tc.ownID, tc.ownID)
	invalid.Steps[0].Arguments = json.RawMessage(`{"ids":"wrong-type"}`)
	if _, _, err := approval.ValidateDraft(invalid, principal, nil); err == nil {
		t.Fatal("approval accepted malformed invitation IDs")
	}
	for _, id := range []string{tc.foreignID, tc.ownID} {
		plan := approvedContractGapDelete(t, approval, principal, tc.tool, id)
		result, err := runProtectedGapTool(t, fixture, tc.session, tc.tool, map[string]any{"ids": []string{id}})
		if id == tc.foreignID {
			if err == nil {
				t.Fatalf("foreign invitation was deleted: %#v", result)
			}
			assertInvitationDeleteState(t, fixture, id, false)
			assertAuditCount(t, fixture, tc.action, 0)
			continue
		}
		if err != nil || result["deleteMembershipInvitations"] != true {
			t.Fatalf("approved deletion failed: result=%#v err=%v", result, err)
		}
		if err := approval.RecordSuccess(plan.ID, tc.tool); err != nil {
			t.Fatal(err)
		}
		assertInvitationDeleteState(t, fixture, id, true)
		assertAuditCount(t, fixture, tc.action, 1)
	}
}

func contractGapApproval(t *testing.T, fixture *securityFixture, tc protectedGapWriteCase) (*ai.ApprovalService, *auth.WorkspacePrincipal) {
	t.Helper()
	all, err := aitools.PreparedSpecs()
	if err != nil {
		t.Fatal(err)
	}
	var specs []ai.ToolSpec
	for _, spec := range all {
		if spec.OperationID == "graphql.mutation.deleteMembershipInvitations" {
			specs = append(specs, spec)
		}
	}
	contract := ai.ContractRecord{OperationID: "graphql.mutation.deleteMembershipInvitations", Protocol: "GRAPHQL_MUTATION", Availability: "CALLABLE", Classifications: []string{"HEADQUARTERS_ADMIN", "FRANCHISE_ADMIN"}}
	catalog, err := ai.NewCatalog([]ai.ContractRecord{contract}, specs)
	if err != nil {
		t.Fatal(err)
	}
	approval, err := ai.NewApprovalService(catalog)
	if err != nil {
		t.Fatal(err)
	}
	copy := *fixture.principals[tc.session]
	copy.Permissions = map[string]struct{}{}
	for _, permission := range tc.permissions {
		copy.Permissions[permission] = struct{}{}
	}
	return approval, &copy
}

func contractGapDeleteDraft(toolID, argumentID, targetID string) ai.PlanDraft {
	arguments, _ := json.Marshal(map[string]any{"ids": []string{argumentID}})
	draft := ai.PlanDraft{Steps: []ai.ApprovedStep{{OperationID: "graphql.mutation.deleteMembershipInvitations", ToolID: toolID,
		Arguments: arguments, TargetIDs: []string{targetID}, MaxCalls: 1, Sequence: 1}}}
	draft.BindPrompt("Delete invitation " + argumentID)
	return draft
}

func approvedContractGapDelete(t *testing.T, service *ai.ApprovalService, principal *auth.WorkspacePrincipal, toolID, id string) ai.Plan {
	t.Helper()
	plan, required, err := service.ValidateDraft(contractGapDeleteDraft(toolID, id, id), principal, nil)
	if err != nil || len(required) != 0 {
		t.Fatalf("valid invitation plan rejected: required=%v err=%v", required, err)
	}
	token, err := service.Issue(plan)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Consume(token, principal); err != nil {
		t.Fatal(err)
	}
	if err := service.AuthorizeStep(plan.ID, toolID, plan.Steps[0].Arguments); err != nil {
		t.Fatal(err)
	}
	return plan
}

func assertInvitationDeleteState(t *testing.T, fixture *securityFixture, id string, deleted bool) {
	t.Helper()
	var item gen.MembershipInvitation
	if err := fixture.db.Unscoped().First(&item, "id = ?", id).Error; err != nil {
		t.Fatal(err)
	}
	actual := item.IsDelete != nil && *item.IsDelete == 2
	if actual != deleted {
		t.Fatalf("invitation %s deletion=%v want=%v", id, actual, deleted)
	}
}
