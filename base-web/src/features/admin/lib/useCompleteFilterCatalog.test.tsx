// @vitest-environment jsdom

import { act, renderHook, waitFor } from '@testing-library/react';
import { StrictMode } from 'react';
import { describe, expect, it, vi } from 'vitest';
import { useCompleteFilterCatalog } from './useCompleteFilterCatalog';

describe('useCompleteFilterCatalog', () => {
  it('首批已是完整目录时立即可用且不追加请求', () => {
    const first = { items: [{ id: 'one' }], total: 1 };
    const fetchPage = vi.fn(async () => first);
    const { result } = renderHook(() => useCompleteFilterCatalog('hq', first, false, undefined, fetchPage));

    expect(result.current.ready).toBe(true);
    expect(result.current.items).toEqual(first.items);
    expect(fetchPage).not.toHaveBeenCalled();
  });

  it('第二页加载完成前不暴露部分目录，完成后提供全部选项', async () => {
    const first = { items: Array.from({ length: 200 }, (_, index) => ({ id: String(index + 1) })), total: 201 };
    const fetchPage = vi.fn(async (page: number) => {
      expect(page).toBe(2);
      return { items: [{ id: '201' }], total: 201 };
    });

    const { result } = renderHook(() => useCompleteFilterCatalog('hq', first, false, undefined, fetchPage));
    expect(result.current.ready).toBe(false);
    expect(result.current.items).toEqual([]);

    await waitFor(() => expect(result.current.ready).toBe(true));
    expect(result.current.items).toHaveLength(201);
    expect(fetchPage).toHaveBeenCalledTimes(1);
  });

  it('后续页失败时停止加载并保持筛选禁用', async () => {
    const first = { items: [{ id: 'one' }], total: 2 };
    const fetchPage = vi.fn(async () => { throw new Error('offline'); });
    const { result } = renderHook(() => useCompleteFilterCatalog('hq', first, false, undefined, fetchPage));

    await waitFor(() => expect(result.current.error).toBeTruthy());
    expect(result.current.loading).toBe(false);
    expect(result.current.ready).toBe(false);
    expect(result.current.items).toEqual([]);
  });
});

describe('目录重试与严格模式', () => {
  it('同一目录首批数据刷新后重新汇总后续页', async () => {
    const fetchPage = vi.fn()
      .mockResolvedValueOnce({ items: [{ id: 'old-tail', name: '旧单位' }], total: 2 })
      .mockResolvedValueOnce({ items: [{ id: 'new-tail', name: '新单位' }], total: 2 });
    const oldFirst = { items: [{ id: 'base', name: '原单位' }], total: 2 };
    const newFirst = { items: [{ id: 'base', name: '新名称' }], total: 2 };
    const { result, rerender } = renderHook(({ first }) => useCompleteFilterCatalog('packages', first, false, undefined, fetchPage), {
      initialProps: { first: oldFirst },
    });
    await waitFor(() => expect(result.current.items.map((item) => item.name)).toEqual(['原单位', '旧单位']));

    rerender({ first: newFirst });
    await waitFor(() => expect(result.current.items.map((item) => item.name)).toEqual(['新名称', '新单位']));
    expect(fetchPage).toHaveBeenCalledTimes(2);
  });

  it('后续页失败后可主动重试并恢复完整目录', async () => {
    const first = { items: [{ id: 'one' }], total: 2 };
    const fetchPage = vi.fn()
      .mockRejectedValueOnce(new Error('offline'))
      .mockResolvedValueOnce({ items: [{ id: 'two' }], total: 2 });
    const { result } = renderHook(() => useCompleteFilterCatalog('hq', first, false, undefined, fetchPage));

    await waitFor(() => expect(result.current.error).toBeTruthy());
    act(() => result.current.retry());
    await waitFor(() => expect(result.current.ready).toBe(true));
    expect(result.current.error).toBeUndefined();
    expect(result.current.items).toEqual([{ id: 'one' }, { id: 'two' }]);
    expect(fetchPage).toHaveBeenCalledTimes(2);
  });

  it('严格模式重建 effect 后仍会完成目录加载', async () => {
    const first = { items: [{ id: 'one' }], total: 2 };
    const fetchPage = vi.fn(async () => ({ items: [{ id: 'two' }], total: 2 }));
    const { result } = renderHook(() => useCompleteFilterCatalog('hq', first, false, undefined, fetchPage), { wrapper: StrictMode });

    await waitFor(() => expect(result.current.ready).toBe(true));
    expect(result.current.items).toHaveLength(2);
  });
});

describe('目录组织隔离', () => {
  it('切换组织后不沿用上一个组织的目录', async () => {
    const first = { items: [{ id: 'one' }], total: 2 };
    const fetchPage = vi.fn(async () => ({ items: [{ id: 'two' }], total: 2 }));
    const { result, rerender } = renderHook(({ scopeKey }) => useCompleteFilterCatalog(scopeKey, first, false, undefined, fetchPage), {
      initialProps: { scopeKey: 'franchise-a' },
    });
    await waitFor(() => expect(result.current.ready).toBe(true));

    rerender({ scopeKey: 'franchise-b' });
    expect(result.current.ready).toBe(false);
    expect(result.current.items).toEqual([]);
    await waitFor(() => expect(result.current.ready).toBe(true));
    expect(fetchPage).toHaveBeenCalledTimes(2);
  });

  it('切到完整目录期间，旧组织未完成的请求不能回填目录', async () => {
    let resolveOld: ((page: { items: { id: string }[]; total: number }) => void) | undefined;
    const fetchPage = vi.fn(() => new Promise<{ items: { id: string }[]; total: number }>((resolve) => { resolveOld = resolve; }));
    const first = { items: [{ id: 'one' }], total: 2 };
    const complete = { items: [{ id: 'other' }], total: 1 };
    const { result, rerender } = renderHook(({ scopeKey }) => useCompleteFilterCatalog(scopeKey, scopeKey === 'a' ? first : complete, false, undefined, fetchPage), {
      initialProps: { scopeKey: 'a' },
    });

    rerender({ scopeKey: 'b' });
    expect(result.current.items).toEqual(complete.items);
    await act(async () => { resolveOld?.({ items: [{ id: 'two' }], total: 2 }); });
    rerender({ scopeKey: 'a' });
    expect(result.current.ready).toBe(false);
    expect(result.current.items).toEqual([]);
  });
});
