import { describe, expect, it } from 'vitest';
import {
  closeAllTabs, closeOtherTabs, closeRightTabs, closeTab, createAdminTabManifest,
  createInitialSnapshot, openTab, refreshTab, sanitizeSnapshot, type AdminTabsSnapshot,
} from './tabOperations';

const manifest = createAdminTabManifest('HEADQUARTERS', [
  'organization:read', 'store:read_all', 'auditLog:read',
]);

function snapshotWithAllPages(): AdminTabsSnapshot {
  let state = createInitialSnapshot(manifest);
  state = openTab(state, '/admin/hq/direct-stores', manifest);
  state = openTab(state, '/admin/hq/franchises', manifest);
  return openTab(state, '/admin/hq/audit', manifest);
}

describe('标签打开与关闭', () => {
    it('统一门店入口只要求跨组织门店读权限，旧路径仍可恢复标签', () => {
    const storesOnly = createAdminTabManifest('HEADQUARTERS', ['store:read_all']);
    const complete = createAdminTabManifest('HEADQUARTERS', ['organization:read', 'store:read_all']);

    expect(storesOnly.items.map((item) => item.path)).toContain('/admin/hq/stores');
    expect(storesOnly.items.map((item) => item.path)).toContain('/admin/hq/franchise-stores');
    expect(storesOnly.items.map((item) => item.path)).not.toContain('/admin/hq/franchises');
    expect(complete.items.map((item) => item.path)).toContain('/admin/hq/franchise-stores');
  });

  it('打开已存在标签只激活，不重复插入', () => {
    const once = openTab(createInitialSnapshot(manifest), '/admin/hq/franchises', manifest);
    const twice = openTab(once, '/admin/hq/franchises', manifest);
    expect(twice.tabs.filter((tab) => tab.path === '/admin/hq/franchises')).toHaveLength(1);
  });

  it('拒绝未登记路径并拒绝关闭固定首页', () => {
    const initial = createInitialSnapshot(manifest);
    expect(openTab(initial, '/admin/obsolete', manifest)).toBe(initial);
    expect(closeTab(initial, '/admin/hq', manifest)).toBe(initial);
  });

  it('关闭活动标签时选择相邻标签', () => {
    const state = { ...snapshotWithAllPages(), activePath: '/admin/hq/franchises' };
    expect(closeTab(state, '/admin/hq/franchises', manifest).activePath).toBe('/admin/hq/audit');
  });

  it('同一加盟门店页按路径复用标签并更新当前查询上下文', () => {
    const initial = createInitialSnapshot(manifest);
    const first = openTab(initial, '/admin/hq/franchise-stores', manifest, '?organizationId=org-a');
    const second = openTab(first, '/admin/hq/franchise-stores', manifest, '?organizationId=org-b');

    expect(second.tabs.filter((tab) => tab.path === '/admin/hq/stores')).toHaveLength(1);
    expect(second.tabs.find((tab) => tab.path === '/admin/hq/stores')?.search)
      .toBe('?organizationId=org-b&type=FRANCHISE');
  });
});

describe('标签批量操作', () => {
  it('关闭其他、右侧和全部始终保留固定首页', () => {
    const state = snapshotWithAllPages();
    expect(closeOtherTabs(state, '/admin/hq/franchises', manifest).tabs.map((tab) => tab.path)).toEqual(['/admin/hq', '/admin/hq/franchises']);
    expect(closeRightTabs(state, '/admin/hq/stores').tabs.map((tab) => tab.path)).toEqual(['/admin/hq', '/admin/hq/stores']);
    expect(closeAllTabs(state, manifest).tabs.map((tab) => tab.path)).toEqual(['/admin/hq']);
  });

  it('刷新只递增目标标签 version', () => {
    const state = snapshotWithAllPages();
    const after = refreshTab(state, '/admin/hq/franchises');
    expect(after.tabs.find((tab) => tab.path === '/admin/hq/franchises')?.version).toBe(1);
    expect(after.tabs.find((tab) => tab.path === '/admin/hq/audit')?.version).toBe(0);
  });
});

describe('持久化快照校验', () => {
  it('丢弃陈旧、重复或越权路由并刷新标题', () => {
    const raw = {
      tabs: [
        { path: '/admin/hq/franchises', title: '旧标题', version: 2 },
        { path: '/admin/hq/franchises', title: '重复项' },
        { path: '/admin/franchise/stores', title: '另一工作台' },
      ],
      activePath: '/admin/franchise/stores',
    };
    const after = sanitizeSnapshot(raw, manifest);
    expect(after.tabs.map((tab) => tab.path)).toEqual(['/admin/hq', '/admin/hq/franchises']);
    expect(after.tabs[1].title).toBe('加盟商');
    expect(after.activePath).toBe('/admin/hq');
  });

  it('任意非法输入都返回可用初始快照', () => {
    for (const raw of [undefined, null, 'corrupted', 42, {}, { tabs: 'invalid' }]) {
      expect(sanitizeSnapshot(raw, manifest)).toEqual(createInitialSnapshot(manifest));
    }
  });

  it('恢复字符串查询上下文并丢弃非法查询值', () => {
    const valid = sanitizeSnapshot({
      tabs: [{ path: '/admin/hq/franchise-stores', search: '?organizationId=org-a' }],
      activePath: '/admin/hq/franchise-stores',
    }, manifest);
    const invalid = sanitizeSnapshot({
      tabs: [{ path: '/admin/hq/franchise-stores', search: { organizationId: 'org-a' } }],
      activePath: '/admin/hq/franchise-stores',
    }, manifest);

    expect(valid.tabs.find((tab) => tab.path === '/admin/hq/stores')?.search)
      .toBe('?organizationId=org-a&type=FRANCHISE');
    expect(invalid.tabs.find((tab) => tab.path === '/admin/hq/stores')?.search).toBe('?type=FRANCHISE');
    expect(valid.activePath).toBe('/admin/hq/stores');
  });

  it('同时恢复旧直营和加盟标签时保留原活动标签的加盟商条件', () => {
    const restored = sanitizeSnapshot({
      tabs: [
        { path: '/admin/hq/direct-stores', search: '?q=直营' },
        { path: '/admin/hq/franchise-stores', search: '?organizationId=org-a' },
      ],
      activePath: '/admin/hq/franchise-stores',
    }, manifest);
    expect(restored.tabs.filter((tab) => tab.path === '/admin/hq/stores')).toHaveLength(1);
    expect(restored.tabs.find((tab) => tab.path === '/admin/hq/stores')?.search)
      .toBe('?organizationId=org-a&type=FRANCHISE');
  });
});
