package tools

import (
	"base-engine/auth"
	"base-engine/src/services/ai"
)

func task12Specs() []ai.ToolSpec {
	return []ai.ToolSpec{
		reviewedSpec("FranchiseCreateRole", "graphql.mutation.createOperatorRole", `mutation FranchiseCreateRole($input: CreateOperatorRoleInput!) {
    createOperatorRole(input: $input) { id name kind organizationId }
  }`, ai.ModeWrite, "HIGH", "operatorRole:create", auth.WorkspaceTypeFranchise, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("FranchiseDeleteRoles", "graphql.mutation.deleteOperatorRoles", `mutation FranchiseDeleteRoles($ids: [ID!]!) { deleteOperatorRoles(id: $ids) }`, ai.ModeWrite, "HIGH", "operatorRole:delete", auth.WorkspaceTypeFranchise, ai.WriteExisting, []string{"ids"}, []string{"ids"}),
		reviewedSpec("FranchiseRoles", "graphql.query.operatorRoles", `query FranchiseRoles($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorRoleFilterType) {
    operatorRoles(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data { id name kind organizationId permissions { id name action module scope } }
      total current_page per_page total_page
    }
  }`, ai.ModeReadOnly, "LOW", "operatorRole:read", auth.WorkspaceTypeFranchise, "", nil, []string{"page", "pageSize", "q", "filter"}),
		reviewedSpec("FranchiseUpdateRole", "graphql.mutation.updateOperatorRole", `mutation FranchiseUpdateRole($id: ID!, $input: UpdateOperatorRoleInput!) {
    updateOperatorRole(id: $id, input: $input) { id name kind organizationId }
  }`, ai.ModeWrite, "HIGH", "operatorRole:update", auth.WorkspaceTypeFranchise, ai.WriteExisting, []string{"id"}, []string{"id", "input"}),
		reviewedSpec("TenantPermissions", "graphql.query.permissions", `query TenantPermissions($page:Int!){permissions(current_page:$page,per_page:50,filter:{scope:TENANT}){data{id name action module scope} total current_page per_page total_page}}`, ai.ModeReadOnly, "LOW", "operatorRole:read", auth.WorkspaceTypeFranchise, "", nil, []string{"page"}),
	}
}
