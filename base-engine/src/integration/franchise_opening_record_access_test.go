package integration_test

import (
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestOpeningRecordCannotLeakThroughRelationIDsOrFilters(t *testing.T) {
	fixture := newSecurityFixture(t)
	fixture.createAll([]gen.FranchiseOpeningRecord{{ID: "opening-a", RecordNumber: "OPEN-1", Source: gen.FranchiseOpeningSourceSystemProvision, OrganizationID: "org-a", InitialAccountID: "account-staff", RecordedByAccountID: "account-shared"}})
	seedHQCustomAdministrator(fixture, "org-reader", []string{"organization:read"})
	seedHQCustomAdministrator(fixture, "account-reader", []string{"account:read"})
	queries := []struct{ session, query string }{
		{"session-hq-org-reader", `query { organization(id: "org-a") { openingRecordsIds } }`},
		{"session-hq-account-reader", `query { account(id: "account-staff") { openingRecordsIds } }`},
		{"session-hq-account-reader", `query { account(id: "account-shared") { recordedOpeningRecordsIds } }`},
		{"session-hq-org-reader", `query { organizations(filter: { openingRecords: { recordNumber: "OPEN-1" } }) { total } }`},
		{"session-hq-account-reader", `query { accounts(filter: { recordedOpeningRecords: { recordNumber: "OPEN-1" } }) { total } }`},
		{"session-hq", `query { organizations(filter: { openingRecords: { recordNumber: "OPEN-1" } }) { total } }`},
		{"session-a-1", `query { franchiseOpeningRecords { total } }`},
	}
	for _, item := range queries {
		assertCode(t, fixture.execute(item.session, item.query), auth.CodePermissionDenied)
	}
}
