package authorization

import (
	"context"
	"errors"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

// confirmFranchiseInitialAccount blocks the generated update path; evidence is required through the dedicated HTTP command.
func confirmFranchiseInitialAccount(ctx context.Context, resolver *gen.GeneratedResolver, organizationID string, input map[string]interface{}) (*gen.Organization, error) {
	return nil, auth.NewError(auth.CodePermissionDenied)
}

func authorizeInitialAccountWriter(principal *auth.WorkspacePrincipal) error {
	if principal == nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	if principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	if err := Authorize(principal, Intent{Action: "hqMembership:read", Mode: AccessRead}); err != nil {
		return err
	}
	return Authorize(principal, Intent{Action: "account:update", Mode: AccessUpdate})
}

func verifyInitialAccountCandidate(db *gorm.DB, organizationID, accountID string) error {
	var account gen.Account
	if err := activeRecordCondition(db.Session(&gorm.Session{NewDB: true})).First(&account, "id = ? AND status = ?", accountID, gen.AccountStatusActive).Error; err != nil {
		return concealInitialAccountWriteError(err)
	}
	var membership gen.OperatorMembership
	if err := activeRecordCondition(db.Session(&gorm.Session{NewDB: true})).First(&membership, "organization_id = ? AND account_id = ? AND status = ?", organizationID, accountID, gen.MembershipStatusActive).Error; err != nil {
		return concealInitialAccountWriteError(err)
	}
	var hqCount int64
	hqOrganizations := db.Session(&gorm.Session{NewDB: true}).Model(&gen.Organization{}).Select("id").Where("type = ?", gen.OrganizationTypeHeadquarters)
	err := db.Session(&gorm.Session{NewDB: true}).Model(&gen.OperatorMembership{}).
		Where("account_id = ? AND organization_id IN (?)", accountID, hqOrganizations).
		Where("is_delete IS NULL OR is_delete = ?", 1).
		Count(&hqCount).Error
	if err != nil {
		return err
	}
	if hqCount != 0 {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return nil
}

func concealInitialAccountWriteError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return err
}
