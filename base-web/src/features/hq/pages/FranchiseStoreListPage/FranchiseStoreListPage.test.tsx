// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { renderToStaticMarkup } from 'react-dom/server';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { FranchiseStoreListPage, FranchiseStoreLoadError } from './FranchiseStoreListPage';
import { FranchiseStoreTable, type FranchiseStoreRow } from './FranchiseStoreTable';
import { franchiseStoreViewKey, selectedFranchiseId } from './pageState';

const mocks = vi.hoisted(() => ({
  queryCalls: [] as Array<{ name: string | undefined; options: Record<string, unknown> | undefined }>,
  search: '',
  refetch: vi.fn(async () => undefined),
  catalogTotal: 0,
  storeRows: [] as FranchiseStoreRow[],
  permissions: [] as string[],
  paymentFetch: vi.fn(),
  fetchMore: vi.fn(async () => ({ data: { organizations: { data: [{ id: 'franchise-201', name: '第201加盟商' }], total: 201 } } })),
}));

vi.mock('@apollo/client/react', () => ({
  useQuery: (document: { definitions?: Array<{ name?: { value?: string } }> }, options?: Record<string, unknown>) => {
    const name = document.definitions?.find((definition) => definition.name)?.name?.value;
    mocks.queryCalls.push({ name, options });
    if (name === 'HqFranchiseStores') {
      return { data: { stores: { data: mocks.storeRows, total: mocks.storeRows.length } }, loading: false, error: undefined, refetch: mocks.refetch };
    }
    const catalog = mocks.catalogTotal === 201
      ? Array.from({ length: 200 }, (_, index) => ({ id: `franchise-${index + 1}`, name: `第${index + 1}加盟商` }))
      : [];
    return { data: { organizations: { data: catalog, total: mocks.catalogTotal } }, loading: false, error: undefined, fetchMore: mocks.fetchMore };
  },
}));

vi.mock('@/features/admin/hooks/useAdminTab', () => ({
  useAdminTab: () => ({ search: mocks.search }),
}));
vi.mock('@/features/auth/store/authStore', () => ({ useAuthStore: (select: (state: unknown) => unknown) => select({ viewer: {
  currentWorkspace: { workspaceType: 'HEADQUARTERS' }, permissions: mocks.permissions,
} }) }));

function franchiseStoreQueryOptions() {
  return mocks.queryCalls.find((item) => item.name === 'HqFranchiseStores')?.options;
}

beforeEach(() => {
  mocks.queryCalls.length = 0;
  mocks.search = '';
  mocks.catalogTotal = 0;
  mocks.storeRows = [];
  mocks.permissions = [];
  mocks.paymentFetch.mockReset().mockResolvedValue({ ok: true, json: async () => ({ channels: [] }) });
  vi.stubGlobal('fetch', mocks.paymentFetch);
  vi.clearAllMocks();
});

afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

const rows: FranchiseStoreRow[] = [
  { id: 'store-1', code: 'ST-001', name: '草稿店', lifecycle: 'DRAFT', businessStatus: 'CLOSED', contactPhone: null, organizationId: 'org-1', organization: { id: 'org-1', name: '甲加盟商' } },
  { id: 'store-2', code: 'ST-002', name: '待审店', lifecycle: 'PENDING_APPROVAL', businessStatus: 'CLOSED', contactPhone: '13800000000', organizationId: 'org-1', organization: { id: 'org-1', name: '甲加盟商' } },
  { id: 'store-3', code: 'ST-003', name: '营业店', lifecycle: 'ACTIVE', businessStatus: 'OPEN', contactPhone: '13900000000', organizationId: 'org-1', organization: { id: 'org-1', name: '甲加盟商' } },
  { id: 'store-4', code: 'ST-004', name: '退回店', lifecycle: 'REJECTED', businessStatus: 'CLOSED', contactPhone: null, organizationId: 'org-1', organization: { id: 'org-1', name: '甲加盟商' } },
];

describe('加盟门店支付配置入口', () => {
  it.each([0, 1, 3])('未启用门店 %s 不打开支付抽屉或发送配置请求', (index) => {
    mocks.storeRows = [rows[index]];
    mocks.permissions = ['paymentConfig:read'];
    render(<MemoryRouter><FranchiseStoreListPage /></MemoryRouter>);
    const button = screen.getByRole('button', { name: '支付配置' }) as HTMLButtonElement;
    expect(button.disabled).toBe(true);
    expect(button.parentElement?.title).toBe('门店审核通过并启用后才能配置支付。');
    fireEvent.click(button);
    expect(screen.queryByRole('dialog')).toBeNull();
    expect(mocks.paymentFetch).not.toHaveBeenCalled();
  });

  it('有读取权限时从该行打开右侧抽屉并请求该门店', async () => {
    mocks.storeRows = [rows[2]];
    mocks.permissions = ['paymentConfig:read'];
    render(<MemoryRouter><FranchiseStoreListPage /></MemoryRouter>);
    fireEvent.click(screen.getByRole('button', { name: '支付配置' }));
    expect(screen.getByRole('dialog').textContent).toContain('营业店');
    await waitFor(() => expect(mocks.paymentFetch).toHaveBeenCalledWith(
      expect.stringContaining('scope=STORE&storeId=store-3'), expect.objectContaining({ method: 'GET' }),
    ));
  });

  it('缺少支付读取权限时不显示行入口', () => {
    mocks.storeRows = [rows[2]];
    render(<MemoryRouter><FranchiseStoreListPage /></MemoryRouter>);
    expect(screen.queryByRole('button', { name: '支付配置' })).toBeNull();
  });
});

