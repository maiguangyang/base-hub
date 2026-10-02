import { beforeEach, describe, expect, it } from 'vitest';
import { useAuthStore, type ViewerSummary } from './authStore';

const viewer: ViewerSummary = {
  account: {
    id: 'account-1', phone: '13800000000', displayName: 'Owner',
    email: null, status: 'ACTIVE', mustChangePassword: false,
  },
  currentWorkspace: {
    workspaceType: 'FRANCHISE', organizationId: 'org-1',
    organizationName: 'One', homePath: '/admin/franchise',
  },
  workspaces: [],
  permissions: ['store:read'],
};

describe('认证状态', () => {
  beforeEach(() => useAuthStore.getState().reset());

  it('只保存 Viewer 摘要并进入 ready 阶段', () => {
    useAuthStore.getState().setViewer(viewer);
    expect(useAuthStore.getState().phase).toBe('ready');
    expect(useAuthStore.getState().viewer?.account.id).toBe('account-1');
  });

  it('序列化状态不含任何凭据或令牌字段', () => {
    useAuthStore.getState().setViewer(viewer);
    const serialized = JSON.stringify(useAuthStore.getState()).toLowerCase();
    for (const secret of ['password', 'temporaryPassword', 'cookie', 'jwt', 'bearer', 'token']) {
      expect(serialized).not.toContain(`"${secret.toLowerCase()}"`);
    }
  });
});
