package authorization

import (
	"context"
	"base-engine/auth"
	"base-engine/gen"
)

func storeScope(ctx context.Context, client *gen.StoreFilterType) (*auth.WorkspacePrincipal, *gen.StoreFilterType, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, nil, err
	}
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		if principal.Has("store:read_all") {
			if err := Authorize(principal, Intent{Action: "store:read_all", Mode: AccessRead}); err != nil {
				return nil, nil, err
			}
			return principal, activeStoreFilter(client), nil
		}
		if err := Authorize(principal, Intent{Action: "hqStore:read", Mode: AccessRead}); err != nil {
			return nil, nil, err
		}
		organizationID, err := requireOrganization(principal)
		if err != nil {
			return nil, nil, err
		}
		return principal, scopedStoreFilter(client, organizationID, nil), nil
	}
	if err := Authorize(principal, Intent{Action: "store:read", Mode: AccessRead}); err != nil {
		return nil, nil, err
	}
	organizationID, err := requireOrganization(principal)
	if err != nil {
		return nil, nil, err
	}
	var storeIDs []string
	if !principal.AllStores {
		storeIDs = mapKeys(principal.StoreIDs)
	}
	return principal, scopedStoreFilter(client, organizationID, storeIDs), nil
}

func scopedStoreFilter(client *gen.StoreFilterType, organizationID string, storeIDs []string) *gen.StoreFilterType {
	scope := &gen.StoreFilterType{OrganizationID: &organizationID}
	if storeIDs != nil {
		scope.IDIn = storeIDs
	}
	filters := []*gen.StoreFilterType{activeStoreFilter(nil), scope}
	if client != nil {
		filters = append(filters, client)
	}
	return &gen.StoreFilterType{And: filters}
}

func activeStoreFilter(client *gen.StoreFilterType) *gen.StoreFilterType {
	active := 1
	null := true
	filter := &gen.StoreFilterType{Or: []*gen.StoreFilterType{{IsDelete: &active}, {IsDeleteNull: &null}}}
	if client == nil {
		return filter
	}
	return &gen.StoreFilterType{And: []*gen.StoreFilterType{filter, client}}
}

func queryStore(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryStoreHandlerOptions) (*gen.Store, error) {
	principal, filter, err := storeScope(ctx, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	item, err := gen.QueryStoreHandler(ctx, resolver, options)
	if err != nil {
		return nil, concealScopedNotFound(err, principal)
	}
	return item, nil
}

func queryStores(ctx context.Context, resolver *gen.GeneratedResolver, options gen.QueryStoresHandlerOptions) (*gen.StoreResultType, error) {
	_, filter, err := storeScope(ctx, options.Filter)
	if err != nil {
		return nil, err
	}
	options.Filter = filter
	return gen.QueryStoresHandler(ctx, resolver, options)
}
