import { useState } from 'react';
import { useMutation, useQuery } from '@apollo/client/react';
import type { HqAdministratorIdentityQuery, HqAdministratorsQuery, MembershipStatus } from '@/__generated__/graphql';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { roleDisplayName } from '@/features/admin/config/roleDisplayName';
import { useAuthStore } from '@/features/auth/store/authStore';
import type { ViewerSummary } from '@/features/auth/store/authStore';
import { getGraphQLErrorCode } from '@/lib/graphql/errors';
import { catalogPage } from '@/features/admin/lib/filterCatalog';
import { useCompleteFilterCatalog } from '@/features/admin/lib/useCompleteFilterCatalog';
import { HQ_ROLES_QUERY } from '../../graphql/roles';
import {
  CHANGE_HQ_ADMINISTRATOR_STATUS_MUTATION,
  DELETE_HQ_ADMINISTRATORS_MUTATION,
  HQ_ADMINISTRATOR_IDENTITY_QUERY,
  HQ_ADMINISTRATORS_QUERY,
  INVITE_HQ_ADMINISTRATOR_MUTATION,
  RESET_HQ_ADMINISTRATOR_PASSWORD_MUTATION,
  UPDATE_HQ_ADMINISTRATOR_MUTATION,
} from '../../graphql/administrators';
import type { AdministratorRow } from './AdministratorTable';
import { hqAdministratorFilter, hqRoleFilter } from '../hqListFilters';
import {
  administratorInviteInput,
  administratorInviteOutcome,
  assignableAdministratorRoles,
  emptyAdministratorValues,
  type AdministratorInviteOutcome,
  type AdministratorValues,
} from './AdministratorFormDialog';

export function useAdministratorPage() {
  const { search } = useAdminTab();
  const viewer = useAuthStore((state) => state.viewer);
  const permissions = hqViewerPermissions(viewer); const organizationId = viewer?.currentWorkspace?.organizationId ?? '';
  const [page, setPage] = useState(1); const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState(() => initialKeyword(search)); const [status, setStatusFilter] = useState<MembershipStatus>(); const [roleId, setRoleId] = useState<string>();
  const [open, setOpen] = useState(false); const [editing, setEditing] = useState<AdministratorRow>();
  const [values, setValues] = useState<AdministratorValues>(emptyAdministratorValues());
  const [outcome, setOutcome] = useState<AdministratorInviteOutcome>(); const [actionError, setActionError] = useState<string>();
  const query = useQuery(HQ_ADMINISTRATORS_QUERY, { variables: { page, pageSize, q: optionalKeyword(keyword), filter: hqAdministratorFilter(status, roleId) } });
  const identity = useQuery(HQ_ADMINISTRATOR_IDENTITY_QUERY, { variables: { accountId: viewerAccountId(viewer) }, skip: viewer === null });
  const [invite, inviting] = useMutation(INVITE_HQ_ADMINISTRATOR_MUTATION);
  const [update, updating] = useMutation(UPDATE_HQ_ADMINISTRATOR_MUTATION);
  const [changeStatus] = useMutation(CHANGE_HQ_ADMINISTRATOR_STATUS_MUTATION);
  const [removeMembership] = useMutation(DELETE_HQ_ADMINISTRATORS_MUTATION);
  const [resetPassword] = useMutation(RESET_HQ_ADMINISTRATOR_PASSWORD_MUTATION);
  const roleDirectory = useHqRoleDirectory(organizationId);
  const rows = administratorRows(query.data); const roleCatalog = roleDirectory.items;
  const roles = assignableAdministratorRoles(roleCatalog, identityIsSuper(identity.data));
  const roleCatalogReady = roleDirectory.ready;
  const edit = (row?: AdministratorRow) => { setActionError(undefined); setEditing(row); setValues(row ? { phone: '', displayName: '', email: '', roleIds: row.roles.map((role) => role.id) } : emptyAdministratorValues()); setOutcome(undefined); setOpen(true); };
  const changeOpen = (next: boolean) => { setOpen(next); if (!next) { setEditing(undefined); setOutcome(undefined); setValues(emptyAdministratorValues()); } };
  const searchRows = (next: string) => { setKeyword(next); setPage(1); }; const changeStatusFilter = (next?: string) => { setStatusFilter(next as MembershipStatus | undefined); setPage(1); };
  const changeRole = (next?: string) => { setRoleId(next); setPage(1); };
  const resetFilters = () => { setKeyword(''); setStatusFilter(undefined); setRoleId(undefined); setPage(1); };
  async function save() {
    setActionError(undefined);
    try {
      if (editing) await update({ variables: { id: editing.id, input: { rolesIds: values.roleIds, storeAccessMode: 'ALL_STORES', storesIds: [] } } });
      else { const result = await invite({ variables: { input: administratorInviteInput(values) } }); if (result.data) setOutcome(administratorInviteOutcome(result.data.inviteOperator)); }
      await query.refetch();
      if (editing) changeOpen(false);
    } catch (cause) { setActionError(administratorMutationErrorMessage(getGraphQLErrorCode(cause))); }
  }
  async function setStatus(row: AdministratorRow, active: boolean) {
    if (!window.confirm(`确认${active ? '启用' : '停用'}管理员 ${row.account.displayName}？`)) return;
    await perform(async () => { await changeStatus({ variables: { input: { membershipId: row.id, status: active ? 'ACTIVE' : 'SUSPENDED' } } }); await query.refetch(); });
  }
  async function reset(row: AdministratorRow) {
    if (!window.confirm(`确认重置 ${row.account.displayName} 的密码？该账号全部在线会话将被撤销。`)) return;
    await perform(async () => { const result = await resetPassword({ variables: { accountId: row.account.id } }); const password = result.data?.resetTemporaryPassword.temporaryPassword; if (password) window.alert(`临时密码仅显示一次：${password}`); });
  }
  async function remove(row: AdministratorRow) {
    if (!window.confirm(`确认删除管理员 ${row.account.displayName}？其总部在线会话将立即失效。`)) return;
    await perform(async () => { await removeMembership({ variables: { ids: [row.id] } }); await query.refetch(); });
  }
  async function perform(action: () => Promise<void>) { setActionError(undefined); try { await action(); } catch (cause) { setActionError(administratorMutationErrorMessage(getGraphQLErrorCode(cause))); } }
  const view = administratorViewState(viewer, query.data, query.loading, rows, inviting.loading, updating.loading, query.error);
  return { rows, roles, permissions, ...view, page, pageSize, keyword, status, roleId, open, editing, values, outcome, actionError, setPage, setPageSize, setValues, edit, changeOpen, save, setStatus, reset, remove, search: searchRows, changeStatusFilter, changeRole, resetFilters, hasFilters: [keyword, status, roleId].some(Boolean), roleFilterOptions: roleCatalog.map((role) => ({ value: String(role.id), label: roleDisplayName(role) })), roleCatalogReady, filterCatalogError: roleDirectory.loading || roleCatalogReady ? undefined : '角色目录未完整加载，暂时无法按角色筛选。', canCreate: permissions.includes('hqMembership:create') };
}

