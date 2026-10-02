import type { FranchiseAuditLogsQuery } from '@/__generated__/graphql';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { AdminStatusIndicator } from '@/features/admin/components/AdminStatusIndicator';
import { toneForAuditResult } from '@/features/admin/config/statusTone';

export type FranchiseAuditRow = NonNullable<FranchiseAuditLogsQuery['auditLogs']>['data'][number];

export function AuditLogTable({ rows }: { rows: FranchiseAuditRow[] }) {
  return <Table><TableHeader><TableRow><TableHead>动作</TableHead><TableHead>资源</TableHead><TableHead>门店</TableHead><TableHead>结果</TableHead></TableRow></TableHeader><TableBody>{rows.map((row) => <TableRow key={row.id}><TableCell>{row.action}</TableCell><TableCell>{row.resourceType} · {row.resourceId ?? '-'}</TableCell><TableCell>{row.storeId ?? '-'}</TableCell><TableCell><AdminStatusIndicator tone={toneForAuditResult(row.resultCode)} label={row.resultCode} /></TableCell></TableRow>)}</TableBody></Table>;
}
