package chatmodel

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

func TestChatCompletionsStreamingToolCalls(t *testing.T) {
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkStreamingRequest(t, r)
		w.Header().Set("Content-Type", "text/event-stream")
		for _, chunk := range streamingToolChunks() {
			_, _ = w.Write([]byte("data: " + chunk + "\n\n"))
		}
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	m, err := newModelWithClient(t.Context(), ModelConfig{Name: "test-model", BaseURL: "https://model.example/v1", APIKey: "test-key"}, client)
	if err != nil {
		t.Fatal(err)
	}
	request := &model.LLMRequest{
		Contents: []*genai.Content{genai.NewContentFromText("lookup", genai.RoleUser)},
		Config: &genai.GenerateContentConfig{Tools: []*genai.Tool{{FunctionDeclarations: []*genai.FunctionDeclaration{{
			Name: "lookup", ParametersJsonSchema: map[string]any{"type": "object"},
		}}}}},
	}
	var partial strings.Builder
	var final *model.LLMResponse
	for response, runErr := range m.GenerateContent(t.Context(), request, true) {
		if runErr != nil {
			t.Fatal(runErr)
		}
		if response.Partial {
			partial.WriteString(response.Content.Parts[0].Text)
		} else {
			final = response
		}
	}
	checkStreamingToolResponse(t, partial.String(), final)
}

func checkStreamingToolResponse(t *testing.T, partial string, final *model.LLMResponse) {
	t.Helper()
	if partial != "hello" || final == nil || !final.TurnComplete || final.UsageMetadata == nil || final.UsageMetadata.TotalTokenCount != 12 {
		t.Fatalf("partial=%q final=%+v", partial, final)
	}
	if len(final.Content.Parts) != 2 || final.Content.Parts[1].FunctionCall == nil || final.Content.Parts[1].FunctionCall.ID != "call-1" || final.Content.Parts[1].FunctionCall.Args["symbol"] != "AAPL" {
		t.Fatalf("tool call = %+v", final.Content)
	}
}

func checkStreamingRequest(t *testing.T, r *http.Request) {
	t.Helper()
	if r.URL.Path != "/v1/chat/completions" {
		t.Errorf("path = %s", r.URL.Path)
	}
	var payload struct {
		Stream        bool `json:"stream"`
		StreamOptions struct {
			IncludeUsage bool `json:"include_usage"`
		} `json:"stream_options"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		t.Errorf("decode request: %v", err)
	}
	if !payload.Stream || !payload.StreamOptions.IncludeUsage {
		t.Errorf("stream and usage must be enabled")
	}
}

func streamingToolChunks() []string {
	return []string{
		`{"choices":[{"index":0,"delta":{"content":"hel"}}]}`,
		`{"choices":[{"index":0,"delta":{"content":"lo","tool_calls":[{"index":0,"id":"call-1","type":"function","function":{"name":"lookup","arguments":"{\"symbol\":"}}]}}]}`,
		`{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"\"AAPL\"}"}}]},"finish_reason":"tool_calls"}]}`,
		`{"choices":[],"usage":{"prompt_tokens":8,"completion_tokens":4,"total_tokens":12}}`,
	}
}

func TestChatCompletionsStreamingUsageInTerminalChoice(t *testing.T) {
	request := &model.LLMRequest{
		Contents: []*genai.Content{genai.NewContentFromText("lookup", genai.RoleUser)},
		Config: &genai.GenerateContentConfig{Tools: []*genai.Tool{{FunctionDeclarations: []*genai.FunctionDeclaration{{
			Name: "lookup", ParametersJsonSchema: map[string]any{"type": "object"},
		}}}}},
	}
	var payload string
	for _, chunk := range streamingToolChunks()[:3] {
		payload += "data: " + chunk + "\n\n"
	}
	payload += `data: {"choices":[{"index":0,"delta":{},"finish_reason":null}],"usage":{"prompt_tokens":2,"completion_tokens":1,"total_tokens":3}}` + "\n\n"
	payload += "data: [DONE]\n\n"
	if err := runStreamError(t, payload, request); err != nil {
		t.Fatalf("terminal choice usage should be accepted: %v", err)
	}
}
