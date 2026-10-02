import { describe, expect, it, vi } from 'vitest';
import { createSessionResetController } from './client';
import { getGraphQLErrorCode, isTerminalSessionCode } from './errors';

describe('GraphQL 错误码与会话重置', () => {
  it('只读取 extensions.code，不解析 message 文本', () => {
    expect(getGraphQLErrorCode({ extensions: { code: 'PERMISSION_DENIED' } })).toBe('PERMISSION_DENIED');
    expect(getGraphQLErrorCode({ message: 'SESSION_REVOKED' })).toBeUndefined();
    expect(getGraphQLErrorCode({ extensions: { code: 42 } })).toBeUndefined();
  });

  it('仅把四种终态码识别为会话终止', () => {
    expect(['AUTH_REQUIRED', 'SESSION_REVOKED', 'ORGANIZATION_SUSPENDED', 'CREDENTIALS_CHANGED'].every(isTerminalSessionCode)).toBe(true);
    expect(isTerminalSessionCode('PERMISSION_DENIED')).toBe(false);
  });

  it('并发终态只销毁连接、缓存和认证状态一次', async () => {
    const disposeSubscription = vi.fn();
    const clearStore = vi.fn(async () => undefined);
    const resetAuth = vi.fn();
    const resetWorkspace = vi.fn();
    const reset = createSessionResetController({ disposeSubscription, clearStore, resetWorkspace, resetAuth });
    await Promise.all([reset('SESSION_REVOKED'), reset('CREDENTIALS_CHANGED')]);
    expect(disposeSubscription).toHaveBeenCalledTimes(1);
    expect(clearStore).toHaveBeenCalledTimes(1);
    expect(resetWorkspace).toHaveBeenCalledTimes(1);
    expect(resetAuth).toHaveBeenCalledTimes(1);
  });
});
