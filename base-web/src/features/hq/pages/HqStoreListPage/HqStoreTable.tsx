import type { HqStoresQuery } from '@/__generated__/graphql';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminStatusIndicator } from '@/features/admin/components/AdminStatusIndicator';
import { AdminStatusSwitch } from '@/features/admin/components/AdminStatusSwitch';
import { toneForStoreLifecycle } from '@/features/admin/config/statusTone';
import { storeLifecycleLabel } from '@/features/admin/config/filterOptions';
import { StorePaymentConfigButton } from '../PaymentConfigPage/StorePaymentConfigButton';

export type HqStoreRow = NonNullable<HqStoresQuery['stores']>['data'][number];

interface Props {
  rows: HqStoreRow[];
  hqOrganizationId: string;
  canUpdate: boolean;
  canDelete: boolean;
  statusBusy: boolean;
  onEdit(row: HqStoreRow): void;
  onDelete(row: HqStoreRow): void;
  onBusinessStatus(row: HqStoreRow, open: boolean): void;
  onConfigurePayment?(row: HqStoreRow): void;
}

export function HqStoreTable(props: Props) {
  return <Table>
    <TableHeader><TableRow>
      <TableHead>编码</TableHead><TableHead>门店名称</TableHead><TableHead>类型</TableHead>
      <TableHead>所属组织</TableHead><TableHead>联系电话</TableHead>
      <TableHead>所在地址</TableHead><TableHead>营业时间</TableHead>
      <TableHead>营业状态</TableHead><TableHead>准入状态</TableHead><TableHead className="admin-operation-column">操作</TableHead>
    </TableRow></TableHeader>
    <TableBody>{props.rows.map((row) => <HqStoreTableRow key={row.id} row={row} props={props} />)}</TableBody>
  </Table>;
}

function HqStoreTableRow({ row, props }: { row: HqStoreRow; props: Props }) {
  const direct = row.organization.type === 'HEADQUARTERS' && row.organizationId === props.hqOrganizationId;
  const canEdit = direct && props.canUpdate;
  const canDelete = direct && props.canDelete;
  return <TableRow>
    <TableCell className="font-mono text-xs">{row.code}</TableCell>
    <TableCell className="font-medium">{row.name}</TableCell>
    <TableCell>{row.organization.type === 'HEADQUARTERS' ? '直营' : '加盟'}</TableCell>
    <TableCell>{row.organization.name}</TableCell>
    <TableCell>{row.contactPhone || '—'}</TableCell>
    <TableCell className="max-w-[200px] truncate" title={addressText(row)}>{addressText(row)}</TableCell>
    <TableCell>{row.businessHours || '—'}</TableCell>
    <TableCell>{row.businessStatus ? <AdminStatusSwitch checked={row.businessStatus === 'OPEN'} label={`${row.name}营业状态`}
      disabled={!canEdit || props.statusBusy} onCheckedChange={(next) => props.onBusinessStatus(row, next)} /> : '未设置'}</TableCell>
    <TableCell><AdminStatusIndicator tone={toneForStoreLifecycle(row.lifecycle)} label={storeLifecycleLabel(row.lifecycle)} /></TableCell>
    <TableCell className="text-right"><HqStoreActions row={row} props={props} canEdit={canEdit} canDelete={canDelete} /></TableCell>
  </TableRow>;
}

function addressText(row: HqStoreRow): string {
  return [row.city, row.district, row.address].filter(Boolean).join(' ') || '—';
}

function HqStoreActions({ row, props, canEdit, canDelete }: {
  row: HqStoreRow; props: Props; canEdit: boolean; canDelete: boolean;
}) {
  return <div className="flex justify-end gap-2">
    {canEdit && <Button variant="outline" size="sm" onClick={() => props.onEdit(row)}>编辑</Button>}
    {canDelete && <Button variant="destructive" size="sm" onClick={() => props.onDelete(row)}>删除</Button>}
    {row.organization.type === 'FRANCHISE' && props.onConfigurePayment && <StorePaymentConfigButton lifecycle={row.lifecycle} onClick={() => props.onConfigurePayment?.(row)} />}
  </div>;
}
