package integration_test

import (
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
)

func TestHQInvitationReaderCannotReadFranchiseInvitations(t *testing.T) {
	fixture := newSecurityFixture(t)
	seedHQCustomAdministrator(fixture, "invitation-reader", []string{"hqMembership:read"})
	fixture.createAll([]gen.MembershipInvitation{{
		ID: "invitation-hq", MembershipID: "membership-hq", InvitedByAccountID: "account-shared", ExpiresAt: time.Now().Add(time.Hour),
	}})

	listed := fixture.execute("session-hq-invitation-reader", `query { membershipInvitations { total data { id membershipId } } }`)
	assertNoErrors(t, listed)
	assertResultTotal(t, listed, "membershipInvitations", 1)
	if !strings.Contains(listed.Body, "invitation-hq") || strings.Contains(listed.Body, "invitation-a") || strings.Contains(listed.Body, "invitation-b") {
		t.Fatalf("HQ invitation list escaped organization scope: %s", listed.Body)
	}

	foreign := fixture.execute("session-hq-invitation-reader", `query { membershipInvitation(id: "invitation-a") { id } }`)
	assertCode(t, foreign, auth.CodePermissionDenied)
	own := fixture.execute("session-hq-invitation-reader", `query { membershipInvitation(id: "invitation-hq") { id membershipId } }`)
	assertNoErrors(t, own)
	if !strings.Contains(own.Body, "invitation-hq") {
		t.Fatalf("HQ invitation hidden from its own organization: %s", own.Body)
	}
}
