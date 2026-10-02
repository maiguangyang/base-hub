import { appConfig } from '@/config/app';

const IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/webp'];
const MAX_BYTES = 5 << 20;

export function validateStoreDocument(file: File): string | undefined {
  if (!IMAGE_TYPES.includes(file.type) || file.size === 0 || file.size > MAX_BYTES) {
    return '请选择不超过 5 MB 的 JPG、PNG 或 WebP 图片。';
  }
  return undefined;
}

export async function uploadStoreDocument(file: File): Promise<string> {
  const error = validateStoreDocument(file);
  if (error) throw new Error(error);
  let response: Response;
  try {
    response = await fetch(`${appConfig.engine.httpUrl}/api/store-documents`, {
      method: 'POST', credentials: 'include', cache: 'no-store', headers: { 'Content-Type': file.type }, body: file,
    });
  } catch {
    throw new Error('证照上传失败，请重试。');
  }
  if (response.status === 400 || response.status === 415) {
    throw new Error('图片内容无效，请更换 JPG、PNG 或 WebP 图片。');
  }
  if (!response.ok) throw new Error('证照上传失败，请重试。');
  const result = await response.json() as { attachmentId?: string };
  if (!result.attachmentId) throw new Error('证照上传失败，请重试。');
  return result.attachmentId;
}

export async function loadStoreDocumentImage(path: string): Promise<Blob> {
  if (!/^\/uploads\/stores\/[^/.][^/]*\/[^/]+\.(jpg|png|webp)$/.test(path)) {
    throw new Error('证照路径无效。');
  }
  const response = await fetch(`${appConfig.engine.httpUrl}${path}`, { credentials: 'include', cache: 'no-store' });
  if (!response.ok) throw new Error('证照预览加载失败，请重试。');
  return response.blob();
}
