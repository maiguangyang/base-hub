import { useEffect, useState } from 'react';
import { ChevronDown } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import type { AiTranscriptEntry } from './useAiSession';

type ToolEntry = Extract<AiTranscriptEntry, { kind: 'tool' }>;
const statusLabels: Record<ToolEntry['status'], string> = {
  RUNNING: '执行中', SUCCESS: '已完成', FAILED: '执行失败', INTERRUPTED: '已中断', UNKNOWN: '状态未确认',
};
const statusDots: Record<ToolEntry['status'], string> = {
  RUNNING: 'bg-primary', SUCCESS: 'bg-success-fg', FAILED: 'bg-destructive', INTERRUPTED: 'bg-warning-fg', UNKNOWN: 'bg-muted-foreground',
};
const statusText: Record<ToolEntry['status'], string> = {
  RUNNING: 'text-primary', SUCCESS: 'text-success-fg', FAILED: 'text-destructive', INTERRUPTED: 'text-warning-fg', UNKNOWN: 'text-muted-foreground',
};

export function AiToolSequence({ entries }: { entries: ToolEntry[] }) {
  const [expanded, setExpanded] = useState(false);
  const [front, setFront] = useState(0);
  const [moving, setMoving] = useState(false);
  const hasNext = front < entries.length - 1;
  useEffect(() => {
    if (expanded || !hasNext || moving) return;
    const wait = entries.length - front > 2 ? 100 : 250;
    const timer = window.setTimeout(() => setMoving(true), wait);
    return () => window.clearTimeout(timer);
  }, [entries.length, expanded, front, hasNext, moving]);
  useEffect(() => {
    if (!moving) return;
    const timer = window.setTimeout(() => { setFront((current) => current + 1); setMoving(false); }, 220);
    return () => window.clearTimeout(timer);
  }, [moving]);
  if (entries.length === 1) return <AiToolPill entry={entries[0]} />;
  function changeExpanded(value: boolean) {
    setExpanded(value);
    setFront(entries.length - 1);
    setMoving(false);
  }
  const phase = hasNext ? moving ? 'moving' : 'waiting' : 'settled';
  return <Collapsible open={expanded} onOpenChange={changeExpanded} data-ai-tool-group={entries.length}>
    <div className="inline-flex max-w-full items-start gap-2" data-ai-tool-stack data-phase={phase}>
      <div className="grid min-w-0 items-center">
        <div className={`ai-tool-stack-face ai-tool-stack-layered relative col-start-1 row-start-1 inline-flex flex-col items-start transition-opacity duration-200 ${expanded ? 'pointer-events-none opacity-0' : 'opacity-100'}`} data-ai-tool-stack-layers data-layer-count={Math.min(entries.length, 3)} aria-hidden={expanded}>
          <div className="ai-tool-stack-front" data-phase={phase}><AiToolPill entry={entries[front]} /></div>
          {hasNext && <div className="ai-tool-stack-incoming mt-1" data-phase={phase}><AiToolPill entry={entries[front + 1]} /></div>}
        </div>
        <div className={`ai-tool-stack-face col-start-1 row-start-1 inline-flex min-h-8 items-center rounded-full border bg-card px-3 text-xs text-muted-foreground shadow-xs transition-opacity duration-200 ${expanded ? 'opacity-100' : 'pointer-events-none opacity-0'}`} aria-hidden={!expanded}>连续工具调用</div>
      </div>
      <CollapsibleTrigger asChild>
        <Button type="button" size="sm" variant="ghost" aria-label={`${expanded ? '收起工具调用' : '展开全部工具调用'}，共 ${entries.length} 次`} aria-expanded={expanded} className="h-8 shrink-0 gap-1 rounded-full px-2 text-muted-foreground hover:bg-transparent hover:text-muted-foreground dark:hover:bg-transparent" title={expanded ? '收起工具调用' : '展开全部工具调用'}>
          <span className="text-xs tabular-nums">{entries.length} 次</span>
          <ChevronDown aria-hidden="true" className={`size-4 transition-transform duration-200 ${expanded ? 'rotate-180' : ''}`} />
        </Button>
      </CollapsibleTrigger>
    </div>
    <CollapsibleContent className="ai-tool-sequence-content overflow-hidden">
      <div className="space-y-2 pt-2">
        {entries.map((entry) => <AiToolPill key={entry.sequence} entry={entry} />)}
      </div>
    </CollapsibleContent>
  </Collapsible>;
}

function AiToolPill({ entry }: { entry: ToolEntry }) {
  const running = entry.status === 'RUNNING';
  return <div className="flex">
    <div data-ai-tool={entry.toolId} data-state={entry.status.toLowerCase()} className={`relative inline-flex min-h-8 max-w-full items-center gap-2 overflow-hidden rounded-full border bg-card px-3 py-1 text-xs shadow-xs ${running ? 'ai-tool-running border-primary/40 text-foreground' : 'border-border text-muted-foreground'}`}>
      <span aria-hidden="true" className={`relative z-10 size-1.5 shrink-0 rounded-full ${statusDots[entry.status]}`} />
      <span className="relative z-10 truncate" title={entry.title}>{entry.title}</span>
      <span className={`relative z-10 shrink-0 ${statusText[entry.status]}`}>{statusLabels[entry.status]}</span>
      {entry.createdResourceId && <span className="relative z-10 truncate opacity-75">新建编号 {entry.createdResourceId}</span>}
    </div>
  </div>;
}