describe('总部加盟门店页', () => {
  it('汇总模式没有组织 ID 时仍执行查询', () => {
    render(<MemoryRouter><FranchiseStoreListPage /></MemoryRouter>);

    expect(franchiseStoreQueryOptions()).toEqual({
      variables: {
        page: 1,
        pageSize: 10,
        q: null,
        filter: { organization: { type: 'FRANCHISE' } },
      },
    });
  });

  it('深链加盟商与关键词共同进入查询变量', () => {
    mocks.search = '?organizationId=org-1&q=%E9%A6%96%E5%B0%94';
    render(<MemoryRouter><FranchiseStoreListPage /></MemoryRouter>);

    expect(franchiseStoreQueryOptions()).toEqual({
      variables: {
        page: 1,
        pageSize: 10,
        q: '首尔',
        filter: { organizationId: 'org-1', organization: { type: 'FRANCHISE' } },
      },
    });
  });

  it('从标签查询上下文读取加盟商标识', () => {
    expect(selectedFranchiseId('?organizationId=org-1')).toBe('org-1');
    expect(selectedFranchiseId('')).toBeUndefined();
    expect(selectedFranchiseId('?organizationId=%20%20')).toBeUndefined();
  });

  it('切换加盟商时生成不同视图键以重置页码与搜索状态', () => {
    expect(franchiseStoreViewKey('?organizationId=org-a')).toBe('org-a');
    expect(franchiseStoreViewKey('?organizationId=org-b')).toBe('org-b');
    expect(franchiseStoreViewKey('')).toBe('all-franchises');
  });

  it('加载失败时提供可访问的错误提示和重试动作', () => {
    const markup = renderToStaticMarkup(<FranchiseStoreLoadError onRetry={() => undefined} />);

    expect(markup).toContain('role="alert"');
    expect(markup).toContain('门店加载失败，请稍后重试。');
    expect(markup).toContain('重试');
  });

  it('只读表格呈现全部生命周期且不暴露经营或审核写操作', () => {
    const markup = renderToStaticMarkup(<FranchiseStoreTable rows={rows} />);

    for (const text of ['所属加盟商', '甲加盟商', 'ST-001', '草稿店', '草稿', '待审核', '已启用', '已退回', '13800000000']) {
      expect(markup).toContain(text);
    }
    expect(markup).toMatch(/aria-checked="true"[^>]*disabled=""[^>]*aria-label="营业店营业状态"/);
    for (const action of ['编辑', '提交审核', '删除', '批准', '退回']) expect(markup).not.toContain(`>${action}<`);
  });
});

it('加盟商超过首批 200 条时加载后续页并启用完整筛选', async () => {
  mocks.catalogTotal = 201;
  render(<MemoryRouter><FranchiseStoreListPage /></MemoryRouter>);

  expect(screen.getByRole('combobox', { name: '全部加盟商' }).hasAttribute('disabled')).toBe(true);
  await waitFor(() => expect(screen.getByRole('combobox', { name: '全部加盟商' }).hasAttribute('disabled')).toBe(false));
  expect(mocks.fetchMore).toHaveBeenCalledWith(expect.objectContaining({ variables: expect.objectContaining({ page: 2 }) }));
});

it('加盟商筛选支持按名称搜索并以选中 ID 查询门店', async () => {
  mocks.catalogTotal = 201;
  render(<MemoryRouter><FranchiseStoreListPage /></MemoryRouter>);

  const trigger = screen.getByRole('combobox', { name: '全部加盟商' });
  await waitFor(() => expect(trigger.hasAttribute('disabled')).toBe(false));
  fireEvent.click(trigger);
  const search = screen.getByPlaceholderText('搜索加盟商');
  fireEvent.change(search, { target: { value: '第201加盟商' } });

  expect(screen.getByRole('option', { name: '第201加盟商' })).toBeTruthy();
  expect(screen.queryByRole('option', { name: '第1加盟商' })).toBeNull();
  fireEvent.keyDown(search, { key: 'ArrowDown' });
  expect(document.activeElement).toBe(screen.getByRole('option', { name: '第201加盟商' }));
  fireEvent.click(screen.getByRole('option', { name: '第201加盟商' }));
  expect(mocks.queryCalls.filter((item) => item.name === 'HqFranchiseStores').at(-1)?.options).toEqual({
    variables: { page: 1, pageSize: 10, q: null, filter: { organizationId: 'franchise-201', organization: { type: 'FRANCHISE' } } },
  });
});
