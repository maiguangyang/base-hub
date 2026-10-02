import { gql } from '@/__generated__';

export const HQ_STORES_QUERY = gql(`
  query HqStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {
    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data {
        id code name lifecycle organizationId
        organization { id name type }
        contactPhone managerName managerPhone province city district address
        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter
      }
      total current_page per_page total_page
    }
  }
`);

export const HQ_STORES_TOTAL_QUERY = gql(`
  query HqStoresTotal { stores(current_page: 1, per_page: 1) { total } }
`);
