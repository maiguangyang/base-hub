import { useApolloClient, useQuery } from '@apollo/client/react';
import { useNavigate } from 'react-router';
import { Button } from '@/components/ui/button';
import { ADMIN_STORE_EDITOR_QUERY } from '@/features/admin/graphql/storeEditor';
import { useAdminTab } from '@/features/admin/hooks/useAdminTab';
import { useAuthStore } from '@/features/auth/store/authStore';
import { getGraphQLErrorCode } from '@/lib/graphql/errors';
import { StoreEditorForm } from './StoreEditorForm';
import { canCreateStore, canEditStoreRecord, type StoreWorkspace } from './storeEditorPolicy';

export function StoreEditorPage({ workspace }: { workspace: StoreWorkspace }) {
  const viewer = useAuthStore((state) => state.viewer);
  const { search } = useAdminTab();
  const id = new URLSearchParams(search).get('id');
  const permissions = viewer?.permissions ?? [];
  const canWrite = hasStoreWritePermission(workspace, id, permissions);
  if (viewer?.currentWorkspace?.workspaceType !== workspace || !viewer.currentWorkspace.organizationId || !canWrite) {
    return <p role="alert">当前工作区无权维护门店。</p>;
  }
  if (id === '') return <p role="alert">门店链接无效，请从列表重新进入。</p>;
  return <StoreEditorWorkspace workspace={workspace} organizationId={viewer.currentWorkspace.organizationId}
    permissions={permissions} id={id} />;
}

function hasStoreWritePermission(workspace: StoreWorkspace, id: string | null, permissions: readonly string[]) {
  if (id === null) return canCreateStore(workspace, permissions);
  return permissions.includes(workspace === 'HEADQUARTERS' ? 'hqStore:update' : 'store:update');
}

interface WorkspaceProps {
  workspace: StoreWorkspace;
  organizationId: string;
  permissions: readonly string[];
  id: string | null;
}

// eslint-disable-next-line complexity -- Distinct loading, permission, query-error and missing-record states are intentional.
function StoreEditorWorkspace({ workspace, organizationId, permissions, id }: WorkspaceProps) {
  const navigate = useNavigate();
  const client = useApolloClient();
  const listPath = workspace === 'HEADQUARTERS' ? '/admin/hq/stores' : '/admin/franchise/stores';
  const { data, loading, error, refetch } = useQuery(ADMIN_STORE_EDITOR_QUERY, {
    variables: { id: id ?? '' }, skip: !id, fetchPolicy: 'network-only',
  });
  const existing = id && data?.store?.id === id ? data.store : undefined;
  const back = () => navigate(listPath);
  if (id && loading) return <p role="status">正在加载门店…</p>;
  if (id && error) {
    const code = getGraphQLErrorCode(error);
    if (code === 'PERMISSION_DENIED' || code === 'STORE_SCOPE_DENIED') {
      return <EditorMessage message="门店不存在或无权访问" onBack={back} />;
    }
    return <div className="flex flex-col gap-3"><p role="alert">门店加载失败，请重试。</p>
      <Button type="button" variant="outline" onClick={() => void refetch().catch(() => undefined)}>重试</Button></div>;
  }
  if (id && !existing) return <EditorMessage message="门店不存在或无权访问" onBack={back} />;
  if (existing && !canEditStoreRecord(workspace, organizationId, permissions, existing)) {
    return <EditorMessage message="当前门店不可编辑" onBack={back} />;
  }
  const saved = () => {
    client.cache.evict({ fieldName: 'stores' });
    client.cache.evict({ fieldName: 'store' });
    client.cache.gc();
    back();
  };
  return <StoreEditorForm key={`${workspace}:${id ?? 'new'}`} workspace={workspace} organizationId={organizationId}
    existing={existing} onBack={back} onSaved={saved} />;
}

function EditorMessage({ message, onBack }: { message: string; onBack(): void }) {
  return <div className="flex flex-col gap-3"><p role="alert">{message}</p>
    <Button type="button" variant="outline" onClick={onBack}>返回门店</Button></div>;
}
