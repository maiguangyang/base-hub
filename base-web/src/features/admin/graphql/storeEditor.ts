import { gql } from '@/__generated__';

export const ADMIN_STORE_EDITOR_QUERY = gql(`
  query AdminStoreEditor($id: ID!) {
    store(id: $id) {
      id name lifecycle rejectionReason organizationId organization { type }
      contactPhone managerName managerPhone province city district address
      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter
      businessLicenseImageUrl otherDocumentImageUrl
    }
  }
`);

export const SET_STORE_DOCUMENT_MUTATION = gql(`
  mutation SetStoreDocument($storeId: ID!, $kind: StoreDocumentKind!, $attachmentId: ID!) {
    setStoreDocument(storeId: $storeId, kind: $kind, attachmentId: $attachmentId) {
      id businessLicenseImageUrl otherDocumentImageUrl
    }
  }
`);

export const REMOVE_STORE_DOCUMENT_MUTATION = gql(`
  mutation RemoveStoreDocument($storeId: ID!, $kind: StoreDocumentKind!) {
    removeStoreDocument(storeId: $storeId, kind: $kind) {
      id businessLicenseImageUrl otherDocumentImageUrl
    }
  }
`);
