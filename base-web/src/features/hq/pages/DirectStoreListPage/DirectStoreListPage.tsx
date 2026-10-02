import { useState } from 'react';
import { useMutation, useQuery } from '@apollo/client/react';
import { useNavigate } from 'react-router';
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';
import { storeBusinessStatusInput } from '@/features/admin/pages/StoreEditorPage/storeInputs';
import { AdminBulkBar } from '@/features/admin/components/AdminBulkBar';
import { AdminFilterSelect } from '@/features/admin/components/AdminFilterSelect';
import { AdminActionError, runAdminAction } from '@/features/admin/components/AdminFormDialogShell';
import { AdminListFilters } from '@/features/admin/components/AdminListFilters';
import { AdminListSearch } from '@/features/admin/components/AdminListSearch';
import { AdminPageHeader } from '@/features/admin/components/AdminPageHeader';
import { AdminPagination } from '@/features/admin/components/AdminPagination';
import { AdminPrimaryActionButton } from '@/features/admin/components/AdminPrimaryActionButton';
import { AdminStatsStrip } from '@/features/admin/components/AdminStatsStrip';
import { AdminTableShell } from '@/features/admin/components/AdminTableShell';
import { AdminToolbar } from '@/features/admin/components/AdminToolbar';
import { storeBusinessStatusOptions, storeLifecycleOptions } from '@/features/admin/config/filterOptions';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { useAuthStore } from '@/features/auth/store/authStore';
import type { ViewerSummary } from '@/features/auth/store/authStore';
import type { HqDirectStoresQuery, StoreBusinessStatus, StoreLifecycle } from '@/__generated__/graphql';
import { DELETE_DIRECT_STORES_MUTATION, HQ_DIRECT_STORES_QUERY, UPDATE_DIRECT_STORE_MUTATION } from '../../graphql/directStores';
import { canDeleteDirectStore, DirectStoreTable, type DirectStoreRow } from './DirectStoreTable';
import { directStoreStats } from './directStoreStats';
import { hqDirectStoreFilter } from '../hqListFilters';

export function DirectStoreListPage() {
  const state = useDirectStorePage();
  const primaryAction = state.canCreate
    ? <AdminPrimaryActionButton onClick={() => state.edit()}>新增</AdminPrimaryActionButton>
    : null;
  return (
    <div><AdminPageHeader title="直营门店" description="管理总部组织直接经营的门店。" /><AdminStatsStrip items={directStoreStats(state.rows, state.total)} />
      <AdminBulkBar count={0} onClear={() => undefined} />
            <AdminActionError message={state.listError ?? (state.pendingDelete ? undefined : state.actionError)} />
      <AdminTableShell toolbar={<AdminToolbar rightSlot={primaryAction}><AdminListFilters hasActiveFilters={state.hasFilters} onReset={state.resetFilters}><AdminListSearch value={state.keyword} placeholder="搜索门店" onSearch={state.search} /><AdminFilterSelect value={state.lifecycle} placeholder="全部准入状态" options={storeLifecycleOptions} onValueChange={state.changeLifecycle} /><AdminFilterSelect value={state.businessStatus} placeholder="全部营业状态" options={storeBusinessStatusOptions} onValueChange={state.changeBusinessStatus} /></AdminListFilters></AdminToolbar>} isEmpty={state.empty} empty="暂无直营门店" pagination={<AdminPagination total={state.total} page={state.page} pageSize={state.pageSize} onPageChange={state.setPage} onPageSizeChange={state.setPageSize} />}>
        <DirectStoreTable rows={state.rows} canUpdate={state.canUpdate} canDelete={state.canDelete} onEdit={state.edit} onDelete={state.openDelete} onBusinessStatus={(row, open) => void state.toggleBusinessStatus(row, open)} statusBusy={state.saving} />
      </AdminTableShell>
      <DeleteStoreDialog state={state} />
    </div>
  );
}

type PageState = ReturnType<typeof useDirectStorePage>;

