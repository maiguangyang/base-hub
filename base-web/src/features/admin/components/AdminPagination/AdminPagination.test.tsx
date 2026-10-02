// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { AdminPagination } from './AdminPagination';

afterEach(cleanup);

describe('AdminPagination', () => {
  it('自身持有柔和分页表面并保持控件白底', () => {
    const { container } = render(
      <AdminPagination
        total={20}
        page={1}
        pageSize={10}
        onPageChange={() => undefined}
        onPageSizeChange={() => undefined}
      />,
    );

    const markup = container.innerHTML;
    expect(markup).toContain('bg-muted/30');
    expect(markup.match(/bg-card/g)?.length).toBe(3);
  });

  it('使用公共下拉选择每页条数并回调数值', () => {
    const onPageSizeChange = vi.fn();
    const { container } = render(<AdminPagination total={100} page={1} pageSize={20}
      onPageChange={vi.fn()} onPageSizeChange={onPageSizeChange} />);

    expect(container.querySelector('select')).toBeNull();
    const trigger = screen.getByRole('combobox', { name: '每页' });
    fireEvent.click(trigger);
    fireEvent.click(screen.getByRole('option', { name: '50' }));

    expect(onPageSizeChange).toHaveBeenCalledWith(50);
  });
});
