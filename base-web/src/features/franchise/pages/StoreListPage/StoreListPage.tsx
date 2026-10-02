import { useState } from 'react';
import { useMutation, useQuery } from '@apollo/client/react';
import { useNavigate } from 'react-router';
import type { FranchiseStoresQuery, StoreLifecycle } from '@/__generated__/graphql';
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';
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
import { storeLifecycleOptions } from '@/features/admin/config/filterOptions';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { useAuthStore } from '@/features/auth/store/authStore';
import { DELETE_FRANCHISE_STORES_MUTATION, FRANCHISE_STORES_QUERY } from '../../graphql/stores';
import { StoreTable, type FranchiseStoreRow } from './StoreTable';
import { storeStats } from './storeStats';
import { useStoreSubmission } from './useStoreSubmission';
import { franchiseStoreFilter } from '../franchiseListFilters';
import { initialSearchKeyword, queryKeyword, viewerOrganizationId, viewerPermissions } from '../pageState';

export function StoreListPage() {
  const state = useStorePage();
  return <div>
    <AdminPageHeader title="我的门店" description="门店先保存为草稿，再提交总部审核。" />
    <AdminStatsStrip items={storeStats(state.rows, state.total)} />
    <AdminBulkBar count={0} onClear={() => undefined} />
    <AdminActionError message={state.listError ?? (state.pendingDelete || state.pendingSubmit ? undefined : state.actionError)} />
    <AdminTableShell toolbar={<StoreFilters state={state} />} isEmpty={state.empty} empty="暂无门店"
      pagination={<AdminPagination total={state.total} page={state.page} pageSize={state.pageSize}
        onPageChange={state.setPage} onPageSizeChange={state.setPageSize} />}>
      <StoreTable rows={state.rows} organizationId={state.organizationId} permissions={state.permissions}
        onEdit={state.edit} onSubmit={state.openSubmit} onDelete={state.openDelete} />
    </AdminTableShell>
    <DeleteStoreDialog state={state} />
    <SubmitStoreDialog state={state} />
  </div>;
}

type PageState = ReturnType<typeof useStorePage>;

function StoreFilters({ state }: { state: PageState }) {
  return <AdminToolbar rightSlot={state.canCreate
    ? <AdminPrimaryActionButton onClick={() => state.edit()}>新增</AdminPrimaryActionButton> : null}>
    <AdminListFilters hasActiveFilters={state.hasFilters} onReset={state.resetFilters}>
      <AdminListSearch value={state.keyword} placeholder="搜索门店" onSearch={state.search} />
      <AdminFilterSelect value={state.lifecycle} placeholder="全部准入状态"
        options={storeLifecycleOptions} onValueChange={state.changeLifecycle} />
    </AdminListFilters>
  </AdminToolbar>;
}

function DeleteStoreDialog({ state }: { state: PageState }) {
  return <AlertDialog open={Boolean(state.pendingDelete)} onOpenChange={(open) => { if (!open && !state.deleting) state.closeDelete(); }}>
    <AlertDialogContent><AlertDialogHeader><AlertDialogTitle>删除门店</AlertDialogTitle>
      <AlertDialogDescription>确认删除“{state.pendingDelete?.name}”？删除后不会在门店列表显示。</AlertDialogDescription>
    </AlertDialogHeader><AdminActionError message={state.actionError} />
    <AlertDialogFooter><AlertDialogCancel disabled={state.deleting}>取消</AlertDialogCancel>
      <AlertDialogAction disabled={state.deleting} onClick={(event) => { event.preventDefault(); void state.remove(); }}>确认删除</AlertDialogAction>
    </AlertDialogFooter></AlertDialogContent>
  </AlertDialog>;
}

function SubmitStoreDialog({ state }: { state: PageState }) {
  return <AlertDialog open={Boolean(state.pendingSubmit)} onOpenChange={(open) => { if (!open) state.closeSubmit(); }}>
    <AlertDialogContent><AlertDialogHeader><AlertDialogTitle>提交门店审核</AlertDialogTitle>
      <AlertDialogDescription>确认提交“{state.pendingSubmit?.name}”给总部审核？提交后将无法编辑门店资料，请确认信息完整、准确。</AlertDialogDescription>
    </AlertDialogHeader><AdminActionError message={state.actionError} />
    <AlertDialogFooter><AlertDialogCancel disabled={state.submitting}>取消</AlertDialogCancel>
      <AlertDialogAction disabled={state.submitting} onClick={(event) => { event.preventDefault(); void state.confirmSubmit(); }}>
        {state.submitting ? '提交中…' : '确认提交'}
      </AlertDialogAction>
    </AlertDialogFooter></AlertDialogContent>
  </AlertDialog>;
}

function useStorePage() {
  const { search } = useAdminTab();
  const navigate = useNavigate();
  const viewer = useAuthStore((current) => current.viewer);
  const permissions = viewerPermissions(viewer);
  const organizationId = viewerOrganizationId(viewer);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState(() => initialSearchKeyword(search));
  const [lifecycle, setLifecycle] = useState<StoreLifecycle>();
  const [pendingDelete, setPendingDelete] = useState<FranchiseStoreRow>();
  const [actionError, setActionError] = useState<string>();
  const { data, loading, error, refetch } = useQuery(FRANCHISE_STORES_QUERY,
    { variables: { page, pageSize, q: queryKeyword(keyword), filter: franchiseStoreFilter(lifecycle) } });
  const submission = useStoreSubmission(refetch, setActionError);
  const [deleteStores, deleteState] = useMutation(DELETE_FRANCHISE_STORES_MUTATION);
  const { rows, total } = storeResult(data);
  const edit = (row?: FranchiseStoreRow) => navigate(`/admin/franchise/stores/manage${row ? `?id=${encodeURIComponent(row.id)}` : ''}`);
  const searchRows = (next: string) => { setKeyword(next); setPage(1); };
  const changeLifecycle = (next?: string) => { setLifecycle(next as StoreLifecycle | undefined); setPage(1); };
  const resetFilters = () => { setKeyword(''); setLifecycle(undefined); setPage(1); };
  const openDelete = (row: FranchiseStoreRow) => { setActionError(undefined); setPendingDelete(row); };
  const closeDelete = () => { setPendingDelete(undefined); setActionError(undefined); };
  async function remove() {
    if (!pendingDelete) return;
    const deleted = await runAdminAction(async () => { await deleteStores({ variables: { ids: [pendingDelete.id] } }); }, setActionError);
    if (!deleted) return;
    setPendingDelete(undefined);
    try { await refetch(); }
    catch { setActionError('门店已删除，但列表刷新失败，请刷新页面。'); }
  }
  return { rows, total, permissions, organizationId, page, pageSize, keyword, lifecycle, pendingDelete,
    actionError, listError: error ? '门店加载失败，请稍后重试。' : undefined,
    ...submission, setPage, setPageSize, edit, remove, openDelete, closeDelete, search: searchRows,
    changeLifecycle, resetFilters, hasFilters: [keyword, lifecycle].some(Boolean),
    canCreate: permissions.includes('store:create'), empty: !loading && !error && rows.length === 0,
    deleting: deleteState.loading };
}

function storeResult(data: FranchiseStoresQuery | undefined) {
  return { rows: data?.stores?.data ?? [], total: data?.stores?.total ?? 0 };
}
