import { useRef } from 'react';
import { Camera, Copy, Undo2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { useAdminToast } from '@/features/admin/components/AdminToast';
import type { AiSessionState, AiTranscriptEntry } from './useAiSession';
import { AiMarkdown } from './AiMarkdown';
import { AiToolSequence } from './AiToolSequence';
import { renderReplyImage, ReplyImageTooLargeError } from './AiReplyImage';

type ToolEntry = Extract<AiTranscriptEntry, { kind: 'tool' }>;
type TranscriptGroup = { kind: 'tools'; firstIndex: number; entries: ToolEntry[] } | { kind: 'message'; index: number; entry: Extract<AiTranscriptEntry, { kind: 'message' }> };

function groupTranscript(transcript: AiTranscriptEntry[]): TranscriptGroup[] {
  const groups: TranscriptGroup[] = [];
  transcript.forEach((entry, index) => {
    const last = groups.at(-1);
    if (entry.kind === 'tool' && last?.kind === 'tools') last.entries.push(entry);
    else if (entry.kind === 'tool') groups.push({ kind: 'tools', firstIndex: index, entries: [entry] });
    else groups.push({ kind: 'message', index, entry });
  });
  return groups;
}

type CompletedReply = { text: string; messageIndexes: number[] };

function completedAssistantReplies(state: AiSessionState): Map<number, CompletedReply> {
  const replies = new Map<number, CompletedReply>();
  let parts: string[] = [];
  let messageIndexes: number[] = [];
  const finish = () => {
    const lastAssistantIndex = messageIndexes.at(-1);
    if (lastAssistantIndex !== undefined) replies.set(lastAssistantIndex, { text: parts.join('\n\n'), messageIndexes });
    parts = [];
    messageIndexes = [];
  };
  state.transcript.forEach((entry, index) => {
    if (entry.kind !== 'message') return;
    if (entry.role === 'user') finish();
    else {
      parts.push(entry.text);
      messageIndexes.push(index);
    }
  });
  if (state.phase === 'completed' || state.phase === 'awaiting_confirmation') finish();
  return replies;
}

export function AiMessageList({ state, onRetract, canRetract }: { state: AiSessionState; onRetract(): void; canRetract: boolean }) {
  const showToast = useAdminToast();
  const listRef = useRef<HTMLDivElement>(null);
  const lastUserIndex = state.transcript.reduce((last, entry, index) => entry.kind === 'message' && entry.role === 'user' ? index : last, -1);
  const assistantReplies = completedAssistantReplies(state);
  async function copy(text: string) {
    let copied = false;
    try {
      if (navigator.clipboard?.writeText) {
        await navigator.clipboard.writeText(text);
        copied = true;
      }
    } catch { /* 尝试兼容复制方式。 */ }
    if (!copied) copied = copyWithSelection(text);
    showToast(copied ? 'success' : 'error', copied ? '消息已复制' : '复制失败，请手动选择消息复制');
  }
  return <TooltipProvider><div ref={listRef} data-ai-conversation aria-live="polite" className="space-y-3">
    {groupTranscript(state.transcript).map((group) => group.kind === 'tools'
      ? <AiToolSequence key={group.firstIndex} entries={group.entries} />
      : <article key={group.index} data-ai-message-index={group.index} className={group.entry.role === 'user' ? 'ml-8' : 'mr-8'}>
      <div data-ai-message-content className={group.entry.role === 'user' ? 'rounded-xl border border-primary/15 bg-primary/10 p-3 text-sm text-foreground whitespace-pre-wrap' : 'rounded-xl bg-muted/60 p-3 text-sm text-foreground'}>{group.entry.role === 'user' ? group.entry.text : <AiMarkdown>{group.entry.text}</AiMarkdown>}</div>
      {(group.entry.role === 'user' || assistantReplies.has(group.index)) && <div className={`mt-1 flex ${group.entry.role === 'user' ? 'justify-end' : 'justify-start'}`}>
        <Tooltip><TooltipTrigger asChild><Button type="button" size="icon-sm" variant="ghost" aria-label="复制消息" className="size-7 text-muted-foreground" onClick={() => void copy(group.entry.role === 'user' ? group.entry.text : assistantReplies.get(group.index)!.text)}><Copy aria-hidden="true" className="size-4" /></Button></TooltipTrigger><TooltipContent side="top">复制消息文字</TooltipContent></Tooltip>
        {assistantReplies.has(group.index) && <Tooltip><TooltipTrigger asChild><Button type="button" size="icon-sm" variant="ghost" aria-label="复制回复图片" className="size-7 text-muted-foreground" onClick={() => void copyReplyImage(listRef.current, assistantReplies.get(group.index)!.messageIndexes, showToast)}><Camera aria-hidden="true" className="size-4" /></Button></TooltipTrigger><TooltipContent side="top">截图并复制本轮回复</TooltipContent></Tooltip>}
        {group.index === lastUserIndex && <Tooltip><TooltipTrigger asChild><Button type="button" size="icon-sm" variant="ghost" aria-label="撤回最后一条消息" disabled={!canRetract} className="size-7 text-muted-foreground" onClick={onRetract}><Undo2 aria-hidden="true" className="size-4" /></Button></TooltipTrigger><TooltipContent side="top">撤回本轮对话；已执行的后台操作不会撤销</TooltipContent></Tooltip>}
      </div>}
    </article>)}
  </div></TooltipProvider>;
}

async function copyReplyImage(list: HTMLDivElement | null, messageIndexes: number[], showToast: ReturnType<typeof useAdminToast>): Promise<void> {
  if (!navigator.clipboard?.write || typeof ClipboardItem === 'undefined') {
    showToast('error', '无法复制图片，请使用支持图片剪贴板的浏览器');
    return;
  }
  let renderError: unknown;
  try {
    let replacedImages = 0;
    // 在点击事件内发起写入，保留浏览器要求的用户操作授权。
    const image = renderReplyImage(list, messageIndexes, () => { replacedImages += 1; }).catch((error: unknown) => {
      renderError = error;
      throw error;
    });
    void image.catch(() => {});
    await navigator.clipboard.write([new ClipboardItem({ 'image/png': image })]);
    showToast('success', replacedImages ? '回复图片已复制，部分图片已用文字代替' : '回复图片已复制');
  } catch (error) {
    const cause = renderError ?? error;
    console.error('复制回复图片失败', cause);
    showToast('error', imageCopyFailureMessage(cause, renderError !== undefined));
  }
}

function imageCopyFailureMessage(error: unknown, renderFailed: boolean): string {
  if (error instanceof ReplyImageTooLargeError) return '回复太长，请使用文字复制';
  if (renderFailed) return '回复图片生成失败，请重试或复制文字';
  if (error instanceof DOMException && error.name === 'NotAllowedError') return '浏览器未允许复制图片，请检查剪贴板权限';
  return '图片写入剪贴板失败，请重试';
}

function copyWithSelection(text: string): boolean {
  const field = document.createElement('textarea');
  field.value = text;
  field.style.position = 'fixed';
  field.style.opacity = '0';
  document.body.appendChild(field);
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
