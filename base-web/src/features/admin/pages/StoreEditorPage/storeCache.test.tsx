// @vitest-environment jsdom

import { ApolloClient, ApolloLink, InMemoryCache, Observable } from '@apollo/client';
import { ApolloProvider } from '@apollo/client/react';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { HQ_STORES_QUERY } from '@/features/hq/graphql/hqStores';
import { FRANCHISE_STORES_QUERY } from '@/features/franchise/graphql/stores';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';
import { StoreEditorPage } from './StoreEditorPage';
import type { StoreLifecycle } from '@/__generated__/graphql';

const navigation = vi.hoisted(() => ({ navigate: vi.fn(), workspace: 'HEADQUARTERS' }));
vi.mock('react-router', async (importOriginal) => ({ ...(await importOriginal<typeof import('react-router')>()), useNavigate: () => navigation.navigate }));
vi.mock('@/features/admin/hooks/useAdminTab', () => ({ useAdminTab: () => ({ search: '' }) }));
vi.mock('@/features/auth/store/authStore', () => ({ useAuthStore: (select: (state: unknown) => unknown) => select({
  viewer: { permissions: navigation.workspace === 'HEADQUARTERS' ? ['hqStore:create'] : ['store:create'],
    currentWorkspace: { workspaceType: navigation.workspace, organizationId: 'org' } },
}) }));

beforeEach(() => { vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }); });
afterEach(() => { cleanup(); navigation.navigate.mockReset(); vi.unstubAllGlobals(); });

const storeFields = {
  contactPhone: null, managerName: null, managerPhone: null, province: null, city: null,
  district: null, address: null, businessHours: null, businessStatus: 'OPEN',
  supportDineIn: true, supportTakeout: true, storeArea: null, tableCount: null, receiptFooter: null,
};

it('shows a franchise INTERNAL_ERROR in a floating toast and allows retry', async () => {
  navigation.workspace = 'FRANCHISE';
  let attempts = 0;
  const client = new ApolloClient({ cache: new InMemoryCache(), link: new ApolloLink(() => new Observable((observer) => {
    attempts += 1;
    if (attempts === 1) observer.next({ data: null, errors: [{ message: 'INTERNAL_ERROR', path: ['createStore'],
      extensions: { code: 'INTERNAL_ERROR' } }] });
    else observer.next({ data: { createStore: { __typename: 'Store', id: 'new', code: 'STR2', name: '新店',
      lifecycle: 'DRAFT', organizationId: 'org' } } });
    observer.complete();
  })) });
  const view = render(<ApolloProvider client={client}><StoreEditorPage workspace="FRANCHISE" /></ApolloProvider>, { wrapper: AdminToastProvider });
  for (const label of [/门店名称/, /联系电话/, /省份/, /城市/, /区县/, /详细地址/]) {
    fireEvent.change(screen.getByLabelText(label), { target: { value: '新店' } });
  }
  fireEvent.submit(document.getElementById('store-editor-form')!);
  const alert = await screen.findByRole('alert');
  expect(alert.textContent).toBe('门店保存失败：服务端处理异常，请稍后重试');
  expect(view.container.contains(alert)).toBe(false);
  expect(alert.parentElement?.className).toContain('fixed');
  expect(alert.className).toContain('bg-danger-bg');
  expect(navigation.navigate).not.toHaveBeenCalled();
  expect((screen.getByLabelText(/门店名称/) as HTMLInputElement).value).toBe('新店');
  fireEvent.click(screen.getByRole('button', { name: '保存' }));
  await waitFor(() => expect(navigation.navigate).toHaveBeenCalledWith('/admin/franchise/stores'));
  expect(attempts).toBe(2);
});

it.each([
  ['HEADQUARTERS', HQ_STORES_QUERY, 'HqStores', 'HqCreateDirectStore'],
  ['FRANCHISE', FRANCHISE_STORES_QUERY, 'FranchiseStores', 'FranchiseCreateStore'],
] as const)('%s save refreshes an already watched filtered list', async (workspace, listQuery, listName, createName) => {
  navigation.workspace = workspace;
  const requests: Array<{ name: string; variables: Record<string, unknown> }> = [];
  let saved = false;
  const link = new ApolloLink((operation) => new Observable((observer) => {
    requests.push({ name: operation.operationName ?? '', variables: operation.variables });
    if (operation.operationName === listName) {
      observer.next({ data: { stores: { __typename: 'StorePagination', total: 1, current_page: 3, per_page: 5,
        total_page: 3, data: [{ __typename: 'Store', id: saved ? 'new' : 'old', name: saved ? '新店' : '旧店',
        code: 'STR1', lifecycle: workspace === 'HEADQUARTERS' ? 'ACTIVE' : 'DRAFT', organizationId: 'org',
          rejectionReason: null, ...storeFields,
          organization: { __typename: 'Organization', id: 'org', name: '机构', type: workspace } }] } } });
    } else if (operation.operationName === createName) {
      saved = true;
      observer.next({ data: { createStore: { __typename: 'Store', id: 'new', code: 'STR2', name: '新店',
        lifecycle: workspace === 'HEADQUARTERS' ? 'ACTIVE' : 'DRAFT', organizationId: 'org', ...storeFields } } });
    } else observer.error(new Error(`unexpected operation: ${operation.operationName}`));
    observer.complete();
  }));
  const client = new ApolloClient({ link, cache: new InMemoryCache() });
  const variables = { page: 3, pageSize: 5, q: '旧', filter: { lifecycle: (workspace === 'HEADQUARTERS' ? 'ACTIVE' : 'DRAFT') as StoreLifecycle } };
  const watcher = client.watchQuery({ query: listQuery, variables, fetchPolicy: 'cache-and-network' });
  const seen: string[] = [];
  const subscription = watcher.subscribe(({ data }) => {
    const name = data?.stores?.data?.[0]?.name;
    if (name) seen.push(name);
  });
  try {
    await waitFor(() => expect(seen).toContain('旧店'));
    render(<ApolloProvider client={client}><StoreEditorPage workspace={workspace} /></ApolloProvider>, { wrapper: AdminToastProvider });
    for (const label of [/门店名称/, /联系电话/, /省份/, /城市/, /区县/, /详细地址/]) {
      fireEvent.change(screen.getByLabelText(label), { target: { value: '新店' } });
    }
    fireEvent.submit(document.getElementById('store-editor-form')!);
    await waitFor(() => expect(navigation.navigate).toHaveBeenCalled());
    await waitFor(() => expect(seen).toContain('新店'));
    const listRequests = requests.filter((request) => request.name === listName);
    expect(listRequests).toHaveLength(2);
    expect(listRequests[1].variables).toEqual(variables);
    expect(requests.filter((request) => request.name === createName)).toHaveLength(1);
  } finally { subscription.unsubscribe(); }
});
