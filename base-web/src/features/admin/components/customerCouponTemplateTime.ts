export type CouponTemplateAvailability = '已停用' | '未生效' | '已结束' | '已发完' | '可发放';

export function parseCouponLocalDateTime(value: string): string | null {
  const trimmed = value.trim();
  if (!trimmed) return null;

  const match = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2})$/.exec(trimmed);
  if (!match) throw new Error('请输入有效的日期和时间');

  const [, year, month, day, hour, minute] = match.map(Number);
  const date = new Date(trimmed);
  if (Number.isNaN(date.getTime())
    || date.getFullYear() !== year
    || date.getMonth() + 1 !== month
    || date.getDate() !== day
    || date.getHours() !== hour
    || date.getMinutes() !== minute) {
    throw new Error('请输入有效的日期和时间');
  }
  return date.toISOString();
}

export function couponTemplateAvailability(template: {
  enabled: boolean;
  effectiveAt: unknown;
  distributionEndsAt?: unknown;
  totalIssueLimit: number;
  issuedCount: number;
}, now = Date.now()): CouponTemplateAvailability {
  if (!template.enabled) return '已停用';
  if (couponTimeMillis(template.effectiveAt) > now) return '未生效';
  if (template.distributionEndsAt && couponTimeMillis(template.distributionEndsAt) <= now) return '已结束';
  if (template.totalIssueLimit > 0 && template.issuedCount >= template.totalIssueLimit) return '已发完';
  return '可发放';
}

export function formatCouponDateTime(value: unknown): string {
  return new Date(couponTimeMillis(value)).toLocaleString('zh-CN', { hour12: false });
}

function couponTimeMillis(value: unknown): number {
  if (typeof value !== 'string' && typeof value !== 'number' && !(value instanceof Date)) return Number.NaN;
  return new Date(value).getTime();
}
