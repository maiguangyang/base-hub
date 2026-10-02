import { useState } from 'react';
import { useQuery } from '@apollo/client/react';
import type { FranchiseAuditLogsQuery } from '@/__generated__/graphql';
import { AdminBulkBar } from '@/features/admin/components/AdminBulkBar';
import { AdminFilterSelect } from '@/features/admin/components/AdminFilterSelect';
import { AdminActionError } from '@/features/admin/components/AdminFormDialogShell';
import { AdminListFilters } from '@/features/admin/components/AdminListFilters';
import { AdminListSearch } from '@/features/admin/components/AdminListSearch';
import { AdminPageHeader } from '@/features/admin/components/AdminPageHeader';
import { AdminPagination } from '@/features/admin/components/AdminPagination';
import { AdminStatsStrip } from '@/features/admin/components/AdminStatsStrip';
import { AdminTableShell } from '@/features/admin/components/AdminTableShell';
import { AdminToolbar } from '@/features/admin/components/AdminToolbar';
import { auditResourceTypeOptions, auditResultOptions } from '@/features/admin/config/filterOptions';
import { catalogPage } from '@/features/admin/lib/filterCatalog';
import { useCompleteFilterCatalog } from '@/features/admin/lib/useCompleteFilterCatalog';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { useAuthStore } from '@/features/auth/store/authStore';
import { FRANCHISE_AUDIT_LOGS_QUERY } from '../../graphql/audit';
import { FRANCHISE_STORES_QUERY } from '../../graphql/stores';
import { franchiseAuditFilter, type FranchiseAuditResultFilter } from '../franchiseListFilters';
import { initialSearchKeyword, queryKeyword, viewerOrganizationId } from '../pageState';
import { AuditLogTable } from './AuditLogTable';
import { auditStats } from './auditStats';

export function AuditLogListPage() {
  const state = useAuditLogPage();
  return (
    <div><AdminPageHeader title="操作审计" description="仅查看当前加盟组织的操作记录。" /><AdminStatsStrip items={auditStats(state.total)} /><AdminBulkBar count={0} onClear={() => undefined} />
      <AdminActionError message={state.listError ?? state.catalogError} />
      <AdminTableShell toolbar={<AdminToolbar><AdminListFilters hasActiveFilters={state.hasFilters} onReset={state.resetFilters}><AdminListSearch value={state.keyword} placeholder="搜索审计记录" onSearch={state.search} /><AdminFilterSelect value={state.result} placeholder="全部执行结果" options={auditResultOptions} onValueChange={state.changeResult} /><AdminFilterSelect value={state.resourceType} placeholder="全部资源类型" options={auditResourceTypeOptions} onValueChange={state.changeResource} /><AdminFilterSelect value={state.storeId} placeholder="全部门店" options={state.storeOptions} disabled={!state.storeCatalogReady} onValueChange={state.changeStore} /></AdminListFilters></AdminToolbar>} isEmpty={state.empty} empty="暂无审计记录" pagination={<AdminPagination total={state.total} page={state.page} pageSize={state.pageSize} onPageChange={state.setPage} onPageSizeChange={state.setPageSize} />}>
        <AuditLogTable rows={state.rows} />
      </AdminTableShell>
    </div>
  );
}

function useAuditLogPage() {
  const { search } = useAdminTab();
  const organizationId = viewerOrganizationId(useAuthStore((state) => state.viewer));
  const [page, setPage] = useState(1); const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState(() => initialSearchKeyword(search));
  const [result, setResult] = useState<FranchiseAuditResultFilter>(); const [resourceType, setResourceType] = useState<string>(); const [storeId, setStoreId] = useState<string>();
  const { data, loading, error } = useQuery(FRANCHISE_AUDIT_LOGS_QUERY, { variables: { page, pageSize, q: queryKeyword(keyword), filter: franchiseAuditFilter(result, resourceType, storeId) } });
  const storesQuery = useQuery(FRANCHISE_STORES_QUERY, { variables: { page: 1, pageSize: 200, q: null, filter: undefined } });
  const resultView = auditResult(data);
  const directory = useCompleteFilterCatalog(`franchise-stores:${organizationId}`, catalogPage(storesQuery.data?.stores), storesQuery.loading, storesQuery.error, async (nextPage) => {
    const next = await storesQuery.fetchMore({ variables: { page: nextPage } });
    const pageResult = catalogPage(next.data?.stores);
    if (!pageResult) throw new Error('catalog unavailable');
    return pageResult;
  });
  const changeResult = (next?: string) => { setResult(next as FranchiseAuditResultFilter); setPage(1); };
  const changeResource = (next?: string) => { setResourceType(next); setPage(1); };
  const changeStore = (next?: string) => { setStoreId(next); setPage(1); };
  const resetFilters = () => { setKeyword(''); setResult(undefined); setResourceType(undefined); setStoreId(undefined); setPage(1); };
  const searchRows = (next: string) => { setKeyword(next); setPage(1); };
  return { ...resultView, storeCatalogReady: directory.ready, storeOptions: directory.items.map((store) => ({ value: String(store.id), label: store.name })), catalogError: directory.loading || directory.ready ? undefined : '门店目录未完整加载，暂时无法按门店筛选。', listError: error ? '审计日志加载失败，请稍后重试。' : undefined, empty: !loading && !error && resultView.rows.length === 0, page, pageSize, keyword, result, resourceType, storeId, setPage, setPageSize, search: searchRows, changeResult, changeResource, changeStore, resetFilters, hasFilters: [keyword, result, resourceType, storeId].some(Boolean) };
}

function auditResult(data: FranchiseAuditLogsQuery | undefined) {
  return { rows: data?.auditLogs?.data ?? [], total: data?.auditLogs?.total ?? 0 };
}
