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

func registerInvitationHandlers(handlers *gen.ResolutionHandlers) {
	handlers.CreateMembershipInvitation = denyCreateInvitation
	handlers.UpdateMembershipInvitation = denyUpdateInvitation
	handlers.DeleteMembershipInvitations = deleteInvitations
	handlers.RecoveryMembershipInvitations = denyRecoveryInvitations
	handlers.QueryMembershipInvitation = queryInvitation
	handlers.QueryMembershipInvitations = queryInvitations
	handlers.MembershipInvitationMembership = invitationMembership
	handlers.MembershipInvitationInvitedByAccount = invitationInvitedByAccount
}

func invitationMembershipIDs(ctx context.Context, resolver *gen.GeneratedResolver, principal *auth.WorkspacePrincipal) ([]string, error) {
	organizationID, err := requireOrganization(principal)
	if err != nil {
		return nil, err
	}
	database := activeRecordCondition(resolver.DB.Query().WithContext(ctx).Model(&gen.OperatorMembership{})).
		Where("organization_id = ?", organizationID)
	var ids []string
	return ids, database.Pluck("id", &ids).Error
}

func invitationScope(ctx context.Context, resolver *gen.GeneratedResolver, client *gen.MembershipInvitationFilterType) (*auth.WorkspacePrincipal, *gen.MembershipInvitationFilterType, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, nil, err
	}
	action := invitationAction(principal, "read")
	if err := Authorize(principal, Intent{Action: action, Mode: AccessRead}); err != nil {
		return nil, nil, err
	}
	ids, err := invitationMembershipIDs(ctx, resolver, principal)
	if err != nil {
		return nil, nil, err
	}
	scope := &gen.MembershipInvitationFilterType{MembershipIDIn: ids}
	filters := []*gen.MembershipInvitationFilterType{scope}
	if client != nil {
		filters = append(filters, client)
	}
	return principal, &gen.MembershipInvitationFilterType{And: filters}, nil
}

func queryInvitation(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryMembershipInvitationHandlerOptions) (*gen.MembershipInvitation, error) {
	principal, filter, err := invitationScope(ctx, resolver, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	item, err := gen.QueryMembershipInvitationHandler(ctx, resolver, options)
	if err != nil {
		return nil, concealScopedNotFound(err, principal)
	}
	return item, nil
}

func queryInvitations(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryMembershipInvitationsHandlerOptions) (*gen.MembershipInvitationResultType, error) {
	_, filter, err := invitationScope(ctx, resolver, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	return gen.QueryMembershipInvitationsHandler(ctx, resolver, options)
}

func denyCreateInvitation(ctx context.Context, _ *gen.GeneratedResolver, _ map[string]interface{}) (*gen.MembershipInvitation, error) {
	return nil, denyEntityWrite(ctx)
}

func deleteInvitations(ctx context.Context, resolver *gen.GeneratedResolver, ids []string, _ *bool) (bool, error) {
	type deletion struct {
		invitation *gen.MembershipInvitation
		membership *gen.OperatorMembership
		principal  *auth.WorkspacePrincipal
		action     string
	}
	deletions := make([]deletion, 0, len(ids))
	for _, id := range ids {
		invitation := &gen.MembershipInvitation{}
		if err := resolverDB(ctx, resolver).First(invitation, "id = ?", id).Error; err != nil {
			return false, auth.NewError(auth.CodePermissionDenied)
		}
		if err := authorizeDeletionState(invitation.IsDelete, AccessDelete); err != nil {
			return false, err
		}
		membership, principal, err := loadAuthorizedMembership(ctx, resolver, invitation.MembershipID, "read", AccessRead)
		if err != nil {
			return false, err
		}
		if err := Authorize(principal, Intent{Action: invitationAction(principal, "delete"), Mode: AccessDelete, ResourceOrganizationID: &membership.OrganizationID}); err != nil {
			return false, err
		}
		deletions = append(deletions, deletion{invitation: invitation, membership: membership, principal: principal, action: invitationAction(principal, "delete")})
	}
	scoped := false
	done, err := gen.DeleteMembershipInvitationsHandler(ctx, resolver, ids, &scoped)
	if err != nil || !done {
		return done, err
	}
	for _, item := range deletions {
		if err := writeGeneratedAudit(ctx, resolver, item.principal, item.membership.OrganizationID, item.action, "membershipInvitation", item.invitation.ID); err != nil {
			return false, err
		}
	}
	return true, nil
}

func invitationAction(principal *auth.WorkspacePrincipal, verb string) string {
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		return "hqMembership:" + verb
	}
	return "membershipInvitation:" + verb
}

func denyUpdateInvitation(ctx context.Context, _ *gen.GeneratedResolver, _ string, _ map[string]interface{}) (*gen.MembershipInvitation, error) {
	return nil, denyEntityWrite(ctx)
}

func denyRecoveryInvitations(ctx context.Context, _ *gen.GeneratedResolver, _ []string) (bool, error) {
	return false, denyEntityWrite(ctx)
}
