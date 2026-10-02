/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package src

import (
	"context"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authentication"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func resetTemporaryPassword(ctx context.Context, services Dependencies, principal *auth.WorkspacePrincipal, accountID string) (string, error) {
	return resetTemporaryPasswordWithGenerator(ctx, services, principal, accountID, auth.GenerateTemporaryPassword)
}

func resetTemporaryPasswordWithGenerator(ctx context.Context, services Dependencies, principal *auth.WorkspacePrincipal, accountID string, generate func() (string, error)) (string, error) {
	if err := authorizePasswordResetTarget(services.DB.WithContext(ctx), principal, accountID); err != nil {
		auth.LogAuthorizationDenied(principal, "hqAdministrator:reset_password", "account", accountID, err)
		return "", err
	}
	password, err := generate()
	if err != nil {
		return "", err
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return "", err
	}
	var revoked []string
	now := time.Now()
	err = services.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if txErr := authorizePasswordResetTarget(tx, principal, accountID); txErr != nil {
			auth.LogAuthorizationDenied(principal, "hqAdministrator:reset_password", "account", accountID, txErr)
			return txErr
		}
		ids, txErr := resetPasswordTransaction(ctx, services, tx, accountID, hash, now)
		revoked = ids
		if txErr != nil {
			return txErr
		}
		return auditPasswordReset(services, tx, principal, accountID)
	})
	if err != nil {
		return "", err
	}
	publishCredentialReset(services, revoked, now)
	return password, nil
}

func authorizePasswordResetTarget(tx *gorm.DB, principal *auth.WorkspacePrincipal, accountID string) error {
	if principal == nil || principal.OrganizationID == nil {
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	membership := &gen.OperatorMembership{}
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("is_delete IS NULL OR is_delete = ?", 1).
		First(membership, "account_id = ? AND organization_id = ? AND status = ?", accountID, *principal.OrganizationID, gen.MembershipStatusActive).Error
	if err != nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return authorization.ValidateHQMembershipTarget(tx, principal, membership)
}

func resetPasswordTransaction(ctx context.Context, services Dependencies, tx *gorm.DB, accountID, hash string, now time.Time) ([]string, error) {
	account := &gen.Account{}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("is_delete IS NULL OR is_delete = ?", 1).First(account, "id = ?", accountID).Error; err != nil {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	account.CredentialVersion++
	if err := tx.Model(account).Updates(map[string]any{"credential_version": account.CredentialVersion, "must_change_password": true}).Error; err != nil {
		return nil, err
	}
	updates := map[string]any{
		"password_hash": hash, "password_changed_at": now,
		"temporary_password_expires_at": authentication.NewTemporaryPasswordExpiry(now),
		"failed_login_count":            0, "locked_until": nil, "updated_at": now,
	}
	result := tx.Model(&authentication.AccountCredential{}).Where("account_id = ?", accountID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	return services.Sessions.RevokeAccountSessions(ctx, tx, accountID, "CREDENTIALS_CHANGED", now)
}

func auditPasswordReset(services Dependencies, tx *gorm.DB, principal *auth.WorkspacePrincipal, accountID string) error {
	if services.Audit == nil {
		return nil
	}
	metadata := audit.MetadataForPrincipal(principal, audit.Metadata{Source: "headquarters"})
	return services.Audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: principal.OrganizationID,
		Action:         "hqAdministrator:reset_password", ResourceType: "account", ResourceID: accountID,
		ResultCode: "SUCCESS", Metadata: metadata,
	})
}

func publishCredentialReset(services Dependencies, sessionIDs []string, occurredAt time.Time) {
	if services.Publisher == nil {
		return
	}
	for _, sessionID := range sessionIDs {
		services.Publisher.PublishSession(sessionID, &gen.SessionEvent{
			Code: gen.SessionEventCodeCredentialsChanged, SessionID: sessionID, OccurredAt: occurredAt,
		})
	}
}
