package authorization

import (
	"context"
	"reflect"
	"strings"

	"base-engine/gen"
)

// Product, store merchandising, and stock entities use dedicated business APIs.
func registerProductHandlers(handlers *gen.ResolutionHandlers) {
	value := reflect.ValueOf(handlers).Elem()
	typ := value.Type()
	for index := 0; index < value.NumField(); index++ {
		if !isProductHandler(typ.Field(index).Name) {
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

func isProductHandler(name string) bool {
	for _, fragment := range []string{
		"Product", "Specification", "StoreListing", "StorePackageOffer", "StorePriceRevision",
		"StoreInventoryBatch", "StoreStock", "StorePromotion", "StoreCouponTemplates",
		"Stocktakes", "StocktakeLines",
	} {
		if strings.Contains(name, fragment) {
			return true
		}
	}
	return false
}
