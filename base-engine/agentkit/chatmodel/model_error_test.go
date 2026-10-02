package chatmodel

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

func TestChatCompletionsRejectsMissingUsage(t *testing.T) {
	request := &model.LLMRequest{Contents: []*genai.Content{genai.NewContentFromText("hello", genai.RoleUser)}}
	for _, payload := range []string{
		`data: {"choices":[{"index":0,"delta":{"content":"ok"},"finish_reason":"stop"}]}\n\ndata: [DONE]\n\n`,
		`data: {"choices":[{"index":0,"delta":{"content":"ok"},"finish_reason":"stop"}]}\n\ndata: {"choices":[],"usage":{"prompt_tokens":0,"completion_tokens":0,"total_tokens":0}}\n\ndata: [DONE]\n\n`,
	} {
		got := runStreamError(t, strings.ReplaceAll(payload, `\n`, "\n"), request)
		if !errors.Is(got, ErrModelProtocol) {
			t.Fatalf("error = %v, want protocol error", got)
		}
	}
}

func TestChatCompletionsRejectsMalformedStreams(t *testing.T) {
	request := &model.LLMRequest{Contents: []*genai.Content{genai.NewContentFromText("hello", genai.RoleUser)}}
	usage := `data: {"choices":[],"usage":{"prompt_tokens":2,"completion_tokens":1,"total_tokens":3}}\n\n`
	cases := map[string]string{
		"malformed JSON":  `data: {invalid}\n\n` + usage + `data: [DONE]\n\n`,
		"missing done":    `data: {"choices":[{"index":0,"delta":{"content":"ok"},"finish_reason":"stop"}]}\n\n` + usage,
		"truncated JSON":  `data: {"choices":[{\n\n` + usage + `data: [DONE]\n\n`,
		"event too large": `data: ` + strings.Repeat("x", maxSSELineBytes) + `\n\n`,
	}
	for name, payload := range cases {
		t.Run(name, func(t *testing.T) {
			got := runStreamError(t, strings.ReplaceAll(payload, `\n`, "\n"), request)
			if !errors.Is(got, ErrModelProtocol) {
				t.Fatalf("error = %v, want protocol error", got)
			}
		})
	}
}

func TestChatCompletionsReportsSafeMalformedFrameShape(t *testing.T) {
	request := &model.LLMRequest{Contents: []*genai.Content{genai.NewContentFromText("hello", genai.RoleUser)}}
	secret := "private-model-output"
	for _, tc := range []struct{ payload, marker string }{
		{"{\"choices\":[{\"content\":\"" + secret + "\",}]}", "JSON syntax at byte"},
		{"{\"choices\":[{\"content\":\"" + secret + "\"}", "JSON incomplete frame length"},
		{"{\"choices\":[]} {\"content\":\"" + secret + "\"}", "JSON multiple values"},
		{"{\"choices\":\"" + secret + "\"}", "JSON field choices type string"},
	} {
		err := runStreamError(t, "data: "+tc.payload+"\n\ndata: [DONE]\n\n", request)
		if !errors.Is(err, ErrModelProtocol) || !strings.Contains(err.Error(), tc.marker) || strings.Contains(err.Error(), secret) {
			t.Fatalf("safe malformed frame diagnosis = %v", err)
		}
	}
}

func TestChatCompletionsRejectsUnknownTool(t *testing.T) {
	request := &model.LLMRequest{Contents: []*genai.Content{genai.NewContentFromText("hello", genai.RoleUser)}}
	payload := `data: {"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call-1","type":"function","function":{"name":"admin_delete","arguments":"{}"}}]},"finish_reason":"tool_calls"}]}\n\n`
	payload += `data: {"choices":[],"usage":{"prompt_tokens":2,"completion_tokens":1,"total_tokens":3}}\n\ndata: [DONE]\n\n`
	got := runStreamError(t, strings.ReplaceAll(payload, `\n`, "\n"), request)
	if !errors.Is(got, ErrModelProtocol) {
		t.Fatalf("error = %v, want protocol error", got)
	}
}

