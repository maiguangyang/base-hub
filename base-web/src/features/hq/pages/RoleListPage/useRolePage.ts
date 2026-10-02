import { useState } from 'react';
import { useMutation, useQuery } from '@apollo/client/react';
import type { HqRolesQuery, RoleKind, SystemPermissionsQuery } from '@/__generated__/graphql';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { useAuthStore } from '@/features/auth/store/authStore';
import type { ViewerSummary } from '@/features/auth/store/authStore';
import { getGraphQLErrorCode } from '@/lib/graphql/errors';
import {
  CREATE_HQ_ROLE_MUTATION,
  DELETE_HQ_ROLES_MUTATION,
  HQ_ROLES_QUERY,
  SYSTEM_PERMISSIONS_QUERY,
  UPDATE_HQ_ROLE_MUTATION,
} from '../../graphql/roles';
import { delegableSystemPermissions, type RoleValues } from './RoleFormDialog';
import type { RoleRow } from './RoleTable';
import { hqRoleFilter } from '../hqListFilters';

export function useRolePage() {
  const { search } = useAdminTab();
  const viewer = useAuthStore((state) => state.viewer);
  const permissions = hqViewerPermissions(viewer); const organizationId = hqOrganizationId(viewer);
  const [page, setPage] = useState(1); const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState(() => initialKeyword(search));
  const [kind, setKind] = useState<RoleKind>();
  const [open, setOpen] = useState(false); const [editing, setEditing] = useState<RoleRow>();
  const [values, setValues] = useState<RoleValues>({ name: '', permissionIds: [] }); const [actionError, setActionError] = useState<string>();
  const query = useQuery(HQ_ROLES_QUERY, { variables: { page, pageSize, q: optionalKeyword(keyword), filter: hqRoleFilter(organizationId, kind) }, skip: !organizationId });
  const catalog = useQuery(SYSTEM_PERMISSIONS_QUERY);
  const [createRole, creating] = useMutation(CREATE_HQ_ROLE_MUTATION);
  const [updateRole, updating] = useMutation(UPDATE_HQ_ROLE_MUTATION);
  const [deleteRoles] = useMutation(DELETE_HQ_ROLES_MUTATION);
  const rows = roleRows(query.data); const permissionCatalog = permissionRows(catalog.data);
  const catalogReady = permissionCatalogReady(catalog.data, catalog.loading, catalog.error);
  const catalogError = !catalog.loading && !catalogReady ? '无法加载完整的系统权限目录，已禁止角色写入。' : undefined;
  const edit = (row?: RoleRow) => { setActionError(undefined); setEditing(row); setValues(row ? { name: row.name, permissionIds: row.permissions.map((permission) => permission.id) } : { name: '', permissionIds: [] }); setOpen(true); };
  const changeOpen = (next: boolean) => { setOpen(next); if (!next) setEditing(undefined); };
  const searchRows = (next: string) => { setKeyword(next); setPage(1); };
  const changeKind = (next?: string) => { setKind(next as RoleKind | undefined); setPage(1); };
  const resetFilters = () => { setKeyword(''); setKind(undefined); setPage(1); };
  async function save() {
    setActionError(undefined);
    try {
      const input = buildRoleMutationInput(values, permissionCatalog, permissions, catalogReady);
      if (editing) await updateRole({ variables: { id: editing.id, input } });
      else await createRole({ variables: { input: { ...input, organizationId, kind: 'CUSTOM' } } });
      changeOpen(false); await query.refetch();
    } catch (cause) { setActionError(roleMutationErrorMessage(getGraphQLErrorCode(cause) ?? (cause instanceof Error ? cause.message : undefined))); }
  }
  async function remove(row: RoleRow) {
    if (!window.confirm(`确认删除角色 ${row.name}？`)) return;
    setActionError(undefined);
    try { await deleteRoles({ variables: { ids: [row.id] } }); await query.refetch(); }
    catch (cause) { setActionError(roleMutationErrorMessage(getGraphQLErrorCode(cause))); }
  }
  const view = roleViewState(query.data, query.loading, rows, creating.loading, updating.loading, query.error);
  return { rows, ...view, permissionCatalog, permissions, page, pageSize, keyword, kind, open, editing, values, actionError, catalogError, setPage, setPageSize, setValues, edit, changeOpen, save, remove, search: searchRows, changeKind, resetFilters, hasFilters: [keyword, kind].some(Boolean), canCreate: catalogReady && permissions.includes('hqRole:create'), canUpdate: catalogReady && permissions.includes('hqRole:update'), canDelete: permissions.includes('hqRole:delete') };
}

export function permissionCatalogReady(data: SystemPermissionsQuery | undefined, loading: boolean, error: unknown): boolean {
  const result = data?.permissions;
  return !loading && !error && result != null && result.data.length === result.total;
}

export function buildRoleMutationInput(values: RoleValues, catalog: ReturnType<typeof permissionRows>, viewerPermissions: readonly string[], catalogReady: boolean) {
  if (!catalogReady) throw new Error('PERMISSION_CATALOG_UNAVAILABLE');
  const name = values.name.trim();
  if (!name) throw new Error('VALIDATION_FAILED');
  const allowed = new Set(delegableSystemPermissions(catalog, viewerPermissions).map((permission) => permission.id));
  return { name, permissionsIds: values.permissionIds.filter((id) => allowed.has(id)) };
}

export function roleMutationErrorMessage(code: string | undefined): string {
  if (code === 'ROLE_IN_USE') return '该角色仍被有效管理员使用，不能删除。';
  if (code === 'PERMISSION_DELEGATION_DENIED') return '所选权限超出当前账号可委派范围。';
  if (code === 'VALIDATION_FAILED') return '请输入角色名称。';
  if (code === 'PERMISSION_CATALOG_UNAVAILABLE') return '系统权限目录未完整加载，请稍后重试。';
  if (code === 'PERMISSION_DENIED') return '当前账号没有执行此操作的权限。';
  return '操作失败，请稍后重试。';
}

function roleRows(data: HqRolesQuery | undefined) { return data?.operatorRoles?.data ?? []; }
function permissionRows(data: SystemPermissionsQuery | undefined) { return data?.permissions?.data ?? []; }
function hqViewerPermissions(viewer: ViewerSummary | null) { return viewer ? viewer.permissions : []; }
function hqOrganizationId(viewer: ViewerSummary | null) { return viewer?.currentWorkspace?.organizationId ?? ''; }
function initialKeyword(search: string) { return new URLSearchParams(search).get('q') ?? ''; }
function optionalKeyword(keyword: string) { return keyword === '' ? null : keyword; }
export function roleViewState(data: HqRolesQuery | undefined, loading: boolean, rows: RoleRow[], creating: boolean, updating: boolean, error?: unknown) {
  return {
    total: data?.operatorRoles?.total ?? 0,
    empty: !loading && !error && rows.length === 0,
    saving: creating || updating,
    listError: error ? '角色列表加载失败，请稍后重试。' : undefined,
  };
}
