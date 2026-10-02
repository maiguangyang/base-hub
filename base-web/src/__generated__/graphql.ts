/* eslint-disable */
/** Internal type. DO NOT USE DIRECTLY. */
type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
/** Internal type. DO NOT USE DIRECTLY. */
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type AccountFilterType = {
  AND?: Array<AccountFilterType> | null | undefined;
  OR?: Array<AccountFilterType> | null | undefined;
  auditLogs?: AuditLogFilterType | null | undefined;
  createdAt?: number | null | undefined;
  createdAt_gt?: number | null | undefined;
  createdAt_gte?: number | null | undefined;
  createdAt_in?: Array<number> | null | undefined;
  createdAt_lt?: number | null | undefined;
  createdAt_lte?: number | null | undefined;
  createdAt_ne?: number | null | undefined;
  createdAt_null?: boolean | null | undefined;
  createdBy?: string | number | null | undefined;
  createdBy_gt?: string | number | null | undefined;
  createdBy_gte?: string | number | null | undefined;
  createdBy_in?: Array<string | number> | null | undefined;
  createdBy_lt?: string | number | null | undefined;
  createdBy_lte?: string | number | null | undefined;
  createdBy_ne?: string | number | null | undefined;
  createdBy_null?: boolean | null | undefined;
  credentialVersion?: number | null | undefined;
  credentialVersion_gt?: number | null | undefined;
  credentialVersion_gte?: number | null | undefined;
  credentialVersion_in?: Array<number> | null | undefined;
  credentialVersion_lt?: number | null | undefined;
  credentialVersion_lte?: number | null | undefined;
  credentialVersion_ne?: number | null | undefined;
  credentialVersion_null?: boolean | null | undefined;
  deletedAt?: number | null | undefined;
  deletedAt_gt?: number | null | undefined;
  deletedAt_gte?: number | null | undefined;
  deletedAt_in?: Array<number> | null | undefined;
  deletedAt_lt?: number | null | undefined;
  deletedAt_lte?: number | null | undefined;
  deletedAt_ne?: number | null | undefined;
  deletedAt_null?: boolean | null | undefined;
  deletedBy?: string | number | null | undefined;
  deletedBy_gt?: string | number | null | undefined;
  deletedBy_gte?: string | number | null | undefined;
  deletedBy_in?: Array<string | number> | null | undefined;
  deletedBy_lt?: string | number | null | undefined;
  deletedBy_lte?: string | number | null | undefined;
  deletedBy_ne?: string | number | null | undefined;
  deletedBy_null?: boolean | null | undefined;
  displayName?: string | null | undefined;
  displayName_gt?: string | null | undefined;
  displayName_gte?: string | null | undefined;
  displayName_in?: Array<string> | null | undefined;
  displayName_like?: string | null | undefined;
  displayName_lt?: string | null | undefined;
  displayName_lte?: string | null | undefined;
  displayName_ne?: string | null | undefined;
  displayName_null?: boolean | null | undefined;
  displayName_prefix?: string | null | undefined;
  displayName_suffix?: string | null | undefined;
  email?: string | null | undefined;
  email_gt?: string | null | undefined;
  email_gte?: string | null | undefined;
  email_in?: Array<string> | null | undefined;
  email_like?: string | null | undefined;
  email_lt?: string | null | undefined;
  email_lte?: string | null | undefined;
  email_ne?: string | null | undefined;
  email_null?: boolean | null | undefined;
  email_prefix?: string | null | undefined;
  email_suffix?: string | null | undefined;
  id?: string | number | null | undefined;
  id_gt?: string | number | null | undefined;
  id_gte?: string | number | null | undefined;
  id_in?: Array<string | number> | null | undefined;
  id_lt?: string | number | null | undefined;
  id_lte?: string | number | null | undefined;
  id_ne?: string | number | null | undefined;
  id_null?: boolean | null | undefined;
  initializedOrganizations?: OrganizationFilterType | null | undefined;
  isDelete?: number | null | undefined;
  isDelete_gt?: number | null | undefined;
  isDelete_gte?: number | null | undefined;
  isDelete_in?: Array<number> | null | undefined;
  isDelete_lt?: number | null | undefined;
  isDelete_lte?: number | null | undefined;
  isDelete_ne?: number | null | undefined;
  isDelete_null?: boolean | null | undefined;
  memberships?: OperatorMembershipFilterType | null | undefined;
  mustChangePassword?: boolean | null | undefined;
  mustChangePassword_gt?: boolean | null | undefined;
  mustChangePassword_gte?: boolean | null | undefined;
  mustChangePassword_in?: Array<boolean> | null | undefined;
  mustChangePassword_lt?: boolean | null | undefined;
  mustChangePassword_lte?: boolean | null | undefined;
  mustChangePassword_ne?: boolean | null | undefined;
  mustChangePassword_null?: boolean | null | undefined;
  openingRecords?: FranchiseOpeningRecordFilterType | null | undefined;
  phone?: string | null | undefined;
  phone_gt?: string | null | undefined;
  phone_gte?: string | null | undefined;
  phone_in?: Array<string> | null | undefined;
  phone_like?: string | null | undefined;
  phone_lt?: string | null | undefined;
  phone_lte?: string | null | undefined;
  phone_ne?: string | null | undefined;
  phone_null?: boolean | null | undefined;
  phone_prefix?: string | null | undefined;
  phone_suffix?: string | null | undefined;
  recordedOpeningRecords?: FranchiseOpeningRecordFilterType | null | undefined;
  reviewedStores?: StoreFilterType | null | undefined;
  sentMembershipInvitations?: MembershipInvitationFilterType | null | undefined;
  sessions?: SessionFilterType | null | undefined;
  state?: number | null | undefined;
  state_gt?: number | null | undefined;
  state_gte?: number | null | undefined;
  state_in?: Array<number> | null | undefined;
  state_lt?: number | null | undefined;
  state_lte?: number | null | undefined;
  state_ne?: number | null | undefined;
  state_null?: boolean | null | undefined;
  status?: AccountStatus | null | undefined;
  status_gt?: AccountStatus | null | undefined;
  status_gte?: AccountStatus | null | undefined;
  status_in?: Array<AccountStatus> | null | undefined;
  status_lt?: AccountStatus | null | undefined;
  status_lte?: AccountStatus | null | undefined;
  status_ne?: AccountStatus | null | undefined;
  status_null?: boolean | null | undefined;
  updatedAt?: number | null | undefined;
  updatedAt_gt?: number | null | undefined;
  updatedAt_gte?: number | null | undefined;
  updatedAt_in?: Array<number> | null | undefined;
  updatedAt_lt?: number | null | undefined;
  updatedAt_lte?: number | null | undefined;
  updatedAt_ne?: number | null | undefined;
  updatedAt_null?: boolean | null | undefined;
  updatedBy?: string | number | null | undefined;
  updatedBy_gt?: string | number | null | undefined;
  updatedBy_gte?: string | number | null | undefined;
  updatedBy_in?: Array<string | number> | null | undefined;
  updatedBy_lt?: string | number | null | undefined;
  updatedBy_lte?: string | number | null | undefined;
  updatedBy_ne?: string | number | null | undefined;
  updatedBy_null?: boolean | null | undefined;
  weight?: number | null | undefined;
  weight_gt?: number | null | undefined;
  weight_gte?: number | null | undefined;
  weight_in?: Array<number> | null | undefined;
  weight_lt?: number | null | undefined;
  weight_lte?: number | null | undefined;
  weight_ne?: number | null | undefined;
  weight_null?: boolean | null | undefined;
};

export type AccountRelationship = {
  credentialVersion?: number | null | undefined;
  displayName?: string | null | undefined;
  email?: string | null | undefined;
  id?: string | number | null | undefined;
  isDelete?: number | null | undefined;
  mustChangePassword?: boolean | null | undefined;
  phone?: string | null | undefined;
  state?: number | null | undefined;
  status?: AccountStatus | null | undefined;
  weight?: number | null | undefined;
};

export type AccountStatus =
  | 'ACTIVE'
  | 'DISABLED';

export type AuditLogFilterType = {
  AND?: Array<AuditLogFilterType> | null | undefined;
  OR?: Array<AuditLogFilterType> | null | undefined;
  action?: string | null | undefined;
  action_gt?: string | null | undefined;
  action_gte?: string | null | undefined;
  action_in?: Array<string> | null | undefined;
  action_like?: string | null | undefined;
  action_lt?: string | null | undefined;
  action_lte?: string | null | undefined;
  action_ne?: string | null | undefined;
  action_null?: boolean | null | undefined;
  action_prefix?: string | null | undefined;
  action_suffix?: string | null | undefined;
  actorAccount?: AccountFilterType | null | undefined;
  actorAccountId?: string | number | null | undefined;
  actorAccountId_gt?: string | number | null | undefined;
  actorAccountId_gte?: string | number | null | undefined;
  actorAccountId_in?: Array<string | number> | null | undefined;
  actorAccountId_lt?: string | number | null | undefined;
  actorAccountId_lte?: string | number | null | undefined;
  actorAccountId_ne?: string | number | null | undefined;
  actorAccountId_null?: boolean | null | undefined;
  createdAt?: number | null | undefined;
  createdAt_gt?: number | null | undefined;
  createdAt_gte?: number | null | undefined;
  createdAt_in?: Array<number> | null | undefined;
  createdAt_lt?: number | null | undefined;
  createdAt_lte?: number | null | undefined;
  createdAt_ne?: number | null | undefined;
  createdAt_null?: boolean | null | undefined;
  createdBy?: string | number | null | undefined;
  createdBy_gt?: string | number | null | undefined;
  createdBy_gte?: string | number | null | undefined;
  createdBy_in?: Array<string | number> | null | undefined;
  createdBy_lt?: string | number | null | undefined;
  createdBy_lte?: string | number | null | undefined;
  createdBy_ne?: string | number | null | undefined;
  createdBy_null?: boolean | null | undefined;
  deletedAt?: number | null | undefined;
  deletedAt_gt?: number | null | undefined;
  deletedAt_gte?: number | null | undefined;
  deletedAt_in?: Array<number> | null | undefined;
  deletedAt_lt?: number | null | undefined;
  deletedAt_lte?: number | null | undefined;
  deletedAt_ne?: number | null | undefined;
  deletedAt_null?: boolean | null | undefined;
  deletedBy?: string | number | null | undefined;
  deletedBy_gt?: string | number | null | undefined;
  deletedBy_gte?: string | number | null | undefined;
  deletedBy_in?: Array<string | number> | null | undefined;
  deletedBy_lt?: string | number | null | undefined;
  deletedBy_lte?: string | number | null | undefined;
  deletedBy_ne?: string | number | null | undefined;
  deletedBy_null?: boolean | null | undefined;
  id?: string | number | null | undefined;
  id_gt?: string | number | null | undefined;
  id_gte?: string | number | null | undefined;
  id_in?: Array<string | number> | null | undefined;
  id_lt?: string | number | null | undefined;
  id_lte?: string | number | null | undefined;
  id_ne?: string | number | null | undefined;
  id_null?: boolean | null | undefined;
  isDelete?: number | null | undefined;
  isDelete_gt?: number | null | undefined;
  isDelete_gte?: number | null | undefined;
  isDelete_in?: Array<number> | null | undefined;
  isDelete_lt?: number | null | undefined;
  isDelete_lte?: number | null | undefined;
  isDelete_ne?: number | null | undefined;
  isDelete_null?: boolean | null | undefined;
  metadataJson?: string | null | undefined;
  metadataJson_gt?: string | null | undefined;
  metadataJson_gte?: string | null | undefined;
  metadataJson_in?: Array<string> | null | undefined;
  metadataJson_like?: string | null | undefined;
  metadataJson_lt?: string | null | undefined;
  metadataJson_lte?: string | null | undefined;
  metadataJson_ne?: string | null | undefined;
  metadataJson_null?: boolean | null | undefined;
  metadataJson_prefix?: string | null | undefined;
  metadataJson_suffix?: string | null | undefined;
  organization?: OrganizationFilterType | null | undefined;
  organizationId?: string | number | null | undefined;
  organizationId_gt?: string | number | null | undefined;
  organizationId_gte?: string | number | null | undefined;
  organizationId_in?: Array<string | number> | null | undefined;
  organizationId_lt?: string | number | null | undefined;
  organizationId_lte?: string | number | null | undefined;
  organizationId_ne?: string | number | null | undefined;
  organizationId_null?: boolean | null | undefined;
  resourceId?: string | null | undefined;
  resourceId_gt?: string | null | undefined;
  resourceId_gte?: string | null | undefined;
  resourceId_in?: Array<string> | null | undefined;
  resourceId_like?: string | null | undefined;
  resourceId_lt?: string | null | undefined;
  resourceId_lte?: string | null | undefined;
  resourceId_ne?: string | null | undefined;
  resourceId_null?: boolean | null | undefined;
  resourceId_prefix?: string | null | undefined;
  resourceId_suffix?: string | null | undefined;
  resourceType?: string | null | undefined;
  resourceType_gt?: string | null | undefined;
  resourceType_gte?: string | null | undefined;
  resourceType_in?: Array<string> | null | undefined;
  resourceType_like?: string | null | undefined;
  resourceType_lt?: string | null | undefined;
  resourceType_lte?: string | null | undefined;
  resourceType_ne?: string | null | undefined;
  resourceType_null?: boolean | null | undefined;
  resourceType_prefix?: string | null | undefined;
  resourceType_suffix?: string | null | undefined;
  resultCode?: string | null | undefined;
  resultCode_gt?: string | null | undefined;
  resultCode_gte?: string | null | undefined;
  resultCode_in?: Array<string> | null | undefined;
  resultCode_like?: string | null | undefined;
  resultCode_lt?: string | null | undefined;
  resultCode_lte?: string | null | undefined;
  resultCode_ne?: string | null | undefined;
  resultCode_null?: boolean | null | undefined;
  resultCode_prefix?: string | null | undefined;
  resultCode_suffix?: string | null | undefined;
  session?: SessionFilterType | null | undefined;
  sessionId?: string | number | null | undefined;
  sessionId_gt?: string | number | null | undefined;
  sessionId_gte?: string | number | null | undefined;
  sessionId_in?: Array<string | number> | null | undefined;
  sessionId_lt?: string | number | null | undefined;
  sessionId_lte?: string | number | null | undefined;
  sessionId_ne?: string | number | null | undefined;
  sessionId_null?: boolean | null | undefined;
  state?: number | null | undefined;
  state_gt?: number | null | undefined;
  state_gte?: number | null | undefined;
  state_in?: Array<number> | null | undefined;
  state_lt?: number | null | undefined;
  state_lte?: number | null | undefined;
  state_ne?: number | null | undefined;
  state_null?: boolean | null | undefined;
  store?: StoreFilterType | null | undefined;
  storeId?: string | number | null | undefined;
  storeId_gt?: string | number | null | undefined;
  storeId_gte?: string | number | null | undefined;
  storeId_in?: Array<string | number> | null | undefined;
  storeId_lt?: string | number | null | undefined;
  storeId_lte?: string | number | null | undefined;
  storeId_ne?: string | number | null | undefined;
  storeId_null?: boolean | null | undefined;
  updatedAt?: number | null | undefined;
  updatedAt_gt?: number | null | undefined;
  updatedAt_gte?: number | null | undefined;
  updatedAt_in?: Array<number> | null | undefined;
  updatedAt_lt?: number | null | undefined;
  updatedAt_lte?: number | null | undefined;
  updatedAt_ne?: number | null | undefined;
  updatedAt_null?: boolean | null | undefined;
  updatedBy?: string | number | null | undefined;
  updatedBy_gt?: string | number | null | undefined;
  updatedBy_gte?: string | number | null | undefined;
  updatedBy_in?: Array<string | number> | null | undefined;
  updatedBy_lt?: string | number | null | undefined;
  updatedBy_lte?: string | number | null | undefined;
  updatedBy_ne?: string | number | null | undefined;
  updatedBy_null?: boolean | null | undefined;
  weight?: number | null | undefined;
  weight_gt?: number | null | undefined;
  weight_gte?: number | null | undefined;
  weight_in?: Array<number> | null | undefined;
  weight_lt?: number | null | undefined;
  weight_lte?: number | null | undefined;
  weight_ne?: number | null | undefined;
  weight_null?: boolean | null | undefined;
};

export type AuditLogRelationship = {
  action?: string | null | undefined;
  actorAccountId?: string | number | null | undefined;
  id?: string | number | null | undefined;
  isDelete?: number | null | undefined;
  metadataJson?: string | null | undefined;
  organizationId?: string | number | null | undefined;
  resourceId?: string | null | undefined;
  resourceType?: string | null | undefined;
  resultCode?: string | null | undefined;
  sessionId?: string | number | null | undefined;
  state?: number | null | undefined;
  storeId?: string | number | null | undefined;
  weight?: number | null | undefined;
};

export type ChangeMembershipStatusInput = {
  membershipId: string | number;
  status: MembershipStatus;
};

export type ChangePasswordInput = {
  currentPassword: string;
  newPassword: string;
};

export type CreateOperatorRoleInput = {
  isDelete?: number | null | undefined;
  kind: RoleKind;
  members?: Array<OperatorMembershipRelationship | null | undefined> | null | undefined;
  membersIds?: Array<string | number> | null | undefined;
  name: string;
  organization?: OrganizationRelationship | null | undefined;
  organizationId: string | number;
  permissions?: Array<PermissionRelationship | null | undefined> | null | undefined;
  permissionsIds?: Array<string | number> | null | undefined;
  state?: number | null | undefined;
  weight?: number | null | undefined;
};

export type CreateStoreInput = {
  address?: string | null | undefined;
  auditLogs?: Array<AuditLogRelationship | null | undefined> | null | undefined;
  auditLogsIds?: Array<string | number> | null | undefined;
  businessHours?: string | null | undefined;
  businessLicenseImageUrl?: string | null | undefined;
  businessStatus?: StoreBusinessStatus | null | undefined;
  city?: string | null | undefined;
  code: string;
  contactPhone?: string | null | undefined;
  district?: string | null | undefined;
  isDelete?: number | null | undefined;
  lifecycle: StoreLifecycle;
  managerName?: string | null | undefined;
  managerPhone?: string | null | undefined;
  members?: Array<OperatorMembershipRelationship | null | undefined> | null | undefined;
  membersIds?: Array<string | number> | null | undefined;
  name: string;
  organization?: OrganizationRelationship | null | undefined;
  organizationId: string | number;
  otherDocumentImageUrl?: string | null | undefined;
  paymentConfigs?: Array<StorePaymentConfigRelationship | null | undefined> | null | undefined;
  paymentConfigsIds?: Array<string | number> | null | undefined;
  province?: string | null | undefined;
  receiptFooter?: string | null | undefined;
  rejectionReason?: string | null | undefined;
  reviewedAt?: unknown;
  reviewedByAccount?: AccountRelationship | null | undefined;
  reviewedByAccountId?: string | number | null | undefined;
  state?: number | null | undefined;
  storeArea?: number | null | undefined;
  submittedAt?: unknown;
  supportDineIn?: boolean | null | undefined;
  supportTakeout?: boolean | null | undefined;
  tableCount?: number | null | undefined;
  weight?: number | null | undefined;
};

