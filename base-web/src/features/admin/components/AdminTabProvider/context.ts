import { createContext } from 'react';

/** 页面读取自身标签路由上下文的唯一入口的数据形状。
 *  keep-alive 渲染栈中所有已打开标签同时挂载，react-router 的 useLocation /
 *  useSearchParams / useMatches 读的是 Router 顶层的全局 LocationContext，
 *  隐藏标签会拿到活动标签的值；useParams 虽按标签隔离，但依赖 renderMatches
 *  建立 RouteContext 这一实现细节。故统一由本上下文注入。见 rules §4.8。 */
export interface AdminTabContext {
  /** 该标签的路由路径。 */
  path: string;
  /** 该标签自身的查询串，与全局 URL 无关。 */
  search: string;
  /** 该标签自身的路由参数，取自其 matches 末项。 */
  params: Readonly<Record<string, string | undefined>>;
  /** 该标签当前是否可见。隐藏标签仍在挂载并运行副作用，
   *  定时器、window 监听与订阅必须据此自行判断是否执行。 */
  isActive: boolean;
}

export const adminTabContext = createContext<AdminTabContext | null>(null);
