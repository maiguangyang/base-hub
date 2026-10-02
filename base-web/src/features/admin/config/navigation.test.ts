import { describe, expect, it } from 'vitest';
import { adminNavigationFor, isWorkspacePath, knownAdminPaths } from './navigation';

it('双工作台门店编辑路径只在各自标签清单中，不进入侧栏', () => {
    expect(isWorkspacePath('HEADQUARTERS', '/admin/hq/stores/manage')).toBe(true);
    expect(isWorkspacePath('FRANCHISE', '/admin/franchise/stores/manage')).toBe(true);
    expect(isWorkspacePath('HEADQUARTERS', '/admin/franchise/stores/manage')).toBe(false);
    expect(isWorkspacePath('FRANCHISE', '/admin/hq/stores/manage')).toBe(false);
    expect(adminNavigationFor('HEADQUARTERS', ['store:read_all']).flatMap((group) => group.items.map((item) => item.path))).not.toContain('/admin/hq/stores/manage');
    expect(adminNavigationFor('FRANCHISE', ['store:read']).flatMap((group) => group.items.map((item) => item.path))).not.toContain('/admin/franchise/stores/manage');
});

describe('工作台导航', () => {
  it('总部与加盟工作台使用完全独立的路由清单', () => {
    const hq = adminNavigationFor('HEADQUARTERS', ['organization:read', 'store:read_all']);
    const franchise = adminNavigationFor('FRANCHISE', ['store:read', 'operatorMembership:read']);
    const hqPaths = hq.flatMap((group) => group.items.map((item) => item.path));
    const franchisePaths = franchise.flatMap((group) => group.items.map((item) => item.path));
    expect(hqPaths).toContain('/admin/hq/franchises');
    expect(hqPaths.some((path) => path.startsWith('/admin/franchise'))).toBe(false);
    expect(franchisePaths).toContain('/admin/franchise/stores');
    expect(franchisePaths.some((path) => path.startsWith('/admin/hq'))).toBe(false);
  });

  it('总部侧边栏只显示一个统一门店入口，旧路径仍属于总部工作台', () => {
    const visible = adminNavigationFor('HEADQUARTERS', ['organization:read', 'store:read_all'])
      .flatMap((group) => group.items.map((item) => item.path));

    expect(visible).toContain('/admin/hq/stores');
    expect(visible).not.toContain('/admin/hq/direct-stores');
    expect(visible).not.toContain('/admin/hq/franchise-stores');
    expect(knownAdminPaths).toContain('/admin/hq/stores');
    expect(knownAdminPaths).toContain('/admin/hq/franchise-stores');
    expect(isWorkspacePath('HEADQUARTERS', '/admin/hq/franchise-stores')).toBe(true);
  });

  it('缺少权限时隐藏菜单，但固定首页始终可见', () => {
    const groups = adminNavigationFor('HEADQUARTERS', []);
    expect(groups.flatMap((group) => group.items.map((item) => item.path))).toEqual(['/admin/hq']);
  });

  it('总部管理员与角色菜单分别受读取权限控制', () => {
    const administrators = adminNavigationFor('HEADQUARTERS', ['hqMembership:read'])
      .flatMap((group) => group.items.map((item) => item.path));
    const roles = adminNavigationFor('HEADQUARTERS', ['hqRole:read'])
      .flatMap((group) => group.items.map((item) => item.path));
    expect(administrators).toContain('/admin/hq/administrators');
    expect(administrators).not.toContain('/admin/hq/roles');
    expect(roles).toContain('/admin/hq/roles');
    expect(roles).not.toContain('/admin/hq/administrators');
  });

  it('AI 模型配置仅对有权限的总部工作区开放', () => {
    const path = '/admin/hq/ai-model';
    const permitted = adminNavigationFor('HEADQUARTERS', ['aiModelConfig:read']).flatMap((group) => group.items.map((item) => item.path));
    const denied = adminNavigationFor('HEADQUARTERS', []).flatMap((group) => group.items.map((item) => item.path));
    expect(permitted).toContain(path);
    expect(denied).not.toContain(path);
    expect(isWorkspacePath('FRANCHISE', path)).toBe(false);
  });

  it('静态路径包含认证页及两个工作台，工作台守卫拒绝交叉路径', () => {
    expect(knownAdminPaths).toEqual(expect.arrayContaining(['/admin', '/admin/initialize', '/admin/login', '/admin/change-password', '/admin/workspaces', '/admin/hq', '/admin/franchise']));
    expect(isWorkspacePath('HEADQUARTERS', '/admin/hq/franchises')).toBe(true);
    expect(isWorkspacePath('HEADQUARTERS', '/admin/franchise/stores')).toBe(false);
  });
});

it('支付配置仅对有读取权限的总部工作区开放', () => {
  const path = '/admin/hq/payment-config';
  const permitted = adminNavigationFor('HEADQUARTERS', ['paymentConfig:read']).flatMap((group) => group.items.map((item) => item.path));
  expect(permitted).toContain(path);
  expect(adminNavigationFor('HEADQUARTERS', []).flatMap((group) => group.items.map((item) => item.path))).not.toContain(path);
  expect(isWorkspacePath('FRANCHISE', path)).toBe(false);
});



