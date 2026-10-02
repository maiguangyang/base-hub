// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, expect, it } from 'vitest';
import { AdminPageHeader } from './AdminPageHeader';

afterEach(cleanup);

it('页头使用无阴影的独立卡片表面并保留与正文的间距', () => {
  render(<AdminPageHeader title="编辑商品" description="维护商品信息、主图、规格及 SKU 包装。" actions={<button type="button">保存</button>} />);

  const titleGroup = screen.getByRole('heading', { name: '编辑商品' }).parentElement;
  const surface = titleGroup?.parentElement;
  const spacing = surface?.parentElement;
  expect(surface?.classList.contains('bg-card')).toBe(true);
  expect(surface?.classList.contains('border')).toBe(true);
  expect(surface?.classList.contains('shadow-sm')).toBe(false);
  expect(spacing?.classList.contains('pb-4')).toBe(true);
});
