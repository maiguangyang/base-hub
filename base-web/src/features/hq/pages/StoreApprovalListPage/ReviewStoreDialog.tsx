import { Textarea } from '@/components/ui/textarea';
import { AdminFormDialogShell } from '@/features/admin/components/AdminFormDialogShell';

interface ReviewStoreDialogProps {
  open: boolean;
  onOpenChange(open: boolean): void;
  storeId: string;
  storeName?: string;
  approved: boolean;
  reason?: string;
  onReasonChange?(value: string): void;
  onSubmit(): void | Promise<void>;
  isSubmitting?: boolean;
}

export function availableReviewActions(permissions: readonly string[]) {
  return { approve: permissions.includes('store:approve'), reject: permissions.includes('store:reject') };
}

export function reviewStoreInput(storeId: string, approved: boolean, reason: string) {
  if (approved) return { storeId, approved: true };
  const rejectionReason = reason.trim();
  if (!rejectionReason) throw new Error('VALIDATION_FAILED');
  return { storeId, approved: false, rejectionReason };
}

export function ReviewStoreDialog(props: ReviewStoreDialogProps) {
  return (
    <AdminFormDialogShell open={props.open} onOpenChange={props.onOpenChange} title={props.approved ? '批准门店' : '退回门店'} description={props.approved ? '批准后门店进入启用状态。' : '请填写明确的退回原因。'} submitLabel={props.approved ? '确认批准' : '确认退回'} cancelLabel="取消" isSubmitting={props.isSubmitting} onSubmit={props.onSubmit} maxWidth="lg">
      {props.approved ? <p className="text-sm">确认批准“{props.storeName ?? '当前门店'}”的门店申请？请确认门店资料已核实。</p> : <label className="flex flex-col gap-5 text-sm"><span>退回原因</span><Textarea placeholder="请输入退回原因" value={props.reason ?? ''} onChange={(event) => props.onReasonChange?.(event.target.value)} rows={4} maxLength={512} required /></label>}
    </AdminFormDialogShell>
  );
}
