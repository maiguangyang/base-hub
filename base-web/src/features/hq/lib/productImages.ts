import { appConfig } from '@/config/app';

export function productImageSource(url?: string | null): string | undefined {
  if (!url) return undefined;
  return url.startsWith('/uploads/') ? `${appConfig.engine.httpUrl}${url}` : url;
}

export async function uploadProductMainImage(file: File): Promise<string> {
  if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type) || file.size > (5 << 20)) {
    throw new Error('请选择不超过 5 MB 的 JPG、PNG 或 WebP 图片');
  }
  const response = await fetch(`${appConfig.engine.httpUrl}/api/product-main-images`, {
    method: 'POST', credentials: 'include', cache: 'no-store', headers: { 'Content-Type': file.type }, body: file,
  });
  if (!response.ok) throw new Error('图片上传失败');
  const result = await response.json() as { attachmentId?: string };
  if (!result.attachmentId) throw new Error('图片上传失败');
  return result.attachmentId;
}
