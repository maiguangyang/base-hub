package src

import (
	"context"
	"errors"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ResetFranchiseInitialPassword 只轮换开通加盟商时记录的账号密码。
func (r *MutationResolver) ResetFranchiseInitialPassword(ctx context.Context, organizationID string) (*gen.TemporaryPasswordPayload, error) {
	principal, err := requireResolverPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return nil, auth.NewError(auth.CodeWorkspaceForbidden)
	}
	if err := authorization.Authorize(principal, authorization.Intent{Action: "account:update", Mode: authorization.AccessUpdate}); err != nil {
		return nil, err
	}
	return resetFranchiseInitialPassword(ctx, r.Services, principal, organizationID)
}

func resetFranchiseInitialPassword(ctx context.Context, services Dependencies, principal *auth.WorkspacePrincipal, organizationID string) (*gen.TemporaryPasswordPayload, error) {
	accountID, err := franchiseInitialAccountTarget(services.DB.WithContext(ctx), organizationID)
	if err != nil {
		return nil, err
	}
	password, err := auth.GenerateTemporaryPassword()
	if err != nil {
		return nil, err
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	var revoked []string
	err = services.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		lockedID, txErr := franchiseInitialAccountTarget(tx.Clauses(clause.Locking{Strength: "UPDATE"}), organizationID)
		if txErr != nil {
			return txErr
		}
		if lockedID != accountID {
			return auth.NewError(auth.CodeConflict)
		}
		revoked, txErr = resetPasswordTransaction(ctx, services, tx, accountID, hash, now)
		if txErr != nil {
			return txErr
		}
		return auditFranchiseInitialPasswordReset(services, tx, principal, organizationID, accountID)
	})
	if err != nil {
		return nil, err
	}
	publishCredentialReset(services, revoked, now)
	return &gen.TemporaryPasswordPayload{AccountID: accountID, TemporaryPassword: password}, nil
}

func franchiseInitialAccountTarget(tx *gorm.DB, organizationID string) (string, error) {
	var organization gen.Organization
	err := tx.Where("is_delete IS NULL OR is_delete = ?", 1).First(&organization, "id = ? AND type = ?", organizationID, gen.OrganizationTypeFranchise).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", auth.NewError(auth.CodePermissionDenied)
		}
		return "", err
	}
	if organization.InitialAccountID == nil {
		return "", auth.NewError(auth.CodePermissionDenied)
	}
	accountID := *organization.InitialAccountID
	if err := verifyOpeningRecordAccount(tx, organizationID, accountID); err != nil {
		return "", err
	}
	var account gen.Account
	err = tx.Session(&gorm.Session{NewDB: true}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("is_delete IS NULL OR is_delete = ?", 1).First(&account, "id = ? AND status = ?", accountID, gen.AccountStatusActive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", auth.NewError(auth.CodePermissionDenied)
		}
		return "", err
	}
	var membership gen.OperatorMembership
	err = tx.Session(&gorm.Session{NewDB: true}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("is_delete IS NULL OR is_delete = ?", 1).
		First(&membership, "organization_id = ? AND account_id = ? AND status = ?", organizationID, accountID, gen.MembershipStatusActive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", auth.NewError(auth.CodeMembershipInactive)
		}
		return "", err
	}
	if err := verifyResetTargetPolicy(tx, organizationID, accountID); err != nil {
		return "", err
	}
	return accountID, nil
}

func verifyOpeningRecordAccount(tx *gorm.DB, organizationID, accountID string) error {
	var opening gen.FranchiseOpeningRecord
	result := tx.Session(&gorm.Session{NewDB: true}).Where("organization_id = ?", organizationID).Limit(1).Find(&opening)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 && opening.InitialAccountID != accountID {
		return auth.NewError(auth.CodeConflict)
	}
	return nil
}

func verifyResetTargetPolicy(tx *gorm.DB, organizationID, accountID string) error {
	var hqCount int64
	hqOrganizations := tx.Session(&gorm.Session{NewDB: true}).Model(&gen.Organization{}).Select("id").Where("type = ?", gen.OrganizationTypeHeadquarters)
	err := tx.Session(&gorm.Session{NewDB: true}).Model(&gen.OperatorMembership{}).
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

func auditFranchiseInitialPasswordReset(services Dependencies, tx *gorm.DB, principal *auth.WorkspacePrincipal, organizationID, accountID string) error {
	if services.Audit == nil {
		return nil
	}
	return services.Audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &organizationID, Action: "franchiseInitialAccount:reset_password",
		ResourceType: "account", ResourceID: accountID, ResultCode: "SUCCESS",
		Metadata: audit.MetadataForPrincipal(principal, audit.Metadata{Source: "headquarters"}),
	})
}
