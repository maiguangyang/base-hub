import type { ReactNode } from 'react';
import { Button } from '@/components/ui/button';

export interface AdminListFiltersProps {
  children: ReactNode;
  hasActiveFilters: boolean;
  onReset: () => void;
}

/** 后台列表统一筛选区，只负责布局和清除入口。 */
export function AdminListFilters({ children, hasActiveFilters, onReset }: AdminListFiltersProps) {
  return (
    <div className="flex min-w-0 flex-1 flex-wrap items-center gap-3">
      {children}
      {hasActiveFilters ? (
        <Button type="button" variant="ghost" className="h-10" onClick={onReset}>重置筛选</Button>
      ) : null}
    </div>
  );
}
