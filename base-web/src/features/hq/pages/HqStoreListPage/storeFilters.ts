import type { OrganizationType, StoreBusinessStatus, StoreFilterType, StoreLifecycle } from '@/__generated__/graphql';

export function hqStoreFilter(
  type?: OrganizationType,
  organizationId?: string,
  lifecycle?: StoreLifecycle,
  businessStatus?: StoreBusinessStatus,
): StoreFilterType {
  return {
    ...(type && { organization: { type } }),
    ...(organizationId && { organizationId }),
    ...(lifecycle && { lifecycle }),
    ...(businessStatus && { businessStatus }),
  };
}
