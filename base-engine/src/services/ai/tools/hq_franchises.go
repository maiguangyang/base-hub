package tools

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/jsonschema-go/jsonschema"
	"base-engine/src/services/ai"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

type franchisePageInput struct {
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
	Q        *string `json:"q,omitempty"`
	Status   *string `json:"status,omitempty"`
}
type franchiseIDInput struct {
	ID string `json:"id"`
}
type organizationIDInput struct {
	OrganizationID string `json:"organizationId"`
}
type provisionInput struct {
	Input struct {
		Code             string  `json:"code"`
		Name             string  `json:"name"`
		OwnerPhone       string  `json:"ownerPhone"`
		OwnerDisplayName string  `json:"ownerDisplayName"`
		OwnerEmail       *string `json:"ownerEmail,omitempty"`
	} `json:"input"`
}
type suspendInput struct {
	Input struct {
		OrganizationID string `json:"organizationId"`
		ReasonCode     string `json:"reasonCode"`
	} `json:"input"`
}
type initialAccountInput struct {
	OrganizationID string `json:"organizationId"`
	AccountID      string `json:"accountId"`
}

func buildHQFranchiseTools(runtime ai.FixedToolRuntime, specs []ai.ToolSpec) ([]tool.Tool, error) {
	builders := []func(ai.ToolSpec, ai.FixedToolRuntime) (tool.Tool, error){
		newFranchiseListTool, newInitialCandidatesTool, newProvisionTool,
		newResetInitialPasswordTool, newSuspendTool, newRestoreTool, newInitialAccountTool,
	}
	if len(specs) != len(builders) {
		return nil, errors.New("INCOMPLETE_HQ_FRANCHISE_TOOLS")
	}
	items := make([]tool.Tool, 0, len(specs))
	for index, build := range builders {
		item, err := build(specs[index], runtime)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

type graphqlData = map[string]json.RawMessage

func newHQGraphQLTool[TIn any](spec ai.ToolSpec, runtime ai.FixedToolRuntime, project func(graphqlData) (ai.SafeToolResult, error)) (tool.Tool, error) {
	schema, err := reviewedInputSchema(spec)
	if err != nil {
		return nil, err
	}
	return ai.NewFixedGraphQLTool[TIn, graphqlData](spec, runtime, project, schema)
}

func newFranchiseListTool(spec ai.ToolSpec, runtime ai.FixedToolRuntime) (tool.Tool, error) {
	return newHQGraphQLTool[franchisePageInput](spec, runtime, func(data graphqlData) (ai.SafeToolResult, error) {
		var page struct {
			Data        []struct{ ID, Code, Name, Status string } `json:"data"`
			Total       int                                       `json:"total"`
			CurrentPage int                                       `json:"current_page"`
			PerPage     int                                       `json:"per_page"`
			TotalPage   int                                       `json:"total_page"`
		}
		if err := json.Unmarshal(data["organizations"], &page); err != nil {
			return ai.SafeToolResult{}, err
		}
		if len(page.Data) > 50 {
			page.Data = page.Data[:50]
		}
		return ai.SafeToolResult{ModelOutput: map[string]any{"organizations": page}}, nil
	})
}

func newInitialCandidatesTool(spec ai.ToolSpec, runtime ai.FixedToolRuntime) (tool.Tool, error) {
	return newHQGraphQLTool[franchiseIDInput](spec, runtime, func(data graphqlData) (ai.SafeToolResult, error) {
		var organization struct {
			ID          string `json:"id"`
			Memberships []struct {
				ID, Status string
				Account    struct{ ID, Phone, DisplayName, Status string }
			} `json:"memberships"`
		}
		if err := json.Unmarshal(data["organization"], &organization); err != nil {
			return ai.SafeToolResult{}, err
		}
		if len(organization.Memberships) > 50 {
			organization.Memberships = organization.Memberships[:50]
		}
		return ai.SafeToolResult{ModelOutput: map[string]any{"organization": organization}}, nil
	})
}

func newProvisionTool(spec ai.ToolSpec, runtime ai.FixedToolRuntime) (tool.Tool, error) {
	return newHQGraphQLTool[provisionInput](spec, runtime, func(data graphqlData) (ai.SafeToolResult, error) {
		var result struct {
			Organization struct {
				ID string `json:"id"`
			} `json:"organization"`
			Membership struct {
				ID      string `json:"id"`
				Account struct {
					ID string `json:"id"`
				} `json:"account"`
			} `json:"membership"`
			TemporaryPassword *string `json:"temporaryPassword"`
			InvitationPending bool    `json:"invitationPending"`
		}
		if err := json.Unmarshal(data["provisionFranchise"], &result); err != nil {
			return ai.SafeToolResult{}, err
		}
		output := ai.SafeToolResult{ModelOutput: map[string]any{"organizationId": result.Organization.ID, "membershipId": result.Membership.ID, "invitationPending": result.InvitationPending}}
		if result.TemporaryPassword != nil && *result.TemporaryPassword != "" {
			output.Secret = &ai.SecretPayload{TargetAccountID: result.Membership.Account.ID, Value: *result.TemporaryPassword}
		}
		return output, nil
	})
}

func newResetInitialPasswordTool(spec ai.ToolSpec, runtime ai.FixedToolRuntime) (tool.Tool, error) {
	return newHQGraphQLTool[organizationIDInput](spec, runtime, func(data graphqlData) (ai.SafeToolResult, error) {
		var result struct{ AccountID, TemporaryPassword string }
		if err := json.Unmarshal(data["resetFranchiseInitialPassword"], &result); err != nil {
			return ai.SafeToolResult{}, err
		}
		return ai.SafeToolResult{ModelOutput: map[string]any{"accountId": result.AccountID}, Secret: &ai.SecretPayload{TargetAccountID: result.AccountID, Value: result.TemporaryPassword}}, nil
	})
}

func newSuspendTool(spec ai.ToolSpec, runtime ai.FixedToolRuntime) (tool.Tool, error) {
	return newHQGraphQLTool[suspendInput](spec, runtime, func(data graphqlData) (ai.SafeToolResult, error) {
		var result struct{ ID, Status, SuspensionReasonCode string }
		if err := json.Unmarshal(data["suspendOrganization"], &result); err != nil {
			return ai.SafeToolResult{}, err
		}
		return ai.SafeToolResult{ModelOutput: map[string]any{"id": result.ID, "status": result.Status, "suspensionReasonCode": result.SuspensionReasonCode}}, nil
	})
}

func newRestoreTool(spec ai.ToolSpec, runtime ai.FixedToolRuntime) (tool.Tool, error) {
	return newHQGraphQLTool[franchiseIDInput](spec, runtime, func(data graphqlData) (ai.SafeToolResult, error) {
		var result struct{ ID, Status string }
		if err := json.Unmarshal(data["restoreOrganization"], &result); err != nil {
			return ai.SafeToolResult{}, err
		}
		return ai.SafeToolResult{ModelOutput: map[string]any{"id": result.ID, "status": result.Status}}, nil
	})
}

func newInitialAccountTool(spec ai.ToolSpec, runtime ai.FixedToolRuntime) (tool.Tool, error) {
	inputSchema := initialAccountInputSchema(spec)
	if inputSchema == nil {
		return nil, errors.New("UNREVIEWED_HTTP_TOOL")
	}
	return functiontool.New(functiontool.Config{Name: spec.Name, Description: spec.Description, InputSchema: inputSchema}, func(ctx agent.Context, input initialAccountInput) (map[string]any, error) {
		variables, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		response, err := runtime.Call(ctx, spec, variables)
		if err != nil {
			return nil, err
		}
		if response.Status != http.StatusOK {
			return nil, errors.New("PROTECTED_CALL_FAILED")
		}
		var result struct{ OrganizationID, AccountID string }
		if err := json.Unmarshal(response.Body, &result); err != nil {
			return nil, err
		}
		if result.OrganizationID != input.OrganizationID || result.AccountID != input.AccountID {
			return nil, errors.New("INVALID_PROTECTED_RESULT")
		}
		return map[string]any{"organizationId": result.OrganizationID, "accountId": result.AccountID}, nil
	})
}

func initialAccountInputSchema(spec ai.ToolSpec) *jsonschema.Schema {
	if spec.OperationID != "http.franchiseInitialAccount.set" || spec.ID != "HqSetFranchiseInitialAccount" {
		return nil
	}
	inputSchema := strictObjectSchema()
	minimumLength := 1
	for _, name := range []string{"organizationId", "accountId"} {
		field := &jsonschema.Schema{Type: "string", MinLength: &minimumLength}
		describeInputSchema(name, field, spec.Title)
		inputSchema.Properties[name] = field
		inputSchema.Required = append(inputSchema.Required, name)
	}
	return inputSchema
}
