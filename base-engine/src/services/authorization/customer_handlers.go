package authorization

import (
	"context"
	"reflect"
	"strings"

	"base-engine/gen"
)

// registerCustomerHandlers denies every generated customer data entry,
// including relationship fields and relationship-ID projections.
func registerCustomerHandlers(handlers *gen.ResolutionHandlers) {
	value := reflect.ValueOf(handlers).Elem()
	typ := value.Type()
	for index := 0; index < value.NumField(); index++ {
		if !strings.Contains(typ.Field(index).Name, "Customer") {
			continue
		}
		field := value.Field(index)
		signature := field.Type()
		field.Set(reflect.MakeFunc(signature, func(args []reflect.Value) []reflect.Value {
			ctx := args[0].Interface().(context.Context)
			results := make([]reflect.Value, signature.NumOut())
			for output := range results {
				results[output] = reflect.Zero(signature.Out(output))
			}
			results[len(results)-1] = reflect.ValueOf(generatedDataAccessDenied(ctx))
			return results
		}))
	}
}
