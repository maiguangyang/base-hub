/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/config"
	"base-engine/src"
)

// TestMiddlewareAllowsExactOriginAndInjectsContext 验证 CORS 与请求上下文的完整契约。
func TestMiddlewareAllowsExactOriginAndInjectsContext(t *testing.T) {
	cfg := middlewareTestConfig()
	resolver := &stubPrincipalResolver{principal: &auth.WorkspacePrincipal{AccountID: "account-1"}}
	called := false
	next := func(writer http.ResponseWriter, request *http.Request) {
		called = true
		assertRequestContext(t, request.Context())
		writer.WriteHeader(http.StatusNoContent)
	}
	handler := New(Dependencies{Config: cfg, PrincipalResolver: resolver, Now: time.Now})(http.HandlerFunc(next))
	request := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	request.RemoteAddr = "192.0.2.4:443"
	request.Header.Set("Origin", "https://admin.example.com")
	request.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: "token"})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if !called || response.Code != http.StatusNoContent {
		t.Fatalf("handler result: called=%v code=%d", called, response.Code)
	}
	if response.Header().Get("Access-Control-Allow-Origin") != "https://admin.example.com" || response.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatalf("invalid CORS headers: %#v", response.Header())
	}
	if response.Header().Get("Access-Control-Allow-Origin") == "*" {
		t.Fatal("wildcard origin returned")
	}
}

// TestMiddlewareWildcardReflectsCredentialedOrigin 验证开放来源时仍回显实际 Origin 以支持 Cookie。
func TestMiddlewareWildcardReflectsCredentialedOrigin(t *testing.T) {
	called := false
	cfg := config.SecurityConfig{AllowedOrigins: map[string]struct{}{"*": {}}}
	handler := New(Dependencies{Config: cfg})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}))
	request := httptest.NewRequest(http.MethodGet, "/api/system-initialization", nil)
	request.Header.Set("Origin", "http://localhost:4322")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if !called || response.Code != http.StatusOK {
		t.Fatalf("handler result: called=%v code=%d", called, response.Code)
	}
	if response.Header().Get("Access-Control-Allow-Origin") != "http://localhost:4322" {
		t.Fatalf("origin header = %q", response.Header().Get("Access-Control-Allow-Origin"))
	}
	if response.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatalf("credentials header = %q", response.Header().Get("Access-Control-Allow-Credentials"))
	}
}

func assertRequestContext(t *testing.T, ctx context.Context) {
	t.Helper()
	principal, ok := auth.PrincipalFromContext(ctx)
	if !ok {
		t.Fatal("principal missing from context")
	}
	if principal.AccountID != "account-1" {
		t.Fatalf("account ID = %s", principal.AccountID)
	}
	if principal.RequestID == "" || principal.RequestID != src.RequestIDFromContext(ctx) {
		t.Fatalf("principal request ID = %q, context request ID = %q", principal.RequestID, src.RequestIDFromContext(ctx))
	}
	if src.RequestIDFromContext(ctx) == "" {
		t.Fatal("request ID missing")
	}
	if src.ResponseWriterFromContext(ctx) == nil {
		t.Fatal("response writer missing")
	}
	if src.RemoteIPFromContext(ctx) != "192.0.2.4" {
		t.Fatal("remote IP missing")
	}
	if src.SessionTokenFromContext(ctx) != "token" {
		t.Fatal("cookie token missing")
	}
}

// TestMiddlewareRejectsUnlistedHTTPAndWebSocketOrigins 验证 HTTP 与升级请求共享来源白名单。
func TestMiddlewareRejectsUnlistedHTTPAndWebSocketOrigins(t *testing.T) {
	handler := New(Dependencies{Config: middlewareTestConfig()})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("unlisted origin reached application")
	}))
	for _, upgrade := range []bool{false, true} {
		request := httptest.NewRequest(http.MethodPost, "/graphql", nil)
		request.Header.Set("Origin", "https://evil.example.com")
		if upgrade {
			request.Header.Set("Connection", "Upgrade")
			request.Header.Set("Upgrade", "websocket")
		}
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusForbidden || response.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Fatalf("upgrade=%v code=%d headers=%#v", upgrade, response.Code, response.Header())
		}
	}
}