export type FranchiseOpeningRecordFilterType = {
  AND?: Array<FranchiseOpeningRecordFilterType> | null | undefined;
  OR?: Array<FranchiseOpeningRecordFilterType> | null | undefined;
  createdAt?: number | null | undefined;
  createdAt_gt?: number | null | undefined;
  createdAt_gte?: number | null | undefined;
  createdAt_in?: Array<number> | null | undefined;
  createdAt_lt?: number | null | undefined;
  createdAt_lte?: number | null | undefined;
  createdAt_ne?: number | null | undefined;
  createdAt_null?: boolean | null | undefined;
  createdBy?: string | number | null | undefined;
  createdBy_gt?: string | number | null | undefined;
  createdBy_gte?: string | number | null | undefined;
  createdBy_in?: Array<string | number> | null | undefined;
  createdBy_lt?: string | number | null | undefined;
  createdBy_lte?: string | number | null | undefined;
  createdBy_ne?: string | number | null | undefined;
  createdBy_null?: boolean | null | undefined;
  deletedAt?: number | null | undefined;
  deletedAt_gt?: number | null | undefined;
  deletedAt_gte?: number | null | undefined;
  deletedAt_in?: Array<number> | null | undefined;
  deletedAt_lt?: number | null | undefined;
  deletedAt_lte?: number | null | undefined;
  deletedAt_ne?: number | null | undefined;
  deletedAt_null?: boolean | null | undefined;
  deletedBy?: string | number | null | undefined;
  deletedBy_gt?: string | number | null | undefined;
  deletedBy_gte?: string | number | null | undefined;
  deletedBy_in?: Array<string | number> | null | undefined;
  deletedBy_lt?: string | number | null | undefined;
  deletedBy_lte?: string | number | null | undefined;
  deletedBy_ne?: string | number | null | undefined;
  deletedBy_null?: boolean | null | undefined;
  id?: string | number | null | undefined;
  id_gt?: string | number | null | undefined;
  id_gte?: string | number | null | undefined;
  id_in?: Array<string | number> | null | undefined;
  id_lt?: string | number | null | undefined;
  id_lte?: string | number | null | undefined;
  id_ne?: string | number | null | undefined;
  id_null?: boolean | null | undefined;
  initialAccount?: AccountFilterType | null | undefined;
  initialAccountId?: string | number | null | undefined;
  initialAccountId_gt?: string | number | null | undefined;
  initialAccountId_gte?: string | number | null | undefined;
  initialAccountId_in?: Array<string | number> | null | undefined;
  initialAccountId_lt?: string | number | null | undefined;
  initialAccountId_lte?: string | number | null | undefined;
  initialAccountId_ne?: string | number | null | undefined;
  initialAccountId_null?: boolean | null | undefined;
  isDelete?: number | null | undefined;
  isDelete_gt?: number | null | undefined;
  isDelete_gte?: number | null | undefined;
  isDelete_in?: Array<number> | null | undefined;
  isDelete_lt?: number | null | undefined;
  isDelete_lte?: number | null | undefined;
  isDelete_ne?: number | null | undefined;
  isDelete_null?: boolean | null | undefined;
  organization?: OrganizationFilterType | null | undefined;
  organizationId?: string | number | null | undefined;
  organizationId_gt?: string | number | null | undefined;
  organizationId_gte?: string | number | null | undefined;
  organizationId_in?: Array<string | number> | null | undefined;
  organizationId_lt?: string | number | null | undefined;
  organizationId_lte?: string | number | null | undefined;
  organizationId_ne?: string | number | null | undefined;
  organizationId_null?: boolean | null | undefined;
  recordNumber?: string | null | undefined;
  recordNumber_gt?: string | null | undefined;
  recordNumber_gte?: string | null | undefined;
  recordNumber_in?: Array<string> | null | undefined;
  recordNumber_like?: string | null | undefined;
  recordNumber_lt?: string | null | undefined;
  recordNumber_lte?: string | null | undefined;
  recordNumber_ne?: string | null | undefined;
  recordNumber_null?: boolean | null | undefined;
  recordNumber_prefix?: string | null | undefined;
  recordNumber_suffix?: string | null | undefined;
  recordedByAccount?: AccountFilterType | null | undefined;
  recordedByAccountId?: string | number | null | undefined;
  recordedByAccountId_gt?: string | number | null | undefined;
  recordedByAccountId_gte?: string | number | null | undefined;
  recordedByAccountId_in?: Array<string | number> | null | undefined;
  recordedByAccountId_lt?: string | number | null | undefined;
  recordedByAccountId_lte?: string | number | null | undefined;
  recordedByAccountId_ne?: string | number | null | undefined;
  recordedByAccountId_null?: boolean | null | undefined;
  source?: FranchiseOpeningSource | null | undefined;
  source_gt?: FranchiseOpeningSource | null | undefined;
  source_gte?: FranchiseOpeningSource | null | undefined;
  source_in?: Array<FranchiseOpeningSource> | null | undefined;
  source_lt?: FranchiseOpeningSource | null | undefined;
  source_lte?: FranchiseOpeningSource | null | undefined;
  source_ne?: FranchiseOpeningSource | null | undefined;
  source_null?: boolean | null | undefined;
  state?: number | null | undefined;
  state_gt?: number | null | undefined;
  state_gte?: number | null | undefined;
  state_in?: Array<number> | null | undefined;
  state_lt?: number | null | undefined;
  state_lte?: number | null | undefined;
  state_ne?: number | null | undefined;
  state_null?: boolean | null | undefined;
  updatedAt?: number | null | undefined;
  updatedAt_gt?: number | null | undefined;
  updatedAt_gte?: number | null | undefined;
  updatedAt_in?: Array<number> | null | undefined;
  updatedAt_lt?: number | null | undefined;
  updatedAt_lte?: number | null | undefined;
  updatedAt_ne?: number | null | undefined;
  updatedAt_null?: boolean | null | undefined;
  updatedBy?: string | number | null | undefined;
  updatedBy_gt?: string | number | null | undefined;
  updatedBy_gte?: string | number | null | undefined;
  updatedBy_in?: Array<string | number> | null | undefined;
  updatedBy_lt?: string | number | null | undefined;
  updatedBy_lte?: string | number | null | undefined;
  updatedBy_ne?: string | number | null | undefined;
  updatedBy_null?: boolean | null | undefined;
  weight?: number | null | undefined;
  weight_gt?: number | null | undefined;
  weight_gte?: number | null | undefined;
  weight_in?: Array<number> | null | undefined;
  weight_lt?: number | null | undefined;
  weight_lte?: number | null | undefined;
  weight_ne?: number | null | undefined;
  weight_null?: boolean | null | undefined;
};

export type FranchiseOpeningSource =
  | 'HISTORICAL_ATTESTATION'
  | 'SYSTEM_PROVISION';

export type FranchisePaymentConfigFilterType = {
  AND?: Array<FranchisePaymentConfigFilterType> | null | undefined;
  OR?: Array<FranchisePaymentConfigFilterType> | null | undefined;
  channel?: string | null | undefined;
  channel_gt?: string | null | undefined;
  channel_gte?: string | null | undefined;
  channel_in?: Array<string> | null | undefined;
  channel_like?: string | null | undefined;
  channel_lt?: string | null | undefined;
  channel_lte?: string | null | undefined;
  channel_ne?: string | null | undefined;
  channel_null?: boolean | null | undefined;
  channel_prefix?: string | null | undefined;
  channel_suffix?: string | null | undefined;
  configState?: string | null | undefined;
  configState_gt?: string | null | undefined;
  configState_gte?: string | null | undefined;
  configState_in?: Array<string> | null | undefined;
  configState_like?: string | null | undefined;
  configState_lt?: string | null | undefined;
  configState_lte?: string | null | undefined;
  configState_ne?: string | null | undefined;
  configState_null?: boolean | null | undefined;
  configState_prefix?: string | null | undefined;
  configState_suffix?: string | null | undefined;
  createdAt?: number | null | undefined;
  createdAt_gt?: number | null | undefined;
  createdAt_gte?: number | null | undefined;
  createdAt_in?: Array<number> | null | undefined;
  createdAt_lt?: number | null | undefined;
  createdAt_lte?: number | null | undefined;
  createdAt_ne?: number | null | undefined;
  createdAt_null?: boolean | null | undefined;
  createdBy?: string | number | null | undefined;
  createdBy_gt?: string | number | null | undefined;
  createdBy_gte?: string | number | null | undefined;
  createdBy_in?: Array<string | number> | null | undefined;
  createdBy_lt?: string | number | null | undefined;
  createdBy_lte?: string | number | null | undefined;
  createdBy_ne?: string | number | null | undefined;
  createdBy_null?: boolean | null | undefined;
  credentialCiphertext?: string | null | undefined;
  credentialCiphertext_gt?: string | null | undefined;
  credentialCiphertext_gte?: string | null | undefined;
  credentialCiphertext_in?: Array<string> | null | undefined;
  credentialCiphertext_like?: string | null | undefined;
  credentialCiphertext_lt?: string | null | undefined;
  credentialCiphertext_lte?: string | null | undefined;
  credentialCiphertext_ne?: string | null | undefined;
  credentialCiphertext_null?: boolean | null | undefined;
  credentialCiphertext_prefix?: string | null | undefined;
  credentialCiphertext_suffix?: string | null | undefined;
  deletedAt?: number | null | undefined;
  deletedAt_gt?: number | null | undefined;
  deletedAt_gte?: number | null | undefined;
  deletedAt_in?: Array<number> | null | undefined;
  deletedAt_lt?: number | null | undefined;
  deletedAt_lte?: number | null | undefined;
  deletedAt_ne?: number | null | undefined;
  deletedAt_null?: boolean | null | undefined;
  deletedBy?: string | number | null | undefined;
  deletedBy_gt?: string | number | null | undefined;
  deletedBy_gte?: string | number | null | undefined;
  deletedBy_in?: Array<string | number> | null | undefined;
  deletedBy_lt?: string | number | null | undefined;
  deletedBy_lte?: string | number | null | undefined;
  deletedBy_ne?: string | number | null | undefined;
  deletedBy_null?: boolean | null | undefined;
  environment?: string | null | undefined;
  environment_gt?: string | null | undefined;
  environment_gte?: string | null | undefined;
  environment_in?: Array<string> | null | undefined;
  environment_like?: string | null | undefined;
  environment_lt?: string | null | undefined;
  environment_lte?: string | null | undefined;
  environment_ne?: string | null | undefined;
  environment_null?: boolean | null | undefined;
  environment_prefix?: string | null | undefined;
  environment_suffix?: string | null | undefined;
  id?: string | number | null | undefined;
  id_gt?: string | number | null | undefined;
  id_gte?: string | number | null | undefined;
  id_in?: Array<string | number> | null | undefined;
  id_lt?: string | number | null | undefined;
  id_lte?: string | number | null | undefined;
  id_ne?: string | number | null | undefined;
  id_null?: boolean | null | undefined;
  isDelete?: number | null | undefined;
  isDelete_gt?: number | null | undefined;
  isDelete_gte?: number | null | undefined;
  isDelete_in?: Array<number> | null | undefined;
  isDelete_lt?: number | null | undefined;
  isDelete_lte?: number | null | undefined;
  isDelete_ne?: number | null | undefined;
  isDelete_null?: boolean | null | undefined;
  keyId?: string | null | undefined;
  keyId_gt?: string | null | undefined;
  keyId_gte?: string | null | undefined;
  keyId_in?: Array<string> | null | undefined;
  keyId_like?: string | null | undefined;
  keyId_lt?: string | null | undefined;
  keyId_lte?: string | null | undefined;
  keyId_ne?: string | null | undefined;
  keyId_null?: boolean | null | undefined;
  keyId_prefix?: string | null | undefined;
  keyId_suffix?: string | null | undefined;
  merchantId?: string | null | undefined;
  merchantId_gt?: string | null | undefined;
  merchantId_gte?: string | null | undefined;
  merchantId_in?: Array<string> | null | undefined;
  merchantId_like?: string | null | undefined;
  merchantId_lt?: string | null | undefined;
  merchantId_lte?: string | null | undefined;
  merchantId_ne?: string | null | undefined;
  merchantId_null?: boolean | null | undefined;
  merchantId_prefix?: string | null | undefined;
  merchantId_suffix?: string | null | undefined;
  organization?: OrganizationFilterType | null | undefined;
  organizationId?: string | number | null | undefined;
  organizationId_gt?: string | number | null | undefined;
  organizationId_gte?: string | number | null | undefined;
  organizationId_in?: Array<string | number> | null | undefined;
  organizationId_lt?: string | number | null | undefined;
  organizationId_lte?: string | number | null | undefined;
  organizationId_ne?: string | number | null | undefined;
  organizationId_null?: boolean | null | undefined;
  ratePpm?: number | null | undefined;
  ratePpm_gt?: number | null | undefined;
  ratePpm_gte?: number | null | undefined;
  ratePpm_in?: Array<number> | null | undefined;
  ratePpm_lt?: number | null | undefined;
  ratePpm_lte?: number | null | undefined;
  ratePpm_ne?: number | null | undefined;
  ratePpm_null?: boolean | null | undefined;
  state?: number | null | undefined;
  state_gt?: number | null | undefined;
  state_gte?: number | null | undefined;
  state_in?: Array<number> | null | undefined;
  state_lt?: number | null | undefined;
  state_lte?: number | null | undefined;
  state_ne?: number | null | undefined;
  state_null?: boolean | null | undefined;
  updatedAt?: number | null | undefined;
  updatedAt_gt?: number | null | undefined;
  updatedAt_gte?: number | null | undefined;
  updatedAt_in?: Array<number> | null | undefined;
  updatedAt_lt?: number | null | undefined;
  updatedAt_lte?: number | null | undefined;
  updatedAt_ne?: number | null | undefined;
  updatedAt_null?: boolean | null | undefined;
  updatedBy?: string | number | null | undefined;
  updatedBy_gt?: string | number | null | undefined;
  updatedBy_gte?: string | number | null | undefined;
  updatedBy_in?: Array<string | number> | null | undefined;
  updatedBy_lt?: string | number | null | undefined;
  updatedBy_lte?: string | number | null | undefined;
  updatedBy_ne?: string | number | null | undefined;
  updatedBy_null?: boolean | null | undefined;
  version?: number | null | undefined;
  version_gt?: number | null | undefined;
  version_gte?: number | null | undefined;
  version_in?: Array<number> | null | undefined;
  version_lt?: number | null | undefined;
  version_lte?: number | null | undefined;
  version_ne?: number | null | undefined;
  version_null?: boolean | null | undefined;
  weight?: number | null | undefined;
  weight_gt?: number | null | undefined;
  weight_gte?: number | null | undefined;
  weight_in?: Array<number> | null | undefined;
  weight_lt?: number | null | undefined;
  weight_lte?: number | null | undefined;
  weight_ne?: number | null | undefined;
  weight_null?: boolean | null | undefined;
};

export type InviteOperatorInput = {
  displayName: string;
  email?: string | null | undefined;
  phone: string;
  roleIds: Array<string | number>;
  storeAccessMode: StoreAccessMode;
  storeIds: Array<string | number>;
};

export type LoginInput = {
  password: string;
  phone: string;
};

export type MembershipInvitationFilterType = {
  AND?: Array<MembershipInvitationFilterType> | null | undefined;
  OR?: Array<MembershipInvitationFilterType> | null | undefined;
  acceptedAt?: unknown;
  acceptedAt_gt?: unknown;
  acceptedAt_gte?: unknown;
  acceptedAt_in?: Array<unknown> | null | undefined;
  acceptedAt_lt?: unknown;
  acceptedAt_lte?: unknown;
  acceptedAt_ne?: unknown;
  acceptedAt_null?: boolean | null | undefined;
  createdAt?: number | null | undefined;
  createdAt_gt?: number | null | undefined;
  createdAt_gte?: number | null | undefined;
  createdAt_in?: Array<number> | null | undefined;
  createdAt_lt?: number | null | undefined;
  createdAt_lte?: number | null | undefined;
  createdAt_ne?: number | null | undefined;
  createdAt_null?: boolean | null | undefined;
  createdBy?: string | number | null | undefined;
  createdBy_gt?: string | number | null | undefined;
  createdBy_gte?: string | number | null | undefined;
  createdBy_in?: Array<string | number> | null | undefined;
  createdBy_lt?: string | number | null | undefined;
  createdBy_lte?: string | number | null | undefined;
  createdBy_ne?: string | number | null | undefined;
  createdBy_null?: boolean | null | undefined;
  deletedAt?: number | null | undefined;
  deletedAt_gt?: number | null | undefined;
  deletedAt_gte?: number | null | undefined;
  deletedAt_in?: Array<number> | null | undefined;
  deletedAt_lt?: number | null | undefined;
  deletedAt_lte?: number | null | undefined;
  deletedAt_ne?: number | null | undefined;
  deletedAt_null?: boolean | null | undefined;
  deletedBy?: string | number | null | undefined;
  deletedBy_gt?: string | number | null | undefined;
  deletedBy_gte?: string | number | null | undefined;
  deletedBy_in?: Array<string | number> | null | undefined;
  deletedBy_lt?: string | number | null | undefined;
  deletedBy_lte?: string | number | null | undefined;
  deletedBy_ne?: string | number | null | undefined;
  deletedBy_null?: boolean | null | undefined;
  expiresAt?: unknown;
  expiresAt_gt?: unknown;
  expiresAt_gte?: unknown;
  expiresAt_in?: Array<unknown> | null | undefined;
  expiresAt_lt?: unknown;
  expiresAt_lte?: unknown;
  expiresAt_ne?: unknown;
  expiresAt_null?: boolean | null | undefined;
  id?: string | number | null | undefined;
  id_gt?: string | number | null | undefined;
  id_gte?: string | number | null | undefined;
  id_in?: Array<string | number> | null | undefined;
  id_lt?: string | number | null | undefined;
  id_lte?: string | number | null | undefined;
  id_ne?: string | number | null | undefined;
  id_null?: boolean | null | undefined;
  invitedByAccount?: AccountFilterType | null | undefined;
  invitedByAccountId?: string | number | null | undefined;
  invitedByAccountId_gt?: string | number | null | undefined;
  invitedByAccountId_gte?: string | number | null | undefined;
  invitedByAccountId_in?: Array<string | number> | null | undefined;
  invitedByAccountId_lt?: string | number | null | undefined;
  invitedByAccountId_lte?: string | number | null | undefined;
  invitedByAccountId_ne?: string | number | null | undefined;
  invitedByAccountId_null?: boolean | null | undefined;
  isDelete?: number | null | undefined;
  isDelete_gt?: number | null | undefined;
  isDelete_gte?: number | null | undefined;
  isDelete_in?: Array<number> | null | undefined;
  isDelete_lt?: number | null | undefined;
  isDelete_lte?: number | null | undefined;
  isDelete_ne?: number | null | undefined;
  isDelete_null?: boolean | null | undefined;
  membership?: OperatorMembershipFilterType | null | undefined;
  membershipId?: string | number | null | undefined;
  membershipId_gt?: string | number | null | undefined;
  membershipId_gte?: string | number | null | undefined;
  membershipId_in?: Array<string | number> | null | undefined;
  membershipId_lt?: string | number | null | undefined;
  membershipId_lte?: string | number | null | undefined;
  membershipId_ne?: string | number | null | undefined;
  membershipId_null?: boolean | null | undefined;
  revokedAt?: unknown;
  revokedAt_gt?: unknown;
  revokedAt_gte?: unknown;
  revokedAt_in?: Array<unknown> | null | undefined;
  revokedAt_lt?: unknown;
  revokedAt_lte?: unknown;
  revokedAt_ne?: unknown;
  revokedAt_null?: boolean | null | undefined;
  state?: number | null | undefined;
  state_gt?: number | null | undefined;
  state_gte?: number | null | undefined;
  state_in?: Array<number> | null | undefined;
  state_lt?: number | null | undefined;
  state_lte?: number | null | undefined;
  state_ne?: number | null | undefined;
  state_null?: boolean | null | undefined;
  updatedAt?: number | null | undefined;
  updatedAt_gt?: number | null | undefined;
  updatedAt_gte?: number | null | undefined;
  updatedAt_in?: Array<number> | null | undefined;
  updatedAt_lt?: number | null | undefined;
  updatedAt_lte?: number | null | undefined;
  updatedAt_ne?: number | null | undefined;
  updatedAt_null?: boolean | null | undefined;
  updatedBy?: string | number | null | undefined;
  updatedBy_gt?: string | number | null | undefined;
  updatedBy_gte?: string | number | null | undefined;
  updatedBy_in?: Array<string | number> | null | undefined;
  updatedBy_lt?: string | number | null | undefined;
  updatedBy_lte?: string | number | null | undefined;
  updatedBy_ne?: string | number | null | undefined;
  updatedBy_null?: boolean | null | undefined;
  weight?: number | null | undefined;
  weight_gt?: number | null | undefined;
  weight_gte?: number | null | undefined;
  weight_in?: Array<number> | null | undefined;
  weight_lt?: number | null | undefined;
  weight_lte?: number | null | undefined;
  weight_ne?: number | null | undefined;
  weight_null?: boolean | null | undefined;
};

