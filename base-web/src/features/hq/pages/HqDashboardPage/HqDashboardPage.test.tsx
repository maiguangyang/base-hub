// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { HqDashboardPage } from './HqDashboardPage';

const mocks = vi.hoisted(() => ({
  queries: [] as Array<{ name: string; skip?: boolean }>,
  permissions: ['store:read_all'] as string[],
  storesError: false,
}));

function mockStoresQuery(skip?: boolean) {
  if (skip) return { data: undefined };
  if (mocks.storesError) return { data: undefined, error: new Error('offline') };
  return { data: { stores: { total: 7 } } };
}

vi.mock('@apollo/client/react', () => ({
  useQuery: (document: { definitions?: Array<{ name?: { value?: string } }> }, options?: { skip?: boolean }) => {
    const name = document.definitions?.find((entry) => entry.name)?.name?.value ?? '';
    mocks.queries.push({ name, skip: options?.skip });
    if (name === 'HqStoresTotal') return mockStoresQuery(options?.skip);
    if (name === 'HqFranchises') return { data: { organizations: { total: 2 } } };
    return { data: { stores: { total: 1 } } };
  },
}));
vi.mock('@/features/admin/hooks/useAdminTab', () => ({ useAdminTab: () => ({}) }));
vi.mock('@/features/auth/store/authStore', () => ({
  useAuthStore: (select: (state: unknown) => unknown) => select({ viewer: { permissions: mocks.permissions } }),
}));

afterEach(() => { cleanup(); mocks.queries.length = 0; mocks.permissions = ['store:read_all']; mocks.storesError = false; });

it('总部概览用全门店总数查询显示门店数量', () => {
  render(<HqDashboardPage />);
  expect(mocks.queries).toContainEqual({ name: 'HqStoresTotal', skip: false });
  expect(screen.getByText('门店总数')).toBeTruthy();
  expect(screen.getByText('7')).toBeTruthy();
});

it('没有跨组织门店权限时不查询总数，也不把本组织数量当成总数', () => {
  mocks.permissions = ['hqStore:read'];
  render(<HqDashboardPage />);
  expect(mocks.queries).toContainEqual({ name: 'HqStoresTotal', skip: true });
  expect(screen.getByText('无权限')).toBeTruthy();
  expect(screen.queryByText('7')).toBeNull();
});

it('门店总数查询失败时显示失败状态而不是零', () => {
  mocks.storesError = true;
  render(<HqDashboardPage />);
  expect(screen.getByText('加载失败')).toBeTruthy();
  expect(screen.queryByText('0')).toBeNull();
});
