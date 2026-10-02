package productcatalog

import (
	"context"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func barcodePackage(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, code string) (*gen.ProductSku, *gen.ProductPackage) {
	t.Helper()
	ctx := context.Background()
	category, err := service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"})
	catalogNoError(t, err)
	product, err := service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID})
	catalogNoError(t, err)
	sku, err := service.CreateSku(ctx, principal, SkuInput{Name: "Spicy", ProductID: product.ID})
	catalogNoError(t, err)
	pack, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "unit", Barcode: code})
	catalogNoError(t, err)
	return sku, pack
}

func TestSetPackageBarcodePreservesIdentity(t *testing.T) {
	service, db, principal := catalogFixture(t)
	_, unit := barcodePackage(t, service, principal, "")
	updated, err := service.SetPackageBarcode(context.Background(), principal, unit.ID, " 0012345678905 ")
	catalogNoError(t, err)
	if updated.Barcode == nil || *updated.Barcode != "0012345678905" || updated.ID != unit.ID || updated.SkuID != unit.SkuID || updated.PackageSetVersion != unit.PackageSetVersion {
		t.Fatalf("barcode update changed identity: %+v", updated)
	}
	var before, after int64
	catalogNoError(t, db.Model(&gen.AuditLog{}).Count(&before).Error)
	_, err = service.SetPackageBarcode(context.Background(), principal, unit.ID, "0012345678905")
	catalogNoError(t, err)
	catalogNoError(t, db.Model(&gen.AuditLog{}).Count(&after).Error)
	if before != after {
		t.Fatal("same value produced another audit")
	}
}

func TestEmptyPackageBarcodeRemainsUnset(t *testing.T) {
	service, _, principal := catalogFixture(t)
	_, unit := barcodePackage(t, service, principal, "")
	if unit.Barcode != nil {
		t.Fatal("empty barcode must remain unset")
	}
}

func TestSetPackageBarcodeClearsValue(t *testing.T) {
	service, db, principal := catalogFixture(t)
	_, unit := barcodePackage(t, service, principal, "original")
	updated, err := service.SetPackageBarcode(context.Background(), principal, unit.ID, " ")
	catalogNoError(t, err)
	if updated.Barcode != nil {
		t.Fatal("cleared barcode must be NULL")
	}
	var found gen.ProductPackage
	catalogNoError(t, db.First(&found, "id = ?", unit.ID).Error)
	if found.Barcode != nil {
		t.Fatal("clearing was not persisted")
	}
}

func TestSetPackageBarcodeRejectsInvalidInput(t *testing.T) {
	service, db, principal := catalogFixture(t)
	_, unit := barcodePackage(t, service, principal, "original")
	for _, code := range []string{strings.Repeat("a", 65), strings.Repeat("袋", 22)} {
		_, err := service.SetPackageBarcode(context.Background(), principal, unit.ID, code)
		if auth.ErrorCode(err) != auth.CodeValidationFailed {
			t.Fatalf("invalid barcode: %v", err)
		}
	}
	var found gen.ProductPackage
	catalogNoError(t, db.First(&found, "id = ?", unit.ID).Error)
	if found.Barcode == nil || *found.Barcode != "original" {
		t.Fatal("invalid input changed barcode")
	}
}

func TestPackageBarcodeVersionReuse(t *testing.T) {
	service, _, principal := catalogFixture(t)
	sku, unit := barcodePackage(t, service, principal, "00123")
	ctx := context.Background()
	draft, err := service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "unit", Barcode: "00123", PackageSetVersion: 2})
	catalogNoError(t, err)
	_, err = service.CreatePackage(ctx, principal, PackageInput{SkuID: sku.ID, Name: "box", Barcode: "00123", PackageSetVersion: 2, ContainsPackageID: &draft.ID, ContainsQuantity: 5})
	if auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("same version duplicate: %v", err)
	}
	_, err = service.PublishPackageSet(ctx, principal, sku.ID, 2)
	catalogNoError(t, err)
	_, err = service.SetPackageBarcode(ctx, principal, unit.ID, "retired")
	if auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("retired barcode changed: %v", err)
	}
}

func TestPackageBarcodeActiveConflict(t *testing.T) {
	service, db, principal := catalogFixture(t)
	sku, unit := barcodePackage(t, service, principal, "first")
	other := gen.ProductSku{ID: "other-sku", ProductID: sku.ProductID, Name: "Other", Enabled: true}
	catalogNoError(t, db.Select("*").Create(&other).Error)
	pack, err := service.CreatePackage(context.Background(), principal, PackageInput{SkuID: other.ID, Name: "unit", Barcode: "second"})
	catalogNoError(t, err)
	_, err = service.SetPackageBarcode(context.Background(), principal, unit.ID, "second")
	if auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("active conflict: %v", err)
	}
	_, err = service.CreatePackage(context.Background(), principal, PackageInput{SkuID: other.ID, Name: "unit", Barcode: "first", PackageSetVersion: 2})
	if auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("foreign draft conflict: %v", err)
	}
	foreign := *principal
	hq := "foreign-hq"
	foreign.OrganizationID = &hq
	_, err = service.SetPackageBarcode(context.Background(), &foreign, pack.ID, "stolen")
	if auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("foreign access: %v", err)
	}
	principal.Permissions = map[string]struct{}{}
	_, err = service.SetPackageBarcode(context.Background(), principal, unit.ID, "denied")
	if auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("missing permission: %v", err)
	}
}
