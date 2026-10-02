package ai

import (
	"context"
	"iter"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

type conversationCaptureModel struct{ seen string }

func (*conversationCaptureModel) Name() string { return "conversation-capture" }

func (capture *conversationCaptureModel) GenerateContent(_ context.Context, request *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		capture.seen = request.Contents[len(request.Contents)-1].Parts[0].Text
		yield(&model.LLMResponse{Content: genai.NewContentFromText("好的", genai.RoleModel), TurnComplete: true,
			UsageMetadata: &genai.GenerateContentResponseUsageMetadata{PromptTokenCount: 20, CandidatesTokenCount: 2, TotalTokenCount: 22}}, nil)
	}
}

func TestAIConversationKeepsPriorTurnsForFollowUp(t *testing.T) {
	catalog, _ := NewCatalog(nil, nil)
	principal := &auth.WorkspacePrincipal{AccountID: "account-1", SessionID: "session-1", WorkspaceType: auth.WorkspaceTypeHeadquarters}
	capture := &conversationCaptureModel{}
	service, err := NewService(ServiceConfig{Model: capture, Catalog: catalog, Prompt: "Test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	input := previewRequest{Prompt: "它们分别是什么？", History: []chatTurn{{Role: "user", Text: "列出门店"}, {Role: "assistant", Text: "有 A 店与 B 店。"}}}
	if err := service.runPreview(t.Context(), principal, "signed", "https://admin.example", input, newEventSink(t.Context())); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(capture.seen, "列出门店") || !strings.Contains(capture.seen, "有 A 店与 B 店") || !strings.Contains(capture.seen, "它们分别是什么？") {
		t.Fatalf("conversation context missing: %q", capture.seen)
	}
}
