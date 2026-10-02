import { useState } from 'react';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { yuanToFen } from '@/features/hq/lib/customerMoney';
import { parseCouponLocalDateTime } from './customerCouponTemplateTime';

export interface CustomerCouponTemplateInput {
  code: string;
  title: string;
  requestKey: string;
  amountFen: number;
  minSpendFen: number;
  effectiveAt: string | null;
  distributionEndsAt: string | null;
  daysAfterActivation: number;
  perMemberLimit: number;
  totalIssueLimit: number | null;
  enabled: boolean;
}

type CouponIssuer = 'headquarters' | 'store';
type TemplateDraft = Record<'title' | 'amount' | 'minimum' | 'effectiveAt' | 'distributionEndsAt' | 'days' | 'perMember' | 'total', string>;
type TemplateErrors = Partial<Record<keyof TemplateDraft, string>>;

const initialDraft: TemplateDraft = {
  title: '', amount: '', minimum: '', effectiveAt: '', distributionEndsAt: '', days: '30', perMember: '1', total: '',
};

export function CustomerCouponTemplateDialog({ issuer = 'headquarters', open, busy, onOpenChange, onSave }: {
  issuer?: CouponIssuer;
  open: boolean;
  busy: boolean;
  onOpenChange(open: boolean): void;
  onSave(input: CustomerCouponTemplateInput): Promise<void>;
}) {
  const form = useCouponTemplateDraft(issuer, onSave);
  return <AdminFormDialogShell
    open={open} onOpenChange={onOpenChange}
    title={issuer === 'store' ? '新增门店券模板' : '新增固定金额券'}
    description="已发放模板的金额、门槛和数量上限不可修改；调整时请新建模板。"
    submitLabel={issuer === 'store' ? '保存模板' : '创建模板'} cancelLabel="取消" isSubmitting={busy}
    onSubmit={form.save}
  >
    <CouponTemplateFields {...form} />
  </AdminFormDialogShell>;
}

function useCouponTemplateDraft(issuer: CouponIssuer, onSave: (input: CustomerCouponTemplateInput) => Promise<void>) {
  const [requestKey] = useState(() => crypto.randomUUID());
  const code = issuer === 'store'
    ? `SC-${requestKey.slice(0, 12).toUpperCase()}`
    : `C${requestKey.replaceAll('-', '').slice(0, 10).toUpperCase()}`;
  const [draft, setDraft] = useState<TemplateDraft>(initialDraft);
  const [errors, setErrors] = useState<TemplateErrors>({});
  const [error, setError] = useState<string>();
  function change(key: keyof TemplateDraft, value: string) {
    setDraft((current) => ({ ...current, [key]: value }));
    setErrors((current) => ({ ...current, [key]: undefined }));
  }
  async function save() {
    setError(undefined);
    const { input, errors: nextErrors } = validateTemplateDraft(draft, code, requestKey);
    setErrors(nextErrors);
    if (!input) return;
    try {
      await onSave(input);
    } catch {
      setError('券模板创建未完成，请核对发放额度后重试。');
    }
  }
  return { draft, errors, error, change, save };
}

