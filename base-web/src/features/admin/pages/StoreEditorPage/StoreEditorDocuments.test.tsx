// @vitest-environment jsdom

import { act, cleanup, fireEvent, render as renderWithoutToast, screen, waitFor } from '@testing-library/react';
import type { ReactNode } from 'react';
import { afterEach, expect, it, vi } from 'vitest';
import type { AdminStoreEditorQuery } from '@/__generated__/graphql';
import { StoreEditorForm } from './StoreEditorForm';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';

const render = (ui: ReactNode) => renderWithoutToast(ui, { wrapper: AdminToastProvider });

const mocks = vi.hoisted(() => ({
  create: vi.fn(), update: vi.fn(), bind: vi.fn(), remove: vi.fn(), upload: vi.fn(), saved: vi.fn(),
}));

vi.mock('@apollo/client/react', () => ({ useMutation: (document: { definitions: Array<{ name?: { value: string } }> }) => {
  const name = document.definitions[0]?.name?.value;
  return [name === 'HqCreateDirectStore' ? mocks.create : name === 'SetStoreDocument' ? mocks.bind
    : name === 'RemoveStoreDocument' ? mocks.remove : mocks.update, { loading: false }];
} }));
vi.mock('./storeDocumentImages', async (importOriginal) => ({
  ...(await importOriginal<typeof import('./storeDocumentImages')>()), uploadStoreDocument: mocks.upload,
}));

afterEach(() => { cleanup(); vi.useRealTimers(); vi.resetAllMocks(); vi.unstubAllGlobals(); });

function fillRequiredFields() {
  for (const label of [/门店名称/, /联系电话/, /省份/, /城市/, /区县/, /详细地址/]) {
    fireEvent.change(screen.getByLabelText(label), { target: { value: '测试' } });
  }
}

it('shows the same save failure again after the previous toast disappears', async () => {
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
  mocks.create.mockRejectedValue(new Error('offline'));
  render(<StoreEditorForm workspace="HEADQUARTERS" organizationId="hq" onBack={vi.fn()} onSaved={mocks.saved} />);
  fillRequiredFields();
  vi.useFakeTimers();
  await act(async () => { fireEvent.submit(document.getElementById('store-editor-form')!); });
  expect(screen.getByRole('alert').textContent).toBe('门店保存失败，请重试');
  act(() => { vi.advanceTimersByTime(3000); });
  expect(screen.queryByRole('alert')).toBeNull();
  await act(async () => { fireEvent.submit(document.getElementById('store-editor-form')!); });
  expect(screen.getByRole('alert').textContent).toBe('门店保存失败，请重试');
  expect(mocks.create).toHaveBeenCalledTimes(2);
  expect(mocks.saved).not.toHaveBeenCalled();
});

it('keeps an incomplete create response visible without reporting success or creating again', async () => {
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
  mocks.create.mockResolvedValue({ data: { createStore: null } });
  render(<StoreEditorForm workspace="HEADQUARTERS" organizationId="hq" onBack={vi.fn()} onSaved={mocks.saved} />);
  fillRequiredFields();
  fireEvent.submit(document.getElementById('store-editor-form')!);
  await waitFor(() => expect(screen.getByRole('alert').textContent).toContain('未获得门店编号'));
  expect(mocks.saved).not.toHaveBeenCalled();
  fireEvent.submit(document.getElementById('store-editor-form')!);
  await waitFor(() => expect((screen.getByRole('button', { name: '保存' }) as HTMLButtonElement).disabled).toBe(false));
  expect(mocks.create).toHaveBeenCalledTimes(1);
  expect(mocks.saved).not.toHaveBeenCalled();
});

it('shows two document slots and retries a partial save without another store creation', async () => {
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
  URL.createObjectURL = vi.fn(() => 'blob:preview'); URL.revokeObjectURL = vi.fn();
  mocks.create.mockResolvedValue({ data: { createStore: { id: 'store-1' } } });
  mocks.upload.mockResolvedValueOnce('attachment-1').mockResolvedValueOnce('attachment-2');
  mocks.bind.mockResolvedValueOnce({ data: {} }).mockRejectedValueOnce(new Error('offline')).mockResolvedValueOnce({ data: {} });
  render(<StoreEditorForm workspace="HEADQUARTERS" organizationId="hq" onBack={vi.fn()} onSaved={mocks.saved} />);
  expect(screen.getByLabelText('上传营业执照')).toBeTruthy();
  expect(screen.getByLabelText('上传其他')).toBeTruthy();
  fireEvent.change(screen.getByLabelText('上传营业执照'), { target: { files: [new File(['a'], 'a.png', { type: 'image/png' })] } });
  fireEvent.change(screen.getByLabelText('上传其他'), { target: { files: [new File(['b'], 'b.png', { type: 'image/png' })] } });
  fillRequiredFields();
  fireEvent.submit(document.getElementById('store-editor-form')!);
  await waitFor(() => expect(screen.getByText(/门店已保存，证照未完成/)).toBeTruthy());
  expect(mocks.create).toHaveBeenCalledTimes(1);
  expect(mocks.bind).toHaveBeenCalledTimes(2);
  fireEvent.submit(document.getElementById('store-editor-form')!);
  await waitFor(() => expect(mocks.saved).toHaveBeenCalledTimes(1));
  expect(mocks.create).toHaveBeenCalledTimes(1);
  expect(mocks.bind).toHaveBeenCalledTimes(3);
});

