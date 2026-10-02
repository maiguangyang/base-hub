import { useState } from 'react';
import { useQuery } from '@apollo/client/react';
import { Link } from 'react-router';
import type { HqFranchiseStoresQuery, StoreBusinessStatus, StoreLifecycle } from '@/__generated__/graphql';
import { Button } from '@/components/ui/button';
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
import { storeBusinessStatusOptions, storeLifecycleOptions } from '@/features/admin/config/filterOptions';
import { useCompleteFilterCatalog } from '@/features/admin/lib/useCompleteFilterCatalog';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { useAuthStore } from '@/features/auth/store/authStore';
import { HQ_FRANCHISE_STORES_QUERY } from '../../graphql/franchiseStores';
import { HQ_FRANCHISES_QUERY } from '../../graphql/franchises';
import { hqFranchiseFilter, hqFranchiseStoreFilter } from '../hqListFilters';
import { franchiseStoreStats } from './franchiseStoreStats';
import { FranchiseSearchSelect } from './FranchiseSearchSelect';
import { FranchiseStoreTable, type FranchiseStoreRow } from './FranchiseStoreTable';
import { franchiseStoreViewKey, initialFranchiseStoreKeyword, selectedFranchiseId } from './pageState';
import { PaymentConfigDrawer } from '../PaymentConfigPage/PaymentConfigDrawer';

export function FranchiseStoreListPage() {
  const { search } = useAdminTab();
  return <FranchiseStoreView key={franchiseStoreViewKey(search)} search={search} />;
}

function FranchiseStoreView({ search }: { search: string }) {
  const state = useFranchiseStorePage(search);
  const canReadPayment = useAuthStore((auth) => auth.viewer?.permissions.includes('paymentConfig:read') ?? false);
  const [paymentRow, setPaymentRow] = useState<FranchiseStoreRow>();
  const title = state.organizationName ? `${state.organizationName} · 门店` : '加盟门店';
  const pagination = state.loadFailed ? undefined : <AdminPagination total={state.total} page={state.page} pageSize={state.pageSize} onPageChange={state.setPage} onPageSizeChange={state.setPageSize} />;
  return (
    <div>
      <AdminPageHeader title={title} description="总部只读查看加盟商门店；门店审核仍在“门店审核”中处理。" actions={<Button asChild variant="outline"><Link to="/admin/hq/franchises">返回加盟商</Link></Button>} />
      <AdminStatsStrip items={franchiseStoreStats(state.total)} />
      <AdminBulkBar count={0} onClear={() => undefined} />
      <AdminActionError message={state.catalogError} />
      <AdminTableShell toolbar={<AdminToolbar><AdminListFilters hasActiveFilters={state.hasFilters} onReset={state.resetFilters}><AdminListSearch value={state.keyword} placeholder="搜索加盟门店" onSearch={state.search} /><FranchiseSearchSelect value={state.organizationId} options={state.franchiseOptions} disabled={!state.catalogReady} onValueChange={state.changeOrganization} /><AdminFilterSelect value={state.lifecycle} placeholder="全部准入状态" options={storeLifecycleOptions} onValueChange={state.changeLifecycle} /><AdminFilterSelect value={state.businessStatus} placeholder="全部营业状态" options={storeBusinessStatusOptions} onValueChange={state.changeBusinessStatus} /></AdminListFilters></AdminToolbar>} isEmpty={state.empty || state.loadFailed} empty={state.loadFailed ? <FranchiseStoreLoadError onRetry={state.retry} /> : state.emptyMessage} pagination={pagination}>
        <FranchiseStoreTable rows={state.rows} onConfigurePayment={canReadPayment ? setPaymentRow : undefined} />
      </AdminTableShell>
      {paymentRow && <PaymentConfigDrawer target={{ scope: 'STORE', storeId: paymentRow.id, name: paymentRow.name }} onClose={() => setPaymentRow(undefined)} />}
    </div>
  );
}

export function FranchiseStoreLoadError({ onRetry }: { onRetry: () => void }) {
  return (
    <div role="alert" className="flex flex-col items-center gap-3">
      <p>门店加载失败，请稍后重试。</p>
      <Button type="button" variant="outline" size="sm" onClick={onRetry}>重试</Button>
    </div>
  );
}

function useFranchiseStorePage(search: string) {
  const initialOrganizationId = selectedFranchiseId(search);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState(() => initialFranchiseStoreKeyword(search));
  const [organizationId, setOrganizationId] = useState(initialOrganizationId);
  const [lifecycle, setLifecycle] = useState<StoreLifecycle>();
  const [businessStatus, setBusinessStatus] = useState<StoreBusinessStatus>();
  const variables = { page, pageSize, q: keyword || null, filter: hqFranchiseStoreFilter(organizationId, lifecycle, businessStatus) };
  const { data, loading, error, refetch } = useQuery(HQ_FRANCHISE_STORES_QUERY, { variables });
  const catalog = useQuery(HQ_FRANCHISES_QUERY, { variables: { page: 1, pageSize: 200, q: null, filter: hqFranchiseFilter() } });
  const { rows, total } = franchiseStoreResult(data);
  const directory = useCompleteFilterCatalog('hq-franchises', franchisePage(catalog.data), catalog.loading, catalog.error, async (nextPage) => {
    const next = await catalog.fetchMore({ variables: { page: nextPage } });
    const pageResult = franchisePage(next.data);
    if (!pageResult) throw new Error('catalog unavailable');
    return pageResult;
  });
  const franchises = directory.items;
  const catalogReady = directory.ready;
  const searchRows = (next: string) => { setKeyword(next); setPage(1); };
  const changeOrganization = (next?: string) => { setOrganizationId(next); setPage(1); };
  const changeLifecycle = (next?: string) => { setLifecycle(next as StoreLifecycle | undefined); setPage(1); };
  const changeBusinessStatus = (next?: string) => { setBusinessStatus(next as StoreBusinessStatus | undefined); setPage(1); };
  const resetFilters = () => { setKeyword(''); setOrganizationId(undefined); setLifecycle(undefined); setBusinessStatus(undefined); setPage(1); };
  const retry = () => { void refetch().catch(() => undefined); };
  return {
    rows, total, page, pageSize, keyword, organizationId, lifecycle, businessStatus, setPage, setPageSize, search: searchRows,
    organizationName: franchises.find((item) => item.id === organizationId)?.name,
    franchiseOptions: franchises.map((item) => ({ value: String(item.id), label: item.name })),
    catalogReady, catalogError: franchiseCatalogError(directory.loading, catalogReady),
    changeOrganization, changeLifecycle, changeBusinessStatus, resetFilters,
    hasFilters: [keyword, organizationId, lifecycle, businessStatus].some(Boolean),
    empty: !loading && rows.length === 0,
    loadFailed: Boolean(error), retry,
    emptyMessage: organizationId ? '该加盟商暂无门店' : '暂无加盟门店',
  };
}

function franchiseCatalogError(loading: boolean, ready: boolean): string | undefined {
  if (loading || ready) return undefined;
  return '加盟商目录未完整加载，暂时无法按加盟商筛选。';
}

function franchiseStoreResult(data: HqFranchiseStoresQuery | undefined) {
  return { rows: data?.stores?.data ?? [], total: data?.stores?.total ?? 0 };
}

function franchisePage(data: { organizations?: { data: Array<{ id: string; name: string }>; total: number } | null } | undefined) {
  const result = data?.organizations;
  return result ? { items: result.data, total: result.total } : undefined;
}
