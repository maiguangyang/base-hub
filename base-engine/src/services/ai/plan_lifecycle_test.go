package ai

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCancelledPlanReleasesStoredPromptAndPlan(t *testing.T) {
	service, principal := approvalFixture(t)
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), MaxCalls: 1, Sequence: 1}}}
	draft.BindPrompt("Create a store with private notes")
	plan, _, err := service.ValidateDraft(draft, principal, nil)
	if err != nil {
		t.Fatal(err)
	}
	service.Cancel(plan.ID)
	if service.store.Has(plan.ID) {
		t.Fatal("cancelled plan and prompt remain in memory")
	}
}

func TestExpiredUnconsumedPlanReleasesStoredPromptAndPlan(t *testing.T) {
	service, principal := approvalFixture(t)
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), MaxCalls: 1, Sequence: 1}}}
	draft.BindPrompt("Create a store with private notes")
	plan, _, err := service.ValidateDraft(draft, principal, nil)
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return plan.ExpiresAt.Add(time.Second) }
	service.expireUnusedPlan(plan.ID)
	if service.store.Has(plan.ID) {
		t.Fatal("expired unused plan and prompt remain in memory")
	}
}

func TestConsumedPlanCanExecuteAfterPreviewTokenExpires(t *testing.T) {
	service, principal := approvalFixture(t)
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), MaxCalls: 1, Sequence: 1}}}
	draft.BindPrompt("Create a store")
	plan, _, err := service.ValidateDraft(draft, principal, nil)
	if err != nil {
		t.Fatal(err)
	}
	token, err := service.Issue(plan)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Consume(token, principal); err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return plan.ExpiresAt.Add(time.Second) }
	if err := service.AuthorizeStep(plan.ID, "create_store", plan.Steps[0].Arguments); err != nil {
		t.Fatalf("consumed plan expired during execution: %v", err)
	}
}
