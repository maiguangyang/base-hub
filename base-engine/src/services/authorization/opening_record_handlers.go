package authorization

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
)

func registerOpeningRecordHandlers(handlers *gen.ResolutionHandlers) {
	handlers.CreateFranchiseOpeningRecord = denyCreateOpeningRecord
	handlers.UpdateFranchiseOpeningRecord = denyUpdateOpeningRecord
	handlers.DeleteFranchiseOpeningRecords = denyDeleteOpeningRecords
	handlers.RecoveryFranchiseOpeningRecords = denyRecoveryOpeningRecords
	handlers.QueryFranchiseOpeningRecord = denyQueryOpeningRecord
	handlers.QueryFranchiseOpeningRecords = denyQueryOpeningRecords
	handlers.AccountOpeningRecords = denyAccountOpeningRecords
	handlers.AccountRecordedOpeningRecords = denyAccountOpeningRecords
	handlers.OrganizationOpeningRecords = denyOrganizationOpeningRecords
	handlers.FranchiseOpeningRecordOrganization = denyOpeningRecordOrganization
	handlers.FranchiseOpeningRecordInitialAccount = denyOpeningRecordAccount
	handlers.FranchiseOpeningRecordRecordedByAccount = denyOpeningRecordAccount
}

func generatedDataAccessDenied(ctx context.Context) error {
	if _, err := requestPrincipal(ctx); err != nil {
		return err
	}
	return auth.NewError(auth.CodePermissionDenied)
}

func denyCreateOpeningRecord(ctx context.Context, _ *gen.GeneratedResolver, _ map[string]interface{}) (*gen.FranchiseOpeningRecord, error) {
	return nil, generatedDataAccessDenied(ctx)
}

func denyUpdateOpeningRecord(ctx context.Context, _ *gen.GeneratedResolver, _ string, _ map[string]interface{}) (*gen.FranchiseOpeningRecord, error) {
	return nil, generatedDataAccessDenied(ctx)
}

func denyDeleteOpeningRecords(ctx context.Context, _ *gen.GeneratedResolver, _ []string, _ *bool) (bool, error) {
	return false, generatedDataAccessDenied(ctx)
}

func denyRecoveryOpeningRecords(ctx context.Context, _ *gen.GeneratedResolver, _ []string) (bool, error) {
	return false, generatedDataAccessDenied(ctx)
}

func denyQueryOpeningRecord(ctx context.Context, _ *gen.GeneratedResolver, _ gen.QueryFranchiseOpeningRecordHandlerOptions) (*gen.FranchiseOpeningRecord, error) {
	return nil, generatedDataAccessDenied(ctx)
}

func denyQueryOpeningRecords(ctx context.Context, _ *gen.GeneratedResolver, _ gen.QueryFranchiseOpeningRecordsHandlerOptions) (*gen.FranchiseOpeningRecordResultType, error) {
	return nil, generatedDataAccessDenied(ctx)
}

func denyAccountOpeningRecords(ctx context.Context, _ *gen.GeneratedResolver, _ *gen.Account) ([]*gen.FranchiseOpeningRecord, error) {
	return nil, generatedDataAccessDenied(ctx)
}

func denyOrganizationOpeningRecords(ctx context.Context, _ *gen.GeneratedResolver, _ *gen.Organization) ([]*gen.FranchiseOpeningRecord, error) {
	return nil, generatedDataAccessDenied(ctx)
}

func denyOpeningRecordOrganization(ctx context.Context, _ *gen.GeneratedResolver, _ *gen.FranchiseOpeningRecord) (*gen.Organization, error) {
	return nil, generatedDataAccessDenied(ctx)
}

func denyOpeningRecordAccount(ctx context.Context, _ *gen.GeneratedResolver, _ *gen.FranchiseOpeningRecord) (*gen.Account, error) {
	return nil, generatedDataAccessDenied(ctx)
}
