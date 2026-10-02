package ai

import (
	"context"
	"errors"
	"iter"

	"google.golang.org/adk/v2/model"
)

type limitedModel struct{ state *runState }

var errModelUnavailable = errors.New("AI_MODEL_UNAVAILABLE")

func (m *limitedModel) Name() string {
	if m.state.service.config.Model != nil {
		return m.state.service.config.Model.Name()
	}
	return "configured_chat_completion"
}

func (m *limitedModel) GenerateContent(ctx context.Context, request *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		m.forwardContent(ctx, request, stream, yield)
	}
}

func (m *limitedModel) forwardContent(ctx context.Context, request *model.LLMRequest, stream bool, yield func(*model.LLMResponse, error) bool) {
	underlying, err := m.beforeCall(ctx)
	if err != nil {
		yield(nil, err)
		return
	}
	metered := false
	for response, callErr := range underlying.GenerateContent(ctx, m.state.advertisedRequest(request), stream) {
		if callErr != nil {
			yield(nil, callErr)
			return
		}
		if !metered && hasTrustedUsage(response) {
			metered = true
			if err := m.addUsage(response); err != nil {
				yield(nil, err)
				return
			}
		}
		if !yield(response, nil) {
			return
		}
	}
	if !metered {
		yield(nil, errors.New("AI_MODEL_USAGE_MISSING"))
	}
}

func hasTrustedUsage(response *model.LLMResponse) bool {
	return response != nil && response.UsageMetadata != nil && response.UsageMetadata.TotalTokenCount > 0
}

func (m *limitedModel) beforeCall(ctx context.Context) (model.LLM, error) {
	state := m.state
	if _, err := state.checkPrincipal(ctx); err != nil {
		return nil, err
	}
	state.mu.Lock()
	if state.modelCalls >= 80 {
		state.mu.Unlock()
		return nil, errors.New("AI_MODEL_CALL_LIMIT")
	}
	state.modelCalls++
	state.mu.Unlock()
	if state.service.config.ModelProvider == nil {
		return state.service.config.Model, nil
	}
	underlying, version, err := state.service.config.ModelProvider(ctx)
	if err != nil || underlying == nil || version == 0 {
		return nil, errModelUnavailable
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.modelVersion != 0 && state.modelVersion != version {
		return nil, errors.New("AI_MODEL_VERSION_CHANGED")
	}
	state.modelVersion = version
	return underlying, nil
}

func (state *runState) checkModelCurrent(ctx context.Context) error {
	if state.service.config.ModelProvider == nil {
		return nil
	}
	underlying, version, err := state.service.config.ModelProvider(ctx)
	if err != nil || underlying == nil || version == 0 {
		return errModelUnavailable
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.modelVersion == 0 || state.modelVersion != version {
		return errors.New("AI_MODEL_VERSION_CHANGED")
	}
	return nil
}

func (m *limitedModel) addUsage(response *model.LLMResponse) error {
	state := m.state
	state.mu.Lock()
	defer state.mu.Unlock()
	usage := response.UsageMetadata
	state.inputTokens += int(usage.PromptTokenCount)
	state.outputTokens += int(usage.CandidatesTokenCount)
	if state.inputTokens > state.service.config.InputTokenBudget || state.outputTokens > state.service.config.OutputTokenBudget {
		return errors.New("AI_TOKEN_BUDGET_EXCEEDED")
	}
	return nil
}