function useHqRoleDirectory(organizationId: string) {
  const query = useQuery(HQ_ROLES_QUERY, { variables: { page: 1, pageSize: 200, q: null, filter: hqRoleFilter(organizationId) }, skip: !organizationId });
  return useCompleteFilterCatalog(`hq-roles:${organizationId}`, catalogPage(query.data?.operatorRoles), query.loading, query.error, async (nextPage) => {
    const next = await query.fetchMore({ variables: { page: nextPage } });
    const result = catalogPage(next.data?.operatorRoles);
    if (!result) throw new Error('catalog unavailable');
    return result;
  });
}

export function administratorMutationErrorMessage(code: string | undefined): string {
  if (code === 'SELF_MEMBERSHIP_CHANGE_DENIED') return '当前账号不能管理自己的总部成员关系。';
  if (code === 'PERMISSION_DELEGATION_DENIED') return '所选角色或目标管理员超出当前账号的授权范围。';
  if (code === 'LAST_HQ_SUPER_ADMIN_REQUIRED') return '总部必须保留至少一位有效的超级管理员。';
  if (code === 'ROLE_IN_USE') return '该角色仍被有效管理员使用。';
  if (code === 'CONFLICT') return '该账号已存在总部成员关系或待处理邀请。';
  if (code === 'VALIDATION_FAILED') return '请填写完整信息并至少选择一个角色。';
  if (code === 'PERMISSION_DENIED') return '当前账号没有执行此操作的权限。';
  return '操作失败，请稍后重试。';
}

function administratorRows(data: HqAdministratorsQuery | undefined) { return data?.operatorMemberships?.data ?? []; }
function hqViewerPermissions(viewer: ViewerSummary | null) { return viewer ? viewer.permissions : []; }
function viewerAccountId(viewer: ViewerSummary | null) { return viewer ? viewer.account.id : ''; }
function initialKeyword(search: string) { return new URLSearchParams(search).get('q') ?? ''; }
function optionalKeyword(keyword: string) { return keyword === '' ? null : keyword; }
export function administratorViewState(viewer: ViewerSummary | null, data: HqAdministratorsQuery | undefined, loading: boolean, rows: AdministratorRow[], inviting: boolean, updating: boolean, error?: unknown) {
  return {
    currentAccountId: viewer ? viewer.account.id : '',
    total: data?.operatorMemberships?.total ?? 0,
    empty: !loading && !error && rows.length === 0,
    saving: inviting || updating,
    listError: error ? '管理员列表加载失败，请稍后重试。' : undefined,
  };
}
function identityIsSuper(data: HqAdministratorIdentityQuery | undefined) {
  return data?.operatorMemberships?.data.some((membership) => membership.roles.some((role) => role.kind === 'HQ_SUPER_ADMIN')) ?? false;
}
