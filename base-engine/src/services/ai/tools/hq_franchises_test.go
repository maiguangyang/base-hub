package tools

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"base-engine/src/services/ai"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool/toolconfirmation"
)

type toolTestContext struct{ agent.Context }

func (toolTestContext) ToolConfirmation() *toolconfirmation.ToolConfirmation { return nil }

func TestHQFranchiseResetPasswordSeparatesSecret(t *testing.T) {
	spec := hqFranchiseSpecs()[3]
	var delivered ai.SecretPayload
	item, err := newResetInitialPasswordTool(spec, ai.FixedToolRuntime{
		Call: func(_ agent.Context, got ai.ToolSpec, arguments json.RawMessage) (ai.FixedResponse, error) {
			if got.ID != spec.ID || string(arguments) != `{"organizationId":"org-a"}` {
				t.Fatalf("call = %s %s", got.ID, arguments)
			}
			return ai.FixedResponse{Status: http.StatusOK, Body: []byte(`{"data":{"resetFranchiseInitialPassword":{"accountId":"account-a","temporaryPassword":"temporary-secret"}}}`)}, nil
		}, DeliverSecret: func(_ agent.Context, payload ai.SecretPayload) error { delivered = payload; return nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	runner := item.(interface {
		Run(agent.Context, any) (map[string]any, error)
	})
	result, err := runner.Run(toolTestContext{}, map[string]any{"organizationId": "org-a"})
	if err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(result)
	if bytes.Contains(encoded, []byte("temporary-secret")) || delivered.Value != "temporary-secret" || delivered.ToolID != spec.ID || delivered.TargetAccountID != "account-a" {
		t.Fatalf("projection or secret delivery = %s %+v", encoded, delivered)
	}
}

func TestHQFranchiseProvisionSecretsAndCreatedIDs(t *testing.T) {
	spec := hqFranchiseSpecs()[2]
	var delivered ai.SecretPayload
	item, err := newProvisionTool(spec, ai.FixedToolRuntime{
		Call: func(_ agent.Context, _ ai.ToolSpec, _ json.RawMessage) (ai.FixedResponse, error) {
			return ai.FixedResponse{Status: http.StatusOK, Body: []byte(`{"data":{"provisionFranchise":{"organization":{"id":"org-new"},"membership":{"id":"member-new","account":{"id":"account-new"}},"temporaryPassword":"temporary-secret","invitationPending":false}}}`)}, nil
		}, DeliverSecret: func(_ agent.Context, payload ai.SecretPayload) error { delivered = payload; return nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	runner := item.(interface {
		Run(agent.Context, any) (map[string]any, error)
	})
	result, err := runner.Run(toolTestContext{}, map[string]any{"input": map[string]any{"code": "new", "name": "New", "ownerPhone": "01000000000", "ownerDisplayName": "Owner"}})
	if err != nil {
		t.Fatal(err)
	}
	if result["organizationId"] != "org-new" || delivered.TargetAccountID != "account-new" || delivered.ToolID != spec.ID {
		t.Fatalf("result = %+v, secret = %+v", result, delivered)
	}
	encoded, _ := json.Marshal(result)
	if bytes.Contains(encoded, []byte(delivered.Value)) {
		t.Fatalf("secret entered model result: %s", encoded)
	}
}

func TestAIGenericFixedVariablesToolSchema(t *testing.T) {
	spec := hqFranchiseSpecs()[5]
	_, err := ai.NewFixedGraphQLTool[map[string]any, graphqlData](spec, ai.FixedToolRuntime{Call: func(agent.Context, ai.ToolSpec, json.RawMessage) (ai.FixedResponse, error) {
		return ai.FixedResponse{}, nil
	}}, func(graphqlData) (ai.SafeToolResult, error) {
		return ai.SafeToolResult{ModelOutput: map[string]any{"id": "a"}}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHQFranchiseWriteRejectsMissingGraphQLResult(t *testing.T) {
	spec := hqFranchiseSpecs()[5]
	for _, body := range []string{`{"data":{}}`, `{"data":{"restoreOrganization":null}}`} {
		item, err := newRestoreTool(spec, ai.FixedToolRuntime{Call: func(agent.Context, ai.ToolSpec, json.RawMessage) (ai.FixedResponse, error) {
			return ai.FixedResponse{Status: http.StatusOK, Body: []byte(body)}, nil
		}})
		if err != nil {
			t.Fatal(err)
		}
		runner := item.(interface {
			Run(agent.Context, any) (map[string]any, error)
		})
		if result, err := runner.Run(toolTestContext{}, map[string]any{"id": "org-a"}); err == nil {
			t.Fatalf("missing result accepted: body=%s result=%+v", body, result)
		}
	}
}
