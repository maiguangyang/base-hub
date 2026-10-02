import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';
import type { AiApprovedStep } from '@/features/admin/lib/aiStream';
import { AiActionPreviewCard } from './AiActionPreviewCard';

export function AiRiskConfirmationDialog({ open, onOpenChange, steps, onConfirm }: { open: boolean; onOpenChange(open: boolean): void; steps: AiApprovedStep[]; onConfirm(): void }) {
  const memberWrite = steps.length === 1 && customerWriteToolIds.has(steps[0].toolId);
  return <AlertDialog open={open} onOpenChange={onOpenChange}>
    <AlertDialogContent className="max-h-[85vh] overflow-y-auto">
      <AlertDialogHeader>
        <AlertDialogTitle>确认执行后台操作</AlertDialogTitle>
        <AlertDialogDescription>{memberWrite
          ? '请逐项确认目标、积分、金额、规则和凭证编号，并自行核实凭证真实性。本次确认只允许 AI 调用这项会员操作一次。'
          : '请核对目标、范围和全部参数。确认一次后，AI 将按此计划连续执行。'}</AlertDialogDescription>
      </AlertDialogHeader>
      <ol className="space-y-2">{steps.map((step) => <AiActionPreviewCard key={step.sequence} step={step} />)}</ol>
      <AlertDialogFooter>
        <AlertDialogCancel>取消</AlertDialogCancel>
        <AlertDialogAction onClick={onConfirm}>执行已确认操作</AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>;
}

const customerWriteToolIds = new Set([
  'HqCreateCustomerMember', 'HqSetCustomerMemberStatus', 'HqRequestCustomerCancellation',
  'HqCompleteCustomerCancellation', 'HqSaveCustomerBenefitPolicy', 'HqGrantCustomerPoints',
  'HqReverseCustomerPoints', 'HqCorrectCustomerPoints', 'HqCreateCustomerCouponTemplate',
  'HqSetCustomerCouponTemplateEnabled', 'HqGrantCustomerCoupon', 'HqRevokeCustomerCoupon',
]);
