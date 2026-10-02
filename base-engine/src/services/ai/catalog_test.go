package ai

import (
	"bytes"
	"encoding/json"
	"testing"

	"base-engine/auth"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool/toolconfirmation"
)

func TestVisibleToolsRespectWorkspacePermissionAndPhase(t *testing.T) {
	inventory := []ContractRecord{
		{OperationID: "graphql.query.stores", Protocol: "GRAPHQL_QUERY", Availability: "CALLABLE", Classifications: []string{"HEADQUARTERS_ADMIN", "FRANCHISE_ADMIN"}},
		{OperationID: "graphql.mutation.createStore", Protocol: "GRAPHQL_MUTATION", Availability: "CALLABLE", Classifications: []string{"FRANCHISE_ADMIN"}},
		{OperationID: "graphql.mutation.login", Protocol: "GRAPHQL_MUTATION", Availability: "NON_CALLABLE", Classifications: []string{"OTHER_EXTERNAL_API"}, NonCallableReason: "SESSION_OR_IDENTITY_OPERATION"},
	}
	specs := []ToolSpec{
		{ID: "hq_stores", OperationID: "graphql.query.stores", Document: "query { stores { id } }", Description: "Read stores", Mode: ModeReadOnly, Risk: "LOW", Permission: "store:read_all", Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeHeadquarters}, OutputFields: []string{"id"}},
		{ID: "franchise_stores", OperationID: "graphql.query.stores", Document: "query { stores { id } }", Description: "Read tenant stores", Mode: ModeReadOnly, Risk: "LOW", Permission: "store:read", Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeFranchise}, OutputFields: []string{"id"}},
		{ID: "franchise_create_store", OperationID: "graphql.mutation.createStore", Document: "mutation { createStore { id } }", Description: "Create store", Mode: ModeWrite, Risk: "HIGH", Permission: "store:create", Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeFranchise}, OutputFields: []string{"id"}},
	}
	catalog, err := NewCatalog(inventory, specs)
	if err != nil {
		t.Fatal(err)
	}
	hq := &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"store:read_all": {}}}
	if visible := catalog.Visible(hq, PhasePreview, nil); len(visible) != 1 || visible[0].ID != "hq_stores" {
		t.Fatalf("HQ preview: %+v", visible)
	}
	franchise := &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeFranchise, Permissions: map[string]struct{}{"store:read": {}, "store:create": {}}}
	if visible := catalog.Visible(franchise, PhasePreview, nil); len(visible) != 1 || visible[0].ID != "franchise_stores" {
		t.Fatalf("franchise preview: %+v", visible)
	}
	if visible := catalog.Visible(franchise, PhaseRun, map[string]struct{}{"franchise_create_store": {}}); len(visible) != 1 || visible[0].ID != "franchise_create_store" {
		t.Fatalf("franchise run: %+v", visible)
	}
	franchise.Permissions = map[string]struct{}{"store:read": {}}
	if visible := catalog.Visible(franchise, PhaseRun, map[string]struct{}{"franchise_create_store": {}}); len(visible) != 0 {
		t.Fatalf("revoked permission: %+v", visible)
	}
	if _, ok := catalog.Lookup("missing"); ok {
		t.Fatal("unknown tool became visible")
	}
}

func TestVisibleToolsRejectUnreviewedOperation(t *testing.T) {
	_, err := NewCatalog([]ContractRecord{{OperationID: "graphql.mutation.login", Availability: "NON_CALLABLE"}}, []ToolSpec{{ID: "login", OperationID: "graphql.mutation.login", Document: "mutation { login }", Description: "Login", Mode: ModeWrite, Risk: "HIGH", Permission: "login", Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeHeadquarters}, OutputFields: []string{"id"}}})
	if err == nil {
		t.Fatal("non-callable operation entered catalog")
	}
}

