package ai

import (
	"context"
	"iter"
	"testing"
	"time"

	"base-engine/auth"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

type observedRequestModel struct{ names []string }

func (m *observedRequestModel) Name() string { return "observed" }
func (m *observedRequestModel) GenerateContent(_ context.Context, request *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		for _, item := range request.Config.Tools {
			for _, declaration := range item.FunctionDeclarations {
				m.names = append(m.names, declaration.Name)
			}
		}
		yield(&model.LLMResponse{Content: genai.NewContentFromText("ok", "model"), TurnComplete: true,
			UsageMetadata: &genai.GenerateContentResponseUsageMetadata{PromptTokenCount: 2, CandidatesTokenCount: 1, TotalTokenCount: 3}}, nil)
	}
}

func TestAIModelInitiallyAdvertisesOnlyToolSelection(t *testing.T) {
	approval, principal := approvalFixture(t)
	underlying := &observedRequestModel{}
	service, err := NewService(ServiceConfig{Model: underlying, Catalog: approval.catalog, Prompt: "test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	state := &runState{service: service, principal: principal, token: "signed", phase: PhasePreview}
	request := &model.LLMRequest{Config: &genai.GenerateContentConfig{Tools: []*genai.Tool{{FunctionDeclarations: []*genai.FunctionDeclaration{
		{Name: "create_store"}, {Name: "select_tools"}, {Name: "propose_plan"},
	}}}}}
	for _, callErr := range (&limitedModel{state: state}).GenerateContent(t.Context(), request, true) {
		if callErr != nil {
			t.Fatal(callErr)
		}
	}
	if len(underlying.names) != 2 || underlying.names[0] != "select_tools" || underlying.names[1] != "propose_plan" {
		t.Fatalf("initial advertised tools = %v", underlying.names)
	}
}

func TestAISelectionToolIsAvailableInBothPhases(t *testing.T) {
	approval, principal := approvalFixture(t)
	service, err := NewService(ServiceConfig{Model: runnerTextModel{}, Catalog: approval.catalog, Prompt: "test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	for _, phase := range []Phase{PhasePreview, PhaseRun} {
		state := &runState{service: service, principal: principal, phase: phase}
		items, err := service.toolsForRun(state)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, item := range items {
			found = found || item.Name() == "select_tools"
		}
		if !found {
			t.Fatalf("selection tool missing in %s", phase)
		}
	}
}

func TestAISelectionAdvertisesOnlyRequestedAuthorizedTools(t *testing.T) {
	approval, principal := approvalFixture(t)
	underlying := &observedRequestModel{}
	service, err := NewService(ServiceConfig{Model: underlying, Catalog: approval.catalog, Prompt: "test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	state := &runState{service: service, principal: principal, token: "signed", phase: PhaseRun,
		approvedIDs: map[string]struct{}{"create_store": {}, "suspend": {}}}
	selector, err := state.selectionTool()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := selector.(interface {
		Run(agent.Context, any) (map[string]any, error)
	}).Run(runtimeTestContext{base: t.Context()}, map[string]any{"names": []string{"create_store"}}); err != nil {
		t.Fatal(err)
	}
	request := &model.LLMRequest{Config: &genai.GenerateContentConfig{Tools: []*genai.Tool{{FunctionDeclarations: []*genai.FunctionDeclaration{
		{Name: "create_store"}, {Name: "suspend"}, {Name: "select_tools"},
	}}}}}
	for _, callErr := range (&limitedModel{state: state}).GenerateContent(t.Context(), request, true) {
		if callErr != nil {
			t.Fatal(callErr)
		}
	}
	if len(underlying.names) != 2 || underlying.names[0] != "create_store" || underlying.names[1] != "select_tools" {
		t.Fatalf("advertised tools after selection = %v", underlying.names)
	}
}

func TestAISelectionKeepsRecentlySelectedToolsAvailable(t *testing.T) {
	approval, principal := approvalFixture(t)
	underlying := &observedRequestModel{}
	service, err := NewService(ServiceConfig{Model: underlying, Catalog: approval.catalog, Prompt: "test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	state := &runState{service: service, principal: principal, token: "signed", phase: PhaseRun,
		approvedIDs: map[string]struct{}{"create_store": {}, "suspend": {}}}
	ctx := runtimeTestContext{base: t.Context()}
	for _, name := range []string{"create_store", "suspend"} {
		if _, err := state.chooseTools(ctx, toolSelectionInput{Names: []string{name}}); err != nil {
			t.Fatal(err)
		}
	}
	request := &model.LLMRequest{Config: &genai.GenerateContentConfig{Tools: []*genai.Tool{{FunctionDeclarations: []*genai.FunctionDeclaration{
		{Name: "create_store"}, {Name: "suspend"}, {Name: "select_tools"},
	}}}}}
	for _, callErr := range (&limitedModel{state: state}).GenerateContent(t.Context(), request, true) {
		if callErr != nil {
			t.Fatal(callErr)
		}
	}
	if len(underlying.names) != 3 || underlying.names[0] != "create_store" || underlying.names[1] != "suspend" || underlying.names[2] != "select_tools" {
		t.Fatalf("advertised tools after two selections = %v", underlying.names)
	}
}

func TestAISelectionEvictsOldestAfterFiveBusinessTools(t *testing.T) {
	_, principal := approvalFixture(t)
	catalog := &Catalog{}
	for _, name := range []string{"one", "two", "three", "four", "five", "six", "seven"} {
		catalog.specs = append(catalog.specs, ToolSpec{ID: name, Name: name, Permission: "store:create", Mode: ModeReadOnly,
			Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeHeadquarters}})
	}
	service := &Service{config: ServiceConfig{Catalog: catalog, ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }}, now: time.Now}
	state := &runState{service: service, principal: principal, phase: PhasePreview}
	ctx := runtimeTestContext{base: t.Context()}
	for _, name := range []string{"one", "two", "three", "four", "five", "six"} {
		if _, err := state.chooseTools(ctx, toolSelectionInput{Names: []string{name}}); err != nil {
			t.Fatal(err)
		}
	}
	if len(state.selectedIDs) != 5 {
		t.Fatalf("selected tool count = %d, want 5", len(state.selectedIDs))
	}
	if _, present := state.selectedIDs["one"]; present {
		t.Fatal("oldest tool remains advertised")
	}
	for _, name := range []string{"two", "seven"} {
		if _, err := state.chooseTools(ctx, toolSelectionInput{Names: []string{name}}); err != nil {
			t.Fatal(err)
		}
	}
	if _, present := state.selectedIDs["three"]; present {
		t.Fatal("least recently selected tool remains advertised")
	}
	if _, present := state.selectedIDs["two"]; !present {
		t.Fatal("reselected tool was evicted")
	}
}

func TestAISelectionRequiresOneBusinessToolAtATime(t *testing.T) {
	approval, principal := approvalFixture(t)
	service, err := NewService(ServiceConfig{Model: runnerTextModel{}, Catalog: approval.catalog, Prompt: "test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	state := &runState{service: service, principal: principal, token: "signed", phase: PhaseRun,
		approvedIDs: map[string]struct{}{"create_store": {}, "suspend": {}}}
	if _, err := state.chooseTools(runtimeTestContext{base: t.Context()}, toolSelectionInput{Names: []string{"create_store", "suspend"}}); err == nil {
		t.Fatal("multiple business tool declarations accepted")
	}
}
