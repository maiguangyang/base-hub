package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

type SecretPayload struct{ ToolID, TargetAccountID, Value string }
type SafeToolResult struct {
	ModelOutput map[string]any
	Secret      *SecretPayload
}
type FixedToolRuntime struct {
	Call          func(agent.Context, ToolSpec, json.RawMessage) (FixedResponse, error)
	DeliverSecret func(agent.Context, SecretPayload) error
}

func NewFixedGraphQLTool[TIn any, TOut any](spec ToolSpec, runtime FixedToolRuntime, project func(TOut) (SafeToolResult, error), inputSchema ...*jsonschema.Schema) (tool.Tool, error) {
	if spec.ID == "" || strings.TrimSpace(spec.Description) == "" || !strings.HasPrefix(spec.OperationID, "graphql.") ||
		strings.TrimSpace(spec.Document) == "" || runtime.Call == nil || project == nil {
		return nil, errors.New("INVALID_TOOL_CONFIG")
	}
	if !validJSONType(reflect.TypeOf((*TOut)(nil)).Elem(), map[reflect.Type]bool{}) {
		return nil, errors.New("INVALID_OUTPUT_SCHEMA")
	}
	name, validName := toolFunctionName(spec)
	if !validName {
		return nil, errors.New("INVALID_TOOL_NAME")
	}
	config := functiontool.Config{Name: name, Description: spec.Description}
	if len(inputSchema) > 0 {
		config.InputSchema = inputSchema[0]
	}
	return functiontool.New(config,
		func(ctx agent.Context, input TIn) (map[string]any, error) {
			return executeFixedGraphQLTool(ctx, spec, runtime, input, project)
		})
}

func toolFunctionName(spec ToolSpec) (string, bool) {
	if spec.Name == "" {
		return spec.ID, true
	}
	return spec.Name, spec.Name == spec.ID
}

func validJSONType(t reflect.Type, seen map[reflect.Type]bool) bool {
	if seen[t] {
		return true
	}
	seen[t] = true
	switch t.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Array:
		return validJSONType(t.Elem(), seen)
	case reflect.Map:
		return t.Key().Kind() == reflect.String && validJSONType(t.Elem(), seen)
	case reflect.Struct:
		return validJSONStruct(t, seen)
	case reflect.Func, reflect.Chan, reflect.UnsafePointer, reflect.Complex64, reflect.Complex128:
		return false
	default:
		return true
	}
}

func validJSONStruct(t reflect.Type, seen map[reflect.Type]bool) bool {
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() && !validJSONType(t.Field(i).Type, seen) {
			return false
		}
	}
	return true
}

func executeFixedGraphQLTool[TIn any, TOut any](ctx agent.Context, spec ToolSpec, runtime FixedToolRuntime, input TIn, project func(TOut) (SafeToolResult, error)) (map[string]any, error) {
	variables, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	response, err := runtime.Call(ctx, spec, variables)
	if err != nil {
		return nil, err
	}
	if response.Status < http.StatusOK || response.Status >= http.StatusMultipleChoices {
		return nil, errors.New("PROTECTED_CALL_FAILED")
	}
	if err := validateGraphQLRoot(response.Body, spec.OperationID); err != nil {
		return nil, err
	}
	var envelope struct {
		Data   TOut              `json:"data"`
		Errors []json.RawMessage `json:"errors"`
	}
	if err := json.Unmarshal(response.Body, &envelope); err != nil {
		return nil, err
	}
	if len(envelope.Errors) > 0 {
		return nil, errors.New("GRAPHQL_OPERATION_FAILED")
	}
	safe, err := project(envelope.Data)
	if err != nil {
		return nil, err
	}
	if err := finishSafeToolResult(ctx, spec, safe, runtime); err != nil {
		return nil, err
	}
	return safe.ModelOutput, nil
}

func validateGraphQLRoot(body []byte, operationID string) error {
	var envelope struct {
		Data map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return err
	}
	root := operationID[strings.LastIndex(operationID, ".")+1:]
	value, ok := envelope.Data[root]
	if !ok || len(value) == 0 || string(value) == "null" {
		return errors.New("EMPTY_TOOL_RESULT")
	}
	return nil
}

func finishSafeToolResult(ctx agent.Context, spec ToolSpec, safe SafeToolResult, runtime FixedToolRuntime) error {
	if err := validateOutputFields(spec, safe.ModelOutput); err != nil {
		return err
	}
	if safe.Secret != nil {
		safe.Secret.ToolID = spec.ID
	}
	return deliverSafeResult(ctx, safe, runtime)
}

func validateOutputFields(spec ToolSpec, output map[string]any) error {
	if output == nil || len(spec.OutputFields) == 0 {
		return errors.New("MISSING_OUTPUT_PROJECTION")
	}
	allowed := make(map[string]struct{}, len(spec.OutputFields))
	for _, field := range spec.OutputFields {
		allowed[field] = struct{}{}
	}
	for field := range output {
		if _, ok := allowed[field]; !ok {
			return errors.New("UNAPPROVED_OUTPUT_FIELD")
		}
	}
	return nil
}

func deliverSafeResult(ctx agent.Context, result SafeToolResult, runtime FixedToolRuntime) error {
	if result.ModelOutput == nil {
		return errors.New("MISSING_MODEL_OUTPUT")
	}
	if result.Secret == nil {
		return nil
	}
	if result.Secret.Value == "" || result.Secret.TargetAccountID == "" || runtime.DeliverSecret == nil {
		return errors.New("INVALID_SECRET_DELIVERY")
	}
	modelJSON, err := json.Marshal(result.ModelOutput)
	if err != nil {
		return err
	}
	secretJSON, err := json.Marshal(result.Secret.Value)
	if err != nil {
		return err
	}
	if bytes.Contains(modelJSON, secretJSON) || bytes.Contains(modelJSON, []byte(result.Secret.Value)) {
		return errors.New("SECRET_IN_MODEL_OUTPUT")
	}
	return runtime.DeliverSecret(ctx, *result.Secret)
}