it('shows an upload failure in its document card and retries from that card', async () => {
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
  URL.createObjectURL = vi.fn(() => 'blob:preview'); URL.revokeObjectURL = vi.fn();
  mocks.upload.mockRejectedValueOnce(new Error('证照上传失败，请重试。')).mockResolvedValueOnce('attachment-1');
  mocks.create.mockResolvedValue({ data: { createStore: { id: 'store-1' } } });
  mocks.bind.mockResolvedValue({ data: {} });
  render(<StoreEditorForm workspace="HEADQUARTERS" organizationId="hq" onBack={vi.fn()} onSaved={mocks.saved} />);
  fireEvent.change(screen.getByLabelText('上传营业执照'), { target: { files: [new File(['a'], 'a.png', { type: 'image/png' })] } });
  fillRequiredFields();
  fireEvent.submit(document.getElementById('store-editor-form')!);
  await waitFor(() => expect(screen.getByRole('button', { name: '重试保存营业执照' })).toBeTruthy());
  expect(screen.getByText('证照上传失败，请重试。').closest('section')?.textContent).toContain('营业执照');
  expect(screen.getByText('证照上传失败，请重试').closest('[role="alert"]')).not.toBeNull();
  expect(mocks.create).not.toHaveBeenCalled();
  fireEvent.click(screen.getByRole('button', { name: '重试保存营业执照' }));
  await waitFor(() => expect(mocks.saved).toHaveBeenCalledTimes(1));
  expect(mocks.create).toHaveBeenCalledTimes(1);
});

it('replaces only the unfinished document after a partial save', async () => {
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
  URL.createObjectURL = vi.fn(() => 'blob:preview'); URL.revokeObjectURL = vi.fn();
  mocks.create.mockResolvedValue({ data: { createStore: { id: 'store-1' } } });
  mocks.upload.mockResolvedValueOnce('attachment-1').mockResolvedValueOnce('attachment-old').mockResolvedValueOnce('attachment-new');
  mocks.bind.mockResolvedValueOnce({ data: {} }).mockRejectedValueOnce(new Error('offline')).mockResolvedValueOnce({ data: {} });
  render(<StoreEditorForm workspace="HEADQUARTERS" organizationId="hq" onBack={vi.fn()} onSaved={mocks.saved} />);
  fireEvent.change(screen.getByLabelText('上传营业执照'), { target: { files: [new File(['a'], 'a.png', { type: 'image/png' })] } });
  fireEvent.change(screen.getByLabelText('上传其他'), { target: { files: [new File(['b'], 'b.png', { type: 'image/png' })] } });
  fillRequiredFields();
  fireEvent.submit(document.getElementById('store-editor-form')!);
  await waitFor(() => expect(screen.getByText(/门店已保存，证照未完成/)).toBeTruthy());
  expect((screen.getByLabelText('上传营业执照') as HTMLInputElement).disabled).toBe(true);
  expect((screen.getByLabelText('上传其他') as HTMLInputElement).disabled).toBe(false);
  fireEvent.change(screen.getByLabelText('上传其他'), { target: { files: [new File(['c'], 'c.png', { type: 'image/png' })] } });
  fireEvent.submit(document.getElementById('store-editor-form')!);
  await waitFor(() => expect(mocks.saved).toHaveBeenCalledTimes(1));
  expect(mocks.upload).toHaveBeenCalledTimes(3);
  expect(mocks.create).toHaveBeenCalledTimes(1);
  expect(mocks.bind.mock.calls[2]?.[0].variables.attachmentId).toBe('attachment-new');
});

