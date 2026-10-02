package ai

import "encoding/json"

type approvedStepView struct {
	ToolID    string         `json:"toolId"`
	Title     string         `json:"title"`
	TargetIDs []string       `json:"targetIds,omitempty"`
	ScopeIDs  []string       `json:"scopeIds,omitempty"`
	Arguments map[string]any `json:"arguments"`
	MaxCalls  int            `json:"maxCalls"`
	Sequence  int            `json:"sequence"`
	Risk      string         `json:"risk"`
}

func (state *runState) previewView(summary string) map[string]any {
	steps := make([]approvedStepView, 0, len(state.previewPlan.Steps))
	highRisk := false
	for _, step := range state.previewPlan.Steps {
		spec, _ := state.service.config.Catalog.Lookup(step.ToolID)
		var arguments map[string]any
		_ = json.Unmarshal(step.Arguments, &arguments)
		if spec.Risk == "HIGH" {
			highRisk = true
		}
		steps = append(steps, approvedStepView{ToolID: step.ToolID, Title: spec.Title, TargetIDs: step.TargetIDs, ScopeIDs: step.ScopeIDs,
			Arguments: arguments, MaxCalls: step.MaxCalls, Sequence: step.Sequence, Risk: spec.Risk})
	}
	requiredInputs := state.requiredInputs
	if requiredInputs == nil {
		requiredInputs = []string{}
	}
	return map[string]any{"summary": summary, "steps": steps, "isHighRisk": highRisk,
		"requiredInputs": requiredInputs, "previewToken": state.previewToken}
}
