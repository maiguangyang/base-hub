package tools

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/src/services/ai"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
)

func TestPaymentConfigToolsRegisterOnlyApprovedOperations(t *testing.T) {
	want := map[string]string{
		"HqPaymentConfigRead":       "http.paymentConfig.read",
		"HqGlobalPaymentState":      "http.paymentConfig.state",
		"HqFranchisePaymentState":   "http.paymentConfig.state",
		"HqStorePaymentState":       "http.paymentConfig.state",
		"HqFranchisePaymentRestore": "http.paymentConfig.restoreInheritance",
		"HqStorePaymentRestore":     "http.paymentConfig.restoreInheritance",
	}
	for _, spec := range Specs() {
		if operation, ok := want[spec.ID]; ok {
			if spec.OperationID != operation || spec.Permission == "" {
				t.Fatalf("invalid payment tool: %+v", spec)
			}
			delete(want, spec.ID)
		}
		if spec.OperationID == "http.paymentConfig.save" {
			t.Fatal("credential save became callable")
		}
	}
	if len(want) != 0 {
		t.Fatalf("missing payment tools: %v", want)
	}
}

func TestPaymentRestoreSchemaRejectsMissingRecordBeforeApproval(t *testing.T) {
	for _, spec := range Specs() {
		if spec.OperationID != "http.paymentConfig.restoreInheritance" {
			continue
		}
		schema, err := paymentInputSchema(spec).Resolve(nil)
		if err != nil {
			t.Fatal(err)
		}
		input := map[string]any{"scope": paymentToolScope(spec.ID), "channel": "WECHAT", "recordId": "", "version": 0}
		if input["scope"] == "FRANCHISE" {
			input["organizationId"] = "org-a"
		} else {
			input["storeId"] = "store-a"
		}
		for _, bad := range []map[string]any{input, {"scope": input["scope"], "channel": "WECHAT", "recordId": "record-a", "version": 0}} {
			if bad["scope"] == "FRANCHISE" {
				bad["organizationId"] = "org-a"
			} else {
				bad["storeId"] = "store-a"
			}
			if err := schema.Validate(bad); err == nil {
				t.Fatalf("invalid restore arguments approved for %s: %v", spec.ID, bad)
			}
		}
		input["recordId"], input["version"] = "record-a", 1
		if err := schema.Validate(input); err != nil {
			t.Fatalf("valid restore rejected: %v", err)
		}
	}
}

func TestPaymentRestoreInvalidDraftRejectedBeforeApproval(t *testing.T) {
	var spec ai.ToolSpec
	for _, candidate := range Specs() {
		if candidate.ID == "HqFranchisePaymentRestore" {
			spec = candidate
			break
		}
	}
	resolved, err := paymentInputSchema(spec).Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	spec.InputSchema = resolved
	contract := ai.ContractRecord{OperationID: spec.OperationID, Protocol: "HTTP_ROUTE", Method: spec.Method, Path: spec.Path, Availability: "CALLABLE", Classifications: []string{"HEADQUARTERS_ADMIN"}}
	catalog, err := ai.NewCatalog([]ai.ContractRecord{contract}, []ai.ToolSpec{spec})
	if err != nil {
		t.Fatal(err)
	}
	approval, err := ai.NewApprovalService(catalog)
	if err != nil {
		t.Fatal(err)
	}
	principal := &auth.WorkspacePrincipal{AccountID: "account-a", SessionID: "session-a", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"paymentConfig:read": {}, "paymentConfig:manage": {}}}
	for _, args := range []string{
		`{"scope":"FRANCHISE","organizationId":"org-a","channel":"WECHAT","recordId":"","version":0}`,
		`{"scope":"FRANCHISE","organizationId":"org-a","channel":"WECHAT","recordId":"record-a","version":0}`,
	} {
		draft := ai.PlanDraft{Steps: []ai.ApprovedStep{{ToolID: spec.ID, OperationID: spec.OperationID, Arguments: json.RawMessage(args), MaxCalls: 1, Sequence: 1}}}
		draft.BindPromptHash("prompt-sha256")
		if _, _, err := approval.ValidateDraft(draft, principal, nil); err == nil {
			t.Fatalf("invalid restore approved: %s", args)
		}
	}
}

func TestPaymentProtectedErrorsPreserveSafeCodes(t *testing.T) {
	for _, example := range []struct {
		status int
		code   string
	}{
		{400, "VALIDATION_FAILED"}, {403, "PERMISSION_DENIED"}, {409, "CONFLICT"}, {503, "CONFIG_UNAVAILABLE"},
	} {
		items, err := Build(ai.FixedToolRuntime{Call: func(agent.Context, ai.ToolSpec, json.RawMessage) (ai.FixedResponse, error) {
			return ai.FixedResponse{Status: example.status, Body: []byte(`{"code":"` + example.code + `","secret":"private-value"}`)}, nil
		}, DeliverSecret: func(agent.Context, ai.SecretPayload) error { return nil }})
		if err != nil {
			t.Fatal(err)
		}
		index := slices.IndexFunc(items, func(item tool.Tool) bool { return item.Name() == "HqPaymentConfigRead" })
		_, runErr := items[index].(interface {
			Run(agent.Context, any) (map[string]any, error)
		}).Run(toolTestContext{}, map[string]any{"scope": "GLOBAL"})
		if runErr == nil || runErr.Error() != example.code || strings.Contains(runErr.Error(), "private-value") {
			t.Fatalf("status %d: %v", example.status, runErr)
		}
	}
}

