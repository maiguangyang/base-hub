package storemerchandising

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"base-engine/auth"
	"gorm.io/gorm"
)

// 调用者持有门店行锁直至入库事务提交，所有商品共享当天序号。
func nextStockBatchNumber(tx *gorm.DB, storeID string, now time.Time) (string, error) {
	zone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return "", err
	}
	prefix := "B" + now.In(zone).Format("20060102") + "-"
	var suffixes []string
	// 前缀匹配遵循数据库排序规则；按字符截取后缀，兼容大小写和重音等价前缀。
	err = tx.Table("store_inventory_batches AS b").
		Joins("JOIN store_listings AS l ON l.id = b.listing_id").
		Where("l.store_id = ? AND b.batch_number LIKE ?", storeID, prefix+"%").
		Select("SUBSTR(b.batch_number, ?) AS suffix", len(prefix)+1).
		Scan(&suffixes).Error
	if err != nil {
		return "", err
	}
	var maximum uint64
	for _, suffix := range suffixes {
		n, parseErr := strconv.ParseUint(suffix, 10, 64)
		if parseErr == nil && n > maximum {
			maximum = n
		}
	}
	return availableStockBatchNumber(tx, storeID, prefix, maximum)
}

func availableStockBatchNumber(tx *gorm.DB, storeID, prefix string, maximum uint64) (string, error) {
	for maximum < math.MaxUint64 {
		maximum++
		number := fmt.Sprintf("%s%04d", prefix, maximum)
		var count int64
		// 用与入库查询相同的数据库比较规则兜底，跳过全角数字等无法解析的等价编号。
		err := tx.Table("store_inventory_batches AS b").
			Joins("JOIN store_listings AS l ON l.id = b.listing_id").
			Where("l.store_id = ? AND b.batch_number = ?", storeID, number).Count(&count).Error
		if err != nil {
			return "", err
		}
		if count == 0 {
			return number, nil
		}
	}
	return "", auth.NewError(auth.CodeConflict)
}
