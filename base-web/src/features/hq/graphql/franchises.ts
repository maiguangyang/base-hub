import { gql } from '@/__generated__';

export const HQ_FRANCHISES_QUERY = gql(`
  query HqFranchises($page: Int!, $pageSize: Int!, $q: String, $filter: OrganizationFilterType!, $canResetPassword: Boolean! = false) {
    organizations(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data { id code name status initialAccountId @include(if: $canResetPassword) initialAccount @include(if: $canResetPassword) { phone } }
      total current_page per_page total_page
    }
  }
`);

export const HQ_FRANCHISE_INITIAL_ACCOUNT_CANDIDATES_QUERY = gql(`
  query HqFranchiseInitialAccountCandidates($id: ID!) {
    organization(id: $id) {
      id memberships { id status account { id phone displayName status } }
    }
  }
`);

export const RESET_FRANCHISE_INITIAL_PASSWORD_MUTATION = gql(`
  mutation HqResetFranchiseInitialPassword($organizationId: ID!) {
    resetFranchiseInitialPassword(organizationId: $organizationId) { accountId temporaryPassword }
  }
`);

export const PROVISION_FRANCHISE_MUTATION = gql(`
  mutation HqProvisionFranchise($input: ProvisionFranchiseInput!) {
    provisionFranchise(input: $input) {
      organization { id code name status }
      membership { id status }
      temporaryPassword invitationPending
    }
  }
`);

export const SUSPEND_ORGANIZATION_MUTATION = gql(`
  mutation HqSuspendOrganization($input: SuspendOrganizationInput!) {
    suspendOrganization(input: $input) { id status suspensionReasonCode }
  }
`);

export const RESTORE_ORGANIZATION_MUTATION = gql(`
  mutation HqRestoreOrganization($id: ID!) {
    restoreOrganization(id: $id) { id status }
  }
`);
