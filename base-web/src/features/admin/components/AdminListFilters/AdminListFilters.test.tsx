// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { AdminListFilters } from './AdminListFilters';

afterEach(cleanup);

describe('AdminListFilters', () => {
  it('没有生效条件时隐藏重置操作', () => {
    render(<AdminListFilters hasActiveFilters={false} onReset={() => undefined}><span>筛选项</span></AdminListFilters>);

    expect(screen.queryByRole('button', { name: '重置筛选' })).toBeNull();
  });

  it('有生效条件时只调用一次重置操作', () => {
    const onReset = vi.fn();
    render(<AdminListFilters hasActiveFilters onReset={onReset}><span>筛选项</span></AdminListFilters>);

    fireEvent.click(screen.getByRole('button', { name: '重置筛选' }));
    expect(onReset).toHaveBeenCalledTimes(1);
  });
});
