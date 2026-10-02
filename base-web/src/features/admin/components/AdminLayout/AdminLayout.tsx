import { useEffect, useState } from 'react';
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar';
import { AdminHeader } from '@/features/admin/components/AdminHeader';
import { AdminSidebar } from '@/features/admin/components/AdminSidebar';
import { AiDrawer } from '@/features/admin/components/AiDrawer/AiDrawer';
import { useAiSession } from '@/features/admin/components/AiDrawer/useAiSession';
import { AdminTabBar } from '@/features/admin/components/AdminTabBar';
import { TabStack } from '@/features/admin/components/TabStack';
import { AdminNotFoundPage } from '@/features/admin/pages/AdminNotFoundPage';
import { useSyncTabWithLocation } from '@/features/admin/routes/useSyncTabWithLocation';
import type { BusinessWorkspaceType } from '@/features/admin/config/navigation';
import { authMessages } from '@/features/auth/pages/authMessages';
import { useAuthStore } from '@/features/auth/store/authStore';
import { bindAdminWorkspace, useAdminTabsStore, workspaceNamespace } from '@/stores/slices/adminTabsSlice';

/** 后台持久外壳。只做组装，各区块由子组件负责，保持本文件薄。 */
export function AdminLayout() {
	const [aiOpen, setAiOpen] = useState(false);
  const viewer = useAuthStore((state) => state.viewer);
  const namespace = useAdminTabsStore((state) => state.namespace);
  const workspace = viewer?.currentWorkspace;
  const workspaceType = workspace?.workspaceType as BusinessWorkspaceType;
  const expectedNamespace = workspace && viewer
    ? workspaceNamespace(viewer.account.id, workspaceType, workspace.organizationId)
    : '';
	const aiSession = useAiSession(expectedNamespace);
  useEffect(() => {
    if (!workspace || !viewer) return;
    void bindAdminWorkspace({
      accountId: viewer.account.id, workspaceType,
      organizationId: workspace.organizationId, permissions: viewer.permissions,
    });
  }, [viewer, workspace, workspaceType]);
  const { isKnownPath } = useSyncTabWithLocation();
  if (!workspace || namespace !== expectedNamespace) return <div className="p-6 text-sm text-muted-foreground">{authMessages.workspace.loading}</div>;

  return (
    <SidebarProvider>
      <AdminSidebar />
      <SidebarInset className="flex h-screen min-w-0 flex-col overflow-hidden">
		<AdminHeader onOpenAi={() => setAiOpen(true)} />
        {isKnownPath ? (
          <>
            <AdminTabBar />
            <TabStack />
          </>
        ) : (
          <AdminNotFoundPage />
        )}
      </SidebarInset>
		<AiDrawer open={aiOpen} onOpenChange={setAiOpen} session={aiSession} />
    </SidebarProvider>
  );
}
