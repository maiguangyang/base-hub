import { toBlob } from 'html-to-image';
import { prepareCaptureLayout, prepareExternalImages, ReplyImageTooLargeError } from './AiReplyImage';

const STRIP_HEIGHT = 4096;
const COLUMN_HEIGHT = 16384;
const MAX_IMAGE_WIDTH = 4096;
const PADDING = 12;
const COLUMN_GAP = 12;

export async function renderConversationImage(list: HTMLDivElement | null, onImageFallback: () => void): Promise<Blob> {
  if (!list || !list.children.length) throw new Error('会话内容已不可用');
  const surface = list.closest<HTMLElement>('.ai-drawer-surface') ?? document.body;
  const backgroundColor = getComputedStyle(surface).backgroundColor;
  const contentWidth = Math.ceil(list.getBoundingClientRect().width) || 320;
  const width = contentWidth + PADDING * 2;
  if (width > MAX_IMAGE_WIDTH) throw new ReplyImageTooLargeError();
  const host = document.createElement('div');
  host.className = 'ai-drawer-surface';
  Object.assign(host.style, {
    position: 'fixed', left: '-10000px', top: '0', boxSizing: 'border-box',
    width: `${width}px`, padding: `${PADDING}px`, backgroundColor,
  });
  const copy = list.cloneNode(true) as HTMLDivElement;
  copy.style.width = `${contentWidth}px`;
  copy.querySelectorAll<HTMLElement>('*').forEach((element) => {
    element.style.animation = 'none';
    element.style.transition = 'none';
  });
  host.appendChild(copy);
  document.body.appendChild(host);
  try {
    prepareCaptureLayout(copy);
    await prepareExternalImages(copy, onImageFallback);
    const totalHeight = Math.ceil(copy.scrollHeight) + PADDING * 2;
    if (totalHeight <= STRIP_HEIGHT) return await captureStrip(host, backgroundColor, totalHeight);
    return await composeConversation(host, copy, backgroundColor, width, totalHeight);
  } finally {
    host.remove();
  }
}

async function captureStrip(host: HTMLElement, backgroundColor: string, height: number): Promise<Blob> {
  host.style.height = `${height}px`;
  const blob = await toBlob(host, {
    width: host.offsetWidth, height, backgroundColor, pixelRatio: 1,
    style: { position: 'static', left: '0', top: '0' },
  });
  if (!blob) throw new Error('图片生成失败');
  return blob;
}

async function composeConversation(host: HTMLElement, copy: HTMLElement, backgroundColor: string, width: number, totalHeight: number): Promise<Blob> {
  const columns = Math.ceil(totalHeight / COLUMN_HEIGHT);
  const outputWidth = columns * width + (columns - 1) * COLUMN_GAP;
  if (outputWidth > MAX_IMAGE_WIDTH) throw new ReplyImageTooLargeError();
  const output = document.createElement('canvas');
  output.width = outputWidth;
  output.height = Math.min(totalHeight, COLUMN_HEIGHT);
  const context = output.getContext('2d');
  if (!context) throw new Error('浏览器无法绘制图片');
  context.fillStyle = backgroundColor;
  context.fillRect(0, 0, output.width, output.height);
  host.style.overflow = 'hidden';
  copy.style.position = 'relative';
  for (let offset = 0; offset < totalHeight;) {
    const column = Math.floor(offset / COLUMN_HEIGHT);
    const top = offset % COLUMN_HEIGHT;
    const height = Math.min(STRIP_HEIGHT, COLUMN_HEIGHT - top, totalHeight - offset);
    copy.style.top = `${-offset}px`;
    const blob = await captureStrip(host, backgroundColor, height);
    const url = URL.createObjectURL(blob);
    try {
      const image = new Image();
      image.src = url;
      await image.decode();
      context.drawImage(image, column * (width + COLUMN_GAP), top);
    } finally {
      URL.revokeObjectURL(url);
    }
    offset += height;
  }
  return new Promise<Blob>((resolve, reject) => output.toBlob((blob) => blob ? resolve(blob) : reject(new Error('图片合成失败')), 'image/png'));
}
