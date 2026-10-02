/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authentication

import (
	"context"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"base-engine/src/services/audit"
	sessionservice "base-engine/src/services/session"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestLoginCreatesDiscoverySession 验证账号无需已有工作空间即可登录发现区。
func TestLoginCreatesDiscoverySession(t *testing.T) {
	service, db, cfg := newAuthenticationFixture(t)
	seedLoginAccount(t, db, "account-1", "13800000000", gen.AccountStatusActive)
	result, err := service.Login(context.Background(), "13800000000", "Correct-Horse-42", "127.0.0.1", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if result.Session.WorkspaceType != gen.WorkspaceTypeDiscovery || result.Session.OrganizationID != nil {
		t.Fatalf("invalid discovery session: %#v", result.Session)
	}
	claims, err := auth.ParseSessionClaims(cfg, result.Token)
	if err != nil || claims.SessionID != result.Session.ID || claims.AccountID != result.Account.ID {
		t.Fatalf("invalid signed claims: %#v, %v", claims, err)
	}
	assertAuthenticationAudit(t, db, "authentication:login", "SUCCESS", 1)
}

// TestLoginLockout 验证第五次失败触发锁定且正确密码也不能绕过。
func TestLoginLockout(t *testing.T) {
	service, db, _ := newAuthenticationFixture(t)
	seedLoginAccount(t, db, "account-1", "13800000000", gen.AccountStatusActive)
	now := time.Now()
	for attempt := 0; attempt < 5; attempt++ {
		if _, err := service.Login(context.Background(), "13800000000", "Wrong-Horse-42", "127.0.0.1", now); auth.ErrorCode(err) != auth.CodePasswordIncorrect {
			t.Fatalf("attempt %d error = %v", attempt+1, err)
		}
	}
	var credential AccountCredential
	if err := db.First(&credential, "account_id = ?", "account-1").Error; err != nil {
		t.Fatal(err)
	}
	if credential.LockedUntil == nil || credential.LockedUntil.Before(now.Add(14*time.Minute)) {
		t.Fatalf("account was not locked: %#v", credential)
	}
	if _, err := service.Login(context.Background(), "13800000000", "Correct-Horse-42", "127.0.0.1", now); auth.ErrorCode(err) != auth.CodeAccountLocked {
		t.Fatalf("locked login error = %v", err)
	}
	assertAuthenticationAudit(t, db, "authentication:login", string(auth.CodePasswordIncorrect), 5)
	assertAuthenticationAudit(t, db, "authentication:login", string(auth.CodeAccountLocked), 1)
}

// TestLoginResetsFailuresAndDistinguishesMissingPhone 验证成功清零且未知手机号有独立错误码。
func TestLoginResetsFailuresAndDistinguishesMissingPhone(t *testing.T) {
	service, db, _ := newAuthenticationFixture(t)
	seedLoginAccount(t, db, "account-1", "13800000000", gen.AccountStatusActive)
	db.Model(&AccountCredential{}).Where("account_id = ?", "account-1").Update("failed_login_count", 3)
	if _, err := service.Login(context.Background(), "13800000000", "Correct-Horse-42", "10.0.0.1", time.Now()); err != nil {
		t.Fatal(err)
	}
	var credential AccountCredential
	db.First(&credential, "account_id = ?", "account-1")
	if credential.FailedLoginCount != 0 || credential.LockedUntil != nil {
		t.Fatalf("failures were not reset: %#v", credential)
	}
	_, unknownErr := service.Login(context.Background(), "13900000000", "Correct-Horse-42", "10.0.0.2", time.Now())
	db.Model(&gen.Account{}).Where("id = ?", "account-1").Update("status", gen.AccountStatusDisabled)
	_, disabledErr := service.Login(context.Background(), "13800000000", "Correct-Horse-42", "10.0.0.3", time.Now())
	db.Model(&gen.Account{}).Where("id = ?", "account-1").Updates(map[string]any{"status": gen.AccountStatusActive, "is_delete": 2})
	_, deletedErr := service.Login(context.Background(), "13800000000", "Correct-Horse-42", "10.0.0.4", time.Now())
	if auth.ErrorCode(unknownErr) != auth.CodePhoneNotFound || auth.ErrorCode(disabledErr) != auth.CodeAccountInactive || auth.ErrorCode(deletedErr) != auth.CodePhoneNotFound {
		t.Fatalf("unexpected login errors: unknown=%v disabled=%v deleted=%v", unknownErr, disabledErr, deletedErr)
	}
	assertAuthenticationAudit(t, db, "authentication:login", string(auth.CodePhoneNotFound), 2)
	assertAuthenticationAudit(t, db, "authentication:login", string(auth.CodeAccountInactive), 1)
}

func TestLoginWrongPasswordReturnsSpecificCode(t *testing.T) {
	service, db, _ := newAuthenticationFixture(t)
	seedLoginAccount(t, db, "account-1", "13800000000", gen.AccountStatusActive)
	_, err := service.Login(context.Background(), "13800000000", "Wrong-Horse-42", "127.0.0.1", time.Now())
	if auth.ErrorCode(err) != auth.CodePasswordIncorrect {
		t.Fatalf("wrong password error = %v", err)
	}
	assertAuthenticationAudit(t, db, "authentication:login", string(auth.CodePasswordIncorrect), 1)
}

func TestDisabledAccountDoesNotReportPasswordIncorrect(t *testing.T) {
	service, db, _ := newAuthenticationFixture(t)
	seedLoginAccount(t, db, "account-1", "13800000000", gen.AccountStatusDisabled)
	_, err := service.Login(context.Background(), "13800000000", "Wrong-Horse-42", "127.0.0.1", time.Now())
	if auth.ErrorCode(err) != auth.CodeAccountInactive {
		t.Fatalf("disabled account error = %v", err)
	}
}

// TestLoginIPLimit 验证同一 IP 一分钟内第 21 次尝试被拒绝。
func TestLoginIPLimit(t *testing.T) {
	service, _, _ := newAuthenticationFixture(t)
	now := time.Now()
	for attempt := 0; attempt < 20; attempt++ {
		_, _ = service.Login(context.Background(), "13900000000", "Wrong-Horse-42", "192.0.2.1", now)
	}
	if _, err := service.Login(context.Background(), "13900000000", "Wrong-Horse-42", "192.0.2.1", now); auth.ErrorCode(err) != auth.CodeRateLimited {
		t.Fatalf("rate limit error = %v", err)
	}
}

// TestLoginLimitSeparatesAccountsBehindProxy 验证共享代理地址不会让不同账号共用同一限流桶。
func TestLoginLimitSeparatesAccountsBehindProxy(t *testing.T) {
	limiter := NewLoginLimiter(1, time.Minute)
	now := time.Now()
	if !limiter.Allow(loginLimitKey("192.0.2.1", "13800000000"), now) {
		t.Fatal("first account was unexpectedly limited")
	}
	if !limiter.Allow(loginLimitKey("192.0.2.1", "13900000000"), now) {
		t.Fatal("second account shared the proxy rate-limit bucket")
	}
	cleanupLimiter := NewLoginLimiter(1, time.Minute)
	cleanupLimiter.Allow(loginLimitKey("198.51.100.1", "old"), now.Add(-2*time.Minute))
	cleanupLimiter.Allow(loginLimitKey("198.51.100.1", "new"), now)
	if len(cleanupLimiter.entries) != 1 {
		t.Fatalf("expired limiter entries were retained: %d", len(cleanupLimiter.entries))
	}
}

// TestExpiredTemporaryPasswordCannotLogin 验证临时密码过期后不能创建会话。
func TestExpiredTemporaryPasswordCannotLogin(t *testing.T) {
	service, db, _ := newAuthenticationFixture(t)
	seedLoginAccount(t, db, "account-1", "13800000000", gen.AccountStatusActive)
	now := time.Now()
	if err := db.Model(&gen.Account{}).Where("id = ?", "account-1").Update("must_change_password", true).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&AccountCredential{}).Where("account_id = ?", "account-1").Update("temporary_password_expires_at", now.Add(-time.Minute)).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.Login(context.Background(), "13800000000", "Correct-Horse-42", "198.51.100.1", now); auth.ErrorCode(err) != auth.CodePasswordIncorrect {
		t.Fatalf("expired temporary password error = %v", err)
	}
}

func TestMissingCredentialCannotReportWrongPhoneOrPassword(t *testing.T) {
	service, db, _ := newAuthenticationFixture(t)
	seedLoginAccount(t, db, "account-1", "13800000000", gen.AccountStatusActive)
	if err := db.Delete(&AccountCredential{}, "account_id = ?", "account-1").Error; err != nil {
		t.Fatal(err)
	}
	_, err := service.Login(context.Background(), "13800000000", "Correct-Horse-42", "198.51.100.2", time.Now())
	if auth.ErrorCode(err) != auth.CodeAccountUnavailable {
		t.Fatalf("missing credential error = %v", err)
	}
}

// TestChangePasswordClearsTemporaryExpiry 验证正式改密后清除临时密码时限。
func TestChangePasswordClearsTemporaryExpiry(t *testing.T) {
	service, db, _ := newAuthenticationFixture(t)
	seedLoginAccount(t, db, "account-1", "13800000000", gen.AccountStatusActive)
	now := time.Now()
	db.Model(&gen.Account{}).Where("id = ?", "account-1").Update("must_change_password", true)
	db.Model(&AccountCredential{}).Where("account_id = ?", "account-1").Update("temporary_password_expires_at", now.Add(time.Hour))
	principal := &auth.WorkspacePrincipal{AccountID: "account-1", SessionID: "temporary-session", WorkspaceType: auth.WorkspaceTypeDiscovery}
	if _, err := service.ChangePassword(context.Background(), principal, "Correct-Horse-42", "Replace-Horse-42!", now); err != nil {
		t.Fatal(err)
	}
	var credential AccountCredential
	var account gen.Account
	db.First(&credential, "account_id = ?", "account-1")
	db.First(&account, "id = ?", "account-1")
	if credential.TemporaryPasswordExpiresAt != nil || account.MustChangePassword {
		t.Fatalf("temporary password state not cleared: account=%#v credential=%#v", account, credential)
	}
	assertAuthenticationAudit(t, db, "account:password_change", "SUCCESS", 1)
}

// TestChangePasswordNotifiesOnlyOtherOldSessions 验证替代当前会话不会被旧连接事件误强退。
func TestChangePasswordNotifiesOnlyOtherOldSessions(t *testing.T) {
	_, db, cfg := newAuthenticationFixture(t)
	seedLoginAccount(t, db, "account-1", "13800000000", gen.AccountStatusActive)
	now := time.Now()
	for _, id := range []string{"current-session", "other-session"} {
		if err := db.Create(&gen.Session{ID: id, AccountID: "account-1", WorkspaceType: gen.WorkspaceTypeDiscovery, CredentialVersion: 1, ExpiresAt: now.Add(time.Hour), LastSeenAt: now}).Error; err != nil {
			t.Fatal(err)
		}
	}
	publisher := sessionservice.NewPublisher()
	currentCtx, cancelCurrent := context.WithCancel(context.Background())
	defer cancelCurrent()
	otherCtx, cancelOther := context.WithCancel(context.Background())
	defer cancelOther()
	currentEvents, _ := publisher.Subscribe(currentCtx, "current-session")
	otherEvents, _ := publisher.Subscribe(otherCtx, "other-session")
	service := NewService(db, cfg, sessionservice.NewService(cfg.SessionDuration), NewLoginLimiter(20, time.Minute), audit.NewService(), publisher)
	principal := &auth.WorkspacePrincipal{AccountID: "account-1", SessionID: "current-session", WorkspaceType: auth.WorkspaceTypeDiscovery}
	if _, err := service.ChangePassword(context.Background(), principal, "Correct-Horse-42", "Replace-Horse-42!", now); err != nil {
		t.Fatal(err)
	}
	select {
	case event := <-otherEvents:
		if event == nil || event.Code != gen.SessionEventCodeCredentialsChanged {
			t.Fatalf("unexpected other-session event: %#v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("other old session was not notified")
	}
	select {
	case event := <-currentEvents:
		t.Fatalf("current rotating session was notified: %#v", event)
	case <-time.After(20 * time.Millisecond):
	}
}

func newAuthenticationFixture(t *testing.T) (*Service, *gorm.DB, config.SecurityConfig) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range []any{&gen.Account{}, &gen.Session{}, &gen.AuditLog{}, &AccountCredential{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	cfg := config.SecurityConfig{
		SigningKey: []byte(strings.Repeat("k", 32)), TokenIssuer: "test",
		SessionDuration: 12 * time.Hour, LoginLockThreshold: 5,
		LoginLockDuration: 15 * time.Minute, LoginIPAttemptsPerMinute: 20,
	}
	sessions := sessionservice.NewService(cfg.SessionDuration)
	return NewService(db, cfg, sessions, NewLoginLimiter(20, time.Minute), audit.NewService(), nil), db, cfg
}

func assertAuthenticationAudit(t *testing.T, db *gorm.DB, action, resultCode string, expected int64) {
	t.Helper()
	var count int64
	if err := db.Model(&gen.AuditLog{}).Where("action = ? AND result_code = ?", action, resultCode).Count(&count).Error; err != nil || count != expected {
		t.Fatalf("audit %s/%s count=%d want=%d err=%v", action, resultCode, count, expected, err)
	}
}

func seedLoginAccount(t *testing.T, db *gorm.DB, id, phone string, status gen.AccountStatus) {
	t.Helper()
	hash, err := auth.HashPassword("Correct-Horse-42")
	if err != nil {
		t.Fatal(err)
	}
	account := gen.Account{ID: id, Phone: phone, DisplayName: "Operator", Status: status, CredentialVersion: 1}
	now := time.Now()
	credential := AccountCredential{
		AccountID: id, PasswordHash: hash, TemporaryPasswordExpiresAt: NewTemporaryPasswordExpiry(now),
		PasswordChangedAt: now, UpdatedAt: now,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&credential).Error; err != nil {
		t.Fatal(err)
	}
}
