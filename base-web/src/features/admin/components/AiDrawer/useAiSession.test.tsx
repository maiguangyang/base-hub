// @vitest-environment jsdom

import { useLayoutEffect } from 'react';
import { act, cleanup, renderHook, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { AIStreamError } from '@/features/admin/lib/aiStream';
import { useAiSession } from './useAiSession';

const mocks = vi.hoisted(() => ({ readAiStream: vi.fn() }));
vi.mock('@/features/admin/lib/aiStream', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@/features/admin/lib/aiStream')>()), readAiStream: mocks.readAiStream,
}));

beforeEach(() => { vi.clearAllMocks(); sessionStorage.clear(); });
afterEach(() => cleanup());

it('按工具开始、完成、最终回复的到达顺序保留显示记录', async () => {
  mocks.readAiStream.mockImplementation(async (_path: string, _body: object, _signal: AbortSignal, onEvent: (event: unknown) => void) => {
    onEvent({ type: 'tool_started', toolId: 'select_tools', title: '选择后台工具', sequence: 1 });
    onEvent({ type: 'tool_finished', toolId: 'select_tools', title: '选择后台工具', sequence: 1, status: 'SUCCESS' });
    onEvent({ type: 'text_delta', text: '当前有 1 个加盟商。' });
    onEvent({ type: 'run_finished', status: 'SUCCESS' });
  });
  const { result } = renderHook(() => useAiSession('account:hq'));
  await act(async () => result.current.startPreview('有多少加盟商？'));
  expect(result.current.state.transcript).toEqual([
    { kind: 'message', role: 'user', text: '有多少加盟商？' },
    { kind: 'tool', toolId: 'select_tools', title: '选择后台工具', sequence: 1, status: 'SUCCESS' },
    { kind: 'message', role: 'assistant', text: '当前有 1 个加盟商。' },
  ]);
});

it('工具未返回完成事件时随失败结束为失败状态', async () => {
  mocks.readAiStream.mockImplementation(async (_path: string, _body: object, _signal: AbortSignal, onEvent: (event: unknown) => void) => {
    onEvent({ type: 'tool_started', toolId: 'select_tools', title: '选择后台工具', sequence: 1 });
    onEvent({ type: 'run_finished', status: 'FAILED' });
  });
  const { result } = renderHook(() => useAiSession('account:hq'));
  await act(async () => result.current.startPreview('查询'));
  expect(result.current.state.transcript[1]).toMatchObject({ kind: 'tool', status: 'FAILED' });
});

describe('普通对话', () => {
  it('普通对话保留上下文并在下一轮发送历史正文', async () => {
    mocks.readAiStream.mockImplementation(async (_path: string, body: { prompt: string }, _signal: AbortSignal, onEvent: (event: unknown) => void) => {
      onEvent({ type: 'text_delta', text: body.prompt === '你好' ? '你好！' : '可以继续帮助你。' });
      onEvent({ type: 'run_finished', status: 'SUCCESS' });
    });
    const { result } = renderHook(() => useAiSession('account:hq'));
    await act(async () => result.current.startPreview('你好'));
    await act(async () => result.current.startPreview('你还能做什么？'));
    expect(result.current.state.messages.map((message) => message.text)).toEqual(['你好', '你好！', '你还能做什么？', '可以继续帮助你。']);
    expect(mocks.readAiStream.mock.calls[1][1]).toEqual({ prompt: '你还能做什么？', history: [{ role: 'user', text: '你好' }, { role: 'assistant', text: '你好！' }] });
  });

  it('撤回最后一轮后保留更早对话供重新提问', async () => {
    mocks.readAiStream.mockImplementation(async (_path: string, body: { prompt: string }, _signal: AbortSignal, onEvent: (event: unknown) => void) => {
      if (body.prompt === '第二问') {
        onEvent({ type: 'tool_started', toolId: 'select_tools', title: '选择后台工具', sequence: 1 });
        onEvent({ type: 'tool_finished', toolId: 'select_tools', title: '选择后台工具', sequence: 1, status: 'SUCCESS' });
      }
      onEvent({ type: 'text_delta', text: `${body.prompt}的回答` });
      onEvent({ type: 'run_finished', status: 'SUCCESS' });
    });
    const { result } = renderHook(() => useAiSession('account:hq'));
    await act(async () => result.current.startPreview('第一问'));
    await act(async () => result.current.startPreview('第二问'));
    let restored: string | undefined;
    act(() => { restored = result.current.retractLastExchange(); });
    expect(restored).toBe('第二问');
    expect(result.current.state.messages.map((message) => message.text)).toEqual(['第一问', '第一问的回答']);
    expect(result.current.state.transcript).toEqual([
      { kind: 'message', role: 'user', text: '第一问' }, { kind: 'message', role: 'assistant', text: '第一问的回答' },
    ]);
    await act(async () => result.current.startPreview('改写的第二问'));
    expect(mocks.readAiStream.mock.calls[2][1]).toEqual({ prompt: '改写的第二问', history: [
      { role: 'user', text: '第一问' }, { role: 'assistant', text: '第一问的回答' },
    ] });
  });

  it('完成事件已到达但连接尚未关闭时仍可撤回', async () => {
    mocks.readAiStream.mockImplementation(async (_path: string, _body: object, signal: AbortSignal, onEvent: (event: unknown) => void) => {
      onEvent({ type: 'text_delta', text: '回复' });
      onEvent({ type: 'run_finished', status: 'SUCCESS' });
      await new Promise<void>((resolve) => signal.addEventListener('abort', () => resolve(), { once: true }));
    });
    const { result } = renderHook(() => useAiSession('account:hq'));
    act(() => { void result.current.startPreview('你好'); });
    await waitFor(() => expect(result.current.state.phase).toBe('completed'));
    expect(result.current.canRetractLastExchange).toBe(true);
    act(() => { result.current.retractLastExchange(); });
    expect(result.current.state.messages).toEqual([]);
  });
});

