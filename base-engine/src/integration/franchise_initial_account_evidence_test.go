package integration_test

import (
	"testing"

	"base-engine/gen"
)

func TestLegacyManualMarkerCanResetWithoutReverification(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedInitialPasswordTarget(t, fixture)
	legacyMetadata := `{}`
	fixture.createAll([]gen.AuditLog{{ID: "legacy-confirmation", Action: "franchiseInitialAccount:confirm", ResourceType: "account", ResourceID: evidenceString("account-initial"), OrganizationID: evidenceString("org-a"), ResultCode: "SUCCESS", MetadataJSON: &legacyMetadata}})
	query := `mutation { resetFranchiseInitialPassword(organizationId: "org-a") { accountId } }`
	response := fixture.execute("session-hq", query)
	if len(response.Errors) != 0 {
		t.Fatalf("existing marker could not reset: %s", response.Body)
	}
}

func evidenceString(value string) *string { return &value }
