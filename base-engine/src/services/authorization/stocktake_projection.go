package authorization

import (
	"context"
	"reflect"
	"strings"

	"base-engine/gen"
)

// Generated filters and sorts can join through Account, Store, and other
// accessible parents into blind or out-of-scope stocktake data. Reject such
// query shapes at every generated root before any SQL is built.
func rejectGeneratedStocktakeProjections(handlers *gen.ResolutionHandlers) {
	value := reflect.ValueOf(handlers).Elem()
	for index := 0; index < value.NumField(); index++ {
		if !strings.HasPrefix(value.Type().Field(index).Name, "Query") {
			continue
		}
		field := value.Field(index)
		field.Set(guardStocktakeQuery(field))
	}
}

func guardStocktakeQuery(field reflect.Value) reflect.Value {
	original, signature := reflect.ValueOf(field.Interface()), field.Type()
	return reflect.MakeFunc(signature, func(args []reflect.Value) []reflect.Value {
		for _, argument := range args[2:] {
			if containsStocktakeProjection(argument, map[uintptr]bool{}) {
				return deniedStocktakeQuery(args[0].Interface().(context.Context), signature)
			}
		}
		return original.Call(args)
	})
}

func deniedStocktakeQuery(ctx context.Context, signature reflect.Type) []reflect.Value {
	results := make([]reflect.Value, signature.NumOut())
	for output := range results {
		results[output] = reflect.Zero(signature.Out(output))
	}
	results[len(results)-1] = reflect.ValueOf(generatedDataAccessDenied(ctx))
	return results
}

func containsStocktakeProjection(value reflect.Value, seen map[uintptr]bool) bool {
	if !value.IsValid() {
		return false
	}
	switch value.Kind() {
	case reflect.Interface, reflect.Pointer:
		return stocktakeProjectionPointer(value, seen)
	case reflect.Slice, reflect.Array:
		return stocktakeProjectionItems(value, seen)
	case reflect.Struct:
		return stocktakeProjectionFields(value, seen)
	}
	return false
}

func stocktakeProjectionPointer(value reflect.Value, seen map[uintptr]bool) bool {
	if value.IsNil() {
		return false
	}
	if value.Kind() == reflect.Pointer {
		address := value.Pointer()
		if seen[address] {
			return false
		}
		seen[address] = true
	}
	return containsStocktakeProjection(value.Elem(), seen)
}

func stocktakeProjectionItems(value reflect.Value, seen map[uintptr]bool) bool {
	for index := 0; index < value.Len(); index++ {
		if containsStocktakeProjection(value.Index(index), seen) {
			return true
		}
	}
	return false
}

func stocktakeProjectionFields(value reflect.Value, seen map[uintptr]bool) bool {
	if value.Type().PkgPath() != reflect.TypeOf(gen.AccountFilterType{}).PkgPath() {
		return false
	}
	for index := 0; index < value.NumField(); index++ {
		field := value.Field(index)
		if strings.Contains(value.Type().Field(index).Name, "Stocktake") && !field.IsZero() {
			return true
		}
		if containsStocktakeProjection(field, seen) {
			return true
		}
	}
	return false
}
