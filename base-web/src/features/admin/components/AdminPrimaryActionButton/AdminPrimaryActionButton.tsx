import type { ReactNode } from 'react';
import { Button } from '@/components/ui/button';

/** 资源型列表页的主新增入口。§4.6 规定它必须通过 AdminToolbar 的 rightSlot 注入。 */
export interface AdminPrimaryActionButtonProps {
  children: ReactNode;
  onClick?: () => void;
}

export function AdminPrimaryActionButton({ children, onClick }: AdminPrimaryActionButtonProps) {
  return (
    <Button type="button" className="h-10 px-4" onClick={onClick}>
      {children}
    </Button>
  );
}
