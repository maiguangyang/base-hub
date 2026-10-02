import { describe, expect, it } from 'vitest';
import { statusToneClasses, toneForCompletion, toneForDeleteState, toneForEnabledState } from './statusTone';

describe('状态 tone 映射', () => {
  it('启用状态映射到 success，禁用映射到 danger', () => {
    expect(toneForEnabledState(1)).toBe('success');
    expect(toneForEnabledState(0)).toBe('danger');
  });

  it('已删除映射到 danger，未删除映射到 neutral', () => {
    expect(toneForDeleteState(1)).toBe('danger');
    expect(toneForDeleteState(0)).toBe('neutral');
  });

  it('已完成映射到 success，未完成映射到 warning', () => {
    expect(toneForCompletion(true)).toBe('success');
    expect(toneForCompletion(false)).toBe('warning');
  });

  it('每个 tone 都给出 bg / border / fg 三层 class，且不含裸色值', () => {
    for (const tone of ['success', 'warning', 'danger', 'info', 'neutral'] as const) {
      const classes = statusToneClasses[tone];
      expect(classes).toContain(`bg-${tone}-bg`);
      expect(classes).toContain(`border-${tone}-border`);
      expect(classes).toContain(`text-${tone}-fg`);
      expect(classes).not.toMatch(/oklch|#[0-9a-f]{3,8}\b|\brgb\(/i);
    }
  });
});
