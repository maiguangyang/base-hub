// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { FranchiseListPage } from './FranchiseListPage';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';

const mocks = vi.hoisted(() => ({ reset: vi.fn(), provision: vi.fn(), suspend: vi.fn(), restore: vi.fn(), save: vi.fn(), refetch: vi.fn(), clipboard: vi.fn(), status: { value: 'ACTIVE' }, initialAccount: { value: { id: 'account-1', phone: '13800000001' } as { id: string; phone: string } | null }, permissions: ['account:update', 'hqMembership:read', 'franchise:provision', 'organization:suspend', 'organization:restore'] }));
const renderPage = () => render(<AdminToastProvider><MemoryRouter><FranchiseListPage /></MemoryRouter></AdminToastProvider>);

vi.mock('@apollo/client/react', () => ({
  useQuery: (document: { definitions?: Array<{ name?: { value?: string } }> }) => {
    const name = document.definitions?.find((item) => item.name)?.name?.value;
    if (name === 'HqFranchiseInitialAccountCandidates') return { data: { organization: { memberships: [{ id: 'member-1', status: 'ACTIVE', account: { id: 'account-1', phone: '13800000001', displayName: '王先生', status: 'ACTIVE' } }, { id: 'member-2', status: 'ACTIVE', account: { id: 'account-2', phone: '13800000002', displayName: '李先生', status: 'ACTIVE' } }] } }, loading: false };
    return { data: { organizations: { data: [{ id: 'org-1', code: 'F001', name: '测试加盟商', status: mocks.status.value, initialAccountId: mocks.initialAccount.value?.id ?? null, initialAccount: mocks.initialAccount.value }], total: 1 } }, loading: false, refetch: mocks.refetch };
  },
  useMutation: (document: { definitions?: Array<{ name?: { value?: string } }> }) => {
    const name = document.definitions?.find((item) => item.name)?.name?.value;
    return [name === 'HqResetFranchiseInitialPassword' ? mocks.reset : name === 'HqProvisionFranchise' ? mocks.provision : name === 'HqSuspendOrganization' ? mocks.suspend : name === 'HqRestoreOrganization' ? mocks.restore : vi.fn(), { loading: false }];
  },
}));
vi.mock('@/features/admin/hooks/useAdminTab', () => ({ useAdminTab: () => ({ search: '' }) }));
vi.mock('@/features/auth/store/authStore', () => ({ useAuthStore: (selector: (state: unknown) => unknown) => selector({ viewer: { currentWorkspace: { workspaceType: 'HEADQUARTERS' }, permissions: mocks.permissions } }) }));

beforeEach(() => { vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }); mocks.reset.mockReset(); mocks.provision.mockReset(); mocks.suspend.mockReset(); mocks.restore.mockReset(); mocks.save.mockReset(); mocks.save.mockResolvedValue({ ok: true, json: async () => ({ organizationId: 'org-1', accountId: 'account-1' }) }); vi.stubGlobal('fetch', mocks.save); mocks.refetch.mockReset(); mocks.refetch.mockResolvedValue({}); mocks.clipboard.mockReset(); mocks.status.value = 'ACTIVE'; mocks.initialAccount.value = { id: 'account-1', phone: '13800000001' }; mocks.permissions = ['account:update', 'hqMembership:read', 'franchise:provision', 'organization:suspend', 'organization:restore']; });
afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.unstubAllGlobals(); });

describe('加盟商支付配置入口', () => {
  it('有读取权限时从该行打开右侧抽屉并请求该加盟商', async () => {
    mocks.permissions = ['paymentConfig:read'];
    mocks.save.mockResolvedValue({ ok: true, json: async () => ({ channels: [] }) });
    renderPage();
    fireEvent.click(screen.getByRole('button', { name: '支付配置' }));
    expect(screen.getByRole('dialog').textContent).toContain('测试加盟商');
    await waitFor(() => expect(mocks.save).toHaveBeenCalledWith(
      expect.stringContaining('scope=FRANCHISE&organizationId=org-1'), expect.objectContaining({ method: 'GET' }),
    ));
  });

  it('缺少支付读取权限时不显示行入口', () => {
    renderPage();
    expect(screen.queryByRole('button', { name: '支付配置' })).toBeNull();
  });
});

