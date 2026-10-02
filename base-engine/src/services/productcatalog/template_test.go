package productcatalog

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func createPackageChain(t *testing.T, service *Service, principal *auth.WorkspacePrincipal) *gen.ProductPackageTemplate {
	t.Helper()
	ctx := context.Background()
	unit := mustCatalog[*gen.ProductPackageTemplate](t)(service.CreatePackageTemplate(ctx, principal, PackageTemplateInput{Name: "pack"}))
	bag := mustCatalog[*gen.ProductPackageTemplate](t)(service.CreatePackageTemplate(ctx, principal, PackageTemplateInput{Name: "bag", ContainsPackageID: &unit.ID, ContainsQuantity: 10}))
	return mustCatalog[*gen.ProductPackageTemplate](t)(service.CreatePackageTemplate(ctx, principal, PackageTemplateInput{Name: "box", ContainsPackageID: &bag.ID, ContainsQuantity: 12}))
}

func packagesForSku(t *testing.T, db *gorm.DB, skuID string) []gen.ProductPackage {
	t.Helper()
	var packages []gen.ProductPackage
	mustCatalogErrorFree(t, db.Where("sku_id = ?", skuID).Find(&packages).Error)
	return packages
}

func TestPackageTemplateRejectsCyclesAndForeignReferences(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	unit, err := service.CreatePackageTemplate(ctx, principal, PackageTemplateInput{Name: "unit"})
	if err != nil {
		t.Fatal(err)
	}
	bag, err := service.CreatePackageTemplate(ctx, principal, PackageTemplateInput{Name: "bag", ContainsPackageID: &unit.ID, ContainsQuantity: 10})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.UpdatePackageTemplate(ctx, principal, unit.ID, PackageTemplateInput{Name: "unit", ContainsPackageID: &bag.ID, ContainsQuantity: 2}); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("cycle error = %v", err)
	}
	foreign := gen.ProductPackageTemplate{ID: "foreign", OrganizationID: "other-hq", Name: "foreign", Enabled: true}
	if err := db.Create(&foreign).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.CreatePackageTemplate(ctx, principal, PackageTemplateInput{Name: "wrong", ContainsPackageID: &foreign.ID, ContainsQuantity: 2}); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("foreign child error = %v", err)
	}
}

func TestPackageTemplateCannotDisableReferencedChild(t *testing.T) {
	service, db, principal := catalogFixture(t)
	ctx := context.Background()
	unit := mustCatalog[*gen.ProductPackageTemplate](t)(service.CreatePackageTemplate(ctx, principal, PackageTemplateInput{Name: "pack"}))
	_, err := service.CreatePackageTemplate(ctx, principal, PackageTemplateInput{Name: "box", ContainsPackageID: &unit.ID, ContainsQuantity: 12})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetPackageTemplateEnabled(ctx, principal, unit.ID, false); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("disable referenced child = %v", err)
	}
	var persisted gen.ProductPackageTemplate
	if err := db.First(&persisted, "id = ?", unit.ID).Error; err != nil || !persisted.Enabled {
		t.Fatalf("referenced child changed: %+v, %v", persisted, err)
	}
}

func TestPackageTemplatesFilterByContainedUnitBeforePagination(t *testing.T) {
	service, _, principal := catalogFixture(t)
	ctx := context.Background()
	unit := mustCatalog[*gen.ProductPackageTemplate](t)(service.CreatePackageTemplate(ctx, principal, PackageTemplateInput{Name: "unit"}))
	bag := mustCatalog[*gen.ProductPackageTemplate](t)(service.CreatePackageTemplate(ctx, principal, PackageTemplateInput{Name: "bag", ContainsPackageID: &unit.ID, ContainsQuantity: 10}))
	mustCatalog[*gen.ProductPackageTemplate](t)(service.CreatePackageTemplate(ctx, principal, PackageTemplateInput{Name: "box", ContainsPackageID: &bag.ID, ContainsQuantity: 12}))

	result, err := service.PackageTemplates(ctx, principal, CatalogFilter{ContainsPackageID: &unit.ID}, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 || len(result.Data) != 1 || result.Data[0].ID != bag.ID {
		t.Fatalf("contained unit filter = %+v", result)
	}
}
