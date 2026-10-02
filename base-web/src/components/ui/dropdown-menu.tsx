import * as React from 'react';
import { CheckIcon, CircleIcon } from 'lucide-react';
import { DropdownMenu as Primitive } from 'radix-ui';
import { cn } from '@/lib/utils';

/** shadcn 菜单根组件，方向由使用者的页面语言决定。 */
export function DropdownMenu(props: React.ComponentProps<typeof Primitive.Root>) {
  return <Primitive.Root data-slot="dropdown-menu" {...props} />;
}

/** 将现有按钮作为菜单触发器，保留其键盘语义。 */
export function DropdownMenuTrigger(props: React.ComponentProps<typeof Primitive.Trigger>) {
  return <Primitive.Trigger data-slot="dropdown-menu-trigger" {...props} />;
}

/** 菜单浮层使用语义色并通过 Portal 脱离页面裁切区域。 */
export function DropdownMenuContent({ className, sideOffset = 4, ...props }: React.ComponentProps<typeof Primitive.Content>) {
  return (
    <Primitive.Portal>
      <Primitive.Content
        data-slot="dropdown-menu-content"
        sideOffset={sideOffset}
        className={cn('z-50 min-w-40 rounded-md border bg-popover p-1 text-popover-foreground shadow-lg outline-none data-[state=open]:animate-in data-[state=closed]:animate-out', className)}
        {...props}
      />
    </Primitive.Portal>
  );
}

/** 通用菜单项；可由语言和主题控件复用。 */
export function DropdownMenuItem({ className, ...props }: React.ComponentProps<typeof Primitive.Item>) {
  return (
    <Primitive.Item
      data-slot="dropdown-menu-item"
      className={cn('relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:opacity-50', className)}
      {...props}
    />
  );
}

export function DropdownMenuCheckboxItem({ className, children, ...props }: React.ComponentProps<typeof Primitive.CheckboxItem>) {
  return <Primitive.CheckboxItem data-slot="dropdown-menu-checkbox-item"
    className={cn('relative flex cursor-default select-none items-center rounded-sm py-1.5 pe-2 ps-8 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:opacity-50', className)}
    {...props}>
    <span className="pointer-events-none absolute start-2 flex size-3.5 items-center justify-center">
      <Primitive.ItemIndicator><CheckIcon className="size-3.5" /></Primitive.ItemIndicator>
    </span>
    {children}
  </Primitive.CheckboxItem>;
}

/** 单选菜单组为主题模式提供原生键盘与选中语义。 */
export function DropdownMenuRadioGroup(props: React.ComponentProps<typeof Primitive.RadioGroup>) {
  return <Primitive.RadioGroup data-slot="dropdown-menu-radio-group" {...props} />;
}

/** 主题选项的指示符跟随文档书写方向。 */
export function DropdownMenuRadioItem({ className, children, ...props }: React.ComponentProps<typeof Primitive.RadioItem>) {
  return (
    <Primitive.RadioItem
      data-slot="dropdown-menu-radio-item"
      className={cn('relative flex cursor-default select-none items-center rounded-sm py-1.5 pe-2 ps-8 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[state=checked]:text-primary data-[disabled]:opacity-50', className)}
      {...props}
    >
      <span className="pointer-events-none absolute start-2 flex size-3.5 items-center justify-center">
        <Primitive.ItemIndicator><CircleIcon className="size-2 fill-current" /></Primitive.ItemIndicator>
      </span>
      {children}
    </Primitive.RadioItem>
  );
}
