// @vitest-environment jsdom

import { cleanup, fireEvent, render as renderWithoutToast, screen, waitFor } from '@testing-library/react';
import type { ReactNode } from 'react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { StoreEditorPage } from './StoreEditorPage';
import { storeSaveErrorMessage } from './StoreEditorForm';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';

const render = (ui: ReactNode) => renderWithoutToast(ui, { wrapper: AdminToastProvider });

const mocks = vi.hoisted(() => ({
  search: '', workspace: 'HEADQUARTERS', organizationId: 'hq',
  permissions: ['hqStore:create', 'hqStore:update'],
  loading: false, error: undefined as Error | undefined, record: undefined as Record<string, unknown> | null | undefined,
  queryOptions: [] as Array<Record<string, unknown>>, navigate: vi.fn(), evict: vi.fn(), gc: vi.fn(), refetch: vi.fn(),
  mutations: new Map<string, ReturnType<typeof vi.fn>>(),
}));

vi.mock('@/features/auth/store/authStore', () => ({ useAuthStore: (select: (state: unknown) => unknown) => select({
  viewer: { permissions: mocks.permissions, currentWorkspace: { workspaceType: mocks.workspace, organizationId: mocks.organizationId } },
}) }));
vi.mock('@/features/admin/hooks/useAdminTab', () => ({ useAdminTab: () => ({ search: mocks.search }) }));
vi.mock('react-router', async (importOriginal) => ({ ...(await importOriginal<typeof import('react-router')>()), useNavigate: () => mocks.navigate }));
vi.mock('@apollo/client/react', () => ({
  useApolloClient: () => ({ cache: { evict: mocks.evict, gc: mocks.gc } }),
  useQuery: (_document: unknown, options: Record<string, unknown>) => {
    mocks.queryOptions.push(options);
    return { data: options.skip ? undefined : { store: mocks.record }, loading: mocks.loading, error: mocks.error, refetch: mocks.refetch };
  },
  useMutation: (document: { definitions: Array<{ name?: { value: string } }> }) => {
    const name = document.definitions[0]?.name?.value ?? '';
    return [mocks.mutations.get(name), { loading: false }];
  },
}));

beforeEach(() => {
  mocks.search = ''; mocks.workspace = 'HEADQUARTERS'; mocks.organizationId = 'hq';
  mocks.permissions = ['hqStore:create', 'hqStore:update']; mocks.loading = false; mocks.error = undefined; mocks.record = undefined;
  mocks.queryOptions.length = 0; mocks.navigate.mockReset(); mocks.evict.mockReset(); mocks.gc.mockReset(); mocks.refetch.mockReset();
  mocks.refetch.mockResolvedValue({ data: {} });
  mocks.mutations.clear();
  for (const name of ['HqCreateDirectStore', 'HqUpdateDirectStore', 'FranchiseCreateStore', 'FranchiseUpdateStore']) {
    mocks.mutations.set(name, vi.fn().mockResolvedValue({ data: name.includes('Create')
      ? { createStore: { id: 'new-store' } } : { updateStore: { id: 'store-2' } } }));
  }
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
});
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

