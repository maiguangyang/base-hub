package authorization

import (
	"context"

	"base-engine/gen"
)

// Payment configuration entities are persisted through the dedicated HQ API only.
func registerPaymentConfigHandlers(h *gen.ResolutionHandlers) {
	registerGlobalPaymentConfigHandlers(h)
	registerFranchisePaymentConfigHandlers(h)
	registerStorePaymentConfigHandlers(h)
}

func registerGlobalPaymentConfigHandlers(h *gen.ResolutionHandlers) {
	h.CreateGlobalPaymentConfig = func(ctx context.Context, _ *gen.GeneratedResolver, _ map[string]interface{}) (*gen.GlobalPaymentConfig, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.UpdateGlobalPaymentConfig = func(ctx context.Context, _ *gen.GeneratedResolver, _ string, _ map[string]interface{}) (*gen.GlobalPaymentConfig, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.DeleteGlobalPaymentConfigs = denyPaymentConfigDelete
	h.RecoveryGlobalPaymentConfigs = denyPaymentConfigRecovery
	h.QueryGlobalPaymentConfig = func(ctx context.Context, _ *gen.GeneratedResolver, _ gen.QueryGlobalPaymentConfigHandlerOptions) (*gen.GlobalPaymentConfig, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.QueryGlobalPaymentConfigs = func(ctx context.Context, _ *gen.GeneratedResolver, _ gen.QueryGlobalPaymentConfigsHandlerOptions) (*gen.GlobalPaymentConfigResultType, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
}

func registerFranchisePaymentConfigHandlers(h *gen.ResolutionHandlers) {
	h.CreateFranchisePaymentConfig = func(ctx context.Context, _ *gen.GeneratedResolver, _ map[string]interface{}) (*gen.FranchisePaymentConfig, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.UpdateFranchisePaymentConfig = func(ctx context.Context, _ *gen.GeneratedResolver, _ string, _ map[string]interface{}) (*gen.FranchisePaymentConfig, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.DeleteFranchisePaymentConfigs = denyPaymentConfigDelete
	h.RecoveryFranchisePaymentConfigs = denyPaymentConfigRecovery
	h.QueryFranchisePaymentConfig = func(ctx context.Context, _ *gen.GeneratedResolver, _ gen.QueryFranchisePaymentConfigHandlerOptions) (*gen.FranchisePaymentConfig, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.QueryFranchisePaymentConfigs = func(ctx context.Context, _ *gen.GeneratedResolver, _ gen.QueryFranchisePaymentConfigsHandlerOptions) (*gen.FranchisePaymentConfigResultType, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.FranchisePaymentConfigOrganization = func(ctx context.Context, _ *gen.GeneratedResolver, _ *gen.FranchisePaymentConfig) (*gen.Organization, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.OrganizationPaymentConfigs = func(ctx context.Context, _ *gen.GeneratedResolver, _ *gen.Organization) ([]*gen.FranchisePaymentConfig, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
}

func registerStorePaymentConfigHandlers(h *gen.ResolutionHandlers) {
	h.CreateStorePaymentConfig = func(ctx context.Context, _ *gen.GeneratedResolver, _ map[string]interface{}) (*gen.StorePaymentConfig, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.UpdateStorePaymentConfig = func(ctx context.Context, _ *gen.GeneratedResolver, _ string, _ map[string]interface{}) (*gen.StorePaymentConfig, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.DeleteStorePaymentConfigs = denyPaymentConfigDelete
	h.RecoveryStorePaymentConfigs = denyPaymentConfigRecovery
	h.QueryStorePaymentConfig = func(ctx context.Context, _ *gen.GeneratedResolver, _ gen.QueryStorePaymentConfigHandlerOptions) (*gen.StorePaymentConfig, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.QueryStorePaymentConfigs = func(ctx context.Context, _ *gen.GeneratedResolver, _ gen.QueryStorePaymentConfigsHandlerOptions) (*gen.StorePaymentConfigResultType, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.StorePaymentConfigStore = func(ctx context.Context, _ *gen.GeneratedResolver, _ *gen.StorePaymentConfig) (*gen.Store, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
	h.StorePaymentConfigs = func(ctx context.Context, _ *gen.GeneratedResolver, _ *gen.Store) ([]*gen.StorePaymentConfig, error) {
		return nil, paymentConfigAccessDenied(ctx)
	}
}

func paymentConfigAccessDenied(ctx context.Context) error { return generatedDataAccessDenied(ctx) }

func denyPaymentConfigDelete(ctx context.Context, _ *gen.GeneratedResolver, _ []string, _ *bool) (bool, error) {
	return false, paymentConfigAccessDenied(ctx)
}

func denyPaymentConfigRecovery(ctx context.Context, _ *gen.GeneratedResolver, _ []string) (bool, error) {
	return false, paymentConfigAccessDenied(ctx)
}
