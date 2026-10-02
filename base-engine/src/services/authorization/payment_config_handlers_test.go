package authorization

import (
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestPaymentConfigNestedStoreWritesAreDenied(t *testing.T) {
	for _, field := range []string{"paymentConfigs", "paymentConfigsIds"} {
		if err := rejectStoreManagedRelations(map[string]interface{}{field: []string{"unsafe"}}); auth.ErrorCode(err) != auth.CodePermissionDenied {
			t.Fatalf("nested payment field %s was accepted: %v", field, err)
		}
	}
}

func TestPaymentConfigGeneratedHandlersAreReplaced(t *testing.T) {
	defaults := gen.DefaultResolutionHandlers()
	if err := ValidateHandlerCoverage(RegisterHandlers(defaults), defaults); err != nil {
		t.Fatal(err)
	}
}
