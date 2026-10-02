import type { AdminStatsItem } from '@/features/admin/components/AdminStatsStrip/types';
import type { HqDirectStoresQuery } from '@/__generated__/graphql';

export function directStoreStats(rows: NonNullable<HqDirectStoresQuery['stores']>['data'], total: number): AdminStatsItem[] {
  return [
    { key: 'total', label: '直营门店', value: String(total), icon: 'list-checks' },
    { key: 'active', label: '已启用', value: String(rows.filter((row) => row.lifecycle === 'ACTIVE').length), tone: 'success', icon: 'check-circle' },
  ];
}
