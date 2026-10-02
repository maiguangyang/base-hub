import { useState } from 'react';

export function useListPagination(initialPageSize = 20) {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(initialPageSize);
  const changePageSize = (size: number) => { setPageSize(size); setPage(1); };
  return { page, setPage, pageSize, changePageSize };
}
