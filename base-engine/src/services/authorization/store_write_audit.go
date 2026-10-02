/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
)

func loadStoresForWrite(ctx context.Context, resolver *gen.GeneratedResolver, ids []string, verb string, mode AccessMode) ([]*gen.Store, error) {
	stores := make([]*gen.Store, 0, len(ids))
	for _, id := range ids {
		store, _, err := loadAuthorizedStore(ctx, resolver, id, verb, mode)
		if err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}
	return stores, nil
}

func auditStoreWrite(ctx context.Context, resolver *gen.GeneratedResolver, principal *auth.WorkspacePrincipal, store *gen.Store, verb string) error {
	return writeGeneratedAudit(ctx, resolver, principal, store.OrganizationID, storeAction(principal, verb), "store", store.ID)
}

func auditStoreWrites(ctx context.Context, resolver *gen.GeneratedResolver, stores []*gen.Store, verb string) error {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return err
	}
	for _, store := range stores {
		if err := auditStoreWrite(ctx, resolver, principal, store, verb); err != nil {
			return err
		}
	}
	return nil
}
