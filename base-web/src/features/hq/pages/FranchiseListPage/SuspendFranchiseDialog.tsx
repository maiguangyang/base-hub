import { useState, type SubmitEvent } from 'react';
import { AlertDialog, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import type { FranchiseRow } from './FranchiseTable';

/** 由列表持有待暂停行；提交结果决定是否关闭，失败时保留输入。 */
interface SuspendFranchiseDialogProps {
  row: FranchiseRow;
  pending: boolean;
  error?: string;
  onConfirm(reason: string): Promise<boolean>;
  onClose(): void;
}

/** 暂停前采集审计原因；服务端确认成功后才关闭对话框。 */
export function SuspendFranchiseDialog({ row, pending, error, onConfirm, onClose }: SuspendFranchiseDialogProps) {
  const [reason, setReason] = useState('');

  /** 原因只在非空且未提交中时发送，成功后再关闭弹窗。 */
  async function submit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    if (pending || !reason.trim()) return;
    if (await onConfirm(reason.trim())) onClose();
  }

  return <AlertDialog open onOpenChange={(open) => { if (!open && !pending) onClose(); }}>
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>暂停{row.name}</AlertDialogTitle>
        <AlertDialogDescription>暂停后将撤销该加盟商当前在线会话，请填写原因后确认。</AlertDialogDescription>
      </AlertDialogHeader>
      <form className="flex flex-col gap-5" onSubmit={(event) => void submit(event)}>
        <label className="flex flex-col gap-2 text-sm">
          <span>暂停原因</span>
          <Textarea className="min-h-24" placeholder="请输入暂停原因" value={reason} onChange={(event) => setReason(event.target.value)} maxLength={64} rows={3} required />
        </label>
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
        <AlertDialogFooter>
          <AlertDialogCancel type="button" disabled={pending}>取消</AlertDialogCancel>
          <Button type="submit" disabled={pending || !reason.trim()}>{pending ? '暂停中…' : '确认暂停'}</Button>
        </AlertDialogFooter>
      </form>
    </AlertDialogContent>
  </AlertDialog>;
}
