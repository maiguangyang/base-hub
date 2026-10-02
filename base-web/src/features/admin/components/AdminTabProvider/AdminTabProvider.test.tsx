import { createElement as h } from 'react';
import { renderToStaticMarkup } from 'react-dom/server';
import { matchRoutes } from 'react-router';
import { describe, expect, it } from 'vitest';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import type { AdminTab } from '@/stores/slices/tabOperations';
import { AdminTabProvider } from './AdminTabProvider';

/** 构造一个标签记录。默认值对应 navigation.ts 中的真实菜单项。 */
function makeTab(overrides: Partial<AdminTab> = {}): AdminTab {
  return { path: '/admin/franchise/staff', search: '', title: '员工管理', affix: false, version: 0, ...overrides };
}

/** 探针组件：把 useAdminTab 读到的上下文序列化输出，供断言检查。 */
function Probe() {
  const ctx = useAdminTab();
  return h('output', null, JSON.stringify(ctx));
}

/** 复刻 TabStack 的注入方式：params 由调用方从 matches 末项取出后传入。
 *  返回解析后的上下文对象——renderToStaticMarkup 会把引号转成 HTML 实体，
 *  直接对字符串断言会把转义问题误判成机制问题。 */
function renderTab(tab: AdminTab, path: string, isActive: boolean): Record<string, unknown> {
  const routes = [
    { path: '/admin/franchise/staff', element: h(Probe) },
    { path: '/admin/franchise/staff/:id', element: h(Probe) },
    { path: '*', element: h(Probe) },
  ];
  const matches = matchRoutes(routes, path);
  const params = matches?.[matches.length - 1]?.params ?? {};
  const html = renderToStaticMarkup(h(AdminTabProvider, { tab, params, isActive, children: h(Probe) }));
  const inner = /<output>(.*)<\/output>/s.exec(html);
  if (!inner) throw new Error(`探针未渲染出上下文: ${html}`);
  const decoded = inner[1].replace(/&quot;/g, '"').replace(/&amp;/g, '&').replace(/&#x27;/g, "'");
  return JSON.parse(decoded) as Record<string, unknown>;
}

describe('AdminTabProvider 与 useAdminTab', () => {
  it('把标签自身的 path 与 search 注入上下文', () => {
    const ctx = renderTab(makeTab({ search: '?state=1' }), '/admin/franchise/staff', true);
    expect(ctx.path).toBe('/admin/franchise/staff');
    expect(ctx.search).toBe('?state=1');
  });

  it('注入该标签自身的路由参数', () => {
    const ctx = renderTab(makeTab({ path: '/admin/franchise/staff/42' }), '/admin/franchise/staff/42', true);
    expect(ctx.params).toEqual({ id: '42' });
  });

  it('多个标签各自拿到自己的参数，不互相串值', () => {
    const a = renderTab(makeTab({ path: '/admin/franchise/staff/42' }), '/admin/franchise/staff/42', false);
    const b = renderTab(makeTab({ path: '/admin/franchise/staff/99' }), '/admin/franchise/staff/99', true);
    expect(a.params).toEqual({ id: '42' });
    expect(b.params).toEqual({ id: '99' });
  });

  it('isActive 如实反映该标签是否可见', () => {
    expect(renderTab(makeTab(), '/admin/franchise/staff', true).isActive).toBe(true);
    expect(renderTab(makeTab(), '/admin/franchise/staff', false).isActive).toBe(false);
  });

  it('无参数路由返回空参数对象而非 undefined', () => {
    expect(renderTab(makeTab(), '/admin/franchise/staff', true).params).toEqual({});
  });

  it('脱离 Provider 使用时明确抛错，不静默返回空上下文', () => {
    expect(() => renderToStaticMarkup(h(Probe))).toThrow('useAdminTab 必须在 AdminTabProvider 内使用');
  });
});
