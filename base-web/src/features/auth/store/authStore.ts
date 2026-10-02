import { create } from 'zustand';
import type { ViewerFieldsFragment } from '@/__generated__/graphql';

export type ViewerSummary = ViewerFieldsFragment;
export type AuthPhase = 'loading' | 'anonymous' | 'password-change' | 'workspace-select' | 'ready';

interface AuthState {
  viewer: ViewerSummary | null;
  phase: AuthPhase;
  setViewer(viewer: ViewerSummary): void;
  reset(): void;
}

/** 认证状态只保存服务端 Viewer 摘要，不保存密码、Cookie 或令牌。 */
export const useAuthStore = create<AuthState>((set) => ({
  viewer: null,
  phase: 'loading',
  setViewer: (viewer) => set({ viewer, phase: phaseForViewer(viewer) }),
  reset: () => set({ viewer: null, phase: 'anonymous' }),
}));

function phaseForViewer(viewer: ViewerSummary): AuthPhase {
  if (viewer.account.mustChangePassword) return 'password-change';
  if (!viewer.currentWorkspace || viewer.currentWorkspace.workspaceType === 'DISCOVERY') return 'workspace-select';
  return 'ready';
}
