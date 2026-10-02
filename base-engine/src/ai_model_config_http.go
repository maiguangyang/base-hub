package src

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/src/services/ai"
	"base-engine/src/services/authentication"
)

const modelConfigBodyLimit = 4096

type modelConfigHTTPHandler struct {
	store        *ai.ModelConfigStore
	security     config.SecurityConfig
	probeLimiter *authentication.LoginLimiter
}

type modelConfigVersionInput struct {
	Version uint64 `json:"version"`
}

// RegisterAIModelConfigRoutes exposes the headquarters-only configuration workflow.
func RegisterAIModelConfigRoutes(router *mux.Router, store *ai.ModelConfigStore, security config.SecurityConfig) {
	h := modelConfigHTTPHandler{store: store, security: security, probeLimiter: authentication.NewLoginLimiter(3, time.Minute)}
	router.HandleFunc("/api/ai/model-config", h.status).Methods(http.MethodGet).Name(ai.HTTPRouteName("http.aiModelConfig.status", nil, ai.ModelConfigStatus{}))
	router.HandleFunc("/api/ai/model-config", h.save).Methods(http.MethodPut).Name(ai.HTTPRouteName("http.aiModelConfig.save", ai.ModelConfigUpdate{}, ai.ModelConfigStatus{}))
	router.HandleFunc("/api/ai/model-config/probe", h.probe).Methods(http.MethodPost).Name(ai.HTTPRouteName("http.aiModelConfig.probe", modelConfigVersionInput{}, ai.ModelConfigStatus{}))
	router.HandleFunc("/api/ai/model-config/activate", h.activate).Methods(http.MethodPost).Name(ai.HTTPRouteName("http.aiModelConfig.activate", modelConfigVersionInput{}, ai.ModelConfigStatus{}))
	router.HandleFunc("/api/ai/model-config/deactivate", h.deactivate).Methods(http.MethodPost).Name(ai.HTTPRouteName("http.aiModelConfig.deactivate", modelConfigVersionInput{}, ai.ModelConfigStatus{}))
}

func (h modelConfigHTTPHandler) status(w http.ResponseWriter, r *http.Request) {
	modelConfigHeaders(w)
	if !h.validOrigin(w, r, false) {
		return
	}
	p, ok := h.principal(w, r)
	if !ok {
		return
	}
	status, err := h.store.Status(r.Context(), p)
	h.respond(w, status, err)
}

func (h modelConfigHTTPHandler) save(w http.ResponseWriter, r *http.Request) {
	modelConfigHeaders(w)
	if !h.validWrite(w, r) {
		return
	}
	p, ok := h.principal(w, r)
	if !ok {
		return
	}
	var input ai.ModelConfigUpdate
	if !decodeModelConfigInput(w, r, &input) {
		return
	}
	status, err := h.store.Save(r.Context(), p, input)
	h.respond(w, status, err)
}

func (h modelConfigHTTPHandler) probe(w http.ResponseWriter, r *http.Request) {
	h.writeVersionAction(w, r, h.store.Probe, true)
}

func (h modelConfigHTTPHandler) activate(w http.ResponseWriter, r *http.Request) {
	h.writeVersionAction(w, r, h.store.Activate, false)
}

func (h modelConfigHTTPHandler) deactivate(w http.ResponseWriter, r *http.Request) {
	h.writeVersionAction(w, r, h.store.Deactivate, false)
}

func (h modelConfigHTTPHandler) writeVersionAction(w http.ResponseWriter, r *http.Request, action func(context.Context, *auth.WorkspacePrincipal, uint64) (ai.ModelConfigStatus, error), limited bool) {
	modelConfigHeaders(w)
	if !h.validWrite(w, r) {
		return
	}
	p, ok := h.principal(w, r)
	if !ok {
		return
	}
	var input modelConfigVersionInput
	if !decodeModelConfigInput(w, r, &input) {
		return
	}
	if input.Version == 0 {
		modelConfigError(w, http.StatusBadRequest, "VALIDATION_FAILED")
		return
	}
	if limited && !h.probeLimiter.Allow(p.AccountID, time.Now()) {
		modelConfigError(w, http.StatusTooManyRequests, "RATE_LIMITED")
		return
	}
	status, err := action(r.Context(), p, input.Version)
	h.respond(w, status, err)
}

