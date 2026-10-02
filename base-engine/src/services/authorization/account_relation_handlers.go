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

func loadScopedAccount(ctx context.Context, resolver *gen.GeneratedResolver, id string) (*gen.Account, error) {
	return queryAccount(ctx, resolver, gen.QueryAccountHandlerOptions{ID: &id})
}

func accountSessions(ctx context.Context, resolver *gen.GeneratedResolver, account *gen.Account) ([]*gen.Session, error) {
	if _, err := loadScopedAccount(ctx, resolver, account.ID); err != nil {
		return nil, err
	}
	principal, _, err := sessionScope(ctx, nil)
	if err != nil {
		return nil, err
	}
	var items []*gen.Session
	database := activeRecordCondition(resolverDB(ctx, resolver)).Where("account_id = ?", account.ID)
	if principal.WorkspaceType != auth.WorkspaceTypeHeadquarters || !principal.Has("session:read") {
		database = database.Where("account_id = ?", principal.AccountID)
	}
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise {
		database = database.Where("organization_id = ?", *principal.OrganizationID)
	}
	return items, database.Find(&items).Error
}

func accountReviewedStores(ctx context.Context, resolver *gen.GeneratedResolver, account *gen.Account) ([]*gen.Store, error) {
	if _, err := loadScopedAccount(ctx, resolver, account.ID); err != nil {
		return nil, err
	}
	principal, _, err := storeScope(ctx, nil)
	if err != nil {
		return nil, err
	}
	var items []*gen.Store
	database := activeRecordCondition(resolverDB(ctx, resolver)).Where("reviewed_by_account_id = ?", account.ID)
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise {
		database = database.Where("organization_id = ?", *principal.OrganizationID)
		if !principal.AllStores {
			database = database.Where("id IN ?", mapKeys(principal.StoreIDs))
		}
	}
	return items, database.Find(&items).Error
}

func accountSentInvitations(ctx context.Context, resolver *gen.GeneratedResolver, account *gen.Account) ([]*gen.MembershipInvitation, error) {
	if _, err := loadScopedAccount(ctx, resolver, account.ID); err != nil {
		return nil, err
	}
	principal, _, err := invitationScope(ctx, resolver, nil)
	if err != nil {
		return nil, err
	}
	membershipIDs, err := invitationMembershipIDs(ctx, resolver, principal)
	if err != nil {
		return nil, err
	}
	var items []*gen.MembershipInvitation
	database := activeRecordCondition(resolverDB(ctx, resolver))
	return items, database.Where("invited_by_account_id = ? AND membership_id IN ?", account.ID, membershipIDs).Find(&items).Error
}

func accountAuditLogs(ctx context.Context, resolver *gen.GeneratedResolver, account *gen.Account) ([]*gen.AuditLog, error) {
	if _, err := loadScopedAccount(ctx, resolver, account.ID); err != nil {
		return nil, err
	}
	principal, _, err := auditScope(ctx, nil)
	if err != nil {
		return nil, err
	}
	var items []*gen.AuditLog
	database := activeRecordCondition(resolverDB(ctx, resolver)).Where("actor_account_id = ?", account.ID)
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise {
		database = database.Where("organization_id = ?", *principal.OrganizationID)
	}
	return items, database.Find(&items).Error
}
