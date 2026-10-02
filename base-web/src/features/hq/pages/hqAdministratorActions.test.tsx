import { describe, expect, it } from 'vitest';
import { renderToStaticMarkup } from 'react-dom/server';
import { AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { AdminRoleSelect } from '@/features/admin/components/AdminRoleSelect/AdminRoleSelect';
import {
  AdministratorFormDialog,
  administratorInviteInput,
  administratorInviteOutcome,
  assignableAdministratorRoles,
} from './AdministratorListPage/AdministratorFormDialog';
import { AdministratorTable, administratorActions, formatAdministratorUpdatedAt, type AdministratorRow } from './AdministratorListPage/AdministratorTable';
import { administratorMutationErrorMessage, administratorViewState } from './AdministratorListPage/useAdministratorPage';
import {
  RoleFormDialog,
  canManageRole,
  delegableSystemPermissions,
  groupPermissionActions,
  systemPermissionsOnly,
} from './RoleListPage/RoleFormDialog';
import { buildRoleMutationInput, permissionCatalogReady, roleViewState } from './RoleListPage/useRolePage';

const permissions = [
  { id: 'read', name: 'read', action: 'hqMembership:read', module: 'membership', scope: 'SYSTEM' as const },
  { id: 'update', name: 'update', action: 'hqMembership:update', module: 'membership', scope: 'SYSTEM' as const },
  { id: 'tenant', name: 'tenant', action: 'store:read', module: 'store', scope: 'TENANT' as const },
];

describe('总部管理员表单与动作', () => {
  it('新增和编辑统一使用表单壳，编辑模式不允许修改账号资料', () => {
    const create = AdministratorFormDialog({ open: true, mode: 'create', onOpenChange: () => undefined, onSubmit: async () => undefined, roles: [] });
    const edit = AdministratorFormDialog({ open: true, mode: 'edit', onOpenChange: () => undefined, onSubmit: async () => undefined, roles: [], account: { displayName: '管理员', phone: '13800000000', email: null } });
    expect(create.type).toBe(AdminFormDialogShell);
    expect(edit.type).toBe(AdminFormDialogShell);
    expect(JSON.stringify(create)).toContain('请输入手机号');
    expect(JSON.stringify(create)).toContain('请输入姓名');
    expect(JSON.stringify(create)).not.toMatch(/编号|编码/);
    expect(JSON.stringify(edit)).not.toContain('请输入手机号');
    expect(JSON.stringify(edit)).not.toContain('请输入姓名');
  });

  it('总部管理员输入固定全部门店，邀请结果互斥', () => {
    expect(administratorInviteInput({ phone: '13800000000', displayName: '管理员', email: '', roleIds: ['role'] })).toEqual({
      phone: '13800000000', displayName: '管理员', email: null, roleIds: ['role'], storeAccessMode: 'ALL_STORES', storeIds: [],
    });
    expect(administratorInviteOutcome({ temporaryPassword: 'Once-123', invitationPending: false })).toEqual({ kind: 'temporary-password', value: 'Once-123' });
    expect(administratorInviteOutcome({ temporaryPassword: null, invitationPending: true })).toEqual({ kind: 'invitation-pending' });
  });

  it('动作权限相互独立且当前账号不可管理', () => {
    expect(administratorActions('target', 'self', 'ACTIVE', ['hqMembership:update'])).toEqual({ edit: true, status: true, resetPassword: false, delete: false });
    expect(administratorActions('target', 'self', 'ACTIVE', ['account:update', 'hqMembership:delete'])).toEqual({ edit: false, status: false, resetPassword: true, delete: true });
    expect(administratorActions('target', 'self', 'SUSPENDED', ['account:update'])).toEqual({ edit: false, status: false, resetPassword: false, delete: false });
    expect(administratorActions('self', 'self', 'ACTIVE', ['hqMembership:update', 'hqMembership:delete', 'account:update'])).toEqual({ edit: false, status: false, resetPassword: false, delete: false });
		expect(administratorActions('target', 'self', 'INVITED', ['hqMembership:update']).status).toBe(false);
		expect(administratorActions('target', 'self', 'LEFT', ['hqMembership:update']).status).toBe(false);
  });

	it('管理员查询失败时显示错误且不伪装为空列表', () => {
		const state = administratorViewState(null, undefined, false, [], false, false, new Error('network'));
		expect(state.listError).toContain('管理员列表');
		expect(state.empty).toBe(false);
	});

  it('列表展示可读的更新时间', () => {
    expect(formatAdministratorUpdatedAt(0)).toBe('-');
    expect(formatAdministratorUpdatedAt(1_700_000_000_000)).not.toBe('-');
  });

  it('只有总部超级管理员能分配超级管理员角色', () => {
    const roles = [{ id: 'super', name: '超级管理员', kind: 'HQ_SUPER_ADMIN' }, { id: 'custom', name: '运营', kind: 'CUSTOM' }];
    expect(assignableAdministratorRoles(roles, false).map((role) => role.id)).toEqual(['custom']);
    expect(assignableAdministratorRoles(roles, true).map((role) => role.id)).toEqual(['super', 'custom']);
  });

  it('稳定错误码映射为中文操作提示', () => {
    expect(administratorMutationErrorMessage('SELF_MEMBERSHIP_CHANGE_DENIED')).toContain('当前账号');
    expect(administratorMutationErrorMessage('PERMISSION_DELEGATION_DENIED')).toContain('超出');
    expect(administratorMutationErrorMessage('LAST_HQ_SUPER_ADMIN_REQUIRED')).toContain('超级管理员');
    expect(administratorMutationErrorMessage('ROLE_IN_USE')).toContain('使用');
  });
});

describe('总部管理员角色名称', () => {
  it('管理员列表和角色选择显示与角色页面相同的中文名称', () => {
    const role = { id: 'super', name: 'role.hqSuperAdministrator', kind: 'HQ_SUPER_ADMIN' as const };
    const row: AdministratorRow = {
      id: 'membership-1', status: 'ACTIVE', storeAccessMode: 'ALL_STORES', accountId: 'account-1', organizationId: 'hq', updatedAt: null,
      account: { id: 'account-1', phone: '13800000000', displayName: '管理员', email: null }, roles: [role],
    };
    const table = renderToStaticMarkup(<AdministratorTable rows={[row]} currentAccountId="account-1" permissions={[]} onEdit={() => undefined} onStatus={() => undefined} onResetPassword={() => undefined} onDelete={() => undefined} />);
    const roleSelect = renderToStaticMarkup(<AdminRoleSelect id="administrator-roles" roles={[role]} value={['super']} onValueChange={() => undefined} />);
    expect(table).toContain('总部超级管理员');
    expect(table).not.toContain(role.name);
    expect(roleSelect).toContain('总部超级管理员');
    expect(roleSelect).not.toContain(role.name);
  });
});

describe('总部角色权限', () => {
	it('角色查询失败时显示错误且不伪装为空列表', () => {
		const state = roleViewState(undefined, false, [], false, false, new Error('network'));
		expect(state.listError).toContain('角色列表');
		expect(state.empty).toBe(false);
	});

  it('角色表单只接受 SYSTEM 权限并按资源动作分组', () => {
    expect(systemPermissionsOnly(permissions).map((item) => item.id)).toEqual(['read', 'update']);
    expect(groupPermissionActions(systemPermissionsOnly(permissions))).toEqual({ hqMembership: { read: 'read', update: 'update' } });
    expect(delegableSystemPermissions(permissions, ['hqMembership:read']).map((item) => item.id)).toEqual(['read']);
  });

  it('权限目录未完整时拒绝构造角色写入', () => {
    expect(permissionCatalogReady(undefined, false, undefined)).toBe(false);
    expect(permissionCatalogReady({ permissions: { data: permissions.slice(0, 1), total: 2 } }, false, undefined)).toBe(false);
    expect(() => buildRoleMutationInput({ name: '运营', permissionIds: ['read'] }, permissions, ['hqMembership:read'], false)).toThrow('PERMISSION_CATALOG_UNAVAILABLE');
    expect(buildRoleMutationInput({ name: ' 运营 ', permissionIds: ['read', 'tenant'] }, permissions, ['hqMembership:read'], true)).toEqual({ name: '运营', permissionsIds: ['read'] });
  });

  it('系统角色只读且自定义角色可按独立权限管理', () => {
    expect(canManageRole('HQ_SUPER_ADMIN')).toBe(false);
    expect(canManageRole('CUSTOM')).toBe(true);
    const dialog = RoleFormDialog({ open: true, onOpenChange: () => undefined, permissions, viewerPermissions: ['hqMembership:read'], onSubmit: async () => undefined });
    expect(dialog.type).toBe(AdminFormDialogShell);
    expect(JSON.stringify(dialog)).toContain('请输入角色名称');
    expect(JSON.stringify(dialog)).not.toMatch(/编号|编码/);
  });

  it('弹窗内可见显示后端业务错误', () => {
    const administrator = AdministratorFormDialog({ open: true, mode: 'create', onOpenChange: () => undefined, onSubmit: async () => undefined, roles: [], error: '邀请失败' });
    const role = RoleFormDialog({ open: true, onOpenChange: () => undefined, permissions, viewerPermissions: ['hqMembership:read'], onSubmit: async () => undefined, error: '角色保存失败' });
    expect(JSON.stringify(administrator)).toContain('邀请失败');
    expect(JSON.stringify(role)).toContain('角色保存失败');
  });
});
