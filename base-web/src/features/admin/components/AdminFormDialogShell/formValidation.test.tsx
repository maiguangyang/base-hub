// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { ProvisionFranchiseDialog } from '@/features/hq/pages/FranchiseListPage/ProvisionFranchiseDialog';
import { AdministratorFormDialog } from '@/features/hq/pages/AdministratorListPage/AdministratorFormDialog';
import { InviteStaffDialog } from '@/features/franchise/pages/StaffListPage/InviteStaffDialog';
import { StoreBasicFields, StoreFacilityFields, StoreLocationFields, StoreOperationFields } from '@/features/admin/pages/StoreEditorPage/StoreFormSections';
import { RoleFormDialog } from '@/features/franchise/pages/RoleListPage/RoleFormDialog';
import { RoleFormDialog as HqRoleFormDialog } from '@/features/hq/pages/RoleListPage/RoleFormDialog';
import { ReviewStoreDialog } from '@/features/hq/pages/StoreApprovalListPage/ReviewStoreDialog';
import { AdminValidatedField } from '@/features/admin/components/AdminValidatedField';

beforeEach(() => { vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} }); });
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

describe('管理表单输入校验', () => {
  const cases = [
    ['开通加盟商', <ProvisionFranchiseDialog open onOpenChange={() => undefined} onSubmit={() => undefined} values={{ name: '加盟商', ownerPhone: '138000000003', ownerDisplayName: '老板', ownerEmail: 'bad' }} />],
    ['新增管理员', <AdministratorFormDialog open mode="create" onOpenChange={() => undefined} onSubmit={() => undefined} roles={[]} values={{ phone: '138000000003', displayName: '管理员', email: 'bad', roleIds: [] }} />],
    ['新增员工', <InviteStaffDialog open onOpenChange={() => undefined} onSubmit={() => undefined} values={{ phone: '138000000003', displayName: '员工', email: 'bad', roleIds: [], storeAccessMode: 'ALL_STORES', storeIds: [] }} />],
  ] as const;

  it.each(cases)('%s 限制手机号码并提示无效号码和邮箱', (_, dialog) => {
    render(dialog);
    const phone = screen.getByLabelText(/手机号/) as HTMLInputElement;
    const email = screen.getByLabelText(/邮箱/) as HTMLInputElement;
    expect(phone.maxLength).toBe(11);
    expect(phone.inputMode).toBe('numeric');
    expect(phone.checkValidity()).toBe(false);
    fireEvent.blur(phone);
    expect(screen.getByText('请输入11位有效手机号')).toBeTruthy();
    expect(email.type).toBe('email');
    expect(email.maxLength).toBe(128);
    fireEvent.blur(email);
    expect(screen.getByText('请输入有效的邮箱地址')).toBeTruthy();
  });

  it('直营店的店长手机号也有输入长度和格式限制', () => {
    render(<StoreBasicFields values={{ name: '店', managerPhone: '138000000003' }} change={() => undefined} />);
    const phone = screen.getByLabelText('店长手机号') as HTMLInputElement;
    expect(phone.maxLength).toBe(11);
    expect(phone.checkValidity()).toBe(false);
    expect((screen.getByLabelText(/联系电话/) as HTMLInputElement).maxLength).toBe(32);
    expect((screen.getByLabelText(/门店名称/) as HTMLInputElement).maxLength).toBe(128);
  });

  it('加盟门店和角色名称遵守服务端字段长度', () => {
    render(<StoreBasicFields values={{ name: '店' }} change={() => undefined} />);
    expect((screen.getByLabelText(/门店名称/) as HTMLInputElement).maxLength).toBe(128);
    cleanup();
    render(<RoleFormDialog open onOpenChange={() => undefined} organizationId="org" permissions={[]} onSubmit={() => undefined} />);
    expect((screen.getByLabelText(/角色名称/) as HTMLInputElement).maxLength).toBe(64);
  });

  it('总部角色名称和门店退回原因遵守服务端字段长度', () => {
    render(<HqRoleFormDialog open onOpenChange={() => undefined} onSubmit={() => undefined} permissions={[]} viewerPermissions={[]} />);
    expect((screen.getByLabelText(/角色名称/) as HTMLInputElement).maxLength).toBe(64);
    cleanup();
    render(<ReviewStoreDialog open onOpenChange={() => undefined} onSubmit={() => undefined} storeId="store" approved={false} />);
    expect((screen.getByLabelText('退回原因') as HTMLInputElement).maxLength).toBe(512);
  });

  it('邮箱在提交前识别服务端不接受的格式', () => {
    render(<AdminValidatedField label="邮箱" placeholder="请输入邮箱" value="Owner@example.com" onChange={() => undefined} kind="email" />);
    const email = screen.getByLabelText('邮箱') as HTMLInputElement;
    expect(email.checkValidity()).toBe(false);
    fireEvent.blur(email);
    expect(screen.getByText('请输入有效的邮箱地址')).toBeTruthy();
  });
});

describe('加盟门店申请资料', () => {
  it('提供与直营店相同的资料字段', () => {
    const props = { values: { name: '店' }, change: () => undefined };
    render(<div data-testid="store-fields"><StoreBasicFields {...props} /><StoreLocationFields {...props} /><StoreOperationFields {...props} /><StoreFacilityFields {...props} /></div>);

    for (const label of ['门店名称', '联系电话', '营业时间', '店长姓名', '店长手机号', '省份', '城市', '区县', '详细地址', '营业状态', '服务模式', '经营面积（㎡）', '桌位数（桌）', '小票底部寄语']) {
      expect(screen.getByTestId('store-fields').textContent).toContain(label);
    }
    expect((screen.getByLabelText(/联系电话/) as HTMLInputElement).required).toBe(true);
    expect((screen.getByLabelText(/详细地址/) as HTMLInputElement).required).toBe(true);
  });
});
