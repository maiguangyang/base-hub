import { gql } from '@/__generated__';

export const HQ_DIRECT_STORES_QUERY = gql(`
  query HqDirectStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {
    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data {
        id code name lifecycle organizationId
        contactPhone managerName managerPhone province city district address
        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter
      }
      total current_page per_page total_page
    }
  }
`);

export const CREATE_DIRECT_STORE_MUTATION = gql(`
  mutation HqCreateDirectStore($input: CreateStoreInput!) {
    createStore(input: $input) {
      id code name lifecycle organizationId
      contactPhone managerName managerPhone province city district address
      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter
    }
  }
`);

export const UPDATE_DIRECT_STORE_MUTATION = gql(`
  mutation HqUpdateDirectStore($id: ID!, $input: UpdateStoreInput!) {
    updateStore(id: $id, input: $input) {
      id code name lifecycle organizationId
      contactPhone managerName managerPhone province city district address
      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter
    }
  }
`);

export const DELETE_DIRECT_STORES_MUTATION = gql(`
  mutation HqDeleteDirectStores($ids: [ID!]!) { deleteStores(id: $ids) }
`);
