package tools

import (
	"base-engine/auth"
	"base-engine/src/services/ai"
)

func customerReadSpecs() []ai.ToolSpec {
	const hq = auth.WorkspaceTypeHeadquarters
	return []ai.ToolSpec{
		reviewedSpec("HqCustomerMembers", "graphql.query.hqCustomerMembers", `query HqCustomerMembers($phone: String, $status: CustomerMemberStatus, $page: Int!, $perPage: Int!) {
  hqCustomerMembers(phone: $phone, status: $status, page: $page, perPage: $perPage) {
    data { id memberNumber phoneMasked status pointsBalance pointsFrozen pendingCouponCount createdAt }
    total currentPage perPage
  }
}`, ai.ModeReadOnly, "LOW", "hqCustomer:read", hq, "", nil, []string{"phone", "status", "page", "perPage"}),
		reviewedSpec("HqCustomerMember", "graphql.query.hqCustomerMember", `query HqCustomerMember($id: ID!) {
  hqCustomerMember(id: $id) { id memberNumber phoneMasked status pointsBalance pointsFrozen pendingCouponCount createdAt }
}`, ai.ModeReadOnly, "LOW", "hqCustomer:read", hq, "", nil, []string{"id"}),
		reviewedSpec("HqCustomerSensitivePhone", "graphql.query.hqCustomerSensitivePhone", `query HqCustomerSensitivePhone($id: ID!) {
  hqCustomerSensitivePhone(id: $id)
}`, ai.ModeReadOnly, "HIGH", "customer:read_sensitive", hq, "", nil, []string{"id"}),
		reviewedSpec("HqCustomerBenefitPolicy", "graphql.query.hqCustomerBenefitPolicy", `query HqCustomerBenefitPolicy {
  hqCustomerBenefitPolicy {
    version discountEnabled discountBasisPoints purchaseEarnEnabled earnAmountFen earnPoints
    redemptionEnabled redeemPoints redeemAmountFen maxRedemptionBasisPoints maxRedemptionPoints
    manualGrantMaxSingle manualGrantMaxDaily promotionWithHqCoupon promotionWithStoreCoupon
    memberPriceWithPromotion pointsWithPromotion pointsWithCoupon
  }
}`, ai.ModeReadOnly, "LOW", "hqCustomerPolicy:read", hq, "", nil, nil),
		reviewedSpec("HqCustomerPointEntries", "graphql.query.hqCustomerPointEntries", `query HqCustomerPointEntries($memberId: ID!, $page: Int!, $perPage: Int!) {
  hqCustomerPointEntries(memberId: $memberId, page: $page, perPage: $perPage) {
    data { id memberId delta source operationKind reasonCode note reversesId createdAt }
    total currentPage perPage
  }
}`, ai.ModeReadOnly, "LOW", "hqCustomerPoints:read", hq, "", nil, []string{"memberId", "page", "perPage"}),
		reviewedSpec("HqCustomerCouponTemplates", "graphql.query.hqCustomerCouponTemplates", `query HqCustomerCouponTemplates($q: String, $enabled: Boolean, $page: Int!, $perPage: Int!) {
  hqCustomerCouponTemplates(q: $q, enabled: $enabled, page: $page, perPage: $perPage) {
	    data { id code title amountFen minSpendFen daysAfterActivation effectiveAt distributionEndsAt perMemberLimit totalIssueLimit issuedCount enabled createdAt }
    total currentPage perPage
  }
}`, ai.ModeReadOnly, "LOW", "hqCustomerCoupon:read", hq, "", nil, []string{"q", "enabled", "page", "perPage"}),
		reviewedSpec("HqCustomerCouponGrants", "graphql.query.hqCustomerCouponGrants", `query HqCustomerCouponGrants($memberId: ID!, $templateId: ID, $status: CustomerCouponGrantStatus, $page: Int!, $perPage: Int!) {
  hqCustomerCouponGrants(memberId: $memberId, templateId: $templateId, status: $status, page: $page, perPage: $perPage) {
    data { id templateId memberId status amountFen minSpendFen daysAfterActivation issuedAt activatedAt expiresAt revokedAt }
    total currentPage perPage
  }
}`, ai.ModeReadOnly, "LOW", "hqCustomerCoupon:read", hq, "", nil, []string{"memberId", "templateId", "status", "page", "perPage"}),
	}
}
