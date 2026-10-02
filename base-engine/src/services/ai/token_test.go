package ai

import (
	"encoding/json"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"base-engine/auth"
)

func TestTokenConcurrentConsumeHasOneWinner(t *testing.T) {
	service, principal := approvalFixture(t)
	plan := issueFixturePlan(t, service, principal)
	token, err := service.Issue(plan)
	if err != nil {
		t.Fatal(err)
	}
	var winners atomic.Int32
	var group sync.WaitGroup
	for range 20 {
		group.Go(func() {
			if _, err := service.Consume(token, principal); err == nil {
				winners.Add(1)
			}
		})
	}
	group.Wait()
	if winners.Load() != 1 {
		t.Fatalf("consume winners = %d", winners.Load())
	}
}

func TestTokenRejectsRevokedPermission(t *testing.T) {
	service, principal := approvalFixture(t)
	plan := issueFixturePlan(t, service, principal)
	token, err := service.Issue(plan)
	if err != nil {
		t.Fatal(err)
	}
	delete(principal.Permissions, "store:create")
	if _, err := service.Consume(token, principal); err == nil {
		t.Fatal("revoked permission accepted")
	}
}

func issueFixturePlan(t *testing.T, service *ApprovalService, principal *auth.WorkspacePrincipal) Plan {
	t.Helper()
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), MaxCalls: 1, Sequence: 1}}}
	draft.BindPromptHash("prompt-sha256")
	plan, _, err := service.ValidateDraft(draft, principal, nil)
	if err != nil {
		t.Fatal(err)
	}
	return plan
}

func TestTokenBindsSessionExpiresAndCannotReplay(t *testing.T) {
	service, principal := approvalFixture(t)
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), MaxCalls: 1, Sequence: 1}}}
	draft.BindPromptHash("prompt-sha256")
	plan, _, err := service.ValidateDraft(draft, principal, nil)
	if err != nil {
		t.Fatal(err)
	}
	token, err := service.Issue(plan)
	if err != nil {
		t.Fatal(err)
	}
	other := *principal
	other.SessionID = "session-2"
	if _, err := service.Consume(token, &other); err == nil {
		t.Fatal("cross-session token accepted")
	}
	if _, err := service.Consume(token, principal); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Consume(token, principal); err == nil {
		t.Fatal("replayed token accepted")
	}
	service.now = func() time.Time { return plan.ExpiresAt.Add(time.Second) }
	if _, err := service.Issue(plan); err == nil {
		t.Fatal("expired plan reissued")
	}
}

func TestTokenRejectsTamperingAndRestart(t *testing.T) {
	service, principal := approvalFixture(t)
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), MaxCalls: 1, Sequence: 1}}}
	draft.BindPromptHash("prompt-sha256")
	plan, _, err := service.ValidateDraft(draft, principal, nil)
	if err != nil {
		t.Fatal(err)
	}
	token, err := service.Issue(plan)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Consume(token+"x", principal); err == nil {
		t.Fatal("tampered token accepted")
	}
	restarted, err := NewApprovalService(service.catalog)
	if err != nil {
		t.Fatal(err)
	}
	restarted.now = service.now
	if _, err := restarted.Consume(token, principal); err == nil {
		t.Fatal("token survived process restart")
	}
}

func TestTokenRejectsUnvalidatedAndMutatedPlans(t *testing.T) {
	service, principal := approvalFixture(t)
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), MaxCalls: 1, Sequence: 1}}}
	draft.BindPromptHash("prompt-sha256")
	plan, _, err := service.ValidateDraft(draft, principal, nil)
	if err != nil {
		t.Fatal(err)
	}
	forged := plan
	forged.ID = "forged"
	if _, err := service.Issue(forged); err == nil {
		t.Fatal("unvalidated plan issued")
	}
	changed := plan
	changed.Steps = append([]ApprovedStep(nil), plan.Steps...)
	changed.Steps[0].Arguments = json.RawMessage(`{"input":{"organizationId":"org-a","name":"Other"}}`)
	if _, err := service.Issue(changed); err == nil {
		t.Fatal("mutated plan issued")
	}
	if _, err := service.Issue(plan); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Issue(plan); err == nil {
		t.Fatal("same plan issued twice")
	}
}
