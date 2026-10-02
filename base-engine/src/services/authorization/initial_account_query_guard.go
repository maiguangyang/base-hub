package authorization

import (
	"context"
	"reflect"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"base-engine/auth"
)

// Generated filters and sorts can expose a protected relation without selecting its field.
func authorizeInitialAccountQueryArgs(ctx context.Context, principal *auth.WorkspacePrincipal) error {
	field := graphql.GetFieldContext(ctx)
	if field == nil || field.Object != "Query" {
		return nil
	}
	args := reflect.ValueOf(field.Args)
	if containsProtectedArgument(args, protectedOpeningRecordArgument) || containsProtectedArgument(args, protectedPaymentConfigArgument) ||
		!(principal.WorkspaceType == auth.WorkspaceTypeHeadquarters && principal.Has("account:update")) && containsProtectedArgument(args, protectedInitialAccountArgument) {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return nil
}

func containsProtectedArgument(root reflect.Value, protected func(string) bool) bool {
	stack := []reflect.Value{root}
	for len(stack) > 0 {
		value := unwrapArgument(stack[len(stack)-1])
		stack = stack[:len(stack)-1]
		if !value.IsValid() {
			continue
		}
		if containsProtectedField(value, protected) {
			return true
		}
		stack = appendArgumentChildren(stack, value)
	}
	return false
}

func unwrapArgument(value reflect.Value) reflect.Value {
	for value.IsValid() && (value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer) {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}
	return value
}

func containsProtectedField(value reflect.Value, protected func(string) bool) bool {
	switch value.Kind() {
	case reflect.Struct:
		for index := 0; index < value.NumField(); index++ {
			if protected(value.Type().Field(index).Name) && !value.Field(index).IsZero() {
				return true
			}
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			if key.Kind() == reflect.String && protected(key.String()) && !value.MapIndex(key).IsZero() {
				return true
			}
		}
	}
	return false
}

func appendArgumentChildren(stack []reflect.Value, value reflect.Value) []reflect.Value {
	switch value.Kind() {
	case reflect.Struct:
		for index := 0; index < value.NumField(); index++ {
			stack = append(stack, value.Field(index))
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			stack = append(stack, value.MapIndex(key))
		}
	case reflect.Slice, reflect.Array:
		for index := 0; index < value.Len(); index++ {
			stack = append(stack, value.Index(index))
		}
	}
	return stack
}

func protectedInitialAccountArgument(name string) bool {
	name = strings.ToLower(strings.ReplaceAll(name, "_", ""))
	return name == "initialaccount" || strings.HasPrefix(name, "initialaccountid") || strings.HasPrefix(name, "initializedorganizations")
}

func protectedOpeningRecordArgument(name string) bool {
	name = strings.ToLower(strings.ReplaceAll(name, "_", ""))
	return strings.HasPrefix(name, "openingrecords") || strings.HasPrefix(name, "recordedopeningrecords")
}

func protectedPaymentConfigArgument(name string) bool {
	name = strings.ToLower(strings.ReplaceAll(name, "_", ""))
	return strings.HasPrefix(name, "paymentconfigs")
}
