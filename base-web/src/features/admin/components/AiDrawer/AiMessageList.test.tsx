// @vitest-environment jsdom

import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';
import { AiMessageList } from './AiMessageList';
import type { AiSessionState } from './useAiSession';

beforeEach(() => { vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }); });
afterEach(() => { cleanup(); vi.unstubAllGlobals(); Reflect.deleteProperty(navigator, 'clipboard'); Reflect.deleteProperty(document, 'execCommand'); });

function renderMessages(state: AiSessionState) {
  return render(<AdminToastProvider><AiMessageList state={state} onRetract={vi.fn()} canRetract /></AdminToastProvider>);
}

function conversation(messages: AiSessionState['messages']): AiSessionState {
  return { phase: 'completed', messages, transcript: messages.map((message) => ({ kind: 'message', ...message })), secrets: [], interrupted: false };
}

it('工具调用在最终回复之前显示，完成后保留原位置', () => {
  const state: AiSessionState = { phase: 'completed', messages: [
    { role: 'user', text: '有多少加盟商？' }, { role: 'assistant', text: '最终结果' },
  ], transcript: [
    { kind: 'message', role: 'user', text: '有多少加盟商？' },
    { kind: 'tool', toolId: 'select_tools', title: '选择后台工具', sequence: 1, status: 'SUCCESS' },
    { kind: 'message', role: 'assistant', text: '最终结果' },
  ], secrets: [], interrupted: false };
  renderMessages(state);
  const tool = screen.getAllByText(/选择后台工具/)[0];
  const answer = screen.getByText('最终结果');
  expect(tool.compareDocumentPosition(answer) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy();
});

it('工具执行时显示带流光的胶囊，完成后停止动画', () => {
  const running: AiSessionState = { ...conversation([]), phase: 'previewing', transcript: [
    { kind: 'tool', toolId: 'select_tools', title: '选择后台工具', sequence: 1, status: 'RUNNING' },
  ] };
  const { rerender } = renderMessages(running);
  const pill = screen.getByText(/选择后台工具/).closest('[data-ai-tool]');
  expect(pill?.className).toContain('rounded-full');
  expect(pill?.className).toContain('border');
  expect(pill?.className).toContain('ai-tool-running');
  rerender(<AdminToastProvider><AiMessageList state={{ ...running, phase: 'completed', transcript: [
    { kind: 'tool', toolId: 'select_tools', title: '选择后台工具', sequence: 1, status: 'SUCCESS' },
  ] }} onRetract={vi.fn()} canRetract /></AdminToastProvider>);
  expect(screen.getByText(/选择后台工具/).closest('[data-ai-tool]')?.className).not.toContain('ai-tool-running');
});

it('连续工具折叠为最新胶囊，展开后按原顺序查看全部', () => {
  vi.useFakeTimers();
  try {
    const state: AiSessionState = { ...conversation([]), transcript: [
      { kind: 'tool', toolId: 'a', title: '选择工具', sequence: 1, status: 'SUCCESS' },
      { kind: 'tool', toolId: 'b', title: '查询加盟商', sequence: 2, status: 'SUCCESS' },
      { kind: 'tool', toolId: 'c', title: '查询门店', sequence: 3, status: 'SUCCESS' },
    ] };
    renderMessages(state);
    expect(screen.getByText('选择工具')).toBeTruthy();
    expect(screen.getByText('查询加盟商')).toBeTruthy();
    act(() => vi.advanceTimersByTime(100));
    act(() => vi.advanceTimersByTime(220));
    act(() => vi.advanceTimersByTime(250));
    act(() => vi.advanceTimersByTime(220));
    expect(screen.getByText('查询门店')).toBeTruthy();
    expect(screen.queryByText('选择工具')).toBeNull();
    fireEvent.click(screen.getByRole('button', { name: '展开全部工具调用，共 3 次' }));
    expect(screen.getByText('选择工具')).toBeTruthy();
    expect(screen.getByText('查询加盟商')).toBeTruthy();
    expect(screen.getByRole('button', { name: '收起工具调用，共 3 次' }).getAttribute('aria-expanded')).toBe('true');
  } finally { vi.useRealTimers(); }
});

it('助手正文会断开连续工具组', () => {
  const state: AiSessionState = { ...conversation([]), transcript: [
    { kind: 'tool', toolId: 'a', title: '选择工具', sequence: 1, status: 'SUCCESS' },
    { kind: 'tool', toolId: 'b', title: '查询加盟商', sequence: 2, status: 'SUCCESS' },
    { kind: 'message', role: 'assistant', text: '中间回复' },
    { kind: 'tool', toolId: 'c', title: '查询门店', sequence: 3, status: 'SUCCESS' },
  ] };
  renderMessages(state);
  expect(screen.getAllByRole('button', { name: '展开全部工具调用，共 2 次' })).toHaveLength(1);
  expect(screen.getByText('中间回复').compareDocumentPosition(screen.getByText('查询门店')) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy();
});

it('新工具先在下方停留 250 毫秒再上移覆盖旧胶囊', () => {
  vi.useFakeTimers();
  try {
    const first: AiSessionState = { ...conversation([]), transcript: [
      { kind: 'tool', toolId: 'a', title: '选择工具', sequence: 1, status: 'SUCCESS' },
    ] };
    const { rerender } = renderMessages(first);
    rerender(<AdminToastProvider><AiMessageList state={{ ...first, transcript: [...first.transcript,
      { kind: 'tool', toolId: 'b', title: '查询加盟商', sequence: 2, status: 'RUNNING' },
    ] }} onRetract={vi.fn()} canRetract /></AdminToastProvider>);
    expect(screen.getByText('查询加盟商').closest('.ai-tool-stack-incoming')?.getAttribute('data-phase')).toBe('waiting');
    act(() => vi.advanceTimersByTime(249));
    expect(screen.getByText('选择工具')).toBeTruthy();
    act(() => vi.advanceTimersByTime(1));
    expect(screen.getByText('查询加盟商').closest('.ai-tool-stack-incoming')?.getAttribute('data-phase')).toBe('moving');
    act(() => vi.advanceTimersByTime(220));
    expect(screen.queryByText('选择工具')).toBeNull();
    expect(screen.getByText('查询加盟商')).toBeTruthy();
  } finally { vi.useRealTimers(); }
});

it('多个工具事件同批到达时仍逐个完成覆盖动画', () => {
  vi.useFakeTimers();
  try {
    const state: AiSessionState = { ...conversation([]), phase: 'previewing', transcript: [
      { kind: 'tool', toolId: 'a', title: '选择工具', sequence: 1, status: 'SUCCESS' },
      { kind: 'tool', toolId: 'b', title: '查询加盟商', sequence: 2, status: 'SUCCESS' },
      { kind: 'tool', toolId: 'c', title: '查询门店', sequence: 3, status: 'RUNNING' },
    ] };
    renderMessages(state);
    expect(screen.getByText('选择工具')).toBeTruthy();
    expect(screen.getByText('查询加盟商')).toBeTruthy();
    expect(screen.queryByText('查询门店')).toBeNull();
    act(() => vi.advanceTimersByTime(100));
    act(() => vi.advanceTimersByTime(220));
    expect(screen.getByText('查询加盟商')).toBeTruthy();
    expect(screen.getByText('查询门店')).toBeTruthy();
    act(() => vi.advanceTimersByTime(250));
    act(() => vi.advanceTimersByTime(220));
    expect(screen.queryByText('查询加盟商')).toBeNull();
    expect(screen.getByText('查询门店')).toBeTruthy();
  } finally { vi.useRealTimers(); }
});

it('等待期间到达第三个工具时改用快速节奏', () => {
  vi.useFakeTimers();
  try {
    const first: AiSessionState = { ...conversation([]), transcript: [
      { kind: 'tool', toolId: 'a', title: '选择工具', sequence: 1, status: 'SUCCESS' },
    ] };
    const second = { kind: 'tool' as const, toolId: 'b', title: '查询加盟商', sequence: 2, status: 'SUCCESS' as const };
    const { rerender } = renderMessages(first);
    rerender(<AdminToastProvider><AiMessageList state={{ ...first, transcript: [...first.transcript, second] }} onRetract={vi.fn()} canRetract /></AdminToastProvider>);
    act(() => vi.advanceTimersByTime(90));
    rerender(<AdminToastProvider><AiMessageList state={{ ...first, transcript: [...first.transcript, second,
      { kind: 'tool', toolId: 'c', title: '查询门店', sequence: 3, status: 'RUNNING' },
    ] }} onRetract={vi.fn()} canRetract /></AdminToastProvider>);
    act(() => vi.advanceTimersByTime(99));
    expect(screen.getByText('查询加盟商').closest('.ai-tool-stack-incoming')?.getAttribute('data-phase')).toBe('waiting');
    act(() => vi.advanceTimersByTime(1));
    expect(screen.getByText('查询加盟商').closest('.ai-tool-stack-incoming')?.getAttribute('data-phase')).toBe('moving');
  } finally { vi.useRealTimers(); }
});

it('每条消息都可复制原文', async () => {
  const writeText = vi.fn().mockResolvedValue(undefined);
  Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText } });
  const state = conversation([
    { role: 'user', text: '你好' }, { role: 'assistant', text: '**你好！**' },
  ]);
  renderMessages(state);
  const buttons = screen.getAllByRole('button', { name: '复制消息' });
  expect(buttons).toHaveLength(2);
  fireEvent.click(buttons[1]);
  await waitFor(() => expect(writeText).toHaveBeenCalledWith('**你好！**'));
});

