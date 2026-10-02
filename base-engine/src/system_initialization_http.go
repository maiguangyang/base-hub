/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package src

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/src/dbup"
	"base-engine/src/services/ai"
	"base-engine/src/services/authentication"
	"gorm.io/gorm"
)

const systemInitializationBodyLimit = 4096

// SystemInitializationDependencies 描述一次性系统初始化 HTTP 边界所需依赖。
type SystemInitializationDependencies struct {
	DB      *gorm.DB
	Config  config.SecurityConfig
	Limiter *authentication.LoginLimiter
	Now     func() time.Time
}

type systemInitializationInput struct {
	Phone                string `json:"phone"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

type systemInitializationStatus struct {
	Initialized bool `json:"initialized"`
}

// RegisterSystemInitializationRoutes 在现有路由器上注册初始化状态与写入接口。
func RegisterSystemInitializationRoutes(router *mux.Router, dependencies SystemInitializationDependencies) {
	if dependencies.Now == nil {
		dependencies.Now = time.Now
	}
	if dependencies.Limiter == nil {
		limit := dependencies.Config.LoginIPAttemptsPerMinute
		if limit < 1 {
			limit = 20
		}
		dependencies.Limiter = authentication.NewLoginLimiter(limit, time.Minute)
	}
	handler := systemInitializationHandler{dependencies: dependencies}
	router.HandleFunc("/api/system-initialization", handler.status).Methods(http.MethodGet).Name(ai.HTTPRouteName("http.systemInitialization.status", nil, systemInitializationStatus{}))
	router.HandleFunc("/api/system-initialization", handler.initialize).Methods(http.MethodPost).Name(ai.HTTPRouteName("http.systemInitialization.initialize", systemInitializationInput{}, systemInitializationStatus{}))
}

type systemInitializationHandler struct {
	dependencies SystemInitializationDependencies
}

func (h systemInitializationHandler) status(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Cache-Control", "no-store")
	initialized, err := dbup.HQInitialized(request.Context(), h.dependencies.DB)
	if err != nil {
		logInitialization(request, auth.CodeInternalError)
		writeInitializationJSON(writer, http.StatusInternalServerError, map[string]any{"code": "INTERNAL_ERROR"})
		return
	}
	writeInitializationJSON(writer, http.StatusOK, systemInitializationStatus{Initialized: initialized})
}

func (h systemInitializationHandler) initialize(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Cache-Control", "no-store")
	if !h.validOrigin(request) {
		h.writeError(writer, request, http.StatusForbidden, auth.CodePermissionDenied)
		return
	}
	if !isJSONContentType(request.Header.Get("Content-Type")) {
		h.writeError(writer, request, http.StatusUnsupportedMediaType, auth.CodeValidationFailed)
		return
	}
	if !h.dependencies.Limiter.Allow(RemoteIPFromContext(request.Context()), h.dependencies.Now()) {
		h.writeError(writer, request, http.StatusTooManyRequests, auth.CodeRateLimited)
		return
	}
	input, err := decodeSystemInitializationInput(writer, request)
	if err != nil {
		h.writeError(writer, request, http.StatusBadRequest, auth.CodeValidationFailed)
		return
	}
	err = dbup.InitializeHQ(request.Context(), h.dependencies.DB, input.Phone, input.Password, input.PasswordConfirmation)
	if err != nil {
		status, code := initializationErrorResponse(err)
		h.writeError(writer, request, status, code)
		return
	}
	logInitialization(request, "SUCCESS")
	writeInitializationJSON(writer, http.StatusCreated, systemInitializationStatus{Initialized: true})
}

func (h systemInitializationHandler) validOrigin(request *http.Request) bool {
	origin := strings.TrimSpace(request.Header.Get("Origin"))
	return origin != "" && h.dependencies.Config.OriginAllowed(origin)
}

func (h systemInitializationHandler) writeError(writer http.ResponseWriter, request *http.Request, status int, code auth.Code) {
	logInitialization(request, code)
	writeInitializationJSON(writer, status, map[string]any{"code": code})
}

func decodeSystemInitializationInput(writer http.ResponseWriter, request *http.Request) (systemInitializationInput, error) {
	input := systemInitializationInput{}
	decoder := json.NewDecoder(http.MaxBytesReader(writer, request.Body, systemInitializationBodyLimit))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return input, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return input, errors.New("INVALID_JSON")
	}
	return input, nil
}

func isJSONContentType(value string) bool {
	mediaType, _, err := mime.ParseMediaType(value)
	return err == nil && mediaType == "application/json"
}

func initializationErrorResponse(err error) (int, auth.Code) {
	switch code := auth.ErrorCode(err); code {
	case auth.CodeValidationFailed, auth.CodePasswordWeak, auth.CodePasswordConfirmMismatch:
		return http.StatusBadRequest, code
	case auth.CodeHQAlreadyBootstrapped:
		return http.StatusConflict, code
	default:
		return http.StatusInternalServerError, auth.CodeInternalError
	}
}

func logInitialization(request *http.Request, result any) {
	log.Printf("request_id=%s action=system_initialize result=%v", RequestIDFromContext(request.Context()), result)
}

func writeInitializationJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
