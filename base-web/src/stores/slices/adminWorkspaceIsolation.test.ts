import { beforeEach, describe, expect, it } from 'vitest';
import { bindAdminWorkspace, useAdminTabsStore, workspaceNamespace } from './adminTabsSlice';

describe('后台标签工作台隔离', () => {
  beforeEach(() => {
    useAdminTabsStore.setState({ tabs: [], activePath: '', erroredPaths: [] });
  });

  it('使用账号、工作台类型和组织组成精确命名空间', () => {
    expect(workspaceNamespace('account-1', 'FRANCHISE', 'org-a')).toBe('account-1:FRANCHISE:org-a');
    expect(workspaceNamespace('account-1', 'HEADQUARTERS', null)).toBe('account-1:HEADQUARTERS:none');
  });

  it('切换组织会先清空旧挂载标签与错误状态，再绑定新的持久化 key', async () => {
    await bindAdminWorkspace({ accountId: 'account-1', workspaceType: 'FRANCHISE', organizationId: 'org-a', permissions: ['store:read'] });
    useAdminTabsStore.getState().openTab('/admin/franchise/stores');
    useAdminTabsStore.getState().markTabError('/admin/franchise/stores');
    await bindAdminWorkspace({ accountId: 'account-1', workspaceType: 'FRANCHISE', organizationId: 'org-b', permissions: [] });
    expect(useAdminTabsStore.persist.getOptions().name).toBe('korean-admin-tabs:account-1:FRANCHISE:org-b');
    expect(useAdminTabsStore.getState().tabs.map((tab) => tab.path)).toEqual(['/admin/franchise']);
    expect(useAdminTabsStore.getState().erroredPaths).toEqual([]);
  });

  it('切换账号改变命名空间并丢弃不属于新权限清单的陈旧路径', async () => {
    await bindAdminWorkspace({ accountId: 'account-1', workspaceType: 'HEADQUARTERS', organizationId: 'hq', permissions: ['organization:read'] });
    useAdminTabsStore.getState().openTab('/admin/hq/franchises');
    await bindAdminWorkspace({ accountId: 'account-2', workspaceType: 'HEADQUARTERS', organizationId: 'hq', permissions: [] });
    expect(useAdminTabsStore.getState().namespace).toBe('account-2:HEADQUARTERS:hq');
    expect(useAdminTabsStore.getState().tabs.map((tab) => tab.path)).toEqual(['/admin/hq']);
  });
});
