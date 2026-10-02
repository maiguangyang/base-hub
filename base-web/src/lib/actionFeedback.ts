export type ActionFeedback = { kind: 'success' | 'error'; message: string };
type Listener = (feedback: ActionFeedback) => void;
const listeners = new Set<Listener>();

/** Transport feedback stays independent of React and survives page navigation. */
export function publishActionFeedback(kind: ActionFeedback['kind'], message: string): void {
  for (const listener of listeners) listener({ kind, message });
}

export function subscribeActionFeedback(listener: Listener): () => void {
  listeners.add(listener);
  return () => { listeners.delete(listener); };
}
