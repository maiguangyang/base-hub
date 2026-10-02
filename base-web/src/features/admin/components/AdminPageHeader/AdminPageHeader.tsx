import type { ReactNode } from 'react';

/** 后台页面标题区。§4.6 规定它必须留在页面组件中，
 *  不得塞回表格组件内部。actions 承载页头级全局控制，
 *  不是资源主新增入口——后者必须走 AdminToolbar 的 rightSlot。 */
export interface AdminPageHeaderProps {
  title: string;
  description?: string;
  actions?: ReactNode;
}

export function AdminPageHeader({ title, description, actions }: AdminPageHeaderProps) {
  return (
    <div className="pb-4">
      <div className="flex flex-wrap items-start justify-between gap-3 rounded-lg border bg-card px-5 py-4">
        <div className="space-y-1">
          <h1 className="text-xl font-semibold tracking-tight">{title}</h1>
          {description ? <p className="text-sm text-muted-foreground">{description}</p> : null}
        </div>
        {actions ? <div className="flex items-center gap-2">{actions}</div> : null}
      </div>
    </div>
  );
}
