import type { HqFranchisesQuery } from '@/__generated__/graphql';
import type { AdminStatsItem } from '@/features/admin/components/AdminStatsStrip/types';

export function franchiseStats(rows: NonNullable<HqFranchisesQuery['organizations']>['data'], total: number): AdminStatsItem[] {
  return [
    { key: 'total', label: '加盟商', value: String(total), icon: 'users' },
    { key: 'active', label: '正常', value: String(rows.filter((row) => row.status === 'ACTIVE').length), tone: 'success', icon: 'user-check' },
    { key: 'suspended', label: '已暂停', value: String(rows.filter((row) => row.status === 'SUSPENDED').length), tone: 'danger', icon: 'timer' },
  ];
}
