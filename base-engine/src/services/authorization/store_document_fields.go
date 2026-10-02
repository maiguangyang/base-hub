package authorization

import "base-engine/auth"

func rejectStoreDocumentFields(input map[string]interface{}) error {
	for _, key := range []string{"businessLicenseImageUrl", "otherDocumentImageUrl"} {
		if _, supplied := input[key]; supplied { return auth.NewError(auth.CodePermissionDenied) }
	}
	return nil
}
