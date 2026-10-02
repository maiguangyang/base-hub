// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { DirectStoreListPage } from '@/features/hq/pages/DirectStoreListPage';
import { FranchiseListPage } from '@/features/hq/pages/FranchiseListPage';
import { StoreApprovalListPage } from '@/features/hq/pages/StoreApprovalListPage';
import { AuditLogListPage as HqAuditLogListPage } from '@/features/hq/pages/AuditLogListPage';
import { StoreListPage } from '@/features/franchise/pages/StoreListPage';
import { RoleListPage as FranchiseRoleListPage } from '@/features/franchise/pages/RoleListPage';
import { AuditLogListPage as FranchiseAuditLogListPage } from '@/features/franchise/pages/AuditLogListPage';

const mocks = vi.hoisted(() => ({ failedQuery: '', workspaceType: 'HEADQUARTERS', franchiseRows: false, mutationFails: false, permissions: [] as string[] }));

vi.mock('@apollo/client/react', () => ({
  useQuery: (document: { definitions?: Array<{ name?: { value?: string } }> }) => {
    const name = document.definitions?.find((definition) => definition.name)?.name?.value;
    if (name === mocks.failedQuery) return { data: undefined, loading: false, error: new Error('offline'), refetch: vi.fn() };
    if (name === 'HqFranchises' && mocks.franchiseRows) return { data: { organizations: { data: [{ id: 'franchise-1', code: 'FR-1', name: '甲加盟商', status: 'SUSPENDED' }], total: 1 } }, loading: false, refetch: vi.fn() };
    if (name === 'HqFranchises') return { data: { organizations: { data: [], total: 0 } }, loading: false };
    if (name === 'FranchiseStores') return { data: { stores: { data: [], total: 0 } }, loading: false };
    if (name === 'TenantPermissions') return { data: { permissions: { data: [], total: 0 } }, loading: false };
    return { data: undefined, loading: false };
  },
  useMutation: () => [vi.fn(async () => { if (mocks.mutationFails) throw new Error('mutation failed'); }), { loading: false }],
}));

vi.mock('@/features/admin/hooks/useAdminTab', () => ({ useAdminTab: () => ({ search: '' }) }));
vi.mock('@/features/auth/store/authStore', () => ({
  useAuthStore: (selector: (state: unknown) => unknown) => selector({
    viewer: {
      account: { id: 'account-1' },
      currentWorkspace: { workspaceType: mocks.workspaceType, organizationId: 'org-1' },
      permissions: mocks.permissions, workspaces: [],
    },
  }),
}));

beforeEach(() => { mocks.failedQuery = ''; mocks.workspaceType = 'HEADQUARTERS'; mocks.franchiseRows = false; mocks.mutationFails = false; mocks.permissions = []; });
afterEach(cleanup);

const cases = [
  ['HqDirectStores', DirectStoreListPage, '直营门店', '暂无直营门店', 'HEADQUARTERS'],
  ['HqFranchises', FranchiseListPage, '加盟商', '暂无加盟商', 'HEADQUARTERS'],
  ['HqStoreApprovals', StoreApprovalListPage, '待审核门店', '暂无待审核门店', 'HEADQUARTERS'],
  ['HqAuditLogs', HqAuditLogListPage, '审计日志', '暂无审计记录', 'HEADQUARTERS'],
  ['FranchiseStores', StoreListPage, '门店', '暂无门店', 'FRANCHISE'],
  ['FranchiseRoles', FranchiseRoleListPage, '角色', '暂无角色', 'FRANCHISE'],
  ['FranchiseAuditLogs', FranchiseAuditLogListPage, '审计日志', '暂无审计记录', 'FRANCHISE'],
] as const;

describe('后台列表主查询错误', () => {
  it.each(cases)('%s 失败时提示错误且不显示空状态', (queryName, Page, title, emptyText, workspaceType) => {
    mocks.failedQuery = queryName;
    mocks.workspaceType = workspaceType;
    render(<MemoryRouter><Page /></MemoryRouter>);

    expect(screen.getByRole('alert').textContent).toContain(`${title}加载失败`);
    expect(screen.queryByText(emptyText)).toBeNull();
  });

  it('操作失败后主列表查询失败时仍显示当前列表错误', async () => {
    mocks.franchiseRows = true;
    mocks.mutationFails = true;
    mocks.permissions = ['organization:restore'];
    const view = render(<MemoryRouter><FranchiseListPage /></MemoryRouter>);

    fireEvent.click(screen.getByRole('switch', { name: '甲加盟商状态' }));
    await waitFor(() => expect(screen.getByRole('alert').textContent).toContain('操作失败'));

    mocks.failedQuery = 'HqFranchises';
    view.rerender(<MemoryRouter><FranchiseListPage /></MemoryRouter>);
    expect(screen.getByRole('alert').textContent).toContain('加盟商加载失败');
  });
});
