package agentkit

import (
	"context"
	"errors"
	"iter"
	"testing"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
	"google.golang.org/genai"
)

type textModel struct{}

func (textModel) Name() string { return "test" }
func (textModel) GenerateContent(context.Context, *model.LLMRequest, bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		if !yield(&model.LLMResponse{Content: genai.NewContentFromText("done", "model"), Partial: true}, nil) {
			return
		}
		yield(&model.LLMResponse{Content: genai.NewContentFromText("done", "model"), TurnComplete: true}, nil)
	}
}

type failingToolModel struct {
	calls    int
	toolName string
}

func (*failingToolModel) Name() string { return "failing-tool" }
func (m *failingToolModel) GenerateContent(context.Context, *model.LLMRequest, bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		m.calls++
		if m.calls == 1 {
			name := m.toolName
			if name == "" {
				name = "read"
			}
			yield(&model.LLMResponse{Content: &genai.Content{Role: genai.RoleModel, Parts: []*genai.Part{{FunctionCall: &genai.FunctionCall{ID: "call-1", Name: name, Args: map[string]any{}}}}}, TurnComplete: true}, nil)
			return
		}
		yield(&model.LLMResponse{Content: genai.NewContentFromText("done", "model"), TurnComplete: true}, nil)
	}
}

func TestRunReportsToolFailureEvenWhenModelContinues(t *testing.T) {
	read, err := functiontool.New(functiontool.Config{Name: "read", Description: "Read data"}, func(agent.Context, struct{}) (map[string]any, error) {
		return nil, errors.New("READ_FAILED")
	})
	if err != nil {
		t.Fatal(err)
	}
	m := &failingToolModel{}
	if _, err := Run(t.Context(), Config{Name: "test_agent", Model: m, Tools: []tool.Tool{read}}, "user-1", "read", func(string) error { return nil }); err == nil {
		t.Fatal("tool failure was hidden by model continuation")
	}
	if m.calls != 2 {
		t.Fatalf("model calls = %d, want 2", m.calls)
	}
}

func TestRunReportsUnknownToolCallEvenWhenModelContinues(t *testing.T) {
	read, err := functiontool.New(functiontool.Config{Name: "read", Description: "Read data"}, func(agent.Context, struct{}) (map[string]any, error) {
		return map[string]any{"ok": true}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	m := &failingToolModel{toolName: "missing"}
	if _, err := Run(t.Context(), Config{Name: "test_agent", Model: m, Tools: []tool.Tool{read}}, "user-1", "read", func(string) error { return nil }); err == nil {
		t.Fatal("unknown tool call was hidden by model continuation")
	}
	if m.calls != 2 {
		t.Fatalf("model calls = %d, want 2", m.calls)
	}
}

func TestRunStreamsVisibleTextOnce(t *testing.T) {
	var deltas []string
	summary, err := Run(t.Context(), Config{Name: "test_agent", Instruction: "Test", Model: textModel{}}, "user-1", "hello", func(delta string) error {
		deltas = append(deltas, delta)
		return nil
	})
	if err != nil || summary != "done" || len(deltas) != 1 || deltas[0] != "done" {
		t.Fatalf("summary=%q deltas=%q err=%v", summary, deltas, err)
	}
}

func TestTextStreamKeepsFinalAnswerAfterThoughtOnlyPartial(t *testing.T) {
	partial := &session.Event{LLMResponse: model.LLMResponse{
		Content: &genai.Content{Role: genai.RoleModel, Parts: []*genai.Part{{Text: "thinking", Thought: true}}}, Partial: true,
	}}
	final := &session.Event{LLMResponse: model.LLMResponse{Content: genai.NewContentFromText("answer", "model"), TurnComplete: true}}
	var stream TextStream
	if delta := stream.Delta(partial); delta != "" {
		t.Fatalf("thought leaked into visible text: %q", delta)
	}
	if delta := stream.Delta(final); delta != "answer" {
		t.Fatalf("final answer was lost: %q", delta)
	}
}

func TestTextStreamAppendsOnlyMissingFinalSuffix(t *testing.T) {
	var stream TextStream
	for _, fragment := range []string{"do", "n"} {
		partial := &session.Event{LLMResponse: model.LLMResponse{Content: genai.NewContentFromText(fragment, "model"), Partial: true}}
		if delta := stream.Delta(partial); delta != fragment {
			t.Fatalf("partial delta = %q, want %q", delta, fragment)
		}
	}
	final := &session.Event{LLMResponse: model.LLMResponse{Content: genai.NewContentFromText("done", "model"), TurnComplete: true}}
	if delta := stream.Delta(final); delta != "e" {
		t.Fatalf("final suffix = %q, want e", delta)
	}
}
