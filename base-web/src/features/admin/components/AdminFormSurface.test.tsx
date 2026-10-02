// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';

beforeEach(() => {
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} });
});
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

it('标准表单控件和 Portal 选项层暴露统一后台表面语义', () => {
  render(<div>
    <Input aria-label="名称" placeholder="请输入名称" disabled />
    <Textarea aria-label="说明" placeholder="请输入说明" readOnly />
    <Select value="OPEN" onValueChange={() => undefined}>
      <SelectTrigger aria-label="营业状态"><SelectValue /></SelectTrigger>
      <SelectContent><SelectItem value="OPEN">正常营业</SelectItem></SelectContent>
    </Select>
  </div>);

  expect(screen.getByRole('textbox', { name: '名称' }).hasAttribute('data-admin-form-surface')).toBe(true);
  expect(screen.getByRole('textbox', { name: '说明' }).hasAttribute('data-admin-form-surface')).toBe(true);
  const trigger = screen.getByRole('combobox', { name: '营业状态' });
  expect(trigger.hasAttribute('data-admin-form-surface')).toBe(true);
  fireEvent.click(trigger);
  expect(document.querySelector('[data-slot="select-content"]')?.hasAttribute('data-admin-form-surface')).toBe(true);
});
