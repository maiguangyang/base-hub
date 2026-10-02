/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package dbup

import (
	"strings"
	"testing"

	"base-engine/gen"
	"base-engine/src/services/ai"
	"base-engine/src/services/authentication"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// TestMigrateSecurityTables 验证秘密表独立于 GraphQL 实体完成迁移。
func TestMigrateSecurityTables(t *testing.T) {
	db := openTestDB(t)
	if err := MigrateSecurityTables(db); err != nil {
		t.Fatal(err)
	}
	if !db.Migrator().HasTable(&authentication.AccountCredential{}) {
		t.Fatal("account credential table was not created")
	}
	if !db.Migrator().HasTable(&ai.StoredModelConfig{}) {
		t.Fatal("AI model configuration table was not created")
	}
}

// TestEnsureGovernanceIndexes 验证租户自然键和一对一邀请均受数据库唯一约束保护。
func TestEnsureGovernanceIndexes(t *testing.T) {
	db := openTestDB(t)
	if err := migrateGovernanceModels(db); err != nil {
		t.Fatal(err)
	}
	assertOpeningRecordSourceIndex(t, db)
	if err := EnsureGovernanceIndexes(db); err != nil {
		t.Fatal(err)
	}
	assertUniqueIndex(t, db, "operator_memberships", "uidx_membership_account_organization")
	assertUniqueIndex(t, db, "operator_roles", "uidx_role_organization_name")
	assertUniqueIndex(t, db, "stores", "uidx_store_organization_code")
	assertUniqueIndex(t, db, "membership_invitations", "uidx_invitation_membership")
	assertUniqueIndex(t, db, "franchise_opening_records", "uidx_opening_record_organization")
	assertUniqueIndex(t, db, "franchise_opening_records", "idx_franchise_opening_records_record_number")
	assertOpeningRecordSourceIndex(t, db)
	seedOpeningRecordParents(t, db)
	first := gen.FranchiseOpeningRecord{ID: "opening-1", RecordNumber: "OPEN-1", Source: gen.FranchiseOpeningSourceHistoricalAttestation, OrganizationID: "org-1", InitialAccountID: "account-1", RecordedByAccountID: "actor-1"}
	if err := db.Create(&first).Error; err != nil {
		t.Fatal(err)
	}
	if err := migrateGovernanceModels(db); err != nil {
		t.Fatalf("generated remigration with opening record: %v", err)
	}
	if err := EnsureGovernanceIndexes(db); err != nil {
		t.Fatalf("repeat migration with existing opening record: %v", err)
	}
	assertUniqueIndex(t, db, "franchise_opening_records", "uidx_opening_record_organization")
	assertUniqueIndex(t, db, "franchise_opening_records", "idx_franchise_opening_records_record_number")
	assertOpeningRecordSourceIndex(t, db)
	assertOpeningRecordForeignKeys(t, db)
	assertOpeningRecordNullsRejected(t, db)
	duplicate := first
	duplicate.ID, duplicate.RecordNumber = "opening-2", "OPEN-2"
	if err := db.Create(&duplicate).Error; err == nil {
		t.Fatal("duplicate organization opening record was accepted")
	}
	duplicate.OrganizationID, duplicate.RecordNumber = "org-2", "OPEN-1"
	if err := db.Create(&duplicate).Error; err == nil {
		t.Fatal("duplicate opening record number was accepted")
	}
	assertGovernanceDuplicatesRejected(t, db)
}

func TestEnsureGovernanceIndexesRejectsOrphanOpeningRecordRelationships(t *testing.T) {
	db := openTestDB(t)
	if err := migrateGovernanceModels(db); err != nil {
		t.Fatal(err)
	}
	if err := EnsureGovernanceIndexes(db); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatal(err)
	}
	seedOpeningRecordParents(t, db)
	assertOpeningRecordForeignKeys(t, db)
	assertOpeningRecordOrphansRejected(t, db)
	assertOpeningRecordNullsRejected(t, db)
}

