import { AdminStatusSwitch } from '@/features/admin/components/AdminStatusSwitch';
import { AdminFormSelect } from '@/features/admin/components/AdminFormSelect';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import { formatPpmPercent, parsePercentToPpm, type ChannelView, type PaymentChannel, type PaymentScope, type WechatCredentials, type AlipayCredentials } from '@/features/hq/api/paymentConfig';
import type { SaveDraft } from './usePaymentConfigPage';

const titles: Record<PaymentChannel, string> = { WECHAT: '微信', ALIPAY: '支付宝' };
const sources: Record<PaymentScope, string> = { GLOBAL: '全局', FRANCHISE: '加盟商', STORE: '门店' };

interface Props {
  channel: PaymentChannel; item: ChannelView; scope: PaymentScope; canManage: boolean; busy: boolean;
  form: ChannelForm; error?: string; onFormChange(form: ChannelForm): void;
  onState(channel: PaymentChannel, state: 'VALID' | 'DISABLED'): Promise<boolean>;
  onRestore(channel: PaymentChannel): Promise<boolean>;
}

export interface ChannelForm {
  merchantId: string; rate: string; environment: '' | 'PRODUCTION' | 'SANDBOX';
  wechat: WechatCredentials; alipay: AlipayCredentials; dirty: boolean;
}

export function emptyForm(item: ChannelView): ChannelForm {
  return { merchantId: '', rate: formatPpmPercent(item.own.ratePpm), environment: '', wechat: emptyWechat(), alipay: emptyAlipay(), dirty: false };
}

export function PaymentChannelCard(props: Props) {
  const { channel, item, scope, canManage } = props;
  const name = titles[channel];
  const effective = item.effective;
  return <Card>
    <CardHeader>
      <div className="flex items-center justify-between gap-3"><CardTitle>{name}支付</CardTitle><Badge variant="outline">{stateLabelFor(effective.state)}</Badge></div>
      <CardDescription>{name} · 生效来源：{effective.sourceScope ? sources[effective.sourceScope] : '未配置'}。此页只做配置与本地校验，不验证实际交易。</CardDescription>
    </CardHeader>
    <CardContent className="flex flex-col gap-5">
      <div className="grid gap-2 text-sm sm:grid-cols-2">
        <p>当前范围：{stateLabelFor(item.own.state)}{item.own.reasonCode === 'DECRYPTION_FAILED' ? ' · 解密失败' : item.own.reasonCode === 'CREDENTIALS_INVALID' ? ' · 凭证无效或已过期' : ''}</p>
        <p>生效配置：{effective.merchantMasked || '未填写商户'} · {formatPpmPercent(effective.ratePpm)}%</p>
        <p>凭证：{item.own.credentialsConfigured ? '已配置；留空沿用现有凭证' : '未配置'}</p>
        {scope !== 'GLOBAL' && <p>本级记录：{item.own.recordId ? '已设置覆盖' : '继承上级'}</p>}
      </div>
      {canManage && <PaymentChannelEditor {...props} />}
    </CardContent>
  </Card>;
}

function PaymentChannelEditor(props: Props) {
  const { channel, item, busy, form, onFormChange } = props;
  const name = titles[channel];
  const update = (patch: Partial<ChannelForm>) => onFormChange({ ...form, ...patch, dirty: true });

  return <div className="flex flex-col gap-5 border-t pt-5">
    <div className="grid gap-4 sm:grid-cols-2">
      <div className="flex flex-col gap-2"><Label htmlFor={`${channel}-merchant`}>{name}商户标识</Label><Input id={`${channel}-merchant`} value={form.merchantId} onChange={(event) => update({ merchantId: event.target.value })} placeholder={item.own.merchantMasked || (channel === 'WECHAT' ? '商户号 mchId' : '应用 ID appId')} disabled={busy} /></div>
      <div className="flex flex-col gap-2"><Label htmlFor={`${channel}-rate`}>{name}预计平台费率（%）</Label><Input id={`${channel}-rate`} inputMode="decimal" value={form.rate} onChange={(event) => update({ rate: event.target.value })} placeholder="例如 0.38" disabled={busy} /></div>
    </div>
    {channel === 'WECHAT' ? <WechatFields value={form.wechat} setValue={(wechat) => update({ wechat })} busy={busy} /> : <AlipayFields value={form.alipay} setValue={(alipay) => update({ alipay })} environment={form.environment} setEnvironment={(environment) => update({ environment })} busy={busy} />}
    {props.error && <p role="alert" className="text-sm text-destructive">{props.error}</p>}
    <PaymentChannelActions {...props} />
  </div>;
}