func TestCatalogRejectsDiscoveryPlaceholderWithoutTools(t *testing.T) {
	_, err := NewCatalog([]ContractRecord{{OperationID: "graphql.query.accounts", Availability: "NON_CALLABLE", Classifications: []string{"OTHER_EXTERNAL_API"}, NonCallableReason: "NOT_REVIEWED_FOR_AI"}}, nil)
	if err == nil {
		t.Fatal("discovery placeholder entered the runtime catalog")
	}
}

func TestCallableContractRequiresRegisteredTool(t *testing.T) {
	_, err := NewCatalog([]ContractRecord{{OperationID: "graphql.query.stores", Availability: "CALLABLE"}}, nil)
	if err == nil {
		t.Fatal("callable operation without an executable tool was accepted")
	}
}

func TestCallableContractRequiresToolInEveryClassifiedWorkspace(t *testing.T) {
	contract := ContractRecord{OperationID: "graphql.query.stores", Protocol: "GRAPHQL_QUERY", Availability: "CALLABLE", Classifications: []string{"HEADQUARTERS_ADMIN", "FRANCHISE_ADMIN"}}
	hq := ToolSpec{ID: "hq_stores", OperationID: contract.OperationID, Document: "query { stores { id } }", Description: "Read stores", Mode: ModeReadOnly, Risk: "LOW", Permission: "store:read_all", Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeHeadquarters}, OutputFields: []string{"stores"}}
	if _, err := NewCatalog([]ContractRecord{contract}, []ToolSpec{hq}); err == nil {
		t.Fatal("franchise workspace has no tool but callable contract was accepted")
	}
	franchise := hq
	franchise.ID, franchise.Permission, franchise.Workspaces = "franchise_stores", "store:read", []auth.WorkspaceType{auth.WorkspaceTypeFranchise}
	if _, err := NewCatalog([]ContractRecord{contract}, []ToolSpec{franchise}); err == nil {
		t.Fatal("headquarters workspace has no tool but callable contract was accepted")
	}
}

// TestCatalogRejectsUnknownAvailability 验证非法状态不能绕过可调用工具覆盖检查。
func TestCatalogRejectsUnknownAvailability(t *testing.T) {
	_, err := NewCatalog([]ContractRecord{{OperationID: "graphql.query.stores", Availability: "CALLABE"}}, nil)
	if err == nil {
		t.Fatal("unknown availability bypassed callable coverage")
	}
}

func TestVisibleToolsRequireAdditionalPermissionsAndFixedField(t *testing.T) {
	contract := ContractRecord{OperationID: "graphql.mutation.createStore", Protocol: "GRAPHQL_MUTATION", Availability: "CALLABLE", Classifications: []string{"FRANCHISE_ADMIN"}}
	spec := ToolSpec{ID: "create_store", OperationID: contract.OperationID, Document: "mutation { createStore { id } }", Description: "Create store", Mode: ModeWrite, Risk: "MEDIUM", Permission: "store:create", AdditionalPermissions: []string{"operatorRole:read"}, Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeFranchise}, OutputFields: []string{"id"}}
	catalog, err := NewCatalog([]ContractRecord{contract}, []ToolSpec{spec})
	if err != nil {
		t.Fatal(err)
	}
	principal := &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeFranchise, Permissions: map[string]struct{}{"store:create": {}}}
	approved := map[string]struct{}{"create_store": {}}
	if len(catalog.Visible(principal, PhaseRun, approved)) != 0 {
		t.Fatal("missing secondary permission was ignored")
	}
	principal.Permissions["operatorRole:read"] = struct{}{}
	if len(catalog.Visible(principal, PhaseRun, approved)) != 1 {
		t.Fatal("valid secondary permission was ignored")
	}
	spec.Document = "mutation { deleteStores(id:[]) }"
	if _, err := NewCatalog([]ContractRecord{contract}, []ToolSpec{spec}); err == nil {
		t.Fatal("mismatched fixed document accepted")
	}
}

type catalogTestContext struct{ agent.Context }

func (catalogTestContext) ToolConfirmation() *toolconfirmation.ToolConfirmation { return nil }

