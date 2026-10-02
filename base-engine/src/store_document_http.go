package src

import (
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/src/services/ai"
	storeservice "base-engine/src/services/store"
)

type storeDocumentHTTPHandler struct {
	store    *storeservice.Service
	security config.SecurityConfig
}

type storeDocumentUploadResult struct {
	AttachmentID string `json:"attachmentId"`
}

func RegisterStoreDocumentRoutes(router *mux.Router, store *storeservice.Service, security config.SecurityConfig) {
	h := storeDocumentHTTPHandler{store: store, security: security}
	router.HandleFunc("/api/store-documents", h.upload).Methods(http.MethodPost).
		Name(ai.HTTPRouteName("http.storeDocument.upload", nil, storeDocumentUploadResult{}))
	router.HandleFunc("/uploads/stores/{storeId}/{filename}", h.read).Methods(http.MethodGet).
		Name(ai.HTTPRouteName("http.storeDocument.read", nil, nil))
}

func (h storeDocumentHTTPHandler) upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	if origin := strings.TrimSpace(r.Header.Get("Origin")); origin == "" || !h.security.OriginAllowed(origin) {
		storeDocumentError(w, http.StatusForbidden, auth.CodePermissionDenied)
		return
	}
	p, err := auth.RequirePrincipal(r.Context())
	if err != nil {
		storeDocumentError(w, http.StatusUnauthorized, auth.CodeAuthRequired)
		return
	}
	if h.store == nil {
		storeDocumentError(w, http.StatusServiceUnavailable, auth.CodeInternalError)
		return
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "image/") {
		storeDocumentError(w, http.StatusUnsupportedMediaType, auth.CodeValidationFailed)
		return
	}
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 5<<20))
	if err != nil {
		storeDocumentError(w, http.StatusBadRequest, auth.CodeValidationFailed)
		return
	}
	id, err := h.store.UploadDocument(r.Context(), p, data)
	if err != nil {
		storeDocumentServiceError(w, err)
		return
	}
	writeInitializationJSON(w, http.StatusCreated, storeDocumentUploadResult{AttachmentID: id})
}

func (h storeDocumentHTTPHandler) read(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if h.store == nil {
		http.NotFound(w, r)
		return
	}
	p, err := auth.RequirePrincipal(r.Context())
	if err != nil {
		http.NotFound(w, r)
		return
	}
	vars := mux.Vars(r)
	file, contentType, err := h.store.OpenDocument(r.Context(), p, vars["storeId"], vars["filename"])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", contentType)
	http.ServeContent(w, r, vars["filename"], stat.ModTime(), file)
}

func storeDocumentServiceError(w http.ResponseWriter, err error) {
	code := auth.ErrorCode(err)
	status := http.StatusBadRequest
	switch code {
	case auth.CodePermissionDenied:
		status = http.StatusForbidden
	case auth.CodeAuthRequired:
		status = http.StatusUnauthorized
	case auth.CodeInternalError:
		status = http.StatusInternalServerError
	}
	storeDocumentError(w, status, code)
}

func storeDocumentError(w http.ResponseWriter, status int, code auth.Code) {
	writeInitializationJSON(w, status, map[string]string{"code": string(code)})
}
