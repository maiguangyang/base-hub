/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package integration_test

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
	"base-engine/src/services/authorization"
	sessionservice "base-engine/src/services/session"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type membershipRole struct{ OperatorMembershipID, OperatorRoleID string }
type permissionRole struct{ PermissionID, OperatorRoleID string }
type membershipStore struct{ OperatorMembershipID, StoreID string }

func (membershipRole) TableName() string  { return "operator_membership_roles" }
func (permissionRole) TableName() string  { return "permission_roles" }
func (membershipStore) TableName() string { return "operator_membership_stores" }

type securityFixture struct {
	t          *testing.T
	db         *gorm.DB
	cfg        config.SecurityConfig
	publisher  sessionservice.Publisher
	resolver   *enginesrc.Resolver
	handler    http.Handler
	cookies    map[string]*http.Cookie
	principals map[string]*auth.WorkspacePrincipal
}

func newSecurityFixture(t *testing.T) *securityFixture {
	return newSecurityFixtureWithPublisher(t, sessionservice.NewPublisher())
}

func newSecurityFixtureWithPublisher(t *testing.T, publisher sessionservice.Publisher) *securityFixture {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	migrateSecurityFixture(t, db)
	cfg := integrationConfig()
	dependencies := enginesrc.NewDependencies(db, cfg, publisher)
	generatedDB := gen.NewDB(db)
	resolver := enginesrc.NewResolver(generatedDB, &gen.EventController{}, dependencies)
	router := gen.GetHTTPServeMux(enginesrc.New(generatedDB, &gen.EventController{}, dependencies), generatedDB)
	fixture := &securityFixture{t: t, db: db, cfg: cfg, publisher: publisher, resolver: resolver, handler: httpmiddleware.New(httpmiddleware.Dependencies{Config: cfg, PrincipalResolver: dependencies.Principal})(router), cookies: map[string]*http.Cookie{}, principals: map[string]*auth.WorkspacePrincipal{}}
	fixture.seed()
	return fixture
}

func migrateSecurityFixture(t *testing.T, db *gorm.DB) {
	t.Helper()
	models := []any{&gen.Account{}, &gen.Session{}, &gen.OperatorMembership{}, &gen.Organization{}, &gen.OperatorRole{}, &gen.Permission{}, &gen.Store{}, &gen.MembershipInvitation{}, &gen.AuditLog{}, &gen.FranchiseOpeningRecord{}, &gen.GlobalPaymentConfig{}, &gen.FranchisePaymentConfig{}, &gen.StorePaymentConfig{}, &membershipRole{}, &permissionRole{}, &membershipStore{}}
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
}

