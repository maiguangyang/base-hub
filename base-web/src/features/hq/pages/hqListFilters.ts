import type {
  AuditLogFilterType,
  MembershipStatus,
  OperatorMembershipFilterType,
  OperatorRoleFilterType,
  OrganizationFilterType,
  OrganizationStatus,
  RoleKind,
  StoreBusinessStatus,
  StoreFilterType,
  StoreLifecycle,
} from '@/__generated__/graphql';

export type AuditResultFilter = 'SUCCESS' | 'FAILURE' | undefined;

export function hqDirectStoreFilter(
  organizationId: string,
  lifecycle?: StoreLifecycle,
  businessStatus?: StoreBusinessStatus,
): StoreFilterType {
  return { organizationId, ...(lifecycle && { lifecycle }), ...(businessStatus && { businessStatus }) };
}

export function hqFranchiseFilter(status?: OrganizationStatus): OrganizationFilterType {
  return { type: 'FRANCHISE', ...(status && { status }) };
}

export function hqFranchiseStoreFilter(
  organizationId?: string,
  lifecycle?: StoreLifecycle,
  businessStatus?: StoreBusinessStatus,
): StoreFilterType {
  return {
    ...(organizationId && { organizationId }),
    organization: { type: 'FRANCHISE' },
    ...(lifecycle && { lifecycle }),
    ...(businessStatus && { businessStatus }),
  };
}

export function hqStoreApprovalFilter(organizationId?: string): StoreFilterType {
  return { lifecycle: 'PENDING_APPROVAL', ...(organizationId && { organizationId }) };
}

export function hqAdministratorFilter(
  status?: MembershipStatus,
  roleId?: string,
): OperatorMembershipFilterType | undefined {
  if (!status && !roleId) return undefined;
  return { ...(status && { status }), ...(roleId && { roles: { id: roleId } }) };
}

export function hqRoleFilter(organizationId: string, kind?: RoleKind): OperatorRoleFilterType {
  return { organizationId, ...(kind && { kind }) };
}

export function hqAuditFilter(
  result?: AuditResultFilter,
  resourceType?: string,
): AuditLogFilterType | undefined {
  if (!result && !resourceType) return undefined;
  return {
    ...(result === 'SUCCESS' && { resultCode: 'SUCCESS' }),
    ...(result === 'FAILURE' && { resultCode_ne: 'SUCCESS' }),
    ...(resourceType && { resourceType }),
  };
}
