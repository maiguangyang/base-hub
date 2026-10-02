package tools

import (
	"base-engine/auth"
	"base-engine/src/services/ai"
)

func task13Specs() []ai.ToolSpec {
	return []ai.ToolSpec{
		reviewedSpec("FranchiseAuditLogs", "graphql.query.auditLogs", `query FranchiseAuditLogs($page: Int!, $pageSize: Int!, $q: String, $filter: AuditLogFilterType) {
    auditLogs(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data { id action resourceType resourceId resultCode actorAccountId storeId createdAt }
      total current_page per_page total_page
    }
  }`, ai.ModeReadOnly, "LOW", "tenantAudit:read", auth.WorkspaceTypeFranchise, "", nil, []string{"page", "pageSize", "q", "filter"}),
		reviewedSpec("FranchiseCreateStore", "graphql.mutation.createStore", `mutation FranchiseCreateStore($input: CreateStoreInput!) {
    createStore(input: $input) { id code name lifecycle organizationId }
  }`, ai.ModeWrite, "MEDIUM", "store:create", auth.WorkspaceTypeFranchise, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("FranchiseDeleteStores", "graphql.mutation.deleteStores", `mutation FranchiseDeleteStores($ids: [ID!]!) { deleteStores(id: $ids) }`, ai.ModeWrite, "HIGH", "store:delete", auth.WorkspaceTypeFranchise, ai.WriteExisting, []string{"ids"}, []string{"ids"}),
		reviewedSpec("FranchiseStores", "graphql.query.stores", `query FranchiseStores($page: Int!, $pageSize: Int!, $q: String, $filter: StoreFilterType) {
    stores(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data { id code name lifecycle rejectionReason organizationId }
      total current_page per_page total_page
    }
  }`, ai.ModeReadOnly, "LOW", "store:read", auth.WorkspaceTypeFranchise, "", nil, []string{"page", "pageSize", "q", "filter"}),
		reviewedSpec("FranchiseSubmitStore", "graphql.mutation.submitStore", `mutation FranchiseSubmitStore($id: ID!) {
    submitStore(id: $id) { id lifecycle submittedAt }
  }`, ai.ModeWrite, "MEDIUM", "store:submit", auth.WorkspaceTypeFranchise, ai.WriteExisting, []string{"id"}, []string{"id"}),
		reviewedSpec("FranchiseUpdateStore", "graphql.mutation.updateStore", `mutation FranchiseUpdateStore($id: ID!, $input: UpdateStoreInput!) {
    updateStore(id: $id, input: $input) { id code name lifecycle organizationId }
  }`, ai.ModeWrite, "MEDIUM", "store:update", auth.WorkspaceTypeFranchise, ai.WriteExisting, []string{"id"}, []string{"id", "input"}),
	}
}
