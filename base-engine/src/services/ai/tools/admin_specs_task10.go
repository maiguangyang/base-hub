package tools

import (
	"base-engine/auth"
	"base-engine/src/services/ai"
)

func task10Specs() []ai.ToolSpec {
	return []ai.ToolSpec{
		reviewedSpec("HqAuditLogs", "graphql.query.auditLogs", `query HqAuditLogs($page: Int!, $pageSize: Int!, $q: String, $filter: AuditLogFilterType) {
    auditLogs(current_page: $page, per_page: $pageSize, q: $q, filter: $filter) {
      data { id action resourceType resourceId resultCode actorAccountId organizationId storeId createdAt }
      total current_page per_page total_page
    }
  }`, ai.ModeReadOnly, "LOW", "auditLog:read", auth.WorkspaceTypeHeadquarters, "", nil, []string{"page", "pageSize", "q", "filter"}),
		reviewedSpec("HqCreateDirectStore", "graphql.mutation.createStore", `mutation HqCreateDirectStore($input: CreateStoreInput!) {
    createStore(input: $input) {
      id code name lifecycle organizationId
      contactPhone managerName managerPhone province city district address
      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter
    }
  }`, ai.ModeWrite, "MEDIUM", "hqStore:create", auth.WorkspaceTypeHeadquarters, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("HqDeleteDirectStores", "graphql.mutation.deleteStores", `mutation HqDeleteDirectStores($ids: [ID!]!) { deleteStores(id: $ids) }`, ai.ModeWrite, "HIGH", "hqStore:delete", auth.WorkspaceTypeHeadquarters, ai.WriteExisting, []string{"ids"}, []string{"ids"}),
		reviewedSpec("HqDirectStores", "graphql.query.stores", `query HqDirectStores($page:Int!,$pageSize:Int!,$q:String){stores(current_page:$page,per_page:$pageSize,q:$q,filter:{organization:{type:HEADQUARTERS}}){data{id code name lifecycle organizationId} total current_page per_page total_page}}`, ai.ModeReadOnly, "LOW", "hqStore:read", auth.WorkspaceTypeHeadquarters, "", nil, []string{"page", "pageSize", "q"}),
		reviewedSpec("HqFranchiseStores", "graphql.query.stores", `query HqFranchiseStores($page:Int!,$pageSize:Int!,$q:String){stores(current_page:$page,per_page:$pageSize,q:$q,filter:{organization:{type:FRANCHISE}}){data{id code name lifecycle organizationId} total current_page per_page total_page}}`, ai.ModeReadOnly, "LOW", "store:read_all", auth.WorkspaceTypeHeadquarters, "", nil, []string{"page", "pageSize", "q"}),
		reviewedSpec("HqApproveStore", "graphql.mutation.reviewStore", `mutation HqApproveStore($storeId:ID!){reviewStore(input:{storeId:$storeId,approved:true}){id lifecycle rejectionReason reviewedAt}}`, ai.ModeWrite, "HIGH", "store:approve", auth.WorkspaceTypeHeadquarters, ai.WriteExisting, []string{"storeId"}, []string{"storeId"}),
		reviewedSpec("HqRejectStore", "graphql.mutation.reviewStore", `mutation HqRejectStore($storeId:ID!,$rejectionReason:String!){reviewStore(input:{storeId:$storeId,approved:false,rejectionReason:$rejectionReason}){id lifecycle rejectionReason reviewedAt}}`, ai.ModeWrite, "HIGH", "store:reject", auth.WorkspaceTypeHeadquarters, ai.WriteExisting, []string{"storeId"}, []string{"storeId", "rejectionReason"}),
		reviewedSpec("HqStoreApprovals", "graphql.query.stores", `query HqStoreApprovals($page:Int!,$pageSize:Int!,$q:String){stores(current_page:$page,per_page:$pageSize,q:$q,filter:{lifecycle:PENDING_APPROVAL}){data{id code name lifecycle submittedAt organizationId organization{id name}} total current_page per_page total_page}}`, ai.ModeReadOnly, "LOW", "store:read_all", auth.WorkspaceTypeHeadquarters, "", nil, []string{"page", "pageSize", "q"}),
		reviewedSpec("HqUpdateDirectStore", "graphql.mutation.updateStore", `mutation HqUpdateDirectStore($id: ID!, $input: UpdateStoreInput!) {
    updateStore(id: $id, input: $input) {
      id code name lifecycle organizationId
      contactPhone managerName managerPhone province city district address
      businessHours businessStatus supportDineIn supportTakeout storeArea tableCount receiptFooter
    }
  }`, ai.ModeWrite, "MEDIUM", "hqStore:update", auth.WorkspaceTypeHeadquarters, ai.WriteExisting, []string{"id"}, []string{"id", "input"}),
	}
}
