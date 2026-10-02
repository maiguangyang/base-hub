import { AdminBulkBar } from '@/features/admin/components/AdminBulkBar';
import { AdminFilterSelect } from '@/features/admin/components/AdminFilterSelect';
import { AdminActionError } from '@/features/admin/components/AdminFormDialogShell';
import { AdminListFilters } from '@/features/admin/components/AdminListFilters';
import { AdminPageHeader } from '@/features/admin/components/AdminPageHeader';
import { AdminListSearch } from '@/features/admin/components/AdminListSearch';
import { AdminPagination } from '@/features/admin/components/AdminPagination';
import { AdminPrimaryActionButton } from '@/features/admin/components/AdminPrimaryActionButton';
import { AdminStatsStrip } from '@/features/admin/components/AdminStatsStrip';
import { AdminTableShell } from '@/features/admin/components/AdminTableShell';
import { AdminToolbar } from '@/features/admin/components/AdminToolbar';
import { hqRoleKindOptions } from '@/features/admin/config/filterOptions';
import { RoleFormDialog } from './RoleFormDialog';
import { RoleTable } from './RoleTable';
import { roleStats } from './roleStats';
import { useRolePage } from './useRolePage';

export function RoleListPage() {
  const state = useRolePage();
  return <div><AdminPageHeader title="角色与权限" description="自定义总部角色使用精确的读取、新增、修改和删除权限。" /><AdminStatsStrip items={roleStats(state.rows, state.total)} /><AdminBulkBar count={0} onClear={() => undefined} />
    <AdminActionError message={state.listError ?? state.actionError ?? state.catalogError} />
    <AdminTableShell toolbar={<AdminToolbar rightSlot={state.canCreate ? <AdminPrimaryActionButton onClick={() => state.edit()}>新增</AdminPrimaryActionButton> : null}><AdminListFilters hasActiveFilters={state.hasFilters} onReset={state.resetFilters}><AdminListSearch value={state.keyword} placeholder="搜索角色名称" onSearch={state.search} /><AdminFilterSelect value={state.kind} placeholder="全部角色类型" options={hqRoleKindOptions} onValueChange={state.changeKind} /></AdminListFilters></AdminToolbar>} isEmpty={state.empty} empty="暂无角色" pagination={<AdminPagination total={state.total} page={state.page} pageSize={state.pageSize} onPageChange={state.setPage} onPageSizeChange={state.setPageSize} />}>
      <RoleTable rows={state.rows} permissions={state.permissions} canUpdate={state.canUpdate} canDelete={state.canDelete} onEdit={state.edit} onDelete={(row) => void state.remove(row)} />
    </AdminTableShell>
    <RoleFormDialog open={state.open} mode={state.editing ? 'edit' : 'create'} onOpenChange={state.changeOpen} permissions={state.permissionCatalog} viewerPermissions={state.permissions} values={state.values} onValuesChange={state.setValues} error={state.actionError ?? state.catalogError} onSubmit={state.save} isSubmitting={state.saving} />
  </div>;
}
