import { gql } from '@/__generated__';

export const HQ_ROLES_QUERY = gql(`
  query HqRoles($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorRoleFilterType!) {
    operatorRoles(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data { id name kind organizationId permissions { id name action module scope } }
      total current_page per_page total_page
    }
  }
`);

export const SYSTEM_PERMISSIONS_QUERY = gql(`
  query SystemPermissions {
    permissions(current_page: 1, per_page: 200, filter: { scope: SYSTEM }) {
      data { id name action module scope }
      total
    }
  }
`);

export const CREATE_HQ_ROLE_MUTATION = gql(`
  mutation HqCreateRole($input: CreateOperatorRoleInput!) {
    createOperatorRole(input: $input) { id name kind organizationId }
  }
`);

export const UPDATE_HQ_ROLE_MUTATION = gql(`
  mutation HqUpdateRole($id: ID!, $input: UpdateOperatorRoleInput!) {
    updateOperatorRole(id: $id, input: $input) { id name kind organizationId }
  }
`);

export const DELETE_HQ_ROLES_MUTATION = gql(`
  mutation HqDeleteRoles($ids: [ID!]!) { deleteOperatorRoles(id: $ids) }
`);
