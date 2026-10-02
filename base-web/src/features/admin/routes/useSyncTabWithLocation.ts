import { useEffect } from 'react';
import { useLocation } from 'react-router';
import { useAdminTabsStore } from '@/stores/slices/adminTabsSlice';

/** 把浏览器地址同步为标签集合中的活动标签。
 *  仅外壳可用——它依赖 useLocation，该 hook 在 keep-alive 栈中对页面不可信。
 *  返回当前路径是否已登记，未登记时由外壳渲染兜底页而非开标签。 */
export function useSyncTabWithLocation(): { isKnownPath: boolean } {
  const { pathname, search } = useLocation();
  const manifest = useAdminTabsStore((state) => state.manifest);
  const isKnownPath = manifest.items.some((item) => item.path === pathname);

  useEffect(() => {
    if (isKnownPath) useAdminTabsStore.getState().openTab(pathname, search);
  }, [pathname, search, isKnownPath]);

  return { isKnownPath };
}
