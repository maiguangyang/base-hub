import { useEffect, useRef, useState, type Dispatch, type SetStateAction } from 'react';
import { AIStreamError, readAiStream, type AiEvent } from '@/features/admin/lib/aiStream';

export type AiPhase = 'idle' | 'previewing' | 'awaiting_confirmation' | 'running' | 'completed' | 'interrupted' | 'failed';
export type AiTranscriptEntry =
  | { kind: 'message'; role: 'user' | 'assistant'; text: string }
  | { kind: 'tool'; toolId: string; title: string; sequence: number; status: 'RUNNING' | 'SUCCESS' | 'FAILED' | 'INTERRUPTED' | 'UNKNOWN'; createdResourceId?: string };
export interface AiSessionState {
  phase: AiPhase;
  messages: { role: 'user' | 'assistant'; text: string }[];
  transcript: AiTranscriptEntry[];
  preview?: Extract<AiEvent, { type: 'preview_ready' }>;
  secrets: Extract<AiEvent, { type: 'secret' }>[];
  errorCode?: string;
  interrupted: boolean;
}

const emptyState = (): AiSessionState => ({ phase: 'idle', messages: [], transcript: [], secrets: [], interrupted: false });
const runMarkerKey = (namespace: string) => `admin.ai.run:${namespace}`;

function storedRunMarker(namespace: string): boolean {
	try { return typeof window !== 'undefined' && Boolean(sessionStorage.getItem(runMarkerKey(namespace))); }
	catch { return false; }
}

function saveRunMarker(namespace: string, runId: string) {
	try { sessionStorage.setItem(runMarkerKey(namespace), JSON.stringify({ runId, phase: 'RUN' })); } catch { /* Storage may be unavailable. */ }
}

function clearRunMarker(namespace: string) {
	try { sessionStorage.removeItem(runMarkerKey(namespace)); } catch { /* Storage may be unavailable. */ }
}

export function useAiSession(namespace: string) {
  const [state, setState] = useState<AiSessionState>(() => storedRunMarker(namespace)
	? { ...emptyState(), phase: 'interrupted', interrupted: true }
	: emptyState());
  const active = useRef<AbortController | null>(null);
  const token = useRef<string | undefined>(undefined);
  const generation = useRef(0);
  const previousNamespace = useRef(namespace);
  const runtime = { active, token, generation, setState, namespace };
	const namespaceIsCurrent = previousNamespace.current === namespace;
	const visibleState = namespaceIsCurrent ? state : emptyState();

  const stop = () => stopSession(runtime);
  const clear = () => clearSession(runtime);
  const canRetractLastExchange = namespaceIsCurrent && visibleState.phase !== 'previewing' && visibleState.phase !== 'running' &&
    !storedRunMarker(namespace) && lastUserIndex(visibleState.messages) >= 0;

  useEffect(() => {
    if (previousNamespace.current !== namespace) {
      previousNamespace.current = namespace;
	  invalidateSession(runtime);
	  setState(storedRunMarker(namespace) ? { ...emptyState(), phase: 'interrupted', interrupted: true } : emptyState());
    }
    return () => { generation.current++; active.current?.abort(); active.current = null; token.current = undefined; };
  }, [namespace]);

  async function startPreview(prompt: string, attestation?: { evidenceReference: string; attested: true }, attachmentId?: string) {
    if (!namespaceIsCurrent || active.current || !prompt.trim()) return;
	const history = visibleState.phase === 'completed' || visibleState.phase === 'awaiting_confirmation' ? recentHistory(visibleState.messages) : [];
	clearRunMarker(namespace);
    const controller = new AbortController();
    active.current = controller;
    token.current = undefined;
    const run = ++generation.current;
    setState((current) => ({ ...emptyState(), phase: 'previewing', messages: [...current.messages, { role: 'user', text: prompt.trim() }], transcript: [...current.transcript, { kind: 'message', role: 'user', text: prompt.trim() }] }));
    await runSessionStream(runtime, 'preview', { prompt: prompt.trim(), ...(history.length ? { history } : {}), ...(attestation ? { attestation } : {}), ...(attachmentId ? { attachmentId } : {}) }, controller, run);
  }

  async function confirmRun() {
    const previewToken = token.current;
    if (!namespaceIsCurrent || !previewToken || active.current || visibleState.phase !== 'awaiting_confirmation') return;
    token.current = undefined;
    const controller = new AbortController();
    active.current = controller;
    const run = ++generation.current;
    setState((current) => ({ ...current, phase: 'running', errorCode: undefined, preview: withoutToken(current.preview) }));
	// The request may consume its one-time approval before the first SSE frame arrives.
	saveRunMarker(namespace, '');
    await runSessionStream(runtime, 'run', { previewToken }, controller, run);
  }

  const retractLastExchange = () => retractSessionExchange(runtime, visibleState, canRetractLastExchange);

  return { state: visibleState, startPreview, confirmRun, stop, clear, retractLastExchange, canRetractLastExchange };
}

