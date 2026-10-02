// @vitest-environment jsdom

import type { ReactNode } from 'react';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { AdminHeader } from './AdminHeader';

const mocks = vi.hoisted(() => ({
  logout: vi.fn(),
  navigate: vi.fn(),
  resetAuth: vi.fn(),
  resetSession: vi.fn(async () => undefined),
}));

vi.mock('@apollo/client/react', () => ({
  useMutation: () => [mocks.logout, { loading: false }],
}));

vi.mock('react-router', () => ({
  useNavigate: () => mocks.navigate,
}));

vi.mock('@/components/shared/ThemeToggle', () => ({ ThemeToggle: () => <button>外观</button> }));
vi.mock('@/components/ui/avatar', () => ({
  Avatar: ({ children }: { children: ReactNode }) => <span>{children}</span>,
  AvatarFallback: ({ children }: { children: ReactNode }) => <span>{children}</span>,
}));
vi.mock('@/components/ui/sidebar', () => ({ SidebarTrigger: () => <button>侧边栏</button> }));
vi.mock('@/features/admin/hooks/useAdminTabs', () => ({ useActiveTabPath: () => '/admin/hq' }));
vi.mock('@/features/auth/store/authStore', () => ({
  useAuthStore: (selector: (state: unknown) => unknown) => selector({
    viewer: {
      account: { id: 'account-1', displayName: '王店长', email: 'owner@example.com' },
      currentWorkspace: { workspaceType: 'HEADQUARTERS', organizationId: 'hq' },
      permissions: [],
    },
    reset: mocks.resetAuth,
  }),
}));
vi.mock('@/lib/graphql/client', () => ({
  getGraphQLRuntime: () => ({ resetSession: mocks.resetSession }),
}));

beforeEach(() => {
  vi.clearAllMocks();
  window.history.replaceState({}, '', '/admin/hq/stores?page=2#active');
  mocks.logout.mockResolvedValue({ data: { logout: true } });
});
afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
});

describe('后台头像账号菜单', () => {
  it('账号面板使用可聚焦分组而不是无切换行为的按钮', () => {
    render(<AdminHeader onOpenAi={vi.fn()} />);

    const trigger = screen.getByRole('group', { name: '账号操作：王店长' });
    expect(trigger.tabIndex).toBe(0);
    expect(screen.queryByRole('button', { name: '账号菜单：王店长' })).toBeNull();
    expect(trigger.className).toContain('size-8');
    expect(trigger.className).not.toContain('size-11');
    const accountMenu = trigger.closest('[data-account-menu]');
    expect(accountMenu?.className).toContain('group/account');
    const actions = screen.getByRole('region', { name: '账号操作' });
    expect(actions.parentElement?.className).toContain('group-hover/account:visible');
    expect(actions.parentElement?.className).toContain('group-focus-within/account:visible');
    expect(actions.parentElement?.className).toContain('pt-2');
    expect(screen.getByRole('button', { name: '修改密码' }).getAttribute('role')).toBeNull();
    fireEvent.pointerEnter(trigger);
    expect(document.body.style.pointerEvents).toBe('');
    expect(screen.getByText('owner@example.com')).toBeTruthy();
  });

  it('点击非聚焦区域也会清除账号区域焦点并退出 CSS 展开条件', () => {
    render(<AdminHeader onOpenAi={vi.fn()} />);

    const trigger = screen.getByRole('group', { name: '账号操作：王店长' });
    const accountMenu = trigger.closest('[data-account-menu]');
    trigger.focus();
    expect(accountMenu?.matches(':focus-within')).toBe(true);

    fireEvent.pointerDown(document.body);

    expect(accountMenu?.contains(document.activeElement)).toBe(false);
    expect(document.body.style.pointerEvents).toBe('');
  });
});

describe('后台账号操作', () => {
  it('右上角 AI 按钮打开全局抽屉', () => {
    const open = vi.fn();
    render(<AdminHeader onOpenAi={open} />);
    fireEvent.click(screen.getByRole('button', { name: '打开 AI 助手' }));
    expect(open).toHaveBeenCalledTimes(1);
  });
  it('修改密码在当前后台骨架内打开弹窗而不跳转路由', async () => {
    render(<AdminHeader onOpenAi={vi.fn()} />);
    fireEvent.pointerEnter(screen.getByRole('group', { name: '账号操作：王店长' }));
    fireEvent.click(screen.getByRole('button', { name: '修改密码' }));

    expect(await screen.findByRole('heading', { name: '修改密码' })).toBeTruthy();
    expect(mocks.navigate).not.toHaveBeenCalled();
  });

  it('服务端退出成功后清理会话并回到登录页', async () => {
    render(<AdminHeader onOpenAi={vi.fn()} />);
    fireEvent.pointerEnter(screen.getByRole('group', { name: '账号操作：王店长' }));
    fireEvent.click(screen.getByRole('button', { name: '退出登录' }));

    await waitFor(() => expect(mocks.logout).toHaveBeenCalledTimes(1));
    expect(mocks.resetSession).toHaveBeenCalledWith('LOGOUT');
    expect(mocks.navigate).toHaveBeenCalledWith('/admin/login', { replace: true });
  });

  it('服务端退出成功后即使本地缓存清理失败也会回到登录页', async () => {
    const diagnostics = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    mocks.resetSession.mockRejectedValueOnce(new Error('cache cleanup failed'));
    render(<AdminHeader onOpenAi={vi.fn()} />);
    fireEvent.click(screen.getByRole('button', { name: '退出登录' }));

    await waitFor(() => expect(mocks.logout).toHaveBeenCalledTimes(1));
    expect(mocks.navigate).toHaveBeenCalledWith('/admin/login', { replace: true });
    expect(screen.queryByRole('alert')).toBeNull();
    expect(diagnostics).toHaveBeenCalledWith(
      '后台会话清理失败，继续执行安全跳转。',
      expect.any(Error),
    );
  });

  it('退出失败时保留当前会话并提供重试提示', async () => {
    mocks.logout.mockRejectedValueOnce(new Error('offline'));
    render(<AdminHeader onOpenAi={vi.fn()} />);
    fireEvent.pointerEnter(screen.getByRole('group', { name: '账号操作：王店长' }));
    fireEvent.click(screen.getByRole('button', { name: '退出登录' }));

    expect((await screen.findByRole('alert')).textContent).toBe('退出失败，请重试。');
    expect(mocks.resetSession).not.toHaveBeenCalled();
    expect(mocks.navigate).not.toHaveBeenCalledWith('/admin/login', { replace: true });
  });
});
