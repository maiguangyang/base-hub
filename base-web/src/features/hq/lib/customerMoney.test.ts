import { describe, expect, it } from 'vitest';
import { fenToYuan, yuanToFen } from './customerMoney';

describe('customer money', () => {
  it('converts decimal yuan exactly', () => {
    expect(yuanToFen('1.23')).toBe(123);
    expect(yuanToFen('0.01')).toBe(1);
    expect(yuanToFen('100')).toBe(10000);
    expect(fenToYuan(123)).toBe('1.23');
  });

  it('rejects precision and GraphQL Int overflow', () => {
    expect(() => yuanToFen('1.234')).toThrow();
    expect(() => yuanToFen('-1')).toThrow();
    expect(() => yuanToFen('21474836.48')).toThrow();
    expect(() => yuanToFen('NaN')).toThrow();
  });
});
