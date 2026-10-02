import { useEffect } from 'react';
import type { ApolloClient } from '@apollo/client';
import type { ViewerSummary } from '@/features/auth/store/authStore';
import { SESSION_EVENTS_SUBSCRIPTION } from '@/features/auth/graphql/auth';
import type { GraphQLRuntime } from '@/lib/graphql/client';
import type { SessionEventCode, SessionEventsSubscription } from '@/__generated__/graphql';
import { currentDestination, loginPath } from '@/features/auth/pages/authFlow';

interface SessionEventActions {
  currentOrganizationId: string | null;
  returnDestination(): string | undefined;
  resetSession(code: string): Promise<void>;
  navigate(path: string): void;
}

const terminalEventCodes: readonly SessionEventCode[] = [
  'SESSION_REVOKED', 'CREDENTIALS_CHANGED', 'ORGANIZATION_SUSPENDED',
];

/** 订阅当前权威 Session 的终态事件，返回显式取消函数供工作台切换使用。 */
export function subscribeSessionEvents(client: ApolloClient, actions: SessionEventActions): () => void {
  let handled = false;
  const subscription = client.subscribe({ query: SESSION_EVENTS_SUBSCRIPTION }).subscribe({
    next: (result) => {
      const event = result.data?.sessionEvents;
      if (handled || !event || !shouldHandleEvent(event, actions.currentOrganizationId)) return;
      handled = true;
      void actions.resetSession(event.code).finally(() => actions.navigate(loginPath(actions.returnDestination())));
    },
  });
  return () => subscription.unsubscribe();
}

/** 在已认证后台根节点挂载一次；依赖变化会先取消旧 Session 订阅。 */
export function useSessionEvents(runtime: GraphQLRuntime, viewer: ViewerSummary | null): void {
  const organizationId = viewer?.currentWorkspace?.organizationId ?? null;
  const sessionReady = Boolean(viewer?.currentWorkspace);
  useEffect(() => {
    if (!sessionReady) return undefined;
    return subscribeSessionEvents(runtime.client, {
      currentOrganizationId: organizationId,
      returnDestination: () => currentDestination(window.location),
      resetSession: runtime.resetSession,
      navigate: (path) => window.location.replace(path),
    });
  }, [organizationId, runtime, sessionReady]);
}

type SessionEventValue = SessionEventsSubscription['sessionEvents'];

function shouldHandleEvent(event: Pick<SessionEventValue, 'code' | 'organizationId'>, currentOrganizationId: string | null): boolean {
  if (!terminalEventCodes.includes(event.code)) return false;
  if (event.code !== 'ORGANIZATION_SUSPENDED') return true;
  return event.organizationId !== null && event.organizationId === currentOrganizationId;
}
