import { ChevronDown } from 'lucide-react';
import { DropdownMenu, DropdownMenuCheckboxItem, DropdownMenuContent, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { Label } from '@/components/ui/label';
import { roleDisplayName } from '@/features/admin/config/roleDisplayName';

interface RoleOption { id: string; name: string; kind?: string }

export function AdminRoleSelect({ id, roles, value, onValueChange }: {
  id: string;
  roles: readonly RoleOption[];
  value: readonly string[];
  onValueChange(value: string[]): void;
}) {
  const selectedNames = roles.filter((role) => value.includes(role.id)).map(roleDisplayName);
  const summary = value.length === 0 ? '请选择角色' : selectedNames.length === value.length ? selectedNames.join('、') : `已选 ${value.length} 个角色`;
  const toggle = (roleId: string) => onValueChange(value.includes(roleId) ? value.filter((id) => id !== roleId) : [...value, roleId]);
  return <div className="flex flex-col gap-2 text-sm">
    <Label htmlFor={id}>角色</Label>
    <DropdownMenu>
      <DropdownMenuTrigger asChild><button id={id} type="button" aria-label={`角色：${summary}`} title={summary}
        data-admin-form-surface=""
        className="flex h-10 w-full items-center justify-between gap-2 rounded-md border border-input px-3 text-left text-sm shadow-xs outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50">
        <span className="truncate">{summary}</span><ChevronDown className="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
      </button></DropdownMenuTrigger>
      <DropdownMenuContent data-admin-form-surface="" align="start" className="max-h-64 w-[var(--radix-dropdown-menu-trigger-width)] overflow-y-auto">
        {roles.length === 0 && <p className="px-2 py-2 text-sm text-muted-foreground">暂无可选角色</p>}
        {roles.map((role) => <DropdownMenuCheckboxItem key={role.id} checked={value.includes(role.id)}
          onCheckedChange={() => toggle(role.id)} onSelect={(event) => event.preventDefault()}>
          {roleDisplayName(role)}
        </DropdownMenuCheckboxItem>)}
      </DropdownMenuContent>
    </DropdownMenu>
  </div>;
}
