// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { StoreListPage } from './StoreListPage';

const mocks = vi.hoisted(() => ({
  updateStore: vi.fn(async () => ({ data: {} })),
  refetch: vi.fn(async () => undefined),
  navigate: vi.fn(),
  deleteStore: vi.fn(async () => ({ data: {} })),
  submitStore: vi.fn(async () => ({ data: {} })),
  permissions: ['store:read', 'store:update', 'store:create'] as string[],
  lifecycle: 'DRAFT',
}));
vi.mock('react-router', async (importOriginal) => ({ ...(await importOriginal<typeof import('react-router')>()), useNavigate: () => mocks.navigate }));

vi.mock('@apollo/client/react', () => ({
  useQuery: () => ({
    data: { stores: { total: 1, data: [{
      id: 'store-1', code: 'STR1', name: '加盟店', lifecycle: mocks.lifecycle, organizationId: 'org-1',
      contactPhone: '020-12345678', managerName: '店长', managerPhone: '13800000000',
      province: '广东省', city: '广州市', district: '海珠区', address: '江南大道1号',
      businessHours: '10:00 - 21:00', businessStatus: 'OPEN', supportDineIn: true, supportTakeout: false,
      storeArea: 86.5, tableCount: 12, receiptFooter: '欢迎光临',
    }] } },
    loading: false, error: undefined, refetch: mocks.refetch,
  }),
  useMutation: (document: { definitions?: Array<{ name?: { value?: string } }> }) => {
    const name = document.definitions?.find((definition) => definition.name)?.name?.value;
    return [name === 'FranchiseDeleteStores' ? mocks.deleteStore : name === 'FranchiseSubmitStore' ? mocks.submitStore : mocks.updateStore, { loading: false }];
  },
}));

vi.mock('@/features/admin/hooks/useAdminTab', () => ({ useAdminTab: () => ({ search: '' }) }));
vi.mock('@/features/auth/store/authStore', () => ({
  useAuthStore: (selector: (state: unknown) => unknown) => selector({
    viewer: { currentWorkspace: { workspaceType: 'FRANCHISE', organizationId: 'org-1' }, permissions: mocks.permissions },
  }),
}));

beforeEach(() => { vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }); vi.clearAllMocks(); mocks.permissions = ['store:read', 'store:update', 'store:create']; mocks.lifecycle = 'DRAFT'; });
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

it.each([['新增', '/admin/franchise/stores/manage'], ['编辑', '/admin/franchise/stores/manage?id=store-1']])('%s进入共用门店页面', (action, path) => {
  render(<StoreListPage />);
  fireEvent.click(screen.getByRole('button', { name: action }));
  expect(mocks.navigate).toHaveBeenCalledWith(path);
  expect(screen.queryByRole('dialog')).toBeNull();
});

it('草稿可提交审核，待审核状态没有编辑与提交按钮', async () => {
  mocks.permissions.push('store:submit');
  const view = render(<StoreListPage />);
  fireEvent.click(screen.getByRole('button', { name: '提交审核' }));
  expect(mocks.submitStore).not.toHaveBeenCalled();
  expect(within(screen.getByRole('alertdialog')).getByText(/确认提交“加盟店”/)).toBeTruthy();
  fireEvent.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '确认提交' }));
  await waitFor(() => expect(mocks.submitStore).toHaveBeenCalledWith({ variables: { id: 'store-1' } }));
  await waitFor(() => expect(screen.queryByRole('alertdialog')).toBeNull());
  mocks.lifecycle = 'PENDING_APPROVAL';
  view.rerender(<StoreListPage />);
  expect(screen.queryByRole('button', { name: '编辑' })).toBeNull();
  expect(screen.queryByRole('button', { name: '提交审核' })).toBeNull();
});

it.each(['DRAFT', 'REJECTED'])('%s 门店取消审核确认时不提交', (lifecycle) => {
  mocks.permissions.push('store:submit');
  mocks.lifecycle = lifecycle;
  render(<StoreListPage />);
  fireEvent.click(screen.getByRole('button', { name: '提交审核' }));
  fireEvent.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '取消' }));
  expect(screen.queryByRole('alertdialog')).toBeNull();
  expect(mocks.submitStore).not.toHaveBeenCalled();
  expect(mocks.refetch).not.toHaveBeenCalled();
});

