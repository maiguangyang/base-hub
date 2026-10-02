package tools

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"base-engine/gen"
	"base-engine/src/services/ai"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

var aiExcludedInputFields = map[string]struct{}{
	"isDelete": {}, "weight": {}, "state": {},
	"organization": {}, "members": {}, "permissions": {}, "auditLogs": {},
	"reviewedByAccount": {}, "reviewedByAccountId": {},
}

var aiFilterFields = map[string]map[string]struct{}{
	"StoreFilterType":                selectedInputFields("id", "code", "name", "lifecycle", "organizationId"),
	"OperatorMembershipFilterType":   selectedInputFields("id", "status", "accountId", "organizationId", "storeAccessMode"),
	"OperatorRoleFilterType":         selectedInputFields("id", "name", "kind", "organizationId"),
	"AuditLogFilterType":             selectedInputFields("action", "resourceType", "resourceId", "resultCode", "actorAccountId", "organizationId", "storeId", "createdAt"),
	"MembershipInvitationFilterType": selectedInputFields("id", "membershipId"),
}

var aiWriteInputFields = map[string]map[string]struct{}{
	"CreateOperatorRoleInput":       selectedInputFields("name", "kind", "organizationId", "permissionsIds"),
	"UpdateOperatorRoleInput":       selectedInputFields("name", "permissionsIds"),
	"UpdateOperatorMembershipInput": selectedInputFields("storeAccessMode", "rolesIds", "storesIds"),
	"CreateStoreInput":              selectedInputFields("code", "name", "lifecycle", "organizationId", "contactPhone", "managerName", "managerPhone", "province", "city", "district", "address", "businessHours", "businessStatus", "supportDineIn", "supportTakeout", "storeArea", "tableCount", "receiptFooter"),
	"UpdateStoreInput":              selectedInputFields("code", "name", "contactPhone", "managerName", "managerPhone", "province", "city", "district", "address", "businessHours", "businessStatus", "supportDineIn", "supportTakeout", "storeArea", "tableCount", "receiptFooter"),
}

func selectedInputFields(names ...string) map[string]struct{} {
	result := make(map[string]struct{}, len(names))
	for _, name := range names {
		result[name] = struct{}{}
	}
	return result
}

func reviewedInputSchema(spec ai.ToolSpec) (*jsonschema.Schema, error) {
	query, err := parser.ParseQuery(&ast.Source{Input: spec.Document})
	if err != nil || len(query.Operations) != 1 {
		return nil, errors.New("INVALID_TOOL_INPUT_DOCUMENT")
	}
	definitions := gen.NewExecutableSchema(gen.Config{}).Schema().Types
	result := strictObjectSchema()
	for _, variable := range query.Operations[0].VariableDefinitions {
		field, err := schemaForGraphQLType(variable.Type, definitions, map[string]bool{})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", spec.ID, err)
		}
		applyReviewedEnum(spec.ID, variable.Variable, field)
		if err := constrainMemberPhone(spec.ID, variable.Variable, field); err != nil {
			return nil, err
		}
		result.Properties[variable.Variable] = field
		describeReviewedField(spec, variable.Variable, field)
		if variable.Type.NonNull && variable.DefaultValue == nil {
			result.Required = append(result.Required, variable.Variable)
		}
		applyToolPageBound(variable.Variable, field)
	}
	if err := removeTrustedSpecFields(spec, result); err != nil {
		return nil, err
	}
	return result, nil
}

func describeReviewedField(spec ai.ToolSpec, name string, field *jsonschema.Schema) {
	describeInputSchema(name, field, spec.Title)
}

func constrainMemberPhone(specID, variable string, field *jsonschema.Schema) error {
	return nil
}

func applyReviewedEnum(specID, variable string, field *jsonschema.Schema) {
}

func removeTrustedSpecFields(spec ai.ToolSpec, schema *jsonschema.Schema) error {
	for _, path := range []string{spec.RequestKeyPath, spec.GeneratedCodePath} {
		if path != "" && !removeTrustedInputField(schema, path) {
			return fmt.Errorf("%s: unknown trusted field %s", spec.ID, path)
		}
	}
	return nil
}