export type MembershipInvitationRelationship = {
  acceptedAt?: unknown;
  expiresAt?: unknown;
  id?: string | number | null | undefined;
  invitedByAccountId?: string | number | null | undefined;
  isDelete?: number | null | undefined;
  membershipId?: string | number | null | undefined;
  revokedAt?: unknown;
  state?: number | null | undefined;
  weight?: number | null | undefined;
};

export type MembershipStatus =
  | 'ACTIVE'
  | 'INVITED'
  | 'LEFT'
  | 'SUSPENDED';

export type OperatorMembershipFilterType = {
  AND?: Array<OperatorMembershipFilterType> | null | undefined;
  OR?: Array<OperatorMembershipFilterType> | null | undefined;
  acceptedAt?: unknown;
  acceptedAt_gt?: unknown;
  acceptedAt_gte?: unknown;
  acceptedAt_in?: Array<unknown> | null | undefined;
  acceptedAt_lt?: unknown;
  acceptedAt_lte?: unknown;
  acceptedAt_ne?: unknown;
  acceptedAt_null?: boolean | null | undefined;
  account?: AccountFilterType | null | undefined;
  accountId?: string | number | null | undefined;
  accountId_gt?: string | number | null | undefined;
  accountId_gte?: string | number | null | undefined;
  accountId_in?: Array<string | number> | null | undefined;
  accountId_lt?: string | number | null | undefined;
  accountId_lte?: string | number | null | undefined;
  accountId_ne?: string | number | null | undefined;
  accountId_null?: boolean | null | undefined;
  createdAt?: number | null | undefined;
  createdAt_gt?: number | null | undefined;
  createdAt_gte?: number | null | undefined;
  createdAt_in?: Array<number> | null | undefined;
  createdAt_lt?: number | null | undefined;
  createdAt_lte?: number | null | undefined;
  createdAt_ne?: number | null | undefined;
  createdAt_null?: boolean | null | undefined;
  createdBy?: string | number | null | undefined;
  createdBy_gt?: string | number | null | undefined;
  createdBy_gte?: string | number | null | undefined;
  createdBy_in?: Array<string | number> | null | undefined;
  createdBy_lt?: string | number | null | undefined;
  createdBy_lte?: string | number | null | undefined;
  createdBy_ne?: string | number | null | undefined;
  createdBy_null?: boolean | null | undefined;
  deletedAt?: number | null | undefined;
  deletedAt_gt?: number | null | undefined;
  deletedAt_gte?: number | null | undefined;
  deletedAt_in?: Array<number> | null | undefined;
  deletedAt_lt?: number | null | undefined;
  deletedAt_lte?: number | null | undefined;
  deletedAt_ne?: number | null | undefined;
  deletedAt_null?: boolean | null | undefined;
  deletedBy?: string | number | null | undefined;
  deletedBy_gt?: string | number | null | undefined;
  deletedBy_gte?: string | number | null | undefined;
  deletedBy_in?: Array<string | number> | null | undefined;
  deletedBy_lt?: string | number | null | undefined;
  deletedBy_lte?: string | number | null | undefined;
  deletedBy_ne?: string | number | null | undefined;
  deletedBy_null?: boolean | null | undefined;
  id?: string | number | null | undefined;
  id_gt?: string | number | null | undefined;
  id_gte?: string | number | null | undefined;
  id_in?: Array<string | number> | null | undefined;
  id_lt?: string | number | null | undefined;
  id_lte?: string | number | null | undefined;
  id_ne?: string | number | null | undefined;
  id_null?: boolean | null | undefined;
  invitations?: MembershipInvitationFilterType | null | undefined;
  invitedAt?: unknown;
  invitedAt_gt?: unknown;
  invitedAt_gte?: unknown;
  invitedAt_in?: Array<unknown> | null | undefined;
  invitedAt_lt?: unknown;
  invitedAt_lte?: unknown;
  invitedAt_ne?: unknown;
  invitedAt_null?: boolean | null | undefined;
  isDelete?: number | null | undefined;
  isDelete_gt?: number | null | undefined;
  isDelete_gte?: number | null | undefined;
  isDelete_in?: Array<number> | null | undefined;
  isDelete_lt?: number | null | undefined;
  isDelete_lte?: number | null | undefined;
  isDelete_ne?: number | null | undefined;
  isDelete_null?: boolean | null | undefined;
  organization?: OrganizationFilterType | null | undefined;
  organizationId?: string | number | null | undefined;
  organizationId_gt?: string | number | null | undefined;
  organizationId_gte?: string | number | null | undefined;
  organizationId_in?: Array<string | number> | null | undefined;
  organizationId_lt?: string | number | null | undefined;
  organizationId_lte?: string | number | null | undefined;
  organizationId_ne?: string | number | null | undefined;
  organizationId_null?: boolean | null | undefined;
  roles?: OperatorRoleFilterType | null | undefined;
  state?: number | null | undefined;
  state_gt?: number | null | undefined;
  state_gte?: number | null | undefined;
  state_in?: Array<number> | null | undefined;
  state_lt?: number | null | undefined;
  state_lte?: number | null | undefined;
  state_ne?: number | null | undefined;
  state_null?: boolean | null | undefined;
  status?: MembershipStatus | null | undefined;
  status_gt?: MembershipStatus | null | undefined;
  status_gte?: MembershipStatus | null | undefined;
  status_in?: Array<MembershipStatus> | null | undefined;
  status_lt?: MembershipStatus | null | undefined;
  status_lte?: MembershipStatus | null | undefined;
  status_ne?: MembershipStatus | null | undefined;
  status_null?: boolean | null | undefined;
  storeAccessMode?: StoreAccessMode | null | undefined;
  storeAccessMode_gt?: StoreAccessMode | null | undefined;
  storeAccessMode_gte?: StoreAccessMode | null | undefined;
  storeAccessMode_in?: Array<StoreAccessMode> | null | undefined;
  storeAccessMode_lt?: StoreAccessMode | null | undefined;
  storeAccessMode_lte?: StoreAccessMode | null | undefined;
  storeAccessMode_ne?: StoreAccessMode | null | undefined;
  storeAccessMode_null?: boolean | null | undefined;
  stores?: StoreFilterType | null | undefined;
  updatedAt?: number | null | undefined;
  updatedAt_gt?: number | null | undefined;
  updatedAt_gte?: number | null | undefined;
  updatedAt_in?: Array<number> | null | undefined;
  updatedAt_lt?: number | null | undefined;
  updatedAt_lte?: number | null | undefined;
  updatedAt_ne?: number | null | undefined;
  updatedAt_null?: boolean | null | undefined;
  updatedBy?: string | number | null | undefined;
  updatedBy_gt?: string | number | null | undefined;
  updatedBy_gte?: string | number | null | undefined;
  updatedBy_in?: Array<string | number> | null | undefined;
  updatedBy_lt?: string | number | null | undefined;
  updatedBy_lte?: string | number | null | undefined;
  updatedBy_ne?: string | number | null | undefined;
  updatedBy_null?: boolean | null | undefined;
  weight?: number | null | undefined;
  weight_gt?: number | null | undefined;
  weight_gte?: number | null | undefined;
  weight_in?: Array<number> | null | undefined;
  weight_lt?: number | null | undefined;
  weight_lte?: number | null | undefined;
  weight_ne?: number | null | undefined;
  weight_null?: boolean | null | undefined;
};

export type OperatorMembershipRelationship = {
  acceptedAt?: unknown;
  accountId?: string | number | null | undefined;
  id?: string | number | null | undefined;
  invitedAt?: unknown;
  isDelete?: number | null | undefined;
  organizationId?: string | number | null | undefined;
  state?: number | null | undefined;
  status?: MembershipStatus | null | undefined;
  storeAccessMode?: StoreAccessMode | null | undefined;
  weight?: number | null | undefined;
};

export type OperatorRoleFilterType = {
  AND?: Array<OperatorRoleFilterType> | null | undefined;
  OR?: Array<OperatorRoleFilterType> | null | undefined;
  createdAt?: number | null | undefined;
  createdAt_gt?: number | null | undefined;
  createdAt_gte?: number | null | undefined;
  createdAt_in?: Array<number> | null | undefined;
  createdAt_lt?: number | null | undefined;
  createdAt_lte?: number | null | undefined;
  createdAt_ne?: number | null | undefined;
  createdAt_null?: boolean | null | undefined;
  createdBy?: string | number | null | undefined;
  createdBy_gt?: string | number | null | undefined;
  createdBy_gte?: string | number | null | undefined;
  createdBy_in?: Array<string | number> | null | undefined;
  createdBy_lt?: string | number | null | undefined;
  createdBy_lte?: string | number | null | undefined;
  createdBy_ne?: string | number | null | undefined;
  createdBy_null?: boolean | null | undefined;
  deletedAt?: number | null | undefined;
  deletedAt_gt?: number | null | undefined;
  deletedAt_gte?: number | null | undefined;
  deletedAt_in?: Array<number> | null | undefined;
  deletedAt_lt?: number | null | undefined;
  deletedAt_lte?: number | null | undefined;
  deletedAt_ne?: number | null | undefined;
  deletedAt_null?: boolean | null | undefined;
  deletedBy?: string | number | null | undefined;
  deletedBy_gt?: string | number | null | undefined;
  deletedBy_gte?: string | number | null | undefined;
  deletedBy_in?: Array<string | number> | null | undefined;
  deletedBy_lt?: string | number | null | undefined;
  deletedBy_lte?: string | number | null | undefined;
  deletedBy_ne?: string | number | null | undefined;
  deletedBy_null?: boolean | null | undefined;
  id?: string | number | null | undefined;
  id_gt?: string | number | null | undefined;
  id_gte?: string | number | null | undefined;
  id_in?: Array<string | number> | null | undefined;
  id_lt?: string | number | null | undefined;
  id_lte?: string | number | null | undefined;
  id_ne?: string | number | null | undefined;
  id_null?: boolean | null | undefined;
  isDelete?: number | null | undefined;
  isDelete_gt?: number | null | undefined;
  isDelete_gte?: number | null | undefined;
  isDelete_in?: Array<number> | null | undefined;
  isDelete_lt?: number | null | undefined;
  isDelete_lte?: number | null | undefined;
  isDelete_ne?: number | null | undefined;
  isDelete_null?: boolean | null | undefined;
  kind?: RoleKind | null | undefined;
  kind_gt?: RoleKind | null | undefined;
  kind_gte?: RoleKind | null | undefined;
  kind_in?: Array<RoleKind> | null | undefined;
  kind_lt?: RoleKind | null | undefined;
  kind_lte?: RoleKind | null | undefined;
  kind_ne?: RoleKind | null | undefined;
  kind_null?: boolean | null | undefined;
  members?: OperatorMembershipFilterType | null | undefined;
  name?: string | null | undefined;
  name_gt?: string | null | undefined;
  name_gte?: string | null | undefined;
  name_in?: Array<string> | null | undefined;
  name_like?: string | null | undefined;
  name_lt?: string | null | undefined;
  name_lte?: string | null | undefined;
  name_ne?: string | null | undefined;
  name_null?: boolean | null | undefined;
  name_prefix?: string | null | undefined;
  name_suffix?: string | null | undefined;
  organization?: OrganizationFilterType | null | undefined;
  organizationId?: string | number | null | undefined;
  organizationId_gt?: string | number | null | undefined;
  organizationId_gte?: string | number | null | undefined;
  organizationId_in?: Array<string | number> | null | undefined;
  organizationId_lt?: string | number | null | undefined;
  organizationId_lte?: string | number | null | undefined;
  organizationId_ne?: string | number | null | undefined;
  organizationId_null?: boolean | null | undefined;
  permissions?: PermissionFilterType | null | undefined;
  state?: number | null | undefined;
  state_gt?: number | null | undefined;
  state_gte?: number | null | undefined;
  state_in?: Array<number> | null | undefined;
  state_lt?: number | null | undefined;
  state_lte?: number | null | undefined;
  state_ne?: number | null | undefined;
  state_null?: boolean | null | undefined;
  updatedAt?: number | null | undefined;
  updatedAt_gt?: number | null | undefined;
  updatedAt_gte?: number | null | undefined;
  updatedAt_in?: Array<number> | null | undefined;
  updatedAt_lt?: number | null | undefined;
  updatedAt_lte?: number | null | undefined;
  updatedAt_ne?: number | null | undefined;
  updatedAt_null?: boolean | null | undefined;
  updatedBy?: string | number | null | undefined;
  updatedBy_gt?: string | number | null | undefined;
  updatedBy_gte?: string | number | null | undefined;
  updatedBy_in?: Array<string | number> | null | undefined;
  updatedBy_lt?: string | number | null | undefined;
  updatedBy_lte?: string | number | null | undefined;
  updatedBy_ne?: string | number | null | undefined;
  updatedBy_null?: boolean | null | undefined;
  weight?: number | null | undefined;
  weight_gt?: number | null | undefined;
  weight_gte?: number | null | undefined;
  weight_in?: Array<number> | null | undefined;
  weight_lt?: number | null | undefined;
  weight_lte?: number | null | undefined;
  weight_ne?: number | null | undefined;
  weight_null?: boolean | null | undefined;
};

export type OperatorRoleRelationship = {
  id?: string | number | null | undefined;
  isDelete?: number | null | undefined;
  kind?: RoleKind | null | undefined;
  name?: string | null | undefined;
  organizationId?: string | number | null | undefined;
  state?: number | null | undefined;
  weight?: number | null | undefined;
};

export type OrganizationFilterType = {
  AND?: Array<OrganizationFilterType> | null | undefined;
  OR?: Array<OrganizationFilterType> | null | undefined;
  auditLogs?: AuditLogFilterType | null | undefined;
  code?: string | null | undefined;
  code_gt?: string | null | undefined;
  code_gte?: string | null | undefined;
  code_in?: Array<string> | null | undefined;
  code_like?: string | null | undefined;
  code_lt?: string | null | undefined;
  code_lte?: string | null | undefined;
  code_ne?: string | null | undefined;
  code_null?: boolean | null | undefined;
  code_prefix?: string | null | undefined;
  code_suffix?: string | null | undefined;
  createdAt?: number | null | undefined;
  createdAt_gt?: number | null | undefined;
  createdAt_gte?: number | null | undefined;
  createdAt_in?: Array<number> | null | undefined;
  createdAt_lt?: number | null | undefined;
  createdAt_lte?: number | null | undefined;
  createdAt_ne?: number | null | undefined;
  createdAt_null?: boolean | null | undefined;
  createdBy?: string | number | null | undefined;
  createdBy_gt?: string | number | null | undefined;
  createdBy_gte?: string | number | null | undefined;
  createdBy_in?: Array<string | number> | null | undefined;
  createdBy_lt?: string | number | null | undefined;
  createdBy_lte?: string | number | null | undefined;
  createdBy_ne?: string | number | null | undefined;
  createdBy_null?: boolean | null | undefined;
  deletedAt?: number | null | undefined;
  deletedAt_gt?: number | null | undefined;
  deletedAt_gte?: number | null | undefined;
  deletedAt_in?: Array<number> | null | undefined;
  deletedAt_lt?: number | null | undefined;
  deletedAt_lte?: number | null | undefined;
  deletedAt_ne?: number | null | undefined;
  deletedAt_null?: boolean | null | undefined;
  deletedBy?: string | number | null | undefined;
  deletedBy_gt?: string | number | null | undefined;
  deletedBy_gte?: string | number | null | undefined;
  deletedBy_in?: Array<string | number> | null | undefined;
  deletedBy_lt?: string | number | null | undefined;
  deletedBy_lte?: string | number | null | undefined;
  deletedBy_ne?: string | number | null | undefined;
  deletedBy_null?: boolean | null | undefined;
  id?: string | number | null | undefined;
  id_gt?: string | number | null | undefined;
  id_gte?: string | number | null | undefined;
  id_in?: Array<string | number> | null | undefined;
  id_lt?: string | number | null | undefined;
  id_lte?: string | number | null | undefined;
  id_ne?: string | number | null | undefined;
  id_null?: boolean | null | undefined;
  initialAccount?: AccountFilterType | null | undefined;
  initialAccountId?: string | number | null | undefined;
  initialAccountId_gt?: string | number | null | undefined;
  initialAccountId_gte?: string | number | null | undefined;
  initialAccountId_in?: Array<string | number> | null | undefined;
  initialAccountId_lt?: string | number | null | undefined;
  initialAccountId_lte?: string | number | null | undefined;
  initialAccountId_ne?: string | number | null | undefined;
  initialAccountId_null?: boolean | null | undefined;
  isDelete?: number | null | undefined;
  isDelete_gt?: number | null | undefined;
  isDelete_gte?: number | null | undefined;
  isDelete_in?: Array<number> | null | undefined;
  isDelete_lt?: number | null | undefined;
  isDelete_lte?: number | null | undefined;
  isDelete_ne?: number | null | undefined;
  isDelete_null?: boolean | null | undefined;
  memberships?: OperatorMembershipFilterType | null | undefined;
  name?: string | null | undefined;
  name_gt?: string | null | undefined;
  name_gte?: string | null | undefined;
  name_in?: Array<string> | null | undefined;
  name_like?: string | null | undefined;
  name_lt?: string | null | undefined;
  name_lte?: string | null | undefined;
  name_ne?: string | null | undefined;
  name_null?: boolean | null | undefined;
  name_prefix?: string | null | undefined;
  name_suffix?: string | null | undefined;
  openingRecords?: FranchiseOpeningRecordFilterType | null | undefined;
  paymentConfigs?: FranchisePaymentConfigFilterType | null | undefined;
  roles?: OperatorRoleFilterType | null | undefined;
  sessions?: SessionFilterType | null | undefined;
  state?: number | null | undefined;
  state_gt?: number | null | undefined;
  state_gte?: number | null | undefined;
  state_in?: Array<number> | null | undefined;
  state_lt?: number | null | undefined;
  state_lte?: number | null | undefined;
  state_ne?: number | null | undefined;
  state_null?: boolean | null | undefined;
  status?: OrganizationStatus | null | undefined;
  status_gt?: OrganizationStatus | null | undefined;
  status_gte?: OrganizationStatus | null | undefined;
  status_in?: Array<OrganizationStatus> | null | undefined;
  status_lt?: OrganizationStatus | null | undefined;
  status_lte?: OrganizationStatus | null | undefined;
  status_ne?: OrganizationStatus | null | undefined;
  status_null?: boolean | null | undefined;
  stores?: StoreFilterType | null | undefined;
  suspendedAt?: unknown;
  suspendedAt_gt?: unknown;
  suspendedAt_gte?: unknown;
  suspendedAt_in?: Array<unknown> | null | undefined;
  suspendedAt_lt?: unknown;
  suspendedAt_lte?: unknown;
  suspendedAt_ne?: unknown;
  suspendedAt_null?: boolean | null | undefined;
  suspensionReasonCode?: string | null | undefined;
  suspensionReasonCode_gt?: string | null | undefined;
  suspensionReasonCode_gte?: string | null | undefined;
  suspensionReasonCode_in?: Array<string> | null | undefined;
  suspensionReasonCode_like?: string | null | undefined;
  suspensionReasonCode_lt?: string | null | undefined;
  suspensionReasonCode_lte?: string | null | undefined;
  suspensionReasonCode_ne?: string | null | undefined;
  suspensionReasonCode_null?: boolean | null | undefined;
  suspensionReasonCode_prefix?: string | null | undefined;
  suspensionReasonCode_suffix?: string | null | undefined;
  type?: OrganizationType | null | undefined;
  type_gt?: OrganizationType | null | undefined;
  type_gte?: OrganizationType | null | undefined;
  type_in?: Array<OrganizationType> | null | undefined;
  type_lt?: OrganizationType | null | undefined;
  type_lte?: OrganizationType | null | undefined;
  type_ne?: OrganizationType | null | undefined;
  type_null?: boolean | null | undefined;
  updatedAt?: number | null | undefined;
  updatedAt_gt?: number | null | undefined;
  updatedAt_gte?: number | null | undefined;
  updatedAt_in?: Array<number> | null | undefined;
  updatedAt_lt?: number | null | undefined;
  updatedAt_lte?: number | null | undefined;
  updatedAt_ne?: number | null | undefined;
  updatedAt_null?: boolean | null | undefined;
  updatedBy?: string | number | null | undefined;
  updatedBy_gt?: string | number | null | undefined;
  updatedBy_gte?: string | number | null | undefined;
  updatedBy_in?: Array<string | number> | null | undefined;
  updatedBy_lt?: string | number | null | undefined;
  updatedBy_lte?: string | number | null | undefined;
  updatedBy_ne?: string | number | null | undefined;
  updatedBy_null?: boolean | null | undefined;
  weight?: number | null | undefined;
  weight_gt?: number | null | undefined;
  weight_gte?: number | null | undefined;
  weight_in?: Array<number> | null | undefined;
  weight_lt?: number | null | undefined;
  weight_lte?: number | null | undefined;
  weight_ne?: number | null | undefined;
  weight_null?: boolean | null | undefined;
};

