import { appConfig } from '@/config/app';

export type PaymentScope = 'GLOBAL' | 'FRANCHISE' | 'STORE';
export type PaymentChannel = 'WECHAT' | 'ALIPAY';
export type PaymentState = 'UNCONFIGURED' | 'VALID' | 'DISABLED' | 'ERROR';

export interface ScopeRef { scope: PaymentScope; organizationId?: string; storeId?: string }
export interface ConfigStatus {
  scope: PaymentScope; channel: PaymentChannel; sourceScope?: PaymentScope; sourceId?: string;
  recordId?: string; merchantMasked?: string; ratePpm: number; state: PaymentState;
  version: number; credentialsConfigured: boolean; reasonCode?: string;
}
export interface ChannelView { own: ConfigStatus; effective: ConfigStatus }
export interface ScopeView { channels: ChannelView[] }

export interface WechatCredentials {
  merchantApiSerial: string; merchantApiCert: string; apiV3Key: string;
  merchantPrivateKey: string; verificationMode: 'PUBLIC_KEY' | 'PLATFORM_CERT';
  wechatPublicKeyId?: string; wechatPublicKey?: string;
}
export interface AlipayCredentials {
  appPrivateKey: string; alipayPublicKey: string;
}
export interface SaveInput extends ScopeRef {
  channel: PaymentChannel; merchantId: string; environment?: 'PRODUCTION' | 'SANDBOX';
  ratePpm: number; recordId: string; version: number;
  wechatCredentials?: WechatCredentials; alipayCredentials?: AlipayCredentials;
}
export interface StateInput extends ScopeRef {
  channel: PaymentChannel; recordId: string; version: number; state: 'VALID' | 'DISABLED';
}
export interface ResetInput extends ScopeRef { channel: PaymentChannel; recordId: string; version: number }

export class PaymentConfigRequestError extends Error {
  constructor(readonly code: string) { super(code); }
}

const url = `${appConfig.engine.httpUrl}/api/payment-config`;

async function request(path: string, method: 'GET' | 'PUT' | 'POST', body?: object): Promise<ScopeView> {
  let response: Response;
  try {
    response = await fetch(`${url}${path}`, {
      method, credentials: 'include', cache: 'no-store',
      headers: { Accept: 'application/json', ...(body ? { 'Content-Type': 'application/json' } : {}) },
      ...(body ? { body: JSON.stringify(body) } : {}),
    });
  } catch { throw new PaymentConfigRequestError('NETWORK_ERROR'); }
  if (!response.ok) {
    const payload = await response.json().catch(() => null) as { code?: string } | null;
    throw new PaymentConfigRequestError(payload?.code ?? 'UNKNOWN');
  }
  return response.json() as Promise<ScopeView>;
}

export const paymentConfigAPI = {
  read: (ref: ScopeRef) => request(`?${new URLSearchParams({ scope: ref.scope, ...(ref.organizationId ? { organizationId: ref.organizationId } : {}), ...(ref.storeId ? { storeId: ref.storeId } : {}) })}`, 'GET'),
  save: (input: SaveInput) => request('', 'PUT', input),
  state: (input: StateInput) => request('/state', 'POST', input),
  restoreInheritance: (input: ResetInput) => request('/restore-inheritance', 'POST', input),
};

export function parsePercentToPpm(raw: string): number {
  const value = raw.trim();
  if (!/^(?:\d{1,2}|100)(?:\.\d{1,4})?$/.test(value)) throw new Error('INVALID_RATE');
  const [whole, fraction = ''] = value.split('.');
  const ppm = Number(whole) * 10000 + Number(fraction.padEnd(4, '0'));
  if (ppm > 1000000) throw new Error('INVALID_RATE');
  return ppm;
}

export function formatPpmPercent(ppm: number): string {
  const whole = Math.floor(ppm / 10000);
  const fraction = String(ppm % 10000).padStart(4, '0').replace(/0+$/, '');
  return fraction ? `${whole}.${fraction}` : String(whole);
}