describe('加盟商重置密码操作', () => {
  it('历史账号初始化必须由操作人确认已人工核对原始开通记录', () => {
    mocks.initialAccount.value = null;
    renderPage();
    fireEvent.click(screen.getByRole('button', { name: '初始化账户' }));
    expect(screen.getByRole('dialog').textContent).toContain('人工核对');
    fireEvent.click(screen.getByRole('combobox', { name: '初始账号' }));
    fireEvent.click(screen.getByRole('option', { name: '王先生 · 13800000001' }));
    fireEvent.change(screen.getByPlaceholderText('请输入已核对的开通记录编号'), { target: { value: 'OPEN-001' } });
    expect(screen.getByRole('button', { name: '确认绑定' }).hasAttribute('disabled')).toBe(true);
    fireEvent.click(screen.getByRole('checkbox', { name: /已人工核对原始开通记录/ }));
    expect(screen.getByRole('button', { name: '确认绑定' }).hasAttribute('disabled')).toBe(false);
  });
  it('历史加盟商核对成员后提交初始账号确认', async () => {
    mocks.initialAccount.value = null;
    renderPage();
    fireEvent.click(screen.getByRole('button', { name: '初始化账户' }));
    fireEvent.click(screen.getByRole('combobox', { name: '初始账号' }));
    fireEvent.change(screen.getByPlaceholderText('搜索姓名或手机号'), { target: { value: '13800000001' } });
    fireEvent.click(screen.getByRole('option', { name: '王先生 · 13800000001' }));
    fireEvent.change(screen.getByPlaceholderText('请输入已核对的开通记录编号'), { target: { value: 'OPEN-001' } });
    fireEvent.click(screen.getByRole('checkbox', { name: /已人工核对原始开通记录/ }));
    fireEvent.click(screen.getByRole('button', { name: '确认绑定' }));
    expect(mocks.save).toHaveBeenCalledWith(expect.stringContaining('/api/franchise-initial-account'), expect.objectContaining({ method: 'POST', credentials: 'include', body: JSON.stringify({ organizationId: 'org-1', accountId: 'account-1', evidenceReference: 'OPEN-001', attested: true }) }));
  });
});

describe('历史账户初始化提交状态', () => {
  it('初始化提交期间只发送一次请求，刷新失败时告知写入已成功', async () => {
    mocks.initialAccount.value = null;
    let resolveSave: (value: unknown) => void = () => undefined;
    mocks.save.mockImplementation(() => new Promise((resolve) => { resolveSave = resolve; }));
    mocks.refetch.mockRejectedValue(new Error('refresh failed'));
    renderPage();
    fireEvent.click(screen.getByRole('button', { name: '初始化账户' }));
    fireEvent.click(screen.getByRole('combobox', { name: '初始账号' }));
    fireEvent.click(screen.getByRole('option', { name: '王先生 · 13800000001' }));
    fireEvent.change(screen.getByPlaceholderText('请输入已核对的开通记录编号'), { target: { value: 'OPEN-001' } });
    fireEvent.click(screen.getByRole('checkbox', { name: /已人工核对原始开通记录/ }));
    fireEvent.click(screen.getByRole('button', { name: '确认绑定' }));
    fireEvent.click(screen.getByRole('button', { name: '确认绑定' }));
    expect(mocks.save).toHaveBeenCalledTimes(1);
    resolveSave({ ok: true, json: async () => ({ organizationId: 'org-1', accountId: 'account-1' }) });
    await waitFor(() => expect(screen.getByRole('alert').textContent).toContain('初始账户已保存，但列表刷新失败'));
    expect(screen.getByRole('button', { name: '确认绑定' }).hasAttribute('disabled')).toBe(true);
    expect(mocks.save).toHaveBeenCalledTimes(1);
  });
  it('初始化请求未结束时按 Escape 不关闭对话框', async () => {
    mocks.initialAccount.value = null;
    mocks.save.mockImplementation(() => new Promise(() => undefined));
    renderPage();
    fireEvent.click(screen.getByRole('button', { name: '初始化账户' }));
    fireEvent.click(screen.getByRole('combobox', { name: '初始账号' }));
    fireEvent.click(screen.getByRole('option', { name: '王先生 · 13800000001' }));
    fireEvent.change(screen.getByPlaceholderText('请输入已核对的开通记录编号'), { target: { value: 'OPEN-001' } });
    fireEvent.click(screen.getByRole('checkbox', { name: /已人工核对原始开通记录/ }));
    fireEvent.click(screen.getByRole('button', { name: '确认绑定' }));
    expect(mocks.save).toHaveBeenCalledTimes(1);
    fireEvent.keyDown(document, { key: 'Escape' });
    expect(screen.getByRole('dialog')).toBeTruthy();
  });
  it('已有初始账户时不显示初始化或核验更正入口', () => {
    renderPage();
    expect(screen.queryByRole('button', { name: '核验/更正账户' })).toBeNull();
    expect(screen.queryByRole('button', { name: '初始化账户' })).toBeNull();
    expect(screen.getByRole('button', { name: '重置密码' })).toBeTruthy();
  });
});

