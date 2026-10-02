// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { useState } from 'react';
import { afterAll, afterEach, beforeAll, describe, expect, it, vi } from 'vitest';
import { RoleFormDialog, type RoleValues } from './RoleFormDialog';

const permissions = [
  { id: 'payment-read', name: 'permission.storePaymentConfig.read', action: 'storePaymentConfig:read', module: 'storePaymentConfig', scope: 'TENANT' as const },
  { id: 'payment-update', name: 'permission.storePaymentConfig.update', action: 'storePaymentConfig:update', module: 'storePaymentConfig', scope: 'TENANT' as const },
  { id: 'store-read', name: 'permission.store.read', action: 'store:read', module: 'store', scope: 'TENANT' as const },
  { id: 'system-read', name: 'permission.session.read', action: 'session:read', module: 'session', scope: 'SYSTEM' as const },
];

beforeAll(() => vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }));
afterEach(cleanup);
afterAll(() => vi.unstubAllGlobals());

function Harness() {
  const [values, setValues] = useState<RoleValues>({ name: '', permissionIds: [] });
  return <RoleFormDialog open onOpenChange={() => undefined} organizationId="org" permissions={permissions} values={values} onValuesChange={setValues} onSubmit={() => undefined} />;
}

describe('加盟商角色权限弹窗', () => {
  it('使用总部的分组选择器和实体中文标题，并排除系统权限', () => {
    render(<Harness />);
    expect(screen.getByTestId('role-permission-scroll')).toBeTruthy();
    expect(screen.getByText('门店支付配置')).toBeTruthy();
    expect(screen.getByRole('checkbox', { name: '查看门店支付配置' })).toBeTruthy();
    expect(screen.queryByText('storePaymentConfig:read')).toBeNull();
    expect(screen.queryByText('登录会话')).toBeNull();
    expect(screen.getByText('(0/3)')).toBeTruthy();
  });

  it('全选和分组选中只改变加盟商权限', () => {
    render(<Harness />);
    fireEvent.click(screen.getByRole('checkbox', { name: '全选 门店支付配置 全部权限' }));
    expect(screen.getByText('(2/3)')).toBeTruthy();
    fireEvent.click(screen.getByRole('checkbox', { name: '全选全部权限' }));
    expect(screen.getByText('(3/3)')).toBeTruthy();
  });

  it('编辑已有角色时显示编辑标题', () => {
    render(<RoleFormDialog open mode="edit" onOpenChange={() => undefined} organizationId="org" permissions={permissions} onSubmit={() => undefined} />);
    expect(screen.getByRole('heading', { name: '编辑角色' })).toBeTruthy();
  });
});