export type OrganizationRelationship = {
  code?: string | null | undefined;
  id?: string | number | null | undefined;
  initialAccountId?: string | number | null | undefined;
  isDelete?: number | null | undefined;
  name?: string | null | undefined;
  state?: number | null | undefined;
  status?: OrganizationStatus | null | undefined;
  suspendedAt?: unknown;
  suspensionReasonCode?: string | null | undefined;
  type?: OrganizationType | null | undefined;
  weight?: number | null | undefined;
};

export type OrganizationStatus =
  | 'ACTIVE'
  | 'SUSPENDED';

export type OrganizationType =
  | 'FRANCHISE'
  | 'HEADQUARTERS';

export type PermissionFilterType = {
  AND?: Array<PermissionFilterType> | null | undefined;
  OR?: Array<PermissionFilterType> | null | undefined;
  action?: string | null | undefined;
  action_gt?: string | null | undefined;
  action_gte?: string | null | undefined;
  action_in?: Array<string> | null | undefined;
  action_like?: string | null | undefined;
  action_lt?: string | null | undefined;
  action_lte?: string | null | undefined;
  action_ne?: string | null | undefined;
  action_null?: boolean | null | undefined;
  action_prefix?: string | null | undefined;
  action_suffix?: string | null | undefined;
  createdAt?: number | null | undefined;
  createdAt_gt?: number | null | undefined;
  createdAt_gte?: number | null | undefined;
  createdAt_in?: Array<number> | null | undefined;
  createdAt_lt?: number | null | undefined;
  createdAt_lte?: number | null | undefined;
  createdAt_ne?: number | null | undefined;
  createdAt_null?: boolean | null | undefined;
  createdBy?: string | number | null | undefined;
  createdBy_gt?: string | number | null | undefined;
  createdBy_gte?: string | number | null | undefined;
  createdBy_in?: Array<string | number> | null | undefined;
  createdBy_lt?: string | number | null | undefined;
  createdBy_lte?: string | number | null | undefined;
  createdBy_ne?: string | number | null | undefined;
  createdBy_null?: boolean | null | undefined;
  deletedAt?: number | null | undefined;
  deletedAt_gt?: number | null | undefined;
  deletedAt_gte?: number | null | undefined;
  deletedAt_in?: Array<number> | null | undefined;
  deletedAt_lt?: number | null | undefined;
  deletedAt_lte?: number | null | undefined;
  deletedAt_ne?: number | null | undefined;
  deletedAt_null?: boolean | null | undefined;
  deletedBy?: string | number | null | undefined;
  deletedBy_gt?: string | number | null | undefined;
  deletedBy_gte?: string | number | null | undefined;
  deletedBy_in?: Array<string | number> | null | undefined;
  deletedBy_lt?: string | number | null | undefined;
  deletedBy_lte?: string | number | null | undefined;
  deletedBy_ne?: string | number | null | undefined;
  deletedBy_null?: boolean | null | undefined;
  id?: string | number | null | undefined;
  id_gt?: string | number | null | undefined;
  id_gte?: string | number | null | undefined;
  id_in?: Array<string | number> | null | undefined;
  id_lt?: string | number | null | undefined;
  id_lte?: string | number | null | undefined;
  id_ne?: string | number | null | undefined;
  id_null?: boolean | null | undefined;
  isDelete?: number | null | undefined;
  isDelete_gt?: number | null | undefined;
  isDelete_gte?: number | null | undefined;
  isDelete_in?: Array<number> | null | undefined;
  isDelete_lt?: number | null | undefined;
  isDelete_lte?: number | null | undefined;
  isDelete_ne?: number | null | undefined;
  isDelete_null?: boolean | null | undefined;
  module?: string | null | undefined;
  module_gt?: string | null | undefined;
  module_gte?: string | null | undefined;
  module_in?: Array<string> | null | undefined;
  module_like?: string | null | undefined;
  module_lt?: string | null | undefined;
  module_lte?: string | null | undefined;
  module_ne?: string | null | undefined;
  module_null?: boolean | null | undefined;
  module_prefix?: string | null | undefined;
  module_suffix?: string | null | undefined;
  name?: string | null | undefined;
  name_gt?: string | null | undefined;
  name_gte?: string | null | undefined;
  name_in?: Array<string> | null | undefined;
  name_like?: string | null | undefined;
  name_lt?: string | null | undefined;
  name_lte?: string | null | undefined;
  name_ne?: string | null | undefined;
  name_null?: boolean | null | undefined;
  name_prefix?: string | null | undefined;
  name_suffix?: string | null | undefined;
  roles?: OperatorRoleFilterType | null | undefined;
  scope?: PermissionScope | null | undefined;
  scope_gt?: PermissionScope | null | undefined;
  scope_gte?: PermissionScope | null | undefined;
  scope_in?: Array<PermissionScope> | null | undefined;
  scope_lt?: PermissionScope | null | undefined;
  scope_lte?: PermissionScope | null | undefined;
  scope_ne?: PermissionScope | null | undefined;
  scope_null?: boolean | null | undefined;
  state?: number | null | undefined;
  state_gt?: number | null | undefined;
  state_gte?: number | null | undefined;
  state_in?: Array<number> | null | undefined;
  state_lt?: number | null | undefined;
  state_lte?: number | null | undefined;
  state_ne?: number | null | undefined;
  state_null?: boolean | null | undefined;
  updatedAt?: number | null | undefined;
  updatedAt_gt?: number | null | undefined;
  updatedAt_gte?: number | null | undefined;
  updatedAt_in?: Array<number> | null | undefined;
  updatedAt_lt?: number | null | undefined;
  updatedAt_lte?: number | null | undefined;
  updatedAt_ne?: number | null | undefined;
  updatedAt_null?: boolean | null | undefined;
  updatedBy?: string | number | null | undefined;
  updatedBy_gt?: string | number | null | undefined;
  updatedBy_gte?: string | number | null | undefined;
  updatedBy_in?: Array<string | number> | null | undefined;
  updatedBy_lt?: string | number | null | undefined;
  updatedBy_lte?: string | number | null | undefined;
  updatedBy_ne?: string | number | null | undefined;
  updatedBy_null?: boolean | null | undefined;
  weight?: number | null | undefined;
  weight_gt?: number | null | undefined;
  weight_gte?: number | null | undefined;
  weight_in?: Array<number> | null | undefined;
  weight_lt?: number | null | undefined;
  weight_lte?: number | null | undefined;
  weight_ne?: number | null | undefined;
  weight_null?: boolean | null | undefined;
};

export type PermissionRelationship = {
  action?: string | null | undefined;
  id?: string | number | null | undefined;
  isDelete?: number | null | undefined;
  module?: string | null | undefined;
  name?: string | null | undefined;
  scope?: PermissionScope | null | undefined;
  state?: number | null | undefined;
  weight?: number | null | undefined;
};

export type PermissionScope =
  | 'SYSTEM'
  | 'TENANT';

export type ProvisionFranchiseInput = {
  code: string;
  name: string;
  ownerDisplayName: string;
  ownerEmail?: string | null | undefined;
  ownerPhone: string;
};

export type ReviewStoreInput = {
  approved: boolean;
  rejectionReason?: string | null | undefined;
  storeId: string | number;
};

export type RoleKind =
  | 'CUSTOM'
  | 'FRANCHISE_OWNER'
  | 'HQ_SUPER_ADMIN';

export type SelectWorkspaceInput = {
  organizationId?: string | number | null | undefined;
  workspaceType: WorkspaceType;
};

export type SessionEventCode =
  | 'CREDENTIALS_CHANGED'
  | 'ORGANIZATION_SUSPENDED'
  | 'SESSION_REVOKED';

export type SessionFilterType = {
  AND?: Array<SessionFilterType> | null | undefined;
  OR?: Array<SessionFilterType> | null | undefined;
  account?: AccountFilterType | null | undefined;
  accountId?: string | number | null | undefined;
  accountId_gt?: string | number | null | undefined;
  accountId_gte?: string | number | null | undefined;
  accountId_in?: Array<string | number> | null | undefined;
  accountId_lt?: string | number | null | undefined;
  accountId_lte?: string | number | null | undefined;
  accountId_ne?: string | number | null | undefined;
  accountId_null?: boolean | null | undefined;
  auditLogs?: AuditLogFilterType | null | undefined;
  createdAt?: number | null | undefined;
  createdAt_gt?: number | null | undefined;
  createdAt_gte?: number | null | undefined;
  createdAt_in?: Array<number> | null | undefined;
  createdAt_lt?: number | null | undefined;
  createdAt_lte?: number | null | undefined;
  createdAt_ne?: number | null | undefined;
  createdAt_null?: boolean | null | undefined;
  createdBy?: string | number | null | undefined;
  createdBy_gt?: string | number | null | undefined;
  createdBy_gte?: string | number | null | undefined;
  createdBy_in?: Array<string | number> | null | undefined;
  createdBy_lt?: string | number | null | undefined;
  createdBy_lte?: string | number | null | undefined;
  createdBy_ne?: string | number | null | undefined;
  createdBy_null?: boolean | null | undefined;
  credentialVersion?: number | null | undefined;
  credentialVersion_gt?: number | null | undefined;
  credentialVersion_gte?: number | null | undefined;
  credentialVersion_in?: Array<number> | null | undefined;
  credentialVersion_lt?: number | null | undefined;
  credentialVersion_lte?: number | null | undefined;
  credentialVersion_ne?: number | null | undefined;
  credentialVersion_null?: boolean | null | undefined;
  deletedAt?: number | null | undefined;
  deletedAt_gt?: number | null | undefined;
  deletedAt_gte?: number | null | undefined;
  deletedAt_in?: Array<number> | null | undefined;
  deletedAt_lt?: number | null | undefined;
  deletedAt_lte?: number | null | undefined;
  deletedAt_ne?: number | null | undefined;
  deletedAt_null?: boolean | null | undefined;
  deletedBy?: string | number | null | undefined;
  deletedBy_gt?: string | number | null | undefined;
  deletedBy_gte?: string | number | null | undefined;
  deletedBy_in?: Array<string | number> | null | undefined;
  deletedBy_lt?: string | number | null | undefined;
  deletedBy_lte?: string | number | null | undefined;
  deletedBy_ne?: string | number | null | undefined;
  deletedBy_null?: boolean | null | undefined;
  expiresAt?: unknown;
  expiresAt_gt?: unknown;
  expiresAt_gte?: unknown;
  expiresAt_in?: Array<unknown> | null | undefined;
  expiresAt_lt?: unknown;
  expiresAt_lte?: unknown;
  expiresAt_ne?: unknown;
  expiresAt_null?: boolean | null | undefined;
  id?: string | number | null | undefined;
  id_gt?: string | number | null | undefined;
  id_gte?: string | number | null | undefined;
  id_in?: Array<string | number> | null | undefined;
  id_lt?: string | number | null | undefined;
  id_lte?: string | number | null | undefined;
  id_ne?: string | number | null | undefined;
  id_null?: boolean | null | undefined;
  isDelete?: number | null | undefined;
  isDelete_gt?: number | null | undefined;
  isDelete_gte?: number | null | undefined;
  isDelete_in?: Array<number> | null | undefined;
  isDelete_lt?: number | null | undefined;
  isDelete_lte?: number | null | undefined;
  isDelete_ne?: number | null | undefined;
  isDelete_null?: boolean | null | undefined;
  lastSeenAt?: unknown;
  lastSeenAt_gt?: unknown;
  lastSeenAt_gte?: unknown;
  lastSeenAt_in?: Array<unknown> | null | undefined;
  lastSeenAt_lt?: unknown;
  lastSeenAt_lte?: unknown;
  lastSeenAt_ne?: unknown;
  lastSeenAt_null?: boolean | null | undefined;
  organization?: OrganizationFilterType | null | undefined;
  organizationId?: string | number | null | undefined;
  organizationId_gt?: string | number | null | undefined;
  organizationId_gte?: string | number | null | undefined;
  organizationId_in?: Array<string | number> | null | undefined;
  organizationId_lt?: string | number | null | undefined;
  organizationId_lte?: string | number | null | undefined;
  organizationId_ne?: string | number | null | undefined;
  organizationId_null?: boolean | null | undefined;
  revocationCode?: string | null | undefined;
  revocationCode_gt?: string | null | undefined;
  revocationCode_gte?: string | null | undefined;
  revocationCode_in?: Array<string> | null | undefined;
  revocationCode_like?: string | null | undefined;
  revocationCode_lt?: string | null | undefined;
  revocationCode_lte?: string | null | undefined;
  revocationCode_ne?: string | null | undefined;
  revocationCode_null?: boolean | null | undefined;
  revocationCode_prefix?: string | null | undefined;
  revocationCode_suffix?: string | null | undefined;
  revokedAt?: unknown;
  revokedAt_gt?: unknown;
  revokedAt_gte?: unknown;
  revokedAt_in?: Array<unknown> | null | undefined;
  revokedAt_lt?: unknown;
  revokedAt_lte?: unknown;
  revokedAt_ne?: unknown;
  revokedAt_null?: boolean | null | undefined;
  state?: number | null | undefined;
  state_gt?: number | null | undefined;
  state_gte?: number | null | undefined;
  state_in?: Array<number> | null | undefined;
  state_lt?: number | null | undefined;
  state_lte?: number | null | undefined;
  state_ne?: number | null | undefined;
  state_null?: boolean | null | undefined;
  updatedAt?: number | null | undefined;
  updatedAt_gt?: number | null | undefined;
  updatedAt_gte?: number | null | undefined;
  updatedAt_in?: Array<number> | null | undefined;
  updatedAt_lt?: number | null | undefined;
  updatedAt_lte?: number | null | undefined;
  updatedAt_ne?: number | null | undefined;
  updatedAt_null?: boolean | null | undefined;
  updatedBy?: string | number | null | undefined;
  updatedBy_gt?: string | number | null | undefined;
  updatedBy_gte?: string | number | null | undefined;
  updatedBy_in?: Array<string | number> | null | undefined;
  updatedBy_lt?: string | number | null | undefined;
  updatedBy_lte?: string | number | null | undefined;
  updatedBy_ne?: string | number | null | undefined;
  updatedBy_null?: boolean | null | undefined;
  weight?: number | null | undefined;
  weight_gt?: number | null | undefined;
  weight_gte?: number | null | undefined;
  weight_in?: Array<number> | null | undefined;
  weight_lt?: number | null | undefined;
  weight_lte?: number | null | undefined;
  weight_ne?: number | null | undefined;
  weight_null?: boolean | null | undefined;
  workspaceType?: WorkspaceType | null | undefined;
  workspaceType_gt?: WorkspaceType | null | undefined;
  workspaceType_gte?: WorkspaceType | null | undefined;
  workspaceType_in?: Array<WorkspaceType> | null | undefined;
  workspaceType_lt?: WorkspaceType | null | undefined;
  workspaceType_lte?: WorkspaceType | null | undefined;
  workspaceType_ne?: WorkspaceType | null | undefined;
  workspaceType_null?: boolean | null | undefined;
};

export type StoreAccessMode =
  | 'ALL_STORES'
  | 'SELECTED_STORES';

export type StoreBusinessStatus =
  | 'CLOSED'
  | 'OPEN';

export type StoreDocumentKind =
  | 'BUSINESS_LICENSE'
  | 'OTHER';