function CouponTemplateFields({ draft, errors, error, change }: ReturnType<typeof useCouponTemplateDraft>) {
  return <div className="flex flex-col gap-5">
      <CouponField id="coupon-title" label="名称" value={draft.title} onChange={(value) => change('title', value)} placeholder="例如：门店开业优惠券" error={errors.title} kind="text" />
      <div data-testid="coupon-amount-threshold-row" className="grid grid-cols-1 gap-5 sm:grid-cols-2">
        <CouponField id="coupon-amount" label="面额（元）" value={draft.amount} onChange={(value) => change('amount', value)} placeholder="例如：10.00" error={errors.amount} kind="money" />
        <CouponField id="coupon-minimum" label="消费门槛（元）" value={draft.minimum} onChange={(value) => change('minimum', value)} placeholder="留空或 0 表示不限制" error={errors.minimum} kind="money" />
      </div>
      <div data-testid="coupon-effective-days-row" className="grid grid-cols-1 gap-5 sm:grid-cols-2">
        <CouponField id="coupon-effective-at" label="生效日期时间" value={draft.effectiveAt} onChange={(value) => change('effectiveAt', value)} placeholder="留空表示立即生效" error={errors.effectiveAt} kind="datetime" />
        <CouponField id="coupon-days" label="激活后有效天数" value={draft.days} onChange={(value) => change('days', value)} placeholder="1～365" error={errors.days} kind="count" />
      </div>
      <div data-testid="coupon-issue-limits-row" className="grid grid-cols-1 gap-5 sm:grid-cols-2">
        <CouponField id="coupon-member-limit" label="每会员最多发放" value={draft.perMember} onChange={(value) => change('perMember', value)} placeholder="例如：1" error={errors.perMember} kind="count" />
        <CouponField id="coupon-total-limit" label="总发放量" value={draft.total} onChange={(value) => change('total', value)} placeholder="留空表示不限制" error={errors.total} kind="count" />
      </div>
      <div data-testid="coupon-distribution-cutoff-row">
        <CouponField id="coupon-distribution-ends-at" label="发放截止日期时间" value={draft.distributionEndsAt} onChange={(value) => change('distributionEndsAt', value)} placeholder="留空表示持续向符合条件的会员补发" error={errors.distributionEndsAt} kind="datetime" />
      </div>
      {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
    </div>;
}

function validateTemplateDraft(draft: TemplateDraft, code: string, requestKey: string) {
  const errors: TemplateErrors = {};
  const title = validateTitle(draft.title, errors);
  const { amountFen, minSpendFen } = validateMoney(draft, errors);
  const { effectiveAt, distributionEndsAt } = validateSchedule(draft, errors);
  const { days, perMember, total } = validateLimits(draft, errors);
  if (Object.keys(errors).length > 0) return { errors };
  return { errors, input: {
    code, title, requestKey, amountFen, minSpendFen,
    effectiveAt, distributionEndsAt,
    daysAfterActivation: days, perMemberLimit: perMember,
    totalIssueLimit: total, enabled: true,
  } satisfies CustomerCouponTemplateInput };
}

function validateTitle(value: string, errors: TemplateErrors): string {
  const title = value.trim();
  if (!title) errors.title = '请输入券名称';
  else if (new TextEncoder().encode(title).length > 128) errors.title = '券名称不能超过 128 字节';
  return title;
}

function validateMoney(draft: TemplateDraft, errors: TemplateErrors) {
  let amountFen = 0;
  let minSpendFen = 0;
  try {
    amountFen = yuanToFen(draft.amount);
    if (amountFen < 1) errors.amount = '面额须大于 0 元';
  } catch { errors.amount = '请输入有效面额（最多两位小数）'; }
  try { minSpendFen = draft.minimum.trim() ? yuanToFen(draft.minimum) : 0; }
  catch { errors.minimum = '请输入有效最低消费金额（最多两位小数）'; }
  return { amountFen, minSpendFen };
}

function validateSchedule(draft: TemplateDraft, errors: TemplateErrors) {
  let effectiveAt: string | null = null;
  let distributionEndsAt: string | null = null;
  try { effectiveAt = parseCouponLocalDateTime(draft.effectiveAt); }
  catch { errors.effectiveAt = '请输入有效的日期和时间'; }
  try { distributionEndsAt = parseCouponLocalDateTime(draft.distributionEndsAt); }
  catch { errors.distributionEndsAt = '请输入有效的日期和时间'; }
  if (!errors.effectiveAt && !errors.distributionEndsAt && distributionEndsAt) {
    if (effectiveAt && Date.parse(distributionEndsAt) <= Date.parse(effectiveAt)) {
      errors.distributionEndsAt = '发放截止时间须晚于生效时间';
    } else if (!effectiveAt && Date.parse(distributionEndsAt) <= Date.now()) {
      errors.distributionEndsAt = '发放截止时间须晚于当前时间';
    }
  }
  return { effectiveAt, distributionEndsAt };
}

function validateLimits(draft: TemplateDraft, errors: TemplateErrors) {
  let days = 0;
  let perMember = 0;
  let total: number | null = null;
  try {
    days = positiveInt(draft.days);
    if (days > 365) throw new Error('out of range');
  } catch { errors.days = '有效天数须为 1～365 天'; }
  try { perMember = positiveInt(draft.perMember); }
  catch { errors.perMember = '每会员发放上限须为正整数'; }
  if (draft.total.trim()) {
    try { total = positiveInt(draft.total); }
    catch { errors.total = '总发放量须为正整数；不限量请留空'; }
  }
  return { days, perMember, total };
}

function positiveInt(value: string): number {
  if (!/^[1-9]\d*$/.test(value)) throw new Error('请输入正整数');
  const parsed = Number(value);
  if (!Number.isSafeInteger(parsed) || parsed > 2_147_483_647) throw new Error('数值超过允许范围');
  return parsed;
}

function CouponField({ id, label, value, onChange, placeholder, kind, error, help }: {
  id: string;
  label: string;
  value: string;
  onChange(value: string): void;
  placeholder: string;
  kind: 'text' | 'money' | 'count' | 'datetime';
  error?: string;
  help?: string;
}) {
  const type = kind === 'datetime' ? 'datetime-local' : 'text';
  const inputMode = kind === 'money' ? 'decimal' : kind === 'count' ? 'numeric' : 'text';
  return <div className="flex flex-col gap-2">
    <Label htmlFor={id}>{label}</Label>
    <Input id={id} type={type} inputMode={inputMode} value={value} onChange={(event) => onChange(event.target.value)}
      placeholder={placeholder} aria-invalid={Boolean(error)} aria-describedby={error ? `${id}-error` : help ? `${id}-help` : undefined} />
    {error && <p id={`${id}-error`} className="text-sm text-destructive">{error}</p>}
    {!error && help && <p id={`${id}-help`} className="text-xs text-muted-foreground">{help}</p>}
  </div>;
}
