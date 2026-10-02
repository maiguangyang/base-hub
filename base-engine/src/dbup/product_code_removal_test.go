package dbup

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestCatalogCodeRemovalPreservesRows(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"products", "product_skus"} {
		if err := db.Exec("CREATE TABLE " + table + " (id TEXT PRIMARY KEY, code TEXT NOT NULL, name TEXT)").Error; err != nil {
			t.Fatal(err)
		}
		if err := db.Exec("INSERT INTO " + table + " VALUES ('kept-id', 'old-code', 'kept-name')").Error; err != nil {
			t.Fatal(err)
		}
	}
	assertCodeRemoval(t, db)
	assertCodeRemoval(t, db)
	for _, table := range []string{"products", "product_skus"} {
		if db.Migrator().HasColumn(table, "code") {
			t.Fatal("obsolete code remains")
		}
		var name string
		if err := db.Table(table).Select("name").Where("id = ?", "kept-id").Scan(&name).Error; err != nil || name != "kept-name" {
			t.Fatalf("row changed: %s %v", name, err)
		}
	}
}

func assertCodeRemoval(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := ensureNoCatalogCodes(db); err != nil {
		t.Fatal(err)
	}
}
