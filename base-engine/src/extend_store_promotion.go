package src

import (
	"context"

	"base-engine/gen"
	"base-engine/src/services/storemerchandising"
)

func optionalInt64(value *int) *int64 {
	if value == nil {
		return nil
	}
	converted := int64(*value)
	return &converted
}

func promotionView(item storemerchandising.PromotionView) (*gen.FranchisePromotionView, error) {
	value := item.Promotion
	version, err := catalogInt(value.Version)
	if err != nil {
		return nil, err
	}
	thresholdFen, err := catalogOptionalInt(value.ThresholdFen)
	if err != nil {
		return nil, err
	}
	thresholdQuantity, err := catalogOptionalInt(value.ThresholdQuantity)
	if err != nil {
		return nil, err
	}
	discountFen, err := catalogOptionalInt(value.DiscountFen)
	if err != nil {
		return nil, err
	}
	discountBasisPoints, err := catalogOptionalInt(value.DiscountBasisPoints)
	if err != nil {
		return nil, err
	}
	fixedPriceFen, err := catalogOptionalInt(value.FixedPriceFen)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchisePromotionView{ID: value.ID, StoreID: value.StoreID, RuleKey: value.RuleKey, Version: version, Latest: item.Latest,
		Kind: value.Kind, TimeZone: value.TimeZone, Enabled: value.Enabled, StartsAt: value.StartsAt, EndsAt: value.EndsAt,
		ThresholdFen: thresholdFen, ThresholdQuantity: thresholdQuantity, DiscountFen: discountFen,
		DiscountBasisPoints: discountBasisPoints, FixedPriceFen: fixedPriceFen,
		StackWithHqCoupon: value.StackWithHqCoupon, StackWithStoreCoupon: value.StackWithStoreCoupon,
		StackWithMemberPrice: value.StackWithMemberPrice, CheckoutUnavailable: true,
		Targets: make([]*gen.FranchisePromotionTargetView, 0, len(item.Targets))}
	for _, target := range item.Targets {
		quantity, err := catalogInt(target.RequiredQuantity)
		if err != nil {
			return nil, err
		}
		result.Targets = append(result.Targets, &gen.FranchisePromotionTargetView{OfferID: target.OfferID, RequiredQuantity: quantity})
	}
	return result, nil
}

func (r *QueryResolver) FranchisePromotions(ctx context.Context, storeID string, kind *gen.StorePromotionKind, enabled *bool, page, perPage int) (*gen.FranchisePromotionPage, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	items, err := r.Services.Merchandising.PromotionsFiltered(ctx, principal, storeID, storemerchandising.PromotionFilter{Kind: kind, Enabled: enabled}, page, perPage)
	if err != nil {
		return nil, err
	}
	total, err := catalogInt(items.Total)
	if err != nil {
		return nil, err
	}
	result := &gen.FranchisePromotionPage{Data: make([]*gen.FranchisePromotionView, 0, len(items.Data)), Total: total, CurrentPage: page, PerPage: perPage}
	for _, item := range items.Data {
		view, err := promotionView(item)
		if err != nil {
			return nil, err
		}
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (r *MutationResolver) FranchiseSavePromotion(ctx context.Context, input gen.FranchiseSavePromotionInput) (*gen.FranchisePromotionView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	targets := make([]storemerchandising.PromotionTargetInput, 0, len(input.Targets))
	for _, item := range input.Targets {
		targets = append(targets, storemerchandising.PromotionTargetInput{OfferID: item.OfferID, RequiredQuantity: int64(item.RequiredQuantity)})
	}
	value := storemerchandising.PromotionInput{StoreID: input.StoreID, PromotionID: input.PromotionID, Kind: input.Kind,
		TimeZone: input.TimeZone, StartsAt: input.StartsAt, EndsAt: input.EndsAt,
		ThresholdFen: optionalInt64(input.ThresholdFen), ThresholdQuantity: optionalInt64(input.ThresholdQuantity),
		DiscountFen: optionalInt64(input.DiscountFen), DiscountBasisPoints: optionalInt64(input.DiscountBasisPoints),
		FixedPriceFen: optionalInt64(input.FixedPriceFen), StackWithHqCoupon: input.StackWithHqCoupon,
		StackWithStoreCoupon: input.StackWithStoreCoupon, StackWithMemberPrice: input.StackWithMemberPrice, Targets: targets}
	item, err := r.Services.Merchandising.SavePromotion(ctx, principal, value)
	if err != nil {
		return nil, err
	}
	view, err := r.Services.Merchandising.Promotion(ctx, principal, item.ID)
	if err != nil {
		return nil, err
	}
	return promotionView(*view)
}

func (r *MutationResolver) FranchiseSetPromotionEnabled(ctx context.Context, id string, enabled bool) (*gen.FranchisePromotionView, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	item, err := r.Services.Merchandising.SetPromotionEnabled(ctx, principal, id, enabled)
	if err != nil {
		return nil, err
	}
	view, err := r.Services.Merchandising.Promotion(ctx, principal, item.ID)
	if err != nil {
		return nil, err
	}
	return promotionView(*view)
}
