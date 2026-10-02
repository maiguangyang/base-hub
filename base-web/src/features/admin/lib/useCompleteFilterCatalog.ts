import { useEffect, useRef, useState } from 'react';
import { completeCatalogPages, type CatalogPage } from './filterCatalog';

interface CatalogState<T> {
  key: string;
  items: T[];
  status: 'loading' | 'ready' | 'error';
  error?: unknown;
}

export function useCompleteFilterCatalog<T extends { id: string }>(
  scopeKey: string,
  first: CatalogPage<T> | undefined,
  loading: boolean,
  queryError: unknown,
  fetchPage: (page: number) => Promise<CatalogPage<T>>,
) {
  const [state, setState] = useState<CatalogState<T>>();
  const [retryCount, setRetryCount] = useState(0);
  const startedKey = useRef<string | undefined>(undefined);
  const generation = useRef(0);
  const catalogKey = JSON.stringify([scopeKey, first]);
  const firstReady = isFirstPageReady(first, loading, queryError);
  const retry = () => { startedKey.current = undefined; setRetryCount((count) => count + 1); };

  useEffect(() => () => { generation.current += 1; startedKey.current = undefined; }, [catalogKey]);
  useEffect(() => {
    if (loading || queryError || !first || firstReady || startedKey.current === catalogKey) return;
    startedKey.current = catalogKey;
    const currentGeneration = ++generation.current;
    setState({ key: catalogKey, items: [], status: 'loading' });
    void completeCatalogPages(first, fetchPage).then(
      (items) => { if (generation.current === currentGeneration) setState({ key: catalogKey, items, status: 'ready' }); },
      (error: unknown) => { if (generation.current === currentGeneration) setState({ key: catalogKey, items: [], status: 'error', error }); },
    );
  }, [catalogKey, first, firstReady, loading, queryError, fetchPage, retryCount]);

  return { ...catalogView(catalogKey, first, firstReady, loading, queryError, state), retry };
}

function isFirstPageReady<T extends { id: string }>(first: CatalogPage<T> | undefined, loading: boolean, error: unknown): boolean {
  if (loading || error || !first) return false;
  return first.items.length === first.total && new Set(first.items.map((item) => item.id)).size === first.total;
}

function catalogView<T extends { id: string }>(scopeKey: string, first: CatalogPage<T> | undefined, firstReady: boolean, loading: boolean, queryError: unknown, state?: CatalogState<T>) {
  const current = state?.key === scopeKey ? state : undefined;
  const ready = firstReady || current?.status === 'ready';
  return {
    items: catalogItems(first, firstReady, current),
    ready,
    loading: catalogLoading(first, ready, loading, queryError, current),
    error: firstReady ? undefined : catalogError(first, loading, queryError, current),
  };
}

function catalogItems<T extends { id: string }>(first: CatalogPage<T> | undefined, firstReady: boolean, current?: CatalogState<T>): T[] {
  if (firstReady) return first?.items ?? [];
  return current?.status === 'ready' ? current.items : [];
}

function catalogLoading<T extends { id: string }>(first: CatalogPage<T> | undefined, ready: boolean, loading: boolean, queryError: unknown, current?: CatalogState<T>): boolean {
  return loading || !ready && !queryError && first !== undefined && current?.status !== 'error';
}

function catalogError<T extends { id: string }>(first: CatalogPage<T> | undefined, loading: boolean, queryError: unknown, current?: CatalogState<T>): unknown {
  if (queryError) return queryError;
  if (current?.status === 'error') return current.error;
  return !loading && !first ? new Error('catalog unavailable') : undefined;
}
