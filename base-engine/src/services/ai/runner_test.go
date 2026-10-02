package ai

import (
	"context"
	"errors"
	"iter"
	"strings"
	"testing"
	"time"

	"base-engine/agentkit"
	"base-engine/auth"
	"base-engine/system_prompt"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/session"
	"google.golang.org/genai"
)

func TestAISystemPromptExplainsOnDemandToolSelection(t *testing.T) {
	if !strings.Contains(system_prompt.AdminAgent, "select_tools") {
		t.Fatal("system prompt does not explain on-demand tool selection")
	}
	if !strings.Contains(system_prompt.AdminAgent, "无需调用工具") {
		t.Fatal("system prompt does not allow direct conversational replies")
	}
}

func TestAISystemPromptLimitsMarkdownToUserVisibleText(t *testing.T) {
	if !strings.Contains(system_prompt.AdminAgent, "Markdown") ||
		!strings.Contains(system_prompt.AdminAgent, "工具调用") ||
		!strings.Contains(system_prompt.AdminAgent, "最终回复") {
		t.Fatal("system prompt must distinguish Markdown replies from structured tool calls")
	}
}

func TestAISystemPromptExplainsApprovedRunPayload(t *testing.T) {
	if !strings.Contains(system_prompt.AdminAgent, "approvedSteps") || !strings.Contains(system_prompt.AdminAgent, "originalUserRequest") {
		t.Fatal("execution prompt does not explain the approved plan payload")
	}
}

func TestAICompactionThresholdUsesSingleCallContextWindow(t *testing.T) {
	service := &Service{config: ServiceConfig{ModelContextTokens: 10000, InputTokenBudget: 2000}}
	got := service.compactionConfig()
	if got.TokenThreshold != 5000 {
		t.Fatalf("compaction threshold = %d, want 5000", got.TokenThreshold)
	}
}

func TestAIStreamFinalSnapshotAddsOnlyMissingSuffix(t *testing.T) {
	partial := &session.Event{LLMResponse: model.LLMResponse{Content: genai.NewContentFromText("don", "model"), Partial: true}}
	final := &session.Event{LLMResponse: model.LLMResponse{Content: genai.NewContentFromText("done", "model"), TurnComplete: true}}
	var state agentkit.TextStream
	first := state.Delta(partial)
	if first != "don" {
		t.Fatalf("partial delta = %q", first)
	}
	if delta := state.Delta(final); delta != "e" {
		t.Fatalf("final snapshot suffix = %q, want e", delta)
	}
}

type runnerTextModel struct{}

func (runnerTextModel) Name() string { return "scripted" }
func (runnerTextModel) GenerateContent(_ context.Context, _ *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		if !yield(&model.LLMResponse{Content: genai.NewContentFromText("done", "model"), Partial: true}, nil) {
			return
		}
		yield(&model.LLMResponse{Content: genai.NewContentFromText("done", "model"), TurnComplete: true,
			UsageMetadata: &genai.GenerateContentResponseUsageMetadata{PromptTokenCount: 2, CandidatesTokenCount: 1, TotalTokenCount: 3}}, nil)
	}
}

type cancelAtFinalModel struct{ cancel context.CancelFunc }

func (cancelAtFinalModel) Name() string { return "scripted-cancel" }
func (m cancelAtFinalModel) GenerateContent(_ context.Context, _ *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		response := &model.LLMResponse{Content: genai.NewContentFromText("done", "model"), TurnComplete: true,
			UsageMetadata: &genai.GenerateContentResponseUsageMetadata{PromptTokenCount: 2, CandidatesTokenCount: 1, TotalTokenCount: 3}}
		m.cancel()
		yield(response, nil)
	}
}

func TestAIRunnerNeverSucceedsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	catalog, _ := NewCatalog(nil, nil)
	principal := &auth.WorkspacePrincipal{AccountID: "a", SessionID: "s", WorkspaceType: auth.WorkspaceTypeHeadquarters}
	service, err := NewService(ServiceConfig{Model: cancelAtFinalModel{cancel: cancel}, Catalog: catalog, Prompt: "Test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	if err := service.runPreview(ctx, principal, "signed", "https://admin.example", previewRequest{Prompt: "List"}, newEventSink(ctx)); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled run = %v", err)
	}
}

