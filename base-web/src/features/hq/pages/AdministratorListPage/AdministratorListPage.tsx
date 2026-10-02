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
import { membershipStatusOptions } from '@/features/admin/config/filterOptions';
import { AdministratorFormDialog } from './AdministratorFormDialog';
import { AdministratorTable } from './AdministratorTable';
import { administratorStats } from './administratorStats';
import { useAdministratorPage } from './useAdministratorPage';

export function AdministratorListPage() {
  const state = useAdministratorPage();
  return <div><AdminPageHeader title="管理员" description="管理总部后台账号、角色和登录状态。" /><AdminStatsStrip items={administratorStats(state.rows, state.total)} /><AdminBulkBar count={0} onClear={() => undefined} />
    <AdminActionError message={state.listError ?? state.actionError ?? state.filterCatalogError} />
    <AdminTableShell toolbar={<AdminToolbar rightSlot={state.canCreate ? <AdminPrimaryActionButton onClick={() => state.edit()}>新增</AdminPrimaryActionButton> : null}><AdminListFilters hasActiveFilters={state.hasFilters} onReset={state.resetFilters}><AdminListSearch value={state.keyword} placeholder="搜索管理员" onSearch={state.search} /><AdminFilterSelect value={state.status} placeholder="全部成员状态" options={membershipStatusOptions} onValueChange={state.changeStatusFilter} /><AdminFilterSelect value={state.roleId} placeholder="全部角色" options={state.roleFilterOptions} disabled={!state.roleCatalogReady} onValueChange={state.changeRole} /></AdminListFilters></AdminToolbar>} isEmpty={state.empty} empty="暂无管理员" pagination={<AdminPagination total={state.total} page={state.page} pageSize={state.pageSize} onPageChange={state.setPage} onPageSizeChange={state.setPageSize} />}>
      <AdministratorTable rows={state.rows} currentAccountId={state.currentAccountId} permissions={state.permissions} onEdit={state.edit} onStatus={(row, active) => void state.setStatus(row, active)} onResetPassword={(row) => void state.reset(row)} onDelete={(row) => void state.remove(row)} />
    </AdminTableShell>
    <AdministratorFormDialog open={state.open} mode={state.editing ? 'edit' : 'create'} onOpenChange={state.changeOpen} values={state.values} onValuesChange={state.setValues} roles={state.roles} account={state.editing?.account} outcome={state.outcome} error={state.actionError} onSubmit={state.save} isSubmitting={state.saving} />
  </div>;
}
