package tools

import (
	"base-engine/auth"
	"base-engine/src/services/ai"
)

func task8Specs() []ai.ToolSpec {
	return []ai.ToolSpec{
		reviewedSpec("HqAdministratorIdentity", "graphql.query.operatorMemberships", `query HqAdministratorIdentity($accountId: ID!) {
    operatorMemberships(current_page: 1, per_page: 1, filter: { accountId: $accountId }) {
      data { id accountId roles { id kind } }
      total
    }
  }`, ai.ModeReadOnly, "LOW", "hqMembership:read", auth.WorkspaceTypeHeadquarters, "", nil, []string{"accountId"}),
		reviewedSpec("HqAdministrators", "graphql.query.operatorMemberships", `query HqAdministrators($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorMembershipFilterType) {
    operatorMemberships(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data {
        id status storeAccessMode accountId organizationId updatedAt
        account { id phone displayName email }
        roles { id name kind }
      }
      total current_page per_page total_page
    }
  }`, ai.ModeReadOnly, "LOW", "hqMembership:read", auth.WorkspaceTypeHeadquarters, "", nil, []string{"page", "pageSize", "q", "filter"}),
		reviewedSpec("HqChangeAdministratorStatus", "graphql.mutation.changeMembershipStatus", `mutation HqChangeAdministratorStatus($input: ChangeMembershipStatusInput!) {
    changeMembershipStatus(input: $input) { id status }
  }`, ai.ModeWrite, "HIGH", "hqMembership:update", auth.WorkspaceTypeHeadquarters, ai.WriteExisting, []string{"input.membershipId"}, []string{"input"}),
		reviewedSpec("HqDeleteAdministrators", "graphql.mutation.deleteOperatorMemberships", `mutation HqDeleteAdministrators($ids: [ID!]!) { deleteOperatorMemberships(id: $ids) }`, ai.ModeWrite, "HIGH", "hqMembership:delete", auth.WorkspaceTypeHeadquarters, ai.WriteExisting, []string{"ids"}, []string{"ids"}),
		reviewedSpec("HqInviteAdministrator", "graphql.mutation.inviteOperator", `mutation HqInviteAdministrator($input: InviteOperatorInput!) {
    inviteOperator(input: $input) {
      membership { id status accountId organizationId }
      temporaryPassword invitationPending
    }
  }`, ai.ModeWrite, "HIGH", "hqMembership:create", auth.WorkspaceTypeHeadquarters, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("HqResetAdministratorPassword", "graphql.mutation.resetTemporaryPassword", `mutation HqResetAdministratorPassword($accountId: ID!) {
    resetTemporaryPassword(accountId: $accountId) { accountId temporaryPassword }
  }`, ai.ModeWrite, "HIGH", "account:update", auth.WorkspaceTypeHeadquarters, ai.WriteExisting, []string{"accountId"}, []string{"accountId"}),
		reviewedSpec("HqUpdateAdministrator", "graphql.mutation.updateOperatorMembership", `mutation HqUpdateAdministrator($id: ID!, $input: UpdateOperatorMembershipInput!) {
    updateOperatorMembership(id: $id, input: $input) { id status rolesIds }
  }`, ai.ModeWrite, "MEDIUM", "hqMembership:update", auth.WorkspaceTypeHeadquarters, ai.WriteExisting, []string{"id"}, []string{"id", "input"}),
	}
}
