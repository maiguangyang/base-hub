package ai

import (
	"encoding/json"
	"errors"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

func validateReadArguments(spec ToolSpec, arguments json.RawMessage) error {
	if spec.Mode != ModeReadOnly {
		return nil
	}
	var values map[string]json.RawMessage
	if err := json.Unmarshal(arguments, &values); err != nil {
		return errors.New("INVALID_READ_ARGUMENTS")
	}
	if values == nil {
		return errors.New("INVALID_READ_ARGUMENTS")
	}
	if err := validateReadVariableNames(spec.Document, values); err != nil {
		return err
	}
	if err := checkPageBound(values, "pageSize", 50); err != nil {
		return err
	}
	if err := checkPageBound(values, "perPage", 50); err != nil {
		return err
	}
	if err := checkPageBound(values, "page", 1000); err != nil {
		return err
	}
	return nil
}

func validateReadVariableNames(document string, values map[string]json.RawMessage) error {
	if document == "" {
		return nil
	}
	definitions, err := readVariableDefinitions(document)
	if err != nil {
		return err
	}
	for name := range values {
		if _, ok := definitions[name]; !ok {
			return errors.New("UNDECLARED_READ_VARIABLE")
		}
	}
	return requireReadVariables(definitions, values)
}

func readVariableDefinitions(document string) (map[string]*ast.VariableDefinition, error) {
	parsed, err := parser.ParseQuery(&ast.Source{Input: document})
	if err != nil || len(parsed.Operations) != 1 {
		return nil, errors.New("INVALID_READ_DOCUMENT")
	}
	definitions := make(map[string]*ast.VariableDefinition)
	for _, definition := range parsed.Operations[0].VariableDefinitions {
		definitions[definition.Variable] = definition
	}
	return definitions, nil
}

func requireReadVariables(definitions map[string]*ast.VariableDefinition, values map[string]json.RawMessage) error {
	for name, definition := range definitions {
		value, present := values[name]
		if definition.Type.NonNull && definition.DefaultValue == nil && (!present || string(value) == "null") {
			return errors.New("MISSING_READ_VARIABLE")
		}
	}
	return nil
}

func checkPageBound(values map[string]json.RawMessage, field string, max int) error {
	raw, exists := values[field]
	if !exists {
		return nil
	}
	var number int
	if json.Unmarshal(raw, &number) != nil || number < 1 || number > max {
		return errors.New("AI_PAGE_OUT_OF_RANGE")
	}
	return nil
}
