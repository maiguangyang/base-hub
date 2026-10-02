import type { ReactNode } from 'react';

/** 列表工作台容器。§4.6 规定后台列表页必须采用页面层 + 工作台层的双层骨架，
 *  工作台层由本组件承载 toolbar、列表内容、空态与分页。 */
export interface AdminTableShellProps {
  toolbar: ReactNode;
  children: ReactNode;
  /** 为真时渲染 empty 而非 children。 */
  isEmpty?: boolean;
  /** 空态内容。可含引导性补充操作，但主新增入口仍在 toolbar 的 rightSlot。 */
  empty?: ReactNode;
  pagination?: ReactNode;
}

export function AdminTableShell({ toolbar, children, isEmpty = false, empty, pagination }: AdminTableShellProps) {
  return (
    <div className="rounded-lg border border-border bg-card [&_[data-slot=table-header]]:bg-muted/30 [&_[data-slot=table-head]]:h-12 [&_[data-slot=table-head]]:px-4 [&_[data-slot=table-head]]:font-semibold [&_[data-slot=table-cell]]:px-4 [&_[data-slot=table-cell]]:py-3 [&_[data-slot=table-row]]:hover:bg-muted/20">
      {toolbar}
      {isEmpty ? <div className="p-10 text-center text-sm text-muted-foreground">{empty}</div> : children}
      {pagination ? <div className="border-t border-border">{pagination}</div> : null}
    </div>
  );
}
