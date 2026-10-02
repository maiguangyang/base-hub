import { SidebarGroup, SidebarGroupContent, SidebarGroupLabel, SidebarMenu } from '@/components/ui/sidebar';
import type { AdminNavGroup } from '@/features/admin/config/navigation';
import { NavItem } from './NavItem';

/** 一个菜单分组。分组标题不可点击，仅作视觉归类，不参与路由。 */
export interface NavGroupProps {
  group: AdminNavGroup;
  activePath: string;
}

export function NavGroup({ group, activePath }: NavGroupProps) {
  return (
    <SidebarGroup>
      <SidebarGroupLabel>{group.title}</SidebarGroupLabel>
      <SidebarGroupContent>
        <SidebarMenu>
          {group.items.map((item) => (
            <NavItem key={item.path} item={item} isActive={item.path === activePath} />
          ))}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  );
}
