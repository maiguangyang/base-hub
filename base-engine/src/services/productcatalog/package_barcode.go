package productcatalog

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	mysqldriver "github.com/go-sql-driver/mysql"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func (s *Service) SetPackageBarcode(ctx context.Context, principal *auth.WorkspacePrincipal, id, barcode string) (*gen.ProductPackage, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	barcode = strings.TrimSpace(barcode)
	if len(barcode) > 64 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var item gen.ProductPackage
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.setPackageBarcodeTx(tx, principal, hqID, id, barcode, &item)
	}, packageWriteOptions(s.db))
	if err != nil {
		return nil, packageWriteError(err)
	}
	return &item, nil
}

func (s *Service) setPackageBarcodeTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID, id, barcode string, item *gen.ProductPackage) error {
	var sku gen.ProductSku
	if err := lockedPackageForHQ(tx, hqID, id, item, &sku); err != nil {
		return err
	}
	current := item.Enabled && item.PackageSetVersion == sku.PublishedPackageSetVersion
	draft := !item.Enabled && item.PackageSetVersion > sku.PublishedPackageSetVersion
	if !current && !draft {
		return auth.NewError(auth.CodeConflict)
	}
	if packageBarcodeEquals(item.Barcode, barcode) {
		return nil
	}
	input := PackageInput{SkuID: item.SkuID, PackageSetVersion: item.PackageSetVersion, Barcode: barcode}
	if err := ensurePackageBarcode(tx, input, item.Enabled, item.ID); err != nil {
		return err
	}
	var value *string
	if barcode != "" {
		value = &barcode
	}
	if err := tx.Model(item).Update("barcode", value).Error; err != nil {
		return err
	}
	item.Barcode = value
	return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productPackage", ResourceID: id, ResultCode: "SUCCESS"})
}

func ensurePackageBarcode(tx *gorm.DB, input PackageInput, enabled bool, excludeID string) error {
	if input.Barcode == "" {
		return nil
	}
	var matches []gen.ProductPackage
	if err := tx.Select("sku_id", "package_set_version", "enabled").Where("barcode = ? AND id <> ?", strings.TrimSpace(input.Barcode), excludeID).Find(&matches).Error; err != nil {
		return err
	}
	for _, match := range matches {
		sameVersion := match.SkuID == input.SkuID && match.PackageSetVersion == input.PackageSetVersion
		activeConflict := match.Enabled && (enabled || match.SkuID != input.SkuID)
		if sameVersion || activeConflict {
			return auth.NewError(auth.CodeConflict)
		}
	}
	return nil
}

func ensurePublishedBarcodes(tx *gorm.DB, packages []gen.ProductPackage, version int64) error {
	for _, pack := range packages {
		if pack.PackageSetVersion != version || pack.Barcode == nil {
			continue
		}
		input := PackageInput{SkuID: pack.SkuID, PackageSetVersion: version, Barcode: *pack.Barcode}
		if err := ensurePackageBarcode(tx, input, false, pack.ID); err != nil {
			return err
		}
	}
	return nil
}

func packageWriteError(err error) error {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return auth.NewError(auth.CodeConflict)
	}
	var mysqlErr *mysqldriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return auth.NewError(auth.CodeConflict)
	}
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) && (sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique || sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey) {
		return auth.NewError(auth.CodeConflict)
	}
	return err
}

// READ COMMITTED keeps barcode checks fresh after waiting for the SKU lock.
// In REPEATABLE READ, the initial ownership lookup can establish an older snapshot.
func packageWriteOptions(db *gorm.DB) *sql.TxOptions {
	if db.Dialector.Name() == "mysql" {
		return &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	}
	return nil
}

func packageBarcodeEquals(current *string, barcode string) bool {
	if current == nil {
		return barcode == ""
	}
	return *current == barcode
}
