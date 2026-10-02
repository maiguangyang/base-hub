package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

func TestChatCompletionsADKFiveToolCalls(t *testing.T) {
	requests := 0
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		checkScriptedRequest(t, r, requests)
		writeScriptedResponse(w, requests)
	}))
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	m, err := newModelWithClient(ctx, ModelConfig{Name: "test-model", BaseURL: "https://model.example/v1", APIKey: "test-key"}, client)
	if err != nil {
		t.Fatal(err)
	}
	r, sessionID, steps := newScriptedRunner(t, ctx, m)
	checkLiveStream(t, ctx, r, sessionID, steps)
	if requests != 6 {
		t.Fatalf("model requests = %d, want 6", requests)
	}
}

func checkScriptedRequest(t *testing.T, r *http.Request, number int) testChatRequest {
	t.Helper()
	var request testChatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		t.Fatal(err)
	}
	if !request.Stream || len(request.Tools) != 1 || request.Tools[0].Function.Name != "next_step" {
		t.Fatalf("request %d is missing streaming tool declaration", number)
	}
	if number == 1 {
		return request
	}
	id := fmt.Sprintf("call-%d", number-1)
	if !requestHasToolExchange(request.Messages, id) {
		t.Fatalf("request %d lost tool call and response ID %s", number, id)
	}
	return request
}

func requestHasToolExchange(messages []testChatMessage, id string) bool {
	call, result := false, false
	for _, message := range messages {
		result = result || message.Role == "tool" && message.ToolCallID == id
		for _, item := range message.ToolCalls {
			call = call || item.ID == id
		}
	}
	return call && result
}

func writeScriptedResponse(w http.ResponseWriter, number int) {
	w.Header().Set("Content-Type", "text/event-stream")
	if number <= 5 {
		fmt.Fprintf(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call-%d\",\"type\":\"function\",\"function\":{\"name\":\"next_step\",\"arguments\":\"{}\"}}]},\"finish_reason\":\"tool_calls\"}]}\n\n", number)
	} else {
		_, _ = w.Write([]byte("data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"done\"},\"finish_reason\":\"stop\"}]}\n\n"))
	}
	_, _ = w.Write([]byte("data: {\"choices\":[],\"usage\":{\"prompt_tokens\":10,\"completion_tokens\":3,\"total_tokens\":13}}\n\ndata: [DONE]\n\n"))
}

func newScriptedRunner(t *testing.T, ctx context.Context, m model.LLM) (*runner.Runner, string, *int) {
	t.Helper()
	steps := 0
	stepTool, err := functiontool.New(functiontool.Config{Name: "next_step", Description: "Advance one step."},
		func(_ agent.Context, _ liveStepInput) (liveStepOutput, error) {
			steps++
			return liveStepOutput{Step: steps}, nil
		})
	if err != nil {
		t.Fatal(err)
	}
	a, err := llmagent.New(llmagent.Config{Name: "scripted_probe", Model: m, Tools: []tool.Tool{stepTool}})
	if err != nil {
		t.Fatal(err)
	}
	svc := session.InMemoryService()
	created, err := svc.Create(ctx, &session.CreateRequest{AppName: "scripted_probe", UserID: "probe"})
	if err != nil {
		t.Fatal(err)
	}
	r, err := runner.New(runner.Config{AppName: "scripted_probe", Agent: a, SessionService: svc})
	if err != nil {
		t.Fatal(err)
	}
	return r, created.Session.ID(), &steps
}
