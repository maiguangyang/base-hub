import type { FranchiseRolesQuery } from '@/__generated__/graphql';
import type { AdminStatsItem } from '@/features/admin/components/AdminStatsStrip/types';

export function roleStats(rows: NonNullable<FranchiseRolesQuery['operatorRoles']>['data'], total: number): AdminStatsItem[] {
  return [
    { key: 'total', label: '角色', value: String(total), icon: 'users' },
    { key: 'custom', label: '自定义角色', value: String(rows.filter((row) => row.kind === 'CUSTOM').length), icon: 'list-checks' },
  ];
}
