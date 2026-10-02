export type StoreWorkspace = 'HEADQUARTERS' | 'FRANCHISE';

export function canCreateStore(workspace: StoreWorkspace, permissions: readonly string[]): boolean {
  return permissions.includes(workspace === 'HEADQUARTERS' ? 'hqStore:create' : 'store:create');
}

export function canEditStoreRecord(workspace: StoreWorkspace, organizationId: string, permissions: readonly string[], record: {
  organizationId: string;
  lifecycle: string;
  organization?: { type: string } | null;
}): boolean {
  if (!organizationId || record.organizationId !== organizationId) return false;
  if (workspace === 'HEADQUARTERS') return permissions.includes('hqStore:update') && record.organization?.type === 'HEADQUARTERS';
  return permissions.includes('store:update') && (record.lifecycle === 'DRAFT' || record.lifecycle === 'REJECTED');
}

export function canSubmitFranchiseStore(lifecycle: string, permissions: readonly string[]): boolean {
  return permissions.includes('store:submit') && (lifecycle === 'DRAFT' || lifecycle === 'REJECTED');
}