// TestMiddlewareOptionsReturnsBeforeGraphQL 验证预检请求不会进入 GraphQL。
func TestMiddlewareOptionsReturnsBeforeGraphQL(t *testing.T) {
	handler := New(Dependencies{Config: middlewareTestConfig()})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("OPTIONS reached GraphQL")
	}))
	request := httptest.NewRequest(http.MethodOptions, "/graphql", nil)
	request.Header.Set("Origin", "https://admin.example.com")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d", response.Code)
	}
	if !strings.Contains(response.Header().Get("Access-Control-Allow-Methods"), "PUT") {
		t.Fatalf("model configuration PUT preflight disallowed: %#v", response.Header())
	}
}

// TestMiddlewarePreservesTerminalAuthenticationCode 验证权威会话终态不会被降级成通用未登录。
func TestMiddlewarePreservesTerminalAuthenticationCode(t *testing.T) {
	resolver := &stubPrincipalResolver{err: auth.NewError(auth.CodeSessionRevoked)}
	handler := New(Dependencies{Config: middlewareTestConfig(), PrincipalResolver: resolver})(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		_, err := auth.RequirePrincipal(request.Context())
		if auth.ErrorCode(err) != auth.CodeSessionRevoked {
			t.Fatalf("authentication code = %s", auth.ErrorCode(err))
		}
	}))
	request := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	request.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: "revoked"})
	handler.ServeHTTP(httptest.NewRecorder(), request)
}

type stubPrincipalResolver struct {
	principal *auth.WorkspacePrincipal
	err       error
	calls     int
}

func (s *stubPrincipalResolver) Resolve(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) {
	s.calls++
	return s.principal, s.err
}

// TestMiddlewareWebSocketRevalidatesPrincipal 验证升级连接上的每次受保护操作都会重新解析 Session。
func TestMiddlewareWebSocketRevalidatesPrincipal(t *testing.T) {
	resolver := &stubPrincipalResolver{principal: &auth.WorkspacePrincipal{AccountID: "account-1"}}
	handler := New(Dependencies{Config: middlewareTestConfig(), PrincipalResolver: resolver})(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		if _, err := auth.RequirePrincipal(request.Context()); err != nil {
			t.Fatalf("initial principal: %v", err)
		}
		resolver.err = auth.NewError(auth.CodeSessionRevoked)
		if _, err := auth.RequirePrincipal(request.Context()); auth.ErrorCode(err) != auth.CodeSessionRevoked {
			t.Fatalf("authentication code = %s", auth.ErrorCode(err))
		}
	}))
	request := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "websocket")
	request.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: "token"})
	handler.ServeHTTP(httptest.NewRecorder(), request)
	if resolver.calls != 3 {
		t.Fatalf("resolver calls = %d, want 3", resolver.calls)
	}
}

// TestMiddlewareHTTPKeepsRequestPrincipal 验证普通 HTTP 请求不会为每个 resolver 重复读取 Session。
func TestMiddlewareHTTPKeepsRequestPrincipal(t *testing.T) {
	resolver := &stubPrincipalResolver{principal: &auth.WorkspacePrincipal{AccountID: "account-1"}}
	handler := New(Dependencies{Config: middlewareTestConfig(), PrincipalResolver: resolver})(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		if _, err := auth.RequirePrincipal(request.Context()); err != nil {
			t.Fatalf("initial principal: %v", err)
		}
		resolver.err = auth.NewError(auth.CodeSessionRevoked)
		if _, err := auth.RequirePrincipal(request.Context()); err != nil {
			t.Fatalf("request principal was unexpectedly reloaded: %v", err)
		}
	}))
	request := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	request.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: "token"})
	handler.ServeHTTP(httptest.NewRecorder(), request)
	if resolver.calls != 1 {
		t.Fatalf("resolver calls = %d, want 1", resolver.calls)
	}
}

// TestMiddlewareNormalizesWebSocketHostForAllowedOrigin 验证合规跨域 WebSocket 握手 Host 会对齐 Origin 以通过 Upgrader 检查。
func TestMiddlewareNormalizesWebSocketHostForAllowedOrigin(t *testing.T) {
	reached := false
	handler := New(Dependencies{Config: middlewareTestConfig()})(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		reached = true
		if request.Host != "admin.example.com" {
			t.Fatalf("request.Host = %q, want admin.example.com", request.Host)
		}
	}))
	request := httptest.NewRequest(http.MethodGet, "/graphql", nil)
	request.Host = "engine.internal:1980"
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "websocket")
	request.Header.Set("Origin", "https://admin.example.com")
	handler.ServeHTTP(httptest.NewRecorder(), request)
	if !reached {
		t.Fatal("handler not reached")
	}
}

func middlewareTestConfig() config.SecurityConfig {
	return config.SecurityConfig{AllowedOrigins: map[string]struct{}{"https://admin.example.com": {}}}
}
