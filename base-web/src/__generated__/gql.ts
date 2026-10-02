/* eslint-disable */
import * as types from './graphql';
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 * Learn more about it here: https://the-guild.dev/graphql/codegen/plugins/presets/preset-client#reducing-bundle-size
 */
type Documents = {
    "\n  query AdminStoreEditor($id: ID!) {\n    store(id: $id) {\n      id name lifecycle rejectionReason organizationId organization { type }\n      contactPhone managerName managerPhone province city district address\n      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      businessLicenseImageUrl otherDocumentImageUrl\n    }\n  }\n": typeof types.AdminStoreEditorDocument,
    "\n  mutation SetStoreDocument($storeId: ID!, $kind: StoreDocumentKind!, $attachmentId: ID!) {\n    setStoreDocument(storeId: $storeId, kind: $kind, attachmentId: $attachmentId) {\n      id businessLicenseImageUrl otherDocumentImageUrl\n    }\n  }\n": typeof types.SetStoreDocumentDocument,
    "\n  mutation RemoveStoreDocument($storeId: ID!, $kind: StoreDocumentKind!) {\n    removeStoreDocument(storeId: $storeId, kind: $kind) {\n      id businessLicenseImageUrl otherDocumentImageUrl\n    }\n  }\n": typeof types.RemoveStoreDocumentDocument,
    "\n  fragment ViewerFields on Viewer {\n    account { id phone displayName email status mustChangePassword }\n    currentWorkspace { workspaceType organizationId organizationName homePath }\n    workspaces { workspaceType organizationId organizationName homePath }\n    permissions\n  }\n": typeof types.ViewerFieldsFragmentDoc,
    "\n  mutation Login($input: LoginInput!) {\n    login(input: $input) {\n      requiresPasswordChange\n      viewer { ...ViewerFields }\n    }\n  }\n": typeof types.LoginDocument,
    "\n  query Viewer {\n    viewer { ...ViewerFields }\n  }\n": typeof types.ViewerDocument,
    "\n  query Workspaces {\n    workspaces { workspaceType organizationId organizationName homePath }\n  }\n": typeof types.WorkspacesDocument,
    "\n  mutation ChangeTemporaryPassword($input: ChangePasswordInput!) {\n    changeTemporaryPassword(input: $input) { ...ViewerFields }\n  }\n": typeof types.ChangeTemporaryPasswordDocument,
    "\n  mutation SelectWorkspace($input: SelectWorkspaceInput!) {\n    selectWorkspace(input: $input) { ...ViewerFields }\n  }\n": typeof types.SelectWorkspaceDocument,
    "\n  mutation Logout { logout }\n": typeof types.LogoutDocument,
    "\n  subscription SessionEvents {\n    sessionEvents { code sessionId organizationId occurredAt }\n  }\n": typeof types.SessionEventsDocument,
    "\n  query PendingMembershipInvitations {\n    pendingMembershipInvitations { id membershipId expiresAt }\n  }\n": typeof types.PendingMembershipInvitationsDocument,
    "\n  mutation AcceptMembershipInvitation($id: ID!) {\n    acceptMembershipInvitation(id: $id) { id status organizationId }\n  }\n": typeof types.AcceptMembershipInvitationDocument,
    "\n  query FranchiseAuditLogs($page: Int!, $pageSize: Int!, $q: String, $filter: AuditLogFilterType) {\n    auditLogs(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id action resourceType resourceId resultCode actorAccountId storeId createdAt }\n      total current_page per_page total_page\n    }\n  }\n": typeof types.FranchiseAuditLogsDocument,
    "\n  query FranchiseRoles($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorRoleFilterType) {\n    operatorRoles(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id name kind organizationId permissions { id name action module scope } }\n      total current_page per_page total_page\n    }\n  }\n": typeof types.FranchiseRolesDocument,
    "\n  query TenantPermissions {\n    permissions(current_page: 1, per_page: 200, filter: { scope: TENANT }) {\n      data { id name action module scope }\n      total\n    }\n  }\n": typeof types.TenantPermissionsDocument,
    "\n  mutation FranchiseCreateRole($input: CreateOperatorRoleInput!) {\n    createOperatorRole(input: $input) { id name kind organizationId }\n  }\n": typeof types.FranchiseCreateRoleDocument,
    "\n  mutation FranchiseUpdateRole($id: ID!, $input: UpdateOperatorRoleInput!) {\n    updateOperatorRole(id: $id, input: $input) { id name kind organizationId }\n  }\n": typeof types.FranchiseUpdateRoleDocument,
    "\n  mutation FranchiseDeleteRoles($ids: [ID!]!) { deleteOperatorRoles(id: $ids) }\n": typeof types.FranchiseDeleteRolesDocument,
    "\n  query FranchiseStaff($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorMembershipFilterType) {\n    operatorMemberships(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id status storeAccessMode accountId organizationId\n        account { id phone displayName email }\n        roles { id name kind }\n        stores { id name lifecycle }\n      }\n      total current_page per_page total_page\n    }\n  }\n": typeof types.FranchiseStaffDocument,
    "\n  mutation FranchiseInviteStaff($input: InviteOperatorInput!) {\n    inviteOperator(input: $input) {\n      membership { id status accountId organizationId }\n      temporaryPassword invitationPending\n    }\n  }\n": typeof types.FranchiseInviteStaffDocument,
    "\n  mutation FranchiseUpdateStaff($id: ID!, $input: UpdateOperatorMembershipInput!) {\n    updateOperatorMembership(id: $id, input: $input) { id status storeAccessMode rolesIds storesIds }\n  }\n": typeof types.FranchiseUpdateStaffDocument,
    "\n  mutation FranchiseChangeMembershipStatus($input: ChangeMembershipStatusInput!) {\n    changeMembershipStatus(input: $input) { id status }\n  }\n": typeof types.FranchiseChangeMembershipStatusDocument,
    "query FranchiseStocktakes($storeId:ID!,$status:StocktakeStatus,$listingId:ID,$batchId:ID,$hasDifference:Boolean,$from:Time,$to:Time,$page:Int!,$perPage:Int!){franchiseStocktakes(storeId:$storeId,status:$status,listingId:$listingId,batchId:$batchId,hasDifference:$hasDifference,from:$from,to:$to,page:$page,perPage:$perPage,includeHistory:false){data{id status startedAt lines{id countedQuantity snapshotQuantity difference needsRecount}} total currentPage perPage}}": typeof types.FranchiseStocktakesDocument,
    "query FranchiseStocktake($storeId:ID!,$id:ID!){franchiseStocktake(storeId:$storeId,id:$id){id storeId status startedAt reviewedAt postedAt canceledAt initiatedByAccountId postedById addLineChoices{batchId packageId packageName packageEnabled} lines{id batchId listingId batchNumber expiresAt packageId packageName packageSetVersion packageEnabled countedQuantity snapshotQuantity difference reasonCode reasonNote countHistory{actorAccountId countedQuantity countedAt} needsRecount}}}": typeof types.FranchiseStocktakeDocument,
    "query FranchiseStocktakeBatchChoices($storeId:ID!,$listingId:ID!,$page:Int!,$perPage:Int!){franchiseStocktakeBatchChoices(storeId:$storeId,listingId:$listingId,page:$page,perPage:$perPage){data{id listingId batchNumber expiresAt} total currentPage perPage}}": typeof types.FranchiseStocktakeBatchChoicesDocument,
    "mutation FranchiseCreateStocktake($input:FranchiseCreateStocktakeInput!){franchiseCreateStocktake(input:$input){id status}}": typeof types.FranchiseCreateStocktakeDocument,
    "mutation FranchiseAddStocktakeLine($storeId:ID!,$id:ID!,$batchId:ID!,$packageId:ID!){franchiseAddStocktakeLine(storeId:$storeId,id:$id,batchId:$batchId,packageId:$packageId){id status}}": typeof types.FranchiseAddStocktakeLineDocument,
    "mutation FranchiseRecordStocktakeLine($storeId:ID!,$id:ID!,$lineId:ID!,$quantity:Int!){franchiseRecordStocktakeLine(storeId:$storeId,id:$id,lineId:$lineId,quantity:$quantity){id status}}": typeof types.FranchiseRecordStocktakeLineDocument,
    "mutation FranchiseSubmitStocktake($storeId:ID!,$id:ID!){franchiseSubmitStocktake(storeId:$storeId,id:$id){id status}}": typeof types.FranchiseSubmitStocktakeDocument,
    "mutation FranchiseSetStocktakeReason($storeId:ID!,$id:ID!,$lineId:ID!,$reasonCode:String!,$note:String){franchiseSetStocktakeReason(storeId:$storeId,id:$id,lineId:$lineId,reasonCode:$reasonCode,note:$note){id status}}": typeof types.FranchiseSetStocktakeReasonDocument,
    "mutation FranchiseReturnStocktake($storeId:ID!,$id:ID!){franchiseReturnStocktake(storeId:$storeId,id:$id){id status}}": typeof types.FranchiseReturnStocktakeDocument,
    "mutation FranchiseCancelStocktake($storeId:ID!,$id:ID!){franchiseCancelStocktake(storeId:$storeId,id:$id){id status}}": typeof types.FranchiseCancelStocktakeDocument,
    "mutation FranchisePostStocktake($storeId:ID!,$id:ID!){franchisePostStocktake(storeId:$storeId,id:$id){id status}}": typeof types.FranchisePostStocktakeDocument,
    "\n  query FranchiseStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle rejectionReason organizationId\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      }\n      total current_page per_page total_page\n    }\n  }\n": typeof types.FranchiseStoresDocument,
    "\n  mutation FranchiseCreateStore($input: CreateStoreInput!) {\n    createStore(input: $input) { id code name lifecycle organizationId }\n  }\n": typeof types.FranchiseCreateStoreDocument,
    "\n  mutation FranchiseUpdateStore($id: ID!, $input: UpdateStoreInput!) {\n    updateStore(id: $id, input: $input) { id code name lifecycle organizationId }\n  }\n": typeof types.FranchiseUpdateStoreDocument,
    "\n  mutation FranchiseSubmitStore($id: ID!) {\n    submitStore(id: $id) { id lifecycle submittedAt }\n  }\n": typeof types.FranchiseSubmitStoreDocument,
    "\n  mutation FranchiseDeleteStores($ids: [ID!]!) { deleteStores(id: $ids) }\n": typeof types.FranchiseDeleteStoresDocument,
    "\n  query HqAdministrators($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorMembershipFilterType) {\n    operatorMemberships(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id status storeAccessMode accountId organizationId updatedAt\n        account { id phone displayName email }\n        roles { id name kind }\n      }\n      total current_page per_page total_page\n    }\n  }\n": typeof types.HqAdministratorsDocument,
    "\n  query HqAdministratorIdentity($accountId: ID!) {\n    operatorMemberships(current_page: 1, per_page: 1, filter: { accountId: $accountId }) {\n      data { id accountId roles { id kind } }\n      total\n    }\n  }\n": typeof types.HqAdministratorIdentityDocument,
    "\n  mutation HqInviteAdministrator($input: InviteOperatorInput!) {\n    inviteOperator(input: $input) {\n      membership { id status accountId organizationId }\n      temporaryPassword invitationPending\n    }\n  }\n": typeof types.HqInviteAdministratorDocument,
    "\n  mutation HqUpdateAdministrator($id: ID!, $input: UpdateOperatorMembershipInput!) {\n    updateOperatorMembership(id: $id, input: $input) { id status rolesIds }\n  }\n": typeof types.HqUpdateAdministratorDocument,
    "\n  mutation HqChangeAdministratorStatus($input: ChangeMembershipStatusInput!) {\n    changeMembershipStatus(input: $input) { id status }\n  }\n": typeof types.HqChangeAdministratorStatusDocument,
    "\n  mutation HqDeleteAdministrators($ids: [ID!]!) { deleteOperatorMemberships(id: $ids) }\n": typeof types.HqDeleteAdministratorsDocument,
    "\n  mutation HqResetAdministratorPassword($accountId: ID!) {\n    resetTemporaryPassword(accountId: $accountId) { accountId temporaryPassword }\n  }\n": typeof types.HqResetAdministratorPasswordDocument,
    "\n  query HqAuditLogs($page: Int!, $pageSize: Int!, $q: String, $filter: AuditLogFilterType) {\n    auditLogs(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id action resourceType resourceId resultCode actorAccountId organizationId storeId createdAt }\n      total current_page per_page total_page\n    }\n  }\n": typeof types.HqAuditLogsDocument,
    "\n  query HqDirectStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle organizationId\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      }\n      total current_page per_page total_page\n    }\n  }\n": typeof types.HqDirectStoresDocument,
    "\n  mutation HqCreateDirectStore($input: CreateStoreInput!) {\n    createStore(input: $input) {\n      id code name lifecycle organizationId\n      contactPhone managerName managerPhone province city district address\n      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n    }\n  }\n": typeof types.HqCreateDirectStoreDocument,
    "\n  mutation HqUpdateDirectStore($id: ID!, $input: UpdateStoreInput!) {\n    updateStore(id: $id, input: $input) {\n      id code name lifecycle organizationId\n      contactPhone managerName managerPhone province city district address\n      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n    }\n  }\n": typeof types.HqUpdateDirectStoreDocument,
    "\n  mutation HqDeleteDirectStores($ids: [ID!]!) { deleteStores(id: $ids) }\n": typeof types.HqDeleteDirectStoresDocument,
    "\n  query HqFranchiseStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id code name lifecycle businessStatus contactPhone organizationId organization { id name } }\n      total current_page per_page total_page\n    }\n  }\n": typeof types.HqFranchiseStoresDocument,
    "\n  query HqFranchises($page: Int!, $pageSize: Int!, $q: String, $filter: OrganizationFilterType!, $canResetPassword: Boolean! = false) {\n    organizations(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id code name status initialAccountId @include(if: $canResetPassword) initialAccount @include(if: $canResetPassword) { phone } }\n      total current_page per_page total_page\n    }\n  }\n": typeof types.HqFranchisesDocument,
    "\n  query HqFranchiseInitialAccountCandidates($id: ID!) {\n    organization(id: $id) {\n      id memberships { id status account { id phone displayName status } }\n    }\n  }\n": typeof types.HqFranchiseInitialAccountCandidatesDocument,
    "\n  mutation HqResetFranchiseInitialPassword($organizationId: ID!) {\n    resetFranchiseInitialPassword(organizationId: $organizationId) { accountId temporaryPassword }\n  }\n": typeof types.HqResetFranchiseInitialPasswordDocument,
    "\n  mutation HqProvisionFranchise($input: ProvisionFranchiseInput!) {\n    provisionFranchise(input: $input) {\n      organization { id code name status }\n      membership { id status }\n      temporaryPassword invitationPending\n    }\n  }\n": typeof types.HqProvisionFranchiseDocument,
    "\n  mutation HqSuspendOrganization($input: SuspendOrganizationInput!) {\n    suspendOrganization(input: $input) { id status suspensionReasonCode }\n  }\n": typeof types.HqSuspendOrganizationDocument,
    "\n  mutation HqRestoreOrganization($id: ID!) {\n    restoreOrganization(id: $id) { id status }\n  }\n": typeof types.HqRestoreOrganizationDocument,
    "\n  query HqStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle organizationId\n        organization { id name type }\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      }\n      total current_page per_page total_page\n    }\n  }\n": typeof types.HqStoresDocument,
    "\n  query HqStoresTotal { stores(current_page: 1, per_page: 1) { total } }\n": typeof types.HqStoresTotalDocument,
    "\n  query HqRoles($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorRoleFilterType!) {\n    operatorRoles(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id name kind organizationId permissions { id name action module scope } }\n      total current_page per_page total_page\n    }\n  }\n": typeof types.HqRolesDocument,
    "\n  query SystemPermissions {\n    permissions(current_page: 1, per_page: 200, filter: { scope: SYSTEM }) {\n      data { id name action module scope }\n      total\n    }\n  }\n": typeof types.SystemPermissionsDocument,
    "\n  mutation HqCreateRole($input: CreateOperatorRoleInput!) {\n    createOperatorRole(input: $input) { id name kind organizationId }\n  }\n": typeof types.HqCreateRoleDocument,
    "\n  mutation HqUpdateRole($id: ID!, $input: UpdateOperatorRoleInput!) {\n    updateOperatorRole(id: $id, input: $input) { id name kind organizationId }\n  }\n": typeof types.HqUpdateRoleDocument,
    "\n  mutation HqDeleteRoles($ids: [ID!]!) { deleteOperatorRoles(id: $ids) }\n": typeof types.HqDeleteRolesDocument,
    "\n  query HqStoreApprovals($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle submittedAt organizationId organization { id name }\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount\n      }\n      total current_page per_page total_page\n    }\n  }\n": typeof types.HqStoreApprovalsDocument,
    "\n  mutation HqReviewStore($input: ReviewStoreInput!) {\n    reviewStore(input: $input) { id lifecycle rejectionReason reviewedAt }\n  }\n": typeof types.HqReviewStoreDocument,
};
const documents: Documents = {
    "\n  query AdminStoreEditor($id: ID!) {\n    store(id: $id) {\n      id name lifecycle rejectionReason organizationId organization { type }\n      contactPhone managerName managerPhone province city district address\n      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      businessLicenseImageUrl otherDocumentImageUrl\n    }\n  }\n": types.AdminStoreEditorDocument,
    "\n  mutation SetStoreDocument($storeId: ID!, $kind: StoreDocumentKind!, $attachmentId: ID!) {\n    setStoreDocument(storeId: $storeId, kind: $kind, attachmentId: $attachmentId) {\n      id businessLicenseImageUrl otherDocumentImageUrl\n    }\n  }\n": types.SetStoreDocumentDocument,
    "\n  mutation RemoveStoreDocument($storeId: ID!, $kind: StoreDocumentKind!) {\n    removeStoreDocument(storeId: $storeId, kind: $kind) {\n      id businessLicenseImageUrl otherDocumentImageUrl\n    }\n  }\n": types.RemoveStoreDocumentDocument,
    "\n  fragment ViewerFields on Viewer {\n    account { id phone displayName email status mustChangePassword }\n    currentWorkspace { workspaceType organizationId organizationName homePath }\n    workspaces { workspaceType organizationId organizationName homePath }\n    permissions\n  }\n": types.ViewerFieldsFragmentDoc,
    "\n  mutation Login($input: LoginInput!) {\n    login(input: $input) {\n      requiresPasswordChange\n      viewer { ...ViewerFields }\n    }\n  }\n": types.LoginDocument,
    "\n  query Viewer {\n    viewer { ...ViewerFields }\n  }\n": types.ViewerDocument,
    "\n  query Workspaces {\n    workspaces { workspaceType organizationId organizationName homePath }\n  }\n": types.WorkspacesDocument,
    "\n  mutation ChangeTemporaryPassword($input: ChangePasswordInput!) {\n    changeTemporaryPassword(input: $input) { ...ViewerFields }\n  }\n": types.ChangeTemporaryPasswordDocument,
    "\n  mutation SelectWorkspace($input: SelectWorkspaceInput!) {\n    selectWorkspace(input: $input) { ...ViewerFields }\n  }\n": types.SelectWorkspaceDocument,
    "\n  mutation Logout { logout }\n": types.LogoutDocument,
    "\n  subscription SessionEvents {\n    sessionEvents { code sessionId organizationId occurredAt }\n  }\n": types.SessionEventsDocument,
    "\n  query PendingMembershipInvitations {\n    pendingMembershipInvitations { id membershipId expiresAt }\n  }\n": types.PendingMembershipInvitationsDocument,
    "\n  mutation AcceptMembershipInvitation($id: ID!) {\n    acceptMembershipInvitation(id: $id) { id status organizationId }\n  }\n": types.AcceptMembershipInvitationDocument,
    "\n  query FranchiseAuditLogs($page: Int!, $pageSize: Int!, $q: String, $filter: AuditLogFilterType) {\n    auditLogs(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id action resourceType resourceId resultCode actorAccountId storeId createdAt }\n      total current_page per_page total_page\n    }\n  }\n": types.FranchiseAuditLogsDocument,
    "\n  query FranchiseRoles($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorRoleFilterType) {\n    operatorRoles(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id name kind organizationId permissions { id name action module scope } }\n      total current_page per_page total_page\n    }\n  }\n": types.FranchiseRolesDocument,
    "\n  query TenantPermissions {\n    permissions(current_page: 1, per_page: 200, filter: { scope: TENANT }) {\n      data { id name action module scope }\n      total\n    }\n  }\n": types.TenantPermissionsDocument,
    "\n  mutation FranchiseCreateRole($input: CreateOperatorRoleInput!) {\n    createOperatorRole(input: $input) { id name kind organizationId }\n  }\n": types.FranchiseCreateRoleDocument,
    "\n  mutation FranchiseUpdateRole($id: ID!, $input: UpdateOperatorRoleInput!) {\n    updateOperatorRole(id: $id, input: $input) { id name kind organizationId }\n  }\n": types.FranchiseUpdateRoleDocument,
    "\n  mutation FranchiseDeleteRoles($ids: [ID!]!) { deleteOperatorRoles(id: $ids) }\n": types.FranchiseDeleteRolesDocument,
    "\n  query FranchiseStaff($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorMembershipFilterType) {\n    operatorMemberships(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id status storeAccessMode accountId organizationId\n        account { id phone displayName email }\n        roles { id name kind }\n        stores { id name lifecycle }\n      }\n      total current_page per_page total_page\n    }\n  }\n": types.FranchiseStaffDocument,
    "\n  mutation FranchiseInviteStaff($input: InviteOperatorInput!) {\n    inviteOperator(input: $input) {\n      membership { id status accountId organizationId }\n      temporaryPassword invitationPending\n    }\n  }\n": types.FranchiseInviteStaffDocument,
    "\n  mutation FranchiseUpdateStaff($id: ID!, $input: UpdateOperatorMembershipInput!) {\n    updateOperatorMembership(id: $id, input: $input) { id status storeAccessMode rolesIds storesIds }\n  }\n": types.FranchiseUpdateStaffDocument,
    "\n  mutation FranchiseChangeMembershipStatus($input: ChangeMembershipStatusInput!) {\n    changeMembershipStatus(input: $input) { id status }\n  }\n": types.FranchiseChangeMembershipStatusDocument,
    "query FranchiseStocktakes($storeId:ID!,$status:StocktakeStatus,$listingId:ID,$batchId:ID,$hasDifference:Boolean,$from:Time,$to:Time,$page:Int!,$perPage:Int!){franchiseStocktakes(storeId:$storeId,status:$status,listingId:$listingId,batchId:$batchId,hasDifference:$hasDifference,from:$from,to:$to,page:$page,perPage:$perPage,includeHistory:false){data{id status startedAt lines{id countedQuantity snapshotQuantity difference needsRecount}} total currentPage perPage}}": types.FranchiseStocktakesDocument,
    "query FranchiseStocktake($storeId:ID!,$id:ID!){franchiseStocktake(storeId:$storeId,id:$id){id storeId status startedAt reviewedAt postedAt canceledAt initiatedByAccountId postedById addLineChoices{batchId packageId packageName packageEnabled} lines{id batchId listingId batchNumber expiresAt packageId packageName packageSetVersion packageEnabled countedQuantity snapshotQuantity difference reasonCode reasonNote countHistory{actorAccountId countedQuantity countedAt} needsRecount}}}": types.FranchiseStocktakeDocument,
    "query FranchiseStocktakeBatchChoices($storeId:ID!,$listingId:ID!,$page:Int!,$perPage:Int!){franchiseStocktakeBatchChoices(storeId:$storeId,listingId:$listingId,page:$page,perPage:$perPage){data{id listingId batchNumber expiresAt} total currentPage perPage}}": types.FranchiseStocktakeBatchChoicesDocument,
    "mutation FranchiseCreateStocktake($input:FranchiseCreateStocktakeInput!){franchiseCreateStocktake(input:$input){id status}}": types.FranchiseCreateStocktakeDocument,
    "mutation FranchiseAddStocktakeLine($storeId:ID!,$id:ID!,$batchId:ID!,$packageId:ID!){franchiseAddStocktakeLine(storeId:$storeId,id:$id,batchId:$batchId,packageId:$packageId){id status}}": types.FranchiseAddStocktakeLineDocument,
    "mutation FranchiseRecordStocktakeLine($storeId:ID!,$id:ID!,$lineId:ID!,$quantity:Int!){franchiseRecordStocktakeLine(storeId:$storeId,id:$id,lineId:$lineId,quantity:$quantity){id status}}": types.FranchiseRecordStocktakeLineDocument,
    "mutation FranchiseSubmitStocktake($storeId:ID!,$id:ID!){franchiseSubmitStocktake(storeId:$storeId,id:$id){id status}}": types.FranchiseSubmitStocktakeDocument,
    "mutation FranchiseSetStocktakeReason($storeId:ID!,$id:ID!,$lineId:ID!,$reasonCode:String!,$note:String){franchiseSetStocktakeReason(storeId:$storeId,id:$id,lineId:$lineId,reasonCode:$reasonCode,note:$note){id status}}": types.FranchiseSetStocktakeReasonDocument,
    "mutation FranchiseReturnStocktake($storeId:ID!,$id:ID!){franchiseReturnStocktake(storeId:$storeId,id:$id){id status}}": types.FranchiseReturnStocktakeDocument,
    "mutation FranchiseCancelStocktake($storeId:ID!,$id:ID!){franchiseCancelStocktake(storeId:$storeId,id:$id){id status}}": types.FranchiseCancelStocktakeDocument,
    "mutation FranchisePostStocktake($storeId:ID!,$id:ID!){franchisePostStocktake(storeId:$storeId,id:$id){id status}}": types.FranchisePostStocktakeDocument,
    "\n  query FranchiseStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle rejectionReason organizationId\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      }\n      total current_page per_page total_page\n    }\n  }\n": types.FranchiseStoresDocument,
    "\n  mutation FranchiseCreateStore($input: CreateStoreInput!) {\n    createStore(input: $input) { id code name lifecycle organizationId }\n  }\n": types.FranchiseCreateStoreDocument,
    "\n  mutation FranchiseUpdateStore($id: ID!, $input: UpdateStoreInput!) {\n    updateStore(id: $id, input: $input) { id code name lifecycle organizationId }\n  }\n": types.FranchiseUpdateStoreDocument,
    "\n  mutation FranchiseSubmitStore($id: ID!) {\n    submitStore(id: $id) { id lifecycle submittedAt }\n  }\n": types.FranchiseSubmitStoreDocument,
    "\n  mutation FranchiseDeleteStores($ids: [ID!]!) { deleteStores(id: $ids) }\n": types.FranchiseDeleteStoresDocument,
    "\n  query HqAdministrators($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorMembershipFilterType) {\n    operatorMemberships(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id status storeAccessMode accountId organizationId updatedAt\n        account { id phone displayName email }\n        roles { id name kind }\n      }\n      total current_page per_page total_page\n    }\n  }\n": types.HqAdministratorsDocument,
    "\n  query HqAdministratorIdentity($accountId: ID!) {\n    operatorMemberships(current_page: 1, per_page: 1, filter: { accountId: $accountId }) {\n      data { id accountId roles { id kind } }\n      total\n    }\n  }\n": types.HqAdministratorIdentityDocument,
    "\n  mutation HqInviteAdministrator($input: InviteOperatorInput!) {\n    inviteOperator(input: $input) {\n      membership { id status accountId organizationId }\n      temporaryPassword invitationPending\n    }\n  }\n": types.HqInviteAdministratorDocument,
    "\n  mutation HqUpdateAdministrator($id: ID!, $input: UpdateOperatorMembershipInput!) {\n    updateOperatorMembership(id: $id, input: $input) { id status rolesIds }\n  }\n": types.HqUpdateAdministratorDocument,
    "\n  mutation HqChangeAdministratorStatus($input: ChangeMembershipStatusInput!) {\n    changeMembershipStatus(input: $input) { id status }\n  }\n": types.HqChangeAdministratorStatusDocument,
    "\n  mutation HqDeleteAdministrators($ids: [ID!]!) { deleteOperatorMemberships(id: $ids) }\n": types.HqDeleteAdministratorsDocument,
    "\n  mutation HqResetAdministratorPassword($accountId: ID!) {\n    resetTemporaryPassword(accountId: $accountId) { accountId temporaryPassword }\n  }\n": types.HqResetAdministratorPasswordDocument,
    "\n  query HqAuditLogs($page: Int!, $pageSize: Int!, $q: String, $filter: AuditLogFilterType) {\n    auditLogs(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id action resourceType resourceId resultCode actorAccountId organizationId storeId createdAt }\n      total current_page per_page total_page\n    }\n  }\n": types.HqAuditLogsDocument,
    "\n  query HqDirectStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle organizationId\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      }\n      total current_page per_page total_page\n    }\n  }\n": types.HqDirectStoresDocument,
    "\n  mutation HqCreateDirectStore($input: CreateStoreInput!) {\n    createStore(input: $input) {\n      id code name lifecycle organizationId\n      contactPhone managerName managerPhone province city district address\n      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n    }\n  }\n": types.HqCreateDirectStoreDocument,
    "\n  mutation HqUpdateDirectStore($id: ID!, $input: UpdateStoreInput!) {\n    updateStore(id: $id, input: $input) {\n      id code name lifecycle organizationId\n      contactPhone managerName managerPhone province city district address\n      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n    }\n  }\n": types.HqUpdateDirectStoreDocument,
    "\n  mutation HqDeleteDirectStores($ids: [ID!]!) { deleteStores(id: $ids) }\n": types.HqDeleteDirectStoresDocument,
    "\n  query HqFranchiseStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id code name lifecycle businessStatus contactPhone organizationId organization { id name } }\n      total current_page per_page total_page\n    }\n  }\n": types.HqFranchiseStoresDocument,
    "\n  query HqFranchises($page: Int!, $pageSize: Int!, $q: String, $filter: OrganizationFilterType!, $canResetPassword: Boolean! = false) {\n    organizations(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id code name status initialAccountId @include(if: $canResetPassword) initialAccount @include(if: $canResetPassword) { phone } }\n      total current_page per_page total_page\n    }\n  }\n": types.HqFranchisesDocument,
    "\n  query HqFranchiseInitialAccountCandidates($id: ID!) {\n    organization(id: $id) {\n      id memberships { id status account { id phone displayName status } }\n    }\n  }\n": types.HqFranchiseInitialAccountCandidatesDocument,
    "\n  mutation HqResetFranchiseInitialPassword($organizationId: ID!) {\n    resetFranchiseInitialPassword(organizationId: $organizationId) { accountId temporaryPassword }\n  }\n": types.HqResetFranchiseInitialPasswordDocument,
    "\n  mutation HqProvisionFranchise($input: ProvisionFranchiseInput!) {\n    provisionFranchise(input: $input) {\n      organization { id code name status }\n      membership { id status }\n      temporaryPassword invitationPending\n    }\n  }\n": types.HqProvisionFranchiseDocument,
    "\n  mutation HqSuspendOrganization($input: SuspendOrganizationInput!) {\n    suspendOrganization(input: $input) { id status suspensionReasonCode }\n  }\n": types.HqSuspendOrganizationDocument,
    "\n  mutation HqRestoreOrganization($id: ID!) {\n    restoreOrganization(id: $id) { id status }\n  }\n": types.HqRestoreOrganizationDocument,
    "\n  query HqStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle organizationId\n        organization { id name type }\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      }\n      total current_page per_page total_page\n    }\n  }\n": types.HqStoresDocument,
    "\n  query HqStoresTotal { stores(current_page: 1, per_page: 1) { total } }\n": types.HqStoresTotalDocument,
    "\n  query HqRoles($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorRoleFilterType!) {\n    operatorRoles(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id name kind organizationId permissions { id name action module scope } }\n      total current_page per_page total_page\n    }\n  }\n": types.HqRolesDocument,
    "\n  query SystemPermissions {\n    permissions(current_page: 1, per_page: 200, filter: { scope: SYSTEM }) {\n      data { id name action module scope }\n      total\n    }\n  }\n": types.SystemPermissionsDocument,
    "\n  mutation HqCreateRole($input: CreateOperatorRoleInput!) {\n    createOperatorRole(input: $input) { id name kind organizationId }\n  }\n": types.HqCreateRoleDocument,
    "\n  mutation HqUpdateRole($id: ID!, $input: UpdateOperatorRoleInput!) {\n    updateOperatorRole(id: $id, input: $input) { id name kind organizationId }\n  }\n": types.HqUpdateRoleDocument,
    "\n  mutation HqDeleteRoles($ids: [ID!]!) { deleteOperatorRoles(id: $ids) }\n": types.HqDeleteRolesDocument,
    "\n  query HqStoreApprovals($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle submittedAt organizationId organization { id name }\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount\n      }\n      total current_page per_page total_page\n    }\n  }\n": types.HqStoreApprovalsDocument,
    "\n  mutation HqReviewStore($input: ReviewStoreInput!) {\n    reviewStore(input: $input) { id lifecycle rejectionReason reviewedAt }\n  }\n": types.HqReviewStoreDocument,
};

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = gql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function gql(source: string): unknown;

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query AdminStoreEditor($id: ID!) {\n    store(id: $id) {\n      id name lifecycle rejectionReason organizationId organization { type }\n      contactPhone managerName managerPhone province city district address\n      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      businessLicenseImageUrl otherDocumentImageUrl\n    }\n  }\n"): (typeof documents)["\n  query AdminStoreEditor($id: ID!) {\n    store(id: $id) {\n      id name lifecycle rejectionReason organizationId organization { type }\n      contactPhone managerName managerPhone province city district address\n      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      businessLicenseImageUrl otherDocumentImageUrl\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation SetStoreDocument($storeId: ID!, $kind: StoreDocumentKind!, $attachmentId: ID!) {\n    setStoreDocument(storeId: $storeId, kind: $kind, attachmentId: $attachmentId) {\n      id businessLicenseImageUrl otherDocumentImageUrl\n    }\n  }\n"): (typeof documents)["\n  mutation SetStoreDocument($storeId: ID!, $kind: StoreDocumentKind!, $attachmentId: ID!) {\n    setStoreDocument(storeId: $storeId, kind: $kind, attachmentId: $attachmentId) {\n      id businessLicenseImageUrl otherDocumentImageUrl\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation RemoveStoreDocument($storeId: ID!, $kind: StoreDocumentKind!) {\n    removeStoreDocument(storeId: $storeId, kind: $kind) {\n      id businessLicenseImageUrl otherDocumentImageUrl\n    }\n  }\n"): (typeof documents)["\n  mutation RemoveStoreDocument($storeId: ID!, $kind: StoreDocumentKind!) {\n    removeStoreDocument(storeId: $storeId, kind: $kind) {\n      id businessLicenseImageUrl otherDocumentImageUrl\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  fragment ViewerFields on Viewer {\n    account { id phone displayName email status mustChangePassword }\n    currentWorkspace { workspaceType organizationId organizationName homePath }\n    workspaces { workspaceType organizationId organizationName homePath }\n    permissions\n  }\n"): (typeof documents)["\n  fragment ViewerFields on Viewer {\n    account { id phone displayName email status mustChangePassword }\n    currentWorkspace { workspaceType organizationId organizationName homePath }\n    workspaces { workspaceType organizationId organizationName homePath }\n    permissions\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation Login($input: LoginInput!) {\n    login(input: $input) {\n      requiresPasswordChange\n      viewer { ...ViewerFields }\n    }\n  }\n"): (typeof documents)["\n  mutation Login($input: LoginInput!) {\n    login(input: $input) {\n      requiresPasswordChange\n      viewer { ...ViewerFields }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query Viewer {\n    viewer { ...ViewerFields }\n  }\n"): (typeof documents)["\n  query Viewer {\n    viewer { ...ViewerFields }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query Workspaces {\n    workspaces { workspaceType organizationId organizationName homePath }\n  }\n"): (typeof documents)["\n  query Workspaces {\n    workspaces { workspaceType organizationId organizationName homePath }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation ChangeTemporaryPassword($input: ChangePasswordInput!) {\n    changeTemporaryPassword(input: $input) { ...ViewerFields }\n  }\n"): (typeof documents)["\n  mutation ChangeTemporaryPassword($input: ChangePasswordInput!) {\n    changeTemporaryPassword(input: $input) { ...ViewerFields }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation SelectWorkspace($input: SelectWorkspaceInput!) {\n    selectWorkspace(input: $input) { ...ViewerFields }\n  }\n"): (typeof documents)["\n  mutation SelectWorkspace($input: SelectWorkspaceInput!) {\n    selectWorkspace(input: $input) { ...ViewerFields }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation Logout { logout }\n"): (typeof documents)["\n  mutation Logout { logout }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  subscription SessionEvents {\n    sessionEvents { code sessionId organizationId occurredAt }\n  }\n"): (typeof documents)["\n  subscription SessionEvents {\n    sessionEvents { code sessionId organizationId occurredAt }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query PendingMembershipInvitations {\n    pendingMembershipInvitations { id membershipId expiresAt }\n  }\n"): (typeof documents)["\n  query PendingMembershipInvitations {\n    pendingMembershipInvitations { id membershipId expiresAt }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation AcceptMembershipInvitation($id: ID!) {\n    acceptMembershipInvitation(id: $id) { id status organizationId }\n  }\n"): (typeof documents)["\n  mutation AcceptMembershipInvitation($id: ID!) {\n    acceptMembershipInvitation(id: $id) { id status organizationId }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query FranchiseAuditLogs($page: Int!, $pageSize: Int!, $q: String, $filter: AuditLogFilterType) {\n    auditLogs(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id action resourceType resourceId resultCode actorAccountId storeId createdAt }\n      total current_page per_page total_page\n    }\n  }\n"): (typeof documents)["\n  query FranchiseAuditLogs($page: Int!, $pageSize: Int!, $q: String, $filter: AuditLogFilterType) {\n    auditLogs(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id action resourceType resourceId resultCode actorAccountId storeId createdAt }\n      total current_page per_page total_page\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query FranchiseRoles($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorRoleFilterType) {\n    operatorRoles(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id name kind organizationId permissions { id name action module scope } }\n      total current_page per_page total_page\n    }\n  }\n"): (typeof documents)["\n  query FranchiseRoles($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorRoleFilterType) {\n    operatorRoles(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id name kind organizationId permissions { id name action module scope } }\n      total current_page per_page total_page\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query TenantPermissions {\n    permissions(current_page: 1, per_page: 200, filter: { scope: TENANT }) {\n      data { id name action module scope }\n      total\n    }\n  }\n"): (typeof documents)["\n  query TenantPermissions {\n    permissions(current_page: 1, per_page: 200, filter: { scope: TENANT }) {\n      data { id name action module scope }\n      total\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation FranchiseCreateRole($input: CreateOperatorRoleInput!) {\n    createOperatorRole(input: $input) { id name kind organizationId }\n  }\n"): (typeof documents)["\n  mutation FranchiseCreateRole($input: CreateOperatorRoleInput!) {\n    createOperatorRole(input: $input) { id name kind organizationId }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation FranchiseUpdateRole($id: ID!, $input: UpdateOperatorRoleInput!) {\n    updateOperatorRole(id: $id, input: $input) { id name kind organizationId }\n  }\n"): (typeof documents)["\n  mutation FranchiseUpdateRole($id: ID!, $input: UpdateOperatorRoleInput!) {\n    updateOperatorRole(id: $id, input: $input) { id name kind organizationId }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation FranchiseDeleteRoles($ids: [ID!]!) { deleteOperatorRoles(id: $ids) }\n"): (typeof documents)["\n  mutation FranchiseDeleteRoles($ids: [ID!]!) { deleteOperatorRoles(id: $ids) }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query FranchiseStaff($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorMembershipFilterType) {\n    operatorMemberships(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id status storeAccessMode accountId organizationId\n        account { id phone displayName email }\n        roles { id name kind }\n        stores { id name lifecycle }\n      }\n      total current_page per_page total_page\n    }\n  }\n"): (typeof documents)["\n  query FranchiseStaff($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorMembershipFilterType) {\n    operatorMemberships(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id status storeAccessMode accountId organizationId\n        account { id phone displayName email }\n        roles { id name kind }\n        stores { id name lifecycle }\n      }\n      total current_page per_page total_page\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation FranchiseInviteStaff($input: InviteOperatorInput!) {\n    inviteOperator(input: $input) {\n      membership { id status accountId organizationId }\n      temporaryPassword invitationPending\n    }\n  }\n"): (typeof documents)["\n  mutation FranchiseInviteStaff($input: InviteOperatorInput!) {\n    inviteOperator(input: $input) {\n      membership { id status accountId organizationId }\n      temporaryPassword invitationPending\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation FranchiseUpdateStaff($id: ID!, $input: UpdateOperatorMembershipInput!) {\n    updateOperatorMembership(id: $id, input: $input) { id status storeAccessMode rolesIds storesIds }\n  }\n"): (typeof documents)["\n  mutation FranchiseUpdateStaff($id: ID!, $input: UpdateOperatorMembershipInput!) {\n    updateOperatorMembership(id: $id, input: $input) { id status storeAccessMode rolesIds storesIds }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation FranchiseChangeMembershipStatus($input: ChangeMembershipStatusInput!) {\n    changeMembershipStatus(input: $input) { id status }\n  }\n"): (typeof documents)["\n  mutation FranchiseChangeMembershipStatus($input: ChangeMembershipStatusInput!) {\n    changeMembershipStatus(input: $input) { id status }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "query FranchiseStocktakes($storeId:ID!,$status:StocktakeStatus,$listingId:ID,$batchId:ID,$hasDifference:Boolean,$from:Time,$to:Time,$page:Int!,$perPage:Int!){franchiseStocktakes(storeId:$storeId,status:$status,listingId:$listingId,batchId:$batchId,hasDifference:$hasDifference,from:$from,to:$to,page:$page,perPage:$perPage,includeHistory:false){data{id status startedAt lines{id countedQuantity snapshotQuantity difference needsRecount}} total currentPage perPage}}"): (typeof documents)["query FranchiseStocktakes($storeId:ID!,$status:StocktakeStatus,$listingId:ID,$batchId:ID,$hasDifference:Boolean,$from:Time,$to:Time,$page:Int!,$perPage:Int!){franchiseStocktakes(storeId:$storeId,status:$status,listingId:$listingId,batchId:$batchId,hasDifference:$hasDifference,from:$from,to:$to,page:$page,perPage:$perPage,includeHistory:false){data{id status startedAt lines{id countedQuantity snapshotQuantity difference needsRecount}} total currentPage perPage}}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "query FranchiseStocktake($storeId:ID!,$id:ID!){franchiseStocktake(storeId:$storeId,id:$id){id storeId status startedAt reviewedAt postedAt canceledAt initiatedByAccountId postedById addLineChoices{batchId packageId packageName packageEnabled} lines{id batchId listingId batchNumber expiresAt packageId packageName packageSetVersion packageEnabled countedQuantity snapshotQuantity difference reasonCode reasonNote countHistory{actorAccountId countedQuantity countedAt} needsRecount}}}"): (typeof documents)["query FranchiseStocktake($storeId:ID!,$id:ID!){franchiseStocktake(storeId:$storeId,id:$id){id storeId status startedAt reviewedAt postedAt canceledAt initiatedByAccountId postedById addLineChoices{batchId packageId packageName packageEnabled} lines{id batchId listingId batchNumber expiresAt packageId packageName packageSetVersion packageEnabled countedQuantity snapshotQuantity difference reasonCode reasonNote countHistory{actorAccountId countedQuantity countedAt} needsRecount}}}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "query FranchiseStocktakeBatchChoices($storeId:ID!,$listingId:ID!,$page:Int!,$perPage:Int!){franchiseStocktakeBatchChoices(storeId:$storeId,listingId:$listingId,page:$page,perPage:$perPage){data{id listingId batchNumber expiresAt} total currentPage perPage}}"): (typeof documents)["query FranchiseStocktakeBatchChoices($storeId:ID!,$listingId:ID!,$page:Int!,$perPage:Int!){franchiseStocktakeBatchChoices(storeId:$storeId,listingId:$listingId,page:$page,perPage:$perPage){data{id listingId batchNumber expiresAt} total currentPage perPage}}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "mutation FranchiseCreateStocktake($input:FranchiseCreateStocktakeInput!){franchiseCreateStocktake(input:$input){id status}}"): (typeof documents)["mutation FranchiseCreateStocktake($input:FranchiseCreateStocktakeInput!){franchiseCreateStocktake(input:$input){id status}}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "mutation FranchiseAddStocktakeLine($storeId:ID!,$id:ID!,$batchId:ID!,$packageId:ID!){franchiseAddStocktakeLine(storeId:$storeId,id:$id,batchId:$batchId,packageId:$packageId){id status}}"): (typeof documents)["mutation FranchiseAddStocktakeLine($storeId:ID!,$id:ID!,$batchId:ID!,$packageId:ID!){franchiseAddStocktakeLine(storeId:$storeId,id:$id,batchId:$batchId,packageId:$packageId){id status}}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "mutation FranchiseRecordStocktakeLine($storeId:ID!,$id:ID!,$lineId:ID!,$quantity:Int!){franchiseRecordStocktakeLine(storeId:$storeId,id:$id,lineId:$lineId,quantity:$quantity){id status}}"): (typeof documents)["mutation FranchiseRecordStocktakeLine($storeId:ID!,$id:ID!,$lineId:ID!,$quantity:Int!){franchiseRecordStocktakeLine(storeId:$storeId,id:$id,lineId:$lineId,quantity:$quantity){id status}}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "mutation FranchiseSubmitStocktake($storeId:ID!,$id:ID!){franchiseSubmitStocktake(storeId:$storeId,id:$id){id status}}"): (typeof documents)["mutation FranchiseSubmitStocktake($storeId:ID!,$id:ID!){franchiseSubmitStocktake(storeId:$storeId,id:$id){id status}}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "mutation FranchiseSetStocktakeReason($storeId:ID!,$id:ID!,$lineId:ID!,$reasonCode:String!,$note:String){franchiseSetStocktakeReason(storeId:$storeId,id:$id,lineId:$lineId,reasonCode:$reasonCode,note:$note){id status}}"): (typeof documents)["mutation FranchiseSetStocktakeReason($storeId:ID!,$id:ID!,$lineId:ID!,$reasonCode:String!,$note:String){franchiseSetStocktakeReason(storeId:$storeId,id:$id,lineId:$lineId,reasonCode:$reasonCode,note:$note){id status}}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "mutation FranchiseReturnStocktake($storeId:ID!,$id:ID!){franchiseReturnStocktake(storeId:$storeId,id:$id){id status}}"): (typeof documents)["mutation FranchiseReturnStocktake($storeId:ID!,$id:ID!){franchiseReturnStocktake(storeId:$storeId,id:$id){id status}}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "mutation FranchiseCancelStocktake($storeId:ID!,$id:ID!){franchiseCancelStocktake(storeId:$storeId,id:$id){id status}}"): (typeof documents)["mutation FranchiseCancelStocktake($storeId:ID!,$id:ID!){franchiseCancelStocktake(storeId:$storeId,id:$id){id status}}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "mutation FranchisePostStocktake($storeId:ID!,$id:ID!){franchisePostStocktake(storeId:$storeId,id:$id){id status}}"): (typeof documents)["mutation FranchisePostStocktake($storeId:ID!,$id:ID!){franchisePostStocktake(storeId:$storeId,id:$id){id status}}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query FranchiseStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle rejectionReason organizationId\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      }\n      total current_page per_page total_page\n    }\n  }\n"): (typeof documents)["\n  query FranchiseStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle rejectionReason organizationId\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      }\n      total current_page per_page total_page\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation FranchiseCreateStore($input: CreateStoreInput!) {\n    createStore(input: $input) { id code name lifecycle organizationId }\n  }\n"): (typeof documents)["\n  mutation FranchiseCreateStore($input: CreateStoreInput!) {\n    createStore(input: $input) { id code name lifecycle organizationId }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation FranchiseUpdateStore($id: ID!, $input: UpdateStoreInput!) {\n    updateStore(id: $id, input: $input) { id code name lifecycle organizationId }\n  }\n"): (typeof documents)["\n  mutation FranchiseUpdateStore($id: ID!, $input: UpdateStoreInput!) {\n    updateStore(id: $id, input: $input) { id code name lifecycle organizationId }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation FranchiseSubmitStore($id: ID!) {\n    submitStore(id: $id) { id lifecycle submittedAt }\n  }\n"): (typeof documents)["\n  mutation FranchiseSubmitStore($id: ID!) {\n    submitStore(id: $id) { id lifecycle submittedAt }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation FranchiseDeleteStores($ids: [ID!]!) { deleteStores(id: $ids) }\n"): (typeof documents)["\n  mutation FranchiseDeleteStores($ids: [ID!]!) { deleteStores(id: $ids) }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query HqAdministrators($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorMembershipFilterType) {\n    operatorMemberships(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id status storeAccessMode accountId organizationId updatedAt\n        account { id phone displayName email }\n        roles { id name kind }\n      }\n      total current_page per_page total_page\n    }\n  }\n"): (typeof documents)["\n  query HqAdministrators($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorMembershipFilterType) {\n    operatorMemberships(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id status storeAccessMode accountId organizationId updatedAt\n        account { id phone displayName email }\n        roles { id name kind }\n      }\n      total current_page per_page total_page\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query HqAdministratorIdentity($accountId: ID!) {\n    operatorMemberships(current_page: 1, per_page: 1, filter: { accountId: $accountId }) {\n      data { id accountId roles { id kind } }\n      total\n    }\n  }\n"): (typeof documents)["\n  query HqAdministratorIdentity($accountId: ID!) {\n    operatorMemberships(current_page: 1, per_page: 1, filter: { accountId: $accountId }) {\n      data { id accountId roles { id kind } }\n      total\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqInviteAdministrator($input: InviteOperatorInput!) {\n    inviteOperator(input: $input) {\n      membership { id status accountId organizationId }\n      temporaryPassword invitationPending\n    }\n  }\n"): (typeof documents)["\n  mutation HqInviteAdministrator($input: InviteOperatorInput!) {\n    inviteOperator(input: $input) {\n      membership { id status accountId organizationId }\n      temporaryPassword invitationPending\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqUpdateAdministrator($id: ID!, $input: UpdateOperatorMembershipInput!) {\n    updateOperatorMembership(id: $id, input: $input) { id status rolesIds }\n  }\n"): (typeof documents)["\n  mutation HqUpdateAdministrator($id: ID!, $input: UpdateOperatorMembershipInput!) {\n    updateOperatorMembership(id: $id, input: $input) { id status rolesIds }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqChangeAdministratorStatus($input: ChangeMembershipStatusInput!) {\n    changeMembershipStatus(input: $input) { id status }\n  }\n"): (typeof documents)["\n  mutation HqChangeAdministratorStatus($input: ChangeMembershipStatusInput!) {\n    changeMembershipStatus(input: $input) { id status }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqDeleteAdministrators($ids: [ID!]!) { deleteOperatorMemberships(id: $ids) }\n"): (typeof documents)["\n  mutation HqDeleteAdministrators($ids: [ID!]!) { deleteOperatorMemberships(id: $ids) }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqResetAdministratorPassword($accountId: ID!) {\n    resetTemporaryPassword(accountId: $accountId) { accountId temporaryPassword }\n  }\n"): (typeof documents)["\n  mutation HqResetAdministratorPassword($accountId: ID!) {\n    resetTemporaryPassword(accountId: $accountId) { accountId temporaryPassword }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query HqAuditLogs($page: Int!, $pageSize: Int!, $q: String, $filter: AuditLogFilterType) {\n    auditLogs(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id action resourceType resourceId resultCode actorAccountId organizationId storeId createdAt }\n      total current_page per_page total_page\n    }\n  }\n"): (typeof documents)["\n  query HqAuditLogs($page: Int!, $pageSize: Int!, $q: String, $filter: AuditLogFilterType) {\n    auditLogs(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id action resourceType resourceId resultCode actorAccountId organizationId storeId createdAt }\n      total current_page per_page total_page\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query HqDirectStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle organizationId\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      }\n      total current_page per_page total_page\n    }\n  }\n"): (typeof documents)["\n  query HqDirectStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle organizationId\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      }\n      total current_page per_page total_page\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqCreateDirectStore($input: CreateStoreInput!) {\n    createStore(input: $input) {\n      id code name lifecycle organizationId\n      contactPhone managerName managerPhone province city district address\n      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n    }\n  }\n"): (typeof documents)["\n  mutation HqCreateDirectStore($input: CreateStoreInput!) {\n    createStore(input: $input) {\n      id code name lifecycle organizationId\n      contactPhone managerName managerPhone province city district address\n      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqUpdateDirectStore($id: ID!, $input: UpdateStoreInput!) {\n    updateStore(id: $id, input: $input) {\n      id code name lifecycle organizationId\n      contactPhone managerName managerPhone province city district address\n      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n    }\n  }\n"): (typeof documents)["\n  mutation HqUpdateDirectStore($id: ID!, $input: UpdateStoreInput!) {\n    updateStore(id: $id, input: $input) {\n      id code name lifecycle organizationId\n      contactPhone managerName managerPhone province city district address\n      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqDeleteDirectStores($ids: [ID!]!) { deleteStores(id: $ids) }\n"): (typeof documents)["\n  mutation HqDeleteDirectStores($ids: [ID!]!) { deleteStores(id: $ids) }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query HqFranchiseStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id code name lifecycle businessStatus contactPhone organizationId organization { id name } }\n      total current_page per_page total_page\n    }\n  }\n"): (typeof documents)["\n  query HqFranchiseStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id code name lifecycle businessStatus contactPhone organizationId organization { id name } }\n      total current_page per_page total_page\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query HqFranchises($page: Int!, $pageSize: Int!, $q: String, $filter: OrganizationFilterType!, $canResetPassword: Boolean! = false) {\n    organizations(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id code name status initialAccountId @include(if: $canResetPassword) initialAccount @include(if: $canResetPassword) { phone } }\n      total current_page per_page total_page\n    }\n  }\n"): (typeof documents)["\n  query HqFranchises($page: Int!, $pageSize: Int!, $q: String, $filter: OrganizationFilterType!, $canResetPassword: Boolean! = false) {\n    organizations(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id code name status initialAccountId @include(if: $canResetPassword) initialAccount @include(if: $canResetPassword) { phone } }\n      total current_page per_page total_page\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query HqFranchiseInitialAccountCandidates($id: ID!) {\n    organization(id: $id) {\n      id memberships { id status account { id phone displayName status } }\n    }\n  }\n"): (typeof documents)["\n  query HqFranchiseInitialAccountCandidates($id: ID!) {\n    organization(id: $id) {\n      id memberships { id status account { id phone displayName status } }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqResetFranchiseInitialPassword($organizationId: ID!) {\n    resetFranchiseInitialPassword(organizationId: $organizationId) { accountId temporaryPassword }\n  }\n"): (typeof documents)["\n  mutation HqResetFranchiseInitialPassword($organizationId: ID!) {\n    resetFranchiseInitialPassword(organizationId: $organizationId) { accountId temporaryPassword }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqProvisionFranchise($input: ProvisionFranchiseInput!) {\n    provisionFranchise(input: $input) {\n      organization { id code name status }\n      membership { id status }\n      temporaryPassword invitationPending\n    }\n  }\n"): (typeof documents)["\n  mutation HqProvisionFranchise($input: ProvisionFranchiseInput!) {\n    provisionFranchise(input: $input) {\n      organization { id code name status }\n      membership { id status }\n      temporaryPassword invitationPending\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqSuspendOrganization($input: SuspendOrganizationInput!) {\n    suspendOrganization(input: $input) { id status suspensionReasonCode }\n  }\n"): (typeof documents)["\n  mutation HqSuspendOrganization($input: SuspendOrganizationInput!) {\n    suspendOrganization(input: $input) { id status suspensionReasonCode }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqRestoreOrganization($id: ID!) {\n    restoreOrganization(id: $id) { id status }\n  }\n"): (typeof documents)["\n  mutation HqRestoreOrganization($id: ID!) {\n    restoreOrganization(id: $id) { id status }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query HqStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle organizationId\n        organization { id name type }\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      }\n      total current_page per_page total_page\n    }\n  }\n"): (typeof documents)["\n  query HqStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle organizationId\n        organization { id name type }\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter\n      }\n      total current_page per_page total_page\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query HqStoresTotal { stores(current_page: 1, per_page: 1) { total } }\n"): (typeof documents)["\n  query HqStoresTotal { stores(current_page: 1, per_page: 1) { total } }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query HqRoles($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorRoleFilterType!) {\n    operatorRoles(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id name kind organizationId permissions { id name action module scope } }\n      total current_page per_page total_page\n    }\n  }\n"): (typeof documents)["\n  query HqRoles($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorRoleFilterType!) {\n    operatorRoles(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data { id name kind organizationId permissions { id name action module scope } }\n      total current_page per_page total_page\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query SystemPermissions {\n    permissions(current_page: 1, per_page: 200, filter: { scope: SYSTEM }) {\n      data { id name action module scope }\n      total\n    }\n  }\n"): (typeof documents)["\n  query SystemPermissions {\n    permissions(current_page: 1, per_page: 200, filter: { scope: SYSTEM }) {\n      data { id name action module scope }\n      total\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqCreateRole($input: CreateOperatorRoleInput!) {\n    createOperatorRole(input: $input) { id name kind organizationId }\n  }\n"): (typeof documents)["\n  mutation HqCreateRole($input: CreateOperatorRoleInput!) {\n    createOperatorRole(input: $input) { id name kind organizationId }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqUpdateRole($id: ID!, $input: UpdateOperatorRoleInput!) {\n    updateOperatorRole(id: $id, input: $input) { id name kind organizationId }\n  }\n"): (typeof documents)["\n  mutation HqUpdateRole($id: ID!, $input: UpdateOperatorRoleInput!) {\n    updateOperatorRole(id: $id, input: $input) { id name kind organizationId }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqDeleteRoles($ids: [ID!]!) { deleteOperatorRoles(id: $ids) }\n"): (typeof documents)["\n  mutation HqDeleteRoles($ids: [ID!]!) { deleteOperatorRoles(id: $ids) }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query HqStoreApprovals($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle submittedAt organizationId organization { id name }\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount\n      }\n      total current_page per_page total_page\n    }\n  }\n"): (typeof documents)["\n  query HqStoreApprovals($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType!) {\n    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {\n      data {\n        id code name lifecycle submittedAt organizationId organization { id name }\n        contactPhone managerName managerPhone province city district address\n        businessHours businessStatus supportDineIn supportTakeout storeArea tableCount\n      }\n      total current_page per_page total_page\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation HqReviewStore($input: ReviewStoreInput!) {\n    reviewStore(input: $input) { id lifecycle rejectionReason reviewedAt }\n  }\n"): (typeof documents)["\n  mutation HqReviewStore($input: ReviewStoreInput!) {\n    reviewStore(input: $input) { id lifecycle rejectionReason reviewedAt }\n  }\n"];

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;