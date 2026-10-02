/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-22
 */
package membership

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestHQInviteCreatesAllStoresAdministrator(t *testing.T) {
	service, db, principal := newHQMembershipFixture(t, true)

	result, err := service.InviteOperator(context.Background(), principal, InviteInput{
		Phone: "13900000000", DisplayName: "HQ Admin", RoleIDs: []string{"custom-role"},
		StoreAccessMode: gen.StoreAccessModeSelectedStores, StoreIDs: []string{"active-store"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Membership.StoreAccessMode != gen.StoreAccessModeAllStores {
		t.Fatalf("store access mode = %s", result.Membership.StoreAccessMode)
	}
	var storeCount int64
	if err := db.Table("operator_membership_stores").Where("operator_membership_id = ?", result.Membership.ID).Count(&storeCount).Error; err != nil {
		t.Fatal(err)
	}
	if storeCount != 0 {
		t.Fatalf("HQ membership retained %d stores", storeCount)
	}
}

func TestHQInviteRejectsMissingOrStrongerRoles(t *testing.T) {
	service, db, principal := newHQMembershipFixture(t, false)
	base := InviteInput{Phone: "13900000001", DisplayName: "HQ Admin", StoreAccessMode: gen.StoreAccessModeAllStores}
	if _, err := service.InviteOperator(context.Background(), principal, base); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("missing role error = %v", err)
	}
	permission := gen.Permission{ID: "system-account-delete", Name: "account delete", Action: "account:delete", Module: "account", Scope: gen.PermissionScopeSystem}
	if err := db.Create(&permission).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&permissionRole{PermissionID: permission.ID, OperatorRoleID: "custom-role"}).Error; err != nil {
		t.Fatal(err)
	}
	base.RoleIDs = []string{"custom-role"}
	if _, err := service.InviteOperator(context.Background(), principal, base); auth.ErrorCode(err) != auth.Code("PERMISSION_DELEGATION_DENIED") {
		t.Fatalf("strong role error = %v", err)
	}
	base.RoleIDs = []string{"hq-super-role"}
	if _, err := service.InviteOperator(context.Background(), principal, base); auth.ErrorCode(err) != auth.Code("PERMISSION_DELEGATION_DENIED") {
		t.Fatalf("super role error = %v", err)
	}
}

func TestHQInviteRejectsTenantPermissionsAndForeignRoles(t *testing.T) {
	service, db, principal := newHQMembershipFixture(t, false)
	base := InviteInput{Phone: "13900000004", DisplayName: "HQ Admin", StoreAccessMode: gen.StoreAccessModeAllStores}
	tenantPermission := gen.Permission{ID: "tenant-order-read", Name: "order read", Action: "order:read", Module: "order", Scope: gen.PermissionScopeTenant}
	foreignRole := gen.OperatorRole{ID: "foreign-role", Name: "Foreign", Kind: gen.RoleKindCustom, OrganizationID: "other-org"}
	if err := db.Create(&tenantPermission).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&foreignRole).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&gen.OperatorRole{ID: "custom-role"}).Association("Permissions").Replace(&tenantPermission); err != nil {
		t.Fatal(err)
	}
	base.RoleIDs = []string{"custom-role"}
	if _, err := service.InviteOperator(context.Background(), principal, base); auth.ErrorCode(err) != auth.CodePermissionDelegationDenied {
		t.Fatalf("tenant permission error = %v", err)
	}
	base.RoleIDs = []string{foreignRole.ID}
	if _, err := service.InviteOperator(context.Background(), principal, base); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("foreign role error = %v", err)
	}
}

