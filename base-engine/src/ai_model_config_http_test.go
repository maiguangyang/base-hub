package src

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"base-engine/src/services/ai"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAIModelConfigHTTPRedactsKeyAndChecksPermission(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:ai-model-http?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&ai.StoredModelConfig{}, &gen.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	security := config.AIModelSecurityConfig{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)}}
	store := ai.NewModelConfigStore(db, security, nil)
	router := mux.NewRouter()
	RegisterAIModelConfigRoutes(router, store, config.SecurityConfig{AllowedOrigins: map[string]struct{}{"https://admin.example": {}}})
	principal := &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"aiModelConfig:read": {}, "aiModelConfig:manage": {}}}
	put := modelConfigHTTPCall(router, principal, http.MethodPut, "/api/ai/model-config", `{"modelName":"example","baseUrl":"https://model.example/v1","apiKey":"private-key","version":0}`, "https://admin.example")
	if put.Code != http.StatusOK || strings.Contains(put.Body.String(), "private-key") {
		t.Fatalf("save status=%d body=%s", put.Code, put.Body.String())
	}
	assertModelConfigHTTPRedaction(t, router, db, principal)
	principal.WorkspaceType = auth.WorkspaceTypeFranchise
	denied := modelConfigHTTPCall(router, principal, http.MethodGet, "/api/ai/model-config", "", "")
	if denied.Code != http.StatusForbidden {
		t.Fatalf("franchise status = %d", denied.Code)
	}
}

func assertModelConfigHTTPRedaction(t *testing.T, router *mux.Router, db *gorm.DB, principal *auth.WorkspacePrincipal) {
	t.Helper()
	get := modelConfigHTTPCall(router, principal, http.MethodGet, "/api/ai/model-config", "", "")
	if get.Code != http.StatusOK || strings.Contains(get.Body.String(), "private-key") {
		t.Fatalf("read status=%d body=%s", get.Code, get.Body.String())
	}
	var status ai.ModelConfigStatus
	if err := json.Unmarshal(get.Body.Bytes(), &status); err != nil || !status.KeyConfigured {
		t.Fatalf("read status = %+v, %v", status, err)
	}
	var audit gen.AuditLog
	if err := db.Where("action = ?", "aiModelConfig:save").First(&audit).Error; err != nil {
		t.Fatal(err)
	}
	if audit.MetadataJSON == nil || strings.Contains(*audit.MetadataJSON, "private-key") {
		t.Fatalf("unsafe audit metadata = %#v", audit.MetadataJSON)
	}
}

