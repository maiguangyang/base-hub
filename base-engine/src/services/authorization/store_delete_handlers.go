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

func deleteStores(ctx context.Context, resolver *gen.GeneratedResolver, ids []string, _ *bool) (bool, error) {
	stores, err := loadStoresForWrite(ctx, resolver, ids, "delete", AccessDelete)
	if err != nil {
		return false, err
	}
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return false, err
	}
	if err := validateStoreDeletionLifecycle(principal, stores); err != nil {
		return false, err
	}
	scoped := false
	done, err := gen.DeleteStoresHandler(ctx, resolver, ids, &scoped)
	if err != nil || !done {
		return done, err
	}
	return done, auditStoreWrites(ctx, resolver, stores, "delete")
}

func recoverStores(context.Context, *gen.GeneratedResolver, []string) (bool, error) {
	return false, auth.NewError(auth.CodePermissionDenied)
}

func validateStoreDeletionLifecycle(principal *auth.WorkspacePrincipal, stores []*gen.Store) error {
	if principal.WorkspaceType != auth.WorkspaceTypeFranchise {
		return nil
	}
	for _, store := range stores {
		if store.Lifecycle != gen.StoreLifecycleDraft && store.Lifecycle != gen.StoreLifecycleRejected {
			return auth.NewError(auth.CodeConflict)
		}
	}
	return nil
}