func TestHQSuperAdminCanAssignSuperRoleWithoutOverwritingAccount(t *testing.T) {
	service, db, principal := newHQMembershipFixture(t, true)
	email := "original@example.com"
	account := gen.Account{ID: "existing-account", Phone: "13900000003", DisplayName: "Original", Email: &email, Status: gen.AccountStatusActive}
	if err := db.Create(&account).Error; err != nil {
		t.Fatal(err)
	}
	replacementEmail := "replacement@example.com"
	result, err := service.InviteOperator(context.Background(), principal, InviteInput{
		Phone: account.Phone, DisplayName: "Replacement", Email: &replacementEmail,
		RoleIDs: []string{"hq-super-role"}, StoreAccessMode: gen.StoreAccessModeSelectedStores,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.InvitationPending || result.TemporaryPassword != nil {
		t.Fatalf("existing account result = %#v", result)
	}
	var reloaded gen.Account
	if err := db.First(&reloaded, "id = ?", account.ID).Error; err != nil {
		t.Fatal(err)
	}
	if reloaded.Phone != account.Phone || reloaded.DisplayName != account.DisplayName || reloaded.Email == nil || *reloaded.Email != email {
		t.Fatalf("existing account was overwritten: %#v", reloaded)
	}
}

func TestHQInviteRejectsNonHeadquartersOrganization(t *testing.T) {
	service, db, principal := newHQMembershipFixture(t, true)
	if err := db.Model(&gen.Organization{}).Where("id = ?", *principal.OrganizationID).
		Update("type", gen.OrganizationTypeFranchise).Error; err != nil {
		t.Fatal(err)
	}
	_, err := service.InviteOperator(context.Background(), principal, InviteInput{
		Phone: "13900000002", DisplayName: "HQ Admin", RoleIDs: []string{"custom-role"},
		StoreAccessMode: gen.StoreAccessModeAllStores,
	})
	if auth.ErrorCode(err) != auth.CodeWorkspaceForbidden {
		t.Fatalf("organization type error = %v", err)
	}
}

func newHQMembershipFixture(t *testing.T, super bool) (*Service, *gorm.DB, *auth.WorkspacePrincipal) {
	t.Helper()
	service, db, principal := newMembershipFixture(t)
	organizationID := *principal.OrganizationID
	if err := db.Model(&gen.Organization{}).Where("id = ?", organizationID).
		Update("type", gen.OrganizationTypeHeadquarters).Error; err != nil {
		t.Fatal(err)
	}
	principal.WorkspaceType = auth.WorkspaceTypeHeadquarters
	principal.Permissions = map[string]struct{}{"hqMembership:create": {}}
	seedHQInviter(t, db, organizationID)
	principal.AccountID = "hq-custom-account"
	principal.MembershipID = membershipStringPointer("hq-custom-membership")
	if super {
		principal.AccountID = "hq-super-account"
		principal.MembershipID = membershipStringPointer("hq-super-membership")
	}
	return service, db, principal
}

func seedHQInviter(t *testing.T, db *gorm.DB, organizationID string) {
	t.Helper()
	values := []any{
		&gen.Account{ID: "hq-super-account", Phone: "13700000000", DisplayName: "Super", Status: gen.AccountStatusActive},
		&gen.OperatorMembership{ID: "hq-super-membership", AccountID: "hq-super-account", OrganizationID: organizationID, Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores},
		&gen.OperatorRole{ID: "hq-super-role", Name: "Super", Kind: gen.RoleKindHqSuperAdmin, OrganizationID: organizationID},
		&membershipRole{OperatorMembershipID: "hq-super-membership", OperatorRoleID: "hq-super-role"},
		&gen.Account{ID: "hq-custom-account", Phone: "13600000000", DisplayName: "Custom", Status: gen.AccountStatusActive},
		&gen.OperatorMembership{ID: "hq-custom-membership", AccountID: "hq-custom-account", OrganizationID: organizationID, Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores},
		&gen.OperatorRole{ID: "hq-custom-actor-role", Name: "Custom Actor", Kind: gen.RoleKindCustom, OrganizationID: organizationID},
		&membershipRole{OperatorMembershipID: "hq-custom-membership", OperatorRoleID: "hq-custom-actor-role"},
	}
	for _, value := range values {
		if err := db.Create(value).Error; err != nil {
			t.Fatal(err)
		}
	}
}
