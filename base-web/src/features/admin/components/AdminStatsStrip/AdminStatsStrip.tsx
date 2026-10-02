import { StatCard } from './StatCard';
import { statsIcons } from './statsIcons';
import type { AdminStatsItem } from './types';

/** 后台列表页的页面级摘要。§4.6 规定它是统一标准，
 *  不再允许保留自定义 summary 卡片作为长期结构特例。 */
export interface AdminStatsStripProps {
  items: AdminStatsItem[];
  /** 每行列数。列表页用 4 列，仪表盘等宽屏场景可传 4 保持一致。 */
  columns?: 3 | 4;
}

export function AdminStatsStrip({ items, columns = 4 }: AdminStatsStripProps) {
  const grid = columns === 3 ? 'sm:grid-cols-2 lg:grid-cols-3' : 'sm:grid-cols-2 lg:grid-cols-4';
  return (
    <div className={`grid gap-3 pb-4 ${grid}`}>
      {items.map((item) => (
        <StatCard key={item.key} item={item} Icon={item.icon ? statsIcons[item.icon] : undefined} />
      ))}
    </div>
  );
}
