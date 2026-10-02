package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"base-engine/src"
	"base-engine/src/services/authentication"
	"base-engine/src/services/authorization"
)

type initialPasswordPayload struct {
	AccountID         string `json:"accountId"`
	TemporaryPassword string `json:"temporaryPassword"`
}

func TestFranchiseInitialPasswordResetsSharedAccount(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedInitialPasswordTarget(t, fixture)
	fixture.createAll([]membershipRole{{"membership-staff-a", "role-owner-a"}})
	response := fixture.execute("session-hq", `mutation { resetFranchiseInitialPassword(organizationId: "org-a") { accountId temporaryPassword } }`)
	if len(response.Errors) != 0 {
		t.Fatalf("reset failed: %s", response.Body)
	}
	var payload initialPasswordPayload
	if err := json.Unmarshal(response.Data["resetFranchiseInitialPassword"], &payload); err != nil {
		t.Fatal(err)
	}
	if payload.AccountID != "account-initial" || !regexp.MustCompile(`^[A-Za-z0-9]{8}$`).MatchString(payload.TemporaryPassword) {
		t.Fatalf("unexpected reset payload: %#v", payload)
	}
	assertInitialPasswordCredential(t, fixture, payload)
	assertInitialPasswordRevocationAndAudit(t, fixture, payload)
}

func assertInitialPasswordCredential(t *testing.T, fixture *securityFixture, payload initialPasswordPayload) {
	t.Helper()
	var account gen.Account
	var credential authentication.AccountCredential
	fixture.db.First(&account, "id = ?", payload.AccountID)
	fixture.db.First(&credential, "account_id = ?", payload.AccountID)
	if account.CredentialVersion != 2 || !account.MustChangePassword || auth.VerifyPassword(credential.PasswordHash, payload.TemporaryPassword) != nil {
		t.Fatalf("credential was not rotated: account=%#v", account)
	}
}

func assertInitialPasswordRevocationAndAudit(t *testing.T, fixture *securityFixture, payload initialPasswordPayload) {
	t.Helper()
	var sessions []gen.Session
	fixture.db.Where("account_id = ?", payload.AccountID).Find(&sessions)
	if len(sessions) != 2 || sessions[0].RevokedAt == nil || sessions[1].RevokedAt == nil {
		t.Fatalf("shared account sessions not all revoked: %#v", sessions)
	}
	var record gen.AuditLog
	if err := fixture.db.First(&record, "action = ?", "franchiseInitialAccount:reset_password").Error; err != nil || record.OrganizationID == nil || *record.OrganizationID != "org-a" || record.ResourceID == nil || *record.ResourceID != payload.AccountID || strings.Contains(*record.MetadataJSON, payload.TemporaryPassword) {
		t.Fatalf("audit invalid: %#v err=%v", record, err)
	}
}

func TestFranchiseInitialPasswordRejectsUnmarkedAndHQAccounts(t *testing.T) {
	fixture := newSecurityFixture(t)
	unmarked := fixture.execute("session-hq", `mutation { resetFranchiseInitialPassword(organizationId: "org-a") { accountId } }`)
	assertCode(t, unmarked, auth.CodePermissionDenied)
	fixture.db.Model(&gen.Organization{}).Where("id = ?", "org-a").Update("initial_account_id", "account-shared")
	hqAccount := fixture.execute("session-hq", `mutation { resetFranchiseInitialPassword(organizationId: "org-a") { accountId } }`)
	assertCode(t, hqAccount, auth.CodePermissionDenied)
	franchise := fixture.execute("session-a-1", `mutation { resetFranchiseInitialPassword(organizationId: "org-a") { accountId } }`)
	assertCode(t, franchise, auth.CodePermissionDenied)
}

