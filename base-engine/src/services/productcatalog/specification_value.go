package productcatalog

import (
	"context"
	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

func (s *Service) CreateSpecificationValue(ctx context.Context, principal *auth.WorkspacePrincipal, specificationID, rawName string) (*gen.SpecificationValue, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(rawName)
	if !validName(name, 128) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	item := &gen.SpecificationValue{ID: uuid.Must(uuid.NewV4()).String(), SpecificationID: specificationID, Name: name, Enabled: true}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var definition gen.SpecificationDefinition
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND organization_id = ? AND enabled = ?", specificationID, hqID, true).First(&definition).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err := availableSpecificationValueName(tx, specificationID, name, ""); err != nil {
			return err
		}
		var maxWeight int64
		if err := tx.Model(&gen.SpecificationValue{}).Where("specification_id = ?", specificationID).
			Select("COALESCE(MAX(weight), 0)").Scan(&maxWeight).Error; err != nil {
			return err
		}
		weight := maxWeight + 1
		item.Weight = &weight
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		return s.templateAudit(tx, principal, hqID, "specificationValue", item.ID)
	})
	return item, err
}

// ReorderSpecificationValues replaces the complete order within one headquarters specification.
func (s *Service) ReorderSpecificationValues(ctx context.Context, principal *auth.WorkspacePrincipal, specificationID string, orderedIDs []string) error {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var definition gen.SpecificationDefinition
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND organization_id = ?", specificationID, hqID).First(&definition).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		var values []gen.SpecificationValue
		if err := tx.Where("specification_id = ?", specificationID).Find(&values).Error; err != nil {
			return err
		}
		if err := validateSpecificationOrder(values, orderedIDs); err != nil {
			return err
		}
		for index, id := range orderedIDs {
			if err := tx.Model(&gen.SpecificationValue{}).Where("id = ? AND specification_id = ?", id, specificationID).
				Update("weight", index+1).Error; err != nil {
				return err
			}
			if err := s.templateAudit(tx, principal, hqID, "specificationValue", id); err != nil {
				return err
			}
		}
		return nil
	})
}

func validateSpecificationOrder(values []gen.SpecificationValue, orderedIDs []string) error {
	if len(values) != len(orderedIDs) {
		return auth.NewError(auth.CodeValidationFailed)
	}
	valid := make(map[string]struct{}, len(values))
	for _, value := range values {
		valid[value.ID] = struct{}{}
	}
	for _, id := range orderedIDs {
		if _, exists := valid[id]; !exists {
			return auth.NewError(auth.CodeValidationFailed)
		}
		delete(valid, id)
	}
	return nil
}

func (s *Service) UpdateSpecificationValue(ctx context.Context, principal *auth.WorkspacePrincipal, id, rawName string) (*gen.SpecificationValue, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(rawName)
	if !validName(name, 128) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var item gen.SpecificationValue
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := specificationValueForHQ(tx, hqID, id, &item); err != nil {
			return err
		}
		if err := availableSpecificationValueName(tx, item.SpecificationID, name, id); err != nil {
			return err
		}
		if err := tx.Model(&item).Update("name", name).Error; err != nil {
			return err
		}
		return s.templateAudit(tx, principal, hqID, "specificationValue", id)
	})
	item.Name = name
	return &item, err
}

func (s *Service) SetSpecificationValueEnabled(ctx context.Context, principal *auth.WorkspacePrincipal, id string, enabled bool) (*gen.SpecificationValue, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	var item gen.SpecificationValue
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := specificationValueForHQ(tx, hqID, id, &item); err != nil {
			return err
		}
		if err := tx.Model(&item).Update("enabled", enabled).Error; err != nil {
			return err
		}
		return s.templateAudit(tx, principal, hqID, "specificationValue", id)
	})
	item.Enabled = enabled
	return &item, err
}

func (s *Service) DeleteSpecificationValue(ctx context.Context, principal *auth.WorkspacePrincipal, id string) error {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessDelete)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item gen.SpecificationValue
		if err := specificationValueForHQ(tx, hqID, id, &item); err != nil {
			return err
		}
		var choices, skus int64
		if err := tx.Model(&gen.ProductSpecificationChoice{}).Where("value_id = ?", id).Count(&choices).Error; err != nil {
			return err
		}
		if err := tx.Model(&gen.ProductSkuSpecificationValue{}).Where("value_id = ?", id).Count(&skus).Error; err != nil {
			return err
		}
		if choices+skus != 0 {
			return auth.NewError(auth.CodeValidationFailed)
		}
		if err := tx.Unscoped().Delete(&item).Error; err != nil {
			return err
		}
		return s.templateAudit(tx, principal, hqID, "specificationValue", id)
	})
}

func specificationValueForHQ(tx *gorm.DB, hqID, id string, item *gen.SpecificationValue) error {
	if err := tx.Joins("JOIN specification_definitions ON specification_definitions.id = specification_values.specification_id").Where("specification_values.id = ? AND specification_definitions.organization_id = ?", id, hqID).First(item).Error; err != nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return nil
}

func availableSpecificationValueName(tx *gorm.DB, specificationID, name, exceptID string) error {
	query := tx.Model(&gen.SpecificationValue{}).Where("specification_id = ? AND name = ?", specificationID, name)
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
