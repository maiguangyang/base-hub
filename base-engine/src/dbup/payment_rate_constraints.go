package dbup

import (
	"fmt"

	"base-engine/gen"
	"gorm.io/gorm"
)

type globalPaymentRateConstraint struct {
	RatePpm int `gorm:"column:rate_ppm;check:chk_global_payment_rate,rate_ppm >= 0 AND rate_ppm <= 1000000"`
}

func (globalPaymentRateConstraint) TableName() string { return "global_payment_configs" }

type franchisePaymentRateConstraint struct {
	RatePpm int `gorm:"column:rate_ppm;check:chk_franchise_payment_rate,rate_ppm >= 0 AND rate_ppm <= 1000000"`
}

func (franchisePaymentRateConstraint) TableName() string { return "franchise_payment_configs" }

type storePaymentRateConstraint struct {
	RatePpm int `gorm:"column:rate_ppm;check:chk_store_payment_rate,rate_ppm >= 0 AND rate_ppm <= 1000000"`
}

func (storePaymentRateConstraint) TableName() string { return "store_payment_configs" }

func ensurePaymentRateConstraints(db *gorm.DB) error {
	for _, item := range []struct {
		model, tableModel any
		name              string
	}{
		{&globalPaymentRateConstraint{}, &gen.GlobalPaymentConfig{}, "chk_global_payment_rate"},
		{&franchisePaymentRateConstraint{}, &gen.FranchisePaymentConfig{}, "chk_franchise_payment_rate"},
		{&storePaymentRateConstraint{}, &gen.StorePaymentConfig{}, "chk_store_payment_rate"},
	} {
		table, err := governanceTableName(db, item.tableModel)
		if err != nil {
			return err
		}
		constraintDB := db.Session(&gorm.Session{NewDB: true}).Table(table)
		if constraintDB.Migrator().HasConstraint(item.model, item.name) {
			continue
		}
		if err := constraintDB.Migrator().CreateConstraint(item.model, item.name); err != nil {
			return fmt.Errorf("create %s constraint: %w", item.name, err)
		}
	}
	return nil
}
