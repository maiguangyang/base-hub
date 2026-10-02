import { matchRoutes } from 'react-router';
import type { ReactElement } from 'react';
import { describe, expect, it } from 'vitest';
import type { ViewerSummary } from '@/features/auth/store/authStore';
import { passwordChangeRedirect, workspaceRedirect } from './AdminApp';
import { adminRoutes, routesForWorkspace } from './adminRoutes';

function viewer(currentWorkspace: ViewerSummary['currentWorkspace']): ViewerSummary {
  return {
    account: {
      id: 'account-1', phone: '13800000000', displayName: 'Owner',
      email: null, status: 'ACTIVE', mustChangePassword: false,
    },
    currentWorkspace,
    workspaces: [],
    permissions: [],
  };
}

it('总部和加盟店共用门店编辑页面组件', () => {
    const hq = routesForWorkspace('HEADQUARTERS').find((route) => route.path === '/admin/hq/stores/manage');
    const franchise = routesForWorkspace('FRANCHISE').find((route) => route.path === '/admin/franchise/stores/manage');
    expect(hq?.element).toBeTruthy();
    expect(franchise?.element).toBeTruthy();
    expect((hq?.element as ReactElement).type).toBe((franchise?.element as ReactElement).type);
});

describe('后台工作台路由', () => {
  it('总部与加盟清单都能匹配且不复用旧 mock 路径', () => {
    expect(matchRoutes(adminRoutes, '/admin/hq/store-approvals')?.at(-1)?.route.path).toBe('/admin/hq/store-approvals');
    expect(matchRoutes(adminRoutes, '/admin/franchise/roles')?.at(-1)?.route.path).toBe('/admin/franchise/roles');
    expect(matchRoutes(adminRoutes, '/admin/obsolete')?.at(-1)?.route.path).toBe('*');
  });

  it('总部管理员和角色路由可被精确匹配', () => {
    expect(matchRoutes(adminRoutes, '/admin/hq/administrators')?.at(-1)?.route.path).toBe('/admin/hq/administrators');
    expect(matchRoutes(adminRoutes, '/admin/hq/roles')?.at(-1)?.route.path).toBe('/admin/hq/roles');
  });

  it('支付配置仅登记在总部工作区路由', () => {
    expect(routesForWorkspace('HEADQUARTERS').some((route) => route.path === '/admin/hq/payment-config')).toBe(true);
    expect(routesForWorkspace('FRANCHISE').some((route) => route.path === '/admin/hq/payment-config')).toBe(false);
  });

  
  it('按工作台返回的路由不包含另一工作台页面', () => {
    expect(routesForWorkspace('HEADQUARTERS').some((route) => route.path?.startsWith('/admin/franchise'))).toBe(false);
    expect(routesForWorkspace('FRANCHISE').some((route) => route.path?.startsWith('/admin/hq'))).toBe(false);
		expect(routesForWorkspace('FRANCHISE').some((route) => route.path === '/admin/hq/administrators')).toBe(false);
  });

  it('已登录深链进入改密或工作台选择时继续携带完整 dt', () => {
    const location = { pathname: '/admin/hq/audit', search: '?page=2', hash: '#latest' };
    expect(workspaceRedirect('password-change', viewer(null), location))
      .toBe('/admin/change-password?dt=%2Fadmin%2Fhq%2Faudit%3Fpage%3D2%23latest');
    expect(workspaceRedirect('workspace-select', viewer(null), location))
      .toBe('/admin/workspaces?dt=%2Fadmin%2Fhq%2Faudit%3Fpage%3D2%23latest');
    const headquarters = { workspaceType: 'HEADQUARTERS' as const, organizationId: 'hq', organizationName: 'HQ', homePath: '/admin/hq' };
    const franchiseLocation = { pathname: '/admin/franchise/stores', search: '?page=3', hash: '' };
    expect(workspaceRedirect('ready', viewer(headquarters), franchiseLocation))
      .toBe('/admin/workspaces?dt=%2Fadmin%2Ffranchise%2Fstores%3Fpage%3D3');
  });

  it('伪工作台前缀不能绕过路径段边界', () => {
		const headquarters = { workspaceType: 'HEADQUARTERS' as const, organizationId: 'hq', organizationName: 'HQ', homePath: '/admin/hq' };
		const franchise = { workspaceType: 'FRANCHISE' as const, organizationId: 'org', organizationName: 'F', homePath: '/admin/franchise' };
		expect(workspaceRedirect('ready', viewer(headquarters), { pathname: '/admin/hqevil', search: '', hash: '' }))
			.toBe('/admin/workspaces');
		expect(workspaceRedirect('ready', viewer(franchise), { pathname: '/admin/franchiseevil', search: '', hash: '' }))
			.toBe('/admin/workspaces');
	});
});

it('总部统一门店页和旧深链均可被精确匹配', () => {
  for (const path of ['/admin/hq/stores', '/admin/hq/franchise-stores', '/admin/hq/direct-stores']) {
    expect(matchRoutes(adminRoutes, path)?.at(-1)?.route.path).toBe(path);
  }
});

it('旧门店网址进入总部工作台时跳到统一门店并保留加盟商条件', () => {
  const headquarters = viewer({ workspaceType: 'HEADQUARTERS', organizationId: 'hq', organizationName: 'HQ', homePath: '/admin/hq' });
  expect(workspaceRedirect('ready', headquarters, {
    pathname: '/admin/hq/franchise-stores', search: '?organizationId=org-a', hash: '',
  })).toBe('/admin/hq/stores?organizationId=org-a&type=FRANCHISE');
});


describe('首次登录改密路由', () => {
  it('仅密码待修改阶段可停留在独立改密页', () => {
    const temporaryViewer = viewer(null);
    temporaryViewer.account.mustChangePassword = true;

    expect(passwordChangeRedirect('password-change', temporaryViewer)).toBeUndefined();
    expect(passwordChangeRedirect('password-change', viewer(null))).toBe('/admin/workspaces');
    expect(passwordChangeRedirect('anonymous', null)).toBe('/admin/login');
    expect(passwordChangeRedirect('workspace-select', viewer(null))).toBe('/admin/workspaces');
    expect(passwordChangeRedirect('ready', viewer({
      workspaceType: 'HEADQUARTERS', organizationId: 'hq', organizationName: 'HQ', homePath: '/admin/hq',
    }))).toBe('/admin/hq');
  });
});
