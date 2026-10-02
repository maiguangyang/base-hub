import { useEffect, useRef, useState } from 'react';
import { LoaderCircle, MessageSquarePlus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from '@/components/ui/sheet';
import type { useAiSession } from './useAiSession';
import { AiActionPreviewCard } from './AiActionPreviewCard';
import { AiMessageList } from './AiMessageList';
import { AiConversationImageCopy } from './AiConversationImageCopy';
import { AiPromptInput } from './AiPromptInput';
import { AiRetractConfirmationDialog } from './AiRetractConfirmationDialog';
import { AiRiskConfirmationDialog } from './AiRiskConfirmationDialog';
import { AiSecretCard } from './AiSecretCard';
import type { AiSessionState } from './useAiSession';

export function AiDrawer({ open, onOpenChange, session }: {
  open: boolean; onOpenChange(open: boolean): void; session: ReturnType<typeof useAiSession>;
}) {
  const actions = useAiDrawerActions(session);
  return <>
    <AiDrawerSheet open={open} onOpenChange={onOpenChange} session={session} actions={actions} />
    <AiDrawerDialogs open={open} state={session.state} actions={actions} />
  </>;
}

type DrawerActions = ReturnType<typeof useAiDrawerActions>;

function AiDrawerSheet({ open, onOpenChange, session, actions }: {
  open: boolean; onOpenChange(open: boolean): void; session: ReturnType<typeof useAiSession>; actions: DrawerActions;
}) {
  const { state } = session;
  const userTurn = state.messages.filter((message) => message.role === 'user').length;
  const { scrollRootRef, contentRef } = useAiScrollFollowing(open, userTurn);
  const busy = state.phase === 'previewing' || state.phase === 'running';
  return <Sheet open={open} onOpenChange={(value) => { if (!value) actions.setRetractOpen(false); onOpenChange(value); }}>
      <SheetContent side="right" className="ai-drawer-surface w-[min(100vw,42rem)] sm:max-w-[42rem] gap-0 bg-white outline-none focus:outline-none focus-visible:outline-none" aria-describedby="ai-drawer-description">
        <SheetHeader className="relative border-b pr-20">
          <SheetTitle>AI 助手</SheetTitle>
          <Button type="button" size="icon-sm" variant="ghost" aria-label="开始新对话" title="开始新对话" className="absolute right-12 top-0" disabled={busy} onClick={actions.newConversation}><MessageSquarePlus aria-hidden="true" /></Button>
          <AiDrawerDescription />
          <AiConversationImageCopy source={() => contentRef.current?.querySelector<HTMLDivElement>('[data-ai-conversation]') ?? null} disabled={state.transcript.length === 0} />
        </SheetHeader>
        <ScrollArea ref={scrollRootRef} className="min-h-0 flex-1 px-4 py-4">
          <div ref={contentRef}>
            <AiMessageList state={state} onRetract={() => actions.setRetractOpen(true)} canRetract={session.canRetractLastExchange && !actions.confirmOpen && !actions.retractOpen} />
            <AiPreviewSection state={state} onExecute={() => state.preview?.isHighRisk ? actions.setConfirmOpen(true) : actions.confirm()} />
            {state.secrets.map((secret, index) => <div key={`${secret.toolId}-${index}`} className="mt-4"><AiSecretCard secret={secret} /></div>)}
            <AiRunStatus state={state} busy={busy} />
          </div>
        </ScrollArea>
        <AiDrawerInput actions={actions} session={session} busy={busy} />
      </SheetContent>
    </Sheet>;
}

function AiDrawerInput({ actions, session, busy }: {
  actions: DrawerActions; session: ReturnType<typeof useAiSession>; busy: boolean;
}) {
  return <AiPromptInput key={actions.inputRevision} prompt={actions.draft} onPromptChange={actions.setDraft} disabled={busy || actions.confirmOpen}
    busy={busy} requiredInputs={session.state.preview?.requiredInputs ?? []}
    onSubmit={actions.send} onStop={session.stop} />;
}

function AiDrawerDialogs({ open, state, actions }: { open: boolean; state: AiSessionState; actions: DrawerActions }) {
  return <>
    {state.preview?.isHighRisk && <AiRiskConfirmationDialog open={actions.confirmOpen} onOpenChange={actions.setConfirmOpen} steps={state.preview.steps} onConfirm={actions.confirm} />}
    <AiRetractConfirmationDialog open={actions.retractOpen && open} onOpenChange={actions.setRetractOpen} onConfirm={actions.retract} />
  </>;
}

function useAiDrawerActions(session: ReturnType<typeof useAiSession>) {
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [retractOpen, setRetractOpen] = useState(false);
  const [draft, setDraft] = useState('');
  const [inputRevision, setInputRevision] = useState(0);
  function confirm() {
    setConfirmOpen(false);
    void session.confirmRun();
  }
  function retract() {
    setRetractOpen(false);
    const prompt = session.retractLastExchange();
    if (prompt !== undefined) {
      setDraft(prompt);
      setInputRevision((current) => current + 1);
    }
  }
  async function send(prompt: string, attestation?: { evidenceReference: string; attested: true }) {
    void session.startPreview(prompt, attestation);
  }
  function newConversation() { session.clear(); setDraft(''); }
  return { confirmOpen, setConfirmOpen, retractOpen, setRetractOpen, draft, setDraft,
    inputRevision, confirm, retract, send, newConversation };
}

function AiDrawerDescription() {
  return <SheetDescription id="ai-drawer-description">
    提问、查询或处理当前工作区事务
  </SheetDescription>;
}

const AUTO_SCROLL_DISTANCE = 96;

function useAiScrollFollowing(open: boolean, userTurn: number) {
  const [scrollRoot, setScrollRoot] = useState<HTMLDivElement | null>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!open) return;
    const viewport = scrollRoot?.querySelector<HTMLElement>('[data-slot="scroll-area-viewport"]');
    const content = contentRef.current;
    if (!viewport || !content) return;
    let following = true;
    let previousTop = viewport.scrollTop;
    const scrollToBottom = () => {
      viewport.scrollTop = viewport.scrollHeight;
      previousTop = viewport.scrollTop;
    };
    const onScroll = () => {
      const top = viewport.scrollTop;
      const distance = viewport.scrollHeight - viewport.clientHeight - top;
      if (distance <= AUTO_SCROLL_DISTANCE) following = true;
      else if (top < previousTop - 1) following = false;
      previousTop = top;
    };
    viewport.addEventListener('scroll', onScroll, { passive: true });
    const observer = typeof ResizeObserver === 'undefined' ? null : new ResizeObserver(() => {
      if (following) scrollToBottom();
    });
    observer?.observe(content);
    scrollToBottom();
    return () => { viewport.removeEventListener('scroll', onScroll); observer?.disconnect(); };
  }, [open, scrollRoot, userTurn]);
  return { scrollRootRef: setScrollRoot, contentRef };
}

