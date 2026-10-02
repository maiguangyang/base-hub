package productcatalog

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

func (s *Service) DeleteProduct(ctx context.Context, principal *auth.WorkspacePrincipal, id string) error {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessDelete)
	if err != nil {
		return err
	}
	var imageURL *string
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockCategoryHierarchy(tx, hqID); err != nil {
			return err
		}
		var item gen.Product
		if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err := deleteUnreferencedProductChildren(tx, id); err != nil {
			return err
		}
		imageURL = item.ImageURL
		if err := tx.Unscoped().Delete(&item).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "product", ResourceID: id, ResultCode: "SUCCESS"})
	})
	if err == nil && imageURL != nil {
		s.cleanupManagedImage(*imageURL, id)
	}
	return err
}

func deleteUnreferencedProductChildren(tx *gorm.DB, productID string) error {
	skuIDs := tx.Model(&gen.ProductSku{}).Select("id").Where("product_id = ?", productID)
	var count int64
	if err := tx.Model(&gen.StoreListing{}).Where("sku_id IN (?)", skuIDs).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if err := tx.Unscoped().Where("sku_id IN (?)", skuIDs).Delete(&gen.ProductSkuSpecificationValue{}).Error; err != nil {
		return err
	}
	if err := tx.Unscoped().Where("product_id = ?", productID).Delete(&gen.ProductSpecificationChoice{}).Error; err != nil {
		return err
	}
	if err := tx.Unscoped().Where("sku_id IN (?)", skuIDs).Delete(&gen.ProductPackage{}).Error; err != nil {
		return err
	}
	return tx.Unscoped().Where("product_id = ?", productID).Delete(&gen.ProductSku{}).Error
}
