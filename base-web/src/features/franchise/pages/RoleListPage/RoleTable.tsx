import type { FranchiseRolesQuery } from '@/__generated__/graphql';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { AdminRolePermissionSummary } from '@/features/admin/components/AdminRolePermissionSummary';
import { roleDisplayName } from '@/features/admin/config/roleDisplayName';
import { canManageRole } from './RoleFormDialog';

export type RoleRow = NonNullable<FranchiseRolesQuery['operatorRoles']>['data'][number];

interface RoleTableProps { rows: RoleRow[]; canUpdate: boolean; canDelete: boolean; onEdit(row: RoleRow): void; onDelete(row: RoleRow): void }

export function RoleTable(props: RoleTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[28%]">角色名称</TableHead>
          <TableHead className="w-40">类型</TableHead>
          <TableHead>权限范围</TableHead>
          <TableHead className="admin-operation-column">操作</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {props.rows.map((row) => {
          const manageable = canManageRole(row.kind);
          return (
            <TableRow key={row.id}>
              <TableCell className="font-medium">{roleDisplayName(row)}</TableCell>
              <TableCell>{manageable ? '自定义' : <Badge variant="outline">系统角色</Badge>}</TableCell>
              <TableCell><AdminRolePermissionSummary kind={row.kind} permissionCount={row.permissions.length} /></TableCell>
              <TableCell className="text-right">
                <div className="flex justify-end gap-2">
                  {manageable && props.canUpdate && <Button variant="outline" size="sm" onClick={() => props.onEdit(row)}>编辑</Button>}
                  {manageable && props.canDelete && <Button variant="destructive" size="sm" onClick={() => props.onDelete(row)}>删除</Button>}
                </div>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
