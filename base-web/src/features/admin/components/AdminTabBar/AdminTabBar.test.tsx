// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { AdminTabBar } from './AdminTabBar';

vi.mock('@/features/admin/hooks/useAdminTabs', () => ({
  useAdminTabList: () => [],
  useActiveTabPath: () => '/admin/hq',
  useErroredTabPaths: () => [],
}));

vi.mock('./hooks/useTabActions', () => ({
  useTabActions: () => ({
    activate: vi.fn(), close: vi.fn(), closeOthers: vi.fn(), closeRight: vi.fn(), closeAll: vi.fn(), refresh: vi.fn(),
  }),
}));

vi.mock('./hooks/useTabOverflow', () => ({
  useTabOverflow: () => ({
    scrollerRef: { current: null },
    canScrollLeft: false,
    canScrollRight: false,
    scrollByStep: vi.fn(),
  }),
}));

afterEach(cleanup);

it('保留横向滚动能力但不显示浏览器原生滚动条', () => {
  render(<AdminTabBar />);

  const scroller = screen.getByRole('button', { name: '向左滚动标签' }).nextElementSibling;
  expect(scroller).not.toBeNull();
  expect(scroller?.classList.contains('overflow-x-auto')).toBe(true);
  expect(scroller?.classList.contains('[scrollbar-width:none]')).toBe(true);
  expect(scroller?.classList.contains('[&::-webkit-scrollbar]:hidden')).toBe(true);
});

it('快捷导航使用区别于页面的次级底色', () => {
  render(<AdminTabBar />);

  const tabBar = screen.getByRole('button', { name: '向左滚动标签' }).parentElement;
  expect(tabBar?.classList.contains('bg-secondary')).toBe(true);
  expect(tabBar?.classList.contains('bg-background')).toBe(false);
});
