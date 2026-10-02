import { afterEach, describe, expect, it, vi } from 'vitest';
import type { ViewerSummary } from '../store/authStore';
import {
  businessWorkspaces, completeLogin, completePasswordChange, currentDestination, loginPath,
  nextAuthPath, returnDestinationFromSearch, safeReturnDestination, validateNewPassword,
} from './authFlow';

function viewer(overrides: Partial<ViewerSummary> = {}): ViewerSummary {
  return {
    account: {
      id: 'account-1', phone: '13800000000', displayName: 'Owner',
      email: null, status: 'ACTIVE', mustChangePassword: false,
    },
    currentWorkspace: null,
    workspaces: [],
    permissions: [],
    ...overrides,
  };
}

afterEach(() => vi.restoreAllMocks());

describe('认证页面流程', () => {
  it('临时密码、待选工作台和已就绪状态分别进入固定路由', () => {
    expect(nextAuthPath(viewer({ account: { ...viewer().account, mustChangePassword: true } }))).toBe('/admin/change-password');
    expect(nextAuthPath(viewer())).toBe('/admin/workspaces');
    const currentWorkspace = { workspaceType: 'HEADQUARTERS' as const, organizationId: 'hq', organizationName: 'HQ', homePath: '/admin/hq' };
    expect(nextAuthPath(viewer({ currentWorkspace }))).toBe('/admin/hq');
  });

  it('密码只校验 8–20 个字符，并要求确认一致', () => {
    expect(validateNewPassword('1234567', '1234567')).toBe('PASSWORD_WEAK');
    expect(validateNewPassword('password', 'password')).toBeUndefined();
    expect(validateNewPassword('密'.repeat(20), '密'.repeat(20))).toBeUndefined();
    expect(validateNewPassword('密'.repeat(21), '密'.repeat(21))).toBe('PASSWORD_WEAK');
    expect(validateNewPassword('Correct-Horse-42', 'different')).toBe('PASSWORD_CONFIRM_MISMATCH');
    expect(validateNewPassword('Correct-Horse-42', 'Correct-Horse-42')).toBeUndefined();
  });

  it('工作台选择不把 DISCOVERY 当作业务工作台', () => {
    const workspaces = [
      { workspaceType: 'DISCOVERY' as const, organizationId: null, organizationName: 'discovery', homePath: '/admin' },
      { workspaceType: 'FRANCHISE' as const, organizationId: 'org-1', organizationName: 'One', homePath: '/admin/franchise' },
    ];
    expect(businessWorkspaces(workspaces)).toHaveLength(1);
  });

  it('完整保留后台路径、查询参数和 hash 并编码到登录地址', () => {
    const destination = currentDestination({ pathname: '/admin/hq/franchises', search: '?page=2&q=%E9%9F%A9%E6%96%87', hash: '#row-1' });
    expect(destination).toBe('/admin/hq/franchises?page=2&q=%E9%9F%A9%E6%96%87#row-1');
    expect(loginPath(destination)).toBe('/admin/login?dt=%2Fadmin%2Fhq%2Ffranchises%3Fpage%3D2%26q%3D%25E9%259F%25A9%25E6%2596%2587%23row-1');
    expect(returnDestinationFromSearch(`?dt=${encodeURIComponent(destination)}`)).toBe(destination);
  });

  it('拒绝站外、协议相对、畸形和认证流程回跳地址', () => {
    for (const value of ['https://evil.example/admin/hq', '//evil.example/admin/hq', '%', '/admin/hq/../login', '/admin/hq/%2e%2e/login', '/admin/login', '/admin/initialize', '/admin/change-password', '/admin/workspaces', '/public']) {
      expect(safeReturnDestination(value)).toBeUndefined();
    }
  });

  it('只在目标与当前工作台匹配时完成回跳', () => {
    const headquarters = { workspaceType: 'HEADQUARTERS' as const, organizationId: 'hq', organizationName: 'HQ', homePath: '/admin/hq' };
    const franchise = { workspaceType: 'FRANCHISE' as const, organizationId: 'org-1', organizationName: 'One', homePath: '/admin/franchise' };
    expect(nextAuthPath(viewer({ currentWorkspace: headquarters }), '/admin/hq/audit?q=x')).toBe('/admin/hq/audit?q=x');
    expect(nextAuthPath(viewer({ currentWorkspace: headquarters }), '/admin/hq?tab=today#sales')).toBe('/admin/hq?tab=today#sales');
    expect(nextAuthPath(viewer({ currentWorkspace: headquarters }), '/admin/franchise/stores')).toBe('/admin/hq');
    expect(nextAuthPath(viewer({ currentWorkspace: franchise }), '/admin/franchise/stores?page=3')).toBe('/admin/franchise/stores?page=3');
    expect(nextAuthPath(viewer({ currentWorkspace: franchise }), '/admin/franchise?tab=stores#active')).toBe('/admin/franchise?tab=stores#active');
  });
});

describe('认证后的 GraphQL runtime 隔离', () => {
  it('改密后先销毁旧 GraphQL runtime，再以新 Cookie 完整导航', async () => {
    const calls: string[] = [];
    const nextViewer = viewer();
    await completePasswordChange(nextViewer, () => calls.push('viewer'), async () => { calls.push('dispose'); }, (path) => calls.push(`replace:${path}`), '/admin/hq/audit?page=2');
    expect(calls).toEqual(['viewer', 'dispose', 'replace:/admin/workspaces?dt=%2Fadmin%2Fhq%2Faudit%3Fpage%3D2']);
  });

  it('重新登录先写入 Viewer、销毁旧 runtime，再完整替换页面', async () => {
    const calls: string[] = [];
    const currentWorkspace = { workspaceType: 'HEADQUARTERS' as const, organizationId: 'hq', organizationName: 'HQ', homePath: '/admin/hq' };
    await completeLogin(viewer({ currentWorkspace }), () => calls.push('viewer'), async () => { calls.push('dispose'); }, (path) => calls.push(`replace:${path}`), '/admin/hq/audit?page=2');
    expect(calls).toEqual(['viewer', 'dispose', 'replace:/admin/hq/audit?page=2']);
  });

  it('本地 runtime 缓存清理失败时仍完成密码更新后的页面替换', async () => {
    const diagnostics = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    const calls: string[] = [];
    const currentWorkspace = { workspaceType: 'HEADQUARTERS' as const, organizationId: 'hq', organizationName: 'HQ', homePath: '/admin/hq' };

    await expect(completePasswordChange(
      viewer({ currentWorkspace }),
      () => calls.push('viewer'),
      async () => { calls.push('dispose'); throw new Error('cache cleanup failed'); },
      (path) => calls.push(`replace:${path}`),
      '/admin/hq/audit?page=2',
    )).resolves.toBeUndefined();

    expect(calls).toEqual(['viewer', 'dispose', 'replace:/admin/hq/audit?page=2']);
    expect(diagnostics).toHaveBeenCalledWith(
      '认证切换后的 GraphQL runtime 清理失败，继续执行安全导航。',
      expect.any(Error),
    );
  });
});
