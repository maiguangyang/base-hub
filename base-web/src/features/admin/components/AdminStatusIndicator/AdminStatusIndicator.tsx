import { Badge } from '@/components/ui/badge';
import { statusToneClasses, type AdminStatusTone } from '@/features/admin/config/statusTone';
import { cn } from '@/lib/utils';

/** 只读状态展示。§4.4 规定审计状态、删除状态、生命周期状态必须用它，
 *  禁止伪装成可交互开关。 */
export interface AdminStatusIndicatorProps {
  tone: AdminStatusTone;
  label: string;
}

export function AdminStatusIndicator({ tone, label }: AdminStatusIndicatorProps) {
  return (
    <Badge variant="outline" className={cn('font-normal', statusToneClasses[tone])}>
      {label}
    </Badge>
  );
}
