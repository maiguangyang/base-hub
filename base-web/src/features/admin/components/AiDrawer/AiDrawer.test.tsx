// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';
import type { AiSessionState } from './useAiSession';
import { AiDrawer } from './AiDrawer';

beforeEach(() => { vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }); });
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

const preview: NonNullable<AiSessionState['preview']> = {
  type: 'preview_ready', summary: '将停用加盟组织', isHighRisk: true, requiredInputs: [], previewToken: 'one-use',
  steps: [{ toolId: 'HqSuspendOrganization', title: '停用加盟商组织', targetIds: ['org-a'], arguments: { input: { organizationId: 'org-a', reasonCode: 'CONTRACT_ENDED' } }, maxCalls: 1, sequence: 1, risk: 'HIGH' }],
};

function fixture(phase: AiSessionState['phase'] = 'awaiting_confirmation') {
  const startPreview = vi.fn(async () => undefined);
  const confirmRun = vi.fn(async () => undefined);
  const session = {
    state: { phase, messages: [], transcript: [], secrets: [], interrupted: false, preview } as AiSessionState,
    startPreview, confirmRun, stop: vi.fn(), clear: vi.fn(), retractLastExchange: vi.fn(), canRetractLastExchange: true,
  };
  return { session, startPreview, confirmRun };
}

function renderDrawer(session: ReturnType<typeof fixture>['session']) {
  return render(<AdminToastProvider><AiDrawer open onOpenChange={vi.fn()} session={session} /></AdminToastProvider>);
}

describe('全局 AI 抽屉', () => {
  it('累计模型额度耗尽时给出明确原因', () => {
    const { session } = fixture('failed');
    session.state.errorCode = 'AI_TOKEN_BUDGET_EXCEEDED';
    renderDrawer(session);
    expect(screen.getByRole('alert').textContent).toContain('本轮模型调用额度已用完');
    expect(screen.getByRole('alert').textContent).toContain('AI_TOKEN_BUDGET_EXCEEDED');
  });

  it('抽屉使用白色背景', () => {
    const { session } = fixture('completed');
    renderDrawer(session);
    expect(screen.getByRole('dialog', { name: 'AI 助手' }).className).toContain('bg-white');
  });

  it('高危计划展示完整参数，取消确认不执行，确认后只发起一次', () => {
    const { session, confirmRun } = fixture();
    renderDrawer(session);
    expect(screen.getByText('1. 停用加盟商组织')).toBeTruthy();
    expect(screen.getByText('CONTRACT_ENDED')).toBeTruthy();
    fireEvent.click(screen.getByRole('button', { name: '确认执行' }));
    expect(screen.getByRole('alertdialog')).toBeTruthy();
    fireEvent.click(screen.getByRole('button', { name: '取消' }));
    expect(confirmRun).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole('button', { name: '确认执行' }));
    fireEvent.click(screen.getByRole('button', { name: '执行已确认操作' }));
    expect(confirmRun).toHaveBeenCalledTimes(1);
  });

  it('没有预览令牌时无法执行写入', () => {
    const { session, confirmRun } = fixture();
    session.state.preview = { ...preview, previewToken: undefined };
    renderDrawer(session);
    expect(screen.queryByRole('button', { name: '确认执行' })).toBeNull();
    expect(confirmRun).not.toHaveBeenCalled();
  });

  it('处理中输入区显示终止按钮并停止当前流', () => {
    const { session } = fixture('previewing');
    renderDrawer(session);
    expect(screen.queryByRole('button', { name: '发送消息' })).toBeNull();
    fireEvent.click(screen.getByRole('button', { name: '终止生成' }));
    expect(session.stop).toHaveBeenCalledTimes(1);
  });
});

describe('AI 对话撤回', () => {
  it('撤回按钮聚焦时说明已执行操作不会撤销', async () => {
    const { session } = fixture('completed');
    session.state.transcript = [{ kind: 'message', role: 'user', text: '原来的问题' }];
    renderDrawer(session);
    fireEvent.focus(screen.getByRole('button', { name: '撤回最后一条消息' }));
    await waitFor(() => expect(screen.getAllByRole('tooltip').some((tooltip) => tooltip.textContent?.includes('撤回本轮对话；已执行的后台操作不会撤销'))).toBe(true));
  });

  it('撤回最后一轮后把原文回填到输入框', () => {
    const { session } = fixture('completed');
    session.state.messages = [{ role: 'user', text: '原来的问题' }, { role: 'assistant', text: '原来的回答' }];
    session.state.transcript = session.state.messages.map((message) => ({ kind: 'message', ...message }));
    session.retractLastExchange.mockReturnValue('原来的问题');
    renderDrawer(session);
    fireEvent.click(screen.getByRole('button', { name: '撤回最后一条消息' }));
    expect(session.retractLastExchange).not.toHaveBeenCalled();
    expect(screen.getByRole('alertdialog', { name: '确认撤回本轮对话' })).toBeTruthy();
    fireEvent.click(screen.getByRole('button', { name: '确认撤回' }));
    expect(session.retractLastExchange).toHaveBeenCalledTimes(1);
    expect((screen.getByRole('textbox', { name: '向 AI 助手发送消息' }) as HTMLTextAreaElement).value).toBe('原来的问题');
  });

  it('取消撤回时保留消息与输入', () => {
    const { session } = fixture('completed');
    session.state.messages = [{ role: 'user', text: '需要保留的问题' }, { role: 'assistant', text: '需要保留的回答' }];
    session.state.transcript = session.state.messages.map((message) => ({ kind: 'message', ...message }));
    renderDrawer(session);
    fireEvent.click(screen.getByRole('button', { name: '撤回最后一条消息' }));
    fireEvent.click(screen.getByRole('button', { name: '取消' }));
    expect(session.retractLastExchange).not.toHaveBeenCalled();
    expect(screen.getByText('需要保留的回答')).toBeTruthy();
  });

  it('撤回后清除上轮操作依据与确认状态', () => {
    const { session } = fixture();
    session.state.preview = { ...preview, requiredInputs: ['evidenceReference'] };
    session.state.messages = [{ role: 'user', text: '处理账号' }, { role: 'assistant', text: '需要依据' }];
    session.state.transcript = session.state.messages.map((message) => ({ kind: 'message', ...message }));
    session.retractLastExchange.mockReturnValue('处理账号');
    const { rerender } = renderDrawer(session);
    fireEvent.change(screen.getByPlaceholderText('请输入历史开通记录编号'), { target: { value: 'record-1' } });
    fireEvent.click(screen.getByRole('checkbox', { name: '我已核实这项认定依据' }));
    expect(screen.getByRole('checkbox', { name: '我已核实这项认定依据' }).getAttribute('aria-checked')).toBe('true');
    fireEvent.click(screen.getByRole('button', { name: '撤回最后一条消息' }));
    fireEvent.click(screen.getByRole('button', { name: '确认撤回' }));
    rerender(<AdminToastProvider><AiDrawer open onOpenChange={vi.fn()} session={session} /></AdminToastProvider>);
    expect((screen.getByPlaceholderText('请输入历史开通记录编号') as HTMLInputElement).value).toBe('');
    expect(screen.getByRole('checkbox', { name: '我已核实这项认定依据' }).getAttribute('aria-checked')).toBe('false');
  });
});
