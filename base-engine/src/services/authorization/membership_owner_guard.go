/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"context"
	"sort"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ensureOwnersRemainAfterDeletion 按组织校验批量删除后仍有活跃老板。
func ensureOwnersRemainAfterDeletion(ctx context.Context, resolver *gen.GeneratedResolver, memberships []*gen.OperatorMembership) error {
	idsByOrganization := make(map[string][]string)
	for _, membership := range memberships {
		idsByOrganization[membership.OrganizationID] = append(idsByOrganization[membership.OrganizationID], membership.ID)
	}
	database := resolverDB(ctx, resolver)
	organizationIDs := make([]string, 0, len(idsByOrganization))
	for organizationID := range idsByOrganization {
		organizationIDs = append(organizationIDs, organizationID)
	}
	sort.Strings(organizationIDs)
	for _, organizationID := range organizationIDs {
		if err := lockOwnerOrganization(database, organizationID); err != nil {
			return err
		}
		roleKind, errorCode, err := privilegedRolePolicy(database, organizationID)
		if err != nil {
			return err
		}
		ids := idsByOrganization[organizationID]
		total, err := countActivePrivileged(database, organizationID, roleKind, nil)
		if err != nil {
			return err
		}
		deleting, err := countActivePrivileged(database, organizationID, roleKind, ids)
		if err != nil {
			return err
		}
		if deleting > 0 && total-deleting < 1 {
			return auth.NewError(errorCode)
		}
	}
	return nil
}

func countActivePrivileged(database *gorm.DB, organizationID string, roleKind gen.RoleKind, membershipIDs []string) (int64, error) {
	query := database.Model(&gen.OperatorMembership{}).
		Distinct("operator_memberships.id").
		Joins("JOIN operator_membership_roles mr ON mr.operator_membership_id = operator_memberships.id").
		Joins("JOIN operator_roles r ON r.id = mr.operator_role_id").
		Where("operator_memberships.organization_id = ?", organizationID).
		Where("operator_memberships.status = ? AND r.kind = ?", gen.MembershipStatusActive, roleKind).
		Where("operator_memberships.is_delete IS NULL OR operator_memberships.is_delete = ?", 1).
		Where("r.is_delete IS NULL OR r.is_delete = ?", 1)
	if membershipIDs != nil {
		query = query.Where("operator_memberships.id IN ?", membershipIDs)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

func ensureOwnerRoleRemains(ctx context.Context, resolver *gen.GeneratedResolver, membership *gen.OperatorMembership, input map[string]interface{}) error {
	value, replacing := input["rolesIds"]
	if !replacing || membership.Status != gen.MembershipStatusActive {
		return nil
	}
	database := resolverDB(ctx, resolver)
	if err := lockOwnerOrganization(database, membership.OrganizationID); err != nil {
		return err
	}
	roleKind, errorCode, err := privilegedRolePolicy(database, membership.OrganizationID)
	if err != nil {
		return err
	}
	current, err := countActivePrivileged(database, membership.OrganizationID, roleKind, []string{membership.ID})
	if err != nil || current == 0 {
		return err
	}
	if replacementIncludesPrivileged(database, membership.OrganizationID, roleKind, stringIDs(value)) {
		return nil
	}
	total, err := countActivePrivileged(database, membership.OrganizationID, roleKind, nil)
	if err != nil {
		return err
	}
	if total <= 1 {
		return auth.NewError(errorCode)
	}
	return nil
}

// EnsurePrivilegedMembershipRemains protects the final active owner or HQ super administrator.
func EnsurePrivilegedMembershipRemains(database *gorm.DB, membership *gen.OperatorMembership) error {
	if err := lockOwnerOrganization(database, membership.OrganizationID); err != nil {
		return err
	}
	roleKind, errorCode, err := privilegedRolePolicy(database, membership.OrganizationID)
	if err != nil {
		return err
	}
	current, err := countActivePrivileged(database, membership.OrganizationID, roleKind, []string{membership.ID})
	if err != nil || current == 0 {
		return err
	}
	total, err := countActivePrivileged(database, membership.OrganizationID, roleKind, nil)
	if err != nil {
		return err
	}
	if total <= 1 {
		return auth.NewError(errorCode)
	}
	return nil
}

func lockOwnerOrganization(database *gorm.DB, organizationID string) error {
	organization := &gen.Organization{}
	return database.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").First(organization, "id = ?", organizationID).Error
}

func replacementIncludesPrivileged(database *gorm.DB, organizationID string, roleKind gen.RoleKind, roleIDs []string) bool {
	if len(roleIDs) == 0 {
		return false
	}
	var count int64
	err := database.Model(&gen.OperatorRole{}).
		Where("id IN ? AND organization_id = ? AND kind = ?", roleIDs, organizationID, roleKind).
		Where("is_delete IS NULL OR is_delete = ?", 1).Count(&count).Error
	return err == nil && count > 0
}

func privilegedRolePolicy(database *gorm.DB, organizationID string) (gen.RoleKind, auth.Code, error) {
	organization := &gen.Organization{}
	if err := database.Select("id", "type").First(organization, "id = ?", organizationID).Error; err != nil {
		return "", "", err
	}
	if organization.Type == gen.OrganizationTypeHeadquarters {
		return gen.RoleKindHqSuperAdmin, auth.CodeLastHQSuperAdminRequired, nil
	}
	return gen.RoleKindFranchiseOwner, auth.CodeLastOwnerRequired, nil
}
