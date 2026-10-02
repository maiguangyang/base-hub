package productcatalog

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryInput struct {
	Name     string
	ParentID *string
}

func (s *Service) CreateCategory(ctx context.Context, principal *auth.WorkspacePrincipal, input CategoryInput) (*gen.ProductCategory, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" || utf8.RuneCountInString(name) > 128 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	item := &gen.ProductCategory{
		ID: uuid.Must(uuid.NewV4()).String(), OrganizationID: hqID,
		Name: name, Enabled: true,
	}
	if input.ParentID != nil {
		item.ParentID = input.ParentID
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockCategoryHierarchy(tx, hqID); err != nil {
			return err
		}
		if err := ensureCategoryNameAvailable(tx, hqID, name, input.ParentID); err != nil {
			return err
		}
		if err := ensureCategoryParent(tx, hqID, "", input.ParentID); err != nil {
			return err
		}
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{
			ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage",
			ResourceType: "productCategory", ResourceID: item.ID, ResultCode: "SUCCESS",
		})
	})
	return item, err
}

func (s *Service) MoveCategory(ctx context.Context, principal *auth.WorkspacePrincipal, categoryID string, parentID *string) error {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return err
	}
	if categoryID == "" || (parentID != nil && *parentID == categoryID) {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var category gen.ProductCategory
		if err := lockedCategoryForHQ(tx, hqID, categoryID, &category); err != nil {
			return err
		}
		if err := ensureCategoryParent(tx, hqID, categoryID, parentID); err != nil {
			return err
		}
		if err := tx.Model(&category).Update("parent_id", parentID).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productCategory", ResourceID: category.ID, ResultCode: "SUCCESS"})
	})
}

func lockCategoryHierarchy(tx *gorm.DB, hqID string) error {
	var organization gen.Organization
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("id = ?", hqID).First(&organization).Error
}

func lockedCategoryForHQ(tx *gorm.DB, hqID, id string, item *gen.ProductCategory) error {
	if err := lockCategoryHierarchy(tx, hqID); err != nil {
		return err
	}
	if err := tx.Where("id = ? AND organization_id = ?", id, hqID).First(item).Error; err != nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return nil
}

func ensureCategoryNameAvailable(tx *gorm.DB, hqID, name string, parentID *string) error {
	return ensureCategoryNameAvailableExcept(tx, hqID, "", name, parentID)
}

func ensureCategoryNameAvailableExcept(tx *gorm.DB, hqID, exceptID, name string, parentID *string) error {
	var existing int64
	query := tx.Model(&gen.ProductCategory{}).Where("organization_id = ? AND name = ?", hqID, name)
	if exceptID != "" {
		query = query.Where("id <> ?", exceptID)
	}
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	if err := query.Count(&existing).Error; err != nil {
		return err
	}
	if existing > 0 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}

func ensureCategoryParent(tx *gorm.DB, hqID, categoryID string, parentID *string) error {
	for next := parentID; next != nil; {
		if *next == categoryID {
			return auth.NewError(auth.CodeValidationFailed)
		}
		var parent gen.ProductCategory
		if err := tx.Where("id = ? AND organization_id = ?", *next, hqID).First(&parent).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		next = parent.ParentID
	}
	return nil
}

func headquartersID(principal *auth.WorkspacePrincipal, action string, mode authorization.AccessMode) (string, error) {
	if principal == nil {
		return "", auth.NewError(auth.CodeAuthRequired)
	}
	if principal.WorkspaceType != auth.WorkspaceTypeHeadquarters || principal.OrganizationID == nil {
		return "", auth.NewError(auth.CodeWorkspaceForbidden)
	}
	if err := authorization.Authorize(principal, authorization.Intent{
		Action: action, Mode: mode, ResourceOrganizationID: principal.OrganizationID,
	}); err != nil {
		return "", err
	}
	return *principal.OrganizationID, nil
}
