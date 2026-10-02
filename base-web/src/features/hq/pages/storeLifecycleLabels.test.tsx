// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, expect, it } from 'vitest';
import { HqStoreTable, type HqStoreRow } from './HqStoreListPage/HqStoreTable';
import { DirectStoreTable } from './DirectStoreListPage/DirectStoreTable';
import { FranchiseStoreTable } from './FranchiseStoreListPage/FranchiseStoreTable';
import { StoreApprovalTable } from './StoreApprovalListPage/StoreApprovalTable';

afterEach(cleanup);

const base = {
  id: 'store', code: 'STR001', name: '测试门店', organizationId: 'org',
  organization: { id: 'org', name: '加盟商', type: 'FRANCHISE' },
  lifecycle: 'DRAFT', businessStatus: 'OPEN', contactPhone: null,
  city: null, district: null, address: null, businessHours: null,
} as HqStoreRow;

it.each([
  ['DRAFT', '草稿'], ['PENDING_APPROVAL', '待审核'], ['ACTIVE', '已启用'], ['REJECTED', '已退回'],
] as const)('shows %s as %s consistently in HQ store tables', (lifecycle, label) => {
  const row = { ...base, lifecycle, submittedAt: null, province: null, managerName: null, managerPhone: null,
    supportDineIn: null, supportTakeout: null, storeArea: null, tableCount: null };
  const tables = [
    <HqStoreTable rows={[row]} hqOrganizationId="hq" canUpdate={false} canDelete={false} statusBusy={false}
      onEdit={() => undefined} onDelete={() => undefined} onBusinessStatus={() => undefined} />,
    <DirectStoreTable rows={[row]} canUpdate={false} canDelete={false} statusBusy={false}
      onEdit={() => undefined} onDelete={() => undefined} onBusinessStatus={() => undefined} />,
    <FranchiseStoreTable rows={[row]} />,
    <StoreApprovalTable rows={[row]} approve={false} reject={false} onReview={() => undefined} />,
  ];
  for (const table of tables) {
    const view = render(table);
    expect(screen.getByText(label)).toBeTruthy();
    expect(screen.queryByText(lifecycle)).toBeNull();
    view.unmount();
  }
});
