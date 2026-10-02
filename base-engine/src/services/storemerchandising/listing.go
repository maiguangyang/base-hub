package storemerchandising

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

func (s *Service) SetListing(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, skuID string, enabled bool) (*gen.StoreListing, error) {
	if storeID == "" || skuID == "" {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var listing gen.StoreListing
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		store, err := storeScope(tx, principal, storeID, "franchiseProduct:manage", authorization.AccessUpdate)
		if err != nil {
			return err
		}
		if enabled {
			if err := publishedSku(tx, skuID); err != nil {
				return err
			}
		}
		if err := upsertListing(tx, storeID, skuID, enabled, &listing); err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: "franchiseProduct:manage",
			ResourceType: "storeListing", ResourceID: listing.ID, ResultCode: "SUCCESS"})
	})
	return &listing, err
}

func upsertListing(tx *gorm.DB, storeID, skuID string, enabled bool, listing *gen.StoreListing) error {
	query := tx.Where("store_id = ? AND sku_id = ?", storeID, skuID).First(listing)
	if query.Error != nil && query.Error != gorm.ErrRecordNotFound {
		return query.Error
	}
	if query.Error == gorm.ErrRecordNotFound {
		*listing = gen.StoreListing{ID: uuid.Must(uuid.NewV4()).String(), StoreID: storeID, SkuID: skuID, Enabled: enabled, SelectedAt: time.Now()}
		return tx.Create(listing).Error
	}
	if listing.Enabled == enabled {
		return nil
	}
	if err := tx.Model(listing).Update("enabled", enabled).Error; err != nil {
		return err
	}
	listing.Enabled = enabled
	return nil
}

func publishedSku(tx *gorm.DB, skuID string) error {
	var sku gen.ProductSku
	if err := tx.Joins("JOIN products ON products.id = product_skus.product_id").Joins("JOIN organizations ON organizations.id = products.organization_id").Where("product_skus.id = ? AND product_skus.enabled = ? AND products.enabled = ? AND organizations.type = ? AND organizations.status = ?", skuID, true, true, gen.OrganizationTypeHeadquarters, gen.OrganizationStatusActive).First(&sku).Error; err != nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return nil
}
