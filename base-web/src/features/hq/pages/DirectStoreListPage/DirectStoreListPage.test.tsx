// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { DirectStoreListPage } from './DirectStoreListPage';

const mocks = vi.hoisted(() => ({ navigate: vi.fn(), refetch: vi.fn().mockResolvedValue({}), deleteStore: vi.fn().mockResolvedValue({}) }));
vi.mock('react-router', async (importOriginal) => ({ ...(await importOriginal<typeof import('react-router')>()), useNavigate: () => mocks.navigate }));
vi.mock('@/features/admin/hooks/useAdminTab', () => ({ useAdminTab: () => ({ search: '' }) }));
vi.mock('@/features/auth/store/authStore', () => ({ useAuthStore: (select: (state: unknown) => unknown) => select({
  viewer: { currentWorkspace: { workspaceType: 'HEADQUARTERS', organizationId: 'hq' },
    permissions: ['hqStore:create', 'hqStore:update', 'hqStore:delete'] },
}) }));
vi.mock('@apollo/client/react', () => ({
  useQuery: () => ({ data: { stores: { data: [{ id: 'direct', name: '直营店', code: 'STR1', lifecycle: 'ACTIVE',
    organizationId: 'hq', businessStatus: 'OPEN' }], total: 1 } }, loading: false, error: undefined, refetch: mocks.refetch }),
  useMutation: (document: { definitions: Array<{ name?: { value?: string } }> }) => [
    document.definitions[0]?.name?.value === 'HqDeleteDirectStores' ? mocks.deleteStore : vi.fn().mockResolvedValue({}), { loading: false },
  ],
}));

beforeEach(() => { vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }); vi.clearAllMocks(); });
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

it.each([['新增', '/admin/hq/stores/manage'], ['编辑', '/admin/hq/stores/manage?id=direct']])('%s进入共享门店页', (label, path) => {
  render(<DirectStoreListPage />);
  fireEvent.click(screen.getByRole('button', { name: label }));
  expect(mocks.navigate).toHaveBeenCalledWith(path);
  expect(screen.queryByRole('dialog')).toBeNull();
});

it('删除直营店时先确认，取消不发送请求，失败则保留确认框', async () => {
  mocks.deleteStore.mockRejectedValueOnce(new Error('offline'));
  render(<DirectStoreListPage />);
  fireEvent.click(screen.getByRole('button', { name: '删除' }));
  expect(screen.getByRole('alertdialog')).toBeTruthy();
  fireEvent.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '取消' }));
  expect(mocks.deleteStore).not.toHaveBeenCalled();
  fireEvent.click(screen.getByRole('button', { name: '删除' }));
  fireEvent.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '确认删除' }));
  await waitFor(() => expect(within(screen.getByRole('alertdialog')).getByRole('alert').textContent).toContain('操作失败'));
  expect(mocks.deleteStore).toHaveBeenCalledTimes(1);
});

it('删除成功后刷新失败仍关闭确认框，不提供再次提交', async () => {
  mocks.refetch.mockRejectedValueOnce(new Error('offline'));
  render(<DirectStoreListPage />);
  fireEvent.click(screen.getByRole('button', { name: '删除' }));
  fireEvent.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '确认删除' }));
  await waitFor(() => expect(screen.queryByRole('alertdialog')).toBeNull());
  expect(screen.getByText('门店已删除，但列表刷新失败，请刷新页面。')).toBeTruthy();
  expect(mocks.deleteStore).toHaveBeenCalledTimes(1);
});