function lastUserIndex(messages: AiSessionState['messages']): number {
  return messages.reduce((last, message, index) => message.role === 'user' ? index : last, -1);
}

function retractSessionExchange(runtime: SessionRuntime, state: AiSessionState, allowed: boolean): string | undefined {
  if (!allowed) return undefined;
  const index = lastUserIndex(state.messages);
  const prompt = state.messages[index]?.text;
  if (!prompt) return undefined;
  invalidateSession(runtime);
  const messages = state.messages.slice(0, index);
  const transcriptIndex = state.transcript.reduce((last, entry, position) => entry.kind === 'message' && entry.role === 'user' ? position : last, -1);
  runtime.setState({ ...emptyState(), phase: messages.length ? 'completed' : 'idle', messages, transcript: state.transcript.slice(0, transcriptIndex) });
  return prompt;
}

function recentHistory(messages: AiSessionState['messages']): AiSessionState['messages'] {
	const selected: AiSessionState['messages'] = [];
	let bytes = 0;
	for (const message of messages.slice(-12).reverse()) {
		const length = new TextEncoder().encode(message.text).length;
		if (length === 0 || length > 2000 || bytes + length > 6000) continue;
		selected.push(message);
		bytes += length;
	}
	return selected.reverse();
}

interface SessionRuntime {
  active: { current: AbortController | null };
  token: { current: string | undefined };
  generation: { current: number };
  setState: Dispatch<SetStateAction<AiSessionState>>;
	namespace: string;
}

function invalidateSession(runtime: SessionRuntime) {
	runtime.generation.current++;
	runtime.active.current?.abort();
	runtime.active.current = null;
	runtime.token.current = undefined;
}

function stopSession(runtime: SessionRuntime) {
	invalidateSession(runtime);
	runtime.setState((current) => ({ ...current, phase: 'interrupted', interrupted: true, preview: withoutToken(current.preview), transcript: settleRunningTools(current.transcript, 'INTERRUPTED') }));
}

function clearSession(runtime: SessionRuntime) {
	invalidateSession(runtime);
	clearRunMarker(runtime.namespace);
	runtime.setState(emptyState());
}

async function runSessionStream(runtime: SessionRuntime, path: 'preview' | 'run', body: object, controller: AbortController, run: number) {
    try {
      await readAiStream(path, body, controller.signal, (event) => {
        if (runtime.generation.current !== run) return;
        if (path === 'run' && event.type === 'run_started' && event.phase === 'RUN') saveRunMarker(runtime.namespace, event.runId);
		if (path === 'run' && event.type === 'run_finished') clearRunMarker(runtime.namespace);
        if (event.type === 'preview_ready') runtime.token.current = event.previewToken;
		if (event.type === 'run_finished' && event.status !== 'SUCCESS') runtime.token.current = undefined;
        runtime.setState((current) => applyAiEvent(current, event, path));
      });
    } catch (cause) {
	  if (runtime.generation.current === run) {
		runtime.token.current = undefined;
		runtime.setState((current) => failedStream(current, cause, path));
	  }
    } finally {
      if (runtime.active.current === controller) runtime.active.current = null;
    }
}

