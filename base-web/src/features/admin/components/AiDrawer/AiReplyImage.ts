import { toBlob } from 'html-to-image';

type ToBlob = typeof import('html-to-image').toBlob;
const MAX_IMAGE_SIZE = 4096;
const PAGE_GAP = 12;

export class ReplyImageTooLargeError extends Error {}

export async function renderReplyImage(list: HTMLDivElement | null, messageIndexes: number[], onImageFallback: () => void): Promise<Blob> {
  const sources = messageIndexes.map((index) => list?.querySelector<HTMLElement>(`[data-ai-message-index="${index}"] > [data-ai-message-content]`));
  if (sources.some((source) => !source)) throw new Error('回复内容已不可用');
  const content = sources as HTMLElement[];
  const surface = content[0].closest<HTMLElement>('.ai-drawer-surface') ?? document.body;
  const backgroundColor = getComputedStyle(surface).backgroundColor;
  const surfaceWidth = surface.getBoundingClientRect().width;
  const availableWidth = surfaceWidth > 48 ? surfaceWidth - 48 : 320;
  const width = Math.min(Math.max(...content.map((source) => source.getBoundingClientRect().width)), availableWidth, MAX_IMAGE_SIZE - 24) || 320;
  const canvas = document.createElement('div');
  Object.assign(canvas.style, {
    position: 'fixed', left: '-10000px', top: '0', boxSizing: 'border-box',
    width: `${Math.ceil(width) + 24}px`, padding: '12px', display: 'flex',
    flexDirection: 'column', gap: '12px', backgroundColor,
  });
  content.forEach((source) => canvas.appendChild(source.cloneNode(true)));
  surface.appendChild(canvas);
  try {
    prepareCaptureLayout(canvas);
    await prepareExternalImages(canvas, onImageFallback);
    if (canvas.scrollHeight > MAX_IMAGE_SIZE) return renderPagedReplyImage(canvas, backgroundColor, toBlob);
    // 导出的克隆节点必须回到原点，否则离屏定位也会被画进图片。
    return await captureBlob(canvas, backgroundColor, toBlob);
  } finally {
    canvas.remove();
  }
}

export function prepareCaptureLayout(root: HTMLElement): void {
  root.style.overflowWrap = 'anywhere';
  root.querySelectorAll<HTMLElement>('pre, pre code').forEach((pre) => {
    pre.style.whiteSpace = 'pre-wrap';
    pre.style.overflowWrap = 'anywhere';
    pre.style.wordBreak = 'break-word';
    pre.style.overflowX = 'visible';
  });
  root.querySelectorAll<HTMLElement>('table').forEach((table) => {
    table.style.width = '100%';
    table.style.tableLayout = 'fixed';
    if (table.parentElement) table.parentElement.style.overflowX = 'visible';
  });
  root.querySelectorAll<HTMLElement>('th, td').forEach((cell) => { cell.style.overflowWrap = 'anywhere'; });
  root.querySelectorAll<HTMLImageElement>('img').forEach((image) => { image.style.maxWidth = '100%'; image.style.height = 'auto'; });
}

async function captureBlob(node: HTMLElement, backgroundColor: string, toBlob: ToBlob): Promise<Blob> {
  const blob = await toBlob(node, { backgroundColor, pixelRatio: 1, style: { position: 'static', left: '0', top: '0' } });
  if (!blob) throw new Error('图片生成失败');
  return blob;
}

class CapturePages {
  readonly pages: HTMLElement[] = [];
  private page: HTMLElement;

  constructor(private readonly source: HTMLElement) {
    this.page = this.openPage();
  }

  private openPage(): HTMLElement {
    const page = this.source.cloneNode(false) as HTMLElement;
    this.source.parentElement?.appendChild(page);
    this.pages.push(page);
    const width = page.offsetWidth;
    if (this.pages.length * width + (this.pages.length - 1) * PAGE_GAP > MAX_IMAGE_SIZE) throw new ReplyImageTooLargeError();
    return page;
  }

  add(message: HTMLElement): void {
    const markdown = message.firstElementChild as HTMLElement | null;
    if (!markdown) {
      this.addWholeMessage(message);
      return;
    }
    let target: HTMLElement | null = null;
    const pending: Node[] = Array.from(markdown.childNodes);
    while (pending.length) {
      const block = pending.shift()!;
      if (!target) target = this.addMessageShell(message, markdown);
      const clone = block.cloneNode(true);
      target.appendChild(clone);
      if (this.page.scrollHeight <= MAX_IMAGE_SIZE) continue;
      target.removeChild(clone);
      if (!target.hasChildNodes()) this.page.lastElementChild?.remove();
      target = null;
      if (this.page.hasChildNodes()) {
        this.page = this.openPage();
        pending.unshift(block);
        continue;
      }
      const parts = splitBlockAtMiddle(block);
      if (!parts) throw new ReplyImageTooLargeError();
      pending.unshift(...parts);
    }
  }

