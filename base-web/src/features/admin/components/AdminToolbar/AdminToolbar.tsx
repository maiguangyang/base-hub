import type { ReactNode } from 'react';

/** 列表工作台的筛选栏。§4.6 规定资源主新增动作必须经 rightSlot 注入，
 *  不得另开 actions 一类的兼容口子。 */
export interface AdminToolbarProps {
  children?: ReactNode;
  rightSlot?: ReactNode;
}

export function AdminToolbar({ children, rightSlot }: AdminToolbarProps) {
  return (
    <div className="m-4 flex flex-wrap items-center gap-3 rounded-lg border border-border bg-muted/30 p-4">
      <div className="min-w-0 flex-1 flex flex-wrap items-center gap-3">{children}</div>
      {rightSlot ? <div className="ml-auto flex items-center gap-2">{rightSlot}</div> : null}
    </div>
  );
}
