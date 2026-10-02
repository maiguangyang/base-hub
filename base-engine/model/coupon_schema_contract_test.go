package model

import "testing"

func TestCustomerCouponTimestampSchema(t *testing.T) {
	modelSchema := readSchemaFile(t, "model.graphql")
	extendSchema := readSchemaFile(t, "extend.graphql")
	assertCouponStorageSchema(t, modelSchema)
	assertCouponAPISchema(t, extendSchema)
	assertCouponJobSchema(t, modelSchema)
}

func assertCouponStorageSchema(t *testing.T, schema string) {
	t.Helper()
	assertContainsAll(t, schema, []string{
		`effectiveAt: Int! @column(gorm: "type:bigint(13);NOT NULL;")`,
		`distributionEndsAt: Int @column(gorm: "type:bigint(13);")`,
		`issuedAt: Int! @column(gorm: "type:bigint(13);NOT NULL;")`,
		`activatedAt: Int @column(gorm: "type:bigint(13);")`,
		`expiresAt: Int @column(gorm: "type:bigint(13);")`,
		`revokedAt: Int @column(gorm: "type:bigint(13);")`,
	})
}

func assertCouponAPISchema(t *testing.T, schema string) {
	t.Helper()
	assertContainsAll(t, schema, []string{
		`effectiveAt: Time!`, `distributionEndsAt: Time`, `issuedAt: Time!`,
		`activatedAt: Time`, `expiresAt: Time`, `revokedAt: Time`,
		"input HqCreateCustomerCouponTemplateInput {\n  code: String!\n  requestKey: String!\n  title: String!\n  amountFen: Int!\n  minSpendFen: Int!\n  daysAfterActivation: Int!\n  effectiveAt: Time\n  distributionEndsAt: Time\n  perMemberLimit: Int!\n  totalIssueLimit: Int",
		`input FranchiseCreateCouponTemplateInput { storeId: ID! code: String! title: String! requestKey: String! amountFen: Int! minSpendFen: Int! daysAfterActivation: Int! effectiveAt: Time distributionEndsAt: Time perMemberLimit: Int! totalIssueLimit: Int enabled: Boolean! }`,
	})
}

func assertCouponJobSchema(t *testing.T, schema string) {
	t.Helper()
	assertContainsAll(t, schema, []string{
		`enum CustomerCouponDistributionJobKind`, `enum CustomerCouponDistributionJobStatus`,
		`type CustomerCouponDistributionJob @entity(title: "优惠券派发任务")`,
		`availableAt: Int! @column(gorm: "type:bigint(13);NOT NULL;")`,
		`leaseExpiresAt: Int @column(gorm: "type:bigint(13);")`,
		`cursorCreatedAt: Int @column(gorm: "type:bigint(13);")`,
	})
}
