import { useContext } from 'react';
import { adminTabContext, type AdminTabContext } from '../components/AdminTabProvider/context';

/** 读取当前标签自身的路由上下文。后台页面读取路由信息的唯一合法入口。
 *  直接使用 react-router 的 useParams / useLocation / useSearchParams /
 *  useMatch / useMatches 已由 ESLint no-restricted-imports 禁止。 */
export function useAdminTab(): AdminTabContext {
  const value = useContext(adminTabContext);
  if (!value) throw new Error('useAdminTab 必须在 AdminTabProvider 内使用');
  return value;
}
