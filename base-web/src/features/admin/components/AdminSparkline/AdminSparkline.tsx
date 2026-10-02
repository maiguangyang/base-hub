import { statusToneClasses, type AdminStatusTone } from '@/features/admin/config/statusTone';

/** KPI 卡内的迷你趋势线。内联 SVG，不引入图表库。 */
export interface AdminSparklineProps {
  /** 趋势数据点，至少 2 个。数值范围不限，内部会归一化。 */
  points: number[];
  /** 线条语气，取自共享状态映射，不接受自定义颜色。 */
  tone: AdminStatusTone;
  /** 无障碍描述，说明这条趋势代表什么。 */
  label: string;
}

/** 把任意数值序列归一化到 viewBox 坐标。等值序列压平到中线，避免除零。 */
function toPath(points: number[], width: number, height: number): string {
  const max = Math.max(...points);
  const min = Math.min(...points);
  const span = max - min || 1;
  const step = width / (points.length - 1 || 1);
  return points
    .map((value, index) => {
      const x = index * step;
      const y = height - ((value - min) / span) * height;
      return `${index === 0 ? 'M' : 'L'}${x.toFixed(1)},${y.toFixed(1)}`;
    })
    .join(' ');
}

export function AdminSparkline({ points, tone, label }: AdminSparklineProps) {
  if (points.length < 2) return null;
  const width = 96;
  const height = 24;
  const line = toPath(points, width, height);
  const stroke = statusToneClasses[tone].split(' ').find((cls) => cls.startsWith('text-')) ?? 'text-foreground';

  return (
    <svg
      viewBox={`0 0 ${width} ${height}`}
      className={`h-6 w-24 overflow-visible ${stroke}`}
      role="img"
      aria-label={label}
      preserveAspectRatio="none"
    >
      <path d={`${line} L${width},${height} L0,${height} Z`} fill="currentColor" opacity="0.12" />
      <path d={line} fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinejoin="round" strokeLinecap="round" vectorEffect="non-scaling-stroke" />
    </svg>
  );
}