function AiPreviewSection({ state, onExecute }: { state: AiSessionState; onExecute(): void }) {
  const preview = state.preview;
  if (!preview) return null;
  const canExecute = state.phase === 'awaiting_confirmation' && Boolean(preview.previewToken);
  return <section className="mt-4 space-y-3" aria-label="操作预览">
    <h3 className="font-medium">操作预览</h3>
    <p className="text-sm text-muted-foreground">{preview.summary}</p>
    <ol className="space-y-2">{preview.steps.map((step) => <AiActionPreviewCard key={step.sequence} step={step} />)}</ol>
    {canExecute && <Button type="button" onClick={onExecute} className="w-full">确认执行</Button>}
  </section>;
}

function AiRunStatus({ state, busy }: { state: AiSessionState; busy: boolean }) {
  if (state.errorCode) return <p role="alert" className="mt-4 rounded-md bg-destructive/10 p-3 text-sm text-destructive">{runErrorMessage(state.errorCode)}（{state.errorCode}）。如涉及写入，请先核对已执行结果，再重新发送。</p>;
  if (state.interrupted) return <p role="alert" className="mt-4 text-sm text-destructive">连接已中断。请核对已执行结果后重新预览。</p>;
  if (busy) return <p className="mt-4 flex items-center gap-2 text-sm text-muted-foreground"><LoaderCircle className="size-4 animate-spin" />正在处理…</p>;
  return null;
}

const runErrorMessages: Record<string, string> = {
  AI_TOKEN_BUDGET_EXCEEDED: '本轮模型调用额度已用完',
  AI_MODEL_CALL_LIMIT: '本轮模型调用次数已达上限',
  AI_TOOL_CALL_LIMIT: '本轮工具调用次数已达上限',
  AI_MODEL_RATE_LIMIT: '模型服务当前请求过多',
  AI_MODEL_AUTH_FAILED: '模型服务拒绝了已配置的密钥',
  AI_MODEL_UNAVAILABLE: '暂时无法连接模型服务',
  AI_MODEL_PROTOCOL: '模型服务返回了无效的流事件或工具调用，请在模型配置中测试连接；反复出现时需修复或更换模型服务',
  AI_MODEL_USAGE_MISSING: '模型服务未返回调用用量',
  AI_MODEL_VERSION_CHANGED: '本轮对话期间模型配置发生了变化',
  AI_SESSION_REVOKED: '当前登录会话已失效',
  AI_RUN_TIMEOUT: '本轮对话处理超时',
  VALIDATION_FAILED: '支付配置参数无效，请核对后重试',
  PERMISSION_DENIED: '当前账号无权操作该支付配置',
  CONFLICT: '支付配置已被修改，请重新查询后重试',
  CONFIG_UNAVAILABLE: '服务端支付配置暂不可用',
};

function runErrorMessage(code: string): string {
  return runErrorMessages[code] ?? '对话未完成';
}
