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
	"base-engine/src/services/paymentconfig"
	"base-engine/src/services/session"
)

func TestPaymentConfigUsesDatabaseSessionAndWorkspace(t *testing.T) {
	fixture := newSecurityFixture(t)
	keys := config.PaymentConfigSecurity{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)}}
	router := mux.NewRouter()
	enginesrc.RegisterPaymentConfigRoutes(router, paymentconfig.NewStore(fixture.db, keys), fixture.cfg)
	deps := enginesrc.NewDependencies(fixture.db, fixture.cfg, session.NewPublisher())
	handler := middleware.New(middleware.Dependencies{Config: fixture.cfg, PrincipalResolver: deps.Principal})(router)
	if response := callPaymentConfig(t, handler, fixture, "session-hq", http.MethodGet, "/api/payment-config?scope=GLOBAL", ""); response.Code != 403 {
		t.Fatalf("missing permission: %d", response.Code)
	}
	fixture.createAll([]gen.Permission{
		{ID: "payment-read", Name: "paymentConfig:read", Action: "paymentConfig:read", Module: "paymentConfig", Scope: gen.PermissionScopeSystem},
		{ID: "payment-manage", Name: "paymentConfig:manage", Action: "paymentConfig:manage", Module: "paymentConfig", Scope: gen.PermissionScopeSystem},
	}, []permissionRole{{"payment-read", "role-hq"}, {"payment-manage", "role-hq"}})
	if response := callPaymentConfig(t, handler, fixture, "session-hq", http.MethodGet, "/api/payment-config?scope=GLOBAL", ""); response.Code != 200 {
		t.Fatalf("HQ read: %d %s", response.Code, response.Body.String())
	}
	if response := callPaymentConfig(t, handler, fixture, "session-a-1", http.MethodGet, "/api/payment-config?scope=GLOBAL", ""); response.Code != 403 {
		t.Fatalf("franchise read: %d", response.Code)
	}
	assertPaymentConfigStoreScope(t, handler, fixture)
	now := time.Now()
	if err := fixture.db.Model(&gen.Session{}).Where("id = ?", "session-hq").Update("revoked_at", now).Error; err != nil {
		t.Fatal(err)
	}
	if response := callPaymentConfig(t, handler, fixture, "session-hq", http.MethodGet, "/api/payment-config?scope=GLOBAL", ""); response.Code != 401 {
		t.Fatalf("revoked session: %d", response.Code)
	}
}

func assertPaymentConfigStoreScope(t *testing.T, handler http.Handler, fixture *securityFixture) {
	t.Helper()
	for _, body := range []string{
		`{"scope":"STORE","storeId":"store-hq","channel":"WECHAT","state":"DISABLED"}`,
		`{"scope":"STORE","storeId":"store-a","organizationId":"org-b","channel":"WECHAT","state":"DISABLED"}`,
	} {
		response := callPaymentConfig(t, handler, fixture, "session-hq", http.MethodPost, "/api/payment-config/state", body)
		if response.Code != 400 {
			t.Fatalf("invalid store target: %d %s", response.Code, response.Body.String())
		}
	}
	valid := `{"scope":"STORE","storeId":"store-a","channel":"WECHAT","state":"DISABLED"}`
	if response := callPaymentConfig(t, handler, fixture, "session-hq", http.MethodPost, "/api/payment-config/state", valid); response.Code != 200 {
		t.Fatalf("HQ store override: %d %s", response.Code, response.Body.String())
	}
	if response := callPaymentConfig(t, handler, fixture, "session-a-1", http.MethodPost, "/api/payment-config/state", valid); response.Code != 403 {
		t.Fatalf("franchise write: %d", response.Code)
	}
}

func callPaymentConfig(t *testing.T, handler http.Handler, fixture *securityFixture, sessionID, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Origin", "https://admin.example.com")
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(fixture.cookies[sessionID])
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
