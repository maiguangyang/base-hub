package productcatalog

import (
	"context"
	"strings"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) UpdateSku(ctx context.Context, principal *auth.WorkspacePrincipal, id string, input SkuInput) (*gen.ProductSku, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if !validName(name, 128) || (input.ShelfLifeDays != nil && *input.ShelfLifeDays <= 0) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var item gen.ProductSku
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := skuForHQ(tx, id, hqID, &item); err != nil {
			return err
		}
		if input.ProductID != item.ProductID {
			return auth.NewError(auth.CodeValidationFailed)
		}
		if err := tx.Model(&item).Updates(map[string]any{"name": name, "ingredients": input.Ingredients, "allergens": input.Allergens,
			"storage_instructions": input.StorageInstructions, "shelf_life_days": input.ShelfLifeDays}).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productSku", ResourceID: id, ResultCode: "SUCCESS"})
	})
	item.Name = name
	item.Ingredients, item.Allergens = input.Ingredients, input.Allergens
	item.StorageInstructions, item.ShelfLifeDays = input.StorageInstructions, input.ShelfLifeDays
	return &item, err
}

func (s *Service) SetSkuEnabled(ctx context.Context, principal *auth.WorkspacePrincipal, id string, enabled bool) (*gen.ProductSku, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	var item gen.ProductSku
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := skuForHQ(tx, id, hqID, &item); err != nil {
			return err
		}
		if enabled {
			if err := ensureSkuCanBeEnabled(tx, &item); err != nil {
				return err
			}
		}
		if err := tx.Model(&item).Updates(map[string]any{"enabled": enabled, "selection_retired": false}).Error; err != nil {
			return err
		}
		if !enabled {
			skuIDs := tx.Model(&gen.ProductSku{}).Select("id").Where("id = ?", id)
			if err := disableOffersForSkus(tx, skuIDs); err != nil {
				return err
			}
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productSku", ResourceID: id, ResultCode: "SUCCESS"})
	})
	item.Enabled = enabled
	return &item, err
}

func (s *Service) RetirePackage(ctx context.Context, principal *auth.WorkspacePrincipal, id string) (*gen.ProductPackage, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	var item gen.ProductPackage
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var sku gen.ProductSku
		if err := lockedPackageForHQ(tx, hqID, id, &item, &sku); err != nil {
			return err
		}
		if err := retirePackageRecord(tx, &item, &sku); err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productPackage", ResourceID: id, ResultCode: "SUCCESS"})
	})
	item.Enabled = false
	return &item, err
}

func lockedPackageForHQ(tx *gorm.DB, hqID, id string, item *gen.ProductPackage, sku *gen.ProductSku) error {
	if err := tx.Joins("JOIN product_skus ON product_skus.id = product_packages.sku_id").Joins("JOIN products ON products.id = product_skus.product_id").
		Where("product_packages.id = ? AND products.organization_id = ?", id, hqID).First(item).Error; err != nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	if err := skuForHQ(tx.Clauses(clause.Locking{Strength: "UPDATE"}), item.SkuID, hqID, sku); err != nil {
		return err
	}
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(item, "id = ?", id).Error
}

func retirePackageRecord(tx *gorm.DB, item *gen.ProductPackage, sku *gen.ProductSku) error {
	draft := !item.Enabled && item.PackageSetVersion > sku.PublishedPackageSetVersion
	parents := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("contains_package_id = ?", item.ID)
	if !draft {
		parents = parents.Where("enabled = ?", true)
	}
	var found []gen.ProductPackage
	if err := parents.Limit(1).Find(&found).Error; err != nil {
		return err
	}
	if len(found) > 0 {
		return auth.NewError(auth.CodeConflict)
	}
	if draft {
		return tx.Delete(item).Error
	}
	if err := tx.Model(item).Update("enabled", false).Error; err != nil {
		return err
	}
	return retirePackageOffers(tx, item.ID)
}

func retirePackageOffers(tx *gorm.DB, packageID string) error {
	offerIDs := tx.Model(&gen.StorePackageOffer{}).Select("id").Where("package_id = ?", packageID)
	promotionIDs := tx.Model(&gen.StorePromotionTarget{}).Select("promotion_id").Where("offer_id IN (?)", offerIDs)
	if err := tx.Model(&gen.StorePromotion{}).Where("id IN (?)", promotionIDs).Update("enabled", false).Error; err != nil {
		return err
	}
	return tx.Model(&gen.StorePackageOffer{}).Where("package_id = ?", packageID).Update("enabled", false).Error
}

func skuForHQ(tx *gorm.DB, id, hqID string, item *gen.ProductSku) error {
	if err := tx.Joins("JOIN products ON products.id = product_skus.product_id").Where("product_skus.id = ? AND products.organization_id = ?", id, hqID).First(item).Error; err != nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return nil
}
