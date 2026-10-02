import type { FranchiseStaffQuery } from '@/__generated__/graphql';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminStatusIndicator } from '@/features/admin/components/AdminStatusIndicator';
import { AdminStatusSwitch } from '@/features/admin/components/AdminStatusSwitch';
import { roleDisplayName } from '@/features/admin/config/roleDisplayName';
import { toneForMembershipStatus } from '@/features/admin/config/statusTone';

export type StaffRow = NonNullable<FranchiseStaffQuery['operatorMemberships']>['data'][number];

interface StaffTableProps {
  rows: StaffRow[];
  canUpdate: boolean;
  onEdit(row: StaffRow): void;
  onStatus(row: StaffRow, active: boolean): void;
}

export function StaffTable(props: StaffTableProps) {
  return (
    <Table><TableHeader><TableRow><TableHead>员工</TableHead><TableHead>角色</TableHead><TableHead>门店范围</TableHead><TableHead>状态</TableHead><TableHead className="admin-operation-column">操作</TableHead></TableRow></TableHeader>
      <TableBody>{props.rows.map((row) => <TableRow key={row.id}><TableCell>{row.account.displayName}<span className="block text-xs text-muted-foreground">{row.account.phone}</span></TableCell><TableCell>{row.roles.map(roleDisplayName).join('、')}</TableCell><TableCell>{row.storeAccessMode === 'ALL_STORES' ? '全部门店' : row.stores.map((store) => store.name).join('、')}</TableCell><TableCell>{row.status === 'INVITED' ? <AdminStatusIndicator tone={toneForMembershipStatus(row.status)} label={row.status} /> : <AdminStatusSwitch checked={row.status === 'ACTIVE'} label={`${row.account.displayName}员工状态`} disabled={!props.canUpdate} onCheckedChange={(next) => props.onStatus(row, next)} />}</TableCell><TableCell className="text-right"><div className="flex justify-end gap-2">{props.canUpdate && <Button variant="outline" size="sm" onClick={() => props.onEdit(row)}>编辑</Button>}</div></TableCell></TableRow>)}</TableBody>
    </Table>
  );
}
