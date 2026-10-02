// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { RoleListPage } from './RoleListPage';
import { useRolePage } from './useRolePage';

vi.mock('./useRolePage', () => ({ useRolePage: vi.fn() }));

afterEach(cleanup);

describe('总部角色列表筛选栏', () => {
  it('仅在点击搜索按钮或按 Enter 后应用输入的关键词', () => {
    const search = vi.fn();
    vi.mocked(useRolePage).mockReturnValue({
      rows: [],
      total: 0,
      page: 1,
      pageSize: 10,
      keyword: '',
      kind: undefined,
      hasFilters: false,
      open: false,
      editing: undefined,
      values: { name: '', permissionIds: [] },
      actionError: undefined,
      catalogError: undefined,
      listError: undefined,
      permissionCatalog: [],
      permissions: [],
      saving: false,
      empty: true,
      canCreate: true,
      canUpdate: true,
      canDelete: true,
      setPage: vi.fn(),
      setPageSize: vi.fn(),
      setValues: vi.fn(),
      edit: vi.fn(),
      changeOpen: vi.fn(),
      save: vi.fn(),
      remove: vi.fn(),
      search,
      changeKind: vi.fn(),
      resetFilters: vi.fn(),
    });

    render(<RoleListPage />);
    fireEvent.change(screen.getByRole('textbox'), { target: { value: '运营' } });

    expect(search).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole('button', { name: '搜索' }));
    expect(search).toHaveBeenCalledWith('运营');

    fireEvent.change(screen.getByRole('textbox'), { target: { value: '财务' } });
    expect(search).toHaveBeenCalledTimes(1);
    fireEvent.submit(screen.getByRole('textbox').closest('form')!);
    expect(search).toHaveBeenLastCalledWith('财务');
  });
});
