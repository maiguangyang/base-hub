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

func (s *Service) UpdateCategory(ctx context.Context, principal *auth.WorkspacePrincipal, id, name string, parentID *string) (*gen.ProductCategory, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)
	if !validName(name, 128) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	if id == "" || (parentID != nil && *parentID == id) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var item gen.ProductCategory
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockedCategoryForHQ(tx, hqID, id, &item); err != nil {
			return err
		}
		if err := ensureCategoryParent(tx, hqID, id, parentID); err != nil {
			return err
		}
		if err := ensureCategoryNameAvailableExcept(tx, hqID, id, name, parentID); err != nil {
			return err
		}
		if err := tx.Model(&item).Updates(map[string]any{"name": name, "parent_id": parentID}).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productCategory", ResourceID: id, ResultCode: "SUCCESS"})
	})
	item.Name, item.ParentID = name, parentID
	return &item, err
}

func (s *Service) SetCategoryEnabled(ctx context.Context, principal *auth.WorkspacePrincipal, id string, enabled bool) (*gen.ProductCategory, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	var item gen.ProductCategory
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err := tx.Model(&item).Update("enabled", enabled).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productCategory", ResourceID: id, ResultCode: "SUCCESS"})
	})
	item.Enabled = enabled
	return &item, err
}

func (s *Service) UpdateProduct(ctx context.Context, principal *auth.WorkspacePrincipal, id string, input ProductInput) (*gen.Product, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if !validName(name, 128) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var item gen.Product
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockCategoryHierarchy(tx, hqID); err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err := validateSelectedCategoryForUpdate(tx, hqID, item.CategoryID, input.CategoryID); err != nil {
			return err
		}
		brand, err := selectedBrandForUpdate(tx, hqID, item.BrandID, input.BrandID)
		if err != nil {
			return err
		}
		if err := s.saveProductVariants(tx, principal, hqID, &item, input); err != nil {
			return err
		}
		if err := tx.Model(&item).Updates(map[string]any{"name": name, "category_id": input.CategoryID, "brand_id": input.BrandID, "description": input.Description}).Error; err != nil {
			return err
		}
		item.Brand = brand
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "product", ResourceID: id, ResultCode: "SUCCESS"})
	})
	item.Name, item.CategoryID = name, input.CategoryID
	item.BrandID, item.Description = input.BrandID, input.Description
	item.DefaultPackageTemplateID = input.DefaultPackageTemplateID
	return &item, err
}

func validateSelectedCategoryForUpdate(tx *gorm.DB, hqID, currentID, requestedID string) error {
	var category gen.ProductCategory
	if err := tx.Where("id = ? AND organization_id = ?", requestedID, hqID).First(&category).Error; err != nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	if !category.Enabled && requestedID != currentID {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return nil
}

func (s *Service) SetProductEnabled(ctx context.Context, principal *auth.WorkspacePrincipal, id string, enabled bool) (*gen.Product, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	var item gen.Product
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err := tx.Model(&item).Update("enabled", enabled).Error; err != nil {
			return err
		}
		if !enabled {
			skuIDs := tx.Model(&gen.ProductSku{}).Select("id").Where("product_id = ?", id)
			if err := disableOffersForSkus(tx, skuIDs); err != nil {
				return err
			}
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "product", ResourceID: id, ResultCode: "SUCCESS"})
	})
	item.Enabled = enabled
	return &item, err
}
