// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { FranchiseListPage } from './FranchiseListPage';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';

const mocks = vi.hoisted(() => ({ suspend: vi.fn(), restore: vi.fn(), refetch: vi.fn(), status: { value: 'ACTIVE' } }));
const renderPage = () => render(<AdminToastProvider><MemoryRouter><FranchiseListPage /></MemoryRouter></AdminToastProvider>);

vi.mock('@apollo/client/react', () => ({
  useQuery: (document: { definitions?: Array<{ name?: { value?: string } }> }) => {
    const name = document.definitions?.find((item) => item.name)?.name?.value;
    if (name === 'HqFranchises') return { data: { organizations: { data: [{ id: 'org-1', code: 'F001', name: '测试加盟商', status: mocks.status.value, initialAccountId: null, initialAccount: null }], total: 1 } }, loading: false, refetch: mocks.refetch };
    return { data: undefined, loading: false };
  },
  useMutation: (document: { definitions?: Array<{ name?: { value?: string } }> }) => {
    const name = document.definitions?.find((item) => item.name)?.name?.value;
    return [name === 'HqSuspendOrganization' ? mocks.suspend : name === 'HqRestoreOrganization' ? mocks.restore : vi.fn(), { loading: false }];
  },
}));
vi.mock('@/features/admin/hooks/useAdminTab', () => ({ useAdminTab: () => ({ search: '' }) }));
vi.mock('@/features/auth/store/authStore', () => ({ useAuthStore: (selector: (state: unknown) => unknown) => selector({ viewer: { permissions: ['organization:suspend', 'organization:restore'] } }) }));

beforeEach(() => { vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }); mocks.suspend.mockReset(); mocks.restore.mockReset(); mocks.refetch.mockReset(); mocks.refetch.mockResolvedValue({}); mocks.status.value = 'ACTIVE'; });
afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.unstubAllGlobals(); });

describe('加盟商状态开关', () => {
  it('关闭时提交暂停原因，并在请求期间禁止重复操作', async () => {
    const nativePrompt = vi.spyOn(window, 'prompt');
    const nativeConfirm = vi.spyOn(window, 'confirm');
    mocks.suspend.mockImplementation(() => new Promise(() => undefined));
    renderPage();
    const status = screen.getByRole('switch', { name: '测试加盟商状态' });
    fireEvent.click(status);
    expect(screen.getByRole('alertdialog').textContent).toContain('暂停测试加盟商');
    expect(mocks.suspend).not.toHaveBeenCalled();
    expect(screen.getByRole('button', { name: '确认暂停' }).hasAttribute('disabled')).toBe(true);
    const reason = screen.getByRole('textbox', { name: '暂停原因' });
    expect(reason).toBeInstanceOf(HTMLTextAreaElement);
    fireEvent.change(reason, { target: { value: 'CONTRACT_ENDED\n暂停运营' } });
    fireEvent.click(screen.getByRole('button', { name: '确认暂停' }));
    expect(mocks.suspend).toHaveBeenCalledWith({ variables: { input: { organizationId: 'org-1', reasonCode: 'CONTRACT_ENDED\n暂停运营' } } });
    expect(nativePrompt).not.toHaveBeenCalled();
    expect(nativeConfirm).not.toHaveBeenCalled();
    await waitFor(() => expect(screen.getByRole('button', { name: '暂停中…' }).hasAttribute('disabled')).toBe(true));
    expect(screen.getByRole('alertdialog')).toBeTruthy();
    fireEvent.click(status);
    expect(mocks.suspend).toHaveBeenCalledTimes(1);
  });
  it('取消暂停时保持开启且不发送请求', () => {
    renderPage();
    fireEvent.click(screen.getByRole('switch', { name: '测试加盟商状态' }));
    fireEvent.click(screen.getByRole('button', { name: '取消' }));
    expect(mocks.suspend).not.toHaveBeenCalled();
    expect(screen.queryByRole('alertdialog')).toBeNull();
    expect(screen.getByRole('switch', { name: '测试加盟商状态' }).getAttribute('aria-checked')).toBe('true');
  });
});

describe('加盟商状态变更结果', () => {
  it('暂停失败时保留原因和对话框以便重试', async () => {
    mocks.suspend.mockRejectedValue(new Error('offline'));
    renderPage();
    const status = screen.getByRole('switch', { name: '测试加盟商状态' });
    fireEvent.click(status);
    fireEvent.change(screen.getByPlaceholderText('请输入暂停原因'), { target: { value: '合同到期' } });
    fireEvent.click(screen.getByRole('button', { name: '确认暂停' }));
    expect((await screen.findByRole('alert')).textContent).toContain('操作失败');
    expect(screen.getByRole('alertdialog')).toBeTruthy();
    expect((screen.getByPlaceholderText('请输入暂停原因') as HTMLTextAreaElement).value).toBe('合同到期');
    expect(status.getAttribute('aria-checked')).toBe('true');
  });
  it('暂停成功并刷新后关闭对话框，开关显示关闭', async () => {
    mocks.suspend.mockResolvedValue({});
    mocks.refetch.mockImplementation(async () => { mocks.status.value = 'SUSPENDED'; return {}; });
    renderPage();
    const status = screen.getByRole('switch', { name: '测试加盟商状态' });
    fireEvent.click(status);
    fireEvent.change(screen.getByPlaceholderText('请输入暂停原因'), { target: { value: '合同到期' } });
    fireEvent.click(screen.getByRole('button', { name: '确认暂停' }));
    await waitFor(() => expect(screen.queryByRole('alertdialog')).toBeNull());
    expect(mocks.refetch).toHaveBeenCalledTimes(1);
    expect(screen.getByRole('switch', { name: '测试加盟商状态' }).getAttribute('aria-checked')).toBe('false');
  });
  it('开启时恢复加盟商，失败则保持关闭并提示错误', async () => {
    mocks.status.value = 'SUSPENDED';
    mocks.restore.mockRejectedValue(new Error('offline'));
    renderPage();
    const status = screen.getByRole('switch', { name: '测试加盟商状态' });
    fireEvent.click(status);
    expect(mocks.restore).toHaveBeenCalledWith({ variables: { id: 'org-1' } });
    expect((await screen.findByRole('alert')).textContent).toContain('操作失败');
    expect(status.getAttribute('aria-checked')).toBe('false');
  });
});
