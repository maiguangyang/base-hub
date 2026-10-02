package authorization

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestStocktakeGeneratedHandlersDenyEveryEntry(t *testing.T) {
	handlers := RegisterHandlers(gen.DefaultResolutionHandlers())
	value := reflect.ValueOf(handlers)
	typ := value.Type()
	ctx := auth.WithPrincipal(context.Background(), &auth.WorkspacePrincipal{AccountID: "operator",
		Permissions: map[string]struct{}{"franchiseStocktake:read": {}, "franchiseStocktake:record": {}, "franchiseStocktake:post": {}}})
	covered := 0
	for index := 0; index < value.NumField(); index++ {
		name := typ.Field(index).Name
		if !strings.Contains(name, "Stocktake") {
			continue
		}
		covered++
		if !isProductHandler(name) {
			t.Errorf("generated stocktake entry not covered: %s", name)
			continue
		}
		fn := value.Field(index)
		args := make([]reflect.Value, fn.Type().NumIn())
		args[0] = reflect.ValueOf(ctx)
		for parameter := 1; parameter < len(args); parameter++ {
			args[parameter] = reflect.Zero(fn.Type().In(parameter))
		}
		outputs := fn.Call(args)
		err, _ := outputs[len(outputs)-1].Interface().(error)
		if auth.ErrorCode(err) != auth.CodePermissionDenied {
			t.Errorf("generated stocktake entry %s returned %v", name, err)
		}
	}
	if covered == 0 {
		t.Fatal("no generated stocktake handlers found")
	}
}
