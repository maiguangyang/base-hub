/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package dbup

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authentication"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNormalizeHQBootstrapErrorAfterTransactionAbort(t *testing.T) {
	db := setupHQInitializationDB(t)
	if err := db.Create(&SecurityBootstrap{Key: hqBootstrapKey}).Error; err != nil {
		t.Fatal(err)
	}
	err := normalizeHQBootstrapError(context.Background(), db, errors.New("raw unique violation"))
	if auth.ErrorCode(err) != auth.CodeHQAlreadyBootstrapped {
		t.Fatalf("error code = %q", auth.ErrorCode(err))
	}
}

func TestInitializeHQCreatesPermanentAdministrator(t *testing.T) {
	t.Run("[Auth.Initialize] 初始化状态随永久总部管理员创建切换", func(t *testing.T) {
		db := setupHQInitializationDB(t)
		initialized, err := HQInitialized(context.Background(), db)
		if err != nil || initialized {
			t.Fatalf("initial status = %v, %v", initialized, err)
		}

		password := "Correct-Horse-42"
		if err := InitializeHQ(context.Background(), db, "13800000000", password, password); err != nil {
			t.Fatal(err)
		}
		initialized, err = HQInitialized(context.Background(), db)
		if err != nil || !initialized {
			t.Fatalf("initialized status = %v, %v", initialized, err)
		}
		assertPermanentHQAdministrator(t, db, "13800000000", password)
	})
}

func TestInitializeHQAllowsOnlyOneCompetingRequest(t *testing.T) {
	t.Run("[Auth.Initialize] 并发初始化只能创建一个最高管理员", func(t *testing.T) {
		db := setupConcurrentHQInitializationDB(t)
		sqlDB, err := db.DB()
		if err != nil {
			t.Fatal(err)
		}
		sqlDB.SetMaxOpenConns(2)
		errorsByRequest := runCompetingInitializations(t, db)
		assertInitializationRace(t, errorsByRequest)
		var count int64
		if err := db.Model(&gen.Account{}).Count(&count).Error; err != nil || count != 1 {
			t.Fatalf("account count = %d, err = %v", count, err)
		}
	})
}

func runCompetingInitializations(t *testing.T, db *gorm.DB) []error {
	t.Helper()
	start := make(chan struct{})
	ready := make(chan struct{}, 2)
	release := make(chan struct{})
	if err := db.Callback().Create().Before("gorm:create").Register("test:bootstrap-claim-barrier", func(tx *gorm.DB) {
		if tx.Statement.Table == "security_bootstraps" {
			ready <- struct{}{}
			<-release
		}
	}); err != nil {
		t.Fatal(err)
	}
	errorsByRequest := make([]error, 2)
	var wait sync.WaitGroup
	for index, phone := range []string{"13800000000", "13900000000"} {
		wait.Add(1)
		go func(index int, phone string) {
			defer wait.Done()
			<-start
			errorsByRequest[index] = InitializeHQ(context.Background(), db, phone, "Correct-Horse-42", "Correct-Horse-42")
		}(index, phone)
	}
	close(start)
	for range 2 {
		select {
		case <-ready:
		case <-time.After(5 * time.Second):
			t.Fatal("competing transactions did not reach the bootstrap marker together")
		}
	}
	close(release)
	wait.Wait()
	return errorsByRequest
}

