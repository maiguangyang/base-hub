// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { toBlob } from 'html-to-image';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';
import { AiDrawer } from './AiDrawer';
import type { AiSessionState } from './useAiSession';

vi.mock('html-to-image', () => ({ toBlob: vi.fn() }));
beforeEach(() => { vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }); });
afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.clearAllMocks(); vi.unstubAllGlobals(); Reflect.deleteProperty(navigator, 'clipboard'); });

it('复制整段会话时包含所有轮次并保留工具展开状态', async () => {
  const captured: HTMLElement[] = [];
  const png = new Blob(['png'], { type: 'image/png' });
  vi.mocked(toBlob).mockImplementation(async (node) => { captured.push(node.cloneNode(true) as HTMLElement); return png; });
  class TestClipboardItem {
    constructor(readonly values: Record<string, Promise<Blob>>) {}
  }
  vi.stubGlobal('ClipboardItem', TestClipboardItem);
  const write = vi.fn(async (items: TestClipboardItem[]) => { expect(await items[0].values['image/png']).toBe(png); });
  Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { write } });
  const transcript: AiSessionState['transcript'] = [
    { kind: 'message', role: 'user', text: '第一问' },
    { kind: 'tool', toolId: 'one', title: '读取门店', sequence: 1, status: 'SUCCESS' },
    { kind: 'tool', toolId: 'two', title: '核对门店', sequence: 2, status: 'SUCCESS' },
    { kind: 'message', role: 'assistant', text: '第一答' },
    { kind: 'message', role: 'user', text: '第二问' },
    { kind: 'message', role: 'assistant', text: '第二答' },
  ];
  const state: AiSessionState = { phase: 'completed', messages: [], transcript, secrets: [], interrupted: false };
  const session = { state, startPreview: vi.fn(), confirmRun: vi.fn(), stop: vi.fn(), clear: vi.fn(), retractLastExchange: vi.fn(), canRetractLastExchange: false };
  render(<AdminToastProvider><AiDrawer open onOpenChange={vi.fn()} session={session} /></AdminToastProvider>);
  fireEvent.click(screen.getByRole('button', { name: '复制整个会话图片' }));
  await waitFor(() => expect(write).toHaveBeenCalledTimes(1));
  expect(captured[0].querySelector('[data-ai-tool-group]')?.getAttribute('data-state')).toBe('closed');
  fireEvent.click(screen.getByRole('button', { name: '展开全部工具调用，共 2 次' }));
  fireEvent.click(screen.getByRole('button', { name: '复制整个会话图片' }));
  await waitFor(() => expect(write).toHaveBeenCalledTimes(2));
  const copy = captured[1];
  expect(copy.textContent).toContain('第一问');
  expect(copy.textContent).toContain('第一答');
  expect(copy.textContent).toContain('第二问');
  expect(copy.textContent).toContain('第二答');
  expect(copy.textContent).toContain('核对门店');
  expect(copy.querySelector('[data-ai-tool-group]')?.getAttribute('data-state')).toBe('open');
  expect(screen.getByRole('status').textContent).toContain('会话图片已复制');
});

it('长图生成期间禁用截图按钮以防重复复制', async () => {
  let finish!: (blob: Blob) => void;
  vi.mocked(toBlob).mockImplementation(() => new Promise<Blob>((resolve) => { finish = resolve; }));
  class TestClipboardItem {
    constructor(readonly values: Record<string, Promise<Blob>>) {}
  }
  vi.stubGlobal('ClipboardItem', TestClipboardItem);
  const write = vi.fn(async (items: TestClipboardItem[]) => items[0].values['image/png']);
  Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { write } });
  const transcript: AiSessionState['transcript'] = [{ kind: 'message', role: 'user', text: '需要截图' }];
  const state: AiSessionState = { phase: 'completed', messages: [], transcript, secrets: [], interrupted: false };
  const session = { state, startPreview: vi.fn(), confirmRun: vi.fn(), stop: vi.fn(), clear: vi.fn(), retractLastExchange: vi.fn(), canRetractLastExchange: false };
  render(<AdminToastProvider><AiDrawer open onOpenChange={vi.fn()} session={session} /></AdminToastProvider>);
  const button = screen.getByRole('button', { name: '复制整个会话图片' }) as HTMLButtonElement;
  fireEvent.click(button);
  expect(button.disabled).toBe(true);
  fireEvent.click(button);
  expect(write).toHaveBeenCalledTimes(1);
  await waitFor(() => expect(toBlob).toHaveBeenCalledTimes(1));
  finish(new Blob(['png'], { type: 'image/png' }));
  await waitFor(() => expect(button.disabled).toBe(false));
});

it('会话截图按钮在聚焦时显示操作提示', async () => {
  const state: AiSessionState = { phase: 'completed', messages: [], transcript: [{ kind: 'message', role: 'user', text: '你好' }], secrets: [], interrupted: false };
  const session = { state, startPreview: vi.fn(), confirmRun: vi.fn(), stop: vi.fn(), clear: vi.fn(), retractLastExchange: vi.fn(), canRetractLastExchange: false };
  render(<AdminToastProvider><AiDrawer open onOpenChange={vi.fn()} session={session} /></AdminToastProvider>);
  fireEvent.focus(screen.getByRole('button', { name: '复制整个会话图片' }));
  expect((await screen.findByRole('tooltip')).textContent).toContain('截图并复制整个会话');
});
