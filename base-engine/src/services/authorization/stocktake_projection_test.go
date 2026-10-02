package authorization

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestStocktakeGeneratedParentQueriesRejectRelations(t *testing.T) {
	principal := &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeFranchise,
		AccountID: "operator", OrganizationID: stringPointer("org"), AllStores: true,
		Permissions: map[string]struct{}{"operatorMembership:read": {}, "store:read": {}}}
	ctx := auth.WithPrincipal(context.Background(), principal)
	handlers := RegisterHandlers(gen.DefaultResolutionHandlers())
	stocktake := &gen.StoreStocktakeFilterType{Lines: &gen.StoreStocktakeLineFilterType{SnapshotQuantityGt: intPointer(0)}}
	sortOrder := gen.ObjectSortType("ASC")
	cases := []struct {
		name string
		run  func() error
	}{
		{"account nested filter", func() error {
			_, err := handlers.QueryAccounts(ctx, nil, gen.QueryAccountsHandlerOptions{Filter: &gen.AccountFilterType{Or: []*gen.AccountFilterType{{CreatedStocktakes: stocktake}}}})
			return err
		}},
		{"account sort", func() error {
			_, err := handlers.QueryAccounts(ctx, nil, gen.QueryAccountsHandlerOptions{Sort: []*gen.AccountSortType{{PostedStocktakesIds: &sortOrder}}})
			return err
		}},
		{"store filter", func() error {
			_, err := handlers.QueryStores(ctx, nil, gen.QueryStoresHandlerOptions{Filter: &gen.StoreFilterType{Stocktakes: stocktake}})
			return err
		}},
		{"store nested sort", func() error {
			_, err := handlers.QueryStores(ctx, nil, gen.QueryStoresHandlerOptions{Sort: []*gen.StoreSortType{{ReviewedByAccount: &gen.AccountSortType{CreatedStocktakes: &gen.StoreStocktakeSortType{}}}}})
			return err
		}},
		{"organization through store", func() error {
			_, err := handlers.QueryOrganizations(ctx, nil, gen.QueryOrganizationsHandlerOptions{Filter: &gen.OrganizationFilterType{Stores: &gen.StoreFilterType{Stocktakes: stocktake}}})
			return err
		}},
		{"membership through account", func() error {
			_, err := handlers.QueryOperatorMemberships(ctx, nil, gen.QueryOperatorMembershipsHandlerOptions{Filter: &gen.OperatorMembershipFilterType{Account: &gen.AccountFilterType{PostedStocktakes: stocktake}}})
			return err
		}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if code := auth.ErrorCode(test.run()); code != auth.CodePermissionDenied {
				t.Fatalf("stocktake projection was accepted: %s", code)
			}
		})
	}
}

func TestStocktakeNestedStoreWritesAreDenied(t *testing.T) {
	for _, field := range []string{"stocktakes", "stocktakesIds"} {
		if code := auth.ErrorCode(rejectStoreManagedRelations(map[string]interface{}{field: []string{"unsafe"}})); code != auth.CodePermissionDenied {
			t.Fatalf("nested %s accepted: %s", field, code)
		}
	}
}

func stringPointer(value string) *string { return &value }
func intPointer(value int) *int          { return &value }
