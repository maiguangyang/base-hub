package tools

import (
	"base-engine/auth"
	"base-engine/src/services/ai"
)

func stocktakeReadSpecs() []ai.ToolSpec {
	const franchise = auth.WorkspaceTypeFranchise
	return []ai.ToolSpec{
		reviewedSpec("FranchiseStocktakeBatchChoices", "graphql.query.franchiseStocktakeBatchChoices", `query FranchiseStocktakeBatchChoices($storeId:ID!,$listingId:ID!,$page:Int!,$perPage:Int!){franchiseStocktakeBatchChoices(storeId:$storeId,listingId:$listingId,page:$page,perPage:$perPage){data{id listingId batchNumber expiresAt} total currentPage perPage}}`, ai.ModeReadOnly, "LOW", "franchiseStocktake:read", franchise, "", nil, []string{"storeId", "listingId", "page", "perPage"}),
		reviewedSpec("FranchiseStocktakes", "graphql.query.franchiseStocktakes", `query FranchiseStocktakes($storeId:ID!,$status:StocktakeStatus,$listingId:ID,$batchId:ID,$hasDifference:Boolean,$from:Time,$to:Time,$page:Int!,$perPage:Int!){franchiseStocktakes(storeId:$storeId,status:$status,listingId:$listingId,batchId:$batchId,hasDifference:$hasDifference,from:$from,to:$to,page:$page,perPage:$perPage,includeHistory:false){data{id storeId status startedAt lines{id batchId packageId countedQuantity snapshotQuantity difference needsRecount}} total currentPage perPage}}`, ai.ModeReadOnly, "LOW", "franchiseStocktake:read", franchise, "", nil, []string{"storeId", "status", "listingId", "batchId", "hasDifference", "from", "to", "page", "perPage"}),
		reviewedSpec("FranchiseStocktake", "graphql.query.franchiseStocktake", `query FranchiseStocktake($storeId:ID!,$id:ID!){franchiseStocktake(storeId:$storeId,id:$id){id storeId status startedAt reviewedAt postedAt initiatedByAccountId postedById lines{id batchId listingId batchNumber packageId packageName packageSetVersion countedQuantity snapshotQuantity difference reasonCode reasonNote countHistory{actorAccountId countedQuantity countedAt} needsRecount}}}`, ai.ModeReadOnly, "LOW", "franchiseStocktake:read", franchise, "", nil, []string{"storeId", "id"}),
	}
}

func stocktakeWriteSpecs() []ai.ToolSpec {
	const franchise = auth.WorkspaceTypeFranchise
	return []ai.ToolSpec{
		reviewedSpec("FranchiseCreateStocktake", "graphql.mutation.franchiseCreateStocktake", `mutation FranchiseCreateStocktake($input:FranchiseCreateStocktakeInput!){franchiseCreateStocktake(input:$input){id storeId status lines{id batchId packageId countedQuantity snapshotQuantity}}}`, ai.ModeWrite, "HIGH", "franchiseStocktake:record", franchise, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("FranchiseAddStocktakeLine", "graphql.mutation.franchiseAddStocktakeLine", `mutation FranchiseAddStocktakeLine($storeId:ID!,$id:ID!,$batchId:ID!,$packageId:ID!){franchiseAddStocktakeLine(storeId:$storeId,id:$id,batchId:$batchId,packageId:$packageId){id status lines{id batchId packageId countedQuantity snapshotQuantity}}}`, ai.ModeWrite, "HIGH", "franchiseStocktake:record", franchise, ai.WriteExisting, []string{"storeId", "id", "batchId", "packageId"}, []string{"storeId", "id", "batchId", "packageId"}),
		reviewedSpec("FranchiseRecordStocktakeLine", "graphql.mutation.franchiseRecordStocktakeLine", `mutation FranchiseRecordStocktakeLine($storeId:ID!,$id:ID!,$lineId:ID!,$quantity:Int!){franchiseRecordStocktakeLine(storeId:$storeId,id:$id,lineId:$lineId,quantity:$quantity){id status lines{id countedQuantity snapshotQuantity needsRecount}}}`, ai.ModeWrite, "HIGH", "franchiseStocktake:record", franchise, ai.WriteExisting, []string{"storeId", "id", "lineId"}, []string{"storeId", "id", "lineId", "quantity"}),
		reviewedSpec("FranchiseSubmitStocktake", "graphql.mutation.franchiseSubmitStocktake", `mutation FranchiseSubmitStocktake($storeId:ID!,$id:ID!){franchiseSubmitStocktake(storeId:$storeId,id:$id){id status lines{id countedQuantity snapshotQuantity difference needsRecount}}}`, ai.ModeWrite, "HIGH", "franchiseStocktake:record", franchise, ai.WriteExisting, []string{"storeId", "id"}, []string{"storeId", "id"}),
		reviewedSpec("FranchiseSetStocktakeReason", "graphql.mutation.franchiseSetStocktakeReason", `mutation FranchiseSetStocktakeReason($storeId:ID!,$id:ID!,$lineId:ID!,$reasonCode:String!,$note:String){franchiseSetStocktakeReason(storeId:$storeId,id:$id,lineId:$lineId,reasonCode:$reasonCode,note:$note){id status lines{id reasonCode reasonNote difference}}}`, ai.ModeWrite, "HIGH", "franchiseStocktake:post", franchise, ai.WriteExisting, []string{"storeId", "id", "lineId"}, []string{"storeId", "id", "lineId", "reasonCode", "note"}),
		reviewedSpec("FranchiseReturnStocktake", "graphql.mutation.franchiseReturnStocktake", `mutation FranchiseReturnStocktake($storeId:ID!,$id:ID!){franchiseReturnStocktake(storeId:$storeId,id:$id){id status lines{id countedQuantity snapshotQuantity}}}`, ai.ModeWrite, "HIGH", "franchiseStocktake:post", franchise, ai.WriteExisting, []string{"storeId", "id"}, []string{"storeId", "id"}),
		reviewedSpec("FranchiseCancelStocktake", "graphql.mutation.franchiseCancelStocktake", `mutation FranchiseCancelStocktake($storeId:ID!,$id:ID!){franchiseCancelStocktake(storeId:$storeId,id:$id){id status canceledAt}}`, ai.ModeWrite, "HIGH", "franchiseStocktake:post", franchise, ai.WriteExisting, []string{"storeId", "id"}, []string{"storeId", "id"}),
		reviewedSpec("FranchisePostStocktake", "graphql.mutation.franchisePostStocktake", `mutation FranchisePostStocktake($storeId:ID!,$id:ID!){franchisePostStocktake(storeId:$storeId,id:$id){id status postedAt postedById lines{id countedQuantity snapshotQuantity difference reasonCode}}}`, ai.ModeWrite, "HIGH", "franchiseStocktake:post", franchise, ai.WriteExisting, []string{"storeId", "id"}, []string{"storeId", "id"}),
	}
}
