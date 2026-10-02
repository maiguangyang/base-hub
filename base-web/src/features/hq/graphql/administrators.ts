import { gql } from '@/__generated__';

export const HQ_ADMINISTRATORS_QUERY = gql(`
  query HqAdministrators($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorMembershipFilterType) {
    operatorMemberships(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data {
        id status storeAccessMode accountId organizationId updatedAt
        account { id phone displayName email }
        roles { id name kind }
      }
      total current_page per_page total_page
    }
  }
`);

export const HQ_ADMINISTRATOR_IDENTITY_QUERY = gql(`
  query HqAdministratorIdentity($accountId: ID!) {
    operatorMemberships(current_page: 1, per_page: 1, filter: { accountId: $accountId }) {
      data { id accountId roles { id kind } }
      total
    }
  }
`);

export const INVITE_HQ_ADMINISTRATOR_MUTATION = gql(`
  mutation HqInviteAdministrator($input: InviteOperatorInput!) {
    inviteOperator(input: $input) {
      membership { id status accountId organizationId }
      temporaryPassword invitationPending
    }
  }
`);

export const UPDATE_HQ_ADMINISTRATOR_MUTATION = gql(`
  mutation HqUpdateAdministrator($id: ID!, $input: UpdateOperatorMembershipInput!) {
    updateOperatorMembership(id: $id, input: $input) { id status rolesIds }
  }
`);

export const CHANGE_HQ_ADMINISTRATOR_STATUS_MUTATION = gql(`
  mutation HqChangeAdministratorStatus($input: ChangeMembershipStatusInput!) {
    changeMembershipStatus(input: $input) { id status }
  }
`);

export const DELETE_HQ_ADMINISTRATORS_MUTATION = gql(`
  mutation HqDeleteAdministrators($ids: [ID!]!) { deleteOperatorMemberships(id: $ids) }
`);

export const RESET_HQ_ADMINISTRATOR_PASSWORD_MUTATION = gql(`
  mutation HqResetAdministratorPassword($accountId: ID!) {
    resetTemporaryPassword(accountId: $accountId) { accountId temporaryPassword }
  }
`);
