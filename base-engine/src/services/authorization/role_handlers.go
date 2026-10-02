/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
)

func registerRoleHandlers(handlers *gen.ResolutionHandlers, dependencies HandlerDependencies) {
	handlers.CreatePermission = denyCreatePermission
	handlers.UpdatePermission = denyUpdatePermission
	handlers.DeletePermissions = denyDeletePermissions
	handlers.RecoveryPermissions = denyRecoveryPermissions
	handlers.QueryPermission = queryPermission
	handlers.QueryPermissions = queryPermissions
	handlers.PermissionRoles = permissionRoles
	handlers.CreateOperatorRole = createRole
	handlers.UpdateOperatorRole = func(ctx context.Context, resolver *gen.GeneratedResolver, id string, input map[string]interface{}) (*gen.OperatorRole, error) {
		return updateRole(ctx, resolver, dependencies.Sessions, id, input)
	}
	handlers.DeleteOperatorRoles = deleteRoles
	handlers.RecoveryOperatorRoles = recoverRoles
	handlers.QueryOperatorRole = queryRole
	handlers.QueryOperatorRoles = queryRoles
	handlers.OperatorRoleOrganization = roleOrganization
	handlers.OperatorRoleMembers = roleMembers
	handlers.OperatorRolePermissions = rolePermissions
}

func denyCreatePermission(ctx context.Context, _ *gen.GeneratedResolver, _ map[string]interface{}) (*gen.Permission, error) {
	return nil, denyEntityWrite(ctx)
}

func denyUpdatePermission(ctx context.Context, _ *gen.GeneratedResolver, _ string, _ map[string]interface{}) (*gen.Permission, error) {
	return nil, denyEntityWrite(ctx)
}

func denyDeletePermissions(ctx context.Context, _ *gen.GeneratedResolver, _ []string, _ *bool) (bool, error) {
	return false, denyEntityWrite(ctx)
}

func denyRecoveryPermissions(ctx context.Context, _ *gen.GeneratedResolver, _ []string) (bool, error) {
	return false, denyEntityWrite(ctx)
}

func permissionScope(ctx context.Context, client *gen.PermissionFilterType) (*auth.WorkspacePrincipal, *gen.PermissionFilterType, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, nil, err
	}
	scope := gen.PermissionScopeTenant
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		scope = gen.PermissionScopeSystem
		if !principal.Has("permission:read") && !principal.Has("hqRole:read") {
			return nil, nil, auth.NewError(auth.CodePermissionDenied)
		}
	} else if !principal.Has("operatorRole:read") {
		return nil, nil, auth.NewError(auth.CodePermissionDenied)
	}
	scopeFilter := &gen.PermissionFilterType{Scope: &scope}
	filters := []*gen.PermissionFilterType{scopeFilter}
	if client != nil {
		filters = append(filters, client)
	}
	return principal, &gen.PermissionFilterType{And: filters}, nil
}

func queryPermission(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryPermissionHandlerOptions) (*gen.Permission, error) {
	principal, filter, err := permissionScope(ctx, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	item, err := gen.QueryPermissionHandler(ctx, resolver, options)
	if err != nil {
		return nil, concealScopedNotFound(err, principal)
	}
	return item, nil
}

func queryPermissions(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryPermissionsHandlerOptions) (*gen.PermissionResultType, error) {
	_, filter, err := permissionScope(ctx, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	return gen.QueryPermissionsHandler(ctx, resolver, options)
}

func permissionRoles(ctx context.Context, resolver *gen.GeneratedResolver, permission *gen.Permission) ([]*gen.OperatorRole, error) {
	principal, _, err := permissionScope(ctx, &gen.PermissionFilterType{ID: &permission.ID})
	if err != nil {
		return nil, err
	}
	expected := workspacePermissionScope(principal.WorkspaceType)
	if permission.Scope != expected {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	var roles []*gen.OperatorRole
	database := resolver.DB.Query().WithContext(ctx).Model(permission).
		Where("operator_roles.is_delete IS NULL OR operator_roles.is_delete = ?", 1)
	if principal.OrganizationID != nil {
		database = database.Where("organization_id = ?", *principal.OrganizationID)
	}
	return roles, database.Association("Roles").Find(&roles)
}

func roleAction(principal *auth.WorkspacePrincipal, verb string) string {
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		return "hqRole:" + verb
	}
	return "operatorRole:" + verb
}

func roleScope(ctx context.Context, client *gen.OperatorRoleFilterType, verb string) (*auth.WorkspacePrincipal, *gen.OperatorRoleFilterType, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, nil, err
	}
	if err := Authorize(principal, Intent{Action: roleAction(principal, verb), Mode: AccessRead}); err != nil {
		return nil, nil, err
	}
	organizationID, err := requireOrganization(principal)
	if err != nil {
		return nil, nil, err
	}
	scope := &gen.OperatorRoleFilterType{OrganizationID: &organizationID}
	filters := []*gen.OperatorRoleFilterType{scope}
	if client != nil {
		filters = append(filters, client)
	}
	return principal, &gen.OperatorRoleFilterType{And: filters}, nil
}

func queryRole(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryOperatorRoleHandlerOptions) (*gen.OperatorRole, error) {
	principal, filter, err := roleScope(ctx, options.Filter, "read")
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	item, err := gen.QueryOperatorRoleHandler(ctx, resolver, options)
	if err != nil {
		return nil, concealScopedNotFound(err, principal)
	}
	return item, nil
}

func queryRoles(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryOperatorRolesHandlerOptions) (*gen.OperatorRoleResultType, error) {
	_, filter, err := roleScope(ctx, options.Filter, "read")
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	return gen.QueryOperatorRolesHandler(ctx, resolver, options)
}

func loadAuthorizedRole(ctx context.Context, resolver *gen.GeneratedResolver, id, verb string, mode AccessMode) (*gen.OperatorRole, *auth.WorkspacePrincipal, error) {
	role := &gen.OperatorRole{}
	if err := resolverDB(ctx, resolver).First(role, "id = ?", id).Error; err != nil {
		return nil, nil, auth.NewError(auth.CodePermissionDenied)
	}
	if err := authorizeDeletionState(role.IsDelete, mode); err != nil {
		return nil, nil, err
	}
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, nil, err
	}
	err = Authorize(principal, Intent{Action: roleAction(principal, verb), Mode: mode, ResourceOrganizationID: &role.OrganizationID})
	if err != nil {
		auth.LogAuthorizationDenied(principal, roleAction(principal, verb), "operatorRole", role.ID, err)
	}
	return role, principal, err
}
