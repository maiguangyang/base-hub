package customer

import (
	"strings"
	"testing"

	"base-engine/gen"
	"gorm.io/gorm"
)

func assertCancellationAuditBasis(t *testing.T, db *gorm.DB) {
	t.Helper()
	var requestAudit gen.AuditLog
	if err := db.Where("action = ?", "hqCustomer:cancel_request").First(&requestAudit).Error; err != nil {
		t.Fatal(err)
	}
	if requestAudit.MetadataJSON == nil || !strings.Contains(*requestAudit.MetadataJSON, "MEMBERSHIP_END") {
		t.Fatalf("cancellation basis missing from audit: %v", requestAudit.MetadataJSON)
	}
}