func (f *securityFixture) seed() {
	organizations := []gen.Organization{
		{ID: "org-hq", Code: "HQ", Name: "HQ", Type: gen.OrganizationTypeHeadquarters, Status: gen.OrganizationStatusActive},
		{ID: "org-a", Code: "A", Name: "A", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive},
		{ID: "org-b", Code: "B", Name: "B", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive},
	}
	accounts := []gen.Account{
		{ID: "account-shared", Phone: "13800000001", DisplayName: "Shared", Status: gen.AccountStatusActive, CredentialVersion: 1},
		{ID: "account-staff", Phone: "13800000002", DisplayName: "Staff", Status: gen.AccountStatusActive, CredentialVersion: 1},
	}
	roles := []gen.OperatorRole{
		{ID: "role-hq", Name: "HQ", Kind: gen.RoleKindHqSuperAdmin, OrganizationID: "org-hq"},
		{ID: "role-owner-a", Name: "Owner A", Kind: gen.RoleKindFranchiseOwner, OrganizationID: "org-a"},
		{ID: "role-custom-a", Name: "Staff A", Kind: gen.RoleKindCustom, OrganizationID: "org-a"},
		{ID: "role-owner-b", Name: "Owner B", Kind: gen.RoleKindFranchiseOwner, OrganizationID: "org-b"},
	}
	memberships := []gen.OperatorMembership{
		{ID: "membership-hq", AccountID: "account-shared", OrganizationID: "org-hq", Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores},
		{ID: "membership-owner-a", AccountID: "account-shared", OrganizationID: "org-a", Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores},
		{ID: "membership-staff-a", AccountID: "account-staff", OrganizationID: "org-a", Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores},
		{ID: "membership-owner-b", AccountID: "account-shared", OrganizationID: "org-b", Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores},
	}
	stores := []gen.Store{
		{ID: "store-hq", Code: "HQ1", Name: "HQ Store", Lifecycle: gen.StoreLifecycleActive, OrganizationID: "org-hq", ReviewedByAccountID: stringPointer("account-shared")},
		{ID: "store-a", Code: "A1", Name: "A Store", Lifecycle: gen.StoreLifecycleActive, OrganizationID: "org-a", ReviewedByAccountID: stringPointer("account-shared")},
		{ID: "store-b", Code: "B1", Name: "B Store", Lifecycle: gen.StoreLifecycleActive, OrganizationID: "org-b", ReviewedByAccountID: stringPointer("account-shared")},
	}
	invitations := []gen.MembershipInvitation{
		{ID: "invitation-a", MembershipID: "membership-staff-a", InvitedByAccountID: "account-shared", ExpiresAt: time.Now().Add(time.Hour)},
		{ID: "invitation-b", MembershipID: "membership-owner-b", InvitedByAccountID: "account-shared", ExpiresAt: time.Now().Add(time.Hour)},
	}
	audits := []gen.AuditLog{
		{ID: "audit-a", ActorAccountID: stringPointer("account-shared"), SessionID: stringPointer("session-a-1"), OrganizationID: stringPointer("org-a"), StoreID: stringPointer("store-a"), Action: "test:a", ResourceType: "test", ResultCode: "SUCCESS"},
		{ID: "audit-b", ActorAccountID: stringPointer("account-shared"), SessionID: stringPointer("session-b"), OrganizationID: stringPointer("org-b"), StoreID: stringPointer("store-b"), Action: "test:b", ResourceType: "test", ResultCode: "SUCCESS"},
	}
	f.createAll(organizations, accounts, roles, memberships, stores, invitations, audits)
	f.seedPermissions()
	f.createAll([]membershipRole{{"membership-hq", "role-hq"}, {"membership-owner-a", "role-owner-a"}, {"membership-staff-a", "role-custom-a"}, {"membership-owner-b", "role-owner-b"}})
	f.seedSessions()
}

func (f *securityFixture) seedPermissions() {
	tenant := []string{"operatorMembership:read", "operatorMembership:create", "operatorMembership:update", "operatorMembership:delete", "operatorRole:read", "operatorRole:create", "operatorRole:update", "operatorRole:delete", "store:read", "store:create", "store:update", "store:delete", "store:submit", "membershipInvitation:read", "membershipInvitation:create", "membershipInvitation:delete", "tenantAudit:read"}
	system := []string{"organization:read", "organization:suspend", "organization:restore", "store:read_all", "hqStore:read", "hqStore:create", "hqStore:update", "hqStore:delete", "hqMembership:read", "hqMembership:create", "hqMembership:update", "hqMembership:delete", "hqRole:read", "hqRole:create", "hqRole:update", "hqRole:delete", "permission:read", "account:read", "account:update", "auditLog:read", "franchise:provision"}
	permissions := make([]gen.Permission, 0, len(tenant)+len(system))
	joins := make([]permissionRole, 0, len(tenant)*3+len(system))
	for index, action := range append(tenant, system...) {
		id := "permission-" + strings.ReplaceAll(action, ":", "-")
		scope := gen.PermissionScopeTenant
		if index >= len(tenant) {
			scope = gen.PermissionScopeSystem
		}
		permissions = append(permissions, gen.Permission{ID: id, Name: action, Action: action, Module: "test", Scope: scope})
		if scope == gen.PermissionScopeTenant {
			for _, roleID := range []string{"role-owner-a", "role-custom-a", "role-owner-b"} {
				joins = append(joins, permissionRole{id, roleID})
			}
		} else {
			joins = append(joins, permissionRole{id, "role-hq"})
		}
	}
	f.createAll(permissions, joins)
}

