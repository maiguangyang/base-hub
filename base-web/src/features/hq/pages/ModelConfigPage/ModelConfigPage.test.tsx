// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { ModelConfigRequestError } from '@/features/hq/api/modelConfig';
import { ModelConfigPage } from './ModelConfigPage';

const mocks = vi.hoisted(() => ({
  workspace: 'HEADQUARTERS', permissions: ['aiModelConfig:read', 'aiModelConfig:manage'],
  status: vi.fn(), save: vi.fn(), probe: vi.fn(), activate: vi.fn(), deactivate: vi.fn(),
}));

vi.mock('@/features/auth/store/authStore', () => ({ useAuthStore: (selector: (state: unknown) => unknown) => selector({ viewer: {
  permissions: mocks.permissions, currentWorkspace: { workspaceType: mocks.workspace },
} }) }));
vi.mock('@/features/hq/api/modelConfig', () => ({
  modelConfigAPI: mocks,
  ModelConfigRequestError: class ModelConfigRequestError extends Error {
    constructor(readonly code: string) { super(code); }
  },
}));

beforeEach(() => {
  mocks.workspace = 'HEADQUARTERS';
  mocks.permissions = ['aiModelConfig:read', 'aiModelConfig:manage'];
  for (const name of ['status', 'save', 'probe', 'activate', 'deactivate'] as const) mocks[name].mockReset();
  mocks.status.mockResolvedValue({ modelName: 'example', baseUrl: 'https://model.example/v1', keyConfigured: true, version: 2, status: 'DRAFT' });
  mocks.save.mockResolvedValue({ modelName: 'changed', baseUrl: 'https://model.example/v1', keyConfigured: true, version: 3, status: 'DRAFT' });
  mocks.probe.mockResolvedValue({ modelName: 'example', baseUrl: 'https://model.example/v1', keyConfigured: true, version: 2, status: 'TESTED' });
  mocks.activate.mockResolvedValue({ modelName: 'example', baseUrl: 'https://model.example/v1', keyConfigured: true, version: 2, status: 'ACTIVE' });
  mocks.deactivate.mockResolvedValue({ modelName: 'example', baseUrl: 'https://model.example/v1', keyConfigured: true, version: 2, status: 'DRAFT' });
});
afterEach(() => cleanup());

describe('总部模型配置页面', () => {
  it('不回显密钥，测试成功后才能启用', async () => {
    render(<ModelConfigPage />);
    expect((await screen.findByLabelText('模型名称') as HTMLInputElement).value).toBe('example');
    expect((screen.getByLabelText('API Key') as HTMLInputElement).value).toBe('');
    const enabled = screen.getByRole('switch', { name: '是否启用' });
    expect(enabled.getAttribute('aria-checked')).toBe('false');
    expect(enabled.hasAttribute('disabled')).toBe(true);
    fireEvent.click(screen.getByRole('button', { name: '测试连接' }));
    await waitFor(() => expect(mocks.probe).toHaveBeenCalledWith(2));
    await waitFor(() => expect(enabled.hasAttribute('disabled')).toBe(false));
    fireEvent.click(enabled);
    await waitFor(() => expect(mocks.activate).toHaveBeenCalledWith(2));
    await waitFor(() => expect(enabled.getAttribute('aria-checked')).toBe('true'));
    fireEvent.click(enabled);
    await waitFor(() => expect(mocks.deactivate).toHaveBeenCalledWith(2));
    await waitFor(() => expect(enabled.getAttribute('aria-checked')).toBe('false'));
    expect(enabled.hasAttribute('disabled')).toBe(true);
  });

  it('修改后要求先保存，并允许留空沿用原密钥', async () => {
    render(<ModelConfigPage />);
    fireEvent.change(await screen.findByLabelText('模型名称'), { target: { value: 'changed' } });
    expect(screen.getByRole('button', { name: '测试连接' }).hasAttribute('disabled')).toBe(true);
    fireEvent.click(screen.getByRole('button', { name: '保存' }));
    await waitFor(() => expect(mocks.save).toHaveBeenCalledWith({ modelName: 'changed', baseUrl: 'https://model.example/v1', apiKey: '', version: 2 }));
  });
});

