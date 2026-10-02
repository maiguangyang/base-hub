package productcatalog

import (
	"context"
	"strings"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

func (s *Service) Brands(ctx context.Context, principal *auth.WorkspacePrincipal, filter CatalogFilter, page, perPage int) (*CatalogPage[*gen.ProductBrand], error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:read", authorization.AccessRead)
	if err != nil {
		return nil, err
	}
	offset, limit, err := pageBounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := applyCatalogFilter(s.db.WithContext(ctx).Model(&gen.ProductBrand{}).Where("organization_id = ?", hqID), "product_brands", filter)
	result := &CatalogPage[*gen.ProductBrand]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Order("name, id").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) CreateBrand(ctx context.Context, principal *auth.WorkspacePrincipal, rawName string) (*gen.ProductBrand, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(rawName)
	if !validName(name, 128) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	item := &gen.ProductBrand{ID: uuid.Must(uuid.NewV4()).String(), Name: name, OrganizationID: hqID, Enabled: true}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockCategoryHierarchy(tx, hqID); err != nil {
			return err
		}
		if err := availableBrandName(tx, hqID, name, ""); err != nil {
			return err
		}
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productBrand", ResourceID: item.ID, ResultCode: "SUCCESS"})
	})
	return item, err
}

func (s *Service) UpdateBrand(ctx context.Context, principal *auth.WorkspacePrincipal, id, rawName string) (*gen.ProductBrand, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(rawName)
	if !validName(name, 128) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var item gen.ProductBrand
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockCategoryHierarchy(tx, hqID); err != nil {
			return err
		}
		if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err := availableBrandName(tx, hqID, name, id); err != nil {
			return err
		}
		if err := tx.Model(&item).Update("name", name).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productBrand", ResourceID: id, ResultCode: "SUCCESS"})
	})
	item.Name = name
	return &item, err
}

func (s *Service) SetBrandEnabled(ctx context.Context, principal *auth.WorkspacePrincipal, id string, enabled bool) (*gen.ProductBrand, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	var item gen.ProductBrand
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err := tx.Model(&item).Update("enabled", enabled).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productBrand", ResourceID: id, ResultCode: "SUCCESS"})
	})
	item.Enabled = enabled
	return &item, err
}

func (s *Service) DeleteBrand(ctx context.Context, principal *auth.WorkspacePrincipal, id string) error {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessDelete)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockCategoryHierarchy(tx, hqID); err != nil {
			return err
		}
		var item gen.ProductBrand
		if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		var count int64
		if err := tx.Model(&gen.Product{}).Where("brand_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return auth.NewError(auth.CodeValidationFailed)
		}
		if err := tx.Unscoped().Delete(&item).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productBrand", ResourceID: id, ResultCode: "SUCCESS"})
	})
}

func availableBrandName(tx *gorm.DB, hqID, name, exceptID string) error {
	query := tx.Model(&gen.ProductBrand{}).Where("organization_id = ? AND name = ?", hqID, name)
	if exceptID != "" {
		query = query.Where("id <> ?", exceptID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}

func selectedBrand(tx *gorm.DB, hqID string, brandID *string) (*gen.ProductBrand, error) {
	if brandID == nil {
		return nil, nil
	}
	var brand gen.ProductBrand
	if err := tx.Where("id = ? AND organization_id = ? AND enabled = ?", *brandID, hqID, true).First(&brand).Error; err != nil {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	return &brand, nil
}

func selectedBrandForUpdate(tx *gorm.DB, hqID string, currentID, nextID *string) (*gen.ProductBrand, error) {
	if nextID == nil {
		return nil, nil
	}
	if currentID == nil || *currentID != *nextID {
		return selectedBrand(tx, hqID, nextID)
	}
	var brand gen.ProductBrand
	if err := tx.Where("id = ? AND organization_id = ?", *nextID, hqID).First(&brand).Error; err != nil {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	return &brand, nil
}
