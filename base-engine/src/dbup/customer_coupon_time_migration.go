package dbup

import (
	"database/sql"
	"fmt"
	"time"

	"base-engine/gen"
	"gorm.io/gorm"
)

const (
	minimumCouponMillis int64 = 1_000_000_000_000
	maximumCouponMillis int64 = 9_999_999_999_999
)

type couponLegacyTimeColumn struct {
	Name     string
	Nullable bool
}

var couponLegacyTimeColumns = []couponLegacyTimeColumn{
	{Name: "issued_at"},
	{Name: "activated_at", Nullable: true},
	{Name: "expires_at", Nullable: true},
	{Name: "revoked_at", Nullable: true},
}

func PrepareCustomerCouponTimeColumns(db *gorm.DB) error {
	if err := validateCouponMigrationDialect(db); err != nil {
		return err
	}
	templateTable, err := couponMigrationTableName(db, &gen.CustomerCouponTemplate{})
	if err != nil {
		return err
	}
	grantTable, err := couponMigrationTableName(db, &gen.CustomerCouponGrant{})
	if err != nil {
		return err
	}
	if db.Migrator().HasTable(templateTable) {
		if err := prepareCouponTemplateEffectiveAt(db, templateTable); err != nil {
			return err
		}
	}
	if !db.Migrator().HasTable(grantTable) {
		return nil
	}
	for _, column := range couponLegacyTimeColumns {
		if err := prepareCouponGrantTimeColumn(db, grantTable, column); err != nil {
			return err
		}
	}
	return nil
}

func prepareCouponTemplateEffectiveAt(db *gorm.DB, table string) error {
	if !db.Migrator().HasColumn(table, "effective_at") {
		if err := db.Exec(fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `effective_at` BIGINT", table)).Error; err != nil {
			return err
		}
	}
	query := fmt.Sprintf("UPDATE `%s` SET `effective_at` = `created_at` WHERE `effective_at` IS NULL OR `effective_at` = 0", table)
	if err := db.Exec(query).Error; err != nil {
		return err
	}
	return validateCouponMillisColumn(db, table, "effective_at", false)
}

func prepareCouponGrantTimeColumn(db *gorm.DB, table string, column couponLegacyTimeColumn) error {
	shadow, legacy := column.Name+"_ms", column.Name+"_legacy"
	hasCurrent := db.Migrator().HasColumn(table, column.Name)
	hasShadow := db.Migrator().HasColumn(table, shadow)
	hasLegacy := db.Migrator().HasColumn(table, legacy)
	if !hasCurrent {
		return recoverInterruptedCouponTimeSwap(db, table, column, hasShadow, hasLegacy)
	}
	integer, err := couponColumnIsInteger(db, table, column.Name)
	if err != nil {
		return err
	}
	if integer {
		return finishCouponIntegerColumn(db, table, column, hasShadow, hasLegacy)
	}
	if column.Name == "issued_at" && db.Dialector.Name() == "mysql" {
		if err := preflightLegacyIssuedAt(db, table); err != nil {
			return err
		}
	}
	return migrateLegacyCouponTimeColumn(db, table, column)
}

func finishCouponIntegerColumn(db *gorm.DB, table string, column couponLegacyTimeColumn, hasShadow, hasLegacy bool) error {
	for name, exists := range map[string]bool{column.Name + "_legacy": hasLegacy, column.Name + "_ms": hasShadow} {
		if exists {
			if err := dropCouponMigrationColumn(db, table, name); err != nil {
				return err
			}
		}
	}
	return validateCouponMillisColumn(db, table, column.Name, column.Nullable)
}

func recoverInterruptedCouponTimeSwap(db *gorm.DB, table string, column couponLegacyTimeColumn, hasShadow, hasLegacy bool) error {
	shadow, legacy := column.Name+"_ms", column.Name+"_legacy"
	if hasShadow {
		if err := renameCouponMigrationColumn(db, table, shadow, column.Name); err != nil {
			return err
		}
		if hasLegacy {
			if err := dropCouponMigrationColumn(db, table, legacy); err != nil {
				return err
			}
		}
		return validateCouponMillisColumn(db, table, column.Name, column.Nullable)
	}
	if hasLegacy {
		if err := renameCouponMigrationColumn(db, table, legacy, column.Name); err != nil {
			return err
		}
		return prepareCouponGrantTimeColumn(db, table, column)
	}
	return nil
}

func renameCouponMigrationColumn(db *gorm.DB, table, from, to string) error {
	return db.Exec(fmt.Sprintf("ALTER TABLE `%s` RENAME COLUMN `%s` TO `%s`", table, from, to)).Error
}

