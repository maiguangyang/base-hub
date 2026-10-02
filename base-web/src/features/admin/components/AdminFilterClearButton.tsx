import { X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

/** 宿主使用 group/filter：有值时挂载，悬停或键盘聚焦时显示。 */
export function AdminFilterClearButton({ label, onClear, disabled, className }: {
  label: string;
  onClear(): void;
  disabled?: boolean;
  className?: string;
}) {
  return <Button type="button" variant="ghost" size="icon" aria-label={`清空${label}`} title={`清空${label}`}
    disabled={disabled}
    className={cn('absolute right-1 top-1/2 size-7 -translate-y-1/2 text-muted-foreground opacity-0 pointer-events-none disabled:opacity-0 group-hover/filter:opacity-100 group-hover/filter:pointer-events-auto group-hover/filter:disabled:opacity-50 group-hover/filter:disabled:pointer-events-none group-has-[:focus-visible]/filter:opacity-100 focus-visible:opacity-100 focus-visible:pointer-events-auto', className)}
    onClick={(event) => { event.stopPropagation(); onClear(); }}>
    <X aria-hidden="true" className="size-3.5" />
  </Button>;
}
