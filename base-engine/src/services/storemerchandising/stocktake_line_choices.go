package storemerchandising

import (
	"base-engine/gen"
	"gorm.io/gorm"
)

func populateStocktakeAddLineChoices(tx *gorm.DB, views []StocktakeView) error {
	batchIDs := countingStocktakeBatchIDs(views)
	if len(batchIDs) == 0 {
		return nil
	}
	var candidates []StocktakeAddLineChoice
	err := tx.Table("store_inventory_batches AS batch").
		Select("batch.id AS batch_id, pack.id AS package_id, pack.name AS package_name, pack.enabled AS package_enabled").
		Joins("JOIN store_listings AS listing ON listing.id = batch.listing_id").
		Joins("JOIN product_packages AS pack ON pack.sku_id = listing.sku_id").
		Joins("LEFT JOIN store_stock_balances AS balance ON balance.batch_id = batch.id AND balance.package_id = pack.id").
		Where("batch.id IN ? AND (pack.enabled = ? OR balance.id IS NOT NULL)", batchIDs, true).
		Order("batch.id, pack.name, pack.id").Scan(&candidates).Error
	if err != nil {
		return err
	}
	for index := range views {
		if views[index].Status == gen.StocktakeStatusCounting {
			views[index].AddLineChoices = append(views[index].AddLineChoices,
				missingStocktakeLineChoices(views[index].Lines, candidates)...)
		}
	}
	return nil
}

func countingStocktakeBatchIDs(views []StocktakeView) []string {
	batchIDs := make([]string, 0)
	seenBatches := make(map[string]bool)
	for _, view := range views {
		if view.Status != gen.StocktakeStatusCounting {
			continue
		}
		for _, line := range view.Lines {
			if !seenBatches[line.BatchID] {
				seenBatches[line.BatchID] = true
				batchIDs = append(batchIDs, line.BatchID)
			}
		}
	}
	return batchIDs
}

func missingStocktakeLineChoices(lines []StocktakeLineView, candidates []StocktakeAddLineChoice) []StocktakeAddLineChoice {
	existing := make(map[string]map[string]bool)
	for _, line := range lines {
		if existing[line.BatchID] == nil {
			existing[line.BatchID] = make(map[string]bool)
		}
		existing[line.BatchID][line.PackageID] = true
	}
	choices := make([]StocktakeAddLineChoice, 0)
	for _, item := range candidates {
		if existing[item.BatchID] != nil && !existing[item.BatchID][item.PackageID] {
			choices = append(choices, item)
		}
	}
	return choices
}
