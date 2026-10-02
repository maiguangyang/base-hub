package ai

import (
	"errors"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

func validateToolDocument(spec ToolSpec, contract ContractRecord) error {
	document, err := parser.ParseQuery(&ast.Source{Input: spec.Document})
	if err != nil || len(document.Operations) != 1 {
		return errors.New("invalid fixed GraphQL document")
	}
	operation := document.Operations[0]
	if !documentOperationMatches(operation.Operation, contract.Protocol) || len(operation.SelectionSet) != 1 {
		return errors.New("fixed document operation mismatch")
	}
	field, ok := operation.SelectionSet[0].(*ast.Field)
	if !ok || field.Name != strings.TrimPrefix(spec.OperationID, operationIDPrefix(contract.Protocol)) {
		return errors.New("fixed document field mismatch")
	}
	return nil
}

func documentOperationMatches(operation ast.Operation, protocol string) bool {
	switch protocol {
	case "GRAPHQL_QUERY":
		return operation == ast.Query
	case "GRAPHQL_MUTATION":
		return operation == ast.Mutation
	case "GRAPHQL_SUBSCRIPTION":
		return operation == ast.Subscription
	default:
		return false
	}
}

func operationIDPrefix(protocol string) string {
	switch protocol {
	case "GRAPHQL_QUERY":
		return "graphql.query."
	case "GRAPHQL_MUTATION":
		return "graphql.mutation."
	case "GRAPHQL_SUBSCRIPTION":
		return "graphql.subscription."
	default:
		return ""
	}
}
