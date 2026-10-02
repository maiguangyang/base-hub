package productcatalog

import (
	"context"
	"strings"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type PackageTemplateInput struct {
	Name              string
	ContainsPackageID *string
	ContainsQuantity  int64
}

func (s *Service) PackageTemplates(ctx context.Context, principal *auth.WorkspacePrincipal, filter CatalogFilter, page, perPage int) (*CatalogPage[*gen.ProductPackageTemplate], error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:read", authorization.AccessRead)
	if err != nil {
		return nil, err
	}
	offset, limit, err := pageBounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := applyCatalogFilter(s.db.WithContext(ctx).Model(&gen.ProductPackageTemplate{}).Where("organization_id = ?", hqID), "product_package_templates", filter)
	if filter.ContainsPackageID != nil {
		query = query.Where("product_package_templates.contains_package_id = ?", *filter.ContainsPackageID)
	}
	result := &CatalogPage[*gen.ProductPackageTemplate]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := query.Order("name, id").Offset(offset).Limit(limit).Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func validatePackageTemplateInput(input PackageTemplateInput) (string, error) {
	name := strings.TrimSpace(input.Name)
	if !validName(name, 64) || input.ContainsPackageID == nil && input.ContainsQuantity != 0 ||
		input.ContainsPackageID != nil && (*input.ContainsPackageID == "" || input.ContainsQuantity < 2) {
		return "", auth.NewError(auth.CodeValidationFailed)
	}
	return name, nil
}

func (s *Service) CreatePackageTemplate(ctx context.Context, principal *auth.WorkspacePrincipal, input PackageTemplateInput) (*gen.ProductPackageTemplate, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	name, err := validatePackageTemplateInput(input)
	if err != nil {
		return nil, err
	}
	item := &gen.ProductPackageTemplate{ID: uuid.Must(uuid.NewV4()).String(), Name: name, OrganizationID: hqID, Enabled: true,
		ContainsPackageID: input.ContainsPackageID}
	if input.ContainsPackageID != nil {
		item.ContainsQuantity = &input.ContainsQuantity
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockCategoryHierarchy(tx, hqID); err != nil {
			return err
		}
		if err := validateTemplateChild(tx, hqID, item.ID, input.ContainsPackageID); err != nil {
			return err
		}
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		return s.templateAudit(tx, principal, hqID, "productPackageTemplate", item.ID)
	})
	return item, err
}

func (s *Service) UpdatePackageTemplate(ctx context.Context, principal *auth.WorkspacePrincipal, id string, input PackageTemplateInput) (*gen.ProductPackageTemplate, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	name, err := validatePackageTemplateInput(input)
	if err != nil {
		return nil, err
	}
	var item gen.ProductPackageTemplate
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockCategoryHierarchy(tx, hqID); err != nil {
			return err
		}
		if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err := validateTemplateChild(tx, hqID, id, input.ContainsPackageID); err != nil {
			return err
		}
		var quantity *int64
		if input.ContainsPackageID != nil {
			quantity = &input.ContainsQuantity
		}
		if err := tx.Model(&item).Updates(map[string]any{"name": name, "contains_package_id": input.ContainsPackageID,
			"contains_quantity": quantity}).Error; err != nil {
			return err
		}
		return s.templateAudit(tx, principal, hqID, "productPackageTemplate", id)
	})
	if err != nil {
		return nil, err
	}
	item.Name, item.ContainsPackageID = name, input.ContainsPackageID
	if input.ContainsPackageID == nil {
		item.ContainsQuantity = nil
	} else {
		item.ContainsQuantity = &input.ContainsQuantity
	}
	return &item, nil
}

func validateTemplateChild(tx *gorm.DB, hqID, selfID string, childID *string) error {
	seen := map[string]bool{}
	for next := childID; next != nil; {
		if *next == selfID || seen[*next] {
			return auth.NewError(auth.CodeValidationFailed)
		}
		seen[*next] = true
		var child gen.ProductPackageTemplate
		if err := tx.Where("id = ? AND organization_id = ? AND enabled = ?", *next, hqID, true).First(&child).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		next = child.ContainsPackageID
	}
	return nil
}

func (s *Service) SetPackageTemplateEnabled(ctx context.Context, principal *auth.WorkspacePrincipal, id string, enabled bool) (*gen.ProductPackageTemplate, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	var item gen.ProductPackageTemplate
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockCategoryHierarchy(tx, hqID); err != nil {
			return err
		}
		if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if !enabled {
			var parents, defaults int64
			if err := tx.Model(&gen.ProductPackageTemplate{}).Where("organization_id = ? AND contains_package_id = ? AND enabled = ?", hqID, id, true).Count(&parents).Error; err != nil {
				return err
			}
			if err := tx.Model(&gen.Product{}).Where("default_package_template_id = ?", id).Count(&defaults).Error; err != nil {
				return err
			}
			if parents+defaults > 0 {
				return auth.NewError(auth.CodeValidationFailed)
			}
		}
		if err := tx.Model(&item).Update("enabled", enabled).Error; err != nil {
			return err
		}
		return s.templateAudit(tx, principal, hqID, "productPackageTemplate", id)
	})
	item.Enabled = enabled
	return &item, err
}

func (s *Service) DeletePackageTemplate(ctx context.Context, principal *auth.WorkspacePrincipal, id string) error {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessDelete)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item gen.ProductPackageTemplate
		if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		var descendants, copies, defaults int64
		if err := tx.Model(&gen.ProductPackageTemplate{}).Where("contains_package_id = ?", id).Count(&descendants).Error; err != nil {
			return err
		}
		if err := tx.Model(&gen.ProductPackage{}).Where("template_id = ?", id).Count(&copies).Error; err != nil {
			return err
		}
		if err := tx.Model(&gen.Product{}).Where("default_package_template_id = ?", id).Count(&defaults).Error; err != nil {
			return err
		}
		if descendants+copies+defaults != 0 {
			return auth.NewError(auth.CodeValidationFailed)
		}
		if err := tx.Unscoped().Delete(&item).Error; err != nil {
			return err
		}
		return s.templateAudit(tx, principal, hqID, "productPackageTemplate", id)
	})
}
