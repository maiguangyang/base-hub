// Package agentkit contains reusable ADK runner primitives for Go services.
package agentkit

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/session/compaction"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/genai"
)

type Config struct {
	Name                 string
	Instruction          string
	Model                model.LLM
	Tools                []tool.Tool
	Compaction           *compaction.Config
	BeforeModelCallbacks []llmagent.BeforeModelCallback
	BeforeToolCallbacks  []llmagent.BeforeToolCallback
	AfterToolCallbacks   []llmagent.AfterToolCallback
}

// Run creates an isolated ADK session and streams user-visible text to emit.
func Run(ctx context.Context, cfg Config, userID, prompt string, emit func(string) error) (string, error) {
	if ctx == nil || cfg.Name == "" || cfg.Model == nil || userID == "" || prompt == "" || emit == nil {
		return "", errors.New("INVALID_AGENT_CONFIG")
	}
	failures := &toolFailureTracker{}
	run, sessionID, err := newRunner(ctx, cfg, userID, failures)
	if err != nil {
		return "", err
	}
	summary, err := streamRun(ctx, run, userID, sessionID, prompt, emit)
	if err != nil {
		return "", err
	}
	if err := failures.Err(); err != nil {
		return "", err
	}
	return summary, nil
}

func newRunner(ctx context.Context, cfg Config, userID string, failures *toolFailureTracker) (*runner.Runner, string, error) {
	a, err := llmagent.New(llmagent.Config{Name: cfg.Name, Model: cfg.Model, Instruction: cfg.Instruction,
		Tools: cfg.Tools, BeforeModelCallbacks: cfg.BeforeModelCallbacks,
		BeforeToolCallbacks: cfg.BeforeToolCallbacks, AfterToolCallbacks: trackToolFailures(cfg.AfterToolCallbacks, failures.record),
		OnToolErrorCallbacks: []llmagent.OnToolErrorCallback{func(_ agent.Context, _ tool.Tool, _ map[string]any, err error) (map[string]any, error) {
			failures.record(err)
			return nil, nil
		}}})
	if err != nil {
		return nil, "", err
	}
	sessions := session.InMemoryService()
	created, err := sessions.Create(ctx, &session.CreateRequest{AppName: cfg.Name, UserID: userID})
	if err != nil {
		return nil, "", err
	}
	run, err := runner.New(runner.Config{AppName: cfg.Name, Agent: a, SessionService: sessions, Compaction: cfg.Compaction})
	if err != nil {
		return nil, "", err
	}
	return run, created.Session.ID(), nil
}

func streamRun(ctx context.Context, run *runner.Runner, userID, sessionID, prompt string, emit func(string) error) (string, error) {
	var summary strings.Builder
	var textState TextStream
	for event, runErr := range run.Run(ctx, userID, sessionID, genai.NewContentFromText(prompt, "user"), agent.RunConfig{StreamingMode: agent.StreamingModeSSE}) {
		if runErr != nil {
			return "", runErr
		}
		delta := textState.Delta(event)
		if delta == "" {
			continue
		}
		summary.WriteString(delta)
		if err := emit(delta); err != nil {
			return "", err
		}
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return summary.String(), nil
}

// TextStream removes thought and tool events and avoids repeating final snapshots.
type TextStream struct{ partialText strings.Builder }

func (state *TextStream) Delta(event *session.Event) string {
	if event == nil || event.Content == nil {
		return ""
	}
	delta, toolEvent := visibleText(event.Content.Parts)
	if toolEvent {
		state.partialText.Reset()
		return ""
	}
	if event.Partial {
		state.partialText.WriteString(delta)
		return delta
	}
	return state.finalDelta(event, delta)
}

func visibleText(parts []*genai.Part) (string, bool) {
	var text strings.Builder
	for _, part := range parts {
		if part == nil || part.Thought {
			continue
		}
		if part.FunctionCall != nil || part.FunctionResponse != nil {
			return "", true
		}
		text.WriteString(part.Text)
	}
	return text.String(), false
}

func (state *TextStream) finalDelta(event *session.Event, delta string) string {
	if !event.IsFinalResponse() {
		return ""
	}
	if state.partialText.Len() > 0 {
		partial := state.partialText.String()
		state.partialText.Reset()
		if strings.HasPrefix(delta, partial) {
			return strings.TrimPrefix(delta, partial)
		}
		return ""
	}
	return delta
}
