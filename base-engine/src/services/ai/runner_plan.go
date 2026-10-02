package ai

import (
	"encoding/json"
	"errors"

	"github.com/google/jsonschema-go/jsonschema"
	"base-engine/system_prompt"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

type proposedStep struct {
	OperationID string         `json:"operationId"`
	ToolID      string         `json:"toolId"`
	TargetIDs   []string       `json:"targetIds,omitempty"`
	ScopeIDs    []string       `json:"scopeIds,omitempty"`
	Arguments   map[string]any `json:"arguments"`
	MaxCalls    int            `json:"maxCalls"`
	Sequence    int            `json:"sequence"`
}
type proposedPlan struct {
	Steps []proposedStep `json:"steps"`
}

func (state *runState) proposalTool() (tool.Tool, error) {
	description, ok := system_prompt.ToolDescription("propose_plan")
	if !ok {
		return nil, errors.New("PROPOSAL_DESCRIPTION_MISSING")
	}
	return functiontool.New(functiontool.Config{Name: "propose_plan", Description: description, InputSchema: proposalInputSchema()},
		func(ctx agent.Context, input proposedPlan) (map[string]any, error) {
			return state.acceptProposal(ctx, input)
		})
}

func proposalInputSchema() *jsonschema.Schema {
	stringField := func(description string) *jsonschema.Schema {
		return &jsonschema.Schema{Type: "string", Description: description}
	}
	ids := func(description string) *jsonschema.Schema {
		return &jsonschema.Schema{Type: "array", Description: description, Items: stringField("已核实的目标或范围 ID。")}
	}
	step := &jsonschema.Schema{Type: "object", Description: "一项按顺序执行的写入操作及其完整批准边界。", Properties: map[string]*jsonschema.Schema{
		"operationId": stringField("对应固定工具的 Engine 操作 ID；必须与工具目录一致。"),
		"toolId":      stringField("要执行的已注册写入工具名称；必须属于当前工作区并有权限。"),
		"targetIds":   ids("已存在目标的 ID 列表；须与 arguments 中的目标一致。"),
		"scopeIds":    ids("创建操作的父级范围 ID 列表；须与 arguments 中的范围一致。"),
		"arguments":   {Type: "object", Description: "本次写入的完整工具参数对象；字段、类型和目标须符合对应工具的参数 Schema。"},
		"maxCalls":    {Type: "integer", Description: "此步骤获准执行的最大调用次数，范围为 1 至 10。"},
		"sequence":    {Type: "integer", Description: "从 1 开始的执行顺序，与步骤在数组中的位置一致。"},
	}, Required: []string{"operationId", "toolId", "arguments", "maxCalls", "sequence"}, AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}}}
	return &jsonschema.Schema{Type: "object", Properties: map[string]*jsonschema.Schema{
		"steps": {Type: "array", Description: "按用户任务顺序排列的全部写入步骤；预览阶段仅提交计划，不执行写入。", Items: step},
	}, Required: []string{"steps"}, AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}}}
}

func (state *runState) acceptProposal(ctx agent.Context, input proposedPlan) (map[string]any, error) {
	state.mu.Lock()
	if state.proposalSet {
		state.mu.Unlock()
		return nil, errors.New("PROPOSAL_ALREADY_SUBMITTED")
	}
	state.proposalSet = true
	state.mu.Unlock()
	latest, err := state.checkPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	draft, err := draftFromProposal(input)
	if err != nil {
		return nil, err
	}
	draft.BindPrompt(state.prompt)
	draft.BindOperatorMessage(state.operatorMessage)
	draft.BindAttachmentID(state.attachmentID)
	plan, required, err := state.service.approval.ValidateDraft(draft, latest, state.attestation)
	if err != nil {
		return nil, err
	}
	if len(required) > 0 {
		state.requiredInputs = required
		return map[string]any{"accepted": false, "requiredInputs": required}, nil
	}
	issued, err := state.service.approval.Issue(plan)
	if err != nil {
		return nil, err
	}
	state.previewPlan, state.previewToken = plan, issued
	return map[string]any{"accepted": true}, nil
}

func draftFromProposal(input proposedPlan) (PlanDraft, error) {
	draft := PlanDraft{Steps: make([]ApprovedStep, 0, len(input.Steps))}
	for _, item := range input.Steps {
		arguments, err := json.Marshal(item.Arguments)
		if err != nil {
			return PlanDraft{}, err
		}
		draft.Steps = append(draft.Steps, ApprovedStep{OperationID: item.OperationID, ToolID: item.ToolID,
			TargetIDs: item.TargetIDs, ScopeIDs: item.ScopeIDs, Arguments: arguments, MaxCalls: item.MaxCalls, Sequence: item.Sequence})
	}
	return draft, nil
}
