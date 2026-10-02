import { useMemo, type ReactNode } from 'react';
import type { AdminTab } from '@/stores/slices/tabOperations';
import { adminTabContext, type AdminTabContext } from './context';

/** 渲染栈向每个标签注入其自身上下文。params 与 isActive 页面无法自行推导，
 *  必须由栈传入；页面不得自行调用 matchRoutes 重算，否则产生第二个真值源。
 *
 *  params 由调用方以稳定引用传入，不接收 matches：matchRoutes 每次调用都返回
 *  新数组，若把它当 memo 依赖，context value 每次渲染都是新对象，
 *  所有已挂载标签的子树都会跟着重渲染，keep-alive 的隔离收益被抵消。 */
export interface AdminTabProviderProps {
  tab: AdminTab;
  params: Readonly<Record<string, string | undefined>>;
  isActive: boolean;
  children: ReactNode;
}

export function AdminTabProvider({ tab, params, isActive, children }: AdminTabProviderProps) {
  const value = useMemo<AdminTabContext>(
    () => ({ path: tab.path, search: tab.search, params, isActive }),
    [tab.path, tab.search, params, isActive],
  );

  return <adminTabContext.Provider value={value}>{children}</adminTabContext.Provider>;
}
