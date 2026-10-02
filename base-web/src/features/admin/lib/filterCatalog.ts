/** 关联筛选目录只有完整加载后才可用，避免把截断列表伪装成全部选项。 */
export function filterCatalogReady(
  loading: boolean,
  error: unknown,
  loadedCount: number,
  total: number | undefined,
): boolean {
  if (loading || error || total === undefined) return false;
  return loadedCount === total;
}

export interface CatalogPage<T extends { id: string }> {
  items: T[];
  total: number;
}

export function catalogPage<T extends { id: string }>(result: { data: T[]; total: number } | null | undefined): CatalogPage<T> | undefined {
  return result ? { items: result.data, total: result.total } : undefined;
}

/** 只在全部分页且 ID 无重复时交付可选目录。 */
export async function completeCatalogPages<T extends { id: string }>(
  first: CatalogPage<T>,
  fetchPage: (page: number) => Promise<CatalogPage<T>>,
): Promise<T[]> {
  const items = [...first.items];
  const ids = new Set(items.map((item) => item.id));
  if (ids.size !== items.length || items.length > first.total) throw new Error('catalog incomplete');
  for (let page = 2; items.length < first.total; page += 1) {
    const next = await fetchPage(page);
    if (next.total !== first.total || next.items.length === 0) throw new Error('catalog incomplete');
    const before = ids.size;
    for (const item of next.items) {
      if (!ids.has(item.id)) { ids.add(item.id); items.push(item); }
    }
    if (ids.size === before || items.length > first.total) throw new Error('catalog incomplete');
  }
  return items;
}
