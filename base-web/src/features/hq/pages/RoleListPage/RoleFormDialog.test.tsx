// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { useState } from 'react';
import { afterAll, afterEach, beforeAll, describe, expect, it, vi } from 'vitest';
import { RoleFormDialog, type RoleValues } from './RoleFormDialog';

const permissions = [
  { id: 'member-read', name: 'permission.hqMembership.read', action: 'hqMembership:read', module: 'hqMembership', scope: 'SYSTEM' as const },
  { id: 'member-create', name: 'permission.hqMembership.create', action: 'hqMembership:create', module: 'hqMembership', scope: 'SYSTEM' as const },
  { id: 'store-read', name: 'permission.hqStore.read', action: 'hqStore:read', module: 'hqStore', scope: 'SYSTEM' as const },
];

beforeAll(() => {
  vi.stubGlobal('ResizeObserver', class {
    observe() {}
    unobserve() {}
    disconnect() {}
  });
});

afterEach(cleanup);
afterAll(() => vi.unstubAllGlobals());

function RoleDialogHarness() {
  const [values, setValues] = useState<RoleValues>({ name: '', permissionIds: [] });
  return <RoleFormDialog
    open
    mode="create"
    onOpenChange={() => undefined}
    permissions={permissions}
    viewerPermissions={permissions.map((permission) => permission.action)}
    values={values}
    onValuesChange={setValues}
    onSubmit={() => undefined}
  />;
}

describe('总部角色权限弹窗', () => {
  it('在视口内呈现紧凑的分组权限选择器', () => {
    render(<RoleDialogHarness />);

    expect(screen.getByRole('dialog').className).toContain('max-h-[calc(100dvh-2rem)]');
    expect(screen.getByRole('heading', { name: '新增角色' })).toBeTruthy();
    expect(screen.getByTestId('role-permission-scroll')).toBeTruthy();
    expect(document.querySelectorAll('fieldset')).toHaveLength(0);
    expect(screen.getByText('管理员')).toBeTruthy();
    expect(screen.getByText('直营门店')).toBeTruthy();
    expect(screen.getByText('查看管理员')).toBeTruthy();
    expect(screen.queryByText('(hqMembership:read)')).toBeNull();
    expect(screen.getByText('(0/3)')).toBeTruthy();
  });

  it('滚动权限列表时将全选栏固定在顶部', () => {
    render(<RoleDialogHarness />);

    const selectAllBar = screen.getByRole('checkbox', { name: '全选全部权限' }).parentElement;

    expect(selectAllBar?.className).toContain('sticky');
    expect(selectAllBar?.className).toContain('top-0');
    expect(selectAllBar?.className).toContain('z-10');
    expect(selectAllBar?.className).toContain('bg-card');
  });

  it('支持全选、分组半选和重复清除', () => {
    render(<RoleDialogHarness />);

    const read = screen.getByRole('checkbox', { name: '查看管理员' });
    const group = screen.getByRole('checkbox', { name: '全选 管理员 全部权限' });
    const all = screen.getByRole('checkbox', { name: '全选全部权限' });

    fireEvent.click(read);
    expect(group.getAttribute('aria-checked')).toBe('mixed');
    expect(all.getAttribute('aria-checked')).toBe('mixed');

    fireEvent.click(group);
    expect(group.getAttribute('aria-checked')).toBe('true');
    fireEvent.click(group);
    expect(group.getAttribute('aria-checked')).toBe('false');

    fireEvent.click(all);
    expect(all.getAttribute('aria-checked')).toBe('true');
    fireEvent.click(all);
    expect(all.getAttribute('aria-checked')).toBe('false');
  });
});
