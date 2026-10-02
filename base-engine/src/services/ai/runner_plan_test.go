package ai

import (
	"encoding/json"
	"testing"

	"google.golang.org/genai"
)

func TestAIProposalToolDescribesStructuredParameters(t *testing.T) {
	item, err := (&runState{}).proposalTool()
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
	var schema map[string]any
	if err := json.Unmarshal(encoded, &schema); err != nil {
		t.Fatal(err)
	}
	steps := schema["properties"].(map[string]any)["steps"].(map[string]any)
	if steps["description"] == nil || steps["type"] != "array" {
		t.Fatalf("proposal steps schema = %s", encoded)
	}
	itemSchema := steps["items"].(map[string]any)
	for _, name := range []string{"operationId", "toolId", "arguments", "maxCalls", "sequence"} {
		field := itemSchema["properties"].(map[string]any)[name].(map[string]any)
		if field["description"] == nil || field["type"] == nil {
			t.Fatalf("proposal %s schema = %s", name, encoded)
		}
	}
}
