import { useCallback, useEffect, useRef, useState } from 'react';
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from '@/components/ui/sheet';
import { useAuthStore } from '@/features/auth/store/authStore';
import type { ScopeRef } from '@/features/hq/api/paymentConfig';
import { PaymentConfigPanel } from './PaymentConfigPanel';

export type PaymentConfigTarget =
  | { scope: 'FRANCHISE'; organizationId: string; name: string }
  | { scope: 'STORE'; storeId: string; name: string };

export function PaymentConfigDrawer({ target, onClose }: { target: PaymentConfigTarget; onClose(): void }) {
  const [open, setOpen] = useState(true);
  const closed = useRef(false);
  const finishClose = useCallback(() => {
    if (closed.current) return;
    closed.current = true;
    onClose();
  }, [onClose]);
  useEffect(() => {
    if (open) return;
    const timer = window.setTimeout(finishClose, 350);
    return () => window.clearTimeout(timer);
  }, [open, finishClose]);
  const viewer = useAuthStore((state) => state.viewer);
  const canRead = viewer?.currentWorkspace?.workspaceType === 'HEADQUARTERS' && viewer.permissions.includes('paymentConfig:read');
  if (!canRead) return null;
  const ref: ScopeRef = target.scope === 'FRANCHISE'
    ? { scope: 'FRANCHISE', organizationId: target.organizationId }
    : { scope: 'STORE', storeId: target.storeId };
  const kind = target.scope === 'FRANCHISE' ? '加盟商' : '加盟门店';
  return <Sheet open={open} onOpenChange={(next) => { if (!next) setOpen(false); }}>
    <SheetContent side="right" data-side="right" onAnimationEnd={(event) => {
      if (!open && event.target === event.currentTarget) finishClose();
    }} className="w-[min(100vw,44rem)] gap-0 bg-white text-black sm:max-w-[44rem]">
      <SheetHeader className="border-b pr-12">
        <SheetTitle className="text-black">{target.name} · 支付配置</SheetTitle>
        <SheetDescription className="text-slate-600">为该{kind}配置微信或支付宝收款资料；保存后仅完成本地校验。</SheetDescription>
      </SheetHeader>
      <div className="min-h-0 flex-1 overflow-y-auto px-4 py-5">
        <PaymentConfigPanel key={`${ref.scope}:${ref.organizationId ?? ref.storeId}`} target={ref} canManage={viewer.permissions.includes('paymentConfig:manage')} />
      </div>
    </SheetContent>
  </Sheet>;
}
