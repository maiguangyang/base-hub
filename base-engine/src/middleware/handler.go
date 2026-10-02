/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package middleware

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/src"
)

// PrincipalResolver 描述中间件所需的权威身份解析能力。
type PrincipalResolver interface {
	Resolve(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error)
}

// Dependencies 是 HTTP 与 WebSocket 握手共享的安全依赖。
type Dependencies struct {
	Config            config.SecurityConfig
	PrincipalResolver PrincipalResolver
	Now               func() time.Time
	NewRequestID      func() string
}

// New 创建严格来源校验与请求上下文注入中间件。
func New(dependencies Dependencies) func(http.Handler) http.Handler {
	applyMiddlewareDefaults(&dependencies)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.URL.Path == "/automigrate" {
				http.NotFound(writer, request)
				return
			}
			if request.URL.Path == "/graphql" {
				writer.Header().Set("Cache-Control", "no-store")
			}
			if !applyOriginPolicy(writer, request, dependencies.Config) {
				return
			}
			if request.Method == http.MethodOptions {
				writer.WriteHeader(http.StatusNoContent)
				return
			}
			if isWebSocketUpgrade(request) {
				normalizeWebSocketHost(request)
			}
			token := sessionCookieToken(request)
			ctx := src.WithRequestMetadata(
				request.Context(), dependencies.NewRequestID(), writer,
				remoteIP(request.RemoteAddr), token,
			)
			if token != "" && dependencies.PrincipalResolver != nil {
				principal, err := dependencies.PrincipalResolver.Resolve(ctx, token, dependencies.Now())
				if err == nil {
					principal = principalWithRequestID(ctx, principal)
					ctx = auth.WithPrincipal(ctx, principal)
				} else {
					ctx = auth.WithAuthenticationError(ctx, err)
				}
				if isWebSocketUpgrade(request) {
					ctx = auth.WithPrincipalResolver(ctx, func(resolveContext context.Context) (*auth.WorkspacePrincipal, error) {
						resolved, resolveErr := dependencies.PrincipalResolver.Resolve(resolveContext, token, dependencies.Now())
						return principalWithRequestID(resolveContext, resolved), resolveErr
					})
				}
			}
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

func principalWithRequestID(ctx context.Context, principal *auth.WorkspacePrincipal) *auth.WorkspacePrincipal {
	if principal != nil {
		principal.RequestID = src.RequestIDFromContext(ctx)
	}
	return principal
}

func isWebSocketUpgrade(request *http.Request) bool {
	if !strings.EqualFold(strings.TrimSpace(request.Header.Get("Upgrade")), "websocket") {
		return false
	}
	for _, token := range strings.Split(request.Header.Get("Connection"), ",") {
		if strings.EqualFold(strings.TrimSpace(token), "upgrade") {
			return true
		}
	}
	return false
}

func applyMiddlewareDefaults(dependencies *Dependencies) {
	if dependencies.Now == nil {
		dependencies.Now = time.Now
	}
	if dependencies.NewRequestID == nil {
		dependencies.NewRequestID = func() string { return uuid.Must(uuid.NewV4()).String() }
	}
}

func applyOriginPolicy(writer http.ResponseWriter, request *http.Request, cfg config.SecurityConfig) bool {
	origin := strings.TrimSpace(request.Header.Get("Origin"))
	if origin != "" && !cfg.OriginAllowed(origin) {
		http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return false
	}
	if origin != "" {
		writer.Header().Set("Access-Control-Allow-Origin", origin)
		writer.Header().Set("Vary", "Origin")
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
	return true
}

func sessionCookieToken(request *http.Request) string {
	cookie, err := request.Cookie(auth.SessionCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func remoteIP(remoteAddress string) string {
	host, _, err := net.SplitHostPort(remoteAddress)
	if err == nil {
		return host
	}
	return remoteAddress
}

func normalizeWebSocketHost(request *http.Request) {
	origin := strings.TrimSpace(request.Header.Get("Origin"))
	if origin == "" {
		return
	}
	parsed, err := url.Parse(origin)
	if err == nil && parsed.Host != "" {
		request.Host = parsed.Host
	}
}
