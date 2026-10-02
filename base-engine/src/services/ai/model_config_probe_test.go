package ai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"google.golang.org/adk/v2/platform"
)

func TestChatCompletionsConfigProbeFiveSteps(t *testing.T) {
	requests := 0
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		request := checkScriptedRequest(t, r, requests)
		if requests > 1 {
			assertProbeToolResult(t, request.Messages, requests-1)
		}
		if requests <= 5 {
			fmt.Fprintf(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call-%d\",\"type\":\"function\",\"function\":{\"name\":\"next_step\",\"arguments\":\"{\\\"value\\\":%d}\"}}]},\"finish_reason\":\"tool_calls\"}]}\n\n", requests, requests)
		} else {
			_, _ = w.Write([]byte("data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"done\"},\"finish_reason\":\"stop\"}]}\n\n"))
		}
		_, _ = w.Write([]byte("data: {\"choices\":[],\"usage\":{\"prompt_tokens\":10,\"completion_tokens\":3,\"total_tokens\":13}}\n\ndata: [DONE]\n\n"))
	}))
	err := probeModelConnectionWithClient(t.Context(), ModelConfig{Name: "example", BaseURL: "https://model.example/v1", APIKey: "secret"}, client)
	if err != nil || requests != 6 {
		t.Fatalf("probe error=%v requests=%d", err, requests)
	}
}

func TestConfigProbePreservesSafeUpstreamStatus(t *testing.T) {
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"private-upstream-detail"}`))
	}))
	err := probeModelConnectionWithClient(t.Context(), ModelConfig{
		Name: "test-model", BaseURL: "https://model.example/v1", APIKey: "private-key",
	}, client)
	var upstream interface{ UpstreamStatusCode() int }
	if !errors.Is(err, ErrModelUnavailable) || !errors.As(err, &upstream) ||
		upstream.UpstreamStatusCode() != http.StatusUnauthorized || strings.Contains(err.Error(), "private-upstream-detail") {
		t.Fatalf("probe upstream failure = %v", err)
	}
}

func TestConfigProbePreservesDeadline(t *testing.T) {
	ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
	defer cancel()
	err := probeModelConnectionWithClient(ctx, ModelConfig{
		Name: "test-model", BaseURL: "https://model.example/v1", APIKey: "test-key",
	}, modelTestClient(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("expired probe made an upstream request")
	})))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("deadline probe error = %v", err)
	}
}

func TestChatCompletionsConfigProbeRejectsExtraToolCall(t *testing.T) {
	requests := 0
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		checkScriptedRequest(t, r, requests)
		if requests <= 6 {
			fmt.Fprintf(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call-%d\",\"type\":\"function\",\"function\":{\"name\":\"next_step\",\"arguments\":\"{\\\"value\\\":%d}\"}}]},\"finish_reason\":\"tool_calls\"}]}\n\n", requests, requests)
		} else {
			_, _ = w.Write([]byte("data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"done\"},\"finish_reason\":\"stop\"}]}\n\n"))
		}
		_, _ = w.Write([]byte("data: {\"choices\":[],\"usage\":{\"prompt_tokens\":10,\"completion_tokens\":3,\"total_tokens\":13}}\n\ndata: [DONE]\n\n"))
	}))
	err := probeModelConnectionWithClient(t.Context(), ModelConfig{Name: "example", BaseURL: "https://model.example/v1", APIKey: "secret"}, client)
	if !errors.Is(err, ErrModelProtocol) || requests != 7 {
		t.Fatalf("extra tool call must fail: error=%v requests=%d", err, requests)
	}
}

func TestChatCompletionsConfigProbeRejectsBatchedSteps(t *testing.T) {
	for _, test := range []struct {
		name       string
		sequential bool
	}{{"sequential", true}, {"concurrent", false}} {
		t.Run(test.name, func(t *testing.T) {
			assertProbeRejectsBatchedSteps(t, test.sequential)
		})
	}
}

func assertProbeRejectsBatchedSteps(t *testing.T, sequential bool) {
	t.Helper()
	requests := 0
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		checkScriptedRequest(t, r, requests)
		if requests == 1 {
			calls := make([]string, probeStepCount)
			for step := 1; step <= probeStepCount; step++ {
				calls[step-1] = fmt.Sprintf(`{"index":%d,"id":"call-%d","type":"function","function":{"name":"next_step","arguments":"{\"value\":%d}"}}`, step-1, step, step)
			}
			fmt.Fprintf(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[%s]},\"finish_reason\":\"tool_calls\"}]}\n\n", strings.Join(calls, ","))
		} else {
			_, _ = w.Write([]byte("data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"done\"},\"finish_reason\":\"stop\"}]}\n\n"))
		}
		_, _ = w.Write([]byte("data: {\"choices\":[],\"usage\":{\"prompt_tokens\":10,\"completion_tokens\":3,\"total_tokens\":13}}\n\ndata: [DONE]\n\n"))
	}))
	ctx := t.Context()
	if sequential {
		ctx = platform.WithTaskRunner(ctx, func(ctx context.Context, tasks []func(context.Context)) {
			for _, task := range tasks {
				task(ctx)
			}
		})
	}
	err := probeModelConnectionWithClient(ctx, ModelConfig{Name: "example", BaseURL: "https://model.example/v1", APIKey: "secret"}, client)
	if !errors.Is(err, ErrModelProtocol) || requests != 2 {
		t.Fatalf("batched calls must fail: error=%v requests=%d", err, requests)
	}
}

func assertProbeToolResult(t *testing.T, messages []testChatMessage, step int) {
	t.Helper()
	id := fmt.Sprintf("call-%d", step)
	want := fmt.Sprintf(`{"value":%d}`, step)
	for _, message := range messages {
		if message.Role == "tool" && message.ToolCallID == id && message.Content == want {
			return
		}
	}
	t.Fatalf("step %d did not return %s", step, want)
}
