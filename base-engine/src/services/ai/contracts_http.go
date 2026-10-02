package ai

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"github.com/vektah/gqlparser/v2/ast"
)

const routeNamePrefix = "contract:"

func HTTPRouteName(operationID string, request, response any) string {
	return routeNamePrefix + operationID + ":" + ContractHash(dtoSignature(reflect.TypeOf(request)), dtoSignature(reflect.TypeOf(response)))
}

func httpContracts(schema *ast.Schema, router *mux.Router) ([]ContractRecord, error) {
	var records []ContractRecord
	err := router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return err
		}
		for _, method := range methods {
			id, dtoHash, err := routeIdentity(route.GetName(), method, path)
			if err != nil {
				return err
			}
			fingerprint := ContractHash(method, path, dtoHash)
			if path == "/graphql" {
				fingerprint = ContractHash(method, path, schemaFingerprint(schema))
			}
			records = append(records, ContractRecord{OperationID: id, Protocol: "HTTP_ROUTE", Method: method,
				Path: path, Fingerprint: fingerprint, Classifications: []string{"OTHER_EXTERNAL_API"},
				Availability: "NON_CALLABLE", NonCallableReason: "NOT_REVIEWED_FOR_AI"})
		}
		return nil
	})
	return records, err
}

func routeIdentity(name, method, path string) (string, string, error) {
	if strings.HasPrefix(name, routeNamePrefix) {
		parts := strings.SplitN(strings.TrimPrefix(name, routeNamePrefix), ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return "", "", fmt.Errorf("invalid route metadata %q", name)
		}
		return parts[0], parts[1], nil
	}
	if name != "" {
		return "", "", fmt.Errorf("unrecognized route metadata %q", name)
	}
	return "http." + method + "." + path, "transport", nil
}

func schemaFingerprint(schema *ast.Schema) string {
	parts := make([]string, 0, 80)
	for _, record := range graphQLContracts(schema) {
		parts = append(parts, record.OperationID+":"+record.Fingerprint)
	}
	sort.Strings(parts)
	return ContractHash(parts...)
}

func dtoSignature(t reflect.Type) string {
	if t == nil {
		return "nil"
	}
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Struct:
		parts := []string{"struct:" + t.PkgPath() + "." + t.Name()}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			parts = append(parts, field.Name+":"+string(field.Tag)+":"+dtoSignature(field.Type))
		}
		sort.Strings(parts[1:])
		return strings.Join(parts, "|")
	case reflect.Slice, reflect.Array:
		return "[]" + dtoSignature(t.Elem())
	case reflect.Map:
		return "map[" + dtoSignature(t.Key()) + "]" + dtoSignature(t.Elem())
	default:
		return t.String()
	}
}