it('同一轮助手回复在结束前不显示复制，结束后只复制一次完整回复', async () => {
  const writeText = vi.fn().mockResolvedValue(undefined);
  Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText } });
  const transcript: AiSessionState['transcript'] = [
    { kind: 'message', role: 'user', text: '有多少加盟商和门店？' },
    { kind: 'message', role: 'assistant', text: '我先查询加盟商。' },
    { kind: 'tool', toolId: 'a', title: '查询加盟商', sequence: 1, status: 'SUCCESS' },
    { kind: 'message', role: 'assistant', text: '当前有 1 家加盟商。' },
  ];
  const state: AiSessionState = { ...conversation([]), phase: 'previewing', transcript };
  const { rerender } = renderMessages(state);
  expect(screen.getAllByRole('button', { name: '复制消息' })).toHaveLength(1);
  rerender(<AdminToastProvider><AiMessageList state={{ ...state, phase: 'completed' }} onRetract={vi.fn()} canRetract /></AdminToastProvider>);
  const buttons = screen.getAllByRole('button', { name: '复制消息' });
  expect(buttons).toHaveLength(2);
  const firstAssistant = screen.getByText('我先查询加盟商。').closest('article');
  expect(within(firstAssistant!).queryByRole('button', { name: '复制消息' })).toBeNull();
  fireEvent.click(buttons[1]);
  await waitFor(() => expect(writeText).toHaveBeenCalledWith('我先查询加盟商。\n\n当前有 1 家加盟商。'));
});

