package dbup

import (
	"fmt"

	"base-engine/gen"
	"gorm.io/gorm"
)

var productIndexes = []governanceIndex{
	{Model: &gen.ProductBrand{}, Name: "uidx_product_brand_org_name", Columns: "organization_id, name"},
	{Model: &gen.SpecificationDefinition{}, Name: "uidx_specification_org_name", Columns: "organization_id, name"},
	{Model: &gen.SpecificationValue{}, Name: "uidx_specification_value_name", Columns: "specification_id, name"},
	{Model: &gen.ProductSpecificationChoice{}, Name: "uidx_product_specification_choice", Columns: "product_id, value_id"},
	{Model: &gen.ProductSkuSpecificationValue{}, Name: "uidx_sku_specification_value", Columns: "sku_id, value_id"},
	{Model: &gen.ProductCategory{}, Name: "uidx_product_category_parent_name", Columns: "organization_id, parent_scope, name"},
	{Model: &gen.ProductPackage{}, Name: "uidx_product_package_active_barcode", Columns: "active_barcode"},
	{Model: &gen.ProductPackage{}, Name: "uidx_product_package_base_version", Columns: "base_sku_scope, package_set_version"},
	{Model: &gen.StoreListing{}, Name: "uidx_store_listing_store_sku", Columns: "store_id, sku_id"},
	{Model: &gen.StorePackageOffer{}, Name: "uidx_store_offer_listing_package", Columns: "listing_id, package_id"},
	{Model: &gen.StoreStockBalance{}, Name: "uidx_store_balance_batch_package", Columns: "batch_id, package_id"},
	{Model: &gen.StoreInventoryBatch{}, Name: "uidx_store_batch_listing_number", Columns: "listing_id, batch_number"},
	{Model: &gen.StorePromotion{}, Name: "uidx_store_promotion_rule_version", Columns: "store_id, rule_key, version"},
	{Model: &gen.StoreStockMovement{}, Name: "uidx_store_movement_request", Columns: "store_id, kind, request_key"},
	{Model: &gen.StoreStocktake{}, Name: "uidx_store_stocktake_request", Columns: "store_id, request_key"},
	{Model: &gen.StoreStocktakeLine{}, Name: "uidx_stocktake_line_target", Columns: "stocktake_id, batch_id, package_id"},
}

// EnsureProductIndexes installs the cross-entity natural-key constraints.
func EnsureProductIndexes(db *gorm.DB) error {
	dialect := db.Dialector.Name()
	if dialect != "sqlite" && dialect != "mysql" {
		return fmt.Errorf("unsupported product index dialect: %s", dialect)
	}
	if err := ensureNoCatalogCodes(db); err != nil {
		return err
	}
	if err := ensureCategoryParentScope(db); err != nil {
		return err
	}
	if err := ensurePackageScopes(db); err != nil {
		return err
	}
	for _, index := range productIndexes {
		if err := ensureGovernanceIndex(db, index); err != nil {
			return err
		}
	}
	return nil
}

func ensurePackageScopes(db *gorm.DB) error {
	table, err := governanceTableName(db, &gen.ProductPackage{})
	if err != nil {
		return err
	}
	legacy := "uidx_product_package_barcode"
	if db.Migrator().HasIndex(&gen.ProductPackage{}, legacy) {
		if err := db.Migrator().DropIndex(&gen.ProductPackage{}, legacy); err != nil {
			return err
		}
	}
	columns := []struct{ name, expression string }{
		{"active_barcode", "CASE WHEN `enabled` THEN `barcode` ELSE NULL END"},
		{"base_sku_scope", "CASE WHEN `contains_package_id` IS NULL THEN `sku_id` ELSE NULL END"},
	}
	for _, column := range columns {
		if db.Migrator().HasColumn(table, column.name) {
			continue
		}
		statement := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` VARCHAR(64) GENERATED ALWAYS AS (%s) VIRTUAL", table, column.name, column.expression)
		if err := db.Exec(statement).Error; err != nil {
			return err
		}
	}
	return nil
}

func ensureCategoryParentScope(db *gorm.DB) error {
	table, err := governanceTableName(db, &gen.ProductCategory{})
	if err != nil {
		return err
	}
	if db.Migrator().HasColumn(table, "parent_scope") {
		return nil
	}
	return db.Exec(fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `parent_scope` VARCHAR(36) GENERATED ALWAYS AS (COALESCE(`parent_id`, '')) VIRTUAL", table)).Error
}
