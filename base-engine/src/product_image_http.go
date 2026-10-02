package src

import (
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/src/services/ai"
	"base-engine/src/services/productcatalog"
)

type productImageHTTPHandler struct {
	service  *productcatalog.Service
	security config.SecurityConfig
}

type productImageUploadResult struct {
	AttachmentID string `json:"attachmentId"`
}

func RegisterProductImageRoutes(router *mux.Router, service *productcatalog.Service, security config.SecurityConfig) {
	h := productImageHTTPHandler{service: service, security: security}
	router.HandleFunc("/api/product-main-images", h.upload).Methods(http.MethodPost).
		Name(ai.HTTPRouteName("http.productMainImage.upload", nil, productImageUploadResult{}))
	router.HandleFunc("/uploads/products/{productId}/{filename}", h.read).Methods(http.MethodGet).
		Name(ai.HTTPRouteName("http.productMainImage.read", nil, nil))
}

func (h productImageHTTPHandler) upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" || !h.security.OriginAllowed(origin) {
		productImageError(w, http.StatusForbidden, auth.CodePermissionDenied)
		return
	}
	principal, err := auth.RequirePrincipal(r.Context())
	if err != nil {
		productImageError(w, http.StatusUnauthorized, auth.CodeAuthRequired)
		return
	}
	if h.service == nil {
		productImageError(w, http.StatusServiceUnavailable, auth.CodeInternalError)
		return
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "image/") {
		productImageError(w, http.StatusUnsupportedMediaType, auth.CodeValidationFailed)
		return
	}
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 5<<20))
	if err != nil {
		productImageError(w, http.StatusBadRequest, auth.CodeValidationFailed)
		return
	}
	id, err := h.service.UploadMainImage(r.Context(), principal, data)
	if err != nil {
		productImageServiceError(w, err)
		return
	}
	writeInitializationJSON(w, http.StatusCreated, productImageUploadResult{AttachmentID: id})
}

func (h productImageHTTPHandler) read(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if h.service == nil {
		http.NotFound(w, r)
		return
	}
	vars := mux.Vars(r)
	file, contentType, err := h.service.OpenMainImage(r.Context(), vars["productId"], vars["filename"])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "no-store")
	stat, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, vars["filename"], stat.ModTime(), file)
}

func productImageServiceError(w http.ResponseWriter, err error) {
	code := auth.ErrorCode(err)
	status := http.StatusInternalServerError
	switch code {
	case auth.CodeAuthRequired:
		status = http.StatusUnauthorized
	case auth.CodeWorkspaceForbidden, auth.CodePermissionDenied:
		status = http.StatusForbidden
	case auth.CodeValidationFailed:
		status = http.StatusBadRequest
	case auth.CodeConflict:
		status = http.StatusConflict
	default:
		code = auth.CodeInternalError
	}
	productImageError(w, status, code)
}

func productImageError(w http.ResponseWriter, status int, code auth.Code) {
	writeInitializationJSON(w, status, map[string]string{"code": string(code)})
}
