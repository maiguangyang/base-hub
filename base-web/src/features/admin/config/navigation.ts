import type { WorkspaceType } from "@/__generated__/graphql";

export type BusinessWorkspaceType = Exclude<WorkspaceType, "DISCOVERY">;

/** 后台菜单项同时定义路由白名单、标签身份与可见权限。 */
export interface AdminNavItem {
  path: string;
  title: string;
  icon: string;
  /** 单个字符串表示一项权限，数组表示必须同时满足的权限。 */
  permission?: string | readonly string[];
  affix?: boolean;
}

export interface AdminNavGroup {
  title: string;
  items: AdminNavItem[];
}

const headquartersNavigation: AdminNavGroup[] = [
  { title: "概览", items: [{ path: "/admin/hq", title: "平台概览", icon: "gauge", affix: true }] },
  {
    title: "组织架构",
    items: [
      { path: "/admin/hq/stores", title: "网点/空间管理", icon: "store", permission: "store:read_all" },
      { path: "/admin/hq/franchises", title: "加盟商", icon: "building", permission: "organization:read" },
      { path: "/admin/hq/store-approvals", title: "审批流管理", icon: "clipboard-check", permission: "store:read_all" },
    ],
  },
  {
    title: "权限管理",
    items: [
      { path: "/admin/hq/administrators", title: "管理员", icon: "users", permission: "hqMembership:read" },
      { path: "/admin/hq/roles", title: "角色与权限", icon: "shield-check", permission: "hqRole:read" },
    ],
  },
  {
    title: "系统安全",
    items: [
      { path: "/admin/hq/audit", title: "审计日志", icon: "scroll-text", permission: "auditLog:read" },
      { path: "/admin/hq/ai-model", title: "AI 模型配置", icon: "bot", permission: "aiModelConfig:read" },
      { path: "/admin/hq/payment-config", title: "接口与支付配置", icon: "credit-card", permission: "paymentConfig:read" },
    ],
  },
];

const franchiseNavigation: AdminNavGroup[] = [
  { title: "概览", items: [{ path: "/admin/franchise", title: "租户工作台", icon: "gauge", affix: true }] },
  {
    title: "成员与安全",
    items: [
      { path: "/admin/franchise/stores", title: "分支空间", icon: "store", permission: "store:read" },
      { path: "/admin/franchise/staff", title: "员工与成员", icon: "users", permission: "operatorMembership:read" },
      { path: "/admin/franchise/roles", title: "角色与权限", icon: "shield-check", permission: "operatorRole:read" },
      { path: "/admin/franchise/audit", title: "操作审计", icon: "scroll-text", permission: "tenantAudit:read" },
    ],
  },
];

const manifests: Record<BusinessWorkspaceType, AdminNavGroup[]> = {
  HEADQUARTERS: headquartersNavigation,
  FRANCHISE: franchiseNavigation,
};

const utilityItems: Partial<Record<BusinessWorkspaceType, AdminNavItem[]>> = {
  HEADQUARTERS: [
    { path: "/admin/hq/stores/manage", title: "维护网点", icon: "store", permission: "store:read_all" },
    { path: "/admin/hq/direct-stores", title: "网点管理", icon: "store", permission: "store:read_all" },
    { path: "/admin/hq/franchise-stores", title: "网点管理", icon: "store", permission: "store:read_all" },
  ],
  FRANCHISE: [
    { path: "/admin/franchise/stores/manage", title: "维护网点", icon: "store", permission: "store:read" },
  ],
};

export const authAdminPaths = ["/admin/initialize", "/admin/login", "/admin/change-password", "/admin/workspaces"];
export const adminNavItems = Object.values(manifests).flatMap((groups) => groups.flatMap((group) => group.items))
  .concat(Object.values(utilityItems).flatMap((items) => items ?? []));
export const knownAdminPaths = ["/admin", ...authAdminPaths, ...adminNavItems.map((item) => item.path)];

/** 检查菜单项的全部必需权限。 */
export function isAdminNavItemAllowed(item: AdminNavItem, allowed: ReadonlySet<string>): boolean {
  const required = typeof item.permission === "string" ? [item.permission] : item.permission ?? [];
  return required.every((permission) => allowed.has(permission));
}

/** 返回按后端权限过滤后的当前工作台导航。 */
export function adminNavigationFor(workspaceType: BusinessWorkspaceType, permissions: readonly string[]): AdminNavGroup[] {
  const allowed = new Set(permissions);
  return manifests[workspaceType]
    .map((group) => ({
      ...group,
      items: group.items.filter((item) => isAdminNavItemAllowed(item, allowed)),
    }))
    .filter((group) => group.items.length > 0);
}

/** 当前工作台全部路由；权限只控制可见性，后端仍是最终授权边界。 */
export function adminItemsFor(workspaceType: BusinessWorkspaceType): AdminNavItem[] {
  return [...manifests[workspaceType].flatMap((group) => group.items), ...(utilityItems[workspaceType] ?? [])];
}

export function adminAffixItemFor(workspaceType: BusinessWorkspaceType): AdminNavItem {
  const affix = adminItemsFor(workspaceType).find((item) => item.affix);
  if (!affix) throw new Error(`工作台 ${workspaceType} 缺少固定首页`);
  return affix;
}

export function findAdminNavItem(path: string): AdminNavItem | undefined {
  return adminNavItems.find((item) => item.path === path);
}

export function isWorkspacePath(workspaceType: BusinessWorkspaceType, path: string): boolean {
  return adminItemsFor(workspaceType).some((item) => item.path === path);
}

/** 将历史门店网址和标签指向总部唯一门店页，同时保留原有筛选条件。 */
export function canonicalHqStoreLocation(path: string, search = ""): { path: string; search: string } {
  const type = path === "/admin/hq/direct-stores" ? "HEADQUARTERS"
    : path === "/admin/hq/franchise-stores" ? "FRANCHISE" : undefined;
  if (!type) return { path, search };
  const params = new URLSearchParams(search);
  params.set("type", type);
  return { path: "/admin/hq/stores", search: `?${params}` };
}
