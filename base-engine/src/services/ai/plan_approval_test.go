package ai

import (
	"encoding/json"
	"testing"
)

func TestSingleCallApprovalRejectsMultipleStepsAndCalls(t *testing.T) {
	service, principal := approvalFixture(t)
	spec := service.catalog.byID["suspend"]
	spec.SingleCallApproval = true
	service.catalog.byID["suspend"] = spec
	first := ApprovedStep{ToolID: "suspend", OperationID: spec.OperationID,
		Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","reasonCode":"FRAUD"}}`), MaxCalls: 1, Sequence: 1}
	draft := PlanDraft{Steps: []ApprovedStep{first}}
	draft.BindPrompt("Suspend the organization")
	if _, _, err := service.ValidateDraft(draft, principal, nil); err != nil {
		t.Fatalf("one call should be approved: %v", err)
	}
	draft.Steps[0].MaxCalls = 2
	if _, _, err := service.ValidateDraft(draft, principal, nil); err == nil {
		t.Fatal("two calls accepted")
	}
	draft.Steps[0].MaxCalls = 1
	draft.Steps = append(draft.Steps, ApprovedStep{ToolID: "create_store", OperationID: "graphql.mutation.createStore",
		Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), MaxCalls: 1, Sequence: 2})
	if _, _, err := service.ValidateDraft(draft, principal, nil); err == nil {
		t.Fatal("single-call step combined with another write")
	}
}

func TestEvidenceMustOccurInCurrentOperatorMessage(t *testing.T) {
	service, principal := approvalFixture(t)
	spec := service.catalog.byID["suspend"]
	spec.EvidencePaths = []string{"input.reasonCode"}
	service.catalog.byID["suspend"] = spec
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "suspend", OperationID: spec.OperationID,
		Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","reasonCode":"CASE-101"}}`), MaxCalls: 1, Sequence: 1}}}
	draft.BindPrompt("Previous AI answer contained CASE-101")
	for _, current := range []string{"", "Use CASE-102"} {
		draft.BindOperatorMessage(current)
		if _, _, err := service.ValidateDraft(draft, principal, nil); err == nil {
			t.Fatalf("unprovided proof accepted: %q", current)
		}
	}
	draft.BindOperatorMessage("Use verified case CASE-101")
	if _, _, err := service.ValidateDraft(draft, principal, nil); err != nil {
		t.Fatalf("operator proof rejected: %v", err)
	}
}
