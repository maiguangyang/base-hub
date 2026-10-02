package src_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	enginesrc "base-engine/src"
	"base-engine/src/services/ai"
	"base-engine/src/services/paymentconfig"
	"gorm.io/gorm"
)

func TestPaymentConfigHTTPBoundary(t *testing.T) {
	db := openResolverTestDB(t)
	keys := config.PaymentConfigSecurity{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)}}
	router := mux.NewRouter()
	enginesrc.RegisterPaymentConfigRoutes(router, paymentconfig.NewStore(db, keys), config.SecurityConfig{AllowedOrigins: map[string]struct{}{"https://admin.example": {}}})
	hq := &auth.WorkspacePrincipal{AccountID: "hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"paymentConfig:read": {}, "paymentConfig:manage": {}}}
	franchise := &auth.WorkspacePrincipal{AccountID: "franchise", WorkspaceType: auth.WorkspaceTypeFranchise, Permissions: hq.Permissions}
	path := "/api/payment-config?scope=GLOBAL"
	if got := paymentHTTPCall(router, nil, http.MethodGet, path, "", ""); got.Code != 401 {
		t.Fatalf("anonymous read: %d", got.Code)
	}
	if got := paymentHTTPCall(router, franchise, http.MethodGet, path, "", ""); got.Code != 403 {
		t.Fatalf("franchise read: %d", got.Code)
	}
	get := paymentHTTPCall(router, hq, http.MethodGet, path, "", "")
	if get.Code != 200 || !strings.Contains(get.Body.String(), `"UNCONFIGURED"`) {
		t.Fatalf("empty read: %d %s", get.Code, get.Body.String())
	}
	assertPaymentHTTPMutations(t, router, hq)
}

func TestPaymentConfigAIProtectedReadAndState(t *testing.T) {
	db := openResolverTestDB(t)
	keys := config.PaymentConfigSecurity{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)}}
	router := mux.NewRouter()
	enginesrc.RegisterPaymentConfigRoutes(router, paymentconfig.NewStore(db, keys), config.SecurityConfig{AllowedOrigins: map[string]struct{}{"https://admin.example": {}}})
	hq := &auth.WorkspacePrincipal{AccountID: "hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"paymentConfig:read": {}, "paymentConfig:manage": {}}}
	assertPaymentAIProtectedRead(t, router, hq)
	assertPaymentAIProtectedState(t, router, hq)
	assertPaymentAIProtectedRestore(t, db, router, hq)
}

func assertPaymentAIProtectedRestore(t *testing.T, db *gorm.DB, router *mux.Router, hq *auth.WorkspacePrincipal) {
	t.Helper()
	if err := db.Create(&gen.Organization{ID: "ai-pay-franchise", Code: "AIPAY", Name: "AI Payment", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}).Error; err != nil {
		t.Fatal(err)
	}
	state := json.RawMessage(`{"scope":"FRANCHISE","organizationId":"ai-pay-franchise","channel":"ALIPAY","state":"DISABLED","recordId":"","version":0}`)
	created := paymentAIProtectedCall(t, router, hq, http.MethodPost, "/api/payment-config/state", state)
	if created.Status != 200 {
		t.Fatalf("create override = %d %s", created.Status, created.Body)
	}
	var view paymentconfig.ScopeView
	if err := json.Unmarshal(created.Body, &view); err != nil {
		t.Fatal(err)
	}
	reset, err := json.Marshal(paymentconfig.ResetInput{ScopeRef: paymentconfig.ScopeRef{Scope: "FRANCHISE", OrganizationID: "ai-pay-franchise"}, Channel: "ALIPAY", RecordID: view.Channels[1].Own.RecordID, Version: view.Channels[1].Own.Version})
	if err != nil {
		t.Fatal(err)
	}
	path := "/api/payment-config/restore-inheritance"
	if got := paymentAIProtectedCall(t, router, hq, http.MethodPost, path, reset); got.Status != 200 {
		t.Fatalf("restore override = %d %s", got.Status, got.Body)
	}
	if got := paymentAIProtectedCall(t, router, hq, http.MethodPost, path, reset); got.Status != 409 {
		t.Fatalf("stale restore = %d", got.Status)
	}
}

