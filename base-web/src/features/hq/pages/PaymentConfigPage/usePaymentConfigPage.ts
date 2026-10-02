import { useEffect, useRef, useState } from 'react';
import { paymentConfigAPI, PaymentConfigRequestError, type ChannelView, type PaymentChannel, type SaveInput, type ScopeRef, type ScopeView } from '@/features/hq/api/paymentConfig';
import { publishActionFeedback } from '@/lib/actionFeedback';

const errors: Record<string, string> = {
  CONFLICT: '配置已被其他管理员修改，请重新加载最新数据后重试。',
  CONFIG_UNAVAILABLE: '服务端尚未配置支付加密密钥 PAYMENT_CONFIG_ENCRYPTION_KEYS。',
  PERMISSION_DENIED: '当前账号没有支付配置权限。',
  ORIGIN_NOT_ALLOWED: '当前后台地址未列入服务端允许来源。',
  VALIDATION_FAILED: '配置格式或目标范围无效，请检查后重试。',
  NETWORK_ERROR: '网络连接失败，请稍后重试。',
};

export type SaveDraft = Pick<SaveInput, 'merchantId' | 'environment' | 'ratePpm' | 'wechatCredentials' | 'alipayCredentials'>;

export function usePaymentConfigPage(ref: ScopeRef) {
  const [view, setView] = useState<ScopeView>();
  const [viewKey, setViewKey] = useState('');
  const [loading, setLoading] = useState(true);
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const [notice, setNotice] = useState<string>();
  const [conflict, setConflict] = useState(false);
  const requestSequence = useRef(0);
  const scopeKey = `${ref.scope}:${ref.organizationId ?? ''}:${ref.storeId ?? ''}`;
  const currentView = viewKey === scopeKey ? view : undefined;

  useEffect(() => {
    let current = true;
    const requestId = ++requestSequence.current;
    setError(undefined); setNotice(undefined); setConflict(false); setView(undefined); setViewKey('');
    setLoading(true);
    paymentConfigAPI.read(ref).then((next) => { if (current && requestId === requestSequence.current) { setView(next); setViewKey(scopeKey); } })
      .catch((cause: unknown) => { if (current && requestId === requestSequence.current) setError(paymentError(cause)); })
      .finally(() => { if (current && requestId === requestSequence.current) setLoading(false); });
    return () => { current = false; };
  }, [scopeKey]);

  async function reload() {
    const requestId = ++requestSequence.current;
    setLoading(true); setError(undefined); setConflict(false);
    try { const next = await paymentConfigAPI.read(ref); if (requestId === requestSequence.current) { setView(next); setViewKey(scopeKey); } }
    catch (cause) { if (requestId === requestSequence.current) setError(paymentError(cause)); }
    finally { if (requestId === requestSequence.current) setLoading(false); }
  }

  async function run(action: () => Promise<ScopeView>, success: string): Promise<boolean> {
    if (pending || loading) return false;
    setPending(true); setError(undefined); setNotice(undefined);
    try { setView(await action()); setViewKey(scopeKey); setConflict(false); setNotice(success); publishActionFeedback('success', success); return true; }
    catch (cause) { setConflict(cause instanceof PaymentConfigRequestError && cause.code === 'CONFLICT'); setError(paymentError(cause)); publishActionFeedback('error', `支付配置操作失败：${paymentError(cause)}`); return false; }
    finally { setPending(false); }
  }

  const own = (channel: PaymentChannel): ChannelView['own'] | undefined => currentView?.channels.find((item) => item.own.channel === channel)?.own;
  const save = (channel: PaymentChannel, draft: SaveDraft) => run(() => paymentConfigAPI.save({ ...ref, channel, ...draft,
    recordId: own(channel)?.recordId ?? '', version: own(channel)?.version ?? 0 }), '配置已保存并完成本地校验。');
  const setState = (channel: PaymentChannel, state: 'VALID' | 'DISABLED') => run(() => paymentConfigAPI.state({ ...ref, channel, state,
    recordId: own(channel)?.recordId ?? '', version: own(channel)?.version ?? 0 }), state === 'VALID' ? '此级配置已启用。' : '此级配置已停用。');
  const restoreInheritance = (channel: PaymentChannel) => run(() => paymentConfigAPI.restoreInheritance({ ...ref, channel,
    recordId: own(channel)?.recordId ?? '', version: own(channel)?.version ?? 0 }), '已恢复继承。');

  return { view: currentView, loading, pending, error, notice, conflict, reload, save, setState, restoreInheritance };
}

function paymentError(cause: unknown): string {
  const code = cause instanceof PaymentConfigRequestError ? cause.code : 'UNKNOWN';
  return errors[code] ?? '操作失败，请稍后重试。';
}
