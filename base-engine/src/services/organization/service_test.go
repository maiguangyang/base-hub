/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package organization

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

// TestProvisionFranchiseNewAndExistingAccounts 验证新账号直开与既有账号邀请两条路径。
func TestProvisionFranchiseNewAndExistingAccounts(t *testing.T) {
	service, db := newOrganizationFixture(t)
	principal := headquartersPrincipal()
	first, err := service.ProvisionFranchise(context.Background(), principal, ProvisionInput{Code: "F001", Name: "First", OwnerPhone: "13800000000", OwnerDisplayName: "Owner"})
	if err != nil {
		t.Fatal(err)
	}
	assertNewOwnerResult(t, first)
	assertTemporaryPasswordExpiry(t, db, first.Membership.AccountID)
	var firstCredential authentication.AccountCredential
	if err := db.First(&firstCredential, "account_id = ?", first.Membership.AccountID).Error; err != nil {
		t.Fatal(err)
	}
	second, err := service.ProvisionFranchise(context.Background(), principal, ProvisionInput{Code: "F002", Name: "Second", OwnerPhone: "13800000000", OwnerDisplayName: "Owner"})
	if err != nil {
		t.Fatal(err)
	}
	assertExistingOwnerResult(t, second)
	assertSharedInitialAccount(t, db, first, second, firstCredential.PasswordHash)
	var invitation gen.MembershipInvitation
	if err := db.First(&invitation, "membership_id = ?", second.Membership.ID).Error; err != nil {
		t.Fatal(err)
	}
	if invitation.ExpiresAt.Before(time.Now().Add(71 * time.Hour)) {
		t.Fatalf("invitation expiry = %v", invitation.ExpiresAt)
	}
}

func TestProvisionFranchiseRejectsHeadquartersAccountAsInitialAccount(t *testing.T) {
	service, db := newOrganizationFixture(t)
	account := gen.Account{ID: "hq-owner", Phone: "13800000000", DisplayName: "HQ", Status: gen.AccountStatusActive}
	hq := gen.Organization{ID: "hq", Code: "HQ", Name: "HQ", Type: gen.OrganizationTypeHeadquarters, Status: gen.OrganizationStatusActive}
	membership := gen.OperatorMembership{ID: "hq-member", AccountID: account.ID, OrganizationID: hq.ID, Status: gen.MembershipStatusActive}
	for _, record := range []any{&account, &hq, &membership} {
		if err := db.Create(record).Error; err != nil {
			t.Fatal(err)
		}
	}
	_, err := service.ProvisionFranchise(context.Background(), headquartersPrincipal(), ProvisionInput{Code: "F001", Name: "First", OwnerPhone: account.Phone, OwnerDisplayName: "Owner"})
	if auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("HQ account provision error = %v", err)
	}
	var count int64
	db.Model(&gen.Organization{}).Where("code = ?", "F001").Count(&count)
	if count != 0 {
		t.Fatal("rejected provision left a franchise organization")
	}
}

func assertSharedInitialAccount(t *testing.T, db *gorm.DB, first, second *ProvisionResult, originalHash string) {
	t.Helper()
	for _, result := range []*ProvisionResult{first, second} {
		var organization gen.Organization
		if err := db.First(&organization, "id = ?", result.Organization.ID).Error; err != nil {
			t.Fatal(err)
		}
		if organization.InitialAccountID == nil || *organization.InitialAccountID != first.Membership.AccountID {
			t.Fatalf("initial account not persisted for %s: %#v", organization.ID, organization.InitialAccountID)
		}
	}
	var secondCredential authentication.AccountCredential
	if err := db.First(&secondCredential, "account_id = ?", first.Membership.AccountID).Error; err != nil {
		t.Fatal(err)
	}
	if secondCredential.PasswordHash != originalHash {
		t.Fatal("second franchise provisioning rotated the shared account password")
	}
}

func assertTemporaryPasswordExpiry(t *testing.T, db *gorm.DB, accountID string) {
	t.Helper()
	var credential authentication.AccountCredential
	if err := db.First(&credential, "account_id = ?", accountID).Error; err != nil {
		t.Fatal(err)
	}
	if credential.TemporaryPasswordExpiresAt == nil || !credential.TemporaryPasswordExpiresAt.After(time.Now()) {
		t.Fatalf("temporary password expiry missing: %#v", credential)
	}
}

func assertNewOwnerResult(t *testing.T, result *ProvisionResult) {
	t.Helper()
	if result.TemporaryPassword == nil {
		t.Fatalf("temporary password missing: %#v", result)
	}
	if result.InvitationPending || result.Membership.Status != gen.MembershipStatusActive {
		t.Fatalf("invalid new-account result: %#v", result)
	}
}

func assertExistingOwnerResult(t *testing.T, result *ProvisionResult) {
	t.Helper()
	if result.TemporaryPassword != nil {
		t.Fatalf("existing password was exposed: %#v", result)
	}
	if !result.InvitationPending || result.Membership.Status != gen.MembershipStatusInvited {
		t.Fatalf("invalid existing-account result: %#v", result)
	}
}

