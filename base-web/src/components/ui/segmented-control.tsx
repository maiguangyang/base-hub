import type { ComponentProps } from 'react';
import { RadioGroup } from 'radix-ui';
import { cn } from '@/lib/utils';

export function SegmentedControl({ className, ...props }: ComponentProps<typeof RadioGroup.Root>) {
  return (
    <RadioGroup.Root
      orientation="horizontal"
      className={cn('flex gap-1 rounded-md border border-input bg-muted p-1', className)}
      {...props}
    />
  );
}

export function SegmentedControlItem({
  className,
  ...props
}: ComponentProps<typeof RadioGroup.Item>) {
  return (
    <RadioGroup.Item
      className={cn(
        'min-h-9 flex-1 rounded px-3 py-2 text-sm font-medium outline-none transition-colors focus-visible:ring-2 focus-visible:ring-ring data-[state=checked]:bg-white data-[state=checked]:shadow-sm dark:data-[state=checked]:bg-background disabled:opacity-50',
        className,
      )}
      {...props}
    />
  );
}
