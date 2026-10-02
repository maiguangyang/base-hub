// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { StoreDocumentCard } from './StoreDocumentCard';

afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

it('selects and previews a business license image', () => {
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
  URL.createObjectURL = vi.fn(() => 'blob:local-preview');
  URL.revokeObjectURL = vi.fn();
  const onChange = vi.fn();
  const file = new File(['image'], 'license.png', { type: 'image/png' });
  const view = render(<StoreDocumentCard label="营业执照" change={{}} onChange={onChange} />);
  const openPicker = vi.spyOn(screen.getByLabelText('上传营业执照'), 'click');
  fireEvent.click(screen.getByRole('button', { name: '选择营业执照图片' }));
  expect(openPicker).toHaveBeenCalledTimes(1);
  expect(screen.queryByText('上传营业执照')).toBeNull();
  fireEvent.change(screen.getByLabelText('上传营业执照'), { target: { files: [file] } });
  expect(onChange).toHaveBeenCalledWith({ file });
  view.rerender(<StoreDocumentCard label="营业执照" change={{ file }} onChange={onChange} />);
  expect(screen.getByRole('button', { name: '更换营业执照图片' })).toBeTruthy();
  expect(screen.getByRole('button', { name: '删除营业执照' })).toBeTruthy();
  fireEvent.click(screen.getByRole('button', { name: '预览营业执照' }));
  expect(screen.getByRole('dialog').querySelector('img')?.getAttribute('src')).toBe('blob:local-preview');
});

it('rejects unsupported and oversized files without changing the draft', () => {
  const onChange = vi.fn();
  render(<StoreDocumentCard label="其他" change={{}} onChange={onChange} />);
  fireEvent.change(screen.getByLabelText('上传其他'), { target: { files: [new File(['x'], 'bad.txt', { type: 'text/plain' })] } });
  expect(onChange).not.toHaveBeenCalled();
  expect(screen.getByRole('alert').textContent).toContain('JPG');
});

it('clears an unsaved selection without requesting a server removal', () => {
  URL.createObjectURL = vi.fn(() => 'blob:local-preview'); URL.revokeObjectURL = vi.fn();
  const onChange = vi.fn();
  const file = new File(['image'], 'license.png', { type: 'image/png' });
  render(<StoreDocumentCard label="营业执照" change={{ file }} onChange={onChange} />);
  fireEvent.click(screen.getByRole('button', { name: '删除营业执照' }));
  expect(onChange).toHaveBeenCalledWith({});
});

it('accepts a valid image dropped on the preview area', () => {
  const onChange = vi.fn();
  const file = new File(['image'], 'other.png', { type: 'image/png' });
  render(<StoreDocumentCard label="其他" change={{}} onChange={onChange} />);
  fireEvent.drop(screen.getByRole('button', { name: '选择其他图片' }), { dataTransfer: { files: [file] } });
  expect(onChange).toHaveBeenCalledWith({ file });
});

it('keeps the hidden file input out of keyboard tab order', () => {
  render(<StoreDocumentCard label="其他" change={{}} onChange={vi.fn()} />);
  expect((screen.getByLabelText('上传其他') as HTMLInputElement).tabIndex).toBe(-1);
  expect(screen.getByRole('button', { name: '选择其他图片' }).getAttribute('tabindex')).toBeNull();
});

it('keeps drag highlighting while moving across children inside the upload area', () => {
  render(<StoreDocumentCard label="其他" change={{}} onChange={vi.fn()} />);
  const button = screen.getByRole('button', { name: '选择其他图片' });
  const area = button.parentElement!;
  fireEvent.dragEnter(area);
  fireEvent.dragEnter(button);
  fireEvent.dragLeave(button);
  expect(area.className).toContain('border-primary');
  fireEvent.dragLeave(area);
  expect(area.className).not.toContain('border-primary');
});

it('rejects a dropped file while the upload area is disabled', () => {
  const onChange = vi.fn();
  const file = new File(['image'], 'other.png', { type: 'image/png' });
  render(<StoreDocumentCard label="其他" change={{}} onChange={onChange} disabled />);
  fireEvent.drop(screen.getByRole('button', { name: '选择其他图片' }), { dataTransfer: { files: [file] } });
  expect(onChange).not.toHaveBeenCalled();
});

it('loads a saved image through an authenticated request and opens its preview', async () => {
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
  URL.createObjectURL = vi.fn(() => 'blob:persisted'); URL.revokeObjectURL = vi.fn();
  const fetchImage = vi.fn().mockResolvedValue({ ok: true, blob: async () => new Blob(['image'], { type: 'image/png' }) });
  vi.stubGlobal('fetch', fetchImage);
  render(<StoreDocumentCard label="营业执照" storedPath="/uploads/stores/store-1/license.png" change={{}} onChange={vi.fn()} />);
  await waitFor(() => expect(screen.getByRole('button', { name: '预览营业执照' }).hasAttribute('disabled')).toBe(false));
  expect(fetchImage.mock.calls[0]?.[1]).toMatchObject({ credentials: 'include', cache: 'no-store' });
  fireEvent.click(screen.getByRole('button', { name: '预览营业执照' }));
  expect(screen.getByRole('dialog').querySelector('img')?.getAttribute('src')).toBe('blob:persisted');
});

it('retries a failed saved-image preview', async () => {
  URL.createObjectURL = vi.fn(() => 'blob:recovered'); URL.revokeObjectURL = vi.fn();
  const fetchImage = vi.fn().mockRejectedValueOnce(new Error('offline'))
    .mockResolvedValueOnce({ ok: true, blob: async () => new Blob(['image'], { type: 'image/png' }) });
  vi.stubGlobal('fetch', fetchImage);
  render(<StoreDocumentCard label="其他" storedPath="/uploads/stores/store-1/other.png" change={{}} onChange={vi.fn()} />);
  await waitFor(() => expect(screen.getByRole('alert').textContent).toContain('预览加载失败'));
  expect(screen.queryByText('正在加载预览')).toBeNull();
  fireEvent.click(screen.getByRole('button', { name: '重试预览' }));
  await waitFor(() => expect(screen.getByRole('button', { name: '预览其他' }).hasAttribute('disabled')).toBe(false));
  expect(fetchImage).toHaveBeenCalledTimes(2);
});
