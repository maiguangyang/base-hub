import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';

/** 从 CSS 文本中取出指定选择器块内声明的全部自定义属性名。 */
function declaredVars(css: string, selector: string): string[] {
  const escaped = selector.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  const block = new RegExp(`${escaped}\\s*\\{([^}]*)\\}`).exec(css);
  if (!block) return [];
  return [...block[1].matchAll(/(--[\w-]+)\s*:/g)].map((m) => m[1]).sort();
}

function declaredValue(css: string, selector: string, property: string): string | undefined {
  const escapedSelector = selector.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  const block = new RegExp(`${escapedSelector}\\s*\\{([^}]*)\\}`).exec(css);
  if (!block) return undefined;
  const escapedProperty = property.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  return new RegExp(`${escapedProperty}\\s*:\\s*([^;]+);`).exec(block[1])?.[1].trim();
}

function relativeLuminance(hex: string): number {
  const channels = /^#([\da-f]{2})([\da-f]{2})([\da-f]{2})$/i.exec(hex);
  if (!channels) throw new Error(`Expected a six-digit hex color, received ${hex}`);
  const [red, green, blue] = channels.slice(1).map((channel) => {
    const value = Number.parseInt(channel, 16) / 255;
    return value <= 0.04045 ? value / 12.92 : ((value + 0.055) / 1.055) ** 2.4;
  });
  return 0.2126 * red + 0.7152 * green + 0.0722 * blue;
}

function contrastRatio(first: string, second: string): number {
  const [lighter, darker] = [relativeLuminance(first), relativeLuminance(second)].sort((a, b) => b - a);
  return (lighter + 0.05) / (darker + 0.05);
}

const adminTheme = readFileSync(new URL('./admin-theme.css', import.meta.url), 'utf8');
const globalCss = readFileSync(new URL('./global.css', import.meta.url), 'utf8');

describe('后台 Ant Design 蓝色主题', () => {
  it('浅色主题以 #1677ff 为品牌基准并使用可读的填充色', () => {
    const selector = ':root[data-surface="admin"]';
    expect(declaredValue(adminTheme, selector, '--primary')).toBe('#0958d9');
    expect(declaredValue(adminTheme, selector, '--ring')).toBe('#1677ff');
    expect(declaredValue(adminTheme, selector, '--chart-1')).toBe('#1677ff');
  });

  it('深色主题使用同色系高亮变体', () => {
    const selector = ':root.dark[data-surface="admin"]';
    for (const property of ['--primary', '--ring', '--chart-1']) {
      expect(declaredValue(adminTheme, selector, property)).toBe('#4096ff');
    }
  });

  it('侧栏恢复原来的低饱和背景并保持选中项可读', () => {
    const light = ':root[data-surface="admin"]';
    const dark = ':root.dark[data-surface="admin"]';
    expect(declaredValue(adminTheme, light, '--sidebar')).toBe('oklch(0.21 0.035 265)');
    expect(declaredValue(adminTheme, dark, '--sidebar')).toBe('oklch(0.165 0.022 265)');
    for (const selector of [light, dark]) {
      const background = declaredValue(adminTheme, selector, '--sidebar-primary');
      const foreground = declaredValue(adminTheme, selector, '--sidebar-primary-foreground');
      expect(background).toBe('#0958d9');
      expect(foreground).toBe('#ffffff');
      expect(contrastRatio(foreground!, background!)).toBeGreaterThanOrEqual(4.5);
    }
  });

  it('页面和交互表面恢复原来的中性色板', () => {
    const light = ':root[data-surface="admin"]';
    const dark = ':root.dark[data-surface="admin"]';
    expect(declaredValue(adminTheme, light, '--background')).toBe('oklch(0.965 0.004 265)');
    expect(declaredValue(adminTheme, light, '--accent')).toBe('oklch(0.94 0.012 265)');
    expect(declaredValue(adminTheme, dark, '--background')).toBe('oklch(0.205 0.012 265)');
    expect(declaredValue(adminTheme, dark, '--accent')).toBe('oklch(0.32 0.03 265)');
  });
});

describe('后台作用域主题', () => {
  it('深色块穷尽声明浅色块的每一个变量', () => {
    const light = declaredVars(adminTheme, '[data-surface="admin"]');
    const dark = declaredVars(adminTheme, '.dark[data-surface="admin"]');
    expect(light.length).toBeGreaterThan(0);
    expect(dark).toEqual(light);
  });

  it('五个语义状态各有 bg / border / fg 三层', () => {
    const light = declaredVars(adminTheme, '[data-surface="admin"]');
    for (const tone of ['success', 'warning', 'danger', 'info', 'neutral']) {
      for (const layer of ['bg', 'border', 'fg']) {
        expect(light).toContain(`--${tone}-${layer}`);
      }
    }
  });

  it('图表分类色定义六个', () => {
    const light = declaredVars(adminTheme, '[data-surface="admin"]');
    for (let i = 1; i <= 6; i++) expect(light).toContain(`--chart-${i}`);
  });

  it('新增 token 全部注册进 global.css 的 @theme inline', () => {
    const theme = /@theme inline\s*\{([^}]*)\}/.exec(globalCss);
    expect(theme).not.toBeNull();
    const mapped = theme === null ? '' : theme[1];
    const adminOnly = declaredVars(adminTheme, '[data-surface="admin"]')
      .filter((name) => !/^--(background|foreground|card|popover|primary|secondary|muted|accent|destructive|border|input|ring)/.test(name));
    expect(adminOnly.length).toBeGreaterThan(0);
    for (const name of adminOnly) {
      expect(mapped).toContain(`--color-${name.slice(2)}: var(${name});`);
    }
  });

  it('作用域选择器必须带 :root 前缀以压过 global.css 的 :root', () => {
    // @import 只能在文件开头，本文件因此加载在 :root 之前。
    // 裸 [data-surface="admin"] 与 :root 特异性相同，后来的 :root 会赢，
    // 导致后台覆盖在浅色模式下静默失效。
    expect(adminTheme).toContain(':root[data-surface="admin"] {');
    expect(adminTheme).toContain(':root.dark[data-surface="admin"] {');
    expect(adminTheme).not.toMatch(/\n\[data-surface="admin"\]\s*\{/);
  });

  it('官网 :root 与 .dark 不含后台专属 token', () => {
    for (const selector of [':root', '.dark']) {
      const vars = declaredVars(globalCss, selector);
      expect(vars).not.toContain('--sidebar');
      expect(vars).not.toContain('--success-fg');
      expect(vars).not.toContain('--chart-1');
    }
  });
});
