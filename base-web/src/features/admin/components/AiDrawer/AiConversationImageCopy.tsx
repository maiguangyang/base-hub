import { useState } from 'react';
import { Camera, LoaderCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { useAdminToast } from '@/features/admin/components/AdminToast';
import { renderConversationImage } from './AiConversationImage';
import { ReplyImageTooLargeError } from './AiReplyImage';

export function AiConversationImageCopy({ source, disabled }: { source(): HTMLDivElement | null; disabled: boolean }) {
  const showToast = useAdminToast();
  const [copying, setCopying] = useState(false);
  async function copyImage() {
    if (!navigator.clipboard?.write || typeof ClipboardItem === 'undefined') {
      showToast('error', '无法复制图片，请使用支持图片剪贴板的浏览器');
      return;
    }
    setCopying(true);
    let renderError: unknown;
    try {
      let replacedImages = 0;
      const image = renderConversationImage(source(), () => { replacedImages += 1; }).catch((error: unknown) => {
        renderError = error;
        throw error;
      });
      void image.catch(() => {});
      const written = navigator.clipboard.write([new ClipboardItem({ 'image/png': image })]);
      void written.catch(() => {});
      await image;
      await written;
      showToast('success', replacedImages ? '会话图片已复制，部分图片已用文字代替' : '会话图片已复制');
    } catch (error) {
      const cause = renderError ?? error;
      console.error('复制会话图片失败', cause);
      showToast('error', imageCopyFailureMessage(cause, renderError !== undefined));
    } finally {
      setCopying(false);
    }
  }
  return <TooltipProvider><Tooltip><TooltipTrigger asChild><Button type="button" size="icon-sm" variant="ghost" aria-label="复制整个会话图片" disabled={disabled || copying} className="absolute right-3 bottom-2 size-9 text-muted-foreground" onClick={() => void copyImage()}>{copying ? <LoaderCircle aria-hidden="true" className="size-5 animate-spin" /> : <Camera aria-hidden="true" className="size-5" />}</Button></TooltipTrigger><TooltipContent side="bottom">截图并复制整个会话</TooltipContent></Tooltip></TooltipProvider>;
}

function imageCopyFailureMessage(error: unknown, renderFailed: boolean): string {
  if (error instanceof ReplyImageTooLargeError) return '会话太长，超出单张图片容量';
  if (renderFailed) return '会话图片生成失败，请重试';
  if (error instanceof DOMException && error.name === 'NotAllowedError') return '浏览器未允许复制图片，请检查剪贴板权限';
  return '图片写入剪贴板失败，请重试';
}
