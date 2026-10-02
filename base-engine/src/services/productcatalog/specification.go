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

func (s *Service) Specifications(ctx context.Context, principal *auth.WorkspacePrincipal, filter CatalogFilter, page, perPage int) (*CatalogPage[*gen.SpecificationDefinition], error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:read", authorization.AccessRead)
	if err != nil {
		return nil, err
	}
	offset, limit, err := pageBounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := applyCatalogFilter(s.db.WithContext(ctx).Model(&gen.SpecificationDefinition{}).Where("organization_id = ?", hqID), "specification_definitions", filter)
	result := &CatalogPage[*gen.SpecificationDefinition]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	return result, query.Order("name, id").Offset(offset).Limit(limit).Find(&result.Data).Error
}

func (s *Service) SpecificationValues(ctx context.Context, principal *auth.WorkspacePrincipal, specificationID *string, filter CatalogFilter, page, perPage int) (*CatalogPage[*gen.SpecificationValue], error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:read", authorization.AccessRead)
	if err != nil {
		return nil, err
	}
	offset, limit, err := pageBounds(page, perPage)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&gen.SpecificationValue{}).Joins("JOIN specification_definitions ON specification_definitions.id = specification_values.specification_id").Where("specification_definitions.organization_id = ?", hqID)
	if specificationID != nil {
		query = query.Where("specification_values.specification_id = ?", *specificationID)
	}
	query = applyCatalogFilter(query, "specification_values", filter)
	result := &CatalogPage[*gen.SpecificationValue]{Page: page, PerPage: perPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	return result, query.Order("specification_values.weight, specification_values.name, specification_values.id").Offset(offset).Limit(limit).Find(&result.Data).Error
}

func (s *Service) CreateSpecification(ctx context.Context, principal *auth.WorkspacePrincipal, rawName string) (*gen.SpecificationDefinition, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(rawName)
	if !validName(name, 128) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	item := &gen.SpecificationDefinition{ID: uuid.Must(uuid.NewV4()).String(), Name: name, OrganizationID: hqID, Enabled: true}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := availableSpecificationName(tx, hqID, name, ""); err != nil {
			return err
		}
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		return s.templateAudit(tx, principal, hqID, "specificationDefinition", item.ID)
	})
	return item, err
}

func (s *Service) UpdateSpecification(ctx context.Context, principal *auth.WorkspacePrincipal, id, rawName string) (*gen.SpecificationDefinition, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(rawName)
	if !validName(name, 128) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var item gen.SpecificationDefinition
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err := availableSpecificationName(tx, hqID, name, id); err != nil {
			return err
		}
		if err := tx.Model(&item).Update("name", name).Error; err != nil {
			return err
		}
		return s.templateAudit(tx, principal, hqID, "specificationDefinition", id)
	})
	item.Name = name
	return &item, err
}

func (s *Service) SetSpecificationEnabled(ctx context.Context, principal *auth.WorkspacePrincipal, id string, enabled bool) (*gen.SpecificationDefinition, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	var item gen.SpecificationDefinition
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err := tx.Model(&item).Update("enabled", enabled).Error; err != nil {
			return err
		}
		return s.templateAudit(tx, principal, hqID, "specificationDefinition", id)
	})
	item.Enabled = enabled
	return &item, err
}

func (s *Service) DeleteSpecification(ctx context.Context, principal *auth.WorkspacePrincipal, id string) error {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessDelete)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item gen.SpecificationDefinition
		if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		var count int64
		if err := tx.Model(&gen.SpecificationValue{}).Where("specification_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count != 0 {
			return auth.NewError(auth.CodeValidationFailed)
		}
		if err := tx.Unscoped().Delete(&item).Error; err != nil {
			return err
		}
		return s.templateAudit(tx, principal, hqID, "specificationDefinition", id)
	})
}

func availableSpecificationName(tx *gorm.DB, hqID, name, exceptID string) error {
	query := tx.Model(&gen.SpecificationDefinition{}).Where("organization_id = ? AND name = ?", hqID, name)
	if exceptID != "" {
		query = query.Where("id <> ?", exceptID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}
