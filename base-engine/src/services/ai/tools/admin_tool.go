package tools

import (
	"encoding/json"
	"errors"

	"base-engine/src/services/ai"
	"google.golang.org/adk/v2/tool"
)

func buildReviewedTool(spec ai.ToolSpec, runtime ai.FixedToolRuntime) (tool.Tool, error) {
	inputSchema, err := reviewedInputSchema(spec)
	if err != nil {
		return nil, err
	}
	return ai.NewFixedGraphQLTool[map[string]any, graphqlData](spec, runtime, func(data graphqlData) (ai.SafeToolResult, error) {
		root := spec.OperationID[len("graphql.query."):]
		if spec.Mode == ai.ModeWrite {
			root = spec.OperationID[len("graphql.mutation."):]
		}
		var value any
		if err := json.Unmarshal(data[root], &value); err != nil {
			return ai.SafeToolResult{}, err
		}
		if value == nil {
			return ai.SafeToolResult{}, errors.New("EMPTY_TOOL_RESULT")
		}
		output := ai.SafeToolResult{ModelOutput: map[string]any{root: value}}
		if object, ok := value.(map[string]any); ok {
			if err := separateToolSecret(object, &output); err != nil {
				return ai.SafeToolResult{}, err
			}
			if id := resultResourceID(object); id != "" && spec.Mode == ai.ModeWrite {
				output.ModelOutput["id"] = id
			}
		}
		limitToolResult(value)
		return output, nil
	}, inputSchema)
}

func separateToolSecret(value map[string]any, output *ai.SafeToolResult) error {
	password, present := value["temporaryPassword"]
	if !present {
		return nil
	}
	delete(value, "temporaryPassword")
	secret, ok := password.(string)
	if password == nil {
		return nil
	}
	if !ok {
		return errors.New("INVALID_SECRET_RESULT")
	}
	if secret == "" {
		return errors.New("EMPTY_SECRET_RESULT")
	}
	accountID := resultAccountID(value)
	if accountID == "" {
		return errors.New("MISSING_SECRET_TARGET")
	}
	output.Secret = &ai.SecretPayload{TargetAccountID: accountID, Value: secret}
	return nil
}

func resultAccountID(value map[string]any) string {
	if id, ok := value["accountId"].(string); ok {
		return id
	}
	if membership, ok := value["membership"].(map[string]any); ok {
		if id, ok := membership["accountId"].(string); ok {
			return id
		}
	}
	return ""
}

func resultResourceID(value map[string]any) string {
	if id, ok := value["id"].(string); ok {
		return id
	}
	if membership, ok := value["membership"].(map[string]any); ok {
		if id, ok := membership["id"].(string); ok {
			return id
		}
	}
	return ""
}

func limitToolResult(value any) {
	object, ok := value.(map[string]any)
	if !ok {
		return
	}
	if data, ok := object["data"].([]any); ok && len(data) > 50 {
		object["data"] = data[:50]
	}
}
