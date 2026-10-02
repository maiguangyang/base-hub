import { LayoutGrid } from 'lucide-react';
import { Sidebar, SidebarContent, SidebarHeader, SidebarRail } from '@/components/ui/sidebar';
import { adminNavigationFor, type BusinessWorkspaceType } from '@/features/admin/config/navigation';
import { useActiveTabPath } from '@/features/admin/hooks/useAdminTabs';
import { useAuthStore } from '@/features/auth/store/authStore';
import { NavGroup } from './NavGroup';

/** 后台侧边栏。折叠态、图标模式、移动端抽屉与 sidebar_state cookie 持久化
 *  全部由 shadcn Sidebar 提供，不自行实现。
 *  品牌区做成独立色块，与导航项拉开层次——图标折叠态下它会收缩为方块。 */
export function AdminSidebar() {
  const activePath = useActiveTabPath();
  const viewer = useAuthStore((state) => state.viewer);
  const workspaceType = viewer?.currentWorkspace?.workspaceType as BusinessWorkspaceType | undefined;
  const navigation = workspaceType ? adminNavigationFor(workspaceType, viewer?.permissions ?? []) : [];

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader className="p-2">
        <div className="flex items-center gap-2 rounded-lg bg-sidebar-accent px-2.5 py-2 text-sidebar-accent-foreground">
          <LayoutGrid className="size-4 shrink-0" aria-hidden="true" />
          <span className="truncate text-sm font-semibold group-data-[collapsible=icon]:hidden">Korean Hub</span>
        </div>
      </SidebarHeader>
      <SidebarContent className="admin-sidebar-scroll gap-0">
        {navigation.map((group) => (
          <NavGroup key={group.title} group={group} activePath={activePath} />
        ))}
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
  );
}
