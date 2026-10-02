import type {
  AuditLogFilterType,
  MembershipStatus,
  OperatorMembershipFilterType,
  OperatorRoleFilterType,
  RoleKind,
  StoreAccessMode,
  StoreFilterType,
  StoreLifecycle,
} from '@/__generated__/graphql';

export type FranchiseAuditResultFilter = 'SUCCESS' | 'FAILURE' | undefined;

export function franchiseStoreFilter(lifecycle?: StoreLifecycle): StoreFilterType | undefined {
  return lifecycle ? { lifecycle } : undefined;
}

export function franchiseStaffFilter(
  status?: MembershipStatus,
  roleId?: string,
  storeAccessMode?: StoreAccessMode,
  storeId?: string,
): OperatorMembershipFilterType | undefined {
  if (![status, roleId, storeAccessMode, storeId].some(Boolean)) return undefined;
  return {
    ...(status && { status }),
    ...(roleId && { roles: { id: roleId } }),
    ...(storeAccessMode && { storeAccessMode }),
    ...(storeId && { stores: { id: storeId } }),
  };
}

export function franchiseRoleFilter(kind?: RoleKind): OperatorRoleFilterType | undefined {
  return kind ? { kind } : undefined;
}

export function franchiseAuditFilter(
  result?: FranchiseAuditResultFilter,
  resourceType?: string,
  storeId?: string,
): AuditLogFilterType | undefined {
  if (![result, resourceType, storeId].some(Boolean)) return undefined;
  return {
    ...(result === 'SUCCESS' && { resultCode: 'SUCCESS' }),
    ...(result === 'FAILURE' && { resultCode_ne: 'SUCCESS' }),
    ...(resourceType && { resourceType }),
    ...(storeId && { storeId }),
  };
}
