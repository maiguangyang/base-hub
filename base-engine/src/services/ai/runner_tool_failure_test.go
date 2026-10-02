package ai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"base-engine/auth"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

func TestAIPreviewFailsWhenReadToolFails(t *testing.T) {
	requests := 0
	client := modelTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		w.Header().Set("Content-Type", "text/event-stream")
		switch requests {
		case 1:
			fmt.Fprint(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"select-1\",\"type\":\"function\",\"function\":{\"name\":\"select_tools\",\"arguments\":\"{\\\"names\\\":[\\\"read_probe\\\"]}\"}}]},\"finish_reason\":\"tool_calls\"}]}\n\n")
		case 2:
			fmt.Fprint(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"read-1\",\"type\":\"function\",\"function\":{\"name\":\"read_probe\",\"arguments\":\"{}\"}}]},\"finish_reason\":\"tool_calls\"}]}\n\n")
		default:
			fmt.Fprint(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"done\"},\"finish_reason\":\"stop\"}]}\n\n")
		}
		fmt.Fprint(w, "data: {\"choices\":[],\"usage\":{\"prompt_tokens\":2,\"completion_tokens\":1,\"total_tokens\":3}}\n\ndata: [DONE]\n\n")
	}))
	m, err := newModelWithClient(t.Context(), ModelConfig{Name: "test-model", BaseURL: "https://model.example/v1", APIKey: "test-key"}, client)
	if err != nil {
		t.Fatal(err)
	}
	spec := ToolSpec{ID: "read_probe", Name: "read_probe", Title: "Read", Mode: ModeReadOnly, Permission: "read", Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeHeadquarters}}
	catalog := &Catalog{specs: []ToolSpec{spec}, byID: map[string]ToolSpec{spec.ID: spec}}
	principal := &auth.WorkspacePrincipal{AccountID: "a", SessionID: "s", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"read": {}}}
	service, err := NewService(ServiceConfig{Model: m, Catalog: catalog, Prompt: "Test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	readTool, err := functiontool.New(functiontool.Config{Name: spec.ID, Description: "Read data"}, func(agent.Context, struct{}) (map[string]any, error) {
		return nil, errors.New("READ_FAILED")
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := service.SetTools([]tool.Tool{readTool}); err != nil {
		t.Fatal(err)
	}
	if err := service.runPreview(t.Context(), principal, "signed", "https://admin.example", previewRequest{Prompt: "Read data"}, newEventSink(t.Context())); err == nil {
		t.Fatal("preview succeeded after a failed read tool")
	}
	if requests != 3 {
		t.Fatalf("model requests = %d, want 3", requests)
	}
}
