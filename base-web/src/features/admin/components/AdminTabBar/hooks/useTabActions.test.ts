import { describe, expect, it } from 'vitest';
import type { AdminTab } from '@/stores/slices/tabOperations';
import { destinationForTab } from './useTabActions';

const tabs: AdminTab[] = [{
  path: '/admin/hq/franchise-stores',
  search: '?organizationId=org-a',
  title: '加盟门店',
  affix: false,
  version: 0,
}];

describe('标签地址恢复', () => {
  it('激活标签时恢复其查询上下文', () => {
    expect(destinationForTab('/admin/hq/franchise-stores', tabs))
      .toBe('/admin/hq/franchise-stores?organizationId=org-a');
  });

  it('未知标签回退到裸路径', () => {
    expect(destinationForTab('/admin/hq', tabs)).toBe('/admin/hq');
  });
});
