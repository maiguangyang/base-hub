package integration_test

import (
	"bytes"
	"image"
	"image/png"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
)

func documentIntegrationPNG(t *testing.T) []byte {
	t.Helper()
	var data bytes.Buffer
	if err := png.Encode(&data, image.NewRGBA(image.Rect(0, 0, 1, 1))); err != nil {
		t.Fatal(err)
	}
	return data.Bytes()
}

func TestStoreDocumentGraphQLContract(t *testing.T) {
	fixture := newSecurityFixture(t)
	fixture.resolver.Services.Stores.SetDocumentRoot(t.TempDir())
	if err := fixture.db.Model(&gen.Store{}).Where("id = ?", "store-a").Update("lifecycle", gen.StoreLifecycleDraft).Error; err != nil {
		t.Fatal(err)
	}
	owner, err := fixture.resolver.Services.Principal.Resolve(t.Context(), fixture.cookies["session-a-1"].Value, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	id, err := fixture.resolver.Services.Stores.UploadDocument(t.Context(), owner, documentIntegrationPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	assertCode(t, fixture.execute("session-a-1", `mutation { updateStore(id: "store-a", input: {businessLicenseImageUrl: "/uploads/stores/store-b/foreign.png"}) { id } }`), auth.CodePermissionDenied)
	assertCode(t, fixture.execute("session-a-1", `mutation { createStore(input: {code: "DIRECT", name: "Direct", lifecycle: DRAFT, organizationId: "org-a", otherDocumentImageUrl: null}) { id } }`), auth.CodePermissionDenied)
	assertCode(t, fixture.execute("session-b", `mutation { setStoreDocument(storeId: "store-a", kind: OTHER, attachmentId: "`+id+`") { id } }`), auth.CodePermissionDenied)
	bound := fixture.execute("session-a-1", `mutation { setStoreDocument(storeId: "store-a", kind: BUSINESS_LICENSE, attachmentId: "`+id+`") { id businessLicenseImageUrl } }`)
	assertNoErrors(t, bound)
	if !strings.Contains(bound.Body, "/uploads/stores/store-a/") {
		t.Fatalf("document path missing: %s", bound.Body)
	}
	assertAuditCount(t, fixture, "store:document:set", 1)
	assertCode(t, fixture.execute("session-a-1", `mutation { setStoreDocument(storeId: "store-a", kind: OTHER, attachmentId: "`+id+`") { id } }`), auth.CodeValidationFailed)
	assertNoErrors(t, fixture.execute("session-a-1", `mutation { removeStoreDocument(storeId: "store-a", kind: BUSINESS_LICENSE) { id businessLicenseImageUrl } }`))
	assertAuditCount(t, fixture, "store:document:remove", 1)
	created := fixture.execute("session-a-1", `mutation { createStore(input: {code: "DOCUMENT-TEST", name: "Document Test", lifecycle: DRAFT, organizationId: "org-a"}) { id createdBy } }`)
	assertNoErrors(t, created)
	if !strings.Contains(created.Body, `"createdBy":"`+owner.AccountID+`"`) {
		t.Fatalf("store creator metadata missing: %s", created.Body)
	}
}
