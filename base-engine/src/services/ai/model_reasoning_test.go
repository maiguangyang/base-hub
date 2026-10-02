package ai

import (
	"fmt"
	"net/http"
	"testing"

	"base-engine/agentkit"
	"google.golang.org/adk/v2/session"
	"google.golang.org/genai"
)

func TestConfigProbeReplaysThinkingBeforeNextToolCall(t *testing.T) {
	requests := 0
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		request := checkScriptedRequest(t, r, requests)
		if requests > 1 {
			for step := 1; step < requests; step++ {
				if !hasReasoningForCall(request.Messages, fmt.Sprintf("call-%d", step), fmt.Sprintf("thought-%d", step)) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
		}
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprintf(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"reasoning_content\":\"thought-%d\"}}]}\n\n", requests)
		if requests <= probeStepCount {
			fmt.Fprintf(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call-%d\",\"type\":\"function\",\"function\":{\"name\":\"next_step\",\"arguments\":\"{\\\"value\\\":%d}\"}}]},\"finish_reason\":\"tool_calls\"}]}\n\n", requests, requests)
		} else {
			_, _ = w.Write([]byte("data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"done\"},\"finish_reason\":\"stop\"}]}\n\n"))
		}
		_, _ = w.Write([]byte("data: {\"choices\":[],\"usage\":{\"prompt_tokens\":10,\"completion_tokens\":3,\"total_tokens\":13}}\n\ndata: [DONE]\n\n"))
	}))
	err := probeModelConnectionWithClient(t.Context(), ModelConfig{Name: "thinking-model", BaseURL: "https://model.example/v1", APIKey: "test-key"}, client)
	if err != nil || requests != probeStepCount+1 {
		t.Fatalf("thinking probe error=%v requests=%d", err, requests)
	}
}

func hasReasoningForCall(messages []testChatMessage, callID, reasoning string) bool {
	for _, message := range messages {
		if message.Role != "assistant" || message.ReasoningContent != reasoning || message.Content == reasoning {
			continue
		}
		for _, call := range message.ToolCalls {
			if call.ID == callID {
				return true
			}
		}
	}
	return false
}

func TestReasoningIsExcludedFromUserVisibleText(t *testing.T) {
	state := &agentkit.TextStream{}
	event := &session.Event{
		Content: &genai.Content{Role: genai.RoleModel, Parts: []*genai.Part{
			{Text: "private reasoning", Thought: true},
			{Text: "visible answer"},
		}},
	}
	if got := state.Delta(event); got != "visible answer" {
		t.Fatalf("visible text = %q", got)
	}
}