describe('模型配置保存与权限', () => {
  it('配置未改变时仍能保存，以便沿用 API Key 轮换主密钥', async () => {
    mocks.status.mockResolvedValue({ modelName: 'example', baseUrl: 'https://model.example/v1', keyConfigured: true, version: 2, status: 'ACTIVE' });
    mocks.save.mockResolvedValue({ modelName: 'example', baseUrl: 'https://model.example/v1', keyConfigured: true, version: 2, status: 'ACTIVE' });
    render(<ModelConfigPage />);
    const save = await screen.findByRole('button', { name: '保存' });
    const enabled = screen.getByRole('switch', { name: '是否启用' });
    expect(save.hasAttribute('disabled')).toBe(false);
    expect(enabled.getAttribute('aria-checked')).toBe('true');
    fireEvent.click(save);
    await waitFor(() => expect(mocks.save).toHaveBeenCalledWith({ modelName: 'example', baseUrl: 'https://model.example/v1', apiKey: '', version: 2 }));
    expect(enabled.getAttribute('aria-checked')).toBe('true');
  });

  it('加盟商工作区不能打开页面', () => {
    mocks.workspace = 'FRANCHISE';
    render(<ModelConfigPage />);
    expect(screen.getByRole('alert').textContent).toContain('无权查看');
    expect(mocks.status).not.toHaveBeenCalled();
  });

  it('来源未配置时提示地址问题，而不是误报账号权限', async () => {
    mocks.save.mockRejectedValue(new ModelConfigRequestError('ORIGIN_NOT_ALLOWED'));
    render(<ModelConfigPage />);
    fireEvent.change(await screen.findByLabelText('模型名称'), { target: { value: 'changed' } });
    await waitFor(() => expect(screen.getByRole('button', { name: '保存' }).hasAttribute('disabled')).toBe(false));
    fireEvent.click(screen.getByRole('button', { name: '保存' }));
    expect((await screen.findByRole('alert')).textContent).toContain('后台地址');
  });

  it('缺少服务端加密密钥时提示实际配置项', async () => {
    mocks.save.mockRejectedValue(new ModelConfigRequestError('CONFIG_UNAVAILABLE'));
    render(<ModelConfigPage />);
    fireEvent.change(await screen.findByLabelText('模型名称'), { target: { value: 'changed' } });
    await waitFor(() => expect(screen.getByRole('button', { name: '保存' }).hasAttribute('disabled')).toBe(false));
    fireEvent.click(screen.getByRole('button', { name: '保存' }));
    expect((await screen.findByRole('alert')).textContent).toContain('AI_MODEL_ENCRYPTION_KEYS');
  });

  it('输入自建 HTTP 地址时提示明文传输风险', async () => {
    render(<ModelConfigPage />);
    const input = await screen.findByLabelText('模型地址');
    fireEvent.change(input, { target: { value: 'http://zsgw.sjdistributor.com:4000/v1' } });
    expect(screen.getByText(/HTTP 会明文传输 API Key/)).toBeTruthy();
  });
});

it('停用模型时保留尚未保存的新密钥', async () => {
  mocks.status.mockResolvedValue({ modelName: 'example', baseUrl: 'https://model.example/v1', keyConfigured: true, version: 2, status: 'ACTIVE' });
  render(<ModelConfigPage />);
  const key = await screen.findByLabelText('API Key') as HTMLInputElement;
  fireEvent.change(key, { target: { value: 'replacement-key' } });
  fireEvent.click(screen.getByRole('switch', { name: '是否启用' }));
  await waitFor(() => expect(mocks.deactivate).toHaveBeenCalledWith(2));
  await waitFor(() => expect(screen.getByRole('switch', { name: '是否启用' }).getAttribute('aria-checked')).toBe('false'));
  expect(key.value).toBe('replacement-key');
});

describe('模型连接测试反馈', () => {
  it.each([
    ['MODEL_UPSTREAM_AUTH', 'API Key'],
    ['MODEL_TARGET_NOT_FOUND', '模型地址和名称'],
    ['MODEL_UPSTREAM_RATE_LIMITED', '供应商配额'],
    ['MODEL_REQUEST_REJECTED', 'Chat Completions'],
    ['MODEL_UPSTREAM_UNAVAILABLE', '暂时不可用'],
    ['MODEL_CONNECTION_FAILED', '网络和 TLS'],
    ['MODEL_PROTOCOL_UNSUPPORTED', 'Token 用量'],
  ])('对 %s 显示针对性提示', async (code, guidance) => {
    mocks.probe.mockRejectedValue(new ModelConfigRequestError(code));
    render(<ModelConfigPage />);
    fireEvent.click(await screen.findByRole('button', { name: '测试连接' }));
    await waitFor(() => expect(screen.getByRole('alert').textContent).toContain(guidance));
  });

  it('测试连接期间显示加载状态并阻止重复提交', async () => {
    let finishProbe!: (value: { modelName: string; baseUrl: string; keyConfigured: boolean; version: number; status: string }) => void;
    mocks.probe.mockImplementation(() => new Promise((resolve) => { finishProbe = resolve; }));
    render(<ModelConfigPage />);
    fireEvent.click(await screen.findByRole('button', { name: '测试连接' }));

    const testingButton = await screen.findByRole('button', { name: '测试中…' });
    expect(testingButton.hasAttribute('disabled')).toBe(true);
    expect(testingButton.querySelector('svg.animate-spin')).toBeTruthy();
    expect(mocks.probe).toHaveBeenCalledTimes(1);

    finishProbe({ modelName: 'example', baseUrl: 'https://model.example/v1', keyConfigured: true, version: 2, status: 'TESTED' });
    await waitFor(() => expect(screen.getByRole('button', { name: '测试连接' }).hasAttribute('disabled')).toBe(true));
    expect(screen.getByRole('switch', { name: '是否启用' }).hasAttribute('disabled')).toBe(false);
  });

  it('测试失败不能开启', async () => {
    mocks.probe.mockRejectedValue(new ModelConfigRequestError('MODEL_PROBE_FAILED'));
    render(<ModelConfigPage />);
    const enabled = await screen.findByRole('switch', { name: '是否启用' });
    fireEvent.click(screen.getByRole('button', { name: '测试连接' }));
    await screen.findByText(/测试未通过/);
    expect(enabled.hasAttribute('disabled')).toBe(true);
    expect(mocks.activate).not.toHaveBeenCalled();
  });

  it('测试通过后编辑未保存的配置时不能开启', async () => {
    mocks.status.mockResolvedValue({ modelName: 'example', baseUrl: 'https://model.example/v1', keyConfigured: true, version: 2, status: 'TESTED' });
    render(<ModelConfigPage />);
    const enabled = await screen.findByRole('switch', { name: '是否启用' });
    await waitFor(() => expect(enabled.hasAttribute('disabled')).toBe(false));
    fireEvent.change(screen.getByLabelText('模型名称'), { target: { value: 'changed' } });
    expect(enabled.hasAttribute('disabled')).toBe(true);
    expect(mocks.activate).not.toHaveBeenCalled();
  });
});
