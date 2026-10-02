import type { AdminStatsItem } from '@/features/admin/components/AdminStatsStrip/types';

export function franchiseStoreStats(total: number): AdminStatsItem[] {
  return [{ key: 'total', label: '加盟门店', value: String(total), icon: 'list-checks' }];
}
