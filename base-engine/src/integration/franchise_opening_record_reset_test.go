package integration_test

import (
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authentication"
)

func TestResetRejectsOpeningRecordThatDisagreesWithMarker(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedInitialPasswordTarget(t, fixture)
	var before authentication.AccountCredential
	if err := fixture.db.First(&before, "account_id = ?", "account-initial").Error; err != nil {
		t.Fatal(err)
	}
	fixture.createAll([]gen.FranchiseOpeningRecord{{ID: "opening-a", RecordNumber: "OPEN-1", Source: gen.FranchiseOpeningSourceHistoricalAttestation, OrganizationID: "org-a", InitialAccountID: "account-staff", RecordedByAccountID: "account-shared"}})
	response := fixture.execute("session-hq", `mutation { resetFranchiseInitialPassword(organizationId: "org-a") { accountId } }`)
	assertCode(t, response, auth.CodeConflict)
	var after authentication.AccountCredential
	if err := fixture.db.First(&after, "account_id = ?", "account-initial").Error; err != nil {
		t.Fatal(err)
	}
	if after.PasswordHash != before.PasswordHash {
		t.Fatal("rejected reset rotated password")
	}
}