func TestAIModelConfigHTTPRejectsMalformedBodyWithoutLeakingIt(t *testing.T) {
	router := mux.NewRouter()
	RegisterAIModelConfigRoutes(router, nil, config.SecurityConfig{AllowedOrigins: map[string]struct{}{"https://admin.example": {}}})
	response := modelConfigHTTPCall(router, &auth.WorkspacePrincipal{}, http.MethodPut, "/api/ai/model-config", `{"apiKey":"private-key","unexpected":true}`, "https://admin.example")
	if response.Code != http.StatusBadRequest || strings.Contains(response.Body.String(), "private-key") {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestAIModelConfigHTTPTimeoutIsSafe(t *testing.T) {
	status, code := modelConfigHTTPError(context.DeadlineExceeded)
	if status != http.StatusGatewayTimeout || code != "MODEL_PROBE_TIMEOUT" {
		t.Fatalf("timeout mapping = %d %s", status, code)
	}
}

func TestAIModelConfigHTTPSeparatesProtocolFailure(t *testing.T) {
	status, code := modelConfigHTTPError(ai.ErrModelProtocol)
	if status != http.StatusBadGateway || code != "MODEL_PROTOCOL_UNSUPPORTED" {
		t.Fatalf("protocol mapping = %d %s", status, code)
	}
}

type testUpstreamStatusFailure struct{ status int }

func (e testUpstreamStatusFailure) Error() string           { return "safe upstream status" }
func (e testUpstreamStatusFailure) Unwrap() error           { return ai.ErrModelUnavailable }
func (e testUpstreamStatusFailure) UpstreamStatusCode() int { return e.status }

func TestAIModelConfigHTTPClassifiesUpstreamStatus(t *testing.T) {
	cases := []struct {
		status int
		code   string
	}{
		{http.StatusUnauthorized, "MODEL_UPSTREAM_AUTH"},
		{http.StatusForbidden, "MODEL_UPSTREAM_AUTH"},
		{http.StatusNotFound, "MODEL_TARGET_NOT_FOUND"},
		{http.StatusTooManyRequests, "MODEL_UPSTREAM_RATE_LIMITED"},
		{http.StatusBadRequest, "MODEL_REQUEST_REJECTED"},
		{http.StatusServiceUnavailable, "MODEL_UPSTREAM_UNAVAILABLE"},
	}
	for _, tc := range cases {
		status, code := modelConfigHTTPError(testUpstreamStatusFailure{status: tc.status})
		if status != http.StatusBadGateway || code != tc.code {
			t.Errorf("upstream %d mapping = %d %s, want %s", tc.status, status, code, tc.code)
		}
	}
}

func TestAIModelConfigHTTPClassifiesTransportFailure(t *testing.T) {
	status, code := modelConfigHTTPError(ai.ErrModelUnavailable)
	if status != http.StatusBadGateway || code != "MODEL_CONNECTION_FAILED" {
		t.Fatalf("transport mapping = %d %s", status, code)
	}
}

func TestAIModelConfigHTTPClassifiesRealProbeFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"private-upstream-detail"}`))
	}))
	defer server.Close()
	err := ai.ProbeModelConnection(t.Context(), ai.ModelConfig{
		Name: "test-model", BaseURL: server.URL + "/v1", APIKey: "private-key",
	})
	status, code := modelConfigHTTPError(err)
	if status != http.StatusBadGateway || code != "MODEL_UPSTREAM_AUTH" || strings.Contains(code, "private") {
		t.Fatalf("real probe mapping = %d %s, error=%v", status, code, err)
	}
}

func TestAIModelConfigHTTPRequiresTrustedOrigin(t *testing.T) {
	router := mux.NewRouter()
	RegisterAIModelConfigRoutes(router, nil, config.SecurityConfig{AllowedOrigins: map[string]struct{}{"https://admin.example": {}}})
	for _, origin := range []string{"", "https://other.example"} {
		response := modelConfigHTTPCall(router, nil, http.MethodPut, "/api/ai/model-config", `{}`, origin)
		if response.Code != http.StatusForbidden || response.Header().Get("Cache-Control") != "no-store" || !strings.Contains(response.Body.String(), "ORIGIN_NOT_ALLOWED") {
			t.Fatalf("origin=%q status=%d", origin, response.Code)
		}
	}
}

func TestAIModelConfigHTTPRejectsWildcardOriginPolicy(t *testing.T) {
	router := mux.NewRouter()
	RegisterAIModelConfigRoutes(router, nil, config.SecurityConfig{AllowedOrigins: map[string]struct{}{"*": {}}})
	for _, method := range []string{http.MethodGet, http.MethodPut} {
		response := modelConfigHTTPCall(router, nil, method, "/api/ai/model-config", `{}`, "https://unlisted.example")
		if response.Code != http.StatusForbidden {
			t.Fatalf("wildcard origin accepted for %s: %d", method, response.Code)
		}
	}
}

func TestAIModelConfigHTTPRateLimitsFailedProbe(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:ai-model-probe-limit?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&ai.StoredModelConfig{}, &gen.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	security := config.AIModelSecurityConfig{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)}}
	store := ai.NewModelConfigStore(db, security, func(context.Context, ai.ModelConfig) error { return ai.ErrModelUnavailable })
	router := mux.NewRouter()
	RegisterAIModelConfigRoutes(router, store, config.SecurityConfig{AllowedOrigins: map[string]struct{}{"https://admin.example": {}}})
	principal := &auth.WorkspacePrincipal{AccountID: "hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"aiModelConfig:manage": {}}}
	put := modelConfigHTTPCall(router, principal, http.MethodPut, "/api/ai/model-config", `{"modelName":"example","baseUrl":"https://model.example/v1","apiKey":"private-key","version":0}`, "https://admin.example")
	if put.Code != http.StatusOK {
		t.Fatalf("save status = %d", put.Code)
	}
	for attempt := 1; attempt <= 4; attempt++ {
		response := modelConfigHTTPCall(router, principal, http.MethodPost, "/api/ai/model-config/probe", `{"version":1}`, "https://admin.example")
		wanted := http.StatusBadGateway
		if attempt == 4 {
			wanted = http.StatusTooManyRequests
		}
		if response.Code != wanted || strings.Contains(response.Body.String(), "private-key") {
			t.Fatalf("attempt=%d status=%d body=%s", attempt, response.Code, response.Body.String())
		}
	}
}

func modelConfigHTTPCall(router *mux.Router, principal *auth.WorkspacePrincipal, method, path, body, origin string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Origin", origin)
	request.Header.Set("Content-Type", "application/json")
	if principal != nil {
		request = request.WithContext(auth.WithPrincipal(context.Background(), principal))
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}
