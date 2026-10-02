import { afterEach, describe, expect, it, vi } from 'vitest';
import { formatPpmPercent, parsePercentToPpm, paymentConfigAPI, PaymentConfigRequestError } from './paymentConfig';

afterEach(() => vi.unstubAllGlobals());

describe('支付配置接口', () => {
	it('精确换算预计手续费率', () => {
		expect(parsePercentToPpm('0.38')).toBe(3800);
		expect(parsePercentToPpm('100')).toBe(1000000);
		expect(formatPpmPercent(3800)).toBe('0.38');
		expect(formatPpmPercent(1000000)).toBe('100');
		for (const value of ['0.12345', '100.0001', '-1', '1e2']) expect(() => parsePercentToPpm(value)).toThrow();
	});

	it('读取只带范围参数，保存携带版本和凭证且不缓存', async () => {
		const fetchMock = vi.fn().mockResolvedValue({ ok: true, json: async () => ({ channels: [] }) });
		vi.stubGlobal('fetch', fetchMock);
		await paymentConfigAPI.read({ scope: 'STORE', storeId: 'shop' });
		await paymentConfigAPI.save({ scope: 'GLOBAL', channel: 'WECHAT', merchantId: '1234567890', ratePpm: 3800, recordId: '', version: 0, wechatCredentials: {
			merchantApiSerial: '2A', merchantApiCert: 'CERTIFICATE', apiV3Key: 'private', merchantPrivateKey: 'PRIVATE KEY', verificationMode: 'PLATFORM_CERT',
		} });
		expect(fetchMock.mock.calls[0][0]).toContain('scope=STORE&storeId=shop');
		expect(fetchMock.mock.calls[0][1].body).toBeUndefined();
		expect(fetchMock.mock.calls[1][1].body).toContain('PRIVATE KEY');
		expect(fetchMock.mock.calls[1][1].credentials).toBe('include');
		expect(fetchMock.mock.calls[1][1].cache).toBe('no-store');
	});

	it('向页面保留 409 错误码', async () => {
		vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false, json: async () => ({ code: 'CONFLICT' }) }));
		await expect(paymentConfigAPI.read({ scope: 'GLOBAL' })).rejects.toEqual(new PaymentConfigRequestError('CONFLICT'));
	});
});
