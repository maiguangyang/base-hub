import { useEffect, useState, type ReactNode, type SubmitEvent } from 'react';
import { X } from 'lucide-react';
import { Dialog } from 'radix-ui';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { actionErrorReason } from '@/lib/graphql/actionErrors';

export type DialogMaxWidth = 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '3xl' | '4xl';

export interface AdminFormDialogShellProps {
  open: boolean;
  onOpenChange(open: boolean): void;
  title: string;
  description?: string;
  children: ReactNode;
  submitLabel: string;
  cancelLabel: string;
  isSubmitting?: boolean;
  submitDisabled?: boolean;
  hideSubmit?: boolean;
  hideCloseButton?: boolean;
  closeOnlyFromFooter?: boolean;
  footerActions?: ReactNode;
  maxWidth?: DialogMaxWidth;
  onSubmit(): void | Promise<void>;
}

const maxWidthClasses: Record<DialogMaxWidth, string> = {
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-xl',
  '2xl': 'max-w-2xl',
  '3xl': 'max-w-3xl',
  '4xl': 'max-w-4xl',
};

function shouldShowCloseButton(props: Pick<AdminFormDialogShellProps, 'closeOnlyFromFooter' | 'hideCloseButton'>) {
  return !props.closeOnlyFromFooter && !props.hideCloseButton;
}



export async function runAdminAction(action: () => void | Promise<void>, setError: (message?: string) => void): Promise<boolean> {
	setError(undefined);
  try {
		await action();
    return true;
  } catch (cause) {
		setError(`操作失败：${actionErrorReason(cause)}`);
    return false;
  }
}

export function AdminActionError({ message }: { message?: string }) {
  return message ? <p role="alert" className="mb-4 text-sm text-destructive">{message}</p> : null;
}

/** 后台新增与编辑表单统一对话框。 */
export function AdminFormDialogShell(props: AdminFormDialogShellProps) {
  const [error, setError] = useState<string>();
  useEffect(() => {
    if (props.open) setError(undefined);
  }, [props.open]);

  async function submit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    await runAdminAction(props.onSubmit, setError);
  }

  return (
    <Dialog.Root open={props.open} onOpenChange={(open) => {
      if (!props.closeOnlyFromFooter || open) props.onOpenChange(open);
    }}>
      <Dialog.Portal>
        <Dialog.Overlay data-testid="admin-form-dialog-overlay" className="fixed inset-0 z-50 bg-black/50" />
        <Dialog.Content
          data-testid="admin-form-dialog-content"
          onEscapeKeyDown={(event) => { if (props.closeOnlyFromFooter) event.preventDefault(); }}
          onPointerDownOutside={(event) => { if (props.closeOnlyFromFooter) event.preventDefault(); }}
          onInteractOutside={(event) => { if (props.closeOnlyFromFooter) event.preventDefault(); }}
          className={cn(
            'fixed left-1/2 top-1/2 z-50 flex max-h-[calc(100dvh-2rem)] w-[calc(100%-2rem)] -translate-x-1/2 -translate-y-1/2 flex-col overflow-hidden rounded-xl border border-border bg-background p-0 shadow-lg',
            maxWidthClasses[props.maxWidth ?? '2xl'] ?? 'max-w-2xl',
          )}
        >
          <header className="shrink-0 px-6 pt-6">
            <Dialog.Title className="text-lg font-semibold">{props.title}</Dialog.Title>
            {props.description && <Dialog.Description className="mt-1 text-sm text-muted-foreground">{props.description}</Dialog.Description>}
          </header>
          <form className="flex min-h-0 flex-1 flex-col" onSubmit={submit}>
            <div data-testid="admin-form-dialog-body" className="min-h-0 flex-1 overflow-y-auto px-6 py-5">
              {props.children}
              {error && <p role="alert" className="mt-4 text-sm text-destructive">{error}</p>}
            </div>
            <footer data-testid="admin-form-dialog-footer" className="flex shrink-0 justify-end gap-2 border-t border-border px-6 py-4">
              {props.footerActions ?? <>{props.closeOnlyFromFooter ?
                <Button type="button" variant="outline" disabled={props.isSubmitting} onClick={() => props.onOpenChange(false)}>{props.cancelLabel}</Button> :
                <Dialog.Close asChild><Button type="button" variant="outline" disabled={props.isSubmitting}>{props.cancelLabel}</Button></Dialog.Close>}
              {!props.hideSubmit && <Button type="submit" disabled={props.isSubmitting || props.submitDisabled}>{props.submitLabel}</Button>}
              </>}
            </footer>
          </form>
          {shouldShowCloseButton(props) && <Dialog.Close className="absolute right-4 top-4 rounded-sm p-1 text-muted-foreground hover:text-foreground" aria-label={props.cancelLabel} disabled={props.isSubmitting}>
            <X className="size-4" />
          </Dialog.Close>}
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
