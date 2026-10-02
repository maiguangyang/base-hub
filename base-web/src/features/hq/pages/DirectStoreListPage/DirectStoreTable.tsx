import type { HqDirectStoresQuery } from '@/__generated__/graphql';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { AdminStatusIndicator } from '@/features/admin/components/AdminStatusIndicator';
import { AdminStatusSwitch } from '@/features/admin/components/AdminStatusSwitch';
import { toneForStoreLifecycle } from '@/features/admin/config/statusTone';
import { storeLifecycleLabel } from '@/features/admin/config/filterOptions';

export type DirectStoreRow = NonNullable<HqDirectStoresQuery['stores']>['data'][number];

interface DirectStoreTableProps {
  rows: DirectStoreRow[];
  canUpdate: boolean;
  canDelete: boolean;
  onEdit(row: DirectStoreRow): void;
  onDelete(row: DirectStoreRow): void;
  onBusinessStatus(row: DirectStoreRow, open: boolean): void;
  statusBusy: boolean;
}

export function DirectStoreTable({ rows, canUpdate, canDelete, onEdit, onDelete, onBusinessStatus, statusBusy }: DirectStoreTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>编码</TableHead>
          <TableHead>门店名称</TableHead>
          <TableHead>联系电话</TableHead>
          <TableHead>所在地址</TableHead>
          <TableHead>营业时间</TableHead>
          <TableHead>营业状态</TableHead>
          <TableHead>准入状态</TableHead>
          <TableHead className="admin-operation-column">操作</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {rows.map((row) => (
          <TableRow key={row.id}>
            <TableCell className="font-mono text-xs">{row.code}</TableCell>
            <TableCell className="font-medium">{row.name}</TableCell>
            <TableCell>{row.contactPhone || '-'}</TableCell>
            <TableCell className="max-w-[200px] truncate" title={[row.province, row.city, row.district, row.address].filter(Boolean).join(' ')}>
              {[row.city, row.district, row.address].filter(Boolean).join(' ') || '-'}
            </TableCell>
            <TableCell>{row.businessHours || '-'}</TableCell>
            <TableCell>
              {row.businessStatus ? <AdminStatusSwitch checked={row.businessStatus === 'OPEN'} label={`${row.name}营业状态`}
                disabled={!canUpdate || statusBusy} onCheckedChange={(next) => onBusinessStatus(row, next)} /> : '未设置'}
            </TableCell>
            <TableCell>
              <AdminStatusIndicator tone={toneForStoreLifecycle(row.lifecycle)} label={storeLifecycleLabel(row.lifecycle)} />
            </TableCell>
            <TableCell className="text-right">
              <div className="flex justify-end gap-2">
                {canUpdate && <Button variant="outline" size="sm" onClick={() => onEdit(row)}>编辑</Button>}
                {canDelete && <Button variant="destructive" size="sm" onClick={() => onDelete(row)}>删除</Button>}
              </div>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

export function canDeleteDirectStore(permissions: string[]): boolean {
  return permissions.includes('hqStore:delete');
}