export type StoreFilterType = {
  AND?: Array<StoreFilterType> | null | undefined;
  OR?: Array<StoreFilterType> | null | undefined;
  address?: string | null | undefined;
  address_gt?: string | null | undefined;
  address_gte?: string | null | undefined;
  address_in?: Array<string> | null | undefined;
  address_like?: string | null | undefined;
  address_lt?: string | null | undefined;
  address_lte?: string | null | undefined;
  address_ne?: string | null | undefined;
  address_null?: boolean | null | undefined;
  address_prefix?: string | null | undefined;
  address_suffix?: string | null | undefined;
  auditLogs?: AuditLogFilterType | null | undefined;
  businessHours?: string | null | undefined;
  businessHours_gt?: string | null | undefined;
  businessHours_gte?: string | null | undefined;
  businessHours_in?: Array<string> | null | undefined;
  businessHours_like?: string | null | undefined;
  businessHours_lt?: string | null | undefined;
  businessHours_lte?: string | null | undefined;
  businessHours_ne?: string | null | undefined;
  businessHours_null?: boolean | null | undefined;
  businessHours_prefix?: string | null | undefined;
  businessHours_suffix?: string | null | undefined;
  businessLicenseImageUrl?: string | null | undefined;
  businessLicenseImageUrl_gt?: string | null | undefined;
  businessLicenseImageUrl_gte?: string | null | undefined;
  businessLicenseImageUrl_in?: Array<string> | null | undefined;
  businessLicenseImageUrl_like?: string | null | undefined;
  businessLicenseImageUrl_lt?: string | null | undefined;
  businessLicenseImageUrl_lte?: string | null | undefined;
  businessLicenseImageUrl_ne?: string | null | undefined;
  businessLicenseImageUrl_null?: boolean | null | undefined;
  businessLicenseImageUrl_prefix?: string | null | undefined;
  businessLicenseImageUrl_suffix?: string | null | undefined;
  businessStatus?: StoreBusinessStatus | null | undefined;
  businessStatus_gt?: StoreBusinessStatus | null | undefined;
  businessStatus_gte?: StoreBusinessStatus | null | undefined;
  businessStatus_in?: Array<StoreBusinessStatus> | null | undefined;
  businessStatus_lt?: StoreBusinessStatus | null | undefined;
  businessStatus_lte?: StoreBusinessStatus | null | undefined;
  businessStatus_ne?: StoreBusinessStatus | null | undefined;
  businessStatus_null?: boolean | null | undefined;
  city?: string | null | undefined;
  city_gt?: string | null | undefined;
  city_gte?: string | null | undefined;
  city_in?: Array<string> | null | undefined;
  city_like?: string | null | undefined;
  city_lt?: string | null | undefined;
  city_lte?: string | null | undefined;
  city_ne?: string | null | undefined;
  city_null?: boolean | null | undefined;
  city_prefix?: string | null | undefined;
  city_suffix?: string | null | undefined;
  code?: string | null | undefined;
  code_gt?: string | null | undefined;
  code_gte?: string | null | undefined;
  code_in?: Array<string> | null | undefined;
  code_like?: string | null | undefined;
  code_lt?: string | null | undefined;
  code_lte?: string | null | undefined;
  code_ne?: string | null | undefined;
  code_null?: boolean | null | undefined;
  code_prefix?: string | null | undefined;
  code_suffix?: string | null | undefined;
  contactPhone?: string | null | undefined;
  contactPhone_gt?: string | null | undefined;
  contactPhone_gte?: string | null | undefined;
  contactPhone_in?: Array<string> | null | undefined;
  contactPhone_like?: string | null | undefined;
  contactPhone_lt?: string | null | undefined;
  contactPhone_lte?: string | null | undefined;
  contactPhone_ne?: string | null | undefined;
  contactPhone_null?: boolean | null | undefined;
  contactPhone_prefix?: string | null | undefined;
  contactPhone_suffix?: string | null | undefined;
  createdAt?: number | null | undefined;
  createdAt_gt?: number | null | undefined;
  createdAt_gte?: number | null | undefined;
  createdAt_in?: Array<number> | null | undefined;
  createdAt_lt?: number | null | undefined;
  createdAt_lte?: number | null | undefined;
  createdAt_ne?: number | null | undefined;
  createdAt_null?: boolean | null | undefined;
  createdBy?: string | number | null | undefined;
  createdBy_gt?: string | number | null | undefined;
  createdBy_gte?: string | number | null | undefined;
  createdBy_in?: Array<string | number> | null | undefined;
  createdBy_lt?: string | number | null | undefined;
  createdBy_lte?: string | number | null | undefined;
  createdBy_ne?: string | number | null | undefined;
  createdBy_null?: boolean | null | undefined;
  deletedAt?: number | null | undefined;
  deletedAt_gt?: number | null | undefined;
  deletedAt_gte?: number | null | undefined;
  deletedAt_in?: Array<number> | null | undefined;
  deletedAt_lt?: number | null | undefined;
  deletedAt_lte?: number | null | undefined;
  deletedAt_ne?: number | null | undefined;
  deletedAt_null?: boolean | null | undefined;
  deletedBy?: string | number | null | undefined;
  deletedBy_gt?: string | number | null | undefined;
  deletedBy_gte?: string | number | null | undefined;
  deletedBy_in?: Array<string | number> | null | undefined;
  deletedBy_lt?: string | number | null | undefined;
  deletedBy_lte?: string | number | null | undefined;
  deletedBy_ne?: string | number | null | undefined;
  deletedBy_null?: boolean | null | undefined;
  district?: string | null | undefined;
  district_gt?: string | null | undefined;
  district_gte?: string | null | undefined;
  district_in?: Array<string> | null | undefined;
  district_like?: string | null | undefined;
  district_lt?: string | null | undefined;
  district_lte?: string | null | undefined;
  district_ne?: string | null | undefined;
  district_null?: boolean | null | undefined;
  district_prefix?: string | null | undefined;
  district_suffix?: string | null | undefined;
  id?: string | number | null | undefined;
  id_gt?: string | number | null | undefined;
  id_gte?: string | number | null | undefined;
  id_in?: Array<string | number> | null | undefined;
  id_lt?: string | number | null | undefined;
  id_lte?: string | number | null | undefined;
  id_ne?: string | number | null | undefined;
  id_null?: boolean | null | undefined;
  isDelete?: number | null | undefined;
  isDelete_gt?: number | null | undefined;
  isDelete_gte?: number | null | undefined;
  isDelete_in?: Array<number> | null | undefined;
  isDelete_lt?: number | null | undefined;
  isDelete_lte?: number | null | undefined;
  isDelete_ne?: number | null | undefined;
  isDelete_null?: boolean | null | undefined;
  lifecycle?: StoreLifecycle | null | undefined;
  lifecycle_gt?: StoreLifecycle | null | undefined;
  lifecycle_gte?: StoreLifecycle | null | undefined;
  lifecycle_in?: Array<StoreLifecycle> | null | undefined;
  lifecycle_lt?: StoreLifecycle | null | undefined;
  lifecycle_lte?: StoreLifecycle | null | undefined;
  lifecycle_ne?: StoreLifecycle | null | undefined;
  lifecycle_null?: boolean | null | undefined;
  managerName?: string | null | undefined;
  managerName_gt?: string | null | undefined;
  managerName_gte?: string | null | undefined;
  managerName_in?: Array<string> | null | undefined;
  managerName_like?: string | null | undefined;
  managerName_lt?: string | null | undefined;
  managerName_lte?: string | null | undefined;
  managerName_ne?: string | null | undefined;
  managerName_null?: boolean | null | undefined;
  managerName_prefix?: string | null | undefined;
  managerName_suffix?: string | null | undefined;
  managerPhone?: string | null | undefined;
  managerPhone_gt?: string | null | undefined;
  managerPhone_gte?: string | null | undefined;
  managerPhone_in?: Array<string> | null | undefined;
  managerPhone_like?: string | null | undefined;
  managerPhone_lt?: string | null | undefined;
  managerPhone_lte?: string | null | undefined;
  managerPhone_ne?: string | null | undefined;
  managerPhone_null?: boolean | null | undefined;
  managerPhone_prefix?: string | null | undefined;
  managerPhone_suffix?: string | null | undefined;
  members?: OperatorMembershipFilterType | null | undefined;
  name?: string | null | undefined;
  name_gt?: string | null | undefined;
  name_gte?: string | null | undefined;
  name_in?: Array<string> | null | undefined;
  name_like?: string | null | undefined;
  name_lt?: string | null | undefined;
  name_lte?: string | null | undefined;
  name_ne?: string | null | undefined;
  name_null?: boolean | null | undefined;
  name_prefix?: string | null | undefined;
  name_suffix?: string | null | undefined;
  organization?: OrganizationFilterType | null | undefined;
  organizationId?: string | number | null | undefined;
  organizationId_gt?: string | number | null | undefined;
  organizationId_gte?: string | number | null | undefined;
  organizationId_in?: Array<string | number> | null | undefined;
  organizationId_lt?: string | number | null | undefined;
  organizationId_lte?: string | number | null | undefined;
  organizationId_ne?: string | number | null | undefined;
  organizationId_null?: boolean | null | undefined;
  otherDocumentImageUrl?: string | null | undefined;
  otherDocumentImageUrl_gt?: string | null | undefined;
  otherDocumentImageUrl_gte?: string | null | undefined;
  otherDocumentImageUrl_in?: Array<string> | null | undefined;
  otherDocumentImageUrl_like?: string | null | undefined;
  otherDocumentImageUrl_lt?: string | null | undefined;
  otherDocumentImageUrl_lte?: string | null | undefined;
  otherDocumentImageUrl_ne?: string | null | undefined;
  otherDocumentImageUrl_null?: boolean | null | undefined;
  otherDocumentImageUrl_prefix?: string | null | undefined;
  otherDocumentImageUrl_suffix?: string | null | undefined;
  paymentConfigs?: StorePaymentConfigFilterType | null | undefined;
  province?: string | null | undefined;
  province_gt?: string | null | undefined;
  province_gte?: string | null | undefined;
  province_in?: Array<string> | null | undefined;
  province_like?: string | null | undefined;
  province_lt?: string | null | undefined;
  province_lte?: string | null | undefined;
  province_ne?: string | null | undefined;
  province_null?: boolean | null | undefined;
  province_prefix?: string | null | undefined;
  province_suffix?: string | null | undefined;
  receiptFooter?: string | null | undefined;
  receiptFooter_gt?: string | null | undefined;
  receiptFooter_gte?: string | null | undefined;
  receiptFooter_in?: Array<string> | null | undefined;
  receiptFooter_like?: string | null | undefined;
  receiptFooter_lt?: string | null | undefined;
  receiptFooter_lte?: string | null | undefined;
  receiptFooter_ne?: string | null | undefined;
  receiptFooter_null?: boolean | null | undefined;
  receiptFooter_prefix?: string | null | undefined;
  receiptFooter_suffix?: string | null | undefined;
  rejectionReason?: string | null | undefined;
  rejectionReason_gt?: string | null | undefined;
  rejectionReason_gte?: string | null | undefined;
  rejectionReason_in?: Array<string> | null | undefined;
  rejectionReason_like?: string | null | undefined;
  rejectionReason_lt?: string | null | undefined;
  rejectionReason_lte?: string | null | undefined;
  rejectionReason_ne?: string | null | undefined;
  rejectionReason_null?: boolean | null | undefined;
  rejectionReason_prefix?: string | null | undefined;
  rejectionReason_suffix?: string | null | undefined;
  reviewedAt?: unknown;
  reviewedAt_gt?: unknown;
  reviewedAt_gte?: unknown;
  reviewedAt_in?: Array<unknown> | null | undefined;
  reviewedAt_lt?: unknown;
  reviewedAt_lte?: unknown;
  reviewedAt_ne?: unknown;
  reviewedAt_null?: boolean | null | undefined;
  reviewedByAccount?: AccountFilterType | null | undefined;
  reviewedByAccountId?: string | number | null | undefined;
  reviewedByAccountId_gt?: string | number | null | undefined;
  reviewedByAccountId_gte?: string | number | null | undefined;
  reviewedByAccountId_in?: Array<string | number> | null | undefined;
  reviewedByAccountId_lt?: string | number | null | undefined;
  reviewedByAccountId_lte?: string | number | null | undefined;
  reviewedByAccountId_ne?: string | number | null | undefined;
  reviewedByAccountId_null?: boolean | null | undefined;
  state?: number | null | undefined;
  state_gt?: number | null | undefined;
  state_gte?: number | null | undefined;
  state_in?: Array<number> | null | undefined;
  state_lt?: number | null | undefined;
  state_lte?: number | null | undefined;
  state_ne?: number | null | undefined;
  state_null?: boolean | null | undefined;
  storeArea?: number | null | undefined;
  storeArea_gt?: number | null | undefined;
  storeArea_gte?: number | null | undefined;
  storeArea_in?: Array<number> | null | undefined;
  storeArea_lt?: number | null | undefined;
  storeArea_lte?: number | null | undefined;
  storeArea_ne?: number | null | undefined;
  storeArea_null?: boolean | null | undefined;
  submittedAt?: unknown;
  submittedAt_gt?: unknown;
  submittedAt_gte?: unknown;
  submittedAt_in?: Array<unknown> | null | undefined;
  submittedAt_lt?: unknown;
  submittedAt_lte?: unknown;
  submittedAt_ne?: unknown;
  submittedAt_null?: boolean | null | undefined;
  supportDineIn?: boolean | null | undefined;
  supportDineIn_gt?: boolean | null | undefined;
  supportDineIn_gte?: boolean | null | undefined;
  supportDineIn_in?: Array<boolean> | null | undefined;
  supportDineIn_lt?: boolean | null | undefined;
  supportDineIn_lte?: boolean | null | undefined;
  supportDineIn_ne?: boolean | null | undefined;
  supportDineIn_null?: boolean | null | undefined;
  supportTakeout?: boolean | null | undefined;
  supportTakeout_gt?: boolean | null | undefined;
  supportTakeout_gte?: boolean | null | undefined;
  supportTakeout_in?: Array<boolean> | null | undefined;
  supportTakeout_lt?: boolean | null | undefined;
  supportTakeout_lte?: boolean | null | undefined;
  supportTakeout_ne?: boolean | null | undefined;
  supportTakeout_null?: boolean | null | undefined;
  tableCount?: number | null | undefined;
  tableCount_gt?: number | null | undefined;
  tableCount_gte?: number | null | undefined;
  tableCount_in?: Array<number> | null | undefined;
  tableCount_lt?: number | null | undefined;
  tableCount_lte?: number | null | undefined;
  tableCount_ne?: number | null | undefined;
  tableCount_null?: boolean | null | undefined;
  updatedAt?: number | null | undefined;
  updatedAt_gt?: number | null | undefined;
  updatedAt_gte?: number | null | undefined;
  updatedAt_in?: Array<number> | null | undefined;
  updatedAt_lt?: number | null | undefined;
  updatedAt_lte?: number | null | undefined;
  updatedAt_ne?: number | null | undefined;
  updatedAt_null?: boolean | null | undefined;
  updatedBy?: string | number | null | undefined;
  updatedBy_gt?: string | number | null | undefined;
  updatedBy_gte?: string | number | null | undefined;
  updatedBy_in?: Array<string | number> | null | undefined;
  updatedBy_lt?: string | number | null | undefined;
  updatedBy_lte?: string | number | null | undefined;
  updatedBy_ne?: string | number | null | undefined;
  updatedBy_null?: boolean | null | undefined;
  weight?: number | null | undefined;
  weight_gt?: number | null | undefined;
  weight_gte?: number | null | undefined;
  weight_in?: Array<number> | null | undefined;
  weight_lt?: number | null | undefined;
  weight_lte?: number | null | undefined;
  weight_ne?: number | null | undefined;
  weight_null?: boolean | null | undefined;
};

export type StoreLifecycle =
  | 'ACTIVE'
  | 'DRAFT'
  | 'PENDING_APPROVAL'
  | 'REJECTED';

export type StorePaymentConfigFilterType = {
  AND?: Array<StorePaymentConfigFilterType> | null | undefined;
  OR?: Array<StorePaymentConfigFilterType> | null | undefined;
  channel?: string | null | undefined;
  channel_gt?: string | null | undefined;
  channel_gte?: string | null | undefined;
  channel_in?: Array<string> | null | undefined;
  channel_like?: string | null | undefined;
  channel_lt?: string | null | undefined;
  channel_lte?: string | null | undefined;
  channel_ne?: string | null | undefined;
  channel_null?: boolean | null | undefined;
  channel_prefix?: string | null | undefined;
  channel_suffix?: string | null | undefined;
  configState?: string | null | undefined;
  configState_gt?: string | null | undefined;
  configState_gte?: string | null | undefined;
  configState_in?: Array<string> | null | undefined;
  configState_like?: string | null | undefined;
  configState_lt?: string | null | undefined;
  configState_lte?: string | null | undefined;
  configState_ne?: string | null | undefined;
  configState_null?: boolean | null | undefined;
  configState_prefix?: string | null | undefined;
  configState_suffix?: string | null | undefined;
  createdAt?: number | null | undefined;
  createdAt_gt?: number | null | undefined;
  createdAt_gte?: number | null | undefined;
  createdAt_in?: Array<number> | null | undefined;
  createdAt_lt?: number | null | undefined;
  createdAt_lte?: number | null | undefined;
  createdAt_ne?: number | null | undefined;
  createdAt_null?: boolean | null | undefined;
  createdBy?: string | number | null | undefined;
  createdBy_gt?: string | number | null | undefined;
  createdBy_gte?: string | number | null | undefined;
  createdBy_in?: Array<string | number> | null | undefined;
  createdBy_lt?: string | number | null | undefined;
  createdBy_lte?: string | number | null | undefined;
  createdBy_ne?: string | number | null | undefined;
  createdBy_null?: boolean | null | undefined;
  credentialCiphertext?: string | null | undefined;
  credentialCiphertext_gt?: string | null | undefined;
  credentialCiphertext_gte?: string | null | undefined;
  credentialCiphertext_in?: Array<string> | null | undefined;
  credentialCiphertext_like?: string | null | undefined;
  credentialCiphertext_lt?: string | null | undefined;
  credentialCiphertext_lte?: string | null | undefined;
  credentialCiphertext_ne?: string | null | undefined;
  credentialCiphertext_null?: boolean | null | undefined;
  credentialCiphertext_prefix?: string | null | undefined;
  credentialCiphertext_suffix?: string | null | undefined;
  deletedAt?: number | null | undefined;
  deletedAt_gt?: number | null | undefined;
  deletedAt_gte?: number | null | undefined;
  deletedAt_in?: Array<number> | null | undefined;
  deletedAt_lt?: number | null | undefined;
  deletedAt_lte?: number | null | undefined;
  deletedAt_ne?: number | null | undefined;
  deletedAt_null?: boolean | null | undefined;
  deletedBy?: string | number | null | undefined;
  deletedBy_gt?: string | number | null | undefined;
  deletedBy_gte?: string | number | null | undefined;
  deletedBy_in?: Array<string | number> | null | undefined;
  deletedBy_lt?: string | number | null | undefined;
  deletedBy_lte?: string | number | null | undefined;
  deletedBy_ne?: string | number | null | undefined;
  deletedBy_null?: boolean | null | undefined;
  environment?: string | null | undefined;
  environment_gt?: string | null | undefined;
  environment_gte?: string | null | undefined;
  environment_in?: Array<string> | null | undefined;
  environment_like?: string | null | undefined;
  environment_lt?: string | null | undefined;
  environment_lte?: string | null | undefined;
  environment_ne?: string | null | undefined;
  environment_null?: boolean | null | undefined;
  environment_prefix?: string | null | undefined;
  environment_suffix?: string | null | undefined;
  id?: string | number | null | undefined;
  id_gt?: string | number | null | undefined;
  id_gte?: string | number | null | undefined;
  id_in?: Array<string | number> | null | undefined;
  id_lt?: string | number | null | undefined;
  id_lte?: string | number | null | undefined;
  id_ne?: string | number | null | undefined;
  id_null?: boolean | null | undefined;
  isDelete?: number | null | undefined;
  isDelete_gt?: number | null | undefined;
  isDelete_gte?: number | null | undefined;
  isDelete_in?: Array<number> | null | undefined;
  isDelete_lt?: number | null | undefined;
  isDelete_lte?: number | null | undefined;
  isDelete_ne?: number | null | undefined;
  isDelete_null?: boolean | null | undefined;
  keyId?: string | null | undefined;
  keyId_gt?: string | null | undefined;
  keyId_gte?: string | null | undefined;
  keyId_in?: Array<string> | null | undefined;
  keyId_like?: string | null | undefined;
  keyId_lt?: string | null | undefined;
  keyId_lte?: string | null | undefined;
  keyId_ne?: string | null | undefined;
  keyId_null?: boolean | null | undefined;
  keyId_prefix?: string | null | undefined;
  keyId_suffix?: string | null | undefined;
  merchantId?: string | null | undefined;
  merchantId_gt?: string | null | undefined;
  merchantId_gte?: string | null | undefined;
  merchantId_in?: Array<string> | null | undefined;
  merchantId_like?: string | null | undefined;
  merchantId_lt?: string | null | undefined;
  merchantId_lte?: string | null | undefined;
  merchantId_ne?: string | null | undefined;
  merchantId_null?: boolean | null | undefined;
  merchantId_prefix?: string | null | undefined;
  merchantId_suffix?: string | null | undefined;
  ratePpm?: number | null | undefined;
  ratePpm_gt?: number | null | undefined;
  ratePpm_gte?: number | null | undefined;
  ratePpm_in?: Array<number> | null | undefined;
  ratePpm_lt?: number | null | undefined;
  ratePpm_lte?: number | null | undefined;
  ratePpm_ne?: number | null | undefined;
  ratePpm_null?: boolean | null | undefined;
  state?: number | null | undefined;
  state_gt?: number | null | undefined;
  state_gte?: number | null | undefined;
  state_in?: Array<number> | null | undefined;
  state_lt?: number | null | undefined;
  state_lte?: number | null | undefined;
  state_ne?: number | null | undefined;
  state_null?: boolean | null | undefined;
  store?: StoreFilterType | null | undefined;
  storeId?: string | number | null | undefined;
  storeId_gt?: string | number | null | undefined;
  storeId_gte?: string | number | null | undefined;
  storeId_in?: Array<string | number> | null | undefined;
  storeId_lt?: string | number | null | undefined;
  storeId_lte?: string | number | null | undefined;
  storeId_ne?: string | number | null | undefined;
  storeId_null?: boolean | null | undefined;
  updatedAt?: number | null | undefined;
  updatedAt_gt?: number | null | undefined;
  updatedAt_gte?: number | null | undefined;
  updatedAt_in?: Array<number> | null | undefined;
  updatedAt_lt?: number | null | undefined;
  updatedAt_lte?: number | null | undefined;
  updatedAt_ne?: number | null | undefined;
  updatedAt_null?: boolean | null | undefined;
  updatedBy?: string | number | null | undefined;
  updatedBy_gt?: string | number | null | undefined;
  updatedBy_gte?: string | number | null | undefined;
  updatedBy_in?: Array<string | number> | null | undefined;
  updatedBy_lt?: string | number | null | undefined;
  updatedBy_lte?: string | number | null | undefined;
  updatedBy_ne?: string | number | null | undefined;
  updatedBy_null?: boolean | null | undefined;
  version?: number | null | undefined;
  version_gt?: number | null | undefined;
  version_gte?: number | null | undefined;
  version_in?: Array<number> | null | undefined;
  version_lt?: number | null | undefined;
  version_lte?: number | null | undefined;
  version_ne?: number | null | undefined;
  version_null?: boolean | null | undefined;
  weight?: number | null | undefined;
  weight_gt?: number | null | undefined;
  weight_gte?: number | null | undefined;
  weight_in?: Array<number> | null | undefined;
  weight_lt?: number | null | undefined;
  weight_lte?: number | null | undefined;
  weight_ne?: number | null | undefined;
  weight_null?: boolean | null | undefined;
};

