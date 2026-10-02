package ai

import (
	"encoding/json"
	"testing"
)

func TestAIReadPageBounds(t *testing.T) {
	spec := ToolSpec{Mode: ModeReadOnly}
	for _, arguments := range []string{`{"page":0,"pageSize":1}`, `{"page":1,"pageSize":50000}`, `{"page":1,"perPage":51}`, `{"page":1001,"pageSize":1}`} {
		if err := validateReadArguments(spec, json.RawMessage(arguments)); err == nil {
			t.Fatalf("accepted %s", arguments)
		}
	}
	if err := validateReadArguments(spec, json.RawMessage(`{"page":1,"pageSize":50}`)); err != nil {
		t.Fatal(err)
	}
}

func TestAIReadRejectsUndeclaredAndMissingVariables(t *testing.T) {
	spec := ToolSpec{Mode: ModeReadOnly, Document: `query Stores($page:Int!){stores(current_page:$page,per_page:50){total}}`}
	for _, arguments := range []string{`{"page":1,"organizationId":"other"}`, `{}`} {
		if err := validateReadArguments(spec, json.RawMessage(arguments)); err == nil {
			t.Fatalf("accepted variables %s", arguments)
		}
	}
	if err := validateReadArguments(spec, json.RawMessage(`{"page":2}`)); err != nil {
		t.Fatal(err)
	}
}
