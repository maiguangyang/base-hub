// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { toBlob } from 'html-to-image';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';
import { AiMessageList } from './AiMessageList';
import type { AiSessionState } from './useAiSession';

vi.mock('html-to-image', () => ({ toBlob: vi.fn(), toSvg: vi.fn() }));

beforeEach(() => { vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }); });
afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.clearAllMocks(); vi.unstubAllGlobals(); Reflect.deleteProperty(navigator, 'clipboard'); });

function renderMessages(state: AiSessionState) {
  return render(<AdminToastProvider><AiMessageList state={state} onRetract={vi.fn()} canRetract /></AdminToastProvider>);
}

function conversation(messages: AiSessionState['messages']): AiSessionState {
  return { phase: 'completed', messages, transcript: messages.map((message) => ({ kind: 'message', ...message })), secrets: [], interrupted: false };
}

it('复制文字和回复截图按钮在聚焦时显示操作提示', async () => {
  renderMessages(conversation([{ role: 'assistant', text: '结果' }]));
  fireEvent.focus(screen.getByRole('button', { name: '复制消息' }));
  expect((await screen.findByRole('tooltip')).textContent).toContain('复制消息文字');
  fireEvent.blur(screen.getByRole('button', { name: '复制消息' }));
  fireEvent.focus(screen.getByRole('button', { name: '复制回复图片' }));
  expect((await screen.findByRole('tooltip')).textContent).toContain('截图并复制本轮回复');
});

it('复制图片包含本轮完整助手正文，但不包含用户消息或工具状态', async () => {
  const png = new Blob(['png'], { type: 'image/png' });
  const renderedContent: string[] = [];
  vi.mocked(toBlob).mockImplementation(async (node) => {
    renderedContent.push(node.textContent ?? '');
    return png;
  });
  class TestClipboardItem {
    constructor(private readonly values: Record<string, Blob | Promise<Blob>>) {}
    async getType(type: string) { return this.values[type]; }
  }
  vi.stubGlobal('ClipboardItem', TestClipboardItem);
  const write = vi.fn(async (items: TestClipboardItem[]) => {
    expect(await items[0].getType('image/png')).toBe(png);
  });
  Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { write } });
  const state: AiSessionState = { ...conversation([]), transcript: [
    { kind: 'message', role: 'user', text: '查询数据' },
    { kind: 'message', role: 'assistant', text: '先查加盟商。' },
    { kind: 'tool', toolId: 'a', title: '查询加盟商', sequence: 1, status: 'SUCCESS' },
    { kind: 'message', role: 'assistant', text: '**结果：** 1 家' },
  ] };
  renderMessages(state);
  const button = screen.getByRole('button', { name: '复制回复图片' });
  fireEvent.click(button);
  await waitFor(() => expect(write).toHaveBeenCalledTimes(1));
  expect(renderedContent[0]).toContain('先查加盟商。');
  expect(renderedContent[0]).toContain('结果： 1 家');
  expect(renderedContent[0]).not.toContain('查询数据');
  expect(renderedContent[0]).not.toContain('查询加盟商');
  expect(screen.getByRole('status').textContent).toContain('回复图片已复制');
});

it('助手回复未结束时不显示复制图片按钮', () => {
  const state = { ...conversation([{ role: 'user' as const, text: '你好' }, { role: 'assistant' as const, text: '正在回复' }]), phase: 'previewing' as const };
  renderMessages(state);
  expect(screen.queryByRole('button', { name: '复制回复图片' })).toBeNull();
});

it('浏览器不支持图片剪贴板时显示失败提示', async () => {
  const state = conversation([{ role: 'assistant', text: '结果' }]);
  renderMessages(state);
  fireEvent.click(screen.getByRole('button', { name: '复制回复图片' }));
  await waitFor(() => expect(screen.getByRole('alert').textContent).toContain('无法复制图片'));
  expect(toBlob).not.toHaveBeenCalled();
});

it('跨域图片无法读取时仍复制回复并标出图片', async () => {
  const png = new Blob(['png'], { type: 'image/png' });
  const fetchImage = vi.fn().mockRejectedValue(new TypeError('CORS blocked'));
  vi.stubGlobal('fetch', fetchImage);
  vi.mocked(toBlob).mockImplementation(async (node) => {
    expect(node.querySelector('img')).toBeNull();
    expect(node.textContent).toContain('远程图片');
    return png;
  });
  class TestClipboardItem {
    constructor(private readonly values: Record<string, Promise<Blob>>) {}
    async getType(type: string) { return this.values[type]; }
  }
  vi.stubGlobal('ClipboardItem', TestClipboardItem);
  const write = vi.fn(async (items: TestClipboardItem[]) => { expect(await items[0].getType('image/png')).toBe(png); });
  Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { write } });
  renderMessages(conversation([{ role: 'assistant', text: '结果：\n\n![远程图片](https://images.example.test/picture.png)' }]));
  fireEvent.click(screen.getByRole('button', { name: '复制回复图片' }));
  await waitFor(() => expect(screen.getByRole('status').textContent).toContain('部分图片已用文字代替'));
  expect(write).toHaveBeenCalledTimes(1);
  expect(fetchImage).toHaveBeenCalled();
});