it('工具组的展开按钮在胶囊外侧，收拢时显示多层轮廓', () => {
  const state: AiSessionState = { ...conversation([]), transcript: [
    { kind: 'tool', toolId: 'a', title: '选择工具', sequence: 1, status: 'SUCCESS' },
    { kind: 'tool', toolId: 'b', title: '查询加盟商', sequence: 2, status: 'SUCCESS' },
  ] };
  renderMessages(state);
  const group = document.querySelector('[data-ai-tool-group]');
  const stack = group?.querySelector('[data-ai-tool-stack]');
  expect(stack?.querySelector('[data-ai-tool-stack-layers]')).not.toBeNull();
  const trigger = screen.getByRole('button', { name: '展开全部工具调用，共 2 次' });
  expect(trigger.parentElement).toBe(stack);
  expect(trigger.closest('[data-ai-tool]')).toBeNull();
  expect(stack?.className).toContain('items-start');
});

it('展开入口显示连续调用次数并向读屏器说明数量', () => {
  const state: AiSessionState = { ...conversation([]), transcript: ['选择工具', '查询加盟商', '查询门店'].map((title, index) =>
    ({ kind: 'tool', toolId: String(index), title, sequence: index + 1, status: 'SUCCESS' })) };
  renderMessages(state);
  const trigger = screen.getByRole('button', { name: '展开全部工具调用，共 3 次' });
  expect(within(trigger).getByText('3 次')).toBeTruthy();
  fireEvent.click(trigger);
  expect(screen.getByRole('button', { name: '收起工具调用，共 3 次' })).toBeTruthy();
});

