package membership

import (
	"database/sql"
	"strings"
	"testing"

	"base-engine/gen"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestInviteAccountQueryLocksExistingAccount(t *testing.T) {
	sqlDB, err := sql.Open("mysql", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	database, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	query := database.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return inviteAccountQuery(tx, "13800000000").First(&gen.Account{})
	})
	if !strings.Contains(query, "FOR UPDATE") {
		t.Fatalf("existing account is not serialized during invite: %s", query)
	}
}