func removeTrustedInputField(schema *jsonschema.Schema, path string) bool {
	parts := strings.Split(path, ".")
	for _, part := range parts[:len(parts)-1] {
		schema = schema.Properties[part]
		if schema == nil {
			return false
		}
	}
	name := parts[len(parts)-1]
	if _, exists := schema.Properties[name]; !exists {
		return false
	}
	delete(schema.Properties, name)
	schema.Required = slices.DeleteFunc(schema.Required, func(field string) bool { return field == name })
	return true
}

func schemaForGraphQLType(typ *ast.Type, definitions map[string]*ast.Definition, seen map[string]bool) (*jsonschema.Schema, error) {
	if typ == nil {
		return nil, errors.New("MISSING_GRAPHQL_TYPE")
	}
	if typ.Elem != nil {
		item, err := schemaForGraphQLType(typ.Elem, definitions, seen)
		return &jsonschema.Schema{Type: "array", Items: item}, err
	}
	if scalar := graphqlScalarSchema(typ.NamedType); scalar != nil {
		return scalar, nil
	}
	definition := definitions[typ.NamedType]
	if definition == nil {
		return nil, fmt.Errorf("UNKNOWN_GRAPHQL_TYPE_%s", typ.NamedType)
	}
	if definition.Kind == ast.Enum {
		return enumInputSchema(definition), nil
	}
	if definition.Kind != ast.InputObject || seen[definition.Name] {
		return nil, fmt.Errorf("INVALID_GRAPHQL_INPUT_%s", typ.NamedType)
	}
	seen[definition.Name] = true
	defer delete(seen, definition.Name)
	return inputObjectSchema(definition, definitions, seen)
}

func graphqlScalarSchema(name string) *jsonschema.Schema {
	switch name {
	case "ID", "String", "Time", "DateTime":
		return &jsonschema.Schema{Type: "string"}
	case "Int":
		return &jsonschema.Schema{Type: "integer"}
	case "Float":
		return &jsonschema.Schema{Type: "number"}
	case "Boolean":
		return &jsonschema.Schema{Type: "boolean"}
	default:
		return nil
	}
}

func inputObjectSchema(definition *ast.Definition, definitions map[string]*ast.Definition, seen map[string]bool) (*jsonschema.Schema, error) {
	result := strictObjectSchema()
	for _, field := range definition.Fields {
		if !allowedInputField(definition.Name, field.Name) {
			continue
		}
		if seen[field.Type.Name()] {
			continue
		}
		child, err := schemaForGraphQLType(field.Type, definitions, seen)
		if err != nil {
			return nil, err
		}
		result.Properties[field.Name] = child
		describeInputSchema(field.Name, child, definition.Name)
		if field.Type.NonNull && field.DefaultValue == nil {
			result.Required = append(result.Required, field.Name)
		}
	}
	return result, nil
}

func allowedInputField(inputType, field string) bool {
	if _, excluded := aiExcludedInputFields[field]; excluded {
		return false
	}
	if allowed, scoped := aiWriteInputFields[inputType]; scoped {
		_, ok := allowed[field]
		return ok
	}
	if allowed, scoped := aiFilterFields[inputType]; scoped {
		_, ok := allowed[field]
		return ok
	}
	return true
}

func enumInputSchema(definition *ast.Definition) *jsonschema.Schema {
	result := &jsonschema.Schema{Type: "string"}
	for _, value := range definition.EnumValues {
		result.Enum = append(result.Enum, value.Name)
	}
	return result
}

func strictObjectSchema() *jsonschema.Schema {
	return &jsonschema.Schema{Type: "object", Properties: map[string]*jsonschema.Schema{}, AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}}}
}

func applyToolPageBound(name string, schema *jsonschema.Schema) {
	minimum := float64(1)
	switch name {
	case "page":
		maximum := float64(1000)
		schema.Minimum, schema.Maximum = &minimum, &maximum
	case "pageSize", "perPage":
		maximum := float64(50)
		schema.Minimum, schema.Maximum = &minimum, &maximum
	}
}
