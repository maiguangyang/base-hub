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

func membershipInvitations(ctx context.Context, resolver *gen.GeneratedResolver, membership *gen.OperatorMembership) ([]*gen.MembershipInvitation, error) {
	loaded, _, err := loadAuthorizedMembership(ctx, resolver, membership.ID, "read", AccessRelation)
	if err != nil {
		return nil, err
	}
	if _, _, err := invitationScope(ctx, resolver, nil); err != nil {
		return nil, err
	}
	var items []*gen.MembershipInvitation
	database := activeRecordCondition(resolverDB(ctx, resolver))
	return items, database.Where("membership_id = ?", loaded.ID).Find(&items).Error
}

func invitationMembership(ctx context.Context, resolver *gen.GeneratedResolver, invitation *gen.MembershipInvitation) (*gen.OperatorMembership, error) {
	loaded, err := queryInvitation(ctx, resolver, gen.QueryMembershipInvitationHandlerOptions{ID: &invitation.ID})
	if err != nil {
		return nil, err
	}
	membership, _, err := loadAuthorizedMembership(ctx, resolver, loaded.MembershipID, "read", AccessRelation)
	return membership, err
}

func invitationInvitedByAccount(ctx context.Context, resolver *gen.GeneratedResolver, invitation *gen.MembershipInvitation) (*gen.Account, error) {
	loaded, err := queryInvitation(ctx, resolver, gen.QueryMembershipInvitationHandlerOptions{ID: &invitation.ID})
	if err != nil {
		return nil, err
	}
	return loadScopedAccount(ctx, resolver, loaded.InvitedByAccountID)
}

func storeReviewedByAccount(ctx context.Context, resolver *gen.GeneratedResolver, store *gen.Store) (*gen.Account, error) {
	loaded, _, err := loadAuthorizedStore(ctx, resolver, store.ID, "read", AccessRelation)
	if err != nil || loaded.ReviewedByAccountID == nil {
		return nil, err
	}
	return loadScopedAccount(ctx, resolver, *loaded.ReviewedByAccountID)
}

func storeAuditLogs(ctx context.Context, resolver *gen.GeneratedResolver, store *gen.Store) ([]*gen.AuditLog, error) {
	loaded, _, err := loadAuthorizedStore(ctx, resolver, store.ID, "read", AccessRelation)
	if err != nil {
		return nil, err
	}
	principal, _, err := auditScope(ctx, nil)
	if err != nil {
		return nil, err
	}
	var items []*gen.AuditLog
	database := activeRecordCondition(resolverDB(ctx, resolver)).Where("store_id = ?", loaded.ID)
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise {
		database = database.Where("organization_id = ?", *principal.OrganizationID)
	}
	return items, database.Find(&items).Error
}

func auditLogActorAccount(ctx context.Context, resolver *gen.GeneratedResolver, auditLog *gen.AuditLog) (*gen.Account, error) {
	loaded, err := loadScopedAuditLog(ctx, resolver, auditLog.ID)
	if err != nil || loaded.ActorAccountID == nil {
		return nil, err
	}
	return loadScopedAccount(ctx, resolver, *loaded.ActorAccountID)
}

func auditLogSession(ctx context.Context, resolver *gen.GeneratedResolver, auditLog *gen.AuditLog) (*gen.Session, error) {
	loaded, err := loadScopedAuditLog(ctx, resolver, auditLog.ID)
	if err != nil || loaded.SessionID == nil {
		return nil, err
	}
	return querySession(ctx, resolver, gen.QuerySessionHandlerOptions{ID: loaded.SessionID})
}

func auditLogOrganization(ctx context.Context, resolver *gen.GeneratedResolver, auditLog *gen.AuditLog) (*gen.Organization, error) {
	loaded, err := loadScopedAuditLog(ctx, resolver, auditLog.ID)
	if err != nil || loaded.OrganizationID == nil {
		return nil, err
	}
	return queryOrganization(ctx, resolver, gen.QueryOrganizationHandlerOptions{ID: loaded.OrganizationID})
}

func auditLogStore(ctx context.Context, resolver *gen.GeneratedResolver, auditLog *gen.AuditLog) (*gen.Store, error) {
	loaded, err := loadScopedAuditLog(ctx, resolver, auditLog.ID)
	if err != nil || loaded.StoreID == nil {
		return nil, err
	}
	return queryStore(ctx, resolver, gen.QueryStoreHandlerOptions{ID: loaded.StoreID})
}

func loadScopedAuditLog(ctx context.Context, resolver *gen.GeneratedResolver, id string) (*gen.AuditLog, error) {
	return queryAuditLog(ctx, resolver, gen.QueryAuditLogHandlerOptions{ID: &id})
}
