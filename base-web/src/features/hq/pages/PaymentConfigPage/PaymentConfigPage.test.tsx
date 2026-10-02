// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import type { ScopeView } from '@/features/hq/api/paymentConfig';
import { PaymentConfigPage } from './PaymentConfigPage';

const mocks = vi.hoisted(() => ({
  workspace: 'HEADQUARTERS', permissions: ['paymentConfig:read', 'paymentConfig:manage'],
  page: {} as Record<string, unknown>,
}));

vi.mock('@/features/auth/store/authStore', () => ({ useAuthStore: (select: (state: unknown) => unknown) => select({ viewer: {
  permissions: mocks.permissions, currentWorkspace: { workspaceType: mocks.workspace },
} }) }));
vi.mock('./usePaymentConfigPage', () => ({ usePaymentConfigPage: () => mocks.page }));

const view: ScopeView = { channels: [
  { own: { scope: 'STORE', channel: 'WECHAT', recordId: 'one', version: 2, merchantMasked: '******7890', ratePpm: 3800, state: 'ERROR', reasonCode: 'DECRYPTION_FAILED', credentialsConfigured: true },
    effective: { scope: 'STORE', sourceScope: 'STORE', channel: 'WECHAT', recordId: 'one', version: 2, merchantMasked: '******7890', ratePpm: 3800, state: 'ERROR', reasonCode: 'DECRYPTION_FAILED', credentialsConfigured: true } },
  { own: { scope: 'STORE', channel: 'ALIPAY', version: 0, ratePpm: 0, state: 'UNCONFIGURED', credentialsConfigured: false },
    effective: { scope: 'GLOBAL', sourceScope: 'GLOBAL', channel: 'ALIPAY', recordId: 'global', version: 1, merchantMasked: '********1234', ratePpm: 2500, state: 'VALID', credentialsConfigured: true } },
] };

beforeEach(() => {
  mocks.workspace = 'HEADQUARTERS'; mocks.permissions = ['paymentConfig:read', 'paymentConfig:manage'];
  mocks.page = { scope: 'STORE', setScope: vi.fn(), organizationId: 'org', setOrganizationId: vi.fn(), storeId: 'shop', setStoreId: vi.fn(),
    franchises: [{ id: 'org', name: '加盟商' }], stores: [{ id: 'shop', name: '门店', organizationId: 'org' }], catalogError: undefined,
    loading: false, pending: false, view, error: undefined, notice: undefined, conflict: false,
    reload: vi.fn(), save: vi.fn().mockResolvedValue(true), setState: vi.fn(), restoreInheritance: vi.fn() };
});
afterEach(() => cleanup());

describe('总部支付配置表单背景', () => {
  it('输入框、多行文本框、下拉触发器和选项面板均为白底深色文字', () => {
    const { container } = render(<PaymentConfigPage />);
    const page = container.firstElementChild as HTMLElement;
    expect(page.className).not.toContain('[&_input]:!bg-white');
    expect(Array.from(container.querySelectorAll('input, textarea, [data-slot="select-trigger"]'))
      .every((element) => element.hasAttribute('data-admin-form-surface'))).toBe(true);

    fireEvent.click(screen.getByRole('combobox', { name: '微信验签模式' }));
    const menu = document.querySelector('[data-slot="select-content"]') as HTMLElement;
    expect(menu.hasAttribute('data-admin-form-surface')).toBe(true);
  });
});

describe('总部支付配置下拉框', () => {
  it('全局配置页不再提供加盟商或门店范围选择', () => {
    render(<PaymentConfigPage />);
    expect(document.querySelectorAll('select')).toHaveLength(0);
    expect(screen.getAllByRole('combobox')).toHaveLength(1);
    expect(screen.queryByRole('combobox', { name: '配置级别' })).toBeNull();
    expect(screen.queryByRole('combobox', { name: '加盟商' })).toBeNull();
    expect(screen.queryByRole('combobox', { name: '加盟门店' })).toBeNull();
  });

  it('微信验签模式切换后显示支付公钥字段', () => {
    render(<PaymentConfigPage />);
    fireEvent.click(screen.getByRole('combobox', { name: '微信验签模式' }));
    fireEvent.click(screen.getByRole('option', { name: '微信支付公钥' }));
    expect(screen.getByLabelText('微信支付公钥 ID')).toBeTruthy();
  });
});

