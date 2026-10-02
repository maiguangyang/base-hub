// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { DirectStoreListPage } from '@/features/hq/pages/DirectStoreListPage';
import { HqStoreListPage } from '@/features/hq/pages/HqStoreListPage/HqStoreListPage';
import { FranchiseListPage } from '@/features/hq/pages/FranchiseListPage';
import { FranchiseStoreListPage } from '@/features/hq/pages/FranchiseStoreListPage';
import { StoreApprovalListPage } from '@/features/hq/pages/StoreApprovalListPage';
import { AdministratorListPage } from '@/features/hq/pages/AdministratorListPage';
import { RoleListPage as HqRoleListPage } from '@/features/hq/pages/RoleListPage';
import { AuditLogListPage as HqAuditLogListPage } from '@/features/hq/pages/AuditLogListPage';
import { StoreListPage } from '@/features/franchise/pages/StoreListPage';
import { StaffListPage } from '@/features/franchise/pages/StaffListPage';
import { RoleListPage as FranchiseRoleListPage } from '@/features/franchise/pages/RoleListPage';
import { AuditLogListPage as FranchiseAuditLogListPage } from '@/features/franchise/pages/AuditLogListPage';

const mocks = vi.hoisted(() => ({
  calls: [] as Array<{ name: string | undefined; variables: Record<string, unknown> | undefined }>,
  workspaceType: 'HEADQUARTERS' as 'HEADQUARTERS' | 'FRANCHISE',
  largeCatalog: '',
  fetchMoreCalls: [] as Array<{ name: string | undefined; page: number }>,
}));

vi.mock('@apollo/client/react', () => {
  function catalogRoot(name: string | undefined) {
    if (name === 'HqFranchises') return 'organizations';
    if (name === 'HqRoles' || name === 'FranchiseRoles') return 'operatorRoles';
    return 'stores';
  }
  function largeCatalogData(name: string | undefined) {
    return { [catalogRoot(name)]: { data: Array.from({ length: 200 }, (_, index) => ({ id: `entry-${index + 1}`, name: `第${index + 1}项` })), total: 201 } };
  }
  function queryData(name: string | undefined, pageSize: unknown) {
    const list = { data: [], total: 100 };
    const listData: Record<string, unknown> = {
      HqFranchises: { organizations: list }, HqRoles: { operatorRoles: list }, FranchiseRoles: { operatorRoles: list },
      HqAdministrators: { operatorMemberships: list }, FranchiseStaff: { operatorMemberships: list },
      HqAuditLogs: { auditLogs: list }, FranchiseAuditLogs: { auditLogs: list },
      SystemPermissions: { permissions: { data: [], total: 0 } }, TenantPermissions: { permissions: { data: [], total: 0 } },
      HqAdministratorIdentity: { operatorMemberships: { data: [], total: 0 } },
    };
    const catalogData: Record<string, unknown> = {
      HqFranchises: { organizations: { data: [{ id: 'org-franchise', name: '甲加盟商' }], total: 1 } },
      HqRoles: { operatorRoles: { data: [{ id: 'role-hq', name: '运营', kind: 'CUSTOM' }], total: 1 } },
      FranchiseRoles: { operatorRoles: { data: [{ id: 'role-tenant', name: '店员', kind: 'CUSTOM' }], total: 1 } },
      FranchiseStores: { stores: { data: [{ id: 'store-tenant', name: '甲门店' }], total: 1 } },
    };
    if (pageSize === 200 && name === mocks.largeCatalog) return largeCatalogData(name);
    return (pageSize === 200 ? catalogData[name ?? ''] : undefined) ?? listData[name ?? ''] ?? { stores: list };
  }
  return {
    useQuery: (document: { definitions?: Array<{ name?: { value?: string } }> }, options?: { variables?: Record<string, unknown> }) => {
      const name = document.definitions?.find((definition) => definition.name)?.name?.value;
      mocks.calls.push({ name, variables: options?.variables });
      const root = catalogRoot(name);
      return { data: queryData(name, options?.variables?.pageSize), loading: false, error: undefined, refetch: vi.fn(), fetchMore: async ({ variables }: { variables: { page: number } }) => {
        mocks.fetchMoreCalls.push({ name, page: variables.page });
        return { data: { [root]: { data: [{ id: 'entry-201', name: '第201项' }], total: 201 } } };
      } };
    },
    useMutation: () => [vi.fn(), { loading: false }],
  };
});

vi.mock('@/features/admin/hooks/useAdminTab', () => ({ useAdminTab: () => ({ search: '?q=咖啡' }) }));
vi.mock('@/features/auth/store/authStore', () => ({
  useAuthStore: (selector: (state: unknown) => unknown) => selector({
    viewer: {
      account: { id: 'account-1' },
      currentWorkspace: { workspaceType: mocks.workspaceType, organizationId: 'org-1' },
      permissions: mocks.workspaceType === 'HEADQUARTERS' ? ['organization:read', 'store:read_all'] : [], workspaces: [],
    },
  }),
}));

