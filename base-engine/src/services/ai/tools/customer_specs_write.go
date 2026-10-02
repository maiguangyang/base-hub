package tools

import (
	"base-engine/auth"
	"base-engine/src/services/ai"
)

func customerWriteSpecs() []ai.ToolSpec {
	specs := customerWriteBaseSpecs()
	evidence := map[string][]string{
		"HqRequestCustomerCancellation":  {"identityEvidence", "basisCode"},
		"HqCompleteCustomerCancellation": {"dispositionReference"},
		"HqGrantCustomerPoints":          {"note"},
		"HqReverseCustomerPoints":        {"note"},
		"HqCorrectCustomerPoints":        {"evidenceReference"},
	}
	for index := range specs {
		specs[index].SingleCallApproval = true
		specs[index].EvidencePaths = evidence[specs[index].ID]
		switch specs[index].ID {
		case "HqCreateCustomerMember", "HqCreateCustomerCouponTemplate":
			specs[index].RequestKeyPath = "input.requestKey"
		case "HqGrantCustomerPoints", "HqReverseCustomerPoints", "HqCorrectCustomerPoints", "HqGrantCustomerCoupon":
			specs[index].RequestKeyPath = "requestKey"
		}
		if specs[index].ID == "HqCreateCustomerCouponTemplate" {
			specs[index].GeneratedCodePath = "input.code"
		}
		switch specs[index].ID {
		case "HqCreateCustomerMember":
			specs[index].Description += " 建档成功后会异步补发当前符合条件的总部券和各门店券。"
		case "HqSetCustomerMemberStatus":
			specs[index].Description += " 只有会员从 SUSPENDED 恢复为 ACTIVE 时，才会重新安排符合条件优惠券的异步补发；ACTIVE 到 ACTIVE 不会重置派发进度。"
		case "HqCreateCustomerCouponTemplate":
			specs[index].Description += " 模板启用后会向现有会员异步派发，后注册会员也会异步补发符合条件的券。"
		case "HqGrantCustomerCoupon":
			specs[index].Description += " 人工发放仍校验模板已生效、未到派发截止时刻、每会员上限和总发放量。总部券可供全部门店使用。"
		}
	}
	return specs
}

