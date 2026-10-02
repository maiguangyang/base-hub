package ai

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"base-engine/config"
	"base-engine/gen"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAIModelConfigUnchangedSaveWaitsForConcurrentUpdateMySQL(t *testing.T) {
	db := disposableAIModelMySQLDB(t)
	security := config.AIModelSecurityConfig{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)}}
	store := NewModelConfigStore(db, security, nil)
	principal := modelConfigPrincipal("aiModelConfig:manage")
	saved, err := store.Save(t.Context(), principal, ModelConfigUpdate{Name: "example", BaseURL: "https://model.example/v1", APIKey: "secret"})
	if err != nil {
		t.Fatal(err)
	}

	writer := db.Begin()
	if writer.Error != nil {
		t.Fatal(writer.Error)
	}
	t.Cleanup(func() { _ = writer.Rollback().Error })
	if err := writer.Model(&StoredModelConfig{}).Where("id = ?", modelConfigRecordID).
		Updates(map[string]any{"name": "changed", "version": saved.Version + 1}).Error; err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, saveErr := store.Save(ctx, principal, ModelConfigUpdate{Name: "example", BaseURL: "https://model.example/v1", ExpectedVersion: saved.Version})
		done <- saveErr
	}()
	waitForAIModelLockWait(t, ctx, db, done)
	assertAIModelSaveConflict(t, ctx, writer, done)
}

func waitForAIModelLockWait(t *testing.T, ctx context.Context, db *gorm.DB, done <-chan error) {
	t.Helper()
	databaseName := db.Migrator().CurrentDatabase()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		var waiting int64
		err := db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM information_schema.innodb_trx AS trx
			JOIN information_schema.processlist AS process ON process.ID = trx.trx_mysql_thread_id
			WHERE trx.trx_state = 'LOCK WAIT' AND process.DB = ?
			AND trx.trx_query LIKE '%stored_model_configs%'`, databaseName).Scan(&waiting).Error
		if err != nil {
			t.Fatalf("inspect MySQL row lock wait: %v", err)
		}
		if waiting > 0 {
			return
		}
		select {
		case saveErr := <-done:
			t.Fatalf("save completed without waiting for row lock: %v", saveErr)
		case <-ctx.Done():
			t.Fatal("save did not wait for row lock")
		case <-ticker.C:
		}
	}
}

func assertAIModelSaveConflict(t *testing.T, ctx context.Context, writer *gorm.DB, done <-chan error) {
	t.Helper()
	if err := writer.Commit().Error; err != nil {
		t.Fatal(err)
	}
	select {
	case saveErr := <-done:
		if !errors.Is(saveErr, ErrModelConfigConflict) {
			t.Fatalf("save after concurrent update = %v", saveErr)
		}
	case <-ctx.Done():
		t.Fatal("save did not resume after concurrent writer committed")
	}
}

func disposableAIModelMySQLDB(t *testing.T) *gorm.DB {
	t.Helper()
	rootDSN := os.Getenv("base_MYSQL_TEST_DSN")
	if rootDSN == "" {
		t.Skip("set base_MYSQL_TEST_DSN for MySQL concurrency test")
	}
	if !strings.HasSuffix(rootDSN, "/") {
		t.Fatal("base_MYSQL_TEST_DSN must end with /")
	}
	admin, err := gorm.Open(mysql.Open(rootDSN+"mysql?parseTime=true"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("base_ai_config_test_%d", time.Now().UnixNano())
	if err := admin.Exec("CREATE DATABASE `" + name + "`").Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := admin.Exec("DROP DATABASE `" + name + "`").Error; err != nil {
			t.Errorf("drop disposable database: %v", err)
		}
	})
	db, err := gorm.Open(mysql.Open(rootDSN+name+"?parseTime=true"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&StoredModelConfig{}, &gen.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	return db
}
