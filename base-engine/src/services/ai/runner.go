package ai

import (
	"context"
	"errors"

	"base-engine/agentkit"
	"base-engine/auth"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/session/compaction"
	"google.golang.org/adk/v2/tool"
)

func (s *Service) runPreview(ctx context.Context, principal *auth.WorkspacePrincipal, token, origin string, input previewRequest, sink *eventSink) error {
	prompt := conversationPrompt(input.History, input.Prompt)
	if input.AttachmentID != "" {
		if s.config.ValidateImageAttachment == nil {
			return errors.New("IMAGE_ATTACHMENT_UNAVAILABLE")
		}
		if err := s.config.ValidateImageAttachment(ctx, principal, input.AttachmentID); err != nil {
			return err
		}
		prompt += "\nAttachment-ID: " + input.AttachmentID
	}
	state := &runState{service: s, principal: principal, token: token, origin: origin, phase: PhasePreview, sink: sink, attestation: input.Attestation, prompt: prompt, operatorMessage: input.Prompt, attachmentID: input.AttachmentID, runID: runIDFromContext(ctx)}
	summary, err := s.runAgent(ctx, state, prompt)
	if err != nil {
		return err
	}
	if state.previewPlan.ID == "" && len(state.requiredInputs) == 0 {
		return nil
	}
	return sink.Emit("preview_ready", state.previewView(summary))
}

func (s *Service) runApproved(ctx context.Context, principal *auth.WorkspacePrincipal, token, origin string, plan Plan, sink *eventSink) error {
	prompt := s.approval.promptFor(plan.ID)
	if prompt == "" {
		return errors.New("AI_PLAN_PROMPT_MISSING")
	}
	defer s.approval.Cancel(plan.ID)
	approved := map[string]struct{}{}
	for _, step := range plan.Steps {
		approved[step.ToolID] = struct{}{}
	}
	state := &runState{service: s, principal: principal, token: token, origin: origin, phase: PhaseRun, planID: plan.ID, approvedIDs: approved, sink: sink, runID: runIDFromContext(ctx), attachmentID: plan.AttachmentID}
	approvedPrompt, err := approvedRunPrompt(prompt, plan)
	if err != nil {
		return err
	}
	_, err = s.runAgent(ctx, state, approvedPrompt)
	if err != nil {
		return err
	}
	return s.approval.Complete(plan.ID)
}

func (s *Service) runAgent(ctx context.Context, state *runState, prompt string) (string, error) {
	selected, err := s.toolsForRun(state)
	if err != nil {
		return "", err
	}
	runContext := context.WithValue(ctx, runStateKey{}, state)
	summary, err := agentkit.Run(runContext, agentkit.Config{Name: "admin_ai", Model: &limitedModel{state: state}, Instruction: s.config.Prompt,
		Tools: selected, Compaction: s.compactionConfig(),
		BeforeModelCallbacks: []llmagent.BeforeModelCallback{state.beforeModelCallback},
		BeforeToolCallbacks:  []llmagent.BeforeToolCallback{state.beforeToolCallback},
		AfterToolCallbacks:   []llmagent.AfterToolCallback{state.afterToolCallback}}, state.principal.AccountID, prompt,
		func(delta string) error { return state.sink.Emit("text_delta", map[string]string{"text": delta}) })
	if err != nil {
		return "", err
	}
	if _, err := state.checkPrincipal(ctx); err != nil {
		return "", err
	}
	return summary, nil
}

func (s *Service) compactionConfig() *compaction.Config {
	return &compaction.Config{TokenThreshold: s.config.ModelContextTokens / 2, EventRetentionSize: 6}
}

func (state *runState) beforeModelCallback(ctx agent.Context, _ *model.LLMRequest) (*model.LLMResponse, error) {
	_, err := state.checkPrincipal(ctx)
	return nil, err
}

func (state *runState) beforeToolCallback(ctx agent.Context, item tool.Tool, _ map[string]any) (map[string]any, error) {
	return nil, state.beforeTool(ctx, item.Name())
}

func (state *runState) afterToolCallback(_ agent.Context, item tool.Tool, _, output map[string]any, err error) (map[string]any, error) {
	createdID, err := state.recordToolOutcome(item.Name(), output, err)
	state.afterTool(item.Name(), createdID, err)
	return nil, err
}

func (state *runState) recordToolOutcome(toolID string, output map[string]any, resultErr error) (string, error) {
	spec, ok := state.service.config.Catalog.Lookup(toolID)
	if state.phase != PhaseRun || !ok || spec.Mode != ModeWrite {
		return "", resultErr
	}
	if resultErr == nil {
		resultErr = state.service.approval.RecordSuccess(state.planID, toolID)
	}
	if resultErr != nil {
		state.service.approval.Cancel(state.planID)
		return "", resultErr
	}
	if spec.WriteKind == WriteCreate {
		return createdResourceID(output), nil
	}
	return "", nil
}
