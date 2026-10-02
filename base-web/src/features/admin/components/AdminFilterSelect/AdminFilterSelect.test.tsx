// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { AdminFilterSelect } from './AdminFilterSelect';

afterEach(cleanup);

const options = [
  { value: 'ACTIVE', label: '启用' },
  { value: 'SUSPENDED', label: '停用' },
];

describe('AdminFilterSelect', () => {
  it('使用占位文案作为可访问名称并提交具体选项', () => {
    const onValueChange = vi.fn();
    render(<AdminFilterSelect value={undefined} placeholder="全部状态" options={options} onValueChange={onValueChange} />);

    const trigger = screen.getByRole('combobox', { name: '全部状态' });
    expect(trigger.textContent).toContain('全部状态');
    fireEvent.click(trigger);
    fireEvent.click(screen.getByRole('option', { name: '启用' }));

    expect(onValueChange).toHaveBeenCalledWith('ACTIVE');
  });

  it.each([undefined, 'ACTIVE'])('value 为 %s 时不把占位文案加入选项', (value) => {
    const onValueChange = vi.fn();
    render(<AdminFilterSelect value={value} placeholder="全部状态" options={options} onValueChange={onValueChange} />);

    fireEvent.click(screen.getByRole('combobox', { name: '全部状态' }));
    expect(screen.queryByRole('option', { name: '全部状态' })).toBeNull();
    expect(screen.getAllByRole('option').map((option) => option.textContent)).toEqual(['启用', '停用']);
    expect(onValueChange).not.toHaveBeenCalled();
  });

  it('页面重置筛选后恢复占位提示', () => {
    const onValueChange = vi.fn();
    const view = render(<AdminFilterSelect value="ACTIVE" placeholder="全部状态" options={options} onValueChange={onValueChange} />);
    expect(screen.getByRole('combobox').textContent).toContain('启用');
    view.rerender(<AdminFilterSelect value={undefined} placeholder="全部状态" options={options} onValueChange={onValueChange} />);
    expect(screen.getByRole('combobox').textContent).toContain('全部状态');
    expect(onValueChange).not.toHaveBeenCalled();
  });
});