func TestFixedGraphQLToolKeepsSecretOutOfModelResult(t *testing.T) {
	type output struct {
		CreateAccount struct {
			AccountID string `json:"accountId"`
			Password  string `json:"password"`
		} `json:"createAccount"`
	}
	spec := ToolSpec{ID: "create_account", OperationID: "graphql.mutation.createAccount", Document: "mutation Create($name:String!){ createAccount(name:$name){accountId password} }", Description: "Create account", Mode: ModeWrite, OutputFields: []string{"accountId"}}
	var delivered SecretPayload
	runtime := FixedToolRuntime{Call: func(_ agent.Context, got ToolSpec, variables json.RawMessage) (FixedResponse, error) {
		if got.Document != spec.Document || string(variables) != `{"name":"New"}` {
			t.Fatalf("call = %+v, %s", got, variables)
		}
		return FixedResponse{Status: 200, Body: []byte(`{"data":{"createAccount":{"accountId":"a-1","password":"secret-1"}}}`)}, nil
	}, DeliverSecret: func(_ agent.Context, secret SecretPayload) error { delivered = secret; return nil }}
	wrapped, err := NewFixedGraphQLTool[struct {
		Name string `json:"name"`
	}, output](spec, runtime, func(result output) (SafeToolResult, error) {
		return SafeToolResult{ModelOutput: map[string]any{"accountId": result.CreateAccount.AccountID}, Secret: &SecretPayload{TargetAccountID: result.CreateAccount.AccountID, Value: result.CreateAccount.Password}}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	runner := wrapped.(interface {
		Run(agent.Context, any) (map[string]any, error)
	})
	result, err := runner.Run(catalogTestContext{}, map[string]any{"name": "New"})
	if err != nil {
		t.Fatal(err)
	}
	if delivered.Value != "secret-1" || result["accountId"] != "a-1" {
		t.Fatalf("secret/result = %+v / %+v", delivered, result)
	}
	if delivered.ToolID != spec.ID {
		t.Fatalf("secret tool id = %q", delivered.ToolID)
	}
	encoded, _ := json.Marshal(result)
	if string(encoded) == "" || bytes.Contains(encoded, []byte(delivered.Value)) {
		t.Fatalf("secret leaked: %s", encoded)
	}
}

func TestFixedGraphQLToolPropagatesSchemaErrors(t *testing.T) {
	spec := ToolSpec{ID: "example", OperationID: "graphql.query.example", Document: "query { example }", Description: "Read example"}
	runtime := FixedToolRuntime{Call: func(agent.Context, ToolSpec, json.RawMessage) (FixedResponse, error) { return FixedResponse{}, nil }}
	if _, err := NewFixedGraphQLTool[string, struct{}](spec, runtime, func(struct{}) (SafeToolResult, error) { return SafeToolResult{}, nil }); err == nil {
		t.Fatal("invalid input schema accepted")
	}
	if _, err := NewFixedGraphQLTool[struct{ Name string }, chan int](spec, runtime, func(chan int) (SafeToolResult, error) { return SafeToolResult{}, nil }); err == nil {
		t.Fatal("invalid output schema accepted")
	}
}

func TestFixedGraphQLToolRejectsUnapprovedProjectionAndSecret(t *testing.T) {
	spec := ToolSpec{OutputFields: []string{"id"}}
	if err := validateOutputFields(spec, map[string]any{"password": "leak"}); err == nil {
		t.Fatal("unapproved projection accepted")
	}
	if err := deliverSafeResult(catalogTestContext{}, SafeToolResult{ModelOutput: map[string]any{"id": "secret-1"}, Secret: &SecretPayload{TargetAccountID: "a", Value: "secret-1"}}, FixedToolRuntime{DeliverSecret: func(agent.Context, SecretPayload) error { t.Fatal("secret delivered after leak"); return nil }}); err == nil {
		t.Fatal("secret leak accepted")
	}
}
