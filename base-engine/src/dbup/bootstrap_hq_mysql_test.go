package dbup

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestHQInitializedQuotesMySQLReservedKeyColumn(t *testing.T) {
	db, capture := newMySQLDryRunDB(t)
	_, _ = HQInitialized(context.Background(), db)
	if !strings.Contains(capture.sql, "`key` =") {
		t.Fatalf("reserved column is not quoted: %s", capture.sql)
	}
}

func TestDuplicateBootstrapCheckQuotesMySQLReservedKeyColumn(t *testing.T) {
	db, capture := newMySQLDryRunDB(t)
	_ = isDuplicateBootstrap(db, errors.New("duplicate"))
	if !strings.Contains(capture.sql, "`key` =") {
		t.Fatalf("reserved column is not quoted: %s", capture.sql)
	}
}

func newMySQLDryRunDB(t *testing.T) (*gorm.DB, *sqlCaptureLogger) {
	t.Helper()
	capture := &sqlCaptureLogger{}
	sqlDB, err := sql.Open("mysql", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB, SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun: true, DisableAutomaticPing: true, Logger: capture,
	})
	if err != nil {
		t.Fatal(err)
	}
	return db, capture
}

type sqlCaptureLogger struct {
	sql string
}

func (capture *sqlCaptureLogger) LogMode(logger.LogLevel) logger.Interface { return capture }
func (*sqlCaptureLogger) Info(context.Context, string, ...any)             {}
func (*sqlCaptureLogger) Warn(context.Context, string, ...any)             {}
func (*sqlCaptureLogger) Error(context.Context, string, ...any)            {}

func (capture *sqlCaptureLogger) Trace(_ context.Context, _ time.Time, sqlText func() (string, int64), _ error) {
	capture.sql, _ = sqlText()
}