// TestProvisionFranchiseRequiresHQAndRollsBackDuplicates 验证总部权限和事务原子性。
func TestProvisionFranchiseRequiresHQAndRollsBackDuplicates(t *testing.T) {
	service, db := newOrganizationFixture(t)
	franchise := &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeFranchise, Permissions: map[string]struct{}{"franchise:provision": {}}}
	if _, err := service.ProvisionFranchise(context.Background(), franchise, ProvisionInput{Code: "F001", Name: "First", OwnerPhone: "13800000000", OwnerDisplayName: "Owner"}); auth.ErrorCode(err) != auth.CodeWorkspaceForbidden {
		t.Fatalf("franchise provision error = %v", err)
	}
	principal := headquartersPrincipal()
	if _, err := service.ProvisionFranchise(context.Background(), principal, ProvisionInput{Code: "F001", Name: "First", OwnerPhone: "13800000000", OwnerDisplayName: "Owner"}); err != nil {
		t.Fatal(err)
	}
	if _, err := service.ProvisionFranchise(context.Background(), principal, ProvisionInput{Code: "F001", Name: "Duplicate", OwnerPhone: "13900000000", OwnerDisplayName: "Other"}); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("duplicate organization code error = %v", err)
	}
	var count int64
	db.Model(&gen.Account{}).Where("phone = ?", "13900000000").Count(&count)
	if count != 0 {
		t.Fatal("duplicate transaction left a partial account")
	}
}

// TestProvisionFranchiseValidation 验证自定义开通入口与实体字段规则保持一致。
func TestProvisionFranchiseValidation(t *testing.T) {
	service, db := newOrganizationFixture(t)
	principal := headquartersPrincipal()
	email := "owner@example.com"
	base := ProvisionInput{Code: "F001", Name: "First", OwnerPhone: "13800000000", OwnerDisplayName: "Owner", OwnerEmail: &email}
	assertInvalidProvisionInputs(t, service, principal, base)
	result, err := service.ProvisionFranchise(context.Background(), principal, ProvisionInput{
		Code: " F001 ", Name: " First ", OwnerPhone: " 13800000000 ", OwnerDisplayName: " Owner ", OwnerEmail: organizationStringPointer(" owner@example.com "),
	})
	if err != nil || result.Organization.Code != "F001" || result.Organization.Name != "First" {
		t.Fatalf("normalized result=%#v err=%v", result, err)
	}
	var account gen.Account
	if err := db.First(&account, "id = ?", result.Membership.AccountID).Error; err != nil || account.Phone != "13800000000" || account.DisplayName != "Owner" || account.Email == nil || *account.Email != "owner@example.com" {
		t.Fatalf("normalized account=%#v err=%v", account, err)
	}
}

func assertInvalidProvisionInputs(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, base ProvisionInput) {
	t.Helper()
	cases := []ProvisionInput{
		{Code: " ", Name: base.Name, OwnerPhone: base.OwnerPhone, OwnerDisplayName: base.OwnerDisplayName},
		{Code: strings.Repeat("C", 33), Name: base.Name, OwnerPhone: base.OwnerPhone, OwnerDisplayName: base.OwnerDisplayName},
		{Code: base.Code, Name: " ", OwnerPhone: base.OwnerPhone, OwnerDisplayName: base.OwnerDisplayName},
		{Code: base.Code, Name: strings.Repeat("N", 129), OwnerPhone: base.OwnerPhone, OwnerDisplayName: base.OwnerDisplayName},
		{Code: base.Code, Name: base.Name, OwnerPhone: "123", OwnerDisplayName: base.OwnerDisplayName},
		{Code: base.Code, Name: base.Name, OwnerPhone: base.OwnerPhone, OwnerDisplayName: " "},
		{Code: base.Code, Name: base.Name, OwnerPhone: base.OwnerPhone, OwnerDisplayName: strings.Repeat("D", 65)},
		{Code: base.Code, Name: base.Name, OwnerPhone: base.OwnerPhone, OwnerDisplayName: base.OwnerDisplayName, OwnerEmail: organizationStringPointer("invalid")},
		{Code: base.Code, Name: base.Name, OwnerPhone: base.OwnerPhone, OwnerDisplayName: base.OwnerDisplayName, OwnerEmail: organizationStringPointer(strings.Repeat("a", 120) + "@example.com")},
	}
	for index, input := range cases {
		if _, err := service.ProvisionFranchise(context.Background(), principal, input); auth.ErrorCode(err) != auth.CodeValidationFailed {
			t.Fatalf("case %d error = %v", index, err)
		}
	}
}

// TestSuspendOrganizationValidation 验证暂停原因拒绝空白和越界值。
func TestSuspendOrganizationValidation(t *testing.T) {
	service, db := newOrganizationFixture(t)
	principal := headquartersPrincipal()
	principal.Permissions["organization:suspend"] = struct{}{}
	organization := gen.Organization{ID: "org-a", Code: "F001", Name: "First", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}
	if err := db.Create(&organization).Error; err != nil {
		t.Fatal(err)
	}
	for _, reason := range []string{" ", strings.Repeat("R", 65)} {
		if err := service.Suspend(context.Background(), principal, organization.ID, reason); auth.ErrorCode(err) != auth.CodeValidationFailed {
			t.Fatalf("reason length %d error = %v", len(reason), err)
		}
	}
}

func organizationStringPointer(value string) *string { return &value }

func newOrganizationFixture(t *testing.T) (*Service, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	models := []any{&gen.Account{}, &gen.Organization{}, &gen.FranchiseOpeningRecord{}, &gen.OperatorMembership{}, &gen.OperatorRole{}, &gen.MembershipInvitation{}, &gen.AuditLog{}, &gen.Session{}, &authentication.AccountCredential{}, &membershipRole{}}
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	return NewService(db, audit.NewService()), db
}

func headquartersPrincipal() *auth.WorkspacePrincipal {
	return &auth.WorkspacePrincipal{AccountID: "hq-account", SessionID: "hq-session", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"franchise:provision": {}}}
}

type membershipRole struct{ OperatorMembershipID, OperatorRoleID string }

func (membershipRole) TableName() string { return "operator_membership_roles" }
