// @vitest-environment jsdom

import { act, cleanup, fireEvent, render, screen } from '@testing-library/react';
import { parse } from 'graphql';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';
import { disposeGraphQLRuntime, getGraphQLRuntime } from './client';

vi.mock('graphql-ws', () => ({ createClient: () => ({ dispose: vi.fn(), subscribe: vi.fn() }) }));
const mutation = parse('mutation HqReviewStore { reviewStore(input: {storeId: "store", approved: true}) { id lifecycle } }');
const response = (body: object) => new Response(JSON.stringify(body), { headers: { 'Content-Type': 'application/json' } });
const success = { data: { reviewStore: { id: 'store', lifecycle: 'ACTIVE' } } };
beforeEach(() => render(<AdminToastProvider><span>后台</span></AdminToastProvider>));
afterEach(async () => { cleanup(); await disposeGraphQLRuntime(); vi.unstubAllGlobals(); });

async function mutate(body: object, context?: object) {
  vi.stubGlobal('fetch', vi.fn(async () => response(body)));
  const runtime = getGraphQLRuntime({ resetAuth: vi.fn(), resetWorkspace: vi.fn() });
  await act(async () => { await runtime.client.mutate({ mutation, context }).catch(() => undefined); });
}

it('真实后台写请求成功后显示批准反馈', async () => {
  await mutate(success);
  expect(screen.getByRole('status').textContent).toContain('门店审核成功');
});

it.each([
  ['PERMISSION_DENIED', '当前账号没有执行此操作的权限'],
  ['VALIDATION_FAILED', '资料不符合要求'],
  ['CONFLICT', '数据或状态已发生变化'],
  ['INTERNAL_ERROR', '服务端处理异常'],
])('%s 错误显示中文原因并且不报成功', async (code, reason) => {
  await mutate({ data: null, errors: [{ message: code, extensions: { code } }] });
  expect(screen.getByRole('alert').textContent).toContain(reason);
  expect(screen.queryByRole('status')).toBeNull();
});

it('部分数据伴随错误也不会显示成功', async () => {
  await mutate({ ...success, errors: [{ message: 'CONFLICT', extensions: { code: 'CONFLICT' } }] });
  expect(screen.getByRole('alert').textContent).toContain('数据或状态已发生变化');
  expect(screen.queryByRole('status')).toBeNull();
});

it('请求连接失败时显示网络原因', async () => {
  vi.stubGlobal('fetch', vi.fn(async () => { throw new TypeError('Failed to fetch'); }));
  const runtime = getGraphQLRuntime({ resetAuth: vi.fn(), resetWorkspace: vi.fn() });
  await act(async () => { await runtime.client.mutate({ mutation }).catch(() => undefined); });
  expect(screen.getByRole('alert').textContent).toContain('网络连接失败');
});

it('服务端返回 false 或空结果时不报成功', async () => {
  vi.stubGlobal('fetch', vi.fn(async () => response({ data: { deleteStores: false } })));
  const runtime = getGraphQLRuntime({ resetAuth: vi.fn(), resetWorkspace: vi.fn() });
  await act(async () => { await runtime.client.mutate({ mutation: parse('mutation HqDeleteDirectStores { deleteStores(id:["store"]) }') }); });
  expect(screen.getByRole('alert').textContent).toContain('操作未完成');
});

it('查询不弹成功提示，分步保存可以自行汇总反馈', async () => {
  await mutate(success, { adminFeedback: false });
  expect(screen.queryByRole('status')).toBeNull();
  vi.stubGlobal('fetch', vi.fn(async () => response({ data: { __typename: 'Query' } })));
  const runtime = getGraphQLRuntime({ resetAuth: vi.fn(), resetWorkspace: vi.fn() });
  await act(async () => { await runtime.client.query({ query: parse('query Probe { __typename }'), fetchPolicy: 'network-only' }); });
  expect(screen.queryByRole('status')).toBeNull();
});

it('连续操作的反馈排队展示，不被后一个操作覆盖', async () => {
  await mutate(success);
  await mutate({ data: null, errors: [{ message: 'CONFLICT', extensions: { code: 'CONFLICT' } }] });
  expect(screen.getByRole('status').textContent).toContain('门店审核成功');
  fireEvent.click(screen.getByRole('button', { name: '关闭提示' }));
  expect(screen.getByRole('alert').textContent).toContain('数据或状态已发生变化');
});