func TestPaymentConfigReadToolProjectsOnlyMaskedStatus(t *testing.T) {
	called := false
	items, err := Build(ai.FixedToolRuntime{Call: func(_ agent.Context, spec ai.ToolSpec, raw json.RawMessage) (ai.FixedResponse, error) {
		called = true
		if spec.OperationID != "http.paymentConfig.read" || string(raw) != `{"scope":"GLOBAL"}` {
			t.Fatalf("wrong protected call: %s %s", spec.OperationID, raw)
		}
		return ai.FixedResponse{Status: http.StatusOK, Body: []byte(`{"channels":[{"own":{"scope":"GLOBAL","channel":"WECHAT","ratePpm":3800,"state":"VALID","version":1,"credentialsConfigured":true,"merchantMasked":"******7890"},"effective":{"scope":"GLOBAL","channel":"WECHAT","ratePpm":3800,"state":"VALID","version":1,"credentialsConfigured":true,"merchantMasked":"******7890"}}]}`)}, nil
	}, DeliverSecret: func(agent.Context, ai.SecretPayload) error { return nil }})
	if err != nil {
		t.Fatal(err)
	}
	index := slices.IndexFunc(items, func(item tool.Tool) bool { return item.Name() == "HqPaymentConfigRead" })
	if index < 0 {
		t.Fatal("read tool missing")
	}
	result, err := items[index].(interface {
		Run(agent.Context, any) (map[string]any, error)
	}).Run(toolTestContext{}, map[string]any{"scope": "GLOBAL"})
	if err != nil || !called {
		t.Fatalf("read failed: %v called=%t", err, called)
	}
	encoded, _ := json.Marshal(result)
	if len(encoded) == 0 || string(encoded) == "{}" || result["validationOnly"] != true {
		t.Fatalf("missing status: %s", encoded)
	}
}

func TestPaymentConfigStateToolForwardsExactApprovedShape(t *testing.T) {
	var sent json.RawMessage
	items, err := Build(ai.FixedToolRuntime{Call: func(_ agent.Context, spec ai.ToolSpec, raw json.RawMessage) (ai.FixedResponse, error) {
		if spec.ID != "HqFranchisePaymentState" {
			t.Fatalf("wrong tool: %s", spec.ID)
		}
		sent = raw
		return ai.FixedResponse{Status: http.StatusOK, Body: []byte(`{"channels":[{"own":{"scope":"FRANCHISE","channel":"ALIPAY","state":"DISABLED","version":1,"ratePpm":0,"credentialsConfigured":false},"effective":{"scope":"FRANCHISE","channel":"ALIPAY","state":"DISABLED","version":1,"ratePpm":0,"credentialsConfigured":false}}]}`)}, nil
	}, DeliverSecret: func(agent.Context, ai.SecretPayload) error { return nil }})
	if err != nil {
		t.Fatal(err)
	}
	index := slices.IndexFunc(items, func(item tool.Tool) bool { return item.Name() == "HqFranchisePaymentState" })
	if index < 0 {
		t.Fatal("state tool missing")
	}
	result, err := items[index].(interface {
		Run(agent.Context, any) (map[string]any, error)
	}).Run(toolTestContext{}, map[string]any{
		"scope": "FRANCHISE", "organizationId": "org-a", "channel": "ALIPAY", "state": "DISABLED", "recordId": "", "version": 0,
	})
	if err != nil || result["channels"] == nil {
		t.Fatalf("state failed: %v %v", result, err)
	}
	var input map[string]any
	if err := json.Unmarshal(sent, &input); err != nil || input["organizationId"] != "org-a" || input["recordId"] != "" || input["scope"] != "FRANCHISE" {
		t.Fatalf("wrong body: %s %v", sent, err)
	}
}

func TestPaymentConfigReadRejectsMismatchedScopeBeforeProtectedCall(t *testing.T) {
	called := false
	items, err := Build(ai.FixedToolRuntime{Call: func(agent.Context, ai.ToolSpec, json.RawMessage) (ai.FixedResponse, error) {
		called = true
		return ai.FixedResponse{}, nil
	}, DeliverSecret: func(agent.Context, ai.SecretPayload) error { return nil }})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range items {
		if item.Name() != "HqPaymentConfigRead" {
			continue
		}
		_, err := item.(interface {
			Run(agent.Context, any) (map[string]any, error)
		}).Run(toolTestContext{}, map[string]any{"scope": "GLOBAL", "storeId": "store-a"})
		if err == nil || called {
			t.Fatalf("bad scope accepted: called=%t err=%v", called, err)
		}
		return
	}
	t.Fatal("read tool missing")
}