func TestChatCompletionsCancellationHasNoFinalResponse(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	payload := `data: {"choices":[{"index":0,"delta":{"content":"first"}}]}\n\n`
	payload += `data: {"choices":[{"index":0,"delta":{"content":"second"},"finish_reason":"stop"}]}\n\n`
	payload += `data: {"choices":[],"usage":{"prompt_tokens":2,"completion_tokens":2,"total_tokens":4}}\n\ndata: [DONE]\n\n`
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.ReplaceAll(payload, `\n`, "\n")))
	}))
	m, err := newModelWithClient(ctx, ModelConfig{Name: "test-model", BaseURL: "https://model.example/v1", APIKey: "test-key"}, client)
	if err != nil {
		t.Fatal(err)
	}
	request := &model.LLMRequest{Contents: []*genai.Content{genai.NewContentFromText("hello", genai.RoleUser)}}
	partial, final, cancelled := false, false, false
	for response, runErr := range m.GenerateContent(ctx, request, true) {
		if response != nil {
			partial = partial || response.Partial
			final = final || response.TurnComplete
			cancel()
		}
		cancelled = cancelled || errors.Is(runErr, context.Canceled)
	}
	if !partial || final || !cancelled {
		t.Fatalf("partial=%t final=%t cancelled=%t", partial, final, cancelled)
	}
}

func TestChatCompletionsCancelledConsumerStopsIterator(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	payload := "data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"first\"}}]}\n\n"
	payload += "data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"second\"},\"finish_reason\":\"stop\"}]}\n\n"
	payload += "data: {\"choices\":[],\"usage\":{\"prompt_tokens\":2,\"completion_tokens\":2,\"total_tokens\":4}}\n\ndata: [DONE]\n\n"
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(payload)) }))
	m, err := newModelWithClient(ctx, ModelConfig{Name: "test-model", BaseURL: "https://model.example/v1", APIKey: "test-key"}, client)
	if err != nil {
		t.Fatal(err)
	}
	request := &model.LLMRequest{Contents: []*genai.Content{genai.NewContentFromText("hello", genai.RoleUser)}}
	seen := false
	for response := range m.GenerateContent(ctx, request, true) {
		if response != nil && response.Partial {
			seen = true
			cancel()
			break
		}
	}
	if !seen {
		t.Fatal("stream emitted no partial response")
	}
}

func runStreamError(t *testing.T, payload string, request *model.LLMRequest) error {
	t.Helper()
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte(payload))
	}))
	m, err := newModelWithClient(t.Context(), ModelConfig{Name: "test-model", BaseURL: "https://model.example/v1", APIKey: "test-key"}, client)
	if err != nil {
		t.Fatal(err)
	}
	var result error
	for _, runErr := range m.GenerateContent(t.Context(), request, true) {
		if runErr != nil {
			result = runErr
		}
	}
	return result
}

func TestChatCompletionsRejectsCrossOriginRedirect(t *testing.T) {
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", "https://other.example/collect")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}))
	m, err := newModelWithClient(t.Context(), ModelConfig{Name: "test-model", BaseURL: "https://model.example/v1", APIKey: "secret"}, client)
	if err != nil {
		t.Fatal(err)
	}
	request := &model.LLMRequest{Contents: []*genai.Content{genai.NewContentFromText("hello", genai.RoleUser)}}
	var result error
	for _, runErr := range m.GenerateContent(context.Background(), request, false) {
		result = runErr
	}
	if !errors.Is(result, ErrModelUnavailable) {
		t.Fatalf("redirect error = %v", result)
	}
}

