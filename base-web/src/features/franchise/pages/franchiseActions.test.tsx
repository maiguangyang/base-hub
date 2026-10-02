import { describe, expect, it } from 'vitest';
import { AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { storeCreateInput } from '@/features/admin/pages/StoreEditorPage/storeInputs';
import { canEditStoreRecord, canSubmitFranchiseStore } from '@/features/admin/pages/StoreEditorPage/storeEditorPolicy';
import { getGraphQLErrorCode } from '@/lib/graphql/errors';
import { InviteStaffDialog, staffAccessInput, staffInviteOutcome } from './StaffListPage/InviteStaffDialog';
import { StaffAccessDialog, staffAccessValues } from './StaffListPage/StaffAccessDialog';
import { staffMutationErrorMessage } from './StaffListPage/StaffListPage';
import { RoleFormDialog, canManageRole, groupPermissionActions, tenantPermissionsOnly } from './RoleListPage/RoleFormDialog';
import { canDeleteStore } from './StoreListPage/StoreTable';

describe('加盟商表单与输入', () => {
  it('员工和角色表单统一使用 AdminFormDialogShell', () => {
    expect(InviteStaffDialog({ open: true, onOpenChange: () => undefined, onSubmit: async () => undefined }).type).toBe(AdminFormDialogShell);
    expect(StaffAccessDialog({ open: true, onOpenChange: () => undefined, onSubmit: async () => undefined }).type).toBe(AdminFormDialogShell);
    expect(RoleFormDialog({ open: true, onOpenChange: () => undefined, organizationId: 'org', permissions: [], onSubmit: async () => undefined }).type).toBe(AdminFormDialogShell);
  });

  it('已有员工可以回填多角色和门店范围', () => {
    expect(staffAccessValues({
      storeAccessMode: 'SELECTED_STORES',
      roles: [{ id: 'role-1' }, { id: 'role-2' }],
      stores: [{ id: 'store-1' }],
    })).toEqual({
      roleIds: ['role-1', 'role-2'],
      storeAccessMode: 'SELECTED_STORES',
      storeIds: ['store-1'],
    });
  });

  it('门店表单内部固定 DRAFT，支持内部编码自动生成，只有草稿或退回状态可以提交', () => {
    expect(storeCreateInput({ code: 'S1', name: 'One' }, 'org', 'FRANCHISE')).toMatchObject({ code: 'S1', name: 'One', organizationId: 'org', lifecycle: 'DRAFT' });
    const auto = storeCreateInput({ name: 'AutoFranchiseStore' }, 'org', 'FRANCHISE');
    expect(auto.name).toBe('AutoFranchiseStore');
    expect(auto.code.startsWith('STR')).toBe(true);
    expect(auto.organizationId).toBe('org');
    expect(auto.lifecycle).toBe('DRAFT');
    expect(canSubmitFranchiseStore('DRAFT', ['store:submit'])).toBe(true);
    expect(canSubmitFranchiseStore('PENDING_APPROVAL', ['store:submit'])).toBe(false);
    expect(canEditStoreRecord('FRANCHISE', 'org', ['store:update'], { organizationId: 'org', lifecycle: 'DRAFT' })).toBe(true);
    expect(canEditStoreRecord('FRANCHISE', 'org', ['store:update'], { organizationId: 'org', lifecycle: 'REJECTED' })).toBe(true);
    expect(canEditStoreRecord('FRANCHISE', 'org', ['store:update'], { organizationId: 'org', lifecycle: 'PENDING_APPROVAL' })).toBe(false);
    expect(canEditStoreRecord('FRANCHISE', 'org', ['store:update'], { organizationId: 'org', lifecycle: 'ACTIVE' })).toBe(false);
	expect(canDeleteStore('DRAFT', ['store:delete'])).toBe(true);
	expect(canDeleteStore('REJECTED', ['store:delete'])).toBe(true);
	expect(canDeleteStore('PENDING_APPROVAL', ['store:delete'])).toBe(false);
	expect(canDeleteStore('ACTIVE', ['store:delete'])).toBe(false);
  });
});

describe('加盟门店申请资料', () => {
  it('草稿保存直营店同款资料并保持草稿状态', () => {
    const input = storeCreateInput({
      name: '  江南西加盟店  ', contactPhone: ' 020-12345678 ', managerName: ' 店长 ', managerPhone: '13800000000',
      province: '广东省', city: '广州市', district: '海珠区', address: '江南大道1号', businessHours: '10:00 - 21:00',
      businessStatus: 'CLOSED', supportDineIn: false, supportTakeout: true, storeArea: '86.5', tableCount: '12', receiptFooter: ' 欢迎光临 ',
    }, 'franchise-org', 'FRANCHISE');

    expect(input).toMatchObject({
      name: '江南西加盟店', organizationId: 'franchise-org', lifecycle: 'DRAFT', contactPhone: '020-12345678',
      managerName: '店长', managerPhone: '13800000000', province: '广东省', city: '广州市', district: '海珠区',
      address: '江南大道1号', businessHours: '10:00 - 21:00', businessStatus: 'CLOSED',
      supportDineIn: false, supportTakeout: true, storeArea: 86.5, tableCount: 12, receiptFooter: '欢迎光临',
    });
  });
});

describe('加盟商权限与角色', () => {
  it('全部门店清空门店选择，指定门店至少需要一个门店', () => {
    expect(staffAccessInput('ALL_STORES', ['store-1'])).toEqual({ storeAccessMode: 'ALL_STORES', storeIds: [] });
    expect(() => staffAccessInput('SELECTED_STORES', [])).toThrow('VALIDATION_FAILED');
    expect(staffAccessInput('SELECTED_STORES', ['store-1'])).toEqual({ storeAccessMode: 'SELECTED_STORES', storeIds: ['store-1'] });
  });

  it('角色编辑只接收 TENANT 权限并按精确 CRUD 动作分组', () => {
    const permissions = [
      { id: '1', action: 'store:read', scope: 'TENANT' as const, module: 'store', name: 'read' },
      { id: '2', action: 'account:read', scope: 'SYSTEM' as const, module: 'account', name: 'read' },
    ];
    expect(tenantPermissionsOnly(permissions).map((item) => item.id)).toEqual(['1']);
    expect(groupPermissionActions(tenantPermissionsOnly(permissions))).toEqual({ store: { read: '1' } });
    expect(canManageRole('FRANCHISE_OWNER')).toBe(false);
    expect(canManageRole('CUSTOM')).toBe(true);
  });
});

describe('加盟商邀请反馈', () => {
  it('员工一次性密码与既有账号邀请状态互斥，最后老板错误只看 code', () => {
    expect(staffInviteOutcome({ temporaryPassword: 'Once-Secret-42!', invitationPending: false }).kind).toBe('temporary-password');
    expect(staffInviteOutcome({ temporaryPassword: null, invitationPending: true })).toEqual({ kind: 'invitation-pending' });
    expect(getGraphQLErrorCode({ message: '任意文本', extensions: { code: 'LAST_OWNER_REQUIRED' } })).toBe('LAST_OWNER_REQUIRED');
    expect(staffMutationErrorMessage('LAST_OWNER_REQUIRED')).toBe('每个加盟商必须保留至少一位有效所有者。');
    expect(staffMutationErrorMessage(undefined)).toBe('操作失败，请稍后重试。');
  });

  it('邀请结果页不允许重复提交', () => {
    const dialog = InviteStaffDialog({
      open: true, onOpenChange: () => undefined, onSubmit: async () => undefined,
      outcome: { kind: 'temporary-password', value: 'Once-Secret-42!' },
    });
    expect(dialog.props.hideSubmit).toBe(true);
  });
});
