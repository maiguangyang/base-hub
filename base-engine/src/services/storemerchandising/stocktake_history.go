package storemerchandising

import (
	"encoding/json"
	"base-engine/gen"
	"gorm.io/gorm"
	"sort"
	"time"
)

func stocktakeCountHistory(tx *gorm.DB, lines []gen.StoreStocktakeLine) (map[string][]StocktakeCountEvent, error) {
	history := make(map[string][]StocktakeCountEvent, len(lines))
	if len(lines) == 0 {
		return history, nil
	}
	ids := make([]string, 0, len(lines))
	for _, line := range lines {
		ids = append(ids, line.ID)
		history[line.ID] = []StocktakeCountEvent{}
	}
	var logs []gen.AuditLog
	if err := tx.Where("resource_type = ? AND resource_id IN ? AND result_code = ?", "storeStocktakeLine", ids, "SUCCESS").
		Order("created_at, id").Find(&logs).Error; err != nil {
		return nil, err
	}
	for _, log := range logs {
		event, ok := stocktakeCountEvent(log)
		if !ok {
			continue
		}
		history[*log.ResourceID] = append(history[*log.ResourceID], event)
	}
	for id := range history {
		sort.SliceStable(history[id], func(i, j int) bool {
			return history[id][i].CountedAt.Before(history[id][j].CountedAt)
		})
	}
	return history, nil
}

func stocktakeCountEvent(log gen.AuditLog) (StocktakeCountEvent, bool) {
	if log.ResourceID == nil || log.ActorAccountID == nil || log.MetadataJSON == nil {
		return StocktakeCountEvent{}, false
	}
	var metadata struct {
		CountedQuantity *int64 `json:"countedQuantity"`
		CountedAt       string `json:"countedAt"`
	}
	if err := json.Unmarshal([]byte(*log.MetadataJSON), &metadata); err != nil || metadata.CountedQuantity == nil {
		return StocktakeCountEvent{}, false
	}
	countedAt, err := time.Parse(time.RFC3339Nano, metadata.CountedAt)
	if err != nil {
		countedAt = time.UnixMilli(log.CreatedAt)
	}
	return StocktakeCountEvent{ActorAccountID: *log.ActorAccountID,
		CountedQuantity: *metadata.CountedQuantity, CountedAt: countedAt}, true
}
