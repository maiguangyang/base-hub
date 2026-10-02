import type { FranchiseStoresQuery } from '@/__generated__/graphql';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminStatusIndicator } from '@/features/admin/components/AdminStatusIndicator';
import { storeBusinessStatusOptions, storeLifecycleLabel } from '@/features/admin/config/filterOptions';
import { toneForStoreBusinessStatus, toneForStoreLifecycle } from '@/features/admin/config/statusTone';
import { canEditStoreRecord, canSubmitFranchiseStore } from '@/features/admin/pages/StoreEditorPage/storeEditorPolicy';

export type FranchiseStoreRow = NonNullable<FranchiseStoresQuery['stores']>['data'][number];

interface StoreTableProps { rows: FranchiseStoreRow[]; organizationId: string; permissions: string[]; onEdit(row: FranchiseStoreRow): void; onSubmit(row: FranchiseStoreRow): void; onDelete(row: FranchiseStoreRow): void }

export function StoreTable(props: StoreTableProps) {
  return <Table className="min-w-[1600px]">
    <TableHeader><TableRow>
      <TableHead>编码</TableHead><TableHead>门店名称</TableHead><TableHead>联系电话</TableHead><TableHead>所在地址</TableHead>
      <TableHead>营业时间</TableHead><TableHead>营业状态</TableHead><TableHead>准入状态</TableHead>
      <TableHead>店长</TableHead><TableHead>店长手机号</TableHead><TableHead>服务模式</TableHead>
      <TableHead>经营面积</TableHead><TableHead>桌位数</TableHead>
      <TableHead className="admin-operation-column">操作</TableHead>
    </TableRow></TableHeader>
    <TableBody>{props.rows.map((row) => <StoreTableRow key={row.id} row={row} props={props} />)}</TableBody>
  </Table>;
}

function text(value: string | null | undefined): string { return value?.trim() || '—'; }

function StoreTableRow({ row, props }: { row: FranchiseStoreRow; props: StoreTableProps }) {
  const address = [row.province, row.city, row.district, row.address].filter(Boolean).join(' ');
  const businessLabel = storeBusinessStatusOptions.find((option) => option.value === row.businessStatus)?.label ?? '未设置';
  return <TableRow className="[&>td]:whitespace-nowrap">
    <TableCell className="font-mono text-xs text-muted-foreground">{row.code}</TableCell>
    <TableCell className="font-medium">{row.name}</TableCell>
    <TableCell>{text(row.contactPhone)}</TableCell>
    <TableCell className="max-w-[280px] truncate" title={address}>{text(address)}</TableCell>
    <TableCell>{text(row.businessHours)}</TableCell>
    <TableCell><AdminStatusIndicator tone={toneForStoreBusinessStatus(row.businessStatus)} label={businessLabel} /></TableCell>
    <TableCell title={row.lifecycle === 'REJECTED' ? row.rejectionReason ?? undefined : undefined}>
      <AdminStatusIndicator tone={toneForStoreLifecycle(row.lifecycle)} label={storeLifecycleLabel(row.lifecycle)} />
    </TableCell>
    <TableCell>{text(row.managerName)}</TableCell>
    <TableCell>{text(row.managerPhone)}</TableCell>
    <TableCell>{serviceMode(row)}</TableCell>
    <TableCell>{row.storeArea == null ? '—' : `${row.storeArea} ㎡`}</TableCell>
    <TableCell>{row.tableCount == null ? '—' : `${row.tableCount} 桌`}</TableCell>
    <TableCell className="text-right"><StoreActions row={row} props={props} /></TableCell>
  </TableRow>;
}

function serviceMode(row: FranchiseStoreRow): string {
  const modes = [row.supportDineIn === true ? '堂食' : '', row.supportTakeout === true ? '外卖' : ''].filter(Boolean);
  return modes.join(' · ') || (row.supportDineIn === false && row.supportTakeout === false ? '不提供堂食/外卖' : '服务模式未设置');
}

function StoreActions({ row, props }: { row: FranchiseStoreRow; props: StoreTableProps }) {
  return <div className="flex justify-end gap-2">
    {canEditStoreRecord('FRANCHISE', props.organizationId, props.permissions, row) && <Button variant="outline" size="sm" onClick={() => props.onEdit(row)}>编辑</Button>}
    {canSubmitFranchiseStore(row.lifecycle, props.permissions) && <Button size="sm" onClick={() => props.onSubmit(row)}>提交审核</Button>}
    {canDeleteStore(row.lifecycle, props.permissions) && <Button variant="destructive" size="sm" onClick={() => props.onDelete(row)}>删除</Button>}
  </div>;
}

export function canDeleteStore(lifecycle: FranchiseStoreRow['lifecycle'], permissions: string[]): boolean {
  return permissions.includes('store:delete') && (lifecycle === 'DRAFT' || lifecycle === 'REJECTED');
}
