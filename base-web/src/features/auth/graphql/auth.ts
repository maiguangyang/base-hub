import { gql } from '@/__generated__';

export const VIEWER_FIELDS_FRAGMENT = gql(`
  fragment ViewerFields on Viewer {
    account { id phone displayName email status mustChangePassword }
    currentWorkspace { workspaceType organizationId organizationName homePath }
    workspaces { workspaceType organizationId organizationName homePath }
    permissions
  }
`);

export const LOGIN_MUTATION = gql(`
  mutation Login($input: LoginInput!) {
    login(input: $input) {
      requiresPasswordChange
      viewer { ...ViewerFields }
    }
  }
`);

export const VIEWER_QUERY = gql(`
  query Viewer {
    viewer { ...ViewerFields }
  }
`);

export const WORKSPACES_QUERY = gql(`
  query Workspaces {
    workspaces { workspaceType organizationId organizationName homePath }
  }
`);

export const CHANGE_TEMPORARY_PASSWORD_MUTATION = gql(`
  mutation ChangeTemporaryPassword($input: ChangePasswordInput!) {
    changeTemporaryPassword(input: $input) { ...ViewerFields }
  }
`);

export const SELECT_WORKSPACE_MUTATION = gql(`
  mutation SelectWorkspace($input: SelectWorkspaceInput!) {
    selectWorkspace(input: $input) { ...ViewerFields }
  }
`);

export const LOGOUT_MUTATION = gql(`
  mutation Logout { logout }
`);

export const SESSION_EVENTS_SUBSCRIPTION = gql(`
  subscription SessionEvents {
    sessionEvents { code sessionId organizationId occurredAt }
  }
`);

export const PENDING_MEMBERSHIP_INVITATIONS_QUERY = gql(`
  query PendingMembershipInvitations {
    pendingMembershipInvitations { id membershipId expiresAt }
  }
`);

export const ACCEPT_MEMBERSHIP_INVITATION_MUTATION = gql(`
  mutation AcceptMembershipInvitation($id: ID!) {
    acceptMembershipInvitation(id: $id) { id status organizationId }
  }
`);
