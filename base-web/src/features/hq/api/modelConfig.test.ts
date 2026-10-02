import { afterEach, describe, expect, it, vi } from 'vitest';
import { modelConfigAPI, ModelConfigRequestError } from './modelConfig';

afterEach(() => vi.unstubAllGlobals());

describe('模型配置接口', () => {
  it('只在保存请求发送密钥，读取与操作请求不带密钥', async () => {
    const fetchMock = vi.fn().mockResolvedValue({ ok: true, json: async () => ({ status: 'DRAFT', version: 1, keyConfigured: true }) });
    vi.stubGlobal('fetch', fetchMock);
    await modelConfigAPI.status();
    await modelConfigAPI.save({ modelName: 'model', baseUrl: 'https://model.example/v1', apiKey: 'private-key', version: 0 });
    await modelConfigAPI.probe(1);
    expect(fetchMock.mock.calls[0][1].body).toBeUndefined();
    expect(fetchMock.mock.calls[1][1].body).toContain('private-key');
    expect(fetchMock.mock.calls[2][1].body).toBe('{"version":1}');
    expect(fetchMock.mock.calls[1][1].credentials).toBe('include');
    expect(fetchMock.mock.calls[1][1].cache).toBe('no-store');
  });

  it('保留安全错误码供页面提示', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false, json: async () => ({ code: 'MODEL_PROBE_FAILED' }) }));
    await expect(modelConfigAPI.probe(1)).rejects.toEqual(new ModelConfigRequestError('MODEL_PROBE_FAILED'));
  });
});