func dropCouponMigrationColumn(db *gorm.DB, table, name string) error {
	return db.Exec(fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`", table, name)).Error
}

func preflightLegacyIssuedAt(db *gorm.DB, table string) error {
	type row struct {
		ID        string
		IssuedAt  time.Time
		CreatedAt int64
	}
	lastID := ""
	for {
		var rows []row
		query := fmt.Sprintf("SELECT `id`, `issued_at`, `created_at` FROM `%s` WHERE `id` > ? ORDER BY `id` LIMIT 1000", table)
		if err := db.Raw(query, lastID).Scan(&rows).Error; err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}
		for _, item := range rows {
			delta := item.IssuedAt.UnixMilli() - item.CreatedAt
			if delta < 0 {
				delta = -delta
			}
			if delta >= 1000 {
				return fmt.Errorf("legacy issued_at timezone mismatch for grant %s", item.ID)
			}
			lastID = item.ID
		}
	}
}

func migrateLegacyCouponTimeColumn(db *gorm.DB, table string, column couponLegacyTimeColumn) error {
	shadow := column.Name + "_ms"
	if !db.Migrator().HasColumn(table, shadow) {
		if err := db.Exec(fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` BIGINT", table, shadow)).Error; err != nil {
			return err
		}
	}
	if err := copyLegacyCouponTimes(db, table, column.Name, shadow); err != nil {
		return err
	}
	if err := verifyCouponTimeParity(db, table, column.Name, shadow, column.Nullable); err != nil {
		return err
	}
	legacy := column.Name + "_legacy"
	if err := swapCouponMigrationColumns(db, table, column.Name, shadow, legacy); err != nil {
		return err
	}
	if err := dropCouponMigrationColumn(db, table, legacy); err != nil {
		return err
	}
	return validateCouponMillisColumn(db, table, column.Name, column.Nullable)
}

func swapCouponMigrationColumns(db *gorm.DB, table, current, shadow, legacy string) error {
	if db.Dialector.Name() == "mysql" {
		query := fmt.Sprintf("ALTER TABLE `%s` RENAME COLUMN `%s` TO `%s`, RENAME COLUMN `%s` TO `%s`", table, current, legacy, shadow, current)
		return db.Exec(query).Error
	}
	if err := renameCouponMigrationColumn(db, table, current, legacy); err != nil {
		return err
	}
	return renameCouponMigrationColumn(db, table, shadow, current)
}

func copyLegacyCouponTimes(db *gorm.DB, table, source, target string) error {
	type legacyRow struct {
		ID    string
		Value sql.NullTime
	}
	lastID := ""
	for {
		var rows []legacyRow
		query := fmt.Sprintf("SELECT `id`, `%s` AS `value` FROM `%s` WHERE `id` > ? AND `%s` IS NULL ORDER BY `id` LIMIT 1000", source, table, target)
		if err := db.Raw(query, lastID).Scan(&rows).Error; err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}
		for _, row := range rows {
			var value any
			if row.Value.Valid {
				value = row.Value.Time.UnixMilli()
			}
			update := fmt.Sprintf("UPDATE `%s` SET `%s` = ? WHERE `id` = ?", table, target)
			if err := db.Exec(update, value, row.ID).Error; err != nil {
				return err
			}
			lastID = row.ID
		}
	}
}

type couponTimeParityRow struct {
	ID     string
	Source sql.NullTime  `gorm:"column:source_value"`
	Target sql.NullInt64 `gorm:"column:target_value"`
}

func verifyCouponTimeParity(db *gorm.DB, table, source, target string, nullable bool) error {
	lastID := ""
	for {
		var rows []couponTimeParityRow
		query := fmt.Sprintf("SELECT `id`, `%s` AS `source_value`, `%s` AS `target_value` FROM `%s` WHERE `id` > ? ORDER BY `id` LIMIT 1000", source, target, table)
		if err := db.Raw(query, lastID).Scan(&rows).Error; err != nil {
			return err
		}
		if len(rows) == 0 {
			break
		}
		for _, row := range rows {
			if !couponTimeParityMatches(row, nullable) {
				return fmt.Errorf("coupon time value parity mismatch for %s at %s", source, row.ID)
			}
			lastID = row.ID
		}
	}
	return validateCouponMillisColumn(db, table, target, nullable)
}

func couponTimeParityMatches(row couponTimeParityRow, nullable bool) bool {
	if !nullable && (!row.Source.Valid || !row.Target.Valid) {
		return false
	}
	if row.Source.Valid != row.Target.Valid {
		return false
	}
	return !row.Source.Valid || row.Source.Time.UnixMilli() == row.Target.Int64
}
