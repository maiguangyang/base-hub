package ai

import (
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

func graphQLContracts(schema *ast.Schema) []ContractRecord {
	var records []ContractRecord
	for _, root := range []struct {
		name, protocol string
		definition     *ast.Definition
	}{{"query", "GRAPHQL_QUERY", schema.Query}, {"mutation", "GRAPHQL_MUTATION", schema.Mutation}, {"subscription", "GRAPHQL_SUBSCRIPTION", schema.Subscription}} {
		if root.definition == nil {
			continue
		}
		for _, field := range root.definition.Fields {
			records = append(records, ContractRecord{
				OperationID:     "graphql." + root.name + "." + field.Name,
				Protocol:        root.protocol,
				Fingerprint:     graphQLFingerprint(schema, field),
				Classifications: []string{"OTHER_EXTERNAL_API"}, Availability: "NON_CALLABLE",
				NonCallableReason: "NOT_REVIEWED_FOR_AI",
			})
		}
	}
	return records
}

func graphQLFingerprint(schema *ast.Schema, field *ast.FieldDefinition) string {
	collector := contractTypeCollector{schema: schema, seenTypes: map[string]bool{}, seenDirectives: map[string]bool{}, parts: []string{fieldSignature(field)}}
	collector.visitType(field.Type.Name())
	for _, arg := range field.Arguments {
		collector.visitType(arg.Type.Name())
		collector.visitDirectives(arg.Directives)
	}
	collector.visitDirectives(field.Directives)
	sort.Strings(collector.parts)
	return ContractHash(collector.parts...)
}

type contractTypeCollector struct {
	schema                    *ast.Schema
	seenTypes, seenDirectives map[string]bool
	parts                     []string
}

func (c *contractTypeCollector) visitDirectives(directives ast.DirectiveList) {
	for _, directive := range directives {
		if c.seenDirectives[directive.Name] {
			continue
		}
		c.seenDirectives[directive.Name] = true
		definition := c.schema.Directives[directive.Name]
		if definition == nil {
			continue
		}
		c.parts = append(c.parts, directiveDefinitionSignature(definition))
		for _, arg := range definition.Arguments {
			c.visitType(arg.Type.Name())
		}
	}
}

func (c *contractTypeCollector) visitType(name string) {
	if c.seenTypes[name] {
		return
	}
	c.seenTypes[name] = true
	definition := c.schema.Types[name]
	if definition == nil {
		return
	}
	c.parts = append(c.parts, typeSignature(definition))
	c.visitDirectives(definition.Directives)
	for _, field := range definition.Fields {
		c.visitType(field.Type.Name())
		c.visitDirectives(field.Directives)
		for _, arg := range field.Arguments {
			c.visitType(arg.Type.Name())
			c.visitDirectives(arg.Directives)
		}
	}
	for _, member := range definition.Types {
		c.visitType(member)
	}
	for _, member := range definition.Interfaces {
		c.visitType(member)
	}
}

func typeSignature(definition *ast.Definition) string {
	parts := []string{string(definition.Kind), definition.Name, directiveSignature(definition.Directives)}
	for _, f := range definition.Fields {
		parts = append(parts, fieldSignature(f))
	}
	for _, v := range definition.EnumValues {
		parts = append(parts, "enum:"+v.Name+directiveSignature(v.Directives))
	}
	for _, member := range definition.Types {
		parts = append(parts, "member:"+member)
	}
	for _, iface := range definition.Interfaces {
		parts = append(parts, "interface:"+iface)
	}
	sort.Strings(parts[3:])
	return strings.Join(parts, "|")
}

func fieldSignature(field *ast.FieldDefinition) string {
	parts := []string{"field:" + field.Name, field.Type.String(), directiveSignature(field.Directives)}
	if field.DefaultValue != nil {
		parts = append(parts, "default:"+field.DefaultValue.String())
	}
	for _, arg := range field.Arguments {
		parts = append(parts, argumentSignature(arg))
	}
	sort.Strings(parts[3:])
	return strings.Join(parts, "|")
}

func argumentSignature(arg *ast.ArgumentDefinition) string {
	part := "arg:" + arg.Name + ":" + arg.Type.String() + directiveSignature(arg.Directives)
	if arg.DefaultValue != nil {
		part += "=" + arg.DefaultValue.String()
	}
	return part
}

func directiveSignature(directives ast.DirectiveList) string {
	parts := make([]string, 0, len(directives))
	for _, directive := range directives {
		args := make([]string, 0, len(directive.Arguments))
		for _, arg := range directive.Arguments {
			args = append(args, arg.Name+"="+arg.Value.String())
		}
		sort.Strings(args)
		parts = append(parts, "@"+directive.Name+"("+strings.Join(args, ",")+")")
	}
	sort.Strings(parts)
	return strings.Join(parts, "")
}

func directiveDefinitionSignature(definition *ast.DirectiveDefinition) string {
	parts := []string{"directive:" + definition.Name}
	for _, arg := range definition.Arguments {
		parts = append(parts, argumentSignature(arg))
	}
	for _, location := range definition.Locations {
		parts = append(parts, "on:"+string(location))
	}
	sort.Strings(parts[1:])
	return strings.Join(parts, "|")
}
