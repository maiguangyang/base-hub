import { describe, expect, it } from 'vitest';
import { hqStoreFilter } from './storeFilters';

describe('总部统一门店筛选', () => {
  it('默认查询全部门店，并组合类型、所属加盟商和状态', () => {
    expect(hqStoreFilter()).toEqual({});
    expect(hqStoreFilter('FRANCHISE', 'franchise-1', 'ACTIVE', 'OPEN')).toEqual({
      organization: { type: 'FRANCHISE' }, organizationId: 'franchise-1', lifecycle: 'ACTIVE', businessStatus: 'OPEN',
    });
    expect(hqStoreFilter('HEADQUARTERS')).toEqual({ organization: { type: 'HEADQUARTERS' } });
  });
});
