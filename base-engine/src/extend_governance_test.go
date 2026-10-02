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
	"gorm.io/gorm"
)

// TestGraphQLGovernanceAuthorization 验证治理 mutation 只能由具备独立权限的总部工作台执行。
func TestGraphQLGovernanceAuthorization(t *testing.T) {
	database := openResolverTestDB(t)
	cfg := resolverTestConfig()
	seedResolverPermission(t, database, "permission-provision", "franchise:provision", gen.PermissionScopeSystem)
	seedResolverPermission(t, database, "permission-organization-create", "organization:create", gen.PermissionScopeSystem)
	hqCookie := seedResolverWorkspace(t, database, cfg, "hq", gen.OrganizationTypeHeadquarters, gen.RoleKindHqSuperAdmin)
	franchiseCookie := seedResolverWorkspace(t, database, cfg, "franchise", gen.OrganizationTypeFranchise, gen.RoleKindFranchiseOwner)
	handler := governanceHandler(database, cfg)

	response := executeResolverGraphQL(t, handler, hqCookie, `mutation {
		provisionFranchise(input: {code: "F001", name: "First", ownerPhone: "13800000001", ownerDisplayName: "Owner"}) {
			organization { id code status }
			membership { id status }
			temporaryPassword
			invitationPending
		}
	}`)
	if response.Code != http.StatusOK || strings.Contains(response.Body.String(), `"errors"`) || !strings.Contains(response.Body.String(), `"temporaryPassword":"`) {
		t.Fatalf("HQ provision response: status=%d body=%s", response.Code, response.Body.String())
	}

	denied := executeResolverGraphQL(t, handler, franchiseCookie, `mutation {
		provisionFranchise(input: {code: "F002", name: "Second", ownerPhone: "13800000002", ownerDisplayName: "Owner"}) { organization { id } }
	}`)
	assertGraphQLErrorCode(t, denied, auth.CodePermissionDenied)

	direct := executeResolverGraphQL(t, handler, hqCookie, `mutation {
		createOrganization(input: {code: "BYPASS", name: "Bypass", type: FRANCHISE, status: ACTIVE}) { id }
	}`)
	assertGraphQLErrorCode(t, direct, auth.CodePermissionDenied)
}

// TestGraphQLTemporaryPasswordResetIsAudited 验证重置事务同时撤销会话并写入无明文审计。
func TestGraphQLTemporaryPasswordResetIsAudited(t *testing.T) {
	database := openResolverTestDB(t)
	cfg := resolverTestConfig()
	seedResolverPermission(t, database, "permission-account-update", "account:update", gen.PermissionScopeSystem)
	hqCookie := seedResolverWorkspace(t, database, cfg, "hq", gen.OrganizationTypeHeadquarters, gen.RoleKindHqSuperAdmin)
	seedResolverWorkspace(t, database, cfg, "franchise", gen.OrganizationTypeFranchise, gen.RoleKindFranchiseOwner)
	seedResetCredential(t, database, "hq-account")
	seedResetCredential(t, database, "franchise-account")
	credential := seedHQPasswordResetTarget(t, database)
	handler := governanceHandler(database, cfg)
	response := executeResolverGraphQL(t, governanceHandler(database, cfg), hqCookie, `mutation {
		resetTemporaryPassword(accountId: "hq-target-account") { accountId temporaryPassword }
	}`)
	if response.Code != http.StatusOK || strings.Contains(response.Body.String(), `"errors"`) {
		t.Fatalf("password reset failed: %s", response.Body.String())
	}
	var auditLog gen.AuditLog
	if err := database.Where("action = ?", "hqAdministrator:reset_password").First(&auditLog).Error; err != nil {
		t.Fatalf("password reset audit missing: %v", err)
	}
	var session gen.Session
	if err := database.First(&session, "id = ?", "hq-target-session").Error; err != nil || session.RevokedAt == nil {
		t.Fatalf("target session was not revoked: %#v err=%v", session, err)
	}
	assertResetTemporaryPasswordExpiry(t, database, &credential, "hq-target-account")
	foreign := executeResolverGraphQL(t, handler, hqCookie, `mutation {
		resetTemporaryPassword(accountId: "franchise-account") { accountId }
	}`)
	assertGraphQLErrorCode(t, foreign, auth.CodePermissionDenied)
	self := executeResolverGraphQL(t, handler, hqCookie, `mutation {
		resetTemporaryPassword(accountId: "hq-account") { accountId }
	}`)
	assertGraphQLErrorCode(t, self, auth.CodeSelfMembershipChangeDenied)
}

func seedResetCredential(t *testing.T, database *gorm.DB, accountID string) {
	t.Helper()
	hash, err := auth.HashPassword("Previous-Horse-42")
	if err != nil {
		t.Fatal(err)
	}
	credential := authentication.AccountCredential{AccountID: accountID, PasswordHash: hash, PasswordChangedAt: time.Now(), UpdatedAt: time.Now()}
	if err := database.Create(&credential).Error; err != nil {
		t.Fatal(err)
	}
}

