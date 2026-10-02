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

func registerAuditHandlers(handlers *gen.ResolutionHandlers) {
	handlers.CreateAuditLog = denyCreateAuditLog
	handlers.UpdateAuditLog = denyUpdateAuditLog
	handlers.DeleteAuditLogs = denyDeleteAuditLogs
	handlers.RecoveryAuditLogs = denyRecoveryAuditLogs
	handlers.QueryAuditLog = queryAuditLog
	handlers.QueryAuditLogs = queryAuditLogs
	handlers.AuditLogActorAccount = auditLogActorAccount
	handlers.AuditLogSession = auditLogSession
	handlers.AuditLogOrganization = auditLogOrganization
	handlers.AuditLogStore = auditLogStore
}

func denyCreateAuditLog(ctx context.Context, _ *gen.GeneratedResolver, _ map[string]interface{}) (*gen.AuditLog, error) {
	return nil, denyEntityWrite(ctx)
}

func denyUpdateAuditLog(ctx context.Context, _ *gen.GeneratedResolver, _ string, _ map[string]interface{}) (*gen.AuditLog, error) {
	return nil, denyEntityWrite(ctx)
}

func denyDeleteAuditLogs(ctx context.Context, _ *gen.GeneratedResolver, _ []string, _ *bool) (bool, error) {
	return false, denyEntityWrite(ctx)
}

func denyRecoveryAuditLogs(ctx context.Context, _ *gen.GeneratedResolver, _ []string) (bool, error) {
	return false, denyEntityWrite(ctx)
}

func auditScope(ctx context.Context, client *gen.AuditLogFilterType) (*auth.WorkspacePrincipal, *gen.AuditLogFilterType, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, nil, err
	}
	action := "tenantAudit:read"
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		action = "auditLog:read"
	}
	if err := Authorize(principal, Intent{Action: action, Mode: AccessRead}); err != nil {
		return nil, nil, err
	}
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		return principal, client, nil
	}
	organizationID, err := requireOrganization(principal)
	if err != nil {
		return nil, nil, err
	}
	scope := &gen.AuditLogFilterType{OrganizationID: &organizationID}
	filters := []*gen.AuditLogFilterType{scope}
	if client != nil {
		filters = append(filters, client)
	}
	return principal, &gen.AuditLogFilterType{And: filters}, nil
}

func queryAuditLog(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryAuditLogHandlerOptions) (*gen.AuditLog, error) {
	principal, filter, err := auditScope(ctx, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	item, err := gen.QueryAuditLogHandler(ctx, resolver, options)
	if err != nil {
		return nil, concealScopedNotFound(err, principal)
	}
	return item, nil
}

func queryAuditLogs(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryAuditLogsHandlerOptions) (*gen.AuditLogResultType, error) {
	_, filter, err := auditScope(ctx, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	return gen.QueryAuditLogsHandler(ctx, resolver, options)
}
