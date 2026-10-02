import type { AdminStatsItem } from '@/features/admin/components/AdminStatsStrip/types';

export function auditStats(total: number): AdminStatsItem[] {
  return [{ key: 'total', label: '审计记录', value: String(total), icon: 'list-checks' }];
}
