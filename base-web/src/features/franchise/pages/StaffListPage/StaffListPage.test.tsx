// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { StaffListPage } from './StaffListPage';

const mocks = vi.hoisted(() => ({
  permissions: ['operatorMembership:read'] as string[],
  listError: undefined as Error | undefined,
  roleResult: { data: [] as Array<{ id: string; name: string; kind?: string }>, total: 0 },
  storeResult: { data: [] as Array<{ id: string; name: string }>, total: 0 },
  refetch: vi.fn(async () => undefined),
  fetchMore: vi.fn(async () => { throw new Error('offline'); }),
}));

vi.mock('@apollo/client/react', () => ({
  useQuery: (document: { definitions?: Array<{ name?: { value?: string } }> }) => {
    const name = document.definitions?.find((definition) => definition.name)?.name?.value;
    if (name === 'FranchiseStaff') {
      return { data: undefined, loading: false, error: mocks.listError, refetch: mocks.refetch };
    }
    if (name === 'FranchiseRoles') {
      return { data: { operatorRoles: mocks.roleResult }, loading: false, error: undefined, fetchMore: mocks.fetchMore };
    }
    return { data: { stores: mocks.storeResult }, loading: false, error: undefined, fetchMore: mocks.fetchMore };
  },
  useMutation: () => [vi.fn(), { loading: false }],
}));

vi.mock('@/features/admin/hooks/useAdminTab', () => ({
  useAdminTab: () => ({ search: '' }),
}));

vi.mock('@/features/auth/store/authStore', () => ({
  useAuthStore: (selector: (state: unknown) => unknown) => selector({
    viewer: {
      account: { id: 'account-1' },
      currentWorkspace: { workspaceType: 'FRANCHISE', organizationId: 'org-a' },
      permissions: mocks.permissions,
      workspaces: [],
    },
  }),
}));

beforeEach(() => {
  mocks.permissions = ['operatorMembership:read'];
  mocks.listError = undefined;
  mocks.roleResult = { data: [], total: 0 };
  mocks.storeResult = { data: [], total: 0 };
  vi.clearAllMocks();
});

afterEach(() => {
  cleanup();
});

describe('加盟商员工列表查询状态', () => {
  it('使用新增员工入口打开新增弹窗', () => {
    mocks.permissions.push('operatorMembership:create');
    render(<StaffListPage />);

    fireEvent.click(screen.getByRole('button', { name: '新增员工' }));
    expect(screen.getByRole('heading', { name: '新增员工' })).toBeTruthy();
    expect(screen.queryByText('邀请员工')).toBeNull();
  });

  it('角色筛选显示系统角色名称', () => {
    mocks.roleResult = { data: [{ id: 'owner-role', name: 'role.franchiseOwner', kind: 'FRANCHISE_OWNER' }], total: 1 };
    render(<StaffListPage />);

    fireEvent.click(screen.getByRole('combobox', { name: '全部角色' }));
    expect(screen.getByRole('option', { name: '加盟商负责人' })).toBeTruthy();
    expect(screen.queryByText('role.franchiseOwner')).toBeNull();
  });

  it('主列表查询失败时显示错误而不是空列表', () => {
    mocks.listError = new Error('offline');
    render(<StaffListPage />);

    expect(screen.getByRole('alert').textContent).toContain('员工列表加载失败，请稍后重试。');
    expect(screen.queryByText('暂无员工')).toBeNull();
  });

  it('关联目录后续页失败时禁用对应筛选器并提示错误', async () => {
    mocks.roleResult = { data: [{ id: 'role-1', name: '店员' }], total: 2 };
    mocks.storeResult = { data: [{ id: 'store-1', name: '一号店' }], total: 2 };
    render(<StaffListPage />);

    expect(screen.getByRole('combobox', { name: '全部角色' }).hasAttribute('disabled')).toBe(true);
    expect(screen.getByRole('combobox', { name: '全部指定门店' }).hasAttribute('disabled')).toBe(true);
    await waitFor(() => expect(screen.getByRole('alert').textContent).toContain('角色或门店目录未完整加载，部分筛选暂不可用。'));
  });
});
