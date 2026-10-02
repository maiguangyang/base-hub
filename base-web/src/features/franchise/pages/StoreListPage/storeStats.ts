import type { FranchiseStoresQuery } from '@/__generated__/graphql';
import type { AdminStatsItem } from '@/features/admin/components/AdminStatsStrip/types';

export function storeStats(rows: NonNullable<FranchiseStoresQuery['stores']>['data'], total: number): AdminStatsItem[] {
  return [
    { key: 'total', label: '门店', value: String(total), icon: 'list-checks' },
    { key: 'active', label: '已启用', value: String(rows.filter((row) => row.lifecycle === 'ACTIVE').length), tone: 'success', icon: 'check-circle' },
    { key: 'pending', label: '待审核', value: String(rows.filter((row) => row.lifecycle === 'PENDING_APPROVAL').length), tone: 'warning', icon: 'timer' },
  ];
}
