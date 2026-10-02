package ai

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"base-engine/auth"
)

func TestPlanLocksCreationScopeAndCompleteArguments(t *testing.T) {
	service, principal := approvalFixture(t)
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), ScopeIDs: []string{"org-a"}, MaxCalls: 1, Sequence: 1}}}
	draft.BindPromptHash("prompt-sha256")
	plan, required, err := service.ValidateDraft(draft, principal, nil)
	if err != nil || len(required) != 0 || len(plan.Steps) != 1 {
		t.Fatalf("plan = %+v, %v, %v", plan, required, err)
	}
	if len(plan.Steps[0].TargetIDs) != 0 || len(plan.Steps[0].ScopeIDs) != 1 {
		t.Fatalf("creation scope = %+v", plan.Steps[0])
	}
	token, err := service.Issue(plan)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Consume(token, principal); err != nil {
		t.Fatal(err)
	}
	assertCreateStepLocked(t, service, plan)
}

func TestPlanRejectsInvalidNestedToolArgumentsBeforeApproval(t *testing.T) {
	service, principal := approvalFixture(t)
	for _, arguments := range []string{
		`{"input":{"organizationId":"org-a","name":42}}`,
		`{"input":{"organizationId":"org-a"}}`,
		`{"input":{"organizationId":"org-a","name":"New","unreviewed":true}}`,
	} {
		draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(arguments), ScopeIDs: []string{"org-a"}, MaxCalls: 1, Sequence: 1}}}
		draft.BindPromptHash("prompt-sha256")
		if _, _, err := service.ValidateDraft(draft, principal, nil); err == nil {
			t.Fatalf("invalid tool arguments were approved: %s", arguments)
		}
	}
}

func TestPlanGlobalCreationUsesWorkspaceScope(t *testing.T) {
	service, principal := approvalFixture(t)
	principal.OrganizationID = nil
	service.catalog.specs[0].ScopeFields = nil
	spec := service.catalog.byID["create_store"]
	spec.ScopeFields = nil
	service.catalog.byID["create_store"] = spec
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(`{"input":{"name":"Global"}}`), MaxCalls: 1, Sequence: 1}}}
	draft.BindPromptHash("prompt-sha256")
	plan, _, err := service.ValidateDraft(draft, principal, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Steps[0].ScopeIDs) != 1 || plan.Steps[0].ScopeIDs[0] != "workspace:HEADQUARTERS" {
		t.Fatalf("scope = %+v", plan.Steps[0].ScopeIDs)
	}
}

func assertCreateStepLocked(t *testing.T, service *ApprovalService, plan Plan) {
	t.Helper()
	if err := service.AuthorizeStep(plan.ID, "create_store", json.RawMessage(`{"input":{"name":"Other","organizationId":"org-a"}}`)); err == nil {
		t.Fatal("changed creation field accepted")
	}
	if err := service.AuthorizeStep(plan.ID, "create_store", json.RawMessage(`{"input":{"name":"New","organizationId":"org-b"}}`)); err == nil {
		t.Fatal("changed parent scope accepted")
	}
	if err := service.AuthorizeStep(plan.ID, "create_store", json.RawMessage(`{"input":{"name":"New","organizationId":"org-a"}}`)); err != nil {
		t.Fatal(err)
	}
	if err := service.AuthorizeStep(plan.ID, "create_store", plan.Steps[0].Arguments); err == nil {
		t.Fatal("repeated write accepted")
	}
}

func TestPlanRejectsModelTargetsAndUnknownArguments(t *testing.T) {
	service, principal := approvalFixture(t)
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "suspend", OperationID: "graphql.mutation.suspendOrganization", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","reasonCode":"FRAUD"}}`), TargetIDs: []string{"org-b"}, MaxCalls: 1, Sequence: 1}}}
	draft.BindPromptHash("prompt-sha256")
	if _, _, err := service.ValidateDraft(draft, principal, nil); err == nil {
		t.Fatal("model supplied target accepted")
	}
	draft.Steps[0].TargetIDs = []string{"org-a"}
	draft.Steps[0].Arguments = json.RawMessage(`{"input":{"organizationId":"org-a","reasonCode":"FRAUD"},"path":"/graphql"}`)
	if _, _, err := service.ValidateDraft(draft, principal, nil); err == nil {
		t.Fatal("arbitrary path accepted")
	}
	draft.Steps[0].Arguments = json.RawMessage(`{"input":{"organizationId":"org-a","reasonCode":"FRAUD"}}`)
	plan, _, err := service.ValidateDraft(draft, principal, nil)
	if err != nil || len(plan.Steps[0].TargetIDs) != 1 {
		t.Fatalf("valid target = %+v, %v", plan, err)
	}
}