func TestChatCompletionsReportsSafeHTTPStatusWithoutResponseBody(t *testing.T) {
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"secret-upstream-body"}`))
	}))
	m, err := newModelWithClient(t.Context(), ModelConfig{Name: "test-model", BaseURL: "https://model.example/v1", APIKey: "secret"}, client)
	if err != nil {
		t.Fatal(err)
	}
	request := &model.LLMRequest{Contents: []*genai.Content{genai.NewContentFromText("hello", genai.RoleUser)}}
	var result error
	for _, runErr := range m.GenerateContent(t.Context(), request, true) {
		result = runErr
	}
	var upstream interface{ UpstreamStatusCode() int }
	if !errors.Is(result, ErrModelUnavailable) || !errors.As(result, &upstream) || upstream.UpstreamStatusCode() != http.StatusBadRequest || strings.Contains(result.Error(), "secret-upstream-body") {
		t.Fatalf("status error = %v", result)
	}
}

func TestChatCompletionsPreservesToolCallIDOnNextTurn(t *testing.T) {
	request := &model.LLMRequest{Contents: []*genai.Content{
		genai.NewContentFromText("lookup", genai.RoleUser),
		{Role: genai.RoleModel, Parts: []*genai.Part{{FunctionCall: &genai.FunctionCall{
			ID: "call-7", Name: "lookup", Args: map[string]any{"symbol": "AAPL"},
		}}}},
		{Role: genai.RoleUser, Parts: []*genai.Part{{FunctionResponse: &genai.FunctionResponse{
			ID: "call-7", Name: "lookup", Response: map[string]any{"price": 12},
		}}}},
	}}
	wire, _, err := buildChatRequest("test-model", request, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(wire.Messages) != 3 || len(wire.Messages[1].ToolCalls) != 1 ||
		wire.Messages[1].ToolCalls[0].ID != "call-7" || wire.Messages[2].ToolCallID != "call-7" {
		t.Fatalf("tool call/response IDs lost: %+v", wire.Messages)
	}
}

func TestChatCompletionsRejectsNonTextSystemInstruction(t *testing.T) {
	request := &model.LLMRequest{
		Contents: []*genai.Content{genai.NewContentFromText("hello", genai.RoleUser)},
		Config: &genai.GenerateContentConfig{SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{InlineData: &genai.Blob{MIMEType: "image/png", Data: []byte{1}}}},
		}},
	}
	_, _, err := buildChatRequest("test-model", request, false)
	if !errors.Is(err, ErrModelProtocol) {
		t.Fatalf("system instruction error = %v", err)
	}
}

func TestChatCompletionsRejectsUnsupportedTool(t *testing.T) {
	request := &model.LLMRequest{
		Contents: []*genai.Content{genai.NewContentFromText("hello", genai.RoleUser)},
		Config:   &genai.GenerateContentConfig{Tools: []*genai.Tool{{CodeExecution: &genai.ToolCodeExecution{}}}},
	}
	_, _, err := buildChatRequest("test-model", request, false)
	if !errors.Is(err, ErrModelProtocol) {
		t.Fatalf("unsupported tool error = %v", err)
	}
}

func TestChatCompletionsRejectsNonAssistantCompletion(t *testing.T) {
	response := `{"choices":[{"index":0,"message":{"role":"tool","content":"fake"},"finish_reason":"stop"}],"usage":{"prompt_tokens":2,"completion_tokens":1,"total_tokens":3}}`
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(response))
	}))
	m, err := newModelWithClient(t.Context(), ModelConfig{Name: "test-model", BaseURL: "https://model.example/v1", APIKey: "test-key"}, client)
	if err != nil {
		t.Fatal(err)
	}
	request := &model.LLMRequest{Contents: []*genai.Content{genai.NewContentFromText("hello", genai.RoleUser)}}
	var result error
	for _, runErr := range m.GenerateContent(t.Context(), request, false) {
		result = runErr
	}
	if !errors.Is(result, ErrModelProtocol) {
		t.Fatalf("completion error = %v", result)
	}
}
