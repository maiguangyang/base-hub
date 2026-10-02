import { useState } from 'react';
import { useMutation, useQuery } from '@apollo/client/react';
import type { HqStoreApprovalsQuery } from '@/__generated__/graphql';
import { AdminBulkBar } from '@/features/admin/components/AdminBulkBar';
import { AdminActionError } from '@/features/admin/components/AdminFormDialogShell';
import { AdminListFilters } from '@/features/admin/components/AdminListFilters';
import { AdminListSearch } from '@/features/admin/components/AdminListSearch';
import { AdminPageHeader } from '@/features/admin/components/AdminPageHeader';
import { AdminPagination } from '@/features/admin/components/AdminPagination';
import { AdminStatsStrip } from '@/features/admin/components/AdminStatsStrip';
import { AdminTableShell } from '@/features/admin/components/AdminTableShell';
import { AdminToolbar } from '@/features/admin/components/AdminToolbar';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { useAuthStore } from '@/features/auth/store/authStore';
import { catalogPage } from '@/features/admin/lib/filterCatalog';
import { useCompleteFilterCatalog } from '@/features/admin/lib/useCompleteFilterCatalog';
import { HQ_FRANCHISES_QUERY } from '../../graphql/franchises';
import { HQ_STORE_APPROVALS_QUERY, REVIEW_STORE_MUTATION } from '../../graphql/storeApprovals';
import { hqFranchiseFilter, hqStoreApprovalFilter } from '../hqListFilters';
import { FranchiseSearchSelect } from '../FranchiseStoreListPage/FranchiseSearchSelect';
import { availableReviewActions, ReviewStoreDialog, reviewStoreInput } from './ReviewStoreDialog';
import { StoreApprovalTable, type StoreApprovalRow } from './StoreApprovalTable';
import { storeApprovalStats } from './storeApprovalStats';

export function StoreApprovalListPage() {
  const state = useStoreApprovalPage();
  return (
    <div><AdminPageHeader title="门店审核" description="独立批准或退回加盟商提交的门店。" /><AdminStatsStrip items={storeApprovalStats(state.total)} /><AdminBulkBar count={0} onClear={() => undefined} />
      <AdminActionError message={state.listError ?? state.filterError} />
      <AdminTableShell toolbar={<AdminToolbar><AdminListFilters hasActiveFilters={state.hasFilters} onReset={state.resetFilters}><AdminListSearch value={state.keyword} placeholder="搜索待审门店" onSearch={state.search} /><FranchiseSearchSelect value={state.organizationId} options={state.franchiseOptions} disabled={!state.catalogReady} onValueChange={state.changeOrganization} /></AdminListFilters></AdminToolbar>} isEmpty={state.empty} empty="暂无待审核门店" pagination={<AdminPagination total={state.total} page={state.page} pageSize={state.pageSize} onPageChange={state.setPage} onPageSizeChange={state.setPageSize} />}>
        <StoreApprovalTable rows={state.rows} {...state.actions} onReview={state.selectReview} />
      </AdminTableShell>
      {state.selection ? <ReviewStoreDialog open onOpenChange={state.changeDialogOpen} storeId={state.selection.row.id} storeName={state.selection.row.name} approved={state.selection.approved} reason={state.reason} onReasonChange={state.setReason} onSubmit={state.submitReview} isSubmitting={state.reviewing} /> : null}
    </div>
  );
}

function useStoreApprovalPage() {
  const { search } = useAdminTab();
  const permissions = useAuthStore((state) => state.viewer?.permissions ?? []);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState(() => new URLSearchParams(search).get('q') ?? '');
  const [organizationId, setOrganizationId] = useState<string>();
  const [selection, setSelection] = useState<{ row: StoreApprovalRow; approved: boolean }>();
  const [reason, setReason] = useState('');
  const { data, loading, error, refetch } = useQuery(HQ_STORE_APPROVALS_QUERY, { variables: { page, pageSize, q: queryKeyword(keyword), filter: hqStoreApprovalFilter(organizationId) } });
  const catalog = useQuery(HQ_FRANCHISES_QUERY, { variables: { page: 1, pageSize: 200, q: null, filter: hqFranchiseFilter() } });
  const [review, reviewState] = useMutation(REVIEW_STORE_MUTATION);
  const { rows, total } = approvalResult(data);
  const actions = availableReviewActions(permissions);
  const directory = useCompleteFilterCatalog('hq-franchises', catalogPage(catalog.data?.organizations), catalog.loading, catalog.error, async (nextPage) => {
    const next = await catalog.fetchMore({ variables: { page: nextPage } });
    const result = catalogPage(next.data?.organizations);
    if (!result) throw new Error('catalog unavailable');
    return result;
  });

  async function submitReview() {
    if (!selection) return;
    await review({ variables: { input: reviewStoreInput(selection.row.id, selection.approved, reason) } });
    setSelection(undefined);
    setReason('');
    await refetch();
  }
  const changeOrganization = (next?: string) => { setOrganizationId(next); setPage(1); };
  const resetFilters = () => { setKeyword(''); setOrganizationId(undefined); setPage(1); };
  const searchRows = (next: string) => { setKeyword(next); setPage(1); };
  return {
    page, pageSize, keyword, organizationId, rows, total, selection, reason, actions,
    catalogReady: directory.ready, filterError: directory.loading || directory.ready ? undefined : '加盟商目录未完整加载，暂时无法按加盟商筛选。', listError: error ? '待审核门店加载失败，请稍后重试。' : undefined, franchiseOptions: directory.items.map((item) => ({ value: String(item.id), label: item.name })),
    setPage, setPageSize, setReason, changeOrganization, resetFilters, search: searchRows,
    hasFilters: [keyword, organizationId].some(Boolean), empty: approvalEmpty(loading, Boolean(error), rows.length),
    selectReview: (row: StoreApprovalRow, approved: boolean) => { setReason(''); setSelection({ row, approved }); },
    changeDialogOpen: (open: boolean) => { if (!open && !reviewState.loading) setSelection(undefined); },
    submitReview, reviewing: reviewState.loading,
  };
}

function approvalResult(data: HqStoreApprovalsQuery | undefined) {
  return { rows: data?.stores?.data ?? [], total: data?.stores?.total ?? 0 };
}

function queryKeyword(keyword: string): string | null { return keyword === '' ? null : keyword; }
function approvalEmpty(loading: boolean, error: boolean, rowCount: number): boolean { return !loading && !error && rowCount === 0; }
