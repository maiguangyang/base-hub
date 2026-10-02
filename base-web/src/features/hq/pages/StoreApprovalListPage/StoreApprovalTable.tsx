import type { HqStoreApprovalsQuery } from '@/__generated__/graphql';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminStatusIndicator } from '@/features/admin/components/AdminStatusIndicator';
import { toneForStoreBusinessStatus, toneForStoreLifecycle } from '@/features/admin/config/statusTone';
import { storeBusinessStatusOptions, storeLifecycleLabel } from '@/features/admin/config/filterOptions';

export type StoreApprovalRow = NonNullable<HqStoreApprovalsQuery['stores']>['data'][number];

interface StoreApprovalTableProps { rows: StoreApprovalRow[]; approve: boolean; reject: boolean; onReview(row: StoreApprovalRow, approved: boolean): void }

export function StoreApprovalTable(props: StoreApprovalTableProps) {
  return (
    <Table className="min-w-[1900px]"><TableHeader><TableRow className="[&>th]:whitespace-nowrap">
      <TableHead>编码</TableHead><TableHead>门店名称</TableHead><TableHead>所属加盟商</TableHead>
      <TableHead>联系电话</TableHead><TableHead>所在地址</TableHead><TableHead>营业时间</TableHead>
      <TableHead>营业状态</TableHead><TableHead>准入状态</TableHead><TableHead>店长</TableHead>
      <TableHead>店长手机号</TableHead><TableHead>服务模式</TableHead><TableHead>经营面积</TableHead>
      <TableHead>桌位数</TableHead><TableHead>提交时间</TableHead><TableHead className="admin-operation-column">操作</TableHead>
    </TableRow></TableHeader>
      <TableBody>{props.rows.map((row) => <ApprovalTableRow key={row.id} row={row} props={props} />)}</TableBody>
    </Table>
  );
}

function ApprovalTableRow({ row, props }: { row: StoreApprovalRow; props: StoreApprovalTableProps }) {
  const address = [row.province, row.city, row.district, row.address].filter(Boolean).join(' ');
  const businessLabel = storeBusinessStatusOptions.find((option) => option.value === row.businessStatus)?.label ?? '未设置';
  return <TableRow className="[&>td]:whitespace-nowrap">
    <TableCell className="font-mono text-xs">{row.code}</TableCell><TableCell className="font-medium">{row.name}</TableCell>
    <TableCell>{row.organization.name}</TableCell><TableCell>{text(row.contactPhone)}</TableCell>
    <TableCell className="max-w-[280px] truncate" title={address}>{text(address)}</TableCell>
    <TableCell>{text(row.businessHours)}</TableCell>
    <TableCell><AdminStatusIndicator tone={toneForStoreBusinessStatus(row.businessStatus)} label={businessLabel} /></TableCell>
    <TableCell><AdminStatusIndicator tone={toneForStoreLifecycle(row.lifecycle)} label={storeLifecycleLabel(row.lifecycle)} /></TableCell>
    <TableCell>{text(row.managerName)}</TableCell><TableCell>{text(row.managerPhone)}</TableCell>
    <TableCell>{serviceMode(row)}</TableCell><TableCell>{row.storeArea == null ? '—' : `${row.storeArea} ㎡`}</TableCell>
    <TableCell>{row.tableCount == null ? '—' : row.tableCount}</TableCell>
    <TableCell>{submittedTime(row.submittedAt)}</TableCell>
    <TableCell className="text-right"><div className="flex justify-end gap-2">
      {props.approve && <Button size="sm" onClick={() => props.onReview(row, true)}>批准</Button>}
      {props.reject && <Button variant="outline" size="sm" onClick={() => props.onReview(row, false)}>退回</Button>}
    </div></TableCell>
  </TableRow>;
}

function text(value: string | null | undefined): string { return value?.trim() || '—'; }

function submittedTime(value: unknown): string {
  if (typeof value !== 'string' || !value) return '—';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? '—' : date.toLocaleString('zh-CN');
}

function serviceMode(row: StoreApprovalRow): string {
  const modes = [row.supportDineIn === true ? '堂食' : '', row.supportTakeout === true ? '外卖' : ''].filter(Boolean);
  return modes.join(' · ') || (row.supportDineIn === false && row.supportTakeout === false ? '不提供堂食/外卖' : '服务模式未设置');
}
