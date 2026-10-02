import type { AdminStatusTone } from '@/features/admin/config/statusTone';

/** KPI 卡的环比变化。方向与语气分开：待处理任务下降是好事，
 *  所以不能由正负号推导颜色，必须由调用方显式指定。 */
export interface AdminStatsDelta {
  value: string;
  direction: 'up' | 'down';
  tone: AdminStatusTone;
}

/** 单个统计项。tone / delta / trend 均为可选。 */
export interface AdminStatsItem {
  key: string;
  label: string;
  value: string;
  tone?: AdminStatusTone;
  delta?: AdminStatsDelta;
  /** 迷你趋势线数据点。本轮为 mock，页面已标注数据来源。 */
  trend?: number[];
  /** lucide 图标键，由 statsIcons 映射到组件。 */
  icon?: string;
}