export type StorePaymentConfigRelationship = {
  channel?: string | null | undefined;
  configState?: string | null | undefined;
  credentialCiphertext?: string | null | undefined;
  environment?: string | null | undefined;
  id?: string | number | null | undefined;
  isDelete?: number | null | undefined;
  keyId?: string | null | undefined;
  merchantId?: string | null | undefined;
  ratePpm?: number | null | undefined;
  state?: number | null | undefined;
  storeId?: string | number | null | undefined;
  version?: number | null | undefined;
  weight?: number | null | undefined;
};

export type StoreRelationship = {
  address?: string | null | undefined;
  businessHours?: string | null | undefined;
  businessLicenseImageUrl?: string | null | undefined;
  businessStatus?: StoreBusinessStatus | null | undefined;
  city?: string | null | undefined;
  code?: string | null | undefined;
  contactPhone?: string | null | undefined;
  district?: string | null | undefined;
  id?: string | number | null | undefined;
  isDelete?: number | null | undefined;
  lifecycle?: StoreLifecycle | null | undefined;
  managerName?: string | null | undefined;
  managerPhone?: string | null | undefined;
  name?: string | null | undefined;
  organizationId?: string | number | null | undefined;
  otherDocumentImageUrl?: string | null | undefined;
  province?: string | null | undefined;
  receiptFooter?: string | null | undefined;
  rejectionReason?: string | null | undefined;
  reviewedAt?: unknown;
  reviewedByAccountId?: string | number | null | undefined;
  state?: number | null | undefined;
  storeArea?: number | null | undefined;
  submittedAt?: unknown;
  supportDineIn?: boolean | null | undefined;
  supportTakeout?: boolean | null | undefined;
  tableCount?: number | null | undefined;
  weight?: number | null | undefined;
};

export type SuspendOrganizationInput = {
  organizationId: string | number;
  reasonCode: string;
};

export type UpdateOperatorMembershipInput = {
  acceptedAt?: unknown;
  account?: AccountRelationship | null | undefined;
  accountId?: string | number | null | undefined;
  invitations?: Array<MembershipInvitationRelationship | null | undefined> | null | undefined;
  invitationsIds?: Array<string | number> | null | undefined;
  invitedAt?: unknown;
  isDelete?: number | null | undefined;
  organization?: OrganizationRelationship | null | undefined;
  organizationId?: string | number | null | undefined;
  roles?: Array<OperatorRoleRelationship | null | undefined> | null | undefined;
  rolesIds?: Array<string | number> | null | undefined;
  state?: number | null | undefined;
  status?: MembershipStatus | null | undefined;
  storeAccessMode?: StoreAccessMode | null | undefined;
  stores?: Array<StoreRelationship | null | undefined> | null | undefined;
  storesIds?: Array<string | number> | null | undefined;
  weight?: number | null | undefined;
};

export type UpdateOperatorRoleInput = {
  isDelete?: number | null | undefined;
  kind?: RoleKind | null | undefined;
  members?: Array<OperatorMembershipRelationship | null | undefined> | null | undefined;
  membersIds?: Array<string | number> | null | undefined;
  name?: string | null | undefined;
  organization?: OrganizationRelationship | null | undefined;
  organizationId?: string | number | null | undefined;
  permissions?: Array<PermissionRelationship | null | undefined> | null | undefined;
  permissionsIds?: Array<string | number> | null | undefined;
  state?: number | null | undefined;
  weight?: number | null | undefined;
};

export type UpdateStoreInput = {
  address?: string | null | undefined;
  auditLogs?: Array<AuditLogRelationship | null | undefined> | null | undefined;
  auditLogsIds?: Array<string | number> | null | undefined;
  businessHours?: string | null | undefined;
  businessLicenseImageUrl?: string | null | undefined;
  businessStatus?: StoreBusinessStatus | null | undefined;
  city?: string | null | undefined;
  code?: string | null | undefined;
  contactPhone?: string | null | undefined;
  district?: string | null | undefined;
  isDelete?: number | null | undefined;
  lifecycle?: StoreLifecycle | null | undefined;
  managerName?: string | null | undefined;
  managerPhone?: string | null | undefined;
  members?: Array<OperatorMembershipRelationship | null | undefined> | null | undefined;
  membersIds?: Array<string | number> | null | undefined;
  name?: string | null | undefined;
  organization?: OrganizationRelationship | null | undefined;
  organizationId?: string | number | null | undefined;
  otherDocumentImageUrl?: string | null | undefined;
  paymentConfigs?: Array<StorePaymentConfigRelationship | null | undefined> | null | undefined;
  paymentConfigsIds?: Array<string | number> | null | undefined;
  province?: string | null | undefined;
  receiptFooter?: string | null | undefined;
  rejectionReason?: string | null | undefined;
  reviewedAt?: unknown;
  reviewedByAccount?: AccountRelationship | null | undefined;
  reviewedByAccountId?: string | number | null | undefined;
  state?: number | null | undefined;
  storeArea?: number | null | undefined;
  submittedAt?: unknown;
  supportDineIn?: boolean | null | undefined;
  supportTakeout?: boolean | null | undefined;
  tableCount?: number | null | undefined;
  weight?: number | null | undefined;
};

export type WorkspaceType =
  | 'DISCOVERY'
  | 'FRANCHISE'
  | 'HEADQUARTERS';

export type AdminStoreEditorQueryVariables = Exact<{
  id: string | number;
}>;


export type AdminStoreEditorQuery = { store: { id: string, name: string, lifecycle: StoreLifecycle, rejectionReason: string | null, organizationId: string, contactPhone: string | null, managerName: string | null, managerPhone: string | null, province: string | null, city: string | null, district: string | null, address: string | null, businessHours: string | null, businessStatus: StoreBusinessStatus | null, supportDineIn: boolean | null, supportTakeout: boolean | null, storeArea: number | null, tableCount: number | null, receiptFooter: string | null, businessLicenseImageUrl: string | null, otherDocumentImageUrl: string | null, organization: { type: OrganizationType } } | null };

export type SetStoreDocumentMutationVariables = Exact<{
  storeId: string | number;
  kind: StoreDocumentKind;
  attachmentId: string | number;
}>;


export type SetStoreDocumentMutation = { setStoreDocument: { id: string, businessLicenseImageUrl: string | null, otherDocumentImageUrl: string | null } };

export type RemoveStoreDocumentMutationVariables = Exact<{
  storeId: string | number;
  kind: StoreDocumentKind;
}>;


export type RemoveStoreDocumentMutation = { removeStoreDocument: { id: string, businessLicenseImageUrl: string | null, otherDocumentImageUrl: string | null } };

export type ViewerFieldsFragment = { permissions: Array<string>, account: { id: string, phone: string, displayName: string, email: string | null, status: AccountStatus, mustChangePassword: boolean }, currentWorkspace: { workspaceType: WorkspaceType, organizationId: string | null, organizationName: string, homePath: string } | null, workspaces: Array<{ workspaceType: WorkspaceType, organizationId: string | null, organizationName: string, homePath: string }> } & { ' $fragmentName'?: 'ViewerFieldsFragment' };

export type LoginMutationVariables = Exact<{
  input: LoginInput;
}>;


export type LoginMutation = { login: { requiresPasswordChange: boolean, viewer: { ' $fragmentRefs'?: { 'ViewerFieldsFragment': ViewerFieldsFragment } } } };

export type ViewerQueryVariables = Exact<{ [key: string]: never; }>;


export type ViewerQuery = { viewer: { ' $fragmentRefs'?: { 'ViewerFieldsFragment': ViewerFieldsFragment } } };

export type WorkspacesQueryVariables = Exact<{ [key: string]: never; }>;


export type WorkspacesQuery = { workspaces: Array<{ workspaceType: WorkspaceType, organizationId: string | null, organizationName: string, homePath: string }> };

export type ChangeTemporaryPasswordMutationVariables = Exact<{
  input: ChangePasswordInput;
}>;


export type ChangeTemporaryPasswordMutation = { changeTemporaryPassword: { ' $fragmentRefs'?: { 'ViewerFieldsFragment': ViewerFieldsFragment } } };

export type SelectWorkspaceMutationVariables = Exact<{
  input: SelectWorkspaceInput;
}>;


export type SelectWorkspaceMutation = { selectWorkspace: { ' $fragmentRefs'?: { 'ViewerFieldsFragment': ViewerFieldsFragment } } };

export type LogoutMutationVariables = Exact<{ [key: string]: never; }>;


export type LogoutMutation = { logout: boolean };

export type SessionEventsSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type SessionEventsSubscription = { sessionEvents: { code: SessionEventCode, sessionId: string, organizationId: string | null, occurredAt: unknown } };

export type PendingMembershipInvitationsQueryVariables = Exact<{ [key: string]: never; }>;


export type PendingMembershipInvitationsQuery = { pendingMembershipInvitations: Array<{ id: string, membershipId: string, expiresAt: unknown }> };

export type AcceptMembershipInvitationMutationVariables = Exact<{
  id: string | number;
}>;


export type AcceptMembershipInvitationMutation = { acceptMembershipInvitation: { id: string, status: MembershipStatus, organizationId: string } };

export type FranchiseAuditLogsQueryVariables = Exact<{
  page: number;
  pageSize: number;
  q?: string | null | undefined;
  filter?: AuditLogFilterType | null | undefined;
}>;


export type FranchiseAuditLogsQuery = { auditLogs: { total: number, current_page: number, per_page: number, total_page: number, data: Array<{ id: string, action: string, resourceType: string, resourceId: string | null, resultCode: string, actorAccountId: string | null, storeId: string | null, createdAt: number }> } | null };

export type FranchiseRolesQueryVariables = Exact<{
  page: number;
  pageSize: number;
  q?: string | null | undefined;
  filter?: OperatorRoleFilterType | null | undefined;
}>;


export type FranchiseRolesQuery = { operatorRoles: { total: number, current_page: number, per_page: number, total_page: number, data: Array<{ id: string, name: string, kind: RoleKind, organizationId: string, permissions: Array<{ id: string, name: string, action: string, module: string, scope: PermissionScope }> }> } | null };

export type TenantPermissionsQueryVariables = Exact<{ [key: string]: never; }>;


export type TenantPermissionsQuery = { permissions: { total: number, data: Array<{ id: string, name: string, action: string, module: string, scope: PermissionScope }> } | null };

export type FranchiseCreateRoleMutationVariables = Exact<{
  input: CreateOperatorRoleInput;
}>;


export type FranchiseCreateRoleMutation = { createOperatorRole: { id: string, name: string, kind: RoleKind, organizationId: string } };

export type FranchiseUpdateRoleMutationVariables = Exact<{
  id: string | number;
  input: UpdateOperatorRoleInput;
}>;


export type FranchiseUpdateRoleMutation = { updateOperatorRole: { id: string, name: string, kind: RoleKind, organizationId: string } };

export type FranchiseDeleteRolesMutationVariables = Exact<{
  ids: Array<string | number> | string | number;
}>;


export type FranchiseDeleteRolesMutation = { deleteOperatorRoles: boolean };

export type FranchiseStaffQueryVariables = Exact<{
  page: number;
  pageSize: number;
  q?: string | null | undefined;
  filter?: OperatorMembershipFilterType | null | undefined;
}>;


export type FranchiseStaffQuery = { operatorMemberships: { total: number, current_page: number, per_page: number, total_page: number, data: Array<{ id: string, status: MembershipStatus, storeAccessMode: StoreAccessMode, accountId: string, organizationId: string, account: { id: string, phone: string, displayName: string, email: string | null }, roles: Array<{ id: string, name: string, kind: RoleKind }>, stores: Array<{ id: string, name: string, lifecycle: StoreLifecycle }> }> } | null };

export type FranchiseInviteStaffMutationVariables = Exact<{
  input: InviteOperatorInput;
}>;


export type FranchiseInviteStaffMutation = { inviteOperator: { temporaryPassword: string | null, invitationPending: boolean, membership: { id: string, status: MembershipStatus, accountId: string, organizationId: string } } };

export type FranchiseUpdateStaffMutationVariables = Exact<{
  id: string | number;
  input: UpdateOperatorMembershipInput;
}>;


export type FranchiseUpdateStaffMutation = { updateOperatorMembership: { id: string, status: MembershipStatus, storeAccessMode: StoreAccessMode, rolesIds: Array<string> | null, storesIds: Array<string> | null } };

export type FranchiseChangeMembershipStatusMutationVariables = Exact<{
  input: ChangeMembershipStatusInput;
}>;


export type FranchiseChangeMembershipStatusMutation = { changeMembershipStatus: { id: string, status: MembershipStatus } };

export type FranchiseStoresQueryVariables = Exact<{
  page: number;
  pageSize: number;
  q?: string | null | undefined;
  filter?: StoreFilterType | null | undefined;
}>;


export type FranchiseStoresQuery = { stores: { total: number, current_page: number, per_page: number, total_page: number, data: Array<{ id: string, code: string, name: string, lifecycle: StoreLifecycle, rejectionReason: string | null, organizationId: string, contactPhone: string | null, managerName: string | null, managerPhone: string | null, province: string | null, city: string | null, district: string | null, address: string | null, businessHours: string | null, businessStatus: StoreBusinessStatus | null, supportDineIn: boolean | null, supportTakeout: boolean | null, storeArea: number | null, tableCount: number | null, receiptFooter: string | null }> } | null };

export type FranchiseCreateStoreMutationVariables = Exact<{
  input: CreateStoreInput;
}>;


export type FranchiseCreateStoreMutation = { createStore: { id: string, code: string, name: string, lifecycle: StoreLifecycle, organizationId: string } };

export type FranchiseUpdateStoreMutationVariables = Exact<{
  id: string | number;
  input: UpdateStoreInput;
}>;


export type FranchiseUpdateStoreMutation = { updateStore: { id: string, code: string, name: string, lifecycle: StoreLifecycle, organizationId: string } };

export type FranchiseSubmitStoreMutationVariables = Exact<{
  id: string | number;
}>;


export type FranchiseSubmitStoreMutation = { submitStore: { id: string, lifecycle: StoreLifecycle, submittedAt: unknown } };

export type FranchiseDeleteStoresMutationVariables = Exact<{
  ids: Array<string | number> | string | number;
}>;


export type FranchiseDeleteStoresMutation = { deleteStores: boolean };

export type HqAdministratorsQueryVariables = Exact<{
  page: number;
  pageSize: number;
  q?: string | null | undefined;
  filter?: OperatorMembershipFilterType | null | undefined;
}>;


export type HqAdministratorsQuery = { operatorMemberships: { total: number, current_page: number, per_page: number, total_page: number, data: Array<{ id: string, status: MembershipStatus, storeAccessMode: StoreAccessMode, accountId: string, organizationId: string, updatedAt: number | null, account: { id: string, phone: string, displayName: string, email: string | null }, roles: Array<{ id: string, name: string, kind: RoleKind }> }> } | null };

export type HqAdministratorIdentityQueryVariables = Exact<{
  accountId: string | number;
}>;


export type HqAdministratorIdentityQuery = { operatorMemberships: { total: number, data: Array<{ id: string, accountId: string, roles: Array<{ id: string, kind: RoleKind }> }> } | null };

export type HqInviteAdministratorMutationVariables = Exact<{
  input: InviteOperatorInput;
}>;


export type HqInviteAdministratorMutation = { inviteOperator: { temporaryPassword: string | null, invitationPending: boolean, membership: { id: string, status: MembershipStatus, accountId: string, organizationId: string } } };

export type HqUpdateAdministratorMutationVariables = Exact<{
  id: string | number;
  input: UpdateOperatorMembershipInput;
}>;


export type HqUpdateAdministratorMutation = { updateOperatorMembership: { id: string, status: MembershipStatus, rolesIds: Array<string> | null } };

export type HqChangeAdministratorStatusMutationVariables = Exact<{
  input: ChangeMembershipStatusInput;
}>;


export type HqChangeAdministratorStatusMutation = { changeMembershipStatus: { id: string, status: MembershipStatus } };

export type HqDeleteAdministratorsMutationVariables = Exact<{
  ids: Array<string | number> | string | number;
}>;


export type HqDeleteAdministratorsMutation = { deleteOperatorMemberships: boolean };

export type HqResetAdministratorPasswordMutationVariables = Exact<{
  accountId: string | number;
}>;


export type HqResetAdministratorPasswordMutation = { resetTemporaryPassword: { accountId: string, temporaryPassword: string } };

export type HqAuditLogsQueryVariables = Exact<{
  page: number;
  pageSize: number;
  q?: string | null | undefined;
  filter?: AuditLogFilterType | null | undefined;
}>;


export type HqAuditLogsQuery = { auditLogs: { total: number, current_page: number, per_page: number, total_page: number, data: Array<{ id: string, action: string, resourceType: string, resourceId: string | null, resultCode: string, actorAccountId: string | null, organizationId: string | null, storeId: string | null, createdAt: number }> } | null };

export type HqDirectStoresQueryVariables = Exact<{
  page: number;
  pageSize: number;
  q?: string | null | undefined;
  filter: StoreFilterType;
}>;


export type HqDirectStoresQuery = { stores: { total: number, current_page: number, per_page: number, total_page: number, data: Array<{ id: string, code: string, name: string, lifecycle: StoreLifecycle, organizationId: string, contactPhone: string | null, managerName: string | null, managerPhone: string | null, province: string | null, city: string | null, district: string | null, address: string | null, businessHours: string | null, businessStatus: StoreBusinessStatus | null, supportDineIn: boolean | null, supportTakeout: boolean | null, storeArea: number | null, tableCount: number | null, receiptFooter: string | null }> } | null };

