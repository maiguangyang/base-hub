package integration_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"base-engine/config"
	"base-engine/gen"
	enginesrc "base-engine/src"
	"base-engine/src/middleware"
	"base-engine/src/services/ai"
	"base-engine/src/services/session"
)

func TestAIModelConfigUsesDatabaseSessionAndWorkspace(t *testing.T) {
	fixture := newSecurityFixture(t)
	if err := fixture.db.AutoMigrate(&ai.StoredModelConfig{}); err != nil {
		t.Fatal(err)
	}
	fixture.createAll([]gen.Permission{
		{ID: "permission-ai-read", Name: "aiModelConfig:read", Action: "aiModelConfig:read", Module: "ai", Scope: gen.PermissionScopeSystem},
		{ID: "permission-ai-manage", Name: "aiModelConfig:manage", Action: "aiModelConfig:manage", Module: "ai", Scope: gen.PermissionScopeSystem},
	}, []permissionRole{{"permission-ai-read", "role-hq"}, {"permission-ai-manage", "role-hq"}})
	aiSecurity := config.AIModelSecurityConfig{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)}}
	router := mux.NewRouter()
	enginesrc.RegisterAIModelConfigRoutes(router, ai.NewModelConfigStore(fixture.db, aiSecurity, nil), fixture.cfg)
	deps := enginesrc.NewDependencies(fixture.db, fixture.cfg, session.NewPublisher())
	handler := middleware.New(middleware.Dependencies{Config: fixture.cfg, PrincipalResolver: deps.Principal})(router)
	body := `{"modelName":"example","baseUrl":"http://zsgw.sjdistributor.com:4000/v1","apiKey":"private-key","version":0}`
	if response := callAIModelConfig(t, handler, fixture, "session-hq", http.MethodPut, body); response.Code != http.StatusOK {
		t.Fatalf("HQ save = %d: %s", response.Code, response.Body.String())
	}
	if response := callAIModelConfig(t, handler, fixture, "session-a-1", http.MethodGet, ""); response.Code != http.StatusForbidden {
		t.Fatalf("franchise read = %d", response.Code)
	}
	now := time.Now()
	if err := fixture.db.Model(&gen.Session{}).Where("id = ?", "session-hq").Update("revoked_at", now).Error; err != nil {
		t.Fatal(err)
	}
	if response := callAIModelConfig(t, handler, fixture, "session-hq", http.MethodGet, ""); response.Code != http.StatusUnauthorized {
		t.Fatalf("revoked session read = %d", response.Code)
	}
}

func callAIModelConfig(t *testing.T, handler http.Handler, fixture *securityFixture, sessionID, method, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, "/api/ai/model-config", strings.NewReader(body))
	request.Header.Set("Origin", "https://admin.example.com")
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(fixture.cookies[sessionID])
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
