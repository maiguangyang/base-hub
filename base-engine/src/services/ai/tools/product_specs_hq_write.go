package tools

import (
	"base-engine/auth"
	"base-engine/src/services/ai"
)

func productHQWriteSpecs() []ai.ToolSpec {
	const hq = auth.WorkspaceTypeHeadquarters
	return append([]ai.ToolSpec{
		reviewedSpec("HqSetProductPackageBarcode", "graphql.mutation.hqSetProductPackageBarcode", `mutation HqSetProductPackageBarcode($id:ID!,$barcode:String!){hqSetProductPackageBarcode(id:$id,barcode:$barcode){id skuId barcode packageSetVersion enabled}}`, ai.ModeWrite, "HIGH", "hqProductCatalog:manage", hq, ai.WriteExisting, []string{"id"}, []string{"id", "barcode"}),
		reviewedSpec("HqCreateProductCategory", "graphql.mutation.hqCreateProductCategory", `mutation HqCreateProductCategory($input:HqCreateProductCategoryInput!){hqCreateProductCategory(input:$input){id name parentId enabled}}`, ai.ModeWrite, "HIGH", "hqProductCatalog:manage", hq, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("HqUpdateProductCategory", "graphql.mutation.hqUpdateProductCategory", `mutation HqUpdateProductCategory($id:ID!,$name:String!,$parentId:ID){hqUpdateProductCategory(id:$id,name:$name,parentId:$parentId){id name parentId enabled}}`, ai.ModeWrite, "HIGH", "hqProductCatalog:manage", hq, ai.WriteExisting, []string{"id"}, []string{"id", "name", "parentId"}),
		reviewedSpec("HqSetProductCategoryEnabled", "graphql.mutation.hqSetProductCategoryEnabled", `mutation HqSetProductCategoryEnabled($id:ID!,$enabled:Boolean!){hqSetProductCategoryEnabled(id:$id,enabled:$enabled){id name enabled}}`, ai.ModeWrite, "HIGH", "hqProductCatalog:manage", hq, ai.WriteExisting, []string{"id"}, []string{"id", "enabled"}),
		reviewedSpec("HqCreateProduct", "graphql.mutation.hqCreateProduct", `mutation HqCreateProduct($input:HqCreateProductInput!){hqCreateProduct(input:$input){id name categoryId enabled}}`, ai.ModeWrite, "HIGH", "hqProductCatalog:manage", hq, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("HqUpdateProduct", "graphql.mutation.hqUpdateProduct", `mutation HqUpdateProduct($id:ID!,$input:HqCreateProductInput!){hqUpdateProduct(id:$id,input:$input){id name categoryId enabled}}`, ai.ModeWrite, "HIGH", "hqProductCatalog:manage", hq, ai.WriteExisting, []string{"id"}, []string{"id", "input"}),
		reviewedSpec("HqSetProductEnabled", "graphql.mutation.hqSetProductEnabled", `mutation HqSetProductEnabled($id:ID!,$enabled:Boolean!){hqSetProductEnabled(id:$id,enabled:$enabled){id name enabled}}`, ai.ModeWrite, "HIGH", "hqProductCatalog:manage", hq, ai.WriteExisting, []string{"id"}, []string{"id", "enabled"}),
		reviewedSpec("HqSetProductSkuEnabled", "graphql.mutation.hqSetProductSkuEnabled", `mutation HqSetProductSkuEnabled($id:ID!,$enabled:Boolean!){hqSetProductSkuEnabled(id:$id,enabled:$enabled){id productId name enabled}}`, ai.ModeWrite, "HIGH", "hqProductCatalog:manage", hq, ai.WriteExisting, []string{"id"}, []string{"id", "enabled"}),
		reviewedSpec("HqCreateProductPackage", "graphql.mutation.hqCreateProductPackage", `mutation HqCreateProductPackage($input:HqCreateProductPackageInput!){hqCreateProductPackage(input:$input){id skuId name barcode packageSetVersion containsPackageId containsQuantity enabled}}`, ai.ModeWrite, "HIGH", "hqProductCatalog:manage", hq, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("HqUpdateProductPackageMetadata", "graphql.mutation.hqUpdateProductPackageMetadata", `mutation HqUpdateProductPackageMetadata($id:ID!,$name:String!,$suggestedPriceFen:Int){hqUpdateProductPackageMetadata(id:$id,name:$name,suggestedPriceFen:$suggestedPriceFen){id skuId name suggestedPriceFen enabled}}`, ai.ModeWrite, "HIGH", "hqProductCatalog:manage", hq, ai.WriteExisting, []string{"id"}, []string{"id", "name", "suggestedPriceFen"}),
		reviewedSpec("HqRetireProductPackage", "graphql.mutation.hqRetireProductPackage", `mutation HqRetireProductPackage($id:ID!){hqRetireProductPackage(id:$id){id skuId name enabled}}`, ai.ModeWrite, "HIGH", "hqProductCatalog:manage", hq, ai.WriteExisting, []string{"id"}, []string{"id"}),
		reviewedSpec("HqPublishProductPackageSet", "graphql.mutation.hqPublishProductPackageSet", `mutation HqPublishProductPackageSet($skuId:ID!,$version:Int!){hqPublishProductPackageSet(skuId:$skuId,version:$version){id skuId packageSetVersion enabled}}`, ai.ModeWrite, "HIGH", "hqProductCatalog:manage", hq, ai.WriteExisting, []string{"skuId"}, []string{"skuId", "version"}),
	}, productHQTemplateWriteSpecs()...)
}
