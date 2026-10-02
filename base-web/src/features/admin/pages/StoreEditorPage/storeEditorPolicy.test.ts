import { describe, expect, it } from 'vitest';
import { canCreateStore, canEditStoreRecord, canSubmitFranchiseStore } from './storeEditorPolicy';

describe('store editor policy', () => {
  it('requires workspace-specific create permission', () => {
    expect(canCreateStore('HEADQUARTERS', ['hqStore:create'])).toBe(true);
    expect(canCreateStore('FRANCHISE', ['store:create'])).toBe(true);
    expect(canCreateStore('HEADQUARTERS', ['store:create'])).toBe(false);
    expect(canCreateStore('FRANCHISE', ['hqStore:create'])).toBe(false);
  });

  it('edits only headquarters-owned direct stores with permission', () => {
    const own = { organizationId: 'hq', organization: { type: 'HEADQUARTERS' }, lifecycle: 'ACTIVE' };
    expect(canEditStoreRecord('HEADQUARTERS', 'hq', ['hqStore:update'], own)).toBe(true);
    expect(canEditStoreRecord('HEADQUARTERS', 'hq', [], own)).toBe(false);
    expect(canEditStoreRecord('HEADQUARTERS', 'hq', ['hqStore:update'], { ...own, organizationId: 'tenant' })).toBe(false);
    expect(canEditStoreRecord('HEADQUARTERS', 'hq', ['hqStore:update'], { ...own, organization: { type: 'FRANCHISE' } })).toBe(false);
  });

  it('edits only own draft or rejected franchise stores and gates submit', () => {
    for (const lifecycle of ['DRAFT', 'REJECTED']) {
      const row = { organizationId: 'tenant', lifecycle };
      expect(canEditStoreRecord('FRANCHISE', 'tenant', ['store:update'], row)).toBe(true);
      expect(canEditStoreRecord('FRANCHISE', 'tenant', [], row)).toBe(false);
      expect(canEditStoreRecord('FRANCHISE', 'other', ['store:update'], row)).toBe(false);
      expect(canSubmitFranchiseStore(lifecycle, ['store:submit'])).toBe(true);
    }
    expect(canEditStoreRecord('FRANCHISE', 'tenant', ['store:update'], { organizationId: 'tenant', lifecycle: 'PENDING_APPROVAL' })).toBe(false);
    expect(canSubmitFranchiseStore('PENDING_APPROVAL', ['store:submit'])).toBe(false);
    expect(canSubmitFranchiseStore('DRAFT', [])).toBe(false);
  });
});