function withoutToken(preview: AiSessionState['preview']): AiSessionState['preview'] {
  return preview ? { ...preview, previewToken: undefined } : undefined;
}

function failedStream(current: AiSessionState, cause: unknown, path: 'preview' | 'run'): AiSessionState {
  const code = cause instanceof AIStreamError ? cause.code : 'AI_STREAM_INTERRUPTED';
  const interrupted = path === 'run' || code === 'AI_STREAM_INTERRUPTED' || code === 'AI_STREAM_CANCELLED';
  return { ...current, phase: interrupted ? 'interrupted' : 'failed', interrupted, errorCode: code, preview: withoutToken(current.preview), transcript: settleRunningTools(current.transcript, interrupted ? 'INTERRUPTED' : 'FAILED') };
}

function applyAiEvent(current: AiSessionState, event: AiEvent, path: 'preview' | 'run'): AiSessionState {
  switch (event.type) {
    case 'text_delta': return { ...current, messages: appendAssistantText(current.messages, event.text), transcript: appendTranscriptText(current.transcript, event.text) };
    case 'tool_started': return { ...current, transcript: [...current.transcript, { kind: 'tool', toolId: event.toolId, title: event.title, sequence: event.sequence, status: 'RUNNING' }] };
    case 'tool_finished': return { ...current, transcript: finishTranscriptTool(current.transcript, event) };
    case 'secret': return { ...current, secrets: [...current.secrets, event] };
    case 'preview_ready': return { ...current, preview: event };
    case 'error': return { ...current, errorCode: event.code };
    case 'run_finished': return finishEvent(current, event, path);
    default: return current;
  }
}

function finishEvent(current: AiSessionState, event: Extract<AiEvent, { type: 'run_finished' }>, path: 'preview' | 'run'): AiSessionState {
  const transcript = settleRunningTools(current.transcript, event.status === 'SUCCESS' ? 'UNKNOWN' : event.status === 'FAILED' ? 'FAILED' : 'INTERRUPTED');
  if (event.status !== 'SUCCESS') return { ...current, phase: 'failed', preview: withoutToken(current.preview), transcript };
  if (path === 'preview' && current.preview?.previewToken) return { ...current, phase: 'awaiting_confirmation', transcript };
  return { ...current, phase: 'completed', transcript };
}

function appendTranscriptText(entries: AiTranscriptEntry[], text: string): AiTranscriptEntry[] {
  const last = entries.at(-1);
  if (last?.kind === 'message' && last.role === 'assistant') return [...entries.slice(0, -1), { ...last, text: last.text + text }];
  return [...entries, { kind: 'message', role: 'assistant', text }];
}

function finishTranscriptTool(entries: AiTranscriptEntry[], event: Extract<AiEvent, { type: 'tool_finished' }>): AiTranscriptEntry[] {
  const index = entries.findLastIndex((entry) => entry.kind === 'tool' && entry.status === 'RUNNING' && entry.toolId === event.toolId && entry.sequence === event.sequence);
  const finished: AiTranscriptEntry = { kind: 'tool', toolId: event.toolId, title: event.title, sequence: event.sequence, status: event.status === 'SUCCESS' ? 'SUCCESS' : 'FAILED', ...(event.createdResourceId ? { createdResourceId: event.createdResourceId } : {}) };
  if (index < 0) return [...entries, finished];
  return entries.map((entry, position) => position === index ? finished : entry);
}

function settleRunningTools(entries: AiTranscriptEntry[], status: 'FAILED' | 'INTERRUPTED' | 'UNKNOWN'): AiTranscriptEntry[] {
  return entries.map((entry) => entry.kind === 'tool' && entry.status === 'RUNNING' ? { ...entry, status } : entry);
}

function appendAssistantText(messages: AiSessionState['messages'], text: string): AiSessionState['messages'] {
  const last = messages.at(-1);
  if (last?.role === 'assistant') return [...messages.slice(0, -1), { ...last, text: last.text + text }];
  return [...messages, { role: 'assistant', text }];
}
