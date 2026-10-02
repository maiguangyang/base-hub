package productcatalog

import (
	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func (s *Service) saveProductVariants(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID string, product *gen.Product, input ProductInput) error {
	variants, selected, err := resolveVariants(tx, hqID, input)
	if err != nil {
		return err
	}
	overrides, err := validateVariantOverrides(variants, input.SkuOverrides)
	if err != nil {
		return err
	}
	templates, err := loadVariantTemplates(tx, hqID, input.DefaultPackageTemplateID, overrides)
	if err != nil {
		return err
	}
	existing, err := loadExistingVariants(tx, product.ID)
	if err != nil {
		return err
	}
	if err := replaceValidatedProductChoices(tx, product, selected, variants, existing); err != nil {
		return err
	}
	for _, variant := range variants {
		if err := s.upsertVariantSku(tx, principal, hqID, product.ID, variant, overrides, input.DefaultPackageTemplateID, templates, existing); err != nil {
			return err
		}
	}
	if err := retireRemovedVariants(tx, existing); err != nil {
		return err
	}
	if err := tx.Model(product).Update("default_package_template_id", input.DefaultPackageTemplateID).Error; err != nil {
		return err
	}
	product.DefaultPackageTemplateID = input.DefaultPackageTemplateID
	return nil
}

func replaceValidatedProductChoices(tx *gorm.DB, product *gen.Product, selected []string, variants []variantChoice, existing map[string]*gen.ProductSku) error {
	if err := validateRestrictedVariants(variants, existing); err != nil {
		return err
	}
	return replaceProductChoices(tx, product, selected)
}

func validateRestrictedVariants(variants []variantChoice, existing map[string]*gen.ProductSku) error {
	for _, variant := range variants {
		if variant.restricted && existing[variant.key] == nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
	}
	return nil
}

func (s *Service) upsertVariantSku(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID, productID string, variant variantChoice, overrides map[string]SkuOverride, defaultID *string, templates map[string][]gen.ProductPackageTemplate, existing map[string]*gen.ProductSku) error {
	override, hasOverride := overrides[variant.key]
	if sku := existing[variant.key]; sku != nil {
		enabled := sku.Enabled || sku.SelectionRetired
		if hasOverride {
			enabled = override.Enabled
		}
		if err := s.applyExistingPackageOverride(tx, principal, hqID, sku, override); err != nil {
			return err
		}
		if err := updateVariantSku(tx, sku, variant.name, enabled); err != nil {
			return err
		}
		delete(existing, variant.key)
		return nil
	}
	return s.createVariantSku(tx, principal, hqID, productID, variant, !hasOverride || override.Enabled, override, defaultID, templates)
}

func validateVariantOverrides(variants []variantChoice, input []SkuOverride) (map[string]SkuOverride, error) {
	allowed := make(map[string]struct{}, len(variants))
	for _, variant := range variants {
		allowed[variant.key] = struct{}{}
	}
	overrides := make(map[string]SkuOverride, len(input))
	for _, override := range input {
		key := variantKey(override.ValueIDs)
		if _, ok := allowed[key]; !ok {
			return nil, auth.NewError(auth.CodeValidationFailed)
		}
		if _, ok := overrides[key]; ok || override.DisableDefaultPackage && override.PackageTemplateID != nil {
			return nil, auth.NewError(auth.CodeValidationFailed)
		}
		overrides[key] = override
	}
	return overrides, nil
}

func loadVariantTemplates(tx *gorm.DB, hqID string, defaultID *string, overrides map[string]SkuOverride) (map[string][]gen.ProductPackageTemplate, error) {
	ids := make(map[string]struct{})
	if defaultID != nil {
		ids[*defaultID] = struct{}{}
	}
	for _, override := range overrides {
		if override.PackageTemplateID != nil {
			ids[*override.PackageTemplateID] = struct{}{}
		}
	}
	templates := make(map[string][]gen.ProductPackageTemplate, len(ids))
	for id := range ids {
		packages, err := collectPackageTemplates(tx, hqID, []string{id})
		if err != nil {
			return nil, err
		}
		templates[id] = packages
	}
	return templates, nil
}

func replaceProductChoices(tx *gorm.DB, product *gen.Product, selected []string) error {
	if err := tx.Where("product_id = ?", product.ID).Delete(&gen.ProductSpecificationChoice{}).Error; err != nil {
		return err
	}
	product.SpecificationChoices = make([]*gen.ProductSpecificationChoice, 0, len(selected))
	for _, valueID := range selected {
		choice := gen.ProductSpecificationChoice{ID: uuid.Must(uuid.NewV4()).String(), ProductID: product.ID, ValueID: valueID}
		if err := tx.Create(&choice).Error; err != nil {
			return err
		}
		product.SpecificationChoices = append(product.SpecificationChoices, &choice)
	}
	return nil
}

func loadExistingVariants(tx *gorm.DB, productID string) (map[string]*gen.ProductSku, error) {
	var existing []gen.ProductSku
	if err := tx.Where("product_id = ?", productID).Find(&existing).Error; err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(existing))
	for _, sku := range existing {
		ids = append(ids, sku.ID)
	}
	var links []gen.ProductSkuSpecificationValue
	if len(ids) > 0 {
		if err := tx.Where("sku_id IN ?", ids).Find(&links).Error; err != nil {
			return nil, err
		}
	}
	valuesBySku := make(map[string][]string, len(existing))
	for _, link := range links {
		valuesBySku[link.SkuID] = append(valuesBySku[link.SkuID], link.ValueID)
	}
	byKey := make(map[string]*gen.ProductSku, len(existing))
	for i := range existing {
		key := variantKey(valuesBySku[existing[i].ID])
		if byKey[key] != nil {
			return nil, auth.NewError(auth.CodeConflict)
		}
		byKey[key] = &existing[i]
	}
	return byKey, nil
}

