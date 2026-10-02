import { useState } from 'react';
import { useMutation, useQuery } from '@apollo/client/react';
import type { FranchiseRolesQuery, RoleKind, TenantPermissionsQuery } from '@/__generated__/graphql';
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
import { franchiseRoleKindOptions } from '@/features/admin/config/filterOptions';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { useAuthStore } from '@/features/auth/store/authStore';
import { CREATE_FRANCHISE_ROLE_MUTATION, DELETE_FRANCHISE_ROLES_MUTATION, FRANCHISE_ROLES_QUERY, TENANT_PERMISSIONS_QUERY, UPDATE_FRANCHISE_ROLE_MUTATION } from '../../graphql/roles';
import { RoleFormDialog, type RoleValues } from './RoleFormDialog';
import { RoleTable, type RoleRow } from './RoleTable';
import { roleStats } from './roleStats';
import { franchiseRoleFilter } from '../franchiseListFilters';
import { initialSearchKeyword, queryKeyword, viewerOrganizationId, viewerPermissions } from '../pageState';

export function RoleListPage() {
  const state = useRolePage();
  return (
    <div><AdminPageHeader title="角色与权限" description="自定义角色仅能分配租户权限，多角色按允许并集生效。" /><AdminStatsStrip items={roleStats(state.rows, state.total)} /><AdminBulkBar count={0} onClear={() => undefined} />
		<AdminActionError message={state.listError ?? state.actionError} />
      <AdminTableShell toolbar={<AdminToolbar rightSlot={state.canCreate ? <AdminPrimaryActionButton onClick={() => state.edit()}>新增</AdminPrimaryActionButton> : null}><AdminListFilters hasActiveFilters={state.hasFilters} onReset={state.resetFilters}><AdminListSearch value={state.keyword} placeholder="搜索角色名称" onSearch={state.search} /><AdminFilterSelect value={state.kind} placeholder="全部角色类型" options={franchiseRoleKindOptions} onValueChange={state.changeKind} /></AdminListFilters></AdminToolbar>} isEmpty={state.empty} empty="暂无角色" pagination={<AdminPagination total={state.total} page={state.page} pageSize={state.pageSize} onPageChange={state.setPage} onPageSizeChange={state.setPageSize} />}>
        <RoleTable rows={state.rows} canUpdate={state.canUpdate} canDelete={state.canDelete} onEdit={state.edit} onDelete={(row) => void state.remove(row)} />
      </AdminTableShell>
      <RoleFormDialog open={state.open} mode={state.editing ? 'edit' : 'create'} onOpenChange={state.setOpen} organizationId={state.organizationId} permissions={state.permissionsCatalog} values={state.values} onValuesChange={state.setValues} onSubmit={state.save} isSubmitting={state.saving} />
    </div>
  );
}

function useRolePage() {
  const { search } = useAdminTab();
  const viewer = useAuthStore((state) => state.viewer);
  const auth = viewerPermissions(viewer);
  const organizationId = viewerOrganizationId(viewer);
  const [page, setPage] = useState(1); const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState(() => initialSearchKeyword(search));
  const [kind, setKind] = useState<RoleKind>();
  const [open, setOpen] = useState(false); const [editing, setEditing] = useState<RoleRow>();
  const [values, setValues] = useState<RoleValues>({ name: '', permissionIds: [] });
  const [actionError, setActionError] = useState<string>();
  const query = useQuery(FRANCHISE_ROLES_QUERY, { variables: { page, pageSize, q: queryKeyword(keyword), filter: franchiseRoleFilter(kind) } });
  const catalog = useQuery(TENANT_PERMISSIONS_QUERY);
  const [createRole, creating] = useMutation(CREATE_FRANCHISE_ROLE_MUTATION);
  const [updateRole, updating] = useMutation(UPDATE_FRANCHISE_ROLE_MUTATION);
  const [deleteRoles] = useMutation(DELETE_FRANCHISE_ROLES_MUTATION);
  const result = roleResult(query.data, catalog.data); const { rows, total } = result;
  const edit = (row?: RoleRow) => { setEditing(row); setValues(row ? { name: row.name, permissionIds: row.permissions.map((permission) => permission.id) } : { name: '', permissionIds: [] }); setOpen(true); };
  const searchRows = (next: string) => { setKeyword(next); setPage(1); };
  const changeKind = (next?: string) => { setKind(next as RoleKind | undefined); setPage(1); };
  const resetFilters = () => { setKeyword(''); setKind(undefined); setPage(1); };
  async function save() { const input = { name: values.name, permissionsIds: values.permissionIds }; if (editing) await updateRole({ variables: { id: editing.id, input } }); else await createRole({ variables: { input: { ...input, organizationId, kind: 'CUSTOM' } } }); setOpen(false); await query.refetch(); }
  async function remove(row: RoleRow) { if (!window.confirm(`确认删除角色 ${row.name}？`)) return; await runAdminAction(async () => { await deleteRoles({ variables: { ids: [row.id] } }); await query.refetch(); }, setActionError); }
  return { rows, total, page, pageSize, keyword, kind, open, editing, values, actionError, listError: query.error ? '角色加载失败，请稍后重试。' : undefined, organizationId, setPage, setPageSize, setOpen, setValues, edit, save, remove, search: searchRows, changeKind, resetFilters, hasFilters: [keyword, kind].some(Boolean), permissionsCatalog: result.permissions, canCreate: auth.includes('operatorRole:create'), canUpdate: auth.includes('operatorRole:update'), canDelete: auth.includes('operatorRole:delete'), empty: !query.loading && !query.error && rows.length === 0, saving: creating.loading || updating.loading };
}

function roleResult(roles: FranchiseRolesQuery | undefined, permissions: TenantPermissionsQuery | undefined) {
  return {
    rows: roles?.operatorRoles?.data ?? [], total: roles?.operatorRoles?.total ?? 0,
    permissions: permissions?.permissions?.data ?? [],
  };
}
