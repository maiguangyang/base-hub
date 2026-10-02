import { Link } from 'react-router';
import { SidebarMenuButton, SidebarMenuItem } from '@/components/ui/sidebar';
import type { AdminNavItem } from '@/features/admin/config/navigation';
import { navIcons } from './navIcons';

/** 单个菜单项。isActive 由标签 store 的 activePath 决定，
 *  与标签栏、面包屑共用同一真值源，三者不可能各说各话。 */
export interface NavItemProps {
  item: AdminNavItem;
  isActive: boolean;
}

/** shadcn 的 SidebarMenuButton 把 hover 与 active 都指向 --sidebar-accent，
 *  两者会长得一样。此处把选中态单独指向 --sidebar-primary，
 *  做成实心品牌色胶囊，hover 仍用 accent 的轻微提亮。 */
const activePill = [
  'data-[active=true]:bg-sidebar-primary',
  'data-[active=true]:text-sidebar-primary-foreground',
  'data-[active=true]:hover:bg-sidebar-primary',
  'data-[active=true]:hover:text-sidebar-primary-foreground',
  'data-[active=true]:shadow-sm',
].join(' ');

export function NavItem({ item, isActive }: NavItemProps) {
  const Icon = navIcons[item.icon];
  return (
    <SidebarMenuItem>
      <SidebarMenuButton asChild isActive={isActive} tooltip={item.title} className={`h-9 ${activePill}`}>
        <Link to={item.path}>
          {Icon ? <Icon aria-hidden="true" /> : null}
          <span>{item.title}</span>
        </Link>
      </SidebarMenuButton>
    </SidebarMenuItem>
  );
}
