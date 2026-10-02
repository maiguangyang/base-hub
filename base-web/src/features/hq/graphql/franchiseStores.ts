import { gql } from '@/__generated__';

export const HQ_FRANCHISE_STORES_QUERY = gql(`
  query HqFranchiseStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {
    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data { id code name lifecycle businessStatus contactPhone organizationId organization { id name } }
      total current_page per_page total_page
    }
  }
`);
