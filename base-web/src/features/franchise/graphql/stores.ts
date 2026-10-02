import { gql } from '@/__generated__';

export const FRANCHISE_STORES_QUERY = gql(`
  query FranchiseStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType) {
    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data {
        id code name lifecycle rejectionReason organizationId
        contactPhone managerName managerPhone province city district address
        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter
      }
      total current_page per_page total_page
    }
  }
`);

export const CREATE_FRANCHISE_STORE_MUTATION = gql(`
  mutation FranchiseCreateStore($input: CreateStoreInput!) {
    createStore(input: $input) { id code name lifecycle organizationId }
  }
`);

export const UPDATE_FRANCHISE_STORE_MUTATION = gql(`
  mutation FranchiseUpdateStore($id: ID!, $input: UpdateStoreInput!) {
    updateStore(id: $id, input: $input) { id code name lifecycle organizationId }
  }
`);

export const SUBMIT_FRANCHISE_STORE_MUTATION = gql(`
  mutation FranchiseSubmitStore($id: ID!) {
    submitStore(id: $id) { id lifecycle submittedAt }
  }
`);

export const DELETE_FRANCHISE_STORES_MUTATION = gql(`
  mutation FranchiseDeleteStores($ids: [ID!]!) { deleteStores(id: $ids) }
`);
