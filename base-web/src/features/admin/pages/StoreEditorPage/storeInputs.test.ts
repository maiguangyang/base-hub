import { describe, expect, it } from 'vitest';
import { defaultStoreValues, storeCreateInput, storeUpdateInput, storeValuesFromRecord } from './storeInputs';

describe('shared store inputs', () => {
  it('creates an active headquarters store and a draft franchise store with generated codes', () => {
    expect(storeCreateInput({ ...defaultStoreValues(), name: ' A ' }, 'hq', 'HEADQUARTERS')).toMatchObject({ name: 'A', organizationId: 'hq', lifecycle: 'ACTIVE' });
    const franchise = storeCreateInput({ ...defaultStoreValues(), name: ' B ' }, 'tenant', 'FRANCHISE');
    expect(franchise).toMatchObject({ name: 'B', organizationId: 'tenant', lifecycle: 'DRAFT' });
    expect(franchise.code.startsWith('STR')).toBe(true);
  });

  it('keeps zero values and normalizes nullable fields in records', () => {
    expect(storeValuesFromRecord({ name: '店', storeArea: 0, tableCount: 0 })).toMatchObject({ storeArea: '0', tableCount: '0' });
    expect(storeValuesFromRecord({ name: '店', contactPhone: null, supportDineIn: null, businessStatus: null })).toMatchObject({ contactPhone: '', supportDineIn: true, businessStatus: 'OPEN' });
  });

  it('updates only mutable fields', () => {
    const update = storeUpdateInput({ ...defaultStoreValues(), name: ' 店 ', storeArea: '0', tableCount: '0' });
    expect(update).toMatchObject({ name: '店', storeArea: 0, tableCount: 0 });
    expect(update).not.toHaveProperty('lifecycle');
    expect(update).not.toHaveProperty('organizationId');
    expect(update).not.toHaveProperty('code');
  });
});
