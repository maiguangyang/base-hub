import { useMemo } from 'react';
import { useNavigate } from 'react-router';
import { useAdminTabsStore } from '@/stores/slices/adminTabsSlice';
import type { AdminTab } from '@/stores/slices/tabOperations';

/** 标签操作的统一入口。 */
export interface TabActions {
  /** 激活标签。地址栏是活动标签的唯一入口，激活一律走路由跳转。 */
  activate: (path: string) => void;
  close: (path: string) => void;
  closeOthers: (path: string) => void;
  closeRight: (path: string) => void;
  closeAll: () => void;
  /** 刷新不改变活动标签，因此无需同步地址栏。 */
  refresh: (path: string) => void;
}

/** 恢复标签自己的查询上下文；未知路径仍可安全回退到裸路径。 */
export function destinationForTab(path: string, tabs: readonly AdminTab[]): string {
  const search = tabs.find((tab) => tab.path === path)?.search ?? '';
  return `${path}${search}`;
}

function currentBrowserDestination(): string {
  return `${window.location.pathname}${window.location.search}`;
}

/**
 * 标签操作的唯一出口，把「改 store」与「同步地址栏」绑成一个不可分割的动作。
 *
 * 任何改变 activePath 的操作都必须同步地址栏。两者分叉时：
 * useSyncTabWithLocation 的 effect 依赖 [pathname, search, isKnownPath]，URL 没变就不重跑，
 * 于是面包屑跟着 activePath 走、地址栏停在旧路径；用户一刷新，
 * useSyncTabWithLocation 会按 URL 重新 openTab，刚关掉的标签自己回来。
 *
 * 因此调用方一律用本 hook，不要直接调用 store 的关闭类动作。
 */
export function useTabActions(): TabActions {
  const navigate = useNavigate();

  return useMemo(() => {
    /** 执行 store 变更后，把地址栏对齐到新的活动标签。 */
    const syncAfter = (mutate: () => void) => {
      mutate();
      const { activePath, tabs } = useAdminTabsStore.getState();
      const destination = destinationForTab(activePath, tabs);
      // 读 window.location 而非订阅 useLocation：这里只在回调中做一次比较，
      // 订阅会让整个标签栏随每次路由变化重渲染。
      if (currentBrowserDestination() !== destination) navigate(destination);
    };

    return {
      // 已经在目标标签上就不跳转：navigate 到相同路径仍会压入一条历史记录，
      // 反复点击当前标签会堆满历史，用户得连按多次后退才能真正离开。
      activate: (path) => {
        const destination = destinationForTab(path, useAdminTabsStore.getState().tabs);
        if (currentBrowserDestination() === destination) return;
        navigate(destination);
      },
      close: (path) => syncAfter(() => useAdminTabsStore.getState().closeTab(path)),
      closeOthers: (path) => syncAfter(() => useAdminTabsStore.getState().closeOthers(path)),
      closeRight: (path) => syncAfter(() => useAdminTabsStore.getState().closeRight(path)),
      closeAll: () => syncAfter(() => useAdminTabsStore.getState().closeAll()),
      refresh: (path) => useAdminTabsStore.getState().refreshTab(path),
    };
  }, [navigate]);
}
