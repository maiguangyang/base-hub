package ai

import (
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"base-engine/gen"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestContractInventoryGraphQLCoverageAndDrift(t *testing.T) {
	schema := gen.NewExecutableSchema(gen.Config{}).Schema()
	if schema.Query.Fields.ForName("accounts") == nil ||
		schema.Mutation.Fields.ForName("createAccount") == nil {
		t.Fatal("account generated root operations missing")
	}
	records, err := DiscoverContracts(schema, mux.NewRouter())
	expected := len(schema.Query.Fields) + len(schema.Mutation.Fields) + len(schema.Subscription.Fields)
	if err != nil || expected == 0 || len(records) != expected {
		t.Fatalf("discovered records = %d, expected %d: %v", len(records), expected, err)
	}
	if err := CompareContracts(schema, mux.NewRouter(), records); err != nil {
		t.Fatal(err)
	}
	mutated := *schema
	mutatedTypes := make(map[string]*ast.Definition, len(schema.Types))
	for key, value := range schema.Types {
		mutatedTypes[key] = value
	}
	mutated.Types = mutatedTypes
	field := *schema.Query.Fields[0]
	field.Type = ast.NamedType("String", nil)
	query := *schema.Query
	query.Fields = append(ast.FieldList(nil), schema.Query.Fields...)
	query.Fields[0] = &field
	mutated.Query = &query
	mutated.Types[query.Name] = &query
	if err := CompareContracts(&mutated, mux.NewRouter(), records); err == nil || !strings.Contains(err.Error(), "graphql.query."+field.Name) {
		t.Fatalf("field drift error = %v", err)
	}
}

func TestContractInventoryHTTPRouteDrift(t *testing.T) {
	schema := gen.NewExecutableSchema(gen.Config{}).Schema()
	router := mux.NewRouter()
	router.HandleFunc("/api/example", nil).Methods("POST").Name(HTTPRouteName("http.example", struct {
		Value string `json:"value"`
	}{}, struct{}{}))
	records, err := DiscoverContracts(schema, router)
	if err != nil {
		t.Fatal(err)
	}
	if err := CompareContracts(schema, router, records); err != nil {
		t.Fatal(err)
	}
	changed := mux.NewRouter()
	changed.HandleFunc("/api/example", nil).Methods("PUT").Name(HTTPRouteName("http.example", struct {
		Value string `json:"value"`
	}{}, struct{}{}))
	if err := CompareContracts(schema, changed, records); err == nil || !strings.Contains(err.Error(), "http.example") {
		t.Fatalf("route drift error = %v", err)
	}
}

func TestContractInventoryNestedEnumAndDirectiveDrift(t *testing.T) {
	schema := gen.NewExecutableSchema(gen.Config{}).Schema()
	storeField := schema.Query.Fields.ForName("stores")
	if storeField == nil {
		t.Fatal("stores field missing")
	}
	baseline := graphQLFingerprint(schema, storeField)
	changed := *schema
	changed.Types = cloneTypeMap(schema.Types)
	store := *schema.Types["Store"]
	store.Fields = append(ast.FieldList(nil), store.Fields...)
	name := *store.Fields.ForName("name")
	name.Type = ast.NamedType("Int", nil)
	for i, field := range store.Fields {
		if field.Name == "name" {
			store.Fields[i] = &name
		}
	}
	changed.Types["Store"] = &store
	if graphQLFingerprint(&changed, storeField) == baseline {
		t.Fatal("nested output field drift missed")
	}
	changed.Types = cloneTypeMap(schema.Types)
	status := *schema.Types["StoreBusinessStatus"]
	status.EnumValues = append(ast.EnumValueList(nil), status.EnumValues...)
	status.EnumValues = append(status.EnumValues, &ast.EnumValueDefinition{Name: "FUTURE_STATUS"})
	changed.Types[status.Name] = &status
	if graphQLFingerprint(&changed, storeField) == baseline {
		t.Fatal("enum drift missed")
	}
	permissionField := schema.Mutation.Fields.ForName("suspendOrganization")
	if permissionField == nil {
		t.Fatal("permission field missing")
	}
	baseline = graphQLFingerprint(schema, permissionField)
	permission := *permissionField
	permission.Directives = append(ast.DirectiveList(nil), permission.Directives...)
	directive := *permission.Directives.ForName("hasPermission")
	directive.Arguments = append(ast.ArgumentList(nil), directive.Arguments...)
	argument := *directive.Arguments.ForName("action")
	value := *argument.Value
	value.Raw = "other:action"
	argument.Value = &value
	directive.Arguments[0] = &argument
	permission.Directives[0] = &directive
	if graphQLFingerprint(schema, &permission) == baseline {
		t.Fatal("permission directive drift missed")
	}
}

func TestContractInventoryDTODrift(t *testing.T) {
	old := HTTPRouteName("http.example", struct {
		Nested struct {
			Name string `json:"name"`
		} `json:"nested"`
	}{}, struct{}{})
	changed := HTTPRouteName("http.example", struct {
		Nested struct {
			Name int `json:"name"`
		} `json:"nested"`
	}{}, struct{}{})
	if old == changed {
		t.Fatal("nested HTTP DTO drift missed")
	}
}

func cloneTypeMap(original map[string]*ast.Definition) map[string]*ast.Definition {
	copy := make(map[string]*ast.Definition, len(original))
	for name, definition := range original {
		copy[name] = definition
	}
	return copy
}
