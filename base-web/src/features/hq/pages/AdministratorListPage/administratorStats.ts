import type { AdministratorRow } from './AdministratorTable';

export function administratorStats(rows: AdministratorRow[], total: number) {
  return [
    { key: 'total', label: '管理员', value: String(total), icon: 'users' },
    { key: 'active', label: '当前页已启用', value: String(rows.filter((row) => row.status === 'ACTIVE').length), icon: 'list-checks' },
  ];
}
