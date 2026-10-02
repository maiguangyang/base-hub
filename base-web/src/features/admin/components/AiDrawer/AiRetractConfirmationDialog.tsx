import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';

export function AiRetractConfirmationDialog({ open, onOpenChange, onConfirm }: { open: boolean; onOpenChange(open: boolean): void; onConfirm(): void }) {
  return <AlertDialog open={open} onOpenChange={onOpenChange}>
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>确认撤回本轮对话</AlertDialogTitle>
        <AlertDialogDescription>本轮消息和工具记录将从当前页面移除，发送内容会回填输入框。已执行的后台操作不会撤销。</AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel>取消</AlertDialogCancel>
        <AlertDialogAction onClick={onConfirm}>确认撤回</AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>;
}
