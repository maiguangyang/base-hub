import type { LucideIcon } from 'lucide-react';
import { TrendingDown, TrendingUp } from 'lucide-react';
import { AdminSparkline } from '@/features/admin/components/AdminSparkline';
import { statusToneClasses, type AdminStatusTone } from '@/features/admin/config/statusTone';
import { cn } from '@/lib/utils';
import type { AdminStatsItem } from './types';

/** 从共享映射取出该语气的文字色，页面与组件都不得手写颜色。 */
function toneText(tone: AdminStatusTone | undefined): string {
  if (!tone) return 'text-foreground';
  return statusToneClasses[tone].split(' ').find((cls) => cls.startsWith('text-')) ?? 'text-foreground';
}

/** 单张 KPI 卡。图标、环比、趋势线均为可选——只读型指标可以只有数值。 */
export interface StatCardProps {
  item: AdminStatsItem;
  Icon?: LucideIcon;
}

export function StatCard({ item, Icon }: StatCardProps) {
  const { label, value, tone, delta, trend } = item;

  return (
    <div className="rounded-lg border border-border bg-card p-4 transition-colors hover:border-ring/40">
      <div className="flex items-start justify-between gap-2">
        <p className="text-sm text-muted-foreground">{label}</p>
        {Icon ? <Icon className="size-4 shrink-0 text-muted-foreground" aria-hidden="true" /> : null}
      </div>
      <div className="mt-2 flex items-end justify-between gap-3">
        <div className="flex items-baseline gap-2">
          <span className={cn('text-2xl font-semibold tabular-nums', toneText(tone))}>{value}</span>
          {delta ? (
            <span className={cn('flex items-center gap-0.5 text-xs tabular-nums', toneText(delta.tone))}>
              {delta.direction === 'up'
                ? <TrendingUp className="size-3" aria-hidden="true" />
                : <TrendingDown className="size-3" aria-hidden="true" />}
              {delta.value}
            </span>
          ) : null}
        </div>
        {trend ? <AdminSparkline points={trend} tone={delta?.tone ?? tone ?? 'info'} label={`${label}趋势`} /> : null}
      </div>
    </div>
  );
}
