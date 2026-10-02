package ai

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"base-engine/config"
	"base-engine/gen"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
	"google.golang.org/genai"
)

type liveStepInput struct{}

type liveStepOutput struct {
	Step int `json:"step"`
}

func liveModelConfig(t *testing.T) ModelConfig {
	t.Helper()
	if os.Getenv("AI_LIVE_TEST_FROM_STORE") == "true" {
		return liveModelConfigFromStore(t)
	}
	if os.Getenv("AI_LIVE_TEST") != "true" {
		t.Skip("live model gate is disabled")
	}
	cfg := ModelConfig{
		Name:    os.Getenv("AI_MODEL_NAME"),
		BaseURL: os.Getenv("AI_OPENAI_BASE_URL"),
		APIKey:  os.Getenv("AI_OPENAI_API_KEY"),
	}
	if cfg.Name == "" || cfg.BaseURL == "" || cfg.APIKey == "" {
		t.Fatal("AI_LIVE_TEST requires model name, base URL, and API key")
	}
	return cfg
}

func liveModelConfigFromStore(t *testing.T) ModelConfig {
	t.Helper()
	security, err := config.LoadAIModelSecurityConfig()
	if err != nil {
		t.Fatal(err)
	}
	db, err := gen.OpenDBFromEnvVars("")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	modelConfig, err := NewModelConfigStore(db.Query(), security, nil).Active(t.Context())
	if err != nil {
		t.Fatalf("active model config: %v", err)
	}
	return modelConfig
}

func TestChatCompletionsLive(t *testing.T) {
	cfg := liveModelConfig(t)
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()
	m, err := NewModel(ctx, cfg)
	if err != nil {
		t.Fatalf("NewModel: %v", err)
	}
	steps := 0
	stepTool, err := functiontool.New(functiontool.Config{
		Name: "next_step", Description: "Return the next sequence number.",
	}, func(_ agent.Context, _ liveStepInput) (liveStepOutput, error) {
		steps++
		return liveStepOutput{Step: steps}, nil
	})
	if err != nil {
		t.Fatalf("functiontool.New: %v", err)
	}
	a, err := llmagent.New(llmagent.Config{
		Name: "chat_completions_probe", Model: m,
		Instruction: "Call next_step exactly five times, one call at a time. Use each returned step before the next call. After step 5, give a short final sentence.",
		Tools:       []tool.Tool{stepTool},
	})
	if err != nil {
		t.Fatalf("llmagent.New: %v", err)
	}
	svc := session.InMemoryService()
	created, err := svc.Create(ctx, &session.CreateRequest{AppName: "ai_probe", UserID: "probe"})
	if err != nil {
		t.Fatalf("session.Create: %v", err)
	}
	r, err := runner.New(runner.Config{AppName: "ai_probe", Agent: a, SessionService: svc})
	if err != nil {
		t.Fatalf("runner.New: %v", err)
	}
	checkLiveStream(t, ctx, r, created.Session.ID(), &steps)
	checkLiveCancellation(t, r, created.Session.ID())
}

func checkLiveStream(t *testing.T, ctx context.Context, r *runner.Runner, sessionID string, steps *int) {
	t.Helper()
	var seen liveObservation
	for ev, err := range r.Run(ctx, "probe", sessionID,
		genai.NewContentFromText("Complete the five steps.", genai.RoleUser),
		agent.RunConfig{StreamingMode: agent.StreamingModeSSE}) {
		if err != nil {
			t.Fatalf("runner.Run: %v", err)
		}
		if ev == nil {
			continue
		}
		seen.observe(ev)
	}
	if *steps != 5 || !seen.partial || !seen.final || !seen.metered {
		t.Fatalf("live probe: steps=%d, observation=%+v", *steps, seen)
	}
}

type liveObservation struct{ partial, final, metered bool }

func (seen *liveObservation) observe(ev *session.Event) {
	seen.partial = seen.partial || hasPartialText(ev)
	seen.final = seen.final || ev.IsFinalResponse()
	seen.metered = seen.metered || ev.UsageMetadata != nil && ev.UsageMetadata.TotalTokenCount > 0
}

func hasPartialText(ev *session.Event) bool {
	if !ev.Partial || ev.Content == nil {
		return false
	}
	for _, part := range ev.Content.Parts {
		if part.Text != "" {
			return true
		}
	}
	return false
}

func checkLiveCancellation(t *testing.T, r *runner.Runner, sessionID string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	emitted, final := false, false
	for ev, err := range r.Run(ctx, "probe", sessionID,
		genai.NewContentFromText("Write a detailed answer in several sentences.", genai.RoleUser),
		agent.RunConfig{StreamingMode: agent.StreamingModeSSE}) {
		if err != nil && !errors.Is(err, context.Canceled) { t.Fatalf("cancellation returned: %v", err) }
		if ev != nil {
			final = final || ev.IsFinalResponse()
			if !emitted {
				emitted = true
				cancel()
			}
		}
	}
	if !emitted || final || !errors.Is(ctx.Err(), context.Canceled) {
		t.Fatalf("cancellation probe: emitted=%t, final=%t, context=%v", emitted, final, ctx.Err())
	}
}