describe('全局 AI 会话', () => {
  it('预览收到批准令牌后只能确认执行一次', async () => {
    mocks.readAiStream.mockImplementation(async (path: string, _body: object, _signal: AbortSignal, onEvent: (event: unknown) => void) => {
      if (path === 'preview') {
        onEvent({ type: 'preview_ready', summary: '将停用组织', steps: [{ toolId: 'HqSuspendOrganization', title: '停用组织', arguments: { id: 'org-a' }, sequence: 1, maxCalls: 1, risk: 'HIGH' }], isHighRisk: true, requiredInputs: [], previewToken: 'one-use' });
      }
      onEvent({ type: 'run_finished', status: 'SUCCESS' });
    });
    const { result } = renderHook(() => useAiSession('account:hq'));
    await act(async () => result.current.startPreview('停用组织'));
    expect(result.current.state.phase).toBe('awaiting_confirmation');
    await act(async () => { void result.current.confirmRun(); void result.current.confirmRun(); });
    expect(mocks.readAiStream.mock.calls.filter(([path]) => path === 'run')).toHaveLength(1);
    expect(mocks.readAiStream.mock.calls[1][1]).toEqual({ previewToken: 'one-use' });
    expect(result.current.state.phase).toBe('completed');
  });

  it('执行流中断后销毁凭据并禁止重放', async () => {
    mocks.readAiStream.mockImplementation(async (path: string, _body: object, _signal: AbortSignal, onEvent: (event: unknown) => void) => {
      if (path === 'preview') {
        onEvent({ type: 'preview_ready', summary: '计划', steps: [], isHighRisk: false, requiredInputs: [], previewToken: 'one-use' });
        onEvent({ type: 'run_finished', status: 'SUCCESS' });
      } else throw new AIStreamError('AI_STREAM_INTERRUPTED');
    });
    const { result } = renderHook(() => useAiSession('account:hq'));
    await act(async () => result.current.startPreview('执行任务'));
    await act(async () => result.current.confirmRun());
    expect(result.current.state.interrupted).toBe(true);
    expect(result.current.state.preview?.previewToken).toBeUndefined();
    await act(async () => result.current.confirmRun());
    expect(mocks.readAiStream.mock.calls.filter(([path]) => path === 'run')).toHaveLength(1);
  });

  it('预览失败即使曾收到令牌也禁止确认', async () => {
    mocks.readAiStream.mockImplementation(async (_path: string, _body: object, _signal: AbortSignal, onEvent: (event: unknown) => void) => {
      onEvent({ type: 'preview_ready', summary: '计划', steps: [], isHighRisk: false, requiredInputs: [], previewToken: 'stale' });
      onEvent({ type: 'run_finished', status: 'FAILED' });
    });
    const { result } = renderHook(() => useAiSession('account:hq'));
    await act(async () => result.current.startPreview('执行任务'));
    await act(async () => result.current.confirmRun());
    expect(result.current.state.phase).toBe('failed');
    expect(mocks.readAiStream).toHaveBeenCalledTimes(1);
  });
});

