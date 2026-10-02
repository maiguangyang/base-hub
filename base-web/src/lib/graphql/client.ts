import {
  ApolloClient, ApolloLink, HttpLink, InMemoryCache, Observable,
} from '@apollo/client/core';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { getMainDefinition } from '@apollo/client/utilities';
import { createClient, type Client as GraphQLWsClient } from 'graphql-ws';
import { appConfig } from '@/config/app';
import { reportNonBlockingError } from '@/lib/diagnostics';
import { getGraphQLErrorCode, isTerminalSessionCode } from './errors';
import { createActionFeedbackLink } from './actionFeedbackLink';

interface SessionResetDependencies {
  disposeSubscription(): void;
  clearStore(): Promise<unknown>;
  resetWorkspace(): void;
  resetAuth(): void;
}

export interface GraphQLRuntime {
  client: ApolloClient;
  dispose(): void;
  resetSession(code: string): Promise<void>;
}

let singleton: GraphQLRuntime | undefined;

/** 合并并发终态信号，确保浏览器状态只清理一次。 */
export function createSessionResetController(dependencies: SessionResetDependencies) {
  let active: Promise<void> | undefined;
  return (code: string): Promise<void> => {
    if (code.length === 0) return Promise.resolve();
    if (active) return active;
    active = resetSession(dependencies).finally(() => { active = undefined; });
    return active;
  };
}

/** 返回浏览器进程内唯一的 Apollo runtime。 */
export function getGraphQLRuntime(state: Pick<SessionResetDependencies, 'resetAuth' | 'resetWorkspace'>): GraphQLRuntime {
  if (!singleton) singleton = createGraphQLRuntime(state);
  return singleton;
}

/** 销毁当前 Subscription 连接和 Apollo 缓存。 */
export async function disposeGraphQLRuntime(): Promise<void> {
  const runtime = singleton;
  singleton = undefined;
  if (!runtime) return;
  runtime.dispose();
  try {
    await runtime.client.clearStore();
  } catch (cause) {
    reportNonBlockingError('GraphQL runtime 缓存清理失败，runtime 已销毁。', cause);
  }
}

function createGraphQLRuntime(state: Pick<SessionResetDependencies, 'resetAuth' | 'resetWorkspace'>): GraphQLRuntime {
  const wsClient = createSubscriptionClient();
  const http = createAbortableFetch();
  const holder: { client?: ApolloClient } = {};
  const resetSessionState = createSessionResetController({
    disposeSubscription: () => { void wsClient.dispose(); },
    clearStore: () => holder.client?.clearStore() ?? Promise.resolve(),
    resetWorkspace: state.resetWorkspace,
    resetAuth: state.resetAuth,
  });
  const errorLink = createTerminalErrorLink(resetSessionState);
  const transport = ApolloLink.split(isSubscription, new GraphQLWsLink(wsClient), new HttpLink({ uri: appConfig.engine.graphQLHttpUrl, credentials: 'include', fetch: http.fetch }));
  const apolloClient = new ApolloClient({ cache: new InMemoryCache(), link: ApolloLink.from([createActionFeedbackLink(), errorLink, transport]) });
  holder.client = apolloClient;
  return { client: apolloClient, resetSession: resetSessionState, dispose: () => { http.abortAll(); void wsClient.dispose(); } };
}

function createAbortableFetch() {
  const controllers = new Set<AbortController>();
  const trackedFetch: typeof fetch = async (input, init) => {
    const controller = new AbortController();
    controllers.add(controller);
    const external = init?.signal;
    const abort = () => controller.abort();
    external?.addEventListener('abort', abort, { once: true });
    try {
      return await fetch(input, { ...init, signal: controller.signal });
    } finally {
      controllers.delete(controller);
      external?.removeEventListener('abort', abort);
    }
  };
  return {
    fetch: trackedFetch,
    abortAll: () => {
      for (const controller of controllers) controller.abort();
      controllers.clear();
    },
  };
}

function createSubscriptionClient(): GraphQLWsClient {
  return createClient({
    url: appConfig.engine.graphQLWebSocketUrl,
    lazy: true,
    retryAttempts: Number.POSITIVE_INFINITY,
    shouldRetry: () => typeof navigator === 'undefined' || navigator.onLine,
  });
}

function createTerminalErrorLink(reset: (code: string) => Promise<void>): ApolloLink {
  return new ApolloLink((operation, forward) => new Observable((observer) => {
    const subscription = forward(operation).subscribe({
      next: (result) => {
        handleResultErrors(result, reset);
        observer.next(result);
      },
      error: (error: unknown) => {
        handleError(error, reset);
        observer.error(error);
      },
      complete: () => observer.complete(),
    });
    return () => subscription.unsubscribe();
  }));
}

function handleResultErrors(result: ApolloLink.Result, reset: (code: string) => Promise<void>): void {
  for (const error of result.errors ?? []) handleError(error, reset);
}

function handleError(error: unknown, reset: (code: string) => Promise<void>): void {
  const code = getGraphQLErrorCode(error);
  if (code && isTerminalSessionCode(code)) void reset(code);
}

function isSubscription(operation: ApolloLink.Operation): boolean {
  const definition = getMainDefinition(operation.query);
  return definition.kind === 'OperationDefinition' && definition.operation === 'subscription';
}

async function resetSession(dependencies: SessionResetDependencies): Promise<void> {
  dependencies.disposeSubscription();
  try {
    await dependencies.clearStore();
  } catch (cause) {
    reportNonBlockingError('会话终态的 Apollo 缓存清理失败，继续重置本地状态。', cause);
  } finally {
    dependencies.resetWorkspace();
    dependencies.resetAuth();
  }
}
