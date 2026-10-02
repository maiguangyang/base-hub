/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package membership

import (
	"context"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authentication"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestInviteOperatorNewAndExistingAccounts 验证新账号激活与既有账号邀请语义。
func TestInviteOperatorNewAndExistingAccounts(t *testing.T) {
	service, db, principal := newMembershipFixture(t)
	input := InviteInput{Phone: "13800000000", DisplayName: "Staff", RoleIDs: []string{"custom-role"}, StoreAccessMode: gen.StoreAccessModeSelectedStores, StoreIDs: []string{"active-store"}}
	first, err := service.InviteOperator(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	if first.TemporaryPassword == nil || first.InvitationPending || first.Membership.Status != gen.MembershipStatusActive {
		t.Fatalf("invalid new staff result: %#v", first)
	}
	assertInviteTemporaryPasswordExpiry(t, db, first.Membership.AccountID)
	otherOrg := gen.Organization{ID: "org-b", Code: "B", Name: "B", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}
	db.Create(&otherOrg)
	otherRole := gen.OperatorRole{ID: "custom-role-b", Name: "Staff", Kind: gen.RoleKindCustom, OrganizationID: otherOrg.ID}
	otherStore := gen.Store{ID: "active-store-b", Code: "B", Name: "Active", OrganizationID: otherOrg.ID, Lifecycle: gen.StoreLifecycleActive}
	db.Create(&otherRole)
	db.Create(&otherStore)
	principal.OrganizationID = &otherOrg.ID
	input.RoleIDs = []string{otherRole.ID}
	input.StoreIDs = []string{otherStore.ID}
	second, err := service.InviteOperator(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	if second.TemporaryPassword != nil || !second.InvitationPending || second.Membership.Status != gen.MembershipStatusInvited {
		t.Fatalf("invalid existing staff result: %#v", second)
	}
}

// TestInviteOperatorReissuesExpiredInvitation 验证失效邀请可复用原成员记录重发。
func TestInviteOperatorReissuesExpiredInvitation(t *testing.T) {
	service, db, principal, input, first, invitation := expiredInvitationFixture(t)
	reissued, err := service.InviteOperator(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	assertReissuedInvitation(t, db, first.Membership.ID, invitation.ID, reissued)
}

func expiredInvitationFixture(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal, InviteInput, *InviteResult, gen.MembershipInvitation) {
	t.Helper()
	service, db, principal := newMembershipFixture(t)
	input := InviteInput{Phone: "13800000000", DisplayName: "Staff", RoleIDs: []string{"custom-role"}, StoreAccessMode: gen.StoreAccessModeSelectedStores, StoreIDs: []string{"active-store"}}
	if _, err := service.InviteOperator(context.Background(), principal, input); err != nil {
		t.Fatal(err)
	}
	organization := gen.Organization{ID: "org-reissue", Code: "R", Name: "Reissue", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}
	role := gen.OperatorRole{ID: "role-reissue", Name: "Staff", Kind: gen.RoleKindCustom, OrganizationID: organization.ID}
	store := gen.Store{ID: "store-reissue", Code: "R", Name: "Store", OrganizationID: organization.ID, Lifecycle: gen.StoreLifecycleActive}
	db.Create(&organization)
	db.Create(&role)
	db.Create(&store)
	principal.OrganizationID = &organization.ID
	input.RoleIDs, input.StoreIDs = []string{role.ID}, []string{store.ID}
	first, err := service.InviteOperator(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	var invitation gen.MembershipInvitation
	if err := db.First(&invitation, "membership_id = ?", first.Membership.ID).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&invitation).Updates(map[string]any{"expires_at": time.Now().Add(-time.Hour), "revoked_at": time.Now()}).Error; err != nil {
		t.Fatal(err)
	}
	return service, db, principal, input, first, invitation
}

func assertReissuedInvitation(t *testing.T, db *gorm.DB, membershipID, invitationID string, result *InviteResult) {
	t.Helper()
	if result.Membership.ID != membershipID || !result.InvitationPending {
		t.Fatalf("reissued result = %#v", result)
	}
	var refreshed gen.MembershipInvitation
	if err := db.First(&refreshed, "membership_id = ?", membershipID).Error; err != nil {
		t.Fatal(err)
	}
	if refreshed.ID != invitationID || refreshed.RevokedAt != nil || refreshed.AcceptedAt != nil || !refreshed.ExpiresAt.After(time.Now()) {
		t.Fatalf("invitation was not refreshed: %#v", refreshed)
	}
}

func assertInviteTemporaryPasswordExpiry(t *testing.T, db *gorm.DB, accountID string) {
	t.Helper()
	var credential authentication.AccountCredential
	if err := db.First(&credential, "account_id = ?", accountID).Error; err != nil {
		t.Fatal(err)
	}
	if credential.TemporaryPasswordExpiresAt == nil || !credential.TemporaryPasswordExpiresAt.After(time.Now()) {
		t.Fatalf("temporary password expiry missing: %#v", credential)
	}
}

// TestInviteOperatorValidatesStoreModeAndOrganization 验证门店范围与角色组织边界。
func TestInviteOperatorValidatesStoreModeAndOrganization(t *testing.T) {
	service, _, principal := newMembershipFixture(t)
	invalid := InviteInput{Phone: "13800000000", DisplayName: "Staff", RoleIDs: []string{"custom-role"}, StoreAccessMode: gen.StoreAccessModeAllStores, StoreIDs: []string{"active-store"}}
	if _, err := service.InviteOperator(context.Background(), principal, invalid); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("invalid mode error = %v", err)
	}
	invalid.StoreAccessMode = gen.StoreAccessModeSelectedStores
	invalid.StoreIDs = []string{"draft-store"}
	if _, err := service.InviteOperator(context.Background(), principal, invalid); auth.ErrorCode(err) != auth.CodeStoreNotActive {
		t.Fatalf("draft store error = %v", err)
	}
}

// TestInviteOperatorValidation 验证邀请入口拒绝无效账号资料并规范化持久化值。
func TestInviteOperatorValidation(t *testing.T) {
	service, db, principal := newMembershipFixture(t)
	base := InviteInput{Phone: "13800000000", DisplayName: "Staff", RoleIDs: []string{"custom-role"}, StoreAccessMode: gen.StoreAccessModeAllStores}
	cases := []InviteInput{
		{Phone: "123", DisplayName: base.DisplayName, RoleIDs: base.RoleIDs, StoreAccessMode: base.StoreAccessMode},
		{Phone: base.Phone, DisplayName: " ", RoleIDs: base.RoleIDs, StoreAccessMode: base.StoreAccessMode},
		{Phone: base.Phone, DisplayName: strings.Repeat("D", 65), RoleIDs: base.RoleIDs, StoreAccessMode: base.StoreAccessMode},
		{Phone: base.Phone, DisplayName: base.DisplayName, Email: membershipStringPointer("invalid"), RoleIDs: base.RoleIDs, StoreAccessMode: base.StoreAccessMode},
		{Phone: base.Phone, DisplayName: base.DisplayName, Email: membershipStringPointer(strings.Repeat("a", 120) + "@example.com"), RoleIDs: base.RoleIDs, StoreAccessMode: base.StoreAccessMode},
	}
	for index, input := range cases {
		if _, err := service.InviteOperator(context.Background(), principal, input); auth.ErrorCode(err) != auth.CodeValidationFailed {
			t.Fatalf("case %d error = %v", index, err)
		}
	}
	result, err := service.InviteOperator(context.Background(), principal, InviteInput{
		Phone: " 13800000000 ", DisplayName: " Staff ", Email: membershipStringPointer(" staff@example.com "), RoleIDs: base.RoleIDs, StoreAccessMode: base.StoreAccessMode,
	})
	if err != nil {
		t.Fatal(err)
	}
	var account gen.Account
	if err := db.First(&account, "id = ?", result.Membership.AccountID).Error; err != nil || account.Phone != "13800000000" || account.DisplayName != "Staff" || account.Email == nil || *account.Email != "staff@example.com" {
		t.Fatalf("normalized account=%#v err=%v", account, err)
	}
}

func membershipStringPointer(value string) *string { return &value }

// TestAcceptInvitationRevalidatesPreparedAccess 验证邀请激活前重新检查角色和门店仍然有效。
func TestAcceptInvitationRevalidatesPreparedAccess(t *testing.T) {
	service, db, principal := newMembershipFixture(t)
	input := InviteInput{Phone: "13800000000", DisplayName: "Staff", RoleIDs: []string{"custom-role"}, StoreAccessMode: gen.StoreAccessModeSelectedStores, StoreIDs: []string{"active-store"}}
	first, err := service.InviteOperator(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	organization := gen.Organization{ID: "org-b", Code: "B", Name: "B", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}
	role := gen.OperatorRole{ID: "role-b", Name: "Staff", Kind: gen.RoleKindCustom, OrganizationID: organization.ID}
	store := gen.Store{ID: "store-b", Code: "B", Name: "Store", OrganizationID: organization.ID, Lifecycle: gen.StoreLifecycleActive}
	db.Create(&organization)
	db.Create(&role)
	db.Create(&store)
	principal.OrganizationID = &organization.ID
	input.RoleIDs, input.StoreIDs = []string{role.ID}, []string{store.ID}
	pending, err := service.InviteOperator(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	if err := service.ChangeMembershipStatus(context.Background(), principal, pending.Membership.ID, gen.MembershipStatusActive); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("invited membership activated without acceptance: %v", err)
	}
	var invitation gen.MembershipInvitation
	db.First(&invitation, "membership_id = ?", pending.Membership.ID)
	db.Model(&role).Update("is_delete", 2)
	invitee := &auth.WorkspacePrincipal{AccountID: first.Membership.AccountID, SessionID: "invitee-session", WorkspaceType: auth.WorkspaceTypeDiscovery}
	if _, err := service.AcceptInvitation(context.Background(), invitee, invitation.ID); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("deleted role accepted: %v", err)
	}
	db.Model(&role).Update("is_delete", 1)
	accepted, err := service.AcceptInvitation(context.Background(), invitee, invitation.ID)
	if err != nil || accepted.Status != gen.MembershipStatusActive {
		t.Fatalf("valid invitation result=%#v err=%v", accepted, err)
	}
	var record gen.AuditLog
	if err := db.Where("action = ?", "invitation:accept").First(&record).Error; err != nil || record.MetadataJSON == nil || !strings.Contains(*record.MetadataJSON, `"targetStatus":"ACTIVE"`) {
		t.Fatalf("acceptance audit is stale: %#v err=%v", record.MetadataJSON, err)
	}
}

func newMembershipFixture(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	models := []any{&gen.Account{}, &gen.Organization{}, &gen.OperatorMembership{}, &gen.OperatorRole{}, &gen.Permission{}, &gen.Store{}, &gen.MembershipInvitation{}, &gen.AuditLog{}, &authentication.AccountCredential{}, &membershipRole{}, &membershipStore{}, &permissionRole{}}
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	organization := gen.Organization{ID: "org-a", Code: "A", Name: "A", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}
	role := gen.OperatorRole{ID: "custom-role", Name: "Staff", Kind: gen.RoleKindCustom, OrganizationID: organization.ID}
	stores := []gen.Store{{ID: "active-store", Code: "A", Name: "Active", OrganizationID: organization.ID, Lifecycle: gen.StoreLifecycleActive}, {ID: "draft-store", Code: "D", Name: "Draft", OrganizationID: organization.ID, Lifecycle: gen.StoreLifecycleDraft}}
	db.Create(&organization)
	db.Create(&role)
	db.Create(&stores)
	principal := &auth.WorkspacePrincipal{AccountID: "owner", SessionID: "session", WorkspaceType: auth.WorkspaceTypeFranchise, OrganizationID: &organization.ID, Permissions: map[string]struct{}{"operatorMembership:create": {}, "operatorMembership:update": {}, "operatorRole:create": {}, "operatorRole:update": {}}, StoreIDs: map[string]struct{}{"active-store": {}}, AllStores: true}
	return NewService(db, audit.NewService()), db, principal
}

type membershipRole struct{ OperatorMembershipID, OperatorRoleID string }

func (membershipRole) TableName() string { return "operator_membership_roles" }

type membershipStore struct{ OperatorMembershipID, StoreID string }

func (membershipStore) TableName() string { return "operator_membership_stores" }

type permissionRole struct{ PermissionID, OperatorRoleID string }

func (permissionRole) TableName() string { return "permission_roles" }