it('工具状态文字按状态着色且保留可读标签', () => {
  const statuses = [
    ['RUNNING', '执行中', 'text-primary', 'bg-primary'],
    ['SUCCESS', '已完成', 'text-success-fg', 'bg-success-fg'],
    ['FAILED', '执行失败', 'text-destructive', 'bg-destructive'],
    ['INTERRUPTED', '已中断', 'text-warning-fg', 'bg-warning-fg'],
    ['UNKNOWN', '状态未确认', 'text-muted-foreground', 'bg-muted-foreground'],
  ] as const;
  const state: AiSessionState = { ...conversation([]), transcript: statuses.map(([status], index) => ({
    kind: 'tool' as const, toolId: String(index), title: `工具 ${index}`, sequence: index + 1, status,
  })).flatMap((entry) => [entry, { kind: 'message' as const, role: 'assistant' as const, text: `分隔 ${entry.sequence}` }]) };
  renderMessages(state);
  statuses.forEach(([status, label, color, dotColor], index) => {
    const pill = document.querySelector(`[data-ai-tool="${index}"]`);
    expect(pill?.getAttribute('data-state')).toBe(status.toLowerCase());
    const text = within(pill as HTMLElement).getByText(label);
    expect(text.className).toContain(color);
    expect(pill?.querySelector('[aria-hidden="true"]')?.className).toContain(dotColor);
  });
});

it('Clipboard API 不可用时仍可复制消息', async () => {
  const execCommand = vi.fn().mockReturnValue(true);
  Object.defineProperty(document, 'execCommand', { configurable: true, value: execCommand });
  renderMessages(conversation([{ role: 'assistant', text: '原文' }]));
  fireEvent.click(screen.getByRole('button', { name: '复制消息' }));
  await waitFor(() => expect(execCommand).toHaveBeenCalledWith('copy'));
  expect(screen.getByRole('status').textContent).toContain('消息已复制');
});

it('仅最后一条用户消息显示撤回按钮', () => {
  const onRetract = vi.fn();
  const state = conversation([
    { role: 'user', text: '第一问' }, { role: 'assistant', text: '第一答' },
    { role: 'user', text: '第二问' }, { role: 'assistant', text: '第二答' },
  ]);
  render(<AdminToastProvider><AiMessageList state={state} onRetract={onRetract} canRetract /></AdminToastProvider>);
  const first = screen.getByText('第一问').closest('article');
  const last = screen.getByText('第二问').closest('article');
  expect(first && within(first).queryByRole('button', { name: '撤回最后一条消息' })).toBeNull();
  expect(last).not.toBeNull();
  fireEvent.click(within(last!).getByRole('button', { name: '撤回最后一条消息' }));
  expect(onRetract).toHaveBeenCalledTimes(1);
});

it('助手消息按 GFM 渲染并保留安全链接', () => {
  const state = conversation([{ role: 'assistant', text: '## 可用功能\n- 查询门店\n- 查看报表\n\n[帮助](https://example.com/help)' }]);
  renderMessages(state);
  expect(screen.getByRole('heading', { name: '可用功能' })).toBeTruthy();
  expect(screen.getAllByRole('listitem')).toHaveLength(2);
  expect(screen.getByRole('link', { name: '帮助' }).getAttribute('href')).toBe('https://example.com/help');
});

it('用户消息使用柔和底色', () => {
  const state = conversation([{ role: 'user', text: '你好' }]);
  renderMessages(state);
  expect(screen.getByText('你好').className).toContain('bg-primary/10');
});

it('GFM 表格与删除线按结构呈现', () => {
  const state = conversation([{ role: 'assistant', text: '| 项目 | 状态 |\n| --- | --- |\n| 门店 | ~~旧~~新 |' }]);
  renderMessages(state);
  expect(screen.getByRole('table')).toBeTruthy();
  expect(screen.getByText('旧').tagName).toBe('DEL');
});
