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

	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	enginesrc "base-engine/src"
	httpmiddleware "base-engine/src/middleware"
	"base-engine/src/services/authentication"
	sessionservice "base-engine/src/services/session"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestGraphQLLoginCreatesSecureCookie 验证真实 gqlgen mutation 返回发现区 Viewer 和安全 Cookie。
func TestGraphQLLoginCreatesSecureCookie(t *testing.T) {
	database := openResolverTestDB(t)
	seedResolverAccount(t, database)
	cfg := resolverTestConfig()
	publisher := sessionservice.NewPublisher()
	dependencies := enginesrc.NewDependencies(database, cfg, publisher)
	generatedDB := gen.NewDB(database)
	router := gen.GetHTTPServeMux(enginesrc.New(generatedDB, &gen.EventController{}, dependencies), generatedDB)
	handler := httpmiddleware.New(httpmiddleware.Dependencies{Config: cfg, PrincipalResolver: dependencies.Principal})(router)
	body := map[string]any{"query": `mutation { login(input: {phone: "13800000000", password: "Correct-Horse-42"}) { requiresPasswordChange viewer { account { id phone } permissions } } }`}
	encoded, _ := json.Marshal(body)
	request := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(encoded))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "https://admin.example.com")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK || strings.Contains(response.Body.String(), `"errors"`) {
		t.Fatalf("GraphQL login failed: status=%d body=%s", response.Code, response.Body.String())
	}
	if response.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("GraphQL response must not be cached: %#v", response.Header())
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != auth.SessionCookieName || !cookies[0].HttpOnly {
		t.Fatalf("invalid session cookie: %#v", cookies)
	}
}

// TestGeneratedAutomigrateRouteIsDisabled 验证生成器迁移入口不能通过运行时服务访问。
func TestGeneratedAutomigrateRouteIsDisabled(t *testing.T) {
	database := openResolverTestDB(t)
	cfg := resolverTestConfig()
	dependencies := enginesrc.NewDependencies(database, cfg, sessionservice.NewPublisher())
	generatedDB := gen.NewDB(database)
	router := gen.GetHTTPServeMux(enginesrc.New(generatedDB, &gen.EventController{}, dependencies), generatedDB)
	handler := httpmiddleware.New(httpmiddleware.Dependencies{Config: cfg, PrincipalResolver: dependencies.Principal})(router)
	request := httptest.NewRequest(http.MethodGet, "/automigrate", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("automigrate route status=%d, want %d", response.Code, http.StatusNotFound)
	}
}

func openResolverTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	models := []any{
		&gen.Account{}, &gen.Session{}, &gen.OperatorMembership{}, &gen.Organization{}, &gen.FranchiseOpeningRecord{},
		&gen.OperatorRole{}, &gen.Permission{}, &gen.Store{}, &gen.MembershipInvitation{},
		&gen.GlobalPaymentConfig{}, &gen.FranchisePaymentConfig{}, &gen.StorePaymentConfig{},
		&gen.AuditLog{}, &authentication.AccountCredential{}, &resolverMembershipRole{},
		&resolverPermissionRole{}, &resolverMembershipStore{},
	}
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	return db
}

type resolverMembershipRole struct{ OperatorMembershipID, OperatorRoleID string }

func (resolverMembershipRole) TableName() string { return "operator_membership_roles" }

type resolverPermissionRole struct{ PermissionID, OperatorRoleID string }

func (resolverPermissionRole) TableName() string { return "permission_roles" }

type resolverMembershipStore struct{ OperatorMembershipID, StoreID string }

func (resolverMembershipStore) TableName() string { return "operator_membership_stores" }

func seedResolverAccount(t *testing.T, db *gorm.DB) {
	t.Helper()
	hash, err := auth.HashPassword("Correct-Horse-42")
	if err != nil {
		t.Fatal(err)
	}
	account := gen.Account{ID: "account-1", Phone: "13800000000", DisplayName: "Owner", Status: gen.AccountStatusActive, MustChangePassword: true, CredentialVersion: 1}
	now := time.Now()
	credential := authentication.AccountCredential{
		AccountID: account.ID, PasswordHash: hash,
		TemporaryPasswordExpiresAt: authentication.NewTemporaryPasswordExpiry(now), PasswordChangedAt: now, UpdatedAt: now,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&credential).Error; err != nil {
		t.Fatal(err)
	}
}

func resolverTestConfig() config.SecurityConfig {
	return config.SecurityConfig{
		SigningKey: []byte(strings.Repeat("k", 32)), TokenIssuer: "test",
		CookieSecure: true, AllowedOrigins: map[string]struct{}{"https://admin.example.com": {}},
		SessionDuration: 12 * time.Hour, LoginLockThreshold: 5,
		LoginLockDuration: 15 * time.Minute, LoginIPAttemptsPerMinute: 20,
	}
}
