import { useState } from 'react';
import { useNavigate } from 'react-router';
import { Button } from '@/components/ui/button';
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';
import { AdminBulkBar } from '@/features/admin/components/AdminBulkBar';
import { AdminFilterSelect } from '@/features/admin/components/AdminFilterSelect';
import { AdminActionError } from '@/features/admin/components/AdminFormDialogShell';
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
import { FranchiseSearchSelect } from '../FranchiseStoreListPage/FranchiseSearchSelect';
import { PaymentConfigDrawer } from '../PaymentConfigPage/PaymentConfigDrawer';
import { HqStoreTable, type HqStoreRow } from './HqStoreTable';
import { useHqStorePage } from './useHqStorePage';

type PageState = ReturnType<typeof useHqStorePage>;
const typeOptions = [{ value: 'HEADQUARTERS', label: '直营' }, { value: 'FRANCHISE', label: '加盟' }];

export function HqStoreListPage() {
  const { search } = useAdminTab();
  return <HqStoreView key={search} search={search} />;
}

function HqStoreView({ search }: { search: string }) {
  const navigate = useNavigate();
  const state = useHqStorePage(search);
  const edit = (row: HqStoreRow) => navigate(`/admin/hq/stores/manage?id=${encodeURIComponent(row.id)}`);
  const [paymentRow, setPaymentRow] = useState<HqStoreRow>();
  const canReadPayment = state.permissions.includes('paymentConfig:read');
  return <div>
    <AdminPageHeader title="门店管理" description="统一查看直营与加盟门店；总部仅管理自己经营的门店。" />
    <AdminStatsStrip items={[{ key: 'total', label: '筛选结果', value: String(state.total), icon: 'list-checks' }]} />
    <AdminBulkBar count={0} onClear={() => undefined} />
    <HqStoreErrors state={state} />
    {state.loadFailed && <Button variant="outline" onClick={state.retry}>重试门店列表</Button>}
    <AdminTableShell toolbar={<StoreFilters state={state} />}
      isEmpty={!state.loading && (state.rows.length === 0 || state.loadFailed)} empty={state.loadFailed ? '门店加载失败' : '暂无门店'}
      pagination={state.loadFailed ? undefined : <AdminPagination total={state.total} page={state.page} pageSize={state.pageSize}
        onPageChange={state.setPage} onPageSizeChange={state.setPageSize} />}>
      <HqStoreTable rows={state.rows} hqOrganizationId={state.hqOrganizationId}
        canUpdate={state.permissions.includes('hqStore:update')} canDelete={state.permissions.includes('hqStore:delete')}
        statusBusy={state.saving} onEdit={edit} onDelete={state.setPendingDelete}
        onBusinessStatus={(row, open) => { void state.toggleBusinessStatus(row, open); }}
        onConfigurePayment={canReadPayment ? setPaymentRow : undefined} />
    </AdminTableShell>
    <DeleteStoreDialog state={state} />
    {paymentRow && <PaymentConfigDrawer target={{ scope: 'STORE', storeId: paymentRow.id, name: paymentRow.name }}
      onClose={() => setPaymentRow(undefined)} />}
  </div>;
}

function HqStoreErrors({ state }: { state: PageState }) {
  const deleteError = state.pendingDelete && state.actionError?.startsWith('删除门店失败');
  const pageError = state.loadFailed ? '门店加载失败，请重试。' : deleteError ? undefined : state.actionError;
  const catalogError = state.permissions.includes('organization:read') && state.catalog.error
    ? '加盟商目录加载失败，暂时无法按加盟商筛选。' : undefined;
  return <><AdminActionError message={pageError} /><AdminActionError message={catalogError} /></>;
}

function StoreFilters({ state }: { state: PageState }) {
  const navigate = useNavigate();
  const showOrganization = state.permissions.includes('organization:read');
  const primaryAction = state.permissions.includes('hqStore:create')
    ? <AdminPrimaryActionButton onClick={() => navigate('/admin/hq/stores/manage')}>新增直营店</AdminPrimaryActionButton> : null;
  return <AdminToolbar rightSlot={primaryAction}><AdminListFilters hasActiveFilters={state.hasFilters} onReset={state.resetFilters}>
    <AdminListSearch value={state.keyword} placeholder="搜索门店" onSearch={state.changeKeyword} />
    <AdminFilterSelect value={state.type} placeholder="全部类型" options={typeOptions} onValueChange={state.changeType} />
    {showOrganization && <FranchiseSearchSelect value={state.organizationId} options={state.catalog.options}
      disabled={state.type === 'HEADQUARTERS' || !state.catalog.ready} onValueChange={state.changeOrganization} />}
    <AdminFilterSelect value={state.lifecycle} placeholder="全部准入状态" options={storeLifecycleOptions} onValueChange={state.changeLifecycle} />
    <AdminFilterSelect value={state.businessStatus} placeholder="全部营业状态" options={storeBusinessStatusOptions} onValueChange={state.changeBusinessStatus} />
    {showOrganization && Boolean(state.catalog.error) && <Button variant="outline" onClick={state.catalog.retry}>重试加盟商目录</Button>}
  </AdminListFilters></AdminToolbar>;
}

function DeleteStoreDialog({ state }: { state: PageState }) {
  return <AlertDialog open={Boolean(state.pendingDelete)} onOpenChange={(open) => { if (!open && !state.saving) state.closeDelete(); }}>
    <AlertDialogContent><AlertDialogHeader><AlertDialogTitle>删除直营门店</AlertDialogTitle>
      <AlertDialogDescription>确认删除“{state.pendingDelete?.name}”？删除后不会在门店列表显示。</AlertDialogDescription>
    </AlertDialogHeader><AdminActionError message={state.actionError?.startsWith('删除门店失败') ? state.actionError : undefined} />
    <AlertDialogFooter><AlertDialogCancel disabled={state.saving}>取消</AlertDialogCancel>
      <AlertDialogAction disabled={state.saving} onClick={(event) => { event.preventDefault(); void state.remove(); }}>确认删除</AlertDialogAction>
    </AlertDialogFooter></AlertDialogContent>
  </AlertDialog>;
}
