package dbup

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"base-engine/gen"
	"gorm.io/gorm"
)

func TestPaymentConfigSchemaRelationships(t *testing.T) {
	schema, err := os.ReadFile("../../model/model.graphql")
	if err != nil {
		t.Fatal(err)
	}
	for _, relation := range []string{
		"paymentConfigs: [FranchisePaymentConfig!]! @relationship(inverse: \"organization\")",
		"paymentConfigs: [StorePaymentConfig!]! @relationship(inverse: \"store\")",
		"organization: Organization! @relationship(inverse: \"paymentConfigs\")",
		"store: Store! @relationship(inverse: \"paymentConfigs\")",
	} {
		if !strings.Contains(string(schema), relation) {
			t.Fatalf("missing typed payment relationship %q", relation)
		}
	}
}

func TestPaymentConfigUniqueIndexes(t *testing.T) {
	db := openTestDB(t)
	if err := migrateGeneratedForSQLite(db,
		&gen.Account{}, &gen.Organization{}, &gen.Store{}, &gen.FranchiseOpeningRecord{},
		&gen.OperatorMembership{}, &gen.OperatorRole{}, &gen.MembershipInvitation{},
		&gen.GlobalPaymentConfig{}, &gen.FranchisePaymentConfig{}, &gen.StorePaymentConfig{},
	); err != nil {
		t.Fatal(err)
	}
	if err := EnsureGovernanceIndexes(db); err != nil {
		t.Fatal(err)
	}
	for _, item := range []struct{ table, index string }{
		{"global_payment_configs", "uidx_global_payment_config_channel"},
		{"franchise_payment_configs", "uidx_franchise_payment_config_org_channel"},
		{"store_payment_configs", "uidx_store_payment_config_store_channel"},
	} {
		assertUniqueIndex(t, db, item.table, item.index)
	}
	assertGlobalPaymentUniqueness(t, db)
	assertFranchisePaymentUniqueness(t, db)
	assertStorePaymentUniqueness(t, db)
}

func TestPaymentConfigDatabaseRejectsOutOfRangeRates(t *testing.T) {
	db := openTestDB(t)
	if err := migrateGeneratedForSQLite(db,
		&gen.Account{}, &gen.Organization{}, &gen.Store{}, &gen.FranchiseOpeningRecord{},
		&gen.OperatorMembership{}, &gen.OperatorRole{}, &gen.MembershipInvitation{},
		&gen.GlobalPaymentConfig{}, &gen.FranchisePaymentConfig{}, &gen.StorePaymentConfig{},
	); err != nil {
		t.Fatal(err)
	}
	if err := EnsureGovernanceIndexes(db); err != nil {
		t.Fatal(err)
	}
	for _, target := range []struct{ table, ownerColumn, ownerID string }{
		{"global_payment_configs", "", ""},
		{"franchise_payment_configs", "organization_id", "franchise-a"},
		{"store_payment_configs", "store_id", "store-a"},
	} {
		for _, rate := range []int{-1, 1000001} {
			row := map[string]any{"id": fmt.Sprintf("%s-%d", target.table, rate), "channel": "WECHAT", "config_state": "DISABLED", "version": 1, "rate_ppm": rate}
			if target.ownerColumn != "" {
				row[target.ownerColumn] = target.ownerID
			}
			if err := db.Table(target.table).Create(row).Error; err == nil {
				t.Fatalf("%s accepted invalid rate %d", target.table, rate)
			}
		}
	}
}

func assertGlobalPaymentUniqueness(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, row := range []gen.GlobalPaymentConfig{
		{ID: "global-wechat", Channel: "WECHAT", ConfigState: "DISABLED", Version: 1},
		{ID: "global-alipay", Channel: "ALIPAY", ConfigState: "DISABLED", Version: 1},
	} {
		if err := db.Create(&row).Error; err != nil {
			t.Fatal(err)
		}
	}
	duplicate := gen.GlobalPaymentConfig{ID: "duplicate", Channel: "WECHAT", ConfigState: "DISABLED", Version: 1}
	if err := db.Create(&duplicate).Error; err == nil {
		t.Fatal("duplicate global channel accepted")
	}
}

func assertFranchisePaymentUniqueness(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, row := range []gen.FranchisePaymentConfig{
		{ID: "franchise-wechat", OrganizationID: "franchise-a", Channel: "WECHAT", ConfigState: "DISABLED", Version: 1},
		{ID: "franchise-alipay", OrganizationID: "franchise-a", Channel: "ALIPAY", ConfigState: "DISABLED", Version: 1},
		{ID: "other-franchise-wechat", OrganizationID: "franchise-b", Channel: "WECHAT", ConfigState: "DISABLED", Version: 1},
	} {
		if err := db.Create(&row).Error; err != nil {
			t.Fatal(err)
		}
	}
	franchiseDuplicate := gen.FranchisePaymentConfig{ID: "franchise-duplicate", OrganizationID: "franchise-a", Channel: "WECHAT", ConfigState: "DISABLED", Version: 1}
	if err := db.Create(&franchiseDuplicate).Error; err == nil {
		t.Fatal("duplicate franchise channel accepted")
	}
}

func assertStorePaymentUniqueness(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, row := range []gen.StorePaymentConfig{
		{ID: "store-wechat", StoreID: "store-a", Channel: "WECHAT", ConfigState: "DISABLED", Version: 1},
		{ID: "store-alipay", StoreID: "store-a", Channel: "ALIPAY", ConfigState: "DISABLED", Version: 1},
		{ID: "other-store-wechat", StoreID: "store-b", Channel: "WECHAT", ConfigState: "DISABLED", Version: 1},
	} {
		if err := db.Create(&row).Error; err != nil {
			t.Fatal(err)
		}
	}
	storeDuplicate := gen.StorePaymentConfig{ID: "store-duplicate", StoreID: "store-a", Channel: "WECHAT", ConfigState: "DISABLED", Version: 1}
	if err := db.Create(&storeDuplicate).Error; err == nil {
		t.Fatal("duplicate store channel accepted")
	}
}
