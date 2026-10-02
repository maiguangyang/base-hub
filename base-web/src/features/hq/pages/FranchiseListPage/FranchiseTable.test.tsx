import { renderToStaticMarkup } from 'react-dom/server';
import { MemoryRouter } from 'react-router';
import { describe, expect, it } from 'vitest';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';
import { FranchiseTable, franchiseStorePath, maskTemporaryPassword, type FranchiseRow } from './FranchiseTable';

const row: FranchiseRow = {
  id: 'org-1',
  code: 'FR-001',
  name: '首尔加盟商',
  status: 'ACTIVE',
  initialAccountId: 'account-1',
  initialAccount: { phone: '13800000001' },
};

function renderTable(canViewStores: boolean, canResetPassword = false): string {
  return renderToStaticMarkup(
    <MemoryRouter>
      <FranchiseTable
        rows={[row]}
        canViewStores={canViewStores}
        canSuspend={false}
        canRestore={false}
        canResetPassword={canResetPassword}
        canConfirmInitialAccount={false}
        resettingOrganizationId={undefined}
        temporaryPasswords={{}}
        onSuspend={() => undefined}
        onRestore={() => undefined}
        onResetPassword={() => undefined}
        onConfirmInitialAccount={() => undefined}
      />
    </MemoryRouter>,
  );
}

describe('总部加盟商列表', () => {
  it('显示登录账户，并把状态开关放在操作列前', () => {
    const markup = renderTable(false, true);
    expect(markup).toContain('13800000001');
    expect(markup.indexOf('登录账户')).toBeLessThan(markup.indexOf('状态'));
    expect(markup.indexOf('状态')).toBeLessThan(markup.indexOf('操作'));
    expect(markup).toContain('role="switch"');
    expect(markup).toContain('aria-checked="true"');
    expect(markup).not.toContain('>启用</span>');
  });

  it('无初始账户读取权限时不暴露登录手机号', () => {
    const markup = renderTable(false);
    expect(markup).toContain('无权限查看');
    expect(markup).not.toContain('13800000001');
  });
});

describe('加盟商列表原有操作', () => {
  it('对深链中的加盟商标识进行 URL 编码', () => {
    expect(franchiseStorePath('org/a?b')).toBe('/admin/hq/stores?organizationId=org%2Fa%3Fb&type=FRANCHISE');
  });

  it('有跨组织门店读取权限时显示加盟商门店深链', () => {
    const markup = renderTable(true);

    expect(markup).toContain('查看门店');
    expect(markup).toContain('/admin/hq/stores?organizationId=org-1&amp;type=FRANCHISE');
  });

  it('缺少跨组织门店读取权限时隐藏门店入口', () => {
    expect(renderTable(false)).not.toContain('查看门店');
  });

  it('临时密码只露出末两位', () => {
    expect(maskTemporaryPassword('Ab12Cd34')).toBe('••••••34');
  });

  it('有权限且初始账号存在时显示重置入口和脱敏密码', () => {
    const markup = renderToStaticMarkup(<AdminToastProvider><MemoryRouter><FranchiseTable rows={[row]} canViewStores={false} canSuspend={false} canRestore={false} canResetPassword canConfirmInitialAccount resettingOrganizationId={undefined} temporaryPasswords={{ 'org-1': { accountId: 'account-1', password: 'Ab12Cd34' } }} onSuspend={() => undefined} onRestore={() => undefined} onResetPassword={() => undefined} onConfirmInitialAccount={() => undefined} /></MemoryRouter></AdminToastProvider>);
    expect(markup).toContain('重置密码');
    expect(markup).toContain('••••••34');
    expect(markup).toContain('复制临时密码');
    expect(markup).not.toContain('Ab12Cd34');
  });

  it('历史加盟商没有初始账号和临时密码时仍可渲染列表', () => {
    const unresolvedRow: FranchiseRow = { ...row, initialAccountId: null, initialAccount: null };
    const markup = renderToStaticMarkup(<MemoryRouter><FranchiseTable rows={[unresolvedRow]} canViewStores={false} canSuspend={false} canRestore={false} canResetPassword canConfirmInitialAccount resettingOrganizationId={undefined} temporaryPasswords={{}} onSuspend={() => undefined} onRestore={() => undefined} onResetPassword={() => undefined} onConfirmInitialAccount={() => undefined} /></MemoryRouter>);
    expect(markup).toContain('账户未初始化，暂不可重置密码');
    expect(markup).toContain('初始化账户');
    expect(markup).toContain('未初始化');
    expect(markup).not.toContain('需先核对该加盟商的初始账号');
    expect(markup).not.toContain('复制临时密码');
  });

  it('初始账号已删除但标记仍在时，不提供再次初始化', () => {
    const deletedAccountRow: FranchiseRow = { ...row, initialAccount: null };
    const markup = renderToStaticMarkup(<MemoryRouter><FranchiseTable rows={[deletedAccountRow]} canViewStores={false} canSuspend={false} canRestore={false} canResetPassword canConfirmInitialAccount resettingOrganizationId={undefined} temporaryPasswords={{}} onSuspend={() => undefined} onRestore={() => undefined} onResetPassword={() => undefined} onConfirmInitialAccount={() => undefined} /></MemoryRouter>);
    expect(markup).not.toContain('初始化账户');
    expect(markup).toContain('重置密码');
    expect(markup).toContain('账号已失效');
  });
});