func customerWriteBaseSpecs() []ai.ToolSpec {
	const hq = auth.WorkspaceTypeHeadquarters
	return []ai.ToolSpec{
		reviewedSpec("HqCreateCustomerMember", "graphql.mutation.hqCreateCustomerMember", `mutation HqCreateCustomerMember($input: HqCreateCustomerMemberInput!) {
  hqCreateCustomerMember(input: $input) { id memberNumber phoneMasked status pointsBalance }
}`, ai.ModeWrite, "HIGH", "hqCustomer:create", hq, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("HqSetCustomerMemberStatus", "graphql.mutation.hqSetCustomerMemberStatus", `mutation HqSetCustomerMemberStatus($id: ID!, $status: CustomerMemberStatus!) {
  hqSetCustomerMemberStatus(id: $id, status: $status) { id memberNumber status }
}`, ai.ModeWrite, "HIGH", "hqCustomer:update", hq, ai.WriteExisting, []string{"id"}, []string{"id", "status"}),
		reviewedSpec("HqRequestCustomerCancellation", "graphql.mutation.hqRequestCustomerCancellation", `mutation HqRequestCustomerCancellation($id: ID!, $identityEvidence: String!, $basisCode: String!) {
  hqRequestCustomerCancellation(id: $id, identityEvidence: $identityEvidence, basisCode: $basisCode) { id memberNumber status }
}`, ai.ModeWrite, "HIGH", "hqCustomer:cancel", hq, ai.WriteExisting, []string{"id"}, []string{"id", "identityEvidence", "basisCode"}),
		reviewedSpec("HqCompleteCustomerCancellation", "graphql.mutation.hqCompleteCustomerCancellation", `mutation HqCompleteCustomerCancellation($id: ID!, $dispositionReference: String!) {
  hqCompleteCustomerCancellation(id: $id, dispositionReference: $dispositionReference) { id memberNumber status pointsBalance pendingCouponCount }
}`, ai.ModeWrite, "HIGH", "hqCustomer:cancel", hq, ai.WriteExisting, []string{"id"}, []string{"id", "dispositionReference"}),
		reviewedSpec("HqSaveCustomerBenefitPolicy", "graphql.mutation.hqSaveCustomerBenefitPolicy", `mutation HqSaveCustomerBenefitPolicy($expectedVersion: Int!, $input: HqCustomerBenefitPolicyInput!) {
  hqSaveCustomerBenefitPolicy(expectedVersion: $expectedVersion, input: $input) {
    version discountEnabled discountBasisPoints purchaseEarnEnabled earnAmountFen earnPoints redemptionEnabled redeemPoints redeemAmountFen maxRedemptionBasisPoints maxRedemptionPoints
    manualGrantMaxSingle manualGrantMaxDaily promotionWithHqCoupon promotionWithStoreCoupon
    memberPriceWithPromotion pointsWithPromotion pointsWithCoupon
  }
}`, ai.ModeWrite, "HIGH", "hqCustomerPolicy:manage", hq, ai.WriteCreate, nil, []string{"expectedVersion", "input"}),
		reviewedSpec("HqGrantCustomerPoints", "graphql.mutation.hqGrantCustomerPoints", `mutation HqGrantCustomerPoints($memberId: ID!, $points: Int!, $reasonCode: CustomerPointReasonCode!, $note: String!, $requestKey: String!) {
  hqGrantCustomerPoints(memberId: $memberId, points: $points, reasonCode: $reasonCode, note: $note, requestKey: $requestKey) {
    entry { id memberId delta reasonCode createdAt } currentBalance
  }
}`, ai.ModeWrite, "HIGH", "hqCustomerPoints:grant", hq, ai.WriteExisting, []string{"memberId"}, []string{"memberId", "points", "reasonCode", "note"}),
		reviewedSpec("HqReverseCustomerPoints", "graphql.mutation.hqReverseCustomerPoints", `mutation HqReverseCustomerPoints($entryId: ID!, $reasonCode: CustomerPointReasonCode!, $note: String!, $requestKey: String!) {
  hqReverseCustomerPoints(entryId: $entryId, reasonCode: $reasonCode, note: $note, requestKey: $requestKey) {
    entry { id memberId delta reasonCode reversesId createdAt } currentBalance
  }
}`, ai.ModeWrite, "HIGH", "hqCustomerPoints:reverse", hq, ai.WriteExisting, []string{"entryId"}, []string{"entryId", "reasonCode", "note"}),
		reviewedSpec("HqCorrectCustomerPoints", "graphql.mutation.hqCorrectCustomerPoints", `mutation HqCorrectCustomerPoints($memberId: ID!, $evidenceReference: String!, $requestKey: String!) {
  hqCorrectCustomerPoints(memberId: $memberId, evidenceReference: $evidenceReference, requestKey: $requestKey) {
    entry { id memberId delta reasonCode createdAt } currentBalance
  }
}`, ai.ModeWrite, "HIGH", "hqCustomerPoints:correct", hq, ai.WriteExisting, []string{"memberId"}, []string{"memberId", "evidenceReference"}),
		reviewedSpec("HqCreateCustomerCouponTemplate", "graphql.mutation.hqCreateCustomerCouponTemplate", `mutation HqCreateCustomerCouponTemplate($input: HqCreateCustomerCouponTemplateInput!) {
  hqCreateCustomerCouponTemplate(input: $input) { id code title amountFen minSpendFen daysAfterActivation effectiveAt distributionEndsAt perMemberLimit totalIssueLimit enabled }
}`, ai.ModeWrite, "HIGH", "hqCustomerCoupon:manage", hq, ai.WriteCreate, nil, []string{"input"}),
		reviewedSpec("HqSetCustomerCouponTemplateEnabled", "graphql.mutation.hqSetCustomerCouponTemplateEnabled", `mutation HqSetCustomerCouponTemplateEnabled($templateId: ID!, $enabled: Boolean!) {
  hqSetCustomerCouponTemplateEnabled(templateId: $templateId, enabled: $enabled) { id code title enabled }
}`, ai.ModeWrite, "HIGH", "hqCustomerCoupon:manage", hq, ai.WriteExisting, []string{"templateId"}, []string{"templateId", "enabled"}),
		reviewedSpec("HqGrantCustomerCoupon", "graphql.mutation.hqGrantCustomerCoupon", `mutation HqGrantCustomerCoupon($templateId: ID!, $memberId: ID!, $requestKey: String!) {
  hqGrantCustomerCoupon(templateId: $templateId, memberId: $memberId, requestKey: $requestKey) { id templateId memberId status amountFen minSpendFen }
}`, ai.ModeWrite, "HIGH", "hqCustomerCoupon:grant", hq, ai.WriteExisting, []string{"templateId", "memberId"}, []string{"templateId", "memberId"}),
		customerRevokeCouponSpec(),
	}
}

func customerRevokeCouponSpec() ai.ToolSpec {
	return reviewedSpec("HqRevokeCustomerCoupon", "graphql.mutation.hqRevokeCustomerCoupon", `mutation HqRevokeCustomerCoupon($grantId: ID!, $reasonCode: String!) {
  hqRevokeCustomerCoupon(grantId: $grantId, reasonCode: $reasonCode) { id templateId memberId status revokedAt }
}`, ai.ModeWrite, "HIGH", "hqCustomerCoupon:revoke", auth.WorkspaceTypeHeadquarters, ai.WriteExisting, []string{"grantId"}, []string{"grantId", "reasonCode"})
}
