package tools

import (
	"fmt"

	"base-engine/auth"
	"base-engine/src/services/ai"
)

// contractGapSpecs covers workspace-scoped GraphQL roots that were previously left
// as discovery placeholders despite being accessible through authorization handlers.
func contractGapSpecs() []ai.ToolSpec {
	type readRoot struct {
		name, root, fields, hqPermission, franchisePermission string
	}
	roots := []readRoot{
		{"Account", "account", "id phone displayName email status", "account:read", "operatorMembership:read"},
		{"AuditLog", "auditLog", "id action resourceType resourceId resultCode actorAccountId organizationId storeId createdAt", "auditLog:read", "tenantAudit:read"},
		{"MembershipInvitation", "membershipInvitation", "id membershipId invitedByAccountId expiresAt acceptedAt revokedAt createdAt", "hqMembership:read", "membershipInvitation:read"},
		{"OperatorMembership", "operatorMembership", "id status storeAccessMode accountId organizationId roles { id name kind } stores { id name lifecycle }", "hqMembership:read", "operatorMembership:read"},
		{"OperatorRole", "operatorRole", "id name kind organizationId permissions { id name action module scope }", "hqRole:read", "operatorRole:read"},
		{"Permission", "permission", "id name action module scope", "hqRole:read", "operatorRole:read"},
		{"Store", "store", "id code name lifecycle organizationId contactPhone managerName managerPhone province city district address businessHours businessStatus", "store:read_all", "store:read"},
	}
	result := make([]ai.ToolSpec, 0, 23)
	for _, item := range roots {
		for _, scope := range []struct {
			prefix, permission string
			workspace          auth.WorkspaceType
		}{{"Hq", item.hqPermission, auth.WorkspaceTypeHeadquarters}, {"Franchise", item.franchisePermission, auth.WorkspaceTypeFranchise}} {
			id := scope.prefix + item.name
			document := fmt.Sprintf("query %s($id: ID!) { %s(id: $id) { %s } }", id, item.root, item.fields)
			result = append(result, reviewedSpec(id, "graphql.query."+item.root, document, ai.ModeReadOnly, "LOW", scope.permission, scope.workspace, "", nil, []string{"id"}))
		}
	}
	hqDelete := reviewedSpec("HqDeleteMembershipInvitations", "graphql.mutation.deleteMembershipInvitations", `mutation HqDeleteMembershipInvitations($ids: [ID!]!) { deleteMembershipInvitations(id: $ids) }`, ai.ModeWrite, "HIGH", "hqMembership:delete", auth.WorkspaceTypeHeadquarters, ai.WriteExisting, []string{"ids"}, []string{"ids"})
	hqDelete.AdditionalPermissions = []string{"hqMembership:read"}
	franchiseDelete := reviewedSpec("FranchiseDeleteMembershipInvitations", "graphql.mutation.deleteMembershipInvitations", `mutation FranchiseDeleteMembershipInvitations($ids: [ID!]!) { deleteMembershipInvitations(id: $ids) }`, ai.ModeWrite, "HIGH", "membershipInvitation:delete", auth.WorkspaceTypeFranchise, ai.WriteExisting, []string{"ids"}, []string{"ids"})
	franchiseDelete.AdditionalPermissions = []string{"operatorMembership:read"}
	result = append(result,
		reviewedSpec("HqAccounts", "graphql.query.accounts", `query HqAccounts($page: Int!, $pageSize: Int!, $q: String) { accounts(current_page: $page, per_page: $pageSize, q: $q) { data { id phone displayName email status } total current_page per_page total_page } }`, ai.ModeReadOnly, "LOW", "account:read", auth.WorkspaceTypeHeadquarters, "", nil, []string{"page", "pageSize", "q"}),
		reviewedSpec("FranchiseAccounts", "graphql.query.accounts", `query FranchiseAccounts($page: Int!, $pageSize: Int!, $q: String) { accounts(current_page: $page, per_page: $pageSize, q: $q) { data { id phone displayName email status } total current_page per_page total_page } }`, ai.ModeReadOnly, "LOW", "operatorMembership:read", auth.WorkspaceTypeFranchise, "", nil, []string{"page", "pageSize", "q"}),
		reviewedSpec("HqMembershipInvitations", "graphql.query.membershipInvitations", `query HqMembershipInvitations($page: Int!, $pageSize: Int!, $filter: MembershipInvitationFilterType) { membershipInvitations(current_page: $page, per_page: $pageSize, filter: $filter) { data { id membershipId invitedByAccountId expiresAt acceptedAt revokedAt createdAt } total current_page per_page total_page } }`, ai.ModeReadOnly, "LOW", "hqMembership:read", auth.WorkspaceTypeHeadquarters, "", nil, []string{"page", "pageSize", "filter"}),
		reviewedSpec("FranchiseMembershipInvitations", "graphql.query.membershipInvitations", `query FranchiseMembershipInvitations($page: Int!, $pageSize: Int!, $filter: MembershipInvitationFilterType) { membershipInvitations(current_page: $page, per_page: $pageSize, filter: $filter) { data { id membershipId invitedByAccountId expiresAt acceptedAt revokedAt createdAt } total current_page per_page total_page } }`, ai.ModeReadOnly, "LOW", "membershipInvitation:read", auth.WorkspaceTypeFranchise, "", nil, []string{"page", "pageSize", "filter"}),
		hqDelete,
		franchiseDelete,
		reviewedSpec("HqDirectStore", "graphql.query.store", `query HqDirectStore($id: ID!) { store(id: $id) { id code name lifecycle organizationId contactPhone managerName managerPhone province city district address } }`, ai.ModeReadOnly, "LOW", "hqStore:read", auth.WorkspaceTypeHeadquarters, "", nil, []string{"id"}),
		reviewedSpec("HqPermissionByGrant", "graphql.query.permission", `query HqPermissionByGrant($id: ID!) { permission(id: $id) { id name action module scope } }`, ai.ModeReadOnly, "LOW", "permission:read", auth.WorkspaceTypeHeadquarters, "", nil, []string{"id"}),
	)
	return result
}