it('reuploads an expired attachment and binds it without creating a second store', async () => {
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
  URL.createObjectURL = vi.fn(() => 'blob:preview'); URL.revokeObjectURL = vi.fn();
  mocks.create.mockResolvedValue({ data: { createStore: { id: 'store-1' } } });
  mocks.upload.mockResolvedValueOnce('expired').mockResolvedValueOnce('fresh');
  mocks.bind.mockRejectedValueOnce({ extensions: { code: 'VALIDATION_FAILED' } }).mockResolvedValueOnce({ data: {} });
  render(<StoreEditorForm workspace="HEADQUARTERS" organizationId="hq" onBack={vi.fn()} onSaved={mocks.saved} />);
  fireEvent.change(screen.getByLabelText('上传营业执照'), { target: { files: [new File(['a'], 'a.png', { type: 'image/png' })] } });
  fillRequiredFields();
  fireEvent.submit(document.getElementById('store-editor-form')!);
  await waitFor(() => expect(mocks.saved).toHaveBeenCalledTimes(1));
  expect(mocks.create).toHaveBeenCalledTimes(1);
  expect(mocks.upload).toHaveBeenCalledTimes(2);
  expect(mocks.bind).toHaveBeenCalledTimes(2);
});

it('keeps an unfinished document editable when its recovery upload fails', async () => {
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
  URL.createObjectURL = vi.fn(() => 'blob:preview'); URL.revokeObjectURL = vi.fn();
  mocks.create.mockResolvedValue({ data: { createStore: { id: 'store-1' } } });
  mocks.upload.mockResolvedValueOnce('expired').mockRejectedValueOnce(new Error('证照上传失败，请重试。'));
  mocks.bind.mockRejectedValueOnce({ extensions: { code: 'VALIDATION_FAILED' } });
  render(<StoreEditorForm workspace="HEADQUARTERS" organizationId="hq" onBack={vi.fn()} onSaved={mocks.saved} />);
  fireEvent.change(screen.getByLabelText('上传营业执照'), { target: { files: [new File(['a'], 'a.png', { type: 'image/png' })] } });
  fillRequiredFields();
  fireEvent.submit(document.getElementById('store-editor-form')!);
  await waitFor(() => expect(screen.getByText(/门店已保存，证照未完成/)).toBeTruthy());
  expect(screen.getByText('证照上传失败，请重试。')).toBeTruthy();
  expect((screen.getByLabelText('上传营业执照') as HTMLInputElement).disabled).toBe(false);
  expect(mocks.create).toHaveBeenCalledTimes(1);
});

it('lets a saved store finish without an optional document that failed to bind', async () => {
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
  URL.createObjectURL = vi.fn(() => 'blob:preview'); URL.revokeObjectURL = vi.fn();
  mocks.create.mockResolvedValue({ data: { createStore: { id: 'store-1' } } });
  mocks.upload.mockResolvedValue('attachment-1');
  mocks.bind.mockRejectedValueOnce(new Error('offline'));
  render(<StoreEditorForm workspace="HEADQUARTERS" organizationId="hq" onBack={vi.fn()} onSaved={mocks.saved} />);
  fireEvent.change(screen.getByLabelText('上传营业执照'), { target: { files: [new File(['a'], 'a.png', { type: 'image/png' })] } });
  fillRequiredFields();
  fireEvent.submit(document.getElementById('store-editor-form')!);
  await waitFor(() => expect(screen.getByText(/门店已保存，证照未完成/)).toBeTruthy());
  fireEvent.click(screen.getByRole('button', { name: '删除营业执照' }));
  fireEvent.submit(document.getElementById('store-editor-form')!);
  await waitFor(() => expect(mocks.saved).toHaveBeenCalledTimes(1));
  expect(mocks.create).toHaveBeenCalledTimes(1);
  expect(mocks.bind).toHaveBeenCalledTimes(1);
});

it('renders the same two slots when a franchise edits an existing store', () => {
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, blob: async () => new Blob(['image'], { type: 'image/png' }) }));
  URL.createObjectURL = vi.fn(() => 'blob:saved'); URL.revokeObjectURL = vi.fn();
  const existing = { id: 'store-1', lifecycle: 'DRAFT', organizationId: 'franchise-1',
    businessLicenseImageUrl: '/uploads/stores/store-1/license.png', otherDocumentImageUrl: null,
  } as NonNullable<AdminStoreEditorQuery['store']>;
  render(<StoreEditorForm workspace="FRANCHISE" organizationId="franchise-1" existing={existing}
    onBack={vi.fn()} onSaved={mocks.saved} />);
  expect(screen.getByLabelText('上传营业执照')).toBeTruthy();
  expect(screen.getByLabelText('上传其他')).toBeTruthy();
  expect(screen.getByRole('button', { name: '更换营业执照图片' })).toBeTruthy();
});
