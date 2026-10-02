package storemerchandising

import "gorm.io/gorm"

func dateFilteredStocktakeQuery(query *gorm.DB, filter StocktakeFilter) *gorm.DB {
	if filter.From != nil {
		query = query.Where("started_at >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("started_at <= ?", *filter.To)
	}
	return query
}
