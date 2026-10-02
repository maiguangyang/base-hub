package productcatalog

import (
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func ensureSkuCanBeEnabled(tx *gorm.DB, sku *gen.ProductSku) error {
	if sku.SelectionRetired {
		return auth.NewError(auth.CodeValidationFailed)
	}
	active, err := skuInCurrentSelection(tx, sku)
	if err != nil {
		return err
	}
	if !active {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}

func skuInCurrentSelection(tx *gorm.DB, sku *gen.ProductSku) (bool, error) {
	var values []gen.SpecificationValue
	err := tx.Model(&gen.SpecificationValue{}).
		Joins("JOIN product_specification_choices ON product_specification_choices.value_id = specification_values.id").
		Where("product_specification_choices.product_id = ?", sku.ProductID).Find(&values).Error
	if err != nil {
		return false, err
	}
	var links []gen.ProductSkuSpecificationValue
	if err := tx.Where("sku_id = ?", sku.ID).Find(&links).Error; err != nil {
		return false, err
	}
	selected := make(map[string]string, len(values))
	specifications := make(map[string]struct{}, len(values))
	for _, value := range values {
		selected[value.ID] = value.SpecificationID
		specifications[value.SpecificationID] = struct{}{}
	}
	if len(links) != len(specifications) {
		return false, nil
	}
	seen := make(map[string]struct{}, len(links))
	for _, link := range links {
		specificationID, ok := selected[link.ValueID]
		if !ok {
			return false, nil
		}
		if _, duplicate := seen[specificationID]; duplicate {
			return false, nil
		}
		seen[specificationID] = struct{}{}
	}
	return true, nil
}