describe('加盟商临时密码', () => {
  it('确认后显示脱敏密码，复制时使用完整值并反馈成功', async () => {
    const nativeConfirm = vi.spyOn(window, 'confirm');
    mocks.reset.mockResolvedValue({ data: { resetFranchiseInitialPassword: { accountId: 'account-1', temporaryPassword: 'Ab12Cd34' } } });
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText: mocks.clipboard.mockResolvedValue(undefined) } });
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: '重置密码' }));
    expect(mocks.reset).not.toHaveBeenCalled();
    expect(screen.getByRole('alertdialog').textContent).toContain('所有加盟商工作台');
    fireEvent.click(screen.getByRole('button', { name: '确认重置' }));
    expect(nativeConfirm).not.toHaveBeenCalled();
    expect(mocks.reset).toHaveBeenCalledWith({ variables: { organizationId: 'org-1' } });
    expect(await screen.findByLabelText('脱敏临时密码')).toHaveProperty('textContent', '••••••34');
    expect(screen.queryByText('Ab12Cd34')).toBeNull();
    fireEvent.click(screen.getByRole('button', { name: '复制临时密码' }));
    expect(await screen.findByRole('status')).toHaveProperty('textContent', '复制成功');
    expect(mocks.clipboard).toHaveBeenCalledWith('Ab12Cd34');
  });

  it('重置请求期间阻止重复提交并在失败时提示', async () => {
    let rejectReset: (reason: Error) => void = () => undefined;
    mocks.reset.mockImplementation(() => new Promise((_, reject) => { rejectReset = reject; }));
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: '重置密码' }));
    fireEvent.click(screen.getByRole('button', { name: '确认重置' }));
    expect(screen.getByRole('button', { name: '重置中…' }).hasAttribute('disabled')).toBe(true);
    fireEvent.click(screen.getByRole('button', { name: '重置中…' }));
    expect(mocks.reset).toHaveBeenCalledTimes(1);
    rejectReset(new Error('server unavailable'));
    expect((await screen.findByRole('alert')).textContent).toContain('重置结果暂未确认');
  });
});

