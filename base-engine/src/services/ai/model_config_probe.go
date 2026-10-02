package ai

import (
	"context"
	"errors"
	"iter"
	"net/http"
	"sync"
	"time"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
	"google.golang.org/genai"
)

type probeInput struct {
	Value int `json:"value"`
}

type probeOutput struct {
	Value int `json:"value"`
}

const probeStepCount = 5

type probeObservation struct {
	mu       sync.Mutex
	rounds   int
	steps    int
	rejected bool
	partial  bool
	final    bool
	metered  bool
}

type observedProbeModel struct {
	model.LLM
	seen *probeObservation
}

func (m observedProbeModel) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		m.seen.beginRound()
		for response, err := range m.LLM.GenerateContent(ctx, req, stream) {
			m.seen.observeModelResponse(response)
			if !yield(response, err) {
				return
			}
		}
	}
}

func ProbeModelConnection(ctx context.Context, cfg ModelConfig) error {
	return probeModelConnectionWithClient(ctx, cfg, nil)
}

func probeModelConnectionWithClient(parent context.Context, cfg ModelConfig, client *http.Client) error {
	ctx, cancel := context.WithTimeout(parent, time.Minute)
	defer cancel()
	m, err := newModelWithClient(ctx, cfg, client)
	if err != nil {
		return err
	}
	observation := &probeObservation{}
	probeTool, err := functiontool.New(functiontool.Config{Name: "next_step", Description: "Return the supplied sequence value."},
		func(_ agent.Context, input probeInput) (probeOutput, error) {
			return observation.acceptStep(input.Value)
		})
	if err != nil {
		return ErrModelProtocol
	}
	r, sessionID, err := newProbeRunner(ctx, observedProbeModel{LLM: m, seen: observation}, probeTool)
	if err != nil {
		return err
	}
	return runProbe(ctx, r, sessionID, observation)
}

func newProbeRunner(ctx context.Context, m model.LLM, probeTool tool.Tool) (*runner.Runner, string, error) {
	a, err := llmagent.New(llmagent.Config{
		Name: "model_config_probe", Model: m,
		Instruction: "Call next_step exactly five times, one at a time, with value 1 through 5 in order. Then answer with one short sentence.",
		Tools:       []tool.Tool{probeTool},
	})
	if err != nil {
		return nil, "", ErrModelProtocol
	}
	svc := session.InMemoryService()
	created, err := svc.Create(ctx, &session.CreateRequest{AppName: "model_config_probe", UserID: "model_config_probe"})
	if err != nil {
		return nil, "", ErrModelUnavailable
	}
	r, err := runner.New(runner.Config{AppName: "model_config_probe", Agent: a, SessionService: svc})
	if err != nil {
		return nil, "", ErrModelProtocol
	}
	return r, created.Session.ID(), nil
}

func runProbe(ctx context.Context, r *runner.Runner, sessionID string, seen *probeObservation) error {
	for ev, err := range r.Run(ctx, "model_config_probe", sessionID,
		genai.NewContentFromText("Perform the five step test.", genai.RoleUser),
		agent.RunConfig{StreamingMode: agent.StreamingModeSSE}) {
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if errors.Is(err, ErrModelUnavailable) || errors.Is(err, ErrModelProtocol) {
				return err
			}
			return ErrModelUnavailable
		}
		if ev == nil {
			continue
		}
		seen.observe(ev)
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if !seen.valid() {
		return ErrModelProtocol
	}
	return nil
}

func (seen *probeObservation) observe(ev *session.Event) {
	seen.mu.Lock()
	defer seen.mu.Unlock()
	seen.partial = seen.partial || probeHasPartialText(ev)
	seen.final = seen.final || ev.IsFinalResponse()
	seen.metered = seen.metered || ev.UsageMetadata != nil && ev.UsageMetadata.TotalTokenCount > 0
}

func (seen *probeObservation) beginRound() {
	seen.mu.Lock()
	defer seen.mu.Unlock()
	seen.rounds++
}

func (seen *probeObservation) observeModelResponse(response *model.LLMResponse) {
	if response == nil || response.Partial || response.Content == nil {
		return
	}
	calls, unknown := 0, false
	for _, part := range response.Content.Parts {
		if part.FunctionCall != nil {
			calls++
			unknown = unknown || part.FunctionCall.Name != "next_step"
		}
	}
	if calls > 1 || unknown {
		seen.mu.Lock()
		seen.rejected = true
		seen.mu.Unlock()
	}
}

func (seen *probeObservation) acceptStep(value int) (probeOutput, error) {
	seen.mu.Lock()
	defer seen.mu.Unlock()
	if seen.rejected || seen.rounds != seen.steps+1 || value != seen.steps+1 || seen.steps >= probeStepCount {
		seen.rejected = true
		return probeOutput{}, ErrModelProtocol
	}
	seen.steps++
	return probeOutput{Value: value}, nil
}

func (seen *probeObservation) valid() bool {
	seen.mu.Lock()
	defer seen.mu.Unlock()
	return seen.steps == probeStepCount && seen.rounds == probeStepCount+1 &&
		!seen.rejected && seen.partial && seen.final && seen.metered
}

func probeHasPartialText(ev *session.Event) bool {
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
