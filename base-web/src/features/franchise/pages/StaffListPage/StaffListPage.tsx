import { useState } from 'react';
import { useMutation, useQuery } from '@apollo/client/react';
import type { FranchiseStaffQuery, MembershipStatus, StoreAccessMode } from '@/__generated__/graphql';
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
import { membershipStatusOptions, storeAccessModeOptions } from '@/features/admin/config/filterOptions';
import { roleDisplayName } from '@/features/admin/config/roleDisplayName';
import { catalogPage } from '@/features/admin/lib/filterCatalog';
import { useCompleteFilterCatalog } from '@/features/admin/lib/useCompleteFilterCatalog';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { useAuthStore } from '@/features/auth/store/authStore';
import { getGraphQLErrorCode } from '@/lib/graphql/errors';
import { FRANCHISE_ROLES_QUERY } from '../../graphql/roles';
import { FRANCHISE_STORES_QUERY } from '../../graphql/stores';
import { CHANGE_MEMBERSHIP_STATUS_MUTATION, FRANCHISE_STAFF_QUERY, INVITE_FRANCHISE_STAFF_MUTATION, UPDATE_FRANCHISE_STAFF_MUTATION } from '../../graphql/staff';
import { emptyInviteValues, InviteStaffDialog, staffAccessInput, staffInviteOutcome, type InviteStaffValues, type StaffInviteOutcome } from './InviteStaffDialog';
import { emptyStaffAccessValues, StaffAccessDialog, staffAccessValues, type StaffAccessValues } from './StaffAccessDialog';
import { StaffTable, type StaffRow } from './StaffTable';
import { staffStats } from './staffStats';
import { franchiseStaffFilter } from '../franchiseListFilters';
import { initialSearchKeyword, queryKeyword, viewerOrganizationId, viewerPermissions } from '../pageState';

export function StaffListPage() {
  const state = useStaffPage();
  return (
    <div><AdminPageHeader title="员工与邀请" description="分配多个角色以及全部或指定门店范围。" /><AdminStatsStrip items={staffStats(state.rows, state.total)} /><AdminBulkBar count={0} onClear={() => undefined} />
      <AdminActionError message={state.listError ?? state.error ?? state.filterCatalogError} />
      <AdminTableShell toolbar={<AdminToolbar rightSlot={state.canCreate ? <AdminPrimaryActionButton onClick={() => state.setOpen(true)}>新增员工</AdminPrimaryActionButton> : null}><AdminListFilters hasActiveFilters={state.hasFilters} onReset={state.resetFilters}><AdminListSearch value={state.keyword} placeholder="搜索员工" onSearch={state.search} /><AdminFilterSelect value={state.statusFilter} placeholder="全部成员状态" options={membershipStatusOptions} onValueChange={state.changeStatusFilter} /><AdminFilterSelect value={state.roleFilterId} placeholder="全部角色" options={state.roleFilterOptions} disabled={!state.roleCatalogReady} onValueChange={state.changeRoleFilter} /><AdminFilterSelect value={state.storeAccessModeFilter} placeholder="全部门店范围" options={storeAccessModeOptions} onValueChange={state.changeStoreAccessModeFilter} /><AdminFilterSelect value={state.storeFilterId} placeholder="全部指定门店" options={state.storeFilterOptions} disabled={!state.storeCatalogReady} onValueChange={state.changeStoreFilter} /></AdminListFilters></AdminToolbar>} isEmpty={state.empty} empty="暂无员工" pagination={<AdminPagination total={state.total} page={state.page} pageSize={state.pageSize} onPageChange={state.setPage} onPageSizeChange={state.setPageSize} />}>
        <StaffTable rows={state.rows} canUpdate={state.canUpdate} onEdit={state.edit} onStatus={(row, active) => void state.changeStatus(row, active)} />
      </AdminTableShell>
      <InviteStaffDialog open={state.open} onOpenChange={state.changeOpen} values={state.values} onValuesChange={state.setValues} roles={state.roles} stores={state.stores} onSubmit={state.invite} isSubmitting={state.inviting} outcome={state.outcome} />
      <StaffAccessDialog open={state.editOpen} onOpenChange={state.changeEditOpen} values={state.accessValues} onValuesChange={state.setAccessValues} roles={state.roles} stores={state.stores} onSubmit={state.updateAccess} isSubmitting={state.updating} />
    </div>
  );
}

