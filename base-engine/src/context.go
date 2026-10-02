/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package src

import (
	"context"
	"net/http"
)

type requestContextKey int

const (
	requestIDKey requestContextKey = iota
	responseWriterKey
	remoteIPKey
	sessionTokenKey
)

// WithRequestMetadata 写入 GraphQL 请求所需的可信传输元数据。
func WithRequestMetadata(ctx context.Context, requestID string, writer http.ResponseWriter, remoteIP, token string) context.Context {
	ctx = context.WithValue(ctx, requestIDKey, requestID)
	ctx = context.WithValue(ctx, responseWriterKey, writer)
	ctx = context.WithValue(ctx, remoteIPKey, remoteIP)
	ctx = context.WithValue(ctx, sessionTokenKey, token)
	return ctx
}

// RequestIDFromContext 返回链路请求编号。
func RequestIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(requestIDKey).(string)
	return value
}

// ResponseWriterFromContext 返回用于设置认证 Cookie 的响应写入器。
func ResponseWriterFromContext(ctx context.Context) http.ResponseWriter {
	value, _ := ctx.Value(responseWriterKey).(http.ResponseWriter)
	return value
}

// RemoteIPFromContext 返回去除端口后的远端 IP。
func RemoteIPFromContext(ctx context.Context) string {
	value, _ := ctx.Value(remoteIPKey).(string)
	return value
}

// SessionTokenFromContext 返回会话 Cookie 中的签名令牌。
func SessionTokenFromContext(ctx context.Context) string {
	value, _ := ctx.Value(sessionTokenKey).(string)
	return value
}
