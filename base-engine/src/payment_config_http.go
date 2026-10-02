package src

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/src/services/ai"
	"base-engine/src/services/paymentconfig"
)

const paymentConfigBodyLimit = 64 << 10

type paymentConfigHTTPHandler struct {
	store    *paymentconfig.Store
	security config.SecurityConfig
}

func RegisterPaymentConfigRoutes(router *mux.Router, store *paymentconfig.Store, security config.SecurityConfig) {
	h := paymentConfigHTTPHandler{store: store, security: security}
	router.HandleFunc("/api/payment-config", h.read).Methods(http.MethodGet).Name(ai.HTTPRouteName("http.paymentConfig.read", nil, paymentconfig.ScopeView{}))
	router.HandleFunc("/api/payment-config", h.save).Methods(http.MethodPut).Name(ai.HTTPRouteName("http.paymentConfig.save", paymentconfig.SaveInput{}, paymentconfig.ScopeView{}))
	router.HandleFunc("/api/payment-config/state", h.state).Methods(http.MethodPost).Name(ai.HTTPRouteName("http.paymentConfig.state", paymentconfig.StateInput{}, paymentconfig.ScopeView{}))
	router.HandleFunc("/api/payment-config/restore-inheritance", h.restore).Methods(http.MethodPost).Name(ai.HTTPRouteName("http.paymentConfig.restoreInheritance", paymentconfig.ResetInput{}, paymentconfig.ScopeView{}))
}

func (h paymentConfigHTTPHandler) read(w http.ResponseWriter, r *http.Request) {
	p, ok := h.preflight(w, r, false)
	if !ok {
		return
	}
	ref := paymentconfig.ScopeRef{Scope: r.URL.Query().Get("scope"), OrganizationID: r.URL.Query().Get("organizationId"), StoreID: r.URL.Query().Get("storeId")}
	view, err := h.store.Read(r.Context(), p, ref)
	h.respond(w, view, err)
}

func (h paymentConfigHTTPHandler) save(w http.ResponseWriter, r *http.Request) {
	p, ok := h.preflight(w, r, true)
	if !ok {
		return
	}
	var input paymentconfig.SaveInput
	if !decodePaymentInput(w, r, &input) {
		return
	}
	view, err := h.store.Save(r.Context(), p, input)
	h.respond(w, view, err)
}

func (h paymentConfigHTTPHandler) state(w http.ResponseWriter, r *http.Request) {
	p, ok := h.preflight(w, r, true)
	if !ok {
		return
	}
	var input paymentconfig.StateInput
	if !decodePaymentInput(w, r, &input) {
		return
	}
	view, err := h.store.SetState(r.Context(), p, input)
	h.respond(w, view, err)
}

func (h paymentConfigHTTPHandler) restore(w http.ResponseWriter, r *http.Request) {
	p, ok := h.preflight(w, r, true)
	if !ok {
		return
	}
	var input paymentconfig.ResetInput
	if !decodePaymentInput(w, r, &input) {
		return
	}
	view, err := h.store.RestoreInheritance(r.Context(), p, input)
	h.respond(w, view, err)
}

func (h paymentConfigHTTPHandler) preflight(w http.ResponseWriter, r *http.Request, write bool) (*auth.WorkspacePrincipal, bool) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	_, allowed := h.security.AllowedOrigins[origin]
	if (write && origin == "") || (origin != "" && !allowed) {
		paymentHTTPError(w, http.StatusForbidden, "ORIGIN_NOT_ALLOWED")
		return nil, false
	}
	principal, err := auth.RequirePrincipal(r.Context())
	if err != nil {
		paymentHTTPError(w, http.StatusUnauthorized, string(auth.ErrorCode(err)))
		return nil, false
	}
	if h.store == nil {
		paymentHTTPError(w, http.StatusServiceUnavailable, "CONFIG_UNAVAILABLE")
		return nil, false
	}
	if write && !isJSONContentType(r.Header.Get("Content-Type")) {
		paymentHTTPError(w, http.StatusUnsupportedMediaType, "VALIDATION_FAILED")
		return nil, false
	}
	return principal, true
}

func decodePaymentInput(w http.ResponseWriter, r *http.Request, value any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, paymentConfigBodyLimit))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		paymentHTTPError(w, http.StatusBadRequest, "VALIDATION_FAILED")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		paymentHTTPError(w, http.StatusBadRequest, "VALIDATION_FAILED")
		return false
	}
	return true
}

func (h paymentConfigHTTPHandler) respond(w http.ResponseWriter, view paymentconfig.ScopeView, err error) {
	if err != nil {
		status, code := paymentConfigErrorStatus(err)
		paymentHTTPError(w, status, code)
		return
	}
	writeInitializationJSON(w, http.StatusOK, view)
}

func paymentHTTPError(w http.ResponseWriter, status int, code string) {
	writeInitializationJSON(w, status, map[string]string{"code": code})
}

func paymentConfigErrorStatus(err error) (int, string) {
	switch {
	case errors.Is(err, paymentconfig.ErrForbidden):
		return http.StatusForbidden, "PERMISSION_DENIED"
	case errors.Is(err, paymentconfig.ErrInvalid):
		return http.StatusBadRequest, "VALIDATION_FAILED"
	case errors.Is(err, paymentconfig.ErrConflict):
		return http.StatusConflict, "CONFLICT"
	case errors.Is(err, paymentconfig.ErrUnavailable):
		return http.StatusServiceUnavailable, "CONFIG_UNAVAILABLE"
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR"
	}
}
