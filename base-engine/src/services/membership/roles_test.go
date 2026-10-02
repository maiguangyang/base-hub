/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package membership

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

// TestSaveCustomRoleRejectsCrossScopePermissions 验证加盟角色不能获得 SYSTEM 权限。
func TestSaveCustomRoleRejectsCrossScopePermissions(t *testing.T) {
	_, db, principal := newMembershipFixture(t)
	systemPermission := gen.Permission{ID: "system-permission", Name: "permission.account.read", Action: "account:read", Module: "account", Scope: gen.PermissionScopeSystem}
	tenantPermission := gen.Permission{ID: "tenant-permission", Name: "permission.store.read", Action: "store:read", Module: "store", Scope: gen.PermissionScopeTenant}
	db.Create(&systemPermission)
	db.Create(&tenantPermission)
	service := NewRoleService(db, nil)
	if _, err := service.SaveCustomRole(context.Background(), principal, SaveRoleInput{Name: "Manager", PermissionIDs: []string{systemPermission.ID}}); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("SYSTEM permission error = %v", err)
	}
	role, err := service.SaveCustomRole(context.Background(), principal, SaveRoleInput{Name: "Manager", PermissionIDs: []string{tenantPermission.ID}})
	if err != nil || role.Kind != gen.RoleKindCustom {
		t.Fatalf("custom role result = %#v, %v", role, err)
	}
}

// TestFinalOwnerCannotBeSuspended 验证最后一个活跃老板不可被停用。
func TestFinalOwnerCannotBeSuspended(t *testing.T) {
	service, db, principal := newMembershipFixture(t)
	ownerRole := gen.OperatorRole{ID: "owner-role", Name: "Owner", Kind: gen.RoleKindFranchiseOwner, OrganizationID: *principal.OrganizationID}
	owner := gen.OperatorMembership{ID: "owner-membership", AccountID: "owner", OrganizationID: *principal.OrganizationID, Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores}
	db.Create(&ownerRole)
	db.Create(&owner)
	db.Create(&membershipRole{OperatorMembershipID: owner.ID, OperatorRoleID: ownerRole.ID})
	err := service.ChangeMembershipStatus(context.Background(), principal, owner.ID, gen.MembershipStatusSuspended)
	if auth.ErrorCode(err) != auth.CodeLastOwnerRequired {
		t.Fatalf("last owner error = %v", err)
	}
}
