package src

import (
	"base-engine/gen"
	"base-engine/src/services/storemerchandising"
)

func stocktakeView(item *storemerchandising.StocktakeView) (*gen.FranchiseStocktakeView, error) {
	result := &gen.FranchiseStocktakeView{ID: item.ID, StoreID: item.StoreID, Status: item.Status,
		RequestKey: item.RequestKey, StartedAt: item.StartedAt, ReviewedAt: item.ReviewedAt,
		PostedAt: item.PostedAt, CanceledAt: item.CanceledAt,
		InitiatedByAccountID: item.InitiatedByAccountID, PostedByID: item.PostedByID,
		Lines:          make([]*gen.FranchiseStocktakeLineView, 0, len(item.Lines)),
		AddLineChoices: make([]*gen.FranchiseStocktakeAddLineChoice, 0, len(item.AddLineChoices))}
	for _, line := range item.Lines {
		view, err := stocktakeLineView(line)
		if err != nil {
			return nil, err
		}
		result.Lines = append(result.Lines, view)
	}
	for _, choice := range item.AddLineChoices {
		result.AddLineChoices = append(result.AddLineChoices, &gen.FranchiseStocktakeAddLineChoice{
			BatchID: choice.BatchID, PackageID: choice.PackageID,
			PackageName: choice.PackageName, PackageEnabled: choice.PackageEnabled})
	}
	return result, nil
}

func stocktakeLineView(line storemerchandising.StocktakeLineView) (*gen.FranchiseStocktakeLineView, error) {
	counted, err := catalogOptionalInt(line.CountedQuantity)
	if err != nil {
		return nil, err
	}
	snapshot, err := catalogOptionalInt(line.SnapshotQuantity)
	if err != nil {
		return nil, err
	}
	var difference *int
	if line.CountedQuantity != nil && line.SnapshotQuantity != nil {
		value, err := catalogInt(*line.CountedQuantity - *line.SnapshotQuantity)
		if err != nil {
			return nil, err
		}
		difference = &value
	}
	version, err := catalogInt(line.PackageSetVersion)
	if err != nil {
		return nil, err
	}
	history := make([]*gen.FranchiseStocktakeCountEvent, 0, len(line.CountHistory))
	for _, entry := range line.CountHistory {
		quantity, err := catalogInt(entry.CountedQuantity)
		if err != nil {
			return nil, err
		}
		history = append(history, &gen.FranchiseStocktakeCountEvent{
			ActorAccountID: entry.ActorAccountID, CountedQuantity: quantity, CountedAt: entry.CountedAt})
	}
	return &gen.FranchiseStocktakeLineView{ID: line.ID, BatchID: line.BatchID,
		ListingID: line.ListingID, BatchNumber: line.BatchNumber, ExpiresAt: line.ExpiresAt,
		PackageID: line.PackageID, PackageName: line.PackageName, PackageEnabled: line.PackageEnabled,
		PackageSetVersion: version, CountedQuantity: counted, SnapshotQuantity: snapshot,
		Difference: difference, ReasonCode: line.ReasonCode, ReasonNote: line.ReasonNote,
		CountHistory: history, NeedsRecount: line.NeedsRecount}, nil
}
