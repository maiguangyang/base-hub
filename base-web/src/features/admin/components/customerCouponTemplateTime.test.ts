import { describe, expect, it } from 'vitest';
import { couponTemplateAvailability, parseCouponLocalDateTime } from './customerCouponTemplateTime';

describe('parseCouponLocalDateTime', () => {
  it('keeps an empty optional date blank', () => {
    expect(parseCouponLocalDateTime('')).toBeNull();
    expect(parseCouponLocalDateTime('   ')).toBeNull();
  });

  it('converts a valid local date and time to ISO', () => {
    const value = '2026-06-15T09:30';
    expect(parseCouponLocalDateTime(value)).toBe(new Date(value).toISOString());
  });

  it.each(['not-a-date', '2026-02-30T10:00', '2026-06-15T24:00'])('rejects malformed or impossible local time %s', (value) => {
    expect(() => parseCouponLocalDateTime(value)).toThrow('请输入有效的日期和时间');
  });

  it.runIf(Intl.DateTimeFormat().resolvedOptions().timeZone === 'America/New_York')('rejects a daylight-saving gap', () => {
    expect(() => parseCouponLocalDateTime('2026-03-08T02:30')).toThrow('请输入有效的日期和时间');
  });
});

describe('couponTemplateAvailability', () => {
  const now = new Date('2026-06-15T10:00:00.000Z').getTime();

  it('uses the required status precedence and exclusive cutoff', () => {
    const counts = { totalIssueLimit: 0, issuedCount: 0 };
    expect(couponTemplateAvailability({ enabled: false, effectiveAt: '2027-01-01T00:00:00.000Z', distributionEndsAt: null, ...counts }, now)).toBe('已停用');
    expect(couponTemplateAvailability({ enabled: true, effectiveAt: '2026-06-15T10:00:00.001Z', distributionEndsAt: null, ...counts }, now)).toBe('未生效');
    expect(couponTemplateAvailability({ enabled: true, effectiveAt: '2026-01-01T00:00:00.000Z', distributionEndsAt: '2026-06-15T10:00:00.000Z', ...counts }, now)).toBe('已结束');
    expect(couponTemplateAvailability({ enabled: true, effectiveAt: '2026-01-01T00:00:00.000Z', distributionEndsAt: null, ...counts }, now)).toBe('可发放');
  });

  it('marks a finite exhausted template as unavailable without affecting unlimited templates', () => {
    expect(couponTemplateAvailability({
      enabled: true,
      effectiveAt: '2026-01-01T00:00:00.000Z',
      distributionEndsAt: null,
      totalIssueLimit: 1,
      issuedCount: 1,
    }, now)).toBe('已发完');
    expect(couponTemplateAvailability({
      enabled: true,
      effectiveAt: '2026-01-01T00:00:00.000Z',
      distributionEndsAt: null,
      totalIssueLimit: 0,
      issuedCount: 10,
    }, now)).toBe('可发放');
  });
});
