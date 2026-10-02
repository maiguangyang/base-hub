import { describe, expect, it, vi } from 'vitest';
import type { ApolloClient } from '@apollo/client';
import type { SessionEventsSubscription } from '@/__generated__/graphql';
import { subscribeSessionEvents } from '@/features/auth/hooks/useSessionEvents';

function subscriptionProbe() {
  let next: ((value: { data?: SessionEventsSubscription }) => void) | undefined;
  const client = { subscribe: () => ({ subscribe: (observer: { next(value: { data?: SessionEventsSubscription }): void }) => { next = observer.next; return { unsubscribe: vi.fn() }; } }) } as unknown as ApolloClient;
  return { client, publish: (organizationId: string) => next?.({ data: { sessionEvents: { code: 'ORGANIZATION_SUSPENDED', sessionId: 'session-a', organizationId, occurredAt: 'now' } } }) };
}

describe('组织暂停的工作台隔离', () => {
  it('只处理当前工作台对应组织的强退事件', async () => {
    const source = subscriptionProbe();
    const resetSession = vi.fn(async () => undefined);
    const navigate = vi.fn();
    subscribeSessionEvents(source.client, { currentOrganizationId: 'org-a', returnDestination: () => '/admin/franchise/stores?q=one', resetSession, navigate });
    source.publish('org-b');
    expect(resetSession).not.toHaveBeenCalled();
    source.publish('org-a');
    await vi.waitFor(() => expect(resetSession).toHaveBeenCalledWith('ORGANIZATION_SUSPENDED'));
    expect(navigate).toHaveBeenCalledWith('/admin/login?dt=%2Fadmin%2Ffranchise%2Fstores%3Fq%3Done');
  });
});
