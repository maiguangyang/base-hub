package src

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	storeservice "base-engine/src/services/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestStoreDocumentHTTPRequiresSessionForPreview(t *testing.T) {
	router, service, p := storeDocumentHTTPFixture(t)
	imageBytes := storeDocumentTestPNG(t)
	request := httptest.NewRequest(http.MethodPost, "/api/store-documents", bytes.NewReader(imageBytes))
	request.Header.Set("Origin", "https://admin.example")
	request.Header.Set("Content-Type", "image/png")
	request = request.WithContext(auth.WithPrincipal(request.Context(), p))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("upload: %d %s", response.Code, response.Body.String())
	}
	var result storeDocumentUploadResult
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	item, err := service.SetDocument(t.Context(), p, "store-1", gen.StoreDocumentKindOther, result.AttachmentID)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		principal *auth.WorkspacePrincipal
		want      int
	}{{nil, 404}, {p, 200}} {
		get := httptest.NewRequest(http.MethodGet, *item.OtherDocumentImageURL, nil)
		if tc.principal != nil {
			get = get.WithContext(auth.WithPrincipal(get.Context(), tc.principal))
		}
		preview := httptest.NewRecorder()
		router.ServeHTTP(preview, get)
		if preview.Code != tc.want {
			t.Fatalf("preview status %d, want %d", preview.Code, tc.want)
		}
		if preview.Header().Get("Cache-Control") != "no-store" {
			t.Fatal("document response may be cached")
		}
	}
}

func storeDocumentHTTPFixture(t *testing.T) (*mux.Router, *storeservice.Service, *auth.WorkspacePrincipal) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&gen.Store{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.Store{ID: "store-1", Code: "S1", Name: "Store", OrganizationID: "org-1", Lifecycle: gen.StoreLifecycleDraft}).Error; err != nil {
		t.Fatal(err)
	}
	service := storeservice.NewService(db, nil)
	service.SetDocumentRoot(t.TempDir())
	router := mux.NewRouter()
	RegisterStoreDocumentRoutes(router, service, config.SecurityConfig{AllowedOrigins: map[string]struct{}{"https://admin.example": {}}})
	org := "org-1"
	p := &auth.WorkspacePrincipal{AccountID: "account-1", SessionID: "session-1", WorkspaceType: auth.WorkspaceTypeFranchise, OrganizationID: &org, AllStores: true, Permissions: map[string]struct{}{"store:create": {}, "store:update": {}, "store:read": {}}}
	return router, service, p
}

func storeDocumentTestPNG(t *testing.T) []byte {
	t.Helper()
	var imageBytes bytes.Buffer
	if err := png.Encode(&imageBytes, image.NewRGBA(image.Rect(0, 0, 1, 1))); err != nil {
		t.Fatal(err)
	}
	return imageBytes.Bytes()
}
