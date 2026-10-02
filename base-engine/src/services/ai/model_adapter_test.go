package ai

import (
	"net/http"
	"net/http/httptest"
)

type modelRoundTrip func(*http.Request) (*http.Response, error)

func (f modelRoundTrip) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }

func modelTestClient(handler http.Handler) *http.Client {
	return &http.Client{Transport: modelRoundTrip(func(request *http.Request) (*http.Response, error) {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		return recorder.Result(), nil
	})}
}

type testChatRequest struct {
	Stream   bool              `json:"stream"`
	Tools    []testChatTool    `json:"tools"`
	Messages []testChatMessage `json:"messages"`
}

type testChatTool struct {
	Function struct {
		Name string `json:"name"`
	} `json:"function"`
}

type testChatMessage struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content"`
	ToolCallID       string `json:"tool_call_id"`
	ToolCalls        []struct {
		ID string `json:"id"`
	} `json:"tool_calls"`
}
