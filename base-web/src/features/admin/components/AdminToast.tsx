import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from 'react';
import { createPortal } from 'react-dom';
import { CircleAlert, CircleCheck, X } from 'lucide-react';
import { subscribeActionFeedback } from '@/lib/actionFeedback';

type ToastKind = 'success' | 'error';
type Toast = { kind: ToastKind; message: string };
type ShowToast = (kind: ToastKind, message: string) => void;

const ToastContext = createContext<ShowToast | null>(null);

export function useAdminToast(): ShowToast {
  const showToast = useContext(ToastContext);
  if (!showToast) throw new Error('useAdminToast must be used inside AdminToastProvider');
  return showToast;
}

export function AdminToastProvider({ children }: { children: ReactNode }) {
  const [queue, setQueue] = useState<Toast[]>([]);
  const toast = queue[0];
  const dismiss = useCallback(() => setQueue((current) => current.slice(1)), []);
  const showToast = useCallback<ShowToast>((kind, message) => setQueue((current) => [...current, { kind, message }]), []);
  useEffect(() => subscribeActionFeedback(({ kind, message }) => showToast(kind, message)), [showToast]);

  useEffect(() => {
    if (!toast) return;
    const timeout = window.setTimeout(dismiss, 3000);
    return () => window.clearTimeout(timeout);
  }, [toast, dismiss]);

  return <ToastContext.Provider value={showToast}>
    {children}
    {toast && createPortal(<ToastNotice toast={toast} onDismiss={dismiss} />, document.body)}
  </ToastContext.Provider>;
}

function ToastNotice({ toast, onDismiss }: { toast: Toast; onDismiss(): void }) {
  const success = toast.kind === 'success';
  const Icon = success ? CircleCheck : CircleAlert;
  return <div className="pointer-events-none fixed inset-x-4 top-8 z-[60] flex justify-center">
    <div role={success ? 'status' : 'alert'} aria-live={success ? 'polite' : 'assertive'} className={`pointer-events-auto flex min-h-13 w-fit max-w-[min(480px,100%)] items-center gap-3 rounded-sm border px-3 py-3 text-sm shadow-lg ${success ? 'border-success-border bg-success-bg text-success-fg' : 'border-danger-border bg-danger-bg text-danger-fg'}`}>
      <Icon aria-hidden="true" className="size-4 shrink-0" />
      <span className="min-w-0 flex-1 break-words">{toast.message.replace(/[。.!！?？,，;；:：、…\s]+$/u, '')}</span>
      <button type="button" aria-label="关闭提示" className="shrink-0" onClick={onDismiss}><X aria-hidden="true" className="size-4" /></button>
    </div>
  </div>;
}
