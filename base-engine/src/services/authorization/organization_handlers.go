/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"context"
	"strings"

	"base-engine/auth"
	"base-engine/gen"
)

func registerOrganizationHandlers(handlers *gen.ResolutionHandlers, dependencies HandlerDependencies) {
	handlers.CreateOrganization = denyCreateOrganization
	handlers.UpdateOrganization = confirmFranchiseInitialAccount
	handlers.DeleteOrganizations = denyDeleteOrganizations
	handlers.RecoveryOrganizations = denyRecoveryOrganizations
	handlers.QueryOrganization = queryOrganization
	handlers.QueryOrganizations = queryOrganizations
	handlers.OrganizationMemberships = organizationMemberships
	handlers.OrganizationInitialAccount = organizationInitialAccount
	handlers.OrganizationStores = organizationStores
	handlers.OrganizationRoles = organizationRoles
	handlers.OrganizationSessions = organizationSessions
	handlers.OrganizationAuditLogs = organizationAuditLogs
	registerMembershipHandlers(handlers, dependencies)
	registerInvitationHandlers(handlers)
}

func denyCreateOrganization(ctx context.Context, _ *gen.GeneratedResolver, _ map[string]interface{}) (*gen.Organization, error) {
	return nil, denyEntityWrite(ctx)
}

func denyDeleteOrganizations(ctx context.Context, _ *gen.GeneratedResolver, _ []string, _ *bool) (bool, error) {
	return false, denyEntityWrite(ctx)
}

func denyRecoveryOrganizations(ctx context.Context, _ *gen.GeneratedResolver, _ []string) (bool, error) {
	return false, denyEntityWrite(ctx)
}

func queryOrganization(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryOrganizationHandlerOptions) (*gen.Organization, error) {
	principal, filter, err := organizationScope(ctx, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	item, err := gen.QueryOrganizationHandler(ctx, resolver, options)
	if err != nil {
		return nil, concealScopedNotFound(err, principal)
	}
	return item, nil
}

func queryOrganizations(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryOrganizationsHandlerOptions) (*gen.OrganizationResultType, error) {
	_, filter, err := organizationScope(ctx, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	return gen.QueryOrganizationsHandler(ctx, resolver, options)
}

func organizationScope(ctx context.Context, client *gen.OrganizationFilterType) (*auth.WorkspacePrincipal, *gen.OrganizationFilterType, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, nil, err
	}
	action := "operatorMembership:read"
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		action = "organization:read"
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
	scope := &gen.OrganizationFilterType{ID: &organizationID}
	filters := []*gen.OrganizationFilterType{scope}
	if client != nil {
		filters = append(filters, client)
	}
	return principal, &gen.OrganizationFilterType{And: filters}, nil
}

func organizationMemberships(ctx context.Context, resolver *gen.GeneratedResolver, organization *gen.Organization) ([]*gen.OperatorMembership, error) {
	if err := authorizeOrganizationRelation(ctx, organization.ID, "operatorMembership:read", "hqMembership:read"); err != nil {
		return nil, err
	}
	var items []*gen.OperatorMembership
	database := activeRecordCondition(resolver.DB.Query().WithContext(ctx))
	return items, database.Where("organization_id = ?", organization.ID).Find(&items).Error
}

func organizationStores(ctx context.Context, resolver *gen.GeneratedResolver, organization *gen.Organization) ([]*gen.Store, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if err := Authorize(principal, Intent{Action: storeReadAction(principal, organization.ID), Mode: AccessRelation, ResourceOrganizationID: &organization.ID}); err != nil {
		return nil, err
	}
	var items []*gen.Store
	database := activeRecordCondition(resolver.DB.Query().WithContext(ctx))
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise && !principal.AllStores {
		database = database.Where("id IN ?", mapKeys(principal.StoreIDs))
	}
	return items, database.Where("organization_id = ?", organization.ID).Find(&items).Error
}

func organizationRoles(ctx context.Context, resolver *gen.GeneratedResolver, organization *gen.Organization) ([]*gen.OperatorRole, error) {
	if err := authorizeOrganizationRelation(ctx, organization.ID, "operatorRole:read", "hqRole:read"); err != nil {
		return nil, err
	}
	var items []*gen.OperatorRole
	database := activeRecordCondition(resolver.DB.Query().WithContext(ctx))
	return items, database.Where("organization_id = ?", organization.ID).Find(&items).Error
}

func authorizeOrganizationRelation(ctx context.Context, organizationID, tenantAction, systemAction string) error {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return err
	}
	action := tenantAction
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		action = systemAction
	}
	return Authorize(principal, Intent{Action: action, Mode: AccessRelation, ResourceOrganizationID: &organizationID})
}