func TestPlanRequiresUserAttestation(t *testing.T) {
	service, principal := approvalFixture(t)
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "initial_account", OperationID: "http.franchiseInitialAccount.set", Arguments: json.RawMessage(`{"organizationId":"org-a","accountId":"account-1"}`), TargetIDs: []string{"account-1"}, MaxCalls: 1, Sequence: 1}}}
	draft.BindPromptHash("prompt-sha256")
	plan, required, err := service.ValidateDraft(draft, principal, nil)
	if err != nil || plan.ID != "" || len(required) != 2 {
		t.Fatalf("attestation gate = %+v, %v, %v", plan, required, err)
	}
	if _, _, err := service.ValidateDraft(draft, principal, &UserAttestation{EvidenceReference: "record-1", Attested: true}); err != nil {
		t.Fatal(err)
	}
	draft.Steps[0].Arguments = json.RawMessage(`{"organizationId":"org-a","accountId":"account-1","attested":true}`)
	if _, _, err := service.ValidateDraft(draft, principal, &UserAttestation{EvidenceReference: "record-1", Attested: true}); err == nil {
		t.Fatal("model attestation accepted")
	}
}

func TestPlanRejectsOrderAndPermissionDrift(t *testing.T) {
	service, principal := approvalFixture(t)
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), MaxCalls: 1, Sequence: 2}}}
	draft.BindPromptHash("prompt-sha256")
	if _, _, err := service.ValidateDraft(draft, principal, nil); err == nil {
		t.Fatal("noninitial sequence accepted")
	}
	draft.Steps[0].Sequence = 1
	delete(principal.Permissions, "store:create")
	if _, _, err := service.ValidateDraft(draft, principal, nil); err == nil {
		t.Fatal("revoked permission accepted")
	}
}

func TestPlanStepOrderAndDuplicateDraft(t *testing.T) {
	service, principal := approvalFixture(t)
	first := ApprovedStep{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), MaxCalls: 2, Sequence: 1}
	second := ApprovedStep{ToolID: "suspend", OperationID: "graphql.mutation.suspendOrganization", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","reasonCode":"FRAUD"}}`), MaxCalls: 1, Sequence: 2}
	draft := PlanDraft{Steps: []ApprovedStep{first, second}}
	draft.BindPromptHash("prompt-sha256")
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
	assertPlanRunOrder(t, service, plan)
	draft.Steps = []ApprovedStep{first, first}
	draft.Steps[1].Sequence = 2
	if _, _, err := service.ValidateDraft(draft, principal, nil); err == nil {
		t.Fatal("duplicate step accepted")
	}
}

func assertPlanRunOrder(t *testing.T, service *ApprovalService, plan Plan) {
	t.Helper()
	if err := service.AuthorizeStep(plan.ID, "suspend", plan.Steps[1].Arguments); err == nil {
		t.Fatal("second step ran first")
	}
	if err := service.AuthorizeStep(plan.ID, "create_store", plan.Steps[0].Arguments); err != nil {
		t.Fatal(err)
	}
	if err := service.RecordSuccess(plan.ID, "create_store"); err != nil {
		t.Fatal(err)
	}
	if err := service.AuthorizeStep(plan.ID, "suspend", plan.Steps[1].Arguments); err != nil {
		t.Fatal(err)
	}
	if err := service.RecordSuccess(plan.ID, "suspend"); err != nil {
		t.Fatal(err)
	}
	if err := service.Complete(plan.ID); err != nil {
		t.Fatal(err)
	}
	if err := service.AuthorizeStep(plan.ID, "create_store", plan.Steps[0].Arguments); err == nil {
		t.Fatal("first step ran after second")
	}
}

