// @vitest-environment jsdom

import { useState } from 'react';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { Label } from '@/components/ui/label';
import { AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { AdminFormSelect } from './AdminFormSelect';

beforeEach(() => { vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }); });
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

it('keeps the placeholder out of the options while preserving label association', () => {
  render(<div><Label htmlFor="choice">分类</Label><AdminFormSelect id="choice" value="" placeholder="请选择分类"
    options={[{ value: 'food', label: '食品' }]} onValueChange={vi.fn()} /></div>);
  expect(screen.getByRole('combobox', { name: '分类' }).textContent).toContain('请选择分类');
  fireEvent.click(screen.getByRole('combobox', { name: '分类' }));
  expect(screen.queryByRole('option', { name: '请选择分类' })).toBeNull();
  expect(screen.getAllByRole('option').map((option) => option.textContent)).toEqual(['食品']);
});

it('maps only an explicitly declared empty choice to an empty string', () => {
  const changes: string[] = [];
  // 显式空值选项用于取消业务关联，占位文案只提供未选择提示。
  function Form() {
    const [value, setValue] = useState('');
    return <div><Label htmlFor="choice">分类</Label><AdminFormSelect id="choice" value={value} placeholder="请选择分类" clearLabel="不选择"
      options={[{ value: 'food', label: '食品' }]} onValueChange={(next) => { changes.push(next); setValue(next); }} /></div>;
  }
  render(<Form />);
  fireEvent.click(screen.getByRole('combobox', { name: '分类' }));
  fireEvent.click(screen.getByRole('option', { name: '食品' }));
  expect(screen.getByRole('combobox', { name: '分类' }).textContent).toContain('食品');
  fireEvent.click(screen.getByRole('combobox', { name: '分类' }));
  expect(screen.queryByRole('option', { name: '请选择分类' })).toBeNull();
  fireEvent.click(screen.getByRole('option', { name: '不选择' }));
  expect(changes).toEqual(['food', '']);
});

it('keeps unavailable options disabled', () => {
  render(<AdminFormSelect id="choice" value="" placeholder="请选择" options={[{ value: 'old', label: '已停用', disabled: true }]} onValueChange={vi.fn()} />);
  fireEvent.click(screen.getByRole('combobox'));
  expect(screen.getByRole('option', { name: '已停用' }).getAttribute('aria-disabled')).toBe('true');
});

it('keeps the portal option panel white in admin forms', () => {
  render(<AdminFormDialogShell open onOpenChange={() => undefined} title="表单" submitLabel="保存" cancelLabel="取消" onSubmit={() => undefined}>
    <Label htmlFor="choice">门店范围</Label>
    <AdminFormSelect id="choice" value="ALL_STORES" placeholder="请选择门店范围"
      options={[{ value: 'ALL_STORES', label: '全部门店' }, { value: 'SELECTED_STORES', label: '指定门店' }]}
      onValueChange={() => undefined} />
  </AdminFormDialogShell>);

  fireEvent.click(screen.getByRole('combobox', { name: '门店范围' }));
  const panel = document.querySelector('[data-slot="select-content"]');
  expect(panel).not.toBeNull();
  expect(document.querySelector('[data-testid="admin-form-dialog-content"]')?.contains(panel)).toBe(false);
  expect(panel?.hasAttribute('data-admin-form-surface')).toBe(true);
  expect(panel?.className).not.toContain('!bg-white');
  expect(panel?.className).not.toContain('!text-black');
});