func stringPointer(value string) *string { return &value }

func (f *securityFixture) seedSessions() {
	now := time.Now()
	definitions := []struct {
		id, account, organization string
		workspace                 gen.WorkspaceType
	}{
		{"session-a-1", "account-shared", "org-a", gen.WorkspaceTypeFranchise}, {"session-a-2", "account-staff", "org-a", gen.WorkspaceTypeFranchise},
		{"session-b", "account-shared", "org-b", gen.WorkspaceTypeFranchise}, {"session-hq", "account-shared", "org-hq", gen.WorkspaceTypeHeadquarters},
	}
	for _, item := range definitions {
		organizationID := item.organization
		session := gen.Session{ID: item.id, AccountID: item.account, OrganizationID: &organizationID, WorkspaceType: item.workspace, CredentialVersion: 1, ExpiresAt: now.Add(time.Hour), LastSeenAt: now}
		f.createAll([]gen.Session{session})
		principal := &auth.WorkspacePrincipal{SessionID: item.id, AccountID: item.account, OrganizationID: &organizationID, WorkspaceType: auth.WorkspaceType(item.workspace)}
		f.principals[item.id] = principal
		claims := auth.SessionClaims{SessionID: item.id, AccountID: item.account, OrganizationID: &organizationID, WorkspaceType: auth.WorkspaceType(item.workspace), CredentialVersion: 1}
		token, err := auth.SignSessionClaims(f.cfg, claims)
		if err != nil {
			f.t.Fatal(err)
		}
		f.cookies[item.id] = &http.Cookie{Name: auth.SessionCookieName, Value: token}
	}
}

func (f *securityFixture) createAll(groups ...any) {
	f.t.Helper()
	for _, group := range groups {
		if err := f.db.Create(group).Error; err != nil {
			f.t.Fatal(err)
		}
	}
}

func integrationConfig() config.SecurityConfig {
	return config.SecurityConfig{SigningKey: []byte(strings.Repeat("k", 32)), TokenIssuer: "integration", CookieSecure: true, AllowedOrigins: map[string]struct{}{"https://admin.example.com": {}}, SessionDuration: 12 * time.Hour, LoginLockThreshold: 5, LoginLockDuration: 15 * time.Minute, LoginIPAttemptsPerMinute: 20}
}

type graphQLResponse struct {
	Data   map[string]json.RawMessage `json:"data"`
	Errors []struct {
		Extensions map[string]any `json:"extensions"`
	} `json:"errors"`
	Body string
}

func (f *securityFixture) execute(sessionID, query string) graphQLResponse {
	f.t.Helper()
	encoded, _ := json.Marshal(map[string]string{"query": query})
	request := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(encoded))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "https://admin.example.com")
	request.AddCookie(f.cookies[sessionID])
	recorder := httptest.NewRecorder()
	f.handler.ServeHTTP(recorder, request)
	response := graphQLResponse{Body: recorder.Body.String()}
	_ = json.Unmarshal(recorder.Body.Bytes(), &response)
	return response
}

func assertCode(t *testing.T, response graphQLResponse, code auth.Code) {
	t.Helper()
	if len(response.Errors) == 0 || response.Errors[0].Extensions["code"] != string(code) {
		t.Fatalf("GraphQL code = %#v, want %s; body=%s", response.Errors, code, response.Body)
	}
}

func assertHandlerCoverage(t *testing.T) {
	t.Helper()
	defaults := gen.DefaultResolutionHandlers()
	registered := authorization.RegisterHandlers(defaults)
	registered.StoreMembers = defaults.StoreMembers
	if err := authorization.ValidateHandlerCoverage(registered, defaults); err == nil {
		t.Fatal("missing relationship override was not detected")
	}
}
