package dbup

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func validateCouponMigrationDialect(db *gorm.DB) error {
	switch db.Dialector.Name() {
	case "sqlite":
		return nil
	case "mysql":
		var version string
		if err := db.Raw("SELECT VERSION()").Scan(&version).Error; err != nil {
			return err
		}
		if strings.Contains(strings.ToLower(version), "mariadb") || !mysqlVersionAtLeast(version, 8, 0, 16) {
			return fmt.Errorf("coupon migration requires Oracle MySQL 8.0.16+, got %s", version)
		}
		return nil
	default:
		return fmt.Errorf("unsupported coupon migration dialect: %s", db.Dialector.Name())
	}
}

func mysqlVersionAtLeast(value string, major, minor, patch int) bool {
	match := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)`).FindStringSubmatch(value)
	if len(match) != 4 {
		return false
	}
	parts := make([]int, 3)
	for index := range parts {
		parts[index], _ = strconv.Atoi(match[index+1])
	}
	return parts[0] > major || parts[0] == major && (parts[1] > minor || parts[1] == minor && parts[2] >= patch)
}

func couponMigrationTableName(db *gorm.DB, model any) (string, error) {
	table, err := governanceTableName(db, model)
	if err != nil {
		return "", err
	}
	if !regexp.MustCompile(`^[A-Za-z0-9_]+$`).MatchString(table) {
		return "", fmt.Errorf("unsafe coupon migration table name %q", table)
	}
	return table, nil
}

func couponColumnIsInteger(db *gorm.DB, table, name string) (bool, error) {
	columns, err := db.Migrator().ColumnTypes(table)
	if err != nil {
		return false, err
	}
	for _, column := range columns {
		if strings.EqualFold(column.Name(), name) {
			typeName := strings.ToLower(column.DatabaseTypeName())
			return strings.Contains(typeName, "int"), nil
		}
	}
	return false, fmt.Errorf("coupon column %s.%s not found", table, name)
}

func validateCouponMillisColumn(db *gorm.DB, table, column string, nullable bool) error {
	predicate := fmt.Sprintf("`%s` < ? OR `%s` > ?", column, column)
	if !nullable {
		predicate = fmt.Sprintf("`%s` IS NULL OR %s", column, predicate)
	}
	var invalid int64
	if err := db.Table(table).Where(predicate, minimumCouponMillis, maximumCouponMillis).Count(&invalid).Error; err != nil {
		return err
	}
	if invalid != 0 {
		return fmt.Errorf("invalid coupon milliseconds in %s.%s", table, column)
	}
	return nil
}