func TestHistoricalFranchiseInitialAccountCanBeInitializedOnce(t *testing.T) {
	fixture := newSecurityFixture(t)
	principal := &auth.WorkspacePrincipal{AccountID: "account-hq", SessionID: "session-hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"account:update": {}, "hqMembership:read": {}}}
	if err := authorization.SetFranchiseInitialAccount(fixture.db, principal, "org-a", "account-staff", "OPEN-2023-001"); err != nil {
		t.Fatal(err)
	}
	var organization gen.Organization
	fixture.db.First(&organization, "id = ?", "org-a")
	if organization.InitialAccountID == nil || *organization.InitialAccountID != "account-staff" {
		t.Fatalf("initial account marker not saved: %#v", organization.InitialAccountID)
	}
	assertHistoricalConfirmationAudit(t, fixture)
	if err := authorization.SetFranchiseInitialAccount(fixture.db, principal, "org-a", "account-staff", "OPEN-2023-001"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("repeat confirmation error = %v", err)
	}
}

func TestHistoricalFranchiseInitialAccountHTTPRequiresEvidence(t *testing.T) {
	fixture := newSecurityFixture(t)
	principal := &auth.WorkspacePrincipal{AccountID: "account-hq", SessionID: "session-hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"account:update": {}, "hqMembership:read": {}}}
	router := mux.NewRouter()
	src.RegisterFranchiseInitialAccountRoute(router, fixture.db, config.SecurityConfig{AllowedOrigins: map[string]struct{}{"http://admin.test": {}}})
	for _, test := range []struct {
		reference, expectedAccountID string
		attested                     bool
		status                       int
	}{{"", "", true, http.StatusBadRequest}, {"OPEN-001", "", false, http.StatusBadRequest}, {"OPEN-001", "account-staff", true, http.StatusBadRequest}, {"OPEN-001", "", true, http.StatusOK}} {
		input := map[string]any{"organizationId": "org-a", "accountId": "account-staff", "evidenceReference": test.reference, "attested": test.attested}
		if test.expectedAccountID != "" {
			input["expectedAccountId"] = test.expectedAccountID
		}
		body, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		request := httptest.NewRequest(http.MethodPost, "/api/franchise-initial-account", strings.NewReader(string(body)))
		request.Header.Set("Origin", "http://admin.test")
		request.Header.Set("Content-Type", "application/json")
		request = request.WithContext(auth.WithPrincipal(request.Context(), principal))
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != test.status {
			t.Fatalf("reference=%q expectedAccountID=%q status=%d body=%s", test.reference, test.expectedAccountID, response.Code, response.Body.String())
		}
	}
	var record gen.AuditLog
	if err := fixture.db.First(&record, "action = ?", "franchiseInitialAccount:confirm").Error; err != nil || record.MetadataJSON == nil || !strings.Contains(*record.MetadataJSON, "OPEN-001") {
		t.Fatalf("HTTP confirmation audit = %#v err=%v", record, err)
	}
}

func TestHistoricalFranchiseInitialAccountCannotBeChangedOrReverified(t *testing.T) {
	fixture := newSecurityFixture(t)
	principal := &auth.WorkspacePrincipal{AccountID: "account-hq", SessionID: "session-hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"account:update": {}, "hqMembership:read": {}}}
	if err := authorization.SetFranchiseInitialAccount(fixture.db, principal, "org-a", "account-staff", "OPEN-001"); err != nil {
		t.Fatal(err)
	}
	fixture.createAll([]gen.Account{{ID: "account-second", Phone: "13800000888", DisplayName: "Second", Status: gen.AccountStatusActive}}, []gen.OperatorMembership{{ID: "membership-second", AccountID: "account-second", OrganizationID: "org-a", Status: gen.MembershipStatusActive}})
	for _, accountID := range []string{"account-staff", "account-second"} {
		if err := authorization.SetFranchiseInitialAccount(fixture.db, principal, "org-a", accountID, "OPEN-002"); auth.ErrorCode(err) != auth.CodeConflict {
			t.Fatalf("reassignment to %s error = %v", accountID, err)
		}
	}
	var organization gen.Organization
	fixture.db.First(&organization, "id = ?", "org-a")
	if organization.InitialAccountID == nil || *organization.InitialAccountID != "account-staff" {
		t.Fatalf("initial account changed: %#v", organization.InitialAccountID)
	}
	var auditCount int64
	fixture.db.Model(&gen.AuditLog{}).Where("action IN ?", []string{"franchiseInitialAccount:correct", "franchiseInitialAccount:verify"}).Count(&auditCount)
	if auditCount != 0 {
		t.Fatalf("unexpected correction or verification audit: %d", auditCount)
	}
}

func TestHistoricalFranchiseInitialAccountConfirmationRejectsUnsafeTargets(t *testing.T) {
	fixture := newSecurityFixture(t)
	for _, query := range []string{
		`mutation { updateOrganization(id: "org-a", input: {name: "Changed"}) { id } }`,
		`mutation { updateOrganization(id: "org-a", input: {initialAccount: {id: "account-shared"}}) { id } }`,
		`mutation { updateOrganization(id: "org-a", input: {initialAccount: {id: "account-staff"}, name: "Changed"}) { id } }`,
		`mutation { updateOrganization(id: "org-a", input: {initialAccount: {id: "account-staff", phone: "13800000002"}}) { id } }`,
		`mutation { updateOrganization(id: "org-b", input: {initialAccount: {id: "account-staff"}}) { id } }`,
	} {
		assertCode(t, fixture.execute("session-hq", query), auth.CodePermissionDenied)
	}
	assertCode(t, fixture.execute("session-a-1", `mutation { updateOrganization(id: "org-a", input: {initialAccount: {id: "account-staff"}}) { id } }`), auth.CodePermissionDenied)
	seedHQCustomAdministrator(fixture, "limited-confirm", []string{"organization:read"})
	assertCode(t, fixture.execute("session-hq-limited-confirm", `mutation { updateOrganization(id: "org-a", input: {initialAccount: {id: "account-staff"}}) { id } }`), auth.CodePermissionDenied)
	seedHQCustomAdministrator(fixture, "no-member-read", []string{"account:update"})
	assertCode(t, fixture.execute("session-hq-no-member-read", `mutation { updateOrganization(id: "org-a", input: {initialAccount: {id: "account-staff"}}) { id } }`), auth.CodePermissionDenied)
}

func TestHistoricalFranchiseAccountCandidatesAreScopedToOrganization(t *testing.T) {
	fixture := newSecurityFixture(t)
	query := `query { organization(id: "org-a") { memberships { id status account { id phone displayName status } } } }`
	response := fixture.execute("session-hq", query)
	if len(response.Errors) != 0 || !strings.Contains(response.Body, `"phone":"13800000002"`) {
		t.Fatalf("HQ cannot inspect candidates: %s", response.Body)
	}
	assertCode(t, fixture.execute("session-b", query), auth.CodePermissionDenied)
}

func TestFranchiseInitialPasswordRejectsInactiveTargets(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedInitialPasswordTarget(t, fixture)
	query := `mutation { resetFranchiseInitialPassword(organizationId: "org-a") { accountId } }`
	fixture.db.Model(&gen.OperatorMembership{}).Where("id = ?", "initial-a").Update("status", gen.MembershipStatusSuspended)
	assertCode(t, fixture.execute("session-hq", query), auth.CodeMembershipInactive)
	fixture.db.Model(&gen.OperatorMembership{}).Where("id = ?", "initial-a").Update("status", gen.MembershipStatusActive)
	fixture.db.Model(&gen.Account{}).Where("id = ?", "account-initial").Update("status", gen.AccountStatusDisabled)
	assertCode(t, fixture.execute("session-hq", query), auth.CodePermissionDenied)
	fixture.db.Model(&gen.Account{}).Where("id = ?", "account-initial").Update("status", gen.AccountStatusActive)
	fixture.db.Model(&gen.OperatorMembership{}).Where("id = ?", "initial-a").Update("is_delete", 2)
	assertCode(t, fixture.execute("session-hq", query), auth.CodeMembershipInactive)
}

func TestFranchiseInitialPasswordRequiresAccountUpdate(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedInitialPasswordTarget(t, fixture)
	seedHQCustomAdministrator(fixture, "limited-reset", []string{"organization:read"})
	response := fixture.execute("session-hq-limited-reset", `mutation { resetFranchiseInitialPassword(organizationId: "org-a") { accountId } }`)
	assertCode(t, response, auth.CodePermissionDenied)
}

func TestFranchiseInitialAccountIDProjectionIsRestricted(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedInitialPasswordTarget(t, fixture)
	query := `query { organization(id: "org-a") { initialAccountId } }`
	assertCode(t, fixture.execute("session-a-1", query), auth.CodePermissionDenied)
	response := fixture.execute("session-hq", query)
	if len(response.Errors) != 0 || !strings.Contains(response.Body, `"initialAccountId":"account-initial"`) {
		t.Fatalf("HQ initial account projection failed: %s", response.Body)
	}
}

func TestFranchiseInitialAccountQueryArgumentsRequireResetPermission(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedInitialPasswordTarget(t, fixture)
	seedHQCustomAdministrator(fixture, "org-reader", []string{"organization:read"})
	seedHQCustomAdministrator(fixture, "account-reader", []string{"account:read"})
	queries := []struct {
		session string
		query   string
	}{
		{"session-hq-org-reader", `query { organizations(filter: { initialAccount: { phone: "13800000999" } }) { total } }`},
		{"session-hq-org-reader", `query { organizations(filter: { OR: [{ initialAccountId_null: false }, { type: FRANCHISE }] }) { total } }`},
		{"session-hq-org-reader", `query { organizations(sort: [{ initialAccountId: ASC }]) { total } }`},
		{"session-hq-account-reader", `query { accounts(filter: { initializedOrganizations: { id: "org-a" } }) { total } }`},
		{"session-a-1", `query { organizations(filter: { initialAccountId: "account-initial" }) { total } }`},
	}
	for _, candidate := range queries {
		assertCode(t, fixture.execute(candidate.session, candidate.query), auth.CodePermissionDenied)
	}
	allowed := fixture.execute("session-hq", `query { organizations(filter: { initialAccountId: "account-initial" }) { total } }`)
	if len(allowed.Errors) != 0 {
		t.Fatalf("authorized initial-account filter rejected: %s", allowed.Body)
	}
}

func seedInitialPasswordTarget(t *testing.T, fixture *securityFixture) {
	t.Helper()
	if err := fixture.db.AutoMigrate(&authentication.AccountCredential{}); err != nil {
		t.Fatal(err)
	}
	hash, err := auth.HashPassword("OldPassword8")
	if err != nil {
		t.Fatal(err)
	}
	fixture.createAll(
		[]gen.Account{{ID: "account-initial", Phone: "13800000999", DisplayName: "Initial", Status: gen.AccountStatusActive, CredentialVersion: 1}},
		[]gen.OperatorMembership{
			{ID: "initial-a", AccountID: "account-initial", OrganizationID: "org-a", Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores},
			{ID: "initial-b", AccountID: "account-initial", OrganizationID: "org-b", Status: gen.MembershipStatusActive, StoreAccessMode: gen.StoreAccessModeAllStores},
		},
		[]authentication.AccountCredential{{AccountID: "account-initial", PasswordHash: hash, PasswordChangedAt: time.Now(), UpdatedAt: time.Now()}},
	)
	for _, id := range []string{"org-a", "org-b"} {
		fixture.db.Model(&gen.Organization{}).Where("id = ?", id).Update("initial_account_id", "account-initial")
		fixture.createAll([]gen.Session{{ID: "initial-session-" + id, AccountID: "account-initial", OrganizationID: &id, WorkspaceType: gen.WorkspaceTypeFranchise, CredentialVersion: 1, ExpiresAt: time.Now().Add(time.Hour), LastSeenAt: time.Now()}})
	}
}
