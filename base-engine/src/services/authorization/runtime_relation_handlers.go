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

func organizationSessions(ctx context.Context, resolver *gen.GeneratedResolver, organization *gen.Organization) ([]*gen.Session, error) {
	if _, err := queryOrganization(ctx, resolver, gen.QueryOrganizationHandlerOptions{ID: &organization.ID}); err != nil {
		return nil, err
	}
	principal, _, err := sessionScope(ctx, nil)
	if err != nil {
		return nil, err
	}
	var items []*gen.Session
	database := activeRecordCondition(resolverDB(ctx, resolver)).Where("organization_id = ?", organization.ID)
	if principal.WorkspaceType != auth.WorkspaceTypeHeadquarters || !principal.Has("session:read") {
		database = database.Where("account_id = ?", principal.AccountID)
	}
	return items, database.Find(&items).Error
}

func organizationAuditLogs(ctx context.Context, resolver *gen.GeneratedResolver, organization *gen.Organization) ([]*gen.AuditLog, error) {
	if _, err := queryOrganization(ctx, resolver, gen.QueryOrganizationHandlerOptions{ID: &organization.ID}); err != nil {
		return nil, err
	}
	if _, _, err := auditScope(ctx, nil); err != nil {
		return nil, err
	}
	var items []*gen.AuditLog
	database := activeRecordCondition(resolverDB(ctx, resolver))
	return items, database.Where("organization_id = ?", organization.ID).Find(&items).Error
}

func sessionAccount(ctx context.Context, resolver *gen.GeneratedResolver, session *gen.Session) (*gen.Account, error) {
	loaded, err := querySession(ctx, resolver, gen.QuerySessionHandlerOptions{ID: &session.ID})
	if err != nil {
		return nil, err
	}
	return loadScopedAccount(ctx, resolver, loaded.AccountID)
}

func sessionOrganization(ctx context.Context, resolver *gen.GeneratedResolver, session *gen.Session) (*gen.Organization, error) {
	loaded, err := querySession(ctx, resolver, gen.QuerySessionHandlerOptions{ID: &session.ID})
	if err != nil || loaded.OrganizationID == nil {
		return nil, err
	}
	return queryOrganization(ctx, resolver, gen.QueryOrganizationHandlerOptions{ID: loaded.OrganizationID})
}

func sessionAuditLogs(ctx context.Context, resolver *gen.GeneratedResolver, session *gen.Session) ([]*gen.AuditLog, error) {
	loaded, err := querySession(ctx, resolver, gen.QuerySessionHandlerOptions{ID: &session.ID})
	if err != nil {
		return nil, err
	}
	principal, _, err := auditScope(ctx, nil)
	if err != nil {
		return nil, err
	}
	var items []*gen.AuditLog
	database := activeRecordCondition(resolverDB(ctx, resolver)).Where("session_id = ?", loaded.ID)
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise {
		database = database.Where("organization_id = ?", *principal.OrganizationID)
	}
	return items, database.Find(&items).Error
}
