import { readFileSync } from 'node:fs';
import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it } from 'vitest';
import { AdminStatusSwitch } from '@/features/admin/components/AdminStatusSwitch';
import { Checkbox } from './checkbox';
import { Select, SelectTrigger } from './select';
import { Switch } from './switch';
import { Textarea } from './textarea';

const themedFocusClasses = [
  'focus-visible:border-ring',
  'focus-visible:ring-[3px]',
  'focus-visible:ring-ring/50',
];

function expectThemedFocus(markup: string) {
  for (const className of themedFocusClasses) {
    expect(markup).toContain(className);
  }
}

describe('主题交互状态', () => {
  it.each([
    ['Textarea', renderToStaticMarkup(<Textarea aria-label="备注" />)],
    [
      'SelectTrigger',
      renderToStaticMarkup(
        <Select>
          <SelectTrigger aria-label="状态" />
        </Select>,
      ),
    ],
    ['Checkbox', renderToStaticMarkup(<Checkbox aria-label="全选" />)],
    ['Switch', renderToStaticMarkup(<Switch aria-label="启用" />)],
  ])('%s 键盘聚焦时使用主题边框和外圈', (_name, markup) => {
    expectThemedFocus(markup);
  });

  it('未勾选复选框固定白底，保留选中与半选状态', () => {
    for (const checked of [false, true, 'indeterminate'] as const) {
      const markup = renderToStaticMarkup(<Checkbox checked={checked} className="bg-muted" aria-label="测试复选框" />);
      expect(markup).toContain(`data-state="${checked === false ? 'unchecked' : checked === true ? 'checked' : 'indeterminate'}"`);
      expect(markup).toContain('data-[state=unchecked]:!bg-white');
      expect(markup).toContain('data-[state=checked]:bg-primary');
    }
  });

  it('后台状态开关使用主题聚焦外圈', () => {
    const markup = renderToStaticMarkup(
      <AdminStatusSwitch checked label="启用" onCheckedChange={() => undefined} />,
    );

    expect(markup).toContain('data-slot="switch"');
    expect(markup).toContain('focus-visible:ring-[3px]');
    expect(markup).toContain('data-[state=checked]:bg-success-fg');
    expect(markup).not.toContain('bg-success-bg');
    expect(markup).not.toMatch(/\sdata-\[state=checked\]:bg-primary(?:\s|")/);
  });
});

describe('后台导航主题消费契约', () => {
  it('侧栏选中项消费 sidebar-primary token', () => {
    const source = readFileSync(
      new URL('../../features/admin/components/AdminSidebar/NavItem.tsx', import.meta.url),
      'utf8',
    );

    expect(source).toContain('data-[active=true]:bg-sidebar-primary');
    expect(source).toContain('data-[active=true]:text-sidebar-primary-foreground');
  });

  it('标签栏活动下划线消费 primary token', () => {
    const source = readFileSync(
      new URL('../../features/admin/components/AdminTabBar/TabItem.tsx', import.meta.url),
      'utf8',
    );

    expect(source).toContain('h-0.5 bg-primary');
  });
});
