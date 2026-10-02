/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package dbup

import (
	"fmt"

	"base-engine/gen"
	"gorm.io/gorm"
)

type governanceIndex struct {
	Model   any
	Name    string
	Columns string
}

var governanceIndexes = []governanceIndex{
	{Model: &gen.OperatorMembership{}, Name: "uidx_membership_account_organization", Columns: "account_id, organization_id"},
	{Model: &gen.OperatorRole{}, Name: "uidx_role_organization_name", Columns: "organization_id, name"},
	{Model: &gen.Store{}, Name: "uidx_store_organization_code", Columns: "organization_id, code"},
	{Model: &gen.MembershipInvitation{}, Name: "uidx_invitation_membership", Columns: "membership_id"},
	{Model: &gen.FranchiseOpeningRecord{}, Name: "uidx_opening_record_organization", Columns: "organization_id"},
	{Model: &gen.GlobalPaymentConfig{}, Name: "uidx_global_payment_config_channel", Columns: "channel"},
	{Model: &gen.FranchisePaymentConfig{}, Name: "uidx_franchise_payment_config_org_channel", Columns: "organization_id, channel"},
	{Model: &gen.StorePaymentConfig{}, Name: "uidx_store_payment_config_store_channel", Columns: "store_id, channel"},
}

// EnsureGovernanceIndexes 创建跨生成字段的租户自然键唯一索引。
func EnsureGovernanceIndexes(db *gorm.DB) error {
	dialect := db.Dialector.Name()
	if dialect != "sqlite" && dialect != "mysql" {
		return fmt.Errorf("unsupported governance index dialect: %s", dialect)
	}
	if err := ensureOpeningRecordRelationships(db); err != nil {
		return err
	}
	if err := ensurePaymentRateConstraints(db); err != nil {
		return err
	}
	for _, index := range governanceIndexes {
		if err := ensureGovernanceIndex(db, index); err != nil {
			return err
		}
	}
	for _, field := range []string{"RecordNumber", "Source"} {
		if err := ensureOpeningRecordGeneratedIndex(db, field); err != nil {
			return err
		}
	}
	return nil
}

func ensureOpeningRecordGeneratedIndex(db *gorm.DB, field string) error {
	if db.Migrator().HasIndex(&gen.FranchiseOpeningRecord{}, field) {
		return nil
	}
	return db.Migrator().CreateIndex(&gen.FranchiseOpeningRecord{}, field)
}

func ensureOpeningRecordRelationships(db *gorm.DB) error {
	recordTable, err := governanceTableName(db, &gen.FranchiseOpeningRecord{})
	if err != nil {
		return err
	}
	organizationTable, err := governanceTableName(db, &gen.Organization{})
	if err != nil {
		return err
	}
	accountTable, err := governanceTableName(db, &gen.Account{})
	if err != nil {
		return err
	}
	if err := checkOpeningRecordReferences(db, recordTable, organizationTable, accountTable); err != nil {
		return err
	}
	return createOpeningRecordConstraints(db, recordTable)
}

func checkOpeningRecordReferences(db *gorm.DB, recordTable, organizationTable, accountTable string) error {
	for _, relation := range []struct{ name, column, table string }{
		{"organization", "organization_id", organizationTable},
		{"initial account", "initial_account_id", accountTable},
		{"recording actor", "recorded_by_account_id", accountTable},
	} {
		var orphan struct{ ID string }
		query := fmt.Sprintf("SELECT r.id FROM `%s` r LEFT JOIN `%s` parent ON parent.id = r.`%s` WHERE r.`%s` IS NULL OR r.`%s` = '' OR parent.id IS NULL LIMIT 1", recordTable, relation.table, relation.column, relation.column, relation.column)
		result := db.Raw(query).Scan(&orphan)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 0 {
			return fmt.Errorf("franchise opening record %s has missing %s relationship; repair the record before migration", orphan.ID, relation.name)
		}
	}
	return nil
}

func createOpeningRecordConstraints(db *gorm.DB, recordTable string) error {
	constraintDB := db.Session(&gorm.Session{NewDB: true}).Table(recordTable)
	constraintDB.Config.DisableForeignKeyConstraintWhenMigrating = false
	constraintDB.Config.IgnoreRelationshipsWhenMigrating = false
	for _, relationship := range []string{"Organization", "InitialAccount", "RecordedByAccount"} {
		if constraintDB.Migrator().HasConstraint(&openingRecordConstraintModel{}, relationship) {
			continue
		}
		if err := constraintDB.Migrator().CreateConstraint(&openingRecordConstraintModel{}, relationship); err != nil {
			return fmt.Errorf("create franchise opening record %s constraint: %w", relationship, err)
		}
	}
	for _, column := range []string{"OrganizationID", "InitialAccountID", "RecordedByAccountID"} {
		if constraintDB.Migrator().HasConstraint(&openingRecordConstraintModel{}, column) {
			continue
		}
		if err := constraintDB.Migrator().CreateConstraint(&openingRecordConstraintModel{}, column); err != nil {
			return fmt.Errorf("create franchise opening record %s required constraint: %w", column, err)
		}
	}
	return nil
}

func governanceTableName(db *gorm.DB, model any) (string, error) {
	statement := &gorm.Statement{DB: db}
	if err := statement.Parse(model); err != nil {
		return "", err
	}
	return statement.Schema.Table, nil
}

// openingRecordConstraintModel supplies the database migration with the
// relationships declared in model/model.graphql without editing generated code.
type openingRecordConstraintModel struct {
	OrganizationID      string           `gorm:"column:organization_id;check:chk_opening_record_organization_id,organization_id IS NOT NULL AND organization_id <> ''"`
	InitialAccountID    string           `gorm:"column:initial_account_id;check:chk_opening_record_initial_account_id,initial_account_id IS NOT NULL AND initial_account_id <> ''"`
	RecordedByAccountID string           `gorm:"column:recorded_by_account_id;check:chk_opening_record_recorded_by_account_id,recorded_by_account_id IS NOT NULL AND recorded_by_account_id <> ''"`
	Organization        gen.Organization `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	InitialAccount      gen.Account      `gorm:"foreignKey:InitialAccountID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	RecordedByAccount   gen.Account      `gorm:"foreignKey:RecordedByAccountID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (openingRecordConstraintModel) TableName() string { return "franchise_opening_records" }

func ensureGovernanceIndex(db *gorm.DB, index governanceIndex) error {
	if db.Migrator().HasIndex(index.Model, index.Name) {
		return nil
	}
	statement := &gorm.Statement{DB: db}
	if err := statement.Parse(index.Model); err != nil {
		return err
	}
	query := fmt.Sprintf("CREATE UNIQUE INDEX `%s` ON `%s` (%s)", index.Name, statement.Schema.Table, index.Columns)
	return db.Exec(query).Error
}
