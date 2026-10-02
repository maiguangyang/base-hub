import { Input } from '@/components/ui/input';
import { AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';
import { RolePermissionSelector } from '@/features/admin/components/RolePermissionSelector/RolePermissionSelector';
import type { PermissionOption } from '@/features/admin/components/RolePermissionSelector/permissionLabels';

export type { PermissionOption } from '@/features/admin/components/RolePermissionSelector/permissionLabels';
export interface RoleValues { name: string; permissionIds: string[] }

interface RoleFormDialogProps {
  open: boolean;
  mode?: 'create' | 'edit';
  onOpenChange(open: boolean): void;
  organizationId: string;
  permissions: PermissionOption[];
  values?: RoleValues;
  onValuesChange?(values: RoleValues): void;
  onSubmit(): void | Promise<void>;
  isSubmitting?: boolean;
}

export function tenantPermissionsOnly<T extends PermissionOption>(permissions: T[]): T[] {
  return permissions.filter((permission) => permission.scope === 'TENANT');
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
  const permissions = tenantPermissionsOnly(props.permissions);
  return (
    <AdminFormDialogShell open={props.open} onOpenChange={props.onOpenChange} title={props.mode === 'edit' ? '编辑角色' : '新增角色'} description={props.mode === 'edit' ? '修改角色名称和权限配置。' : '创建一个新的加盟商角色并分配权限。'} submitLabel="保存" cancelLabel="取消" isSubmitting={props.isSubmitting} onSubmit={props.onSubmit} maxWidth="3xl">
      <div className="flex flex-col gap-5"><label className="flex flex-col gap-2 text-sm font-medium"><span>角色名称 *</span><Input placeholder="请输入角色名称，如：门店店长" value={values.name} onChange={(event) => props.onValuesChange?.({ ...values, name: event.target.value })} maxLength={64} required /></label>
        <RolePermissionSelector permissions={permissions} selectedIDs={values.permissionIds} onChange={(permissionIds) => props.onValuesChange?.({ ...values, permissionIds })} />
      </div>
    </AdminFormDialogShell>
  );
}
