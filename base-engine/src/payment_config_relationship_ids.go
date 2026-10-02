package src

import (
	"context"

	"base-engine/gen"
)

// Generated relationship-ID projections bypass ResolutionHandlers; block them here.
func (r *OrganizationResolver) PaymentConfigsIds(ctx context.Context, _ *gen.Organization) ([]string, error) {
	return nil, denyOpeningRecordProjection(ctx)
}

func (r *StoreResolver) PaymentConfigsIds(ctx context.Context, _ *gen.Store) ([]string, error) {
	return nil, denyOpeningRecordProjection(ctx)
}
