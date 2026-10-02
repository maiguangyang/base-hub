import { describe, expect, it } from 'vitest';
import {
  hqAdministratorFilter,
  hqAuditFilter,
  hqDirectStoreFilter,
  hqFranchiseFilter,
  hqFranchiseStoreFilter,
  hqRoleFilter,
  hqStoreApprovalFilter,
} from './hqListFilters';

describe('总部列表 GraphQL 筛选映射', () => {
  it('直营门店组合总部边界、准入状态和营业状态', () => {
    expect(hqDirectStoreFilter('hq-1', 'ACTIVE', 'OPEN')).toEqual({
      organizationId: 'hq-1', lifecycle: 'ACTIVE', businessStatus: 'OPEN',
    });
    expect(hqDirectStoreFilter('hq-1')).toEqual({ organizationId: 'hq-1' });
  });

  it('加盟商始终限制组织类型并按状态筛选', () => {
    expect(hqFranchiseFilter('SUSPENDED')).toEqual({ type: 'FRANCHISE', status: 'SUSPENDED' });
    expect(hqFranchiseFilter()).toEqual({ type: 'FRANCHISE' });
  });

  it('加盟门店支持汇总模式和指定加盟商模式', () => {
    expect(hqFranchiseStoreFilter()).toEqual({ organization: { type: 'FRANCHISE' } });
    expect(hqFranchiseStoreFilter('org-1', 'ACTIVE', 'OPEN')).toEqual({
      organizationId: 'org-1', organization: { type: 'FRANCHISE' }, lifecycle: 'ACTIVE', businessStatus: 'OPEN',
    });
  });

  it('待审核门店保留固定生命周期并可限定加盟商', () => {
    expect(hqStoreApprovalFilter('org-1')).toEqual({ lifecycle: 'PENDING_APPROVAL', organizationId: 'org-1' });
    expect(hqStoreApprovalFilter()).toEqual({ lifecycle: 'PENDING_APPROVAL' });
  });

  it('管理员按成员状态和角色关系组合筛选', () => {
    expect(hqAdministratorFilter('ACTIVE', 'role-1')).toEqual({ status: 'ACTIVE', roles: { id: 'role-1' } });
    expect(hqAdministratorFilter()).toBeUndefined();
  });

  it('总部角色保留组织边界并按类型筛选', () => {
    expect(hqRoleFilter('hq-1', 'CUSTOM')).toEqual({ organizationId: 'hq-1', kind: 'CUSTOM' });
    expect(hqRoleFilter('hq-1')).toEqual({ organizationId: 'hq-1' });
  });

  it('审计结果将失败映射成非成功结果码', () => {
    expect(hqAuditFilter('SUCCESS', 'store')).toEqual({ resultCode: 'SUCCESS', resourceType: 'store' });
    expect(hqAuditFilter('FAILURE')).toEqual({ resultCode_ne: 'SUCCESS' });
    expect(hqAuditFilter()).toBeUndefined();
  });
});