function useStaffPage() {
  const { search } = useAdminTab();
  const viewer = useAuthStore((state) => state.viewer);
  const organizationId = viewerOrganizationId(viewer);
  const permissions = viewerPermissions(viewer);
  const [page, setPage] = useState(1); const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState(() => initialSearchKeyword(search));
  const [statusFilter, setStatusFilter] = useState<MembershipStatus>(); const [roleFilterId, setRoleFilterId] = useState<string>();
  const [storeAccessModeFilter, setStoreAccessModeFilter] = useState<StoreAccessMode>(); const [storeFilterId, setStoreFilterId] = useState<string>();
  const [open, setOpen] = useState(false); const [values, setValues] = useState<InviteStaffValues>(emptyInviteValues());
  const [editing, setEditing] = useState<StaffRow>();
  const [accessValues, setAccessValues] = useState<StaffAccessValues>(emptyStaffAccessValues());
  const [outcome, setOutcome] = useState<StaffInviteOutcome>();
  const [error, setError] = useState<string>();
  const query = useQuery(FRANCHISE_STAFF_QUERY, { variables: { page, pageSize, q: queryKeyword(keyword), filter: franchiseStaffFilter(statusFilter, roleFilterId, storeAccessModeFilter, storeFilterId) } });
  const roleQuery = useQuery(FRANCHISE_ROLES_QUERY, { variables: { page: 1, pageSize: 200, q: null, filter: undefined } });
  const storeQuery = useQuery(FRANCHISE_STORES_QUERY, { variables: { page: 1, pageSize: 200, q: null, filter: undefined } });
  const roleDirectory = useCompleteFilterCatalog(`franchise-roles:${organizationId}`, catalogPage(roleQuery.data?.operatorRoles), roleQuery.loading, roleQuery.error, async (nextPage) => {
    const next = await roleQuery.fetchMore({ variables: { page: nextPage } });
    const result = catalogPage(next.data?.operatorRoles);
    if (!result) throw new Error('catalog unavailable');
    return result;
  });
  const storeDirectory = useCompleteFilterCatalog(`franchise-stores:${organizationId}`, catalogPage(storeQuery.data?.stores), storeQuery.loading, storeQuery.error, async (nextPage) => {
    const next = await storeQuery.fetchMore({ variables: { page: nextPage } });
    const result = catalogPage(next.data?.stores);
    if (!result) throw new Error('catalog unavailable');
    return result;
  });
  const [inviteMutation, inviteState] = useMutation(INVITE_FRANCHISE_STAFF_MUTATION);
  const [updateMutation, updateState] = useMutation(UPDATE_FRANCHISE_STAFF_MUTATION);
  const [statusMutation] = useMutation(CHANGE_MEMBERSHIP_STATUS_MUTATION);
  const { rows, total } = staffRows(query.data);
  const roleCatalogReady = roleDirectory.ready;
  const storeCatalogReady = storeDirectory.ready;
  const catalogsReady = [roleCatalogReady, storeCatalogReady].every(Boolean);
  const changeOpen = (next: boolean) => { setOpen(next); if (!next) { setOutcome(undefined); setValues(emptyInviteValues()); } };
  const changeEditOpen = (next: boolean) => { if (!next) setEditing(undefined); };
  const edit = (row: StaffRow) => { setEditing(row); setAccessValues(staffAccessValues(row)); };
  const searchRows = (next: string) => { setKeyword(next); setPage(1); };
  const changeStatusFilter = (next?: string) => { setStatusFilter(next as MembershipStatus | undefined); setPage(1); };
  const changeRoleFilter = (next?: string) => { setRoleFilterId(next); setPage(1); };
  const changeStoreAccessModeFilter = (next?: string) => { setStoreAccessModeFilter(next as StoreAccessMode | undefined); setPage(1); };
  const changeStoreFilter = (next?: string) => { setStoreFilterId(next); setPage(1); };
  const resetFilters = () => { setKeyword(''); setStatusFilter(undefined); setRoleFilterId(undefined); setStoreAccessModeFilter(undefined); setStoreFilterId(undefined); setPage(1); };
  async function invite() { setError(undefined); try { const access = staffAccessInput(values.storeAccessMode, values.storeIds); const result = await inviteMutation({ variables: { input: { phone: values.phone, displayName: values.displayName, email: values.email || null, roleIds: values.roleIds, ...access } } }); if (result.data) setOutcome(staffInviteOutcome(result.data.inviteOperator)); await query.refetch(); } catch (cause) { setError(staffMutationErrorMessage(getGraphQLErrorCode(cause))); } }
  async function updateAccess() { if (!editing) return; setError(undefined); try { const access = staffAccessInput(accessValues.storeAccessMode, accessValues.storeIds); await updateMutation({ variables: { id: editing.id, input: { rolesIds: accessValues.roleIds, storesIds: access.storeIds, storeAccessMode: access.storeAccessMode } } }); setEditing(undefined); await query.refetch(); } catch (cause) { setError(staffMutationErrorMessage(getGraphQLErrorCode(cause))); } }
  async function changeStatus(row: StaffRow, active: boolean) { setError(undefined); try { await statusMutation({ variables: { input: { membershipId: row.id, status: active ? 'ACTIVE' : 'SUSPENDED' } } }); await query.refetch(); } catch (cause) { setError(staffMutationErrorMessage(getGraphQLErrorCode(cause))); } }
  return { rows, total, page, pageSize, keyword, statusFilter, roleFilterId, storeAccessModeFilter, storeFilterId, open, editOpen: Boolean(editing), values, accessValues, outcome, error, listError: query.error ? '员工列表加载失败，请稍后重试。' : undefined, setPage, setPageSize, setOpen, setValues, setAccessValues, changeOpen, changeEditOpen, edit, invite, updateAccess, changeStatus, search: searchRows, changeStatusFilter, changeRoleFilter, changeStoreAccessModeFilter, changeStoreFilter, resetFilters, hasFilters: [keyword, statusFilter, roleFilterId, storeAccessModeFilter, storeFilterId].some(Boolean), roles: roleDirectory.items, stores: storeDirectory.items, roleFilterOptions: roleDirectory.items.map((role) => ({ value: String(role.id), label: roleDisplayName(role) })), storeFilterOptions: storeDirectory.items.map((store) => ({ value: String(store.id), label: store.name })), roleCatalogReady, storeCatalogReady, filterCatalogError: [roleDirectory.loading, storeDirectory.loading, catalogsReady].some(Boolean) ? undefined : '角色或门店目录未完整加载，部分筛选暂不可用。', canCreate: permissions.includes('operatorMembership:create'), canUpdate: permissions.includes('operatorMembership:update'), empty: !query.loading && !query.error && rows.length === 0, inviting: inviteState.loading, updating: updateState.loading };
}

export function staffMutationErrorMessage(code: string | undefined): string {
  if (code === 'LAST_OWNER_REQUIRED') return '每个加盟商必须保留至少一位有效所有者。';
  if (code === 'PERMISSION_DENIED') return '当前账号没有执行此操作的权限。';
  if (code === 'VALIDATION_FAILED') return '提交内容不符合要求，请检查后重试。';
  return '操作失败，请稍后重试。';
}

function staffRows(data: FranchiseStaffQuery | undefined) {
  return { rows: data?.operatorMemberships?.data ?? [], total: data?.operatorMemberships?.total ?? 0 };
}
