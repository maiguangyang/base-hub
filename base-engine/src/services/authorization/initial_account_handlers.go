package authorization

import (
	"context"
	"errors"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func authorizeInitialAccountRead(ctx context.Context) error {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return err
	}
	if principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	return Authorize(principal, Intent{Action: "account:update", Mode: AccessRead})
}

func organizationInitialAccount(ctx context.Context, resolver *gen.GeneratedResolver, organization *gen.Organization) (*gen.Account, error) {
	if err := authorizeInitialAccountRead(ctx); err != nil {
		return nil, err
	}
	if organization.InitialAccountID == nil {
		return nil, nil
	}
	var account gen.Account
	err := activeRecordCondition(resolverDB(ctx, resolver)).First(&account, "id = ?", *organization.InitialAccountID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &account, err
}

func accountInitializedOrganizations(ctx context.Context, resolver *gen.GeneratedResolver, account *gen.Account) ([]*gen.Organization, error) {
	if err := authorizeInitialAccountRead(ctx); err != nil {
		return nil, err
	}
	var organizations []*gen.Organization
	err := activeRecordCondition(resolverDB(ctx, resolver)).Where("initial_account_id = ?", account.ID).Find(&organizations).Error
	return organizations, err
}