function fillRequiredFields() {
  for (const label of [/门店名称/, /联系电话/, /省份/, /城市/, /区县/, /详细地址/]) {
    fireEvent.change(screen.getByLabelText(label), { target: { value: '测试' } });
  }
}

  it('uses one full page for HQ create and invalidates lists on save', async () => {
    render(<StoreEditorPage workspace="HEADQUARTERS" />);
    expect(mocks.queryOptions.at(-1)).toMatchObject({ skip: true, fetchPolicy: 'network-only' });
    expect(screen.getByRole('heading', { name: '新增直营店' })).toBeTruthy();
    const save = screen.getByRole('button', { name: '保存' });
    expect(save.getAttribute('form')).toBe('store-editor-form');
    expect(save.getAttribute('type')).toBe('submit');
    expect(screen.queryByLabelText(/编码|编号/)).toBeNull();
    fillRequiredFields();
    fireEvent.submit(document.getElementById('store-editor-form')!);
    await waitFor(() => expect(mocks.mutations.get('HqCreateDirectStore')).toHaveBeenCalledWith({ variables: { input: expect.objectContaining({ organizationId: 'hq', lifecycle: 'ACTIVE' }) } }));
    expect(mocks.evict).toHaveBeenCalledWith({ fieldName: 'stores' });
    expect(mocks.evict).toHaveBeenCalledWith({ fieldName: 'store' });
    expect(mocks.navigate).toHaveBeenCalledWith('/admin/hq/stores');
  });

  it('keeps the business-status field white for headquarters create and edit forms', () => {
    const view = render(<StoreEditorPage workspace="HEADQUARTERS" />);
    const createTrigger = screen.getByText('营业状态').parentElement?.querySelector('[data-slot="select-trigger"]');
    expect(createTrigger?.hasAttribute('data-admin-form-surface')).toBe(true);
    expect(createTrigger?.className).not.toContain('bg-white');
    expect(createTrigger?.className).not.toContain('text-black');

    mocks.search = '?id=store-2';
    mocks.record = { id: 'store-2', name: '直营店', organizationId: 'hq', lifecycle: 'ACTIVE', organization: { type: 'HEADQUARTERS' } };
    view.rerender(<StoreEditorPage workspace="HEADQUARTERS" />);
    const editTrigger = screen.getByText('营业状态').parentElement?.querySelector('[data-slot="select-trigger"]');
    expect(editTrigger?.hasAttribute('data-admin-form-surface')).toBe(true);
    expect(editTrigger?.className).not.toContain('bg-white');
    expect(editTrigger?.className).not.toContain('text-black');
  });

  it('loads exactly the requested editable franchise store and shows rejection reason', () => {
    mocks.workspace = 'FRANCHISE'; mocks.organizationId = 'tenant'; mocks.permissions = ['store:update']; mocks.search = '?id=store-2';
    mocks.record = { id: 'store-2', name: '门店', organizationId: 'tenant', lifecycle: 'REJECTED', rejectionReason: '地址不完整' };
    render(<StoreEditorPage workspace="FRANCHISE" />);
    expect(mocks.queryOptions.at(-1)).toMatchObject({ variables: { id: 'store-2' }, skip: false, fetchPolicy: 'network-only' });
    expect(screen.getByRole('heading', { name: '编辑门店' })).toBeTruthy();
    expect(screen.getByRole('alert').textContent).toContain('地址不完整');
    for (const label of ['门店名称', '联系电话', '营业时间', '店长姓名', '店长手机号', '省份', '城市', '区县', '详细地址', '营业状态', '服务模式', '经营面积（㎡）', '桌位数（桌）', '小票底部寄语']) {
      expect(document.getElementById('store-editor-form')?.textContent).toContain(label);
    }
  });

  it('creates a franchise draft and updates a matching record with the franchise mutation', async () => {
    mocks.workspace = 'FRANCHISE'; mocks.organizationId = 'tenant'; mocks.permissions = ['store:create', 'store:update'];
    const view = render(<StoreEditorPage workspace="FRANCHISE" />);
    fillRequiredFields();
    fireEvent.submit(document.getElementById('store-editor-form')!);
    await waitFor(() => expect(mocks.mutations.get('FranchiseCreateStore')).toHaveBeenCalledWith({ variables: { input: expect.objectContaining({ organizationId: 'tenant', lifecycle: 'DRAFT' }) } }));
    expect(mocks.navigate).toHaveBeenCalledWith('/admin/franchise/stores');
    mocks.navigate.mockClear();
    mocks.search = '?id=store-2';
    mocks.record = { id: 'store-2', name: '门店', organizationId: 'tenant', lifecycle: 'DRAFT' };
    view.rerender(<StoreEditorPage workspace="FRANCHISE" />);
    fireEvent.submit(document.getElementById('store-editor-form')!);
    await waitFor(() => expect(mocks.mutations.get('FranchiseUpdateStore')).toHaveBeenCalledWith({ variables: { id: 'store-2', input: expect.not.objectContaining({ lifecycle: expect.anything() }) } }));
  });

  it('distinguishes loading, missing, failed requests and denied edits', () => {
    mocks.search = '?id=missing'; mocks.loading = true;
    const view = render(<StoreEditorPage workspace="HEADQUARTERS" />);
    expect(screen.getByRole('status').textContent).toContain('正在加载');
    mocks.loading = false; mocks.record = null; view.rerender(<StoreEditorPage workspace="HEADQUARTERS" />);
    expect(screen.getByRole('alert').textContent).toContain('不存在或无权');
    mocks.error = new Error('offline'); view.rerender(<StoreEditorPage workspace="HEADQUARTERS" />);
    fireEvent.click(screen.getByRole('button', { name: '重试' }));
    expect(mocks.refetch).toHaveBeenCalled();
    mocks.error = undefined; mocks.record = { id: 'missing', name: '店', organizationId: 'tenant', lifecycle: 'ACTIVE', organization: { type: 'FRANCHISE' } };
    view.rerender(<StoreEditorPage workspace="HEADQUARTERS" />);
    expect(screen.getByRole('alert').textContent).toContain('不可编辑');
    expect(screen.queryByRole('button', { name: '保存' })).toBeNull();
  });

  it('uses scoped-denial missing state and resets edited values when ID changes', () => {
    mocks.search = '?id=a';
    mocks.error = Object.assign(new Error('denied'), { errors: [{ extensions: { code: 'PERMISSION_DENIED' } }] });
    const view = render(<StoreEditorPage workspace="HEADQUARTERS" />);
    expect(screen.getByRole('alert').textContent).toContain('不存在或无权访问');
    expect(screen.queryByRole('button', { name: '重试' })).toBeNull();
    mocks.error = undefined;
    mocks.record = { id: 'a', name: '原店', organizationId: 'hq', lifecycle: 'ACTIVE', organization: { type: 'HEADQUARTERS' } };
    view.rerender(<StoreEditorPage workspace="HEADQUARTERS" />);
    fireEvent.change(screen.getByLabelText(/门店名称/), { target: { value: '未保存的修改' } });
    mocks.search = '?id=b';
    mocks.record = { id: 'b', name: '新店', organizationId: 'hq', lifecycle: 'ACTIVE', organization: { type: 'HEADQUARTERS' } };
    view.rerender(<StoreEditorPage workspace="HEADQUARTERS" />);
    expect((screen.getByLabelText(/门店名称/) as HTMLInputElement).value).toBe('新店');
  });

  it('keeps a failed save on the form without cache eviction or navigation', async () => {
    mocks.mutations.get('HqCreateDirectStore')?.mockRejectedValueOnce(Object.assign(new Error('raw secret'), { errors: [{ extensions: { code: 'CONFLICT' } }] }));
    render(<StoreEditorPage workspace="HEADQUARTERS" />);
    fillRequiredFields();
    fireEvent.submit(document.getElementById('store-editor-form')!);
    await waitFor(() => expect(screen.getByRole('alert').textContent).toContain('状态已变化'));
    expect(screen.getByRole('alert').textContent).not.toContain('raw secret');
    expect(mocks.evict).not.toHaveBeenCalled();
    expect(mocks.navigate).not.toHaveBeenCalled();
  });

  it('does not silently replace cleared required business hours when saving', () => {
    mocks.search = '?id=store-2';
    mocks.record = { id: 'store-2', name: '门店', organizationId: 'hq', organization: { type: 'HEADQUARTERS' }, lifecycle: 'ACTIVE', businessHours: '10:00 - 20:00' };
    render(<StoreEditorPage workspace="HEADQUARTERS" />);
    fillRequiredFields();
    fireEvent.click(screen.getByRole('button', { name: '清空时间' }));
    fireEvent.submit(document.getElementById('store-editor-form')!);
    expect(mocks.mutations.get('HqUpdateDirectStore')).not.toHaveBeenCalled();
    expect(screen.getByRole('alert').textContent).toContain('请选择营业时间');
  });

  it('blocks wrong workspace and missing write permission', () => {
    mocks.workspace = 'FRANCHISE';
    const view = render(<StoreEditorPage workspace="HEADQUARTERS" />);
    expect(screen.getByRole('alert').textContent).toContain('无权');
    expect(mocks.queryOptions).toHaveLength(0);
    mocks.workspace = 'HEADQUARTERS'; mocks.permissions = [];
    view.rerender(<StoreEditorPage workspace="HEADQUARTERS" />);
    expect(screen.getByRole('alert').textContent).toContain('无权');
  });
