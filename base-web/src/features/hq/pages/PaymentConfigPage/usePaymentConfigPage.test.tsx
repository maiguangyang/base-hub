// @vitest-environment jsdom

import { act, cleanup, renderHook, waitFor } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import type { ScopeRef, ScopeView } from '@/features/hq/api/paymentConfig';
import { usePaymentConfigPage } from './usePaymentConfigPage';

const mocks = vi.hoisted(() => ({
  read: vi.fn<(ref: ScopeRef) => Promise<ScopeView>>(),
}));

vi.mock('@/features/hq/api/paymentConfig', async (importOriginal) => ({
  ...await importOriginal<typeof import('@/features/hq/api/paymentConfig')>(),
  paymentConfigAPI: { read: mocks.read },
}));
vi.mock('@apollo/client/react', () => ({
  useQuery: () => ({ data: undefined, loading: false, error: undefined, fetchMore: vi.fn() }),
}));
vi.mock('@/features/admin/lib/useCompleteFilterCatalog', () => ({
  useCompleteFilterCatalog: () => ({ items: [], error: undefined }),
}));

afterEach(() => { cleanup(); vi.clearAllMocks(); });

it('按固定目标读取配置，并忽略切换目标前的旧响应', async () => {
  let resolveOldReload!: (value: ScopeView) => void;
  const oldReload = new Promise<ScopeView>((resolve) => { resolveOldReload = resolve; });
  const globalView: ScopeView = { channels: [] };
  const franchiseView: ScopeView = { channels: [] };
  const globalRef: ScopeRef = { scope: 'GLOBAL' };
  const franchiseRef: ScopeRef = { scope: 'FRANCHISE', organizationId: 'org' };
  let globalReads = 0;
  mocks.read.mockImplementation((ref) => {
    if (ref.scope === 'FRANCHISE') return Promise.resolve(franchiseView);
    globalReads += 1;
    return globalReads === 1 ? Promise.resolve(globalView) : oldReload;
  });

  const { result, rerender } = renderHook(({ ref }) => usePaymentConfigPage(ref), { initialProps: { ref: globalRef } });
  await waitFor(() => expect(result.current.view).toBe(globalView));
  expect(mocks.read).toHaveBeenCalledWith(globalRef);
  let reload!: Promise<void>;
  act(() => { reload = result.current.reload(); });
  rerender({ ref: franchiseRef });
  await waitFor(() => expect(result.current.view).toBe(franchiseView));
  expect(mocks.read).toHaveBeenCalledWith(franchiseRef);
  await act(async () => { resolveOldReload(globalView); await reload; });
  expect(result.current.view).toBe(franchiseView);
  expect(result.current.loading).toBe(false);
});
