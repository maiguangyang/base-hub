import { useQuery } from '@apollo/client/react';
import { AdminPageHeader } from '@/features/admin/components/AdminPageHeader';
import { AdminStatsStrip } from '@/features/admin/components/AdminStatsStrip';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { useAuthStore } from '@/features/auth/store/authStore';
import { HQ_STORES_TOTAL_QUERY } from '../../graphql/hqStores';
import { HQ_FRANCHISES_QUERY } from '../../graphql/franchises';
import { HQ_STORE_APPROVALS_QUERY } from '../../graphql/storeApprovals';
import { hqFranchiseFilter, hqStoreApprovalFilter } from '../hqListFilters';

function storeTotalValue(canReadAllStores: boolean, total: number | null | undefined, loading: boolean, error: unknown): string {
  if (!canReadAllStores) return '无权限';
  if (error) return '加载失败';
  if (loading) return '加载中';
  return total == null ? '加载失败' : String(total);
}

export function HqDashboardPage() {
  useAdminTab();
  const canReadAllStores = useAuthStore((state) => state.viewer?.permissions.includes('store:read_all') ?? false);
  const common = { page: 1, pageSize: 1, q: null };
  const franchises = useQuery(HQ_FRANCHISES_QUERY, { variables: { ...common, filter: hqFranchiseFilter() } });
  const stores = useQuery(HQ_STORES_TOTAL_QUERY, { skip: !canReadAllStores });
  const approvals = useQuery(HQ_STORE_APPROVALS_QUERY, { variables: { ...common, filter: hqStoreApprovalFilter() } });
  const storeTotal = storeTotalValue(canReadAllStores, stores.data?.stores?.total, stores.loading, stores.error);
  const stats = [
    { key: 'franchises', label: '加盟商', value: String(franchises.data?.organizations?.total ?? 0), icon: 'users' },
    { key: 'stores', label: '门店总数', value: storeTotal, icon: 'list-checks' },
    { key: 'approvals', label: '待审门店', value: String(approvals.data?.stores?.total ?? 0), tone: 'warning' as const, icon: 'timer' },
  ];
  return <div><AdminPageHeader title="总部概览" description="所有指标均来自当前权威数据。" /><AdminStatsStrip items={stats} columns={3} /></div>;
}