function PaymentChannelActions(props: Props) {
  const { channel, item, scope, busy } = props;
  const name = titles[channel];
  const enabled = item.effective.state === 'VALID';
  const canEnable = item.own.state === 'DISABLED' && item.own.credentialsConfigured;
  return <div className="flex flex-col gap-4 border-t pt-5">
    <div className="flex items-center justify-between gap-3">
      <span className="text-sm">{name}支付启用状态</span>
      <AdminStatusSwitch checked={enabled} label={`${name}支付启用状态`} disabled={busy || props.form.dirty || (!enabled && !canEnable)} onCheckedChange={(next) => void props.onState(channel, next ? 'VALID' : 'DISABLED')} />
    </div>
    {props.form.dirty && <p className="text-xs text-muted-foreground">请先保存当前支付方式的修改，再切换启用状态。</p>}
    {scope !== 'GLOBAL' && item.own.recordId && <div><Button type="button" variant="outline" disabled={busy || props.form.dirty} onClick={() => void props.onRestore(channel)}>恢复继承</Button></div>}
  </div>;
}

export function buildDraft(channel: PaymentChannel, item: ChannelView, form: ChannelForm): SaveDraft | string {
  const { merchantId, rate, environment, wechat, alipay } = form;
  let ratePpm: number;
  try { ratePpm = parsePercentToPpm(rate); } catch { return '费率请输入 0–100 的百分比，最多四位小数。'; }
  const credentials = channel === 'WECHAT' ? hasWechatInput(wechat) : hasAlipayInput(alipay);
  const error = validateDraft(channel, item, merchantId, environment, credentials, alipay);
  if (error) return error;
  return { merchantId: merchantId.trim(), ratePpm, ...credentialDraft(channel, credentials, environment, wechat, alipay) };
}

function validateDraft(channel: PaymentChannel, item: ChannelView, merchantId: string, environment: string, credentials: boolean, alipay: AlipayCredentials): string | undefined {
  const alipayError = channel === 'ALIPAY' ? validateAlipayInput(credentials, environment, alipay) : undefined;
  if (alipayError) return alipayError;
  if (item.own.state === 'ERROR' && !credentials) return '当前凭证无效，请重新填写完整凭证。';
  if (!item.own.credentialsConfigured && !credentials) return '请填写完整的支付凭证。';
  if (!item.own.recordId && !merchantId.trim()) return '请填写商户标识。';
  return undefined;
}

function validateAlipayInput(credentials: boolean, environment: string, alipay: AlipayCredentials): string | undefined {
  if (credentials && !environment) return '请明确选择支付宝正式或沙箱环境。';
  if (credentials && (!alipay.appPrivateKey.trim() || !alipay.alipayPublicKey.trim())) return '请同时填写支付宝应用私钥和支付宝公钥。';
  return undefined;
}

function credentialDraft(channel: PaymentChannel, provided: boolean, environment: '' | 'PRODUCTION' | 'SANDBOX', wechat: WechatCredentials, alipay: AlipayCredentials): Partial<SaveDraft> {
  if (channel === 'ALIPAY') return provided && environment ? { environment, alipayCredentials: alipay } : {};
  return provided ? { wechatCredentials: wechat } : {};
}

