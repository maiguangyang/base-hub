package src

import (
	"context"

	"base-engine/gen"
)

func (r *MutationResolver) SetStoreDocument(ctx context.Context, storeID string, kind gen.StoreDocumentKind, attachmentID string) (*gen.Store, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	return r.Services.Stores.SetDocument(ctx, principal, storeID, kind, attachmentID)
}

func (r *MutationResolver) RemoveStoreDocument(ctx context.Context, storeID string, kind gen.StoreDocumentKind) (*gen.Store, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	return r.Services.Stores.RemoveDocument(ctx, principal, storeID, kind)
}
