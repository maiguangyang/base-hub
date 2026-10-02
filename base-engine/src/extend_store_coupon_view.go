package src

import "base-engine/gen"

func storeCouponTemplateView(item *gen.CustomerCouponTemplate) (*gen.FranchiseCouponTemplateView, error) {
	hqView, err := couponTemplateGraphQL(item)
	if err != nil {
		return nil, err
	}
	storeID := ""
	if item.ApplicableStoreID != nil {
		storeID = *item.ApplicableStoreID
	}
	return &gen.FranchiseCouponTemplateView{ID: item.ID, StoreID: storeID, Code: item.Code, Title: item.Title,
		AmountFen: hqView.AmountFen, MinSpendFen: hqView.MinSpendFen, DaysAfterActivation: hqView.DaysAfterActivation,
		EffectiveAt: hqView.EffectiveAt, DistributionEndsAt: hqView.DistributionEndsAt,
		PerMemberLimit: hqView.PerMemberLimit, TotalIssueLimit: hqView.TotalIssueLimit, IssuedCount: hqView.IssuedCount,
		Enabled: item.Enabled, IssuerScope: item.IssuerScope}, nil
}
func storeCouponGrantView(item *gen.CustomerCouponGrant) (*gen.FranchiseCouponGrantView, error) {
	hqView, err := couponGrantGraphQL(item)
	if err != nil {
		return nil, err
	}
	return &gen.FranchiseCouponGrantView{ID: item.ID, TemplateID: item.TemplateID, MemberID: item.MemberID,
		Status: item.Status, AmountFen: hqView.AmountFen, MinSpendFen: hqView.MinSpendFen,
		DaysAfterActivation: hqView.DaysAfterActivation, IssuedAt: hqView.IssuedAt, RevokedAt: hqView.RevokedAt}, nil
}
