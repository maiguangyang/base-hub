/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package organization

import (
	"context"
	"testing"
	"time"

	"base-engine/gen"
	sessionservice "base-engine/src/services/session"
)

// TestSuspendRevokesOnlyOrganizationSessionsAndPublishesAfterCommit 验证组织暂停的持久化与强退边界。
func TestSuspendRevokesOnlyOrganizationSessionsAndPublishesAfterCommit(t *testing.T) {
	_, db := newOrganizationFixture(t)
	publisher := sessionservice.NewPublisher()
	service := NewService(db, nil, publisher)
	principal := headquartersPrincipal()
	principal.Permissions["organization:suspend"] = struct{}{}
	organization := gen.Organization{ID: "target-org", Code: "T", Name: "Target", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}
	db.Create(&organization)
	now := time.Now()
	organizationID := organization.ID
	target := gen.Session{ID: "target-session", AccountID: "a", OrganizationID: &organizationID, WorkspaceType: gen.WorkspaceTypeFranchise, CredentialVersion: 1, ExpiresAt: now.Add(time.Hour), LastSeenAt: now}
	otherOrganizationID := "other-org"
	other := gen.Session{ID: "other-session", AccountID: "a", OrganizationID: &otherOrganizationID, WorkspaceType: gen.WorkspaceTypeFranchise, CredentialVersion: 1, ExpiresAt: now.Add(time.Hour), LastSeenAt: now}
	db.Create(&target)
	db.Create(&other)
	targetEvents, cancelTarget := subscribeOrganizationTest(t, publisher, target.ID)
	defer cancelTarget()
	otherEvents, cancelOther := subscribeOrganizationTest(t, publisher, other.ID)
	defer cancelOther()
	if err := service.Suspend(context.Background(), principal, organization.ID, "CONTRACT_ENDED"); err != nil {
		t.Fatal(err)
	}
	var persisted gen.Session
	db.First(&persisted, "id = ?", target.ID)
	if persisted.RevokedAt == nil || persisted.RevocationCode == nil || *persisted.RevocationCode != "ORGANIZATION_SUSPENDED" {
		t.Fatalf("target session not revoked: %#v", persisted)
	}
	assertOrganizationEvent(t, targetEvents)
	select {
	case event := <-otherEvents:
		t.Fatalf("unrelated session received event: %#v", event)
	case <-time.After(50 * time.Millisecond):
	}
}

func subscribeOrganizationTest(t *testing.T, publisher *sessionservice.LocalPublisher, sessionID string) (<-chan *gen.SessionEvent, context.CancelFunc) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	channel, err := publisher.Subscribe(ctx, sessionID)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	return channel, cancel
}

func assertOrganizationEvent(t *testing.T, channel <-chan *gen.SessionEvent) {
	t.Helper()
	select {
	case event := <-channel:
		if event.Code != gen.SessionEventCodeOrganizationSuspended {
			t.Fatalf("event = %#v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("event timeout")
	}
}
