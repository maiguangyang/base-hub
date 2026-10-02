import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useMutation, useQuery } from '@apollo/client/react';
import { useFragment } from '@/__generated__';
import type { PendingMembershipInvitationsQuery } from '@/__generated__/graphql';
import { Button } from '@/components/ui/button';
import { AuthPageShell } from '@/features/auth/components/AuthPageShell';
import { ACCEPT_MEMBERSHIP_INVITATION_MUTATION, PENDING_MEMBERSHIP_INVITATIONS_QUERY, SELECT_WORKSPACE_MUTATION, VIEWER_FIELDS_FRAGMENT, WORKSPACES_QUERY } from '@/features/auth/graphql/auth';
import { useAuthStore } from '@/features/auth/store/authStore';
import { getGraphQLErrorCode } from '@/lib/graphql/errors';
import { disposeGraphQLRuntime } from '@/lib/graphql/client';
import { bindAdminWorkspace } from '@/stores/slices/adminTabsSlice';
import { browserReturnDestination, businessWorkspaces, nextAuthPath, type WorkspaceSummary } from '../authFlow';
import { authErrorMessage, authMessages } from '../authMessages';

/** 账号工作台选择页；只有一个业务工作台时自动进入。 */
export function WorkspaceSelectPage() {
  const setViewer = useAuthStore((state) => state.setViewer);
  const { data, loading: querying, refetch: refetchWorkspaces } = useQuery(WORKSPACES_QUERY);
  const invitations = useQuery(PENDING_MEMBERSHIP_INVITATIONS_QUERY);
  const [selectWorkspace, { loading: selecting }] = useMutation(SELECT_WORKSPACE_MUTATION); const [acceptInvitation, { loading: accepting }] = useMutation(ACCEPT_MEMBERSHIP_INVITATION_MUTATION);
  const [error, setError] = useState<string>();
  const workspaces = useMemo(() => businessWorkspaces(data?.workspaces ?? []), [data?.workspaces]);
  const choose = useCallback(async (workspace: WorkspaceSummary) => {
    setError(undefined);
    try {
      const result = await selectWorkspace({ variables: { input: { workspaceType: workspace.workspaceType, organizationId: workspace.organizationId } } });
      const viewer = useFragment(VIEWER_FIELDS_FRAGMENT, result.data?.selectWorkspace);
      if (!viewer) throw new Error('AUTH_REQUIRED');
      setViewer(viewer);
      const nextPath = nextAuthPath(viewer, browserReturnDestination());
      const selected = viewer.currentWorkspace;
      await disposeGraphQLRuntime();
      if (selected && selected.workspaceType !== 'DISCOVERY') {
        await bindAdminWorkspace({
          accountId: viewer.account.id, workspaceType: selected.workspaceType,
          organizationId: selected.organizationId, permissions: viewer.permissions,
        });
      }
      window.location.replace(nextPath);
    } catch (cause) {
      setError(authErrorMessage(getGraphQLErrorCode(cause)));
    }
  }, [selectWorkspace, setViewer]);
  const accept = useCallback(async (id: string) => {
    setError(undefined);
    try {
      await acceptInvitation({ variables: { id } });
      await Promise.all([invitations.refetch(), refetchWorkspaces()]);
    } catch (cause) {
      setError(authErrorMessage(getGraphQLErrorCode(cause)));
    }
  }, [acceptInvitation, invitations, refetchWorkspaces]);

  useAutomaticWorkspaceChoice(workspaces, querying, selecting, choose);

  return (
    <AuthPageShell title={authMessages.workspace.title} description={authMessages.workspace.description}>
      <div className="space-y-3">
        {querying && <p className="text-sm text-muted-foreground">{authMessages.common.loading}</p>}
        {!querying && workspaces.length === 0 && <p className="text-sm text-muted-foreground">{authMessages.workspace.empty}</p>}
        <WorkspaceChoices workspaces={workspaces} disabled={selecting} onChoose={choose} />
        <InvitationList invitations={invitations.data?.pendingMembershipInvitations ?? []} disabled={accepting} onAccept={accept} />
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
      </div>
    </AuthPageShell>
  );
}

function useAutomaticWorkspaceChoice(workspaces: WorkspaceSummary[], querying: boolean, selecting: boolean, choose: (workspace: WorkspaceSummary) => Promise<void>) {
  const attempted = useRef<string | undefined>(undefined);
  useEffect(() => {
    if (querying || selecting || workspaces.length !== 1) return;
    const workspace = workspaces[0];
    const key = `${workspace.workspaceType}:${workspace.organizationId ?? 'none'}`;
    if (attempted.current === key) return;
    attempted.current = key;
    void choose(workspace);
  }, [choose, querying, selecting, workspaces]);
}

function WorkspaceChoices({ workspaces, disabled, onChoose }: { workspaces: WorkspaceSummary[]; disabled: boolean; onChoose(workspace: WorkspaceSummary): Promise<void> }) {
  return workspaces.map((workspace) => (
    <Button key={`${workspace.workspaceType}:${workspace.organizationId ?? 'none'}`} variant="outline" className="h-auto w-full justify-start px-4 py-3" disabled={disabled} onClick={() => void onChoose(workspace)}>
      <span className="text-left"><strong className="block">{workspace.workspaceType === 'HEADQUARTERS' ? authMessages.workspace.hq : workspace.organizationName}</strong><span className="text-xs text-muted-foreground">{workspace.workspaceType === 'HEADQUARTERS' ? authMessages.workspace.hq : authMessages.workspace.franchise}</span></span>
    </Button>
  ));
}

type Invitation = PendingMembershipInvitationsQuery['pendingMembershipInvitations'][number];

function InvitationList({ invitations, disabled, onAccept }: { invitations: Invitation[]; disabled: boolean; onAccept(id: string): Promise<void> }) {
  if (invitations.length === 0) return null;
  return <section className="space-y-2 border-t pt-4"><h2 className="font-medium">{authMessages.workspace.invitations}</h2><p className="text-xs text-muted-foreground">{authMessages.workspace.invitationDescription}</p>{invitations.map((invitation) => <div key={invitation.id} className="flex items-center justify-between rounded-md border p-3"><span className="text-sm">{authMessages.workspace.invitationPending}</span><Button type="button" disabled={disabled} onClick={() => void onAccept(invitation.id)}>{authMessages.workspace.acceptInvitation}</Button></div>)}</section>;
}
