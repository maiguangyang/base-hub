package dbup

import (
	"fmt"
	"regexp"

	"base-engine/gen"
	"gorm.io/gorm"
)

var customerIndexes = []governanceIndex{
	{Model: &gen.CustomerMember{}, Name: "uidx_customer_member_organization_phone", Columns: "organization_id, phone"},
	{Model: &gen.CustomerMember{}, Name: "uidx_customer_member_request", Columns: "organization_id, request_key"},
	{Model: &gen.CustomerBenefitPolicy{}, Name: "uidx_customer_policy_organization", Columns: "organization_id"},
	{Model: &gen.CustomerDailyPointGrantBudget{}, Name: "uidx_customer_budget_organization_date", Columns: "organization_id, business_date"},
	{Model: &gen.CustomerPointEntry{}, Name: "uidx_customer_point_request", Columns: "source_organization_id, operation_kind, request_key"},
	{Model: &gen.CustomerPointEntry{}, Name: "uidx_customer_point_reversal", Columns: "reverses_id"},
	{Model: &gen.CustomerCouponTemplate{}, Name: "uidx_customer_coupon_template_issuer_code", Columns: "organization_id, issuer_key, code"},
	{Model: &gen.CustomerCouponTemplate{}, Name: "uidx_customer_coupon_template_issuer_request", Columns: "organization_id, issuer_key, request_key"},
	{Model: &gen.CustomerCouponGrant{}, Name: "uidx_customer_coupon_request", Columns: "template_id, member_id, request_key"},
	{Model: &gen.CustomerCouponGrant{}, Name: "uidx_customer_coupon_issuer_request", Columns: "issuer_request_digest"},
	{Model: &gen.CustomerCouponDistributionJob{}, Name: "uidx_customer_coupon_distribution_job_request", Columns: "request_key"},
}

var customerDistributionIndexes = []governanceIndex{
	{Model: &gen.CustomerCouponDistributionJob{}, Name: "idx_customer_coupon_distribution_job_claim", Columns: "status, available_at, lease_expires_at"},
	{Model: &gen.CustomerMember{}, Name: "idx_customer_member_distribution", Columns: "organization_id, status, created_at, id"},
	{Model: &gen.CustomerCouponTemplate{}, Name: "idx_customer_coupon_template_catchup", Columns: "organization_id, enabled, created_at, id"},
}

type couponCheck struct {
	Model            any
	Name, Expression string
}

var customerCouponChecks = []couponCheck{
	{&gen.CustomerCouponTemplate{}, "chk_coupon_template_effective_ms", "`effective_at` BETWEEN 1000000000000 AND 9999999999999"},
	{&gen.CustomerCouponTemplate{}, "chk_coupon_template_distribution_end_ms", "`distribution_ends_at` IS NULL OR `distribution_ends_at` BETWEEN 1000000000000 AND 9999999999999"},
	{&gen.CustomerCouponTemplate{}, "chk_coupon_template_days", "`days_after_activation` BETWEEN 1 AND 365"},
	{&gen.CustomerCouponTemplate{}, "chk_coupon_template_total_limit", "`total_issue_limit` >= 0"},
	{&gen.CustomerCouponTemplate{}, "chk_coupon_template_distribution_order", "`distribution_ends_at` IS NULL OR `distribution_ends_at` > `effective_at`"},
	{&gen.CustomerCouponGrant{}, "chk_coupon_grant_issued_ms", "`issued_at` BETWEEN 1000000000000 AND 9999999999999"},
	{&gen.CustomerCouponGrant{}, "chk_coupon_grant_activated_ms", "`activated_at` IS NULL OR `activated_at` BETWEEN 1000000000000 AND 9999999999999"},
	{&gen.CustomerCouponGrant{}, "chk_coupon_grant_expires_ms", "`expires_at` IS NULL OR `expires_at` BETWEEN 1000000000000 AND 9999999999999"},
	{&gen.CustomerCouponGrant{}, "chk_coupon_grant_revoked_ms", "`revoked_at` IS NULL OR `revoked_at` BETWEEN 1000000000000 AND 9999999999999"},
	{&gen.CustomerCouponDistributionJob{}, "chk_coupon_job_available_ms", "`available_at` BETWEEN 1000000000000 AND 9999999999999"},
	{&gen.CustomerCouponDistributionJob{}, "chk_coupon_job_lease_ms", "`lease_expires_at` IS NULL OR `lease_expires_at` BETWEEN 1000000000000 AND 9999999999999"},
	{&gen.CustomerCouponDistributionJob{}, "chk_coupon_job_cursor_ms", "`cursor_created_at` IS NULL OR `cursor_created_at` BETWEEN 1000000000000 AND 9999999999999"},
	{&gen.CustomerCouponDistributionJob{}, "chk_coupon_job_target", "(`kind` = 'TEMPLATE_FANOUT' AND `template_id` IS NOT NULL AND `member_id` IS NULL) OR (`kind` = 'MEMBER_CATCHUP' AND `member_id` IS NOT NULL AND `template_id` IS NULL)"},
}

