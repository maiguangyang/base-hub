package dbup

import (
	"context"
	"strings"
	"testing"
	"time"

	"base-engine/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type couponMigrationQueryCounter struct {
	logger.Interface
	issuedAtReads int
}

func (counter *couponMigrationQueryCounter) Trace(ctx context.Context, begin time.Time, query func() (string, int64), err error) {
	sql, _ := query()
	if strings.Contains(sql, "SELECT `id`, `issued_at`, `created_at`") {
		counter.issuedAtReads++
	}
}

func TestPrepareCustomerCouponTimeColumnsHandlesAbsentAndTemplateOnlyStates(t *testing.T) {
	db := openTestDB(t)
	if err := PrepareCustomerCouponTimeColumns(db); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`CREATE TABLE customer_coupon_templates (id TEXT PRIMARY KEY, created_at INTEGER NOT NULL)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`INSERT INTO customer_coupon_templates (id, created_at) VALUES ('template', 1700000000123)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := PrepareCustomerCouponTimeColumns(db); err != nil {
		t.Fatal(err)
	}
	var effectiveAt int64
	if err := db.Raw(`SELECT effective_at FROM customer_coupon_templates WHERE id = 'template'`).Scan(&effectiveAt).Error; err != nil || effectiveAt != 1_700_000_000_123 {
		t.Fatalf("effective_at = %d, %v", effectiveAt, err)
	}
	if err := PrepareCustomerCouponTimeColumns(db); err != nil {
		t.Fatalf("rerun: %v", err)
	}
}

func TestPrepareCustomerCouponTimeColumnsConvertsLegacyGrantTimes(t *testing.T) {
	db := openTestDB(t)
	if err := db.Exec(`CREATE TABLE customer_coupon_grants (id TEXT PRIMARY KEY, issued_at DATETIME NOT NULL, activated_at DATETIME, expires_at DATETIME, revoked_at DATETIME, created_at INTEGER NOT NULL)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`INSERT INTO customer_coupon_grants (id, issued_at, activated_at, created_at) VALUES ('grant', '2023-11-14 22:13:20.123', '2023-11-15 22:13:20.456', 1700000000123)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := PrepareCustomerCouponTimeColumns(db); err != nil {
		t.Fatal(err)
	}
	var row struct {
		IssuedAt, ActivatedAt int64
		ExpiresAt, RevokedAt  *int64
	}
	if err := db.Raw(`SELECT issued_at, activated_at, expires_at, revoked_at FROM customer_coupon_grants WHERE id = 'grant'`).Scan(&row).Error; err != nil {
		t.Fatal(err)
	}
	if row.IssuedAt != 1_700_000_000_123 || row.ActivatedAt != 1_700_086_400_456 || row.ExpiresAt != nil || row.RevokedAt != nil {
		t.Fatalf("migrated row = %+v", row)
	}
	if err := PrepareCustomerCouponTimeColumns(db); err != nil {
		t.Fatalf("grant migration rerun: %v", err)
	}
}

func TestPrepareCustomerCouponTimeColumnsRecoversShadowSwap(t *testing.T) {
	db := openTestDB(t)
	if err := db.Exec(`CREATE TABLE customer_coupon_grants (id TEXT PRIMARY KEY, issued_at_legacy DATETIME NOT NULL, issued_at_ms BIGINT, created_at INTEGER NOT NULL)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`INSERT INTO customer_coupon_grants (id, issued_at_legacy, issued_at_ms, created_at) VALUES ('grant', '2023-11-14 22:13:20.123', 1700000000123, 1700000000123)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := PrepareCustomerCouponTimeColumns(db); err != nil {
		t.Fatal(err)
	}
	var issuedAt int64
	if err := db.Raw(`SELECT issued_at FROM customer_coupon_grants WHERE id = 'grant'`).Scan(&issuedAt).Error; err != nil || issuedAt != 1_700_000_000_123 {
		t.Fatalf("recovered issued_at = %d, %v", issuedAt, err)
	}
	if db.Migrator().HasColumn("customer_coupon_grants", "issued_at_legacy") || db.Migrator().HasColumn("customer_coupon_grants", "issued_at_ms") {
		t.Fatal("shadow swap columns were not cleaned up")
	}
}

func TestPrepareCustomerCouponTimeColumnsRespectsTablePrefix(t *testing.T) {
	db := openPrefixedCouponMigrationDB(t)
	if err := PrepareCustomerCouponTimeColumns(db); err != nil {
		t.Fatal(err)
	}
	assertPrefixedCouponMigrationValues(t, db)
}

func openPrefixedCouponMigrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: "prefix_"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`CREATE TABLE prefix_customer_coupon_templates (id TEXT PRIMARY KEY, created_at INTEGER NOT NULL)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`INSERT INTO prefix_customer_coupon_templates (id, created_at) VALUES ('template', 1700000000123)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`CREATE TABLE prefix_customer_coupon_grants (id TEXT PRIMARY KEY, issued_at DATETIME NOT NULL, activated_at DATETIME, expires_at DATETIME, revoked_at DATETIME, created_at INTEGER NOT NULL)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`INSERT INTO prefix_customer_coupon_grants (id, issued_at, created_at) VALUES ('grant', '2023-11-14 22:13:20.123', 1700000000123)`).Error; err != nil {
		t.Fatal(err)
	}
	return db
}

func assertPrefixedCouponMigrationValues(t *testing.T, db *gorm.DB) {
	t.Helper()
	var effectiveAt int64
	if err := db.Raw(`SELECT effective_at FROM prefix_customer_coupon_templates WHERE id = 'template'`).Scan(&effectiveAt).Error; err != nil || effectiveAt != 1_700_000_000_123 {
		t.Fatalf("prefixed effective_at = %d, %v", effectiveAt, err)
	}
	var issuedAt int64
	if err := db.Raw(`SELECT issued_at FROM prefix_customer_coupon_grants WHERE id = 'grant'`).Scan(&issuedAt).Error; err != nil || issuedAt != 1_700_000_000_123 {
		t.Fatalf("prefixed issued_at = %d, %v", issuedAt, err)
	}
}

func TestVerifyCouponTimeParityRejectsValueMismatch(t *testing.T) {
	db := openTestDB(t)
	if err := db.Exec(`CREATE TABLE customer_coupon_grants (id TEXT PRIMARY KEY, issued_at DATETIME NOT NULL, issued_at_ms BIGINT NOT NULL)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`INSERT INTO customer_coupon_grants (id, issued_at, issued_at_ms) VALUES ('grant', '2023-11-14 22:13:20.123', 1700000000124)`).Error; err != nil {
		t.Fatal(err)
	}
	if err := verifyCouponTimeParity(db, "customer_coupon_grants", "issued_at", "issued_at_ms", false); err == nil {
		t.Fatal("mismatched legacy and millisecond values were accepted")
	}
}

