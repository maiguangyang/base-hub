import { gql } from '@/__generated__';

export const HQ_STORE_APPROVALS_QUERY = gql(`
  query HqStoreApprovals($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {
    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data {
        id code name lifecycle submittedAt organizationId organization { id name }
        contactPhone managerName managerPhone province city district address
        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount
      }
      total current_page per_page total_page
    }
  }
`);

export const REVIEW_STORE_MUTATION = gql(`
  mutation HqReviewStore($input: ReviewStoreInput!) {
    reviewStore(input: $input) { id lifecycle rejectionReason reviewedAt }
  }
`);