export type HqCreateDirectStoreMutationVariables = Exact<{
  input: CreateStoreInput;
}>;


export type HqCreateDirectStoreMutation = { createStore: { id: string, code: string, name: string, lifecycle: StoreLifecycle, organizationId: string, contactPhone: string | null, managerName: string | null, managerPhone: string | null, province: string | null, city: string | null, district: string | null, address: string | null, businessHours: string | null, businessStatus: StoreBusinessStatus | null, supportDineIn: boolean | null, supportTakeout: boolean | null, storeArea: number | null, tableCount: number | null, receiptFooter: string | null } };

export type HqUpdateDirectStoreMutationVariables = Exact<{
  id: string | number;
  input: UpdateStoreInput;
}>;


export type HqUpdateDirectStoreMutation = { updateStore: { id: string, code: string, name: string, lifecycle: StoreLifecycle, organizationId: string, contactPhone: string | null, managerName: string | null, managerPhone: string | null, province: string | null, city: string | null, district: string | null, address: string | null, businessHours: string | null, businessStatus: StoreBusinessStatus | null, supportDineIn: boolean | null, supportTakeout: boolean | null, storeArea: number | null, tableCount: number | null, receiptFooter: string | null } };

export type HqDeleteDirectStoresMutationVariables = Exact<{
  ids: Array<string | number> | string | number;
}>;


export type HqDeleteDirectStoresMutation = { deleteStores: boolean };

export type HqFranchiseStoresQueryVariables = Exact<{
  page: number;
  pageSize: number;
  q?: string | null | undefined;
  filter: StoreFilterType;
}>;


export type HqFranchiseStoresQuery = { stores: { total: number, current_page: number, per_page: number, total_page: number, data: Array<{ id: string, code: string, name: string, lifecycle: StoreLifecycle, businessStatus: StoreBusinessStatus | null, contactPhone: string | null, organizationId: string, organization: { id: string, name: string } }> } | null };

export type HqFranchisesQueryVariables = Exact<{
  page: number;
  pageSize: number;
  q?: string | null | undefined;
  filter: OrganizationFilterType;
  canResetPassword?: boolean;
}>;


export type HqFranchisesQuery = { organizations: { total: number, current_page: number, per_page: number, total_page: number, data: Array<{ id: string, code: string, name: string, status: OrganizationStatus, initialAccountId?: string | null, initialAccount?: { phone: string } | null }> } | null };

export type HqFranchiseInitialAccountCandidatesQueryVariables = Exact<{
  id: string | number;
}>;


export type HqFranchiseInitialAccountCandidatesQuery = { organization: { id: string, memberships: Array<{ id: string, status: MembershipStatus, account: { id: string, phone: string, displayName: string, status: AccountStatus } }> } | null };

export type HqResetFranchiseInitialPasswordMutationVariables = Exact<{
  organizationId: string | number;
}>;


export type HqResetFranchiseInitialPasswordMutation = { resetFranchiseInitialPassword: { accountId: string, temporaryPassword: string } };

export type HqProvisionFranchiseMutationVariables = Exact<{
  input: ProvisionFranchiseInput;
}>;


export type HqProvisionFranchiseMutation = { provisionFranchise: { temporaryPassword: string | null, invitationPending: boolean, organization: { id: string, code: string, name: string, status: OrganizationStatus }, membership: { id: string, status: MembershipStatus } } };

export type HqSuspendOrganizationMutationVariables = Exact<{
  input: SuspendOrganizationInput;
}>;


export type HqSuspendOrganizationMutation = { suspendOrganization: { id: string, status: OrganizationStatus, suspensionReasonCode: string | null } };

export type HqRestoreOrganizationMutationVariables = Exact<{
  id: string | number;
}>;


export type HqRestoreOrganizationMutation = { restoreOrganization: { id: string, status: OrganizationStatus } };

export type HqStoresQueryVariables = Exact<{
  page: number;
  pageSize: number;
  q?: string | null | undefined;
  filter: StoreFilterType;
}>;


export type HqStoresQuery = { stores: { total: number, current_page: number, per_page: number, total_page: number, data: Array<{ id: string, code: string, name: string, lifecycle: StoreLifecycle, organizationId: string, contactPhone: string | null, managerName: string | null, managerPhone: string | null, province: string | null, city: string | null, district: string | null, address: string | null, businessHours: string | null, businessStatus: StoreBusinessStatus | null, supportDineIn: boolean | null, supportTakeout: boolean | null, storeArea: number | null, tableCount: number | null, receiptFooter: string | null, organization: { id: string, name: string, type: OrganizationType } }> } | null };

export type HqStoresTotalQueryVariables = Exact<{ [key: string]: never; }>;


export type HqStoresTotalQuery = { stores: { total: number } | null };

export type HqRolesQueryVariables = Exact<{
  page: number;
  pageSize: number;
  q?: string | null | undefined;
  filter: OperatorRoleFilterType;
}>;


export type HqRolesQuery = { operatorRoles: { total: number, current_page: number, per_page: number, total_page: number, data: Array<{ id: string, name: string, kind: RoleKind, organizationId: string, permissions: Array<{ id: string, name: string, action: string, module: string, scope: PermissionScope }> }> } | null };

export type SystemPermissionsQueryVariables = Exact<{ [key: string]: never; }>;


export type SystemPermissionsQuery = { permissions: { total: number, data: Array<{ id: string, name: string, action: string, module: string, scope: PermissionScope }> } | null };

export type HqCreateRoleMutationVariables = Exact<{
  input: CreateOperatorRoleInput;
}>;


export type HqCreateRoleMutation = { createOperatorRole: { id: string, name: string, kind: RoleKind, organizationId: string } };

export type HqUpdateRoleMutationVariables = Exact<{
  id: string | number;
  input: UpdateOperatorRoleInput;
}>;


export type HqUpdateRoleMutation = { updateOperatorRole: { id: string, name: string, kind: RoleKind, organizationId: string } };

export type HqDeleteRolesMutationVariables = Exact<{
  ids: Array<string | number> | string | number;
}>;


export type HqDeleteRolesMutation = { deleteOperatorRoles: boolean };

export type HqStoreApprovalsQueryVariables = Exact<{
  page: number;
  pageSize: number;
  q?: string | null | undefined;
  filter: StoreFilterType;
}>;


export type HqStoreApprovalsQuery = { stores: { total: number, current_page: number, per_page: number, total_page: number, data: Array<{ id: string, code: string, name: string, lifecycle: StoreLifecycle, submittedAt: unknown, organizationId: string, contactPhone: string | null, managerName: string | null, managerPhone: string | null, province: string | null, city: string | null, district: string | null, address: string | null, businessHours: string | null, businessStatus: StoreBusinessStatus | null, supportDineIn: boolean | null, supportTakeout: boolean | null, storeArea: number | null, tableCount: number | null, organization: { id: string, name: string } }> } | null };

export type HqReviewStoreMutationVariables = Exact<{
  input: ReviewStoreInput;
}>;


export type HqReviewStoreMutation = { reviewStore: { id: string, lifecycle: StoreLifecycle, rejectionReason: string | null, reviewedAt: unknown } };

