package integration_test

import (
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
)

func TestHistoricalInitializationWritesTypedOpeningRecord(t *testing.T) {
	fixture := newSecurityFixture(t)
	principal := &auth.WorkspacePrincipal{AccountID: "account-shared", SessionID: "session-hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"account:update": {}, "hqMembership:read": {}}}
	if err := authorization.SetFranchiseInitialAccount(fixture.db, principal, "org-a", "account-staff", "OPEN-001"); err != nil {
		t.Fatal(err)
	}
	var record gen.FranchiseOpeningRecord
	if err := fixture.db.Where("organization_id = ?", "org-a").First(&record).Error; err != nil {
		t.Fatal(err)
	}
	if record.RecordNumber != "OPEN-001" || record.InitialAccountID != "account-staff" || record.RecordedByAccountID != "account-shared" || record.Source != gen.FranchiseOpeningSourceHistoricalAttestation {
		t.Fatalf("opening record was not linked to org/account/operator: %#v", record)
	}
}

func TestHistoricalOpeningRecordNumberCannotBeReused(t *testing.T) {
	fixture := newSecurityFixture(t)
	fixture.createAll([]gen.OperatorMembership{{ID: "membership-staff-b", AccountID: "account-staff", OrganizationID: "org-b", Status: gen.MembershipStatusActive}})
	principal := &auth.WorkspacePrincipal{AccountID: "account-shared", SessionID: "session-hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"account:update": {}, "hqMembership:read": {}}}
	if err := authorization.SetFranchiseInitialAccount(fixture.db, principal, "org-a", "account-staff", "OPEN-001"); err != nil {
		t.Fatal(err)
	}
	err := authorization.SetFranchiseInitialAccount(fixture.db, principal, "org-b", "account-staff", "OPEN-001")
	if auth.ErrorCode(err) != auth.CodeOpeningRecordNumberConflict {
		t.Fatalf("reused record number error = %v", err)
	}
	var organization gen.Organization
	if err := fixture.db.First(&organization, "id = ?", "org-b").Error; err != nil {
		t.Fatal(err)
	}
	if organization.InitialAccountID != nil {
		t.Fatal("failed initialization set marker")
	}
}
