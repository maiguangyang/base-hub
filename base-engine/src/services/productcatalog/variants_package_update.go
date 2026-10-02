package productcatalog

import (
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) applyExistingPackageOverride(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID string, sku *gen.ProductSku, override SkuOverride) error {
	if override.PackageTemplateID == nil && !override.DisableDefaultPackage {
		return nil
	}
	var packages []gen.ProductPackage
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("sku_id = ?", sku.ID).Find(&packages).Error; err != nil {
		return err
	}
	maxVersion, _, err := summarizePackageVersions(packages)
	if err != nil {
		return err
	}
	if maxVersion > sku.PublishedPackageSetVersion {
		return auth.NewError(auth.CodeConflict)
	}
	if override.DisableDefaultPackage {
		return retirePreviousPackageVersions(tx, packages, 0)
	}
	templates, err := collectPackageTemplates(tx, hqID, []string{*override.PackageTemplateID})
	if err != nil {
		return err
	}
	if activePackagesMatchTemplate(packages, sku.PublishedPackageSetVersion, templates) {
		return nil
	}
	return s.switchExistingPackageTemplate(tx, principal, hqID, sku, templates)
}

func (s *Service) switchExistingPackageTemplate(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID string, sku *gen.ProductSku, templates []gen.ProductPackageTemplate) error {
	if !sku.Enabled {
		if err := tx.Model(sku).Update("enabled", true).Error; err != nil {
			return err
		}
		sku.Enabled = true
	}
	version := sku.PublishedPackageSetVersion + 1
	if err := s.copyTemplatePackages(tx, principal, hqID, sku, version, templates); err != nil {
		return err
	}
	if sku.PublishedPackageSetVersion == 0 {
		sku.PublishedPackageSetVersion = version
		return nil
	}
	var base gen.ProductPackage
	if err := s.publishPackageSetTx(tx, principal, hqID, sku.ID, version, &base); err != nil {
		return err
	}
	sku.PublishedPackageSetVersion = version
	return nil
}