func updateVariantSku(tx *gorm.DB, sku *gen.ProductSku, name string, enabled bool) error {
	wasEnabled := sku.Enabled
	if err := tx.Model(sku).Updates(map[string]any{"name": name, "enabled": enabled, "selection_retired": false}).Error; err != nil {
		return err
	}
	if !enabled && wasEnabled {
		return disableOffersForSkus(tx, tx.Model(&gen.ProductSku{}).Select("id").Where("id = ?", sku.ID))
	}
	return nil
}

func (s *Service) createVariantSku(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID, productID string, variant variantChoice, enabled bool, override SkuOverride, defaultID *string, templates map[string][]gen.ProductPackageTemplate) error {
	sku := gen.ProductSku{ID: uuid.Must(uuid.NewV4()).String(), ProductID: productID, Name: variant.name, Enabled: true}
	if err := tx.Create(&sku).Error; err != nil {
		return err
	}
	for _, value := range variant.values {
		link := gen.ProductSkuSpecificationValue{ID: uuid.Must(uuid.NewV4()).String(), SkuID: sku.ID, ValueID: value.ID}
		if err := tx.Create(&link).Error; err != nil {
			return err
		}
	}
	templateID := defaultID
	if override.DisableDefaultPackage {
		templateID = nil
	}
	if override.PackageTemplateID != nil {
		templateID = override.PackageTemplateID
	}
	if templateID != nil {
		if err := s.copyTemplatePackages(tx, principal, hqID, &sku, 1, templates[*templateID]); err != nil {
			return err
		}
	}
	if !enabled {
		return updateVariantSku(tx, &sku, variant.name, false)
	}
	return nil
}

func retireRemovedVariants(tx *gorm.DB, existing map[string]*gen.ProductSku) error {
	for _, sku := range existing {
		if !sku.Enabled {
			continue
		}
		if err := tx.Model(sku).Updates(map[string]any{"enabled": false, "selection_retired": true}).Error; err != nil {
			return err
		}
		if err := disableOffersForSkus(tx, tx.Model(&gen.ProductSku{}).Select("id").Where("id = ?", sku.ID)); err != nil {
			return err
		}
	}
	return nil
}
