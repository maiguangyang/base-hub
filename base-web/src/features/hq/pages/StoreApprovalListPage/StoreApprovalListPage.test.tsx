// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { print } from 'graphql';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { HQ_STORE_APPROVALS_QUERY } from '../../graphql/storeApprovals';
import { StoreApprovalListPage } from './StoreApprovalListPage';
import { StoreApprovalTable, type StoreApprovalRow } from './StoreApprovalTable';

const row = {
  id: 'store-1', code: 'STR001', name: '江南西店', lifecycle: 'PENDING_APPROVAL',
  organizationId: 'org-1', organization: { id: 'org-1', name: '测试加盟商' },
  contactPhone: '13800138000', managerName: '张店长', managerPhone: '13900139000',
  province: '广东省', city: '广州市', district: '天河区', address: '天河路228号',
  businessHours: '09:00 - 22:00', businessStatus: 'OPEN', supportDineIn: true, supportTakeout: true,
  storeArea: 120, tableCount: 15, submittedAt: '2026-09-30T01:00:00Z',
} as StoreApprovalRow;
const mocks = vi.hoisted(() => ({
  review: vi.fn(async () => ({ data: {} })), refetch: vi.fn(async () => undefined),
  permissions: ['store:approve', 'store:reject'],
}));

vi.mock('@apollo/client/react', () => ({
  useQuery: (document: { definitions: Array<{ name?: { value: string } }> }) => document.definitions[0].name?.value === 'HqStoreApprovals'
    ? { data: { stores: { data: [row], total: 1 } }, loading: false, refetch: mocks.refetch }
    : { data: { organizations: { data: [], total: 0 } }, loading: false },
  useMutation: () => [mocks.review, { loading: false }],
}));
vi.mock('@/features/admin/hooks/useAdminTab', () => ({ useAdminTab: () => ({ search: '' }) }));
vi.mock('@/features/admin/lib/useCompleteFilterCatalog', () => ({
  useCompleteFilterCatalog: () => ({ ready: true, loading: false, items: [] }),
}));
vi.mock('@/features/auth/store/authStore', () => ({
  useAuthStore: (selector: (state: unknown) => unknown) => selector({ viewer: { permissions: mocks.permissions } }),
}));

beforeEach(() => { vi.clearAllMocks(); });
afterEach(cleanup);

it('审核列表按独立列展示联系、地址和经营信息，并查询对应字段', () => {
  render(<StoreApprovalTable rows={[row]} approve reject onReview={() => undefined} />);
  expect(screen.getAllByRole('columnheader').map((cell) => cell.textContent)).toEqual([
    '编码', '门店名称', '所属加盟商', '联系电话', '所在地址', '营业时间', '营业状态', '准入状态',
    '店长', '店长手机号', '服务模式', '经营面积', '桌位数', '提交时间', '操作',
  ]);
  for (const value of ['STR001', '江南西店', '测试加盟商', '13800138000', '张店长', '13900139000',
    '广东省 广州市 天河区 天河路228号', '09:00 - 22:00', '营业中', '待审核', '堂食 · 外卖', '120 ㎡', '15']) {
    expect(screen.getByText(value)).toBeTruthy();
  }
  const source = print(HQ_STORE_APPROVALS_QUERY);
  for (const field of ['contactPhone', 'managerName', 'managerPhone', 'province', 'city', 'district', 'address',
    'businessHours', 'businessStatus', 'supportDineIn', 'supportTakeout', 'storeArea', 'tableCount']) {
    expect(source).toContain(field);
  }
});

it('批准前显示门店名称并二次确认，取消不调用审核', async () => {
  render(<StoreApprovalListPage />);
  fireEvent.click(screen.getByRole('button', { name: '批准' }));
  const dialog = screen.getByRole('dialog');
  expect(within(dialog).getByText(/确认批准“江南西店”/)).toBeTruthy();
  expect(mocks.review).not.toHaveBeenCalled();
  fireEvent.click(within(dialog).getAllByRole('button', { name: '取消' })[0]);
  expect(screen.queryByRole('dialog')).toBeNull();
  fireEvent.click(screen.getByRole('button', { name: '批准' }));
  fireEvent.click(within(screen.getByRole('dialog')).getByRole('button', { name: '确认批准' }));
  await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull());
  expect(mocks.review).toHaveBeenCalledTimes(1);
  expect(mocks.review).toHaveBeenCalledWith({ variables: { input: { storeId: 'store-1', approved: true } } });
});

it('退回原因使用必填多行文本框，保留换行提交', async () => {
  render(<StoreApprovalListPage />);
  fireEvent.click(screen.getByRole('button', { name: '退回' }));
  const field = screen.getByRole('textbox', { name: '退回原因' }) as HTMLTextAreaElement;
  expect(field.tagName).toBe('TEXTAREA');
  expect(field.required).toBe(true);
  expect(field.maxLength).toBe(512);
  expect(field.rows).toBeGreaterThan(1);
  expect(field.placeholder).toBe('请输入退回原因');
  fireEvent.change(field, { target: { value: '请补充营业执照\n请核对门店地址' } });
  fireEvent.click(within(screen.getByRole('dialog')).getByRole('button', { name: '确认退回' }));
  await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull());
  expect(mocks.review).toHaveBeenCalledWith({ variables: { input: {
    storeId: 'store-1', approved: false, rejectionReason: '请补充营业执照\n请核对门店地址',
  } } });
});
