// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { ProvisionFranchiseDialog } from '@/features/hq/pages/FranchiseListPage/ProvisionFranchiseDialog';
import { AdministratorFormDialog } from '@/features/hq/pages/AdministratorListPage/AdministratorFormDialog';
import { InviteStaffDialog } from '@/features/franchise/pages/StaffListPage/InviteStaffDialog';
import { CopyTemporaryPasswordButton } from './CopyTemporaryPasswordButton';
import { AdminToastProvider } from './AdminToast';

const password = 'Abc12345';
const clipboardDescriptor = Object.getOwnPropertyDescriptor(navigator, 'clipboard');
const execCommandDescriptor = Object.getOwnPropertyDescriptor(document, 'execCommand');
const renderWithToast = (component: React.ReactElement) => render(<AdminToastProvider>{component}</AdminToastProvider>);

afterEach(() => {
  cleanup();
  if (clipboardDescriptor) Object.defineProperty(navigator, 'clipboard', clipboardDescriptor);
  else Reflect.deleteProperty(navigator, 'clipboard');
  if (execCommandDescriptor) Object.defineProperty(document, 'execCommand', execCommandDescriptor);
  else Reflect.deleteProperty(document, 'execCommand');
  vi.restoreAllMocks();
});

describe('临时密码复制反馈', () => {
  const dialogs = [
    ['开通加盟商', <ProvisionFranchiseDialog open onOpenChange={() => undefined} onSubmit={() => undefined} outcome={{ kind: 'temporary-password', value: password }} />],
    ['新增管理员', <AdministratorFormDialog open mode="create" onOpenChange={() => undefined} onSubmit={() => undefined} roles={[]} outcome={{ kind: 'temporary-password', value: password }} />],
    ['新增员工', <InviteStaffDialog open onOpenChange={() => undefined} onSubmit={() => undefined} outcome={{ kind: 'temporary-password', value: password }} />],
  ] as const;

  it.each(dialogs)('%s 复制成功后显示反馈', async (_, dialog) => {
    const writeText = vi.fn().mockResolvedValue(undefined);
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText } });
    renderWithToast(dialog);
    fireEvent.click(screen.getByRole('button', { name: '复制临时密码' }));
    const status = await screen.findByRole('status');
    expect(status).toHaveProperty('textContent', '复制成功');
    expect(status.className).toContain('bg-success-bg');
    expect(status.querySelector('svg')).not.toBeNull();
    expect(screen.getByRole('dialog').querySelector('[role="status"]')).toBeNull();
    expect(writeText).toHaveBeenCalledWith(password);
  });

  it('剪贴板接口不可用时使用浏览器兼容复制方式', async () => {
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: undefined });
    const execCommand = vi.fn(() => {
      const selected = document.activeElement as HTMLTextAreaElement;
      return selected.value === password && Boolean(selected.closest('[role="dialog"]'));
    });
    Object.defineProperty(document, 'execCommand', { configurable: true, value: execCommand });
    renderWithToast(dialogs[0][1]);
    fireEvent.click(screen.getByRole('button', { name: '复制临时密码' }));
    expect(await screen.findByRole('status')).toHaveProperty('textContent', '复制成功');
    expect(execCommand).toHaveBeenCalledWith('copy');
  });

  it('两种复制方式都失败时提示手动复制', async () => {
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText: vi.fn().mockRejectedValue(new Error('denied')) } });
    Object.defineProperty(document, 'execCommand', { configurable: true, value: vi.fn(() => false) });
    renderWithToast(<CopyTemporaryPasswordButton password={password} />);
    fireEvent.click(screen.getByRole('button', { name: '复制临时密码' }));
    expect(await screen.findByRole('alert')).toHaveProperty('textContent', '复制失败，请手动选择上方密码复制');
  });

  it('脱敏列表复制失败时提示重试', async () => {
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText: vi.fn().mockRejectedValue(new Error('denied')) } });
    Object.defineProperty(document, 'execCommand', { configurable: true, value: vi.fn(() => false) });
    renderWithToast(<CopyTemporaryPasswordButton password={password} failureMessage="复制失败，请重试" iconOnly />);
    fireEvent.click(screen.getByRole('button', { name: '复制临时密码' }));
    expect(await screen.findByRole('alert')).toHaveProperty('textContent', '复制失败，请重试');
    expect(screen.getByRole('alert').className).toContain('bg-danger-bg');
  });
});
