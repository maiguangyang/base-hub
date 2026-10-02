package ai

import (
	"context"
	"encoding/json"
	"iter"
	"testing"
	"time"

	"base-engine/auth"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

type capturingRunModel struct{ prompts []string }

type approvedPromptPayload struct {
	OriginalUserRequest string         `json:"originalUserRequest"`
	ApprovedSteps       []ApprovedStep `json:"approvedSteps"`
}

func findApprovedPrompt(prompts []string) approvedPromptPayload {
	for _, text := range prompts {
		var prompt approvedPromptPayload
		if json.Unmarshal([]byte(text), &prompt) == nil && len(prompt.ApprovedSteps) > 0 {
			return prompt
		}
	}
	return approvedPromptPayload{}
}

func (*capturingRunModel) Name() string { return "capturing-run" }

func (m *capturingRunModel) GenerateContent(_ context.Context, request *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		for _, content := range request.Contents {
			if content.Role == genai.RoleUser {
				for _, part := range content.Parts {
					m.prompts = append(m.prompts, part.Text)
				}
			}
		}
		yield(&model.LLMResponse{Content: genai.NewContentFromText("done", "model"), TurnComplete: true,
			UsageMetadata: &genai.GenerateContentResponseUsageMetadata{PromptTokenCount: 2, CandidatesTokenCount: 1, TotalTokenCount: 3}}, nil)
	}
}

func TestApprovedRunReceivesExactApprovedSteps(t *testing.T) {
	approval, principal := approvalFixture(t)
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "create_store", OperationID: "graphql.mutation.createStore", Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), MaxCalls: 1, Sequence: 1}}}
	draft.BindPrompt("Create the store")
	plan, _, err := approval.ValidateDraft(draft, principal, nil)
	if err != nil {
		t.Fatal(err)
	}
	issued, err := approval.Issue(plan)
	if err != nil {
		t.Fatal(err)
	}
	plan, err = approval.Consume(issued, principal)
	if err != nil {
		t.Fatal(err)
	}
	m := &capturingRunModel{}
	service := &Service{config: ServiceConfig{Model: m, Catalog: approval.catalog, Prompt: "Test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000}, approval: approval, now: time.Now}
	runErr := service.runApproved(t.Context(), principal, "signed", "https://admin.example", plan, newEventSink(t.Context()))
	if len(m.prompts) == 0 {
		t.Fatalf("model was not called: %v", runErr)
	}
	prompt := findApprovedPrompt(m.prompts)
	if prompt.OriginalUserRequest != "Create the store" || len(prompt.ApprovedSteps) != 1 ||
		prompt.ApprovedSteps[0].ToolID != "create_store" || string(prompt.ApprovedSteps[0].Arguments) != string(plan.Steps[0].Arguments) {
		t.Fatalf("approved plan missing from run prompts: %#v", m.prompts)
	}
}

func TestApprovedRunPromptUsesStableStepFields(t *testing.T) {
	prompt, err := approvedRunPrompt("Create", Plan{Steps: []ApprovedStep{{ToolID: "create_store", Arguments: json.RawMessage(`{"id":"a"}`), MaxCalls: 1, Sequence: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	var payload struct {
		ApprovedSteps []map[string]json.RawMessage `json:"approvedSteps"`
	}
	if err := json.Unmarshal([]byte(prompt), &payload); err != nil || len(payload.ApprovedSteps) != 1 {
		t.Fatalf("approved run prompt = %s: %v", prompt, err)
	}
	step := payload.ApprovedSteps[0]
	if string(step["toolId"]) != `"create_store"` || string(step["arguments"]) != `{"id":"a"}` {
		t.Fatalf("approved step fields = %s", prompt)
	}
}
