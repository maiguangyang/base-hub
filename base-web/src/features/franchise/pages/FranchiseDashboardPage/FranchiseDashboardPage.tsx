import { useQuery } from '@apollo/client/react';
import { AdminPageHeader } from '@/features/admin/components/AdminPageHeader';
import { AdminStatsStrip } from '@/features/admin/components/AdminStatsStrip';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { FRANCHISE_ROLES_QUERY } from '../../graphql/roles';
import { FRANCHISE_STAFF_QUERY } from '../../graphql/staff';
import { FRANCHISE_STORES_QUERY } from '../../graphql/stores';

export function FranchiseDashboardPage() {
  useAdminTab();
  const variables = { page: 1, pageSize: 1, q: null, filter: undefined };
  const stores = useQuery(FRANCHISE_STORES_QUERY, { variables });
  const staff = useQuery(FRANCHISE_STAFF_QUERY, { variables });
  const roles = useQuery(FRANCHISE_ROLES_QUERY, { variables });
  const stats = [
    { key: 'stores', label: '门店', value: String(stores.data?.stores?.total ?? 0), icon: 'list-checks' },
    { key: 'staff', label: '成员', value: String(staff.data?.operatorMemberships?.total ?? 0), icon: 'users' },
    { key: 'roles', label: '角色', value: String(roles.data?.operatorRoles?.total ?? 0), icon: 'user-check' },
  ];
  return <div><AdminPageHeader title="加盟商概览" description="数据仅来自当前加盟组织。" /><AdminStatsStrip items={stats} columns={3} /></div>;
}