func TestPreflightLegacyIssuedAtReadsBoundedBatches(t *testing.T) {
	db := openTestDB(t)
	if err := db.Exec(`CREATE TABLE customer_coupon_grants (id TEXT PRIMARY KEY, issued_at DATETIME NOT NULL, created_at INTEGER NOT NULL)`).Error; err != nil {
		t.Fatal(err)
	}
	seed := `WITH RECURSIVE rows(value) AS (SELECT 1 UNION ALL SELECT value + 1 FROM rows WHERE value < 1001)
		INSERT INTO customer_coupon_grants (id, issued_at, created_at)
		SELECT printf('%04d', value), '2023-11-14 22:13:20.123', 1700000000123 FROM rows`
	if err := db.Exec(seed).Error; err != nil {
		t.Fatal(err)
	}
	counter := &couponMigrationQueryCounter{Interface: logger.Discard}
	if err := preflightLegacyIssuedAt(db.Session(&gorm.Session{Logger: counter}), "customer_coupon_grants"); err != nil {
		t.Fatal(err)
	}
	if counter.issuedAtReads < 2 {
		t.Fatalf("legacy preflight reads = %d, want bounded batches", counter.issuedAtReads)
	}
}

func TestEnsureCustomerCouponConstraintsAndHistoricalJobs(t *testing.T) {
	db := openTestDB(t)
	if err := migrateGeneratedForSQLite(db, &gen.CustomerMember{}, &gen.CustomerBenefitPolicy{},
		&gen.CustomerDailyPointGrantBudget{}, &gen.CustomerPointEntry{}, &gen.CustomerCouponTemplate{},
		&gen.CustomerCouponGrant{}, &gen.CustomerCouponDistributionJob{}); err != nil {
		t.Fatal(err)
	}
	if err := EnsureCustomerIndexes(db); err != nil {
		t.Fatal(err)
	}
	assertCouponConstraintRejections(t, db)
	nowMillis, future := seedHistoricalCouponTemplates(t, db)
	if err := EnsureCustomerCouponDistributionJobs(db, nowMillis); err != nil {
		t.Fatal(err)
	}
	if err := EnsureCustomerCouponDistributionJobs(db, nowMillis+10_000); err != nil {
		t.Fatal(err)
	}
	assertHistoricalCouponJobs(t, db, nowMillis, future)
}

