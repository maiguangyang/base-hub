package tools

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/jsonschema-go/jsonschema"
	"base-engine/src/services/ai"
	"base-engine/src/services/paymentconfig"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

type paymentReadInput struct {
	Scope          string `json:"scope"`
	OrganizationID string `json:"organizationId,omitempty"`
	StoreID        string `json:"storeId,omitempty"`
}
type paymentStateInput struct {
	Scope          string `json:"scope"`
	OrganizationID string `json:"organizationId,omitempty"`
	StoreID        string `json:"storeId,omitempty"`
	Channel        string `json:"channel"`
	State          string `json:"state"`
	RecordID       string `json:"recordId"`
	Version        uint64 `json:"version"`
}
type paymentRestoreInput struct {
	Scope          string `json:"scope"`
	OrganizationID string `json:"organizationId,omitempty"`
	StoreID        string `json:"storeId,omitempty"`
	Channel        string `json:"channel"`
	RecordID       string `json:"recordId"`
	Version        uint64 `json:"version"`
}

func paymentInputSchema(spec ai.ToolSpec) *jsonschema.Schema {
	scope := paymentToolScope(spec.ID)
	if scope == "" {
		return nil
	}
	result := strictObjectSchema()
	addPaymentField(result, "scope", &jsonschema.Schema{Type: "string", Enum: scopeValues(scope)}, true, spec.Title)
	if scope == "ANY" {
		addPaymentField(result, "organizationId", &jsonschema.Schema{Type: "string"}, false, spec.Title)
		addPaymentField(result, "storeId", &jsonschema.Schema{Type: "string"}, false, spec.Title)
		return result
	}
	if scope == "FRANCHISE" {
		addPaymentField(result, "organizationId", &jsonschema.Schema{Type: "string"}, true, spec.Title)
	}
	if scope == "STORE" {
		addPaymentField(result, "storeId", &jsonschema.Schema{Type: "string"}, true, spec.Title)
	}
	addPaymentField(result, "channel", &jsonschema.Schema{Type: "string", Enum: []any{"WECHAT", "ALIPAY"}}, true, spec.Title)
	minimum := float64(0)
	record := &jsonschema.Schema{Type: "string"}
	if spec.OperationID == "http.paymentConfig.restoreInheritance" {
		minimum = 1
		minLength := 1
		record.MinLength = &minLength
	}
	addPaymentField(result, "recordId", record, true, spec.Title)
	addPaymentField(result, "version", &jsonschema.Schema{Type: "integer", Minimum: &minimum}, true, spec.Title)
	if spec.OperationID == "http.paymentConfig.restoreInheritance" {
		result.Properties["recordId"].Description = "待恢复继承的支付配置记录 ID；必须来自当前范围与渠道的现有配置。"
		result.Properties["version"].Description = "待恢复继承的现有支付配置并发版本；必须大于 0。"
	}
	if spec.OperationID == "http.paymentConfig.state" {
		addPaymentField(result, "state", &jsonschema.Schema{Type: "string", Enum: []any{"VALID", "DISABLED"}}, true, spec.Title)
	}
	return result
}

func paymentToolScope(id string) string {
	switch id {
	case "HqPaymentConfigRead":
		return "ANY"
	case "HqGlobalPaymentState":
		return "GLOBAL"
	case "HqFranchisePaymentState", "HqFranchisePaymentRestore":
		return "FRANCHISE"
	case "HqStorePaymentState", "HqStorePaymentRestore":
		return "STORE"
	default:
		return ""
	}
}

func scopeValues(scope string) []any {
	if scope == "ANY" {
		return []any{"GLOBAL", "FRANCHISE", "STORE"}
	}
	return []any{scope}
}

func addPaymentField(schema *jsonschema.Schema, name string, field *jsonschema.Schema, required bool, title string) {
	describeInputSchema(name, field, title)
	schema.Properties[name] = field
	if required {
		schema.Required = append(schema.Required, name)
	}
}

func buildPaymentTool(spec ai.ToolSpec, runtime ai.FixedToolRuntime) (tool.Tool, error) {
	schema := paymentInputSchema(spec)
	if schema == nil {
		return nil, errors.New("UNREVIEWED_HTTP_TOOL")
	}
	switch spec.OperationID {
	case "http.paymentConfig.read":
		return newPaymentHTTPTool[paymentReadInput](spec, schema, runtime, validPaymentRead)
	case "http.paymentConfig.state":
		return newPaymentHTTPTool[paymentStateInput](spec, schema, runtime, func(input paymentStateInput) bool {
			return validPaymentScope(input.Scope, input.OrganizationID, input.StoreID)
		})
	case "http.paymentConfig.restoreInheritance":
		return newPaymentHTTPTool[paymentRestoreInput](spec, schema, runtime, func(input paymentRestoreInput) bool {
			return validPaymentScope(input.Scope, input.OrganizationID, input.StoreID)
		})
	default:
		return nil, errors.New("UNREVIEWED_HTTP_TOOL")
	}
}

func validPaymentRead(input paymentReadInput) bool {
	return validPaymentScope(input.Scope, input.OrganizationID, input.StoreID)
}
func validPaymentScope(scope, organizationID, storeID string) bool {
	switch scope {
	case "GLOBAL":
		return organizationID == "" && storeID == ""
	case "FRANCHISE":
		return organizationID != "" && storeID == ""
	case "STORE":
		return storeID != "" && organizationID == ""
	default:
		return false
	}
}

func newPaymentHTTPTool[T any](spec ai.ToolSpec, schema *jsonschema.Schema, runtime ai.FixedToolRuntime, valid func(T) bool) (tool.Tool, error) {
	return functiontool.New(functiontool.Config{Name: spec.Name, Description: spec.Description, InputSchema: schema},
		func(ctx agent.Context, input T) (map[string]any, error) {
			if !valid(input) {
				return nil, errors.New("INVALID_PAYMENT_SCOPE")
			}
			variables, err := json.Marshal(input)
			if err != nil {
				return nil, err
			}
			response, err := runtime.Call(ctx, spec, variables)
			if err != nil {
				return nil, err
			}
			if response.Status != http.StatusOK {
				return nil, paymentProtectedError(response)
			}
			var view paymentconfig.ScopeView
			if err := json.Unmarshal(response.Body, &view); err != nil || len(view.Channels) == 0 {
				return nil, errors.New("INVALID_PROTECTED_RESULT")
			}
			return map[string]any{"channels": view.Channels, "validationOnly": true}, nil
		})
}

func paymentProtectedError(response ai.FixedResponse) error {
	var body struct {
		Code string `json:"code"`
	}
	if json.Unmarshal(response.Body, &body) == nil {
		switch {
		case response.Status == http.StatusBadRequest && body.Code == "VALIDATION_FAILED",
			response.Status == http.StatusForbidden && body.Code == "PERMISSION_DENIED",
			response.Status == http.StatusConflict && body.Code == "CONFLICT",
			response.Status == http.StatusServiceUnavailable && body.Code == "CONFIG_UNAVAILABLE":
			return errors.New(body.Code)
		}
	}
	return errors.New("PROTECTED_CALL_FAILED")
}
