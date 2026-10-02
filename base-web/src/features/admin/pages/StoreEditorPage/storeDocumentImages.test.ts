import { afterEach, expect, it, vi } from 'vitest';
import { uploadStoreDocument } from './storeDocumentImages';

afterEach(() => vi.unstubAllGlobals());

it('asks for another image when the server rejects the image contents', async () => {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false, status: 400 }));
  const image = new File(['invalid image'], 'license.png', { type: 'image/png' });
  await expect(uploadStoreDocument(image)).rejects.toThrow('图片内容无效，请更换 JPG、PNG 或 WebP 图片。');
});

it('shows a retryable upload message for a network failure', async () => {
  vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('Failed to fetch')));
  const image = new File(['image'], 'license.png', { type: 'image/png' });
  await expect(uploadStoreDocument(image)).rejects.toThrow('证照上传失败，请重试。');
});
