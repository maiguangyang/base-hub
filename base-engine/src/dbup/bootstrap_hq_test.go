/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package dbup

import (
	"context"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authentication"
	"gorm.io/gorm"
)

// TestBootstrapHQCreatesInitialAdministratorOnce 验证首个总部账号原子创建且不可重复初始化。
func TestBootstrapHQCreatesInitialAdministratorOnce(t *testing.T) {
	db := openTestDB(t)
	if err := migrateGeneratedForSQLite(db, &gen.Account{}, &gen.Organization{}, &gen.OperatorMembership{}, &gen.Permission{}, &gen.OperatorRole{}, &gen.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&testMembershipRole{}); err != nil {
		t.Fatal(err)
	}
	if err := MigrateSecurityTables(db); err != nil {
		t.Fatal(err)
	}
	if err := InitRoles(db); err != nil {
		t.Fatal(err)
	}
	result, err := BootstrapHQ(context.Background(), db, "13800000000", "Root")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.TemporaryPassword) != 8 || strings.Contains(result.TemporaryPassword, result.AccountID) {
		t.Fatalf("unsafe temporary password result: %#v", result)
	}
	assertBootstrapPersistence(t, db, result)
	if _, err := BootstrapHQ(context.Background(), db, "13900000000", "Other"); auth.ErrorCode(err) != auth.CodeHQAlreadyBootstrapped {
		t.Fatalf("second bootstrap error = %v", err)
	}
}

type testMembershipRole struct {
	OperatorMembershipID string `gorm:"primaryKey"`
	OperatorRoleID       string `gorm:"primaryKey"`
}

func (testMembershipRole) TableName() string {
	return "operator_membership_roles"
}

func assertBootstrapPersistence(t *testing.T, db *gorm.DB, result BootstrapHQResult) {
	t.Helper()
	assertBootstrapAccount(t, db, result)
	assertBootstrapMembership(t, db, result.AccountID)
	assertBootstrapCredential(t, db, result)
	assertBootstrapAudit(t, db, result.TemporaryPassword)
}

func assertBootstrapAccount(t *testing.T, db *gorm.DB, result BootstrapHQResult) {
	t.Helper()
	var account gen.Account
	if err := db.First(&account, "id = ?", result.AccountID).Error; err != nil {
		t.Fatal(err)
	}
	if !account.MustChangePassword || account.Status != gen.AccountStatusActive {
		t.Fatalf("invalid account state: %#v", account)
	}
}

func assertBootstrapMembership(t *testing.T, db *gorm.DB, accountID string) {
	t.Helper()
	var membership gen.OperatorMembership
	if err := db.Preload("Roles").Where("account_id = ?", accountID).First(&membership).Error; err != nil {
		t.Fatal(err)
	}
	if membership.Status != gen.MembershipStatusActive || len(membership.Roles) != 1 || membership.Roles[0].Kind != gen.RoleKindHqSuperAdmin {
		t.Fatalf("invalid headquarters membership: %#v", membership)
	}
}

func assertBootstrapCredential(t *testing.T, db *gorm.DB, result BootstrapHQResult) {
	t.Helper()
	var credential authentication.AccountCredential
	if err := db.First(&credential, "account_id = ?", result.AccountID).Error; err != nil {
		t.Fatal(err)
	}
	if strings.Contains(credential.PasswordHash, result.TemporaryPassword) || auth.VerifyPassword(credential.PasswordHash, result.TemporaryPassword) != nil {
		t.Fatal("temporary password was not stored as bcrypt only")
	}
	if credential.TemporaryPasswordExpiresAt == nil || !credential.TemporaryPasswordExpiresAt.After(time.Now()) {
		t.Fatalf("temporary password expiry missing: %#v", credential)
	}
}

func assertBootstrapAudit(t *testing.T, db *gorm.DB, temporaryPassword string) {
	t.Helper()
	var audit gen.AuditLog
	if err := db.First(&audit, "action = ?", "hq:bootstrap").Error; err != nil {
		t.Fatal(err)
	}
	if audit.MetadataJSON != nil && strings.Contains(*audit.MetadataJSON, temporaryPassword) {
		t.Fatal("audit metadata leaked temporary password")
	}
}
