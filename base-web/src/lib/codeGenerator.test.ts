import { describe, expect, it } from 'vitest';
import { generateEntityCode } from './codeGenerator';

describe('generateEntityCode', () => {
  it('生成包含指定前缀的有效内部编码', () => {
    const code = generateEntityCode('STR');
    expect(code.startsWith('STR')).toBe(true);
    expect(code.length).toBeGreaterThanOrEqual(10);
    expect(code.length).toBeLessThanOrEqual(32);
    expect(/^[A-Z0-9]+$/.test(code)).toBe(true);
  });

  it('连续生成具有唯一性', () => {
    const codes = new Set(Array.from({ length: 50 }, () => generateEntityCode('TEST')));
    expect(codes.size).toBe(50);
  });

  it('无前缀时也能生成合规编码', () => {
    const code = generateEntityCode();
    expect(code.length).toBeGreaterThanOrEqual(8);
    expect(code.length).toBeLessThanOrEqual(32);
  });
});