func paymentAIProtectedCall(t *testing.T, router *mux.Router, principal *auth.WorkspacePrincipal, method, path string, body json.RawMessage) ai.FixedResponse {
	t.Helper()
	response, err := ai.ProtectedCall(auth.WithPrincipal(t.Context(), principal), router, ai.FixedRequest{Method: method, Path: path, Variables: body, Cookie: "session=test", Origin: "https://admin.example"})
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func paymentAIRead(t *testing.T, router *mux.Router, principal *auth.WorkspacePrincipal) ai.FixedResponse {
	t.Helper()
	return paymentAIProtectedCall(t, router, principal, http.MethodGet, "/api/payment-config", json.RawMessage(`{"scope":"GLOBAL"}`))
}

func assertPaymentAIProtectedRead(t *testing.T, router *mux.Router, hq *auth.WorkspacePrincipal) {
	t.Helper()
	if got := paymentAIRead(t, router, &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeFranchise, Permissions: hq.Permissions}); got.Status != 403 {
		t.Fatalf("franchise read = %d", got.Status)
	}
	if got := paymentAIRead(t, router, &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeHeadquarters}); got.Status != 403 {
		t.Fatalf("missing permission read = %d", got.Status)
	}
	if got := paymentAIRead(t, router, hq); got.Status != 200 || !bytes.Contains(got.Body, []byte(`"UNCONFIGURED"`)) || bytes.Contains(got.Body, []byte("credentialCiphertext")) {
		t.Fatalf("protected read = %d %s", got.Status, got.Body)
	}
}

func assertPaymentAIProtectedState(t *testing.T, router *mux.Router, hq *auth.WorkspacePrincipal) {
	t.Helper()
	state := json.RawMessage(`{"scope":"GLOBAL","channel":"WECHAT","state":"DISABLED","recordId":"","version":0}`)
	readonly := &auth.WorkspacePrincipal{AccountID: "hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"paymentConfig:read": {}}}
	if got := paymentAIProtectedCall(t, router, readonly, http.MethodPost, "/api/payment-config/state", state); got.Status != 403 {
		t.Fatalf("read-only write = %d", got.Status)
	}
	if got := paymentAIProtectedCall(t, router, hq, http.MethodPost, "/api/payment-config/state", state); got.Status != 200 {
		t.Fatalf("state write = %d", got.Status)
	}
	if got := paymentAIProtectedCall(t, router, hq, http.MethodPost, "/api/payment-config/state", state); got.Status != 409 {
		t.Fatalf("stale state = %d", got.Status)
	}
	if got := paymentAIRead(t, router, hq); got.Status != 200 || !bytes.Contains(got.Body, []byte(`"DISABLED"`)) {
		t.Fatalf("state not visible: %d %s", got.Status, got.Body)
	}
}

func assertPaymentHTTPMutations(t *testing.T, router *mux.Router, hq *auth.WorkspacePrincipal) {
	t.Helper()
	bad := paymentHTTPCall(router, hq, http.MethodPost, "/api/payment-config/state", `{"scope":"GLOBAL","channel":"WECHAT","state":"DISABLED","secret":"PRIVATE KEY"}`, "https://admin.example")
	if bad.Code != 400 || strings.Contains(bad.Body.String(), "PRIVATE KEY") {
		t.Fatalf("unknown field: %d %s", bad.Code, bad.Body.String())
	}
	created := paymentHTTPCall(router, hq, http.MethodPost, "/api/payment-config/state", `{"scope":"GLOBAL","channel":"WECHAT","state":"DISABLED"}`, "https://admin.example")
	if created.Code != 200 || strings.Contains(created.Body.String(), "PRIVATE KEY") {
		t.Fatalf("disabled override: %d %s", created.Code, created.Body.String())
	}
	var view paymentconfig.ScopeView
	if err := json.Unmarshal(created.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if len(view.Channels) != 2 || view.Channels[0].Own.State != "DISABLED" || view.Channels[1].Own.State != "UNCONFIGURED" {
		t.Fatalf("channel view: %#v", view)
	}
	stale := paymentHTTPCall(router, hq, http.MethodPost, "/api/payment-config/state", `{"scope":"GLOBAL","channel":"WECHAT","state":"VALID","recordId":"wrong","version":1}`, "https://admin.example")
	assertPaymentHTTPStatus(t, stale, 409)
	globalReset := paymentHTTPCall(router, hq, http.MethodPost, "/api/payment-config/restore-inheritance", `{"scope":"GLOBAL","channel":"WECHAT","recordId":"wrong","version":1}`, "https://admin.example")
	assertPaymentHTTPStatus(t, globalReset, 400)
}

func TestPaymentConfigHTTPUnavailableAndErrorState(t *testing.T) {
	db := openResolverTestDB(t)
	secret := "broken-ciphertext"
	if err := db.Create(&gen.GlobalPaymentConfig{ID: "damaged", Channel: "WECHAT", ConfigState: "VALID", Version: 1, KeyID: &secret, CredentialCiphertext: &secret}).Error; err != nil {
		t.Fatal(err)
	}
	router := mux.NewRouter()
	enginesrc.RegisterPaymentConfigRoutes(router, paymentconfig.NewStore(db, config.PaymentConfigSecurity{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)}}), config.SecurityConfig{})
	hq := &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"paymentConfig:read": {}}}
	response := paymentHTTPCall(router, hq, http.MethodGet, "/api/payment-config?scope=GLOBAL", "", "")
	if response.Code != 200 || !strings.Contains(response.Body.String(), `"DECRYPTION_FAILED"`) {
		t.Fatalf("damaged row: %d %s", response.Code, response.Body.String())
	}
	unavailable := mux.NewRouter()
	enginesrc.RegisterPaymentConfigRoutes(unavailable, paymentconfig.NewStore(db, config.PaymentConfigSecurity{}), config.SecurityConfig{})
	response = paymentHTTPCall(unavailable, hq, http.MethodGet, "/api/payment-config?scope=GLOBAL", "", "")
	if response.Code != 503 || !strings.Contains(response.Body.String(), `"CONFIG_UNAVAILABLE"`) {
		t.Fatalf("missing keyring: %d %s", response.Code, response.Body.String())
	}
}

