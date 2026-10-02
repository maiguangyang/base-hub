// @vitest-environment jsdom

import { useState } from 'react';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { AdminFilterSelect } from './AdminFilterSelect';
import { AdminListSearch } from './AdminListSearch';
import { AdminDateFilter } from './AdminDateFilter';
import { FranchiseSearchSelect } from '@/features/hq/pages/FranchiseStoreListPage/FranchiseSearchSelect';
import { AdminSearchSelect } from './AdminSearchSelect';

beforeEach(() => vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }));
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

const options = [{ value: 'ACTIVE', label: '启用' }];

it('单选清空后恢复占位提示，并且不打开下拉或提交外围表单', () => {
  const submit = vi.fn();
  function Filter() {
    const [value, setValue] = useState<string>();
    return <form onSubmit={submit}><AdminFilterSelect value={value} placeholder="全部状态" options={options} onValueChange={setValue} /></form>;
  }
  render(<Filter />);
  expect(screen.queryByRole('button', { name: '清空全部状态' })).toBeNull();
  fireEvent.click(screen.getByRole('combobox'));
  fireEvent.click(screen.getByRole('option', { name: '启用' }));
  const clear = screen.getByRole('button', { name: '清空全部状态' });
  expect(clear.closest('[role="combobox"]')).toBeNull();
  fireEvent.click(clear);
  expect(screen.getByRole('combobox').textContent).toContain('全部状态');
  expect(screen.queryByRole('option')).toBeNull();
  expect(screen.queryByRole('button', { name: '清空全部状态' })).toBeNull();
  expect(submit).not.toHaveBeenCalled();
});

it('禁用筛选不能通过关闭按钮更改值', () => {
  const change = vi.fn();
  render(<AdminFilterSelect value="ACTIVE" placeholder="全部状态" options={options} disabled onValueChange={change} />);
  const clear = screen.getByRole('button', { name: '清空全部状态' }) as HTMLButtonElement;
  expect(clear.disabled).toBe(true);
  fireEvent.click(clear);
  expect(change).not.toHaveBeenCalled();
});

it.each(['草稿', ''])('搜索关闭按钮同时清空已提交关键词和草稿：%s', (draft) => {
  const search = vi.fn();
  render(<AdminListSearch value="旧关键词" placeholder="搜索角色" onSearch={search} />);
  const input = screen.getByRole('textbox') as HTMLInputElement;
  fireEvent.change(input, { target: { value: draft } });
  fireEvent.click(screen.getByRole('button', { name: '清空搜索角色' }));
  expect(input.value).toBe('');
  expect(search).toHaveBeenCalledExactlyOnceWith('');
});

it('未提交的搜索草稿也可清空', () => {
  const search = vi.fn();
  render(<AdminListSearch value="" placeholder="搜索角色" onSearch={search} />);
  expect(screen.queryByRole('button', { name: '清空搜索角色' })).toBeNull();
  fireEvent.change(screen.getByRole('textbox'), { target: { value: '草稿' } });
  fireEvent.click(screen.getByRole('button', { name: '清空搜索角色' }));
  expect((screen.getByRole('textbox') as HTMLInputElement).value).toBe('');
  expect(screen.queryByRole('button', { name: '清空搜索角色' })).toBeNull();
});

it('日期关闭按钮输出空字符串', () => {
  function Filter() {
    const [value, setValue] = useState('2026-09-30');
    return <AdminDateFilter label="开始日期" value={value} onChange={setValue} />;
  }
  render(<Filter />);
  fireEvent.click(screen.getByRole('button', { name: '清空开始日期' }));
  expect((screen.getByLabelText('开始日期') as HTMLInputElement).value).toBe('');
  expect(screen.queryByRole('button', { name: '清空开始日期' })).toBeNull();
});

it('可搜索加盟商筛选清空后不会打开选项', () => {
  function Filter() {
    const [value, setValue] = useState<string | undefined>('ACTIVE');
    return <FranchiseSearchSelect value={value} options={options} disabled={false} onValueChange={setValue} />;
  }
  render(<Filter />);
  fireEvent.click(screen.getByRole('button', { name: '清空全部加盟商' }));
  expect(screen.getByRole('combobox').textContent).toContain('全部加盟商');
  expect(screen.queryByRole('listbox')).toBeNull();
  expect(screen.queryByRole('button', { name: '清空全部加盟商' })).toBeNull();
});

it('共用搜索选择框在表单中默认不增加筛选关闭按钮', () => {
  render(<AdminSearchSelect value="ACTIVE" options={options} label="初始账号" placeholder="请选择账号"
    searchPlaceholder="搜索账号" emptyMessage="没有账号" onValueChange={vi.fn()} />);
  expect(screen.queryByRole('button', { name: '清空初始账号' })).toBeNull();
});
