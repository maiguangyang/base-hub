package integration_test

import (
	"strings"
	"testing"

	"base-engine/gen"
)

func assertHistoricalConfirmationAudit(t *testing.T, fixture *securityFixture) {
	t.Helper()
	var record gen.AuditLog
	if err := fixture.db.First(&record, "action = ?", "franchiseInitialAccount:confirm").Error; err != nil {
		t.Fatal(err)
	}
	if record.ResourceID == nil || *record.ResourceID != "account-staff" || record.ActorAccountID == nil || *record.ActorAccountID != "account-hq" {
		t.Fatalf("confirmation audit identity invalid: %#v", record)
	}
	if record.MetadataJSON == nil || !strings.Contains(*record.MetadataJSON, "OPEN-2023-001") || !strings.Contains(*record.MetadataJSON, `"source":"historical_manual_attestation"`) {
		t.Fatalf("confirmation audit attestation invalid: %#v", record)
	}
}
