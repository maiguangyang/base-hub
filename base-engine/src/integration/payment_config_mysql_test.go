package integration

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"base-engine/src/dbup"
	"base-engine/src/services/paymentconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestPaymentConfigMySQLUniqueness(t *testing.T) {
	db := paymentConfigMySQLTestDB(t)
	var version string
	if err := db.Raw("SELECT VERSION()").Scan(&version).Error; err != nil {
		t.Fatal(err)
	}
	t.Logf("MySQL version: %s", version)
	if err := migratePaymentConfigMySQL(db); err != nil {
		t.Fatal(err)
	}
	for _, row := range []struct{ table, index string }{
		{"global_payment_configs", "uidx_global_payment_config_channel"},
		{"franchise_payment_configs", "uidx_franchise_payment_config_org_channel"},
		{"store_payment_configs", "uidx_store_payment_config_store_channel"},
	} {
		var nonUnique int
		if err := db.Raw("SELECT NON_UNIQUE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND INDEX_NAME = ? LIMIT 1", row.table, row.index).Scan(&nonUnique).Error; err != nil {
			t.Fatal(err)
		}
		if !db.Migrator().HasIndex(row.table, row.index) || nonUnique != 0 {
			t.Fatalf("missing unique index %s", row.index)
		}
	}
	assertPaymentConfigConcurrentCreate(t, db)
	assertPaymentConfigMySQLRateConstraints(t, db)
	assertPaymentConfigMySQLLifecycle(t, db)
}

func assertPaymentConfigMySQLRateConstraints(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, target := range []struct{ table, ownerColumn, ownerID string }{
		{"global_payment_configs", "", ""},
		{"franchise_payment_configs", "organization_id", "franchise-rate"},
		{"store_payment_configs", "store_id", "store-rate"},
	} {
		for _, rate := range []int{-1, 1000001} {
			row := map[string]any{"id": fmt.Sprintf("%s-%d", target.table, rate), "channel": "ALIPAY", "config_state": "DISABLED", "version": 1, "rate_ppm": rate}
			if target.ownerColumn != "" {
				row[target.ownerColumn] = target.ownerID
			}
			if err := db.Table(target.table).Create(row).Error; err == nil {
				t.Fatalf("%s accepted invalid rate %d", target.table, rate)
			}
		}
	}
}

func assertPaymentConfigMySQLLifecycle(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Create(&gen.Organization{ID: "franchise-lifecycle", Code: "LIFE", Name: "Lifecycle", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.Store{ID: "store-lifecycle", Code: "LIFE", Name: "Lifecycle", Lifecycle: gen.StoreLifecycleActive, OrganizationID: "franchise-lifecycle"}).Error; err != nil {
		t.Fatal(err)
	}
	keys := config.PaymentConfigSecurity{ActiveKeyID: "test", Keys: map[string][]byte{"test": []byte(strings.Repeat("k", 32))}}
	service := paymentconfig.NewStore(db, keys)
	principal := &auth.WorkspacePrincipal{AccountID: "test", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"paymentConfig:read": {}, "paymentConfig:manage": {}}}
	ref := paymentconfig.ScopeRef{Scope: "STORE", StoreID: "store-lifecycle"}
	create := paymentconfig.StateInput{ScopeRef: ref, Channel: "ALIPAY", State: "DISABLED"}
	first, err := service.SetState(context.Background(), principal, create)
	if err != nil {
		t.Fatal(err)
	}
	old := first.Channels[1].Own
	reset := paymentconfig.ResetInput{ScopeRef: ref, Channel: "ALIPAY", RecordID: old.RecordID, Version: old.Version}
	if _, err := service.RestoreInheritance(context.Background(), principal, reset); err != nil {
		t.Fatal(err)
	}
	second, err := service.SetState(context.Background(), principal, create)
	if err != nil || second.Channels[1].Own.RecordID == old.RecordID {
		t.Fatalf("recreate: %#v %v", second, err)
	}
	if _, err := service.RestoreInheritance(context.Background(), principal, reset); !errors.Is(err, paymentconfig.ErrConflict) {
		t.Fatalf("stale reset: %v", err)
	}
	if _, err := service.Save(context.Background(), principal, paymentconfig.SaveInput{ScopeRef: ref, Channel: "ALIPAY", RecordID: old.RecordID, Version: old.Version}); !errors.Is(err, paymentconfig.ErrConflict) {
		t.Fatalf("stale save: %v", err)
	}
}

func paymentConfigMySQLTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	rootDSN := os.Getenv("PAYMENT_TEST_MYSQL_DSN")
	if rootDSN == "" {
		t.Skip("PAYMENT_TEST_MYSQL_DSN is required for the disposable MySQL test")
	}
	if !strings.HasSuffix(rootDSN, "/") {
		t.Fatal("PAYMENT_TEST_MYSQL_DSN must end in / with no database selected")
	}
	admin, err := gorm.Open(mysql.Open(rootDSN+"mysql?parseTime=true"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("payment_config_test_%d", time.Now().UnixNano())
	if err := admin.Exec("CREATE DATABASE `" + name + "`").Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := admin.Exec("DROP DATABASE `" + name + "`").Error; err != nil {
			t.Errorf("drop test database: %v", err)
		}
	})
	db, err := gorm.Open(mysql.Open(rootDSN+name+"?parseTime=true"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func migratePaymentConfigMySQL(db *gorm.DB) error {
	names := make([]string, 0, len(gen.TableMap))
	for name := range gen.TableMap {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if err := db.AutoMigrate(gen.TableMap[name]); err != nil {
			return err
		}
	}
	return dbup.EnsureGovernanceIndexes(db)
}

func assertPaymentConfigConcurrentCreate(t *testing.T, db *gorm.DB) {
	t.Helper()
	var wait sync.WaitGroup
	results := make(chan error, 2)
	for _, id := range []string{"first", "second"} {
		wait.Add(1)
		go func(id string) {
			defer wait.Done()
			results <- db.Create(&gen.FranchisePaymentConfig{ID: id, OrganizationID: "franchise-1", Channel: "WECHAT", ConfigState: "DISABLED", Version: 1}).Error
		}(id)
	}
	wait.Wait()
	close(results)
	success, failure := 0, 0
	for err := range results {
		if err == nil {
			success++
		} else {
			failure++
		}
	}
	if success != 1 || failure != 1 {
		t.Fatalf("concurrent inserts: success=%d failure=%d", success, failure)
	}
}
