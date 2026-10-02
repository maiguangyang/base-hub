import { useEffect, useState } from 'react';
import { modelConfigAPI, ModelConfigRequestError, type ModelConfigStatus } from '@/features/hq/api/modelConfig';
import { publishActionFeedback } from '@/lib/actionFeedback';

const errorMessages: Record<string, string> = {
  CONFLICT: '配置已被其他管理员更新，请刷新后重试。',
  CONFIG_UNAVAILABLE: 'Engine 尚未配置模型加密密钥 AI_MODEL_ENCRYPTION_KEYS，请先完成服务端配置。',
  MODEL_PROBE_FAILED: '测试未通过。请检查模型服务是否支持流式输出和连续工具调用。',
  MODEL_UPSTREAM_AUTH: '模型服务拒绝鉴权，请检查 API Key 和账号权限。',
  MODEL_TARGET_NOT_FOUND: '模型服务未找到接口或模型，请检查模型地址和名称。',
  MODEL_UPSTREAM_RATE_LIMITED: '模型服务触发限流，请检查供应商配额后重试。',
  MODEL_REQUEST_REJECTED: '模型服务拒绝请求，请检查是否支持 Chat Completions、流式输出和工具调用。',
  MODEL_UPSTREAM_UNAVAILABLE: '模型服务暂时不可用，请稍后重试。',
  MODEL_CONNECTION_FAILED: '无法连接模型服务，请检查模型地址、网络和 TLS 配置。',
  MODEL_PROTOCOL_UNSUPPORTED: '模型响应不符合当前接入要求，请检查流式输出、工具调用和 Token 用量。',
  MODEL_PROBE_TIMEOUT: '模型服务响应超时，请稍后重试。',
  RATE_LIMITED: '测试过于频繁，请稍后重试。',
  PERMISSION_DENIED: '当前账号没有配置权限。',
  ORIGIN_NOT_ALLOWED: '当前后台地址未列入服务端允许来源，请检查 Engine 的 ALLOWED_ORIGINS 配置。',
  VALIDATION_FAILED: '请检查模型名称、HTTP/HTTPS 地址和密钥。',
  NETWORK_ERROR: '网络连接失败，请检查网络后重试。',
};

export function useModelConfigPage() {
  const [status, setStatus] = useState<ModelConfigStatus>();
  const [modelName, setModelName] = useState('');
  const [baseUrl, setBaseUrl] = useState('');
  const [apiKey, setApiKey] = useState('');
  const [pendingAction, setPendingAction] = useState<'save' | 'probe' | 'activate' | 'deactivate'>();
  const pending = pendingAction !== undefined;
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>();
  const [notice, setNotice] = useState<string>();

  useEffect(() => {
    let current = true;
    modelConfigAPI.status().then((value) => {
      if (!current) return;
      setStatus(value);
      setModelName(value.modelName);
      setBaseUrl(value.baseUrl);
    }).catch(() => { if (current) setError('模型配置加载失败，请刷新页面。'); })
      .finally(() => { if (current) setLoading(false); });
    return () => { current = false; };
  }, []);

  async function run(kind: NonNullable<typeof pendingAction>, action: () => Promise<ModelConfigStatus>, success: string) {
    if (pending || loading) return;
    setPendingAction(kind);
    setError(undefined);
    setNotice(undefined);
    try {
      setStatus(await action());
      if (kind === 'save') setApiKey('');
      setNotice(success);
      publishActionFeedback('success', success);
    } catch (cause) {
      const code = cause instanceof ModelConfigRequestError ? cause.code : 'UNKNOWN';
      const message = errorMessages[code] ?? '操作失败，请稍后重试。';
      setError(message);
      publishActionFeedback('error', `模型配置操作失败：${message}`);
    } finally {
      setPendingAction(undefined);
    }
  }

  const version = status?.version ?? 0;
  const save = () => run('save', () => modelConfigAPI.save({ modelName: modelName.trim(), baseUrl: baseUrl.trim(), apiKey, version }), '已保存。');
  const probe = () => run('probe', () => modelConfigAPI.probe(version), '测试通过，可以启用模型。');
  const activate = () => run('activate', () => modelConfigAPI.activate(version), '模型已启用。');
  const deactivate = () => run('deactivate', () => modelConfigAPI.deactivate(version), '模型已停用。再次启用前需要重新测试。');

  return { status, modelName, setModelName, baseUrl, setBaseUrl, apiKey, setApiKey, pending, testing: pendingAction === 'probe', loading, error, notice, save, probe, activate, deactivate };
}