func TestPaymentConfigHTTPRestoreAndValidation(t *testing.T) {
	db := openResolverTestDB(t)
	if err := db.Create(&gen.Organization{ID: "pay-franchise", Code: "PAY", Name: "Payment", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}).Error; err != nil {
		t.Fatal(err)
	}
	router := mux.NewRouter()
	keys := config.PaymentConfigSecurity{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)}}
	enginesrc.RegisterPaymentConfigRoutes(router, paymentconfig.NewStore(db, keys), config.SecurityConfig{AllowedOrigins: map[string]struct{}{"https://admin.example": {}}})
	hq := &auth.WorkspacePrincipal{AccountID: "hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"paymentConfig:read": {}, "paymentConfig:manage": {}}}
	statePath := "/api/payment-config/state"
	body := `{"scope":"FRANCHISE","organizationId":"pay-franchise","channel":"WECHAT","state":"DISABLED"}`
	for _, bad := range []struct {
		method, path, body, origin string
		status                     int
	}{
		{http.MethodPost, statePath, body, "", 403},
		{http.MethodPost, statePath, body, "https://other.example", 403},
		{http.MethodPost, statePath, strings.Repeat("x", 65537), "https://admin.example", 400},
		{http.MethodPut, "/api/payment-config", `{"scope":"GLOBAL","channel":"WECHAT","merchantId":"private","wechatCredentials":{"merchantPrivateKey":"PRIVATE KEY"}}`, "https://admin.example", 400},
	} {
		response := paymentHTTPCall(router, hq, bad.method, bad.path, bad.body, bad.origin)
		if response.Code != bad.status || strings.Contains(response.Body.String(), "PRIVATE KEY") {
			t.Fatalf("invalid input: %d %s", response.Code, response.Body.String())
		}
	}
	created := paymentHTTPCall(router, hq, http.MethodPost, statePath, body, "https://admin.example")
	if created.Code != 200 {
		t.Fatalf("create: %d %s", created.Code, created.Body.String())
	}
	var view paymentconfig.ScopeView
	if err := json.Unmarshal(created.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	reset, err := json.Marshal(paymentconfig.ResetInput{ScopeRef: paymentconfig.ScopeRef{Scope: "FRANCHISE", OrganizationID: "pay-franchise"}, Channel: "WECHAT", RecordID: view.Channels[0].Own.RecordID, Version: view.Channels[0].Own.Version})
	if err != nil {
		t.Fatal(err)
	}
	path := "/api/payment-config/restore-inheritance"
	assertPaymentHTTPStatus(t, paymentHTTPCall(router, hq, http.MethodPost, path, string(reset), "https://admin.example"), 200)
	assertPaymentHTTPStatus(t, paymentHTTPCall(router, hq, http.MethodPost, statePath, body, "https://admin.example"), 200)
	assertPaymentHTTPStatus(t, paymentHTTPCall(router, hq, http.MethodPost, path, string(reset), "https://admin.example"), 409)
}

func assertPaymentHTTPStatus(t *testing.T, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	if response.Code != want {
		t.Fatalf("status=%d, want=%d, body=%s", response.Code, want, response.Body.String())
	}
}

func paymentHTTPCall(router *mux.Router, principal *auth.WorkspacePrincipal, method, path, body, origin string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	if origin != "" {
		request.Header.Set("Origin", origin)
	}
	if principal != nil {
		request = request.WithContext(auth.WithPrincipal(request.Context(), principal))
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}
