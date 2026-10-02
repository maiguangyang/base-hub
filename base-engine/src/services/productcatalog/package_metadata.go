package productcatalog

import (
	"context"
	"strings"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

// UpdatePackageMetadata changes presentation and suggested price only; conversion identity remains immutable.
func (s *Service) UpdatePackageMetadata(ctx context.Context, principal *auth.WorkspacePrincipal, id, name string, price *int64) (*gen.ProductPackage, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)
	if !validName(name, 64) || price != nil && *price < 0 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var item gen.ProductPackage
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var sku gen.ProductSku
		if err := lockedPackageForHQ(tx, hqID, id, &item, &sku); err != nil {
			return err
		}
		if !item.Enabled && item.PackageSetVersion <= sku.PublishedPackageSetVersion {
			return auth.NewError(auth.CodeConflict)
		}
		if err := tx.Model(&item).Updates(map[string]any{"name": name, "suggested_price_fen": price}).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productPackage", ResourceID: id, ResultCode: "SUCCESS"})
	})
	item.Name, item.SuggestedPriceFen = name, price
	return &item, err
}
