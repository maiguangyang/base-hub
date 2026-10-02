import { ChevronLeft, ChevronRight } from 'lucide-react';
import { useActiveTabPath, useAdminTabList, useErroredTabPaths } from '@/features/admin/hooks/useAdminTabs';
import { TabContextMenu } from './TabContextMenu';
import { TabItem } from './TabItem';
import { useTabActions } from './hooks/useTabActions';
import { useTabOverflow } from './hooks/useTabOverflow';

/** 多标签工作区的标签栏。只做编排，单个标签与右键菜单各自独立。
 *  全部标签操作走 useTabActions，由它保证地址栏跟随活动标签。 */
export function AdminTabBar() {
  const tabs = useAdminTabList();
  const activePath = useActiveTabPath();
  const erroredPaths = useErroredTabPaths();
  const actions = useTabActions();
  const { scrollerRef, canScrollLeft, canScrollRight, scrollByStep } = useTabOverflow(activePath);

  return (
    <div className="flex h-10 shrink-0 items-center border-b border-border bg-secondary">
      <button type="button" aria-label="向左滚动标签" disabled={!canScrollLeft}
        className="h-full px-1.5 text-muted-foreground disabled:opacity-30" onClick={() => scrollByStep(-1)}>
        <ChevronLeft className="size-4" />
      </button>
      {/* 保留横向滚动供箭头、触控板和活动标签定位使用，但隐藏会挤占标签高度的原生滚动条。
          overflow-y 仍须显式 hidden：overflow-x 为非 visible 时，另一轴的 visible 会计算成 auto。 */}
      <div
        ref={scrollerRef}
        className="flex h-full flex-1 overflow-x-auto overflow-y-hidden [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
      >
        {tabs.map((tab) => (
          <TabContextMenu key={tab.path} tab={tab}>
            {/* 包裹节点承接 ContextMenuTrigger 的 asChild，须与内层同为 shrink-0，
                否则标签多时外层被压缩而内层拒绝收缩，布局行为不确定。 */}
            <div className="shrink-0">
              <TabItem
                tab={tab}
                isActive={tab.path === activePath}
                errored={erroredPaths.includes(tab.path)}
                onActivate={actions.activate}
                onClose={actions.close}
              />
            </div>
          </TabContextMenu>
        ))}
      </div>
      <button type="button" aria-label="向右滚动标签" disabled={!canScrollRight}
        className="h-full px-1.5 text-muted-foreground disabled:opacity-30" onClick={() => scrollByStep(1)}>
        <ChevronRight className="size-4" />
      </button>
    </div>
  );
}
