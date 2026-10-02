package ai

import (
	"context"
	"errors"
	"sync"

	"base-engine/auth"
)

type runStateKey struct{}

type runState struct {
	service                                          *Service
	principal                                        *auth.WorkspacePrincipal
	token, origin                                    string
	phase                                            Phase
	planID                                           string
	runID                                            string
	approvedIDs                                      map[string]struct{}
	selectedIDs                                      map[string]struct{}
	selectedOrder                                    []string
	attestation                                      *UserAttestation
	prompt                                           string
	operatorMessage                                  string
	attachmentID                                     string
	sink                                             *eventSink
	mu                                               sync.Mutex
	modelVersion                                     uint64
	modelCalls, toolCalls, inputTokens, outputTokens int
	proposalSet                                      bool
	previewPlan                                      Plan
	previewToken                                     string
	requiredInputs                                   []string
}

func stateFromContext(ctx context.Context) (*runState, error) {
	state, ok := ctx.Value(runStateKey{}).(*runState)
	if !ok || state == nil {
		return nil, errors.New("AI_RUN_CONTEXT_MISSING")
	}
	return state, nil
}

func (state *runState) checkPrincipal(ctx context.Context) (*auth.WorkspacePrincipal, error) {
	principal, err := state.service.config.ResolvePrincipal(ctx, state.token, state.service.now())
	if err != nil || !sameRunPrincipal(state.principal, principal) {
		return nil, errors.New("AI_SESSION_REVOKED")
	}
	return principal, nil
}

func sameRunPrincipal(first, latest *auth.WorkspacePrincipal) bool {
	if first == nil || latest == nil || first.AccountID != latest.AccountID || first.SessionID != latest.SessionID || first.WorkspaceType != latest.WorkspaceType {
		return false
	}
	firstOrg, latestOrg := "", ""
	if first.OrganizationID != nil {
		firstOrg = *first.OrganizationID
	}
	if latest.OrganizationID != nil {
		latestOrg = *latest.OrganizationID
	}
	return firstOrg == latestOrg
}

func (state *runState) beforeTool(ctx context.Context, toolID string) error {
	latest, err := state.checkPrincipal(ctx)
	if err != nil {
		return err
	}
	state.mu.Lock()
	if state.toolCalls >= 100 {
		state.mu.Unlock()
		return errors.New("AI_TOOL_CALL_LIMIT")
	}
	state.toolCalls++
	sequence := state.toolCalls
	state.mu.Unlock()
	if toolID != "propose_plan" && toolID != "select_tools" {
		spec, exists := state.service.config.Catalog.Lookup(toolID)
		if !exists || !specAllowed(spec, latest, state.phase, state.approvedIDs) {
			return errors.New("AI_TOOL_NOT_AUTHORIZED")
		}
		state.mu.Lock()
		_, selected := state.selectedIDs[toolID]
		state.mu.Unlock()
		if !selected {
			return errors.New("AI_TOOL_NOT_SELECTED")
		}
	} else if toolID == "propose_plan" && state.phase != PhasePreview {
		return errors.New("AI_TOOL_NOT_AUTHORIZED")
	}
	return state.sink.Emit("tool_started", map[string]any{"toolId": toolID, "title": state.toolTitle(toolID), "sequence": sequence})
}

func (state *runState) toolTitle(toolID string) string {
	if toolID == "propose_plan" {
		return "提交写入计划"
	}
	if toolID == "select_tools" {
		return "选择后台工具"
	}
	spec, _ := state.service.config.Catalog.Lookup(toolID)
	return spec.Title
}

func (state *runState) afterTool(toolID, createdID string, resultErr error) {
	state.mu.Lock()
	sequence := state.toolCalls
	state.mu.Unlock()
	status := "SUCCESS"
	if resultErr != nil {
		status = "FAILED"
	}
	data := map[string]any{"toolId": toolID, "title": state.toolTitle(toolID), "sequence": sequence, "status": status}
	if createdID != "" {
		data["createdResourceId"] = createdID
	}
	_ = state.sink.Emit("tool_finished", data)
}

func createdResourceID(output map[string]any) string {
	for _, field := range []string{"id", "organizationId", "storeId", "membershipId"} {
		if value, ok := output[field].(string); ok && value != "" {
			return value
		}
	}
	return ""
}
