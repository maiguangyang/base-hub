package src

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
)

func (r *QueryResolver) franchiseOfferWithRevisions(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string, item *gen.StorePackageOffer) (*gen.FranchiseOfferView, error) {
	view, err := offerView(item)
	if err != nil {
		return nil, err
	}
	view.PriceRevisions = []*gen.FranchisePriceRevisionView{}
	if !principal.Has("franchiseProduct:read") {
		return view, nil
	}
	revisions, err := r.Services.Merchandising.OfferPriceRevisions(ctx, principal, storeID, item.ListingID, item.ID)
	if err != nil {
		return nil, err
	}
	view.PriceRevisions = make([]*gen.FranchisePriceRevisionView, 0, len(revisions))
	for _, revision := range revisions {
		oldPrice, err := catalogOptionalInt(revision.PreviousPriceFen)
		if err != nil {
			return nil, err
		}
		newPrice, err := catalogInt(revision.PriceFen)
		if err != nil {
			return nil, err
		}
		view.PriceRevisions = append(view.PriceRevisions, &gen.FranchisePriceRevisionView{
			ID:          revision.ID,
			OldPriceFen: oldPrice, NewPriceFen: newPrice,
			ReasonCode: revision.ReasonCode, EffectiveAt: revision.EffectiveAt,
		})
	}
	return view, nil
}
