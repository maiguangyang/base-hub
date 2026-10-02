import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { expect, it } from 'vitest';

const featureRoot = fileURLToPath(new URL('../../', import.meta.url));
const switchPages = [
  'hq/pages/AdministratorListPage/AdministratorTable.tsx',
  'hq/pages/DirectStoreListPage/DirectStoreTable.tsx',
  'hq/pages/FranchiseStoreListPage/FranchiseStoreTable.tsx',
  'franchise/pages/StaffListPage/StaffTable.tsx',
];

it('shows editable binary statuses as shared switches outside the action column', () => {
  const missing = switchPages.filter((page) => !readFileSync(`${featureRoot}/${page}`, 'utf8').includes('<AdminStatusSwitch'));
  expect(missing).toEqual([]);
  for (const page of switchPages) {
    const source = readFileSync(`${featureRoot}/${page}`, 'utf8');
    expect(source, page).not.toMatch(/<Button[^>]*>\{[^}]*\? '停用' : '启用'\}<\/Button>/);
  }
});
