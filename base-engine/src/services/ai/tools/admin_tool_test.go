package tools

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"base-engine/src/services/ai"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/genai"
)

func TestAIRegisteredToolParametersHaveTypesAndDescriptions(t *testing.T) {
	items, err := Build(ai.FixedToolRuntime{
		Call: func(agent.Context, ai.ToolSpec, json.RawMessage) (ai.FixedResponse, error) {
			return ai.FixedResponse{}, nil
		},
		DeliverSecret: func(agent.Context, ai.SecretPayload) error { return nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range items {
		declaration := item.(interface {
			Declaration() *genai.FunctionDeclaration
		}).Declaration()
		encoded, err := json.Marshal(declaration.ParametersJsonSchema)
		if err != nil {
			t.Fatal(err)
		}
		var schema map[string]any
		if err := json.Unmarshal(encoded, &schema); err != nil {
			t.Fatal(err)
		}
		checkParameterDescriptions(t, declaration.Name, schema)
	}
}

func TestAIRegisteredToolsRejectWrongArgumentTypesBeforeCallingEngine(t *testing.T) {
	for _, example := range []struct {
		id   string
		args map[string]any
	}{
		{"HqFranchises", map[string]any{"page": "first", "pageSize": 50}},
		{"FranchiseInviteStaff", map[string]any{"input": map[string]any{"phone": "01000000000", "displayName": "Staff", "roleIds": "role-1", "storeAccessMode": "ALL_STORES", "storeIds": []string{}}}},
		{"HqSetFranchiseInitialAccount", map[string]any{"organizationId": "org-1", "accountId": 42}},
		{"HqGlobalPaymentState", map[string]any{"scope": "GLOBAL", "channel": "WECHAT", "state": "DISABLED", "recordId": "", "version": "zero"}},
		{"HqPaymentConfigRead", map[string]any{"scope": "GLOBAL", "merchantPrivateKey": "secret"}},
	} {
		t.Run(example.id, func(t *testing.T) {
			called := false
			items, err := Build(ai.FixedToolRuntime{
				Call: func(agent.Context, ai.ToolSpec, json.RawMessage) (ai.FixedResponse, error) {
					called = true
					return ai.FixedResponse{}, nil
				},
				DeliverSecret: func(agent.Context, ai.SecretPayload) error { return nil },
			})
			if err != nil {
				t.Fatal(err)
			}
			for _, item := range items {
				if item.Name() != example.id {
					continue
				}
				_, err := item.(interface {
					Run(agent.Context, any) (map[string]any, error)
				}).Run(toolTestContext{}, example.args)
				if err == nil || called {
					t.Fatalf("wrong type accepted: error=%v called=%t", err, called)
				}
				return
			}
			t.Fatal("tool not registered")
		})
	}
}

func checkParameterDescriptions(t *testing.T, path string, schema map[string]any) {
	t.Helper()
	if schema["type"] == "" || schema["type"] == nil {
		t.Fatalf("%s has no type", path)
	}
	properties, _ := schema["properties"].(map[string]any)
	for name, value := range properties {
		child := value.(map[string]any)
		if _, documented := inputFieldDescriptions[name]; !documented {
			t.Fatalf("%s.%s has no reviewed parameter meaning", path, name)
		}
		description, _ := child["description"].(string)
		if strings.TrimSpace(description) == "" {
			t.Fatalf("%s.%s has no description", path, name)
		}
		checkParameterDescriptions(t, path+"."+name, child)
	}
	if item, ok := schema["items"].(map[string]any); ok {
		checkParameterDescriptions(t, path+"[]", item)
	}
}

func TestAIRegisteredToolsHaveTitleNameAndDetailedDescription(t *testing.T) {
	field, hasTitle := reflect.TypeOf(ai.ToolSpec{}).FieldByName("Title")
	if !hasTitle || field.Type.Kind() != reflect.String {
		t.Fatal("tool directory has no title field")
	}
	for _, spec := range Specs() {
		title := reflect.ValueOf(spec).FieldByName("Title").String()
		if strings.TrimSpace(title) == "" || strings.TrimSpace(spec.Name) == "" || spec.Name != spec.ID ||
			strings.TrimSpace(spec.Description) == "" || spec.Description == title {
			t.Fatalf("%s has incomplete title/name/description: title=%q description=%q", spec.ID, title, spec.Description)
		}
	}
}

func TestAIUnknownHTTPToolHasNoInitialAccountSchema(t *testing.T) {
	if schema := initialAccountInputSchema(ai.ToolSpec{ID: "UnknownHTTP", OperationID: "http.unknown"}); schema != nil {
		t.Fatal("unreviewed HTTP tool reused the initial account input schema")
	}
}

func TestAIUnknownHTTPToolCannotBeConstructed(t *testing.T) {
	_, err := newInitialAccountTool(ai.ToolSpec{ID: "UnknownHTTP", Name: "UnknownHTTP", OperationID: "http.unknown", Description: "unknown"}, ai.FixedToolRuntime{})
	if err == nil {
		t.Fatal("unreviewed HTTP tool was constructed")
	}
}

func TestAIReviewedInviteSeparatesTemporaryPassword(t *testing.T) {
	spec := specByID(t, "HqInviteAdministrator")
	var delivered ai.SecretPayload
	item, err := buildReviewedTool(spec, ai.FixedToolRuntime{
		Call: func(_ agent.Context, got ai.ToolSpec, _ json.RawMessage) (ai.FixedResponse, error) {
			if got.ID != spec.ID {
				t.Fatalf("unexpected tool %s", got.ID)
			}
			return ai.FixedResponse{Status: http.StatusOK, Body: []byte(`{"data":{"inviteOperator":{"membership":{"id":"member-1","accountId":"account-1"},"temporaryPassword":"temporary-secret","invitationPending":false}}}`)}, nil
		}, DeliverSecret: func(_ agent.Context, payload ai.SecretPayload) error { delivered = payload; return nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	runner := item.(interface {
		Run(agent.Context, any) (map[string]any, error)
	})
	result, err := runner.Run(toolTestContext{}, map[string]any{"input": map[string]any{"phone": "01000000000", "displayName": "Owner", "roleIds": []string{"role-a"}, "storeAccessMode": "ALL_STORES", "storeIds": []string{}}})
	if err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(result)
	if bytes.Contains(encoded, []byte("temporary-secret")) || delivered.Value != "temporary-secret" || delivered.TargetAccountID != "account-1" || delivered.ToolID != spec.ID {
		t.Fatalf("projection/secret = %s %+v", encoded, delivered)
	}
}

func TestAIReviewedListCapsModelOutput(t *testing.T) {
	spec := specByID(t, "FranchiseStores")
	stores := make([]map[string]string, 70)
	for index := range stores {
		stores[index] = map[string]string{"id": "store"}
	}
	body, _ := json.Marshal(map[string]any{"data": map[string]any{"stores": map[string]any{"data": stores, "total": 70}}})
	item, err := buildReviewedTool(spec, ai.FixedToolRuntime{Call: func(agent.Context, ai.ToolSpec, json.RawMessage) (ai.FixedResponse, error) {
		return ai.FixedResponse{Status: http.StatusOK, Body: body}, nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	runner := item.(interface {
		Run(agent.Context, any) (map[string]any, error)
	})
	result, err := runner.Run(toolTestContext{}, map[string]any{"page": 1, "pageSize": 50})
	if err != nil {
		t.Fatal(err)
	}
	page := result["stores"].(map[string]any)
	if got := len(page["data"].([]any)); got != 50 {
		t.Fatalf("model rows = %d", got)
	}
}

func TestAIPermissionToolsUseFixedFiftyRowPages(t *testing.T) {
	for _, id := range []string{"SystemPermissions", "TenantPermissions"} {
		spec := specByID(t, id)
		if !bytes.Contains([]byte(spec.Document), []byte("per_page:50")) || bytes.Contains([]byte(spec.Document), []byte("$pageSize")) {
			t.Fatalf("%s page size can change: %s", id, spec.Document)
		}
	}
}

func TestAIReviewedToolDeclaresItsInputFields(t *testing.T) {
	spec := specByID(t, "HqInviteAdministrator")
	item, err := buildReviewedTool(spec, ai.FixedToolRuntime{Call: func(agent.Context, ai.ToolSpec, json.RawMessage) (ai.FixedResponse, error) {
		return ai.FixedResponse{}, nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	declaration := item.(interface {
		Declaration() *genai.FunctionDeclaration
	}).Declaration()
	encoded, err := json.Marshal(declaration.ParametersJsonSchema)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{`"input"`, `"phone"`, `"displayName"`, `"roleIds"`, `"storeAccessMode"`, `"storeIds"`} {
		if !bytes.Contains(encoded, []byte(field)) {
			t.Fatalf("tool input missing %s: %s", field, encoded)
		}
	}
}

func TestAIReviewedFilterSchemaStaysSmall(t *testing.T) {
	spec := specByID(t, "FranchiseStores")
	schema, err := reviewedInputSchema(spec)
	if err != nil {
		t.Fatal(err)
	}
	filter := schema.Properties["filter"]
	if filter == nil || len(filter.Properties) > 20 {
		t.Fatalf("store filter properties = %d", len(filter.Properties))
	}
}

func TestAIReviewedWriteSchemaExcludesGovernanceFields(t *testing.T) {
	for _, id := range []string{"HqCreateDirectStore", "FranchiseUpdateStaff", "HqCreateRole"} {
		spec := specByID(t, id)
		schema, err := reviewedInputSchema(spec)
		if err != nil {
			t.Fatal(err)
		}
		input := schema.Properties["input"]
		if input == nil {
			t.Fatalf("%s input missing", id)
		}
		for _, field := range []string{"isDelete", "auditLogsIds", "roles", "stores", "invitations", "account", "submittedAt", "reviewedAt"} {
			if _, exposed := input.Properties[field]; exposed {
				t.Fatalf("%s exposed %s", id, field)
			}
		}
	}
}

func specByID(t *testing.T, id string) ai.ToolSpec {
	t.Helper()
	for _, spec := range Specs() {
		if spec.ID == id {
			return spec
		}
	}
	t.Fatalf("missing spec %s", id)
	return ai.ToolSpec{}
}
