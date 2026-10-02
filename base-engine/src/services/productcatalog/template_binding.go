package productcatalog

import (
	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func collectPackageTemplates(tx *gorm.DB, hqID string, roots []string) ([]gen.ProductPackageTemplate, error) {
	collector := packageTemplateCollector{tx: tx, hqID: hqID, state: map[string]uint8{}}
	for _, id := range roots {
		if err := collector.visit(id); err != nil {
			return nil, err
		}
	}
	return collector.ordered, nil
}

type packageTemplateCollector struct {
	tx      *gorm.DB
	hqID    string
	state   map[string]uint8
	ordered []gen.ProductPackageTemplate
}

func (c *packageTemplateCollector) visit(id string) error {
	if c.state[id] == 2 {
		return nil
	}
	if id == "" || c.state[id] == 1 || len(c.state) >= 100 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	c.state[id] = 1
	var item gen.ProductPackageTemplate
	if err := c.tx.Where("id = ? AND organization_id = ? AND enabled = ?", id, c.hqID, true).First(&item).Error; err != nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	if err := c.visitContained(&item); err != nil {
		return err
	}
	c.state[id] = 2
	c.ordered = append(c.ordered, item)
	return nil
}

func (c *packageTemplateCollector) visitContained(item *gen.ProductPackageTemplate) error {
	if item.ContainsPackageID == nil {
		if item.ContainsQuantity != nil {
			return auth.NewError(auth.CodeValidationFailed)
		}
		return nil
	}
	if item.ContainsQuantity == nil || *item.ContainsQuantity < 2 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return c.visit(*item.ContainsPackageID)
}

func (s *Service) copyTemplatePackages(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID string, sku *gen.ProductSku, version int64, templates []gen.ProductPackageTemplate) error {
	created := map[string]string{}
	for _, template := range templates {
		item := &gen.ProductPackage{ID: uuid.Must(uuid.NewV4()).String(), SkuID: sku.ID, TemplateID: &template.ID,
			Name: template.Name, PackageSetVersion: version, Enabled: true}
		input := PackageInput{SkuID: sku.ID, Name: template.Name, PackageSetVersion: version}
		if template.ContainsPackageID != nil {
			child, ok := created[*template.ContainsPackageID]
			if !ok {
				return auth.NewError(auth.CodeValidationFailed)
			}
			item.ContainsPackageID, item.ContainsQuantity = &child, template.ContainsQuantity
			input.ContainsPackageID, input.ContainsQuantity = &child, *template.ContainsQuantity
		}
		if err := s.createPackageTx(tx, principal, hqID, input, item); err != nil {
			return err
		}
		created[template.ID] = item.ID
	}
	return nil
}
