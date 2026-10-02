// @vitest-environment jsdom

import { act, cleanup, renderHook, waitFor } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { subscribeActionFeedback } from '@/lib/actionFeedback';
import { PaymentConfigRequestError } from '../api/paymentConfig';
import { ModelConfigRequestError } from '../api/modelConfig';
import { usePaymentConfigPage } from './PaymentConfigPage/usePaymentConfigPage';
import { useModelConfigPage } from './ModelConfigPage/useModelConfigPage';

const mocks = vi.hoisted(() => {
  const view = { channels: [] };
  const status = { modelName: 'model', baseUrl: 'https://model.example', keyConfigured: true, version: 1, status: 'DRAFT' };
  return {
    payment: { read: vi.fn(async () => view), save: vi.fn(async () => view), state: vi.fn(async () => view), restoreInheritance: vi.fn(async () => view) },
    model: { status: vi.fn(async () => status), save: vi.fn(async () => status), probe: vi.fn(async () => status), activate: vi.fn(async () => status), deactivate: vi.fn(async () => status) },
  };
});
vi.mock('../api/paymentConfig', async (original) => ({ ...await original<object>(), paymentConfigAPI: mocks.payment }));
vi.mock('../api/modelConfig', async (original) => ({ ...await original<object>(), modelConfigAPI: mocks.model }));
afterEach(() => { cleanup(); vi.clearAllMocks(); });

it('支付配置保存、启停及恢复继承均有反馈，失败包含原因', async () => {
  const listener = vi.fn(); const unsubscribe = subscribeActionFeedback(listener);
  try {
    const { result } = renderHook(() => usePaymentConfigPage({ scope: 'GLOBAL' }));
    await waitFor(() => expect(result.current.loading).toBe(false));
    await act(async () => { await result.current.save('WECHAT', { merchantId: 'merchant', ratePpm: 0 }); });
    expect(listener).toHaveBeenLastCalledWith(expect.objectContaining({ kind: 'success', message: expect.stringContaining('配置已保存') }));
    await act(async () => { await result.current.setState('WECHAT', 'VALID'); });
    expect(listener).toHaveBeenLastCalledWith(expect.objectContaining({ message: '此级配置已启用。' }));
    await act(async () => { await result.current.restoreInheritance('WECHAT'); });
    expect(listener).toHaveBeenLastCalledWith(expect.objectContaining({ message: '已恢复继承。' }));
    mocks.payment.state.mockRejectedValueOnce(new PaymentConfigRequestError('PERMISSION_DENIED'));
    await act(async () => { await result.current.setState('WECHAT', 'DISABLED'); });
    expect(listener).toHaveBeenLastCalledWith(expect.objectContaining({ kind: 'error', message: expect.stringContaining('没有支付配置权限') }));
  } finally { unsubscribe(); }
});

it('模型配置保存、测试和启停有反馈，失败包含原因', async () => {
  const listener = vi.fn(); const unsubscribe = subscribeActionFeedback(listener);
  try {
    const { result } = renderHook(useModelConfigPage);
    await waitFor(() => expect(result.current.loading).toBe(false));
    await act(async () => { await result.current.save(); });
    expect(listener).toHaveBeenLastCalledWith(expect.objectContaining({ kind: 'success', message: '已保存。' }));
    await act(async () => { await result.current.probe(); });
    expect(listener).toHaveBeenLastCalledWith(expect.objectContaining({ message: expect.stringContaining('测试通过') }));
    await act(async () => { await result.current.activate(); });
    expect(listener).toHaveBeenLastCalledWith(expect.objectContaining({ message: '模型已启用。' }));
    await act(async () => { await result.current.deactivate(); });
    expect(listener).toHaveBeenLastCalledWith(expect.objectContaining({ message: expect.stringContaining('模型已停用') }));
    mocks.model.probe.mockRejectedValueOnce(new ModelConfigRequestError('MODEL_UPSTREAM_AUTH'));
    await act(async () => { await result.current.probe(); });
    expect(listener).toHaveBeenLastCalledWith(expect.objectContaining({ kind: 'error', message: expect.stringContaining('拒绝鉴权') }));
  } finally { unsubscribe(); }
});
