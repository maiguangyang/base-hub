import { matchRoutes, renderMatches } from 'react-router';
import { AdminTabErrorBoundary } from '@/features/admin/components/AdminTabErrorBoundary';
import { AdminTabProvider } from '@/features/admin/components/AdminTabProvider';
import { useActiveTabPath, useAdminTabList } from '@/features/admin/hooks/useAdminTabs';
import { adminRoutes } from '@/features/admin/routes/adminRoutes';

/** 路径 -> 路由匹配结果与参数的缓存。
 *  matchRoutes 每次调用都返回新数组，在渲染期直接调用会让下游的 useMemo 全部失效，
 *  导致任何一次 store 变化都重渲染所有已挂载标签的子树，抵消 keep-alive 的隔离收益。
 *  adminRoutes 在运行期不变，因此按路径缓存一次即可获得稳定引用。 */
const resolvedByPath = new Map<string, { matches: ReturnType<typeof matchRoutes>; params: Readonly<Record<string, string | undefined>> }>();

/** 取得该路径的匹配结果，首次计算后复用同一引用。 */
function resolvePath(path: string) {
  const cached = resolvedByPath.get(path);
  if (cached) return cached;
  const matches = matchRoutes(adminRoutes, path);
  const entry = { matches, params: matches?.[matches.length - 1]?.params ?? {} };
  resolvedByPath.set(path, entry);
  return entry;
}

/**
 * 后台 keep-alive 渲染栈。
 *
 * ⚠️ 此处刻意不使用 <Outlet/>，不要改回去。
 *
 * 全部已打开标签同时挂载，非活动标签用 hidden（等价 display:none）隐藏而非卸载，
 * 以此保留滚动位置、筛选条件与表单草稿。改回 <Outlet/> 等于删除整个多标签工作区：
 * 编译通过、页面正常、只是每次切换标签都重新挂载，问题不会以报错形式出现。
 * 若确需改动，先读 docs/plans/2026-09-18-admin-console-shell-design.md 第三节的实测结论。
 *
 * 实测依据（Chrome 153 / WebKit / Firefox 155 一致）：
 * - display:none 恢复显示后浏览器会还原 scrollTop
 * - 隐藏期间读 scrollTop 与 offsetWidth 均返回 0
 * 因此滚动容器必须下沉到每个标签内部，且测量类逻辑只能在可见时执行。
 */
export function TabStack() {
  const tabs = useAdminTabList();
  const activePath = useActiveTabPath();

  return (
    <div className="flex-1 overflow-hidden">
      {tabs.map((tab) => {
        const entry = resolvePath(tab.path);
        const matches = entry.matches;
        // sanitizeSnapshot 已剔除失效标签，此处为双保险。
        if (!matches) return null;
        const isActive = tab.path === activePath;
        return (
          <div
            key={`${tab.path}#${tab.version}`}
            hidden={!isActive}
            className="h-full overflow-y-auto px-4 py-4 lg:px-6"
          >
            <AdminTabProvider tab={tab} params={entry.params} isActive={isActive}>
              <AdminTabErrorBoundary tabPath={tab.path}>
                {renderMatches(matches)}
              </AdminTabErrorBoundary>
            </AdminTabProvider>
          </div>
        );
      })}
    </div>
  );
}
