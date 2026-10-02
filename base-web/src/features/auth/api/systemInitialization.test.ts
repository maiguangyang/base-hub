import { afterEach, describe, expect, it, vi } from 'vitest';
import {
  getSystemInitializationStatus, initializeSystem, SystemInitializationError,
} from './systemInitialization';

describe('系统初始化 API', () => {
  afterEach(() => vi.unstubAllGlobals());

  it('读取初始化状态时使用配置地址、凭据和 AbortSignal', async () => {
    const signal = new AbortController().signal;
    const fetchMock = vi.fn(async () => new Response(JSON.stringify({ initialized: false }), {
      status: 200, headers: { 'Content-Type': 'application/json' },
    }));
    vi.stubGlobal('fetch', fetchMock);

    await expect(getSystemInitializationStatus(signal)).resolves.toBe(false);
    expect(fetchMock).toHaveBeenCalledWith('http://localhost:1980/api/system-initialization', {
      credentials: 'include', headers: { Accept: 'application/json' }, signal,
    });
  });

  it('提交初始化数据时不增加额外字段并解析成功结果', async () => {
    const fetchMock = vi.fn(async () => new Response(JSON.stringify({ initialized: true }), {
      status: 201, headers: { 'Content-Type': 'application/json' },
    }));
    vi.stubGlobal('fetch', fetchMock);
    const input = { phone: '13800000000', password: 'Correct-Horse-42', passwordConfirmation: 'Correct-Horse-42' };

    await expect(initializeSystem(input)).resolves.toBeUndefined();
    expect(fetchMock).toHaveBeenCalledWith('http://localhost:1980/api/system-initialization', expect.objectContaining({
      method: 'POST', credentials: 'include', body: JSON.stringify(input),
      headers: { Accept: 'application/json', 'Content-Type': 'application/json' },
    }));
  });

  it('把失败响应收敛为稳定错误码', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => new Response(JSON.stringify({ code: 'PASSWORD_WEAK' }), {
      status: 400, headers: { 'Content-Type': 'application/json' },
    })));

    await expect(initializeSystem({ phone: '', password: '', passwordConfirmation: '' }))
      .rejects.toEqual(new SystemInitializationError('PASSWORD_WEAK', 400));
  });

  it('损坏的失败响应不会暴露响应正文', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => new Response('database secret', { status: 500 })));

    await expect(getSystemInitializationStatus()).rejects.toEqual(new SystemInitializationError('INTERNAL_ERROR', 500));
  });
});
