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

func registerIdentityHandlers(handlers *gen.ResolutionHandlers) {
	handlers.CreateAccount = denyCreateAccount
	handlers.UpdateAccount = denyUpdateAccount
	handlers.DeleteAccounts = denyDeleteAccounts
	handlers.RecoveryAccounts = denyRecoveryAccounts
	handlers.QueryAccount = queryAccount
	handlers.QueryAccounts = queryAccounts
	handlers.AccountMemberships = accountMemberships
	handlers.AccountInitializedOrganizations = accountInitializedOrganizations
	handlers.AccountSessions = accountSessions
	handlers.AccountReviewedStores = accountReviewedStores
	handlers.AccountSentMembershipInvitations = accountSentInvitations
	handlers.AccountAuditLogs = accountAuditLogs
	handlers.CreateSession = denyCreateSession
	handlers.UpdateSession = denyUpdateSession
	handlers.DeleteSessions = denyDeleteSessions
	handlers.RecoverySessions = denyRecoverySessions
	handlers.QuerySession = querySession
	handlers.QuerySessions = querySessions
	handlers.SessionAccount = sessionAccount
	handlers.SessionOrganization = sessionOrganization
	handlers.SessionAuditLogs = sessionAuditLogs
}

func denyCreateAccount(ctx context.Context, _ *gen.GeneratedResolver, _ map[string]interface{}) (*gen.Account, error) {
	return nil, denyEntityWrite(ctx)
}

func denyUpdateAccount(ctx context.Context, _ *gen.GeneratedResolver, _ string, _ map[string]interface{}) (*gen.Account, error) {
	return nil, denyEntityWrite(ctx)
}

func denyDeleteAccounts(ctx context.Context, _ *gen.GeneratedResolver, _ []string, _ *bool) (bool, error) {
	return false, denyEntityWrite(ctx)
}

func denyRecoveryAccounts(ctx context.Context, _ *gen.GeneratedResolver, _ []string) (bool, error) {
	return false, denyEntityWrite(ctx)
}

func queryAccount(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryAccountHandlerOptions) (*gen.Account, error) {
	principal, filter, err := accountScope(ctx, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	item, err := gen.QueryAccountHandler(ctx, resolver, options)
	if err != nil {
		return nil, concealScopedNotFound(err, principal)
	}
	return item, nil
}

func queryAccounts(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryAccountsHandlerOptions) (*gen.AccountResultType, error) {
	_, filter, err := accountScope(ctx, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	return gen.QueryAccountsHandler(ctx, resolver, options)
}

func accountScope(ctx context.Context, client *gen.AccountFilterType) (*auth.WorkspacePrincipal, *gen.AccountFilterType, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, nil, err
	}
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		if err := Authorize(principal, Intent{Action: "account:read", Mode: AccessRead}); err != nil {
			return nil, nil, err
		}
		return principal, client, nil
	}
	if err := Authorize(principal, Intent{Action: "operatorMembership:read", Mode: AccessRead}); err != nil {
		return nil, nil, err
	}
	organizationID, err := requireOrganization(principal)
	if err != nil {
		return nil, nil, err
	}
	scope := &gen.AccountFilterType{Memberships: &gen.OperatorMembershipFilterType{OrganizationID: &organizationID}}
	return principal, &gen.AccountFilterType{And: compactAccountFilters(client, scope)}, nil
}

func compactAccountFilters(filters ...*gen.AccountFilterType) []*gen.AccountFilterType {
	result := make([]*gen.AccountFilterType, 0, len(filters))
	for _, filter := range filters {
		if filter != nil {
			result = append(result, filter)
		}
	}
	return result
}

func accountMemberships(ctx context.Context, resolver *gen.GeneratedResolver, account *gen.Account) ([]*gen.OperatorMembership, error) {
	principal, _, err := accountScope(ctx, nil)
	if err != nil {
		return nil, err
	}
	database := activeRecordCondition(resolver.DB.Query().WithContext(ctx)).Where("account_id = ?", account.ID)
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise {
		database = database.Where("organization_id = ?", *principal.OrganizationID)
	}
	var memberships []*gen.OperatorMembership
	return memberships, database.Find(&memberships).Error
}

func denyCreateSession(ctx context.Context, _ *gen.GeneratedResolver, _ map[string]interface{}) (*gen.Session, error) {
	return nil, denyEntityWrite(ctx)
}

func denyUpdateSession(ctx context.Context, _ *gen.GeneratedResolver, _ string, _ map[string]interface{}) (*gen.Session, error) {
	return nil, denyEntityWrite(ctx)
}

func denyDeleteSessions(ctx context.Context, _ *gen.GeneratedResolver, _ []string, _ *bool) (bool, error) {
	return false, denyEntityWrite(ctx)
}

func denyRecoverySessions(ctx context.Context, _ *gen.GeneratedResolver, _ []string) (bool, error) {
	return false, denyEntityWrite(ctx)
}

func querySession(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QuerySessionHandlerOptions) (*gen.Session, error) {
	principal, filter, err := sessionScope(ctx, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	item, err := gen.QuerySessionHandler(ctx, resolver, options)
	if err != nil {
		return nil, concealScopedNotFound(err, principal)
	}
	return item, nil
}

func querySessions(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QuerySessionsHandlerOptions) (*gen.SessionResultType, error) {
	_, filter, err := sessionScope(ctx, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	return gen.QuerySessionsHandler(ctx, resolver, options)
}

func sessionScope(ctx context.Context, client *gen.SessionFilterType) (*auth.WorkspacePrincipal, *gen.SessionFilterType, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, nil, err
	}
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters && principal.Has("session:read") {
		return principal, client, nil
	}
	scope := &gen.SessionFilterType{AccountID: &principal.AccountID}
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise {
		organizationID, err := requireOrganization(principal)
		if err != nil {
			return nil, nil, err
		}
		scope.OrganizationID = &organizationID
	}
	filters := []*gen.SessionFilterType{scope}
	if client != nil {
		filters = append(filters, client)
	}
	return principal, &gen.SessionFilterType{And: filters}, nil
}

func concealScopedNotFound(err error, _ *auth.WorkspacePrincipal) error {
	if _, ok := err.(*gen.NotFoundError); ok {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return err
}