it('审核提交失败时保留确认框，可重试', async () => {
  mocks.permissions.push('store:submit');
  mocks.submitStore.mockRejectedValueOnce(new Error('offline'));
  render(<StoreListPage />);
  fireEvent.click(screen.getByRole('button', { name: '提交审核' }));
  fireEvent.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '确认提交' }));
  await waitFor(() => expect(within(screen.getByRole('alertdialog')).getByRole('alert').textContent).toContain('操作失败'));
  expect(mocks.refetch).not.toHaveBeenCalled();
  fireEvent.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '确认提交' }));
  await waitFor(() => expect(screen.queryByRole('alertdialog')).toBeNull());
  expect(mocks.submitStore).toHaveBeenCalledTimes(2);
});

it('提交期间锁定确认和取消，避免重复请求', async () => {
  mocks.permissions.push('store:submit');
  let resolveSubmit!: (value: { data: object }) => void;
  mocks.submitStore.mockImplementationOnce(() => new Promise((resolve) => { resolveSubmit = resolve; }));
  render(<StoreListPage />);
  fireEvent.click(screen.getByRole('button', { name: '提交审核' }));
  const dialog = screen.getByRole('alertdialog');
  fireEvent.click(within(dialog).getByRole('button', { name: '确认提交' }));
  expect((within(dialog).getByRole('button', { name: '提交中…' }) as HTMLButtonElement).disabled).toBe(true);
  expect((within(dialog).getByRole('button', { name: '取消' }) as HTMLButtonElement).disabled).toBe(true);
  fireEvent.click(within(dialog).getByRole('button', { name: '提交中…' }));
  expect(mocks.submitStore).toHaveBeenCalledTimes(1);
  resolveSubmit({ data: {} });
  await waitFor(() => expect(screen.queryByRole('alertdialog')).toBeNull());
});

it('提交成功但刷新失败时关闭确认框并提示，避免再次提交', async () => {
  mocks.permissions.push('store:submit');
  mocks.refetch.mockRejectedValueOnce(new Error('offline'));
  render(<StoreListPage />);
  fireEvent.click(screen.getByRole('button', { name: '提交审核' }));
  fireEvent.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '确认提交' }));
  await waitFor(() => expect(screen.queryByRole('alertdialog')).toBeNull());
  expect(screen.getByText('门店已提交审核，但列表刷新失败，请刷新页面。')).toBeTruthy();
  expect(mocks.submitStore).toHaveBeenCalledTimes(1);
});

it('删除前确认，取消时不发送请求；失败时保留确认框', async () => {
  mocks.permissions.push('store:delete');
  mocks.deleteStore.mockRejectedValueOnce(new Error('offline'));
  render(<StoreListPage />);
  fireEvent.click(screen.getByRole('button', { name: '删除' }));
  expect(screen.getByRole('alertdialog')).toBeTruthy();
  fireEvent.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '取消' }));
  expect(mocks.deleteStore).not.toHaveBeenCalled();
  fireEvent.click(screen.getByRole('button', { name: '删除' }));
  fireEvent.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '确认删除' }));
  await waitFor(() => expect(within(screen.getByRole('alertdialog')).getByRole('alert').textContent).toContain('操作失败'));
  expect(mocks.deleteStore).toHaveBeenCalledTimes(1);
});

it('删除成功后关闭确认框，即使列表刷新失败也不会再次删除', async () => {
  mocks.permissions.push('store:delete');
  mocks.refetch.mockRejectedValueOnce(new Error('offline'));
  render(<StoreListPage />);
  fireEvent.click(screen.getByRole('button', { name: '删除' }));
  fireEvent.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '确认删除' }));
  await waitFor(() => expect(screen.queryByRole('alertdialog')).toBeNull());
  expect(screen.getByText('门店已删除，但列表刷新失败，请刷新页面。')).toBeTruthy();
  expect(mocks.deleteStore).toHaveBeenCalledTimes(1);
});