func TestPlanCannotFinishWithUnexecutedStepsOrWriteAfterCancel(t *testing.T) {
	service, principal := approvalFixture(t)
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "suspend", OperationID: "graphql.mutation.suspendOrganization", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","reasonCode":"FRAUD"}}`), MaxCalls: 1, Sequence: 1}}}
	draft.BindPromptHash("prompt-sha256")
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
	if err := service.Complete(plan.ID); err == nil {
		t.Fatal("empty execution reported success")
	}
	service.Cancel(plan.ID)
	if err := service.AuthorizeStep(plan.ID, "suspend", plan.Steps[0].Arguments); err == nil {
		t.Fatal("cancelled plan accepted a write")
	}
}

func approvalFixture(t *testing.T) (*ApprovalService, *auth.WorkspacePrincipal) {
	t.Helper()
	contracts := []ContractRecord{
		{OperationID: "graphql.mutation.createStore", Protocol: "GRAPHQL_MUTATION", Availability: "CALLABLE", Classifications: []string{"HEADQUARTERS_ADMIN"}},
		{OperationID: "graphql.mutation.suspendOrganization", Protocol: "GRAPHQL_MUTATION", Availability: "CALLABLE", Classifications: []string{"HEADQUARTERS_ADMIN"}},
		{OperationID: "http.franchiseInitialAccount.set", Protocol: "HTTP_ROUTE", Path: "/api/franchise-initial-account", Availability: "CALLABLE", Classifications: []string{"HEADQUARTERS_ADMIN"}},
	}
	specs := []ToolSpec{
		{ID: "create_store", OperationID: contracts[0].OperationID, Document: "mutation { createStore { id } }", Description: "Create store", Mode: ModeWrite, Risk: "MEDIUM", Permission: "store:create", Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeHeadquarters}, OutputFields: []string{"id"}, WriteKind: WriteCreate, ScopeFields: []string{"input.organizationId"}, ArgumentFields: []string{"input"}},
		{ID: "suspend", OperationID: contracts[1].OperationID, Document: "mutation { suspendOrganization { id } }", Description: "Suspend", Mode: ModeWrite, Risk: "HIGH", Permission: "organization:suspend", Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeHeadquarters}, OutputFields: []string{"id"}, WriteKind: WriteExisting, TargetFields: []string{"input.organizationId"}, ArgumentFields: []string{"input"}},
		{ID: "initial_account", OperationID: contracts[2].OperationID, Path: "/api/franchise-initial-account", Description: "Set initial account", Mode: ModeWrite, Risk: "HIGH", Permission: "account:update", Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeHeadquarters}, OutputFields: []string{"accountId"}, WriteKind: WriteExisting, TargetFields: []string{"accountId"}, ArgumentFields: []string{"organizationId", "accountId"}, RequiredAttestation: true},
	}
	for index := range specs {
		var schema *jsonschema.Schema
		if index == 2 {
			schema = &jsonschema.Schema{Type: "object", Properties: map[string]*jsonschema.Schema{"organizationId": {Type: "string"}, "accountId": {Type: "string"}}, Required: []string{"organizationId", "accountId"}, AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}}}
		} else {
			schema = &jsonschema.Schema{Type: "object", Properties: map[string]*jsonschema.Schema{"input": {Type: "object", Properties: map[string]*jsonschema.Schema{"organizationId": {Type: "string"}, "name": {Type: "string"}, "reasonCode": {Type: "string"}}, AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}}}}, Required: []string{"input"}, AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}}}
			if index == 0 {
				schema.Properties["input"].Required = []string{"name"}
			} else {
				schema.Properties["input"].Required = []string{"organizationId", "reasonCode"}
			}
		}
		resolved, err := schema.Resolve(nil)
		if err != nil {
			t.Fatal(err)
		}
		specs[index].InputSchema = resolved
	}
	catalog, err := NewCatalog(contracts, specs)
	if err != nil {
		t.Fatal(err)
	}
	service, err := NewApprovalService(catalog)
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return time.Unix(1_800_000_000, 0) }
	organizationID := "org-a"
	principal := &auth.WorkspacePrincipal{AccountID: "account-1", SessionID: "session-1", WorkspaceType: auth.WorkspaceTypeHeadquarters, OrganizationID: &organizationID, Permissions: map[string]struct{}{"store:create": {}, "organization:suspend": {}, "account:update": {}}}
	return service, principal
}
