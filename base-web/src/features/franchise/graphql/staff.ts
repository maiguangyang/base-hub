import { gql } from '@/__generated__';

export const FRANCHISE_STAFF_QUERY = gql(`
  query FranchiseStaff($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorMembershipFilterType) {
    operatorMemberships(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data {
        id status storeAccessMode accountId organizationId
        account { id phone displayName email }
        roles { id name kind }
        stores { id name lifecycle }
      }
      total current_page per_page total_page
    }
  }
`);

export const INVITE_FRANCHISE_STAFF_MUTATION = gql(`
  mutation FranchiseInviteStaff($input: InviteOperatorInput!) {
    inviteOperator(input: $input) {
      membership { id status accountId organizationId }
      temporaryPassword invitationPending
    }
  }
`);

export const UPDATE_FRANCHISE_STAFF_MUTATION = gql(`
  mutation FranchiseUpdateStaff($id: ID!, $input: UpdateOperatorMembershipInput!) {
    updateOperatorMembership(id: $id, input: $input) { id status storeAccessMode rolesIds storesIds }
  }
`);

export const CHANGE_MEMBERSHIP_STATUS_MUTATION = gql(`
  mutation FranchiseChangeMembershipStatus($input: ChangeMembershipStatusInput!) {
    changeMembershipStatus(input: $input) { id status }
  }
`);
