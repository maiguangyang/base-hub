package src

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/src/services/audit"
	"base-engine/src/services/productcatalog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestProductImageHTTPRequiresOriginAndOwnedSession(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:product-image-http?mode=memory&cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, IgnoreRelationshipsWhenMigrating: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	service := productcatalog.NewService(db, audit.NewService())
	service.SetImageRoot(t.TempDir())
	router := mux.NewRouter()
	RegisterProductImageRoutes(router, service, config.SecurityConfig{AllowedOrigins: map[string]struct{}{"https://admin.example": {}}})
	hqID := "hq"
	principal := &auth.WorkspacePrincipal{AccountID: "admin", SessionID: "session", WorkspaceType: auth.WorkspaceTypeHeadquarters,
		OrganizationID: &hqID, Permissions: map[string]struct{}{"hqProductCatalog:manage": {}, "hqProductCatalog:read": {}}}
	imageBytes := httpTestPNG(t)
	denied := imageHTTPCall(router, principal, imageBytes, "https://other.example")
	if denied.Code != http.StatusForbidden {
		t.Fatalf("foreign origin status = %d", denied.Code)
	}
	denied = imageHTTPCall(router, nil, imageBytes, "https://admin.example")
	if denied.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous status = %d", denied.Code)
	}
	response := imageHTTPCall(router, principal, imageBytes, "https://admin.example")
	if response.Code != http.StatusCreated {
		t.Fatalf("upload status = %d: %s", response.Code, response.Body.String())
	}
	var upload productImageUploadResult
	if err := json.Unmarshal(response.Body.Bytes(), &upload); err != nil || upload.AttachmentID == "" {
		t.Fatalf("upload = %+v, %v", upload, err)
	}
	if err := service.ValidateMainImageAttachment(context.Background(), principal, upload.AttachmentID); err != nil {
		t.Fatal(err)
	}
	missing := httptest.NewRecorder()
	router.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/uploads/products/hq/anything.png", nil))
	if missing.Code != http.StatusNotFound {
		t.Fatalf("unbound image status = %d", missing.Code)
	}
	if missing.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("image response must disable content sniffing: %q", missing.Header().Get("X-Content-Type-Options"))
	}
}

func imageHTTPCall(router *mux.Router, principal *auth.WorkspacePrincipal, body []byte, origin string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, "/api/product-main-images", bytes.NewReader(body))
	request.Header.Set("Content-Type", "image/png")
	request.Header.Set("Origin", origin)
	if principal != nil {
		request = request.WithContext(auth.WithPrincipal(request.Context(), principal))
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func httpTestPNG(t *testing.T) []byte {
	t.Helper()
	var output bytes.Buffer
	if err := png.Encode(&output, image.NewRGBA(image.Rect(0, 0, 1, 1))); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}