describe('AI 工作区隔离', () => {
  it('执行刷新后只根据非敏感标记提示核对结果', async () => {
		mocks.readAiStream.mockImplementation(async (path: string, _body: object, _signal: AbortSignal, onEvent: (event: unknown) => void) => {
			if (path === 'preview') {
				onEvent({ type: 'preview_ready', summary: '计划', steps: [], isHighRisk: false, requiredInputs: [], previewToken: 'one-use' });
				onEvent({ type: 'run_finished', status: 'SUCCESS' });
			} else {
				onEvent({ type: 'run_started', runId: 'run-7', phase: 'RUN' });
				throw new AIStreamError('AI_STREAM_INTERRUPTED');
			}
		});
		const first = renderHook(() => useAiSession('account:hq'));
		await act(async () => first.result.current.startPreview('处理任务'));
		await act(async () => first.result.current.confirmRun());
		expect(sessionStorage.getItem('admin.ai.run:account:hq')).toBe('{"runId":"run-7","phase":"RUN"}');
		first.unmount();
		const refreshed = renderHook(() => useAiSession('account:hq'));
		expect(refreshed.result.current.state.phase).toBe('interrupted');
		expect(refreshed.result.current.state.preview).toBeUndefined();
	});

  it('工作区切换立即清除秘密并取消连接', async () => {
    let signal: AbortSignal | undefined;
    mocks.readAiStream.mockImplementation(async (_path: string, _body: object, given: AbortSignal, onEvent: (event: unknown) => void) => {
      signal = given;
      onEvent({ type: 'secret', toolId: 'HqResetAdministratorPassword', targetAccountId: 'account-a', value: 'temporary' });
      await new Promise<void>(() => undefined);
    });
    const { result, rerender } = renderHook(({ namespace }) => useAiSession(namespace), { initialProps: { namespace: 'account:hq' } });
    act(() => { void result.current.startPreview('reset'); });
    await waitFor(() => expect(result.current.state.secrets).toHaveLength(1));
    rerender({ namespace: 'account:franchise' });
    await waitFor(() => expect(result.current.state.secrets).toHaveLength(0));
    expect(signal?.aborted).toBe(true);
  });

	it('执行请求发出前先记录可能已消费的任务', async () => {
		mocks.readAiStream.mockImplementation(async (path: string, _body: object, _signal: AbortSignal, onEvent: (event: unknown) => void) => {
			if (path === 'preview') {
				onEvent({ type: 'preview_ready', summary: '计划', steps: [], isHighRisk: false, requiredInputs: [], previewToken: 'one-use' });
				onEvent({ type: 'run_finished', status: 'SUCCESS' });
			} else {
				throw new AIStreamError('AI_NETWORK_ERROR');
			}
		});
		const { result } = renderHook(() => useAiSession('account:hq'));
		await act(async () => result.current.startPreview('处理任务'));
		await act(async () => result.current.confirmRun());
		expect(result.current.state.interrupted).toBe(true);
		expect(sessionStorage.getItem('admin.ai.run:account:hq')).toBe('{"runId":"","phase":"RUN"}');
	});
});

it('切换工作区时清除包含完整手机号的对话', async () => {
  mocks.readAiStream.mockImplementation(async (_path: string, _body: object, _signal: AbortSignal, onEvent: (event: unknown) => void) => {
    onEvent({ type: 'text_delta', text: '会员完整手机号：13800138000' });
    onEvent({ type: 'run_finished', status: 'SUCCESS' });
  });
  const { result, rerender } = renderHook(({ namespace }) => useAiSession(namespace), { initialProps: { namespace: 'account:hq' } });
  await act(async () => result.current.startPreview('查询手机号'));
  expect(result.current.state.messages.some((message) => message.text.includes('13800138000'))).toBe(true);
  rerender({ namespace: 'account:franchise' });
  await waitFor(() => expect(result.current.state.messages).toEqual([]));
  expect(result.current.state.transcript).toEqual([]);
  expect(sessionStorage.getItem('admin.ai.run:account:hq')).toBeNull();
});

it('新工作区首次渲染不暴露旧工作区的完整手机号', async () => {
  mocks.readAiStream.mockImplementation(async (_path: string, _body: object, _signal: AbortSignal, onEvent: (event: unknown) => void) => {
    onEvent({ type: 'text_delta', text: '会员完整手机号：13800138000' });
    onEvent({ type: 'run_finished', status: 'SUCCESS' });
  });
  const rendered: Array<{ namespace: string; messages: string[] }> = [];
  const { result, rerender } = renderHook(({ namespace }) => {
    const session = useAiSession(namespace);
    useLayoutEffect(() => {
      rendered.push({ namespace, messages: session.state.messages.map((message) => message.text) });
    });
    return session;
  }, { initialProps: { namespace: 'account:hq' } });
  await act(async () => result.current.startPreview('查询手机号'));
  expect(result.current.state.messages.some((message) => message.text.includes('13800138000'))).toBe(true);
  rerender({ namespace: 'account:franchise' });
  expect(rendered.find((snapshot) => snapshot.namespace === 'account:franchise')?.messages).toEqual([]);
});
