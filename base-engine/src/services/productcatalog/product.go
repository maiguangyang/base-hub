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
)

type ProductInput struct {
	Name, CategoryID         string
	BrandID, Description     *string
	Selections               []SpecificationSelection
	DefaultPackageTemplateID *string
	SkuOverrides             []SkuOverride
}
type SkuInput struct {
	Name, ProductID                             string
	Ingredients, Allergens, StorageInstructions *string
	ShelfLifeDays                               *int64
}

func (s *Service) CreateProduct(ctx context.Context, principal *auth.WorkspacePrincipal, input ProductInput) (*gen.Product, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if !validName(name, 128) || input.CategoryID == "" {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	item := &gen.Product{ID: uuid.Must(uuid.NewV4()).String(), Name: name,
		OrganizationID: hqID, CategoryID: input.CategoryID, BrandID: input.BrandID, Description: input.Description, Enabled: true}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockCategoryHierarchy(tx, hqID); err != nil {
			return err
		}
		brand, err := selectedBrand(tx, hqID, input.BrandID)
		if err != nil {
			return err
		}
		if err := insertProductRecord(tx, item); err != nil {
			return err
		}
		if err := s.saveProductVariants(tx, principal, hqID, item, input); err != nil {
			return err
		}
		item.Brand = brand
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "product", ResourceID: item.ID, ResultCode: "SUCCESS"})
	})
	return item, err
}

func insertProductRecord(tx *gorm.DB, item *gen.Product) error {
	var category gen.ProductCategory
	if err := tx.Where("id = ? AND organization_id = ? AND enabled = ?", item.CategoryID, item.OrganizationID, true).First(&category).Error; err != nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return tx.Create(item).Error
}

func (s *Service) CreateSku(ctx context.Context, principal *auth.WorkspacePrincipal, input SkuInput) (*gen.ProductSku, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if !validName(name, 128) || input.ProductID == "" || (input.ShelfLifeDays != nil && *input.ShelfLifeDays <= 0) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	item := &gen.ProductSku{ID: uuid.Must(uuid.NewV4()).String(), Name: name, ProductID: input.ProductID,
		Ingredients: input.Ingredients, Allergens: input.Allergens,
		StorageInstructions: input.StorageInstructions, ShelfLifeDays: input.ShelfLifeDays, Enabled: true}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product gen.Product
		if err := tx.Where("id = ? AND organization_id = ? AND enabled = ?", input.ProductID, hqID, true).First(&product).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "productSku", ResourceID: item.ID, ResultCode: "SUCCESS"})
	})
	return item, err
}

func validName(value string, max int) bool {
	return value != "" && utf8.RuneCountInString(value) <= max
}