describe('临时密码会话有效性', () => {
  it('再次重置结果未知时立刻清除上次密码和复制入口', async () => {
    mocks.reset.mockResolvedValueOnce({ data: { resetFranchiseInitialPassword: { accountId: 'account-1', temporaryPassword: 'Ab12Cd34' } } });
    let rejectReset: (reason: Error) => void = () => undefined;
    mocks.reset.mockImplementationOnce(() => new Promise((_, reject) => { rejectReset = reject; }));
    renderPage();
    fireEvent.click(screen.getByRole('button', { name: '重置密码' }));
    fireEvent.click(screen.getByRole('button', { name: '确认重置' }));
    expect(await screen.findByLabelText('脱敏临时密码')).toHaveProperty('textContent', '••••••34');

    fireEvent.click(screen.getByRole('button', { name: '重置密码' }));
    fireEvent.click(screen.getByRole('button', { name: '确认重置' }));
    await waitFor(() => expect(screen.queryByRole('button', { name: '复制临时密码' })).toBeNull());
    rejectReset(new Error('network disconnected after commit'));
    expect((await screen.findByRole('alert')).textContent).toContain('重置结果暂未确认');
    expect(screen.queryByLabelText('脱敏临时密码')).toBeNull();
  });

  it('重新打开列表后不保留之前的临时密码', async () => {
    mocks.reset.mockResolvedValue({ data: { resetFranchiseInitialPassword: { accountId: 'account-1', temporaryPassword: 'Ab12Cd34' } } });
    renderPage();
    fireEvent.click(screen.getByRole('button', { name: '重置密码' }));
    fireEvent.click(screen.getByRole('button', { name: '确认重置' }));
    expect(await screen.findByLabelText('脱敏临时密码')).toHaveProperty('textContent', '••••••34');

    cleanup();
    renderPage();
    expect(screen.queryByLabelText('脱敏临时密码')).toBeNull();
    expect(screen.queryByRole('button', { name: '复制临时密码' })).toBeNull();
  });

  it('取消重置不会提交请求', () => {
    renderPage();
    fireEvent.click(screen.getByRole('button', { name: '重置密码' }));
    fireEvent.click(screen.getByRole('button', { name: '取消' }));
    expect(mocks.reset).not.toHaveBeenCalled();
  });
});

describe('初始账户请求结果未知', () => {
  it('开通记录编号被占用时提示更正编号', async () => {
    mocks.initialAccount.value = null;
    mocks.save.mockResolvedValue({ ok: false, status: 409, json: async () => ({ code: 'OPENING_RECORD_NUMBER_CONFLICT' }) });
    renderPage();
    fireEvent.click(screen.getByRole('button', { name: '初始化账户' }));
    fireEvent.click(screen.getByRole('combobox', { name: '初始账号' }));
    fireEvent.click(screen.getByRole('option', { name: '王先生 · 13800000001' }));
    fireEvent.change(screen.getByPlaceholderText('请输入已核对的开通记录编号'), { target: { value: 'OPEN-001' } });
    fireEvent.click(screen.getByRole('checkbox', { name: /已人工核对原始开通记录/ }));
    fireEvent.click(screen.getByRole('button', { name: '确认绑定' }));
    expect((await within(screen.getByRole('dialog')).findByRole('alert')).textContent).toContain('开通记录编号已被使用');
  });
  it('网络中断后要求刷新核对并阻止重复提交', async () => {
    mocks.initialAccount.value = null;
    mocks.save.mockRejectedValue(new TypeError('network'));
    renderPage();
    fireEvent.click(screen.getByRole('button', { name: '初始化账户' }));
    fireEvent.click(screen.getByRole('combobox', { name: '初始账号' }));
    fireEvent.click(screen.getByRole('option', { name: '王先生 · 13800000001' }));
    fireEvent.change(screen.getByPlaceholderText('请输入已核对的开通记录编号'), { target: { value: 'OPEN-001' } });
    fireEvent.click(screen.getByRole('checkbox', { name: /已人工核对原始开通记录/ }));
    fireEvent.click(screen.getByRole('button', { name: '确认绑定' }));
    expect((await within(screen.getByRole('dialog')).findByRole('alert')).textContent).toContain('请求结果暂未确认');
    expect(screen.getByRole('button', { name: '确认绑定' }).hasAttribute('disabled')).toBe(true);
    expect(mocks.save).toHaveBeenCalledTimes(1);
  });
});

describe('加盟商开通反馈', () => {
  it('总部账号不能复用时给出可操作的提示', async () => {
    mocks.provision.mockRejectedValue({ extensions: { code: 'PERMISSION_DENIED' } });
    renderPage();
    fireEvent.click(screen.getByRole('button', { name: '开通加盟商' }));
    fireEvent.change(screen.getByPlaceholderText('请输入组织名称'), { target: { value: '测试加盟商' } });
    fireEvent.change(screen.getByPlaceholderText('请输入老板手机号'), { target: { value: '13800000000' } });
    fireEvent.change(screen.getByPlaceholderText('请输入老板姓名'), { target: { value: '王先生' } });
    fireEvent.click(screen.getByRole('button', { name: '确认开通' }));
    expect((await screen.findByRole('alert')).textContent).toContain('该手机号对应的账号无法用于开通加盟商');
  });
});
