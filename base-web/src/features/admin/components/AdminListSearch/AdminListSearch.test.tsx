// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { AdminListSearch } from './AdminListSearch';

afterEach(cleanup);

describe('AdminListSearch', () => {
  it('仅在点击搜索时提交去除首尾空格后的草稿', () => {
    const onSearch = vi.fn();
    render(<AdminListSearch value="" placeholder="搜索角色" onSearch={onSearch} />);

    fireEvent.change(screen.getByRole('textbox', { name: '搜索角色' }), {
      target: { value: '  运营  ' },
    });
    expect(onSearch).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole('button', { name: '搜索' }));
    expect(onSearch).toHaveBeenCalledWith('运营');
  });

  it('支持 Enter 提交，并在外部值变化时重新同步草稿', () => {
    const onSearch = vi.fn();
    const view = render(<AdminListSearch value="旧值" placeholder="搜索角色" onSearch={onSearch} />);
    const input = screen.getByRole('textbox', { name: '搜索角色' });

    fireEvent.change(input, { target: { value: '新值' } });
    fireEvent.submit(input.closest('form')!);
    expect(onSearch).toHaveBeenLastCalledWith('新值');

    view.rerender(<AdminListSearch value="已重置" placeholder="搜索角色" onSearch={onSearch} />);
    expect((screen.getByRole('textbox', { name: '搜索角色' }) as HTMLInputElement).value).toBe('已重置');
  });

  it('嵌入现有表单时不生成嵌套 form，并保留 Enter 搜索', () => {
    const onSearch = vi.fn();
    const view = render(<form><AdminListSearch embedded value="" placeholder="搜索销售项" onSearch={onSearch} /></form>);
    const input = screen.getByRole('textbox', { name: '搜索销售项' });

    expect(view.container.querySelectorAll('form')).toHaveLength(1);
    fireEvent.change(input, { target: { value: '  整箱  ' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    expect(onSearch).toHaveBeenCalledWith('整箱');
  });
});
