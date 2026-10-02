import type { FranchiseStaffQuery } from '@/__generated__/graphql';
import type { AdminStatsItem } from '@/features/admin/components/AdminStatsStrip/types';

export function staffStats(rows: NonNullable<FranchiseStaffQuery['operatorMemberships']>['data'], total: number): AdminStatsItem[] {
  return [
    { key: 'total', label: '成员', value: String(total), icon: 'users' },
    { key: 'active', label: '活跃', value: String(rows.filter((row) => row.status === 'ACTIVE').length), tone: 'success', icon: 'user-check' },
    { key: 'invited', label: '待接受邀请', value: String(rows.filter((row) => row.status === 'INVITED').length), tone: 'warning', icon: 'timer' },
  ];
}
