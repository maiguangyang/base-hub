package customer

import (
	"context"
	"strings"
	"testing"

	"base-engine/gen"
	"base-engine/src/services/audit"
)

func TestAICustomerCreationAuditExcludesPhone(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	ctx := audit.WithAIInvocation(context.Background(), "run-1", "HqCreateCustomerMember")
	if _, err := service.CreateMember(ctx, principal, customerCreateInput("ai-create-1")); err != nil {
		t.Fatal(err)
	}
	var record gen.AuditLog
	if err := db.Where("action = ?", "hqCustomer:create").First(&record).Error; err != nil || record.MetadataJSON == nil {
		t.Fatalf("missing AI audit: %v", err)
	}
	if strings.Contains(*record.MetadataJSON, "+8613800000001") || strings.Contains(*record.MetadataJSON, "13800000001") ||
		!strings.Contains(*record.MetadataJSON, "HqCreateCustomerMember") {
		t.Fatal("AI audit must link the tool without storing the phone")
	}
}
