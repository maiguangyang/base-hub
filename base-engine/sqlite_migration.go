package main

import (
	"fmt"
	"sort"

	"base-engine/gen"
	"gorm.io/gorm"
)

var generatedSharedIndexColumns = []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"}

func autoMigrateEngineDB(database *gen.DB) error {
	switch database.Query().Dialector.Name() {
	case "sqlite":
		return autoMigrateSQLite(database.Query())
	case "mysql":
		return autoMigrateMySQL(database.Query())
	default:
		return database.AutoMigrate()
	}
}

// The generated models omit indexes owned by EnsureGovernanceIndexes. MySQL's
// AutoMigrate otherwise drops single-column unique indexes, including one used
// as the opening-record foreign key's backing index.
type franchiseOpeningRecord struct {
	gen.FranchiseOpeningRecord `gorm:"embedded"`
	OrganizationID             string `gorm:"column:organization_id;type:varchar(36);comment:'organization_id';default:null;uniqueIndex:uidx_opening_record_organization"`
}

type membershipInvitation struct {
	gen.MembershipInvitation `gorm:"embedded"`
	MembershipID             string `gorm:"column:membership_id;type:varchar(36);comment:'membership_id';default:null;uniqueIndex:uidx_invitation_membership"`
}

type globalPaymentConfigMigration struct {
	gen.GlobalPaymentConfig `gorm:"embedded"`
	Channel                 string `gorm:"column:channel;type:varchar(16);NOT NULL;uniqueIndex:uidx_global_payment_config_channel"`
}

func (globalPaymentConfigMigration) TableName() string { return "global_payment_configs" }

type franchisePaymentConfigMigration struct {
	gen.FranchisePaymentConfig `gorm:"embedded"`
	OrganizationID             string `gorm:"column:organization_id;type:varchar(36);uniqueIndex:uidx_franchise_payment_config_org_channel,priority:1"`
	Channel                    string `gorm:"column:channel;type:varchar(16);NOT NULL;uniqueIndex:uidx_franchise_payment_config_org_channel,priority:2"`
}

func (franchisePaymentConfigMigration) TableName() string { return "franchise_payment_configs" }

type storePaymentConfigMigration struct {
	gen.StorePaymentConfig `gorm:"embedded"`
	StoreID                string `gorm:"column:store_id;type:varchar(36);uniqueIndex:uidx_store_payment_config_store_channel,priority:1"`
	Channel                string `gorm:"column:channel;type:varchar(16);NOT NULL;uniqueIndex:uidx_store_payment_config_store_channel,priority:2"`
}

func (storePaymentConfigMigration) TableName() string { return "store_payment_configs" }

func autoMigrateMySQL(db *gorm.DB) error {
	names := make([]string, 0, len(gen.TableMap))
	for name := range gen.TableMap {
		names = append(names, name)
	}
	sort.Strings(names)
	models := make([]any, 0, len(names))
	for _, name := range names {
		switch name {
		case "franchise_opening_records":
			models = append(models, &franchiseOpeningRecord{})
		case "membership_invitations":
			models = append(models, &membershipInvitation{})
		case "global_payment_configs":
			models = append(models, &globalPaymentConfigMigration{})
		case "franchise_payment_configs":
			models = append(models, &franchisePaymentConfigMigration{})
		case "store_payment_configs":
			models = append(models, &storePaymentConfigMigration{})
		default:
			models = append(models, gen.TableMap[name])
		}
	}
	return db.AutoMigrate(models...)
}

func autoMigrateSQLite(db *gorm.DB) error {
	migrationDB := db.Session(&gorm.Session{NewDB: true})
	migrationDB.Config.DisableForeignKeyConstraintWhenMigrating = true
	migrationDB.Config.IgnoreRelationshipsWhenMigrating = true
	names := make([]string, 0, len(gen.TableMap))
	for name := range gen.TableMap {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if err := dropGeneratedSharedIndexes(migrationDB); err != nil {
			return err
		}
		if err := migrationDB.AutoMigrate(gen.TableMap[name]); err != nil {
			return err
		}
	}
	if err := dropGeneratedSharedIndexes(migrationDB); err != nil {
		return err
	}
	return nil
}

func ensureSQLiteSharedIndexes(db *gorm.DB) error {
	if db.Dialector.Name() != "sqlite" {
		return nil
	}
	for _, model := range gen.TableMap {
		if err := createScopedSharedIndexes(db, model); err != nil {
			return err
		}
	}
	return nil
}

func dropGeneratedSharedIndexes(db *gorm.DB) error {
	for _, column := range generatedSharedIndexColumns {
		if err := db.Exec("DROP INDEX IF EXISTS `" + column + "`").Error; err != nil {
			return err
		}
	}
	return nil
}

func createScopedSharedIndexes(db *gorm.DB, model any) error {
	statement := &gorm.Statement{DB: db}
	if err := statement.Parse(model); err != nil {
		return err
	}
	table := statement.Schema.Table
	for _, column := range generatedSharedIndexColumns {
		if !db.Migrator().HasColumn(model, column) {
			continue
		}
		index := "idx_" + table + "_" + column
		if db.Migrator().HasIndex(model, index) {
			continue
		}
		query := fmt.Sprintf("CREATE INDEX `%s` ON `%s` (`%s`)", index, table, column)
		if err := db.Exec(query).Error; err != nil {
			return err
		}
	}
	return nil
}