describe('总部支付配置渠道标签', () => {
  it('微信和支付宝以标签切换，每次只显示当前渠道表单', () => {
    render(<PaymentConfigPage />);
    expect(screen.getByRole('tab', { name: '微信支付' })).toBeTruthy();
    expect(screen.getByRole('tab', { name: '支付宝支付' })).toBeTruthy();
    fireEvent.change(screen.getByRole('textbox', { name: '微信商户私钥' }), { target: { value: 'UNSAVED KEY' } });
    expect(screen.queryByRole('textbox', { name: '支付宝应用私钥' })).toBeNull();
    fireEvent.mouseDown(screen.getByRole('tab', { name: '支付宝支付' }), { button: 0 });
    expect(screen.getByRole('textbox', { name: '支付宝应用私钥' })).toBeTruthy();
    expect(screen.queryByRole('textbox', { name: '微信商户私钥' })).toBeNull();
    fireEvent.mouseDown(screen.getByRole('tab', { name: '微信支付' }), { button: 0 });
    expect((screen.getByRole('textbox', { name: '微信商户私钥' }) as HTMLTextAreaElement).value).toBe('UNSAVED KEY');
  });
});

describe('总部支付配置页面', () => {
  it('加盟工作区无权打开，只有读权限时不能编辑', () => {
    mocks.workspace = 'FRANCHISE'; render(<PaymentConfigPage />);
    expect(screen.getByRole('alert').textContent).toContain('无权查看');
    cleanup(); mocks.workspace = 'HEADQUARTERS'; mocks.permissions = ['paymentConfig:read']; render(<PaymentConfigPage />);
    expect(screen.queryByRole('button', { name: '保存配置' })).toBeNull();
  });

  it('错误覆盖不回退，支付宝独立继承，秘密字段保持空白', () => {
    render(<PaymentConfigPage />);
    expect(screen.getByText(/解密失败/)).toBeTruthy();
    fireEvent.mouseDown(screen.getByRole('tab', { name: '支付宝支付' }), { button: 0 });
    expect(screen.getByText(/支付宝.*全局/)).toBeTruthy();
    expect((screen.getByLabelText('支付宝应用私钥') as HTMLTextAreaElement).value).toBe('');
    expect(screen.getByLabelText('支付宝应用私钥').className).toContain('font-mono');
    expect(screen.queryByText('可收款')).toBeNull();
    expect(screen.queryByText('支付已开通')).toBeNull();
  });

  it('平台证书模式不要求手动上传平台证书', () => {
    render(<PaymentConfigPage />);
    expect(screen.queryByLabelText('微信平台证书')).toBeNull();
    expect(screen.getByLabelText('微信商户 API 证书')).toBeTruthy();
  });

  it('费率超出四位小数时阻止保存', () => {
    render(<PaymentConfigPage />);
    fireEvent.change(screen.getByLabelText('微信预计平台费率（%）'), { target: { value: '0.12345' } });
    fireEvent.click(screen.getByRole('button', { name: '保存配置' }));
    expect((mocks.page.save as ReturnType<typeof vi.fn>)).not.toHaveBeenCalled();
  });

  it('精确提交 0.38% 且保存后清空写入凭证', async () => {
    render(<PaymentConfigPage />);
    fireEvent.change(screen.getByLabelText('微信预计平台费率（%）'), { target: { value: '0.38' } });
    fireEvent.change(screen.getByLabelText('微信商户私钥'), { target: { value: 'PRIVATE KEY' } });
    fireEvent.click(screen.getByRole('button', { name: '保存配置' }));
    await waitFor(() => expect(mocks.page.save).toHaveBeenCalledWith('WECHAT', expect.objectContaining({ ratePpm: 3800 })));
    await waitFor(() => expect((screen.getByLabelText('微信商户私钥') as HTMLTextAreaElement).value).toBe(''));
  });

  it('版本冲突提供重新加载操作', () => {
    mocks.page = { ...mocks.page, error: '配置已被其他管理员修改', conflict: true };
    render(<PaymentConfigPage />);
    fireEvent.click(screen.getByRole('button', { name: '重新加载最新数据' }));
    expect(mocks.page.reload).toHaveBeenCalled();
  });
});