it('允许跨域读取的图片会保留在复制图片中', async () => {
  const png = new Blob(['png'], { type: 'image/png' });
  const fetchImage = vi.fn().mockResolvedValue({ ok: true, blob: async () => png });
  vi.stubGlobal('fetch', fetchImage);
  vi.mocked(toBlob).mockImplementation(async (node) => {
    expect(node.querySelector('img')?.getAttribute('src')).toMatch(/^data:image\/png;base64,/);
    return png;
  });
  class TestClipboardItem {
    constructor(readonly values: Record<string, Promise<Blob>>) {}
  }
  vi.stubGlobal('ClipboardItem', TestClipboardItem);
  const write = vi.fn(async (items: TestClipboardItem[]) => { expect(await items[0].values['image/png']).toBe(png); });
  Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { write } });
  renderMessages(conversation([{ role: 'assistant', text: '![可用图片](https://images.example.test/allowed.png)' }]));
  fireEvent.click(screen.getByRole('button', { name: '复制回复图片' }));
  await waitFor(() => expect(screen.getByRole('status').textContent).toContain('回复图片已复制'));
  expect(screen.getByRole('status').textContent).not.toContain('文字代替');
  expect(fetchImage).toHaveBeenCalled();
});

it('外部图片响应头返回后正文卡住仍在五秒内中止并用文字替代', async () => {
  vi.useFakeTimers();
  try {
    const png = new Blob(['png'], { type: 'image/png' });
    let signal: AbortSignal | undefined;
    const fetchImage = vi.fn(async (_url: string, options: { signal: AbortSignal }) => {
      signal = options.signal;
      return {
        ok: true,
        blob: () => new Promise<Blob>((_resolve, reject) => {
          signal?.addEventListener('abort', () => reject(new DOMException('aborted', 'AbortError')));
        }),
      };
    });
    vi.stubGlobal('fetch', fetchImage);
    vi.mocked(toBlob).mockResolvedValue(png);
    class TestClipboardItem {
      constructor(readonly values: Record<string, Promise<Blob>>) {}
    }
    vi.stubGlobal('ClipboardItem', TestClipboardItem);
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { write: vi.fn(async (items: TestClipboardItem[]) => items[0].values['image/png']) } });
    renderMessages(conversation([{ role: 'assistant', text: '![远程图片](https://images.example.test/stalled.png)' }]));
    fireEvent.click(screen.getByRole('button', { name: '复制回复图片' }));
    await vi.advanceTimersByTimeAsync(5001);
    expect(fetchImage).toHaveBeenCalled();
    expect(signal?.aborted).toBe(true);
  } finally {
    vi.useRealTimers();
  }
});

it('图片渲染失败时给出生成阶段提示并保留错误诊断', async () => {
  const failure = new Error('SVG decode failed');
  vi.mocked(toBlob).mockRejectedValue(failure);
  const log = vi.spyOn(console, 'error').mockImplementation(() => {});
  class TestClipboardItem {
    constructor(readonly values: Record<string, Promise<Blob>>) {}
  }
  vi.stubGlobal('ClipboardItem', TestClipboardItem);
  Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { write: vi.fn(async (items: TestClipboardItem[]) => items[0].values['image/png']) } });
  renderMessages(conversation([{ role: 'assistant', text: '统计结果' }]));
  fireEvent.click(screen.getByRole('button', { name: '复制回复图片' }));
  await waitFor(() => expect(screen.getByRole('alert').textContent).toContain('图片生成失败'));
  expect(log).toHaveBeenCalledWith('复制回复图片失败', failure);
});

it('浏览器拒绝剪贴板写入时提示检查权限', async () => {
  vi.mocked(toBlob).mockResolvedValue(new Blob(['png'], { type: 'image/png' }));
  vi.spyOn(console, 'error').mockImplementation(() => {});
  class TestClipboardItem {
    constructor(readonly values: Record<string, Promise<Blob>>) {}
  }
  vi.stubGlobal('ClipboardItem', TestClipboardItem);
  Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { write: vi.fn().mockRejectedValue(new DOMException('denied', 'NotAllowedError')) } });
  renderMessages(conversation([{ role: 'assistant', text: '统计结果' }]));
  fireEvent.click(screen.getByRole('button', { name: '复制回复图片' }));
  await waitFor(() => expect(screen.getByRole('alert').textContent).toContain('剪贴板权限'));
});

it('超过单张图片容量时提示使用文字复制', async () => {
  vi.spyOn(HTMLElement.prototype, 'scrollHeight', 'get').mockImplementation(function (this: HTMLElement) {
    return this.style.left === '-10000px' && this.style.display === 'flex' ? 30000 : 0;
  });
  vi.spyOn(HTMLElement.prototype, 'offsetWidth', 'get').mockImplementation(function (this: HTMLElement) {
    return this.style.left === '-10000px' && this.style.display === 'flex' ? 631 : 0;
  });
  class TestClipboardItem {
    constructor(readonly values: Record<string, Promise<Blob>>) {}
  }
  vi.stubGlobal('ClipboardItem', TestClipboardItem);
  Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { write: vi.fn(async (items: TestClipboardItem[]) => items[0].values['image/png']) } });
  renderMessages(conversation([{ role: 'assistant', text: '很长的回复' }]));
  fireEvent.click(screen.getByRole('button', { name: '复制回复图片' }));
  await waitFor(() => expect(screen.getByRole('alert').textContent).toContain('回复太长，请使用文字复制'));
});