func (h modelConfigHTTPHandler) validWrite(w http.ResponseWriter, r *http.Request) bool {
	if !h.validOrigin(w, r, true) {
		return false
	}
	if !isJSONContentType(r.Header.Get("Content-Type")) {
		modelConfigError(w, http.StatusUnsupportedMediaType, "VALIDATION_FAILED")
		return false
	}
	return true
}

func (h modelConfigHTTPHandler) validOrigin(w http.ResponseWriter, r *http.Request, required bool) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	_, allowed := h.security.AllowedOrigins[origin]
	if (required && origin == "") || (origin != "" && !allowed) {
		modelConfigError(w, http.StatusForbidden, "ORIGIN_NOT_ALLOWED")
		return false
	}
	return true
}

func (h modelConfigHTTPHandler) principal(w http.ResponseWriter, r *http.Request) (*auth.WorkspacePrincipal, bool) {
	p, err := auth.RequirePrincipal(r.Context())
	if err != nil {
		modelConfigError(w, http.StatusUnauthorized, string(auth.ErrorCode(err)))
		return nil, false
	}
	return p, true
}

func (h modelConfigHTTPHandler) respond(w http.ResponseWriter, status ai.ModelConfigStatus, err error) {
	if err != nil {
		httpStatus, code := modelConfigHTTPError(err)
		modelConfigError(w, httpStatus, code)
		return
	}
	writeInitializationJSON(w, http.StatusOK, status)
}

func decodeModelConfigInput(w http.ResponseWriter, r *http.Request, value any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, modelConfigBodyLimit))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		modelConfigError(w, http.StatusBadRequest, "VALIDATION_FAILED")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		modelConfigError(w, http.StatusBadRequest, "VALIDATION_FAILED")
		return false
	}
	return true
}

func modelConfigHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
}

func modelConfigError(w http.ResponseWriter, status int, code string) {
	writeInitializationJSON(w, status, map[string]string{"code": code})
}

func modelConfigHTTPError(err error) (int, string) {
	if status, ok := modelConfigUpstreamStatus(err); ok {
		return http.StatusBadGateway, modelUpstreamFailureCode(status)
	}
	return modelConfigBaseHTTPError(err)
}

func modelConfigBaseHTTPError(err error) (int, string) {
	switch {
	case errors.Is(err, ai.ErrModelConfigForbidden):
		return http.StatusForbidden, "PERMISSION_DENIED"
	case errors.Is(err, ai.ErrModelConfigInvalid):
		return http.StatusBadRequest, "VALIDATION_FAILED"
	case errors.Is(err, ai.ErrModelConfigConflict):
		return http.StatusConflict, "CONFLICT"
	case errors.Is(err, ai.ErrModelConfigUnavailable):
		return http.StatusServiceUnavailable, "CONFIG_UNAVAILABLE"
	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusGatewayTimeout, "MODEL_PROBE_TIMEOUT"
	case errors.Is(err, ai.ErrModelProtocol):
		return http.StatusBadGateway, "MODEL_PROTOCOL_UNSUPPORTED"
	case errors.Is(err, ai.ErrModelUnavailable):
		return http.StatusBadGateway, "MODEL_CONNECTION_FAILED"
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR"
	}
}

func modelConfigUpstreamStatus(err error) (int, bool) {
	var upstream interface{ UpstreamStatusCode() int }
	if !errors.Is(err, ai.ErrModelUnavailable) || !errors.As(err, &upstream) {
		return 0, false
	}
	return upstream.UpstreamStatusCode(), true
}

func modelUpstreamFailureCode(status int) string {
	switch status {
	case http.StatusUnauthorized, http.StatusForbidden:
		return "MODEL_UPSTREAM_AUTH"
	case http.StatusNotFound:
		return "MODEL_TARGET_NOT_FOUND"
	case http.StatusTooManyRequests:
		return "MODEL_UPSTREAM_RATE_LIMITED"
	}
	if status >= 500 {
		return "MODEL_UPSTREAM_UNAVAILABLE"
	}
	return "MODEL_REQUEST_REJECTED"
}
