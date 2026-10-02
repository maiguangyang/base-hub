import { useAdminTabsStore } from '@/stores/slices/adminTabsSlice';
import type { AdminTab } from '@/stores/slices/tabOperations';

/** 标签列表。原子 selector，避免订阅整个 store 造成无关重渲染。 */
export function useAdminTabList(): AdminTab[] {
  return useAdminTabsStore((state) => state.tabs);
}

/** 当前活动标签路径。 */
export function useActiveTabPath(): string {
  return useAdminTabsStore((state) => state.activePath);
}

/** 错误边界已触发的标签路径。标签栏据此给出可见标识——
 *  隐藏标签的 fallback 同样不可见，不标识用户在切过去之前完全无感知。 */
export function useErroredTabPaths(): string[] {
  return useAdminTabsStore((state) => state.erroredPaths);
}
