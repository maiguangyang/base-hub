import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { ScrollArea } from '@/components/ui/scroll-area';
import { permissionDisplayName, permissionGroupLabel, type PermissionOption } from './permissionLabels';

interface PermissionGroup { key: string; label: string; permissions: PermissionOption[] }

const permissionCheckboxClass = 'cursor-pointer bg-white data-[state=checked]:border-primary data-[state=checked]:bg-white data-[state=checked]:text-primary data-[state=indeterminate]:border-primary data-[state=indeterminate]:text-primary';

function permissionGroups(permissions: PermissionOption[]): PermissionGroup[] {
  const groups = new Map<string, PermissionGroup>();
  for (const permission of permissions) {
    const key = permission.module || permission.action.split(':')[0] || 'other';
    const group = groups.get(key) ?? { key, label: permissionGroupLabel(key), permissions: [] };
    group.permissions.push(permission);
    groups.set(key, group);
  }
  return [...groups.values()];
}

function selectionState(selected: Set<string>, ids: string[]): boolean | 'indeterminate' {
  const count = ids.filter((id) => selected.has(id)).length;
  if (count === 0) return false;
  return count === ids.length ? true : 'indeterminate';
}

interface RolePermissionSelectorProps {
  permissions: PermissionOption[];
  selectedIDs: string[];
  onChange(selectedIDs: string[]): void;
}

export function RolePermissionSelector({ permissions, selectedIDs, onChange }: RolePermissionSelectorProps) {
  const groups = permissionGroups(permissions);
  const permissionIDs = permissions.map((permission) => permission.id);
  const selected = new Set(selectedIDs.filter((id) => permissionIDs.includes(id)));
  const toggle = (ids: string[]) => {
    const next = new Set(selected);
    const remove = ids.every((id) => next.has(id));
    for (const id of ids) {
      if (remove) next.delete(id);
      else next.add(id);
    }
    onChange(permissionIDs.filter((id) => next.has(id)));
  };

  return <section className="flex flex-col gap-2" aria-labelledby="role-permission-label">
    <h3 id="role-permission-label" className="text-sm font-medium">权限分配 *</h3>
    <ScrollArea data-testid="role-permission-scroll" className="h-80 min-h-56 max-h-[45dvh] rounded-lg border border-border bg-card">
      {permissions.length > 0 ? <>
        <div className="sticky top-0 z-10 flex min-h-12 items-center gap-2 border-b border-border bg-card px-4 py-3 pr-5">
          <Checkbox id="select-all-role-permissions" className={permissionCheckboxClass} aria-label="全选全部权限" checked={selectionState(selected, permissionIDs)} onCheckedChange={() => toggle(permissionIDs)} />
          <Label htmlFor="select-all-role-permissions" className="cursor-pointer font-semibold">全选</Label>
          <span className="text-xs text-muted-foreground">({selected.size}/{permissions.length})</span>
        </div>
        <div className="flex flex-col gap-3 p-4 pr-5">
          {groups.map((group) => <PermissionGroupSection key={group.key} group={group} selected={selected} onToggle={toggle} />)}
        </div>
      </> : <p className="p-4 py-8 text-center text-sm text-muted-foreground">暂无可分配权限</p>}
    </ScrollArea>
  </section>;
}

function PermissionGroupSection({ group, selected, onToggle }: { group: PermissionGroup; selected: Set<string>; onToggle(ids: string[]): void }) {
  const ids = group.permissions.map((permission) => permission.id);
  const selectedCount = ids.filter((id) => selected.has(id)).length;
  const groupID = `role-permission-group-${group.key}`;
  return <div className="flex flex-col gap-1">
    <div className="flex min-h-8 items-center gap-2">
      <Checkbox id={groupID} className={permissionCheckboxClass} aria-label={`全选 ${group.label} 全部权限`} checked={selectionState(selected, ids)} onCheckedChange={() => onToggle(ids)} />
      <Label htmlFor={groupID} className="cursor-pointer font-semibold">{group.label}</Label>
      <span className="text-xs text-muted-foreground">({selectedCount}/{ids.length})</span>
    </div>
    <div className="grid grid-cols-1 gap-x-8 gap-y-0 pl-6 sm:grid-cols-2">
      {group.permissions.map((permission) => <PermissionCheckbox key={permission.id} permission={permission} checked={selected.has(permission.id)} onToggle={onToggle} />)}
    </div>
  </div>;
}

function PermissionCheckbox({ permission, checked, onToggle }: { permission: PermissionOption; checked: boolean; onToggle(ids: string[]): void }) {
  const checkboxID = `role-permission-${permission.id}`;
  const displayName = permissionDisplayName(permission);
  return <div className="flex min-h-8 items-center gap-2">
    <Checkbox id={checkboxID} className={permissionCheckboxClass} aria-label={displayName} checked={checked} onCheckedChange={() => onToggle([permission.id])} />
    <Label htmlFor={checkboxID} className="min-w-0 cursor-pointer font-normal">{displayName}</Label>
  </div>;
}
