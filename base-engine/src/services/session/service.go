/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package session

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

const RevocationCodeAuthorityChanged = "AUTHORITY_CHANGED"

// Service 管理数据库权威会话的生命周期。
type Service struct {
	duration time.Duration
}

// NewService 创建会话服务。
func NewService(duration time.Duration) *Service {
	return &Service{duration: duration}
}

// CreateDiscovery 创建尚未选择组织的登录会话。
func (s *Service) CreateDiscovery(ctx context.Context, tx *gorm.DB, account *gen.Account, now time.Time) (*gen.Session, error) {
	session := &gen.Session{
		ID: uuid.Must(uuid.NewV4()).String(), AccountID: account.ID,
		WorkspaceType:     gen.WorkspaceTypeDiscovery,
		CredentialVersion: account.CredentialVersion,
		ExpiresAt:         now.Add(s.duration), LastSeenAt: now,
	}
	return session, tx.WithContext(ctx).Create(session).Error
}

// Revoke 撤销一个仍有效的权威会话。
func (s *Service) Revoke(ctx context.Context, tx *gorm.DB, sessionID, code string, now time.Time) error {
	result := tx.WithContext(ctx).Model(&gen.Session{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Updates(map[string]any{"revoked_at": now, "revocation_code": code})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return auth.NewError(auth.CodeSessionRevoked)
	}
	return nil
}

// RevokeOrganizationSessions 撤销组织内所有活跃会话并返回精确 ID 集合。
func (s *Service) RevokeOrganizationSessions(ctx context.Context, tx *gorm.DB, organizationID, code string, now time.Time) ([]string, error) {
	return s.revokeMatching(ctx, tx, "organization_id = ?", organizationID, code, now)
}

// RevokeAccountSessions 撤销账号的全部活跃会话并返回精确 ID 集合。
func (s *Service) RevokeAccountSessions(ctx context.Context, tx *gorm.DB, accountID, code string, now time.Time) ([]string, error) {
	return s.revokeMatching(ctx, tx, "account_id = ?", accountID, code, now)
}

// RevokeWorkspaceSessions 撤销账号在指定组织工作区内的活跃会话。
func (s *Service) RevokeWorkspaceSessions(ctx context.Context, tx *gorm.DB, accountID, organizationID, code string, now time.Time) ([]string, error) {
	var sessionIDs []string
	query := tx.WithContext(ctx).Model(&gen.Session{}).
		Where("account_id = ? AND organization_id = ?", accountID, organizationID).
		Where("revoked_at IS NULL AND expires_at > ?", now)
	if err := query.Pluck("id", &sessionIDs).Error; err != nil || len(sessionIDs) == 0 {
		return sessionIDs, err
	}
	err := tx.WithContext(ctx).Model(&gen.Session{}).Where("id IN ?", sessionIDs).
		Updates(map[string]any{"revoked_at": now, "revocation_code": code}).Error
	return sessionIDs, err
}

func (s *Service) revokeMatching(ctx context.Context, tx *gorm.DB, predicate string, value any, code string, now time.Time) ([]string, error) {
	var sessionIDs []string
	query := tx.WithContext(ctx).Model(&gen.Session{}).
		Where(predicate+" AND revoked_at IS NULL AND expires_at > ?", value, now)
	if err := query.Pluck("id", &sessionIDs).Error; err != nil {
		return nil, err
	}
	if len(sessionIDs) == 0 {
		return sessionIDs, nil
	}
	err := tx.WithContext(ctx).Model(&gen.Session{}).Where("id IN ?", sessionIDs).
		Updates(map[string]any{"revoked_at": now, "revocation_code": code}).Error
	return sessionIDs, err
}
