import { ArrowDown, ArrowUp, ArrowUpDown } from 'lucide-react';
import { TableHead } from '@/components/ui/table';
import type { SortDirection } from '@/features/admin/hooks/useAdminTableState';

/** 可排序表头。aria-sort 必须如实反映当前排序态，屏幕阅读器据此播报。 */
export interface AdminSortableHeadProps {
  label: string;
  /** 本列是否为当前排序列。 */
  active: boolean;
  direction: SortDirection;
  onToggle: () => void;
}

export function AdminSortableHead({ label, active, direction, onToggle }: AdminSortableHeadProps) {
  const sort = active && direction ? (direction === 'asc' ? 'ascending' : 'descending') : 'none';
  const Icon = !active || !direction ? ArrowUpDown : direction === 'asc' ? ArrowUp : ArrowDown;

  return (
    <TableHead aria-sort={sort}>
      <button
        type="button"
        onClick={onToggle}
        className="flex items-center gap-1 rounded-sm outline-none hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring"
      >
        {label}
        <Icon className={active && direction ? 'size-3.5 text-foreground' : 'size-3.5 opacity-40'} aria-hidden="true" />
      </button>
    </TableHead>
  );
}