export const ViewerFieldsFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ViewerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Viewer"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"account"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"mustChangePassword"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentWorkspace"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"workspaceType"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationName"}},{"kind":"Field","name":{"kind":"Name","value":"homePath"}}]}},{"kind":"Field","name":{"kind":"Name","value":"workspaces"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"workspaceType"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationName"}},{"kind":"Field","name":{"kind":"Name","value":"homePath"}}]}},{"kind":"Field","name":{"kind":"Name","value":"permissions"}}]}}]} as unknown as DocumentNode<ViewerFieldsFragment, unknown>;
export const AdminStoreEditorDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"AdminStoreEditor"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"store"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}},{"kind":"Field","name":{"kind":"Name","value":"rejectionReason"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organization"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"type"}}]}},{"kind":"Field","name":{"kind":"Name","value":"contactPhone"}},{"kind":"Field","name":{"kind":"Name","value":"managerName"}},{"kind":"Field","name":{"kind":"Name","value":"managerPhone"}},{"kind":"Field","name":{"kind":"Name","value":"province"}},{"kind":"Field","name":{"kind":"Name","value":"city"}},{"kind":"Field","name":{"kind":"Name","value":"district"}},{"kind":"Field","name":{"kind":"Name","value":"address"}},{"kind":"Field","name":{"kind":"Name","value":"businessHours"}},{"kind":"Field","name":{"kind":"Name","value":"businessStatus"}},{"kind":"Field","name":{"kind":"Name","value":"supportDineIn"}},{"kind":"Field","name":{"kind":"Name","value":"supportTakeout"}},{"kind":"Field","name":{"kind":"Name","value":"storeArea"}},{"kind":"Field","name":{"kind":"Name","value":"tableCount"}},{"kind":"Field","name":{"kind":"Name","value":"receiptFooter"}},{"kind":"Field","name":{"kind":"Name","value":"businessLicenseImageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"otherDocumentImageUrl"}}]}}]}}]} as unknown as DocumentNode<AdminStoreEditorQuery, AdminStoreEditorQueryVariables>;
export const SetStoreDocumentDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"SetStoreDocument"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"storeId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"kind"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"StoreDocumentKind"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"attachmentId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"setStoreDocument"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"storeId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"storeId"}}},{"kind":"Argument","name":{"kind":"Name","value":"kind"},"value":{"kind":"Variable","name":{"kind":"Name","value":"kind"}}},{"kind":"Argument","name":{"kind":"Name","value":"attachmentId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"attachmentId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"businessLicenseImageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"otherDocumentImageUrl"}}]}}]}}]} as unknown as DocumentNode<SetStoreDocumentMutation, SetStoreDocumentMutationVariables>;
export const RemoveStoreDocumentDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"RemoveStoreDocument"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"storeId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"kind"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"StoreDocumentKind"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"removeStoreDocument"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"storeId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"storeId"}}},{"kind":"Argument","name":{"kind":"Name","value":"kind"},"value":{"kind":"Variable","name":{"kind":"Name","value":"kind"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"businessLicenseImageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"otherDocumentImageUrl"}}]}}]}}]} as unknown as DocumentNode<RemoveStoreDocumentMutation, RemoveStoreDocumentMutationVariables>;
export const LoginDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Login"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"LoginInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"login"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"requiresPasswordChange"}},{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"ViewerFields"}}]}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ViewerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Viewer"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"account"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"mustChangePassword"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentWorkspace"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"workspaceType"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationName"}},{"kind":"Field","name":{"kind":"Name","value":"homePath"}}]}},{"kind":"Field","name":{"kind":"Name","value":"workspaces"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"workspaceType"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationName"}},{"kind":"Field","name":{"kind":"Name","value":"homePath"}}]}},{"kind":"Field","name":{"kind":"Name","value":"permissions"}}]}}]} as unknown as DocumentNode<LoginMutation, LoginMutationVariables>;
export const ViewerDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"ViewerFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ViewerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Viewer"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"account"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"mustChangePassword"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentWorkspace"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"workspaceType"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationName"}},{"kind":"Field","name":{"kind":"Name","value":"homePath"}}]}},{"kind":"Field","name":{"kind":"Name","value":"workspaces"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"workspaceType"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationName"}},{"kind":"Field","name":{"kind":"Name","value":"homePath"}}]}},{"kind":"Field","name":{"kind":"Name","value":"permissions"}}]}}]} as unknown as DocumentNode<ViewerQuery, ViewerQueryVariables>;
export const WorkspacesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Workspaces"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"workspaces"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"workspaceType"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationName"}},{"kind":"Field","name":{"kind":"Name","value":"homePath"}}]}}]}}]} as unknown as DocumentNode<WorkspacesQuery, WorkspacesQueryVariables>;
export const ChangeTemporaryPasswordDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"ChangeTemporaryPassword"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ChangePasswordInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"changeTemporaryPassword"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"ViewerFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ViewerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Viewer"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"account"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"mustChangePassword"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentWorkspace"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"workspaceType"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationName"}},{"kind":"Field","name":{"kind":"Name","value":"homePath"}}]}},{"kind":"Field","name":{"kind":"Name","value":"workspaces"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"workspaceType"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationName"}},{"kind":"Field","name":{"kind":"Name","value":"homePath"}}]}},{"kind":"Field","name":{"kind":"Name","value":"permissions"}}]}}]} as unknown as DocumentNode<ChangeTemporaryPasswordMutation, ChangeTemporaryPasswordMutationVariables>;
export const SelectWorkspaceDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"SelectWorkspace"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"SelectWorkspaceInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"selectWorkspace"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"ViewerFields"}}]}}]}},{"kind":"FragmentDefinition","name":{"kind":"Name","value":"ViewerFields"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"Viewer"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"account"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"mustChangePassword"}}]}},{"kind":"Field","name":{"kind":"Name","value":"currentWorkspace"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"workspaceType"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationName"}},{"kind":"Field","name":{"kind":"Name","value":"homePath"}}]}},{"kind":"Field","name":{"kind":"Name","value":"workspaces"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"workspaceType"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationName"}},{"kind":"Field","name":{"kind":"Name","value":"homePath"}}]}},{"kind":"Field","name":{"kind":"Name","value":"permissions"}}]}}]} as unknown as DocumentNode<SelectWorkspaceMutation, SelectWorkspaceMutationVariables>;
export const LogoutDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Logout"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"logout"}}]}}]} as unknown as DocumentNode<LogoutMutation, LogoutMutationVariables>;
export const SessionEventsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"subscription","name":{"kind":"Name","value":"SessionEvents"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"sessionEvents"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"sessionId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"occurredAt"}}]}}]}}]} as unknown as DocumentNode<SessionEventsSubscription, SessionEventsSubscriptionVariables>;
export const PendingMembershipInvitationsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"PendingMembershipInvitations"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"pendingMembershipInvitations"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"membershipId"}},{"kind":"Field","name":{"kind":"Name","value":"expiresAt"}}]}}]}}]} as unknown as DocumentNode<PendingMembershipInvitationsQuery, PendingMembershipInvitationsQueryVariables>;
export const AcceptMembershipInvitationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"AcceptMembershipInvitation"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"acceptMembershipInvitation"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}}]}}]}}]} as unknown as DocumentNode<AcceptMembershipInvitationMutation, AcceptMembershipInvitationMutationVariables>;
export const FranchiseAuditLogsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"FranchiseAuditLogs"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"q"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"AuditLogFilterType"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"auditLogs"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}},{"kind":"Argument","name":{"kind":"Name","value":"q"},"value":{"kind":"Variable","name":{"kind":"Name","value":"q"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"action"}},{"kind":"Field","name":{"kind":"Name","value":"resourceType"}},{"kind":"Field","name":{"kind":"Name","value":"resourceId"}},{"kind":"Field","name":{"kind":"Name","value":"resultCode"}},{"kind":"Field","name":{"kind":"Name","value":"actorAccountId"}},{"kind":"Field","name":{"kind":"Name","value":"storeId"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"current_page"}},{"kind":"Field","name":{"kind":"Name","value":"per_page"}},{"kind":"Field","name":{"kind":"Name","value":"total_page"}}]}}]}}]} as unknown as DocumentNode<FranchiseAuditLogsQuery, FranchiseAuditLogsQueryVariables>;
export const FranchiseRolesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"FranchiseRoles"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"q"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"OperatorRoleFilterType"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"operatorRoles"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}},{"kind":"Argument","name":{"kind":"Name","value":"q"},"value":{"kind":"Variable","name":{"kind":"Name","value":"q"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"kind"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"permissions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"action"}},{"kind":"Field","name":{"kind":"Name","value":"module"}},{"kind":"Field","name":{"kind":"Name","value":"scope"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"current_page"}},{"kind":"Field","name":{"kind":"Name","value":"per_page"}},{"kind":"Field","name":{"kind":"Name","value":"total_page"}}]}}]}}]} as unknown as DocumentNode<FranchiseRolesQuery, FranchiseRolesQueryVariables>;
export const TenantPermissionsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"TenantPermissions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"permissions"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"IntValue","value":"1"}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"IntValue","value":"200"}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"scope"},"value":{"kind":"EnumValue","value":"TENANT"}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"action"}},{"kind":"Field","name":{"kind":"Name","value":"module"}},{"kind":"Field","name":{"kind":"Name","value":"scope"}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}}]}}]}}]} as unknown as DocumentNode<TenantPermissionsQuery, TenantPermissionsQueryVariables>;
export const FranchiseCreateRoleDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"FranchiseCreateRole"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateOperatorRoleInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createOperatorRole"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"kind"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}}]}}]}}]} as unknown as DocumentNode<FranchiseCreateRoleMutation, FranchiseCreateRoleMutationVariables>;
export const FranchiseUpdateRoleDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"FranchiseUpdateRole"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateOperatorRoleInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateOperatorRole"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"kind"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}}]}}]}}]} as unknown as DocumentNode<FranchiseUpdateRoleMutation, FranchiseUpdateRoleMutationVariables>;
export const FranchiseDeleteRolesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"FranchiseDeleteRoles"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"ids"}},"type":{"kind":"NonNullType","type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteOperatorRoles"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"ids"}}}]}]}}]} as unknown as DocumentNode<FranchiseDeleteRolesMutation, FranchiseDeleteRolesMutationVariables>;
export const FranchiseStaffDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"FranchiseStaff"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"q"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"OperatorMembershipFilterType"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"operatorMemberships"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}},{"kind":"Argument","name":{"kind":"Name","value":"q"},"value":{"kind":"Variable","name":{"kind":"Name","value":"q"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"storeAccessMode"}},{"kind":"Field","name":{"kind":"Name","value":"accountId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"account"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"email"}}]}},{"kind":"Field","name":{"kind":"Name","value":"roles"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"kind"}}]}},{"kind":"Field","name":{"kind":"Name","value":"stores"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"current_page"}},{"kind":"Field","name":{"kind":"Name","value":"per_page"}},{"kind":"Field","name":{"kind":"Name","value":"total_page"}}]}}]}}]} as unknown as DocumentNode<FranchiseStaffQuery, FranchiseStaffQueryVariables>;
export const FranchiseInviteStaffDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"FranchiseInviteStaff"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"InviteOperatorInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"inviteOperator"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"membership"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"accountId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}}]}},{"kind":"Field","name":{"kind":"Name","value":"temporaryPassword"}},{"kind":"Field","name":{"kind":"Name","value":"invitationPending"}}]}}]}}]} as unknown as DocumentNode<FranchiseInviteStaffMutation, FranchiseInviteStaffMutationVariables>;
export const FranchiseUpdateStaffDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"FranchiseUpdateStaff"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateOperatorMembershipInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateOperatorMembership"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"storeAccessMode"}},{"kind":"Field","name":{"kind":"Name","value":"rolesIds"}},{"kind":"Field","name":{"kind":"Name","value":"storesIds"}}]}}]}}]} as unknown as DocumentNode<FranchiseUpdateStaffMutation, FranchiseUpdateStaffMutationVariables>;
export const FranchiseChangeMembershipStatusDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"FranchiseChangeMembershipStatus"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ChangeMembershipStatusInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"changeMembershipStatus"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}}]}}]} as unknown as DocumentNode<FranchiseChangeMembershipStatusMutation, FranchiseChangeMembershipStatusMutationVariables>;
export const FranchiseStoresDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"FranchiseStores"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"q"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"StoreFilterType"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"stores"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}},{"kind":"Argument","name":{"kind":"Name","value":"q"},"value":{"kind":"Variable","name":{"kind":"Name","value":"q"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}},{"kind":"Field","name":{"kind":"Name","value":"rejectionReason"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"contactPhone"}},{"kind":"Field","name":{"kind":"Name","value":"managerName"}},{"kind":"Field","name":{"kind":"Name","value":"managerPhone"}},{"kind":"Field","name":{"kind":"Name","value":"province"}},{"kind":"Field","name":{"kind":"Name","value":"city"}},{"kind":"Field","name":{"kind":"Name","value":"district"}},{"kind":"Field","name":{"kind":"Name","value":"address"}},{"kind":"Field","name":{"kind":"Name","value":"businessHours"}},{"kind":"Field","name":{"kind":"Name","value":"businessStatus"}},{"kind":"Field","name":{"kind":"Name","value":"supportDineIn"}},{"kind":"Field","name":{"kind":"Name","value":"supportTakeout"}},{"kind":"Field","name":{"kind":"Name","value":"storeArea"}},{"kind":"Field","name":{"kind":"Name","value":"tableCount"}},{"kind":"Field","name":{"kind":"Name","value":"receiptFooter"}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"current_page"}},{"kind":"Field","name":{"kind":"Name","value":"per_page"}},{"kind":"Field","name":{"kind":"Name","value":"total_page"}}]}}]}}]} as unknown as DocumentNode<FranchiseStoresQuery, FranchiseStoresQueryVariables>;
export const FranchiseCreateStoreDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"FranchiseCreateStore"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateStoreInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createStore"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}}]}}]}}]} as unknown as DocumentNode<FranchiseCreateStoreMutation, FranchiseCreateStoreMutationVariables>;
export const FranchiseUpdateStoreDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"FranchiseUpdateStore"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateStoreInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateStore"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}}]}}]}}]} as unknown as DocumentNode<FranchiseUpdateStoreMutation, FranchiseUpdateStoreMutationVariables>;
export const FranchiseSubmitStoreDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"FranchiseSubmitStore"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"submitStore"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}},{"kind":"Field","name":{"kind":"Name","value":"submittedAt"}}]}}]}}]} as unknown as DocumentNode<FranchiseSubmitStoreMutation, FranchiseSubmitStoreMutationVariables>;
export const FranchiseDeleteStoresDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"FranchiseDeleteStores"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"ids"}},"type":{"kind":"NonNullType","type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteStores"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"ids"}}}]}]}}]} as unknown as DocumentNode<FranchiseDeleteStoresMutation, FranchiseDeleteStoresMutationVariables>;
export const HqAdministratorsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"HqAdministrators"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"q"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"OperatorMembershipFilterType"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"operatorMemberships"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}},{"kind":"Argument","name":{"kind":"Name","value":"q"},"value":{"kind":"Variable","name":{"kind":"Name","value":"q"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"storeAccessMode"}},{"kind":"Field","name":{"kind":"Name","value":"accountId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}},{"kind":"Field","name":{"kind":"Name","value":"account"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"email"}}]}},{"kind":"Field","name":{"kind":"Name","value":"roles"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"kind"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"current_page"}},{"kind":"Field","name":{"kind":"Name","value":"per_page"}},{"kind":"Field","name":{"kind":"Name","value":"total_page"}}]}}]}}]} as unknown as DocumentNode<HqAdministratorsQuery, HqAdministratorsQueryVariables>;
export const HqAdministratorIdentityDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"HqAdministratorIdentity"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"accountId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"operatorMemberships"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"IntValue","value":"1"}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"IntValue","value":"1"}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"accountId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"accountId"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"accountId"}},{"kind":"Field","name":{"kind":"Name","value":"roles"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"kind"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}}]}}]}}]} as unknown as DocumentNode<HqAdministratorIdentityQuery, HqAdministratorIdentityQueryVariables>;
export const HqInviteAdministratorDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqInviteAdministrator"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"InviteOperatorInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"inviteOperator"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"membership"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"accountId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}}]}},{"kind":"Field","name":{"kind":"Name","value":"temporaryPassword"}},{"kind":"Field","name":{"kind":"Name","value":"invitationPending"}}]}}]}}]} as unknown as DocumentNode<HqInviteAdministratorMutation, HqInviteAdministratorMutationVariables>;
export const HqUpdateAdministratorDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqUpdateAdministrator"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateOperatorMembershipInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateOperatorMembership"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"rolesIds"}}]}}]}}]} as unknown as DocumentNode<HqUpdateAdministratorMutation, HqUpdateAdministratorMutationVariables>;
export const HqChangeAdministratorStatusDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqChangeAdministratorStatus"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ChangeMembershipStatusInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"changeMembershipStatus"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}}]}}]} as unknown as DocumentNode<HqChangeAdministratorStatusMutation, HqChangeAdministratorStatusMutationVariables>;
export const HqDeleteAdministratorsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqDeleteAdministrators"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"ids"}},"type":{"kind":"NonNullType","type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteOperatorMemberships"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"ids"}}}]}]}}]} as unknown as DocumentNode<HqDeleteAdministratorsMutation, HqDeleteAdministratorsMutationVariables>;
export const HqResetAdministratorPasswordDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqResetAdministratorPassword"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"accountId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"resetTemporaryPassword"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"accountId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"accountId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"accountId"}},{"kind":"Field","name":{"kind":"Name","value":"temporaryPassword"}}]}}]}}]} as unknown as DocumentNode<HqResetAdministratorPasswordMutation, HqResetAdministratorPasswordMutationVariables>;
export const HqAuditLogsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"HqAuditLogs"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"q"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"AuditLogFilterType"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"auditLogs"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}},{"kind":"Argument","name":{"kind":"Name","value":"q"},"value":{"kind":"Variable","name":{"kind":"Name","value":"q"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"action"}},{"kind":"Field","name":{"kind":"Name","value":"resourceType"}},{"kind":"Field","name":{"kind":"Name","value":"resourceId"}},{"kind":"Field","name":{"kind":"Name","value":"resultCode"}},{"kind":"Field","name":{"kind":"Name","value":"actorAccountId"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"storeId"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"current_page"}},{"kind":"Field","name":{"kind":"Name","value":"per_page"}},{"kind":"Field","name":{"kind":"Name","value":"total_page"}}]}}]}}]} as unknown as DocumentNode<HqAuditLogsQuery, HqAuditLogsQueryVariables>;
export const HqDirectStoresDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"HqDirectStores"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"q"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"StoreFilterType"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"stores"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}},{"kind":"Argument","name":{"kind":"Name","value":"q"},"value":{"kind":"Variable","name":{"kind":"Name","value":"q"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"contactPhone"}},{"kind":"Field","name":{"kind":"Name","value":"managerName"}},{"kind":"Field","name":{"kind":"Name","value":"managerPhone"}},{"kind":"Field","name":{"kind":"Name","value":"province"}},{"kind":"Field","name":{"kind":"Name","value":"city"}},{"kind":"Field","name":{"kind":"Name","value":"district"}},{"kind":"Field","name":{"kind":"Name","value":"address"}},{"kind":"Field","name":{"kind":"Name","value":"businessHours"}},{"kind":"Field","name":{"kind":"Name","value":"businessStatus"}},{"kind":"Field","name":{"kind":"Name","value":"supportDineIn"}},{"kind":"Field","name":{"kind":"Name","value":"supportTakeout"}},{"kind":"Field","name":{"kind":"Name","value":"storeArea"}},{"kind":"Field","name":{"kind":"Name","value":"tableCount"}},{"kind":"Field","name":{"kind":"Name","value":"receiptFooter"}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"current_page"}},{"kind":"Field","name":{"kind":"Name","value":"per_page"}},{"kind":"Field","name":{"kind":"Name","value":"total_page"}}]}}]}}]} as unknown as DocumentNode<HqDirectStoresQuery, HqDirectStoresQueryVariables>;
export const HqCreateDirectStoreDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqCreateDirectStore"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateStoreInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createStore"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"contactPhone"}},{"kind":"Field","name":{"kind":"Name","value":"managerName"}},{"kind":"Field","name":{"kind":"Name","value":"managerPhone"}},{"kind":"Field","name":{"kind":"Name","value":"province"}},{"kind":"Field","name":{"kind":"Name","value":"city"}},{"kind":"Field","name":{"kind":"Name","value":"district"}},{"kind":"Field","name":{"kind":"Name","value":"address"}},{"kind":"Field","name":{"kind":"Name","value":"businessHours"}},{"kind":"Field","name":{"kind":"Name","value":"businessStatus"}},{"kind":"Field","name":{"kind":"Name","value":"supportDineIn"}},{"kind":"Field","name":{"kind":"Name","value":"supportTakeout"}},{"kind":"Field","name":{"kind":"Name","value":"storeArea"}},{"kind":"Field","name":{"kind":"Name","value":"tableCount"}},{"kind":"Field","name":{"kind":"Name","value":"receiptFooter"}}]}}]}}]} as unknown as DocumentNode<HqCreateDirectStoreMutation, HqCreateDirectStoreMutationVariables>;
export const HqUpdateDirectStoreDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqUpdateDirectStore"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateStoreInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateStore"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"contactPhone"}},{"kind":"Field","name":{"kind":"Name","value":"managerName"}},{"kind":"Field","name":{"kind":"Name","value":"managerPhone"}},{"kind":"Field","name":{"kind":"Name","value":"province"}},{"kind":"Field","name":{"kind":"Name","value":"city"}},{"kind":"Field","name":{"kind":"Name","value":"district"}},{"kind":"Field","name":{"kind":"Name","value":"address"}},{"kind":"Field","name":{"kind":"Name","value":"businessHours"}},{"kind":"Field","name":{"kind":"Name","value":"businessStatus"}},{"kind":"Field","name":{"kind":"Name","value":"supportDineIn"}},{"kind":"Field","name":{"kind":"Name","value":"supportTakeout"}},{"kind":"Field","name":{"kind":"Name","value":"storeArea"}},{"kind":"Field","name":{"kind":"Name","value":"tableCount"}},{"kind":"Field","name":{"kind":"Name","value":"receiptFooter"}}]}}]}}]} as unknown as DocumentNode<HqUpdateDirectStoreMutation, HqUpdateDirectStoreMutationVariables>;
export const HqDeleteDirectStoresDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqDeleteDirectStores"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"ids"}},"type":{"kind":"NonNullType","type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteStores"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"ids"}}}]}]}}]} as unknown as DocumentNode<HqDeleteDirectStoresMutation, HqDeleteDirectStoresMutationVariables>;
export const HqFranchiseStoresDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"HqFranchiseStores"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"q"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"StoreFilterType"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"stores"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}},{"kind":"Argument","name":{"kind":"Name","value":"q"},"value":{"kind":"Variable","name":{"kind":"Name","value":"q"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}},{"kind":"Field","name":{"kind":"Name","value":"businessStatus"}},{"kind":"Field","name":{"kind":"Name","value":"contactPhone"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organization"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"current_page"}},{"kind":"Field","name":{"kind":"Name","value":"per_page"}},{"kind":"Field","name":{"kind":"Name","value":"total_page"}}]}}]}}]} as unknown as DocumentNode<HqFranchiseStoresQuery, HqFranchiseStoresQueryVariables>;
export const HqFranchisesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"HqFranchises"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"q"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"OrganizationFilterType"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"canResetPassword"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Boolean"}}},"defaultValue":{"kind":"BooleanValue","value":false}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"organizations"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}},{"kind":"Argument","name":{"kind":"Name","value":"q"},"value":{"kind":"Variable","name":{"kind":"Name","value":"q"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"initialAccountId"},"directives":[{"kind":"Directive","name":{"kind":"Name","value":"include"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"if"},"value":{"kind":"Variable","name":{"kind":"Name","value":"canResetPassword"}}}]}]},{"kind":"Field","name":{"kind":"Name","value":"initialAccount"},"directives":[{"kind":"Directive","name":{"kind":"Name","value":"include"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"if"},"value":{"kind":"Variable","name":{"kind":"Name","value":"canResetPassword"}}}]}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"phone"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"current_page"}},{"kind":"Field","name":{"kind":"Name","value":"per_page"}},{"kind":"Field","name":{"kind":"Name","value":"total_page"}}]}}]}}]} as unknown as DocumentNode<HqFranchisesQuery, HqFranchisesQueryVariables>;
export const HqFranchiseInitialAccountCandidatesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"HqFranchiseInitialAccountCandidates"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"organization"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"memberships"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"account"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"phone"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}}]}}]}}]}}]} as unknown as DocumentNode<HqFranchiseInitialAccountCandidatesQuery, HqFranchiseInitialAccountCandidatesQueryVariables>;
export const HqResetFranchiseInitialPasswordDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqResetFranchiseInitialPassword"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"organizationId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"resetFranchiseInitialPassword"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"organizationId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"organizationId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"accountId"}},{"kind":"Field","name":{"kind":"Name","value":"temporaryPassword"}}]}}]}}]} as unknown as DocumentNode<HqResetFranchiseInitialPasswordMutation, HqResetFranchiseInitialPasswordMutationVariables>;
export const HqProvisionFranchiseDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqProvisionFranchise"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ProvisionFranchiseInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"provisionFranchise"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"organization"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}},{"kind":"Field","name":{"kind":"Name","value":"membership"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}},{"kind":"Field","name":{"kind":"Name","value":"temporaryPassword"}},{"kind":"Field","name":{"kind":"Name","value":"invitationPending"}}]}}]}}]} as unknown as DocumentNode<HqProvisionFranchiseMutation, HqProvisionFranchiseMutationVariables>;
export const HqSuspendOrganizationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqSuspendOrganization"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"SuspendOrganizationInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"suspendOrganization"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}},{"kind":"Field","name":{"kind":"Name","value":"suspensionReasonCode"}}]}}]}}]} as unknown as DocumentNode<HqSuspendOrganizationMutation, HqSuspendOrganizationMutationVariables>;
export const HqRestoreOrganizationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqRestoreOrganization"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"restoreOrganization"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"status"}}]}}]}}]} as unknown as DocumentNode<HqRestoreOrganizationMutation, HqRestoreOrganizationMutationVariables>;
export const HqStoresDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"HqStores"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"q"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"StoreFilterType"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"stores"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}},{"kind":"Argument","name":{"kind":"Name","value":"q"},"value":{"kind":"Variable","name":{"kind":"Name","value":"q"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organization"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"type"}}]}},{"kind":"Field","name":{"kind":"Name","value":"contactPhone"}},{"kind":"Field","name":{"kind":"Name","value":"managerName"}},{"kind":"Field","name":{"kind":"Name","value":"managerPhone"}},{"kind":"Field","name":{"kind":"Name","value":"province"}},{"kind":"Field","name":{"kind":"Name","value":"city"}},{"kind":"Field","name":{"kind":"Name","value":"district"}},{"kind":"Field","name":{"kind":"Name","value":"address"}},{"kind":"Field","name":{"kind":"Name","value":"businessHours"}},{"kind":"Field","name":{"kind":"Name","value":"businessStatus"}},{"kind":"Field","name":{"kind":"Name","value":"supportDineIn"}},{"kind":"Field","name":{"kind":"Name","value":"supportTakeout"}},{"kind":"Field","name":{"kind":"Name","value":"storeArea"}},{"kind":"Field","name":{"kind":"Name","value":"tableCount"}},{"kind":"Field","name":{"kind":"Name","value":"receiptFooter"}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"current_page"}},{"kind":"Field","name":{"kind":"Name","value":"per_page"}},{"kind":"Field","name":{"kind":"Name","value":"total_page"}}]}}]}}]} as unknown as DocumentNode<HqStoresQuery, HqStoresQueryVariables>;
export const HqStoresTotalDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"HqStoresTotal"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"stores"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"IntValue","value":"1"}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"IntValue","value":"1"}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}}]}}]}}]} as unknown as DocumentNode<HqStoresTotalQuery, HqStoresTotalQueryVariables>;
export const HqRolesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"HqRoles"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"q"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"OperatorRoleFilterType"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"operatorRoles"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}},{"kind":"Argument","name":{"kind":"Name","value":"q"},"value":{"kind":"Variable","name":{"kind":"Name","value":"q"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"kind"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"permissions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"action"}},{"kind":"Field","name":{"kind":"Name","value":"module"}},{"kind":"Field","name":{"kind":"Name","value":"scope"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"current_page"}},{"kind":"Field","name":{"kind":"Name","value":"per_page"}},{"kind":"Field","name":{"kind":"Name","value":"total_page"}}]}}]}}]} as unknown as DocumentNode<HqRolesQuery, HqRolesQueryVariables>;
export const SystemPermissionsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"SystemPermissions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"permissions"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"IntValue","value":"1"}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"IntValue","value":"200"}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"scope"},"value":{"kind":"EnumValue","value":"SYSTEM"}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"action"}},{"kind":"Field","name":{"kind":"Name","value":"module"}},{"kind":"Field","name":{"kind":"Name","value":"scope"}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}}]}}]}}]} as unknown as DocumentNode<SystemPermissionsQuery, SystemPermissionsQueryVariables>;
export const HqCreateRoleDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqCreateRole"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateOperatorRoleInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createOperatorRole"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"kind"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}}]}}]}}]} as unknown as DocumentNode<HqCreateRoleMutation, HqCreateRoleMutationVariables>;
export const HqUpdateRoleDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqUpdateRole"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"UpdateOperatorRoleInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateOperatorRole"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"kind"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}}]}}]}}]} as unknown as DocumentNode<HqUpdateRoleMutation, HqUpdateRoleMutationVariables>;
export const HqDeleteRolesDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqDeleteRoles"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"ids"}},"type":{"kind":"NonNullType","type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteOperatorRoles"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"ids"}}}]}]}}]} as unknown as DocumentNode<HqDeleteRolesMutation, HqDeleteRolesMutationVariables>;
export const HqStoreApprovalsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"HqStoreApprovals"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"q"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filter"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"StoreFilterType"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"stores"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"current_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"per_page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}},{"kind":"Argument","name":{"kind":"Name","value":"q"},"value":{"kind":"Variable","name":{"kind":"Name","value":"q"}}},{"kind":"Argument","name":{"kind":"Name","value":"filter"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filter"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"data"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"code"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}},{"kind":"Field","name":{"kind":"Name","value":"submittedAt"}},{"kind":"Field","name":{"kind":"Name","value":"organizationId"}},{"kind":"Field","name":{"kind":"Name","value":"organization"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"contactPhone"}},{"kind":"Field","name":{"kind":"Name","value":"managerName"}},{"kind":"Field","name":{"kind":"Name","value":"managerPhone"}},{"kind":"Field","name":{"kind":"Name","value":"province"}},{"kind":"Field","name":{"kind":"Name","value":"city"}},{"kind":"Field","name":{"kind":"Name","value":"district"}},{"kind":"Field","name":{"kind":"Name","value":"address"}},{"kind":"Field","name":{"kind":"Name","value":"businessHours"}},{"kind":"Field","name":{"kind":"Name","value":"businessStatus"}},{"kind":"Field","name":{"kind":"Name","value":"supportDineIn"}},{"kind":"Field","name":{"kind":"Name","value":"supportTakeout"}},{"kind":"Field","name":{"kind":"Name","value":"storeArea"}},{"kind":"Field","name":{"kind":"Name","value":"tableCount"}}]}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"current_page"}},{"kind":"Field","name":{"kind":"Name","value":"per_page"}},{"kind":"Field","name":{"kind":"Name","value":"total_page"}}]}}]}}]} as unknown as DocumentNode<HqStoreApprovalsQuery, HqStoreApprovalsQueryVariables>;
export const HqReviewStoreDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"HqReviewStore"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ReviewStoreInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"reviewStore"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"lifecycle"}},{"kind":"Field","name":{"kind":"Name","value":"rejectionReason"}},{"kind":"Field","name":{"kind":"Name","value":"reviewedAt"}}]}}]}}]} as unknown as DocumentNode<HqReviewStoreMutation, HqReviewStoreMutationVariables>;