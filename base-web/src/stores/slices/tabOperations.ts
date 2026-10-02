import {
  adminAffixItemFor, adminItemsFor, canonicalHqStoreLocation, isAdminNavItemAllowed, type AdminNavItem, type BusinessWorkspaceType,
} from '@/features/admin/config/navigation';

export interface AdminTab {
  path: string;
  search: string;
  title: string;
  affix: boolean;
  version: number;
}

export interface AdminTabsSnapshot {
  tabs: AdminTab[];
  activePath: string;
}

export interface AdminTabManifest {
  items: AdminNavItem[];
  affix: AdminNavItem;
}

/** 按当前工作台和权限构造标签白名单。 */
export function createAdminTabManifest(workspaceType: BusinessWorkspaceType, permissions: readonly string[]): AdminTabManifest {
  const allowed = new Set(permissions);
  return {
    items: adminItemsFor(workspaceType).filter((item) => isAdminNavItemAllowed(item, allowed)),
    affix: adminAffixItemFor(workspaceType),
  };
}

function createTab(path: string, manifest: AdminTabManifest, search = ''): AdminTab | undefined {
  const canonical = canonicalHqStoreLocation(path, search);
  const item = manifest.items.find((candidate) => candidate.path === canonical.path);
  if (!item) return undefined;
  return { path: item.path, search: canonical.search, title: item.title, affix: item.affix === true, version: 0 };
}

export function createInitialSnapshot(manifest: AdminTabManifest): AdminTabsSnapshot {
  const affix = createTab(manifest.affix.path, manifest);
  if (!affix) throw new Error('固定标签必须位于当前工作台清单');
  return { tabs: [affix], activePath: affix.path };
}

export function openTab(state: AdminTabsSnapshot, path: string, manifest: AdminTabManifest, search = ''): AdminTabsSnapshot {
  ({ path, search } = canonicalHqStoreLocation(path, search));
  const existing = state.tabs.find((tab) => tab.path === path);
  if (existing) {
    if (state.activePath === path && existing.search === search) return state;
    const tabs = existing.search === search
      ? state.tabs
      : state.tabs.map((tab) => tab.path === path ? { ...tab, search } : tab);
    return { ...state, tabs, activePath: path };
  }
  const tab = createTab(path, manifest, search);
  return tab ? { tabs: [...state.tabs, tab], activePath: path } : state;
}

export function closeTab(state: AdminTabsSnapshot, path: string, manifest: AdminTabManifest): AdminTabsSnapshot {
  const index = state.tabs.findIndex((tab) => tab.path === path);
  if (index === -1 || state.tabs[index].affix) return state;
  const tabs = state.tabs.filter((tab) => tab.path !== path);
  if (state.activePath !== path) return { ...state, tabs };
  const next = tabs[index] ?? tabs[tabs.length - 1];
  return { tabs, activePath: next?.path ?? manifest.affix.path };
}

export function closeOtherTabs(state: AdminTabsSnapshot, path: string, manifest: AdminTabManifest): AdminTabsSnapshot {
  const tabs = state.tabs.filter((tab) => tab.affix || tab.path === path);
  const activePath = tabs.some((tab) => tab.path === path) ? path : manifest.affix.path;
  return { tabs, activePath };
}

export function closeRightTabs(state: AdminTabsSnapshot, path: string): AdminTabsSnapshot {
  const index = state.tabs.findIndex((tab) => tab.path === path);
  if (index === -1) return state;
  const tabs = state.tabs.filter((tab, position) => position <= index || tab.affix);
  return { tabs, activePath: tabs.some((tab) => tab.path === state.activePath) ? state.activePath : path };
}

export function closeAllTabs(state: AdminTabsSnapshot, manifest: AdminTabManifest): AdminTabsSnapshot {
  const kept = state.tabs.filter((tab) => tab.affix);
  const tabs = kept.length > 0 ? kept : createInitialSnapshot(manifest).tabs;
  return { tabs, activePath: tabs[0].path };
}

export function refreshTab(state: AdminTabsSnapshot, path: string): AdminTabsSnapshot {
  const tabs = state.tabs.map((tab) => (tab.path === path ? { ...tab, version: tab.version + 1 } : tab));
  return { ...state, tabs };
}

export function sanitizeSnapshot(raw: unknown, manifest: AdminTabManifest): AdminTabsSnapshot {
  const rawActivePath = readActivePath(raw);
  const tabs = rebuildTabs(readCandidates(raw), manifest, rawActivePath);
  if (!tabs.some((tab) => tab.affix)) {
    const affix = createTab(manifest.affix.path, manifest);
    if (affix) tabs.unshift(affix);
  }
  if (tabs.length === 0) return createInitialSnapshot(manifest);
  const requested = rawActivePath ? canonicalHqStoreLocation(rawActivePath).path : undefined;
  const activePath = requested && tabs.some((tab) => tab.path === requested) ? requested : tabs[0].path;
  return { tabs, activePath };
}

function readCandidates(raw: unknown): Partial<AdminTab>[] {
  if (typeof raw !== 'object' || raw === null) return [];
  const tabs = (raw as { tabs?: unknown }).tabs;
  if (!Array.isArray(tabs)) return [];
  return tabs.filter((tab): tab is Partial<AdminTab> => typeof tab === 'object' && tab !== null);
}

function readActivePath(raw: unknown): string | undefined {
  if (typeof raw !== 'object' || raw === null) return undefined;
  const value = (raw as { activePath?: unknown }).activePath;
  return typeof value === 'string' ? value : undefined;
}

function rebuildTabs(candidates: Partial<AdminTab>[], manifest: AdminTabManifest, activePath?: string): AdminTab[] {
  const tabs: AdminTab[] = [];
  for (const candidate of candidates) {
    const search = typeof candidate.search === 'string' ? candidate.search : '';
    const rebuilt = typeof candidate.path === 'string' ? createTab(candidate.path, manifest, search) : undefined;
    if (!rebuilt) continue;
    const next = { ...rebuilt, version: typeof candidate.version === 'number' ? candidate.version : 0 };
    const index = tabs.findIndex((tab) => tab.path === rebuilt.path);
    if (index < 0) tabs.push(next);
    else if (candidate.path === activePath) tabs[index] = next;
  }
  return tabs;
}
