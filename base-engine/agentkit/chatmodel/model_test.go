package chatmodel

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

func TestModelConfig(t *testing.T) {
	valid := ModelConfig{
		Name:    "example-model",
		BaseURL: "https://example.invalid/v1",
		APIKey:  "test-key",
	}
	cases := []struct {
		name   string
		change func(*ModelConfig)
	}{
		{"missing name", func(c *ModelConfig) { c.Name = "" }},
		{"blank name", func(c *ModelConfig) { c.Name = " \t" }},
		{"missing base URL", func(c *ModelConfig) { c.BaseURL = "" }},
		{"missing API key", func(c *ModelConfig) { c.APIKey = "" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := valid
			tc.change(&cfg)
			got, err := NewModel(context.Background(), cfg)
			if got != nil || !errors.Is(err, ErrModelConfig) {
				t.Fatalf("NewModel() = (%v, %v), want nil and ErrModelConfig", got, err)
			}
		})
	}
}

func TestChatCompletionsNonStreaming(t *testing.T) {
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkNonStreamingRequest(t, r)
		_, _ = w.Write([]byte(`{"choices":[{"index":0,"message":{"role":"assistant","content":"done"},"finish_reason":"stop"}],"usage":{"prompt_tokens":3,"completion_tokens":1,"total_tokens":4}}`))
	}))
	m, err := newModelWithClient(t.Context(), ModelConfig{Name: "test-model", BaseURL: "https://model.example/v1", APIKey: "test-key"}, client)
	if err != nil {
		t.Fatal(err)
	}
	var got *model.LLMResponse
	for response, runErr := range m.GenerateContent(t.Context(), &model.LLMRequest{
		Contents: []*genai.Content{genai.NewContentFromText("hello", genai.RoleUser)},
	}, false) {
		if runErr != nil {
			t.Fatal(runErr)
		}
		got = response
	}
	checkNonStreamingResponse(t, got)
}

func checkNonStreamingResponse(t *testing.T, got *model.LLMResponse) {
	t.Helper()
	if got == nil || got.Content == nil || len(got.Content.Parts) != 1 || got.Content.Parts[0].Text != "done" || !got.TurnComplete || got.UsageMetadata == nil || got.UsageMetadata.TotalTokenCount != 4 {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func checkNonStreamingRequest(t *testing.T, r *http.Request) {
	t.Helper()
	if r.URL.Path != "/v1/chat/completions" || r.Header.Get("Authorization") != "Bearer test-key" {
		t.Errorf("unexpected path or authorization")
	}
	var payload struct {
		Model    string `json:"model"`
		Stream   bool   `json:"stream"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		t.Errorf("decode request: %v", err)
	}
	if payload.Model != "test-model" || payload.Stream || len(payload.Messages) != 1 || payload.Messages[0].Content != "hello" {
		t.Errorf("unexpected request: %+v", payload)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func modelTestClient(handler http.Handler) *http.Client {
	return &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, r)
		return recorder.Result(), nil
	})}
}
