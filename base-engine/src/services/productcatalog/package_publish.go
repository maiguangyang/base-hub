package productcatalog

import (
	"context"
	"math"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) PublishPackageSet(ctx context.Context, principal *auth.WorkspacePrincipal, skuID string, version int64) (*gen.ProductPackage, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	if skuID == "" || version < 1 || version > math.MaxInt32 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var base gen.ProductPackage
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.publishPackageSetTx(tx, principal, hqID, skuID, version, &base)
	}, packageWriteOptions(s.db))
	base.Enabled = err == nil
	return &base, packageWriteError(err)
}

func (s *Service) publishPackageSetTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID, skuID string, version int64, base *gen.ProductPackage) error {
	var sku gen.ProductSku
	if err := skuForHQ(tx.Clauses(clause.Locking{Strength: "UPDATE"}), skuID, hqID, &sku); err != nil {
		return err
	}
	var packages []gen.ProductPackage
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("sku_id = ?", skuID).Find(&packages).Error; err != nil {
		return err
	}
	if version <= sku.PublishedPackageSetVersion {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if err := validatePublishPackageSet(packages, version, base); err != nil {
		return err
	}
	if err := ensurePublishedBarcodes(tx, packages, version); err != nil {
		return err
	}
	if err := retirePreviousPackageVersions(tx, packages, version); err != nil {
		return err
	}
	if err := activatePackageSet(tx, &sku, version); err != nil {
		return err
	}
	return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productPackage", ResourceID: base.ID, ResultCode: "SUCCESS"})
}

func retirePreviousPackageVersions(tx *gorm.DB, packages []gen.ProductPackage, version int64) error {
	for _, pack := range packages {
		if !pack.Enabled || pack.PackageSetVersion == version {
			continue
		}
		if err := tx.Model(&pack).Update("enabled", false).Error; err != nil {
			return err
		}
		if err := retirePackageOffers(tx, pack.ID); err != nil {
			return err
		}
	}
	return nil
}

func activatePackageSet(tx *gorm.DB, sku *gen.ProductSku, version int64) error {
	if err := tx.Model(&gen.ProductPackage{}).Where("sku_id = ? AND package_set_version = ?", sku.ID, version).Update("enabled", true).Error; err != nil {
		return err
	}
	return tx.Model(sku).Update("published_package_set_version", version).Error
}

func validatePublishPackageSet(packages []gen.ProductPackage, version int64, base *gen.ProductPackage) error {
	byID, target, maxVersion, activeVersion, err := indexPackageSet(packages, version, base)
	if err != nil {
		return err
	}
	if version != maxVersion || version <= activeVersion || base.ID == "" {
		return auth.NewError(auth.CodeValidationFailed)
	}
	for _, pack := range target {
		if err := validatePackageChainToBase(pack, byID, version, base.ID); err != nil {
			return err
		}
	}
	return nil
}

func indexPackageSet(packages []gen.ProductPackage, version int64, base *gen.ProductPackage) (map[string]gen.ProductPackage, []gen.ProductPackage, int64, int64, error) {
	byID := make(map[string]gen.ProductPackage, len(packages))
	barcodes := make(map[string]struct{}, len(packages))
	var target []gen.ProductPackage
	var maxVersion, activeVersion int64
	for _, pack := range packages {
		byID[pack.ID] = pack
		if pack.PackageSetVersion > maxVersion {
			maxVersion = pack.PackageSetVersion
		}
		if pack.Enabled {
			if activeVersion != 0 && activeVersion != pack.PackageSetVersion {
				return nil, nil, 0, 0, auth.NewError(auth.CodeConflict)
			}
			activeVersion = pack.PackageSetVersion
		}
		if pack.PackageSetVersion != version {
			continue
		}
		if err := recordTargetPackage(pack, base, barcodes); err != nil {
			return nil, nil, 0, 0, err
		}
		target = append(target, pack)
	}
	return byID, target, maxVersion, activeVersion, nil
}

func recordTargetPackage(pack gen.ProductPackage, base *gen.ProductPackage, barcodes map[string]struct{}) error {
	if pack.Barcode != nil {
		if _, exists := barcodes[*pack.Barcode]; exists {
			return auth.NewError(auth.CodeValidationFailed)
		}
		barcodes[*pack.Barcode] = struct{}{}
	}
	if pack.ContainsPackageID == nil {
		if base.ID != "" {
			return auth.NewError(auth.CodeValidationFailed)
		}
		*base = pack
	}
	return nil
}

func validatePackageChainToBase(pack gen.ProductPackage, byID map[string]gen.ProductPackage, version int64, baseID string) error {
	seen := map[string]struct{}{}
	for pack.ContainsPackageID != nil {
		if _, exists := seen[pack.ID]; exists {
			return auth.NewError(auth.CodeValidationFailed)
		}
		seen[pack.ID] = struct{}{}
		child, exists := byID[*pack.ContainsPackageID]
		if !exists || child.PackageSetVersion != version || !validPackageMultiplier(1, pack.ContainsQuantity) {
			return auth.NewError(auth.CodeValidationFailed)
		}
		pack = child
	}
	if pack.ID != baseID {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}
