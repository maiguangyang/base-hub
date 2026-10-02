import { Input } from '@/components/ui/input';
import { AdminActionError, AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { RolePermissionSelector } from '@/features/admin/components/RolePermissionSelector/RolePermissionSelector';
import type { PermissionOption } from '@/features/admin/components/RolePermissionSelector/permissionLabels';

export type { PermissionOption } from '@/features/admin/components/RolePermissionSelector/permissionLabels';
export interface RoleValues { name: string; permissionIds: string[] }

interface RoleFormDialogProps {
  open: boolean;
  mode?: 'create' | 'edit';
  onOpenChange(open: boolean): void;
  permissions: PermissionOption[];
  viewerPermissions: readonly string[];
  values?: RoleValues;
  onValuesChange?(values: RoleValues): void;
  onSubmit(): void | Promise<void>;
  isSubmitting?: boolean;
  error?: string;
}

export function systemPermissionsOnly<T extends PermissionOption>(permissions: T[]): T[] {
  return permissions.filter((permission) => permission.scope === 'SYSTEM');
}

export function delegableSystemPermissions<T extends PermissionOption>(permissions: T[], viewerPermissions: readonly string[]): T[] {
  const allowed = new Set(viewerPermissions);
  return systemPermissionsOnly(permissions).filter((permission) => allowed.has(permission.action));
}

export function groupPermissionActions(permissions: PermissionOption[]): Record<string, Record<string, string>> {
  const groups: Record<string, Record<string, string>> = {};
  for (const permission of permissions) {
    const [resource, action] = permission.action.split(':');
    if (!resource || !action) continue;
    groups[resource] ??= {};
    groups[resource][action] = permission.id;
  }
  return groups;
}

export function canManageRole(kind: string): boolean { return kind === 'CUSTOM'; }

export function RoleFormDialog(props: RoleFormDialogProps) {
  const values = props.values ?? { name: '', permissionIds: [] };
  const permissions = delegableSystemPermissions(props.permissions, props.viewerPermissions);
  const title = props.mode === 'edit' ? '编辑角色' : '新增角色';
  const description = props.mode === 'edit' ? '修改角色名称和权限配置。' : '创建一个新的总部角色并分配权限。';
  return <AdminFormDialogShell open={props.open} onOpenChange={props.onOpenChange} title={title} description={description} submitLabel="保存" cancelLabel="取消" isSubmitting={props.isSubmitting} onSubmit={props.onSubmit} maxWidth="3xl">
    <div className="flex flex-col gap-5"><AdminActionError message={props.error} /><label className="flex flex-col gap-2 text-sm font-medium"><span>角色名称 *</span><Input placeholder="请输入角色名称，如：运营人员" value={values.name} onChange={(event) => props.onValuesChange?.({ ...values, name: event.target.value })} maxLength={64} required /></label>
      <RolePermissionSelector permissions={permissions} selectedIDs={values.permissionIds} onChange={(permissionIds) => props.onValuesChange?.({ ...values, permissionIds })} />
    </div>
  </AdminFormDialogShell>;
}
