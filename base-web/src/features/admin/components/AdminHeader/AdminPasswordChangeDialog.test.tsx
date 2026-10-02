// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { AdminPasswordChangeDialog } from './AdminPasswordChangeDialog';

const mocks = vi.hoisted(() => ({
  changePassword: vi.fn(),
  selectWorkspace: vi.fn(),
  setViewer: vi.fn(),
  disposeRuntime: vi.fn(async () => undefined),
  completePasswordChange: vi.fn(async () => undefined),
}));

const currentViewer = {
  account: { id: 'account-1', phone: '13800000000', displayName: '王店长', email: 'owner@example.com', status: 'ACTIVE' as const, mustChangePassword: false },
  currentWorkspace: { workspaceType: 'HEADQUARTERS' as const, organizationId: 'hq', organizationName: '总部', homePath: '/admin/hq' },
  workspaces: [],
  permissions: ['dashboard:read'],
};

const discoveryViewer = { ...currentViewer, currentWorkspace: null };
const restoredViewer = { ...currentViewer };

vi.mock('@apollo/client/react', () => ({
  useMutation: (document: { definitions?: Array<{ name?: { value?: string } }> }) => {
    const name = document.definitions?.find((definition) => definition.name)?.name?.value;
    return name === 'ChangeTemporaryPassword'
      ? [mocks.changePassword, { loading: false }]
      : [mocks.selectWorkspace, { loading: false }];
  },
}));

vi.mock('@/__generated__', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/__generated__')>();
  return { ...actual, useFragment: (_fragment: unknown, value: unknown) => value };
});

vi.mock('@/features/auth/store/authStore', () => ({
  useAuthStore: (selector: (state: unknown) => unknown) => selector({ viewer: currentViewer, setViewer: mocks.setViewer }),
}));

vi.mock('@/lib/graphql/client', () => ({ disposeGraphQLRuntime: mocks.disposeRuntime }));

vi.mock('@/features/auth/pages/authFlow', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/features/auth/pages/authFlow')>();
  return { ...actual, completePasswordChange: mocks.completePasswordChange };
});

beforeEach(() => {
  vi.clearAllMocks();
  window.history.replaceState({}, '', '/admin/hq/stores?page=2#active');
  mocks.changePassword.mockResolvedValue({ data: { changeTemporaryPassword: discoveryViewer } });
  mocks.selectWorkspace.mockResolvedValue({ data: { selectWorkspace: restoredViewer } });
});
afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
});

describe('后台骨架内修改密码', () => {
  it('改密后恢复当前 workspace 并返回原后台地址', async () => {
    render(<AdminPasswordChangeDialog open onOpenChange={vi.fn()} />);
    fireEvent.change(screen.getByLabelText('当前密码'), { target: { value: 'Old-Password-42!' } });
    fireEvent.change(screen.getByLabelText('新密码'), { target: { value: 'New-Password-42!' } });
    fireEvent.change(screen.getByLabelText('确认新密码'), { target: { value: 'New-Password-42!' } });
    fireEvent.click(screen.getByRole('button', { name: '保存新密码' }));

    await waitFor(() => expect(mocks.changePassword).toHaveBeenCalledWith({
      variables: { input: { currentPassword: 'Old-Password-42!', newPassword: 'New-Password-42!' } },
    }));
    expect(mocks.selectWorkspace).toHaveBeenCalledWith({
      variables: { input: { workspaceType: 'HEADQUARTERS', organizationId: 'hq' } },
    });
    expect(mocks.completePasswordChange).toHaveBeenCalledWith(
      restoredViewer,
      mocks.setViewer,
      mocks.disposeRuntime,
      expect.any(Function),
      '/admin/hq/stores?page=2#active',
    );
  });

  it('改密成功但 workspace 恢复失败时进入安全的工作台选择流程', async () => {
    const diagnostics = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    mocks.selectWorkspace.mockRejectedValueOnce(new Error('offline'));
    render(<AdminPasswordChangeDialog open onOpenChange={vi.fn()} />);
    fireEvent.change(screen.getByLabelText('当前密码'), { target: { value: 'Old-Password-42!' } });
    fireEvent.change(screen.getByLabelText('新密码'), { target: { value: 'New-Password-42!' } });
    fireEvent.change(screen.getByLabelText('确认新密码'), { target: { value: 'New-Password-42!' } });
    fireEvent.click(screen.getByRole('button', { name: '保存新密码' }));

    await waitFor(() => expect(mocks.completePasswordChange).toHaveBeenCalledWith(
      discoveryViewer,
      mocks.setViewer,
      mocks.disposeRuntime,
      expect.any(Function),
      '/admin/hq/stores?page=2#active',
    ));
    expect(diagnostics).toHaveBeenCalledWith(
      '改密后的 workspace 恢复失败，继续进入工作台选择流程。',
      expect.any(Error),
    );
  });
});