func assertInitializationRace(t *testing.T, errorsByRequest []error) {
	t.Helper()
	var successes, conflicts int
	for _, initErr := range errorsByRequest {
		if initErr == nil {
			successes++
		} else if auth.ErrorCode(initErr) == auth.CodeHQAlreadyBootstrapped {
			conflicts++
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("successes = %d, conflicts = %d, errors = %#v", successes, conflicts, errorsByRequest)
	}
}

func TestInitializeHQExcludesCLIRebootstrap(t *testing.T) {
	t.Run("[Auth.Initialize] Web 初始化后 CLI 不能再次创建管理员", func(t *testing.T) {
		db := setupHQInitializationDB(t)
		if err := InitializeHQ(context.Background(), db, "13800000000", "Correct-Horse-42", "Correct-Horse-42"); err != nil {
			t.Fatal(err)
		}
		if _, err := BootstrapHQ(context.Background(), db, "13900000000", "Other"); auth.ErrorCode(err) != auth.CodeHQAlreadyBootstrapped {
			t.Fatalf("CLI bootstrap error = %v", err)
		}
	})
}

func TestInitializeHQRejectsRebootstrapBeforeCredentialValidation(t *testing.T) {
	db := setupHQInitializationDB(t)
	password := "Correct-Horse-42"
	if err := InitializeHQ(context.Background(), db, "13800000000", password, password); err != nil {
		t.Fatal(err)
	}
	err := InitializeHQ(context.Background(), db, "invalid", "weak", "different")
	if auth.ErrorCode(err) != auth.CodeHQAlreadyBootstrapped {
		t.Fatalf("error code = %q", auth.ErrorCode(err))
	}
}

func TestInitializeHQRejectsInvalidCredentialsWithoutClaimingBootstrap(t *testing.T) {
	tests := []struct {
		name         string
		phone        string
		password     string
		confirmation string
		code         auth.Code
	}{
		{name: "[Auth.Initialize] 非法手机号不能占用初始化权", phone: "123", password: "Correct-Horse-42", confirmation: "Correct-Horse-42", code: auth.CodeValidationFailed},
		{name: "[Auth.Initialize] 两次密码不一致不能占用初始化权", phone: "13800000000", password: "Correct-Horse-42", confirmation: "Different-Horse-42", code: auth.CodePasswordConfirmMismatch},
		{name: "[Auth.Initialize] 弱密码不能占用初始化权", phone: "13800000000", password: "weak", confirmation: "weak", code: auth.CodePasswordWeak},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := setupHQInitializationDB(t)
			err := InitializeHQ(context.Background(), db, test.phone, test.password, test.confirmation)
			if auth.ErrorCode(err) != test.code {
				t.Fatalf("error code = %q, want %q", auth.ErrorCode(err), test.code)
			}
			initialized, statusErr := HQInitialized(context.Background(), db)
			if statusErr != nil || initialized {
				t.Fatalf("invalid request claimed bootstrap: %v, %v", initialized, statusErr)
			}
		})
	}
}

func setupHQInitializationDB(t *testing.T) *gorm.DB {
	t.Helper()
	return setupHQInitializationSchema(t, openTestDB(t))
}

func setupConcurrentHQInitializationDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "bootstrap.db") + "?_busy_timeout=5000&_journal_mode=WAL"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return setupHQInitializationSchema(t, db)
}

func setupHQInitializationSchema(t *testing.T, db *gorm.DB) *gorm.DB {
	t.Helper()
	models := []any{
		&gen.Account{}, &gen.Organization{}, &gen.OperatorMembership{},
		&gen.Permission{}, &gen.OperatorRole{}, &gen.AuditLog{},
	}
	if err := migrateGeneratedForSQLite(db, models...); err != nil {
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
	return db
}

func assertPermanentHQAdministrator(t *testing.T, db *gorm.DB, phone, password string) {
	t.Helper()
	var account gen.Account
	if err := db.First(&account, "phone = ?", phone).Error; err != nil {
		t.Fatal(err)
	}
	if account.DisplayName != phone || account.MustChangePassword || account.Status != gen.AccountStatusActive {
		t.Fatalf("invalid account: %#v", account)
	}
	assertBootstrapMembership(t, db, account.ID)
	var credential authentication.AccountCredential
	if err := db.First(&credential, "account_id = ?", account.ID).Error; err != nil {
		t.Fatal(err)
	}
	if credential.TemporaryPasswordExpiresAt != nil || auth.VerifyPassword(credential.PasswordHash, password) != nil {
		t.Fatalf("invalid permanent credential: %#v", credential)
	}
	assertBootstrapAudit(t, db, password)
}