beforeEach(() => { mocks.calls.length = 0; mocks.fetchMoreCalls.length = 0; mocks.workspaceType = 'HEADQUARTERS'; mocks.largeCatalog = ''; });
afterEach(cleanup);

function latestVariables(name: string): Record<string, unknown> | undefined {
  return [...mocks.calls].reverse().find((call) => call.name === name)?.variables;
}

function chooseFilter(label: string, option: string) {
  fireEvent.click(screen.getByRole('combobox', { name: label }));
  fireEvent.click(screen.getByRole('option', { name: option }));
}

function clearIndividualFilters(query: string, selections: readonly (readonly [string, string])[], active: Record<string, unknown>, reset: unknown, extraVariables: Record<string, unknown>) {
  // 清空关键词时保留全部结构化条件，并回到第一页。
  fireEvent.click(screen.getByRole('button', { name: '下一页' }));
  const searchClear = screen.getAllByRole('button').find((button) => button.getAttribute('aria-label')?.startsWith('清空搜索'))!;
  fireEvent.click(searchClear);
  expect(latestVariables(query)).toEqual({ page: 1, pageSize: 20, q: null, filter: active, ...extraVariables });
  const expectedFilter = { ...active };
  const filterKeys: Record<string, string> = {
    全部类型: 'organization', 全部加盟商: 'organizationId', 全部准入状态: 'lifecycle', 全部营业状态: 'businessStatus',
    全部组织状态: 'status', 全部成员状态: 'status', 全部角色: 'roles', 全部角色类型: 'kind',
    全部执行结果: 'resultCode_ne', 全部资源类型: 'resourceType', 全部门店范围: 'storeAccessMode',
    全部指定门店: 'stores', 全部门店: 'storeId',
  };
  for (const [label] of [...selections].reverse()) {
    fireEvent.click(screen.getByRole('button', { name: '下一页' }));
    fireEvent.click(screen.getByRole('button', { name: `清空${label}` }));
    delete expectedFilter[filterKeys[label]];
    expect(latestVariables(query)).toEqual({ page: 1, pageSize: 20, q: null,
      filter: Object.keys(expectedFilter).length ? expectedFilter : reset, ...extraVariables });
  }
}

const cases = [
  { name: '总部统一门店', Page: HqStoreListPage, query: 'HqStores', workspace: 'HEADQUARTERS', selections: [['全部类型', '加盟'], ['全部加盟商', '甲加盟商'], ['全部准入状态', '已启用'], ['全部营业状态', '营业中']], active: { organization: { type: 'FRANCHISE' }, organizationId: 'org-franchise', lifecycle: 'ACTIVE', businessStatus: 'OPEN' }, reset: {} },
  { name: '总部直营门店', Page: DirectStoreListPage, query: 'HqDirectStores', workspace: 'HEADQUARTERS', selections: [['全部准入状态', '已启用'], ['全部营业状态', '营业中']], active: { organizationId: 'org-1', lifecycle: 'ACTIVE', businessStatus: 'OPEN' }, reset: { organizationId: 'org-1' } },
  { name: '总部加盟商', Page: FranchiseListPage, query: 'HqFranchises', workspace: 'HEADQUARTERS', selections: [['全部组织状态', '已暂停']], active: { type: 'FRANCHISE', status: 'SUSPENDED' }, reset: { type: 'FRANCHISE' } },
  { name: '总部加盟门店', Page: FranchiseStoreListPage, query: 'HqFranchiseStores', workspace: 'HEADQUARTERS', selections: [['全部加盟商', '甲加盟商'], ['全部准入状态', '已启用'], ['全部营业状态', '营业中']], active: { organizationId: 'org-franchise', organization: { type: 'FRANCHISE' }, lifecycle: 'ACTIVE', businessStatus: 'OPEN' }, reset: { organization: { type: 'FRANCHISE' } } },
  { name: '总部门店审核', Page: StoreApprovalListPage, query: 'HqStoreApprovals', workspace: 'HEADQUARTERS', selections: [['全部加盟商', '甲加盟商']], active: { lifecycle: 'PENDING_APPROVAL', organizationId: 'org-franchise' }, reset: { lifecycle: 'PENDING_APPROVAL' } },
  { name: '总部管理员', Page: AdministratorListPage, query: 'HqAdministrators', workspace: 'HEADQUARTERS', selections: [['全部成员状态', '正常'], ['全部角色', '运营']], active: { status: 'ACTIVE', roles: { id: 'role-hq' } }, reset: undefined },
  { name: '总部角色', Page: HqRoleListPage, query: 'HqRoles', workspace: 'HEADQUARTERS', selections: [['全部角色类型', '自定义角色']], active: { organizationId: 'org-1', kind: 'CUSTOM' }, reset: { organizationId: 'org-1' } },
  { name: '总部审计', Page: HqAuditLogListPage, query: 'HqAuditLogs', workspace: 'HEADQUARTERS', selections: [['全部执行结果', '失败'], ['全部资源类型', '门店']], active: { resultCode_ne: 'SUCCESS', resourceType: 'store' }, reset: undefined },
  { name: '加盟商门店', Page: StoreListPage, query: 'FranchiseStores', workspace: 'FRANCHISE', selections: [['全部准入状态', '已启用']], active: { lifecycle: 'ACTIVE' }, reset: undefined },
  { name: '加盟商员工', Page: StaffListPage, query: 'FranchiseStaff', workspace: 'FRANCHISE', selections: [['全部成员状态', '正常'], ['全部角色', '店员'], ['全部门店范围', '指定门店'], ['全部指定门店', '甲门店']], active: { status: 'ACTIVE', roles: { id: 'role-tenant' }, storeAccessMode: 'SELECTED_STORES', stores: { id: 'store-tenant' } }, reset: undefined },
  { name: '加盟商角色', Page: FranchiseRoleListPage, query: 'FranchiseRoles', workspace: 'FRANCHISE', selections: [['全部角色类型', '自定义角色']], active: { kind: 'CUSTOM' }, reset: undefined },
  { name: '加盟商审计', Page: FranchiseAuditLogListPage, query: 'FranchiseAuditLogs', workspace: 'FRANCHISE', selections: [['全部执行结果', '失败'], ['全部资源类型', '门店'], ['全部门店', '甲门店']], active: { resultCode_ne: 'SUCCESS', resourceType: 'store', storeId: 'store-tenant' }, reset: undefined },
] as const;

