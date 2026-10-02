import { useRef, useState, type SyntheticEvent } from 'react';
import { ArrowUp, Plus, Square } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Checkbox } from '@/components/ui/checkbox';

export function AiPromptInput({ prompt, onPromptChange, disabled, busy, requiredInputs, attachment, allowAttachment = true, onAttachment, onSubmit, onStop }: { prompt: string; onPromptChange(value: string): void; disabled: boolean; busy: boolean; requiredInputs: string[]; attachment?: File; allowAttachment?: boolean; onAttachment(file?: File): void; onSubmit(prompt: string, attestation?: { evidenceReference: string; attested: true }): void | Promise<void>; onStop(): void }) {
  const [evidence, setEvidence] = useState('');
  const [attested, setAttested] = useState(false);
  const needsEvidence = requiredInputs.includes('evidenceReference');
  const canSend = !disabled && Boolean(prompt.trim()) && (!needsEvidence || Boolean(evidence.trim() && attested));
  async function submit(event: SyntheticEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!prompt.trim() || disabled || needsEvidence && (!evidence.trim() || !attested)) return;
    try {
      await onSubmit(prompt.trim(), needsEvidence ? { evidenceReference: evidence.trim(), attested: true } : undefined);
      onPromptChange('');
    } catch { /* The parent displays upload errors and leaves the draft intact. */ }
  }
  return <form onSubmit={submit} className="flex flex-col gap-5 border-t p-3 sm:p-4">
    {needsEvidence && <div className="flex flex-col gap-5 rounded-md border p-3 text-sm">
      <label htmlFor="ai-evidence">历史账号认定依据</label>
      <input id="ai-evidence" value={evidence} onChange={(event) => setEvidence(event.target.value)} placeholder="请输入历史开通记录编号" className="w-full rounded-md border px-3 py-2" />
      <label className="flex items-start gap-2"><Checkbox checked={attested} onCheckedChange={(value) => setAttested(value === true)} />我已核实这项认定依据</label>
    </div>}
    <AiAttachmentPreview allowAttachment={allowAttachment} attachment={attachment} onAttachment={onAttachment} />
    <div className="rounded-2xl border border-border bg-white text-slate-900 shadow-sm">
      <label htmlFor="ai-prompt" className="sr-only">向 AI 助手发送消息</label>
      <Textarea id="ai-prompt" value={prompt} onChange={(event) => onPromptChange(event.target.value)} placeholder="向 AI 助手提问，或描述要处理的后台事务…" disabled={disabled} className="min-h-20 max-h-48 resize-none border-0 bg-transparent px-3 pt-3 pb-1 shadow-none focus-visible:border-0 focus-visible:ring-0" onKeyDown={(event) => { if (event.key === 'Enter' && !event.shiftKey) { event.preventDefault(); event.currentTarget.form?.requestSubmit(); } }} />
      <div className="flex items-center justify-between gap-2 px-2 pb-2">
        <div className="flex min-w-0 items-center gap-2">
          <AiAttachmentButton allowAttachment={allowAttachment} onAttachment={onAttachment} disabled={disabled} />
          <span className="truncate text-xs text-muted-foreground">写入操作需确认</span>
        </div>
        {busy
          ? <Button type="button" size="icon-sm" aria-label="终止生成" title="终止生成" className="rounded-full" onClick={onStop}><Square aria-hidden="true" className="size-3 fill-current stroke-0" /></Button>
          : <Button type="submit" size="icon-sm" aria-label="发送消息" title="发送消息" className={`rounded-full ${canSend ? '' : 'bg-primary/40 hover:bg-primary/40 disabled:opacity-100'}`} disabled={!canSend}><ArrowUp aria-hidden="true" /></Button>}
      </div>
    </div>
  </form>;
}

function AiAttachmentPreview({ allowAttachment, attachment, onAttachment }: {
  allowAttachment: boolean; attachment?: File; onAttachment(file?: File): void;
}) {
  if (!allowAttachment || !attachment) return null;
  return <div className="flex items-center justify-between rounded-md border px-3 py-2 text-xs"><span className="truncate">图片：{attachment.name}</span>
    <Button type="button" size="sm" variant="ghost" onClick={() => onAttachment(undefined)}>移除附件</Button></div>;
}

function AiAttachmentButton({ allowAttachment, onAttachment, disabled }: {
  allowAttachment: boolean; onAttachment(file?: File): void; disabled: boolean;
}) {
  const fileInput = useRef<HTMLInputElement>(null);
  if (!allowAttachment) return null;
  return <><input ref={fileInput} type="file" accept="image/jpeg,image/png,image/webp" className="sr-only" aria-label="选择商品图片"
    onChange={(event) => { onAttachment(event.target.files?.[0]); event.target.value = ''; }} />
    <Button type="button" size="icon-sm" variant="ghost" aria-label="添加图片附件" title="添加图片附件"
      onClick={() => fileInput.current?.click()} disabled={disabled}><Plus aria-hidden="true" /></Button></>;
}
