package ai

import (
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

func (state *runState) advertisedRequest(request *model.LLMRequest) *model.LLMRequest {
	if request == nil || request.Config == nil || len(request.Config.Tools) == 0 {
		return request
	}
	state.mu.Lock()
	selected := make(map[string]struct{}, len(state.selectedIDs))
	for name := range state.selectedIDs {
		selected[name] = struct{}{}
	}
	state.mu.Unlock()
	copyRequest := *request
	copyConfig := *request.Config
	copyConfig.Tools = nil
	for _, item := range request.Config.Tools {
		if filtered := filteredModelTool(item, selected); filtered != nil {
			copyConfig.Tools = append(copyConfig.Tools, filtered)
		}
	}
	copyRequest.Config = &copyConfig
	return &copyRequest
}

func filteredModelTool(item *genai.Tool, selected map[string]struct{}) *genai.Tool {
	if item == nil {
		return nil
	}
	filtered := *item
	filtered.FunctionDeclarations = nil
	for _, declaration := range item.FunctionDeclarations {
		if declaration != nil && advertisedToolName(declaration.Name, selected) {
			filtered.FunctionDeclarations = append(filtered.FunctionDeclarations, declaration)
		}
	}
	if len(filtered.FunctionDeclarations) == 0 {
		return nil
	}
	return &filtered
}

func advertisedToolName(name string, selected map[string]struct{}) bool {
	if name == "select_tools" || name == "propose_plan" {
		return true
	}
	_, ok := selected[name]
	return ok
}