func seedHQPasswordResetTarget(t *testing.T, database *gorm.DB) authentication.AccountCredential {
	t.Helper()
	hash, err := auth.HashPassword("Previous-Horse-42")
	if err != nil {
		t.Fatal(err)
	}
	organizationID := "hq-organization"
	values := []any{
		&gen.Account{ID: "hq-target-account", Phone: "13900000009", DisplayName: "Target", Status: gen.AccountStatusActive, CredentialVersion: 1},
		&gen.OperatorMembership{ID: "hq-target-membership", AccountID: "hq-target-account", OrganizationID: organizationID, Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores},
		&gen.OperatorRole{ID: "hq-target-role", Name: "Target", Kind: gen.RoleKindCustom, OrganizationID: organizationID},
		&resolverMembershipRole{OperatorMembershipID: "hq-target-membership", OperatorRoleID: "hq-target-role"},
		&gen.Session{ID: "hq-target-session", AccountID: "hq-target-account", OrganizationID: &organizationID, WorkspaceType: gen.WorkspaceTypeHeadquarters, CredentialVersion: 1, ExpiresAt: time.Now().Add(time.Hour), LastSeenAt: time.Now()},
	}
	for _, value := range values {
		if err := database.Create(value).Error; err != nil {
			t.Fatal(err)
		}
	}
	credential := authentication.AccountCredential{AccountID: "hq-target-account", PasswordHash: hash, PasswordChangedAt: time.Now(), UpdatedAt: time.Now()}
	if err := database.Create(&credential).Error; err != nil {
		t.Fatal(err)
	}
	return credential
}

func assertResetTemporaryPasswordExpiry(t *testing.T, database *gorm.DB, credential *authentication.AccountCredential, accountID string) {
	t.Helper()
	err := database.First(credential, "account_id = ?", accountID).Error
	if err != nil || credential.TemporaryPasswordExpiresAt == nil || !credential.TemporaryPasswordExpiresAt.After(time.Now()) {
		t.Fatalf("temporary password expiry missing: %#v err=%v", credential, err)
	}
}

func governanceHandler(database *gorm.DB, cfg config.SecurityConfig) http.Handler {
	publisher := sessionservice.NewPublisher()
	dependencies := enginesrc.NewDependencies(database, cfg, publisher)
	generatedDB := gen.NewDB(database)
	router := gen.GetHTTPServeMux(enginesrc.New(generatedDB, &gen.EventController{}, dependencies), generatedDB)
	return httpmiddleware.New(httpmiddleware.Dependencies{Config: cfg, PrincipalResolver: dependencies.Principal})(router)
}

func seedResolverPermission(t *testing.T, db *gorm.DB, id, action string, scope gen.PermissionScope) {
	t.Helper()
	permission := gen.Permission{ID: id, Name: action, Action: action, Module: "test", Scope: scope}
	if err := db.Create(&permission).Error; err != nil {
		t.Fatal(err)
	}
}

func seedResolverWorkspace(t *testing.T, db *gorm.DB, cfg config.SecurityConfig, id string, organizationType gen.OrganizationType, roleKind gen.RoleKind) *http.Cookie {
	t.Helper()
	organizationID := id + "-organization"
	accountID := id + "-account"
	membershipID := id + "-membership"
	roleID := id + "-role"
	workspaceType := gen.WorkspaceTypeFranchise
	if organizationType == gen.OrganizationTypeHeadquarters {
		workspaceType = gen.WorkspaceTypeHeadquarters
	}
	values := []any{
		&gen.Organization{ID: organizationID, Code: strings.ToUpper(id), Name: id, Type: organizationType, Status: gen.OrganizationStatusActive},
		&gen.Account{ID: accountID, Phone: "13" + id, DisplayName: id, Status: gen.AccountStatusActive, CredentialVersion: 1},
		&gen.OperatorMembership{ID: membershipID, AccountID: accountID, OrganizationID: organizationID, Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores},
		&gen.OperatorRole{ID: roleID, Name: id, Kind: roleKind, OrganizationID: organizationID},
		&resolverMembershipRole{OperatorMembershipID: membershipID, OperatorRoleID: roleID},
		&gen.Session{ID: id + "-session", AccountID: accountID, OrganizationID: &organizationID, WorkspaceType: workspaceType, CredentialVersion: 1, ExpiresAt: time.Now().Add(time.Hour), LastSeenAt: time.Now()},
	}
	for _, value := range values {
		if err := db.Create(value).Error; err != nil {
			t.Fatal(err)
		}
	}
	claims := auth.SessionClaims{SessionID: id + "-session", AccountID: accountID, WorkspaceType: auth.WorkspaceType(workspaceType), OrganizationID: &organizationID, CredentialVersion: 1}
	token, err := auth.SignSessionClaims(cfg, claims)
	if err != nil {
		t.Fatal(err)
	}
	return &http.Cookie{Name: auth.SessionCookieName, Value: token}
}

func executeResolverGraphQL(t *testing.T, handler http.Handler, cookie *http.Cookie, query string) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "https://admin.example.com")
	request.AddCookie(cookie)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func assertGraphQLErrorCode(t *testing.T, response *httptest.ResponseRecorder, code auth.Code) {
	t.Helper()
	var payload struct {
		Errors []struct {
			Extensions map[string]any `json:"extensions"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.Errors) == 0 || payload.Errors[0].Extensions["code"] != string(code) {
		t.Fatalf("GraphQL error code = %#v, want %s; body=%s", payload.Errors, code, response.Body.String())
	}
}
