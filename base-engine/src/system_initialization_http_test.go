/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package src_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	enginesrc "base-engine/src"
	"base-engine/src/dbup"
	httpmiddleware "base-engine/src/middleware"
	"base-engine/src/services/authentication"
	"gorm.io/gorm"
)

func TestSystemInitializationHTTPStatusAndCreation(t *testing.T) {
	t.Run("[Auth.Initialize.HTTP] 状态检查与首次初始化使用无缓存 JSON 契约", func(t *testing.T) {
		handler, db := systemInitializationHandler(t, 20)
		assertInitializationStatus(t, handler, false)

		response := performInitializationRequest(handler, "13800000000", "Correct-Horse-42", "Correct-Horse-42", "https://admin.example.com")
		if response.Code != http.StatusCreated || response.Header().Get("Content-Type") != "application/json" {
			t.Fatalf("create response: status=%d headers=%#v body=%s", response.Code, response.Header(), response.Body.String())
		}
		if bytes.Contains(response.Body.Bytes(), []byte("Correct-Horse-42")) || bytes.Contains(response.Body.Bytes(), []byte("13800000000")) {
			t.Fatalf("response leaked credentials: %s", response.Body.String())
		}
		assertInitializationStatus(t, handler, true)
		var count int64
		if err := db.Table("accounts").Count(&count).Error; err != nil || count != 1 {
			t.Fatalf("account count = %d, err = %v", count, err)
		}
	})
}

func TestSystemInitializationHTTPRejectsInvalidRequests(t *testing.T) {
	t.Run("[Auth.Initialize.HTTP] 非法来源与请求体不能进入初始化事务", func(t *testing.T) {
		tests := []struct {
			name        string
			body        string
			contentType string
			origin      string
			status      int
			code        string
		}{
			{name: "missing origin", body: `{}`, contentType: "application/json", status: http.StatusForbidden, code: "PERMISSION_DENIED"},
			{name: "unlisted origin", body: `{}`, contentType: "application/json", origin: "https://evil.example.com", status: http.StatusForbidden},
			{name: "non JSON", body: `{}`, contentType: "text/plain", origin: "https://admin.example.com", status: http.StatusUnsupportedMediaType, code: "VALIDATION_FAILED"},
			{name: "malformed JSON", body: `{`, contentType: "application/json", origin: "https://admin.example.com", status: http.StatusBadRequest, code: "VALIDATION_FAILED"},
			{name: "unknown field", body: `{"phone":"13800000000","password":"Correct-Horse-42","passwordConfirmation":"Correct-Horse-42","admin":true}`, contentType: "application/json", origin: "https://admin.example.com", status: http.StatusBadRequest, code: "VALIDATION_FAILED"},
			{name: "trailing JSON", body: `{"phone":"13800000000"}{}`, contentType: "application/json", origin: "https://admin.example.com", status: http.StatusBadRequest, code: "VALIDATION_FAILED"},
			{name: "oversized JSON", body: `{"phone":"` + strings.Repeat("1", 5000) + `"}`, contentType: "application/json", origin: "https://admin.example.com", status: http.StatusBadRequest, code: "VALIDATION_FAILED"},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				handler, _ := systemInitializationHandler(t, 20)
				response := performRawInitializationRequest(handler, test.body, test.contentType, test.origin)
				if response.Code != test.status {
					t.Fatalf("status = %d, want %d, body=%s", response.Code, test.status, response.Body.String())
				}
				if test.code != "" && initializationErrorCode(t, response) != test.code {
					t.Fatalf("code = %s, body=%s", initializationErrorCode(t, response), response.Body.String())
				}
			})
		}
	})
}

