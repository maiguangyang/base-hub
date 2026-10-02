import { AdminSearchSelect } from '@/features/admin/components/AdminSearchSelect';
import type { AdminFilterOption } from '@/features/admin/components/AdminFilterSelect';

interface FranchiseSearchSelectProps {
  value: string | undefined;
  options: readonly AdminFilterOption[];
  disabled: boolean;
  onValueChange(value: string | undefined): void;
}

export function FranchiseSearchSelect(props: FranchiseSearchSelectProps) {
  return <AdminSearchSelect {...props} label="全部加盟商" placeholder="全部加盟商" searchPlaceholder="搜索加盟商" emptyMessage="没有匹配的加盟商" clearLabel="全部加盟商" compact clearable />;
}
