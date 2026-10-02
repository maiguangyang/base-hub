import type { HqAdministratorsQuery } from '@/__generated__/graphql';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminStatusIndicator } from '@/features/admin/components/AdminStatusIndicator';
import { AdminStatusSwitch } from '@/features/admin/components/AdminStatusSwitch';
import { roleDisplayName } from '@/features/admin/config/roleDisplayName';
import { toneForMembershipStatus } from '@/features/admin/config/statusTone';

export type AdministratorRow = NonNullable<HqAdministratorsQuery['operatorMemberships']>['data'][number];
export interface AdministratorActionSet { edit: boolean; status: boolean; resetPassword: boolean; delete: boolean }

export function administratorActions(accountId: string, currentAccountId: string, status: string, permissions: readonly string[]): AdministratorActionSet {
  if (accountId === currentAccountId) return { edit: false, status: false, resetPassword: false, delete: false };
  const manageableStatus = status === 'ACTIVE' || status === 'SUSPENDED';
  return {
    edit: permissions.includes('hqMembership:update'),
    status: manageableStatus && permissions.includes('hqMembership:update'),
    resetPassword: status === 'ACTIVE' && permissions.includes('account:update'),
    delete: permissions.includes('hqMembership:delete'),
  };
}

interface AdministratorTableProps {
  rows: AdministratorRow[];
  currentAccountId: string;
  permissions: readonly string[];
  onEdit(row: AdministratorRow): void;
  onStatus(row: AdministratorRow, active: boolean): void;
  onResetPassword(row: AdministratorRow): void;
  onDelete(row: AdministratorRow): void;
}

export function AdministratorTable(props: AdministratorTableProps) {
  return <Table><TableHeader><TableRow><TableHead>管理员</TableHead><TableHead>角色</TableHead><TableHead>状态</TableHead><TableHead>更新时间</TableHead><TableHead className="admin-operation-column">操作</TableHead></TableRow></TableHeader><TableBody>
    {props.rows.map((row) => {
      const actions = administratorActions(row.account.id, props.currentAccountId, row.status, props.permissions);
      return <TableRow key={row.id}><TableCell><span className="font-medium">{row.account.displayName}</span><span className="block text-xs text-muted-foreground">{row.account.phone}</span></TableCell><TableCell>{row.roles.map(roleDisplayName).join('、')}</TableCell><TableCell>{row.status === 'INVITED' ? <AdminStatusIndicator tone={toneForMembershipStatus(row.status)} label={row.status} /> : <AdminStatusSwitch checked={row.status === 'ACTIVE'} label={`${row.account.displayName}管理员状态`} disabled={!actions.status} onCheckedChange={(next) => props.onStatus(row, next)} />}</TableCell><TableCell>{formatAdministratorUpdatedAt(row.updatedAt)}</TableCell><TableCell className="text-right"><div className="flex flex-wrap justify-end gap-2">
        {actions.edit && <Button variant="outline" size="sm" onClick={() => props.onEdit(row)}>编辑角色</Button>}
        {actions.resetPassword && <Button variant="outline" size="sm" onClick={() => props.onResetPassword(row)}>重置密码</Button>}
        {actions.delete && <Button variant="destructive" size="sm" onClick={() => props.onDelete(row)}>删除</Button>}
      </div></TableCell></TableRow>;
    })}
  </TableBody></Table>;
}

export function formatAdministratorUpdatedAt(value?: number | null): string {
  if (!value) return '-';
  return new Date(value).toLocaleString('zh-CN');
}
