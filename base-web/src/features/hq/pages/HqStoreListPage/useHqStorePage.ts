import { useState } from 'react';
import { useMutation, useQuery } from '@apollo/client/react';
import { storeBusinessStatusInput } from '@/features/admin/pages/StoreEditorPage/storeInputs';
import type { OrganizationType, StoreBusinessStatus, StoreLifecycle } from '@/__generated__/graphql';
import { useCompleteFilterCatalog } from '@/features/admin/lib/useCompleteFilterCatalog';
import { useAuthStore } from '@/features/auth/store/authStore';
import { HQ_FRANCHISES_QUERY } from '../../graphql/franchises';
import { HQ_STORES_QUERY } from '../../graphql/hqStores';
import { DELETE_DIRECT_STORES_MUTATION, UPDATE_DIRECT_STORE_MUTATION } from '../../graphql/directStores';
import { hqFranchiseFilter } from '../hqListFilters';
import type { HqStoreRow } from './HqStoreTable';
import { hqStoreFilter } from './storeFilters';

function useFranchiseChoices(enabled: boolean) {
  const query = useQuery(HQ_FRANCHISES_QUERY, {
    variables: { page: 1, pageSize: 200, q: null, filter: hqFranchiseFilter() }, skip: !enabled,
  });
  const first = query.data?.organizations;
  const catalog = useCompleteFilterCatalog('hq-store-franchises',
    first ? { items: first.data, total: first.total } : undefined,
    enabled && query.loading, enabled ? query.error : undefined,
    async (page) => {
      const result = await query.fetchMore({ variables: { page } });
      const next = result.data?.organizations;
      if (!next) throw new Error('catalog unavailable');
      return { items: next.data, total: next.total };
    });
  return {
    ...catalog,
    retry: () => { catalog.retry(); void query.refetch().catch(() => undefined); },
    options: catalog.items.map((item) => ({ value: item.id, label: item.name })),
  };
}

function useHqStoreActions(refetch: () => Promise<unknown>) {
  const [pendingDelete, setPendingDelete] = useState<HqStoreRow>();
  const [actionError, setActionError] = useState<string>();
  const [updateStore, updating] = useMutation(UPDATE_DIRECT_STORE_MUTATION);
  const [deleteStores, deleting] = useMutation(DELETE_DIRECT_STORES_MUTATION);
  const refreshAfterWrite = async (completed: string) => {
    try { await refetch(); }
    catch { setActionError(`${completed}，但门店列表刷新失败，请刷新页面。`); }
  };
  const remove = async () => {
    if (!pendingDelete) return;
    try {
      await deleteStores({ variables: { ids: [pendingDelete.id] } });
    } catch { setActionError('删除门店失败，请重试。'); return; }
    setPendingDelete(undefined); setActionError(undefined);
    await refreshAfterWrite('门店已删除');
  };
  const closeDelete = () => {
    setPendingDelete(undefined);
    if (actionError?.startsWith('删除门店失败')) setActionError(undefined);
  };
  const toggleBusinessStatus = async (row: HqStoreRow, open: boolean) => {
    try {
      await updateStore({ variables: { id: row.id, input: storeBusinessStatusInput(open) } });
    } catch { setActionError('更新营业状态失败，请重试。'); return; }
    setActionError(undefined);
    await refreshAfterWrite('营业状态已更新');
  };
  return { pendingDelete, setPendingDelete, actionError, remove, closeDelete, toggleBusinessStatus,
    saving: updating.loading || deleting.loading };
}

export function useHqStorePage(search: string) {
  const { permissions, hqOrganizationId } = useHqStoreContext();
  const initial = initialStoreFilters(search);
  const [type, setType] = useState<OrganizationType | undefined>(initial.type);
  const [organizationId, setOrganizationId] = useState<string | undefined>(initial.organizationId);
  const [keyword, setKeyword] = useState(initial.keyword);
  const [lifecycle, setLifecycle] = useState<StoreLifecycle>();
  const [businessStatus, setBusinessStatus] = useState<StoreBusinessStatus>();
  const [page, setPage] = useState(1); const [pageSize, setPageSize] = useState(10);
  const variables = { page, pageSize, q: keyword || null, filter: hqStoreFilter(type, organizationId, lifecycle, businessStatus) };
  const query = useQuery(HQ_STORES_QUERY, { variables });
  const catalog = useFranchiseChoices(permissions.includes('organization:read'));
  const actions = useHqStoreActions(() => query.refetch());
  const resetPage = () => setPage(1);
  const changeType = (value?: string) => { setType(value as OrganizationType | undefined); setOrganizationId(undefined); resetPage(); };
  const changeOrganization = (value?: string) => { setOrganizationId(value); if (value) setType('FRANCHISE'); resetPage(); };
  const resetFilters = () => { setKeyword(''); setType(undefined); setOrganizationId(undefined); setLifecycle(undefined); setBusinessStatus(undefined); resetPage(); };
  return { ...actions, catalog, hqOrganizationId, permissions, type, organizationId, keyword, lifecycle, businessStatus,
    page, pageSize, setPage, setPageSize, changeType, changeOrganization, resetFilters,
    changeKeyword: (value: string) => { setKeyword(value); resetPage(); },
    changeLifecycle: (value?: string) => { setLifecycle(value as StoreLifecycle | undefined); resetPage(); },
    changeBusinessStatus: (value?: string) => { setBusinessStatus(value as StoreBusinessStatus | undefined); resetPage(); },
    hasFilters: [type, organizationId, keyword, lifecycle, businessStatus].some(Boolean),
    rows: query.data?.stores?.data ?? [], total: query.data?.stores?.total ?? 0,
    loading: query.loading, loadFailed: Boolean(query.error), retry: () => { void query.refetch().catch(() => undefined); },
  };
}

function useHqStoreContext() {
  const viewer = useAuthStore((state) => state.viewer);
  return { permissions: viewer?.permissions ?? [], hqOrganizationId: viewer?.currentWorkspace?.organizationId ?? '' };
}

function initialStoreFilters(search: string) {
  const params = new URLSearchParams(search);
  const rawType = params.get('type');
  let type: OrganizationType | undefined;
  if (rawType === 'HEADQUARTERS' || rawType === 'FRANCHISE') type = rawType;
  return { type, organizationId: params.get('organizationId') || undefined, keyword: params.get('q') ?? '' };
}
