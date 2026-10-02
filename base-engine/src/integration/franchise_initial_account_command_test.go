package integration_test

import (
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
)

func TestInitialAccountCommandRejectsUnauthorizedAndUnsafeCandidates(t *testing.T) {
	fixture := newSecurityFixture(t)
	full := &auth.WorkspacePrincipal{AccountID: "account-hq", SessionID: "session-hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"account:update": {}, "hqMembership:read": {}}}
	missingAccountUpdate := &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"hqMembership:read": {}}}
	missingMembershipRead := &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"account:update": {}}}
	franchise := &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeFranchise, Permissions: full.Permissions}
	tests := []struct {
		principal                           *auth.WorkspacePrincipal
		organizationID, accountID, evidence string
		want                                auth.Code
	}{
		{missingAccountUpdate, "org-a", "account-staff", "OPEN-001", auth.CodePermissionDenied},
		{missingMembershipRead, "org-a", "account-staff", "OPEN-001", auth.CodePermissionDenied},
		{franchise, "org-a", "account-staff", "OPEN-001", auth.CodeWorkspaceForbidden},
		{full, "org-hq", "account-staff", "OPEN-001", auth.CodePermissionDenied},
		{full, "org-a", "account-shared", "OPEN-001", auth.CodePermissionDenied},
		{full, "org-b", "account-staff", "OPEN-001", auth.CodePermissionDenied},
		{full, "org-a", "account-staff", "", auth.CodeValidationFailed},
		{full, "org-a", "account-staff", strings.Repeat("测", 129), auth.CodeValidationFailed},
	}
	for index, test := range tests {
		err := authorization.SetFranchiseInitialAccount(fixture.db, test.principal, test.organizationID, test.accountID, test.evidence)
		if auth.ErrorCode(err) != test.want {
			t.Fatalf("case %d err=%v want=%s", index, err, test.want)
		}
	}
	var organization gen.Organization
	if err := fixture.db.First(&organization, "id = ?", "org-a").Error; err != nil {
		t.Fatal(err)
	}
	if organization.InitialAccountID != nil {
		t.Fatalf("rejected command changed marker: %#v", organization.InitialAccountID)
	}
}
