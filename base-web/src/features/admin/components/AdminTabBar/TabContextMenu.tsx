import type { ReactNode } from 'react';
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuTrigger } from '@/components/ui/context-menu';
import type { AdminTab } from '@/stores/slices/tabOperations';
import { useTabActions } from './hooks/useTabActions';

/** 标签右键菜单。固定标签的「关闭当前」禁用——store 侧也会拒绝，
 *  此处禁用只是让不可用状态在界面上可见。 */
export interface TabContextMenuProps {
  tab: AdminTab;
  children: ReactNode;
}

export function TabContextMenu({ tab, children }: TabContextMenuProps) {
  // 关闭类动作必须走 useTabActions，它保证地址栏跟随 activePath。
  // 直接调用 store 会让 URL 与活动标签分叉，刷新后关掉的标签复活。
  const actions = useTabActions();

  return (
    <ContextMenu>
      <ContextMenuTrigger asChild>{children}</ContextMenuTrigger>
      <ContextMenuContent>
        <ContextMenuItem onSelect={() => actions.refresh(tab.path)}>刷新当前</ContextMenuItem>
        <ContextMenuItem disabled={tab.affix} onSelect={() => actions.close(tab.path)}>关闭当前</ContextMenuItem>
        <ContextMenuItem onSelect={() => actions.closeOthers(tab.path)}>关闭其他</ContextMenuItem>
        <ContextMenuItem onSelect={() => actions.closeRight(tab.path)}>关闭右侧</ContextMenuItem>
        <ContextMenuItem onSelect={() => actions.closeAll()}>关闭全部</ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  );
}
