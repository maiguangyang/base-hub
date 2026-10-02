package src

import (
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/productcatalog"
)

func productInput(input gen.HqCreateProductInput) productcatalog.ProductInput {
	result := productcatalog.ProductInput{Name: input.Name, CategoryID: input.CategoryID, BrandID: input.BrandID,
		Description: input.Description, DefaultPackageTemplateID: input.DefaultPackageTemplateID}
	for _, selection := range input.Selections {
		if selection == nil {
			continue
		}
		result.Selections = append(result.Selections, productcatalog.SpecificationSelection{
			SpecificationID: selection.SpecificationID, ValueIDs: selection.ValueIds})
	}
	for _, override := range input.SkuOverrides {
		if override == nil {
			continue
		}
		result.SkuOverrides = append(result.SkuOverrides, productcatalog.SkuOverride{
			ValueIDs: override.ValueIds, Enabled: override.Enabled,
			PackageTemplateID: override.PackageTemplateID, DisableDefaultPackage: override.DisableDefaultPackage})
	}
	return result
}

func validateProductInputPointers(input gen.HqCreateProductInput) error {
	for _, selection := range input.Selections {
		if selection == nil {
			return auth.NewError(auth.CodeValidationFailed)
		}
	}
	for _, override := range input.SkuOverrides {
		if override == nil {
			return auth.NewError(auth.CodeValidationFailed)
		}
	}
	return nil
}
