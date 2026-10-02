// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { TabItem } from './TabItem';

afterEach(cleanup);

it('活动标签使用卡片底色和品牌色文字强化当前位置', () => {
  render(<TabItem
    tab={{ path: '/admin/hq/products/manage', search: '', title: '维护商品', affix: false, version: 0 }}
    isActive
    errored={false}
    onActivate={vi.fn()}
    onClose={vi.fn()}
  />);

  const tab = screen.getByRole('button', { name: '维护商品' }).parentElement;
  expect(tab?.classList.contains('bg-card')).toBe(true);
  expect(tab?.classList.contains('text-primary')).toBe(true);
  expect(tab?.classList.contains('font-medium')).toBe(true);
});
