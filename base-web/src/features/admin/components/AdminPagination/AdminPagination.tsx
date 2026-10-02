import { ChevronLeft, ChevronRight } from 'lucide-react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

/** 列表分页控件。受控组件，页码与每页条数均由调用方持有。 */
export interface AdminPaginationProps {
  /** 过滤后的总条数，用于显示区间与计算总页数。 */
  total: number;
  /** 当前页，从 1 开始。 */
  page: number;
  pageSize: number;
  onPageChange: (next: number) => void;
  onPageSizeChange: (next: number) => void;
}

/** 可选的每页条数。后台列表常用区间，不做成配置项以免增加无消费者的灵活度。 */
const pageSizeOptions = [10, 20, 50];

export function AdminPagination({ total, page, pageSize, onPageChange, onPageSizeChange }: AdminPaginationProps) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const from = total === 0 ? 0 : (page - 1) * pageSize + 1;
  const to = Math.min(page * pageSize, total);

  return (
    <div className="flex flex-wrap items-center justify-between gap-3 bg-muted/30 px-4 py-3 text-sm text-muted-foreground">
      <span className="tabular-nums">共 {total} 条，当前 {from}–{to}</span>
      <div className="flex items-center gap-3">
        <div className="flex items-center gap-2">
          <span>每页</span>
          <Select value={String(pageSize)} onValueChange={(value) => onPageSizeChange(Number(value))}>
            <SelectTrigger aria-label="每页" className="h-8 w-[4.5rem] bg-card px-2 py-1 shadow-none">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {pageSizeOptions.map((size) => <SelectItem key={size} value={String(size)}>{size}</SelectItem>)}
            </SelectContent>
          </Select>
        </div>
        <div className="flex items-center gap-1">
          <button type="button" aria-label="上一页" disabled={page <= 1}
            className="flex size-8 items-center justify-center rounded-md border border-input bg-card disabled:opacity-40 hover:bg-accent"
            onClick={() => onPageChange(page - 1)}>
            <ChevronLeft className="size-4" />
          </button>
          <span className="min-w-16 text-center tabular-nums text-foreground">{page} / {totalPages}</span>
          <button type="button" aria-label="下一页" disabled={page >= totalPages}
            className="flex size-8 items-center justify-center rounded-md border border-input bg-card disabled:opacity-40 hover:bg-accent"
            onClick={() => onPageChange(page + 1)}>
            <ChevronRight className="size-4" />
          </button>
        </div>
      </div>
    </div>
  );
}
