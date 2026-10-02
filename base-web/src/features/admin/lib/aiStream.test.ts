import { afterEach, describe, expect, it, vi } from 'vitest';
import { AIStreamError, readAiStream, type AiEvent } from './aiStream';

afterEach(() => vi.unstubAllGlobals());

function streamResponse(chunks: Uint8Array[], status = 200): Response {
  const body = new ReadableStream<Uint8Array>({
    start(controller) { for (const chunk of chunks) controller.enqueue(chunk); controller.close(); },
  });
  return new Response(body, { status, headers: { 'Content-Type': 'text/event-stream' } });
}

describe('AI SSE 流', () => {
  it('跨 UTF-8 边界解析事件且只发送同源 Cookie 请求', async () => {
    const bytes = new TextEncoder().encode('event: text_delta\ndata: {"text":"你好"}\n\n: keepalive\n\nevent: run_finished\ndata: {"status":"SUCCESS"}\n\n');
    const split = bytes.indexOf(0xe5) + 1;
    const fetchMock = vi.fn().mockResolvedValue(streamResponse([bytes.slice(0, split), bytes.slice(split)]));
    vi.stubGlobal('fetch', fetchMock);
    const events: AiEvent[] = [];
    await readAiStream('preview', { prompt: '列出门店' }, new AbortController().signal, (event) => events.push(event));
    expect(events.map((event) => event.type)).toEqual(['text_delta', 'run_finished']);
    expect(events[0]).toEqual({ type: 'text_delta', text: '你好' });
    expect(fetchMock.mock.calls[0][1]).toMatchObject({ method: 'POST', credentials: 'include', cache: 'no-store' });
    expect(fetchMock.mock.calls[0][0]).toContain('/api/ai/preview');
  });

  it('连接中断时拒绝把执行视为成功', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(streamResponse([new TextEncoder().encode('event: run_started\ndata: {"runId":"r","phase":"RUN"}\n\n')])));
    await expect(readAiStream('run', { previewToken: 'token' }, new AbortController().signal, () => undefined))
      .rejects.toEqual(new AIStreamError('AI_STREAM_INTERRUPTED'));
  });

  it('只暴露服务端结构化错误码', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('{"code":"TOKEN_EXPIRED_OR_USED"}', { status: 400, headers: { 'Content-Type': 'application/json' } })));
    await expect(readAiStream('run', { previewToken: 'old' }, new AbortController().signal, () => undefined))
      .rejects.toEqual(new AIStreamError('TOKEN_EXPIRED_OR_USED'));
  });

  it('拒绝缺少字段或伪造事件类型的流帧', async () => {
	for (const frame of [
		'event: preview_ready\ndata: {"summary":"ok","isHighRisk":false}\n\nevent: run_finished\ndata: {"status":"SUCCESS"}\n\n',
		'event: text_delta\ndata: {"type":"run_finished","status":"SUCCESS","text":"oops"}\n\nevent: run_finished\ndata: {"status":"SUCCESS"}\n\n',
		'event: run_finished\ndata: {"status":["SUCCESS"]}\n\n',
	]) {
		vi.stubGlobal('fetch', vi.fn().mockResolvedValue(streamResponse([new TextEncoder().encode(frame)])));
		await expect(readAiStream('preview', { prompt: '查询' }, new AbortController().signal, () => undefined))
			.rejects.toEqual(new AIStreamError('AI_STREAM_PROTOCOL'));
	}
  });
});
