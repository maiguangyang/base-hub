package model

import (
	"strings"
	"testing"
)

func TestProductVariantSchemaContract(t *testing.T) {
	modelSchema := readSchemaFile(t, "model.graphql")
	extendSchema := readSchemaFile(t, "extend.graphql")
	assertContainsAll(t, modelSchema, []string{
		`type SpecificationDefinition @entity(title: "商品规格")`,
		`type SpecificationValue @entity(title: "商品规格值")`,
		`type ProductSpecificationChoice @entity(title: "商品已选规格值")`,
		`type ProductSkuSpecificationValue @entity(title: "SKU 规格值")`,
		`specificationDefinitions: [SpecificationDefinition!]! @relationship(inverse: "organization")`,
		`values: [SpecificationValue!]! @relationship(inverse: "specification")`,
		`specification: SpecificationDefinition! @relationship(inverse: "values")`,
		`defaultPackageTemplate: ProductPackageTemplate @relationship(inverse: "defaultProducts")`,
	})
	assertContainsAll(t, extendSchema, []string{
		`hqSpecifications(q: String, enabled: Boolean, page: Int!, perPage: Int!): HqSpecificationPage!`,
		`selections: [HqProductSpecificationSelectionInput!]!`,
		`skuOverrides: [HqProductSkuOverrideInput!]!`,
	})
	if strings.Contains(modelSchema, "type ProductSkuTemplate @entity") || strings.Contains(extendSchema, "hqCreateProductSkuTemplate") {
		t.Fatal("old sales specification template is still exposed")
	}
	if strings.Contains(modelSchema, `code: String! @column(gorm: "type:varchar(13);NOT NULL;uniqueIndex;")`) {
		t.Fatal("obsolete catalog code exposed")
	}
	assertNoScalarEntityRelationshipIDs(t, modelSchema)
}
