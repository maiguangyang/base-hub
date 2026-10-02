import type { RoleRow } from './RoleTable';

export function roleStats(rows: RoleRow[], total: number) {
  return [
    { key: 'total', label: '角色', value: String(total), icon: 'shield-check' },
    { key: 'custom', label: '当前页自定义角色', value: String(rows.filter((row) => row.kind === 'CUSTOM').length), icon: 'users' },
  ];
}
