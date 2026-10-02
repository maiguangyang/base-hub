package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"base-engine/auth"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/tool/toolconfirmation"
)

type runtimeTestContext struct {
	agent.Context
	base context.Context
}

func (c runtimeTestContext) Deadline() (time.Time, bool)                        { return c.base.Deadline() }
func (c runtimeTestContext) Done() <-chan struct{}                              { return c.base.Done() }
func (c runtimeTestContext) Err() error                                         { return c.base.Err() }
func (c runtimeTestContext) Value(key any) any                                  { return c.base.Value(key) }
func (runtimeTestContext) ToolConfirmation() *toolconfirmation.ToolConfirmation { return nil }

func TestAIFixedRuntimeRejectsPreviewWriteAndArgumentDrift(t *testing.T) {
	service, approval, principal, calls := fixedRuntimeFixture(t)
	spec, _ := approval.catalog.Lookup("create_store")
	args := json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`)
	sink := newEventSink(t.Context())
	state := &runState{service: service, principal: principal, token: "signed", origin: "https://admin.example", phase: PhasePreview, sink: sink}
	ctx := runtimeTestContext{base: context.WithValue(t.Context(), runStateKey{}, state)}
	runtime := service.FixedToolRuntime()
	if _, err := runtime.Call(ctx, spec, args); err == nil {
		t.Fatal("preview write accepted")
	}
	plan := issueFixturePlan(t, approval, principal)
	token, err := approval.Issue(plan)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := approval.Consume(token, principal); err != nil {
		t.Fatal(err)
	}
	state.phase, state.planID = PhaseRun, plan.ID
	state.approvedIDs = map[string]struct{}{"create_store": {}}
	if _, err := runtime.Call(ctx, spec, json.RawMessage(`{"input":{"organizationId":"org-a","name":"Other"}}`)); err == nil {
		t.Fatal("changed write accepted")
	}
	if *calls != 0 {
		t.Fatal("handler called for rejected writes")
	}
	if _, err := runtime.Call(ctx, spec, args); err != nil {
		t.Fatal(err)
	}
	if *calls != 1 {
		t.Fatalf("handler calls = %d", *calls)
	}
}

func TestAIFixedRuntimeInjectsTrustedRequestKeyAfterApproval(t *testing.T) {
	service, approval, principal, _ := fixedRuntimeFixture(t)
	spec := approval.catalog.byID["create_store"]
	spec.RequestKeyPath = "input.requestKey"
	approval.catalog.byID["create_store"] = spec
	var captured map[string]any
	service.protectedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var envelope struct {
			Variables map[string]any `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
			t.Fatal(err)
		}
		captured = envelope.Variables
		_, _ = w.Write([]byte(`{"data":{"createStore":{"id":"store-a"}}}`))
	})
	plan := issueFixturePlan(t, approval, principal)
	token, err := approval.Issue(plan)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := approval.Consume(token, principal); err != nil {
		t.Fatal(err)
	}
	state := &runState{service: service, principal: principal, token: "signed", origin: "https://admin.example",
		phase: PhaseRun, planID: plan.ID, approvedIDs: map[string]struct{}{"create_store": {}}}
	ctx := runtimeTestContext{base: context.WithValue(t.Context(), runStateKey{}, state)}
	if _, err := service.FixedToolRuntime().Call(ctx, spec, plan.Steps[0].Arguments); err != nil {
		t.Fatal(err)
	}
	input := captured["input"].(map[string]any)
	if input["name"] != "New" || input["requestKey"] == "" {
		t.Fatalf("trusted variables = %v", captured)
	}
	if _, err := service.FixedToolRuntime().Call(ctx, spec, plan.Steps[0].Arguments); err == nil {
		t.Fatal("same approval executed twice")
	}
}

func TestAIFixedRuntimeRejectsModelDeactivationBeforeWrite(t *testing.T) {
	service, approval, principal, calls := fixedRuntimeFixture(t)
	plan := issueFixturePlan(t, approval, principal)
	issued, err := approval.Issue(plan)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := approval.Consume(issued, principal); err != nil {
		t.Fatal(err)
	}
	service.config.ModelProvider = func(context.Context) (model.LLM, uint64, error) { return serviceTestModel{}, 2, nil }
	state := &runState{service: service, principal: principal, token: "signed", origin: "https://admin.example", phase: PhaseRun,
		planID: plan.ID, approvedIDs: map[string]struct{}{"create_store": {}}, modelVersion: 1}
	ctx := runtimeTestContext{base: context.WithValue(t.Context(), runStateKey{}, state)}
	spec, _ := approval.catalog.Lookup("create_store")
	if _, err := service.FixedToolRuntime().Call(ctx, spec, plan.Steps[0].Arguments); err == nil {
		t.Fatal("write accepted after model configuration changed")
	}
	if *calls != 0 {
		t.Fatalf("protected handler calls = %d", *calls)
	}
}

