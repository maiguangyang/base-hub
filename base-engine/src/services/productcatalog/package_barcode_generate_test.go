package productcatalog

import (
	"bytes"
	"context"
	"base-engine/auth"
	"base-engine/gen"
	"regexp"
	"testing"
)

func TestGeneratePackageBarcodeDoesNotPersist(t *testing.T) {
	service, db, principal := catalogFixture(t)
	var before, after int64
	catalogNoError(t, db.Model(&gen.ProductPackage{}).Count(&before).Error)
	code, err := service.GeneratePackageBarcode(context.Background(), principal)
	catalogNoError(t, err)
	if !regexp.MustCompile(`^KH[0-9]{12}$`).MatchString(code) {
		t.Fatalf("format: %q", code)
	}
	catalogNoError(t, db.Model(&gen.ProductPackage{}).Count(&after).Error)
	if before != after {
		t.Fatal("candidate generation persisted a package")
	}
	principal.Permissions = map[string]struct{}{}
	_, err = service.GeneratePackageBarcode(context.Background(), principal)
	if auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("unauthorized generation: %v", err)
	}
}

func TestGeneratePackageBarcodeCollisionAndEntropyFailure(t *testing.T) {
	service, _, principal := catalogFixture(t)
	barcodePackage(t, service, principal, "KH000000000000")
	_, err := generatePackageBarcode(context.Background(), service.db, bytes.NewReader(make([]byte, 128)))
	if auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("collision: %v", err)
	}
	_, err = generatePackageBarcode(context.Background(), service.db, bytes.NewReader(nil))
	if err == nil {
		t.Fatal("entropy failure ignored")
	}
}
