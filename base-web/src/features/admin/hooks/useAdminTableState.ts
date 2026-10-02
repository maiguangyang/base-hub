import { useMemo, useState } from 'react';

/** 列表排序方向。null 表示未排序，回到数据原始顺序。 */
export type SortDirection = 'asc' | 'desc' | null;

/** 列表页共用的筛选、排序、分页与选择状态。
 *  用户管理与任务管理两页共用，故提升为共享 hook。 */
export interface AdminTableState<T> {
  keyword: string;
  setKeyword: (next: string) => void;
  sortKey: keyof T | null;
  sortDirection: SortDirection;
  /** 点击表头：无序 → 升序 → 降序 → 无序。 */
  toggleSort: (key: keyof T) => void;
  page: number;
  setPage: (next: number) => void;
  pageSize: number;
  setPageSize: (next: number) => void;
  /** 已选中行的 id 集合。批量操作以此为准。 */
  selected: Set<string>;
  toggleSelect: (id: string) => void;
  /** 全选/取消全选当前页。 */
  toggleSelectAll: (ids: string[]) => void;
  clearSelection: () => void;
  /** 经过筛选与排序后的完整结果，用于计算总数。 */
  filtered: T[];
  /** 当前页要渲染的那一段。 */
  pageRows: T[];
}

/** 按当前排序键与方向排序；未排序时保持原始顺序。 */
function sortRows<T>(rows: T[], key: keyof T | null, direction: SortDirection): T[] {
  if (!key || !direction) return rows;
  return [...rows].sort((a, b) => {
    const left = a[key];
    const right = b[key];
    if (left === right) return 0;
    const result = left > right ? 1 : -1;
    return direction === 'asc' ? result : -result;
  });
}

export function useAdminTableState<T extends { id: string }>(
  rows: T[],
  matches: (row: T, keyword: string) => boolean,
): AdminTableState<T> {
  const [keyword, setKeywordRaw] = useState('');
  const [sortKey, setSortKey] = useState<keyof T | null>(null);
  const [sortDirection, setSortDirection] = useState<SortDirection>(null);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSizeRaw] = useState(10);
  const [selected, setSelected] = useState<Set<string>>(new Set());

  const filtered = useMemo(() => {
    const trimmed = keyword.trim();
    const visible = trimmed === '' ? rows : rows.filter((row) => matches(row, trimmed));
    return sortRows(visible, sortKey, sortDirection);
  }, [rows, keyword, sortKey, sortDirection, matches]);

  const pageRows = useMemo(
    () => filtered.slice((page - 1) * pageSize, page * pageSize),
    [filtered, page, pageSize],
  );

  /** 筛选或页长变化后停留在旧页码会显示空白，统一回到第一页。 */
  const setKeyword = (next: string) => { setKeywordRaw(next); setPage(1); };
  const setPageSize = (next: number) => { setPageSizeRaw(next); setPage(1); };

  const toggleSort = (key: keyof T) => {
    if (sortKey !== key) { setSortKey(key); setSortDirection('asc'); return; }
    if (sortDirection === 'asc') { setSortDirection('desc'); return; }
    setSortKey(null);
    setSortDirection(null);
  };

  const toggleSelect = (id: string) => setSelected((prev) => {
    const next = new Set(prev);
    if (!next.delete(id)) next.add(id);
    return next;
  });

  const toggleSelectAll = (ids: string[]) => setSelected((prev) => (
    ids.every((id) => prev.has(id)) ? new Set() : new Set(ids)
  ));

  return {
    keyword, setKeyword, sortKey, sortDirection, toggleSort,
    page, setPage, pageSize, setPageSize,
    selected, toggleSelect, toggleSelectAll, clearSelection: () => setSelected(new Set()),
    filtered, pageRows,
  };
}