describe('总部支付配置编辑状态', () => {
  it('重新加载后同步最新费率并清空未提交的凭证', () => {
    const page = render(<PaymentConfigPage />);
    fireEvent.change(screen.getByLabelText('微信预计平台费率（%）'), { target: { value: '0.99' } });
    fireEvent.change(screen.getByLabelText('微信商户私钥'), { target: { value: 'STALE KEY' } });
    mocks.page = { ...mocks.page, view: { channels: [
      { ...view.channels[0], own: { ...view.channels[0].own, version: 3, ratePpm: 4200 } },
      view.channels[1],
    ] } };
    page.rerender(<PaymentConfigPage />);
    expect((screen.getByLabelText('微信预计平台费率（%）') as HTMLInputElement).value).toBe('0.42');
    expect((screen.getByLabelText('微信商户私钥') as HTMLTextAreaElement).value).toBe('');
  });

  it('更换支付宝凭证时必须显式选择环境', async () => {
    render(<PaymentConfigPage />);
    fireEvent.mouseDown(screen.getByRole('tab', { name: '支付宝支付' }), { button: 0 });
    fireEvent.change(screen.getByLabelText('支付宝商户标识'), { target: { value: '1234567890123456' } });
    fireEvent.change(screen.getByLabelText('支付宝应用私钥'), { target: { value: 'PRIVATE KEY' } });
    fireEvent.click(screen.getByRole('button', { name: '保存配置' }));
    expect(mocks.page.save).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole('combobox', { name: '支付宝环境' }));
    expect(screen.queryByRole('option', { name: '更换凭证时请选择环境' })).toBeNull();
    expect(screen.getAllByRole('option').map((option) => option.textContent)).toEqual(['正式环境', '沙箱环境']);
    fireEvent.click(screen.getByRole('option', { name: '沙箱环境' }));
    fireEvent.click(screen.getByRole('button', { name: '保存配置' }));
    expect(mocks.page.save).not.toHaveBeenCalled();
    expect(screen.getByText('请同时填写支付宝应用私钥和支付宝公钥。')).toBeTruthy();
    fireEvent.change(screen.getByLabelText('支付宝公钥'), { target: { value: 'PUBLIC KEY' } });
    fireEvent.click(screen.getByRole('button', { name: '保存配置' }));
    await waitFor(() => expect(mocks.page.save).toHaveBeenCalledWith('ALIPAY', expect.objectContaining({ environment: 'SANDBOX' })));
  });

  it('支付宝普通公钥模式只填写应用私钥与支付宝公钥', () => {
    render(<PaymentConfigPage />);
    fireEvent.mouseDown(screen.getByRole('tab', { name: '支付宝支付' }), { button: 0 });
    expect(screen.getByLabelText('支付宝应用私钥')).toBeTruthy();
    expect(screen.getByLabelText('支付宝公钥')).toBeTruthy();
    expect(screen.queryByLabelText('支付宝应用公钥证书')).toBeNull();
    expect(screen.queryByLabelText('支付宝根证书')).toBeNull();
    expect(screen.queryByLabelText('支付宝公钥证书')).toBeNull();
  });
});

