import type { HqFranchiseStoresQuery } from '@/__generated__/graphql';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminStatusIndicator } from '@/features/admin/components/AdminStatusIndicator';
import { AdminStatusSwitch } from '@/features/admin/components/AdminStatusSwitch';
import { toneForStoreLifecycle } from '@/features/admin/config/statusTone';
import { storeLifecycleLabel } from '@/features/admin/config/filterOptions';
import { StorePaymentConfigButton } from '../PaymentConfigPage/StorePaymentConfigButton';

export type FranchiseStoreRow = NonNullable<HqFranchiseStoresQuery['stores']>['data'][number];

export function FranchiseStoreTable({ rows, onConfigurePayment }: { rows: FranchiseStoreRow[]; onConfigurePayment?(row: FranchiseStoreRow): void }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>编码</TableHead>
          <TableHead>门店名称</TableHead>
          <TableHead>所属加盟商</TableHead>
          <TableHead>联系电话</TableHead>
          <TableHead>营业状态</TableHead>
          <TableHead>准入状态</TableHead>
          {onConfigurePayment && <TableHead className="admin-operation-column">操作</TableHead>}
        </TableRow>
      </TableHeader>
      <TableBody>
        {rows.map((row) => (
          <TableRow key={row.id}>
            <TableCell className="font-mono text-xs">{row.code}</TableCell>
            <TableCell className="font-medium">{row.name}</TableCell>
            <TableCell>{row.organization.name}</TableCell>
            <TableCell>{row.contactPhone || '—'}</TableCell>
            <TableCell>{row.businessStatus ? <AdminStatusSwitch checked={row.businessStatus === 'OPEN'} label={`${row.name}营业状态`} disabled onCheckedChange={() => undefined} /> : '未设置'}</TableCell>
            <TableCell><AdminStatusIndicator tone={toneForStoreLifecycle(row.lifecycle)} label={storeLifecycleLabel(row.lifecycle)} /></TableCell>
            {onConfigurePayment && <TableCell className="text-right"><StorePaymentConfigButton lifecycle={row.lifecycle} onClick={() => onConfigurePayment(row)} /></TableCell>}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
