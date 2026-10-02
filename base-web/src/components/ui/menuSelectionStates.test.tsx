// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import {
  ContextMenu,
  ContextMenuCheckboxItem,
  ContextMenuContent,
  ContextMenuRadioGroup,
  ContextMenuRadioItem,
  ContextMenuTrigger,
} from './context-menu';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from './dropdown-menu';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './select';

afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
});

function expectCheckedPrimaryPortal(element: HTMLElement, container: HTMLElement) {
  expect(element.dataset.state).toBe('checked');
  expect(element.className.split(' ')).toContain('data-[state=checked]:text-primary');
  expect(container.contains(element)).toBe(false);
}

describe('菜单持久选中状态', () => {
  it('SelectItem 在真实 Portal 中呈现主题选中态', () => {
    const { container } = render(
      <Select value="ACTIVE">
        <SelectTrigger aria-label="准入状态"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem value="ACTIVE">已启用</SelectItem>
          <SelectItem value="SUSPENDED">已停用</SelectItem>
        </SelectContent>
      </Select>,
    );

    fireEvent.click(screen.getByRole('combobox', { name: '准入状态' }));
    expectCheckedPrimaryPortal(screen.getByRole('option', { name: '已启用' }), container);
  });

  it('DropdownMenuRadioItem 在真实 Portal 中呈现主题选中态', () => {
    const { container } = render(
      <DropdownMenu open>
        <DropdownMenuTrigger>外观</DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuRadioGroup value="light">
            <DropdownMenuRadioItem value="light">浅色</DropdownMenuRadioItem>
            <DropdownMenuRadioItem value="dark">深色</DropdownMenuRadioItem>
          </DropdownMenuRadioGroup>
        </DropdownMenuContent>
      </DropdownMenu>,
    );

    expectCheckedPrimaryPortal(screen.getByRole('menuitemradio', { name: '浅色' }), container);
  });

  it('ContextMenu 复选和单选项在真实 Portal 中呈现主题选中态', () => {
    vi.spyOn(console, 'warn').mockImplementation(() => undefined);
    const { container } = render(
      <ContextMenu open>
        <ContextMenuTrigger>标签页</ContextMenuTrigger>
        <ContextMenuContent>
          <ContextMenuCheckboxItem checked>固定</ContextMenuCheckboxItem>
          <ContextMenuRadioGroup value="current">
            <ContextMenuRadioItem value="current">当前</ContextMenuRadioItem>
          </ContextMenuRadioGroup>
        </ContextMenuContent>
      </ContextMenu>,
    );

    expectCheckedPrimaryPortal(screen.getByRole('menuitemcheckbox', { name: '固定' }), container);
    expectCheckedPrimaryPortal(screen.getByRole('menuitemradio', { name: '当前' }), container);
  });
});