func TestSystemInitializationHTTPMapsBusinessErrorsAndRateLimit(t *testing.T) {
	t.Run("[Auth.Initialize.HTTP] 业务错误与限流映射为稳定状态码", func(t *testing.T) {
		handler, _ := systemInitializationHandler(t, 20)
		mismatch := performInitializationRequest(handler, "13800000000", "Correct-Horse-42", "Different-Horse-42", "https://admin.example.com")
		if mismatch.Code != http.StatusBadRequest || initializationErrorCode(t, mismatch) != "PASSWORD_CONFIRM_MISMATCH" {
			t.Fatalf("mismatch response: status=%d body=%s", mismatch.Code, mismatch.Body.String())
		}
		created := performInitializationRequest(handler, "13800000000", "Correct-Horse-42", "Correct-Horse-42", "https://admin.example.com")
		conflict := performInitializationRequest(handler, "13900000000", "Correct-Horse-42", "Correct-Horse-42", "https://admin.example.com")
		if created.Code != http.StatusCreated || conflict.Code != http.StatusConflict || initializationErrorCode(t, conflict) != "HQ_ALREADY_BOOTSTRAPPED" {
			t.Fatalf("creation/conflict responses: created=%d conflict=%d body=%s", created.Code, conflict.Code, conflict.Body.String())
		}

		t.Run("rate limit uses an isolated database", func(t *testing.T) {
			limitedHandler, _ := systemInitializationHandler(t, 1)
			first := performInitializationRequest(limitedHandler, "13800000000", "weak", "weak", "https://admin.example.com")
			second := performInitializationRequest(limitedHandler, "13800000000", "weak", "weak", "https://admin.example.com")
			if first.Code != http.StatusBadRequest || second.Code != http.StatusTooManyRequests || initializationErrorCode(t, second) != "RATE_LIMITED" {
				t.Fatalf("rate responses: first=%d second=%d body=%s", first.Code, second.Code, second.Body.String())
			}
		})
	})
}

func TestSystemInitializationHTTPHidesDatabaseErrors(t *testing.T) {
	t.Run("[Auth.Initialize.HTTP] 数据库错误只暴露稳定内部错误码", func(t *testing.T) {
		handler, db := systemInitializationHandler(t, 20)
		sqlDB, err := db.DB()
		if err != nil {
			t.Fatal(err)
		}
		if err := sqlDB.Close(); err != nil {
			t.Fatal(err)
		}
		request := httptest.NewRequest(http.MethodGet, "/api/system-initialization", nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusInternalServerError || initializationErrorCode(t, response) != "INTERNAL_ERROR" {
			t.Fatalf("response: status=%d body=%s", response.Code, response.Body.String())
		}
		if response.Header().Get("Cache-Control") != "no-store" || strings.Contains(response.Body.String(), "sql") {
			t.Fatalf("unsafe error response: headers=%#v body=%s", response.Header(), response.Body.String())
		}
	})
}

func systemInitializationHandler(t *testing.T, limit int) (http.Handler, *gorm.DB) {
	t.Helper()
	db := openResolverTestDB(t)
	if err := dbup.MigrateSecurityTables(db); err != nil {
		t.Fatal(err)
	}
	if err := dbup.InitRoles(db); err != nil {
		t.Fatal(err)
	}
	cfg := resolverTestConfig()
	router := mux.NewRouter()
	enginesrc.RegisterSystemInitializationRoutes(router, enginesrc.SystemInitializationDependencies{
		DB: db, Config: cfg,
		Limiter: authentication.NewLoginLimiter(limit, time.Minute), Now: time.Now,
	})
	return httpmiddleware.New(httpmiddleware.Dependencies{Config: cfg})(router), db
}

func assertInitializationStatus(t *testing.T, handler http.Handler, want bool) {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, "/api/system-initialization", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	var body struct {
		Initialized bool `json:"initialized"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusOK || body.Initialized != want || response.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("status response: code=%d body=%s headers=%#v", response.Code, response.Body.String(), response.Header())
	}
}

func performInitializationRequest(handler http.Handler, phone, password, confirmation, origin string) *httptest.ResponseRecorder {
	body, _ := json.Marshal(map[string]string{
		"phone": phone, "password": password, "passwordConfirmation": confirmation,
	})
	return performRawInitializationRequest(handler, string(body), "application/json", origin)
}

func performRawInitializationRequest(handler http.Handler, body, contentType, origin string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, "/api/system-initialization", strings.NewReader(body))
	request.RemoteAddr = "192.0.2.4:443"
	request.Header.Set("Content-Type", contentType)
	if origin != "" {
		request.Header.Set("Origin", origin)
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func initializationErrorCode(t *testing.T, response *httptest.ResponseRecorder) string {
	t.Helper()
	var body struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	return body.Code
}
