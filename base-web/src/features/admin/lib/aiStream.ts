export type AiEvent =
  | { type: 'run_started'; runId: string; phase: 'PREVIEW' | 'RUN' }
  | { type: 'text_delta'; text: string }
  | { type: 'tool_started'; toolId: string; title: string; sequence: number }
  | { type: 'tool_finished'; toolId: string; title: string; sequence: number; status: string; createdResourceId?: string }
  | { type: 'preview_ready'; summary: string; steps: AiApprovedStep[]; isHighRisk: boolean; requiredInputs: string[]; previewToken?: string }
  | { type: 'secret'; toolId: string; targetAccountId: string; value: string }
  | { type: 'run_finished'; status: 'SUCCESS' | 'FAILED' | 'CANCELLED' }
  | { type: 'error'; code: string };

export interface AiApprovedStep {
  toolId: string;
  title: string;
  targetIds?: string[];
  scopeIds?: string[];
  arguments: Record<string, unknown>;
  maxCalls: number;
  sequence: number;
  risk: string;
}

export class AIStreamError extends Error {
  constructor(readonly code: string) { super(code); }
}

export async function readAiStream(
  path: 'preview' | 'run', body: object, signal: AbortSignal, onEvent: (event: AiEvent) => void,
): Promise<void> {
  let response: Response;
  try {
    response = await fetch(`${appConfig.engine.httpUrl}/api/ai/${path}`, {
      method: 'POST', credentials: 'include', cache: 'no-store', signal,
      headers: { Accept: 'text/event-stream', 'Content-Type': 'application/json' }, body: JSON.stringify(body),
    });
  } catch {
    throw new AIStreamError(signal.aborted ? 'AI_STREAM_CANCELLED' : 'AI_NETWORK_ERROR');
  }
  if (!response.ok) {
    const payload = await response.json().catch(() => null) as { code?: unknown } | null;
    const code = typeof payload?.code === 'string' && /^[A-Z][A-Z0-9_]{1,80}$/.test(payload.code) ? payload.code : 'AI_REQUEST_FAILED';
    throw new AIStreamError(code);
  }
  if (!response.body || !response.headers.get('Content-Type')?.startsWith('text/event-stream')) {
    throw new AIStreamError('AI_STREAM_PROTOCOL');
  }
  await readSseEvents(response.body, signal, onEvent);
}
import { appConfig } from '@/config/app';
import { readSseEvents } from './aiStreamParser';
