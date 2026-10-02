/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"context"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"base-engine/src/dbup"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestResolveWorkspacePrincipal 验证 Session、组织、成员、角色、权限与门店范围均从数据库解析。
func TestResolveWorkspacePrincipal(t *testing.T) {
	fixture := newPrincipalFixture(t)
	t.Run("missing cookie", fixture.testMissingCookie)
	t.Run("discovery", fixture.testDiscovery)
	t.Run("owner dynamic tenant permissions", fixture.testOwnerPermissions)
	t.Run("hq dynamic system permissions", fixture.testHQPermissions)
	t.Run("custom exact union and live selected stores", fixture.testCustomScope)
	fixture.assertStateErrors(t)
}

func (f *principalFixture) testMissingCookie(t *testing.T) {
	_, err := f.resolver.Resolve(context.Background(), "", f.now)
	assertCode(t, err, auth.CodeAuthRequired)
}

func (f *principalFixture) testDiscovery(t *testing.T) {
	principal := f.resolve(t, "discovery-session")
	if principal.WorkspaceType != auth.WorkspaceTypeDiscovery {
		t.Fatalf("workspace = %s", principal.WorkspaceType)
	}
	if len(principal.Permissions) != 0 {
		t.Fatalf("discovery permissions: %#v", principal.Permissions)
	}
}

func (f *principalFixture) testOwnerPermissions(t *testing.T) {
	principal := f.resolve(t, "owner-session")
	if !principal.Has("store:create") || !principal.Has("store:submit") || principal.Has("account:create") {
		t.Fatalf("invalid owner permissions: %#v", principal.Permissions)
	}
}

func (f *principalFixture) testHQPermissions(t *testing.T) {
	principal := f.resolve(t, "hq-session")
	if !principal.Has("account:create") || !principal.Has("store:approve") || principal.Has("store:submit") {
		t.Fatalf("invalid HQ permissions: %#v", principal.Permissions)
	}
}

func (f *principalFixture) testCustomScope(t *testing.T) {
	principal := f.resolve(t, "custom-session")
	if !principal.Has("store:read") || principal.Has("store:create") || principal.Has("account:read") {
		t.Fatalf("invalid custom permissions: %#v", principal.Permissions)
	}
	if !principal.HasStore("active-store") || principal.HasStore("draft-store") {
		t.Fatalf("invalid store scope: %#v", principal.StoreIDs)
	}
}

type principalFixture struct {
	t        *testing.T
	db       *gorm.DB
	cfg      config.SecurityConfig
	resolver *PrincipalResolver
	now      time.Time
}

func newPrincipalFixture(t *testing.T) *principalFixture {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	models := []any{&gen.Account{}, &gen.Session{}, &gen.Organization{}, &gen.OperatorMembership{}, &gen.OperatorRole{}, &gen.Permission{}, &gen.Store{}, &membershipRole{}, &permissionRole{}, &membershipStore{}}
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		dropGeneratedIndexes(db)
	}
	if err := dbup.InitPermissions(db); err != nil {
		t.Fatal(err)
	}
	cfg := config.SecurityConfig{SigningKey: []byte(strings.Repeat("k", 32)), TokenIssuer: "test", SessionDuration: 12 * time.Hour}
	fixture := &principalFixture{t: t, db: db, cfg: cfg, resolver: NewPrincipalResolver(db, cfg), now: time.Now()}
	fixture.seed()
	return fixture
}

func (f *principalFixture) seed() {
	create := func(value any) {
		if err := f.db.Create(value).Error; err != nil {
			f.t.Fatal(err)
		}
	}
	create(&gen.Organization{ID: "franchise-org", Code: "F1", Name: "Franchise", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive})
	create(&gen.Organization{ID: "hq-org", Code: "HQ", Name: "HQ", Type: gen.OrganizationTypeHeadquarters, Status: gen.OrganizationStatusActive})
	create(&gen.Store{ID: "active-store", Code: "A", Name: "Active", OrganizationID: "franchise-org", Lifecycle: gen.StoreLifecycleActive})
	create(&gen.Store{ID: "draft-store", Code: "D", Name: "Draft", OrganizationID: "franchise-org", Lifecycle: gen.StoreLifecycleDraft})
	f.seedIdentity("owner", "owner-session", "franchise-org", gen.RoleKindFranchiseOwner, gen.StoreAccessModeAllStores)
	f.seedIdentity("hq", "hq-session", "hq-org", gen.RoleKindHqSuperAdmin, gen.StoreAccessModeAllStores)
	f.seedIdentity("custom", "custom-session", "franchise-org", gen.RoleKindCustom, gen.StoreAccessModeSelectedStores)
	var permission gen.Permission
	f.db.Where("action = ?", "store:read").First(&permission)
	create(&permissionRole{OperatorRoleID: "custom-role", PermissionID: permission.ID})
	f.db.Where("action = ?", "account:read").First(&permission)
	create(&permissionRole{OperatorRoleID: "custom-role", PermissionID: permission.ID})
	create(&membershipStore{OperatorMembershipID: "custom-membership", StoreID: "active-store"})
	create(&membershipStore{OperatorMembershipID: "custom-membership", StoreID: "draft-store"})
	create(&gen.Account{ID: "discovery", Phone: "13000000000", DisplayName: "Discovery", Status: gen.AccountStatusActive, CredentialVersion: 1})
	create(&gen.Session{ID: "discovery-session", AccountID: "discovery", WorkspaceType: gen.WorkspaceTypeDiscovery, CredentialVersion: 1, ExpiresAt: f.now.Add(time.Hour), LastSeenAt: f.now})
}

