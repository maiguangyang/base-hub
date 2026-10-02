package integration_test

import (
	"fmt"
	"base-engine/auth"
	"base-engine/gen"
	"regexp"
	"testing"
	"time"
)

func TestGeneratedBarcodeProtectedToolAndRevokedSession(t *testing.T) {
	fixture := barcodeRuntimeFixture(t)
	var before, after int64
	if err := fixture.db.Model(&gen.ProductPackage{}).Count(&before).Error; err != nil {
		t.Fatal(err)
	}
	result, err := runProtectedGapTool(t, fixture, "session-hq", "HqGenerateProductPackageBarcode", map[string]any{})
	if err != nil || !regexp.MustCompile(`KH[0-9]{12}`).MatchString(fmt.Sprint(result)) {
		t.Fatalf("candidate: %s %v", result, err)
	}
	if err := fixture.db.Model(&gen.ProductPackage{}).Count(&after).Error; err != nil {
		t.Fatal(err)
	}
	if before != after {
		t.Fatal("generation wrote packaging data")
	}
	if err := fixture.db.Model(&gen.Session{}).Where("id = ?", "session-hq").Update("revoked_at", time.Now()).Error; err != nil {
		t.Fatal(err)
	}
	assertCode(t, fixture.execute("session-hq", `query {hqGenerateProductPackageBarcode}`), auth.CodeSessionRevoked)
	if _, err := runProtectedGapTool(t, fixture, "session-hq", "HqGenerateProductPackageBarcode", map[string]any{}); err == nil {
		t.Fatal("revoked session generated barcode")
	}
}
