import { gql } from '@/__generated__';

export const FRANCHISE_ROLES_QUERY = gql(`
  query FranchiseRoles($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorRoleFilterType) {
    operatorRoles(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data { id name kind organizationId permissions { id name action module scope } }
      total current_page per_page total_page
    }
  }
`);

export const TENANT_PERMISSIONS_QUERY = gql(`
  query TenantPermissions {
    permissions(current_page: 1, per_page: 200, filter: { scope: TENANT }) {
      data { id name action module scope }
      total
    }
  }
`);

export const CREATE_FRANCHISE_ROLE_MUTATION = gql(`
  mutation FranchiseCreateRole($input: CreateOperatorRoleInput!) {
    createOperatorRole(input: $input) { id name kind organizationId }
  }
`);

export const UPDATE_FRANCHISE_ROLE_MUTATION = gql(`
  mutation FranchiseUpdateRole($id: ID!, $input: UpdateOperatorRoleInput!) {
    updateOperatorRole(id: $id, input: $input) { id name kind organizationId }
  }
`);

export const DELETE_FRANCHISE_ROLES_MUTATION = gql(`
  mutation FranchiseDeleteRoles($ids: [ID!]!) { deleteOperatorRoles(id: $ids) }
`);
