package tools

import (
	"base-engine/auth"
	"base-engine/src/services/ai"
)

func hqFranchiseSpecs() []ai.ToolSpec {
	return []ai.ToolSpec{
		hqSpec("HqFranchises", "graphql.query.organizations", `query HqFranchises($page:Int!,$pageSize:Int!,$q:String,$status:OrganizationStatus){organizations(current_page:$page,per_page:$pageSize,q:$q,filter:{type:FRANCHISE,status:$status}){data{id code name status} total current_page per_page total_page}}`, ai.ModeReadOnly, "LOW", "organization:read", []string{"organizations"}),
		hqSpec("HqFranchiseInitialAccountCandidates", "graphql.query.organization", `query HqFranchiseInitialAccountCandidates($id:ID!){organization(id:$id,filter:{type:FRANCHISE}){id memberships{id status account{id phone displayName status}}}}`, ai.ModeReadOnly, "LOW", "organization:read", []string{"organization"}),
		franchiseWriteSpec("HqProvisionFranchise", "graphql.mutation.provisionFranchise", `mutation HqProvisionFranchise($input:ProvisionFranchiseInput!){provisionFranchise(input:$input){organization{id code name status} membership{id status account{id}} temporaryPassword invitationPending}}`, "franchise:provision", ai.WriteCreate, nil, []string{"input"}, []string{"organizationId", "membershipId", "invitationPending"}),
		franchiseWriteSpec("HqResetFranchiseInitialPassword", "graphql.mutation.resetFranchiseInitialPassword", `mutation HqResetFranchiseInitialPassword($organizationId:ID!){resetFranchiseInitialPassword(organizationId:$organizationId){accountId temporaryPassword}}`, "account:update", ai.WriteExisting, []string{"organizationId"}, []string{"organizationId"}, []string{"accountId"}),
		franchiseWriteSpec("HqSuspendOrganization", "graphql.mutation.suspendOrganization", `mutation HqSuspendOrganization($input:SuspendOrganizationInput!){suspendOrganization(input:$input){id status suspensionReasonCode}}`, "organization:suspend", ai.WriteExisting, []string{"input.organizationId"}, []string{"input"}, []string{"id", "status", "suspensionReasonCode"}),
		franchiseWriteSpec("HqRestoreOrganization", "graphql.mutation.restoreOrganization", `mutation HqRestoreOrganization($id:ID!){restoreOrganization(id:$id){id status}}`, "organization:restore", ai.WriteExisting, []string{"id"}, []string{"id"}, []string{"id", "status"}),
		initialAccountSpec(),
	}
}

func hqSpec(id, operation, document string, mode ai.ToolMode, risk, permission string, output []string) ai.ToolSpec {
	title, description := toolMetadata(id, mode, auth.WorkspaceTypeHeadquarters)
	return ai.ToolSpec{ID: id, Name: id, Title: title, OperationID: operation, Document: document, Description: description,
		Mode: mode, Risk: risk, Permission: permission, Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeHeadquarters}, OutputFields: output}
}

func franchiseWriteSpec(id, operation, document, permission string, kind ai.WriteKind, targets, arguments, output []string) ai.ToolSpec {
	spec := hqSpec(id, operation, document, ai.ModeWrite, "HIGH", permission, output)
	spec.WriteKind, spec.TargetFields, spec.ArgumentFields = kind, targets, arguments
	return spec
}

func initialAccountSpec() ai.ToolSpec {
	spec := hqSpec("HqSetFranchiseInitialAccount", "http.franchiseInitialAccount.set", "", ai.ModeWrite, "HIGH", "account:update", []string{"organizationId", "accountId"})
	spec.Path = "/api/franchise-initial-account"
	spec.Method = "POST"
	spec.WriteKind = ai.WriteExisting
	spec.TargetFields = []string{"organizationId", "accountId"}
	spec.ArgumentFields = []string{"organizationId", "accountId"}
	spec.AdditionalPermissions = []string{"hqMembership:read"}
	spec.RequiredAttestation = true
	return spec
}
