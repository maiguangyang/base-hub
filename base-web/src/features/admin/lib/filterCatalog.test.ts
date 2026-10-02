import { describe, expect, it } from 'vitest';
import { completeCatalogPages, filterCatalogReady } from './filterCatalog';

describe('filterCatalogReady', () => {
  it('只在加载完成、无错误且条目总数一致时返回 true', () => {
    expect(filterCatalogReady(false, undefined, 2, 2)).toBe(true);
    expect(filterCatalogReady(true, undefined, 2, 2)).toBe(false);
    expect(filterCatalogReady(false, new Error('failed'), 2, 2)).toBe(false);
    expect(filterCatalogReady(false, undefined, 2, 3)).toBe(false);
    expect(filterCatalogReady(false, undefined, 0, undefined)).toBe(false);
  });
});

describe('completeCatalogPages', () => {
  it('按页取齐超过首批上限的目录', async () => {
    const first = Array.from({ length: 200 }, (_, index) => ({ id: String(index + 1) }));
    const fetchPage = async (page: number) => {
      expect(page).toBe(2);
      return { items: [{ id: '201' }], total: 201 };
    };

    const items = await completeCatalogPages({ items: first, total: 201 }, fetchPage);
    expect(items).toHaveLength(201);
    expect(items.at(-1)?.id).toBe('201');
  });

  it('页间目录发生变化时拒绝不完整选项', async () => {
    const fetchPage = async () => ({ items: [{ id: '1' }], total: 2 });
    await expect(completeCatalogPages({ items: [{ id: '1' }], total: 2 }, fetchPage)).rejects.toThrow('catalog incomplete');
  });
});