  private addWholeMessage(message: HTMLElement): void {
    this.page.appendChild(message.cloneNode(true));
    if (this.page.scrollHeight <= MAX_IMAGE_SIZE) return;
    this.page.lastElementChild?.remove();
    if (!this.page.hasChildNodes()) throw new ReplyImageTooLargeError();
    this.page = this.openPage();
    this.page.appendChild(message.cloneNode(true));
    if (this.page.scrollHeight > MAX_IMAGE_SIZE) throw new ReplyImageTooLargeError();
  }

  private addMessageShell(message: HTMLElement, markdown: HTMLElement): HTMLElement {
    const bubble = message.cloneNode(false) as HTMLElement;
    const body = markdown.cloneNode(false) as HTMLElement;
    bubble.appendChild(body);
    this.page.appendChild(bubble);
    return body;
  }

  remove(): void { this.pages.forEach((page) => page.remove()); }
}

function splitBlockAtMiddle(block: Node): [Node, Node] | null {
  if (!(block instanceof Element)) return null;
  const walker = document.createTreeWalker(block, NodeFilter.SHOW_TEXT);
  const textNodes: Text[] = [];
  let text = walker.nextNode();
  while (text) {
    textNodes.push(text as Text);
    text = walker.nextNode();
  }
  const total = textNodes.reduce((count, node) => count + node.length, 0);
  if (total < 2) return null;
  let middle = Math.floor(total / 2);
  const splitNode = textNodes.find((node) => {
    if (middle <= node.length) return true;
    middle -= node.length;
    return false;
  });
  if (!splitNode) return null;
  const first = document.createRange();
  first.selectNodeContents(block);
  first.setEnd(splitNode, middle);
  const second = document.createRange();
  second.selectNodeContents(block);
  second.setStart(splitNode, middle);
  const left = block.cloneNode(false);
  const right = block.cloneNode(false);
  left.appendChild(first.cloneContents());
  right.appendChild(second.cloneContents());
  return [left, right];
}

async function renderPagedReplyImage(content: HTMLElement, backgroundColor: string, toBlob: ToBlob): Promise<Blob> {
  const pages = new CapturePages(content);
  try {
    Array.from(content.children).forEach((message) => pages.add(message as HTMLElement));
    return await composePageImages(pages.pages, backgroundColor, toBlob);
  } finally {
    pages.remove();
  }
}

async function composePageImages(pages: HTMLElement[], backgroundColor: string, toBlob: ToBlob): Promise<Blob> {
  const width = pages[0].offsetWidth;
  const output = document.createElement('canvas');
  output.width = pages.length * width + (pages.length - 1) * PAGE_GAP;
  output.height = Math.max(...pages.map((page) => page.scrollHeight));
  const context = output.getContext('2d');
  if (!context) throw new Error('浏览器无法绘制图片');
  context.fillStyle = backgroundColor;
  context.fillRect(0, 0, output.width, output.height);
  for (let index = 0; index < pages.length; index += 1) {
    const blob = await captureBlob(pages[index], backgroundColor, toBlob);
    const url = URL.createObjectURL(blob);
    try {
      const image = new Image();
      image.src = url;
      await image.decode();
      context.drawImage(image, index * (width + PAGE_GAP), 0);
    } finally {
      URL.revokeObjectURL(url);
    }
  }
  return new Promise<Blob>((resolve, reject) => output.toBlob((blob) => blob ? resolve(blob) : reject(new Error('图片合成失败')), 'image/png'));
}

export async function prepareExternalImages(root: HTMLElement, onFallback: () => void): Promise<void> {
  await Promise.all(Array.from(root.querySelectorAll('img')).map(async (image) => {
    const url = image.currentSrc || image.src;
    if (!url.startsWith('http') || new URL(url).origin === window.location.origin) return;
    try {
      await inlineExternalImage(image, url);
    } catch {
      const placeholder = document.createElement('span');
      placeholder.className = 'inline-block rounded-md border border-border bg-background px-2 py-1 text-muted-foreground';
      placeholder.textContent = `图片无法嵌入：${image.alt || '无描述图片'}`;
      image.replaceWith(placeholder);
      onFallback();
    }
  }));
}

async function inlineExternalImage(image: HTMLImageElement, url: string): Promise<void> {
  const controller = new AbortController();
  const timeout = window.setTimeout(() => controller.abort(), 5000);
  try {
    const response = await fetch(url, { signal: controller.signal });
    if (!response.ok) throw new Error('图片不可用');
    const blob = await response.blob();
    if (!blob.type.startsWith('image/')) throw new Error('不是图片');
    const dataUrl = await readImageDataUrl(blob, controller.signal);
    image.srcset = '';
    image.src = dataUrl;
  } finally {
    window.clearTimeout(timeout);
  }
}

function readImageDataUrl(blob: Blob, signal: AbortSignal): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    const abort = () => reader.abort();
    signal.addEventListener('abort', abort, { once: true });
    reader.onload = () => typeof reader.result === 'string' ? resolve(reader.result) : reject(new Error('图片读取失败'));
    reader.onerror = () => reject(new Error('图片读取失败'));
    reader.onabort = () => reject(new Error('图片读取超时'));
    reader.onloadend = () => signal.removeEventListener('abort', abort);
    reader.readAsDataURL(blob);
  });
}
