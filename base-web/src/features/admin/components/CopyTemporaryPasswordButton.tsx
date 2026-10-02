import { useRef, type RefObject } from 'react';
import { Copy } from 'lucide-react';
import { useAdminToast } from './AdminToast';

function copyWithSelection(password: string, container: HTMLElement): boolean {
  const field = document.createElement('textarea');
  field.value = password;
  field.style.position = 'fixed';
  field.style.opacity = '0';
  container.appendChild(field);
  try {
    field.focus();
    field.select();
    return document.execCommand?.('copy') ?? false;
  } catch {
    return false;
  } finally {
    field.remove();
  }
}

export function CopyTemporaryPasswordButton({ password, failureMessage = '复制失败，请手动选择上方密码复制', iconOnly = false }: { password: string; failureMessage?: string; iconOnly?: boolean }) {
  const buttonRef = useRef<HTMLButtonElement>(null);
  const showToast = useAdminToast();

  async function copy() {
    let copied = false;
    try {
      if (navigator.clipboard?.writeText) {
        await navigator.clipboard.writeText(password);
        copied = true;
      }
    } catch { /* 尝试兼容复制方式。 */ }
    const container = buttonRef.current?.parentElement;
    if (!copied && container) copied = copyWithSelection(password, container);
    showToast(copied ? 'success' : 'error', copied ? '复制成功' : failureMessage);
  }

  return <>
    <CopyTrigger buttonRef={buttonRef} iconOnly={iconOnly} onClick={() => void copy()} />
  </>;
}

function CopyTrigger({ buttonRef, iconOnly, onClick }: { buttonRef: RefObject<HTMLButtonElement | null>; iconOnly: boolean; onClick(): void }) {
  return <button ref={buttonRef} type="button" aria-label={iconOnly ? '复制临时密码' : undefined} title={iconOnly ? '复制临时密码' : undefined} className={iconOnly ? 'inline-flex size-6 items-center justify-center rounded text-muted-foreground hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring' : 'text-sm text-primary underline'} onClick={onClick}>{iconOnly ? <Copy aria-hidden="true" className="size-4" /> : '复制临时密码'}</button>;
}
