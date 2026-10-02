package storemerchandising

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) SetPrice(ctx context.Context, principal *auth.WorkspacePrincipal, listingID, packageID string, priceFen int64, reasonCode string) (*gen.StorePackageOffer, error) {
	reasonCode = strings.TrimSpace(reasonCode)
	if !validPriceInput(listingID, packageID, priceFen, reasonCode) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var offer gen.StorePackageOffer
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.setPriceTx(tx, principal, listingID, packageID, priceFen, reasonCode, &offer)
	})
	return &offer, err
}

func validPriceInput(listingID, packageID string, priceFen int64, reasonCode string) bool {
	return listingID != "" && packageID != "" && priceFen >= 0 && priceFen <= math.MaxInt32 && len(reasonCode) <= 64
}

func (s *Service) setPriceTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, listingID, packageID string, priceFen int64, reasonCode string, offer *gen.StorePackageOffer) error {
	var listing gen.StoreListing
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&listing, "id = ?", listingID).Error; err != nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	store, err := storeScope(tx, principal, listing.StoreID, "franchiseProduct:manage", authorization.AccessUpdate)
	if err != nil {
		return err
	}
	if !listing.Enabled {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if err := publishedSku(tx, listing.SkuID); err != nil {
		return err
	}
	var pack gen.ProductPackage
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND sku_id = ? AND enabled = ?", packageID, listing.SkuID, true).First(&pack).Error; err != nil {
		return auth.NewError(auth.CodeValidationFailed)
	}
	previous, changed, err := upsertOfferPrice(tx, listingID, packageID, priceFen, offer)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	if err := appendPriceRevision(tx, offer.ID, previous, priceFen, reasonCode); err != nil {
		return err
	}
	return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: "franchiseProduct:manage",
		ResourceType: "storePackageOffer", ResourceID: offer.ID, ResultCode: "SUCCESS"})
}

func upsertOfferPrice(tx *gorm.DB, listingID, packageID string, priceFen int64, offer *gen.StorePackageOffer) (*int64, bool, error) {
	found := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("listing_id = ? AND package_id = ?", listingID, packageID).First(offer)
	if found.Error != nil && found.Error != gorm.ErrRecordNotFound {
		return nil, false, found.Error
	}
	if found.Error == gorm.ErrRecordNotFound {
		*offer = gen.StorePackageOffer{ID: uuid.Must(uuid.NewV4()).String(), ListingID: listingID, PackageID: packageID, PriceFen: priceFen, Enabled: true}
		return nil, true, tx.Create(offer).Error
	}
	if offer.PriceFen == priceFen && offer.Enabled {
		return nil, false, nil
	}
	previous := offer.PriceFen
	if err := tx.Model(offer).Updates(map[string]any{"price_fen": priceFen, "enabled": true}).Error; err != nil {
		return nil, false, err
	}
	offer.PriceFen, offer.Enabled = priceFen, true
	return &previous, true, nil
}

func appendPriceRevision(tx *gorm.DB, offerID string, previous *int64, priceFen int64, reasonCode string) error {
	revision := gen.StorePriceRevision{ID: uuid.Must(uuid.NewV4()).String(), OfferID: offerID,
		PreviousPriceFen: previous, PriceFen: priceFen, EffectiveAt: time.Now(), ReasonCode: reasonCode}
	return tx.Create(&revision).Error
}
