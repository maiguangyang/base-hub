package storemerchandising

import (
	"fmt"
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func assertBatchBalanced(tx *gorm.DB, batchID string) error {
	type sum struct {
		PackageID string `gorm:"column:package_id"`
		Quantity  int64  `gorm:"column:quantity"`
	}
	var sums []sum
	table := tx.NamingStrategy.TableName("StoreStockMovement")
	query := fmt.Sprintf(`SELECT package_id, SUM(delta) AS quantity FROM (
		SELECT target_package_id AS package_id, target_quantity AS delta FROM %s WHERE batch_id = ? AND target_package_id IS NOT NULL
		UNION ALL
		SELECT source_package_id AS package_id, -source_quantity AS delta FROM %s WHERE batch_id = ? AND source_package_id IS NOT NULL
	) ledger GROUP BY package_id`, table, table)
	if err := tx.Raw(query, batchID, batchID).Scan(&sums).Error; err != nil {
		return err
	}
	expected := make(map[string]int64, len(sums))
	for _, item := range sums {
		expected[item.PackageID] = item.Quantity
	}
	var balances []gen.StoreStockBalance
	if err := tx.Where("batch_id = ?", batchID).Find(&balances).Error; err != nil {
		return err
	}
	for _, item := range balances {
		if item.Quantity < 0 || expected[item.PackageID] != item.Quantity {
			return auth.NewError(auth.CodeConflict)
		}
		delete(expected, item.PackageID)
	}
	for _, quantity := range expected {
		if quantity != 0 {
			return auth.NewError(auth.CodeConflict)
		}
	}
	return nil
}
