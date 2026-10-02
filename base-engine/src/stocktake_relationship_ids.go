package src

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
)

// Generated relationship-ID resolvers bypass ResolutionHandlers. The dedicated
// franchise stocktake API is the only read path for these records.
func denyStocktakeRelationshipIDs(ctx context.Context) ([]string, error) {
	if _, err := auth.RequirePrincipal(ctx); err != nil {
		return nil, err
	}
	return nil, auth.NewError(auth.CodePermissionDenied)
}

func (r *AccountResolver) CreatedStocktakesIds(ctx context.Context, _ *gen.Account) ([]string, error) {
	return denyStocktakeRelationshipIDs(ctx)
}

func (r *AccountResolver) PostedStocktakesIds(ctx context.Context, _ *gen.Account) ([]string, error) {
	return denyStocktakeRelationshipIDs(ctx)
}

func (r *StoreResolver) StocktakesIds(ctx context.Context, _ *gen.Store) ([]string, error) {
	return denyStocktakeRelationshipIDs(ctx)
}

func (r *ProductPackageResolver) StocktakeLinesIds(ctx context.Context, _ *gen.ProductPackage) ([]string, error) {
	return denyStocktakeRelationshipIDs(ctx)
}

func (r *StoreInventoryBatchResolver) StocktakeLinesIds(ctx context.Context, _ *gen.StoreInventoryBatch) ([]string, error) {
	return denyStocktakeRelationshipIDs(ctx)
}

func (r *StoreStocktakeResolver) LinesIds(ctx context.Context, _ *gen.StoreStocktake) ([]string, error) {
	return denyStocktakeRelationshipIDs(ctx)
}

func (r *StoreStocktakeLineResolver) MovementsIds(ctx context.Context, _ *gen.StoreStocktakeLine) ([]string, error) {
	return denyStocktakeRelationshipIDs(ctx)
}
