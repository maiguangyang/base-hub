import { parse } from 'graphql';
import { afterEach, describe, expect, it, vi } from 'vitest';

const { createClientMock } = vi.hoisted(() => ({
  createClientMock: vi.fn<(options: unknown) => { dispose(): void; subscribe(): void }>(
    () => ({ dispose: vi.fn(), subscribe: vi.fn() }),
  ),
}));

vi.mock('graphql-ws', () => ({ createClient: createClientMock }));

import { createSessionResetController, disposeGraphQLRuntime, getGraphQLRuntime } from './client';

afterEach(async () => {
  await disposeGraphQLRuntime();
  createClientMock.mockClear();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe('GraphQL 连接配置', () => {
  it('HTTP 请求使用集中配置的 Engine 地址', async () => {
    const fetchMock = vi.fn<typeof fetch>(async () => new Response(JSON.stringify({ data: { __typename: 'Query' } }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    }));
    vi.stubGlobal('fetch', fetchMock);
    const runtime = getGraphQLRuntime({ resetAuth: vi.fn(), resetWorkspace: vi.fn() });

    await runtime.client.query({
      query: parse('query ConfigProbe { __typename }'),
      fetchPolicy: 'network-only',
    });

    expect(fetchMock.mock.calls[0]?.[0]).toBe('http://localhost:1980/graphql');
  });

  it('Subscription 使用集中配置的 Engine 地址', () => {
    getGraphQLRuntime({ resetAuth: vi.fn(), resetWorkspace: vi.fn() });

    expect(createClientMock).toHaveBeenCalledTimes(1);
    expect(createClientMock).toHaveBeenCalledWith(expect.objectContaining({
      url: 'ws://localhost:1980/graphql',
    }));
  });
});

describe('GraphQL runtime 清理', () => {
  it('Apollo 缓存清理失败也会完成本地会话终态重置', async () => {
    const diagnostics = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    const calls: string[] = [];
    const resetSession = createSessionResetController({
      disposeSubscription: () => calls.push('dispose'),
      clearStore: async () => { calls.push('clear'); throw new Error('cache cleanup failed'); },
      resetWorkspace: () => calls.push('workspace'),
      resetAuth: () => calls.push('auth'),
    });

    await expect(resetSession('LOGOUT')).resolves.toBeUndefined();
    expect(calls).toEqual(['dispose', 'clear', 'workspace', 'auth']);
    expect(diagnostics).toHaveBeenCalledWith(
      '会话终态的 Apollo 缓存清理失败，继续重置本地状态。',
      expect.any(Error),
    );
  });

  it('销毁 runtime 时 Apollo 缓存清理失败不会阻断认证后导航', async () => {
    const diagnostics = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    const runtime = getGraphQLRuntime({ resetAuth: vi.fn(), resetWorkspace: vi.fn() });
    vi.spyOn(runtime.client, 'clearStore').mockRejectedValueOnce(new Error('cache cleanup failed'));

    await expect(disposeGraphQLRuntime()).resolves.toBeUndefined();
    expect(diagnostics).toHaveBeenCalledWith(
      'GraphQL runtime 缓存清理失败，runtime 已销毁。',
      expect.any(Error),
    );
  });
});
