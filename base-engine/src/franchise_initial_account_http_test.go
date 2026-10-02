package src

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
)

func TestInitialAccountHTTPErrorClassification(t *testing.T) {
	status, code := initialAccountHTTPError(errors.New("database unavailable"))
	if status != http.StatusInternalServerError || code != auth.CodeInternalError {
		t.Fatalf("internal error mapped to %d %s", status, code)
	}
	status, code = initialAccountHTTPError(auth.NewError(auth.CodeOpeningRecordNumberConflict))
	if status != http.StatusConflict || code != auth.CodeOpeningRecordNumberConflict {
		t.Fatalf("number conflict mapped to %d %s", status, code)
	}
}

func TestFranchiseInitialAccountRouteRequiresTrustedOriginAndPrincipal(t *testing.T) {
	router := mux.NewRouter()
	RegisterFranchiseInitialAccountRoute(router, nil, config.SecurityConfig{AllowedOrigins: map[string]struct{}{"http://admin.test": {}}})
	for _, test := range []struct {
		origin string
		status int
	}{{"", http.StatusForbidden}, {"http://other.test", http.StatusForbidden}, {"http://admin.test", http.StatusUnauthorized}} {
		request := httptest.NewRequest(http.MethodPost, "/api/franchise-initial-account", strings.NewReader(`{"organizationId":"org-a"}`))
		request.Header.Set("Origin", test.origin)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != test.status {
			t.Fatalf("origin=%q status=%d want=%d", test.origin, response.Code, test.status)
		}
		if response.Header().Get("Cache-Control") != "no-store" {
			t.Fatal("sensitive response can be cached")
		}
	}
}
