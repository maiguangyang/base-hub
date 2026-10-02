package chatmodel

import (
	"strings"
	"testing"

	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

func TestNonStreamingThinkingToolResponseCanBeReplayed(t *testing.T) {
	payload := `{"choices":[{"index":0,"message":{"role":"assistant","content":"","reasoning_content":"private reasoning","tool_calls":[{"id":"call-1","type":"function","function":{"name":"lookup","arguments":"{}"}}]},"finish_reason":"tool_calls"}],"usage":{"prompt_tokens":8,"completion_tokens":4,"total_tokens":12}}`
	var got *model.LLMResponse
	err := readChatCompletion(strings.NewReader(payload), map[string]struct{}{"lookup": {}}, func(response *model.LLMResponse, _ error) bool {
		got = response
		return true
	})
	if err != nil || got == nil || len(got.Content.Parts) != 2 || !got.Content.Parts[0].Thought {
		t.Fatalf("decoded reasoning response = %+v, error=%v", got, err)
	}
	wire, _, err := buildChatRequest("thinking-model", &model.LLMRequest{Contents: []*genai.Content{
		genai.NewContentFromText("lookup", genai.RoleUser), got.Content,
		{Role: genai.RoleUser, Parts: []*genai.Part{{FunctionResponse: &genai.FunctionResponse{
			ID: "call-1", Name: "lookup", Response: map[string]any{"ok": true},
		}}}},
	}}, false)
	if err != nil || len(wire.Messages) != 3 || wire.Messages[1].ReasoningContent != "private reasoning" {
		t.Fatalf("replayed reasoning messages = %+v, error=%v", wire.Messages, err)
	}
}

func TestStreamingUsageOnFinishingToolChoice(t *testing.T) {
	payload := "data: {\"choices\":[{\"index\":0,\"delta\":{\"reasoning_content\":\"thinking\"}}]}\n\n"
	payload += "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call-1\",\"type\":\"function\",\"function\":{\"name\":\"lookup\",\"arguments\":\"{}\"}}]},\"finish_reason\":\"tool_calls\"}],\"usage\":{\"prompt_tokens\":8,\"completion_tokens\":4,\"total_tokens\":12}}\n\n"
	payload += "data: [DONE]\n\n"
	request := &model.LLMRequest{
		Contents: []*genai.Content{genai.NewContentFromText("lookup", genai.RoleUser)},
		Config: &genai.GenerateContentConfig{Tools: []*genai.Tool{{FunctionDeclarations: []*genai.FunctionDeclaration{{
			Name: "lookup", ParametersJsonSchema: map[string]any{"type": "object"},
		}}}}},
	}
	if err := runStreamError(t, payload, request); err != nil {
		t.Fatalf("finishing choice with usage rejected: %v", err)
	}
}

func TestTerminalUsageChoicePreservesReasoning(t *testing.T) {
	payload := "data: {\"choices\":[{\"index\":0,\"delta\":{\"reasoning_content\":\"before-\"}}]}\n\n"
	payload += "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call-1\",\"type\":\"function\",\"function\":{\"name\":\"lookup\",\"arguments\":\"{}\"}}]},\"finish_reason\":\"tool_calls\"}]}\n\n"
	payload += "data: {\"choices\":[{\"index\":0,\"delta\":{\"reasoning_content\":\"after\"},\"finish_reason\":null}],\"usage\":{\"prompt_tokens\":8,\"completion_tokens\":4,\"total_tokens\":12}}\n\n"
	payload += "data: [DONE]\n\n"
	var final *model.LLMResponse
	err := readChatStream(t.Context(), strings.NewReader(payload), map[string]struct{}{"lookup": {}}, func(response *model.LLMResponse, _ error) bool {
		final = response
		return true
	})
	if err != nil || final == nil || len(final.Content.Parts) != 2 ||
		!final.Content.Parts[0].Thought || final.Content.Parts[0].Text != "before-after" {
		t.Fatalf("terminal reasoning response = %+v, error=%v", final, err)
	}
}