func TestAISecretEventIncludesToolID(t *testing.T) {
	service, _, principal, _ := fixedRuntimeFixture(t)
	sink := newEventSink(t.Context())
	state := &runState{service: service, principal: principal, phase: PhaseRun, sink: sink}
	ctx := runtimeTestContext{base: context.WithValue(t.Context(), runStateKey{}, state)}
	if err := service.deliverSecret(ctx, SecretPayload{ToolID: "reset_account", TargetAccountID: "a-1", Value: "temporary"}); err != nil {
		t.Fatal(err)
	}
	event := <-sink.events
	data, ok := event.Data.(map[string]string)
	if event.Name != "secret" || !ok || data["toolId"] != "reset_account" {
		t.Fatalf("secret event = %+v", event)
	}
}

func TestAICustomerPhoneStaysOutOfToolEventsAndDiagnostics(t *testing.T) {
	service, _, principal, _ := fixedRuntimeFixture(t)
	const phone = "13800138000"
	spec := ToolSpec{ID: "HqCustomerSensitivePhone", OperationID: "graphql.query.hqCustomerSensitivePhone",
		Document: "query HqCustomerSensitivePhone($id: ID!) { hqCustomerSensitivePhone(id: $id) }",
		Mode:     ModeReadOnly, Permission: "customer:read_sensitive", Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeHeadquarters}}
	service.config.Catalog.byID[spec.ID] = spec
	principal.Permissions[spec.Permission] = struct{}{}
	service.protectedHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"data":{"hqCustomerSensitivePhone":"` + phone + `"}}`))
	})
	sink := newEventSink(t.Context())
	state := &runState{service: service, principal: principal, token: "signed", origin: "https://admin.example",
		phase: PhasePreview, selectedIDs: map[string]struct{}{spec.ID: {}}, sink: sink}
	ctx := runtimeTestContext{base: context.WithValue(t.Context(), runStateKey{}, state)}
	if err := state.beforeTool(ctx, spec.ID); err != nil {
		t.Fatal(err)
	}
	result, err := service.FixedToolRuntime().Call(ctx, spec, json.RawMessage(`{"id":"member-1"}`))
	if err != nil || !bytes.Contains(result.Body, []byte(phone)) {
		t.Fatalf("authorized phone query failed: %v", err)
	}
	state.afterTool(spec.ID, "", nil)
	for range 2 {
		event := <-sink.events
		encoded, err := json.Marshal(event)
		if err != nil || bytes.Contains(encoded, []byte(phone)) {
			t.Fatalf("phone leaked in tool event: %s, %v", encoded, err)
		}
	}
	if code := safeRunErrorCode(errors.New("upstream response contained " + phone)); code == "" || bytes.Contains([]byte(code), []byte(phone)) {
		t.Fatalf("phone leaked in diagnostic code: %s", code)
	}
}

func TestAIPreviewIncludesBusinessToolTitle(t *testing.T) {
	service, approval, principal, _ := fixedRuntimeFixture(t)
	plan := issueFixturePlan(t, approval, principal)
	state := &runState{service: service, previewPlan: plan}
	view := state.previewView("Create store")
	encoded, err := json.Marshal(view)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(encoded, []byte(`"title":"Create store"`)) {
		t.Fatalf("preview title absent: %s", encoded)
	}
}

func TestAIToolProgressIncludesBusinessToolTitle(t *testing.T) {
	service, _, principal, _ := fixedRuntimeFixture(t)
	sink := newEventSink(t.Context())
	state := &runState{service: service, principal: principal, token: "signed", phase: PhaseRun,
		approvedIDs: map[string]struct{}{"create_store": {}}, selectedIDs: map[string]struct{}{"create_store": {}}, sink: sink}
	if err := state.beforeTool(t.Context(), "create_store"); err != nil {
		t.Fatal(err)
	}
	event := <-sink.events
	data, ok := event.Data.(map[string]any)
	if event.Name != "tool_started" || !ok || data["title"] != "Create store" {
		t.Fatalf("tool progress title = %+v", event)
	}
}

func TestAIToolCompletionIncludesBusinessToolTitle(t *testing.T) {
	service, _, principal, _ := fixedRuntimeFixture(t)
	sink := newEventSink(t.Context())
	state := &runState{service: service, principal: principal, sink: sink}
	state.afterTool("create_store", "store-1", nil)
	event := <-sink.events
	data, ok := event.Data.(map[string]any)
	if event.Name != "tool_finished" || !ok || data["title"] != "Create store" {
		t.Fatalf("tool completion title = %+v", event)
	}
}

func TestAIBusinessToolRequiresCurrentSelection(t *testing.T) {
	service, _, principal, _ := fixedRuntimeFixture(t)
	state := &runState{service: service, principal: principal, token: "signed", phase: PhaseRun,
		approvedIDs: map[string]struct{}{"create_store": {}}, sink: newEventSink(t.Context())}
	if err := state.beforeTool(t.Context(), "create_store"); err == nil {
		t.Fatal("business tool ran without selection")
	}
}

func fixedRuntimeFixture(t *testing.T) (*Service, *ApprovalService, *auth.WorkspacePrincipal, *int) {
	t.Helper()
	approval, principal := approvalFixture(t)
	service, err := NewService(ServiceConfig{Model: serviceTestModel{}, Catalog: approval.catalog, Prompt: "test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	service.approval = approval
	calls := new(int)
	if err := service.SetProtectedHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*calls++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"id":"done"}}`))
	})); err != nil {
		t.Fatal(err)
	}
	return service, approval, principal, calls
}
