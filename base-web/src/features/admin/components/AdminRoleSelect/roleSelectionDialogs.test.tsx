// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { AdministratorFormDialog } from '@/features/hq/pages/AdministratorListPage/AdministratorFormDialog';
import { InviteStaffDialog } from '@/features/franchise/pages/StaffListPage/InviteStaffDialog';

beforeEach(() => vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }));
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

it('selects multiple administrator roles from a dropdown', () => {
  const onValuesChange = vi.fn();
  render(<AdministratorFormDialog open mode="create" onOpenChange={() => undefined} onSubmit={() => undefined}
    roles={[{ id: 'role-1', name: '运营', kind: 'CUSTOM' }, { id: 'role-2', name: '财务', kind: 'CUSTOM' }]}
    values={{ phone: '', displayName: '', email: '', roleIds: ['role-1'] }} onValuesChange={onValuesChange} />);
  expect(screen.getByRole('button', { name: '角色：运营' }).hasAttribute('data-admin-form-surface')).toBe(true);
  fireEvent.keyDown(screen.getByRole('button', { name: '角色：运营' }), { key: 'ArrowDown' });
  expect(screen.getByRole('menu').hasAttribute('data-admin-form-surface')).toBe(true);
  expect(screen.getByRole('menuitemcheckbox', { name: '运营' }).getAttribute('aria-checked')).toBe('true');
  fireEvent.click(screen.getByRole('menuitemcheckbox', { name: '财务' }));
  expect(onValuesChange).toHaveBeenCalledWith(expect.objectContaining({ roleIds: ['role-1', 'role-2'] }));
  expect(screen.getByRole('menuitemcheckbox', { name: '运营' })).toBeTruthy();
});

it('uses the role dropdown and new employee wording in the franchise create dialog', () => {
  const onValuesChange = vi.fn();
  render(<InviteStaffDialog open onOpenChange={() => undefined} onSubmit={() => undefined}
    roles={[{ id: 'owner', name: 'role.franchiseOwner', kind: 'FRANCHISE_OWNER' }, { id: 'clerk', name: '店员', kind: 'CUSTOM' }]}
    values={{ phone: '', displayName: '', email: '', roleIds: ['owner'], storeAccessMode: 'ALL_STORES', storeIds: [] }}
    onValuesChange={onValuesChange} />);
  expect(screen.getByRole('heading', { name: '新增员工' })).toBeTruthy();
  expect(screen.getByRole('button', { name: '新增员工' })).toBeTruthy();
  fireEvent.keyDown(screen.getByRole('button', { name: '角色：加盟商负责人' }), { key: 'ArrowDown' });
  expect(screen.getByRole('menuitemcheckbox', { name: '加盟商负责人' }).getAttribute('aria-checked')).toBe('true');
  fireEvent.click(screen.getByRole('menuitemcheckbox', { name: '店员' }));
  expect(onValuesChange).toHaveBeenCalledWith(expect.objectContaining({ roleIds: ['owner', 'clerk'] }));
});

it('announces selected role names when the dropdown is closed', () => {
  render(<AdministratorFormDialog open mode="create" onOpenChange={() => undefined} onSubmit={() => undefined}
    roles={[{ id: 'role-1', name: '运营', kind: 'CUSTOM' }, { id: 'role-2', name: '财务', kind: 'CUSTOM' }]}
    values={{ phone: '', displayName: '', email: '', roleIds: ['role-1', 'role-2'] }} />);

  expect(screen.getByRole('button', { name: '角色：运营、财务' })).toBeTruthy();
});
