/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authentication

import (
	"context"
	"errors"
	"time"

	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"base-engine/src/services/audit"
	sessionservice "base-engine/src/services/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var dummyPasswordHash = mustDummyPasswordHash()

// LoginResult 返回登录发现区所需的账号、权威会话与签名令牌。
type LoginResult struct {
	Account *gen.Account
	Session *gen.Session
	Token   string
}

// Service 实现常量成本密码校验、锁定和会话创建。
type Service struct {
	db        *gorm.DB
	config    config.SecurityConfig
	sessions  *sessionservice.Service
	limiter   *LoginLimiter
	audit     *audit.Service
	publisher sessionservice.Publisher
}

// NewService 创建认证服务。
func NewService(db *gorm.DB, cfg config.SecurityConfig, sessions *sessionservice.Service, limiter *LoginLimiter, auditService *audit.Service, publisher sessionservice.Publisher) *Service {
	return &Service{db: db, config: cfg, sessions: sessions, limiter: limiter, audit: auditService, publisher: publisher}
}

// Login 校验凭据并在同一事务中创建发现区会话。
func (s *Service) Login(ctx context.Context, phone, password, ip string, now time.Time) (*LoginResult, error) {
	if !s.limiter.Allow(loginLimitKey(ip, phone), now) {
		return nil, auth.NewError(auth.CodeRateLimited)
	}
	var result *LoginResult
	var outcomeErr error
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		loginResult, loginErr := s.loginTransaction(ctx, tx, phone, password, now)
		result = loginResult
		if auth.ErrorCode(loginErr) != "" {
			outcomeErr = loginErr
			return nil
		}
		return loginErr
	})
	if err != nil {
		return nil, err
	}
	return result, outcomeErr
}

func (s *Service) loginTransaction(ctx context.Context, tx *gorm.DB, phone, password string, now time.Time) (*LoginResult, error) {
	account, credential, found, err := loadLoginRecords(tx, phone)
	if err != nil {
		return nil, err
	}
	if !found {
		_ = bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(password))
		code := auth.CodeAccountUnavailable
		if account.ID == "" {
			code = auth.CodePhoneNotFound
		}
		if err := s.writeAuthenticationAudit(tx, account, nil, string(code)); err != nil {
			return nil, err
		}
		return nil, auth.NewError(code)
	}
	if err := s.rejectLockedCredential(tx, account, credential, now); err != nil {
		return nil, err
	}
	passwordErr := auth.VerifyPassword(credential.PasswordHash, password)
	if account.Status != gen.AccountStatusActive {
		return nil, s.recordAuditedLoginFailure(tx, account, credential, now, auth.CodeAccountInactive)
	}
	if passwordErr != nil || temporaryPasswordExpired(account, credential, now) {
		return nil, s.recordAuditedLoginFailure(tx, account, credential, now, auth.CodePasswordIncorrect)
	}
	return s.completeLogin(ctx, tx, account, credential, now)
}

func (s *Service) completeLogin(ctx context.Context, tx *gorm.DB, account *gen.Account, credential *AccountCredential, now time.Time) (*LoginResult, error) {
	if err := resetLoginFailures(tx, credential); err != nil {
		return nil, err
	}
	created, err := s.sessions.CreateDiscovery(ctx, tx, account, now)
	if err != nil {
		return nil, err
	}
	token, err := s.signSession(account, created)
	if err != nil {
		return nil, err
	}
	if err := s.writeAuthenticationAudit(tx, account, &created.ID, "SUCCESS"); err != nil {
		return nil, err
	}
	return &LoginResult{Account: account, Session: created, Token: token}, nil
}

func (s *Service) rejectLockedCredential(tx *gorm.DB, account *gen.Account, credential *AccountCredential, now time.Time) error {
	if credential.LockedUntil == nil || !credential.LockedUntil.After(now) {
		return nil
	}
	if err := s.writeAuthenticationAudit(tx, account, nil, string(auth.CodeAccountLocked)); err != nil {
		return err
	}
	return auth.NewError(auth.CodeAccountLocked)
}

func (s *Service) recordAuditedLoginFailure(tx *gorm.DB, account *gen.Account, credential *AccountCredential, now time.Time, code auth.Code) error {
	loginErr := s.recordLoginFailure(tx, credential, now, code)
	if err := s.writeAuthenticationAudit(tx, account, nil, string(code)); err != nil {
		return err
	}
	return loginErr
}

func temporaryPasswordExpired(account *gen.Account, credential *AccountCredential, now time.Time) bool {
	if !account.MustChangePassword {
		return false
	}
	return credential.TemporaryPasswordExpiresAt == nil || !credential.TemporaryPasswordExpiresAt.After(now)
}

func loadLoginRecords(tx *gorm.DB, phone string) (*gen.Account, *AccountCredential, bool, error) {
	account := &gen.Account{}
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("phone = ? AND (is_delete IS NULL OR is_delete = ?)", phone, 1).
		First(account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return account, &AccountCredential{}, false, nil
	}
	if err != nil {
		return nil, nil, false, err
	}
	credential := &AccountCredential{}
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(credential, "account_id = ?", account.ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return account, credential, false, nil
	}
	return account, credential, err == nil, err
}

func (s *Service) recordLoginFailure(tx *gorm.DB, credential *AccountCredential, now time.Time, code auth.Code) error {
	failed := credential.FailedLoginCount + 1
	updates := map[string]any{"failed_login_count": failed, "updated_at": now}
	if failed >= s.config.LoginLockThreshold {
		updates["locked_until"] = now.Add(s.config.LoginLockDuration)
	}
	if err := tx.Model(credential).Updates(updates).Error; err != nil {
		return err
	}
	return auth.NewError(code)
}

func resetLoginFailures(tx *gorm.DB, credential *AccountCredential) error {
	return tx.Model(credential).Updates(map[string]any{
		"failed_login_count": 0, "locked_until": nil, "updated_at": time.Now(),
	}).Error
}

func (s *Service) signSession(account *gen.Account, session *gen.Session) (string, error) {
	claims := auth.SessionClaims{
		SessionID: session.ID, AccountID: account.ID,
		WorkspaceType:  auth.WorkspaceType(session.WorkspaceType),
		OrganizationID: session.OrganizationID, CredentialVersion: account.CredentialVersion,
	}
	return auth.SignSessionClaims(s.config, claims)
}

func mustDummyPasswordHash() []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte("Dummy-Password-42!"), 12)
	if err != nil {
		panic(err)
	}
	return hash
}
