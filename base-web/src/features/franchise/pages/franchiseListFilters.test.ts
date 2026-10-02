import { describe, expect, it } from 'vitest';
import {
  franchiseAuditFilter,
  franchiseRoleFilter,
  franchiseStaffFilter,
  franchiseStoreFilter,
} from './franchiseListFilters';

describe('加盟商列表 GraphQL 筛选映射', () => {
  it('我的门店只按页面展示的准入状态筛选', () => {
    expect(franchiseStoreFilter('ACTIVE')).toEqual({ lifecycle: 'ACTIVE' });
    expect(franchiseStoreFilter()).toBeUndefined();
  });

  it('员工筛选组合成员状态、角色、门店范围和指定门店', () => {
    const filter = franchiseStaffFilter('ACTIVE', 'role-1', 'SELECTED_STORES', 'store-1');
    expect(filter).toEqual({
      status: 'ACTIVE', roles: { id: 'role-1' }, storeAccessMode: 'SELECTED_STORES', stores: { id: 'store-1' },
    });
    expect(filter).not.toHaveProperty('organizationId');
    expect(franchiseStaffFilter()).toBeUndefined();
  });

  it('角色按类型筛选且不发送组织覆盖', () => {
    const filter = franchiseRoleFilter('CUSTOM');
    expect(filter).toEqual({ kind: 'CUSTOM' });
    expect(filter).not.toHaveProperty('organizationId');
  });

  it('审计支持执行结果、资源类型和门店组合', () => {
    expect(franchiseAuditFilter('SUCCESS', 'store', 'store-1')).toEqual({
      resultCode: 'SUCCESS', resourceType: 'store', storeId: 'store-1',
    });
    expect(franchiseAuditFilter('FAILURE')).toEqual({ resultCode_ne: 'SUCCESS' });
    expect(franchiseAuditFilter()).toBeUndefined();
  });
});
