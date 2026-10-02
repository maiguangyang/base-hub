package ai

import (
	"errors"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

type toolSelectionInput struct {
	Names []string `json:"names"`
}

const maxSelectedBusinessTools = 5

func (state *runState) selectionTool() (tool.Tool, error) {
	visible := state.service.config.Catalog.Visible(state.principal, state.phase, state.approvedIDs)
	choices := make([]any, 0, len(visible))
	labels := make([]string, 0, len(visible))
	for _, spec := range visible {
		choices = append(choices, spec.Name)
		labels = append(labels, spec.Name+"："+spec.Title)
	}
	minItems, maxItems := 1, 1
	schema := &jsonschema.Schema{Type: "object", Properties: map[string]*jsonschema.Schema{
		"names": {Type: "array", Description: "从当前工作区可用工具中选择下一步需要的 1 个工具名称。后续可再次选择其他工具。",
			Items:    &jsonschema.Schema{Type: "string", Enum: choices, Description: "当前可用工具的 name，必须来自下面的工具清单。"},
			MinItems: &minItems, MaxItems: &maxItems, UniqueItems: true},
	}, Required: []string{"names"}, AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}}}
	description := "先调用此工具选择下一步需要的后台工具；每次只能选择 1 个。选择后才能调用业务工具；需换工具时再次调用。当前可用：" + strings.Join(labels, "；")
	return functiontool.New(functiontool.Config{Name: "select_tools", Description: description, InputSchema: schema}, state.chooseTools)
}

func (state *runState) chooseTools(ctx agent.Context, input toolSelectionInput) (map[string]any, error) {
	principal, err := state.checkPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if len(input.Names) != 1 {
		return nil, errors.New("INVALID_TOOL_SELECTION")
	}
	allowed := make(map[string]struct{})
	for _, spec := range state.service.config.Catalog.Visible(principal, state.phase, state.approvedIDs) {
		allowed[spec.Name] = struct{}{}
	}
	selected := make(map[string]struct{}, len(input.Names))
	for _, name := range input.Names {
		if _, ok := allowed[name]; !ok {
			return nil, errors.New("TOOL_SELECTION_NOT_AUTHORIZED")
		}
		if _, duplicate := selected[name]; duplicate {
			return nil, errors.New("DUPLICATE_TOOL_SELECTION")
		}
		selected[name] = struct{}{}
	}
	state.mu.Lock()
	for name := range selected {
		state.rememberSelectedToolLocked(name)
	}
	state.mu.Unlock()
	return map[string]any{"selected": input.Names}, nil
}

func (state *runState) rememberSelectedToolLocked(name string) {
	if state.selectedIDs == nil {
		state.selectedIDs = make(map[string]struct{})
	}
	for index, previous := range state.selectedOrder {
		if previous == name {
			state.selectedOrder = append(state.selectedOrder[:index], state.selectedOrder[index+1:]...)
			break
		}
	}
	state.selectedOrder = append(state.selectedOrder, name)
	state.selectedIDs[name] = struct{}{}
	if len(state.selectedOrder) > maxSelectedBusinessTools {
		delete(state.selectedIDs, state.selectedOrder[0])
		state.selectedOrder = state.selectedOrder[1:]
	}
}
