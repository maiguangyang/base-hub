import { useState } from 'react';
import { useQuery } from '@apollo/client/react';
import { AdminBulkBar } from '@/features/admin/components/AdminBulkBar';
import { AdminActionError } from '@/features/admin/components/AdminFormDialogShell';
import { AdminFilterSelect } from '@/features/admin/components/AdminFilterSelect';
import { AdminListFilters } from '@/features/admin/components/AdminListFilters';
import { AdminListSearch } from '@/features/admin/components/AdminListSearch';
import { AdminPageHeader } from '@/features/admin/components/AdminPageHeader';
import { AdminPagination } from '@/features/admin/components/AdminPagination';
import { AdminStatsStrip } from '@/features/admin/components/AdminStatsStrip';
import { AdminTableShell } from '@/features/admin/components/AdminTableShell';
import { AdminToolbar } from '@/features/admin/components/AdminToolbar';
import { auditResourceTypeOptions, auditResultOptions } from '@/features/admin/config/filterOptions';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { HQ_AUDIT_LOGS_QUERY } from '../../graphql/audit';
import { AuditLogTable } from './AuditLogTable';
import { auditStats } from './auditStats';
import { hqAuditFilter, type AuditResultFilter } from '../hqListFilters';

export function AuditLogListPage() {
  const { search } = useAdminTab();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState(() => new URLSearchParams(search).get('q') ?? '');
  const [result, setResult] = useState<AuditResultFilter>();
  const [resourceType, setResourceType] = useState<string>();
  const { data, loading, error } = useQuery(HQ_AUDIT_LOGS_QUERY, { variables: { page, pageSize, q: keyword || null, filter: hqAuditFilter(result, resourceType) } });
  const rows = data?.auditLogs?.data ?? [];
  const total = data?.auditLogs?.total ?? 0;
  const changeResult = (next?: string) => { setResult(next as AuditResultFilter); setPage(1); };
  const changeResource = (next?: string) => { setResourceType(next); setPage(1); };
  const resetFilters = () => { setKeyword(''); setResult(undefined); setResourceType(undefined); setPage(1); };
  const view = auditListView(loading, error, rows.length);
  return (
    <div><AdminPageHeader title="审计日志" description="跨组织只读查看治理与安全操作记录。" /><AdminStatsStrip items={auditStats(total)} /><AdminBulkBar count={0} onClear={() => undefined} />
      <AdminActionError message={view.errorMessage} />
      <AdminTableShell toolbar={<AdminToolbar><AdminListFilters hasActiveFilters={[keyword, result, resourceType].some(Boolean)} onReset={resetFilters}><AdminListSearch value={keyword} placeholder="搜索审计记录" onSearch={(next) => { setKeyword(next); setPage(1); }} /><AdminFilterSelect value={result} placeholder="全部执行结果" options={auditResultOptions} onValueChange={changeResult} /><AdminFilterSelect value={resourceType} placeholder="全部资源类型" options={auditResourceTypeOptions} onValueChange={changeResource} /></AdminListFilters></AdminToolbar>} isEmpty={view.empty} empty="暂无审计记录" pagination={<AdminPagination total={total} page={page} pageSize={pageSize} onPageChange={setPage} onPageSizeChange={setPageSize} />}>
        <AuditLogTable rows={rows} />
      </AdminTableShell>
    </div>
  );
}

function auditListView(loading: boolean, error: unknown, rowCount: number) {
  return { errorMessage: error ? '审计日志加载失败，请稍后重试。' : undefined, empty: !loading && !error && rowCount === 0 };
}
