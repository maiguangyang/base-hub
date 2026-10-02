package model

import "testing"

func TestStoreDocumentSchemaContract(t *testing.T) {
	modelSchema := readSchemaFile(t, "model.graphql")
	extendSchema := readSchemaFile(t, "extend.graphql")
	assertContainsAll(t, modelSchema, []string{
		`businessLicenseImageUrl: String @column`,
		`otherDocumentImageUrl: String @column`,
	})
	assertContainsAll(t, extendSchema, []string{
		`enum StoreDocumentKind {`,
		`BUSINESS_LICENSE`,
		`OTHER`,
		`setStoreDocument(storeId: ID!, kind: StoreDocumentKind!, attachmentId: ID!): Store!`,
		`removeStoreDocument(storeId: ID!, kind: StoreDocumentKind!): Store!`,
	})
}
