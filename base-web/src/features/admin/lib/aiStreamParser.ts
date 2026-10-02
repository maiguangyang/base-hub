import { AIStreamError, type AiEvent } from './aiStream';

const frameLimit = 64 * 1024;
const eventNames = new Set(['run_started', 'text_delta', 'tool_started', 'tool_finished', 'preview_ready', 'secret', 'run_finished', 'error']);

export async function readSseEvents(stream: ReadableStream<Uint8Array>, signal: AbortSignal, onEvent: (event: AiEvent) => void): Promise<void> {
	let finished = false;
	try { finished = await pumpSseEvents(stream, onEvent); }
	catch (cause) {
		if (cause instanceof AIStreamError) throw cause;
		throw new AIStreamError(signal.aborted ? 'AI_STREAM_CANCELLED' : 'AI_STREAM_INTERRUPTED');
	}
	if (!finished) throw new AIStreamError(signal.aborted ? 'AI_STREAM_CANCELLED' : 'AI_STREAM_INTERRUPTED');
}

async function pumpSseEvents(stream: ReadableStream<Uint8Array>, onEvent: (event: AiEvent) => void): Promise<boolean> {
  const reader = stream.getReader();
  const decoder = new TextDecoder();
  let buffer = '';
  let finished = false;
  try {
    while (true) {
      const result = await reader.read();
      if (result.done) break;
      const next = consumeSseChunk(buffer + decoder.decode(result.value, { stream: true }), onEvent);
      buffer = next.buffer;
      finished ||= next.finished;
    }
  } finally {
    reader.releaseLock();
  }
  return finished;
}

function consumeSseChunk(raw: string, onEvent: (event: AiEvent) => void): { buffer: string; finished: boolean } {
	let buffer = raw.replaceAll('\r\n', '\n');
	let finished = false;
	while (buffer.includes('\n\n')) {
		const end = buffer.indexOf('\n\n');
		const frame = buffer.slice(0, end);
		if (frame.length > frameLimit) throw new AIStreamError('AI_STREAM_FRAME_TOO_LARGE');
		buffer = buffer.slice(end + 2);
		const event = parseSseFrame(frame);
		if (event) onEvent(event);
		if (event?.type === 'run_finished') finished = true;
	}
	if (buffer.length > frameLimit) throw new AIStreamError('AI_STREAM_FRAME_TOO_LARGE');
	return { buffer, finished };
}

function parseSseFrame(frame: string): AiEvent | null {
  const { name, data } = parseFrameFields(frame);
  if (!name && data.length === 0) return null;
  if (!eventNames.has(name) || data.length === 0) throw new AIStreamError('AI_STREAM_PROTOCOL');
  let payload: unknown;
  try { payload = JSON.parse(data.join('\n')); } catch { throw new AIStreamError('AI_STREAM_PROTOCOL'); }
  if (!payload || typeof payload !== 'object' || Array.isArray(payload)) throw new AIStreamError('AI_STREAM_PROTOCOL');
  if (!validEventPayload(name, payload as Record<string, unknown>)) throw new AIStreamError('AI_STREAM_PROTOCOL');
  return { type: name, ...payload } as AiEvent;
}

function validEventPayload(name: string, value: Record<string, unknown>): boolean {
	if ('type' in value) return false;
	return Boolean(eventValidators[name]?.(value));
}

const eventValidators: Record<string, (value: Record<string, unknown>) => boolean> = {
	run_started: (v) => isString(v.runId) && (v.phase === 'PREVIEW' || v.phase === 'RUN'),
	text_delta: (v) => isString(v.text),
	tool_started: (v) => isString(v.toolId) && isString(v.title) && isSequence(v.sequence),
	tool_finished: (v) => isString(v.toolId) && isString(v.title) && isSequence(v.sequence) && isString(v.status) && optionalString(v.createdResourceId),
	preview_ready: validPreviewPayload,
	secret: (v) => isString(v.toolId) && isString(v.targetAccountId) && isString(v.value),
	run_finished: (v) => v.status === 'SUCCESS' || v.status === 'FAILED' || v.status === 'CANCELLED',
	error: (v) => isString(v.code),
};

function validPreviewPayload(value: Record<string, unknown>): boolean {
	return isString(value.summary) && typeof value.isHighRisk === 'boolean' && stringArray(value.requiredInputs) &&
		Array.isArray(value.steps) && value.steps.every(validStep) && optionalString(value.previewToken);
}

function validStep(value: unknown): boolean {
	if (!value || typeof value !== 'object' || Array.isArray(value)) return false;
	const step = value as Record<string, unknown>;
	return isString(step.toolId) && isString(step.title) && isSequence(step.sequence) && isSequence(step.maxCalls) && isString(step.risk) &&
		validStepDetails(step);
}

function validStepDetails(step: Record<string, unknown>): boolean {
	return objectValue(step.arguments) && optionalStringArray(step.targetIds) && optionalStringArray(step.scopeIds);
}

function isString(value: unknown): value is string { return typeof value === 'string'; }
function isSequence(value: unknown): value is number { return Number.isSafeInteger(value) && (value as number) >= 1; }
function optionalString(value: unknown): boolean { return value === undefined || isString(value); }
function stringArray(value: unknown): boolean { return Array.isArray(value) && value.every(isString); }
function optionalStringArray(value: unknown): boolean { return value === undefined || stringArray(value); }
function objectValue(value: unknown): boolean { return Boolean(value && typeof value === 'object' && !Array.isArray(value)); }

function parseFrameFields(frame: string): { name: string; data: string[] } {
	let name = '';
	const data: string[] = [];
	for (const line of frame.split('\n')) {
		if (line.startsWith('event:')) name = line.slice(6).trim();
		if (line.startsWith('data:')) data.push(line.slice(5).trimStart());
	}
	return { name, data };
}
