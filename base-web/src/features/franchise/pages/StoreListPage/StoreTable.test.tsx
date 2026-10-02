// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { StoreTable, type FranchiseStoreRow } from './StoreTable';

afterEach(cleanup);

const base: FranchiseStoreRow = {
  id: 'store-1', code: 'STR001', name: '江南西店', organizationId: 'org-1', lifecycle: 'DRAFT',
  rejectionReason: null, contactPhone: '020-12345678', managerName: '张店长', managerPhone: '13800000000',
  province: '广东省', city: '广州市', district: '海珠区', address: '江南大道1号',
  businessHours: '09:00 - 22:00', businessStatus: 'OPEN', supportDineIn: true, supportTakeout: true,
  storeArea: 86.5, tableCount: 12, receiptFooter: null,
};

function renderRow(overrides: Partial<FranchiseStoreRow> = {}) {
  const row = { ...base, ...overrides };
  const actions = { onEdit: vi.fn(), onSubmit: vi.fn(), onDelete: vi.fn() };
  render(<StoreTable rows={[row]} organizationId="org-1" permissions={['store:update', 'store:submit', 'store:delete']} {...actions} />);
  return { row, ...actions };
}

it.each([
  ['DRAFT', '草稿'], ['PENDING_APPROVAL', '待审核'], ['ACTIVE', '已启用'], ['REJECTED', '已退回'],
] as const)('renders %s as %s in the table', (lifecycle, label) => {
  renderRow({ lifecycle });
  expect(screen.getByText(label)).toBeTruthy();
  expect(screen.queryByText(lifecycle)).toBeNull();
});

it('shows each store field in its own single-line column and preserves actions', () => {
  const actions = renderRow();
  expect(screen.getAllByRole('columnheader').map((header) => header.textContent)).toEqual([
    '编码', '门店名称', '联系电话', '所在地址', '营业时间', '营业状态', '准入状态',
    '店长', '店长手机号', '服务模式', '经营面积', '桌位数', '操作',
  ]);
  for (const text of ['江南西店', 'STR001', '020-12345678', '张店长', '13800000000',
    '广东省 广州市 海珠区 江南大道1号', '营业中', '09:00 - 22:00', '堂食 · 外卖', '86.5 ㎡', '12 桌']) {
    expect(screen.getByText(text)).toBeTruthy();
  }
  const cells = within(screen.getAllByRole('row')[1]).getAllByRole('cell');
  expect(cells).toHaveLength(13);
  expect(cells.slice(0, 12).map((cell) => cell.textContent)).toEqual([
    'STR001', '江南西店', '020-12345678', '广东省 广州市 海珠区 江南大道1号',
    '09:00 - 22:00', '营业中', '草稿', '张店长', '13800000000', '堂食 · 外卖', '86.5 ㎡', '12 桌',
  ]);
  expect(cells.every((cell) => !cell.querySelector('.flex-col'))).toBe(true);
  fireEvent.click(screen.getByRole('button', { name: '编辑' }));
  fireEvent.click(screen.getByRole('button', { name: '提交审核' }));
  fireEvent.click(screen.getByRole('button', { name: '删除' }));
  expect(actions.onEdit).toHaveBeenCalledWith(actions.row);
  expect(actions.onSubmit).toHaveBeenCalledWith(actions.row);
  expect(actions.onDelete).toHaveBeenCalledWith(actions.row);
});

it('shows rejection reason, closed business status and zero capacity accurately', () => {
  renderRow({ lifecycle: 'REJECTED', rejectionReason: '地址不完整', businessStatus: 'CLOSED',
    supportDineIn: false, supportTakeout: true, storeArea: 0, tableCount: 0 });
  expect(screen.getByText('已退回').closest('td')?.title).toBe('地址不完整');
  for (const text of ['已停业', '外卖', '0 ㎡', '0 桌']) expect(screen.getByText(text)).toBeTruthy();
  expect(screen.queryByText('堂食')).toBeNull();
});

it('does not invent business or service settings for incomplete stores', () => {
  renderRow({ contactPhone: null, managerName: null, managerPhone: null, province: null, city: null,
    district: null, address: null, businessHours: null, businessStatus: null,
    supportDineIn: null, supportTakeout: null, storeArea: null, tableCount: null });
  const table = within(screen.getByRole('table'));
  expect(table.getByText('未设置')).toBeTruthy();
  expect(table.getByText('服务模式未设置')).toBeTruthy();
  expect(table.queryByText('营业中')).toBeNull();
  expect(table.queryByText('堂食 · 外卖')).toBeNull();
  expect(table.getAllByText('—').length).toBeGreaterThan(0);
});
