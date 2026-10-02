import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it } from 'vitest';
import { AdminPrimaryActionButton } from './AdminPrimaryActionButton';
import { AdminTableShell } from './AdminTableShell';
import { AdminToolbar } from './AdminToolbar';

describe('后台列表共享组件视觉基线', () => {
  it('筛选栏保留统一的间距、边框和背景', () => {
    const element = AdminToolbar({ children: null });

    expect(element.props.className).toContain('m-4');
    expect(element.props.className).toContain('rounded-lg');
    expect(element.props.className).toContain('border-border');
    expect(element.props.className).toContain('bg-muted/30');
    expect(element.props.className).toContain('p-4');
  });

  it('表格壳保留统一的表头、单元格与悬停样式', () => {
    const element = AdminTableShell({ toolbar: null, children: null });

    expect(element.props.className).toContain('[&_[data-slot=table-header]]:bg-muted/30');
    expect(element.props.className).toContain('[&_[data-slot=table-head]]:h-12');
    expect(element.props.className).toContain('[&_[data-slot=table-head]]:px-4');
    expect(element.props.className).toContain('[&_[data-slot=table-cell]]:px-4');
    expect(element.props.className).toContain('[&_[data-slot=table-cell]]:py-3');
    expect(element.props.className).toContain('[&_[data-slot=table-row]]:hover:bg-muted/20');
  });

  it('主新增按钮忽略页面传入的视觉覆盖并保留统一尺寸', () => {
    const element = AdminPrimaryActionButton({ children: '新增', className: 'h-20' } as never);

    expect(element.props.className).toBe('h-10 px-4');
  });

  it('主新增按钮只显示文字，不渲染加号图标', () => {
    const markup = renderToStaticMarkup(<AdminPrimaryActionButton>新增直营店</AdminPrimaryActionButton>);

    expect(markup).toContain('新增直营店');
    expect(markup).not.toContain('<svg');
  });
});
