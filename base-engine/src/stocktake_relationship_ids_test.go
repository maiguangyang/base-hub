package src

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestStocktakeGeneratedRelationshipIDsAreDenied(t *testing.T) {
	ctx := auth.WithPrincipal(context.Background(), &auth.WorkspacePrincipal{AccountID: "operator"})
	for _, test := range []struct {
		name string
		read func() ([]string, error)
	}{
		{"account created", func() ([]string, error) { return (&AccountResolver{}).CreatedStocktakesIds(ctx, &gen.Account{}) }},
		{"account posted", func() ([]string, error) { return (&AccountResolver{}).PostedStocktakesIds(ctx, &gen.Account{}) }},
		{"store", func() ([]string, error) { return (&StoreResolver{}).StocktakesIds(ctx, &gen.Store{}) }},
		{"package", func() ([]string, error) {
			return (&ProductPackageResolver{}).StocktakeLinesIds(ctx, &gen.ProductPackage{})
		}},
		{"batch", func() ([]string, error) {
			return (&StoreInventoryBatchResolver{}).StocktakeLinesIds(ctx, &gen.StoreInventoryBatch{})
		}},
		{"stocktake", func() ([]string, error) { return (&StoreStocktakeResolver{}).LinesIds(ctx, &gen.StoreStocktake{}) }},
		{"line", func() ([]string, error) {
			return (&StoreStocktakeLineResolver{}).MovementsIds(ctx, &gen.StoreStocktakeLine{})
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			ids, err := test.read()
			if len(ids) != 0 || auth.ErrorCode(err) != auth.CodePermissionDenied {
				t.Fatalf("generated relation returned %v, %v", ids, err)
			}
		})
	}
}
