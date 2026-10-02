/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-22
 */
package integration_test

import (
	"context"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	sessionservice "base-engine/src/services/session"
	"gorm.io/gorm"
)

func TestHQMembershipRoleAssignmentRevokesOnlyTargetWorkspace(t *testing.T) {
	fixture, events := newRevocationFixture(t, "assign")
	response := fixture.execute("session-hq", `mutation {
		updateOperatorMembership(id: "membership-hq-assign", input: {rolesIds: ["role-hq-assign"]}) { id }
	}`)
	assertNoErrors(t, response)
	assertRevokedSession(t, fixture.db, "session-hq-assign")
	assertActiveSession(t, fixture.db, "session-other-assign")
	assertRevocationEvent(t, events, "session-hq-assign", "org-hq")
	assertCode(t, fixture.execute("session-hq-assign", `query { viewer { account { id } } }`), auth.CodeSessionRevoked)
}

func TestHQSuspensionAndDeletionRevokeSessions(t *testing.T) {
	for _, test := range []struct {
		name, suffix, mutation string
	}{
		{name: "suspend", suffix: "suspend", mutation: `mutation { changeMembershipStatus(input: {membershipId: "membership-hq-suspend", status: SUSPENDED}) { id } }`},
		{name: "delete", suffix: "delete", mutation: `mutation { deleteOperatorMemberships(id: ["membership-hq-delete"]) }`},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture, events := newRevocationFixture(t, test.suffix)
			assertNoErrors(t, fixture.execute("session-hq", test.mutation))
			assertRevokedSession(t, fixture.db, "session-hq-"+test.suffix)
			assertRevocationEvent(t, events, "session-hq-"+test.suffix, "org-hq")
		})
	}
}

func TestHQRolePermissionChangeRevokesActiveMembers(t *testing.T) {
	fixture, events := newRevocationFixture(t, "role-change")
	response := fixture.execute("session-hq", `mutation {
		updateOperatorRole(id: "role-hq-role-change", input: {permissionsIds: ["permission-hqMembership-read"]}) { id }
	}`)
	assertNoErrors(t, response)
	assertRevokedSession(t, fixture.db, "session-hq-role-change")
	assertRevocationEvent(t, events, "session-hq-role-change", "org-hq")
}

func TestHQSessionRevocationPersistsWithoutPublisher(t *testing.T) {
	fixture := newSecurityFixtureWithPublisher(t, nil)
	seedHQCustomAdministrator(fixture, "no-publisher", []string{"hqMembership:read"})
	response := fixture.execute("session-hq", `mutation {
		updateOperatorMembership(id: "membership-hq-no-publisher", input: {rolesIds: ["role-hq-no-publisher"]}) { id }
	}`)
	assertNoErrors(t, response)
	assertRevokedSession(t, fixture.db, "session-hq-no-publisher")
	assertCode(t, fixture.execute("session-hq-no-publisher", `query { viewer { account { id } } }`), auth.CodeSessionRevoked)
}

func newRevocationFixture(t *testing.T, suffix string) (*securityFixture, <-chan *gen.SessionEvent) {
	t.Helper()
	publisher := &commitCheckingPublisher{local: sessionservice.NewPublisher()}
	fixture := newSecurityFixtureWithPublisher(t, publisher)
	publisher.db = fixture.db
	seedHQCustomAdministrator(fixture, suffix, []string{"hqMembership:read"})
	accountID := "account-hq-" + suffix
	otherOrganizationID := "org-a"
	otherSessionID := "session-other-" + suffix
	fixture.createAll([]gen.Session{{
		ID: otherSessionID, AccountID: accountID, OrganizationID: &otherOrganizationID,
		WorkspaceType: gen.WorkspaceTypeFranchise, CredentialVersion: 1,
		ExpiresAt: time.Now().Add(time.Hour), LastSeenAt: time.Now(),
	}})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	events, err := publisher.Subscribe(ctx, "session-hq-"+suffix)
	if err != nil {
		t.Fatal(err)
	}
	return fixture, events
}

type commitCheckingPublisher struct {
	local *sessionservice.LocalPublisher
	db    *gorm.DB
}

func (p *commitCheckingPublisher) Subscribe(ctx context.Context, sessionID string) (<-chan *gen.SessionEvent, error) {
	return p.local.Subscribe(ctx, sessionID)
}

func (p *commitCheckingPublisher) PublishSession(sessionID string, event *gen.SessionEvent) {
	var stored gen.Session
	if p.db == nil || p.db.First(&stored, "id = ?", sessionID).Error != nil || stored.RevokedAt == nil {
		panic("session event published before revocation committed")
	}
	p.local.PublishSession(sessionID, event)
}

func assertRevokedSession(t *testing.T, db *gorm.DB, sessionID string) {
	t.Helper()
	var stored gen.Session
	if err := db.First(&stored, "id = ?", sessionID).Error; err != nil || stored.RevokedAt == nil {
		t.Fatalf("session %s not revoked: %#v, %v", sessionID, stored, err)
	}
}

func assertActiveSession(t *testing.T, db *gorm.DB, sessionID string) {
	t.Helper()
	var stored gen.Session
	if err := db.First(&stored, "id = ?", sessionID).Error; err != nil || stored.RevokedAt != nil {
		t.Fatalf("session %s unexpectedly revoked: %#v, %v", sessionID, stored, err)
	}
}

func assertRevocationEvent(t *testing.T, events <-chan *gen.SessionEvent, sessionID, organizationID string) {
	t.Helper()
	select {
	case event := <-events:
		if event == nil || event.Code != gen.SessionEventCodeSessionRevoked || event.SessionID != sessionID || event.OrganizationID == nil || *event.OrganizationID != organizationID {
			t.Fatalf("unexpected session event: %#v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("session revocation event timeout")
	}
}
