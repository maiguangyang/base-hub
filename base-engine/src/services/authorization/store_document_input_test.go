package authorization

import (
	"testing"

	"base-engine/auth"
)

func TestStoreDocumentPathsCannotUseOrdinaryMutations(t *testing.T) {
	for _, field := range []string{"businessLicenseImageUrl", "otherDocumentImageUrl"} {
		for _, value := range []interface{}{"/uploads/stores/foreign/image.png", nil} {
			if err := rejectStoreManagedRelations(map[string]interface{}{field: value}); auth.ErrorCode(err) != auth.CodePermissionDenied {
				t.Fatalf("%s=%v accepted: %v", field, value, err)
			}
		}
	}
}