describe('两端列表筛选接线', () => {
  it('门店审核的加盟商筛选可按名称搜索并传入选中 ID', () => {
    render(<MemoryRouter><StoreApprovalListPage /></MemoryRouter>);

    fireEvent.click(screen.getByRole('combobox', { name: '全部加盟商' }));
    fireEvent.change(screen.getByPlaceholderText('搜索加盟商'), { target: { value: '甲加盟商' } });
    fireEvent.click(screen.getByRole('option', { name: '甲加盟商' }));

    expect(latestVariables('HqStoreApprovals')?.filter).toEqual({ lifecycle: 'PENDING_APPROVAL', organizationId: 'org-franchise' });
  });

  it.each(cases)('$name 的搜索、结构化条件、分页及重置共用服务端查询', ({ name, Page, query, workspace, selections, active, reset }) => {
    mocks.workspaceType = workspace;
    render(<MemoryRouter><Page /></MemoryRouter>);

    fireEvent.click(screen.getByRole('combobox', { name: '每页' }));
    fireEvent.click(screen.getByRole('option', { name: '20' }));
    fireEvent.click(screen.getByRole('button', { name: '下一页' }));
    expect(latestVariables(query)).toMatchObject({ page: 2, pageSize: 20, q: '咖啡' });

    fireEvent.change(screen.getByRole('textbox'), { target: { value: ' 新词 ' } });
    fireEvent.click(screen.getByRole('button', { name: '搜索' }));
    expect(latestVariables(query)).toMatchObject({ page: 1, pageSize: 20, q: '新词' });

    for (const [label, option] of selections) {
      fireEvent.click(screen.getByRole('button', { name: '下一页' }));
      chooseFilter(label, option);
      expect(latestVariables(query)).toMatchObject({ page: 1, pageSize: 20, q: '新词' });
    }
    const extraVariables = name === '总部加盟商' ? { canResetPassword: false } : {};
    expect(latestVariables(query)).toEqual({ page: 1, pageSize: 20, q: '新词', filter: active, ...extraVariables });

    clearIndividualFilters(query, selections, active, reset, extraVariables);

    fireEvent.change(screen.getByRole('textbox'), { target: { value: '重置前关键词' } });
    fireEvent.click(screen.getByRole('button', { name: '搜索' }));

    fireEvent.click(screen.getByRole('button', { name: '重置筛选' }));
    expect(latestVariables(query)).toEqual({ page: 1, pageSize: 20, q: null, filter: reset, ...extraVariables });
  });

  it.each([
    ['总部审核加盟商', StoreApprovalListPage, 'HEADQUARTERS', 'HqFranchises', '全部加盟商'],
    ['总部管理员角色', AdministratorListPage, 'HEADQUARTERS', 'HqRoles', '全部角色'],
    ['加盟商员工角色', StaffListPage, 'FRANCHISE', 'FranchiseRoles', '全部角色'],
    ['加盟商员工门店', StaffListPage, 'FRANCHISE', 'FranchiseStores', '全部指定门店'],
    ['加盟商审计门店', FranchiseAuditLogListPage, 'FRANCHISE', 'FranchiseStores', '全部门店'],
  ] as const)('%s 目录超过 200 条后仍可筛选', async (_name, Page, workspace, catalog, label) => {
    mocks.workspaceType = workspace;
    mocks.largeCatalog = catalog;
    render(<MemoryRouter><Page /></MemoryRouter>);

    expect(screen.getByRole('combobox', { name: label }).hasAttribute('disabled')).toBe(true);
    await waitFor(() => expect(screen.getByRole('combobox', { name: label }).hasAttribute('disabled')).toBe(false));
    expect(mocks.fetchMoreCalls).toContainEqual({ name: catalog, page: 2 });
  });
});