// EnsureCustomerIndexes installs the database constraints that generated
// relationship fields cannot express in handwritten schema.
func EnsureCustomerIndexes(db *gorm.DB) error {
	dialect := db.Dialector.Name()
	if dialect != "sqlite" && dialect != "mysql" {
		return fmt.Errorf("unsupported customer index dialect: %s", dialect)
	}
	if err := ensureCouponIssuerKey(db); err != nil {
		return err
	}
	if err := ensureCustomerIndexSets(db); err != nil {
		return err
	}
	if err := ensureCustomerCouponChecks(db); err != nil {
		return err
	}
	return dropLegacyCustomerCouponIndexes(db)
}

func ensureCustomerIndexSets(db *gorm.DB) error {
	for _, index := range customerIndexes {
		if err := ensureGovernanceIndex(db, index); err != nil {
			return err
		}
	}
	for _, index := range customerDistributionIndexes {
		if err := ensureCustomerLookupIndex(db, index); err != nil {
			return err
		}
	}
	return nil
}

func dropLegacyCustomerCouponIndexes(db *gorm.DB) error {
	for _, old := range []string{"uidx_customer_coupon_template_code", "uidx_customer_coupon_template_request"} {
		if db.Migrator().HasIndex(&gen.CustomerCouponTemplate{}, old) {
			if err := db.Migrator().DropIndex(&gen.CustomerCouponTemplate{}, old); err != nil {
				return err
			}
		}
	}
	return nil
}

func ensureCustomerLookupIndex(db *gorm.DB, index governanceIndex) error {
	if db.Migrator().HasIndex(index.Model, index.Name) {
		return nil
	}
	table, err := governanceTableName(db, index.Model)
	if err != nil {
		return err
	}
	return db.Exec(fmt.Sprintf("CREATE INDEX `%s` ON `%s` (%s)", index.Name, table, index.Columns)).Error
}

func ensureCustomerCouponChecks(db *gorm.DB) error {
	for _, check := range customerCouponChecks {
		if err := ensureCustomerCouponCheck(db, check); err != nil {
			return err
		}
	}
	return nil
}

func ensureCustomerCouponCheck(db *gorm.DB, check couponCheck) error {
	table, err := governanceTableName(db, check.Model)
	if err != nil {
		return err
	}
	if db.Dialector.Name() == "sqlite" {
		expression := regexp.MustCompile("`([^`]+)`").ReplaceAllString(check.Expression, "`NEW`.`$1`")
		insert := fmt.Sprintf("CREATE TRIGGER IF NOT EXISTS `%s_insert` BEFORE INSERT ON `%s` WHEN NOT (%s) BEGIN SELECT RAISE(ABORT, '%s'); END", check.Name, table, expression, check.Name)
		update := fmt.Sprintf("CREATE TRIGGER IF NOT EXISTS `%s_update` BEFORE UPDATE ON `%s` WHEN NOT (%s) BEGIN SELECT RAISE(ABORT, '%s'); END", check.Name, table, expression, check.Name)
		if err := db.Exec(insert).Error; err != nil {
			return err
		}
		return db.Exec(update).Error
	}
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM information_schema.TABLE_CONSTRAINTS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND CONSTRAINT_NAME = ? AND CONSTRAINT_TYPE = 'CHECK'", table, check.Name).Scan(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return nil
	}
	return db.Exec(fmt.Sprintf("ALTER TABLE `%s` ADD CONSTRAINT `%s` CHECK (%s)", table, check.Name, check.Expression)).Error
}

func ensureCouponIssuerKey(db *gorm.DB) error {
	table, err := governanceTableName(db, &gen.CustomerCouponTemplate{})
	if err != nil {
		return err
	}
	if db.Migrator().HasColumn(table, "issuer_key") {
		return nil
	}
	query := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `issuer_key` VARCHAR(36) GENERATED ALWAYS AS (CASE WHEN `issuer_scope` = 'STORE' THEN COALESCE(`applicable_store_id`, '') ELSE 'HQ' END) VIRTUAL", table)
	return db.Exec(query).Error
}