func TestAIRunnerStreamsReadOnlyPreview(t *testing.T) {
	catalog, _ := NewCatalog(nil, nil)
	principal := &auth.WorkspacePrincipal{AccountID: "account-1", SessionID: "session-1", WorkspaceType: auth.WorkspaceTypeHeadquarters}
	service, err := NewService(ServiceConfig{Model: runnerTextModel{}, Catalog: catalog, Prompt: "Test system prompt", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	sink := newEventSink(t.Context())
	if err := service.runPreview(t.Context(), principal, "signed", "https://admin.example", previewRequest{Prompt: "List stores"}, sink); err != nil {
		t.Fatal(err)
	}
	var names []string
	for len(sink.events) > 0 {
		names = append(names, (<-sink.events).Name)
	}
	if len(names) != 1 || names[0] != "text_delta" {
		t.Fatalf("events = %v", names)
	}
}

func TestAIRunnerCountsModelUsageAndVersion(t *testing.T) {
	catalog, _ := NewCatalog(nil, nil)
	principal := &auth.WorkspacePrincipal{AccountID: "a", SessionID: "s", WorkspaceType: auth.WorkspaceTypeHeadquarters}
	version := uint64(1)
	service, err := NewService(ServiceConfig{ModelProvider: func(context.Context) (model.LLM, uint64, error) { return runnerTextModel{}, version, nil }, Catalog: catalog, Prompt: "Test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	state := &runState{service: service, principal: principal, token: "signed"}
	wrapped := &limitedModel{state: state}
	for _, err := range wrapped.GenerateContent(t.Context(), &model.LLMRequest{}, true) {
		if err != nil {
			t.Fatal(err)
		}
	}
	if state.modelCalls != 1 || state.inputTokens != 2 || state.outputTokens != 1 {
		t.Fatalf("usage = %+v", state)
	}
	version = 2
	failed := false
	for _, err := range wrapped.GenerateContent(t.Context(), &model.LLMRequest{}, true) {
		failed = failed || err != nil
	}
	if !failed {
		t.Fatal("model version change accepted")
	}
}

func TestAIRunnerRejectsMissingModelVersion(t *testing.T) {
	catalog, _ := NewCatalog(nil, nil)
	principal := &auth.WorkspacePrincipal{AccountID: "a", SessionID: "s", WorkspaceType: auth.WorkspaceTypeHeadquarters}
	service, err := NewService(ServiceConfig{ModelProvider: func(context.Context) (model.LLM, uint64, error) { return runnerTextModel{}, 0, nil }, Catalog: catalog, Prompt: "Test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	state := &runState{service: service, principal: principal}
	for _, callErr := range (&limitedModel{state: state}).GenerateContent(t.Context(), &model.LLMRequest{}, true) {
		if callErr == nil || !errors.Is(callErr, errModelUnavailable) {
			t.Fatalf("missing version accepted: %v", callErr)
		}
		return
	}
	t.Fatal("missing version returned no error")
}

func TestAIRunnerRejectsSessionRevokedOnLastResponse(t *testing.T) {
	catalog, _ := NewCatalog(nil, nil)
	principal := &auth.WorkspacePrincipal{AccountID: "a", SessionID: "s", WorkspaceType: auth.WorkspaceTypeHeadquarters}
	checks := 0
	service, err := NewService(ServiceConfig{Model: runnerTextModel{}, Catalog: catalog, Prompt: "Test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) {
		checks++
		if checks >= 3 {
			return nil, errors.New("revoked")
		}
		return principal, nil
	}, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	sink := newEventSink(t.Context())
	err = service.runPreview(t.Context(), principal, "signed", "https://admin.example", previewRequest{Prompt: "List"}, sink)
	if err == nil {
		t.Fatal("revoked session accepted")
	}
}

func TestAIRunnerEnforcesModelAndTokenBudgets(t *testing.T) {
	catalog, _ := NewCatalog(nil, nil)
	principal := &auth.WorkspacePrincipal{AccountID: "a", SessionID: "s", WorkspaceType: auth.WorkspaceTypeHeadquarters}
	service, err := NewService(ServiceConfig{Model: runnerTextModel{}, Catalog: catalog, Prompt: "Test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	state := &runState{service: service, principal: principal}
	wrapped := &limitedModel{state: state}
	for range 80 {
		for _, callErr := range wrapped.GenerateContent(t.Context(), &model.LLMRequest{}, true) {
			if callErr != nil {
				t.Fatal(callErr)
			}
		}
	}
	if !modelCallFails(t, wrapped) {
		t.Fatal("81st model call accepted")
	}
	state.modelCalls, state.inputTokens = 0, service.config.InputTokenBudget-1
	if !modelCallFails(t, wrapped) {
		t.Fatal("input token budget exceeded without error")
	}
}

func modelCallFails(t *testing.T, wrapped *limitedModel) bool {
	t.Helper()
	for _, callErr := range wrapped.GenerateContent(t.Context(), &model.LLMRequest{}, true) {
		if callErr != nil {
			return true
		}
	}
	return false
}

func TestAIRunnerEnforcesToolCallBudget(t *testing.T) {
	catalog, _ := NewCatalog(nil, nil)
	principal := &auth.WorkspacePrincipal{AccountID: "a", SessionID: "s", WorkspaceType: auth.WorkspaceTypeHeadquarters}
	service, err := NewService(ServiceConfig{Model: runnerTextModel{}, Catalog: catalog, Prompt: "Test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	state := &runState{service: service, principal: principal, phase: PhasePreview, sink: newEventSink(t.Context())}
	for range 100 {
		if err := state.beforeTool(t.Context(), "propose_plan"); err != nil {
			t.Fatal(err)
		}
		<-state.sink.events
	}
	if err := state.beforeTool(t.Context(), "propose_plan"); err == nil {
		t.Fatal("101st tool call accepted")
	}
}
