package ai

import (
	"context"
	"errors"
	"net/http"
)

func safeRunErrorCode(err error) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "AI_RUN_TIMEOUT"
	}
	if errors.Is(err, ErrModelProtocol) {
		return "AI_MODEL_PROTOCOL"
	}
	var upstream interface{ UpstreamStatusCode() int }
	if errors.As(err, &upstream) {
		return safeUpstreamErrorCode(upstream.UpstreamStatusCode())
	}
	if errors.Is(err, ErrModelUnavailable) || errors.Is(err, errModelUnavailable) {
		return "AI_MODEL_UNAVAILABLE"
	}
	for cause := err; cause != nil; cause = errors.Unwrap(cause) {
		switch cause.Error() {
		case "AI_TOKEN_BUDGET_EXCEEDED", "AI_MODEL_CALL_LIMIT", "AI_TOOL_CALL_LIMIT", "AI_MODEL_USAGE_MISSING", "AI_MODEL_VERSION_CHANGED", "AI_SESSION_REVOKED", "VALIDATION_FAILED", "PERMISSION_DENIED", "CONFLICT", "CONFIG_UNAVAILABLE":
			return cause.Error()
		}
	}
	return "AI_RUN_FAILED"
}

func safeUpstreamErrorCode(status int) string {
	switch status {
	case http.StatusTooManyRequests:
		return "AI_MODEL_RATE_LIMIT"
	case http.StatusUnauthorized, http.StatusForbidden:
		return "AI_MODEL_AUTH_FAILED"
	default:
		return "AI_MODEL_UNAVAILABLE"
	}
}
