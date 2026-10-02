/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"context"
	"fmt"

	"base-engine/auth"
	"base-engine/gen"
)

func secureRoleInput(ctx context.Context, resolver *gen.GeneratedResolver, principal *auth.WorkspacePrincipal, input map[string]interface{}, organizationID string) (map[string]interface{}, error) {
	secured, err := ensureInputOrganization(input, organizationID)
	if err != nil {
		return nil, err
	}
	for _, key := range []string{"permissions", "members", "membersIds"} {
		if secured[key] != nil {
			return nil, auth.NewError(auth.CodePermissionDenied)
		}
	}
	ids := stringIDs(secured["permissionsIds"])
	if len(ids) == 0 {
		return secured, nil
	}
	var count int64
	expected := workspacePermissionScope(principal.WorkspaceType)
	err = resolverDB(ctx, resolver).Model(&gen.Permission{}).
		Where("id IN ? AND scope = ?", ids, expected).
		Where("is_delete IS NULL OR is_delete = ?", 1).Count(&count).Error
	if err != nil || count != int64(len(ids)) {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	if err := ValidateHQPermissionDelegation(resolverDB(ctx, resolver), principal, ids); err != nil {
		return nil, err
	}
	return secured, nil
}

func stringIDs(value interface{}) []string {
	switch ids := value.(type) {
	case []string:
		return ids
	case []interface{}:
		result := make([]string, 0, len(ids))
		for _, id := range ids {
			result = append(result, fmt.Sprint(id))
		}
		return result
	default:
		return nil
	}
}

func roleOrganization(ctx context.Context, resolver *gen.GeneratedResolver, role *gen.OperatorRole) (*gen.Organization, error) {
	if _, _, err := loadAuthorizedRole(ctx, resolver, role.ID, "read", AccessRelation); err != nil {
		return nil, err
	}
	item := &gen.Organization{}
	return item, activeRecordCondition(resolver.DB.Query().WithContext(ctx)).First(item, "id = ?", role.OrganizationID).Error
}

func roleMembers(ctx context.Context, resolver *gen.GeneratedResolver, role *gen.OperatorRole) ([]*gen.OperatorMembership, error) {
	loaded, _, err := loadAuthorizedRole(ctx, resolver, role.ID, "read", AccessRelation)
	if err != nil {
		return nil, err
	}
	var items []*gen.OperatorMembership
	err = resolver.DB.Query().WithContext(ctx).Model(loaded).
		Where("organization_id = ? AND (operator_memberships.is_delete IS NULL OR operator_memberships.is_delete = ?)", loaded.OrganizationID, 1).
		Association("Members").Find(&items)
	return items, err
}

func rolePermissions(ctx context.Context, resolver *gen.GeneratedResolver, role *gen.OperatorRole) ([]*gen.Permission, error) {
	loaded, principal, err := loadAuthorizedRole(ctx, resolver, role.ID, "read", AccessRelation)
	if err != nil {
		return nil, err
	}
	var items []*gen.Permission
	if scope := dynamicPermissionScope(loaded.Kind, principal.WorkspaceType); scope != "" {
		err = activeRecordCondition(resolver.DB.Query().WithContext(ctx)).
			Where("scope = ?", scope).Order("weight ASC, created_at ASC").Find(&items).Error
		return items, err
	}
	expected := workspacePermissionScope(principal.WorkspaceType)
	err = resolver.DB.Query().WithContext(ctx).Model(loaded).
		Where("scope = ? AND (permissions.is_delete IS NULL OR permissions.is_delete = ?)", expected, 1).
		Association("Permissions").Find(&items)
	return items, err
}