func seedOpeningRecordParents(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, account := range []gen.Account{{ID: "account-1", Phone: "13800000001", DisplayName: "Initial", Status: gen.AccountStatusActive}, {ID: "actor-1", Phone: "13800000002", DisplayName: "Actor", Status: gen.AccountStatusActive}} {
		if err := db.Create(&account).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, organization := range []gen.Organization{{ID: "org-1", Code: "F001", Name: "Franchise", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}, {ID: "org-2", Code: "F002", Name: "Franchise 2", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}} {
		if err := db.Create(&organization).Error; err != nil {
			t.Fatal(err)
		}
	}
}

func assertOpeningRecordForeignKeys(t *testing.T, db *gorm.DB) {
	t.Helper()
	var foreignKeys []struct{ Table string }
	if err := db.Raw("PRAGMA foreign_key_list('franchise_opening_records')").Scan(&foreignKeys).Error; err != nil {
		t.Fatal(err)
	}
	if len(foreignKeys) != 3 {
		t.Fatalf("opening record foreign keys = %v", foreignKeys)
	}
}

func assertOpeningRecordSourceIndex(t *testing.T, db *gorm.DB) {
	t.Helper()
	if !db.Migrator().HasIndex(&gen.FranchiseOpeningRecord{}, "idx_franchise_opening_records_source") {
		t.Fatal("opening source index missing")
	}
}

func assertOpeningRecordOrphansRejected(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, relation := range []struct {
		name, organizationID, initialAccountID, recordedByAccountID string
	}{
		{"organization", "missing-org", "account-1", "actor-1"},
		{"initial account", "org-1", "missing-account", "actor-1"},
		{"recording actor", "org-1", "account-1", "missing-actor"},
	} {
		record := gen.FranchiseOpeningRecord{ID: relation.name, RecordNumber: relation.name, Source: gen.FranchiseOpeningSourceHistoricalAttestation, OrganizationID: relation.organizationID, InitialAccountID: relation.initialAccountID, RecordedByAccountID: relation.recordedByAccountID}
		if err := db.Create(&record).Error; err == nil {
			t.Fatalf("orphan %s relationship was accepted", relation.name)
		}
	}
}

func assertOpeningRecordNullsRejected(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, column := range []string{"organization_id", "initial_account_id", "recorded_by_account_id"} {
		query := "INSERT INTO franchise_opening_records (id, record_number, source, organization_id, initial_account_id, recorded_by_account_id) VALUES (?, ?, ?, ?, ?, ?)"
		values := []any{"null-" + column, "NULL-" + column, gen.FranchiseOpeningSourceHistoricalAttestation, "org-1", "account-1", "actor-1"}
		switch column {
		case "organization_id":
			values[3] = nil
		case "initial_account_id":
			values[4] = nil
		case "recorded_by_account_id":
			values[5] = nil
		}
		if err := db.Exec(query, values...).Error; err == nil {
			t.Fatalf("null %s relationship was accepted", column)
		}
	}
}

func TestEnsureGovernanceIndexesRespectsTablePrefix(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
		NamingStrategy:                           schema.NamingStrategy{TablePrefix: "prefix_"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := migrateGovernanceModels(db); err != nil {
		t.Fatal(err)
	}
	if err := EnsureGovernanceIndexes(db); err != nil {
		t.Fatal(err)
	}
	var foreignKeys []struct{ Table string }
	if err := db.Raw("PRAGMA foreign_key_list('prefix_franchise_opening_records')").Scan(&foreignKeys).Error; err != nil {
		t.Fatal(err)
	}
	if len(foreignKeys) != 3 {
		t.Fatalf("prefixed opening record foreign keys = %v", foreignKeys)
	}
	if err := EnsureGovernanceIndexes(db); err != nil {
		t.Fatalf("repeat prefixed migration: %v", err)
	}
}

func TestEnsureGovernanceIndexesReportsPreexistingOrphan(t *testing.T) {
	db := openTestDB(t)
	if err := migrateGovernanceModels(db); err != nil {
		t.Fatal(err)
	}
	orphan := gen.FranchiseOpeningRecord{ID: "opening-orphan", RecordNumber: "OPEN-ORPHAN", Source: gen.FranchiseOpeningSourceHistoricalAttestation, OrganizationID: "missing-org", InitialAccountID: "missing-account", RecordedByAccountID: "missing-actor"}
	if err := db.Create(&orphan).Error; err != nil {
		t.Fatal(err)
	}
	if err := EnsureGovernanceIndexes(db); err == nil || !strings.Contains(err.Error(), "opening-orphan") || !strings.Contains(err.Error(), "organization") {
		t.Fatalf("expected actionable orphan error, got %v", err)
	}
}

func migrateGeneratedForSQLite(db *gorm.DB, models ...any) error {
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			if err := db.Exec("DROP INDEX IF EXISTS `" + index + "`").Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func migrateGovernanceModels(db *gorm.DB) error {
	return migrateGeneratedForSQLite(db,
		&gen.Account{}, &gen.Organization{}, &gen.FranchiseOpeningRecord{},
		&gen.OperatorMembership{}, &gen.OperatorRole{}, &gen.Store{}, &gen.MembershipInvitation{},
		&gen.GlobalPaymentConfig{}, &gen.FranchisePaymentConfig{}, &gen.StorePaymentConfig{},
	)
}

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func assertUniqueIndex(t *testing.T, db *gorm.DB, table, index string) {
	t.Helper()
	if !db.Migrator().HasIndex(table, index) {
		t.Fatalf("missing unique index %s", index)
	}
}

func assertGovernanceDuplicatesRejected(t *testing.T, db *gorm.DB) {
	t.Helper()
	memberships := []gen.OperatorMembership{
		{ID: "membership-1", AccountID: "account-1", OrganizationID: "org-1"},
		{ID: "membership-2", AccountID: "account-1", OrganizationID: "org-1"},
	}
	roles := []gen.OperatorRole{
		{ID: "role-1", OrganizationID: "org-1", Name: "manager"},
		{ID: "role-2", OrganizationID: "org-1", Name: "manager"},
	}
	stores := []gen.Store{
		{ID: "store-1", OrganizationID: "org-1", Code: "S001"},
		{ID: "store-2", OrganizationID: "org-1", Code: "S001"},
	}
	invitations := []gen.MembershipInvitation{
		{ID: "invitation-1", MembershipID: "membership-1", InvitedByAccountID: "account-1"},
		{ID: "invitation-2", MembershipID: "membership-1", InvitedByAccountID: "account-1"},
	}
	for _, records := range []any{memberships, roles, stores, invitations} {
		if err := db.Create(records).Error; err == nil {
			t.Fatalf("duplicate natural key accepted for %T", records)
		}
	}
}