function WechatFields({ value, setValue, busy }: { value: WechatCredentials; setValue(value: WechatCredentials): void; busy: boolean }) {
  const update = (key: keyof WechatCredentials, next: string) => setValue({ ...value, [key]: next });
  return <div className="flex flex-col gap-4">
    <div className="grid gap-4 sm:grid-cols-2">
      <div className="flex flex-col gap-2"><Label htmlFor="wechat-serial">微信商户 API 证书序列号</Label><Input id="wechat-serial" value={value.merchantApiSerial} onChange={(event) => update('merchantApiSerial', event.target.value)} placeholder="证书十六进制序列号" disabled={busy} /></div>
      <div className="flex flex-col gap-2"><Label htmlFor="wechat-api-key">微信 API V3 Key</Label><Input id="wechat-api-key" type="password" autoComplete="off" value={value.apiV3Key} onChange={(event) => update('apiV3Key', event.target.value)} placeholder="32 字节 API V3 Key" disabled={busy} /></div>
    </div>
    <PemField id="wechat-cert" label="微信商户 API 证书" value={value.merchantApiCert} onChange={(next) => update('merchantApiCert', next)} placeholder="粘贴 -----BEGIN CERTIFICATE----- 开头的商户 API 证书" busy={busy} />
    <PemField id="wechat-private" label="微信商户私钥" value={value.merchantPrivateKey} onChange={(next) => update('merchantPrivateKey', next)} placeholder="粘贴 -----BEGIN PRIVATE KEY----- 或 -----BEGIN RSA PRIVATE KEY----- 开头的商户私钥" busy={busy} />
    <div className="flex flex-col gap-2">
      <Label htmlFor="wechat-mode">微信验签模式</Label>
      <Select value={value.verificationMode} onValueChange={(next) => setValue({ ...value, verificationMode: next as WechatCredentials['verificationMode'], wechatPublicKeyId: '', wechatPublicKey: '' })} disabled={busy}>
        <SelectTrigger id="wechat-mode" className="h-10"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem value="PLATFORM_CERT">平台证书自动获取</SelectItem>
          <SelectItem value="PUBLIC_KEY">微信支付公钥</SelectItem>
        </SelectContent>
      </Select>
    </div>
    {value.verificationMode === 'PUBLIC_KEY' && <><div className="flex flex-col gap-2"><Label htmlFor="wechat-public-id">微信支付公钥 ID</Label><Input id="wechat-public-id" value={value.wechatPublicKeyId} onChange={(event) => update('wechatPublicKeyId', event.target.value)} placeholder="PUB_KEY_ID_ 开头的完整 ID" disabled={busy} /></div><PemField id="wechat-public" label="微信支付公钥" value={value.wechatPublicKey ?? ''} onChange={(next) => update('wechatPublicKey', next)} placeholder="粘贴 -----BEGIN PUBLIC KEY----- 开头的微信支付公钥" busy={busy} /></>}
  </div>;
}

/** 更换支付宝凭证时要求明确选择环境，占位提示不作为可提交选项。 */
function AlipayFields({ value, setValue, environment, setEnvironment, busy }: { value: AlipayCredentials; setValue(value: AlipayCredentials): void; environment: '' | 'PRODUCTION' | 'SANDBOX'; setEnvironment(value: '' | 'PRODUCTION' | 'SANDBOX'): void; busy: boolean }) {
  const update = (key: keyof AlipayCredentials, next: string) => setValue({ ...value, [key]: next });
  return <div className="flex flex-col gap-4">
    <div className="flex flex-col gap-2">
      <Label htmlFor="alipay-environment">支付宝环境</Label>
      <AdminFormSelect id="alipay-environment" value={environment} placeholder="更换凭证时请选择环境" disabled={busy}
        options={[{ value: 'PRODUCTION', label: '正式环境' }, { value: 'SANDBOX', label: '沙箱环境' }]}
        onValueChange={(next) => setEnvironment(next as 'PRODUCTION' | 'SANDBOX')} />
    </div>
    <PemField id="alipay-private" label="支付宝应用私钥" value={value.appPrivateKey} onChange={(next) => update('appPrivateKey', next)} placeholder="粘贴支付宝应用私钥原文或 PEM 内容" busy={busy} />
    <PemField id="alipay-public-key" label="支付宝公钥" value={value.alipayPublicKey} onChange={(next) => update('alipayPublicKey', next)} placeholder="粘贴支付宝公钥原文或 PEM 内容" busy={busy} />
  </div>;
}

function PemField({ id, label, value, onChange, placeholder, busy }: { id: string; label: string; value: string; onChange(value: string): void; placeholder: string; busy: boolean }) {
  return <div className="flex flex-col gap-2"><Label htmlFor={id}>{label}</Label><Textarea id={id} className="min-h-28 font-mono text-xs" value={value} onChange={(event) => onChange(event.target.value)} placeholder={placeholder} autoComplete="off" disabled={busy} /></div>;
}

function stateLabelFor(state: ChannelView['own']['state']): string {
  return { UNCONFIGURED: '未配置', VALID: '本地校验通过', DISABLED: '停用', ERROR: '配置错误' }[state];
}

function emptyWechat(): WechatCredentials { return { merchantApiSerial: '', merchantApiCert: '', apiV3Key: '', merchantPrivateKey: '', verificationMode: 'PLATFORM_CERT', wechatPublicKeyId: '', wechatPublicKey: '' }; }
function emptyAlipay(): AlipayCredentials { return { appPrivateKey: '', alipayPublicKey: '' }; }
function hasWechatInput(value: WechatCredentials): boolean { return Boolean(value.merchantApiSerial || value.merchantApiCert || value.apiV3Key || value.merchantPrivateKey || value.wechatPublicKeyId || value.wechatPublicKey); }
function hasAlipayInput(value: AlipayCredentials): boolean { return Boolean(value.appPrivateKey || value.alipayPublicKey); }