func (f *principalFixture) seedIdentity(id, sessionID, organizationID string, kind gen.RoleKind, access gen.StoreAccessMode) {
	create := func(value any) {
		if err := f.db.Create(value).Error; err != nil {
			f.t.Fatal(err)
		}
	}
	create(&gen.Account{ID: id, Phone: "13" + id, DisplayName: id, Status: gen.AccountStatusActive, CredentialVersion: 1})
	create(&gen.OperatorMembership{ID: id + "-membership", AccountID: id, OrganizationID: organizationID, Status: gen.MembershipStatusActive, StoreAccessMode: access})
	create(&gen.OperatorRole{ID: id + "-role", Name: id, Kind: kind, OrganizationID: organizationID})
	create(&membershipRole{OperatorMembershipID: id + "-membership", OperatorRoleID: id + "-role"})
	organization := organizationID
	workspace := gen.WorkspaceTypeFranchise
	if organizationID == "hq-org" {
		workspace = gen.WorkspaceTypeHeadquarters
	}
	create(&gen.Session{ID: sessionID, AccountID: id, OrganizationID: &organization, WorkspaceType: workspace, CredentialVersion: 1, ExpiresAt: f.now.Add(time.Hour), LastSeenAt: f.now})
}

func (f *principalFixture) resolve(t *testing.T, sessionID string) *auth.WorkspacePrincipal {
	t.Helper()
	var session gen.Session
	if err := f.db.First(&session, "id = ?", sessionID).Error; err != nil {
		t.Fatal(err)
	}
	claims := auth.SessionClaims{SessionID: session.ID, AccountID: session.AccountID, WorkspaceType: auth.WorkspaceType(session.WorkspaceType), OrganizationID: session.OrganizationID, CredentialVersion: session.CredentialVersion}
	token, err := auth.SignSessionClaims(f.cfg, claims)
	if err != nil {
		t.Fatal(err)
	}
	principal, err := f.resolver.Resolve(context.Background(), token, f.now)
	if err != nil {
		t.Fatal(err)
	}
	return principal
}

func (f *principalFixture) assertStateErrors(t *testing.T) {
	t.Helper()
	f.db.Model(&gen.Session{}).Where("id = ?", "owner-session").Update("revoked_at", f.now)
	_, err := f.resolveError("owner-session")
	assertCode(t, err, auth.CodeSessionRevoked)
	f.db.Model(&gen.Session{}).Where("id = ?", "owner-session").Update("revoked_at", nil)
	f.db.Model(&gen.Account{}).Where("id = ?", "owner").Update("credential_version", 2)
	_, err = f.resolveError("owner-session")
	assertCode(t, err, auth.CodeCredentialsChanged)
	f.db.Model(&gen.Account{}).Where("id = ?", "owner").Update("credential_version", 1)
	f.db.Model(&gen.Organization{}).Where("id = ?", "franchise-org").Update("status", gen.OrganizationStatusSuspended)
	_, err = f.resolveError("owner-session")
	assertCode(t, err, auth.CodeOrganizationSuspended)
	f.db.Model(&gen.Organization{}).Where("id = ?", "franchise-org").Update("status", gen.OrganizationStatusActive)
	f.db.Model(&gen.OperatorMembership{}).Where("id = ?", "owner-membership").Update("status", gen.MembershipStatusSuspended)
	_, err = f.resolveError("owner-session")
	assertCode(t, err, auth.CodeMembershipInactive)
}

func (f *principalFixture) resolveError(sessionID string) (*auth.WorkspacePrincipal, error) {
	var session gen.Session
	f.db.First(&session, "id = ?", sessionID)
	claims := auth.SessionClaims{SessionID: session.ID, AccountID: session.AccountID, WorkspaceType: auth.WorkspaceType(session.WorkspaceType), OrganizationID: session.OrganizationID, CredentialVersion: session.CredentialVersion}
	token, _ := auth.SignSessionClaims(f.cfg, claims)
	return f.resolver.Resolve(context.Background(), token, f.now)
}

func assertCode(t *testing.T, err error, code auth.Code) {
	t.Helper()
	if auth.ErrorCode(err) != code {
		t.Fatalf("error code = %q, want %q (%v)", auth.ErrorCode(err), code, err)
	}
}

func dropGeneratedIndexes(db *gorm.DB) {
	for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
		db.Exec("DROP INDEX IF EXISTS `" + index + "`")
	}
}

type membershipRole struct{ OperatorMembershipID, OperatorRoleID string }

func (membershipRole) TableName() string { return "operator_membership_roles" }

type permissionRole struct{ PermissionID, OperatorRoleID string }

func (permissionRole) TableName() string { return "permission_roles" }

type membershipStore struct{ OperatorMembershipID, StoreID string }

func (membershipStore) TableName() string { return "operator_membership_stores" }
