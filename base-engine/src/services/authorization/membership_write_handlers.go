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

func createMembership(ctx context.Context, resolver *gen.GeneratedResolver, input map[string]interface{}) (*gen.OperatorMembership, error) {
	return nil, auth.NewError(auth.CodePermissionDenied)
}

func updateMembership(ctx context.Context, resolver *gen.GeneratedResolver, sessions *sessionservice.Service, id string, input map[string]interface{}) (*gen.OperatorMembership, error) {
	membership, principal, err := loadAuthorizedMembership(ctx, resolver, id, "update", AccessUpdate)
	if err != nil {
		return nil, err
	}
	if err := ValidateHQMembershipTarget(resolverDB(ctx, resolver), principal, membership); err != nil {
		auth.LogAuthorizationDenied(principal, membershipAuditAction(principal, "assign_roles"), "operatorMembership", membership.ID, err)
		return nil, err
	}
	secured, err := ensureInputOrganization(input, membership.OrganizationID)
	if err != nil {
		return nil, err
	}
	secured, err = secureMembershipInput(ctx, resolver, principal, secured, membership.OrganizationID)
	if err != nil {
		auth.LogAuthorizationDenied(principal, membershipAuditAction(principal, "assign_roles"), "operatorMembership", membership.ID, err)
		return nil, err
	}
	if err := enforceMembershipStoreMode(ctx, resolver, membership, secured); err != nil {
		return nil, err
	}
	if err := ensureOwnerRoleRemains(ctx, resolver, membership, secured); err != nil {
		return nil, err
	}
	for _, key := range []string{"status", "acceptedAt", "invitedAt", "accountId"} {
		delete(secured, key)
	}
	roleIDs := stringIDs(secured["rolesIds"])
	updated, err := gen.UpdateOperatorMembershipHandler(ctx, resolver, id, secured)
	if err != nil {
		return nil, err
	}
	return finishMembershipUpdate(ctx, resolver, sessions, principal, membership, updated, roleIDs)
}

func finishMembershipUpdate(ctx context.Context, resolver *gen.GeneratedResolver, sessions *sessionservice.Service, principal *auth.WorkspacePrincipal, membership, updated *gen.OperatorMembership, roleIDs []string) (*gen.OperatorMembership, error) {
	metadata := audit.Metadata{RoleIDs: roleIDs}
	if err := writeGeneratedAuditMetadata(ctx, resolver, principal, membership.OrganizationID, membershipAuditAction(principal, "assign_roles"), "operatorMembership", membership.ID, metadata); err != nil {
		return nil, err
	}
	if err := revokeMembershipSessions(ctx, resolver, sessions, principal, membership); err != nil {
		return nil, err
	}
	return updated, nil
}

func deleteMemberships(ctx context.Context, resolver *gen.GeneratedResolver, sessions *sessionservice.Service, ids []string, _ *bool) (bool, error) {
	memberships := make([]*gen.OperatorMembership, 0, len(ids))
	var principal *auth.WorkspacePrincipal
	for _, id := range ids {
		membership, actor, err := loadAuthorizedMembership(ctx, resolver, id, "delete", AccessDelete)
		if err != nil {
			return false, err
		}
		principal = actor
		if err := ValidateHQMembershipTarget(resolverDB(ctx, resolver), actor, membership); err != nil {
			auth.LogAuthorizationDenied(actor, membershipAuditAction(actor, "delete"), "operatorMembership", membership.ID, err)
			return false, err
		}
		memberships = append(memberships, membership)
	}
	if err := ensureOwnersRemainAfterDeletion(ctx, resolver, memberships); err != nil {
		return false, err
	}
	scoped := false
	done, err := gen.DeleteOperatorMembershipsHandler(ctx, resolver, ids, &scoped)
	if err != nil || !done {
		return done, err
	}
	if err := auditMemberships(ctx, resolver, memberships, "delete"); err != nil {
		return false, err
	}
	for _, membership := range memberships {
		if err := revokeMembershipSessions(ctx, resolver, sessions, principal, membership); err != nil {
			return false, err
		}
	}
	return true, nil
}

func revokeMembershipSessions(ctx context.Context, resolver *gen.GeneratedResolver, sessions *sessionservice.Service, principal *auth.WorkspacePrincipal, membership *gen.OperatorMembership) error {
	if sessions == nil || principal == nil || principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return nil
	}
	now := time.Now()
	ids, err := sessions.RevokeWorkspaceSessions(ctx, resolverDB(ctx, resolver), membership.AccountID, membership.OrganizationID, sessionservice.RevocationCodeAuthorityChanged, now)
	if err == nil {
		sessionservice.QueueRevokedSessions(ctx, ids, membership.OrganizationID, now)
	}
	return err
}

func recoverMemberships(ctx context.Context, resolver *gen.GeneratedResolver, ids []string) (bool, error) {
	return false, denyEntityWrite(ctx)
}

func auditMemberships(ctx context.Context, resolver *gen.GeneratedResolver, memberships []*gen.OperatorMembership, action string) error {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return err
	}
	action = membershipAuditAction(principal, action)
	for _, membership := range memberships {
		if err := writeGeneratedAudit(ctx, resolver, principal, membership.OrganizationID, action, "operatorMembership", membership.ID); err != nil {
			return err
		}
	}
	return nil
}

func membershipAuditAction(principal *auth.WorkspacePrincipal, verb string) string {
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		return "hqAdministrator:" + verb
	}
	if verb == "assign_roles" {
		return "operatorMembership:update"
	}
	return "operatorMembership:" + verb
}
