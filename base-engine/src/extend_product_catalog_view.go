package src

import (
	"math"

	"base-engine/auth"
	"base-engine/gen"
)

func catalogInt(value int64) (int, error) {
	if value < math.MinInt32 || value > math.MaxInt32 {
		return 0, auth.NewError(auth.CodeInternalError)
	}
	return int(value), nil
}
func catalogOptionalInt(value *int64) (*int, error) {
	if value == nil {
		return nil, nil
	}
	result, err := catalogInt(*value)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
func categoryView(item *gen.ProductCategory) (*gen.HqProductCategoryView, error) {
	sortOrder, err := catalogInt(item.SortOrder)
	if err != nil {
		return nil, err
	}
	return &gen.HqProductCategoryView{ID: item.ID, Name: item.Name, ParentID: item.ParentID, Enabled: item.Enabled, SortOrder: sortOrder}, nil
}
func productView(item *gen.Product) *gen.HqProductView {
	var brand *string
	if item.Brand != nil {
		brand = &item.Brand.Name
	}
	selectedValueIDs := make([]string, 0, len(item.SpecificationChoices))
	for _, choice := range item.SpecificationChoices {
		selectedValueIDs = append(selectedValueIDs, choice.ValueID)
	}
	return &gen.HqProductView{ID: item.ID, Name: item.Name, CategoryID: item.CategoryID, BrandID: item.BrandID, Brand: brand,
		Description: item.Description, ImageURL: item.ImageURL, DefaultPackageTemplateID: item.DefaultPackageTemplateID, SelectedValueIds: selectedValueIDs, Enabled: item.Enabled}
}
func skuView(item *gen.ProductSku) (*gen.HqProductSkuView, error) {
	life, err := catalogOptionalInt(item.ShelfLifeDays)
	if err != nil {
		return nil, err
	}
	version, err := catalogInt(item.PublishedPackageSetVersion)
	if err != nil {
		return nil, err
	}
	valueIDs := make([]string, 0, len(item.SpecificationValues))
	for _, link := range item.SpecificationValues {
		valueIDs = append(valueIDs, link.ValueID)
	}
	return &gen.HqProductSkuView{ID: item.ID, ProductID: item.ProductID, Name: item.Name, ValueIds: valueIDs,
		Ingredients: item.Ingredients, Allergens: item.Allergens, StorageInstructions: item.StorageInstructions,
		ShelfLifeDays: life, PublishedPackageSetVersion: version, Enabled: item.Enabled, SelectionRetired: item.SelectionRetired}, nil
}
func packageView(item *gen.ProductPackage) (*gen.HqProductPackageView, error) {
	version, err := catalogInt(item.PackageSetVersion)
	if err != nil {
		return nil, err
	}
	quantity, err := catalogOptionalInt(item.ContainsQuantity)
	if err != nil {
		return nil, err
	}
	price, err := catalogOptionalInt(item.SuggestedPriceFen)
	if err != nil {
		return nil, err
	}
	return &gen.HqProductPackageView{ID: item.ID, SkuID: item.SkuID, Name: item.Name, Barcode: item.Barcode,
		PackageSetVersion: version, ContainsPackageID: item.ContainsPackageID, ContainsQuantity: quantity,
		SuggestedPriceFen: price, Enabled: item.Enabled}, nil
}
