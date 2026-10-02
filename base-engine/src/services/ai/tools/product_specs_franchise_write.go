package tools

import (
	"base-engine/auth"
	"base-engine/src/services/ai"
)

func productFranchiseWriteSpecs() []ai.ToolSpec {
	const franchise = auth.WorkspaceTypeFranchise
	return []ai.ToolSpec{
		reviewedSpec("FranchiseSetStoreListing", "graphql.mutation.franchiseSetStoreListing", `mutation FranchiseSetStoreListing($storeId:ID!,$skuId:ID!,$enabled:Boolean!){franchiseSetStoreListing(storeId:$storeId,skuId:$skuId,enabled:$enabled){id storeId skuId enabled}}`, ai.ModeWrite, "HIGH", "franchiseProduct:manage", franchise, ai.WriteExisting, []string{"storeId", "skuId"}, []string{"storeId", "skuId", "enabled"}),
		reviewedSpec("FranchiseSetStorePrice", "graphql.mutation.franchiseSetStorePrice", `mutation FranchiseSetStorePrice($listingId:ID!,$packageId:ID!,$priceFen:Int!,$reasonCode:String!){franchiseSetStorePrice(listingId:$listingId,packageId:$packageId,priceFen:$priceFen,reasonCode:$reasonCode){id listingId packageId priceFen}}`, ai.ModeWrite, "HIGH", "franchiseProduct:manage", franchise, ai.WriteExisting, []string{"listingId", "packageId"}, []string{"listingId", "packageId", "priceFen", "reasonCode"}),
		reviewedSpec("FranchiseReceiveStock", "graphql.mutation.franchiseReceiveStock", `mutation FranchiseReceiveStock($input:FranchiseReceiveStockInput!){franchiseReceiveStock(input:$input){id storeId batchId kind targetQuantity targetPackageId requestKey}}`, ai.ModeWrite, "HIGH", "franchiseStock:manage", franchise, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("FranchiseUnpackStock", "graphql.mutation.franchiseUnpackStock", `mutation FranchiseUnpackStock($input:FranchiseUnpackStockInput!){franchiseUnpackStock(input:$input){id storeId batchId kind sourceQuantity targetQuantity factorSnapshot}}`, ai.ModeWrite, "HIGH", "franchiseStock:manage", franchise, ai.WriteExisting, []string{"input.storeId", "input.batchId", "input.sourcePackageId"}, []string{"input"}),
		reviewedSpec("FranchiseAdjustStock", "graphql.mutation.franchiseAdjustStock", `mutation FranchiseAdjustStock($input:FranchiseAdjustStockInput!){franchiseAdjustStock(input:$input){id storeId batchId kind sourceQuantity targetQuantity reasonCode}}`, ai.ModeWrite, "HIGH", "franchiseStock:manage", franchise, ai.WriteExisting, []string{"input.storeId", "input.batchId", "input.packageId"}, []string{"input"}),
		reviewedSpec("FranchiseSavePromotion", "graphql.mutation.franchiseSavePromotion", `mutation FranchiseSavePromotion($input:FranchiseSavePromotionInput!){franchiseSavePromotion(input:$input){id storeId ruleKey version kind enabled checkoutUnavailable}}`, ai.ModeWrite, "HIGH", "franchisePromotion:manage", franchise, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("FranchiseSetPromotionEnabled", "graphql.mutation.franchiseSetPromotionEnabled", `mutation FranchiseSetPromotionEnabled($id:ID!,$enabled:Boolean!){franchiseSetPromotionEnabled(id:$id,enabled:$enabled){id storeId version enabled checkoutUnavailable}}`, ai.ModeWrite, "HIGH", "franchisePromotion:manage", franchise, ai.WriteExisting, []string{"id"}, []string{"id", "enabled"}),
		reviewedSpec("FranchiseCreateCouponTemplate", "graphql.mutation.franchiseCreateCouponTemplate", `mutation FranchiseCreateCouponTemplate($input:FranchiseCreateCouponTemplateInput!){franchiseCreateCouponTemplate(input:$input){id storeId code title amountFen minSpendFen daysAfterActivation effectiveAt distributionEndsAt perMemberLimit totalIssueLimit issuedCount enabled issuerScope}}`, ai.ModeWrite, "HIGH", "franchiseCoupon:manage", franchise, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("FranchiseSetCouponTemplateEnabled", "graphql.mutation.franchiseSetCouponTemplateEnabled", `mutation FranchiseSetCouponTemplateEnabled($storeId:ID!,$templateId:ID!,$enabled:Boolean!){franchiseSetCouponTemplateEnabled(storeId:$storeId,templateId:$templateId,enabled:$enabled){id storeId enabled}}`, ai.ModeWrite, "HIGH", "franchiseCoupon:manage", franchise, ai.WriteExisting, []string{"storeId", "templateId"}, []string{"storeId", "templateId", "enabled"}),
		reviewedSpec("FranchiseGrantCoupon", "graphql.mutation.franchiseGrantCoupon", `mutation FranchiseGrantCoupon($storeId:ID!,$templateId:ID!,$memberId:ID!,$requestKey:String!){franchiseGrantCoupon(storeId:$storeId,templateId:$templateId,memberId:$memberId,requestKey:$requestKey){id templateId memberId status amountFen}}`, ai.ModeWrite, "HIGH", "franchiseCoupon:grant", franchise, ai.WriteExisting, []string{"storeId", "templateId", "memberId"}, []string{"storeId", "templateId", "memberId"}),
		reviewedSpec("FranchiseRevokeCoupon", "graphql.mutation.franchiseRevokeCoupon", `mutation FranchiseRevokeCoupon($storeId:ID!,$grantId:ID!,$reasonCode:String!){franchiseRevokeCoupon(storeId:$storeId,grantId:$grantId,reasonCode:$reasonCode){id templateId memberId status revokedAt}}`, ai.ModeWrite, "HIGH", "franchiseCoupon:revoke", franchise, ai.WriteExisting, []string{"storeId", "grantId"}, []string{"storeId", "grantId", "reasonCode"}),
	}
}

func productWriteSpecs() []ai.ToolSpec {
	specs := append(productHQWriteSpecs(), productFranchiseWriteSpecs()...)
	specs = append(specs, stocktakeWriteSpecs()...)
	for index := range specs {
		specs[index].SingleCallApproval = true
		switch specs[index].ID {
		case "FranchiseReceiveStock", "FranchiseSavePromotion", "FranchiseCreateCouponTemplate":
			specs[index].ScopeFields = []string{"input.storeId"}
		case "FranchiseCreateStocktake":
			specs[index].ScopeFields = []string{"input.storeId"}
		case "FranchiseAddStocktakeLine", "FranchiseRecordStocktakeLine", "FranchiseSubmitStocktake", "FranchiseSetStocktakeReason", "FranchiseReturnStocktake", "FranchiseCancelStocktake", "FranchisePostStocktake":
			specs[index].ScopeFields = []string{"storeId"}
		}
		switch specs[index].ID {
		case "FranchiseReceiveStock", "FranchiseUnpackStock", "FranchiseAdjustStock", "FranchiseCreateCouponTemplate", "FranchiseCreateStocktake":
			specs[index].RequestKeyPath = "input.requestKey"
		case "FranchiseGrantCoupon":
			specs[index].RequestKeyPath = "requestKey"
		}
		switch specs[index].ID {
		case "FranchiseCreateCouponTemplate":
			specs[index].Description += " 模板启用后会向总部会员池中的现有会员异步派发，后注册会员也会异步补发符合条件的券；门店券只能在创建门店使用。"
		case "FranchiseGrantCoupon":
			specs[index].Description += " 人工发放仍校验模板已生效、未到派发截止时刻、每会员上限和总发放量；门店券只能在创建门店使用。"
		}
	}
	return specs
}
