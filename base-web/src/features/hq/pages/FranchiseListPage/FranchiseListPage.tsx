import { useState } from 'react';
import { useQuery } from '@apollo/client/react';
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';
import type { OrganizationStatus } from '@/__generated__/graphql';
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
import { organizationStatusOptions } from '@/features/admin/config/filterOptions';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { useAuthStore } from '@/features/auth/store/authStore';
import { HQ_FRANCHISES_QUERY } from '../../graphql/franchises';
import { FranchiseTable } from './FranchiseTable';
import type { FranchiseRow } from './FranchiseTable';
import { ConfirmInitialAccountDialog } from './ConfirmInitialAccountDialog';
import { franchiseStats } from './franchiseStats';
import { ProvisionFranchiseDialog } from './ProvisionFranchiseDialog';
import { useFranchiseListActions } from './useFranchiseListActions';
import { SuspendFranchiseDialog } from './SuspendFranchiseDialog';
import { hqFranchiseFilter } from '../hqListFilters';
import { PaymentConfigDrawer } from '../PaymentConfigPage/PaymentConfigDrawer';

export function FranchiseListPage() {
  const { search } = useAdminTab();
  const permissions = useAuthStore((state) => state.viewer?.permissions ?? []);
  const canResetPassword = permissions.includes('account:update');
  const canConfirmInitialAccount = permissions.includes('hqMembership:read');
  const [confirmingRow, setConfirmingRow] = useState<FranchiseRow>();
  const [resetRow, setResetRow] = useState<FranchiseRow>();
  const [suspendingRow, setSuspendingRow] = useState<FranchiseRow>();
  const [paymentRow, setPaymentRow] = useState<FranchiseRow>();
  const [page, setPage] = useState(1); const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState(() => new URLSearchParams(search).get('q') ?? ''); const [status, setStatus] = useState<OrganizationStatus>();
  const { data, loading, error, refetch } = useQuery(HQ_FRANCHISES_QUERY, { variables: { page, pageSize, q: keyword || null, filter: hqFranchiseFilter(status), canResetPassword } });
  const { open, values, outcome, provisionError, actionError, temporaryPasswords, resettingOrganizationId, changingStatusOrganizationId, provisionState, setOpen, changeValues, changeOpen, submit, suspendRow, restoreRow, resetPasswordRow } = useFranchiseListActions(refetch);
  const rows = data?.organizations?.data ?? [];
  const total = data?.organizations?.total ?? 0;

  const changeStatusFilter = (next?: string) => { setStatus(next as OrganizationStatus | undefined); setPage(1); }; const resetFilters = () => { setKeyword(''); setStatus(undefined); setPage(1); };
  const view = franchiseListView(loading, error, rows.length, actionError);

  return (
    <div><AdminPageHeader title="加盟商" description="总部只执行开通、暂停与恢复，不代替加盟商经营。" /><AdminStatsStrip items={franchiseStats(rows, total)} /><AdminBulkBar count={0} onClear={() => undefined} />
			<AdminActionError message={view.errorMessage} />
		<AdminTableShell toolbar={<AdminToolbar rightSlot={permissions.includes('franchise:provision') ? <AdminPrimaryActionButton onClick={() => setOpen(true)}>开通加盟商</AdminPrimaryActionButton> : null}><AdminListFilters hasActiveFilters={[keyword, status].some(Boolean)} onReset={resetFilters}><AdminListSearch value={keyword} placeholder="搜索加盟商" onSearch={(next) => { setKeyword(next); setPage(1); }} /><AdminFilterSelect value={status} placeholder="全部组织状态" options={organizationStatusOptions} onValueChange={changeStatusFilter} /></AdminListFilters></AdminToolbar>} isEmpty={view.empty} empty="暂无加盟商" pagination={<AdminPagination total={total} page={page} pageSize={pageSize} onPageChange={setPage} onPageSizeChange={setPageSize} />}>
        <FranchiseTable rows={rows} canViewStores={permissions.includes('store:read_all')} canSuspend={permissions.includes('organization:suspend')} canRestore={permissions.includes('organization:restore')} canResetPassword={canResetPassword} canConfirmInitialAccount={canConfirmInitialAccount} resettingOrganizationId={resettingOrganizationId} changingStatusOrganizationId={changingStatusOrganizationId} temporaryPasswords={temporaryPasswords} onSuspend={setSuspendingRow} onRestore={(row) => void restoreRow(row)} onResetPassword={setResetRow} onConfirmInitialAccount={setConfirmingRow} onConfigurePayment={paymentAction(permissions, setPaymentRow)} />
      </AdminTableShell>
      <ProvisionFranchiseDialog open={open} onOpenChange={changeOpen} values={values} onValuesChange={changeValues} onSubmit={submit} isSubmitting={provisionState.loading} outcome={outcome} error={provisionError} />
      {confirmingRow && <ConfirmInitialAccountDialog row={confirmingRow} onClose={() => setConfirmingRow(undefined)} onConfirmed={refetch} />}
      <SuspendDialogSlot row={suspendingRow} pending={changingStatusOrganizationId !== undefined} error={actionError} onConfirm={suspendRow} onClose={() => setSuspendingRow(undefined)} />
      <ResetInitialPasswordDialog row={resetRow} onClose={() => setResetRow(undefined)} onConfirm={(row) => { setResetRow(undefined); void resetPasswordRow(row); }} />
      <PaymentDrawerSlot row={paymentRow} onClose={() => setPaymentRow(undefined)} />
    </div>
  );
}

function paymentAction(permissions: string[], setPaymentRow: (row: FranchiseRow) => void) {
  return permissions.includes('paymentConfig:read') ? setPaymentRow : undefined;
}

function PaymentDrawerSlot({ row, onClose }: { row?: FranchiseRow; onClose(): void }) {
  if (!row) return null;
  return <PaymentConfigDrawer target={{ scope: 'FRANCHISE', organizationId: row.id, name: row.name }} onClose={onClose} />;
}

/** 弹窗只在选中待暂停加盟商后挂载，关闭时清空其原因草稿。 */
function SuspendDialogSlot({ row, pending, error, onConfirm, onClose }: { row?: FranchiseRow; pending: boolean; error?: string; onConfirm(row: FranchiseRow, reason: string): Promise<boolean>; onClose(): void }) {
  if (!row) return null;
  return <SuspendFranchiseDialog row={row} pending={pending} error={error} onConfirm={(reason) => onConfirm(row, reason)} onClose={onClose} />;
}

function ResetInitialPasswordDialog({ row, onClose, onConfirm }: { row?: FranchiseRow; onClose(): void; onConfirm(row: FranchiseRow): void }) {
  return <AlertDialog open={Boolean(row)} onOpenChange={(open) => { if (!open) onClose(); }}>
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>重置初始账号密码</AlertDialogTitle>
        <AlertDialogDescription>确认重置「{row?.name}」的初始账号密码？该账号在所有加盟商工作台的在线会话都会失效。</AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel>取消</AlertDialogCancel>
        <AlertDialogAction onClick={() => { if (row) onConfirm(row); }}>确认重置</AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>;
}

function franchiseListView(loading: boolean, error: unknown, rowCount: number, actionError?: string) {
  return {
    errorMessage: error ? '加盟商加载失败，请稍后重试。' : actionError,
    empty: !loading && !error && rowCount === 0,
  };
}
