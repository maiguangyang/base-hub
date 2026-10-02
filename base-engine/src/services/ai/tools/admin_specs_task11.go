package tools

import (
	"base-engine/auth"
	"base-engine/src/services/ai"
)

func task11Specs() []ai.ToolSpec {
	return []ai.ToolSpec{
		reviewedSpec("FranchiseChangeMembershipStatus", "graphql.mutation.changeMembershipStatus", `mutation FranchiseChangeMembershipStatus($input: ChangeMembershipStatusInput!) {
    changeMembershipStatus(input: $input) { id status }
  }`, ai.ModeWrite, "HIGH", "operatorMembership:update", auth.WorkspaceTypeFranchise, ai.WriteExisting, []string{"input.membershipId"}, []string{"input"}),
		reviewedSpec("FranchiseInviteStaff", "graphql.mutation.inviteOperator", `mutation FranchiseInviteStaff($input: InviteOperatorInput!) {
    inviteOperator(input: $input) {
      membership { id status accountId organizationId }
      temporaryPassword invitationPending
    }
  }`, ai.ModeWrite, "HIGH", "operatorMembership:create", auth.WorkspaceTypeFranchise, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("FranchiseStaff", "graphql.query.operatorMemberships", `query FranchiseStaff($page: Int!, $pageSize: Int!, $q: String, $filter: OperatorMembershipFilterType) {
    operatorMemberships(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data {
        id status storeAccessMode accountId organizationId
        account { id phone displayName email }
        roles { id name kind }
        stores { id name lifecycle }
      }
      total current_page per_page total_page
    }
  }`, ai.ModeReadOnly, "LOW", "operatorMembership:read", auth.WorkspaceTypeFranchise, "", nil, []string{"page", "pageSize", "q", "filter"}),
		reviewedSpec("FranchiseUpdateStaff", "graphql.mutation.updateOperatorMembership", `mutation FranchiseUpdateStaff($id: ID!, $input: UpdateOperatorMembershipInput!) {
    updateOperatorMembership(id: $id, input: $input) { id status storeAccessMode rolesIds storesIds }
  }`, ai.ModeWrite, "MEDIUM", "operatorMembership:update", auth.WorkspaceTypeFranchise, ai.WriteExisting, []string{"id"}, []string{"id", "input"}),
	}
}