function DeleteStoreDialog({ state }: { state: PageState }) {
  return <AlertDialog open={Boolean(state.pendingDelete)} onOpenChange={(open) => { if (!open && !state.deleting) state.closeDelete(); }}>
    <AlertDialogContent><AlertDialogHeader><AlertDialogTitle>删除直营门店</AlertDialogTitle>
      <AlertDialogDescription>确认删除“{state.pendingDelete?.name}”？删除后不会在门店列表显示。</AlertDialogDescription>
    </AlertDialogHeader><AdminActionError message={state.actionError} />
    <AlertDialogFooter><AlertDialogCancel disabled={state.deleting}>取消</AlertDialogCancel>
      <AlertDialogAction disabled={state.deleting} onClick={(event) => { event.preventDefault(); void state.remove(); }}>确认删除</AlertDialogAction>
    </AlertDialogFooter></AlertDialogContent>
  </AlertDialog>;
}

function useDirectStorePage() {
  const { search } = useAdminTab();
  const navigate = useNavigate();
  const viewer = useAuthStore((state) => state.viewer);
  const { permissions, organizationId } = directStoreContext(viewer);
  const [page, setPage] = useState(1); const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState(() => initialKeyword(search));
  const [lifecycle, setLifecycle] = useState<StoreLifecycle>(); const [businessStatus, setBusinessStatus] = useState<StoreBusinessStatus>();
  const [pendingDelete, setPendingDelete] = useState<DirectStoreRow>();
  const [actionError, setActionError] = useState<string>();
  const variables = { page, pageSize, q: keyword === '' ? null : keyword, filter: hqDirectStoreFilter(organizationId, lifecycle, businessStatus) };
  const { data, loading, error, refetch } = useQuery(HQ_DIRECT_STORES_QUERY, { variables, skip: !organizationId });
  const [updateStore, updateState] = useMutation(UPDATE_DIRECT_STORE_MUTATION);
  const [deleteStores, deleteState] = useMutation(DELETE_DIRECT_STORES_MUTATION);
  const { rows, total } = directStoreResult(data);

  const edit = (row?: DirectStoreRow) => navigate(`/admin/hq/stores/manage${row ? `?id=${encodeURIComponent(row.id)}` : ''}`);
  const openDelete = (row: DirectStoreRow) => { setActionError(undefined); setPendingDelete(row); };
  const closeDelete = () => { setActionError(undefined); setPendingDelete(undefined); };
  async function remove() {
    if (!pendingDelete) return;
    const deleted = await runAdminAction(async () => { await deleteStores({ variables: { ids: [pendingDelete.id] } }); }, setActionError);
    if (!deleted) return;
    setPendingDelete(undefined);
    try { await refetch(); }
    catch { setActionError('门店已删除，但列表刷新失败，请刷新页面。'); }
  }
  async function toggleBusinessStatus(row: DirectStoreRow, open: boolean) {
    await runAdminAction(async () => {
      await updateStore({ variables: { id: row.id, input: storeBusinessStatusInput(open) } });
      await refetch();
    }, setActionError);
  }
  const searchRows = (next: string) => { setKeyword(next); setPage(1); };
  const changeLifecycle = (next?: string) => { setLifecycle(next as StoreLifecycle | undefined); setPage(1); };
  const changeBusinessStatus = (next?: string) => { setBusinessStatus(next as StoreBusinessStatus | undefined); setPage(1); };
  const resetFilters = () => { setKeyword(''); setLifecycle(undefined); setBusinessStatus(undefined); setPage(1); };
  return {
    rows, total, page, pageSize, keyword, lifecycle, businessStatus, organizationId, pendingDelete, actionError,
    listError: error ? '直营门店加载失败，请稍后重试。' : undefined,
    setPage, setPageSize, edit, remove, openDelete, closeDelete, toggleBusinessStatus, search: searchRows,
    changeLifecycle, changeBusinessStatus, resetFilters, hasFilters: [keyword, lifecycle, businessStatus].some(Boolean),
    canCreate: permissions.includes('hqStore:create'),
    canUpdate: permissions.includes('hqStore:update'),
    canDelete: canDeleteDirectStore(permissions),
    empty: !loading && !error && rows.length === 0,
    saving: updateState.loading || deleteState.loading, deleting: deleteState.loading,
  };
}

function directStoreContext(viewer: ViewerSummary | null) {
  return {
    permissions: viewer ? viewer.permissions : [],
    organizationId: viewer?.currentWorkspace?.organizationId ?? '',
  };
}

function initialKeyword(search: string): string {
  return new URLSearchParams(search).get('q') ?? '';
}

function directStoreResult(data: HqDirectStoresQuery | undefined) {
  return { rows: data?.stores?.data ?? [], total: data?.stores?.total ?? 0 };
}
