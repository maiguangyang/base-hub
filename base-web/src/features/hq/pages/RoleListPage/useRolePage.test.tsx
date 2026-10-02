// @vitest-environment jsdom

import { act, cleanup, renderHook } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { useRolePage } from './useRolePage';

const mocks = vi.hoisted(() => ({
  queryCalls: [] as Array<{ name: string | undefined; options: Record<string, unknown> | undefined }>,
  search: '?q=运营',
  refetch: vi.fn(async () => undefined),
}));

vi.mock('@apollo/client/react', () => ({
  useQuery: (document: { definitions?: Array<{ name?: { value?: string } }> }, options?: Record<string, unknown>) => {
    const name = document.definitions?.find((definition) => definition.name)?.name?.value;
    mocks.queryCalls.push({ name, options });
    if (name === 'HqRoles') {
      return { data: { operatorRoles: { data: [], total: 0 } }, loading: false, error: undefined, refetch: mocks.refetch };
    }
    return { data: { permissions: { data: [], total: 0 } }, loading: false, error: undefined };
  },
  useMutation: () => [vi.fn(), { loading: false }],
}));

vi.mock('@/features/admin/hooks/useAdminTab', () => ({
  useAdminTab: () => ({ search: mocks.search }),
}));

vi.mock('@/features/auth/store/authStore', () => ({
  useAuthStore: (selector: (state: unknown) => unknown) => selector({
    viewer: {
      account: { id: 'account-hq' },
      currentWorkspace: { workspaceType: 'HEADQUARTERS', organizationId: 'org-hq' },
      permissions: ['hqRole:read'],
      workspaces: [],
    },
  }),
}));

function latestRoleVariables() {
  const call = [...mocks.queryCalls].reverse().find((item) => item.name === 'HqRoles');
  return call?.options?.variables;
}

beforeEach(() => {
  mocks.queryCalls.length = 0;
  mocks.search = '?q=运营';
  vi.clearAllMocks();
});

afterEach(cleanup);

describe('总部角色列表查询状态', () => {
  it('关键词与结构化条件共同进入同一次服务端查询', () => {
    const { result } = renderHook(() => useRolePage());

    act(() => result.current.setPage(4));
    act(() => result.current.changeKind('CUSTOM'));

    expect(latestRoleVariables()).toEqual({
      page: 1,
      pageSize: 10,
      q: '运营',
      filter: { organizationId: 'org-hq', kind: 'CUSTOM' },
    });
  });

  it('统一重置回到第一页并保留页长', () => {
    const { result } = renderHook(() => useRolePage());

    act(() => result.current.setPageSize(50));
    act(() => result.current.setPage(3));
    act(() => result.current.changeKind('CUSTOM'));
    act(() => result.current.resetFilters());

    expect(latestRoleVariables()).toEqual({
      page: 1,
      pageSize: 50,
      q: null,
      filter: { organizationId: 'org-hq' },
    });
  });
});
