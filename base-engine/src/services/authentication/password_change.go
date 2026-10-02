/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authentication

import (
	"context"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

const credentialsChangedCode = "CREDENTIALS_CHANGED"

// ChangePassword 更换密码、撤销旧会话并创建替代发现区会话。
func (s *Service) ChangePassword(ctx context.Context, principal *auth.WorkspacePrincipal, currentPassword, newPassword string, now time.Time) (*LoginResult, error) {
	if principal == nil {
		return nil, auth.NewError(auth.CodeAuthRequired)
	}
	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return nil, auth.NewError(auth.CodePasswordWeak)
	}
	var result *LoginResult
	var revokedSessionIDs []string
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		changed, revoked, txErr := s.changePasswordTransaction(ctx, tx, principal, currentPassword, hash, now)
		result, revokedSessionIDs = changed, revoked
		return txErr
	})
	if err != nil {
		return nil, err
	}
	s.publishCredentialChanges(revokedSessionIDs, principal.SessionID, now)
	return result, nil
}

func (s *Service) changePasswordTransaction(ctx context.Context, tx *gorm.DB, principal *auth.WorkspacePrincipal, currentPassword, hash string, now time.Time) (*LoginResult, []string, error) {
	account := &gen.Account{}
	if err := tx.First(account, "id = ?", principal.AccountID).Error; err != nil {
		return nil, nil, auth.NewError(auth.CodeAuthRequired)
	}
	credential := &AccountCredential{}
	if err := tx.First(credential, "account_id = ?", account.ID).Error; err != nil {
		return nil, nil, auth.NewError(auth.CodeInvalidCredentials)
	}
	if auth.VerifyPassword(credential.PasswordHash, currentPassword) != nil {
		return nil, nil, auth.NewError(auth.CodeInvalidCredentials)
	}
	account.CredentialVersion++
	if err := tx.Model(account).Updates(map[string]any{"credential_version": account.CredentialVersion, "must_change_password": false}).Error; err != nil {
		return nil, nil, err
	}
	if err := tx.Model(credential).Updates(map[string]any{"password_hash": hash, "password_changed_at": now, "temporary_password_expires_at": nil, "failed_login_count": 0, "locked_until": nil, "updated_at": now}).Error; err != nil {
		return nil, nil, err
	}
	revoked, err := s.sessions.RevokeAccountSessions(ctx, tx, account.ID, credentialsChangedCode, now)
	if err != nil {
		return nil, nil, err
	}
	replacement, err := s.sessions.CreateDiscovery(ctx, tx, account, now)
	if err != nil {
		return nil, nil, err
	}
	token, err := s.signSession(account, replacement)
	if err == nil {
		err = s.writePasswordChangeAudit(tx, principal.SessionID, account)
	}
	return &LoginResult{Account: account, Session: replacement, Token: token}, revoked, err
}

func (s *Service) publishCredentialChanges(sessionIDs []string, excludedSessionID string, occurredAt time.Time) {
	if s.publisher == nil {
		return
	}
	for _, sessionID := range sessionIDs {
		if sessionID == excludedSessionID {
			continue
		}
		s.publisher.PublishSession(sessionID, &gen.SessionEvent{
			Code: gen.SessionEventCodeCredentialsChanged, SessionID: sessionID, OccurredAt: occurredAt,
		})
	}
}
