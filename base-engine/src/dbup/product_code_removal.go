package dbup

import (
	"fmt"
	"base-engine/gen"
	"gorm.io/gorm"
)

// Remove only obsolete catalog codes; UUIDs and business rows remain intact.
func ensureNoCatalogCodes(db *gorm.DB) error {
	for _, model := range []any{&gen.Product{}, &gen.ProductSku{}} {
		if !db.Migrator().HasTable(model) || !db.Migrator().HasColumn(model, "code") {
			continue
		}
		table, err := governanceTableName(db, model)
		if err != nil {
			return err
		}
		index := "idx_" + table + "_code"
		if db.Migrator().HasIndex(model, index) {
			if err := db.Migrator().DropIndex(model, index); err != nil {
				return err
			}
		}
		if err := db.Exec(fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `code`", table)).Error; err != nil {
			return err
		}
	}
	return nil
}
