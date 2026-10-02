import { useEffect } from 'react';
import { ApolloProvider } from '@apollo/client/react';
import { useQuery } from '@apollo/client/react';
import { BrowserRouter, Navigate, Route, Routes, useLocation } from 'react-router';
import { useFragment } from '@/__generated__';
import { AdminLayout } from '@/features/admin/components/AdminLayout';
import { AdminToastProvider } from '@/features/admin/components/AdminToast';
import { canonicalHqStoreLocation, type BusinessWorkspaceType } from '@/features/admin/config/navigation';
import { VIEWER_FIELDS_FRAGMENT, VIEWER_QUERY } from '@/features/auth/graphql/auth';
import { LoginPage } from '@/features/auth/pages/LoginPage';
import { PasswordChangePage } from '@/features/auth/pages/PasswordChangePage';
import { WorkspaceSelectPage } from '@/features/auth/pages/WorkspaceSelectPage';
import { useAuthStore, type AuthPhase, type ViewerSummary } from '@/features/auth/store/authStore';
import { useSessionEvents } from '@/features/auth/hooks/useSessionEvents';
import { authPath, currentDestination, loginPath, nextAuthPath } from '@/features/auth/pages/authFlow';
import { SystemInitializationGate } from '@/features/auth/components/SystemInitializationGate';
import { getGraphQLRuntime, type GraphQLRuntime } from '@/lib/graphql/client';
import { resetAdminWorkspace } from '@/stores/slices/adminTabsSlice';

/** 后台 SPA 根组件。路由为绝对路径（/admin/...），不需要 basename。 */
export function AdminApp() {
  const resetAuth = useAuthStore((state) => state.reset);
  const runtime = getGraphQLRuntime({ resetAuth, resetWorkspace: resetAdminWorkspace });
  useEffect(() => {
    // client:only 下 Astro 不输出 HTML，首屏骨架由 AdminBaseLayout 提供，挂载后移除。
    document.getElementById('admin-boot-skeleton')?.remove();
    // 慢速加载时骨架的超时兜底会先插入一段提示，此刻已挂载成功，该提示是错的，一并清除。
    document.getElementById('admin-boot-notice')?.remove();
  }, []);

  return (
    <ApolloProvider client={runtime.client}>
      <AdminToastProvider>
        <BrowserRouter>
          <SystemInitializationGate>
            <AdminRoutes runtime={runtime} />
          </SystemInitializationGate>
        </BrowserRouter>
      </AdminToastProvider>
    </ApolloProvider>
  );
}

function AdminRoutes({ runtime }: { runtime: GraphQLRuntime }) {
  const phase = useAuthStore((state) => state.phase);
  const setViewer = useAuthStore((state) => state.setViewer);
  const reset = useAuthStore((state) => state.reset);
  const { data, error } = useQuery(VIEWER_QUERY, { fetchPolicy: 'network-only' });
  const viewer = useFragment(VIEWER_FIELDS_FRAGMENT, data?.viewer);
  const activeViewer = useAuthStore((state) => state.viewer);
  useSessionEvents(runtime, activeViewer);
  useEffect(() => {
    if (viewer) setViewer(viewer);
    else if (error) reset();
  }, [error, reset, setViewer, viewer]);
  if (phase === 'loading') return <div className="p-6 text-sm text-muted-foreground">正在验证登录状态…</div>;
  return (
    <Routes>
      <Route path="/admin/login" element={<LoginPage />} />
      <Route path="/admin/change-password" element={<PasswordChangeRoute />} />
      <Route path="/admin/workspaces" element={<WorkspaceSelectPage />} />
      <Route path="*" element={<WorkspaceRoute />} />
    </Routes>
  );
}

function PasswordChangeRoute() {
  const { phase, viewer } = useAuthStore();
  const redirect = passwordChangeRedirect(phase, viewer);
  if (redirect) return <Navigate to={redirect} replace />;
  return <PasswordChangePage />;
}

/** 独立改密页只服务首次登录临时密码流程；普通用户必须使用后台壳内弹窗。 */
export function passwordChangeRedirect(phase: AuthPhase, viewer: ViewerSummary | null): string | undefined {
  if (phase === 'anonymous' || !viewer) return '/admin/login';
  if (viewer.account.mustChangePassword) return undefined;
  return nextAuthPath(viewer);
}

function WorkspaceRoute() {
  const location = useLocation();
  const { phase, viewer } = useAuthStore();
  const redirect = workspaceRedirect(phase, viewer, location);
  if (redirect) return <Navigate to={redirect} replace />;
  return <AdminLayout />;
}

type LocationParts = Pick<Location, 'pathname' | 'search' | 'hash'>;

/** 根据认证阶段和当前工作台生成安全跳转，并保留合法业务深链。 */
export function workspaceRedirect(phase: AuthPhase, viewer: ViewerSummary | null, location: LocationParts): string | undefined {
  const destination = currentDestination(location);
  if (phase === 'anonymous' || !viewer) return loginPath(destination);
  if (phase === 'password-change') return authPath('/admin/change-password', destination);
  if (phase === 'workspace-select' || !viewer.currentWorkspace) return authPath('/admin/workspaces', destination);
  const workspaceType = viewer.currentWorkspace.workspaceType as BusinessWorkspaceType;
  return workspacePathRedirect(workspaceType, location, destination);
}

function workspacePathRedirect(workspaceType: BusinessWorkspaceType, location: LocationParts, destination: string): string | undefined {
  const storeRedirect = legacyHqStoreRedirect(workspaceType, location);
  if (storeRedirect) return storeRedirect;
  const expectedPrefix = workspaceType === 'HEADQUARTERS' ? '/admin/hq' : '/admin/franchise';
  const otherPrefix = workspaceType === 'HEADQUARTERS' ? '/admin/franchise' : '/admin/hq';
  if (matchesWorkspacePath(location.pathname, otherPrefix) || !matchesWorkspacePath(location.pathname, expectedPrefix)) {
    return authPath('/admin/workspaces', destination);
  }
  return undefined;
}

function legacyHqStoreRedirect(workspaceType: BusinessWorkspaceType, location: LocationParts): string | undefined {
  if (workspaceType !== 'HEADQUARTERS') return undefined;
  const canonical = canonicalHqStoreLocation(location.pathname, location.search);
  if (canonical.path === location.pathname) return undefined;
  return `${canonical.path}${canonical.search}${location.hash}`;
}

function matchesWorkspacePath(pathname: string, prefix: string): boolean {
	return pathname === prefix || pathname.startsWith(`${prefix}/`);
}
