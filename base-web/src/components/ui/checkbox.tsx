import * as React from 'react';
import { CheckIcon, MinusIcon } from 'lucide-react';
import { Checkbox as CheckboxPrimitive } from 'radix-ui';
import { cn } from '@/lib/utils';

export function Checkbox({
  className,
  checked,
  defaultChecked,
  onCheckedChange,
  ...props
}: React.ComponentProps<typeof CheckboxPrimitive.Root>) {
  const [internalChecked, setInternalChecked] = React.useState<boolean | 'indeterminate'>(
    defaultChecked ?? false,
  );
  const resolvedChecked = checked ?? internalChecked;

  return (
    <CheckboxPrimitive.Root
      data-slot="checkbox"
      className={cn(
        'peer size-4 shrink-0 rounded-[4px] border border-input bg-white shadow-xs transition-shadow outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:border-primary data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground',
        className,
        // WEB-UI-014：未勾选固定白底，不跟随页面背景或业务背景类。
        'data-[state=unchecked]:!bg-white',
      )}
      checked={checked}
      defaultChecked={defaultChecked}
      onCheckedChange={(nextChecked) => {
        if (checked === undefined) {
          setInternalChecked(nextChecked);
        }
        onCheckedChange?.(nextChecked);
      }}
      {...props}
    >
      <CheckboxPrimitive.Indicator
        data-slot="checkbox-indicator"
        className="grid place-content-center text-current transition-none"
      >
        {resolvedChecked === 'indeterminate' ? (
          <MinusIcon className="size-3.5" />
        ) : (
          <CheckIcon className="size-3.5" />
        )}
      </CheckboxPrimitive.Indicator>
    </CheckboxPrimitive.Root>
  );
}
