package tools

import (
	"base-engine/auth"
	"base-engine/src/services/ai"
)

func task9Specs() []ai.ToolSpec {
	return []ai.ToolSpec{
		reviewedSpec("HqCreateRole", "graphql.mutation.createOperatorRole", `mutation HqCreateRole($input: CreateOperatorRoleInput!) {
    createOperatorRole(input: $input) { id name kind organizationId }
  }`, ai.ModeWrite, "HIGH", "hqRole:create", auth.WorkspaceTypeHeadquarters, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("HqDeleteRoles", "graphql.mutation.deleteOperatorRoles", `mutation HqDeleteRoles($ids: [ID!]!) { deleteOperatorRoles(id: $ids) }`, ai.ModeWrite, "HIGH", "hqRole:delete", auth.WorkspaceTypeHeadquarters, ai.WriteExisting, []string{"ids"}, []string{"ids"}),
		reviewedSpec("HqRoles", "graphql.query.operatorRoles", `query HqRoles($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorRoleFilterType!) {
    operatorRoles(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data { id name kind organizationId permissions { id name action module scope } }
      total current_page per_page total_page
    }
  }`, ai.ModeReadOnly, "LOW", "hqRole:read", auth.WorkspaceTypeHeadquarters, "", nil, []string{"page", "pageSize", "q", "filter"}),
		reviewedSpec("HqUpdateRole", "graphql.mutation.updateOperatorRole", `mutation HqUpdateRole($id: ID!, $input: UpdateOperatorRoleInput!) {
    updateOperatorRole(id: $id, input: $input) { id name kind organizationId }
  }`, ai.ModeWrite, "HIGH", "hqRole:update", auth.WorkspaceTypeHeadquarters, ai.WriteExisting, []string{"id"}, []string{"id", "input"}),
		reviewedSpec("SystemPermissions", "graphql.query.permissions", `query SystemPermissions($page:Int!){permissions(current_page:$page,per_page:50,filter:{scope:SYSTEM}){data{id name action module scope} total current_page per_page total_page}}`, ai.ModeReadOnly, "LOW", "hqRole:read", auth.WorkspaceTypeHeadquarters, "", nil, []string{"page"}),
	}
}
