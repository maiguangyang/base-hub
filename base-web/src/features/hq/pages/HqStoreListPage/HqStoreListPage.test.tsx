// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { MemoryRouter } from 'react-router';
import { HqStoreListPage } from './HqStoreListPage';

const mocks = vi.hoisted(() => ({
  update: vi.fn(async () => ({})),
  remove: vi.fn(async (): Promise<object> => { throw new Error('offline'); }),
  refetch: vi.fn(async () => { throw new Error('offline'); }),
  catalogError: false,
  deleteLoading: false,
  navigate: vi.fn(),
}));
vi.mock('react-router', async (importOriginal) => ({ ...(await importOriginal<typeof import('react-router')>()), useNavigate: () => mocks.navigate }));
vi.mock('@apollo/client/react', () => ({
  useQuery: (document: { definitions?: Array<{ name?: { value?: string } }> }) => {
    const name = document.definitions?.find((entry) => entry.name)?.name?.value;
    if (name === 'HqStores') return {
      data: { stores: { data: [{ id: 'direct', code: 'D1', name: '直营一店', organizationId: 'hq',
        organization: { id: 'hq', name: '总部', type: 'HEADQUARTERS' }, lifecycle: 'ACTIVE', businessStatus: 'OPEN',
        contactPhone: '13800000000', province: '上海', city: '上海', district: '徐汇', address: '漕溪路1号',
        businessHours: '09:00 - 22:00' }], total: 1 } },
      loading: false, error: undefined, refetch: mocks.refetch,
    };
    return { data: mocks.catalogError ? undefined : { organizations: { data: [], total: 0 } }, loading: false,
      error: mocks.catalogError ? new Error('offline') : undefined,
      fetchMore: vi.fn(), refetch: vi.fn() };
  },
  useMutation: (document: { definitions?: Array<{ name?: { value?: string } }> }) => {
    const name = document.definitions?.find((entry) => entry.name)?.name?.value;
    return [name === 'HqDeleteDirectStores' ? mocks.remove : mocks.update,
      { loading: name === 'HqDeleteDirectStores' && mocks.deleteLoading }];
  },
}));
vi.mock('@/features/admin/hooks/useAdminTab', () => ({ useAdminTab: () => ({ search: '' }) }));
vi.mock('@/features/auth/store/authStore', () => ({ useAuthStore: (select: (state: unknown) => unknown) => select({ viewer: {
  currentWorkspace: { workspaceType: 'HEADQUARTERS', organizationId: 'hq' },
  permissions: ['store:read_all', 'hqStore:create', 'hqStore:update', 'hqStore:delete', 'organization:read'],
} }) }));

beforeEach(() => { vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }); });
afterEach(() => { cleanup(); vi.clearAllMocks(); vi.unstubAllGlobals(); mocks.catalogError = false; mocks.deleteLoading = false; });

it('营业状态写入成功但列表刷新失败时如实提示操作已完成', async () => {
  render(<MemoryRouter><HqStoreListPage /></MemoryRouter>);
  fireEvent.click(screen.getByRole('switch', { name: '直营一店营业状态' }));
  await waitFor(() => expect(mocks.update).toHaveBeenCalled());
  await waitFor(() => expect(screen.getByText('营业状态已更新，但门店列表刷新失败，请刷新页面。')).toBeTruthy());
});

it.each([['新增直营店', '/admin/hq/stores/manage'], ['编辑', '/admin/hq/stores/manage?id=direct']])('%s进入共用门店页面', (action, path) => {
  render(<MemoryRouter><HqStoreListPage /></MemoryRouter>);
  fireEvent.click(screen.getByRole('button', { name: action }));
  expect(mocks.navigate).toHaveBeenCalledWith(path);
  expect(screen.queryByRole('dialog')).toBeNull();
});

it('删除失败时保留确认框供重试', async () => {
  render(<MemoryRouter><HqStoreListPage /></MemoryRouter>);
  fireEvent.click(screen.getByRole('button', { name: '删除' }));
  fireEvent.click(screen.getByRole('button', { name: '确认删除' }));
  await waitFor(() => expect(mocks.remove).toHaveBeenCalled());
  await waitFor(() => expect(screen.getByText('删除门店失败，请重试。')).toBeTruthy());
  expect(screen.getByRole('alertdialog')).toBeTruthy();
});

it('关闭删除失败的确认框后不在页面或再次打开的确认框中保留旧错误', async () => {
  render(<MemoryRouter><HqStoreListPage /></MemoryRouter>);
  fireEvent.click(screen.getByRole('button', { name: '删除' }));
  fireEvent.click(screen.getByRole('button', { name: '确认删除' }));
  await waitFor(() => expect(screen.getByText('删除门店失败，请重试。')).toBeTruthy());
  fireEvent.click(screen.getByRole('button', { name: '取消' }));
  await waitFor(() => expect(screen.queryByRole('alertdialog')).toBeNull());
  expect(screen.queryByText('删除门店失败，请重试。')).toBeNull();
  fireEvent.click(screen.getByRole('button', { name: '删除' }));
  expect(within(screen.getByRole('alertdialog')).queryByText('删除门店失败，请重试。')).toBeNull();
});

it('删除请求处理中按 Escape 不关闭确认框，迟到的失败仍在框内显示', async () => {
  let rejectDelete: (reason?: unknown) => void = () => undefined;
  mocks.remove.mockImplementationOnce(() => new Promise<object>((_resolve, reject) => { rejectDelete = reject; }));
  const view = <MemoryRouter><HqStoreListPage /></MemoryRouter>;
  const { rerender } = render(view);
  fireEvent.click(screen.getByRole('button', { name: '删除' }));
  fireEvent.click(screen.getByRole('button', { name: '确认删除' }));
  await waitFor(() => expect(mocks.remove).toHaveBeenCalled());
  mocks.deleteLoading = true;
  rerender(<MemoryRouter><HqStoreListPage /></MemoryRouter>);
  fireEvent.keyDown(document, { key: 'Escape' });
  const remainedOpen = screen.queryByRole('alertdialog') !== null;
  rejectDelete(new Error('offline'));
  mocks.deleteLoading = false;
  rerender(<MemoryRouter><HqStoreListPage /></MemoryRouter>);
  await waitFor(() => expect(screen.getByText('删除门店失败，请重试。')).toBeTruthy());
  expect(remainedOpen).toBe(true);
  expect(within(screen.getByRole('alertdialog')).getByText('删除门店失败，请重试。')).toBeTruthy();
});

it('删除成功后关闭确认框，刷新失败时说明删除已经完成', async () => {
  mocks.remove.mockResolvedValueOnce({});
  render(<MemoryRouter><HqStoreListPage /></MemoryRouter>);
  fireEvent.click(screen.getByRole('button', { name: '删除' }));
  fireEvent.click(screen.getByRole('button', { name: '确认删除' }));
  await waitFor(() => expect(screen.queryByRole('alertdialog')).toBeNull());
  expect(screen.getByText('门店已删除，但门店列表刷新失败，请刷新页面。')).toBeTruthy();
});

it('加盟商目录加载失败时说明筛选不可用并允许重试', () => {
  mocks.catalogError = true;
  render(<MemoryRouter><HqStoreListPage /></MemoryRouter>);
  expect(screen.getByText('加盟商目录加载失败，暂时无法按加盟商筛选。')).toBeTruthy();
  expect(screen.getByRole('button', { name: '重试加盟商目录' })).toBeTruthy();
});