describe('支付配置统一保存与启停', () => {
  it('只有一个保存按钮，按渠道提交两个有修改的表单', async () => {
    render(<PaymentConfigPage />);
    expect(screen.getAllByRole('button', { name: '保存配置' })).toHaveLength(1);
    fireEvent.change(screen.getByLabelText('微信商户私钥'), { target: { value: 'WECHAT KEY' } });
    fireEvent.mouseDown(screen.getByRole('tab', { name: '支付宝支付' }), { button: 0 });
    fireEvent.change(screen.getByLabelText('支付宝商户标识'), { target: { value: '1234567890123456' } });
    fireEvent.change(screen.getByLabelText('支付宝应用私钥'), { target: { value: 'ALIPAY KEY' } });
    fireEvent.change(screen.getByLabelText('支付宝公钥'), { target: { value: 'ALIPAY PUBLIC KEY' } });
    fireEvent.click(screen.getByRole('combobox', { name: '支付宝环境' }));
    fireEvent.click(screen.getByRole('option', { name: '沙箱环境' }));
    fireEvent.click(screen.getByRole('button', { name: '保存配置' }));
    await waitFor(() => expect(mocks.page.save).toHaveBeenCalledTimes(2));
    expect((mocks.page.save as ReturnType<typeof vi.fn>).mock.calls.map((call) => call[0])).toEqual(['WECHAT', 'ALIPAY']);
  });

  it('第二个渠道保存失败时提示部分成功并保留失败渠道草稿', async () => {
    (mocks.page.save as ReturnType<typeof vi.fn>).mockResolvedValueOnce(true).mockResolvedValueOnce(false);
    render(<PaymentConfigPage />);
    fireEvent.change(screen.getByLabelText('微信商户私钥'), { target: { value: 'WECHAT KEY' } });
    fireEvent.mouseDown(screen.getByRole('tab', { name: '支付宝支付' }), { button: 0 });
    fireEvent.change(screen.getByLabelText('支付宝商户标识'), { target: { value: '1234567890123456' } });
    fireEvent.change(screen.getByLabelText('支付宝应用私钥'), { target: { value: 'ALIPAY KEY' } });
    fireEvent.change(screen.getByLabelText('支付宝公钥'), { target: { value: 'ALIPAY PUBLIC KEY' } });
    fireEvent.click(screen.getByRole('combobox', { name: '支付宝环境' }));
    fireEvent.click(screen.getByRole('option', { name: '沙箱环境' }));
    fireEvent.click(screen.getByRole('button', { name: '保存配置' }));
    await waitFor(() => expect(screen.getByText(/已保存 1 个支付方式/)).toBeTruthy());
    expect((screen.getByLabelText('支付宝应用私钥') as HTMLTextAreaElement).value).toBe('ALIPAY KEY');
    fireEvent.mouseDown(screen.getByRole('tab', { name: '微信支付' }), { button: 0 });
    expect((screen.getByLabelText('微信商户私钥') as HTMLTextAreaElement).value).toBe('');
  });
});

describe('支付配置开关', () => {
  it('当前渠道有未保存修改时暂不允许切换启停', () => {
    render(<PaymentConfigPage />);
    fireEvent.mouseDown(screen.getByRole('tab', { name: '支付宝支付' }), { button: 0 });
    fireEvent.change(screen.getByLabelText('支付宝商户标识'), { target: { value: '1234567890123456' } });
    expect(screen.getByRole('switch', { name: '支付宝支付启用状态' }).hasAttribute('disabled')).toBe(true);
  });

  it('继承的有效配置可通过独立一行的开关停用', () => {
    render(<PaymentConfigPage />);
    fireEvent.mouseDown(screen.getByRole('tab', { name: '支付宝支付' }), { button: 0 });
    const toggle = screen.getByRole('switch', { name: '支付宝支付启用状态' });
    expect(toggle.getAttribute('aria-checked')).toBe('true');
    expect(toggle.parentElement?.contains(screen.getByRole('button', { name: '保存配置' }))).toBe(false);
    fireEvent.click(toggle);
    expect(mocks.page.setState).toHaveBeenCalledWith('ALIPAY', 'DISABLED');
  });

  it('有凭证的停用配置可启用，未配置时无法凭空启用', () => {
    mocks.page = { ...mocks.page, view: { channels: [
      { ...view.channels[0], own: { ...view.channels[0].own, state: 'DISABLED' }, effective: { ...view.channels[0].effective, state: 'DISABLED' } },
      { ...view.channels[1], effective: { ...view.channels[1].effective, state: 'UNCONFIGURED' } },
    ] } };
    render(<PaymentConfigPage />);
    const wechatToggle = screen.getByRole('switch', { name: '微信支付启用状态' });
    fireEvent.click(wechatToggle);
    expect(mocks.page.setState).toHaveBeenCalledWith('WECHAT', 'VALID');
    fireEvent.mouseDown(screen.getByRole('tab', { name: '支付宝支付' }), { button: 0 });
    expect(screen.getByRole('switch', { name: '支付宝支付启用状态' }).hasAttribute('disabled')).toBe(true);
  });
});