func assertHistoricalCouponJobs(t *testing.T, db *gorm.DB, nowMillis, future int64) {
	t.Helper()
	var jobs []gen.CustomerCouponDistributionJob
	if err := db.Order("request_key").Find(&jobs).Error; err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 2 || jobs[0].AvailableAt != future || jobs[1].AvailableAt != nowMillis {
		t.Fatalf("historical jobs = %+v", jobs)
	}
	var grants int64
	if err := db.Model(&gen.CustomerCouponGrant{}).Count(&grants).Error; err != nil || grants != 0 {
		t.Fatalf("migration grants = %d, %v", grants, err)
	}
}

func seedHistoricalCouponTemplates(t *testing.T, db *gorm.DB) (int64, int64) {
	t.Helper()
	nowMillis := int64(1_800_000_000_000)
	future := nowMillis + 5_000
	for _, template := range []gen.CustomerCouponTemplate{
		{ID: "enabled-now", OrganizationID: "hq", Code: "NOW", RequestKey: "now", Title: "Now", AmountFen: 1, DaysAfterActivation: 1, EffectiveAt: nowMillis - 1, PerMemberLimit: 1, Enabled: true},
		{ID: "enabled-future", OrganizationID: "hq", Code: "FUTURE", RequestKey: "future", Title: "Future", AmountFen: 1, DaysAfterActivation: 365, EffectiveAt: future, PerMemberLimit: 1, Enabled: true},
		{ID: "disabled", OrganizationID: "hq", Code: "OFF", RequestKey: "off", Title: "Off", AmountFen: 1, DaysAfterActivation: 1, EffectiveAt: nowMillis, PerMemberLimit: 1, Enabled: false},
	} {
		if err := db.Create(&template).Error; err != nil {
			t.Fatal(err)
		}
	}
	return nowMillis, future
}

func assertCouponConstraintRejections(t *testing.T, db *gorm.DB) {
	t.Helper()
	templateInsert := `INSERT INTO customer_coupon_templates (id, code, request_key, title, amount_fen, min_spend_fen, days_after_activation, effective_at, per_member_limit, total_issue_limit, issued_count, enabled, issuer_scope, organization_id, created_at) VALUES (?, ?, ?, ?, 1, 0, ?, ?, 1, ?, 0, 1, 'HEADQUARTERS', 'hq', 1800000000000)`
	for _, item := range []struct {
		id          string
		days, total int64
		effectiveAt int64
	}{
		{"bad-time", 1, 0, 999_999_999_999},
		{"bad-days", 366, 0, 1_800_000_000_000},
		{"bad-total", 1, -1, 1_800_000_000_000},
	} {
		if err := db.Exec(templateInsert, item.id, item.id, item.id, item.id, item.days, item.effectiveAt, item.total).Error; err == nil {
			t.Fatalf("invalid template accepted: %s", item.id)
		}
	}
	if err := db.Exec(`INSERT INTO customer_coupon_templates (id, code, request_key, title, amount_fen, min_spend_fen, days_after_activation, effective_at, distribution_ends_at, per_member_limit, total_issue_limit, issued_count, enabled, issuer_scope, organization_id, created_at) VALUES ('bad-order', 'bad-order', 'bad-order', 'bad-order', 1, 0, 1, 1800000000000, 1800000000000, 1, 0, 0, 1, 'HEADQUARTERS', 'hq', 1800000000000)`).Error; err == nil {
		t.Fatal("invalid cutoff ordering accepted")
	}
	if err := db.Exec(`INSERT INTO customer_coupon_grants (id, status, amount_fen, min_spend_fen, days_after_activation, issued_at, request_key, member_id, template_id, created_at) VALUES ('bad-grant', 'PENDING_ACTIVATION', 1, 0, 1, 999999999999, 'bad-grant', 'member', 'template', 1800000000000)`).Error; err == nil {
		t.Fatal("invalid grant timestamp accepted")
	}
	if err := db.Exec(`INSERT INTO customer_coupon_distribution_jobs (id, kind, status, request_key, available_at, attempts, template_id, member_id, created_at) VALUES ('bad-job', 'TEMPLATE_FANOUT', 'PENDING', 'bad-job', 1800000000000, 0, NULL, 'member', 1800000000000)`).Error; err == nil {
		t.Fatal("invalid job target accepted")
	}
	if err := db.Exec(`INSERT INTO customer_coupon_distribution_jobs (id, kind, status, request_key, available_at, attempts, template_id, created_at) VALUES ('bad-job-time', 'TEMPLATE_FANOUT', 'PENDING', 'bad-job-time', 999999999999, 0, 'template', 1800000000000)`).Error; err == nil {
		t.Fatal("invalid job timestamp accepted")
	}
}
