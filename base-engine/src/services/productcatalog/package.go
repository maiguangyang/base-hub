package productcatalog

import (
	"context"
	"math"
	"strings"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PackageInput struct {
	SkuID, Name, Barcode string
	ContainsPackageID    *string
	ContainsQuantity     int64
	SuggestedPriceFen    *int64
	PackageSetVersion    int64
}

func (s *Service) CreatePackage(ctx context.Context, principal *auth.WorkspacePrincipal, input PackageInput) (*gen.ProductPackage, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	name, barcode := strings.TrimSpace(input.Name), strings.TrimSpace(input.Barcode)
	if err := validatePackageInput(input, name, barcode); err != nil {
		return nil, err
	}
	if input.PackageSetVersion == 0 {
		input.PackageSetVersion = 1
	}
	item := &gen.ProductPackage{ID: uuid.Must(uuid.NewV4()).String(), Name: name, SkuID: input.SkuID, PackageSetVersion: input.PackageSetVersion,
		ContainsPackageID: input.ContainsPackageID, SuggestedPriceFen: input.SuggestedPriceFen}
	if barcode != "" {
		item.Barcode = &barcode
	}
	input.Barcode = barcode
	if input.ContainsPackageID != nil {
		item.ContainsQuantity = &input.ContainsQuantity
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.createPackageTx(tx, principal, hqID, input, item)
	}, packageWriteOptions(s.db))
	return item, packageWriteError(err)
}

func (s *Service) createPackageTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID string, input PackageInput, item *gen.ProductPackage) error {
	var sku gen.ProductSku
	if err := skuForHQ(tx.Clauses(clause.Locking{Strength: "UPDATE"}), input.SkuID, hqID, &sku); err != nil {
		return err
	}
	enabled, err := packageCreationState(tx, &sku, input.PackageSetVersion)
	if err != nil {
		return err
	}
	item.Enabled = enabled
	if err := validatePackageCreate(tx, hqID, input, enabled); err != nil {
		return err
	}
	if err := tx.Select("*").Create(item).Error; err != nil {
		return err
	}
	if enabled && sku.PublishedPackageSetVersion == 0 {
		if err := tx.Model(&sku).Update("published_package_set_version", input.PackageSetVersion).Error; err != nil {
			return err
		}
	}
	return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productPackage", ResourceID: item.ID, ResultCode: "SUCCESS"})
}

func packageCreationState(tx *gorm.DB, sku *gen.ProductSku, version int64) (bool, error) {
	var packages []gen.ProductPackage
	if err := tx.Select("package_set_version", "enabled").Where("sku_id = ?", sku.ID).Find(&packages).Error; err != nil {
		return false, err
	}
	maxVersion, activeVersion, err := summarizePackageVersions(packages)
	if err != nil {
		return false, err
	}
	return validateNewPackageVersion(version, maxVersion, activeVersion, sku.PublishedPackageSetVersion)
}

func summarizePackageVersions(packages []gen.ProductPackage) (int64, int64, error) {
	var maxVersion, activeVersion int64
	for _, pack := range packages {
		if pack.PackageSetVersion > maxVersion {
			maxVersion = pack.PackageSetVersion
		}
		if pack.Enabled {
			if activeVersion != 0 && activeVersion != pack.PackageSetVersion {
				return 0, 0, auth.NewError(auth.CodeConflict)
			}
			activeVersion = pack.PackageSetVersion
		}
	}
	return maxVersion, activeVersion, nil
}

func validateNewPackageVersion(version, maxVersion, activeVersion, publishedVersion int64) (bool, error) {
	if maxVersion == 0 {
		return initialPackageVersionState(version, publishedVersion)
	}
	if invalidPackageVersionProgression(version, maxVersion, activeVersion, publishedVersion) {
		return false, auth.NewError(auth.CodeValidationFailed)
	}
	if version <= publishedVersion && version != activeVersion {
		return false, auth.NewError(auth.CodeValidationFailed)
	}
	return version == activeVersion, nil
}

func initialPackageVersionState(version, publishedVersion int64) (bool, error) {
	if publishedVersion == 0 && version == 1 {
		return true, nil
	}
	if publishedVersion > 0 && version == publishedVersion+1 {
		return false, nil
	}
	return false, auth.NewError(auth.CodeValidationFailed)
}

func invalidPackageVersionProgression(version, maxVersion, activeVersion, publishedVersion int64) bool {
	if version < maxVersion || version > maxVersion+1 {
		return true
	}
	return version > maxVersion && (maxVersion > publishedVersion || activeVersion != 0 && activeVersion != maxVersion)
}

func validatePackageInput(input PackageInput, name, barcode string) error {
	if !validPackageBasics(input, name, barcode) || !validPackageContainment(input) {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}

func validPackageBasics(input PackageInput, name, barcode string) bool {
	if !validName(name, 64) || len(barcode) > 64 || input.SkuID == "" {
		return false
	}
	if input.SuggestedPriceFen != nil && *input.SuggestedPriceFen < 0 {
		return false
	}
	return input.PackageSetVersion >= 0 && input.PackageSetVersion <= math.MaxInt32
}

func validPackageContainment(input PackageInput) bool {
	if input.ContainsPackageID == nil {
		return input.ContainsQuantity == 0
	}
	return input.ContainsQuantity > 1
}

func validatePackageCreate(tx *gorm.DB, hqID string, input PackageInput, enabled bool) error {
	var sku gen.ProductSku
	if err := tx.Joins("JOIN products ON products.id = product_skus.product_id").Where("product_skus.id = ? AND products.organization_id = ? AND product_skus.enabled = ?", input.SkuID, hqID, true).First(&sku).Error; err != nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	if err := ensurePackageBarcode(tx, input, enabled, ""); err != nil {
		return err
	}
	if input.ContainsPackageID == nil {
		return ensureUniqueBasePackage(tx, input)
	}
	return ensurePackageChain(tx, input, enabled)
}

func ensureUniqueBasePackage(tx *gorm.DB, input PackageInput) error {
	var count int64
	err := tx.Model(&gen.ProductPackage{}).Where("sku_id = ? AND package_set_version = ? AND contains_package_id IS NULL", input.SkuID, input.PackageSetVersion).Count(&count).Error
	if err != nil {
		return err
	}
	if count != 0 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}

func ensurePackageChain(tx *gorm.DB, input PackageInput, enabled bool) error {
	var contained gen.ProductPackage
	if err := tx.Where("id = ? AND sku_id = ?", *input.ContainsPackageID, input.SkuID).First(&contained).Error; err != nil {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if contained.PackageSetVersion != input.PackageSetVersion || contained.Enabled != enabled {
		return auth.NewError(auth.CodeValidationFailed)
	}
	factor := input.ContainsQuantity
	for contained.ContainsPackageID != nil {
		if !validPackageMultiplier(factor, contained.ContainsQuantity) {
			return auth.NewError(auth.CodeValidationFailed)
		}
		factor *= *contained.ContainsQuantity
		var next gen.ProductPackage
		if err := tx.Where("id = ? AND sku_id = ?", *contained.ContainsPackageID, input.SkuID).First(&next).Error; err != nil {
			return auth.NewError(auth.CodeValidationFailed)
		}
		if next.PackageSetVersion != input.PackageSetVersion || next.Enabled != enabled {
			return auth.NewError(auth.CodeValidationFailed)
		}
		contained = next
	}
	return nil
}

func validPackageMultiplier(factor int64, quantity *int64) bool {
	return quantity != nil && *quantity > 1 && factor <= math.MaxInt64 / *quantity
}
