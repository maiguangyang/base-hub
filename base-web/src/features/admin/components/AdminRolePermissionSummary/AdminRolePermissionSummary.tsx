import { Badge } from '@/components/ui/badge';

interface AdminRolePermissionSummaryProps {
  kind: string;
  permissionCount: number;
}

function permissionScopeLabel(kind: string, permissionCount: number) {
  if (kind === 'HQ_SUPER_ADMIN') return '全部系统权限';
  if (kind === 'FRANCHISE_OWNER') return '全部加盟权限';
  return permissionCount > 0 ? '自定义权限' : '未分配权限';
}

export function AdminRolePermissionSummary({ kind, permissionCount }: AdminRolePermissionSummaryProps) {
  return (
    <div className="flex items-center gap-2">
      <Badge variant="outline" className="bg-muted/40">
        {permissionScopeLabel(kind, permissionCount)}
      </Badge>
      <span className="text-xs text-muted-foreground">{permissionCount} 项</span>
    </div>
  );
}
