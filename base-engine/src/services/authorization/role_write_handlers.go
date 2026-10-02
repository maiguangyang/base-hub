/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"context"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	sessionservice "base-engine/src/services/session"
)

func createRole(ctx context.Context, resolver *gen.GeneratedResolver, input map[string]interface{}) (*gen.OperatorRole, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	organizationID, err := requireOrganization(principal)
	if err != nil {
		return nil, err
	}
	if err := Authorize(principal, Intent{Action: roleAction(principal, "create"), Mode: AccessCreate, ResourceOrganizationID: &organizationID}); err != nil {
		auth.LogAuthorizationDenied(principal, roleAction(principal, "create"), "operatorRole", "", err)
		return nil, err
	}
	secured, err := secureRoleInput(ctx, resolver, principal, input, organizationID)
	if err != nil {
		auth.LogAuthorizationDenied(principal, roleAction(principal, "create"), "operatorRole", "", err)
		return nil, err
	}
	secured["kind"] = gen.RoleKindCustom
	role, err := gen.CreateOperatorRoleHandler(ctx, resolver, secured)
	if err != nil {
		return nil, err
	}
	metadata := audit.Metadata{PermissionIDs: stringIDs(secured["permissionsIds"])}
	return role, writeGeneratedAuditMetadata(ctx, resolver, principal, role.OrganizationID, roleAction(principal, "create"), "operatorRole", role.ID, metadata)
}

func updateRole(ctx context.Context, resolver *gen.GeneratedResolver, sessions *sessionservice.Service, id string, input map[string]interface{}) (*gen.OperatorRole, error) {
	role, principal, err := loadAuthorizedRole(ctx, resolver, id, "update", AccessUpdate)
	if err != nil {
		return nil, err
	}
	if role.Kind != gen.RoleKindCustom {
		err := auth.NewError(auth.CodePermissionDenied)
		auth.LogAuthorizationDenied(principal, roleAction(principal, "update"), "operatorRole", role.ID, err)
		return nil, err
	}
	if err := ValidateHQRoleUpdateDelegation(resolverDB(ctx, resolver), principal, role, input); err != nil {
		auth.LogAuthorizationDenied(principal, roleAction(principal, "update"), "operatorRole", role.ID, err)
		return nil, err
	}
	secured, err := secureRoleInput(ctx, resolver, principal, input, role.OrganizationID)
	if err != nil {
		auth.LogAuthorizationDenied(principal, roleAction(principal, "update"), "operatorRole", role.ID, err)
		return nil, err
	}
	delete(secured, "kind")
	updated, err := gen.UpdateOperatorRoleHandler(ctx, resolver, id, secured)
	if err != nil {
		return nil, err
	}
	metadata := audit.Metadata{PermissionIDs: stringIDs(secured["permissionsIds"])}
	if err := writeGeneratedAuditMetadata(ctx, resolver, principal, role.OrganizationID, roleAction(principal, "update"), "operatorRole", role.ID, metadata); err != nil {
		return nil, err
	}
	if _, changed := input["permissionsIds"]; changed {
		if err := revokeRoleMemberSessions(ctx, resolver, sessions, principal, role); err != nil {
			return nil, err
		}
	}
	return updated, nil
}

func revokeRoleMemberSessions(ctx context.Context, resolver *gen.GeneratedResolver, sessions *sessionservice.Service, principal *auth.WorkspacePrincipal, role *gen.OperatorRole) error {
	if sessions == nil || principal == nil || principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return nil
	}
	var memberships []*gen.OperatorMembership
	err := resolverDB(ctx, resolver).Table("operator_memberships m").
		Joins("JOIN operator_membership_roles mr ON mr.operator_membership_id = m.id").
		Where("mr.operator_role_id = ? AND m.organization_id = ?", role.ID, role.OrganizationID).
		Where("m.status = ? AND (m.is_delete IS NULL OR m.is_delete = ?)", gen.MembershipStatusActive, 1).
		Find(&memberships).Error
	if err != nil {
		return err
	}
	for _, membership := range memberships {
		if err := revokeRoleMembershipSessions(ctx, resolver, sessions, membership); err != nil {
			return err
		}
	}
	return nil
}

func revokeRoleMembershipSessions(ctx context.Context, resolver *gen.GeneratedResolver, sessions *sessionservice.Service, membership *gen.OperatorMembership) error {
	now := time.Now()
	ids, err := sessions.RevokeWorkspaceSessions(ctx, resolverDB(ctx, resolver), membership.AccountID, membership.OrganizationID, sessionservice.RevocationCodeAuthorityChanged, now)
	if err == nil {
		sessionservice.QueueRevokedSessions(ctx, ids, membership.OrganizationID, now)
	}
	return err
}

func deleteRoles(ctx context.Context, resolver *gen.GeneratedResolver, ids []string, _ *bool) (bool, error) {
	roles, err := loadManagedRoles(ctx, resolver, ids, "delete", AccessDelete)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if err := EnsureRoleUnused(resolverDB(ctx, resolver), role); err != nil {
			return false, err
		}
	}
	scoped := false
	done, err := gen.DeleteOperatorRolesHandler(ctx, resolver, ids, &scoped)
	if err != nil || !done {
		return done, err
	}
	return done, auditRoles(ctx, resolver, roles, "delete")
}

func recoverRoles(ctx context.Context, resolver *gen.GeneratedResolver, ids []string) (bool, error) {
	return false, denyEntityWrite(ctx)
}

func loadManagedRoles(ctx context.Context, resolver *gen.GeneratedResolver, ids []string, verb string, mode AccessMode) ([]*gen.OperatorRole, error) {
	roles := make([]*gen.OperatorRole, 0, len(ids))
	for _, id := range ids {
		role, principal, err := loadAuthorizedRole(ctx, resolver, id, verb, mode)
		if err != nil {
			return nil, err
		}
		if role.Kind != gen.RoleKindCustom {
			err := auth.NewError(auth.CodePermissionDenied)
			auth.LogAuthorizationDenied(principal, roleAction(principal, verb), "operatorRole", role.ID, err)
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func auditRoles(ctx context.Context, resolver *gen.GeneratedResolver, roles []*gen.OperatorRole, action string) error {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return err
	}
	action = roleAction(principal, action)
	for _, role := range roles {
		if err := writeGeneratedAudit(ctx, resolver, principal, role.OrganizationID, action, "operatorRole", role.ID); err != nil {
			return err
		}
	}
	return nil
}
