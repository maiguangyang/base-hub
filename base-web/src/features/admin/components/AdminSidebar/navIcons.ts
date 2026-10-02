import { BadgePercent, Bot, Boxes, Building2, ClipboardCheck, CreditCard, FolderTree, Gauge, Gift, IdCard, Package, Ruler, ScrollText, Settings2, ShieldCheck, Store, Tag, Users, type LucideIcon } from 'lucide-react';

/** 菜单图标键到 lucide 组件的映射。
 *  navigation.ts 必须保持无 React 依赖（它同时被 .astro 的 getStaticPaths 与
 *  Node 契约检查脚本导入），因此映射只能放在渲染侧。新增菜单项时同步登记。 */
export const navIcons: Record<string, LucideIcon> = {
  gauge: Gauge,
  users: Users,
  building: Building2,
  'clipboard-check': ClipboardCheck,
  'scroll-text': ScrollText,
  'shield-check': ShieldCheck,
  store: Store,
  bot: Bot,
  'credit-card': CreditCard,
  gift: Gift,
  'id-card': IdCard,
  'settings-2': Settings2,
  package: Package,
  tag: Tag,
  'folder-tree': FolderTree,
  ruler: Ruler,
  boxes: Boxes,
  'badge-percent': BadgePercent,
};
