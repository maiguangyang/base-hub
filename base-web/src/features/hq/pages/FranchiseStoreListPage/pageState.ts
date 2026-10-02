export function selectedFranchiseId(search: string): string | undefined {
  const value = new URLSearchParams(search).get('organizationId')?.trim();
  return value || undefined;
}

/** 同一路径切换加盟商时强制重建页面私有状态，避免沿用上一组织的分页和搜索。 */
export function franchiseStoreViewKey(search: string): string {
  return selectedFranchiseId(search) ?? 'all-franchises';
}

export function initialFranchiseStoreKeyword(search: string): string {
  return new URLSearchParams(search).get('q') ?? '';
}
