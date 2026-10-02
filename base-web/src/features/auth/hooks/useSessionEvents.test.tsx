import { describe, expect, it, vi } from 'vitest';
import type { ApolloClient } from '@apollo/client';
import type { SessionEventsSubscription } from '@/__generated__/graphql';
import { subscribeSessionEvents } from './useSessionEvents';

function fakeSessionEvents() {
  let observer: { next(value: { data?: SessionEventsSubscription }): void } | undefined;
  const unsubscribe = vi.fn();
  const client = {
    subscribe: () => ({ subscribe: (next: typeof observer) => { observer = next; return { unsubscribe }; } }),
  } as unknown as ApolloClient;
  return {
    client, unsubscribe,
    publish: (event: SessionEventsSubscription['sessionEvents']) => observer?.next({ data: { sessionEvents: event } }),
  };
}

describe('会话事件订阅', () => {
  it('终态事件仅重置当前浏览器会话并跳转登录，重复事件保持幂等', async () => {
    const source = fakeSessionEvents();
    const resetSession = vi.fn(async () => undefined);
    const navigate = vi.fn();
    let destination = '/admin/hq';
    subscribeSessionEvents(source.client, { currentOrganizationId: 'org-a', returnDestination: () => destination, resetSession, navigate });
    destination = '/admin/hq/audit?page=2#latest';
    const event = { code: 'SESSION_REVOKED' as const, sessionId: 'session-a', organizationId: 'org-a', occurredAt: 'now' };
    source.publish(event); source.publish({ ...event, code: 'CREDENTIALS_CHANGED' });
    await vi.waitFor(() => expect(resetSession).toHaveBeenCalledTimes(1));
    expect(navigate).toHaveBeenCalledWith('/admin/login?dt=%2Fadmin%2Fhq%2Faudit%3Fpage%3D2%23latest');
  });

  it('卸载或切换工作台时取消原订阅', () => {
    const source = fakeSessionEvents();
    const stop = subscribeSessionEvents(source.client, { currentOrganizationId: 'org-a', returnDestination: () => '/admin/hq', resetSession: vi.fn(), navigate: vi.fn() });
    stop();
    expect(source.unsubscribe).toHaveBeenCalledTimes(1);
  });
});
