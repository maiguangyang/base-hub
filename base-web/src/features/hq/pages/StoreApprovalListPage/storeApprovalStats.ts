import type { AdminStatsItem } from '@/features/admin/components/AdminStatsStrip/types';

export function storeApprovalStats(total: number): AdminStatsItem[] {
  return [{ key: 'pending', label: '待审核门店', value: String(total), tone: 'warning', icon: 'timer' }];
}