it('maps store save errors to safe Chinese messages', () => {
  expect(storeSaveErrorMessage('PERMISSION_DENIED')).toContain('无权');
  expect(storeSaveErrorMessage('STORE_SCOPE_DENIED')).toContain('无权');
  expect(storeSaveErrorMessage('CONFLICT')).toContain('状态已变化');
  expect(storeSaveErrorMessage('VALIDATION_FAILED')).toContain('不符合要求');
  expect(storeSaveErrorMessage(undefined)).toBe('门店保存失败，请重试。');
});

it('rejects an empty edit ID instead of opening a create form', () => {
  mocks.search = '?id=';
  render(<StoreEditorPage workspace="HEADQUARTERS" />);
  expect(screen.getByRole('alert').textContent).toContain('链接无效');
  expect(screen.queryByRole('button', { name: '保存' })).toBeNull();
  expect(mocks.queryOptions).toHaveLength(0);
});

it('sends headquarters edits through the existing HQ update mutation', async () => {
  mocks.search = '?id=hq-1';
  mocks.record = { id: 'hq-1', name: '直营店', organizationId: 'hq', lifecycle: 'ACTIVE', organization: { type: 'HEADQUARTERS' } };
  render(<StoreEditorPage workspace="HEADQUARTERS" />);
  fireEvent.submit(document.getElementById('store-editor-form')!);
  await waitFor(() => expect(mocks.mutations.get('HqUpdateDirectStore')).toHaveBeenCalledWith({
    variables: { id: 'hq-1', input: expect.objectContaining({ name: '直营店' }) },
  }));
  expect(mocks.mutations.get('HqCreateDirectStore')).not.toHaveBeenCalled();
});
