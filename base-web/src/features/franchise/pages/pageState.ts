import type { ViewerSummary } from '@/features/auth/store/authStore';

export function initialSearchKeyword(search: string): string {
  return new URLSearchParams(search).get('q') ?? '';
}

export function queryKeyword(keyword: string): string | null {
  return keyword === '' ? null : keyword;
}

export function viewerPermissions(viewer: ViewerSummary | null): string[] {
  return viewer ? viewer.permissions : [];
}

export function viewerOrganizationId(viewer: ViewerSummary | null): string {
  return viewer?.currentWorkspace?.organizationId ?? '';
}
