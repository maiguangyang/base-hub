package productcatalog

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func testPNG(t *testing.T) []byte {
	t.Helper()
	var output bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	if err := png.Encode(&output, img); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}

func TestMainImageUploadBindReplaceRemoveAndDelete(t *testing.T) {
	service, db, principal := catalogFixture(t)
	service.SetImageRoot(t.TempDir())
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID}))
	first := mustCatalog[string](t)(service.UploadMainImage(ctx, principal, testPNG(t)))
	bound := mustCatalog[*gen.Product](t)(service.SetProductMainImage(ctx, principal, product.ID, first))
	assertImageExists(t, service, product.ID, bound.ImageURL)
	if _, err := service.SetProductMainImage(ctx, principal, product.ID, first); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("reuse = %v", err)
	}
	replaced := testImageReplacement(t, service, db, principal, product.ID, bound.ImageURL)
	mustCatalog[*gen.Product](t)(service.RemoveProductMainImage(ctx, principal, product.ID))
	assertImageMissing(t, service, product.ID, replaced.ImageURL)
	testInvalidImageUploads(t, service, principal)
	third := mustCatalog[string](t)(service.UploadMainImage(ctx, principal, testPNG(t)))
	finalImage := mustCatalog[*gen.Product](t)(service.SetProductMainImage(ctx, principal, product.ID, third))
	mustCatalogErrorFree(t, service.DeleteProduct(ctx, principal, product.ID))
	assertImageMissing(t, service, product.ID, finalImage.ImageURL)
}

func testImageReplacement(t *testing.T, service *Service, db *gorm.DB, principal *auth.WorkspacePrincipal, productID string, oldURL *string) *gen.Product {
	t.Helper()
	ctx := context.Background()
	second := mustCatalog[string](t)(service.UploadMainImage(ctx, principal, testPNG(t)))
	other := *principal
	other.SessionID = "other-session"
	if _, err := service.SetProductMainImage(ctx, &other, productID, second); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("cross-session = %v", err)
	}
	replaced := mustCatalog[*gen.Product](t)(service.SetProductMainImage(ctx, principal, productID, second))
	if replaced.ImageURL == nil || *replaced.ImageURL == *oldURL {
		t.Fatalf("replacement = %+v", replaced)
	}
	var previousRefs int64
	mustCatalogErrorFree(t, db.Model(&gen.Product{}).Where("image_url = ?", *oldURL).Count(&previousRefs).Error)
	if previousRefs != 0 {
		t.Fatalf("old image has %d references", previousRefs)
	}
	assertImageMissing(t, service, productID, oldURL)
	return replaced
}

func assertImageExists(t *testing.T, service *Service, productID string, url *string) {
	t.Helper()
	if url == nil {
		t.Fatal("missing image URL")
	}
	path, ok := managedImagePath(service.ImageRoot(), productID, *url)
	if !ok {
		t.Fatal("unmanaged image URL")
	}
	mustCatalogErrorFree(t, func() error { _, err := os.Stat(path); return err }())
}

func assertImageMissing(t *testing.T, service *Service, productID string, url *string) {
	t.Helper()
	if url == nil {
		t.Fatal("missing image URL")
	}
	path, ok := managedImagePath(service.ImageRoot(), productID, *url)
	if !ok {
		t.Fatal("unmanaged image URL")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("image file = %v", err)
	}
}

func testInvalidImageUploads(t *testing.T, service *Service, principal *auth.WorkspacePrincipal) {
	t.Helper()
	ctx := context.Background()
	if _, err := service.UploadMainImage(ctx, principal, []byte("not an image")); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("invalid bytes = %v", err)
	}
	if _, err := service.UploadMainImage(ctx, principal, bytes.Repeat([]byte{'x'}, 5<<20+1)); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("large bytes = %v", err)
	}
}

func TestImageSweepRemovesExpiredAndOrphanFiles(t *testing.T) {
	service, _, principal := catalogFixture(t)
	service.SetImageRoot(t.TempDir())
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID}))
	attachment := mustCatalog[string](t)(service.UploadMainImage(ctx, principal, testPNG(t)))
	bound := mustCatalog[*gen.Product](t)(service.SetProductMainImage(ctx, principal, product.ID, attachment))
	path, ok := managedImagePath(service.ImageRoot(), product.ID, *bound.ImageURL)
	if !ok {
		t.Fatal("unmanaged bound image")
	}
	orphan := filepath.Join(filepath.Dir(path), "00000000-0000-4000-8000-000000000000.png")
	mustCatalogErrorFree(t, os.WriteFile(orphan, testPNG(t), 0600))
	stale := filepath.Join(service.ImageRoot(), "pending", "stale.png")
	mustCatalogErrorFree(t, os.WriteFile(stale, testPNG(t), 0600))
	old := time.Now().Add(-2 * time.Hour)
	mustCatalogErrorFree(t, os.Chtimes(stale, old, old))
	mustCatalogErrorFree(t, service.SweepImageFiles(ctx))
	assertImageExists(t, service, product.ID, bound.ImageURL)
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatalf("orphan = %v", err)
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatalf("stale = %v", err)
	}
}

func TestImageUploadDoesNotRescanTheWholeCatalogEveryTime(t *testing.T) {
	service, _, principal := catalogFixture(t)
	service.SetImageRoot(t.TempDir())
	ctx := context.Background()
	category := mustCatalog[*gen.ProductCategory](t)(service.CreateCategory(ctx, principal, CategoryInput{Name: "Food"}))
	product := mustCatalog[*gen.Product](t)(service.CreateProduct(ctx, principal, ProductInput{Name: "Noodles", CategoryID: category.ID}))
	attachment := mustCatalog[string](t)(service.UploadMainImage(ctx, principal, testPNG(t)))
	bound := mustCatalog[*gen.Product](t)(service.SetProductMainImage(ctx, principal, product.ID, attachment))
	path, _ := managedImagePath(service.ImageRoot(), product.ID, *bound.ImageURL)
	orphan := filepath.Join(filepath.Dir(path), "00000000-0000-4000-8000-000000000000.png")
	mustCatalogErrorFree(t, os.WriteFile(orphan, testPNG(t), 0600))
	mustCatalog[string](t)(service.UploadMainImage(ctx, principal, testPNG(t)))
	if _, err := os.Stat(orphan); err != nil {
		t.Fatalf("a second upload should defer the full sweep: %v", err)
	}
	mustCatalogErrorFree(t, service.SweepImageFiles(ctx))
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatalf("manual sweep must remove the orphan: %v", err)
	}
}

func TestImageBindRejectsProductPathTraversal(t *testing.T) {
	service, _, principal := catalogFixture(t)
	service.SetImageRoot(t.TempDir())
	ctx := context.Background()
	attachment := mustCatalog[string](t)(service.UploadMainImage(ctx, principal, testPNG(t)))
	if _, err := service.SetProductMainImage(ctx, principal, "../escape", attachment); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("traversal error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(service.ImageRoot(), "escape")); !os.IsNotExist(err) {
		t.Fatalf("outside directory = %v", err)
	}
}
