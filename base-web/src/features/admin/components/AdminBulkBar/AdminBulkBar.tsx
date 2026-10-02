import type { ReactNode } from 'react';

/** 批量操作条。有选中项时出现在工具栏与表格之间，
 *  避免用户为了批量操作逐行重复点击（§ Bulk Actions）。 */
export interface AdminBulkBarProps {
  count: number;
  onClear: () => void;
  /** 具体的批量动作按钮，由页面提供——动作语义因资源而异。 */
  children?: ReactNode;
}

export function AdminBulkBar({ count, onClear, children }: AdminBulkBarProps) {
  if (count === 0) return null;
  return (
    <div className="flex flex-wrap items-center gap-3 border-b border-border bg-accent/40 px-3 py-2 text-sm">
      <span className="tabular-nums">已选 {count} 项</span>
      <div className="flex items-center gap-2">{children}</div>
      <button
        type="button"
        onClick={onClear}
        className="ml-auto rounded-sm text-muted-foreground outline-none hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring"
      >
        取消选择
      </button>
    </div>
  );
}