func registerMembershipHandlers(handlers *gen.ResolutionHandlers, dependencies HandlerDependencies) {
	handlers.CreateOperatorMembership = createMembership
	handlers.UpdateOperatorMembership = func(ctx context.Context, resolver *gen.GeneratedResolver, id string, input map[string]interface{}) (*gen.OperatorMembership, error) {
		return updateMembership(ctx, resolver, dependencies.Sessions, id, input)
	}
	handlers.DeleteOperatorMemberships = func(ctx context.Context, resolver *gen.GeneratedResolver, ids []string, unscoped *bool) (bool, error) {
		return deleteMemberships(ctx, resolver, dependencies.Sessions, ids, unscoped)
	}
	handlers.RecoveryOperatorMemberships = recoverMemberships
	handlers.QueryOperatorMembership = queryMembership
	handlers.QueryOperatorMemberships = queryMemberships
	handlers.OperatorMembershipAccount = membershipAccount
	handlers.OperatorMembershipOrganization = membershipOrganization
	handlers.OperatorMembershipRoles = membershipRoles
	handlers.OperatorMembershipStores = membershipStores
	handlers.OperatorMembershipInvitations = membershipInvitations
}

func membershipAction(principal *auth.WorkspacePrincipal, verb string) string {
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		return "hqMembership:" + verb
	}
	return "operatorMembership:" + verb
}

func membershipScope(ctx context.Context, client *gen.OperatorMembershipFilterType, verb string) (*auth.WorkspacePrincipal, *gen.OperatorMembershipFilterType, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, nil, err
	}
	if err := Authorize(principal, Intent{Action: membershipAction(principal, verb), Mode: AccessRead}); err != nil {
		return nil, nil, err
	}
	organizationID, err := requireOrganization(principal)
	if err != nil {
		return nil, nil, err
	}
	scope := &gen.OperatorMembershipFilterType{OrganizationID: &organizationID}
	filters := []*gen.OperatorMembershipFilterType{scope}
	if client != nil {
		filters = append(filters, client)
	}
	return principal, &gen.OperatorMembershipFilterType{And: filters}, nil
}

func queryMembership(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryOperatorMembershipHandlerOptions) (*gen.OperatorMembership, error) {
	principal, filter, err := membershipScope(ctx, options.Filter, "read")
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	item, err := gen.QueryOperatorMembershipHandler(ctx, resolver, options)
	if err != nil {
		return nil, concealScopedNotFound(err, principal)
	}
	return item, nil
}

func queryMemberships(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryOperatorMembershipsHandlerOptions) (*gen.OperatorMembershipResultType, error) {
	clientFilter, roleID, storeID := splitMembershipRelationshipFilter(options.Filter)
	principal, filter, err := membershipScope(ctx, clientFilter, "read")
	if err != nil {
		return nil, err
	}
	if options.Q != nil && strings.TrimSpace(*options.Q) == "" {
		options.Q = nil
	}
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters && options.Q != nil {
		accountFilter, err := membershipAccountIdentityFilter(ctx, resolver, *options.Q)
		if err != nil {
			return nil, err
		}
		filter = &gen.OperatorMembershipFilterType{And: []*gen.OperatorMembershipFilterType{filter, accountFilter}}
		options.Q = nil
	}
	options.Filter = filter
	result, err := gen.QueryOperatorMembershipsHandler(ctx, resolver, options)
	if err != nil || roleID == nil && storeID == nil {
		return result, err
	}
	result.Filter = &membershipRelationshipFilter{base: result.Filter, roleID: roleID, storeID: storeID}
	return result, nil
}

func membershipAccountIdentityFilter(ctx context.Context, resolver *gen.GeneratedResolver, q string) (*gen.OperatorMembershipFilterType, error) {
	pattern := "%" + strings.TrimSpace(q) + "%"
	accountIDs := make([]string, 0)
	database := activeRecordCondition(resolverDB(ctx, resolver).Model(&gen.Account{}))
	err := database.Where("display_name LIKE ? OR phone LIKE ? OR email LIKE ?", pattern, pattern, pattern).
		Pluck("id", &accountIDs).Error
	if err != nil {
		return nil, err
	}
	if len(accountIDs) == 0 {
		missing := ""
		return &gen.OperatorMembershipFilterType{AccountID: &missing}, nil
	}
	return &gen.OperatorMembershipFilterType{AccountIDIn: accountIDs}, nil
}

func loadAuthorizedMembership(ctx context.Context, resolver *gen.GeneratedResolver, id, verb string, mode AccessMode) (*gen.OperatorMembership, *auth.WorkspacePrincipal, error) {
	membership := &gen.OperatorMembership{}
	if err := resolverDB(ctx, resolver).First(membership, "id = ?", id).Error; err != nil {
		return nil, nil, auth.NewError(auth.CodePermissionDenied)
	}
	if err := authorizeDeletionState(membership.IsDelete, mode); err != nil {
		return nil, nil, err
	}
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, nil, err
	}
	err = Authorize(principal, Intent{Action: membershipAction(principal, verb), Mode: mode, ResourceOrganizationID: &membership.OrganizationID})
	if err != nil {
		auth.LogAuthorizationDenied(principal, membershipAction(principal, verb), "operatorMembership", membership.ID, err)
	}
	return membership, principal, err
}
