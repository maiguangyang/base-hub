/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-22
 */
package authorization

import (
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

// ValidateHQRoleDelegation ensures assigned HQ roles do not exceed the actor's authority.
func ValidateHQRoleDelegation(database *gorm.DB, principal *auth.WorkspacePrincipal, organizationID string, roleIDs []string) error {
	if len(roleIDs) == 0 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	roles, err := loadHQRoles(database, organizationID, roleIDs)
	if err != nil {
		return err
	}
	super, err := IsHQSuperAdmin(database, principal)
	if err != nil {
		return err
	}
	return validateDelegableRoles(database, principal, roles, super)
}

func validateDelegableRoles(database *gorm.DB, principal *auth.WorkspacePrincipal, roles []*gen.OperatorRole, super bool) error {
	for _, role := range roles {
		if role.Kind == gen.RoleKindHqSuperAdmin && !super {
			return auth.NewError(auth.CodePermissionDelegationDenied)
		}
		if role.Kind != gen.RoleKindCustom && role.Kind != gen.RoleKindHqSuperAdmin {
			return auth.NewError(auth.CodePermissionDelegationDenied)
		}
		if !super {
			if err := validateRolePermissionSubset(database, principal, role); err != nil {
				return err
			}
		}
	}
	return nil
}

// IsHQSuperAdmin checks the actor's persisted active role rather than inferring from permissions.
func IsHQSuperAdmin(database *gorm.DB, principal *auth.WorkspacePrincipal) (bool, error) {
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeHeadquarters || principal.MembershipID == nil || principal.OrganizationID == nil {
		return false, nil
	}
	var count int64
	err := database.Table("operator_membership_roles mr").
		Joins("JOIN operator_memberships m ON m.id = mr.operator_membership_id").
		Joins("JOIN operator_roles r ON r.id = mr.operator_role_id").
		Where("mr.operator_membership_id = ? AND m.account_id = ?", *principal.MembershipID, principal.AccountID).
		Where("m.organization_id = ? AND m.status = ?", *principal.OrganizationID, gen.MembershipStatusActive).
		Where("r.organization_id = ? AND r.kind = ?", *principal.OrganizationID, gen.RoleKindHqSuperAdmin).
		Where("m.is_delete IS NULL OR m.is_delete = ?", 1).
		Where("r.is_delete IS NULL OR r.is_delete = ?", 1).Count(&count).Error
	return count > 0, err
}

// ValidateHQMembershipTarget prevents self-management and privilege escalation.
func ValidateHQMembershipTarget(database *gorm.DB, principal *auth.WorkspacePrincipal, membership *gen.OperatorMembership) error {
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return nil
	}
	if principal.MembershipID != nil && *principal.MembershipID == membership.ID {
		return auth.NewError(auth.CodeSelfMembershipChangeDenied)
	}
	actorSuper, err := IsHQSuperAdmin(database, principal)
	if err != nil {
		return err
	}
	roles, err := loadMembershipRoles(database, membership)
	if err != nil {
		return err
	}
	return validateManagedRoles(database, principal, roles, actorSuper)
}

func validateManagedRoles(database *gorm.DB, principal *auth.WorkspacePrincipal, roles []*gen.OperatorRole, actorSuper bool) error {
	for _, role := range roles {
		if role.Kind == gen.RoleKindHqSuperAdmin && !actorSuper {
			return auth.NewError(auth.CodePermissionDelegationDenied)
		}
		if !actorSuper {
			if err := validateRolePermissionSubset(database, principal, role); err != nil {
				return err
			}
		}
	}
	return nil
}

// ValidateHQPermissionDelegation bounds custom-role permissions to the actor's authority.
func ValidateHQPermissionDelegation(database *gorm.DB, principal *auth.WorkspacePrincipal, permissionIDs []string) error {
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeHeadquarters || len(permissionIDs) == 0 {
		return nil
	}
	super, err := IsHQSuperAdmin(database, principal)
	if err != nil || super {
		return err
	}
	var permissions []*gen.Permission
	err = database.Where("id IN ? AND scope = ?", permissionIDs, gen.PermissionScopeSystem).
		Where("is_delete IS NULL OR is_delete = ?", 1).Find(&permissions).Error
	if err != nil {
		return err
	}
	if len(permissions) != len(permissionIDs) {
		return auth.NewError(auth.CodePermissionDelegationDenied)
	}
	for _, permission := range permissions {
		if !principal.Has(permission.Action) {
			return auth.NewError(auth.CodePermissionDelegationDenied)
		}
	}
	return nil
}

// ValidateHQRoleUpdateDelegation checks the persisted permissions when an update leaves them unchanged.
func ValidateHQRoleUpdateDelegation(database *gorm.DB, principal *auth.WorkspacePrincipal, role *gen.OperatorRole, input map[string]interface{}) error {
	if value, replacing := input["permissionsIds"]; replacing && value != nil {
		return nil
	}
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return nil
	}
	super, err := IsHQSuperAdmin(database, principal)
	if err != nil || super {
		return err
	}
	return validateRolePermissionSubset(database, principal, role)
}

// EnsureRoleUnused rejects deleting a role assigned to an active membership record.
func EnsureRoleUnused(database *gorm.DB, role *gen.OperatorRole) error {
	var count int64
	err := database.Table("operator_membership_roles mr").
		Joins("JOIN operator_memberships m ON m.id = mr.operator_membership_id").
		Where("mr.operator_role_id = ?", role.ID).
		Where("m.status = ?", gen.MembershipStatusActive).
		Where("m.is_delete IS NULL OR m.is_delete = ?", 1).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return auth.NewError(auth.CodeRoleInUse)
	}
	return nil
}

func loadHQRoles(database *gorm.DB, organizationID string, roleIDs []string) ([]*gen.OperatorRole, error) {
	roles := make([]*gen.OperatorRole, 0, len(roleIDs))
	err := database.Where("id IN ? AND organization_id = ?", roleIDs, organizationID).
		Where("is_delete IS NULL OR is_delete = ?", 1).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	if len(roles) != len(roleIDs) {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	return roles, nil
}

func validateRolePermissionSubset(database *gorm.DB, principal *auth.WorkspacePrincipal, role *gen.OperatorRole) error {
	var permissions []*gen.Permission
	if err := database.Model(role).Association("Permissions").Find(&permissions); err != nil {
		return err
	}
	for _, permission := range permissions {
		if permission.Scope != gen.PermissionScopeSystem || !principal.Has(permission.Action) {
			return auth.NewError(auth.CodePermissionDelegationDenied)
		}
	}
	return nil
}

func loadMembershipRoles(database *gorm.DB, membership *gen.OperatorMembership) ([]*gen.OperatorRole, error) {
	var roles []*gen.OperatorRole
	err := database.Where("operator_roles.is_delete IS NULL OR operator_roles.is_delete = ?", 1).
		Model(membership).Association("Roles").Find(&roles)
	return roles, err
}
