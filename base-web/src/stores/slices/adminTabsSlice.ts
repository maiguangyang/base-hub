import { create } from 'zustand';
import { createJSONStorage, devtools, persist, type StateStorage } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';
import type { BusinessWorkspaceType } from '@/features/admin/config/navigation';
import {
  closeAllTabs, closeOtherTabs, closeRightTabs, closeTab, createAdminTabManifest,
  createInitialSnapshot, openTab, refreshTab, sanitizeSnapshot,
  type AdminTabManifest, type AdminTabsSnapshot,
} from './tabOperations';

interface WorkspaceBinding {
  accountId: string;
  workspaceType: BusinessWorkspaceType;
  organizationId: string | null;
  permissions: readonly string[];
}

interface AdminTabsState extends AdminTabsSnapshot {
  namespace: string;
  manifest: AdminTabManifest;
  erroredPaths: string[];
  openTab(path: string, search?: string): void;
  closeTab(path: string): void;
  closeOthers(path: string): void;
  closeRight(path: string): void;
  closeAll(): void;
  refreshTab(path: string): void;
  markTabError(path: string): void;
}

const initialManifest = createAdminTabManifest('HEADQUARTERS', []);
const memoryEntries = new Map<string, string>();
const memoryStorage: StateStorage = {
  getItem: (name) => memoryEntries.get(name) ?? null,
  setItem: (name, value) => { memoryEntries.set(name, value); },
  removeItem: (name) => { memoryEntries.delete(name); },
};
const tabStorage = createJSONStorage(() => {
  if (typeof window !== 'undefined' && typeof window.localStorage?.setItem === 'function') return window.localStorage;
  return memoryStorage;
});

export const useAdminTabsStore = create<AdminTabsState>()(
  persist(
    devtools(
      immer((set) => ({
        ...createInitialSnapshot(initialManifest),
        namespace: 'unbound', manifest: initialManifest, erroredPaths: [],
        openTab: (path, search) => set((state) => openTab(state, path, state.manifest, search)),
        closeTab: (path) => set((state) => ({ ...closeTab(state, path, state.manifest), erroredPaths: clearError(state.erroredPaths, path) })),
        closeOthers: (path) => set((state) => closeOtherTabs(state, path, state.manifest)),
        closeRight: (path) => set((state) => closeRightTabs(state, path)),
        closeAll: () => set((state) => closeAllTabs(state, state.manifest)),
        refreshTab: (path) => set((state) => ({ ...refreshTab(state, path), erroredPaths: clearError(state.erroredPaths, path) })),
        markTabError: (path) => set((state) => state.erroredPaths.includes(path) ? state : { erroredPaths: [...state.erroredPaths, path] }),
      })),
    ),
    {
      name: 'korean-admin-tabs:unbound', version: 2, skipHydration: true,
      storage: tabStorage,
      partialize: (state) => ({ tabs: state.tabs, activePath: state.activePath }),
      merge: (persisted, current) => ({ ...current, ...sanitizeSnapshot(persisted, current.manifest), erroredPaths: [] }),
    },
  ),
);

/** 生成工作台隔离键。 */
export function workspaceNamespace(accountId: string, workspaceType: BusinessWorkspaceType, organizationId: string | null): string {
  return `${accountId}:${workspaceType}:${organizationId ?? 'none'}`;
}

/** 先清空旧 keep-alive 状态，再从新命名空间恢复合法标签。 */
export async function bindAdminWorkspace(binding: WorkspaceBinding): Promise<void> {
  const namespace = workspaceNamespace(binding.accountId, binding.workspaceType, binding.organizationId);
  if (useAdminTabsStore.getState().namespace === namespace) return;
  const manifest = createAdminTabManifest(binding.workspaceType, binding.permissions);
  useAdminTabsStore.persist.setOptions({ name: `korean-admin-tabs:${namespace}` });
  useAdminTabsStore.setState({ ...createInitialSnapshot(manifest), namespace, manifest, erroredPaths: [] });
  await useAdminTabsStore.persist.rehydrate();
}

/** 清除当前工作台的内存与持久化标签，避免终止会话后残留租户上下文。 */
export function resetAdminWorkspace(): void {
  const state = useAdminTabsStore.getState();
  if (state.namespace !== 'unbound' && typeof window !== 'undefined') {
    window.localStorage?.removeItem(`korean-admin-tabs:${state.namespace}`);
  }
  useAdminTabsStore.persist.setOptions({ name: 'korean-admin-tabs:unbound' });
  useAdminTabsStore.setState({ ...createInitialSnapshot(initialManifest), namespace: 'unbound', manifest: initialManifest, erroredPaths: [] });
}

function clearError(paths: string[], path: string): string[] {
  return paths.filter((item) => item !== path);
}
